package core

import (
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type EstimatedTimetableUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner            *Partner
	referenceGenerator *idgen.IdentifierGenerator
	remoteObjectidKind string
	origin             string

	updateEvents *CollectUpdateEvents
}

func NewEstimatedTimetableUpdateEventBuilder(partner *Partner) EstimatedTimetableUpdateEventBuilder {
	return EstimatedTimetableUpdateEventBuilder{
		partner:            partner,
		referenceGenerator: partner.IdentifierGenerator(idgen.REFERENCE_IDENTIFIER),
		remoteObjectidKind: partner.RemoteObjectIDKind(),
		origin:             string(partner.Slug()),
		updateEvents:       NewCollectUpdateEvents(),
	}
}

func (builder *EstimatedTimetableUpdateEventBuilder) buildUpdateEvents(estimatedJourneyVersionFrame *sxml.XMLEstimatedJourneyVersionFrame) {
	// Lines
	lineObjectId := model.NewObjectID(builder.remoteObjectidKind, estimatedJourneyVersionFrame.LineRef())

	_, ok := builder.updateEvents.Lines[estimatedJourneyVersionFrame.LineRef()]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin:   builder.origin,
			ObjectId: lineObjectId,
			// Name:     estimatedJourneyVersionFrame.PublishedLineName(),
		}

		builder.updateEvents.Lines[estimatedJourneyVersionFrame.LineRef()] = lineEvent
		builder.updateEvents.LineRefs[estimatedJourneyVersionFrame.LineRef()] = struct{}{}
	}

	// VehicleJourneys
	vjObjectId := model.NewObjectID(builder.remoteObjectidKind, estimatedJourneyVersionFrame.DatedVehicleJourneyRef())

	_, ok = builder.updateEvents.VehicleJourneys[estimatedJourneyVersionFrame.DatedVehicleJourneyRef()]
	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:       builder.origin,
			ObjectId:     vjObjectId,
			LineObjectId: lineObjectId,
			OriginRef:    estimatedJourneyVersionFrame.OriginRef(),
			// OriginName:      estimatedJourneyVersionFrame.OriginName(),
			DirectionType:  builder.directionRef(estimatedJourneyVersionFrame.DirectionRef()),
			DestinationRef: estimatedJourneyVersionFrame.DestinationRef(),
			// DestinationName: estimatedJourneyVersionFrame.DestinationName(),
			Monitored: true,
			// Occupancy: model.NormalizedOccupancyName(estimatedJourneyVersionFrame.Occupancy()),

			ObjectidKind: builder.remoteObjectidKind,
			// SiriXML:      &estimatedJourneyVersionFrame.XMLMonitoredVehicleJourney,
		}

		builder.updateEvents.VehicleJourneys[estimatedJourneyVersionFrame.DatedVehicleJourneyRef()] = vjEvent
	}

	for _, call := range estimatedJourneyVersionFrame.EstimatedCalls() {
		builder.handleCall(vjObjectId, estimatedJourneyVersionFrame.RecordedAt(), estimatedJourneyVersionFrame.DatedVehicleJourneyRef(), call)
	}
	for _, call := range estimatedJourneyVersionFrame.RecordedCalls() {
		builder.handleCall(vjObjectId, estimatedJourneyVersionFrame.RecordedAt(), estimatedJourneyVersionFrame.DatedVehicleJourneyRef(), call)
	}
}

func (builder *EstimatedTimetableUpdateEventBuilder) handleCall(vjObjectId model.ObjectID, recordedAt time.Time, datedVehicleJourneyRef string, call *sxml.XMLCall) {
	// StopAreas
	stopAreaObjectId := model.NewObjectID(builder.remoteObjectidKind, call.StopPointRef())

	_, ok := builder.updateEvents.StopAreas[call.StopPointRef()]
	if !ok {
		// CollectedAlways is false by default
		event := &model.StopAreaUpdateEvent{
			Origin:   builder.origin,
			ObjectId: stopAreaObjectId,
			Name:     call.StopPointName(),
		}

		builder.updateEvents.StopAreas[call.StopPointRef()] = event
		builder.updateEvents.MonitoringRefs[call.StopPointRef()] = struct{}{}
	}

	// StopVisits
	stopVisitId := fmt.Sprintf("%s-%s", datedVehicleJourneyRef, strconv.Itoa(call.Order()))
	stopVisitObjectId := model.NewObjectID(builder.remoteObjectidKind, stopVisitId)

	_, ok = builder.updateEvents.StopVisits[call.StopPointRef()][stopVisitId]
	if !ok {
		svEvent := &model.StopVisitUpdateEvent{
			Origin:                 builder.origin,
			ObjectId:               stopVisitObjectId,
			StopAreaObjectId:       stopAreaObjectId,
			VehicleJourneyObjectId: vjObjectId,
			// DataFrameRef:           call.DataFrameRef(),
			PassageOrder:    call.Order(),
			Monitored:       true,
			VehicleAtStop:   call.VehicleAtStop(),
			ArrivalStatus:   model.SetStopVisitArrivalStatus(call.ArrivalStatus()),
			DepartureStatus: model.SetStopVisitDepartureStatus(call.DepartureStatus()),
			RecordedAt:      recordedAt,
			Schedules:       model.NewStopVisitSchedules(),

			ObjectidKind: builder.remoteObjectidKind,
		}

		if !call.AimedDepartureTime().IsZero() || !call.AimedArrivalTime().IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, call.AimedDepartureTime(), call.AimedArrivalTime())
		}
		if !call.ExpectedDepartureTime().IsZero() || !call.ExpectedArrivalTime().IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, call.ExpectedDepartureTime(), call.ExpectedArrivalTime())
		}
		if !call.ActualDepartureTime().IsZero() || !call.ActualArrivalTime().IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, call.ActualDepartureTime(), call.ActualArrivalTime())
		}

		if builder.updateEvents.StopVisits[call.StopPointRef()] == nil {
			builder.updateEvents.StopVisits[call.StopPointRef()] = make(map[string]*model.StopVisitUpdateEvent)
		}
		builder.updateEvents.StopVisits[call.StopPointRef()][stopVisitId] = svEvent

	}
}

func (builder *EstimatedTimetableUpdateEventBuilder) directionRef(direction string) (dir string) {
	in, out, err := builder.partner.PartnerSettings.SIRIDirectionType()
	if err {
		return direction
	}

	switch direction {
	case in:
		dir = model.VEHICLE_DIRECTION_INBOUND
	case out:
		dir = model.VEHICLE_DIRECTION_OUTBOUND
	default:
		dir = direction
	}

	return dir
}

func (builder *EstimatedTimetableUpdateEventBuilder) SetUpdateEvents(estimatedJourneyVersionFrames []*sxml.XMLEstimatedJourneyVersionFrame) {
	for _, estimatedJourneyVersionFrame := range estimatedJourneyVersionFrames {
		builder.buildUpdateEvents(estimatedJourneyVersionFrame)
	}
}

func (builder *EstimatedTimetableUpdateEventBuilder) UpdateEvents() CollectUpdateEvents {
	return *builder.updateEvents
}
