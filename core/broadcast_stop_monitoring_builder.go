package core

import (
	"strings"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type BroadcastStopMonitoringBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner *Partner

	StopVisitTypes string
	MonitoringRef  string

	referenceGenerator            *idgen.IdentifierGenerator
	stopAreareferenceGenerator    *idgen.IdentifierGenerator
	dataFrameGenerator            *idgen.IdentifierGenerator
	remoteObjectidKind            string
	vjRemoteObjectidKinds         []string
	noDestinationRefRewritingFrom []string
	noDataFrameRefRewritingFrom   []string
	rewriteJourneyPatternRef      bool
}

func NewBroadcastStopMonitoringBuilder(partner *Partner, connectorName string) *BroadcastStopMonitoringBuilder {
	return &BroadcastStopMonitoringBuilder{
		partner:                       partner,
		referenceGenerator:            partner.IdentifierGenerator(idgen.REFERENCE_IDENTIFIER),
		stopAreareferenceGenerator:    partner.IdentifierGenerator(idgen.REFERENCE_STOP_AREA_IDENTIFIER),
		dataFrameGenerator:            partner.IdentifierGenerator(idgen.DATA_FRAME_IDENTIFIER),
		remoteObjectidKind:            partner.RemoteObjectIDKind(connectorName),
		vjRemoteObjectidKinds:         partner.VehicleJourneyRemoteObjectIDKindWithFallback(connectorName),
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
	lineObjectId, ok := line.ObjectID(builder.remoteObjectidKind)
	if !ok {
		logger.Log.Printf("Ignore CancelledStopVisit %s with Line without correct ObjectID", stopVisit.Id())
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
		LineRef:                lineObjectId.Value(),
		DatedVehicleJourneyRef: datedVehicleJourneyRef,
		DataFrameRef:           builder.dataFrameRef(stopVisit, vehicleJourney.Origin),
		PublishedLineName:      line.Name,
	}

	return cancelledStopVisit
}

func (builder *BroadcastStopMonitoringBuilder) BuildMonitoredStopVisit(stopVisit *model.StopVisit) *siri.SIRIMonitoredStopVisit {
	stopPointRef, stopPointRefObjectId, ok := builder.stopPointRef(stopVisit.StopAreaId)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %v without StopArea or with StopArea without correct ObjectID", stopVisit.Id())
		return nil
	}

	vehicleJourney, ok := builder.partner.Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %s without Vehiclejourney", stopVisit.Id())
		return nil
	}
	line, ok := builder.partner.Model().Lines().Find(vehicleJourney.LineId)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %s without Line", stopVisit.Id())
		return nil
	}
	lineObjectId, ok := line.ObjectID(builder.remoteObjectidKind)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %s with Line without correct ObjectID", stopVisit.Id())
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
		StopPointRef:           stopPointRefObjectId,
		StopPointName:          stopPointRef.Name,
		VehicleJourneyName:     vehicleJourney.Name,
		OriginName:             vehicleJourney.OriginName,
		DestinationName:        vehicleJourney.DestinationName,
		DirectionType:          builder.directionType(vehicleJourney.DirectionType),
		Monitored:              vehicleJourney.Monitored,
		LineRef:                lineObjectId.Value(),
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
		monitoredStopVisit.AimedArrivalTime = stopVisit.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).ArrivalTime()
		monitoredStopVisit.ExpectedArrivalTime = stopVisit.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).ArrivalTime()
		if monitoredStopVisit.Monitored {
			monitoredStopVisit.ActualArrivalTime = stopVisit.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime()
		}
	}

	if stopVisit.DepartureStatus != model.STOP_VISIT_DEPARTURE_CANCELLED && builder.StopVisitTypes != "arrivals" {
		monitoredStopVisit.AimedDepartureTime = stopVisit.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime()
		monitoredStopVisit.ExpectedDepartureTime = stopVisit.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).DepartureTime()
		if monitoredStopVisit.Monitored {
			monitoredStopVisit.ActualDepartureTime = stopVisit.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime()
		}
	}

	vehicleJourneyRefCopy := vehicleJourney.References.Copy()
	stopVisitRefCopy := stopVisit.References.Copy()

	builder.resolveVJReferences(vehicleJourneyRefCopy, vehicleJourney.Origin)

	builder.resolveOperator(stopVisitRefCopy)

	monitoredStopVisit.Attributes["StopVisitAttributes"] = stopVisit.Attributes
	monitoredStopVisit.References["StopVisitReferences"] = stopVisitRefCopy.GetSiriReferences()

	monitoredStopVisit.Attributes["VehicleJourneyAttributes"] = vehicleJourney.Attributes
	monitoredStopVisit.References["VehicleJourney"] = vehicleJourneyRefCopy.GetSiriReferences()

	return monitoredStopVisit
}

func (builder *BroadcastStopMonitoringBuilder) directionType(direction string) string {
	var dir string

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
	stopPointRefObjectId, ok := stopPointRef.ObjectID(builder.remoteObjectidKind)
	if ok {
		return stopPointRef, stopPointRefObjectId.Value(), true
	}
	referent, ok := stopPointRef.Referent()
	if ok {
		referentObjectId, ok := referent.ObjectID(builder.remoteObjectidKind)
		if ok {
			return referent, referentObjectId.Value(), true
		}
	}
	return &model.StopArea{}, "", false
}

func (builder *BroadcastStopMonitoringBuilder) getItemIdentifier(stopVisit *model.StopVisit) (string, bool) {
	var itemIdentifier string

	stopVisitId, ok := stopVisit.ObjectID(builder.remoteObjectidKind)
	if ok {
		itemIdentifier = stopVisitId.Value()
	} else {
		defaultObjectID, ok := stopVisit.ObjectID("_default")
		if !ok {
			logger.Log.Printf("Ignore StopVisit %s without default ObjectID", stopVisit.Id())
			return "", false
		}
		itemIdentifier = builder.referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: "Item", Id: defaultObjectID.Value()})
	}
	return itemIdentifier, true
}

func (builder *BroadcastStopMonitoringBuilder) datedVehicleJourneyRef(vehicleJourney *model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.ObjectIDWithFallback(builder.vjRemoteObjectidKinds)

	var datedVehicleJourneyRef string
	if ok {
		datedVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultObjectID, ok := vehicleJourney.ObjectID("_default")
		if !ok {
			return "", false
		}
		datedVehicleJourneyRef = builder.referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultObjectID.Value()})
	}
	return datedVehicleJourneyRef, true
}

func (builder *BroadcastStopMonitoringBuilder) resolveOperator(references model.References) {
	operatorRef, ok := references.Get("OperatorRef")
	if !ok {
		return
	}
	operator, ok := builder.partner.Model().Operators().FindByObjectId(*operatorRef.ObjectId)
	if !ok {
		return
	}
	obj, ok := operator.ObjectID(builder.remoteObjectidKind)
	if !ok {
		return
	}
	ref, _ := references.Get("OperatorRef")
	ref.ObjectId.SetValue(obj.Value())
}

func (builder *BroadcastStopMonitoringBuilder) resolveVJReferences(references model.References, origin string) {
	for _, refType := range []string{"RouteRef", "JourneyPatternRef", "DatedVehicleJourneyRef"} {
		if refType == "JourneyPatternRef" && !builder.rewriteJourneyPatternRef {
			continue
		}
		reference, ok := references.Get(refType)
		if !ok {
			continue
		}
		reference.ObjectId.SetValue(builder.referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: refType[:len(refType)-3], Id: reference.GetSha1()}))
	}
	for _, refType := range []string{"PlaceRef", "OriginRef", "DestinationRef"} {
		reference, ok := references.Get(refType)
		if !ok || reference.ObjectId == nil || (refType == "DestinationRef" && builder.noDestinationRefRewrite(origin)) {
			continue
		}
		builder.resolveStopAreaRef(&reference)
	}
}

func (builder *BroadcastStopMonitoringBuilder) resolveStopAreaRef(reference *model.Reference) {
	stopArea, ok := builder.partner.Model().StopAreas().FindByObjectId(*reference.ObjectId)
	if ok {
		obj, ok := stopArea.ReferentOrSelfObjectId(builder.remoteObjectidKind)
		if ok {
			reference.ObjectId.SetValue(obj.Value())
			return
		}
	}
	reference.ObjectId.SetValue(builder.stopAreareferenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Id: reference.GetSha1()}))
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
	return builder.dataFrameGenerator.NewIdentifier(idgen.IdentifierAttributes{Id: modelDate.String()})
}
