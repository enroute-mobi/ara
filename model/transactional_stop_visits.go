package model

type TransactionalStopVisits struct {
	UUIDConsumer

	model           Model
	saved           map[StopVisitId]*StopVisit
	savedByObjectId map[string]map[string]StopVisitId
	deleted         map[StopVisitId]*StopVisit
}

func NewTransactionalStopVisits(model Model) *TransactionalStopVisits {
	stopVisits := TransactionalStopVisits{model: model}
	stopVisits.resetCaches()
	return &stopVisits
}

func (manager *TransactionalStopVisits) resetCaches() {
	manager.saved = make(map[StopVisitId]*StopVisit)
	manager.savedByObjectId = make(map[string]map[string]StopVisitId)
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
	for _, stopAera := range manager.deleted {
		manager.model.StopVisits().Delete(stopAera)
	}
	for _, stopAera := range manager.saved {
		manager.model.StopVisits().Save(stopAera)
	}
	return nil
}

func (manager *TransactionalStopVisits) Rollback() error {
	manager.resetCaches()
	return nil
}
