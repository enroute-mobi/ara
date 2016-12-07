package model

import "time"

type StopVisitDepartureStatus int

const (
	STOP_VISIT_DEPARTURE_ONTIME StopVisitDepartureStatus = iota
	STOP_VISIT_DEPARTURE_EARLY
	STOP_VISIT_DEPARTURE_DELAYED
	STOP_VISIT_DEPARTURE_CANCELLED
	STOP_VISIT_DEPARTURE_NOREPORT
)

type StopVisitArrivalStatus int

const (
	STOP_VISIT_ARRIVAL_ONTIME StopVisitArrivalStatus = iota
	STOP_VISIT_ARRIVAL_EARLY
	STOP_VISIT_ARRIVAL_DELAYED
	STOP_VISIT_ARRIVAL_CANCELLED
	STOP_VISIT_ARRIVAL_NOREPORT
	STOP_VISIT_ARRIVAL_MISSED
	STOP_VISIT_ARRIVAL_NOT_EXPECTED
)

type StopVisitScheduleType int

const (
	STOP_VISIT_SCHEDULE_AIMED StopVisitScheduleType = iota
	STOP_VISIT_SCHEDULE_EXPECTED
	STOP_VISIT_SCHEDULE_ARRIVAL
)

type StopVisitSchedule struct {
	kind          StopVisitScheduleType
	departureTime time.Time
	arrivalTime   time.Time
}

type StopVisitUpdateEvent struct {
	id                  string
	created_at          time.Time
	stop_visit_objectid ObjectID
	schedules           map[StopVisitScheduleType]StopVisitSchedule
	departureStatus     StopVisitDepartureStatus
	arrivalStatuts      StopVisitArrivalStatus
}

func (event *StopVisitUpdateEvent) Id() string {
	return event.id
}

func (event *StopVisitUpdateEvent) CreatedAt() time.Time {
	return event.created_at
}

func (event *StopVisitUpdateEvent) StopVisitObjectId() ObjectID {
	return event.stop_visit_objectid
}

func (event *StopVisitUpdateEvent) Schedules() map[StopVisitScheduleType]StopVisitSchedule {
	return event.schedules
}

func (event *StopVisitUpdateEvent) DepartureStatus() StopVisitDepartureStatus {
	return event.departureStatus
}

func (event *StopVisitUpdateEvent) ArrivalStatus() StopVisitArrivalStatus {
	return event.arrivalStatuts
}
