package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	ig "bitbucket.org/enroute-mobi/ara/core/identifier_generator"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type VehicleMonitoringRequestCollector interface {
	RequestVehicleUpdate(request *VehicleUpdateRequest)
}

type SIRIVehicleMonitoringRequestCollector struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

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

func (connector *SIRIVehicleMonitoringRequestCollector) RequestVehicleUpdate(request *VehicleUpdateRequest) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	line, ok := tx.Model().Lines().Find(request.LineId())
	if !ok {
		logger.Log.Debugf("VehicleUpdateRequest in VehicleMonitoringRequestCollector for unknown Line %v", request.LineId())
		return
	}

	objectidKind := connector.partner.RemoteObjectIDKind()
	objectid, ok := line.ObjectID(objectidKind)
	if !ok {
		logger.Log.Debugf("Requested line %v doesn't have and objectId of kind %v", request.LineId(), objectidKind)
		return
	}

	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	siriVehicleMonitoringRequest := &siri.SIRIGetVehicleMonitoringRequest{
		RequestorRef: connector.Partner().RequestorRef(),
	}
	siriVehicleMonitoringRequest.MessageIdentifier = connector.Partner().IdentifierGenerator(ig.MESSAGE_IDENTIFIER).NewMessageIdentifier()
	siriVehicleMonitoringRequest.LineRef = objectid.Value()
	siriVehicleMonitoringRequest.RequestTimestamp = connector.Clock().Now()

	logSIRIVehicleMonitoringRequest(message, siriVehicleMonitoringRequest)

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

	// Log MonitoringRefs
	logVehicleRefs(message, updateEvents.VehicleRefs)

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
		Type:      "VehicleMonitoringRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (factory *SIRIVehicleMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRIVehicleMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIVehicleMonitoringRequestCollector(partner)
}

func logSIRIVehicleMonitoringRequest(message *audit.BigQueryMessage, request *siri.SIRIGetVehicleMonitoringRequest) {
	message.RequestIdentifier = request.MessageIdentifier

	xml, err := request.BuildXML()
	if err != nil {
		return
	}
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLVehicleMonitoringResponse(message *audit.BigQueryMessage, response *siri.XMLVehicleMonitoringResponse) {
	message.ResponseIdentifier = response.ResponseMessageIdentifier()
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}

func logVehicleRefs(message *audit.BigQueryMessage, refs map[string]struct{}) {
	refSlice := make([]string, len(refs))
	i := 0
	for monitoringRef := range refs {
		refSlice[i] = monitoringRef
		i++
	}

	if message != nil {
		message.Vehicles = refSlice
	}
}
