package core

import (
	"fmt"
	"strconv"
	"strings"

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

// WIP
func (connector *TestStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	stopAreaUpdateEvent := model.NewStopAreaUpdateEvent(connector.NewUUID(), request.StopAreaId())
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
	siriStopMonitoringRequestCollector.stopAreaUpdateSubscriber = manager.BroadcastStopAreaUpdateEvent

	return siriStopMonitoringRequestCollector
}

func (connector *SIRIStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	logStashEvent := make(audit.LogStashEvent)
	startTime := connector.Clock().Now()

	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	stopArea, ok := connector.Partner().Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		return
	}

	objectidKind := connector.partner.Setting("remote_objectid_kind")
	objectid, ok := stopArea.ObjectID(objectidKind)
	if !ok {
		return
	}

	siriStopMonitoringRequest := &siri.SIRIStopMonitoringRequest{
		MessageIdentifier: connector.SIRIPartner().NewMessageIdentifier(),
		MonitoringRef:     objectid.Value(),
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  connector.Clock().Now(),
	}

	logSIRIStopMonitoringRequest(logStashEvent, siriStopMonitoringRequest)

	xmlStopMonitoringResponse, err := connector.SIRIPartner().SOAPClient().StopMonitoring(siriStopMonitoringRequest)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["response"] = fmt.Sprintf("Error during CheckStatus: %v", err)
		return
	}

	logXMLStopMonitoringResponse(logStashEvent, xmlStopMonitoringResponse)

	if !xmlStopMonitoringResponse.Status() {
		return
	}

	// WIP
	stopAreaUpdateEvent := model.NewStopAreaUpdateEvent(connector.NewUUID(), stopArea.Id())

	connector.setStopVisitUpdateEvents(stopAreaUpdateEvent, xmlStopMonitoringResponse)

	collectedStopVisitObjectIDs := []model.ObjectID{}
	for _, stopVisit := range connector.Partner().Model().StopVisits().FindByStopAreaId(stopArea.Id()) {
		if stopVisit.IsCollected() == true {
			objectId, ok := stopVisit.ObjectID(objectidKind)
			if ok {
				collectedStopVisitObjectIDs = append(collectedStopVisitObjectIDs, objectId)
			}
		}
	}

	connector.findAndSetStopVisitNotCollectedEvent(stopAreaUpdateEvent, collectedStopVisitObjectIDs)
	logStopVisitUpdateEvents(logStashEvent, stopAreaUpdateEvent)

	connector.broadcastStopAreaUpdateEvent(stopAreaUpdateEvent)
}

func (connector *SIRIStopMonitoringRequestCollector) broadcastStopAreaUpdateEvent(event *model.StopAreaUpdateEvent) {
	if connector.stopAreaUpdateSubscriber != nil {
		connector.stopAreaUpdateSubscriber(event)
	}
}

func (connector *SIRIStopMonitoringRequestCollector) findAndSetStopVisitNotCollectedEvent(event *model.StopAreaUpdateEvent, collectedStopVisitObjectIDs []model.ObjectID) {
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

func (connector *SIRIStopMonitoringRequestCollector) setStopVisitUpdateEvents(event *model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringResponse) {
	xmlStopVisitEvents := xmlResponse.XMLMonitoredStopVisits()
	if len(xmlStopVisitEvents) == 0 {
		return
	}

	for _, xmlStopVisitEvent := range xmlStopVisitEvents {
		stopVisitEvent := &model.StopVisitUpdateEvent{
			Id:                     connector.NewUUID(),
			Created_at:             connector.Clock().Now(),
			RecordedAt:             xmlStopVisitEvent.RecordedAt(),
			VehicleAtStop:          xmlStopVisitEvent.VehicleAtStop(),
			StopVisitObjectid:      model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.ItemIdentifier()),
			StopAreaObjectId:       model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.StopPointRef()),
			Schedules:              make(model.StopVisitSchedules),
			DepartureStatus:        model.StopVisitDepartureStatus(xmlStopVisitEvent.DepartureStatus()),
			ArrivalStatuts:         model.StopVisitArrivalStatus(xmlStopVisitEvent.ArrivalStatus()),
			DatedVehicleJourneyRef: xmlStopVisitEvent.DatedVehicleJourneyRef(),
			DestinationRef:         xmlStopVisitEvent.DestinationRef(),
			OriginRef:              xmlStopVisitEvent.OriginRef(),
			DestinationName:        xmlStopVisitEvent.DestinationName(),
			OriginName:             xmlStopVisitEvent.OriginName(),
			Attributes:             NewSIRIStopVisitUpdateAttributes(xmlStopVisitEvent, connector.partner.Setting("remote_objectid_kind")),
		}
		stopVisitEvent.Schedules = model.NewStopVisitSchedules()
		if !xmlStopVisitEvent.AimedDepartureTime().IsZero() || !xmlStopVisitEvent.AimedArrivalTime().IsZero() {
			stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, xmlStopVisitEvent.AimedDepartureTime(), xmlStopVisitEvent.AimedArrivalTime())
		}
		if !xmlStopVisitEvent.ExpectedDepartureTime().IsZero() || !xmlStopVisitEvent.ExpectedArrivalTime().IsZero() {
			stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, xmlStopVisitEvent.ExpectedDepartureTime(), xmlStopVisitEvent.ExpectedArrivalTime())
		}
		if !xmlStopVisitEvent.ActualDepartureTime().IsZero() || !xmlStopVisitEvent.ActualArrivalTime().IsZero() {
			stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, xmlStopVisitEvent.ActualDepartureTime(), xmlStopVisitEvent.ActualArrivalTime())
		}
		event.StopVisitUpdateEvents = append(event.StopVisitUpdateEvents, stopVisitEvent)
	}
}

func (connector *SIRIStopMonitoringRequestCollector) SetStopAreaUpdateSubscriber(stopAreaUpdateSubscriber StopAreaUpdateSubscriber) {
	connector.stopAreaUpdateSubscriber = stopAreaUpdateSubscriber
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

func logSIRIStopMonitoringRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIStopMonitoringRequest) {
	logStashEvent["Connector"] = "StopMonitoringRequestCollector"
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
	logStashEvent["status"] = strconv.FormatBool(response.Status())
	if !response.Status() {
		logStashEvent["errorType"] = response.ErrorType()
		logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber())
		logStashEvent["errorText"] = response.ErrorText()
		logStashEvent["errorDescription"] = response.ErrorDescription()
	}
}

func logStopVisitUpdateEvents(logStashEvent audit.LogStashEvent, stopAreaUpdateEvent *model.StopAreaUpdateEvent) {
	var idArray []string
	for _, stopVisitUpdateEvent := range stopAreaUpdateEvent.StopVisitUpdateEvents {
		idArray = append(idArray, stopVisitUpdateEvent.Id)
	}
	logStashEvent["StopVisitUpdateEventIds"] = strings.Join(idArray, ", ")
}
