package model

type TransactionalStopVisits struct {
	UUIDConsumer

	model                   Model
	saved                   map[StopVisitId]*StopVisit
	savedByObjectId         map[string]map[string]StopVisitId
	savedByVehicleJourneyId map[VehicleJourneyId][]StopVisitId
	deleted                 map[StopVisitId]*StopVisit
}

func NewTransactionalStopVisits(model Model) *TransactionalStopVisits {
	stopVisits := TransactionalStopVisits{model: model}
	stopVisits.resetCaches()
	return &stopVisits
}

func (manager *TransactionalStopVisits) resetCaches() {
	manager.saved = make(map[StopVisitId]*StopVisit)
	manager.savedByObjectId = make(map[string]map[string]StopVisitId)
	manager.savedByVehicleJourneyId = make(map[VehicleJourneyId][]StopVisitId)
	manager.deleted = make(map[StopVisitId]*StopVisit)
}

func (manager *TransactionalStopVisits) New() StopVisit {
	return *NewStopVisit(manager.model)
}

func (manager *TransactionalStopVisits) Find(id StopVisitId) (StopVisit, bool) {
	stopVisit, ok := manager.saved[id]
	if ok {
		return *stopVisit, ok
	}

	return manager.model.StopVisits().Find(id)
}

func (manager *TransactionalStopVisits) FindByObjectId(objectid ObjectID) (StopVisit, bool) {
	valueMap, ok := manager.savedByObjectId[objectid.Kind()]
	if !ok {
		return manager.model.StopVisits().FindByObjectId(objectid)
	}
	id, ok := valueMap[objectid.Value()]
	if !ok {
		return manager.model.StopVisits().FindByObjectId(objectid)
	}
	return *manager.saved[id], true
}

func (manager *TransactionalStopVisits) FindByVehicleJourneyId(id VehicleJourneyId) (stopVisits []StopVisit) {
	// Check saved StopVisits
	stopVisitIds, ok := manager.savedByVehicleJourneyId[id]
	if ok {
		for _, stopVisitId := range stopVisitIds {
			stopVisits = append(stopVisits, *manager.saved[stopVisitId])
		}
	}
	// Check model StopVisits
	for _, modelStopVisit := range manager.model.StopVisits().FindByVehicleJourneyId(id) {
		_, ok := manager.saved[modelStopVisit.Id()]
		if !ok {
			stopVisits = append(stopVisits, modelStopVisit)
		}
	}
	return
}

// Temp
func (manager *TransactionalStopVisits) FindByStopAreaId(id StopAreaId) (stopVisits []StopVisit) {
	// Check saved StopVisits
	for _, stopVisit := range manager.saved {
		if stopVisit.stopAreaId == id {
			stopVisits = append(stopVisits, *stopVisit)
		}
	}
	// Check model StopVisits
	for _, modelStopVisit := range manager.model.StopVisits().FindByStopAreaId(id) {
		_, ok := manager.saved[modelStopVisit.Id()]
		if !ok {
			stopVisits = append(stopVisits, modelStopVisit)
		}
	}
	return
}

func (manager *TransactionalStopVisits) FindAll() (stopVisits []StopVisit) {
	for _, stopVisit := range manager.saved {
		stopVisits = append(stopVisits, *stopVisit)
	}
	savedStopVisits := manager.model.StopVisits().FindAll()
	for _, stopVisit := range savedStopVisits {
		_, ok := manager.saved[stopVisit.Id()]
		if !ok {
			stopVisits = append(stopVisits, stopVisit)
		}
	}
	return
}

func (manager *TransactionalStopVisits) Save(stopVisit *StopVisit) bool {
	if stopVisit.Id() == "" {
		stopVisit.id = StopVisitId(manager.NewUUID())
	}
	manager.saved[stopVisit.Id()] = stopVisit
	manager.savedByVehicleJourneyId[stopVisit.vehicleJourneyId] = append(manager.savedByVehicleJourneyId[stopVisit.vehicleJourneyId], stopVisit.id)
	for _, objectid := range stopVisit.ObjectIDs() {
		_, ok := manager.savedByObjectId[objectid.Kind()]
		if !ok {
			manager.savedByObjectId[objectid.Kind()] = make(map[string]StopVisitId)
		}
		manager.savedByObjectId[objectid.Kind()][objectid.Value()] = stopVisit.Id()
	}
	return true
}

func (manager *TransactionalStopVisits) Delete(stopVisit *StopVisit) bool {
	manager.deleted[stopVisit.Id()] = stopVisit
	return true
}

func (manager *TransactionalStopVisits) Commit() error {
	for _, stopVisit := range manager.deleted {
		manager.model.StopVisits().Delete(stopVisit)
	}
	for _, stopVisit := range manager.saved {
		manager.model.StopVisits().Save(stopVisit)
	}
	return nil
}

func (manager *TransactionalStopVisits) Rollback() error {
	manager.resetCaches()
	return nil
}
