package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type EstimatedTimetableSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestLineUpdate(request *LineUpdateRequest)
	HandleNotifyEstimatedTimetable(delivery *sxml.XMLNotifyEstimatedTimetable)
}

type SIRIEstimatedTimetableSubscriptionCollector struct {
	connector

	deletedSubscriptions         *DeletedSubscriptions
	estimatedTimetableSubscriber SIRIEstimatedTimetableSubscriber
	updateSubscriber             UpdateSubscriber
}

type SIRIEstimatedTimetableSubscriptionCollectorFactory struct{}

func (factory *SIRIEstimatedTimetableSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIEstimatedTimetableSubscriptionCollector(partner)
}

func (factory *SIRIEstimatedTimetableSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIEstimatedTimetableSubscriptionCollector(partner *Partner) *SIRIEstimatedTimetableSubscriptionCollector {
	connector := &SIRIEstimatedTimetableSubscriptionCollector{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind()
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent
	connector.estimatedTimetableSubscriber = NewSIRIEstimatedTimetableSubscriber(connector)

	return connector
}

func (connector *SIRIEstimatedTimetableSubscriptionCollector) Stop() {
	connector.estimatedTimetableSubscriber.Stop()
}

func (connector *SIRIEstimatedTimetableSubscriptionCollector) Start() {
	connector.deletedSubscriptions = NewDeletedSubscriptions()
	connector.estimatedTimetableSubscriber.Start()
}

func (connector *SIRIEstimatedTimetableSubscriptionCollector) RequestLineUpdate(request *LineUpdateRequest) {
	line, ok := connector.partner.Model().Lines().Find(request.LineId())
	if !ok {
		logger.Log.Debugf("LineUpdateRequest in EstimatedTimetable SubscriptionCollector for unknown Line %v", request.LineId())
		return
	}

	lineObjectid, ok := line.ObjectID(connector.remoteObjectidKind)
	if !ok {
		logger.Log.Debugf("Requested line %v doesn't have and objectId of kind %v", request.LineId(), connector.remoteObjectidKind)
		return
	}

	// Try to find a Subscription with the resource
	subscriptions := connector.partner.Subscriptions().FindByResourceId(lineObjectid.String(), EstimatedTimetableCollect)
	if len(subscriptions) > 0 {
		for _, subscription := range subscriptions {
			resource := subscription.Resource(lineObjectid)
			if resource == nil { // Should never happen
				logger.Log.Debugf("Can't find resource in subscription after Subscriptions#FindByResourceId")
				return
			}
			if !resource.SubscribedAt.IsZero() {
				resource.SubscribedUntil = connector.Clock().Now().Add(2 * time.Minute)
			}
		}
		return
	}

	// Else we find or create a subscription to add the resource
	newSubscription := connector.partner.Subscriptions().FindOrCreateByKind(EstimatedTimetableCollect)
	ref := model.Reference{
		ObjectId: &lineObjectid,
		Type:     "Line",
	}

	newSubscription.CreateAddNewResource(ref)
}

func (connector *SIRIEstimatedTimetableSubscriptionCollector) SetEstimatedTimetableSubscriber(estimatedTimetableSubscriber SIRIEstimatedTimetableSubscriber) {
	connector.estimatedTimetableSubscriber = estimatedTimetableSubscriber
}

func (connector *SIRIEstimatedTimetableSubscriptionCollector) HandleNotifyEstimatedTimetable(notify *sxml.XMLNotifyEstimatedTimetable) {
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})

	for _, delivery := range notify.EstimatedTimetableDeliveries() {
		subscriptionId := delivery.SubscriptionRef()

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Debugf("Partner %s sent a StopVisitNotify response to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			if !connector.deletedSubscriptions.AlreadySend(subscriptionId) {
				subToDelete[subscriptionId] = struct{}{}
			}
			continue
		}
		if subscription.Kind() != EstimatedTimetableCollect {
			logger.Log.Debugf("Partner %s sent a StopVisitNotify response to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind EstimatedTimetableCollect"
			continue
		}

		// builder := NewEstimatedTimetableUpdateEventBuilder(connector.partner)
		// builder.SetUpdateEvents(delivery.XMLMonitoredStopVisits())
		// builder.SetStopVisitCancellationEvents(delivery)
		// updateEvents := builder.UpdateEvents()

		// connector.broadcastUpdateEvents(&updateEvents)
	}

	for subId := range subToDelete {
		CancelSubscription(subId, "EstimatedTimetableSubscriptionCollector", connector)
	}
}

func (connector *SIRIEstimatedTimetableSubscriptionCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {
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
}
