package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type GeneralMessageSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestAllSituationsUpdate()
	RequestSituationUpdate(kind string, requestedId model.ObjectID)
	HandleNotifyGeneralMessage(notify *siri.XMLNotifyGeneralMessage)
}

type SIRIGeneralMessageSubscriptionCollector struct {
	uuid.UUIDConsumer
	clock.ClockConsumer

	connector

	deletedSubscriptions      *DeletedSubscriptions
	generalMessageSubscriber  SIRIGeneralMessageSubscriber
	situationUpdateSubscriber SituationUpdateSubscriber
}

type SIRIGeneralMessageSubscriptionCollectorFactory struct{}

func (factory *SIRIGeneralMessageSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageSubscriptionCollector(partner)
}

func (factory *SIRIGeneralMessageSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIGeneralMessageSubscriptionCollector(partner *Partner) *SIRIGeneralMessageSubscriptionCollector {
	connector := &SIRIGeneralMessageSubscriptionCollector{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind()
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.situationUpdateSubscriber = manager.BroadcastSituationUpdateEvent
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
	obj := model.NewObjectID("generalMessageCollect", "all")
	connector.RequestSituationUpdate("all", obj)
}

func (connector *SIRIGeneralMessageSubscriptionCollector) RequestSituationUpdate(kind string, requestedObjectId model.ObjectID) {
	// Try to find a Subscription with the resource
	subscriptions := connector.partner.Subscriptions().FindByResourceId(requestedObjectId.String(), "GeneralMessageCollect")
	if len(subscriptions) > 0 {
		for _, subscription := range subscriptions {
			resource := subscription.Resource(requestedObjectId)
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
	newSubscription := connector.partner.Subscriptions().FindOrCreateByKind("GeneralMessageCollect")
	ref := model.Reference{
		ObjectId: &requestedObjectId,
	}
	switch kind {
	case SITUATION_UPDATE_REQUEST_LINE:
		ref.Type = "Line"
	case SITUATION_UPDATE_REQUEST_STOP_AREA:
		ref.Type = "StopArea"
	}

	newSubscription.CreateAddNewResource(ref)
}

func (connector *SIRIGeneralMessageSubscriptionCollector) HandleNotifyGeneralMessage(notify *siri.XMLNotifyGeneralMessage) {
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})

	situationUpdateEvents := &[]*model.SituationUpdateEvent{}
	builder := NewGeneralMessageUpdateEventBuilder(connector.partner)

	for _, delivery := range notify.GeneralMessagesDeliveries() {
		subscriptionId := delivery.SubscriptionRef()
		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Printf("Partner %s sent a NotifyGeneralMessage to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			if !connector.deletedSubscriptions.AlreadySend(subscriptionId) {
				subToDelete[subscriptionId] = struct{}{}
			}
			continue
		}

		if subscription.Kind() != "GeneralMessageCollect" {
			logger.Log.Printf("Partner %s sent a NotifyGeneralMessage to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind StopMonitoringCollect"
			continue
		}
		connector.cancelGeneralMessage(delivery)

		builder.SetGeneralMessageDeliveryUpdateEvents(situationUpdateEvents, delivery, notify.ProducerRef())

		connector.broadcastSituationUpdateEvent(*situationUpdateEvents)
	}

	for subId := range subToDelete {
		connector.cancelSubscription(subId)
	}
}

func (connector *SIRIGeneralMessageSubscriptionCollector) cancelSubscription(subId string) {
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	request := &siri.SIRIDeleteSubscriptionRequest{
		RequestTimestamp:  connector.Clock().Now(),
		SubscriptionRef:   subId,
		RequestorRef:      connector.partner.ProducerRef(),
		MessageIdentifier: connector.Partner().NewMessageIdentifier(),
	}

	logSIRIDeleteSubscriptionRequest(message, request, "GeneralMessageSubscriptionCollector", connector.Partner().SIRIEnvelopeType())
	startTime := connector.Clock().Now()
	response, err := connector.Partner().SIRIClient().DeleteSubscription(request)

	responseTime := connector.Clock().Since(startTime)
	message.ProcessingTime = responseTime.Seconds()

	if err != nil {
		logger.Log.Debugf("Error while terminating subcription with id : %v error : %v", subId, err.Error())
		e := fmt.Sprintf("Error during DeleteSubscription: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLDeleteSubscriptionResponse(message, response)
}

func (connector *SIRIGeneralMessageSubscriptionCollector) cancelGeneralMessage(xmlResponse *siri.XMLGeneralMessageDelivery) {
	xmlGmCancellations := xmlResponse.XMLGeneralMessagesCancellations()

	if len(xmlGmCancellations) == 0 {
		return
	}

	for _, cancellation := range xmlGmCancellations {
		obj := model.NewObjectID(connector.remoteObjectidKind, cancellation.InfoMessageIdentifier())
		situation, ok := connector.partner.Model().Situations().FindByObjectId(obj)
		if ok {
			logger.Log.Debugf("Deleting situation %v cause of cancellation", situation.Id())
			connector.partner.Model().Situations().Delete(&situation)
		}
	}
}

func (connector *SIRIGeneralMessageSubscriptionCollector) SetGeneralMessageSubscriber(generalMessageSubscriber SIRIGeneralMessageSubscriber) {
	connector.generalMessageSubscriber = generalMessageSubscriber
}

func (connector *SIRIGeneralMessageSubscriptionCollector) SetSituationUpdateSubscriber(situationUpdateSubscriber SituationUpdateSubscriber) {
	connector.situationUpdateSubscriber = situationUpdateSubscriber
}

func (connector *SIRIGeneralMessageSubscriptionCollector) broadcastSituationUpdateEvent(event []*model.SituationUpdateEvent) {
	if connector.situationUpdateSubscriber != nil {
		connector.situationUpdateSubscriber(event)
	}
}

func (connector *SIRIGeneralMessageSubscriptionCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}
