package core

import (
	"fmt"

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

	SIRIConnector
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

	xmlStopMonitoringResponse, err := connector.SIRIPartner().SOAPClient().StopMonitoring(siriStopMonitoringRequest)
	if err != nil {
		return nil, err
	}

	// WIP
	stopAreaUpdateEvent := model.NewStopAreaUpdateEvent(connector.NewUUID())
	connector.setStopVisitUpdateEvents(stopAreaUpdateEvent, xmlStopMonitoringResponse)

	return stopAreaUpdateEvent, nil
}

func (connector *SIRIStopMonitoringRequestCollector) setStopVisitUpdateEvents(event *model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringResponse) {
	xmlStopVisitEvents := xmlResponse.XMLMonitoredStopVisits()
	if len(xmlStopVisitEvents) == 0 {
		return
	}
	for _, xmlStopVisitEvent := range xmlStopVisitEvents {
		stopVisitEvent := &model.StopVisitUpdateEvent{
			Id:                  connector.NewUUID(),
			Created_at:          connector.Clock().Now(),
			Stop_visit_objectid: model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlStopVisitEvent.ItemIdentifier()),
			Schedules:           make(model.StopVisitSchedules),
			DepartureStatus:     model.StopVisitDepartureStatus(xmlStopVisitEvent.DepartureStatus()),
			ArrivalStatuts:      model.StopVisitArrivalStatus(xmlStopVisitEvent.ArrivalStatus()),
			StopVisitAttributes: NewSIRIStopMonitoringStopVisitAttributes(xmlStopVisitEvent, connector.partner.Setting("remote_objectid_kind")),
		}
		stopVisitEvent.Schedules = model.NewStopVisitSchedules()
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, xmlStopVisitEvent.AimedDepartureTime(), xmlStopVisitEvent.AimedArrivalTime())
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, xmlStopVisitEvent.ExpectedDepartureTime(), xmlStopVisitEvent.ExpectedArrivalTime())
		stopVisitEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, xmlStopVisitEvent.ActualDepartureTime(), xmlStopVisitEvent.ActualArrivalTime())
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
