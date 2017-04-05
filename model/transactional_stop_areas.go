package model

type TransactionalStopAreas struct {
	UUIDConsumer

	model   Model
	saved   map[StopAreaId]*StopArea
	deleted map[StopAreaId]*StopArea
}

func NewTransactionalStopAreas(model Model) *TransactionalStopAreas {
	stopAreas := TransactionalStopAreas{model: model}
	stopAreas.resetCaches()
	return &stopAreas
}

func (manager *TransactionalStopAreas) resetCaches() {
	manager.saved = make(map[StopAreaId]*StopArea)
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
	for _, stopArea := range manager.saved {
		stopAreaObjectId, _ := stopArea.ObjectID(objectid.Kind())
		if stopAreaObjectId.Value() == objectid.Value() {
			return *stopArea, true
		}
	}
	return manager.model.StopAreas().FindByObjectId(objectid)
}

func (manager *TransactionalStopAreas) FindAll() []StopArea {
	stopAreas := []StopArea{}
	for _, savedStopArea := range manager.saved {
		stopAreas = append(stopAreas, *savedStopArea)
	}
	modelStopAreas := manager.model.StopAreas().FindAll()
	for _, stopArea := range modelStopAreas {
		_, ok := manager.saved[stopArea.Id()]
		if !ok {
			stopAreas = append(stopAreas, stopArea)
		}
	}
	return stopAreas
}

func (manager *TransactionalStopAreas) Save(stopArea *StopArea) bool {
	if stopArea.Id() == "" {
		stopArea.id = StopAreaId(manager.NewUUID())
	}
	manager.saved[stopArea.Id()] = stopArea
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
