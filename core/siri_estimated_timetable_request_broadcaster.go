package core

import (
	"fmt"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type EstimatedTimetableBroadcaster interface {
	RequestLine(*sxml.XMLGetEstimatedTimetable, *audit.BigQueryMessage) *siri.SIRIEstimatedTimeTableResponse
}

type SIRIEstimatedTimetableBroadcaster struct {
	connector

	dataFrameGenerator    *idgen.IdentifierGenerator
	vjRemoteObjectidKinds []string
}

type SIRIEstimatedTimetableBroadcasterFactory struct{}

func NewSIRIEstimatedTimetableBroadcaster(partner *Partner) *SIRIEstimatedTimetableBroadcaster {
	connector := &SIRIEstimatedTimetableBroadcaster{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER)
	connector.dataFrameGenerator = partner.IdentifierGenerator(idgen.DATA_FRAME_IDENTIFIER)
	connector.vjRemoteObjectidKinds = partner.VehicleJourneyRemoteObjectIDKindWithFallback(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER)
	connector.partner = partner
	return connector
}

func (connector *SIRIEstimatedTimetableBroadcaster) RequestLine(request *sxml.XMLGetEstimatedTimetable, message *audit.BigQueryMessage) *siri.SIRIEstimatedTimeTableResponse {
	response := &siri.SIRIEstimatedTimeTableResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
	}

	response.SIRIEstimatedTimetableDelivery = connector.getEstimatedTimetableDelivery(&request.XMLEstimatedTimetableRequest)

	if !response.SIRIEstimatedTimetableDelivery.Status {
		message.Status = "Error"
		message.ErrorDetails = response.SIRIEstimatedTimetableDelivery.ErrorString()
	}
	message.Lines = request.Lines()
	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	return response
}

func (connector *SIRIEstimatedTimetableBroadcaster) getEstimatedTimetableDelivery(request *sxml.XMLEstimatedTimetableRequest) siri.SIRIEstimatedTimetableDelivery {
	currentTime := connector.Clock().Now()

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
		lineObjectId := model.NewObjectID(connector.remoteObjectidKind, lineId)
		line, ok := connector.partner.Model().Lines().FindByObjectId(lineObjectId)
		if !ok {
			logger.Log.Debugf("Cannot find requested line Estimated Time Table with id %v at %v", lineObjectId.String(), connector.Clock().Now())
			continue
		}

		journeyFrame := &siri.SIRIEstimatedJourneyVersionFrame{
			RecordedAtTime: currentTime,
		}

		// SIRIEstimatedVehicleJourney
		vjs := connector.partner.Model().VehicleJourneys().FindByLineId(line.Id())
		for i := range vjs {
			// Handle vehicleJourney Objectid
			vehicleJourneyId, ok := vjs[i].ObjectIDWithFallback(connector.vjRemoteObjectidKinds)
			var datedVehicleJourneyRef string
			if ok {
				datedVehicleJourneyRef = vehicleJourneyId.Value()
			} else {
				defaultObjectID, ok := vjs[i].ObjectID("_default")
				if !ok {
					logger.Log.Debugf("Vehicle journey with id %v does not have a proper objectid at %v", vehicleJourneyId, connector.Clock().Now())
					continue
				}
				referenceGenerator := connector.Partner().IdentifierGenerator(idgen.REFERENCE_IDENTIFIER)
				datedVehicleJourneyRef = referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultObjectID.Value()})
			}

			estimatedVehicleJourney := &siri.SIRIEstimatedVehicleJourney{
				LineRef:                lineObjectId.Value(),
				DirectionType:          vjs[i].DirectionType,
				DatedVehicleJourneyRef: datedVehicleJourneyRef,
				DataFrameRef:           connector.dataFrameRef(),
				PublishedLineName:      connector.publishedLineName(line),
				Attributes:             make(map[string]string),
				References:             make(map[string]string),
			}
			estimatedVehicleJourney.References = connector.getEstimatedVehicleJourneyReferences(vjs[i], vjs[i].Origin)
			estimatedVehicleJourney.Attributes = vjs[i].Attributes

			// SIRIEstimatedCall
			svs := connector.partner.Model().StopVisits().FindFollowingByVehicleJourneyId(vjs[i].Id())
			for i := range svs {
				if !selector(svs[i]) {
					continue
				}

				// Handle StopPointRef
				stopArea, stopAreaId, ok := connector.stopPointRef(svs[i].StopAreaId)
				if !ok {
					logger.Log.Printf("Ignore Stopvisit %v without StopArea or with StopArea without correct ObjectID", svs[i].Id())
					continue
				}

				connector.resolveOperatorRef(estimatedVehicleJourney.References, svs[i])

				estimatedCall := &siri.SIRIEstimatedCall{
					ArrivalStatus:      string(svs[i].ArrivalStatus),
					DepartureStatus:    string(svs[i].DepartureStatus),
					AimedArrivalTime:   svs[i].Schedules.Schedule("aimed").ArrivalTime(),
					AimedDepartureTime: svs[i].Schedules.Schedule("aimed").DepartureTime(),
					Order:              svs[i].PassageOrder,
					StopPointRef:       stopAreaId,
					StopPointName:      stopArea.Name,
					DestinationDisplay: svs[i].Attributes["DestinationDisplay"],
					VehicleAtStop:      svs[i].VehicleAtStop,
				}

				estimatedCall.UseVisitNumber = connector.UseVisitNumber()

				if stopArea.Monitored {
					estimatedCall.ExpectedArrivalTime = svs[i].Schedules.Schedule("expected").ArrivalTime()
					estimatedCall.ExpectedDepartureTime = svs[i].Schedules.Schedule("expected").DepartureTime()
				} else if connector.Partner().SendProducerUnavailableError() {
					delivery.Status = false
					delivery.ErrorType = "OtherError"
					delivery.ErrorNumber = 1
					delivery.ErrorText = fmt.Sprintf("Erreur [PRODUCER_UNAVAILABLE] : %v indisponible", strings.Join(stopArea.Origins.PartnersKO(), ", "))
				}

				estimatedVehicleJourney.EstimatedCalls = append(estimatedVehicleJourney.EstimatedCalls, estimatedCall)
			}
			if len(estimatedVehicleJourney.EstimatedCalls) != 0 {
				journeyFrame.EstimatedVehicleJourneys = append(journeyFrame.EstimatedVehicleJourneys, estimatedVehicleJourney)
			}
		}
		if len(journeyFrame.EstimatedVehicleJourneys) != 0 {
			delivery.EstimatedJourneyVersionFrames = append(delivery.EstimatedJourneyVersionFrames, journeyFrame)
		}
	}
	return delivery
}

func (connector *SIRIEstimatedTimetableBroadcaster) UseVisitNumber() bool {
	switch connector.partner.PartnerSettings.SIRIPassageOrder() {
	case "visit_number":
		return true
	default:
		return false
	}
}

func (connector *SIRIEstimatedTimetableBroadcaster) stopPointRef(stopAreaId model.StopAreaId) (*model.StopArea, string, bool) {
	stopPointRef, ok := connector.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return &model.StopArea{}, "", false
	}
	stopPointRefObjectId, ok := stopPointRef.ObjectID(connector.remoteObjectidKind)
	if ok {
		return stopPointRef, stopPointRefObjectId.Value(), true
	}
	referent, ok := stopPointRef.Referent()
	if ok {
		referentObjectId, ok := referent.ObjectID(connector.remoteObjectidKind)
		if ok {
			return referent, referentObjectId.Value(), true
		}
	}
	return &model.StopArea{}, "", false
}

func (connector *SIRIEstimatedTimetableBroadcaster) publishedLineName(line *model.Line) string {
	var pln string

	switch connector.partner.PartnerSettings.SIRILinePublishedName() {
	case "number":
		if line.Number != "" {
			pln = line.Number
		} else {
			pln = line.Name
		}
	default:
		pln = line.Name
	}

	return pln
}

func (connector *SIRIEstimatedTimetableBroadcaster) getEstimatedVehicleJourneyReferences(vehicleJourney *model.VehicleJourney, origin string) map[string]string {
	references := make(map[string]string)

	for _, refType := range []string{"OriginRef", "DestinationRef"} {
		ref, ok := vehicleJourney.Reference(refType)
		if !ok || ref == (model.Reference{}) || ref.ObjectId == nil {
			continue
		}
		if refType == "DestinationRef" && connector.noDestinationRefRewrite(origin) {
			references[refType] = ref.ObjectId.Value()
			continue
		}
		if foundStopArea, ok := connector.partner.Model().StopAreas().FindByObjectId(*ref.ObjectId); ok {
			obj, ok := foundStopArea.ReferentOrSelfObjectId(connector.remoteObjectidKind)
			if ok {
				references[refType] = obj.Value()
				continue
			}
		}
		generator := connector.Partner().IdentifierGenerator(idgen.REFERENCE_STOP_AREA_IDENTIFIER)
		defaultObjectID := model.NewObjectID(connector.remoteObjectidKind, generator.NewIdentifier(idgen.IdentifierAttributes{Id: ref.GetSha1()}))
		references[refType] = defaultObjectID.Value()
	}

	return references
}

func (connector *SIRIEstimatedTimetableBroadcaster) noDestinationRefRewrite(origin string) bool {
	for _, o := range connector.Partner().NoDestinationRefRewritingFrom() {
		if origin == strings.TrimSpace(o) {
			return true
		}
	}
	return false
}

func (connector *SIRIEstimatedTimetableBroadcaster) dataFrameRef() string {
	modelDate := connector.partner.Model().Date()
	return connector.dataFrameGenerator.NewIdentifier(idgen.IdentifierAttributes{Id: modelDate.String()})
}

func (connector *SIRIEstimatedTimetableBroadcaster) resolveOperatorRef(refs map[string]string, stopVisit *model.StopVisit) {
	if _, ok := refs["OperatorRef"]; ok {
		return
	}

	operatorRef, ok := stopVisit.Reference("OperatorRef")
	if !ok || operatorRef == (model.Reference{}) || operatorRef.ObjectId == nil {
		return
	}
	operator, ok := connector.partner.Model().Operators().FindByObjectId(*operatorRef.ObjectId)
	if !ok {
		refs["OperatorRef"] = operatorRef.ObjectId.Value()
		return
	}
	obj, ok := operator.ObjectID(connector.remoteObjectidKind)
	if !ok {
		refs["OperatorRef"] = operatorRef.ObjectId.Value()
		return
	}
	refs["OperatorRef"] = obj.Value()
}

func (factory *SIRIEstimatedTimetableBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIEstimatedTimetableBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIEstimatedTimetableBroadcaster(partner)
}
