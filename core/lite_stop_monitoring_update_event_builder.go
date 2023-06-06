package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/slite"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type LiteStopMonitoringUpdateEventBuilder struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	originStopAreaObjectId model.ObjectID
	partner                *Partner
	remoteObjectidKind     string

	stopMonitoringUpdateEvents *CollectUpdateEvents
}

func NewLiteStopMonitoringUpdateEventBuilder(partner *Partner, originStopAreaObjectId model.ObjectID) LiteStopMonitoringUpdateEventBuilder {
	return LiteStopMonitoringUpdateEventBuilder{
		originStopAreaObjectId:     originStopAreaObjectId,
		partner:                    partner,
		remoteObjectidKind:         partner.RemoteObjectIDKind(),
		stopMonitoringUpdateEvents: NewCollectUpdateEvents(),
	}
}

func (builder *LiteStopMonitoringUpdateEventBuilder) buildUpdateEvents(StopVisitEvent *slite.MonitoredStopVisit) {
	origin := string(builder.partner.Slug())
	stopPointRef := StopVisitEvent.GetStopPointRef()

	// StopAreas
	stopAreaObjectId := model.NewObjectID(builder.remoteObjectidKind, stopPointRef)

	_, ok := builder.stopMonitoringUpdateEvents.StopAreas[stopPointRef]
	if !ok {
		// CollectedAlways is false by default
		event := &model.StopAreaUpdateEvent{
			Origin:   origin,
			ObjectId: stopAreaObjectId,
			Name:     StopVisitEvent.MonitoredVehicleJourney.MonitoredCall.StopPointName,
		}
		if builder.originStopAreaObjectId.Value() != "" && stopAreaObjectId.String() != builder.originStopAreaObjectId.String() {
			event.ParentObjectId = builder.originStopAreaObjectId
		}

		builder.stopMonitoringUpdateEvents.StopAreas[stopPointRef] = event
		builder.stopMonitoringUpdateEvents.MonitoringRefs[stopPointRef] = struct{}{}
	}

	// Lines
	lineObjectId := model.NewObjectID(builder.remoteObjectidKind, StopVisitEvent.MonitoredVehicleJourney.LineRef)

	_, ok = builder.stopMonitoringUpdateEvents.Lines[StopVisitEvent.MonitoredVehicleJourney.LineRef]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin:   origin,
			ObjectId: lineObjectId,
		}

		builder.stopMonitoringUpdateEvents.Lines[StopVisitEvent.MonitoredVehicleJourney.LineRef] = lineEvent
	}

	// VehicleJourneys
	vjObjectId := model.NewObjectID(builder.remoteObjectidKind,
		StopVisitEvent.
			MonitoredVehicleJourney.
			FramedVehicleJourneyRef.
			DatedVehicleJourneyRef)

	_, ok = builder.
		stopMonitoringUpdateEvents.
		VehicleJourneys[StopVisitEvent.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef]

	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:          origin,
			ObjectId:        vjObjectId,
			LineObjectId:    lineObjectId,
			DestinationRef:  StopVisitEvent.MonitoredVehicleJourney.DestinationRef,
			DestinationName: StopVisitEvent.MonitoredVehicleJourney.DestinationName,
			Monitored:       true,

			ObjectidKind: builder.remoteObjectidKind,
		}

		builder.stopMonitoringUpdateEvents.VehicleJourneys[StopVisitEvent.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef] = vjEvent
	}

	// StopVisits
	stopVisitObjectId := model.NewObjectID(builder.remoteObjectidKind, StopVisitEvent.GetItemIdentifier())

	_, ok = builder.stopMonitoringUpdateEvents.StopVisits[stopPointRef][StopVisitEvent.GetItemIdentifier()]
	if !ok {
		svEvent := &model.StopVisitUpdateEvent{
			Origin:                 origin,
			ObjectId:               stopVisitObjectId,
			StopAreaObjectId:       stopAreaObjectId,
			VehicleJourneyObjectId: vjObjectId,
			DataFrameRef:           StopVisitEvent.MonitoredVehicleJourney.FramedVehicleJourneyRef.DataFrameRef,
			PassageOrder:           StopVisitEvent.MonitoredVehicleJourney.MonitoredCall.Order,
			Monitored:              true,
			VehicleAtStop:          StopVisitEvent.MonitoredVehicleJourney.MonitoredCall.VehicleAtStop,
			ArrivalStatus:          model.SetStopVisitArrivalStatus(StopVisitEvent.MonitoredVehicleJourney.MonitoredCall.ArrivalStatus),
			DepartureStatus:        model.SetStopVisitDepartureStatus(StopVisitEvent.MonitoredVehicleJourney.MonitoredCall.DepartureStatus),
			RecordedAt:             StopVisitEvent.RecordedAtTime,
			Schedules:              model.NewStopVisitSchedules(),

			ObjectidKind: builder.remoteObjectidKind,
		}

		monitoredCall := StopVisitEvent.MonitoredVehicleJourney.MonitoredCall

		aimedDerpatureTime := monitoredCall.AimedDepartureTime
		aimedArrivalTime := monitoredCall.AimedArrivalTime
		if !aimedDerpatureTime.IsZero() || !aimedArrivalTime.IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_AIMED, aimedDerpatureTime, aimedArrivalTime)
		}

		expectedArrivalTime := monitoredCall.ExpectedArrivalTime
		expectedDepartureTime := monitoredCall.ExpectedDepartureTime
		if !expectedDepartureTime.IsZero() || !expectedArrivalTime.IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_EXPECTED, expectedDepartureTime, expectedArrivalTime)
		}

		actualArrivalTime := monitoredCall.ActualArrivalTime
		actualDepartureTime := monitoredCall.ActualDepartureTime
		if !actualDepartureTime.IsZero() || !actualArrivalTime.IsZero() {
			svEvent.Schedules.SetSchedule(model.STOP_VISIT_SCHEDULE_ACTUAL, actualDepartureTime, actualArrivalTime)
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
