package model

type TransactionalVehicles struct {
	UUIDConsumer

	model   Model
	saved   map[VehicleId]*Vehicle
	deleted map[VehicleId]*Vehicle
}

func NewTransactionalVehicles(model Model) *TransactionalVehicles {
	vehicles := TransactionalVehicles{model: model}
	vehicles.resetCaches()
	return &vehicles
}

func (manager *TransactionalVehicles) resetCaches() {
	manager.saved = make(map[VehicleId]*Vehicle)
	manager.deleted = make(map[VehicleId]*Vehicle)
}

func (manager *TransactionalVehicles) New() Vehicle {
	return *NewVehicle(manager.model)
}

func (manager *TransactionalVehicles) Find(id VehicleId) (Vehicle, bool) {
	vehicle, ok := manager.saved[id]
	if ok {
		return *vehicle, ok
	}

	return manager.model.Vehicles().Find(id)
}

func (manager *TransactionalVehicles) FindByObjectId(objectid ObjectID) (Vehicle, bool) {
	for _, vehicle := range manager.saved {
		vehicleObjectId, _ := vehicle.ObjectID(objectid.Kind())
		if vehicleObjectId.Value() == objectid.Value() {
			return *vehicle, true
		}
	}
	return manager.model.Vehicles().FindByObjectId(objectid)
}

func (manager *TransactionalVehicles) FindByLineId(id LineId) (vehicles []Vehicle) {
	// Check saved Vehicles
	for _, vehicle := range manager.saved {
		if vehicle.lineId == id {
			vehicles = append(vehicles, *vehicle)
		}
	}

	// Check model Vehicles
	for _, modelVehicle := range manager.model.Vehicles().FindByLineId(id) {
		_, ok := manager.saved[modelVehicle.Id()]
		if !ok {
			vehicles = append(vehicles, modelVehicle)
		}
	}
	return
}

func (manager *TransactionalVehicles) FindAll() []Vehicle {
	vehicles := []Vehicle{}
	for _, savedVehicle := range manager.saved {
		vehicles = append(vehicles, *savedVehicle)
	}
	modelVehicles := manager.model.Vehicles().FindAll()
	for _, vehicle := range modelVehicles {
		_, ok := manager.saved[vehicle.Id()]
		if !ok {
			vehicles = append(vehicles, vehicle)
		}
	}
	return vehicles
}

func (manager *TransactionalVehicles) Save(vehicle *Vehicle) bool {
	if vehicle.Id() == "" {
		vehicle.id = VehicleId(manager.NewUUID())
	}
	manager.saved[vehicle.Id()] = vehicle
	return true
}

func (manager *TransactionalVehicles) Delete(vehicle *Vehicle) bool {
	manager.deleted[vehicle.Id()] = vehicle
	return true
}

// WIP: Handle errors
func (manager *TransactionalVehicles) Commit() error {
	for _, stopAera := range manager.deleted {
		manager.model.Vehicles().Delete(stopAera)
	}
	for _, stopAera := range manager.saved {
		manager.model.Vehicles().Save(stopAera)
	}
	return nil
}

func (manager *TransactionalVehicles) Rollback() error {
	manager.resetCaches()
	return nil
}
