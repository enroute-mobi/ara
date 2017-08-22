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
	RequestStopArea(request *siri.XMLGetStopMonitoring) *siri.SIRIStopMonitoringResponse
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

func (connector *SIRIStopMonitoringRequestBroadcaster) getStopMonitoringDelivery(tx *model.Transaction, logStashEvent audit.LogStashEvent, request *siri.XMLStopMonitoringRequest) siri.SIRIStopMonitoringDelivery {
	objectidKind := connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_REQUEST_BROADCASTER)
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
		lineSelectorObjectid := model.NewObjectID(objectidKind, request.LineRef())
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

	// Prepare Id Array
	var stopVisitArray []string

	// Find Descendants
	stopAreas := tx.Model().StopAreas().FindFamily(stopArea.Id())

	// Fill StopVisits
	referenceGenerator := connector.SIRIPartner().IdentifierGenerator("reference_identifier")
	for _, stopVisit := range tx.Model().StopVisits().FindFollowingByStopAreaIds(stopAreas) {
		if request.MaximumStopVisits() > 0 && len(stopVisitArray) >= request.MaximumStopVisits() {
			break
		}
		if !selector(stopVisit) {
			continue
		}

		var itemIdentifier string
		stopVisitId, ok := stopVisit.ObjectID(objectidKind)
		if ok {
			itemIdentifier = stopVisitId.Value()
		} else {
			defaultObjectID, ok := stopVisit.ObjectID("_default")
			if !ok {
				continue
			}
			itemIdentifier = referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "Item", Default: defaultObjectID.Value()})
		}

		stopVisitArray = append(stopVisitArray, itemIdentifier)

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
		if _, ok := line.ObjectID(objectidKind); !ok {
			logger.Log.Printf("Ignore StopVisit %s with Line without correct ObjectID", stopVisit.Id())
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
			dataVehicleJourneyRef = referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "VehicleJourney", Default: defaultObjectID.Value()})
		}

		modelDate := tx.Model().Date()

		lineObjectId, _ := line.ObjectID(objectidKind)

		stopPointRef, _ := tx.Model().StopAreas().Find(stopVisit.StopAreaId)
		stopPointRefObjectId, _ := stopPointRef.ObjectID(objectidKind)

		dataFrameGenerator := connector.SIRIPartner().IdentifierGenerator("data_frame_identifier")
		monitoredStopVisit := &siri.SIRIMonitoredStopVisit{
			ItemIdentifier: itemIdentifier,
			MonitoringRef:  objectid.Value(),
			StopPointRef:   stopPointRefObjectId.Value(),
			StopPointName:  stopPointRef.Name,

			VehicleJourneyName:     vehicleJourney.Name,
			LineRef:                lineObjectId.Value(),
			DatedVehicleJourneyRef: dataVehicleJourneyRef,
			DataFrameRef:           dataFrameGenerator.NewIdentifier(IdentifierAttributes{Id: modelDate.String()}),
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
		connector.resolveOperator(stopVisitRefCopy)

		connector.reformatReferences(vehicleJourney.ToFormat(), vehicleJourneyRefCopy)

		monitoredStopVisit.Attributes["StopVisitAttributes"] = stopVisit.Attributes
		monitoredStopVisit.References["StopVisitReferences"] = stopVisitRefCopy

		monitoredStopVisit.Attributes["VehicleJourneyAttributes"] = vehicleJourney.Attributes
		monitoredStopVisit.References["VehicleJourney"] = vehicleJourneyRefCopy

		delivery.MonitoredStopVisits = append(delivery.MonitoredStopVisits, monitoredStopVisit)
	}

	logStashEvent["StopVisitIds"] = strings.Join(stopVisitArray, ", ")

	return delivery
}

func (connector *SIRIStopMonitoringRequestBroadcaster) RequestStopArea(request *siri.XMLGetStopMonitoring) *siri.SIRIStopMonitoringResponse {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLStopMonitoringRequest(logStashEvent, &request.XMLStopMonitoringRequest)
	logStashEvent["requestorRef"] = request.RequestorRef()

	response := &siri.SIRIStopMonitoringResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
	}

	response.SIRIStopMonitoringDelivery = connector.getStopMonitoringDelivery(tx, logStashEvent, &request.XMLStopMonitoringRequest)

	logSIRIStopMonitoringDelivery(logStashEvent, response.SIRIStopMonitoringDelivery)
	logSIRIStopMonitoringResponse(logStashEvent, response)

	return response
}

func (connector *SIRIStopMonitoringRequestBroadcaster) resolveOperator(references model.References) {
	operatorRef, _ := references["OperatorRef"]
	operator, ok := connector.Partner().Model().Operators().Find(model.OperatorId(operatorRef.Id))
	if !ok {
		return
	}

	obj, ok := operator.ObjectID(connector.Partner().Setting("remote_objectid_kind"))
	if !ok {
		return
	}
	references["OperatorRef"].ObjectId.SetValue(obj.Value())
}

func (connector *SIRIStopMonitoringRequestBroadcaster) resolveVehiculeJourneyReferences(references model.References, manager model.StopAreas) {
	toResolve := []string{"PlaceRef", "OriginRef", "DestinationRef"}

	for _, ref := range toResolve {
		if references[ref] == (model.Reference{}) {
			continue
		}
		if foundStopArea, ok := manager.Find(model.StopAreaId(references[ref].Id)); ok {
			obj, ok := foundStopArea.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_REQUEST_BROADCASTER))
			if ok {
				tmp := references[ref]
				tmp.ObjectId = &obj
				references[ref] = tmp
				continue
			}
		}
		generator := connector.SIRIPartner().IdentifierGenerator("reference_stop_area_identifier")
		tmp := references[ref]
		tmp.ObjectId.SetValue(generator.NewIdentifier(IdentifierAttributes{Default: tmp.GetSha1()}))
	}
}

func (connector *SIRIStopMonitoringRequestBroadcaster) reformatReferences(toReformat []string, references model.References) {
	generator := connector.SIRIPartner().IdentifierGenerator("reference_identifier")
	for _, ref := range toReformat {
		if references[ref] != (model.Reference{}) {
			tmp := references[ref]
			tmp.ObjectId.SetValue(generator.NewIdentifier(IdentifierAttributes{Type: ref[:len(ref)-3], Default: tmp.GetSha1()}))
		}
	}
}

func (connector *SIRIStopMonitoringRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringRequestBroadcaster"
	return event
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
	logStashEvent["stopVisitTypes"] = request.StopVisitTypes()
	logStashEvent["lineRef"] = request.LineRef()
	logStashEvent["maximumStopVisits"] = strconv.Itoa(request.MaximumStopVisits())
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["startTime"] = request.StartTime().String()
	logStashEvent["previewInterval"] = request.PreviewInterval().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIStopMonitoringDelivery(logStashEvent audit.LogStashEvent, delivery siri.SIRIStopMonitoringDelivery) {
	logStashEvent["requestMessageRef"] = delivery.RequestMessageRef
	logStashEvent["responseTimestamp"] = delivery.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(delivery.Status)
	if !delivery.Status {
		logStashEvent["errorType"] = delivery.ErrorType
		if delivery.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber)
		}
		logStashEvent["errorText"] = delivery.ErrorText
	}
}

func logSIRIStopMonitoringResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIStopMonitoringResponse) {
	logStashEvent["address"] = response.Address
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
