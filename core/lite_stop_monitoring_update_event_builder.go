package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/slite"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type LiteStopMonitoringUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	originStopAreaCode model.Code
	partner            *Partner
	remoteCodeSpace    string

	stopMonitoringUpdateEvents *CollectUpdateEvents
}

func NewLiteStopMonitoringUpdateEventBuilder(partner *Partner, originStopAreaCode model.Code) LiteStopMonitoringUpdateEventBuilder {
	return LiteStopMonitoringUpdateEventBuilder{
		originStopAreaCode:         originStopAreaCode,
		partner:                    partner,
		remoteCodeSpace:            partner.RemoteCodeSpace(),
		stopMonitoringUpdateEvents: NewCollectUpdateEvents(),
	}
}

func (builder *LiteStopMonitoringUpdateEventBuilder) buildUpdateEvents(StopVisitEvent *slite.MonitoredStopVisit) {
	// When Order is not defined, we should ignore the MonitoredStopVisit
	// see ARA-1240 "Special cases"
	if !StopVisitEvent.HasOrder() {
		return
	}

	origin := string(builder.partner.Slug())
	stopPointRef := StopVisitEvent.GetStopPointRef()

	// StopAreas
	stopAreaCode := model.NewCode(builder.remoteCodeSpace, stopPointRef)

	_, ok := builder.stopMonitoringUpdateEvents.StopAreas[stopPointRef]
	if !ok {
		// CollectedAlways is false by default
		event := &model.StopAreaUpdateEvent{
			Origin: origin,
			Code:   stopAreaCode,
			Name:   StopVisitEvent.MonitoredVehicleJourney.MonitoredCall.StopPointName,
		}
		if builder.originStopAreaCode.Value() != "" && stopAreaCode.String() != builder.originStopAreaCode.String() {
			event.ParentCode = builder.originStopAreaCode
		}

		builder.stopMonitoringUpdateEvents.StopAreas[stopPointRef] = event
		builder.stopMonitoringUpdateEvents.MonitoringRefs[stopPointRef] = struct{}{}
	}

	// Lines
	lineCode := model.NewCode(builder.remoteCodeSpace, StopVisitEvent.MonitoredVehicleJourney.LineRef)

	_, ok = builder.stopMonitoringUpdateEvents.Lines[StopVisitEvent.MonitoredVehicleJourney.LineRef]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin: origin,
			Code:   lineCode,
		}

		lineRef := StopVisitEvent.MonitoredVehicleJourney.LineRef
		builder.stopMonitoringUpdateEvents.Lines[lineRef] = lineEvent
		builder.stopMonitoringUpdateEvents.LineRefs[lineRef] = struct{}{}
	}

	// VehicleJourneys
	vehicleJourneyRef := StopVisitEvent.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef
	vehicleJourneyCode := model.NewCode(builder.remoteCodeSpace, vehicleJourneyRef)

	_, ok = builder.
		stopMonitoringUpdateEvents.
		VehicleJourneys[vehicleJourneyRef]

	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:          origin,
			Code:            vehicleJourneyCode,
			LineCode:        lineCode,
			DestinationRef:  StopVisitEvent.MonitoredVehicleJourney.DestinationRef,
			DestinationName: StopVisitEvent.MonitoredVehicleJourney.DestinationName,
			Monitored:       true,

			CodeSpace: builder.remoteCodeSpace,
		}

		builder.stopMonitoringUpdateEvents.VehicleJourneys[vehicleJourneyRef] = vjEvent
		builder.stopMonitoringUpdateEvents.VehicleJourneyRefs[vehicleJourneyRef] = struct{}{}
	}

	// StopVisits
	stopVisitCode := model.NewCode(builder.remoteCodeSpace, StopVisitEvent.GetItemIdentifier())

	monitoredCall := StopVisitEvent.MonitoredVehicleJourney.MonitoredCall
	_, ok = builder.stopMonitoringUpdateEvents.StopVisits[stopPointRef][StopVisitEvent.GetItemIdentifier()]
	if !ok {
		svEvent := &model.StopVisitUpdateEvent{
			Origin:             origin,
			Code:               stopVisitCode,
			StopAreaCode:       stopAreaCode,
			VehicleJourneyCode: vehicleJourneyCode,
			DataFrameRef:       StopVisitEvent.MonitoredVehicleJourney.FramedVehicleJourneyRef.DataFrameRef,
			PassageOrder:       *monitoredCall.Order,
			VehicleAtStop:      monitoredCall.VehicleAtStop,
			ArrivalStatus:      model.SetStopVisitArrivalStatus(monitoredCall.ArrivalStatus),
			DepartureStatus:    model.SetStopVisitDepartureStatus(monitoredCall.DepartureStatus),
			RecordedAt:         StopVisitEvent.RecordedAtTime,
			Schedules:          schedules.NewStopVisitSchedules(),
			Monitored:          StopVisitEvent.GetMonitored(),

			CodeSpace: builder.remoteCodeSpace,
		}

		aimedDerpatureTime := monitoredCall.AimedDepartureTime
		aimedArrivalTime := monitoredCall.AimedArrivalTime
		if !aimedDerpatureTime.IsZero() || !aimedArrivalTime.IsZero() {
			svEvent.Schedules.SetSchedule(schedules.Aimed, aimedDerpatureTime, aimedArrivalTime)
		}

		expectedArrivalTime := monitoredCall.ExpectedArrivalTime
		expectedDepartureTime := monitoredCall.ExpectedDepartureTime
		if !expectedDepartureTime.IsZero() || !expectedArrivalTime.IsZero() {
			svEvent.Schedules.SetSchedule(schedules.Expected, expectedDepartureTime, expectedArrivalTime)
		}

		actualArrivalTime := monitoredCall.ActualArrivalTime
		actualDepartureTime := monitoredCall.ActualDepartureTime
		if !actualDepartureTime.IsZero() || !actualArrivalTime.IsZero() {
			svEvent.Schedules.SetSchedule(schedules.Actual, actualDepartureTime, actualArrivalTime)
		}

		if builder.stopMonitoringUpdateEvents.StopVisits[stopPointRef] == nil {
			builder.stopMonitoringUpdateEvents.StopVisits[stopPointRef] = make(map[string]*model.StopVisitUpdateEvent)
		}
		builder.stopMonitoringUpdateEvents.StopVisits[stopPointRef][StopVisitEvent.GetItemIdentifier()] = svEvent

	}
}

func (builder *LiteStopMonitoringUpdateEventBuilder) SetUpdateEvents(stopVisits []slite.MonitoredStopVisit) {
	for _, StopVisitEvent := range stopVisits {
		builder.buildUpdateEvents(&StopVisitEvent)
	}
}

func (builder *LiteStopMonitoringUpdateEventBuilder) UpdateEvents() CollectUpdateEvents {
	return *builder.stopMonitoringUpdateEvents
}
