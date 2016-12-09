package model

import "time"

type StopVisitDepartureStatus string

const (
	STOP_VISIT_DEPARTURE_ONTIME    StopVisitDepartureStatus = "ontime"
	STOP_VISIT_DEPARTURE_EARLY     StopVisitDepartureStatus = "early"
	STOP_VISIT_DEPARTURE_DELAYED   StopVisitDepartureStatus = "delayed"
	STOP_VISIT_DEPARTURE_CANCELLED StopVisitDepartureStatus = "cancelled"
	STOP_VISIT_DEPARTURE_NOREPORT  StopVisitDepartureStatus = "noreport"
	STOP_VISIT_DEPARTURE_UNDEFINED StopVisitDepartureStatus = ""
)

type StopVisitArrivalStatus string

const (
	STOP_VISIT_ARRIVAL_ARRIVED      StopVisitArrivalStatus = "arrived"
	STOP_VISIT_ARRIVAL_ONTIME       StopVisitArrivalStatus = "ontime"
	STOP_VISIT_ARRIVAL_EARLY        StopVisitArrivalStatus = "early"
	STOP_VISIT_ARRIVAL_DELAYED      StopVisitArrivalStatus = "delayed"
	STOP_VISIT_ARRIVAL_CANCELLED    StopVisitArrivalStatus = "cancelled"
	STOP_VISIT_ARRIVAL_NOREPORT     StopVisitArrivalStatus = "noreport"
	STOP_VISIT_ARRIVAL_MISSED       StopVisitArrivalStatus = "missed"
	STOP_VISIT_ARRIVAL_NOT_EXPECTED StopVisitArrivalStatus = "notExpected"
	STOP_VISIT_ARRIVAL_UNDEFINED    StopVisitArrivalStatus = ""
)

type StopVisitScheduleType int

const (
	STOP_VISIT_SCHEDULE_AIMED StopVisitScheduleType = iota
	STOP_VISIT_SCHEDULE_EXPECTED
	STOP_VISIT_SCHEDULE_ACTUAL
)

type StopVisitSchedule struct {
	kind          StopVisitScheduleType
	departureTime time.Time
	arrivalTime   time.Time
}

func (schedule *StopVisitSchedule) Kind() StopVisitScheduleType {
	return schedule.kind
}

func (schedule *StopVisitSchedule) DepartureTime() time.Time {
	return schedule.departureTime
}

func (schedule *StopVisitSchedule) ArrivalTime() time.Time {
	return schedule.arrivalTime
}

type StopVisitSchedules map[StopVisitScheduleType]StopVisitSchedule

func NewStopVisitSchedules() StopVisitSchedules {
	return map[StopVisitScheduleType]StopVisitSchedule{
		STOP_VISIT_SCHEDULE_AIMED:    StopVisitSchedule{},
		STOP_VISIT_SCHEDULE_EXPECTED: StopVisitSchedule{},
		STOP_VISIT_SCHEDULE_ACTUAL:   StopVisitSchedule{},
	}
}

func (event *StopVisitUpdateEvent) SetSchedule(kind StopVisitScheduleType, departureTime time.Time, arrivalTime time.Time) {
	event.Schedules[kind] = StopVisitSchedule{
		kind:          kind,
		departureTime: departureTime,
		arrivalTime:   arrivalTime,
	}
}

type StopVisitUpdateEvent struct {
	Id                  string
	Created_at          time.Time
	Stop_visit_objectid ObjectID
	Schedules           StopVisitSchedules
	DepartureStatus     StopVisitDepartureStatus
	ArrivalStatuts      StopVisitArrivalStatus
}
