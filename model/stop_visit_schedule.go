package model

import (
	"encoding/json"
	"time"
)

type StopVisitScheduleType string

const (
	STOP_VISIT_SCHEDULE_AIMED    StopVisitScheduleType = "aimed"
	STOP_VISIT_SCHEDULE_EXPECTED StopVisitScheduleType = "expected"
	STOP_VISIT_SCHEDULE_ACTUAL   StopVisitScheduleType = "actual"
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

func (schedule *StopVisitSchedule) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Kind":          schedule.kind,
		"DepartureTime": schedule.departureTime,
		"ArrivalTime":   schedule.arrivalTime,
	})
}

type StopVisitSchedules map[StopVisitScheduleType]*StopVisitSchedule

func NewStopVisitSchedules() StopVisitSchedules {
	return map[StopVisitScheduleType]*StopVisitSchedule{
		STOP_VISIT_SCHEDULE_AIMED:    &StopVisitSchedule{},
		STOP_VISIT_SCHEDULE_EXPECTED: &StopVisitSchedule{},
		STOP_VISIT_SCHEDULE_ACTUAL:   &StopVisitSchedule{},
	}
}

func (schedules StopVisitSchedules) SetSchedule(kind StopVisitScheduleType, departureTime time.Time, arrivalTime time.Time) {
	schedules[kind] = &StopVisitSchedule{
		kind:          kind,
		departureTime: departureTime,
		arrivalTime:   arrivalTime,
	}
}
