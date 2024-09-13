package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
	"golang.org/x/exp/maps"
)

type GeneralMessageSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestAllSituationsUpdate()
	RequestSituationUpdate(kind string, requestedId model.Code)
	HandleNotifyGeneralMessage(notify *sxml.XMLNotifyGeneralMessage) *CollectedRefs
}

type SIRIGeneralMessageSubscriptionCollector struct {
	connector

	deletedSubscriptions     *DeletedSubscriptions
	generalMessageSubscriber SIRIGeneralMessageSubscriber
	updateSubscriber         UpdateSubscriber
}

type SIRIGeneralMessageSubscriptionCollectorFactory struct{}

func (factory *SIRIGeneralMessageSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageSubscriptionCollector(partner)
}

func (factory *SIRIGeneralMessageSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIGeneralMessageSubscriptionCollector(partner *Partner) *SIRIGeneralMessageSubscriptionCollector {
	connector := &SIRIGeneralMessageSubscriptionCollector{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace()
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent
	connector.generalMessageSubscriber = NewSIRIGeneralMessageSubscriber(connector)

	return connector
}

func (connector *SIRIGeneralMessageSubscriptionCollector) Stop() {
	connector.generalMessageSubscriber.Stop()
}

func (connector *SIRIGeneralMessageSubscriptionCollector) Start() {
	connector.deletedSubscriptions = NewDeletedSubscriptions()
	connector.generalMessageSubscriber.Start()
}

func (connector *SIRIGeneralMessageSubscriptionCollector) RequestAllSituationsUpdate() {
	obj := model.NewCode(GeneralMessageCollect, "all")
	connector.RequestSituationUpdate("all", obj)
}

func (connector *SIRIGeneralMessageSubscriptionCollector) RequestSituationUpdate(kind string, requestedCode model.Code) {
	// Try to find a Subscription with the resource
	subscriptions := connector.partner.Subscriptions().FindByResourceId(requestedCode.String(), GeneralMessageCollect)
	if len(subscriptions) > 0 {
		for _, subscription := range subscriptions {
			resource := subscription.Resource(requestedCode)
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
	newSubscription := connector.partner.Subscriptions().FindOrCreateByKind(GeneralMessageCollect)
	ref := model.Reference{
		Code: &requestedCode,
	}
	switch kind {
	case SITUATION_UPDATE_REQUEST_LINE:
		ref.Type = "Line"
	case SITUATION_UPDATE_REQUEST_STOP_AREA:
		ref.Type = "StopArea"
	}

	newSubscription.CreateAndAddNewResource(ref)
}

func (connector *SIRIGeneralMessageSubscriptionCollector) HandleNotifyGeneralMessage(notify *sxml.XMLNotifyGeneralMessage) (collectedRefs *CollectedRefs) {
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})

	updateEvents := NewCollectUpdateEvents()
	builder := NewGeneralMessageUpdateEventBuilder(connector.partner)

	collectedRefs = NewCollectedRefs()
	for _, delivery := range notify.GeneralMessagesDeliveries() {
		subscriptionId := delivery.SubscriptionRef()

		if subscriptionId == "" {
			logger.Log.Debugf("Partner %s sent a NotifyGeneralMessage with an empty SubscriptionRef\n", connector.Partner().Slug())
			continue
		}

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Printf("Partner %s sent a NotifyGeneralMessage to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			if !connector.deletedSubscriptions.AlreadySend(subscriptionId) {
				subToDelete[subscriptionId] = struct{}{}
			}
			continue
		}

		if subscription.Kind() != GeneralMessageCollect {
			logger.Log.Printf("Partner %s sent a NotifyGeneralMessage to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind StopMonitoringCollect"
			continue
		}
		connector.cancelGeneralMessage(delivery)

		builder.SetGeneralMessageDeliveryUpdateEvents(updateEvents, delivery, notify.ProducerRef())

		maps.Copy(collectedRefs.LineRefs, builder.LineRefs)
		maps.Copy(collectedRefs.MonitoringRefs, builder.MonitoringRefs)

		connector.broadcastSituationUpdateEvent(updateEvents)
	}

	for subId := range subToDelete {
		CancelSubscription(subId, "GeneralMessageSubscriptionCollector", connector)
	}

	return collectedRefs
}

func (connector *SIRIGeneralMessageSubscriptionCollector) cancelGeneralMessage(xmlResponse *sxml.XMLGeneralMessageDelivery) {
	xmlGmCancellations := xmlResponse.XMLGeneralMessagesCancellations()

	if len(xmlGmCancellations) == 0 {
		return
	}

	for _, cancellation := range xmlGmCancellations {
		obj := model.NewCode(connector.remoteCodeSpace, cancellation.InfoMessageIdentifier())
		situation, ok := connector.partner.Model().Situations().FindByCode(obj)
		if ok {
			logger.Log.Debugf("Updating situation %v progress to closed because of cancellation", situation.Id())
			situation.RecordedAt = cancellation.RecordedAtTime()
			situation.Progress = model.SituationProgressClosed
			connector.partner.Model().Situations().Save(&situation)
		}
	}
}

func (connector *SIRIGeneralMessageSubscriptionCollector) SetGeneralMessageSubscriber(generalMessageSubscriber SIRIGeneralMessageSubscriber) {
	connector.generalMessageSubscriber = generalMessageSubscriber
}

func (connector *SIRIGeneralMessageSubscriptionCollector) SetSituationUpdateSubscriber(updateSubscriber UpdateSubscriber) {
	connector.updateSubscriber = updateSubscriber
}

func (connector *SIRIGeneralMessageSubscriptionCollector) broadcastSituationUpdateEvent(events *CollectUpdateEvents) {
	if connector.updateSubscriber == nil {
		return
	}

	for _, e := range events.Situations {
		connector.updateSubscriber(e)
	}
}
