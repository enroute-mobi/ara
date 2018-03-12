package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type BroadcastStopMonitoringBuilder struct {
	model.ClockConsumer
	model.UUIDConsumer

	StopVisitTypes string
	MonitoringRef  string

	tx                         *model.Transaction
	siriPartner                *SIRIPartner
	referenceGenerator         *IdentifierGenerator
	stopAreareferenceGenerator *IdentifierGenerator
	dataFrameGenerator         *IdentifierGenerator
	remoteObjectidKind         string
}

func NewBroadcastStopMonitoringBuilder(tx *model.Transaction, siriPartner *SIRIPartner, connector string) *BroadcastStopMonitoringBuilder {
	return &BroadcastStopMonitoringBuilder{
		tx:                         tx,
		siriPartner:                siriPartner,
		referenceGenerator:         siriPartner.IdentifierGenerator("reference_identifier"),
		stopAreareferenceGenerator: siriPartner.IdentifierGenerator("reference_stop_area_identifier"),
		dataFrameGenerator:         siriPartner.IdentifierGenerator("data_frame_identifier"),
		remoteObjectidKind:         siriPartner.Partner().RemoteObjectIDKind(connector),
	}
}

func (builder *BroadcastStopMonitoringBuilder) BuildCancelledStopVisit(stopVisit model.StopVisit) *siri.SIRICancelledStopVisit {
	vehicleJourney, ok := builder.tx.Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
	if !ok {
		logger.Log.Printf("Ignore CancelledStopVisit %s without Vehiclejourney", stopVisit.Id())
		return nil
	}
	line, ok := builder.tx.Model().Lines().Find(vehicleJourney.LineId)
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

	dataVehicleJourneyRef, ok := builder.dataVehicleJourneyRef(vehicleJourney)
	if !ok {
		return nil
	}

	modelDate := builder.tx.Model().Date()

	cancelledStopVisit := &siri.SIRICancelledStopVisit{
		RecordedAtTime:         stopVisit.RecordedAt,
		ItemRef:                itemIdentifier,
		MonitoringRef:          builder.MonitoringRef,
		LineRef:                lineObjectId.Value(),
		DatedVehicleJourneyRef: dataVehicleJourneyRef,
		DataFrameRef:           builder.dataFrameGenerator.NewIdentifier(IdentifierAttributes{Id: modelDate.String()}),
		PublishedLineName:      line.Name,
	}

	return cancelledStopVisit
}

func (builder *BroadcastStopMonitoringBuilder) BuildMonitoredStopVisit(stopVisit model.StopVisit) *siri.SIRIMonitoredStopVisit {
	stopPointRef, ok := builder.tx.Model().StopAreas().Find(stopVisit.StopAreaId)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %s without StopArea", stopVisit.Id())
		return nil
	}
	stopPointRefObjectId, ok := stopPointRef.ObjectID(builder.remoteObjectidKind)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %s with StopArea without correct ObjectID", stopVisit.Id())
		return nil
	}
	vehicleJourney, ok := builder.tx.Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
	if !ok {
		logger.Log.Printf("Ignore StopVisit %s without Vehiclejourney", stopVisit.Id())
		return nil
	}
	line, ok := builder.tx.Model().Lines().Find(vehicleJourney.LineId)
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

	schedules := stopVisit.Schedules

	dataVehicleJourneyRef, ok := builder.dataVehicleJourneyRef(vehicleJourney)
	if !ok {
		return nil
	}

	modelDate := builder.tx.Model().Date()

	monitoredStopVisit := &siri.SIRIMonitoredStopVisit{
		ItemIdentifier:         itemIdentifier,
		MonitoringRef:          builder.MonitoringRef,
		StopPointRef:           stopPointRefObjectId.Value(),
		StopPointName:          stopPointRef.Name,
		VehicleJourneyName:     vehicleJourney.Name,
		OriginName:             vehicleJourney.OriginName,
		DestinationName:        vehicleJourney.DestinationName,
		Monitored:              vehicleJourney.Monitored,
		LineRef:                lineObjectId.Value(),
		DatedVehicleJourneyRef: dataVehicleJourneyRef,
		DataFrameRef:           builder.dataFrameGenerator.NewIdentifier(IdentifierAttributes{Id: modelDate.String()}),
		RecordedAt:             stopVisit.RecordedAt,
		PublishedLineName:      line.Name,
		DepartureStatus:        string(stopVisit.DepartureStatus),
		ArrivalStatus:          string(stopVisit.ArrivalStatus),
		Order:                  stopVisit.PassageOrder,
		VehicleAtStop:          stopVisit.VehicleAtStop,
		Attributes:             make(map[string]map[string]string),
		References:             make(map[string]map[string]model.Reference),
	}
	if !stopPointRef.Monitored {
		monitoredStopVisit.Monitored = false
	}

	if stopVisit.ArrivalStatus != model.STOP_VISIT_ARRIVAL_CANCELLED && builder.StopVisitTypes != "departures" {
		monitoredStopVisit.AimedArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).ArrivalTime()
		monitoredStopVisit.ExpectedArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).ArrivalTime()
		if monitoredStopVisit.Monitored {
			monitoredStopVisit.ActualArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime()
		}
	}

	if stopVisit.DepartureStatus != model.STOP_VISIT_DEPARTURE_CANCELLED && builder.StopVisitTypes != "arrivals" {
		monitoredStopVisit.AimedDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime()
		monitoredStopVisit.ExpectedDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).DepartureTime()
		if monitoredStopVisit.Monitored {
			monitoredStopVisit.ActualDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime()
		}
	}

	vehicleJourneyRefCopy := vehicleJourney.References.Copy()
	stopVisitRefCopy := stopVisit.References.Copy()

	builder.resolveVJReferences(vehicleJourneyRefCopy)

	builder.resolveOperator(stopVisitRefCopy)

	monitoredStopVisit.Attributes["StopVisitAttributes"] = stopVisit.Attributes
	monitoredStopVisit.References["StopVisitReferences"] = stopVisitRefCopy.GetReferences()

	monitoredStopVisit.Attributes["VehicleJourneyAttributes"] = vehicleJourney.Attributes
	monitoredStopVisit.References["VehicleJourney"] = vehicleJourneyRefCopy.GetReferences()

	return monitoredStopVisit
}

func (builder *BroadcastStopMonitoringBuilder) getItemIdentifier(stopVisit model.StopVisit) (string, bool) {
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
		itemIdentifier = builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "Item", Default: defaultObjectID.Value()})
	}
	return itemIdentifier, true
}

func (builder *BroadcastStopMonitoringBuilder) dataVehicleJourneyRef(vehicleJourney model.VehicleJourney) (string, bool) {
	vehicleJourneyId, ok := vehicleJourney.ObjectID(builder.remoteObjectidKind)

	var dataVehicleJourneyRef string
	if ok {
		dataVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultObjectID, ok := vehicleJourney.ObjectID("_default")
		if !ok {
			return "", false
		}
		dataVehicleJourneyRef = builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "VehicleJourney", Default: defaultObjectID.Value()})
	}
	return dataVehicleJourneyRef, true
}

func (builder *BroadcastStopMonitoringBuilder) resolveOperator(references model.References) {
	operatorRef, ok := references.Get("OperatorRef")
	if !ok {
		return
	}
	operator, ok := builder.tx.Model().Operators().FindByObjectId(*operatorRef.ObjectId)
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

func (builder *BroadcastStopMonitoringBuilder) resolveVJReferences(references model.References) {
	for _, refType := range []string{"RouteRef", "JourneyPatternRef", "DatedVehicleJourneyRef"} {
		reference, ok := references.Get(refType)
		if !ok {
			continue
		}
		reference.ObjectId.SetValue(builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: refType[:len(refType)-3], Default: reference.GetSha1()}))
	}
	for _, refType := range []string{"PlaceRef", "OriginRef", "DestinationRef"} {
		reference, ok := references.Get(refType)
		if !ok || reference.ObjectId == nil {
			continue
		}
		builder.resolveStopAreaRef(&reference)
	}
}

func (builder *BroadcastStopMonitoringBuilder) resolveStopAreaRef(reference *model.Reference) {
	stopArea, ok := builder.tx.Model().StopAreas().FindByObjectId(*reference.ObjectId)
	if ok {
		obj, ok := stopArea.ObjectID(builder.remoteObjectidKind)
		if ok {
			reference.ObjectId.SetValue(obj.Value())
			return
		}
	}
	reference.ObjectId.SetValue(builder.stopAreareferenceGenerator.NewIdentifier(IdentifierAttributes{Default: reference.GetSha1()}))
}
