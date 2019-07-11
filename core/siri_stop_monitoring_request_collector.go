package core

import (
	"fmt"
	"strconv"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringRequestCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
}

type TestStopMonitoringRequestCollector struct {
	model.UUIDConsumer
}

type TestStopMonitoringRequestCollectorFactory struct{}

type SIRIStopMonitoringRequestCollector struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	stopAreaUpdateSubscriber StopAreaUpdateSubscriber
}

type SIRIStopMonitoringRequestCollectorFactory struct{}

func NewTestStopMonitoringRequestCollector() *TestStopMonitoringRequestCollector {
	return &TestStopMonitoringRequestCollector{}
}

func (connector *TestStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	stopAreaUpdateEvent := model.NewLegacyStopAreaUpdateEvent(connector.NewUUID(), request.StopAreaId())
	stopAreaUpdateEvent.StopVisitUpdateEvents = append(stopAreaUpdateEvent.StopVisitUpdateEvents, &model.StopVisitUpdateEvent{})
}

func (factory *TestStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestStopMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewTestStopMonitoringRequestCollector()
}

func NewSIRIStopMonitoringRequestCollector(partner *Partner) *SIRIStopMonitoringRequestCollector {
	siriStopMonitoringRequestCollector := &SIRIStopMonitoringRequestCollector{}
	siriStopMonitoringRequestCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriStopMonitoringRequestCollector.stopAreaUpdateSubscriber = manager.BroadcastLegacyStopAreaUpdateEvent

	return siriStopMonitoringRequestCollector
}

func (connector *SIRIStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	stopArea, ok := tx.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("StopAreaUpdateRequest in StopMonitoringRequestCollector for unknown StopArea %v", request.StopAreaId())
		return
	}

	objectidKind := connector.partner.Setting("remote_objectid_kind")
	objectid, ok := stopArea.ObjectID(objectidKind)
	if !ok {
		logger.Log.Debugf("Requested stopArea %v doesn't have and objectId of kind %v", request.StopAreaId(), objectidKind)
		return
	}

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)
	monitoringRefMap := make(map[string]struct{})
	startTime := connector.Clock().Now()

	siriStopMonitoringRequest := &siri.SIRIGetStopMonitoringRequest{
		RequestorRef: connector.SIRIPartner().RequestorRef(),
	}
	siriStopMonitoringRequest.MessageIdentifier = connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier()
	siriStopMonitoringRequest.MonitoringRef = objectid.Value()
	siriStopMonitoringRequest.RequestTimestamp = connector.Clock().Now()
	siriStopMonitoringRequest.StopVisitTypes = "all"

	logSIRIStopMonitoringRequest(logStashEvent, siriStopMonitoringRequest)

	xmlStopMonitoringResponse, err := connector.SIRIPartner().SOAPClient().StopMonitoring(siriStopMonitoringRequest)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = fmt.Sprintf("Error during StopMonitoring request: %v", err)
		return
	}

	logXMLStopMonitoringResponse(logStashEvent, xmlStopMonitoringResponse)

	stopAreaUpdateEvents := make(map[string]*model.LegacyStopAreaUpdateEvent)
	builder := newStopVisitUpdateEventBuilder(connector.partner, objectid)

	for _, delivery := range xmlStopMonitoringResponse.StopMonitoringDeliveries() {
		if connector.Partner().LogRequestStopMonitoringDeliveries() {
			deliveryLogStashEvent := connector.newLogStashEvent()
			logXMLRequestStopMonitoringDelivery(deliveryLogStashEvent, xmlStopMonitoringResponse.ResponseMessageIdentifier(), delivery)
			audit.CurrentLogStash().WriteEvent(deliveryLogStashEvent)
		}

		if !delivery.Status() {
			continue
		}
		builder.setStopVisitUpdateEvents(stopAreaUpdateEvents, delivery.XMLMonitoredStopVisits())
	}

	for _, event := range stopAreaUpdateEvents {
		monitoringRefMap[event.StopAreaAttributes.ObjectId.Value()] = struct{}{}
		event.SetId(connector.NewUUID())
		monitoredStopVisits := []model.ObjectID{}

		collectedStopArea, ok := tx.Model().StopAreas().FindByObjectId(event.StopAreaAttributes.ObjectId)
		if !ok {
			connector.broadcastLegacyStopAreaUpdateEvent(event)
			continue
		}

		event.StopAreaId = collectedStopArea.Id()

		for _, sv := range tx.Model().StopVisits().FindByStopAreaId(collectedStopArea.Id()) {
			if sv.IsCollected() {
				objectid, ok := sv.ObjectID(objectidKind)
				if ok {
					monitoredStopVisits = append(monitoredStopVisits, objectid)
				}
			}
		}
		connector.findAndSetStopVisitNotCollectedEvent(event, monitoredStopVisits)
		connector.broadcastLegacyStopAreaUpdateEvent(event)
	}
	logMonitoringRefsFromMap(logStashEvent, monitoringRefMap)
}

func (connector *SIRIStopMonitoringRequestCollector) broadcastLegacyStopAreaUpdateEvent(event *model.LegacyStopAreaUpdateEvent) {
	if connector.stopAreaUpdateSubscriber != nil {
		connector.stopAreaUpdateSubscriber(event)
	}
}

func (connector *SIRIStopMonitoringRequestCollector) findAndSetStopVisitNotCollectedEvent(event *model.LegacyStopAreaUpdateEvent, collectedStopVisitObjectIDs []model.ObjectID) {
	objId := make(map[model.ObjectID]bool)

	for _, stopVisitEvent := range event.StopVisitUpdateEvents {
		objId[stopVisitEvent.StopVisitObjectid] = true
	}

	for _, stopVisitObjectID := range collectedStopVisitObjectIDs {
		if _, ok := objId[stopVisitObjectID]; !ok {
			logger.Log.Debugf("Send StopVisitNotCollectedEvent for %v", stopVisitObjectID)
			event.StopVisitNotCollectedEvents = append(event.StopVisitNotCollectedEvents, &model.StopVisitNotCollectedEvent{StopVisitObjectId: stopVisitObjectID})
		}
	}
}

func (connector *SIRIStopMonitoringRequestCollector) SetStopAreaUpdateSubscriber(stopAreaUpdateSubscriber StopAreaUpdateSubscriber) {
	connector.stopAreaUpdateSubscriber = stopAreaUpdateSubscriber
}

func (connector *SIRIStopMonitoringRequestCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringRequestCollector"
	return event
}

func (factory *SIRIStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func (factory *SIRIStopMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringRequestCollector(partner)
}

func logSIRIStopMonitoringRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIGetStopMonitoringRequest) {
	logStashEvent["siriType"] = "StopMonitoringRequest"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier
	logStashEvent["monitoringRef"] = request.MonitoringRef
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	xml, err := request.BuildXML()
	if err != nil {
		logStashEvent["requestXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["requestXML"] = xml
}

func logXMLStopMonitoringResponse(logStashEvent audit.LogStashEvent, response *siri.XMLStopMonitoringResponse) {
	logStashEvent["address"] = response.Address()
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["responseXML"] = response.RawXML()
	status := "true"
	errorCount := 0
	for _, delivery := range response.StopMonitoringDeliveries() {
		if !delivery.Status() {
			status = "false"
			errorCount++
		}
	}
	logStashEvent["status"] = status
	logStashEvent["errorCount"] = strconv.Itoa(errorCount)
}

func logXMLRequestStopMonitoringDelivery(logStashEvent audit.LogStashEvent, parent string, delivery *siri.XMLStopMonitoringDelivery) {
	logStashEvent["siriType"] = "StopMonitoringRequestDelivery"
	logStashEvent["parentMessageIdentifier"] = parent
	logStashEvent["monitoringRef"] = delivery.MonitoringRef()
	logStashEvent["requestMessageRef"] = delivery.RequestMessageRef()
	logStashEvent["responseTimestamp"] = delivery.ResponseTimestamp().String()

	logStashEvent["status"] = strconv.FormatBool(delivery.Status())
	if !delivery.Status() {
		logStashEvent["errorType"] = delivery.ErrorType()
		if delivery.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber())
		}
		logStashEvent["errorText"] = delivery.ErrorText()
		logStashEvent["errorDescription"] = delivery.ErrorDescription()
	}
}
