package model

type TransactionalStopAreas struct {
	UUIDConsumer

	model           Model
	saved           map[StopAreaId]*StopArea
	savedByObjectId map[string]map[string]StopAreaId
	deleted         map[StopAreaId]*StopArea
}

func NewTransactionalStopAreas(model Model) *TransactionalStopAreas {
	stopAreas := TransactionalStopAreas{model: model}
	stopAreas.resetCaches()
	return &stopAreas
}

func (manager *TransactionalStopAreas) resetCaches() {
	manager.saved = make(map[StopAreaId]*StopArea)
	manager.savedByObjectId = make(map[string]map[string]StopAreaId)
	manager.deleted = make(map[StopAreaId]*StopArea)
}

func (manager *TransactionalStopAreas) New() StopArea {
	return *NewStopArea(manager.model)
}

func (manager *TransactionalStopAreas) Find(id StopAreaId) (StopArea, bool) {
	stopArea, ok := manager.saved[id]
	if ok {
		return *stopArea, ok
	}

	return manager.model.StopAreas().Find(id)
}

func (manager *TransactionalStopAreas) FindByObjectId(objectid ObjectID) (StopArea, bool) {
	valueMap, ok := manager.savedByObjectId[objectid.Kind()]
	if !ok {
		return manager.model.StopAreas().FindByObjectId(objectid)
	}
	id, ok := valueMap[objectid.Value()]
	if !ok {
		return manager.model.StopAreas().FindByObjectId(objectid)
	}
	return *manager.saved[id], true
}

func (manager *TransactionalStopAreas) FindAll() (stopAreas []StopArea) {
	for _, stopArea := range manager.saved {
		stopAreas = append(stopAreas, *stopArea)
	}
	savedStopAreas := manager.model.StopAreas().FindAll()
	for _, stopArea := range savedStopAreas {
		_, ok := manager.saved[stopArea.Id()]
		if !ok {
			stopAreas = append(stopAreas, stopArea)
		}
	}
	return
}

func (manager *TransactionalStopAreas) Save(stopArea *StopArea) bool {
	if stopArea.Id() == "" {
		stopArea.id = StopAreaId(manager.NewUUID())
	}
	manager.saved[stopArea.Id()] = stopArea
	for _, objectid := range stopArea.ObjectIDs() {
		_, ok := manager.savedByObjectId[objectid.Kind()]
		if !ok {
			manager.savedByObjectId[objectid.Kind()] = make(map[string]StopAreaId)
		}
		manager.savedByObjectId[objectid.Kind()][objectid.Value()] = stopArea.Id()
	}
	return true
}

func (manager *TransactionalStopAreas) Delete(stopArea *StopArea) bool {
	manager.deleted[stopArea.Id()] = stopArea
	return true
}

// WIP: Handle errors
func (manager *TransactionalStopAreas) Commit() error {
	for _, stopAera := range manager.deleted {
		manager.model.StopAreas().Delete(stopAera)
	}
	for _, stopAera := range manager.saved {
		manager.model.StopAreas().Save(stopAera)
	}
	return nil
}

func (manager *TransactionalStopAreas) Rollback() error {
	manager.resetCaches()
	return nil
}
