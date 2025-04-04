package schedules

import (
	"encoding/json"
	"sync"
	"time"
)

type StopVisitScheduleType string

const (
	Aimed    StopVisitScheduleType = "aimed"
	Expected StopVisitScheduleType = "expected"
	Actual   StopVisitScheduleType = "actual"
)

var ScheduleOrderArray = []StopVisitScheduleType{
	Actual,
	Expected,
	Aimed,
}

var stopVisitScheduleTypes = []StopVisitScheduleType{Aimed, Expected, Actual}

type StopVisitSchedule struct {
	departureTime time.Time
	arrivalTime   time.Time
	kind          StopVisitScheduleType
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
		DepartureTime time.Time
		ArrivalTime   time.Time
		Kind          StopVisitScheduleType
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
	byType map[StopVisitScheduleType]*StopVisitSchedule
	sync.RWMutex
}

func NewStopVisitSchedules() *StopVisitSchedules {
	return &StopVisitSchedules{byType: make(map[StopVisitScheduleType]*StopVisitSchedule, 3)}
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
	return cpy
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

func (schedules *StopVisitSchedules) SetDepartureTimeIfNotDefined(kind StopVisitScheduleType, departureTime time.Time) {
	schedules.Lock()
	defer schedules.Unlock()

	_, ok := schedules.byType[kind]
	if !ok {
		schedules.byType[kind] = &StopVisitSchedule{kind: kind}
	} else if !schedules.byType[kind].departureTime.IsZero() {
		return
	}
	schedules.byType[kind].SetDepartureTime(departureTime)
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

func (schedules *StopVisitSchedules) SetArrivalTimeIfNotDefined(kind StopVisitScheduleType, arrivalTime time.Time) {
	schedules.Lock()
	defer schedules.Unlock()

	_, ok := schedules.byType[kind]
	if !ok {
		schedules.byType[kind] = &StopVisitSchedule{kind: kind}
	} else if !schedules.byType[kind].arrivalTime.IsZero() {
		return
	}
	schedules.byType[kind].SetArrivalTime(arrivalTime)
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
	defer schedules.RUnlock()
	return schedules.schedule(kind)
}

func (schedules *StopVisitSchedules) schedule(kind StopVisitScheduleType) *StopVisitSchedule {
	schedule, ok := schedules.byType[kind]
	if !ok {
		return &StopVisitSchedule{}
	}
	return schedule
}

func (schedules *StopVisitSchedules) ArrivalTimeFromKind(kinds []StopVisitScheduleType) time.Time {
	if kinds == nil {
		kinds = ScheduleOrderArray
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
		kinds = ScheduleOrderArray
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

func (schedules *StopVisitSchedules) ReferenceTime() time.Time {
	schedules.RLock()
	defer schedules.RUnlock()
	if t := schedules.referenceArrivalTime(); !t.IsZero() {
		return t
	}
	return schedules.referenceDepartureTime()
}

func (schedules *StopVisitSchedules) ReferenceArrivalTime() time.Time {
	schedules.RLock()
	defer schedules.RUnlock()
	return schedules.referenceArrivalTime()
}

func (schedules *StopVisitSchedules) referenceArrivalTime() time.Time {
	for _, kind := range ScheduleOrderArray {
		if s := schedules.schedule(kind).arrivalTime; !s.IsZero() {
			return s
		}
	}
	return time.Time{}
}

func (schedules *StopVisitSchedules) ReferenceDepartureTime() time.Time {
	schedules.RLock()
	defer schedules.RUnlock()
	return schedules.referenceDepartureTime()
}

func (schedules *StopVisitSchedules) referenceDepartureTime() time.Time {
	for _, kind := range ScheduleOrderArray {
		if s := schedules.schedule(kind).departureTime; !s.IsZero() {
			return s
		}
	}
	return time.Time{}
}

func (schedules *StopVisitSchedules) SetDefaultAimedTimes() {
	schedules.Lock()
	defer schedules.Unlock()

	_, ok := schedules.byType[Expected]
	if !ok {
		schedules.byType[Expected] = &StopVisitSchedule{kind: Expected}
		return
	}

	_, ok = schedules.byType[Aimed]
	if !ok {
		schedules.byType[Aimed] = &StopVisitSchedule{kind: Aimed}
		schedules.byType[Aimed].arrivalTime = schedules.byType[Expected].arrivalTime
		schedules.byType[Aimed].departureTime = schedules.byType[Expected].departureTime
		return
	}

	if schedules.byType[Aimed].arrivalTime.IsZero() {
		schedules.byType[Aimed].arrivalTime = schedules.byType[Expected].arrivalTime
	}
	if schedules.byType[Aimed].departureTime.IsZero() {
		schedules.byType[Aimed].departureTime = schedules.byType[Expected].departureTime
	}
}
