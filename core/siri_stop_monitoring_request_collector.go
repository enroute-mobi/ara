package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type StopMonitoringRequestCollector interface {
	state.Startable

	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
}

type TestStopMonitoringRequestCollector struct {
	connector
}

type TestStopMonitoringRequestCollectorFactory struct{}

type SIRIStopMonitoringRequestCollector struct {
	connector

	updateSubscriber UpdateSubscriber
}

type SIRIStopMonitoringRequestCollectorFactory struct{}

func NewTestStopMonitoringRequestCollector() *TestStopMonitoringRequestCollector {
	return &TestStopMonitoringRequestCollector{}
}

func (connector *TestStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
}

func (factory *TestStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) {} // Always valid

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

func (connector *SIRIStopMonitoringRequestCollector) Start() {
	connector.remoteObjectidKind = connector.partner.RemoteObjectIDKind()
}

func (connector *SIRIStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	stopArea, ok := connector.partner.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("StopAreaUpdateRequest in StopMonitoringRequestCollector for unknown StopArea %v", request.StopAreaId())
		return
	}

	objectidKind := connector.remoteObjectidKind
	objectid, ok := stopArea.ObjectID(objectidKind)
	if !ok {
		logger.Log.Debugf("Requested stopArea %v doesn't have and objectId of kind %v", request.StopAreaId(), objectidKind)
		return
	}

	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	siriStopMonitoringRequest := &siri.SIRIGetStopMonitoringRequest{
		RequestorRef: connector.Partner().RequestorRef(),
	}
	siriStopMonitoringRequest.MessageIdentifier = connector.Partner().NewMessageIdentifier()
	siriStopMonitoringRequest.MonitoringRef = objectid.Value()
	siriStopMonitoringRequest.RequestTimestamp = connector.Clock().Now()
	siriStopMonitoringRequest.StopVisitTypes = "all"

	connector.logSIRIStopMonitoringRequest(message, siriStopMonitoringRequest)

	xmlStopMonitoringResponse, err := connector.Partner().SIRIClient().StopMonitoring(siriStopMonitoringRequest)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during StopMonitoring request: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLStopMonitoringResponse(message, xmlStopMonitoringResponse)

	builder := NewStopMonitoringUpdateEventBuilder(connector.partner, objectid)

	for _, delivery := range xmlStopMonitoringResponse.StopMonitoringDeliveries() {
		if !delivery.Status() {
			continue
		}
		builder.SetUpdateEvents(delivery.XMLMonitoredStopVisits())
	}

	updateEvents := builder.UpdateEvents()

	// Log Models
	message.StopAreas = updateEvents.GetStopAreas()
	message.Lines = updateEvents.GetLines()
	message.VehicleJourneys = updateEvents.GetVehicleJourneys()

	// Broadcast all events
	connector.broadcastUpdateEvents(&updateEvents)

	// Set all StopVisits not in the response not collected
	monitoredStopVisits := []model.ObjectID{}

	for stopPointRef, events := range updateEvents.StopVisits {
		sa, ok := connector.partner.Model().StopAreas().FindByObjectId(model.NewObjectID(objectidKind, stopPointRef))
		if !ok {
			continue
		}

		svs := connector.partner.Model().StopVisits().FindMonitoredByOriginByStopAreaId(sa.Id(), string(connector.Partner().Slug()))
		for i := range svs {
			objectid, ok := svs[i].ObjectID(objectidKind)
			if ok {
				monitoredStopVisits = append(monitoredStopVisits, objectid)
			}
		}

		connector.broadcastNotCollectedEvents(events, monitoredStopVisits)
	}
}

func (connector *SIRIStopMonitoringRequestCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {
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

func (connector *SIRIStopMonitoringRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "GetStopMonitoringRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (factory *SIRIStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRIStopMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringRequestCollector(partner)
}

func (connector *SIRIStopMonitoringRequestCollector) logSIRIStopMonitoringRequest(message *audit.BigQueryMessage, request *siri.SIRIGetStopMonitoringRequest) {
	message.RequestIdentifier = request.MessageIdentifier

	xml, err := request.BuildXML(connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLStopMonitoringResponse(message *audit.BigQueryMessage, response *sxml.XMLStopMonitoringResponse) {
	for _, delivery := range response.StopMonitoringDeliveries() {
		if !delivery.Status() {
			message.Status = "Error"
		}
	}
	message.ResponseIdentifier = response.ResponseMessageIdentifier()
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
