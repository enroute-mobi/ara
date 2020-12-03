package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopMonitoringRequestCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
}

type TestStopMonitoringRequestCollector struct {
	uuid.UUIDConsumer
}

type TestStopMonitoringRequestCollectorFactory struct{}

type SIRIStopMonitoringRequestCollector struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	siriConnector

	updateSubscriber UpdateSubscriber
}

type SIRIStopMonitoringRequestCollectorFactory struct{}

func NewTestStopMonitoringRequestCollector() *TestStopMonitoringRequestCollector {
	return &TestStopMonitoringRequestCollector{}
}

func (connector *TestStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
}

func (factory *TestStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestStopMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewTestStopMonitoringRequestCollector()
}

func NewSIRIStopMonitoringRequestCollector(partner *Partner) *SIRIStopMonitoringRequestCollector {
	connector := &SIRIStopMonitoringRequestCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent

	return connector
}

func (connector *SIRIStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	stopArea, ok := tx.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("StopAreaUpdateRequest in StopMonitoringRequestCollector for unknown StopArea %v", request.StopAreaId())
		return
	}

	objectidKind := connector.partner.Setting(REMOTE_OBJECTID_KIND)
	objectid, ok := stopArea.ObjectID(objectidKind)
	if !ok {
		logger.Log.Debugf("Requested stopArea %v doesn't have and objectId of kind %v", request.StopAreaId(), objectidKind)
		return
	}

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)
	startTime := connector.Clock().Now()

	siriStopMonitoringRequest := &siri.SIRIGetStopMonitoringRequest{
		RequestorRef: connector.SIRIPartner().RequestorRef(),
	}
	siriStopMonitoringRequest.MessageIdentifier = connector.Partner().IdentifierGenerator(MESSAGE_IDENTIFIER).NewMessageIdentifier()
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

	builder := NewStopMonitoringUpdateEventBuilder(connector.partner, objectid)

	for _, delivery := range xmlStopMonitoringResponse.StopMonitoringDeliveries() {
		if connector.Partner().LogRequestStopMonitoringDeliveries() {
			deliveryLogStashEvent := connector.newLogStashEvent()
			logXMLRequestStopMonitoringDelivery(deliveryLogStashEvent, xmlStopMonitoringResponse.ResponseMessageIdentifier(), delivery)
			audit.CurrentLogStash().WriteEvent(deliveryLogStashEvent)
		}

		if !delivery.Status() {
			continue
		}
		builder.SetUpdateEvents(delivery.XMLMonitoredStopVisits())
	}

	updateEvents := builder.UpdateEvents()
	logger.Log.Printf("%v", updateEvents)

	// Log MonitoringRefs
	logMonitoringRefsFromMap(logStashEvent, updateEvents.MonitoringRefs)

	// Broadcast all events
	connector.broadcastUpdateEvents(&updateEvents)

	// Set all StopVisits not in the response not collected
	monitoredStopVisits := []model.ObjectID{}

	for stopPointRef, events := range updateEvents.StopVisits {
		sa, ok := tx.Model().StopAreas().FindByObjectId(model.NewObjectID(objectidKind, stopPointRef))
		if !ok {
			continue
		}

		for _, sv := range tx.Model().StopVisits().FindMonitoredByOriginByStopAreaId(sa.Id(), string(connector.Partner().Slug())) {
			objectid, ok := sv.ObjectID(objectidKind)
			if ok {
				monitoredStopVisits = append(monitoredStopVisits, objectid)
			}
		}

		connector.broadcastNotCollectedEvents(events, monitoredStopVisits)
	}
}

func (connector *SIRIStopMonitoringRequestCollector) broadcastUpdateEvents(events *StopMonitoringUpdateEvents) {
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

func (connector *SIRIStopMonitoringRequestCollector) broadcastUpdateEvent(event model.UpdateEvent) {
	if connector.updateSubscriber != nil {
		connector.updateSubscriber(event)
	}
}

func (connector *SIRIStopMonitoringRequestCollector) broadcastNotCollectedEvents(events map[string]*model.StopVisitUpdateEvent, collectedStopVisitObjectIDs []model.ObjectID) {
	for _, stopVisitObjectID := range collectedStopVisitObjectIDs {
		if _, ok := events[stopVisitObjectID.Value()]; !ok {
			logger.Log.Debugf("Send StopVisitNotCollectedEvent for %v", stopVisitObjectID)
			connector.broadcastUpdateEvent(model.NewNotCollectedUpdateEvent(stopVisitObjectID))
		}
	}
}

func (connector *SIRIStopMonitoringRequestCollector) SetUpdateSubscriber(updateSubscriber UpdateSubscriber) {
	connector.updateSubscriber = updateSubscriber
}

func (connector *SIRIStopMonitoringRequestCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringRequestCollector"
	return event
}

func (factory *SIRIStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	ok = ok && apiPartner.ValidatePresenceOfSetting(REMOTE_URL)
	ok = ok && apiPartner.ValidatePresenceOfSetting(REMOTE_CREDENTIAL)
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
