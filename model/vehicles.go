package model

import (
	"encoding/json"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"cloud.google.com/go/civil"
)

type VehicleId string

var vehicleLineExtractor = func(instance ModelInstance) string { return string((instance.(*Vehicle)).LineId) }
var vehicleVjExtractor = func(instance ModelInstance) string { return string((instance.(*Vehicle)).VehicleJourneyId) }
var vehicleNextStopVisitExtractor = func(instance ModelInstance) string {
	return string((instance.(*Vehicle)).NextStopVisitId)
}

type Vehicle struct {
	RecordedAtTime time.Time
	ValidUntilTime time.Time
	model          Model
	CodeConsumer
	RawAttributes    RawAttributes
	StopAreaId       StopAreaId       `json:",omitempty"`
	Occupancy        string           `json:",omitempty"`
	LineId           LineId           `json:",omitempty"`
	VehicleJourneyId VehicleJourneyId `json:",omitempty"`
	DriverRef        string           `json:",omitempty"`
	id               VehicleId
	NextStopVisitId  StopVisitId `json:",omitempty"`
	LinkDistance     float64     `json:",omitempty"`
	Percentage       float64     `json:",omitempty"`
	Longitude        float64     `json:",omitempty"`
	Latitude         float64     `json:",omitempty"`
	Bearing          float64     `json:",omitempty"`
}

func NewVehicle(model Model) *Vehicle {
	vehicle := &Vehicle{
		model:         model,
		RawAttributes: NewRawAttributes(),
	}
	vehicle.InitCodes()
	return vehicle
}

func (vehicle *Vehicle) ModelId() string {
	return string(vehicle.id)
}

func (vehicle *Vehicle) copy() *Vehicle {
	return &Vehicle{
		CodeConsumer:     vehicle.CodeConsumer.Copy(),
		model:            vehicle.model,
		id:               vehicle.id,
		LineId:           vehicle.LineId,
		StopAreaId:       vehicle.StopAreaId,
		VehicleJourneyId: vehicle.VehicleJourneyId,
		Longitude:        vehicle.Longitude,
		Latitude:         vehicle.Latitude,
		Bearing:          vehicle.Bearing,
		LinkDistance:     vehicle.LinkDistance,
		Percentage:       vehicle.Percentage,
		DriverRef:        vehicle.DriverRef,
		ValidUntilTime:   vehicle.ValidUntilTime,
		Occupancy:        vehicle.Occupancy,
		RecordedAtTime:   vehicle.RecordedAtTime,
		RawAttributes:    vehicle.RawAttributes.Copy(),
		NextStopVisitId:  vehicle.NextStopVisitId,
	}
}

func (vehicle *Vehicle) Id() VehicleId {
	return vehicle.id
}

func (vehicle *Vehicle) GetLineId() LineId {
	return vehicle.LineId
}

func (vehicle *Vehicle) Save() bool {
	return vehicle.model.Vehicles().Save(vehicle)
}

func (vehicle *Vehicle) VehicleJourney() *VehicleJourney {
	vehicleJourney, ok := vehicle.model.VehicleJourneys().Find(vehicle.VehicleJourneyId)
	if !ok {
		return nil
	}
	return vehicleJourney
}

func (vehicle *Vehicle) MarshalJSON() ([]byte, error) {
	type Alias Vehicle
	aux := struct {
		Codes         Codes         `json:",omitempty"`
		RawAttributes RawAttributes `json:",omitempty"`
		*Alias
		Id VehicleId
	}{
		Id:    vehicle.id,
		Alias: (*Alias)(vehicle),
	}

	if !vehicle.Codes().Empty() {
		aux.Codes = vehicle.Codes()
	}
	if !vehicle.RawAttributes.IsEmpty() {
		aux.RawAttributes = vehicle.RawAttributes
	}

	return json.Marshal(&aux)
}

func (vehicle *Vehicle) UnmarshalJSON(data []byte) error {
	type Alias Vehicle
	aux := &struct {
		Codes map[string]string
		*Alias
	}{
		Alias: (*Alias)(vehicle),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Codes != nil {
		vehicle.SetCodesFromMap(aux.Codes)
	}

	return nil
}

type memoryVehicles struct {
	memoryManager
	clock.ClockConsumer
	IndexHandler

	mutex        *sync.RWMutex
	byIdentifier map[VehicleId]*Vehicle

	broadcastEvent func(event VehicleBroadcastEvent)
}

type Vehicles interface {
	ModelManager[VehicleId, *Vehicle]
	CodeHandler[*Vehicle]
	Broadcaster[VehicleBroadcastEvent]

	FindByLineId(LineId) []*Vehicle
	FindByVehicleJourneyId(VehicleJourneyId) (*Vehicle, bool)
	FindByNextStopVisitId(StopVisitId) (*Vehicle, bool)
}

func NewMemoryVehicles() Vehicles {
	v := &memoryVehicles{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[VehicleId]*Vehicle),
	}
	v.InitIndexes()
	v.AddIndex(ByLine, OneToMany, vehicleLineExtractor)
	v.AddIndex(ByVehicleJourney, OneToOne, vehicleVjExtractor)
	v.AddIndex(ByNextStopVisit, OneToOne, vehicleNextStopVisitExtractor)

	return v
}

func (manager *memoryVehicles) SetBroadcaster(f func(VehicleBroadcastEvent), _ ...string) {
	manager.broadcastEvent = f
}

func (manager *memoryVehicles) New() *Vehicle {
	return NewVehicle(manager.model)
}

func (manager *memoryVehicles) Find(id VehicleId) (*Vehicle, bool) {
	manager.mutex.RLock()
	vehicle, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return vehicle.copy(), true
	}
	return &Vehicle{}, false
}

func (manager *memoryVehicles) FindByCode(code Code) (*Vehicle, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.ByCode().Find(code)
	if ok {
		return manager.byIdentifier[VehicleId(id)].copy(), true
	}
	return &Vehicle{}, false
}

func (manager *memoryVehicles) CodeExists(code Code) bool {
	manager.mutex.RLock()
	_, ok := manager.ByCode().Find(code)
	manager.mutex.RUnlock()

	return ok
}

func (manager *memoryVehicles) FindByLineId(id LineId) (vehicles []*Vehicle) {
	manager.mutex.RLock()

	ids, _ := manager.FindBy(ByLine, string(id))

	for _, id := range ids {
		v := manager.byIdentifier[VehicleId(id)]
		vehicles = append(vehicles, v.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *memoryVehicles) FindByVehicleJourneyId(vjId VehicleJourneyId) (*Vehicle, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.FindOneBy(ByVehicleJourney, string(vjId))
	if ok {
		// The index is kept in sync by Index/Deindex; the guard is belt-and-
		// suspenders so a stale id yields a miss, not a nil-deref or wrong vehicle.
		if vehicle, found := manager.byIdentifier[VehicleId(id)]; found && vehicle.VehicleJourneyId == vjId {
			return vehicle.copy(), true
		}
	}
	return &Vehicle{}, false
}

func (manager *memoryVehicles) FindAll() (vehicles []*Vehicle) {
	manager.mutex.RLock()

	for _, vehicle := range manager.byIdentifier {
		vehicles = append(vehicles, vehicle.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *memoryVehicles) FindByNextStopVisitId(stopVisitId StopVisitId) (*Vehicle, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.FindOneBy(ByNextStopVisit, string(stopVisitId))
	if ok {
		// The index is kept in sync by Index/Deindex; the guard is belt-and-
		// suspenders so a stale id yields a miss, not a nil-deref or wrong vehicle.
		if vehicle, found := manager.byIdentifier[VehicleId(id)]; found && vehicle.NextStopVisitId == stopVisitId {
			return vehicle.copy(), true
		}
	}
	return &Vehicle{}, false
}

func (manager *memoryVehicles) Save(vehicle *Vehicle) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if vehicle.id == "" {
		vehicle.id = VehicleId(manager.NewUUID())
		manager.sendBQMessage(vehicle)
	} else if v, ok := manager.byIdentifier[vehicle.Id()]; ok {
		r, err := Equal(v, vehicle)
		if err != nil {
			logger.Log.Debugf("Error while comparing two vehicles: %v", err)
		} else if !r.Equal {
			manager.sendBQMessage(vehicle)
		}
	}

	vehicle.model = manager.model
	manager.byIdentifier[vehicle.Id()] = vehicle
	// Index maintains ByLine, ByVehicleJourney and ByNextStopVisit, including
	// evicting the previous key when a vehicle advances its next stop.
	manager.Index(vehicle)

	event := VehicleBroadcastEvent{
		ModelId:   string(vehicle.id),
		ModelType: "Vehicle",
	}

	if manager.broadcastEvent != nil {
		manager.broadcastEvent(event)
	}

	return true
}

func (manager *memoryVehicles) sendBQMessage(v *Vehicle) {
	if manager.model == nil {
		return
	}
	vehicleEvent := &audit.BigQueryVehicleEvent{
		Timestamp:      manager.Clock().Now(),
		ID:             string(v.id),
		Codes:          v.CodeSlice(),
		Longitude:      v.Longitude,
		Latitude:       v.Latitude,
		Bearing:        v.Bearing,
		Occupancy:      v.Occupancy,
		RecordedAtTime: civil.DateTimeOf(v.RecordedAtTime),
	}

	audit.CurrentBigQuery(manager.model.Referential()).WriteEvent(vehicleEvent)
}

func (manager *memoryVehicles) Delete(vehicle *Vehicle) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	delete(manager.byIdentifier, vehicle.Id())
	// Deindex removes the vehicle from ByLine, ByVehicleJourney and ByNextStopVisit.
	manager.Deindex(string(vehicle.id))

	return true
}
