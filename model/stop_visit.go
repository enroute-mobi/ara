package model

import (
	"encoding/json"
	"time"
)

type StopVisitId string

type StopVisitAttributes struct {
	ObjectId         ObjectID
	StopAreaObjectId ObjectID

	VehicleJourneyObjectId ObjectID
	PassageOrder           int

	ArrivalStatus   StopVisitArrivalStatus
	DepartureStatus StopVisitDepartureStatus
	RecordedAt      time.Time
	Schedules       StopVisitSchedules
	VehicleAtStop   string
	Attributes      map[string]string
}

type StopVisit struct {
	ObjectIDConsumer
	model Model

	id               StopVisitId
	StopAreaId       StopAreaId
	VehicleJourneyId VehicleJourneyId
	Attributes       map[string]string
	References       map[string]Reference

	ArrivalStatus   StopVisitArrivalStatus
	DepartureStatus StopVisitDepartureStatus
	RecordedAt      time.Time
	Schedules       StopVisitSchedules
	VehicleAtStop   string

	PassageOrder int
}

func NewStopVisit(model Model) *StopVisit {
	stopVisit := &StopVisit{
		model:      model,
		Schedules:  NewStopVisitSchedules(),
		Attributes: make(map[string]string),
		References: make(map[string]Reference),
	}
	stopVisit.objectids = make(ObjectIDs)
	return stopVisit
}

func (stopVisit *StopVisit) Id() StopVisitId {
	return stopVisit.id
}

func (stopVisit *StopVisit) StopArea() StopArea {
	stopArea, _ := stopVisit.model.StopAreas().Find(stopVisit.StopAreaId)
	return stopArea
}

func (stopVisit *StopVisit) VehicleJourney() *VehicleJourney {
	vehicleJourney, _ := stopVisit.model.VehicleJourneys().Find(stopVisit.VehicleJourneyId)
	return &vehicleJourney
}

/* type ResponseInterface map[string]interface{}

func (orderMap *ResponseInterface) ToJson(order []string) string {
	orderedJson := &bytes.Buffer{}
	orderedJson.Write([]byte{'{', '\n'})
	l := len(order)
	for i, key := range order {
		if (*orderMap)[key] == nil {
			(*orderMap)[key] = ""
		}
		fmt.Fprintf(orderedJson, "\t\"%s\": \"%v\"", key, (*orderMap)[key])
		if i < l { // putting the ',' only if not last
			orderedJson.WriteByte(',')
		}
		orderedJson.WriteByte('\n')
	}
	orderedJson.Write([]byte{'}', '\n'})
	return orderedJson.String()
} */

func (stopVisit *StopVisit) MarshalJSON() ([]byte, error) {
	scheduleSlice := []StopVisitSchedule{}
	for _, schedule := range stopVisit.Schedules {
		scheduleSlice = append(scheduleSlice, *schedule)
	}

	stopVisitMap := map[string]interface{}{
		"Id":               stopVisit.id,
		"StopAreaId":       stopVisit.StopAreaId,
		"VehicleJourneyId": stopVisit.VehicleJourneyId,
		"VehicleAtStop":    stopVisit.VehicleAtStop,
		"PassageOrder":     stopVisit.PassageOrder,
		"RecordedAt":       stopVisit.RecordedAt,
		"Schedules":        scheduleSlice,
		"DepartureStatus":  stopVisit.DepartureStatus,
		"ArrivalStatus":    stopVisit.ArrivalStatus,
		"Attributes":       stopVisit.Attributes,
		"References":       stopVisit.References,
	}
	if !stopVisit.ObjectIDs().Empty() {
		stopVisitMap["ObjectIDs"] = stopVisit.ObjectIDs()
	}
	return json.Marshal(stopVisitMap)
}

func (stopVisit *StopVisit) UnmarshalJSON(data []byte) error {
	type Alias StopVisit
	aux := &struct {
		ObjectIDs        map[string]string
		Reference        map[string]Reference
		StopAreaId       string
		VehicleJourneyId string
		PassageOrder     int
		Schedules        []StopVisitSchedule
		*Alias
	}{
		Alias: (*Alias)(stopVisit),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ObjectIDs != nil {
		stopVisit.ObjectIDConsumer.objectids = NewObjectIDsFromMap(aux.ObjectIDs)
	}

	if aux.Schedules != nil {
		stopVisit.Schedules = NewStopVisitSchedules()
		for _, schedule := range aux.Schedules {
			stopVisit.Schedules.SetSchedule(schedule.Kind(), schedule.DepartureTime(), schedule.ArrivalTime())
		}
	}

	if aux.StopAreaId != "" {
		stopVisit.StopAreaId = StopAreaId(aux.StopAreaId)
	}
	if aux.VehicleJourneyId != "" {
		stopVisit.VehicleJourneyId = VehicleJourneyId(aux.VehicleJourneyId)
	}
	if aux.PassageOrder > 0 {
		stopVisit.PassageOrder = aux.PassageOrder
	}

	return nil
}

func (stopVisit *StopVisit) Attribute(key string) (string, bool) {
	value, present := stopVisit.Attributes[key]
	return value, present
}

func (stopVisit *StopVisit) Save() (ok bool) {
	ok = stopVisit.model.StopVisits().Save(stopVisit)
	return
}

func (stopVisit *StopVisit) Reference(key string) (Reference, bool) {
	value, present := stopVisit.References[key]
	return value, present
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
		if stopVisit.StopAreaId == id {
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
	manager.byVehicleJourneyId[stopVisit.VehicleJourneyId] = append(manager.byVehicleJourneyId[stopVisit.VehicleJourneyId], stopVisit.id)
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
	for i, stopVisitId := range manager.byVehicleJourneyId[stopVisit.VehicleJourneyId] {
		if stopVisitId == stopVisit.id {
			manager.byVehicleJourneyId[stopVisit.VehicleJourneyId] = append(manager.byVehicleJourneyId[stopVisit.VehicleJourneyId][:i], manager.byVehicleJourneyId[stopVisit.VehicleJourneyId][i+1:]...)
			break
		}
	}
	return true
}
