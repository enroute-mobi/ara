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

	var itemIdentifier string
	stopVisitId, ok := stopVisit.ObjectID(builder.remoteObjectidKind)
	if ok {
		itemIdentifier = stopVisitId.Value()
	} else {
		defaultObjectID, ok := stopVisit.ObjectID("_default")
		if !ok {
			logger.Log.Printf("Ignore StopVisit %s without default ObjectID", stopVisit.Id())
			return nil
		}
		itemIdentifier = builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "Item", Default: defaultObjectID.Value()})
	}

	schedules := stopVisit.Schedules

	vehicleJourneyId, ok := vehicleJourney.ObjectID(builder.remoteObjectidKind)
	var dataVehicleJourneyRef string
	if ok {
		dataVehicleJourneyRef = vehicleJourneyId.Value()
	} else {
		defaultObjectID, ok := vehicleJourney.ObjectID("_default")
		if !ok {
			return nil
		}
		dataVehicleJourneyRef = builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "VehicleJourney", Default: defaultObjectID.Value()})
	}

	modelDate := builder.tx.Model().Date()

	monitoredStopVisit := &siri.SIRIMonitoredStopVisit{
		ItemIdentifier: itemIdentifier,
		MonitoringRef:  builder.MonitoringRef,
		StopPointRef:   stopPointRefObjectId.Value(),
		StopPointName:  stopPointRef.Name,

		VehicleJourneyName:     vehicleJourney.Name,
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
		if monitoredStopVisit.Monitored {
			monitoredStopVisit.ExpectedArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).ArrivalTime()
			monitoredStopVisit.ActualArrivalTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).ArrivalTime()
		}
	}

	if stopVisit.DepartureStatus != model.STOP_VISIT_DEPARTURE_CANCELLED && builder.StopVisitTypes != "arrivals" {
		monitoredStopVisit.AimedDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED).DepartureTime()
		if monitoredStopVisit.Monitored {
			monitoredStopVisit.ExpectedDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED).DepartureTime()
			monitoredStopVisit.ActualDepartureTime = schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL).DepartureTime()
		}
	}

	vehicleJourneyRefCopy := vehicleJourney.References.Copy()
	stopVisitRefCopy := stopVisit.References.Copy()

	builder.resolveVJReferences(vehicleJourneyRefCopy)

	builder.resolveOperator(stopVisitRefCopy)

	monitoredStopVisit.Attributes["StopVisitAttributes"] = stopVisit.Attributes
	monitoredStopVisit.References["StopVisitReferences"] = stopVisitRefCopy

	monitoredStopVisit.Attributes["VehicleJourneyAttributes"] = vehicleJourney.Attributes
	monitoredStopVisit.References["VehicleJourney"] = vehicleJourneyRefCopy

	return monitoredStopVisit
}

func (builder *BroadcastStopMonitoringBuilder) resolveOperator(references model.References) {
	operatorRef, ok := references["OperatorRef"]
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
	references["OperatorRef"].ObjectId.SetValue(obj.Value())
}

func (builder *BroadcastStopMonitoringBuilder) resolveVJReferences(references model.References) {
	for _, refType := range []string{"RouteRef", "JourneyPatternRef", "DatedVehicleJourneyRef"} {
		reference, ok := references[refType]
		if !ok {
			continue
		}
		reference.ObjectId.SetValue(builder.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: refType[:len(refType)-3], Default: reference.GetSha1()}))
	}
	for _, refType := range []string{"PlaceRef", "OriginRef", "DestinationRef"} {
		reference, ok := references[refType]
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
