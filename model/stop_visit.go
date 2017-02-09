package model

import "encoding/json"

type StopVisitId string

type StopVisitAttributes struct {
	ObjectId         ObjectID
	StopAreaObjectId ObjectID

	VehicleJourneyObjectId ObjectID
	PassageOrder           int

	DepartureStatus StopVisitDepartureStatus
	ArrivalStatus   StopVisitArrivalStatus
	Schedules       StopVisitSchedules
}

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

func (stopVisit *StopVisit) Schedules() (scheduleSlice []StopVisitSchedule) {
	if len(stopVisit.schedules) == 0 {
		return []StopVisitSchedule{}
	}
	for _, schedule := range stopVisit.schedules {
		scheduleSlice = append(scheduleSlice, *schedule)
	}
	return
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

func (stopVisit *StopVisit) MarshalJSON() ([]byte, error) {
	stopVisitMap := map[string]interface{}{
		"Id":              stopVisit.id,
		"StopArea":        stopVisit.stopAreaId,
		"VehicleJourney":  stopVisit.vehicleJourneyId,
		"PassageOrder":    stopVisit.passageOrder,
		"Schedules":       stopVisit.Schedules(),
		"DepartureStatus": stopVisit.departureStatus,
		"ArrivalStatus":   stopVisit.arrivalStatus,
	}
	if stopVisit.ObjectIDs() != nil {
		stopVisitMap["ObjectIDs"] = stopVisit.ObjectIDs()
	}
	return json.Marshal(stopVisitMap)
}

func (stopVisit *StopVisit) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ObjectIDs map[string]string
	}{}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ObjectIDs != nil {
		stopVisit.ObjectIDConsumer.objectids = NewObjectIDsFromMap(aux.ObjectIDs)
	}

	return nil
}

func (stopVisit *StopVisit) Save() (ok bool) {
	ok = stopVisit.model.StopVisits().Save(stopVisit)
	return
}

type MemoryStopVisits struct {
	UUIDConsumer

	model Model

	byIdentifier       map[StopVisitId]*StopVisit
	byObjectId         map[string]map[string]StopVisitId
	byVehicleJourneyId map[VehicleJourneyId][]StopVisitId
}

type StopVisits interface {
	UUIDInterface

	New() StopVisit
	Find(id StopVisitId) (StopVisit, bool)
	FindByObjectId(objectid ObjectID) (StopVisit, bool)
	FindByVehicleJourneyId(id VehicleJourneyId) []StopVisit
	FindByStopAreaId(id StopAreaId) []StopVisit
	FindAll() []StopVisit
	Save(stopVisit *StopVisit) bool
	Delete(stopVisit *StopVisit) bool
}

func NewMemoryStopVisits() *MemoryStopVisits {
	return &MemoryStopVisits{
		byIdentifier:       make(map[StopVisitId]*StopVisit),
		byObjectId:         make(map[string]map[string]StopVisitId),
		byVehicleJourneyId: make(map[VehicleJourneyId][]StopVisitId),
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

func (manager *MemoryStopVisits) FindByVehicleJourneyId(id VehicleJourneyId) (stopVisits []StopVisit) {
	stopVisitIds, ok := manager.byVehicleJourneyId[id]
	if !ok {
		return []StopVisit{}
	}
	for _, stopVisitId := range stopVisitIds {
		stopVisits = append(stopVisits, *manager.byIdentifier[stopVisitId])
	}
	return
}

// Temp
func (manager *MemoryStopVisits) FindByStopAreaId(id StopAreaId) (stopVisits []StopVisit) {
	for _, stopVisit := range manager.byIdentifier {
		if stopVisit.stopAreaId == id {
			stopVisits = append(stopVisits, *stopVisit)
		}
	}
	return
}

func (manager *MemoryStopVisits) FindAll() (stopVisits []StopVisit) {
	if len(manager.byIdentifier) == 0 {
		return []StopVisit{}
	}
	for _, stopVisit := range manager.byIdentifier {
		stopVisits = append(stopVisits, *stopVisit)
	}
	return
}

func (manager *MemoryStopVisits) Save(stopVisit *StopVisit) bool {
	if stopVisit.id == "" {
		stopVisit.id = StopVisitId(manager.NewUUID())
	}
	stopVisit.model = manager.model
	manager.byIdentifier[stopVisit.id] = stopVisit
	manager.byVehicleJourneyId[stopVisit.vehicleJourneyId] = append(manager.byVehicleJourneyId[stopVisit.vehicleJourneyId], stopVisit.id)
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
	// Delete in byObjectId
	for _, objectid := range stopVisit.ObjectIDs() {
		valueMap := manager.byObjectId[objectid.Kind()]
		delete(valueMap, objectid.Value())
	}
	// Delete in byVehicleJourneyId
	for i, stopVisitId := range manager.byVehicleJourneyId[stopVisit.vehicleJourneyId] {
		if stopVisitId == stopVisit.id {
			manager.byVehicleJourneyId[stopVisit.vehicleJourneyId] = append(manager.byVehicleJourneyId[stopVisit.vehicleJourneyId][:i], manager.byVehicleJourneyId[stopVisit.vehicleJourneyId][i+1:]...)
			break
		}
	}
	return true
}
