package core

import (
	"fmt"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringRequestBroadcaster interface {
	RequestStopArea(request *siri.XMLStopMonitoringRequest) (*siri.SIRIStopMonitoringResponse, error)
}

type SIRIStopMonitoringRequestBroadcaster struct {
	model.ClockConsumer

	siriConnector
}

type SIRIStopMonitoringRequestBroadcasterFactory struct{}

func NewSIRIStopMonitoringRequestBroadcaster(partner *Partner) *SIRIStopMonitoringRequestBroadcaster {
	siriStopMonitoringRequestBroadcaster := &SIRIStopMonitoringRequestBroadcaster{}
	siriStopMonitoringRequestBroadcaster.partner = partner
	return siriStopMonitoringRequestBroadcaster
}

func (connector *SIRIStopMonitoringRequestBroadcaster) RequestStopArea(request *siri.XMLStopMonitoringRequest) (*siri.SIRIStopMonitoringResponse, error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	objectidKind := connector.Partner().Setting("remote_objectid_kind")
	objectid := model.NewObjectID(objectidKind, request.MonitoringRef())
	stopArea, ok := tx.Model().StopAreas().FindByObjectId(objectid)
	if !ok {
		return nil, fmt.Errorf("StopArea not found")
	}

	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLStopMonitoringRequest(logStashEvent, request)

	response := new(siri.SIRIStopMonitoringResponse)
	response.Address = connector.Partner().Setting("address")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = request.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = connector.Clock().Now()

	// Fill StopVisits
	for _, stopVisit := range tx.Model().StopVisits().FindByStopAreaId(stopArea.Id()) {
		stopVisitId, ok := stopVisit.ObjectID(objectidKind)
		if !ok {
			continue
		}
		schedules := stopVisit.Schedules
		monitoredStopVisit := &siri.SIRIMonitoredStopVisit{
			ItemIdentifier: stopVisitId.Value(),
			StopPointRef:   objectid.Value(),
			StopPointName:  stopArea.Name,
			// DatedVehicleJourneyRef: stopVisit
			// LineRef                string
			// PublishedLineName      string
			DepartureStatus:       string(stopVisit.DepartureStatus),
			ArrivalStatus:         string(stopVisit.ArrivalStatus),
			Order:                 stopVisit.PassageOrder,
			AimedArrivalTime:      schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).ArrivalTime(),
			ExpectedArrivalTime:   schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).ArrivalTime(),
			ActualArrivalTime:     schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime(),
			AimedDepartureTime:    schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime(),
			ExpectedDepartureTime: schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).DepartureTime(),
			ActualDepartureTime:   schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime(),
			Attributes:            make(map[string]map[string]string),
		}
		fmt.Printf("StopVisit attributes: %#v\n", stopVisit.Attributes)
		monitoredStopVisit.Attributes["StopVisitAttributes"] = stopVisit.Attributes
		fmt.Printf("VehicleJourney attributes: %#v\n", stopVisit.VehicleJourney().Attributes)
		monitoredStopVisit.Attributes["VehicleJourneyAttributes"] = stopVisit.VehicleJourney().Attributes

		response.MonitoredStopVisits = append(response.MonitoredStopVisits, monitoredStopVisit)
	}

	logSIRIStopMonitoringResponse(logStashEvent, response)

	return response, nil
}

func (factory *SIRIStopMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRIStopMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringRequestBroadcaster(partner)
}

func logXMLStopMonitoringRequest(logStashEvent audit.LogStashEvent, request *siri.XMLStopMonitoringRequest) {
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["monitoringRef"] = request.MonitoringRef()
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIStopMonitoringResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIStopMonitoringResponse) {
	logStashEvent["address"] = response.Address
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
