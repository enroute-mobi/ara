package model

type TransactionalVehicleJourneys struct {
	UUIDConsumer

	model   Model
	saved   map[VehicleJourneyId]*VehicleJourney
	deleted map[VehicleJourneyId]*VehicleJourney
}

func NewTransactionalVehicleJourneys(model Model) *TransactionalVehicleJourneys {
	vehicleJourneys := TransactionalVehicleJourneys{model: model}
	vehicleJourneys.resetCaches()
	return &vehicleJourneys
}

func (manager *TransactionalVehicleJourneys) resetCaches() {
	manager.saved = make(map[VehicleJourneyId]*VehicleJourney)
	manager.deleted = make(map[VehicleJourneyId]*VehicleJourney)
}

func (manager *TransactionalVehicleJourneys) New() VehicleJourney {
	return *NewVehicleJourney(manager.model)
}

func (manager *TransactionalVehicleJourneys) Find(id VehicleJourneyId) (VehicleJourney, bool) {
	vehicleJourney, ok := manager.saved[id]
	if ok {
		return *vehicleJourney, ok
	}

	return manager.model.VehicleJourneys().Find(id)
}

func (manager *TransactionalVehicleJourneys) FindByObjectId(objectid ObjectID) (VehicleJourney, bool) {
	for _, vehicleJourney := range manager.saved {
		vehicleJourneyObjectId, _ := vehicleJourney.ObjectID(objectid.Kind())
		if vehicleJourneyObjectId.Value() == objectid.Value() {
			return *vehicleJourney, true
		}
	}
	return manager.model.VehicleJourneys().FindByObjectId(objectid)
}

func (manager *TransactionalVehicleJourneys) FindAll() []VehicleJourney {
	vehicleJourneys := []VehicleJourney{}
	for _, vehicleJourney := range manager.saved {
		vehicleJourneys = append(vehicleJourneys, *vehicleJourney)
	}
	savedVehicleJourneys := manager.model.VehicleJourneys().FindAll()
	for _, vehicleJourney := range savedVehicleJourneys {
		_, ok := manager.saved[vehicleJourney.Id()]
		if !ok {
			vehicleJourneys = append(vehicleJourneys, vehicleJourney)
		}
	}
	return vehicleJourneys
}

func (manager *TransactionalVehicleJourneys) Save(vehicleJourney *VehicleJourney) bool {
	if vehicleJourney.Id() == "" {
		vehicleJourney.id = VehicleJourneyId(manager.NewUUID())
	}
	manager.saved[vehicleJourney.Id()] = vehicleJourney
	return true
}

func (manager *TransactionalVehicleJourneys) Delete(vehicleJourney *VehicleJourney) bool {
	manager.deleted[vehicleJourney.Id()] = vehicleJourney
	return true
}

func (manager *TransactionalVehicleJourneys) Commit() error {
	for _, vehicleJourney := range manager.deleted {
		manager.model.VehicleJourneys().Delete(vehicleJourney)
	}
	for _, vehicleJourney := range manager.saved {
		manager.model.VehicleJourneys().Save(vehicleJourney)
	}
	return nil
}

func (manager *TransactionalVehicleJourneys) Rollback() error {
	manager.resetCaches()
	return nil
}
