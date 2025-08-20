package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
	"golang.org/x/exp/maps"
)

type FacilityMonitoringSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestFacilityUpdate(request *FacilityUpdateRequest)
	HandleNotifyFacilityMonitoring(delivery *sxml.XMLNotifyFacilityMonitoring) *CollectedRefs
}

type SIRIFacilityMonitoringSubscriptionCollector struct {
	connector

	deletedSubscriptions         *DeletedSubscriptions
	facilityMonitoringSubscriber SIRIFacilityMonitoringSubscriber
	updateSubscriber             UpdateSubscriber
}

type SIRIFacilityMonitoringSubscriptionCollectorFactory struct{}

func (factory *SIRIFacilityMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIFacilityMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIFacilityMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIFacilityMonitoringSubscriptionCollector(partner *Partner) *SIRIFacilityMonitoringSubscriptionCollector {
	connector := &SIRIFacilityMonitoringSubscriptionCollector{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace()
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent
	connector.facilityMonitoringSubscriber = NewSIRIFacilityMonitoringSubscriber(connector)

	return connector
}

func (connector *SIRIFacilityMonitoringSubscriptionCollector) Stop() {
	connector.facilityMonitoringSubscriber.Stop()
}

func (connector *SIRIFacilityMonitoringSubscriptionCollector) Start() {
	connector.deletedSubscriptions = NewDeletedSubscriptions()
	connector.facilityMonitoringSubscriber.Start()
}

func (connector *SIRIFacilityMonitoringSubscriptionCollector) RequestFacilityUpdate(request *FacilityUpdateRequest) {
	facility, ok := connector.partner.Model().Facilities().Find(request.FacilityId())
	if !ok {
		logger.Log.Debugf("FacilityUpdateRequest in FacilityMonitoring SubscriptionCollector for unknown facility %v", request.FacilityId())
		return
	}

	facilityCode, ok := facility.Code(connector.remoteCodeSpace)
	if !ok {
		logger.Log.Debugf("Requested facility %v doesn't have a code with codeSpace %v", request.FacilityId(), connector.remoteCodeSpace)
		return
	}

	// Try to find a Subscription with the resource
	subscriptions := connector.partner.Subscriptions().FindByResourceId(facilityCode.String(), "FacilityMonitoringCollect")
	if len(subscriptions) > 0 {
		for _, subscription := range subscriptions {
			resource := subscription.Resource(facilityCode)
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
	newSubscription := connector.partner.Subscriptions().New("FacilityMonitoringCollect")
	ref := model.Reference{
		Code: &facilityCode,
		Type: "Facility",
	}

	newSubscription.CreateAndAddNewResource(ref)
}

func (connector *SIRIFacilityMonitoringSubscriptionCollector) SetFacilityMonitoringSubscriber(facilityMonitoringSubscriber SIRIFacilityMonitoringSubscriber) {
	connector.facilityMonitoringSubscriber = facilityMonitoringSubscriber
}

func (connector *SIRIFacilityMonitoringSubscriptionCollector) HandleNotifyFacilityMonitoring(notify *sxml.XMLNotifyFacilityMonitoring) (collectedRefs *CollectedRefs) {
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})
	var updateEvents CollectUpdateEvents

	collectedRefs = NewCollectedRefs()
	for _, delivery := range notify.FacilityMonitoringDeliveries() {
		subscriptionId := delivery.SubscriptionRef()

		if subscriptionId == "" {
			logger.Log.Debugf("Partner %s sent a NotifyFacilityMonitoring with an empty SubscriptionRef\n", connector.Partner().Slug())
			continue
		}

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Debugf("Partner %s sent a NotifyFacilityMonitoring to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			if !connector.deletedSubscriptions.AlreadySend(subscriptionId) {
				subToDelete[subscriptionId] = struct{}{}
			}
			continue
		}
		if subscription.Kind() != FacilityMonitoringCollect {
			logger.Log.Debugf("Partner %s sent a NotifyFacilityMonitoring to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind FacilityMonitoringCollect"
			continue
		}

		builder := NewFacilityMonitoringUpdateEventBuilder(connector.partner)
		builder.SetUpdateEvents(delivery.FacilityConditions())

		updateEvents = builder.UpdateEvents()

		maps.Copy(collectedRefs.FacilityRefs, updateEvents.FacilityRefs)

		connector.broadcastUpdateEvents(&updateEvents)
	}

	for subId := range subToDelete {
		CancelSubscription(subId, "FacilityMonitoringSubscriptionCollector", connector)
	}

	return collectedRefs
}

func (connector *SIRIFacilityMonitoringSubscriptionCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {
	if connector.updateSubscriber == nil {
		return
	}

	for _, e := range events.Facilities {
		connector.updateSubscriber(e)
	}
}
