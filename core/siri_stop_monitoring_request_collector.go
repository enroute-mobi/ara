package core

import (
	"fmt"
	"strings"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringRequestCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest) (*model.StopAreaUpdateEvent, error)
}

type TestStopMonitoringRequestCollector struct {
	model.UUIDConsumer
}

type TestStopMonitoringRequestCollectorFactory struct{}

type SIRIStopMonitoringRequestCollector struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector
}

type SIRIStopMonitoringRequestCollectorFactory struct{}

func NewTestStopMonitoringRequestCollector() *TestStopMonitoringRequestCollector {
	return &TestStopMonitoringRequestCollector{}
}

// WIP
func (connector *TestStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) (*model.StopAreaUpdateEvent, error) {
	stopAreaUpdateEvent := model.NewStopAreaUpdateEvent(connector.NewUUID())
	stopAreaUpdateEvent.StopVisitUpdateEvents = append(stopAreaUpdateEvent.StopVisitUpdateEvents, &model.StopVisitUpdateEvent{})
	return stopAreaUpdateEvent, nil
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
	return siriStopMonitoringRequestCollector
}

func (connector *SIRIStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) (*model.StopAreaUpdateEvent, error) {
	logStashEvent := make(audit.LogStashEvent)
	startTime := connector.Clock().Now()

	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	stopArea, ok := connector.Partner().Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		return nil, fmt.Errorf("StopArea not found")
	}
	objectid, ok := stopArea.ObjectID(connector.partner.Setting("remote_objectid_kind"))
	if !ok {
		return nil, fmt.Errorf("StopArea %s doesn't have an ojbectID of type %s", stopArea.Id(), connector.partner.Setting("remote_objectid_kind"))
	}

	siriStopMonitoringRequest := &siri.SIRIStopMonitoringRequest{
		MessageIdentifier: connector.SIRIPartner().NewMessageIdentifier(),
		MonitoringRef:     objectid.Value(),
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  connector.Clock().Now(),
	}

	logStopMonitoringRequest(logStashEvent, siriStopMonitoringRequest)

	xmlStopMonitoringResponse, err := connector.SIRIPartner().SOAPClient().StopMonitoring(siriStopMonitoringRequest)
	logStashEvent["responseTi	me"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["response"] = fmt.Sprintf("Error during CheckStatus: %v", err)
		return nil, err
	}

	logStopMonitoringResponse(logStashEvent, xmlStopMonitoringResponse)

	// WIP
	stopAreaUpdateEvent := model.NewStopAreaUpdateEvent(connector.NewUUID())
	connector.setStopVisitUpdateEvents(stopAreaUpdateEvent, xmlStopMonitoringResponse)

	logStopVisitUpdateEvents(logStashEvent, stopAreaUpdateEvent)

	return stopAreaUpdateEvent, nil
}

func (connector *SIRIStopMonitoringRequestCollector) setStopVisitUpdateEvents(event *model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringResponse) {
	xmlStopVisitEvents := xmlResponse.XMLMonitoredStopVisits()
	if len(xmlStopVisitEvents) == 0 {
		return
	}
	for _, xmlStopVisitEvent := range xmlStopVisitEvents {
		stopVisitEvent := &model.StopVisitUpdateEvent{
			Id:                connector.NewUUID(),
			Created_at:        connector.Clock().Now(),
			RecordedAt:        xmlStopVisitEvent.RecordedAt(),
			StopVisitObjectid: model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.ItemIdentifier()),
			StopAreaObjectId:  model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.StopPointRef()),
			Schedules:         make(model.StopVisitSchedules),
			DepartureStatus:   model.StopVisitDepartureStatus(xmlStopVisitEvent.DepartureStatus()),
			ArrivalStatuts:    model.StopVisitArrivalStatus(xmlStopVisitEvent.ArrivalStatus()),
			Attributes:        NewSIRIStopVisitUpdateAttributes(xmlStopVisitEvent, connector.partner.Setting("remote_objectid_kind")),
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

func (factory *SIRIStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func (factory *SIRIStopMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringRequestCollector(partner)
}

func logStopMonitoringRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIStopMonitoringRequest) {
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

func logStopMonitoringResponse(logStashEvent audit.LogStashEvent, response *siri.XMLStopMonitoringResponse) {
	logStashEvent["address"] = response.Address()
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["responseXML"] = response.RawXML()
}

func logStopVisitUpdateEvents(logStashEvent audit.LogStashEvent, stopAreaUpdateEvent *model.StopAreaUpdateEvent) {
	var idArray []string
	for _, stopVisitUpdateEvent := range stopAreaUpdateEvent.StopVisitUpdateEvents {
		idArray = append(idArray, stopVisitUpdateEvent.Id)
	}
	logStashEvent["StopVisitUpdateEventIds"] = strings.Join(idArray, ", ")
}
