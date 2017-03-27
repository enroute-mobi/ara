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

func (schedule *StopVisitSchedule) SetArrivalTime(t time.Time) time.Time {
	schedule.arrivalTime = t
	return t
}

func (schedule *StopVisitSchedule) SetDepartureTime(t time.Time) time.Time {
	schedule.departureTime = t
	return t
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

func (schedules *StopVisitSchedules) Merge(newSchedules StopVisitSchedules) {
	for key, value := range newSchedules {
		(*schedules)[key] = value
	}
}

func (schedules StopVisitSchedules) SetDepartureTime(kind StopVisitScheduleType, departureTime time.Time) {
	_, ok := schedules[kind]
	if !ok {
		schedules[kind] = &StopVisitSchedule{kind: kind}
	}
	schedules[kind].SetDepartureTime(departureTime)
}

func (schedules StopVisitSchedules) SetArrivalTime(kind StopVisitScheduleType, arrivalTime time.Time) {
	_, ok := schedules[kind]
	if !ok {
		schedules[kind] = &StopVisitSchedule{kind: kind}
	}
	schedules[kind].SetArrivalTime(arrivalTime)
}

func (schedules StopVisitSchedules) SetSchedule(kind StopVisitScheduleType, departureTime time.Time, arrivalTime time.Time) {
	_, ok := schedules[kind]
	if !ok {
		schedules[kind] = &StopVisitSchedule{kind: kind}
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

func (schedules StopVisitSchedules) ArrivalTimeFromKind(kinds []StopVisitScheduleType) time.Time {
	if kinds == nil {
		kinds = []StopVisitScheduleType{"actual", "expected", "aimed"}
	}
	for _, kind := range kinds {
		if value, ok := schedules[kind]; ok {
			return value.ArrivalTime()
		}
	}
	return time.Time{}
}

func (schedules StopVisitSchedules) DepartureTimeFromKind(kinds []StopVisitScheduleType) time.Time {
	if kinds == nil {
		kinds = []StopVisitScheduleType{"actual", "expected", "aimed"}
	}
	for _, kind := range kinds {
		if value, ok := schedules[kind]; ok {
			return value.DepartureTime()
		}
	}
	return time.Time{}
}
