package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
	"golang.org/x/exp/maps"
)

type StopMonitoringSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
	HandleNotifyStopMonitoring(delivery *sxml.XMLNotifyStopMonitoring) *CollectedRefs
}

type SIRIStopMonitoringSubscriptionCollector struct {
	connector

	deletedSubscriptions     *DeletedSubscriptions
	stopMonitoringSubscriber SIRIStopMonitoringSubscriber
	updateSubscriber         UpdateSubscriber
}

type SIRIStopMonitoringSubscriptionCollectorFactory struct{}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIStopMonitoringSubscriptionCollector(partner *Partner) *SIRIStopMonitoringSubscriptionCollector {
	connector := &SIRIStopMonitoringSubscriptionCollector{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace()
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent
	connector.stopMonitoringSubscriber = NewSIRIStopMonitoringSubscriber(connector)

	return connector
}

func (connector *SIRIStopMonitoringSubscriptionCollector) Stop() {
	connector.stopMonitoringSubscriber.Stop()
}

func (connector *SIRIStopMonitoringSubscriptionCollector) Start() {
	connector.deletedSubscriptions = NewDeletedSubscriptions()
	connector.stopMonitoringSubscriber.Start()
}

func (connector *SIRIStopMonitoringSubscriptionCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	stopArea, ok := connector.partner.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("StopAreaUpdateRequest in StopMonitoring SubscriptionCollector for unknown StopArea %v", request.StopAreaId())
		return
	}

	stopAreaCode, ok := stopArea.Code(connector.remoteCodeSpace)
	if !ok {
		logger.Log.Debugf("Requested stopArea %v doesn't have a code with codeSpace %v", request.StopAreaId(), connector.remoteCodeSpace)
		return
	}

	// Try to find a Subscription with the resource
	subscriptions := connector.partner.Subscriptions().FindByResourceId(stopAreaCode.String(), StopMonitoringCollect)
	if len(subscriptions) > 0 {
		for _, subscription := range subscriptions {
			resource := subscription.Resource(stopAreaCode)
			if resource == nil { // Should never happen
				logger.Log.Debugf("Can't find resource in subscription after Subscriptions#FindByResourceId")
				return
			}
			if !resource.SubscribedAt().IsZero() {
				resource.SubscribedUntil = connector.Clock().Now().Add(2 * time.Minute)
			}
		}
		return
	}

	// Else we find or create a subscription to add the resource
	newSubscription := connector.partner.Subscriptions().FindOrCreateByKind(StopMonitoringCollect)
	ref := model.Reference{
		Code: &stopAreaCode,
		Type:     "StopArea",
	}

	newSubscription.CreateAndAddNewResource(ref)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) SetStopMonitoringSubscriber(stopMonitoringSubscriber SIRIStopMonitoringSubscriber) {
	connector.stopMonitoringSubscriber = stopMonitoringSubscriber
}

func (connector *SIRIStopMonitoringSubscriptionCollector) HandleNotifyStopMonitoring(notify *sxml.XMLNotifyStopMonitoring) (collectedRefs *CollectedRefs) {
	// subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})
	var updateEvents CollectUpdateEvents

	collectedRefs = NewCollectedRefs()
	for _, delivery := range notify.StopMonitoringDeliveries() {
		subscriptionId := delivery.SubscriptionRef()

		if subscriptionId == "" {
			logger.Log.Debugf("Partner %s sent a NotifyStopMonitoring with an empty SubscriptionRef\n", connector.Partner().Slug())
			continue
		}

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Debugf("Partner %s sent a NotifyStopMonitoring to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			// subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			if !connector.deletedSubscriptions.AlreadySend(subscriptionId) {
				subToDelete[subscriptionId] = struct{}{}
			}
			continue
		}
		if subscription.Kind() != StopMonitoringCollect {
			logger.Log.Debugf("Partner %s sent a NotifyStopMonitoring to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			// subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind StopMonitoringCollect"
			continue
		}

		originStopAreaCode := model.Code{}
		resource := subscription.UniqueResource()
		if resource != nil {
			originStopAreaCode = *resource.Reference.Code
		} else if delivery.MonitoringRef() != "" {
			originStopAreaCode = model.NewCode(connector.remoteCodeSpace, delivery.MonitoringRef())
		}

		builder := NewStopMonitoringUpdateEventBuilder(connector.partner, originStopAreaCode)
		builder.SetUpdateEvents(delivery.XMLMonitoredStopVisits())
		builder.SetStopVisitCancellationEvents(delivery)
		updateEvents = builder.UpdateEvents()

		maps.Copy(collectedRefs.LineRefs, updateEvents.LineRefs)
		maps.Copy(collectedRefs.MonitoringRefs, updateEvents.MonitoringRefs)
		maps.Copy(collectedRefs.VehicleJourneyRefs, updateEvents.VehicleJourneyRefs)

		connector.broadcastUpdateEvents(&updateEvents)
	}

	for subId := range subToDelete {
		CancelSubscription(subId, "StopMonitoringSubscriptionCollector", connector)
	}

	return collectedRefs
}

func (connector *SIRIStopMonitoringSubscriptionCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {
	if connector.updateSubscriber == nil {
		return
	}
	for _, e := range events.StopAreas {
		connector.updateSubscriber(e)
	}
	for _, e := range events.Lines {
		connector.updateSubscriber(e)
	}
	for _, e := range events.VehicleJourneys {
		connector.updateSubscriber(e)
	}
	for _, es := range events.StopVisits { // Stopvisits are map[MonitoringRef]map[ItemIdentifier]event
		for _, e := range es {
			connector.updateSubscriber(e)
		}
	}
	for _, e := range events.Cancellations {
		connector.updateSubscriber(e)
	}
}
