package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type VehicleMonitoringRequestCollector interface {
	state.Startable

	RequestVehicleUpdate(request *VehicleUpdateRequest)
}

type SIRIVehicleMonitoringRequestCollector struct {
	connector

	updateSubscriber UpdateSubscriber
}

type SIRIVehicleMonitoringRequestCollectorFactory struct{}

func NewSIRIVehicleMonitoringRequestCollector(partner *Partner) *SIRIVehicleMonitoringRequestCollector {
	connector := &SIRIVehicleMonitoringRequestCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent

	return connector
}

func (connector *SIRIVehicleMonitoringRequestCollector) Start() {
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace()
}

func (connector *SIRIVehicleMonitoringRequestCollector) RequestVehicleUpdate(request *VehicleUpdateRequest) {
	line, ok := connector.partner.Model().Lines().Find(request.LineId())
	if !ok {
		logger.Log.Debugf("VehicleUpdateRequest in VehicleMonitoringRequestCollector for unknown Line %v", request.LineId())
		return
	}

	codeSpace := connector.remoteCodeSpace
	code, ok := line.Code(codeSpace)
	if !ok {
		logger.Log.Debugf("Requested line %v doesn't have a code with codeSpace %v", request.LineId(), codeSpace)
		return
	}

	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	siriVehicleMonitoringRequest := &siri.SIRIGetVehicleMonitoringRequest{
		RequestorRef: connector.Partner().RequestorRef(),
	}
	siriVehicleMonitoringRequest.MessageIdentifier = connector.Partner().NewMessageIdentifier()
	siriVehicleMonitoringRequest.LineRef = code.Value()
	siriVehicleMonitoringRequest.RequestTimestamp = connector.Clock().Now()

	connector.logSIRIVehicleMonitoringRequest(message, siriVehicleMonitoringRequest)

	xmlVehicleMonitoringResponse, err := connector.Partner().SIRIClient().VehicleMonitoring(siriVehicleMonitoringRequest)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during VehicleMonitoring request: %v", err)
		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLVehicleMonitoringResponse(message, xmlVehicleMonitoringResponse)

	builder := NewVehicleMonitoringUpdateEventBuilder(connector.partner)

	for _, delivery := range xmlVehicleMonitoringResponse.VehicleMonitoringDeliveries() {
		if !delivery.Status() {
			continue
		}
		builder.SetUpdateEvents(delivery.VehicleActivities())
	}

	updateEvents := builder.UpdateEvents()

	// Log VehicleRefs
	message.Lines = updateEvents.GetLines()
	message.StopAreas = updateEvents.GetStopAreas()
	message.VehicleJourneys = updateEvents.GetVehicleJourneys()
	message.Vehicles = updateEvents.GetVehicles()

	// Broadcast all events
	connector.broadcastUpdateEvents(&updateEvents)
}

func (connector *SIRIVehicleMonitoringRequestCollector) broadcastUpdateEvents(events *VehicleMonitoringUpdateEvents) {
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
	for _, e := range events.Vehicles {
		connector.updateSubscriber(e)
	}
}

func (connector *SIRIVehicleMonitoringRequestCollector) SetUpdateSubscriber(updateSubscriber UpdateSubscriber) {
	connector.updateSubscriber = updateSubscriber
}

func (connector *SIRIVehicleMonitoringRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.VEHICLE_MONITORING_REQUEST,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (factory *SIRIVehicleMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRIVehicleMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIVehicleMonitoringRequestCollector(partner)
}

func (connector *SIRIVehicleMonitoringRequestCollector) logSIRIVehicleMonitoringRequest(message *audit.BigQueryMessage, request *siri.SIRIGetVehicleMonitoringRequest) {
	message.RequestIdentifier = request.MessageIdentifier

	xml, err := request.BuildXML(connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLVehicleMonitoringResponse(message *audit.BigQueryMessage, response *sxml.XMLVehicleMonitoringResponse) {
	for _, delivery := range response.VehicleMonitoringDeliveries() {
		if !delivery.Status() {
			message.Status = "Error"
		}
	}
	message.ResponseIdentifier = response.ResponseMessageIdentifier()
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
