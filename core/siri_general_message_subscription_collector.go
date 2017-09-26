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
	model.Stopable
	model.Startable

	RequestSituationUpdate(request *SituationUpdateRequest)
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

func (connector *SIRIGeneralMessageSubscriptionCollector) RequestSituationUpdate(request *SituationUpdateRequest) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	line, ok := tx.Model().Lines().Find(request.LineId())
	if !ok {
		logger.Log.Debugf("SituationUpdateRequest in GeneralMessageSubscriptionCollector for unknown Line %v", request.LineId())
		return
	}

	objectidKind := connector.Partner().Setting("remote_objectid_kind")
	lineObjectid, ok := line.ObjectID(objectidKind)
	if !ok {
		logger.Log.Debugf("Requested line %v doesn't have and objectId of kind %v", request.LineId(), objectidKind)
		return
	}

	subscription, _ := connector.partner.Subscriptions().FindOrCreateByKind("GeneralMessageCollect")

	// Check if we find the resource
	resource := subscription.Resource(lineObjectid)
	if resource != nil {
		if !resource.SubscribedAt.IsZero() {
			resource.SubscribedUntil = resource.SubscribedUntil.Add(1 * time.Minute)
		}
		return
	}

	// Else we create a new resource
	ref := model.Reference{
		ObjectId: &lineObjectid,
		Type:     "Line",
	}

	subscription.CreateAddNewResource(ref)
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
		if subscription.Kind() != "GeneralMessageCollect" {
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
		RequestorRef:     connector.partner.ProducerRef(),
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
	logStashEvent["type"] = "NotifyGeneralMessageCollected"
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
