package model

import (
	"encoding/json"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"cloud.google.com/go/civil"
)

type VehicleId ModelId

type Vehicle struct {
	ObjectIDConsumer

	model Model

	id               VehicleId
	LineId           LineId           `json:",omitempty"`
	VehicleJourneyId VehicleJourneyId `json:",omitempty"`

	Longitude float64 `json:",omitempty"`
	Latitude  float64 `json:",omitempty"`
	Bearing   float64 `json:",omitempty"`

	RecordedAtTime time.Time
}

func NewVehicle(model Model) *Vehicle {
	vehicle := &Vehicle{
		model: model,
	}
	vehicle.objectids = make(ObjectIDs)
	return vehicle
}

func (vehicle *Vehicle) modelId() ModelId {
	return ModelId(vehicle.id)
}

func (vehicle *Vehicle) Id() VehicleId {
	return vehicle.id
}

func (vehicle *Vehicle) Save() (ok bool) {
	ok = vehicle.model.Vehicles().Save(vehicle)
	return
}

func (vehicle *Vehicle) VehicleJourney() *VehicleJourney {
	vehicleJourney, ok := vehicle.model.VehicleJourneys().Find(vehicle.VehicleJourneyId)
	if !ok {
		return nil
	}
	return &vehicleJourney
}

func (vehicle *Vehicle) MarshalJSON() ([]byte, error) {
	type Alias Vehicle
	aux := struct {
		Id        VehicleId
		ObjectIDs ObjectIDs `json:",omitempty"`
		*Alias
	}{
		Id:    vehicle.id,
		Alias: (*Alias)(vehicle),
	}

	if !vehicle.ObjectIDs().Empty() {
		aux.ObjectIDs = vehicle.ObjectIDs()
	}

	return json.Marshal(&aux)
}

func (vehicle *Vehicle) UnmarshalJSON(data []byte) error {
	type Alias Vehicle
	aux := &struct {
		ObjectIDs map[string]string
		*Alias
	}{
		Alias: (*Alias)(vehicle),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ObjectIDs != nil {
		vehicle.ObjectIDConsumer.objectids = NewObjectIDsFromMap(aux.ObjectIDs)
	}

	return nil
}

type MemoryVehicles struct {
	uuid.UUIDConsumer
	clock.ClockConsumer

	model *MemoryModel

	mutex        *sync.RWMutex
	byIdentifier map[VehicleId]*Vehicle
	byObjectId   *ObjectIdIndex
}

type Vehicles interface {
	uuid.UUIDInterface

	New() Vehicle
	Find(id VehicleId) (Vehicle, bool)
	FindByObjectId(objectid ObjectID) (Vehicle, bool)
	FindByLineId(id LineId) []Vehicle
	FindAll() []Vehicle
	Save(vehicle *Vehicle) bool
	Delete(vehicle *Vehicle) bool
}

func NewMemoryVehicles() *MemoryVehicles {
	return &MemoryVehicles{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[VehicleId]*Vehicle),
		byObjectId:   NewObjectIdIndex(),
	}
}

func (manager *MemoryVehicles) New() Vehicle {
	vehicle := NewVehicle(manager.model)
	return *vehicle
}

func (manager *MemoryVehicles) Find(id VehicleId) (Vehicle, bool) {
	if id == "" {
		return Vehicle{}, false
	}

	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	vehicle, ok := manager.byIdentifier[id]
	if ok {
		return *vehicle, true
	} else {
		return Vehicle{}, false
	}
}

func (manager *MemoryVehicles) FindByObjectId(objectid ObjectID) (Vehicle, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byObjectId.Find(objectid)
	if ok {
		return *manager.byIdentifier[VehicleId(id)], true
	}
	return Vehicle{}, false
}

func (manager *MemoryVehicles) FindByLineId(id LineId) (vehicles []Vehicle) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, vehicle := range manager.byIdentifier {
		if vehicle.LineId == id {
			vehicles = append(vehicles, *vehicle)
		}
	}
	return
}

func (manager *MemoryVehicles) FindAll() (vehicles []Vehicle) {
	manager.mutex.RLock()

	if len(manager.byIdentifier) == 0 {
		manager.mutex.RUnlock()
		return
	}
	for _, vehicle := range manager.byIdentifier {
		vehicles = append(vehicles, *vehicle)
	}
	manager.mutex.RUnlock()

	return
}

func (manager *MemoryVehicles) Save(vehicle *Vehicle) bool {
	manager.mutex.Lock()

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
	manager.byObjectId.Index(vehicle)

	manager.mutex.Unlock()

	return true
}

func (manager *MemoryVehicles) sendBQMessage(v *Vehicle) {
	if manager.model == nil {
		return
	}
	vehicleEvent := &audit.BigQueryVehicleEvent{
		Timestamp:      manager.Clock().Now(),
		ID:             string(v.id),
		ObjectIDs:      v.ObjectIDSlice(),
		Longitude:      v.Longitude,
		Latitude:       v.Latitude,
		Bearing:        v.Bearing,
		RecordedAtTime: civil.DateTimeOf(v.RecordedAtTime),
	}
	audit.CurrentBigQuery(manager.model.Referential()).WriteVehicleEvent(vehicleEvent)
}

func (manager *MemoryVehicles) Delete(vehicle *Vehicle) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, vehicle.Id())
	manager.byObjectId.Delete(ModelId(vehicle.id))

	return true
}
