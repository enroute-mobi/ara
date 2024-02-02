package core

import (
	"fmt"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type EstimatedTimetableUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	partner         *Partner
	remoteCodeSpace string
	origin          string

	updateEvents *CollectUpdateEvents
}

func NewEstimatedTimetableUpdateEventBuilder(partner *Partner) EstimatedTimetableUpdateEventBuilder {
	return EstimatedTimetableUpdateEventBuilder{
		partner:         partner,
		remoteCodeSpace: partner.RemoteCodeSpace(),
		origin:          string(partner.Slug()),
		updateEvents:    NewCollectUpdateEvents(),
	}
}

func (builder *EstimatedTimetableUpdateEventBuilder) buildUpdateEvents(estimatedJourneyVersionFrame *sxml.XMLEstimatedJourneyVersionFrame) {
	for _, estimatedVehicleJourney := range estimatedJourneyVersionFrame.EstimatedVehicleJourneys() {
		// Lines
		lineCode := model.NewCode(builder.remoteCodeSpace, estimatedVehicleJourney.LineRef())

		_, ok := builder.updateEvents.Lines[estimatedVehicleJourney.LineRef()]
		if !ok {
			// CollectedAlways is false by default
			lineEvent := &model.LineUpdateEvent{
				Origin: builder.origin,
				Code:   lineCode,
				// Name:     estimatedVehicleJourney.PublishedLineName(),
			}

			builder.updateEvents.Lines[estimatedVehicleJourney.LineRef()] = lineEvent
			builder.updateEvents.LineRefs[estimatedVehicleJourney.LineRef()] = struct{}{}
		}

		// VehicleJourneys
		vjCode := model.NewCode(builder.remoteCodeSpace, estimatedVehicleJourney.DatedVehicleJourneyRef())

		_, ok = builder.updateEvents.VehicleJourneys[estimatedVehicleJourney.DatedVehicleJourneyRef()]
		if !ok {
			vjEvent := &model.VehicleJourneyUpdateEvent{
				Origin:    builder.origin,
				Code:      vjCode,
				LineCode:  lineCode,
				OriginRef: estimatedVehicleJourney.OriginRef(),
				// OriginName:      estimatedVehicleJourney.OriginName(),
				DirectionType:  builder.directionRef(estimatedVehicleJourney.DirectionRef()),
				DestinationRef: estimatedVehicleJourney.DestinationRef(),
				// DestinationName: estimatedVehicleJourney.DestinationName(),
				Monitored: true,
				// Occupancy: model.NormalizedOccupancyName(estimatedVehicleJourney.Occupancy()),

				CodeSpace: builder.remoteCodeSpace,
				// SiriXML:      &estimatedVehicleJourney.XMLMonitoredVehicleJourney,
			}

			builder.updateEvents.VehicleJourneys[estimatedVehicleJourney.DatedVehicleJourneyRef()] = vjEvent
			builder.updateEvents.VehicleJourneyRefs[estimatedVehicleJourney.DatedVehicleJourneyRef()] = struct{}{}
		}

		for _, call := range estimatedVehicleJourney.EstimatedCalls() {
			builder.handleCall(vjCode, estimatedJourneyVersionFrame.RecordedAt(), estimatedVehicleJourney.DatedVehicleJourneyRef(), call)
		}
		for _, call := range estimatedVehicleJourney.RecordedCalls() {
			builder.handleCall(vjCode, estimatedJourneyVersionFrame.RecordedAt(), estimatedVehicleJourney.DatedVehicleJourneyRef(), call)
		}
	}
}

func (builder *EstimatedTimetableUpdateEventBuilder) handleCall(vjCode model.Code, recordedAt time.Time, datedVehicleJourneyRef string, call *sxml.XMLCall) {
	// StopAreas
	stopAreaCode := model.NewCode(builder.remoteCodeSpace, call.StopPointRef())

	_, ok := builder.updateEvents.StopAreas[call.StopPointRef()]
	if !ok {
		// CollectedAlways is false by default
		event := &model.StopAreaUpdateEvent{
			Origin: builder.origin,
			Code:   stopAreaCode,
			Name:   call.StopPointName(),
		}

		builder.updateEvents.StopAreas[call.StopPointRef()] = event
		builder.updateEvents.MonitoringRefs[call.StopPointRef()] = struct{}{}
	}

	// StopVisits
	stopVisitId := fmt.Sprintf("%s-%s", datedVehicleJourneyRef, strconv.Itoa(call.Order()))
	stopVisitCode := model.NewCode(builder.remoteCodeSpace, stopVisitId)

	_, ok = builder.updateEvents.StopVisits[call.StopPointRef()][stopVisitId]
	if !ok {
		svEvent := &model.StopVisitUpdateEvent{
			Origin:             builder.origin,
			Code:               stopVisitCode,
			StopAreaCode:       stopAreaCode,
			VehicleJourneyCode: vjCode,
			// DataFrameRef:           call.DataFrameRef(),
			PassageOrder:    call.Order(),
			Monitored:       true,
			VehicleAtStop:   call.VehicleAtStop(),
			ArrivalStatus:   model.SetStopVisitArrivalStatus(call.ArrivalStatus()),
			DepartureStatus: model.SetStopVisitDepartureStatus(call.DepartureStatus()),
			RecordedAt:      recordedAt,
			Schedules:       model.NewStopVisitSchedules(),

			CodeSpace: builder.remoteCodeSpace,
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
