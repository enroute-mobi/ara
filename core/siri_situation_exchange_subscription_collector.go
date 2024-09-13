package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
	"golang.org/x/exp/maps"
)

type SituationExchangeSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestAllSituationsUpdate()
	RequestSituationUpdate(kind string, requestedId model.Code)
	HandleNotifySituationExchange(notify *sxml.XMLNotifySituationExchange) *CollectedRefs
}

type SIRISituationExchangeSubscriptionCollector struct {
	connector

	deletedSubscriptions        *DeletedSubscriptions
	situationExchangeSubscriber SIRISituationExchangeSubscriber
	updateSubscriber            UpdateSubscriber
}

type SIRISituationExchangeSubscriptionCollectorFactory struct{}

func (factory *SIRISituationExchangeSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRISituationExchangeSubscriptionCollector(partner)
}

func (factory *SIRISituationExchangeSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRISituationExchangeSubscriptionCollector(partner *Partner) *SIRISituationExchangeSubscriptionCollector {
	connector := &SIRISituationExchangeSubscriptionCollector{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace()
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent
	connector.situationExchangeSubscriber = NewSIRISituationExchangeSubscriber(connector)

	return connector
}

func (connector *SIRISituationExchangeSubscriptionCollector) Stop() {
	connector.situationExchangeSubscriber.Stop()
}

func (connector *SIRISituationExchangeSubscriptionCollector) Start() {
	connector.deletedSubscriptions = NewDeletedSubscriptions()
	connector.situationExchangeSubscriber.Start()
}

func (connector *SIRISituationExchangeSubscriptionCollector) RequestAllSituationsUpdate() {
	obj := model.NewCode(SituationExchangeCollect, "all")
	connector.RequestSituationUpdate("all", obj)
}

func (connector *SIRISituationExchangeSubscriptionCollector) RequestSituationUpdate(kind string, requestedCode model.Code) {
	// Try to find a Subscription with the resource
	subscriptions := connector.partner.Subscriptions().FindByResourceId(requestedCode.String(), SituationExchangeCollect)
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
	newSubscription := connector.partner.Subscriptions().FindOrCreateByKind(SituationExchangeCollect)
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

func (connector *SIRISituationExchangeSubscriptionCollector) HandleNotifySituationExchange(notify *sxml.XMLNotifySituationExchange) (collectedRefs *CollectedRefs) {
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})

	updateEvents := NewCollectUpdateEvents()
	builder := NewSituationExchangeUpdateEventBuilder(connector.partner)

	collectedRefs = NewCollectedRefs()
	for _, delivery := range notify.SituationExchangesDeliveries() {
		subscriptionId := delivery.SubscriptionRef()

		if subscriptionId == "" {
			logger.Log.Debugf("Partner %s sent a NotifySituationExchange with an empty SubscriptionRef\n", connector.Partner().Slug())
			continue
		}

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Printf("Partner %s sent a NotifySituationExchange to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			if !connector.deletedSubscriptions.AlreadySend(subscriptionId) {
				subToDelete[subscriptionId] = struct{}{}
			}
			continue
		}

		if subscription.Kind() != SituationExchangeCollect {
			logger.Log.Printf("Partner %s sent a NotifySituationExchange to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind SituationExchangeCollect"
			continue
		}

		builder.SetSituationExchangeDeliveryUpdateEvents(updateEvents, delivery, notify.ProducerRef())

		maps.Copy(collectedRefs.LineRefs, builder.LineRefs)
		maps.Copy(collectedRefs.MonitoringRefs, builder.MonitoringRefs)

		connector.broadcastSituationUpdateEvent(updateEvents)
	}

	for subId := range subToDelete {
		CancelSubscription(subId, "SituationExchangeSubscriptionCollector", connector)
	}

	return collectedRefs
}

func (connector *SIRISituationExchangeSubscriptionCollector) SetSituationExchangeSubscriber(situationExchangeSubscriber SIRISituationExchangeSubscriber) {
	connector.situationExchangeSubscriber = situationExchangeSubscriber
}

func (connector *SIRISituationExchangeSubscriptionCollector) SetSituationUpdateSubscriber(updateSubscriber UpdateSubscriber) {
	connector.updateSubscriber = updateSubscriber
}

func (connector *SIRISituationExchangeSubscriptionCollector) broadcastSituationUpdateEvent(event *CollectUpdateEvents) {
	if connector.updateSubscriber == nil {
		return
	}
	for _, e := range event.Situations {
		connector.updateSubscriber(e)
	}
}
