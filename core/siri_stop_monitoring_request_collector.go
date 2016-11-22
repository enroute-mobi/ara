package core

import (
	"fmt"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringRequestCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest) (*StopAreaUpdateEvent, error)
}

type SIRIStopMonitoringRequestCollector struct {
	model.ClockConsumer

	partner       *SIRIPartner
	objectid_kind string
}

type SIRIStopMonitoringRequestCollectorFactory struct{}

type StopAreaUpdateEvent struct {
	// WIP
}

func NewStopAreaUpdateEvent(response *siri.XMLStopMonitoringResponse) *StopAreaUpdateEvent {
	return &StopAreaUpdateEvent{}
}

func NewSIRIStopMonitoringRequestCollector(partner *SIRIPartner) *SIRIStopMonitoringRequestCollector {
	return &SIRIStopMonitoringRequestCollector{
		partner:       partner,
		objectid_kind: partner.Partner().Setting("remote_objectid_kind"),
	}
}

func (connector *SIRIStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) (*StopAreaUpdateEvent, error) {
	stopArea, ok := connector.partner.Partner().Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		return nil, fmt.Errorf("StopArea not found")
	}
	objectid, ok := stopArea.ObjectID(connector.objectid_kind)
	if !ok {
		return nil, fmt.Errorf("stopArea doesn't have an ojbectID of type %s", connector.objectid_kind)
	}

	siriStopMonitoringRequest := &siri.SIRIStopMonitoringRequest{
		MessageIdentifier: connector.partner.NewMessageIdentifier(),
		MonitoringRef:     objectid.Value(),
		RequestorRef:      connector.partner.RequestorRef(),
		RequestTimestamp:  connector.Clock().Now(),
	}

	xmlStopMonitoringResponse, err := connector.partner.SOAPClient().StopMonitoring(siriStopMonitoringRequest)
	if err != nil {
		return nil, err
	}

	// WIP
	stopAreaUpdateEvent := NewStopAreaUpdateEvent(xmlStopMonitoringResponse)

	return stopAreaUpdateEvent, nil
}

func (factory *SIRIStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := true
	if !apiPartner.IsSettingDefined("remote_objectid_kind") {
		apiPartner.Errors = append(apiPartner.Errors, "StopMonitoringRequestCollector needs partner to have 'remote_objectid_kind' setting defined")
		ok = false
	}
	if !apiPartner.IsSettingDefined("remote_url") {
		apiPartner.Errors = append(apiPartner.Errors, "StopMonitoringRequestCollector needs partner to have 'remote_url' setting defined")
		ok = false
	}
	if !apiPartner.IsSettingDefined("remote_credential") {
		apiPartner.Errors = append(apiPartner.Errors, "StopMonitoringRequestCollector needs partner to have 'remote_credential' setting defined")
		ok = false
	}
	return ok
}

func (factory *SIRIStopMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringRequestCollector(NewSIRIPartner(partner))
}
