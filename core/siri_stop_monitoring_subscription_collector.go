package core

import (
	"strconv"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringSubscriptionCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
	HandleNotifyStopMonitoring(delivery *siri.XMLStopMonitoringResponse)
}

type SIRIStopMonitoringSubscriptionCollector struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	Partner                  Partner
	stopAreaUpdateSubscriber StopAreaUpdateSubscriber
}

type SIRIStopMonitoringSubscriptionCollectorFactory struct{}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func NewSIRIStopMonitoringSubscriptionCollector(partner *Partner) *SIRIStopMonitoringSubscriptionCollector {
	siriStopMonitoringSubscriptionCollector := &SIRIStopMonitoringSubscriptionCollector{}
	siriStopMonitoringSubscriptionCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriStopMonitoringSubscriptionCollector.stopAreaUpdateSubscriber = manager.BroadcastStopAreaUpdateEvent

	return siriStopMonitoringSubscriptionCollector
}

func (connector *SIRIStopMonitoringSubscriptionCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	subscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")

	for _, sr := range subscription.resourcesByObjectID {
		if sr.Reference.ObjectId.Value() == string(request.StopAreaId()) {
			sr.SubscribedUntil = sr.SubscribedUntil.Add(1 * time.Minute)
			return
		}
	}

	objId := model.NewObjectID("StopMonitoring", string(request.StopAreaId()))
	ref := model.Reference{
		ObjectId: &objId,
		Id:       string(request.StopAreaId()),
		Type:     "StopArea",
	}

	subscription.CreateAddNewResource(ref)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) SetStopAreaUpdateSubscriber(stopAreaUpdateSubscriber StopAreaUpdateSubscriber) {
	connector.stopAreaUpdateSubscriber = stopAreaUpdateSubscriber
}

func (connector *SIRIStopMonitoringSubscriptionCollector) broadcastStopAreaUpdateEvent(event *model.StopAreaUpdateEvent) {
	if connector.stopAreaUpdateSubscriber != nil {
		connector.stopAreaUpdateSubscriber(event)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) HandleNotifyStopMonitoring(delivery *siri.XMLStopMonitoringResponse) {
	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLStopMonitoringDelivery(logStashEvent, delivery)

	// TEMP, DON'T WORK
	stopAreaUpdateEvent := model.NewStopAreaUpdateEvent(connector.NewUUID(), "")

	builder := newStopVisitUpdateEventBuilder(connector.partner)
	builder.setStopVisitUpdateEvents(stopAreaUpdateEvent, delivery)

	connector.setStopVisitCancellationEvents(stopAreaUpdateEvent, delivery)

	logStopVisitUpdateEvents(logStashEvent, stopAreaUpdateEvent)

	connector.broadcastStopAreaUpdateEvent(stopAreaUpdateEvent)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) setStopVisitCancellationEvents(event *model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringResponse) {
	xmlStopVisitCancellationEvents := xmlResponse.XMLMonitoredStopVisitCancellations()
	if len(xmlStopVisitCancellationEvents) == 0 {
		return
	}

	for _, xmlStopVisitCancellationEvent := range xmlStopVisitCancellationEvents {
		stopVisitCancellationEvent := &model.StopVisitNotCollectedEvent{
			StopVisitObjectId: model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlStopVisitCancellationEvent.ItemRef()),
		}
		event.StopVisitNotCollectedEvents = append(event.StopVisitNotCollectedEvents, stopVisitCancellationEvent)
	}
}

func logXMLStopMonitoringDelivery(logStashEvent audit.LogStashEvent, delivery *siri.XMLStopMonitoringResponse) {
	logStashEvent["address"] = delivery.Address()
	logStashEvent["producerRef"] = delivery.ProducerRef()
	logStashEvent["requestMessageRef"] = delivery.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = delivery.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = delivery.ResponseTimestamp().String()
	logStashEvent["responseXML"] = delivery.RawXML()
	logStashEvent["status"] = strconv.FormatBool(delivery.Status())
	if !delivery.Status() {
		logStashEvent["errorType"] = delivery.ErrorType()
		logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber())
		logStashEvent["errorText"] = delivery.ErrorText()
		logStashEvent["errorDescription"] = delivery.ErrorDescription()
	}
}
