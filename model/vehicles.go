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
	RecordedAtTime time.Time
	ValidUntilTime time.Time `json:",omitempty"`
	model          Model
	CodeConsumer
	Attributes       Attributes
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
		model:      model,
		Attributes: NewAttributes(),
	}
	vehicle.codes = make(Codes)
	return vehicle
}

func (vehicle *Vehicle) modelId() ModelId {
	return ModelId(vehicle.id)
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
		Attributes:       vehicle.Attributes.Copy(),
		NextStopVisitId:  vehicle.NextStopVisitId,
	}
}

func (vehicle *Vehicle) Id() VehicleId {
	return vehicle.id
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
		Codes      Codes      `json:",omitempty"`
		Attributes Attributes `json:",omitempty"`
		*Alias
		Id VehicleId
	}{
		Id:    vehicle.id,
		Alias: (*Alias)(vehicle),
	}

	if !vehicle.Codes().Empty() {
		aux.Codes = vehicle.Codes()
	}
	if !vehicle.Attributes.IsEmpty() {
		aux.Attributes = vehicle.Attributes
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
		vehicle.CodeConsumer.codes = NewCodesFromMap(aux.Codes)
	}

	return nil
}

type MemoryVehicles struct {
	uuid.UUIDConsumer
	clock.ClockConsumer

	model *MemoryModel

	mutex             *sync.RWMutex
	byIdentifier      map[VehicleId]*Vehicle
	byCode            *CodeIndex
	byNextStopVisitId map[StopVisitId]VehicleId

	broadcastEvent func(event VehicleBroadcastEvent)
}

type Vehicles interface {
	uuid.UUIDInterface

	New() *Vehicle
	Find(VehicleId) (*Vehicle, bool)
	FindByCode(Code) (*Vehicle, bool)
	FindByLineId(LineId) []*Vehicle
	FindByNextStopVisitId(StopVisitId) (*Vehicle, bool)
	FindAll() []*Vehicle
	Save(*Vehicle) bool
	Delete(*Vehicle) bool
}

func NewMemoryVehicles() *MemoryVehicles {
	return &MemoryVehicles{
		mutex:             &sync.RWMutex{},
		byIdentifier:      make(map[VehicleId]*Vehicle),
		byCode:            NewCodeIndex(),
		byNextStopVisitId: make(map[StopVisitId]VehicleId),
	}
}

func (manager *MemoryVehicles) New() *Vehicle {
	return NewVehicle(manager.model)
}

func (manager *MemoryVehicles) Find(id VehicleId) (*Vehicle, bool) {
	manager.mutex.RLock()
	vehicle, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return vehicle.copy(), true
	}
	return &Vehicle{}, false
}

func (manager *MemoryVehicles) FindByCode(code Code) (*Vehicle, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byCode.Find(code)
	if ok {
		return manager.byIdentifier[VehicleId(id)].copy(), true
	}
	return &Vehicle{}, false
}

func (manager *MemoryVehicles) FindByLineId(id LineId) (vehicles []*Vehicle) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, vehicle := range manager.byIdentifier {
		if vehicle.LineId == id {
			vehicles = append(vehicles, vehicle.copy())
		}
	}
	return
}

func (manager *MemoryVehicles) FindAll() (vehicles []*Vehicle) {
	manager.mutex.RLock()

	for _, vehicle := range manager.byIdentifier {
		vehicles = append(vehicles, vehicle.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryVehicles) FindByNextStopVisitId(stopVisitId StopVisitId) (*Vehicle, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	vehicleId, ok := manager.byNextStopVisitId[stopVisitId]
	if ok {
		vehicle, ok := manager.byIdentifier[vehicleId]
		if ok {
			if vehicle.NextStopVisitId == stopVisitId {
				return vehicle.copy(), true
			}
		}
		// clean the index
		manager.mutex.RUnlock()
		manager.mutex.Lock()
		delete(manager.byNextStopVisitId, stopVisitId)
		manager.mutex.Unlock()
		manager.mutex.RLock()
	}
	return &Vehicle{}, false
}

func (manager *MemoryVehicles) Save(vehicle *Vehicle) bool {
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
	manager.byCode.Index(vehicle)

	if vehicle.NextStopVisitId != StopVisitId("") {
		manager.byNextStopVisitId[vehicle.NextStopVisitId] = vehicle.Id()
	}

	event := VehicleBroadcastEvent{
		ModelId:   string(vehicle.id),
		ModelType: "Vehicle",
	}

	if manager.broadcastEvent != nil {
		manager.broadcastEvent(event)
	}

	return true
}

func (manager *MemoryVehicles) sendBQMessage(v *Vehicle) {
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

func (manager *MemoryVehicles) Delete(vehicle *Vehicle) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	delete(manager.byIdentifier, vehicle.Id())
	manager.byCode.Delete(ModelId(vehicle.id))

	return true
}
