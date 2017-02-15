package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringRequestHandler struct {
	xmlRequest *siri.XMLStopMonitoringRequest
}

func (handler *SIRIStopMonitoringRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIStopMonitoringRequestHandler) ConnectorType() string {
	return "siri-stop-monitoring-request-collector"
}

func (handler *SIRIStopMonitoringRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("StopMonitoring %s\n", handler.xmlRequest.MessageIdentifier())

	tx := connector.(*core.SIRIStopMonitoringRequestCollector).Partner().Referential().NewTransaction()
	defer tx.Close()

	objectidKind := connector.(*core.SIRIStopMonitoringRequestCollector).Partner().Setting("remote_objectid_kind")
	objectid := model.NewObjectID(objectidKind, handler.xmlRequest.MonitoringRef())
	stopArea, ok := tx.Model().StopAreas().FindByObjectId(objectid)
	if !ok {
		siriError("NotFound", "StopArea not found", rw)
		return
	}

	response := new(siri.SIRIStopMonitoringResponse)
	response.Address = connector.(*core.SIRIStopMonitoringRequestCollector).Partner().Setting("address")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = handler.xmlRequest.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.(*core.SIRIStopMonitoringRequestCollector).SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = model.DefaultClock().Now()

	// Fill StopVisits
	for _, stopVisit := range tx.Model().StopVisits().FindByStopAreaId(stopArea.Id()) {
		stopVisitId, ok := stopVisit.ObjectID(objectidKind)
		if !ok {
			siriError("InternalServiceError", "", rw)
			return
		}
		schedules := stopVisit.Schedules()
		monitoredStopVisit := &siri.SIRIMonitoredStopVisit{
			ItemIdentifier: stopVisitId.Value(),
			StopPointRef:   objectid.Value(),
			StopPointName:  stopArea.Name,
			// DatedVehicleJourneyRef: stopVisit
			// LineRef                string
			// PublishedLineName      string
			DepartureStatus:       string(stopVisit.DepartureStatus()),
			ArrivalStatus:         string(stopVisit.ArrivalStatus()),
			Order:                 stopVisit.PassageOrder(),
			AimedArrivalTime:      schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).ArrivalTime(),
			ExpectedArrivalTime:   schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).ArrivalTime(),
			ActualArrivalTime:     schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime(),
			AimedDepartureTime:    schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime(),
			ExpectedDepartureTime: schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).DepartureTime(),
			ActualDepartureTime:   schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime(),
		}
		response.MonitoredStopVisits = append(response.MonitoredStopVisits, monitoredStopVisit)
	}

	xmlResponse := response.BuildXML()

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	_, err := soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
	}
}
