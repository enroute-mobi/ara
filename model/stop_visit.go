package model

import "encoding/json"

type StopVisitId string

type StopVisit struct {
	ObjectIDConsumer
	model Model

	id              StopVisitId
	stopAreaId      StopAreaId
	schedules       StopVisitSchedules
	departureStatus StopVisitDepartureStatus
	arrivalStatus   StopVisitArrivalStatus

	// WIP: Add VehiculeJourney ref and passage order
}

func NewStopVisit(model Model) *StopVisit {
	stopVisit := &StopVisit{model: model}
	stopVisit.objectids = make(ObjectIDs)
	return stopVisit
}

func (stopVisit *StopVisit) Id() StopVisitId {
	return stopVisit.id
}

func (stopVisit *StopVisit) StopArea() StopArea {
	stopArea, _ := stopVisit.model.StopAreas().Find(stopVisit.stopAreaId)
	return stopArea
}

func (stopVisit *StopVisit) Schedules() StopVisitSchedules {
	return stopVisit.schedules
}

func (stopVisit *StopVisit) DepartureStatus() StopVisitDepartureStatus {
	return stopVisit.departureStatus
}

func (stopVisit *StopVisit) ArrivalStatus() StopVisitArrivalStatus {
	return stopVisit.arrivalStatus
}

// WIP
func (stopVisit *StopVisit) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Id": stopVisit.id,
	})
}

func (stopVisit *StopVisit) Save() (ok bool) {
	ok = stopVisit.model.StopVisits().Save(stopVisit)
	return
}

type MemoryStopVisits struct {
	UUIDConsumer

	model Model

	byIdentifier map[StopVisitId]*StopVisit
}

type StopVisits interface {
	UUIDInterface

	New() StopVisit
	Find(id StopVisitId) (StopVisit, bool)
	FindAll() []StopVisit
	Save(stopVisit *StopVisit) bool
	Delete(stopVisit *StopVisit) bool
}

func NewMemoryStopVisits() *MemoryStopVisits {
	return &MemoryStopVisits{
		byIdentifier: make(map[StopVisitId]*StopVisit),
	}
}

func (manager *MemoryStopVisits) New() StopVisit {
	stopVisit := NewStopVisit(manager.model)
	return *stopVisit
}

func (manager *MemoryStopVisits) Find(id StopVisitId) (StopVisit, bool) {
	stopVisit, ok := manager.byIdentifier[id]
	if ok {
		return *stopVisit, true
	} else {
		return StopVisit{}, false
	}
}

func (manager *MemoryStopVisits) FindAll() (stopVisits []StopVisit) {
	for _, stopVisit := range manager.byIdentifier {
		stopVisits = append(stopVisits, *stopVisit)
	}
	return
}

func (manager *MemoryStopVisits) Save(stopVisit *StopVisit) bool {
	if stopVisit.Id() == "" {
		stopVisit.id = StopVisitId(manager.NewUUID())
	}
	stopVisit.model = manager.model
	manager.byIdentifier[stopVisit.Id()] = stopVisit
	return true
}

func (manager *MemoryStopVisits) Delete(stopVisit *StopVisit) bool {
	delete(manager.byIdentifier, stopVisit.Id())
	return true
}
