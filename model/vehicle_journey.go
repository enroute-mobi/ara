package model

import "encoding/json"

type VehicleJourneyId string

type VehicleJourneyAttributes struct {
	ObjectId     *ObjectID
	LineObjectId *ObjectID
}

type VehicleJourney struct {
	ObjectIDConsumer
	model Model

	id     VehicleJourneyId
	lineId LineId
}

func NewVehicleJourney(model Model) *VehicleJourney {
	vehicleJourney := &VehicleJourney{model: model}
	vehicleJourney.objectids = make(ObjectIDs)
	return vehicleJourney
}

func (vehicleJourney *VehicleJourney) Id() VehicleJourneyId {
	return vehicleJourney.id
}

func (vehicleJourney *VehicleJourney) Line() Line {
	line, _ := vehicleJourney.model.Lines().Find(vehicleJourney.lineId)
	return line
}

// WIP
func (vehicleJourney *VehicleJourney) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Id": vehicleJourney.id,
	})
}

func (vehicleJourney *VehicleJourney) Save() (ok bool) {
	ok = vehicleJourney.model.VehicleJourneys().Save(vehicleJourney)
	return
}

type MemoryVehicleJourneys struct {
	UUIDConsumer

	model Model

	byIdentifier map[VehicleJourneyId]*VehicleJourney
	byObjectId   map[string]map[string]VehicleJourneyId
}

type VehicleJourneys interface {
	UUIDInterface

	New() VehicleJourney
	Find(id VehicleJourneyId) (VehicleJourney, bool)
	FindByObjectId(objectid ObjectID) (VehicleJourney, bool)
	FindAll() []VehicleJourney
	Save(vehicleJourney *VehicleJourney) bool
	Delete(vehicleJourney *VehicleJourney) bool
}

func NewMemoryVehicleJourneys() *MemoryVehicleJourneys {
	return &MemoryVehicleJourneys{
		byIdentifier: make(map[VehicleJourneyId]*VehicleJourney),
		byObjectId:   make(map[string]map[string]VehicleJourneyId),
	}
}

func (manager *MemoryVehicleJourneys) New() VehicleJourney {
	vehicleJourney := NewVehicleJourney(manager.model)
	return *vehicleJourney
}

func (manager *MemoryVehicleJourneys) Find(id VehicleJourneyId) (VehicleJourney, bool) {
	vehicleJourney, ok := manager.byIdentifier[id]
	if ok {
		return *vehicleJourney, true
	} else {
		return VehicleJourney{}, false
	}
}

func (manager *MemoryVehicleJourneys) FindByObjectId(objectid ObjectID) (VehicleJourney, bool) {
	valueMap, ok := manager.byObjectId[objectid.Kind()]
	if !ok {
		return VehicleJourney{}, false
	}
	id, ok := valueMap[objectid.Value()]
	if !ok {
		return VehicleJourney{}, false
	}
	return *manager.byIdentifier[id], true
}

func (manager *MemoryVehicleJourneys) FindAll() (vehicleJourneys []VehicleJourney) {
	if len(manager.byIdentifier) == 0 {
		return []VehicleJourney{}
	}
	for _, vehicleJourney := range manager.byIdentifier {
		vehicleJourneys = append(vehicleJourneys, *vehicleJourney)
	}
	return
}

func (manager *MemoryVehicleJourneys) Save(vehicleJourney *VehicleJourney) bool {
	if vehicleJourney.Id() == "" {
		vehicleJourney.id = VehicleJourneyId(manager.NewUUID())
	}
	vehicleJourney.model = manager.model
	manager.byIdentifier[vehicleJourney.Id()] = vehicleJourney
	for _, objectid := range vehicleJourney.ObjectIDs() {
		_, ok := manager.byObjectId[objectid.Kind()]
		if !ok {
			manager.byObjectId[objectid.Kind()] = make(map[string]VehicleJourneyId)
		}
		manager.byObjectId[objectid.Kind()][objectid.Value()] = vehicleJourney.Id()
	}
	return true
}

func (manager *MemoryVehicleJourneys) Delete(vehicleJourney *VehicleJourney) bool {
	delete(manager.byIdentifier, vehicleJourney.Id())
	for _, objectid := range vehicleJourney.ObjectIDs() {
		valueMap := manager.byObjectId[objectid.Kind()]
		delete(valueMap, objectid.Value())
	}
	return true
}
