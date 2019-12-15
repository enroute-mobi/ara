package core

import (
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/edwig/audit"
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type GeneralMessageSubscriptionCollector interface {
	model.Stopable
	model.Startable

	RequestAllSituationsUpdate()
	RequestSituationUpdate(kind string, requestedId model.ObjectID)
	HandleNotifyGeneralMessage(notify *siri.XMLNotifyGeneralMessage)
}

type SIRIGeneralMessageSubscriptionCollector struct {
	model.UUIDConsumer
	model.ClockConsumer

	siriConnector

	generalMessageSubscriber  SIRIGeneralMessageSubscriber
	situationUpdateSubscriber SituationUpdateSubscriber
}

type SIRIGeneralMessageSubscriptionCollectorFactory struct{}

func (factory *SIRIGeneralMessageSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageSubscriptionCollector(partner)
}

func (factory *SIRIGeneralMessageSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func NewSIRIGeneralMessageSubscriptionCollector(partner *Partner) *SIRIGeneralMessageSubscriptionCollector {
	connector := &SIRIGeneralMessageSubscriptionCollector{}
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
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})

	logXMLGeneralMessageDelivery(logStashEvent, notify)

	situationUpdateEvents := &[]*model.SituationUpdateEvent{}
	builder := NewGeneralMessageUpdateEventBuilder(connector.partner)

	for _, delivery := range notify.GeneralMessagesDeliveries() {
		subscriptionId := delivery.SubscriptionRef()
		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Printf("Partner %s sent a NotifyGeneralMessage to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			subToDelete[delivery.SubscriptionRef()] = struct{}{}
			continue
		}

		if subscription.Kind() != "GeneralMessageCollect" {
			logger.Log.Printf("Partner %s sent a NotifyGeneralMessage to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind StopMonitoringCollect"
			continue
		}
		connector.cancelGeneralMessage(delivery)

		builder.SetGeneralMessageDeliveryUpdateEvents(situationUpdateEvents, delivery, notify.ProducerRef())

		if len(subscriptionErrors) != 0 {
			logSubscriptionErrorsFromMap(logStashEvent, subscriptionErrors)
		}
		connector.broadcastSituationUpdateEvent(*situationUpdateEvents)
	}

	for subId := range subToDelete {
		connector.cancelSubscription(subId)
	}
}

func (connector *SIRIGeneralMessageSubscriptionCollector) cancelSubscription(subId string) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	request := &siri.SIRIDeleteSubscriptionRequest{
		RequestTimestamp:  connector.Clock().Now(),
		SubscriptionRef:   subId,
		RequestorRef:      connector.partner.ProducerRef(),
		MessageIdentifier: connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier(),
	}

	logSIRIDeleteSubscriptionRequest(logStashEvent, request, "GeneralMessageSubscriptionCollector")
	startTime := connector.Clock().Now()
	response, err := connector.SIRIPartner().SOAPClient().DeleteSubscription(request)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logger.Log.Debugf("Error while terminating subcription with id : %v error : %v", subId, err.Error())
		logStashEvent["status"] = "false"
		logStashEvent["response"] = fmt.Sprintf("Error during DeleteSubscription: %v", err)
		return
	}

	logXMLDeleteSubscriptionResponse(logStashEvent, response)
}

func (connector *SIRIGeneralMessageSubscriptionCollector) cancelGeneralMessage(xmlResponse *siri.XMLGeneralMessageDelivery) {
	xmlGmCancellations := xmlResponse.XMLGeneralMessagesCancellations()
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	if len(xmlGmCancellations) == 0 {
		return
	}

	for _, cancellation := range xmlGmCancellations {
		obj := model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), cancellation.InfoMessageIdentifier())
		situation, ok := tx.Model().Situations().FindByObjectId(obj)
		if ok {
			logger.Log.Debugf("Deleting situation %v cause of cancellation", situation.Id())
			tx.Model().Situations().Delete(&situation)
		}
	}
	tx.Commit()
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

func (connector *SIRIGeneralMessageSubscriptionCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "GeneralMessageSubscriptionCollector"
	return event
}

func logXMLGeneralMessageDelivery(logStashEvent audit.LogStashEvent, notify *siri.XMLNotifyGeneralMessage) {
	logStashEvent["siriType"] = "CollectedNotifyGeneralMessage"
	logStashEvent["address"] = notify.Address()
	logStashEvent["producerRef"] = notify.ProducerRef()
	logStashEvent["requestMessageRef"] = notify.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = notify.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = notify.ResponseTimestamp().String()
	logStashEvent["responseXML"] = notify.RawXML()

	status := "true"
	errorCount := 0
	for _, delivery := range notify.GeneralMessagesDeliveries() {
		if !delivery.Status() {
			status = "false"
			errorCount++
		}
	}
	logStashEvent["status"] = status
	logStashEvent["errorCount"] = strconv.Itoa(errorCount)
}
