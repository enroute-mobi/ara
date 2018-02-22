package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type EstimatedTimetableBroadcaster interface {
	RequestLine(request *siri.XMLGetEstimatedTimetable) *siri.SIRIEstimatedTimeTableResponse
}

type SIRIEstimatedTimetableBroadcaster struct {
	model.ClockConsumer

	siriConnector
}

type SIRIEstimatedTimetableBroadcasterFactory struct{}

func NewSIRIEstimatedTimetableBroadcaster(partner *Partner) *SIRIEstimatedTimetableBroadcaster {
	broadcaster := &SIRIEstimatedTimetableBroadcaster{}
	broadcaster.partner = partner
	return broadcaster
}

func (connector *SIRIEstimatedTimetableBroadcaster) RequestLine(request *siri.XMLGetEstimatedTimetable) *siri.SIRIEstimatedTimeTableResponse {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLEstimatedTimetableRequest(logStashEvent, &request.XMLEstimatedTimetableRequest)
	logStashEvent["requestorRef"] = request.RequestorRef()

	response := &siri.SIRIEstimatedTimeTableResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
	}

	response.SIRIEstimatedTimetableDelivery = connector.getEstimatedTimetableDelivery(tx, &request.XMLEstimatedTimetableRequest, logStashEvent)

	logSIRIEstimatedTimetableResponse(logStashEvent, response)

	return response
}

func (connector *SIRIEstimatedTimetableBroadcaster) getEstimatedTimetableDelivery(tx *model.Transaction, request *siri.XMLEstimatedTimetableRequest, logStashEvent audit.LogStashEvent) siri.SIRIEstimatedTimetableDelivery {
	currentTime := connector.Clock().Now()
	monitoringRefs := []string{}
	lineRefs := []string{}

	delivery := siri.SIRIEstimatedTimetableDelivery{
		RequestMessageRef: request.MessageIdentifier(),
		ResponseTimestamp: currentTime,
		Status:            true,
	}

	selectors := []model.StopVisitSelector{}

	if request.PreviewInterval() != 0 {
		duration := request.PreviewInterval()
		now := connector.Clock().Now()
		if !request.StartTime().IsZero() {
			now = request.StartTime()
		}
		selectors = append(selectors, model.StopVisitSelectorByTime(now, now.Add(duration)))
	}
	selector := model.CompositeStopVisitSelector(selectors)

	// SIRIEstimatedJourneyVersionFrame
	for _, lineId := range request.Lines() {
		lineObjectId := model.NewObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER), lineId)
		line, ok := tx.Model().Lines().FindByObjectId(lineObjectId)
		if !ok {
			logger.Log.Debugf("Cannot find requested line Estimated Time Table with id %v at %v", lineObjectId.String(), connector.Clock().Now())
			continue
		}

		journeyFrame := &siri.SIRIEstimatedJourneyVersionFrame{
			RecordedAtTime: currentTime,
		}

		// SIRIEstimatedVehicleJourney
		for _, vehicleJourney := range tx.Model().VehicleJourneys().FindByLineId(line.Id()) {
			// Handle vehicleJourney Objectid
			vehicleJourneyId, ok := vehicleJourney.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER))
			var datedVehicleJourneyRef string
			if ok {
				datedVehicleJourneyRef = vehicleJourneyId.Value()
			} else {
				defaultObjectID, ok := vehicleJourney.ObjectID("_default")
				if !ok {
					logger.Log.Debugf("Vehicle journey with id %v does not have a proper objectid at %v", vehicleJourneyId, connector.Clock().Now())
					continue
				}
				referenceGenerator := connector.SIRIPartner().IdentifierGenerator("reference_identifier")
				datedVehicleJourneyRef = referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "VehicleJourney", Default: defaultObjectID.Value()})
			}

			estimatedVehicleJourney := &siri.SIRIEstimatedVehicleJourney{
				LineRef:                lineObjectId.Value(),
				DatedVehicleJourneyRef: datedVehicleJourneyRef,
				Attributes:             make(map[string]string),
				References:             make(map[string]model.Reference),
			}
			lineRefs = append(lineRefs, estimatedVehicleJourney.LineRef)
			estimatedVehicleJourney.References = connector.getEstimatedVehicleJourneyReferences(vehicleJourney, tx)
			estimatedVehicleJourney.Attributes = vehicleJourney.Attributes

			// SIRIEstimatedCall
			for _, stopVisit := range tx.Model().StopVisits().FindFollowingByVehicleJourneyId(vehicleJourney.Id()) {
				if !selector(stopVisit) {
					continue
				}

				// Handle StopPointRef
				stopArea, ok := tx.Model().StopAreas().Find(stopVisit.StopAreaId)
				if !ok {
					logger.Log.Debugf("Cant find stopArea with id %v for stopVisit %v at %v", stopVisit.StopAreaId, stopVisit.Id(), connector.Clock().Now())
					continue
				}

				stopAreaId, ok := stopArea.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER))
				if !ok {
					logger.Log.Debugf("Stop Area with id %v does not have a proper objectid at %v", stopArea.Id(), connector.Clock().Now())
					continue
				}

				monitoringRefs = append(monitoringRefs, stopAreaId.Value())
				estimatedCall := &siri.SIRIEstimatedCall{
					ArrivalStatus:         string(stopVisit.ArrivalStatus),
					DepartureStatus:       string(stopVisit.DepartureStatus),
					AimedArrivalTime:      stopVisit.Schedules.Schedule("aimed").ArrivalTime(),
					ExpectedArrivalTime:   stopVisit.Schedules.Schedule("expected").ArrivalTime(),
					AimedDepartureTime:    stopVisit.Schedules.Schedule("aimed").DepartureTime(),
					ExpectedDepartureTime: stopVisit.Schedules.Schedule("expected").DepartureTime(),
					Order:              stopVisit.PassageOrder,
					StopPointRef:       stopAreaId.Value(),
					StopPointName:      stopArea.Name,
					DestinationDisplay: stopVisit.Attributes["DestinationDisplay"],
					VehicleAtStop:      stopVisit.VehicleAtStop,
				}

				estimatedVehicleJourney.EstimatedCalls = append(estimatedVehicleJourney.EstimatedCalls, estimatedCall)
			}

			journeyFrame.EstimatedVehicleJourneys = append(journeyFrame.EstimatedVehicleJourneys, estimatedVehicleJourney)
		}

		delivery.EstimatedJourneyVersionFrames = append(delivery.EstimatedJourneyVersionFrames, journeyFrame)
	}

	logSIRIEstimatedTimetableDelivery(logStashEvent, delivery, monitoringRefs, lineRefs)

	return delivery
}

func (connector *SIRIEstimatedTimetableBroadcaster) getEstimatedVehicleJourneyReferences(vehicleJourney model.VehicleJourney, tx *model.Transaction) map[string]model.Reference {
	references := make(map[string]model.Reference)

	for _, refType := range []string{"OriginRef", "DestinationRef"} {
		ref, ok := vehicleJourney.Reference(refType)
		if !ok || ref == (model.Reference{}) || ref.ObjectId == nil {
			continue
		}
		if foundStopArea, ok := tx.Model().StopAreas().FindByObjectId(*ref.ObjectId); ok {
			obj, ok := foundStopArea.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER))
			if ok {
				references[refType] = *model.NewReference(obj)
				continue
			}
		}
		generator := connector.SIRIPartner().IdentifierGenerator("reference_stop_area_identifier")
		defaultObjectID := model.NewObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER), generator.NewIdentifier(IdentifierAttributes{Default: ref.GetSha1()}))
		references[refType] = *model.NewReference(defaultObjectID)
	}
	return references
}

func (connector *SIRIEstimatedTimetableBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "EstimatedTimetableRequestBroadcaster"
	return event
}

func (factory *SIRIEstimatedTimetableBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRIEstimatedTimetableBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIEstimatedTimetableBroadcaster(partner)
}

func logXMLEstimatedTimetableRequest(logStashEvent audit.LogStashEvent, request *siri.XMLEstimatedTimetableRequest) {
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestedLines"] = strings.Join(request.Lines(), ",")
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIEstimatedTimetableDelivery(logStashEvent audit.LogStashEvent, delivery siri.SIRIEstimatedTimetableDelivery, monitoringRefs, lineRefs []string) {

	logStashEvent["requestMessageRef"] = strings.Join(monitoringRefs, ",")
	logStashEvent["requestMessageRef"] = delivery.RequestMessageRef
	logStashEvent["responseTimestamp"] = delivery.ResponseTimestamp.String()
	logStashEvent["lineRef"] = strings.Join(lineRefs, ",")
	logStashEvent["monitoringRef"] = strings.Join(monitoringRefs, ",")
	logStashEvent["status"] = strconv.FormatBool(delivery.Status)
	if !delivery.Status {
		logStashEvent["errorType"] = delivery.ErrorType
		if delivery.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber)
		}
		logStashEvent["errorText"] = delivery.ErrorText
	}
}

func logSIRIEstimatedTimetableResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIEstimatedTimeTableResponse) {
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
