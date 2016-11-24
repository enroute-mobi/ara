package core

import (
	"fmt"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringRequestCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest) (*StopAreaUpdateEvent, error)
}

type TestStopMonitoringRequestCollector struct {
}

type TestStopMonitoringRequestCollectorFactory struct{}

type SIRIStopMonitoringRequestCollector struct {
	model.ClockConsumer

	SIRIConnector

	objectid_kind string
}

type SIRIStopMonitoringRequestCollectorFactory struct{}

type StopAreaUpdateEvent struct {
	// WIP
}

func NewStopAreaUpdateEvent(response *siri.XMLStopMonitoringResponse) *StopAreaUpdateEvent {
	return &StopAreaUpdateEvent{}
}

func NewTestStopMonitoringRequestCollector() *TestStopMonitoringRequestCollector {
	return &TestStopMonitoringRequestCollector{}
}

// WIP
func (connector *TestStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) (*StopAreaUpdateEvent, error) {
	stopAreaUpdateEvent := NewStopAreaUpdateEvent(&siri.XMLStopMonitoringResponse{})
	return stopAreaUpdateEvent, nil
}

func (factory *TestStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestStopMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewTestStopMonitoringRequestCollector()
}

func NewSIRIStopMonitoringRequestCollector(partner *Partner) *SIRIStopMonitoringRequestCollector {
	siriStopMonitoringRequestCollector := &SIRIStopMonitoringRequestCollector{
		objectid_kind: partner.Setting("remote_objectid_kind"),
	}
	siriStopMonitoringRequestCollector.partner = partner
	return siriStopMonitoringRequestCollector
}

func (connector *SIRIStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) (*StopAreaUpdateEvent, error) {
	stopArea, ok := connector.Partner().Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		return nil, fmt.Errorf("StopArea not found")
	}
	objectid, ok := stopArea.ObjectID(connector.objectid_kind)
	if !ok {
		return nil, fmt.Errorf("stopArea doesn't have an ojbectID of type %s", connector.objectid_kind)
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
	stopAreaUpdateEvent := NewStopAreaUpdateEvent(xmlStopMonitoringResponse)

	return stopAreaUpdateEvent, nil
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