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
	jsonSchedule := map[string]interface{}{
		"Kind": schedule.kind,
	}
	if !schedule.departureTime.IsZero() {
		jsonSchedule["DepartureTime"] = schedule.departureTime
	}
	if !schedule.arrivalTime.IsZero() {
		jsonSchedule["ArrivalTime"] = schedule.arrivalTime
	}
	return json.Marshal(jsonSchedule)
}

func (schedule *StopVisitSchedule) UnmarshalJSON(data []byte) error {

	aux := &struct {
		Kind          StopVisitScheduleType
		DepartureTime time.Time
		ArrivalTime   time.Time
	}{}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	schedule.kind = aux.Kind
	schedule.departureTime = aux.DepartureTime
	schedule.arrivalTime = aux.ArrivalTime

	return nil
}

type StopVisitSchedules map[StopVisitScheduleType]*StopVisitSchedule

func NewStopVisitSchedules() StopVisitSchedules {
	schedules := make(StopVisitSchedules)
	return schedules
}

func (schedules StopVisitSchedules) SetSchedule(kind StopVisitScheduleType, departureTime time.Time, arrivalTime time.Time) {
	_, ok := schedules[kind]
	if !ok {
		schedules[kind] = &StopVisitSchedule{}
	}
	schedules[kind] = &StopVisitSchedule{
		kind:          kind,
		departureTime: departureTime,
		arrivalTime:   arrivalTime,
	}
}

func (schedules StopVisitSchedules) Schedule(kind StopVisitScheduleType) *StopVisitSchedule {
	schedule, ok := schedules[kind]
	if !ok {
		return &StopVisitSchedule{}
	}
	return schedule
}
