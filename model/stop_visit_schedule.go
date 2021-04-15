package model

import (
	"encoding/json"
	"sync"
	"time"
)

type StopVisitScheduleType string

const (
	STOP_VISIT_SCHEDULE_AIMED    StopVisitScheduleType = "aimed"
	STOP_VISIT_SCHEDULE_EXPECTED StopVisitScheduleType = "expected"
	STOP_VISIT_SCHEDULE_ACTUAL   StopVisitScheduleType = "actual"
)

var stopVisitScheduleTypes = [3]StopVisitScheduleType{STOP_VISIT_SCHEDULE_AIMED, STOP_VISIT_SCHEDULE_EXPECTED, STOP_VISIT_SCHEDULE_ACTUAL}

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

type StopVisitSchedules struct {
	sync.RWMutex

	byType map[StopVisitScheduleType]*StopVisitSchedule
}

func NewStopVisitSchedules() StopVisitSchedules {
	return StopVisitSchedules{byType: make(map[StopVisitScheduleType]*StopVisitSchedule)}
}

func (schedules *StopVisitSchedules) Copy() *StopVisitSchedules {
	cpy := NewStopVisitSchedules()

	schedules.RLock()
	for key, value := range schedules.byType {
		cpy.byType[key] = &StopVisitSchedule{
			kind:          value.Kind(),
			arrivalTime:   value.ArrivalTime(),
			departureTime: value.DepartureTime(),
		}
	}
	schedules.RUnlock()
	return &cpy
}

func (schedules *StopVisitSchedules) Merge(newSchedules *StopVisitSchedules) {
	schedules.Lock()
	newSchedules.RLock()
	for key, value := range newSchedules.byType {
		schedules.byType[key] = &StopVisitSchedule{
			kind:          value.Kind(),
			arrivalTime:   value.ArrivalTime(),
			departureTime: value.DepartureTime(),
		}
	}
	schedules.Unlock()
	newSchedules.RUnlock()
}

func (schedules *StopVisitSchedules) Include(scs *StopVisitSchedules) bool {
	schedules.RLock()
	scs.RLock()

	for k, v := range scs.byType {
		if !compareSchedules(v, schedules.byType[k]) {
			schedules.RUnlock()
			scs.RUnlock()
			return false
		}
	}

	schedules.RUnlock()
	scs.RUnlock()
	return true

}

func (schedules *StopVisitSchedules) Eq(scs *StopVisitSchedules) bool {
	schedules.RLock()
	scs.RLock()

	for i := range stopVisitScheduleTypes {
		if schedules.byType[stopVisitScheduleTypes[i]] == nil && scs.byType[stopVisitScheduleTypes[i]] == nil {
			continue
		}
		if !compareSchedules(schedules.byType[stopVisitScheduleTypes[i]], scs.byType[stopVisitScheduleTypes[i]]) {
			schedules.RUnlock()
			scs.RUnlock()
			return false
		}
	}

	schedules.RUnlock()
	scs.RUnlock()
	return true
}

func compareSchedules(sc1, sc2 *StopVisitSchedule) bool {
	if sc1 == nil || sc2 == nil {
		return false
	}
	if sc1.kind != sc2.kind {
		return false
	}
	if !sc1.arrivalTime.Equal(sc2.arrivalTime) || !sc1.departureTime.Equal(sc2.departureTime) {
		return false
	}
	return true
}

func (schedules *StopVisitSchedules) SetDepartureTime(kind StopVisitScheduleType, departureTime time.Time) {
	schedules.Lock()
	_, ok := schedules.byType[kind]
	if !ok {
		schedules.byType[kind] = &StopVisitSchedule{kind: kind}
	}
	schedules.byType[kind].SetDepartureTime(departureTime)
	schedules.Unlock()
}

func (schedules *StopVisitSchedules) SetArrivalTime(kind StopVisitScheduleType, arrivalTime time.Time) {
	schedules.Lock()
	_, ok := schedules.byType[kind]
	if !ok {
		schedules.byType[kind] = &StopVisitSchedule{kind: kind}
	}
	schedules.byType[kind].SetArrivalTime(arrivalTime)
	schedules.Unlock()
}

func (schedules *StopVisitSchedules) SetSchedule(kind StopVisitScheduleType, departureTime time.Time, arrivalTime time.Time) {
	schedules.Lock()
	schedules.byType[kind] = &StopVisitSchedule{
		kind:          kind,
		departureTime: departureTime,
		arrivalTime:   arrivalTime,
	}
	schedules.Unlock()
}

func (schedules *StopVisitSchedules) Schedule(kind StopVisitScheduleType) *StopVisitSchedule {
	schedules.RLock()
	schedule, ok := schedules.byType[kind]
	schedules.RUnlock()
	if !ok {
		return &StopVisitSchedule{}
	}
	return schedule
}

func (schedules *StopVisitSchedules) ArrivalTimeFromKind(kinds []StopVisitScheduleType) time.Time {
	if kinds == nil {
		kinds = []StopVisitScheduleType{"actual", "expected", "aimed"}
	}
	schedules.RLock()
	for _, kind := range kinds {
		if value, ok := schedules.byType[kind]; ok {
			schedules.RUnlock()
			return value.ArrivalTime()
		}
	}
	schedules.RUnlock()
	return time.Time{}
}

func (schedules *StopVisitSchedules) DepartureTimeFromKind(kinds []StopVisitScheduleType) time.Time {
	if kinds == nil {
		kinds = []StopVisitScheduleType{"actual", "expected", "aimed"}
	}
	schedules.RLock()
	for _, kind := range kinds {
		if value, ok := schedules.byType[kind]; ok {
			schedules.RUnlock()
			return value.DepartureTime()
		}
	}
	schedules.RUnlock()
	return time.Time{}
}

func (schedules *StopVisitSchedules) ToSlice() (scheduleSlice []StopVisitSchedule) {
	schedules.RLock()
	for _, schedule := range schedules.byType {
		scheduleSlice = append(scheduleSlice, *schedule)
	}
	schedules.RUnlock()
	return
}
