package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type FacilityMonitoringRequestCollector interface {
	state.Startable

	RequestFacilityUpdate(request *FacilityUpdateRequest)
}

type SIRIFacilityMonitoringRequestCollector struct {
	connector

	updateSubscriber UpdateSubscriber
}

type SIRIFacilityMonitoringRequestCollectorFactory struct{}

func NewSIRIFacilityMonitoringRequestCollector(partner *Partner) *SIRIFacilityMonitoringRequestCollector {
	connector := &SIRIFacilityMonitoringRequestCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent

	return connector
}

func (connector *SIRIFacilityMonitoringRequestCollector) Start() {
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace()
}

func (connector *SIRIFacilityMonitoringRequestCollector) RequestFacilityUpdate(request *FacilityUpdateRequest) {
	facility, ok := connector.partner.Model().Facilities().Find(request.FacilityId())
	if !ok {
		logger.Log.Debugf("FacilityUpdateRequest in FacilityMonitoringRequestCollector for unknown Facility %v", request.FacilityId())
		return
	}

	codeSpace := connector.remoteCodeSpace
	code, ok := facility.Code(codeSpace)
	if !ok {
		logger.Log.Debugf("Requested facility %v doesn't have a code with codeSpace %v", request.FacilityId(), codeSpace)
		return
	}

	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	siriFacilityMonitoringRequest := &siri.SIRIGetFacilityMonitoringRequest{
		RequestorRef: connector.Partner().RequestorRef(),
	}
	siriFacilityMonitoringRequest.MessageIdentifier = connector.Partner().NewMessageIdentifier()
	siriFacilityMonitoringRequest.FacilityRef = code.Value()
	siriFacilityMonitoringRequest.RequestTimestamp = connector.Clock().Now()

	connector.logSIRIFacilityMonitoringRequest(message, siriFacilityMonitoringRequest)

	xmlFacilityMonitoringResponse, err := connector.Partner().SIRIClient().FacilityMonitoring(siriFacilityMonitoringRequest)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during FacilityMonitoring request: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLFacilityMonitoringResponse(message, xmlFacilityMonitoringResponse)

	builder := NewFacilityMonitoringUpdateEventBuilder(connector.partner)

	for _, delivery := range xmlFacilityMonitoringResponse.FacilityMonitoringDeliveries() {
		// if !delivery.Status() {
		// 	continue
		// }
		//
		builder.SetUpdateEvents(delivery.FacilityConditions())
	}

	updateEvents := builder.UpdateEvents()

	// Log Models
	message.Facilities = updateEvents.GetFacilities()

	// Broadcast all events
	connector.broadcastUpdateEvents(&updateEvents)
}

func (connector *SIRIFacilityMonitoringRequestCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {
	if connector.updateSubscriber == nil {
		return
	}
	for _, e := range events.Facilities {
		connector.updateSubscriber(e)
	}
}

// func (connector *SIRIFacilityMonitoringRequestCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {
// 	if connector.updateSubscriber == nil {
// 		return
// 	}
// 	for _, e := range events.Facilitys {
// 		connector.updateSubscriber(e)
// 	}
// 	for _, e := range events.Lines {
// 		connector.updateSubscriber(e)
// 	}
// 	for _, e := range events.VehicleJourneys {
// 		connector.updateSubscriber(e)
// 	}
// 	for _, es := range events.StopVisits { // Stopvisits are map[MonitoringRef]map[ItemIdentifier]event
// 		for _, e := range es {
// 			connector.updateSubscriber(e)
// 		}
// 	}
// }

// func (connector *SIRIFacilityMonitoringRequestCollector) broadcastUpdateEvent(event model.UpdateEvent) {
// 	if connector.updateSubscriber != nil {
// 		connector.updateSubscriber(event)
// 	}
// }

// func (connector *SIRIFacilityMonitoringRequestCollector) SetUpdateSubscriber(updateSubscriber UpdateSubscriber) {
// 	connector.updateSubscriber = updateSubscriber
// }

func (connector *SIRIFacilityMonitoringRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.FACILITY_MONITORING_REQUEST,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (factory *SIRIFacilityMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRIFacilityMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIFacilityMonitoringRequestCollector(partner)
}

func (connector *SIRIFacilityMonitoringRequestCollector) logSIRIFacilityMonitoringRequest(message *audit.BigQueryMessage, request *siri.SIRIGetFacilityMonitoringRequest) {
	message.RequestIdentifier = request.MessageIdentifier

	xml, err := request.BuildXML(connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLFacilityMonitoringResponse(message *audit.BigQueryMessage, response *sxml.XMLFacilityMonitoringResponse) {
	message.ResponseIdentifier = response.ResponseMessageIdentifier()
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
