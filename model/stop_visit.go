package model

import "encoding/json"

type StopVisitId string

type StopVisit struct {
	ObjectIDConsumer
	model Model

	id               StopVisitId
	stopAreaId       StopAreaId
	vehicleJourneyId VehicleJourneyId

	schedules       StopVisitSchedules
	departureStatus StopVisitDepartureStatus
	arrivalStatus   StopVisitArrivalStatus
	passageOrder    int
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

func (stopVisit *StopVisit) VehicleJourney() VehicleJourney {
	vehicleJourney, _ := stopVisit.model.VehicleJourneys().Find(stopVisit.vehicleJourneyId)
	return vehicleJourney
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

func (stopVisit *StopVisit) PassageOrder() int {
	return stopVisit.passageOrder
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
	byObjectId   map[string]map[string]StopVisitId
}

type StopVisits interface {
	UUIDInterface

	New() StopVisit
	Find(id StopVisitId) (StopVisit, bool)
	FindByObjectId(objectid ObjectID) (StopVisit, bool)
	FindAll() []StopVisit
	Save(stopVisit *StopVisit) bool
	Delete(stopVisit *StopVisit) bool
}

func NewMemoryStopVisits() *MemoryStopVisits {
	return &MemoryStopVisits{
		byIdentifier: make(map[StopVisitId]*StopVisit),
		byObjectId:   make(map[string]map[string]StopVisitId),
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

func (manager *MemoryStopVisits) FindByObjectId(objectid ObjectID) (StopVisit, bool) {
	valueMap, ok := manager.byObjectId[objectid.Kind()]
	if !ok {
		return StopVisit{}, false
	}
	id, ok := valueMap[objectid.Value()]
	if !ok {
		return StopVisit{}, false
	}
	return *manager.byIdentifier[id], true
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
	for _, objectid := range stopVisit.ObjectIDs() {
		_, ok := manager.byObjectId[objectid.Kind()]
		if !ok {
			manager.byObjectId[objectid.Kind()] = make(map[string]StopVisitId)
		}
		manager.byObjectId[objectid.Kind()][objectid.Value()] = stopVisit.Id()
	}
	return true
}

func (manager *MemoryStopVisits) Delete(stopVisit *StopVisit) bool {
	delete(manager.byIdentifier, stopVisit.Id())
	for _, objectid := range stopVisit.ObjectIDs() {
		valueMap := manager.byObjectId[objectid.Kind()]
		delete(valueMap, objectid.Value())
	}
	return true
}