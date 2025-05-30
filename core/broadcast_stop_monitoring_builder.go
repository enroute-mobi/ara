package core

import (
	"strings"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type BroadcastStopMonitoringBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner *Partner

	StopVisitTypes string
	MonitoringRef  string

	remoteCodeSpace               string
	vjRemoteCodeSpaces            []string
	noDestinationRefRewritingFrom []string
	noDataFrameRefRewritingFrom   []string
	rewriteJourneyPatternRef      bool
}

func NewBroadcastStopMonitoringBuilder(partner *Partner, connectorName string) *BroadcastStopMonitoringBuilder {
	return &BroadcastStopMonitoringBuilder{
		partner:                       partner,
		remoteCodeSpace:               partner.RemoteCodeSpace(connectorName),
		vjRemoteCodeSpaces:            partner.VehicleJourneyRemoteCodeSpaceWithFallback(connectorName),
		noDestinationRefRewritingFrom: partner.NoDestinationRefRewritingFrom(),
		noDataFrameRefRewritingFrom:   partner.NoDataFrameRefRewritingFrom(),
		rewriteJourneyPatternRef:      partner.RewriteJourneyPatternRef(),
	}
}

func (builder *BroadcastStopMonitoringBuilder) BuildCancelledStopVisit(stopVisit *model.StopVisit) *siri.SIRICancelledStopVisit {
	vehicleJourney, ok := builder.partner.Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
	if !ok {
		logger.Log.Printf("Ignore CancelledStopVisit %s without Vehiclejourney", stopVisit.Id())
		return nil
	}
	line, ok := builder.partner.Model().Lines().Find(vehicleJourney.LineId)
	if !ok {
		logger.Log.Printf("Ignore CancelledStopVisit %s without Line", stopVisit.Id())
		return nil
	}
	lineCode, ok := line.ReferentOrSelfCode(builder.remoteCodeSpace)
	if !ok {
		logger.Log.Printf("Ignore CancelledStopVisit %s with Line without correct Code", stopVisit.Id())
		return nil
	}

	itemIdentifier, ok := builder.getItemIdentifier(stopVisit)
	if !ok {
		return nil
	}

	datedVehicleJourneyRef, ok := builder.datedVehicleJourneyRef(vehicleJourney)
	if !ok {
		return nil
	}

	cancelledStopVisit := &siri.SIRICancelledStopVisit{
		RecordedAtTime:         stopVisit.RecordedAt,
		ItemRef:                itemIdentifier,
		MonitoringRef:          builder.MonitoringRef,
		VisitNumber:            stopVisit.PassageOrder,
		DirectionRef:           stopVisit.VehicleJourney().DirectionType,
		LineRef:                lineCode.Value(),
		DatedVehicleJourneyRef: datedVehicleJourneyRef,
		DataFrameRef:           builder.dataFrameRef(stopVisit, vehicleJourney.Origin),
		PublishedLineName:      line.Name,
	}

	return cancelledStopVisit
}

func (builder *BroadcastStopMonitoringBuilder) BuildMonitoredStopVisit(stopVisit *model.StopVisit) *siri.SIRIMonitoredStopVisit {
	stopPointRef, stopPointRefCode, ok := builder.stopPointRef(stopVisit.StopAreaId)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %v without StopArea or with StopArea without correct Code", stopVisit.Id())
		return nil
	}

	vehicleJourney, ok := builder.partner.Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %s without Vehiclejourney", stopVisit.Id())
		return nil
	}
	line, lineCode, ok := builder.lineAndCode(vehicleJourney.LineId)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %s without Line or with Line without correct Code", stopVisit.Id())
		return nil
	}

	itemIdentifier, ok := builder.getItemIdentifier(stopVisit)
	if !ok {
		return nil
	}

	datedVehicleJourneyRef, ok := builder.datedVehicleJourneyRef(vehicleJourney)
	if !ok {
		return nil
	}
	var useVisitNumber = builder.useVisitNumber()

	monitoredStopVisit := &siri.SIRIMonitoredStopVisit{
		ItemIdentifier:         itemIdentifier,
		MonitoringRef:          builder.MonitoringRef,
		StopPointRef:           stopPointRefCode,
		StopPointName:          stopPointRef.Name,
		VehicleJourneyName:     vehicleJourney.Name,
		OriginName:             vehicleJourney.OriginName,
		DestinationName:        vehicleJourney.DestinationName,
		DirectionType:          builder.directionType(vehicleJourney.DirectionType),
		Monitored:              vehicleJourney.Monitored,
		LineRef:                lineCode,
		DatedVehicleJourneyRef: datedVehicleJourneyRef,
		DataFrameRef:           builder.dataFrameRef(stopVisit, vehicleJourney.Origin),
		Occupancy:              vehicleJourney.Occupancy,
		RecordedAt:             stopVisit.RecordedAt,
		PublishedLineName:      line.Name,
		DepartureStatus:        string(stopVisit.DepartureStatus),
		ArrivalStatus:          string(stopVisit.ArrivalStatus),
		Order:                  stopVisit.PassageOrder,
		VehicleAtStop:          stopVisit.VehicleAtStop,
		Attributes:             make(map[string]map[string]string),
		References:             make(map[string]map[string]string),
	}

	monitoredStopVisit.UseVisitNumber = useVisitNumber

	if !stopPointRef.Monitored {
		monitoredStopVisit.Monitored = false
	}

	if stopVisit.ArrivalStatus != model.STOP_VISIT_ARRIVAL_CANCELLED && builder.StopVisitTypes != "departures" {
		monitoredStopVisit.AimedArrivalTime = stopVisit.Schedules.Schedule(schedules.Aimed).ArrivalTime()
		monitoredStopVisit.ExpectedArrivalTime = stopVisit.Schedules.Schedule(schedules.Expected).ArrivalTime()
		if monitoredStopVisit.Monitored {
			monitoredStopVisit.ActualArrivalTime = stopVisit.Schedules.Schedule(schedules.Actual).ArrivalTime()
		}
	}

	if stopVisit.DepartureStatus != model.STOP_VISIT_DEPARTURE_CANCELLED && builder.StopVisitTypes != "arrivals" {
		monitoredStopVisit.AimedDepartureTime = stopVisit.Schedules.Schedule(schedules.Aimed).DepartureTime()
		monitoredStopVisit.ExpectedDepartureTime = stopVisit.Schedules.Schedule(schedules.Expected).DepartureTime()
		if monitoredStopVisit.Monitored {
			monitoredStopVisit.ActualDepartureTime = stopVisit.Schedules.Schedule(schedules.Actual).DepartureTime()
		}
	}

	vehicleJourneyRefCopy := vehicleJourney.References.Copy()
	stopVisitRefCopy := stopVisit.References.Copy()

	builder.resolveVJReferences(vehicleJourneyRefCopy, vehicleJourney.Origin)

	builder.resolveOperator(stopVisitRefCopy)

	monitoredStopVisit.Attributes["StopVisitAttributes"] = stopVisit.Attributes
	monitoredStopVisit.References["StopVisitReferences"] = stopVisitRefCopy.GetSiriReferences()

	vehicle, ok := builder.partner.Model().Vehicles().FindByNextStopVisitId(stopVisit.Id())
	if ok {
		monitoredStopVisit.HasVehicleInformation = true
		if vehicle.Latitude != 0 || vehicle.Longitude != 0 {
			vehicleLocation := &siri.SIRIVehicleLocation{}
			vehicleLocation.Latitude = vehicle.Latitude
			vehicleLocation.Longitude = vehicle.Longitude
			monitoredStopVisit.SIRIVehicleLocation = *vehicleLocation
		}
		// override occupancy
		monitoredStopVisit.Occupancy = vehicle.Occupancy
		monitoredStopVisit.Bearing = vehicle.Bearing
	}

	monitoredStopVisit.Attributes["VehicleJourneyAttributes"] = vehicleJourney.Attributes
	monitoredStopVisit.References["VehicleJourney"] = vehicleJourneyRefCopy.GetSiriReferences()

	return monitoredStopVisit
}

func (builder *BroadcastStopMonitoringBuilder) directionType(direction string) (dir string) {
	in, out, err := builder.partner.PartnerSettings.SIRIDirectionType()
	if err {
		return direction
	}

	switch direction {
	case model.VEHICLE_DIRECTION_INBOUND:
		dir = in
	case model.VEHICLE_DIRECTION_OUTBOUND:
		dir = out
	default:
		dir = direction
	}

	return dir
}

func (builder *BroadcastStopMonitoringBuilder) useVisitNumber() bool {
	switch builder.partner.PartnerSettings.SIRIPassageOrder() {
	case "visit_number":
		return true
	default:
		return false
	}
}

func (builder *BroadcastStopMonitoringBuilder) stopPointRef(stopAreaId model.StopAreaId) (*model.StopArea, string, bool) {
	stopPointRef, ok := builder.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return &model.StopArea{}, "", false
	}
	code, ok := stopPointRef.Code(builder.remoteCodeSpace)
	if ok {
		return stopPointRef, code.Value(), true
	}
	referent, ok := stopPointRef.Referent()
	if ok {
		code, ok := referent.Code(builder.remoteCodeSpace)
		if ok {
			return referent, code.Value(), true
		}
	}
	return &model.StopArea{}, "", false
}

func (builder *BroadcastStopMonitoringBuilder) lineAndCode(lineId model.LineId) (*model.Line, string, bool) {
	line, ok := builder.partner.Model().Lines().Find(lineId)
	if !ok {
		return &model.Line{}, "", false
	}
	code, ok := line.Code(builder.remoteCodeSpace)
	if ok {
		return line, code.Value(), true
	}
	referent, ok := line.Referent()
	if ok {
		code, ok := referent.Code(builder.remoteCodeSpace)
		if ok {
			return referent, code.Value(), true
		}
	}
	return &model.Line{}, "", false
}

func (builder *BroadcastStopMonitoringBuilder) getItemIdentifier(stopVisit *model.StopVisit) (string, bool) {
	var itemIdentifier string

	stopVisitId, ok := stopVisit.Code(builder.remoteCodeSpace)
	if ok {
		itemIdentifier = stopVisitId.Value()
	} else {
		defaultCode, ok := stopVisit.Code(model.Default)
		if !ok {
			logger.Log.Printf("Ignore StopVisit %s without default Code", stopVisit.Id())
			return "", false
		}
		itemIdentifier = builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "Item", Id: defaultCode.Value()})
	}
	return itemIdentifier, true
}

func (builder *BroadcastStopMonitoringBuilder) datedVehicleJourneyRef(vehicleJourney *model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.CodeWithFallback(builder.vjRemoteCodeSpaces)

	var datedVehicleJourneyRef string
	if ok {
		datedVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultCode, ok := vehicleJourney.Code(model.Default)
		if !ok {
			return "", false
		}
		datedVehicleJourneyRef = builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultCode.Value()})
	}
	return datedVehicleJourneyRef, true
}

func (builder *BroadcastStopMonitoringBuilder) resolveOperator(references model.References) {
	operatorRef, ok := references.Get(siri_attributes.OperatorRef)
	if !ok {
		return
	}
	operator, ok := builder.partner.Model().Operators().FindByCode(*operatorRef.Code)
	if !ok {
		return
	}
	obj, ok := operator.Code(builder.remoteCodeSpace)
	if !ok {
		return
	}
	ref, _ := references.Get(siri_attributes.OperatorRef)
	ref.Code.SetValue(obj.Value())
}

func (builder *BroadcastStopMonitoringBuilder) resolveVJReferences(references model.References, origin string) {
	if builder.rewriteJourneyPatternRef {
		if reference, ok := references.Get("JourneyPatternRef"); ok {
			reference.Code.SetValue(builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "JourneyPattern", Id: reference.GetSha1()}))
		}
	}

	if reference, ok := references.Get("RouteRef"); ok {
		reference.Code.SetValue(builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "Route", Id: reference.GetSha1()}))
	}

	if reference, ok := references.Get("DatedVehicleJourneyRef"); ok {
		reference.Code.SetValue(builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "DatedVehicleJourney", Id: reference.GetSha1()}, "VehicleJourney"))
	}

	for _, refType := range []string{"PlaceRef", "OriginRef", "DestinationRef"} {
		reference, ok := references.Get(refType)
		if !ok || reference.Code == nil || (refType == "DestinationRef" && builder.noDestinationRefRewrite(origin)) {
			continue
		}
		builder.resolveStopAreaRef(&reference)
	}
}

func (builder *BroadcastStopMonitoringBuilder) resolveStopAreaRef(reference *model.Reference) {
	stopArea, ok := builder.partner.Model().StopAreas().FindByCode(*reference.Code)
	if ok {
		obj, ok := stopArea.ReferentOrSelfCode(builder.remoteCodeSpace)
		if ok {
			reference.Code.SetValue(obj.Value())
			return
		}
	}
	reference.Code.SetValue(builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "StopArea", Id: reference.GetSha1()}))
}

func (builder *BroadcastStopMonitoringBuilder) noDestinationRefRewrite(origin string) bool {
	for _, o := range builder.noDestinationRefRewritingFrom {
		if origin == strings.TrimSpace(o) {
			return true
		}
	}
	return false
}

func (builder *BroadcastStopMonitoringBuilder) noDataFrameRefRewrite(origin string) bool {
	for _, o := range builder.noDataFrameRefRewritingFrom {
		if origin == strings.TrimSpace(o) {
			return true
		}
	}
	return false
}

func (builder *BroadcastStopMonitoringBuilder) dataFrameRef(sv *model.StopVisit, origin string) string {
	if sv.DataFrameRef != "" && builder.noDataFrameRefRewrite(origin) {
		return sv.DataFrameRef
	}
	modelDate := builder.partner.Model().Date()
	return builder.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "DataFrame", Id: modelDate.String()})
}
