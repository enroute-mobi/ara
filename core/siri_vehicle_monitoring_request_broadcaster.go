package core

import (
	"fmt"
	"sort"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type VehicleMonitoringRequestBroadcaster interface {
	RequestVehicles(*sxml.XMLGetVehicleMonitoring, *audit.BigQueryMessage) *siri.SIRIVehicleMonitoringResponse
}

type SIRIVehicleMonitoringRequestBroadcaster struct {
	state.Startable

	connector

	vjRemoteCodeSpaces      []string
	vehicleRemoteCodeSpaces []string
}

type SIRIVehicleMonitoringRequestBroadcasterFactory struct{}

func NewSIRIVehicleMonitoringRequestBroadcaster(partner *Partner) *SIRIVehicleMonitoringRequestBroadcaster {
	connector := &SIRIVehicleMonitoringRequestBroadcaster{}

	connector.partner = partner
	return connector
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) Start() {
	connector.vjRemoteCodeSpaces = connector.partner.VehicleJourneyRemoteCodeSpaceWithFallback(SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER)
	connector.vehicleRemoteCodeSpaces = connector.partner.VehicleRemoteCodeSpaceWithFallback(SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER)
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace(SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER)
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) RequestVehicles(request *sxml.XMLGetVehicleMonitoring, message *audit.BigQueryMessage) (siriResponse *siri.SIRIVehicleMonitoringResponse) {
	lineRef := request.LineRef()
	vehicleRef := request.VehicleRef()

	messageIdentifier := request.MessageIdentifier()
	message.RequestIdentifier = messageIdentifier

	siriResponse = &siri.SIRIVehicleMonitoringResponse{
		ResponseTimestamp:         connector.Clock().Now(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
		RequestMessageRef:         messageIdentifier,
	}

	delivery := &siri.SIRIVehicleMonitoringDelivery{
		Version:            "2.0:FR-IDF-2.4",
		ResponseTimestamp:  connector.Clock().Now(),
		RequestMessageRef:  messageIdentifier,
		LineRefs:           make(map[string]struct{}),
		VehicleJourneyRefs: make(map[string]struct{}),
		VehicleRefs:        make(map[string]struct{}),
	}

	invalidFiltering := (lineRef != "" && vehicleRef != "") || (lineRef == "" && vehicleRef == "")
	if invalidFiltering {
		delivery.ErrorCondition = &siri.ErrorCondition{
			ErrorType: "InvalidDataReferencesError",
			ErrorText: "VehicleMonitoringRequest must have one LineRef OR one VehicleRef",
		}
		message.Status = "Error"
		message.ErrorDetails = delivery.ErrorCondition.ErrorText
		if vehicleRef != "" {
			delivery.VehicleRefs[vehicleRef] = struct{}{}
		}
		if lineRef != "" {
			delivery.LineRefs[lineRef] = struct{}{}
		}
	} else if lineRef != "" {
		connector.getVehiclesWithLineRef(lineRef, delivery, message, siriResponse)
	} else if vehicleRef != "" {
		connector.getVehicle(vehicleRef, delivery, message, siriResponse)
	}

	if connector.partner.PartnerSettings.SortPaylodForTest() {
		sort.Sort(siri.SortByVehicleMonitoringRef{VehicleActivities: delivery.VehicleActivity})
	}

	message.Lines = GetModelReferenceSlice(delivery.LineRefs)
	message.Vehicles = GetModelReferenceSlice(delivery.VehicleRefs)
	message.VehicleJourneys = GetModelReferenceSlice(delivery.VehicleJourneyRefs)

	siriResponse.SIRIVehicleMonitoringDelivery = *delivery

	return siriResponse
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) getVehicle(vehicleRef string, delivery *siri.SIRIVehicleMonitoringDelivery, message *audit.BigQueryMessage, siriResponse *siri.SIRIVehicleMonitoringResponse) {
	code := model.NewCode(connector.remoteCodeSpace, vehicleRef)
	vehicle, ok := connector.partner.Model().Vehicles().FindByCode(code)
	if !ok {
		delivery.ErrorCondition = &siri.ErrorCondition{
			ErrorType: "InvalidDataReferencesError",
			ErrorText: fmt.Sprintf("Vehicle %v not found", code.Value()),
		}
		message.Status = "Error"
		message.ErrorDetails = delivery.ErrorCondition.ErrorText
		delivery.VehicleRefs[vehicleRef] = struct{}{}

		return
	}

	delivery.Status = true

	line, ok := connector.partner.Model().Lines().Find(vehicle.LineId)
	if !ok {
		return
	}
	lineCode, ok := line.Code(code.CodeSpace())
	if !ok {
		return
	}
	connector.buildVehicleActivity(delivery, line, lineCode.Value(), vehicle)
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) getVehiclesWithLineRef(lineRef string, delivery *siri.SIRIVehicleMonitoringDelivery, message *audit.BigQueryMessage, siriResponse *siri.SIRIVehicleMonitoringResponse) {
	code := model.NewCode(connector.remoteCodeSpace, lineRef)
	line, ok := connector.partner.Model().Lines().FindByCode(code)
	if !ok {
		delivery.ErrorCondition = &siri.ErrorCondition{
			ErrorType: "InvalidDataReferencesError",
			ErrorText: fmt.Sprintf("Line %v not found", code.Value()),
		}
		message.Status = "Error"
		message.ErrorDetails = delivery.ErrorCondition.ErrorText
		delivery.LineRefs[lineRef] = struct{}{}

		return
	}

	delivery.Status = true

	vs := connector.partner.Model().Vehicles().FindByLineId(line.Id())
	for i := range vs {
		connector.buildVehicleActivity(delivery, line, lineRef, vs[i])
	}
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) buildVehicleActivity(delivery *siri.SIRIVehicleMonitoringDelivery, line *model.Line, lineRef string, vehicle *model.Vehicle) {
	vehicleId, ok := vehicle.CodeWithFallback(connector.vehicleRemoteCodeSpaces)
	if !ok {
		return
	}

	vj := vehicle.VehicleJourney()
	if vj == nil {
		return
	}
	dvj, ok := connector.datedVehicleJourneyRef(vj)
	if !ok {
		return
	}

	refs := vj.References.Copy()

	activity := &siri.SIRIVehicleActivity{
		RecordedAtTime:       vehicle.RecordedAtTime,
		ValidUntilTime:       vehicle.ValidUntilTime,
		VehicleMonitoringRef: vehicleId.Value(),
		ProgressBetweenStops: connector.handleProgressBetweenStops(vehicle),
	}

	monitoredVehicleJourney := &siri.SIRIMonitoredVehicleJourney{
		LineRef:            lineRef,
		PublishedLineName:  line.Name,
		DirectionName:      vj.Attributes["DirectionName"],
		DirectionType:      vj.DirectionType,
		OriginName:         vj.OriginName,
		DestinationName:    vj.DestinationName,
		Monitored:          vj.Monitored,
		Bearing:            vehicle.Bearing,
		DriverRef:          vehicle.DriverRef,
		Occupancy:          vehicle.Occupancy,
		OriginRef:          connector.handleRef("OriginRef", vj.Origin, refs),
		DestinationRef:     connector.handleRef("DestinationRef", vj.Origin, refs),
		JourneyPatternRef:  connector.handleJourneyPatternRef(refs),
		JourneyPatternName: connector.handleJourneyPatternName(refs),
		VehicleLocation:    connector.handleVehicleLocation(vehicle),
	}

	if vehicle.NextStopVisitId != model.StopVisitId("") {
		nextStopVisit, ok := connector.Partner().Model().StopVisits().Find(vehicle.NextStopVisitId)
		if ok {
			stopArea, stopAreaCode, ok := connector.stopPointRef(nextStopVisit.StopArea().Id())
			if ok {
				monitoredCall := &siri.MonitoredCall{
					StopPointRef:          stopAreaCode,
					StopPointName:         stopArea.Name,
					VehicleAtStop:         nextStopVisit.VehicleAtStop,
					DestinationDisplay:    nextStopVisit.Attributes["DestinationDisplay"],
					ExpectedArrivalTime:   nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Expected}),
					ExpectedDepartureTime: nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Expected}),
					DepartureStatus:       string(nextStopVisit.DepartureStatus),
					Order:                 &nextStopVisit.PassageOrder,
					AimedArrivalTime:      nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Aimed}),
					AimedDepartureTime:    nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Aimed}),
					ArrivalStatus:         string(nextStopVisit.ArrivalStatus),
					ActualArrivalTime:     nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Actual}),
					ActualDepartureTime:   nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Actual}),
				}
				monitoredVehicleJourney.MonitoredCall = monitoredCall
			}
		}
	}

	framedVehicleJourneyRef := &siri.SIRIFramedVehicleJourneyRef{}
	modelDate := connector.partner.Model().Date()
	framedVehicleJourneyRef.DataFrameRef =
		connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "DataFrame", Id: modelDate.String()})
	framedVehicleJourneyRef.DatedVehicleJourneyRef = dvj

	monitoredVehicleJourney.FramedVehicleJourneyRef = framedVehicleJourneyRef
	activity.MonitoredVehicleJourney = monitoredVehicleJourney
	delivery.VehicleActivity = append(delivery.VehicleActivity, activity)

	// Logging
	delivery.LineRefs[lineRef] = struct{}{}
	delivery.VehicleJourneyRefs[dvj] = struct{}{}
	delivery.VehicleRefs[vehicleId.Value()] = struct{}{}
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) stopPointRef(stopAreaId model.StopAreaId) (*model.StopArea, string, bool) {
	stopPointRef, ok := connector.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return &model.StopArea{}, "", false
	}

	stopPointRefCode, ok := stopPointRef.Code(connector.remoteCodeSpace)
	if ok {
		return stopPointRef, stopPointRefCode.Value(), true
	}
	referent, ok := stopPointRef.Referent()
	if ok {
		referentCode, ok := referent.Code(connector.remoteCodeSpace)
		if ok {
			return referent, referentCode.Value(), true
		}
	}
	return &model.StopArea{}, "", false
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) datedVehicleJourneyRef(vehicleJourney *model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.CodeWithFallback(connector.vjRemoteCodeSpaces)

	var dataVehicleJourneyRef string
	if ok {
		dataVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultCode, ok := vehicleJourney.Code(model.Default)
		if !ok {
			return "", false
		}
		dataVehicleJourneyRef =
			connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultCode.Value()})
	}
	return dataVehicleJourneyRef, true
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleVehicleLocation(v *model.Vehicle) *siri.SIRIVehicleLocation {
	var lat = v.Latitude
	var lon = v.Longitude
	if lat != 0. || lon != 0. {
		vehicleLocation := &siri.SIRIVehicleLocation{
			Longitude: lon,
			Latitude:  lat,
		}
		return vehicleLocation
	}
	return nil
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleProgressBetweenStops(v *model.Vehicle) *siri.SIRIProgressBetweenStops {
	var dist = v.LinkDistance
	var percent = v.Percentage
	if dist != 0. || percent != 0. {
		progressBetweenStops := &siri.SIRIProgressBetweenStops{
			LinkDistance: dist,
			Percentage:   percent,
		}
		return progressBetweenStops
	}
	return nil
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleJourneyPatternRef(refs model.References) string {
	journeyPatternRef, ok := refs.Get("JourneyPatternRef")
	if ok {
		if connector.remoteCodeSpace == journeyPatternRef.Code.CodeSpace() {
			return journeyPatternRef.Code.Value()
		}
	}

	return ""
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleJourneyPatternName(refs model.References) string {
	journeyPatternName, ok := refs.Get(siri_attributes.JourneyPatternName)
	if ok {
		if connector.remoteCodeSpace == journeyPatternName.Code.CodeSpace() {
			return journeyPatternName.Code.Value()
		}
	}

	return ""
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) handleRef(refType, origin string, references model.References) string {
	reference, ok := references.Get(refType)
	if !ok || reference.Code == nil || (refType == "DestinationRef" && connector.noDestinationRefRewritingFrom(origin)) {
		return ""
	}
	return connector.resolveStopAreaRef(reference)
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) noDestinationRefRewritingFrom(origin string) bool {
	ndrrf := connector.Partner().NoDestinationRefRewritingFrom()
	for _, o := range ndrrf {
		if origin == strings.TrimSpace(o) {
			return true
		}
	}
	return false
}

func (connector *SIRIVehicleMonitoringRequestBroadcaster) resolveStopAreaRef(reference model.Reference) string {
	stopArea, ok := connector.partner.Model().StopAreas().FindByCode(*reference.Code)
	if ok {
		obj, ok := stopArea.ReferentOrSelfCode(connector.remoteCodeSpace)
		if ok {
			return obj.Value()
		}
	}
	return connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "StopArea", Id: reference.GetSha1()})
}

func (factory *SIRIVehicleMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIVehicleMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIVehicleMonitoringRequestBroadcaster(partner)
}
