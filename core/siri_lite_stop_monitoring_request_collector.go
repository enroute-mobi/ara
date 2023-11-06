package core

import (
	"encoding/json"
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/slite"
)

type LiteStopMonitoringRequestCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
}

type SIRILiteStopMonitoringRequestCollector struct {
	connector

	updateSubscriber UpdateSubscriber
}

type SIRILiteStopMonitoringRequestCollectorFactory struct{}

func NewSIRILiteStopMonitoringRequestCollector(partner *Partner) *SIRILiteStopMonitoringRequestCollector {
	connector := &SIRILiteStopMonitoringRequestCollector{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind()
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent

	return connector
}

func (connector *SIRILiteStopMonitoringRequestCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	stopArea, ok := connector.partner.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("StopAreaUpdateRequest in LiteStopMonitoringRequestCollector for unknown StopArea %v", request.StopAreaId())
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

	dest, err := connector.Partner().SIRILiteClient().StopMonitoring(objectid.Value())
	if err != nil {
		e := fmt.Sprintf("Error during LiteStopMonitoring request: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	message.RequestRawMessage = fmt.Sprintf("MonitoringRef=%s", objectid.Value())
	logSIRILiteStopMonitoringResponse(message, dest)

	builder := NewLiteStopMonitoringUpdateEventBuilder(connector.partner, objectid)
	for _, delivery := range dest.Siri.ServiceDelivery.StopMonitoringDelivery {
		if delivery.Status == "false" {
			continue
		}
		builder.SetUpdateEvents(delivery.MonitoredStopVisit)
	}
	updateEvents := builder.UpdateEvents()

	//  Log Models
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

func (connector *SIRILiteStopMonitoringRequestCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {

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

func (connector *SIRILiteStopMonitoringRequestCollector) broadcastUpdateEvent(event model.UpdateEvent) {
	if connector.updateSubscriber != nil {
		connector.updateSubscriber(event)
	}
}

func (connector *SIRILiteStopMonitoringRequestCollector) broadcastNotCollectedEvents(events map[string]*model.StopVisitUpdateEvent, collectedStopVisitObjectIDs []model.ObjectID) {
	for _, stopVisitObjectID := range collectedStopVisitObjectIDs {
		if _, ok := events[stopVisitObjectID.Value()]; !ok {
			logger.Log.Debugf("Send StopVisitNotCollectedEvent for %v", stopVisitObjectID)
			connector.broadcastUpdateEvent(model.NewNotCollectedUpdateEvent(stopVisitObjectID))
		}
	}
}

func (connector *SIRILiteStopMonitoringRequestCollector) SetUpdateSubscriber(updateSubscriber UpdateSubscriber) {
	connector.updateSubscriber = updateSubscriber
}

func (connector *SIRILiteStopMonitoringRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "GetStopMonitoringRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (factory *SIRILiteStopMonitoringRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRILiteStopMonitoringRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRILiteStopMonitoringRequestCollector(partner)
}

func logSIRILiteStopMonitoringResponse(message *audit.BigQueryMessage, response *slite.SIRILiteStopMonitoring) {
	for _, delivery := range response.Siri.ServiceDelivery.StopMonitoringDelivery {
		if delivery.Status == "false " {
			message.Status = "Error"
		}
	}
	message.ResponseIdentifier = response.Siri.ServiceDelivery.ResponseMessageIdentifier
	b, err := json.Marshal(response)
	if err != nil {
		return
	}
	message.ResponseRawMessage = string(b)
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
