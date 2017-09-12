package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type GeneralMessageSubscriptionCollector interface {
	RequestSituationUpdate(request *SituationUpdateRequest)
	HandleNotifyGeneralMessage(notify *siri.XMLNotifyGeneralMessage)
}

type SIRIGeneralMessageSubscriptionCollector struct {
	model.UUIDConsumer
	model.ClockConsumer

	siriConnector

	situationUpdateSubscriber SituationUpdateSubscriber
}

type SIRIGeneralMessageSubscriptionCollectorFactory struct{}

func NewSIRIGeneralMessageSubscriptionCollector(partner *Partner) *SIRIGeneralMessageSubscriptionCollector {
	siriGeneralMessageSubscriptionCollector := &SIRIGeneralMessageSubscriptionCollector{}
	siriGeneralMessageSubscriptionCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriGeneralMessageSubscriptionCollector.situationUpdateSubscriber = manager.BroadcastSituationUpdateEvent

	return siriGeneralMessageSubscriptionCollector
}

func (connector *SIRIGeneralMessageSubscriptionCollector) RequestSituationUpdate(request *SituationUpdateRequest) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	subscription, found := connector.partner.Subscriptions().FindOrCreateByKind("GeneralMessage")
	if found {
		return
	}

	gmRequest := &siri.SIRIGeneralMessageSubscriptionRequest{
		ConsumerAddress:   connector.Partner().Address(),
		MessageIdentifier: connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier(),
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  connector.Clock().Now(),
	}

	entry := &siri.SIRIGeneralMessageSubscriptionRequestEntry{
		SubscriberRef:          connector.SIRIPartner().RequestorRef(),
		SubscriptionIdentifier: connector.SIRIPartner().IdentifierGenerator("subscription_identifier").NewMessageIdentifier(),
		InitialTerminationTime: connector.Clock().Now().Add(48 * time.Hour),
	}
	entry.MessageIdentifier = connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier()
	entry.RequestTimestamp = connector.Clock().Now()

	gmRequest.Entries = append(gmRequest.Entries, entry)

	logSIRIGeneralMessageSubscriptionRequest(logStashEvent, gmRequest)

	response, err := connector.SIRIPartner().SOAPClient().GeneralMessageSubscription(gmRequest)
	if err != nil {
		logger.Log.Debugf("Error while subscribing %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["response"] = fmt.Sprintf("Error during GeneralMessageSubscriptionRequest: %v", err)
		subscription.Delete()
		return
	}
	logXMLGeneralMessageSubscriptionResponse(logStashEvent, response)

	responseStatus := response.ResponseStatus()
	logXMLResponseStatus(logStashEvent, &responseStatus)

	if !responseStatus.Status() {
		logger.Log.Debugf("Subscription status false for General Message Subscription %v %v ", responseStatus.ErrorType(), responseStatus.ErrorText())
		subscription.Delete()
		return
	}
}

func (connector *SIRIGeneralMessageSubscriptionCollector) HandleNotifyGeneralMessage(notify *siri.XMLNotifyGeneralMessage) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLGeneralMessageDelivery(logStashEvent, notify)

	situationUpdateEvents := &[]*model.SituationUpdateEvent{}

	for _, delivery := range notify.GeneralMessagesDeliveries() {
		subscriptionId := delivery.SubscriptionRef()
		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))

		if ok == false {
			logger.Log.Printf("Partner %s sent a NotifyGeneralMessage response to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			connector.cancelSubscription(delivery.SubscriptionRef())
			continue
		}
		if subscription.Kind() != "GeneralMessage" {
			logger.Log.Printf("Partner %s sent a NotifyGeneralMessage response to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			continue
		}
		connector.cancelGeneralMessage(delivery)
		connector.setGeneralMessageUpdateEvents(situationUpdateEvents, delivery)

		logSituationUpdateEvents(logStashEvent, *situationUpdateEvents)

		connector.broadcastSituationUpdateEvent(*situationUpdateEvents)
	}
}

func (connector *SIRIGeneralMessageSubscriptionCollector) cancelSubscription(subId string) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	request := &siri.SIRIDeleteSubscriptionRequest{
		RequestTimestamp: connector.Clock().Now(),
		SubscriptionRef:  subId,
		RequestorRef:     connector.partner.RemoteObjectIDKind(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER),
	}

	response, err := connector.SIRIPartner().SOAPClient().DeleteSubscription(request)

	if err != nil {
		logger.Log.Debugf("Error while terminating subcription with id : %v error : ", subId, err.Error())
		logStashEvent["status"] = "false"
		logStashEvent["response"] = fmt.Sprintf("Error during DeleteSubscription: %v", err)
		return
	}

	logXMLDeleteSubscriptionResponse(logStashEvent, response) //siri_stop_monitoring_subscription_collector
}

func (connector *SIRIGeneralMessageSubscriptionCollector) setGeneralMessageUpdateEvents(events *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageDelivery) {
	builder := NewGeneralMessageUpdateEventBuilder(connector.partner)
	builder.SetGeneralMessageDeliveryUpdateEvents(events, xmlResponse)
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

func logSIRIGeneralMessageSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIGeneralMessageSubscriptionRequest) {
	logStashEvent["Type"] = "GeneralMessageSubscriptionRequest"
	logStashEvent["consumerAddress"] = request.ConsumerAddress
	logStashEvent["messageIdentifier"] = request.MessageIdentifier
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	xml, err := request.BuildXML()
	if err != nil {
		logStashEvent["requestXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["requestXML"] = xml
}

func logXMLGeneralMessageDelivery(logStashEvent audit.LogStashEvent, notify *siri.XMLNotifyGeneralMessage) {
	logStashEvent["address"] = notify.Address()
	logStashEvent["producerRef"] = notify.ProducerRef()
	logStashEvent["requestMessageRef"] = notify.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = notify.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = notify.ResponseTimestamp().String()
	logStashEvent["responseXML"] = notify.RawXML()
	logStashEvent["status"] = strconv.FormatBool(notify.Status())
	if !notify.Status() {
		logStashEvent["errorType"] = notify.ErrorType()
		if notify.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(notify.ErrorNumber())
		}
		logStashEvent["errorText"] = notify.ErrorText()
		logStashEvent["errorDescription"] = notify.ErrorDescription()
	}
}

func logXMLGeneralMessageSubscriptionResponse(logStashEvent audit.LogStashEvent, response *siri.XMLGeneralMessageSubscriptionResponse) {
	logStashEvent["address"] = response.Address()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responderRef"] = response.ResponderRef()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["serviceStartedTime"] = response.ServiceStartedTime().String()
	logStashEvent["responseXML"] = response.RawXML()
}

func logXMLResponseStatus(logStashEvent audit.LogStashEvent, responseStatus *siri.XMLResponseStatus) {
	logStashEvent["subscriberRef"] = responseStatus.SubscriberRef()
	logStashEvent["subscriptionRef"] = responseStatus.SubscriptionRef()
	logStashEvent["status"] = strconv.FormatBool(responseStatus.Status())
	if !responseStatus.Status() {
		logStashEvent["errorType"] = responseStatus.ErrorType()
		if responseStatus.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(responseStatus.ErrorNumber())
		}
		logStashEvent["errorText"] = responseStatus.ErrorText()
		logStashEvent["errorDescription"] = responseStatus.ErrorDescription()
	}
}
