package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringRequestBroadcaster interface {
	RequestStopArea(request *siri.XMLStopMonitoringRequest) *siri.SIRIStopMonitoringResponse
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

func (connector *SIRIStopMonitoringRequestBroadcaster) RemoteObjectIDKind() string {
	if connector.partner.Setting("siri-stop-monitoring-request-broadcaster.remote_objectid_kind") != "" {
		return connector.partner.Setting("siri-stop-monitoring-request-broadcaster.remote_objectid_kind")
	}
	return connector.partner.Setting("remote_objectid_kind")
}

func (connector *SIRIStopMonitoringRequestBroadcaster) getStopMonitoringDelivery(tx *model.Transaction, logStashEvent audit.LogStashEvent, request *siri.XMLStopMonitoringSubRequest) siri.SIRIStopMonitoringDelivery {
	// SMRB
	objectidKind := connector.RemoteObjectIDKind()
	objectid := model.NewObjectID(objectidKind, request.MonitoringRef())
	stopArea, ok := tx.Model().StopAreas().FindByObjectId(objectid)
	if !ok {
		return siri.SIRIStopMonitoringDelivery{
			RequestMessageRef: request.MessageIdentifier(),
			Status:            false,
			ResponseTimestamp: connector.Clock().Now(),
			ErrorType:         "InvalidDataReferencesError",
			ErrorText:         fmt.Sprintf("StopArea not found: '%s'", objectid.Value()),
		}
	}

	if !stopArea.CollectedAlways {
		stopArea.CollectedUntil = connector.Clock().Now().Add(15 * time.Minute)
		logger.Log.Printf("StopArea %s will be collected until %v", stopArea.Id(), stopArea.CollectedUntil)
		stopArea.Save()
	}

	delivery := siri.SIRIStopMonitoringDelivery{
		RequestMessageRef: request.MessageIdentifier(),
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	// Prepare StopVisit Selectors
	selectors := []model.StopVisitSelector{}
	if request.LineRef() != "" {
		lineSelectorObjectid := model.NewObjectID(connector.Partner().Setting("remote_objectid_kind"), request.LineRef())
		selectors = append(selectors, model.StopVisitSelectorByLine(lineSelectorObjectid))
	}
	if request.PreviewInterval() != 0 {
		duration := request.PreviewInterval()
		now := connector.Clock().Now()
		if !request.StartTime().IsZero() {
			now = request.StartTime()
		}
		selectors = append(selectors, model.StopVisitSelectorByTime(now, now.Add(duration)))
	}
	selector := model.CompositeStopVisitSelector(selectors)

	// Prepare Id Array for logstash
	var idArray []string

	// Find Descendants
	stopAreas := tx.Model().StopAreas().FindFamily(stopArea.Id())

	// Fill StopVisits
	for _, stopVisit := range tx.Model().StopVisits().FindFollowingByStopAreaIds(stopAreas) {
		if request.MaximumStopVisits() > 0 && len(idArray) >= request.MaximumStopVisits() {
			break
		}
		if !selector(stopVisit) {
			continue
		}

		idArray = append(idArray, string(stopVisit.Id()))

		var itemIdentifier string
		stopVisitId, ok := stopVisit.ObjectID(objectidKind)
		if ok {
			itemIdentifier = stopVisitId.Value()
		} else {
			defaultObjectID, ok := stopVisit.ObjectID("_default")
			if !ok {
				continue
			}
			itemIdentifier = fmt.Sprintf("RATPDev:Item::%s:LOC", defaultObjectID.HashValue())
		}

		schedules := stopVisit.Schedules

		vehicleJourney := stopVisit.VehicleJourney()
		if vehicleJourney == nil {
			logger.Log.Printf("Ignore StopVisit %s without Vehiclejourney", stopVisit.Id())
			continue
		}
		line := vehicleJourney.Line()
		if line == nil {
			logger.Log.Printf("Ignore StopVisit %s without Line", stopVisit.Id())
			continue
		}

		vehicleJourneyId, ok := vehicleJourney.ObjectID(objectidKind)
		var dataVehicleJourneyRef string
		if ok {
			dataVehicleJourneyRef = vehicleJourneyId.Value()
		} else {
			defaultObjectID, ok := vehicleJourney.ObjectID("_default")
			if !ok {
				continue
			}
			dataVehicleJourneyRef = fmt.Sprintf("RATPDev:VehicleJourney::%s:LOC", defaultObjectID.HashValue())
		}

		modelDate := tx.Model().Date()

		lineObjectId, _ := line.ObjectID(objectidKind)

		stopPointRef, _ := tx.Model().StopAreas().Find(stopVisit.StopAreaId)
		stopPointRefObjectId, _ := stopPointRef.ObjectID(objectidKind)

		monitoredStopVisit := &siri.SIRIMonitoredStopVisit{
			ItemIdentifier: itemIdentifier,
			MonitoringRef:  objectid.Value(),
			StopPointRef:   stopPointRefObjectId.Value(),
			StopPointName:  stopPointRef.Name,

			VehicleJourneyName:     vehicleJourney.Name,
			LineRef:                lineObjectId.Value(),
			DatedVehicleJourneyRef: dataVehicleJourneyRef,
			DataFrameRef:           fmt.Sprintf("RATPDev:DataFrame::%s:LOC", modelDate.String()),
			RecordedAt:             stopVisit.RecordedAt,
			PublishedLineName:      line.Name,
			DepartureStatus:        string(stopVisit.DepartureStatus),
			ArrivalStatus:          string(stopVisit.ArrivalStatus),
			Order:                  stopVisit.PassageOrder,
			VehicleAtStop:          stopVisit.VehicleAtStop,
			Attributes:             make(map[string]map[string]string),
			References:             make(map[string]map[string]model.Reference),
		}

		if stopVisit.ArrivalStatus != model.STOP_VISIT_ARRIVAL_CANCELLED && request.StopVisitTypes() != "departures" {
			monitoredStopVisit.AimedArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).ArrivalTime()
			monitoredStopVisit.ExpectedArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).ArrivalTime()
			monitoredStopVisit.ActualArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime()
		}

		if stopVisit.DepartureStatus != model.STOP_VISIT_DEPARTURE_CANCELLED && request.StopVisitTypes() != "arrivals" {
			monitoredStopVisit.AimedDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime()
			monitoredStopVisit.ExpectedDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).DepartureTime()
			monitoredStopVisit.ActualDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime()
		}

		vehicleJourneyRefCopy := vehicleJourney.References.Copy()
		stopVisitRefCopy := stopVisit.References.Copy()

		connector.resolveVehiculeJourneyReferences(vehicleJourneyRefCopy, tx.Model().StopAreas())

		connector.reformatReferences(vehicleJourney.ToFormat(), vehicleJourneyRefCopy)
		connector.reformatReferences(stopVisit.ToFormat(), stopVisitRefCopy)

		monitoredStopVisit.Attributes["StopVisitAttributes"] = stopVisit.Attributes
		monitoredStopVisit.References["StopVisitReferences"] = stopVisitRefCopy

		monitoredStopVisit.Attributes["VehicleJourneyAttributes"] = vehicleJourney.Attributes
		monitoredStopVisit.References["VehicleJourney"] = vehicleJourneyRefCopy

		delivery.MonitoredStopVisits = append(delivery.MonitoredStopVisits, monitoredStopVisit)
	}

	logStashEvent["StopVisitIds"] = strings.Join(idArray, ", ")

	return delivery
}

func (connector *SIRIStopMonitoringRequestBroadcaster) RequestStopArea(request *siri.XMLStopMonitoringRequest) *siri.SIRIStopMonitoringResponse {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLStopMonitoringRequest(logStashEvent, request)

	response := new(siri.SIRIStopMonitoringResponse)
	response.Address = connector.Partner().Setting("local_url")
	response.ProducerRef = connector.Partner().Setting("remote_credential")
	if response.ProducerRef == "" {
		response.ProducerRef = "Edwig"
	}
	response.ResponseMessageIdentifier = connector.SIRIPartner().NewResponseMessageIdentifier()

	response.SIRIStopMonitoringDelivery = connector.getStopMonitoringDelivery(tx, logStashEvent, &request.XMLStopMonitoringSubRequest)

	logSIRIStopMonitoringResponse(logStashEvent, response)

	return response
}

func (connector *SIRIStopMonitoringRequestBroadcaster) resolveVehiculeJourneyReferences(references model.References, manager model.StopAreas) {
	toResolve := []string{"PlaceRef", "OriginRef", "DestinationRef"}

	for _, ref := range toResolve {
		if references[ref] == (model.Reference{}) {
			continue
		}
		if foundStopArea, ok := manager.Find(model.StopAreaId(references[ref].Id)); ok {
			obj, ok := foundStopArea.ObjectID(connector.RemoteObjectIDKind())
			if ok {
				tmp := references[ref]
				tmp.ObjectId = &obj
				references[ref] = tmp
			}
		} else {
			tmp := references[ref]
			tmp.ObjectId.SetValue(tmp.Getformat(ref, tmp.GetSha1()))
		}
	}
}

func (connector *SIRIStopMonitoringRequestBroadcaster) reformatReferences(toReformat []string, references model.References) {
	for _, ref := range toReformat {
		if references[ref] != (model.Reference{}) {
			tmp := references[ref]
			tmp.ObjectId.SetValue(tmp.Getformat(ref, tmp.GetSha1()))
		}
	}
}

func (connector *SIRIStopMonitoringRequestBroadcaster) reformatStopVisitReferences(references model.References) {
	toReformat := []string{"OperatorRef"}

	for _, ref := range toReformat {
		if references[ref] != (model.Reference{}) {
			tmp := references[ref]
			tmp.ObjectId.SetValue(tmp.Getformat(ref, tmp.GetSha1()))
		}
	}
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
	logStashEvent["Connector"] = "StopMonitoringRequestBroadcaster"
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
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	if !response.Status {
		logStashEvent["errorType"] = response.ErrorType
		logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		logStashEvent["errorText"] = response.ErrorText
	}
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
