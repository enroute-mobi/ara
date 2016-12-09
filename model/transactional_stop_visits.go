package model

type TransactionalStopVisits struct {
	UUIDConsumer

	model   Model
	saved   map[StopVisitId]*StopVisit
	deleted map[StopVisitId]*StopVisit
}

func NewTransactionalStopVisits(model Model) *TransactionalStopVisits {
	stopVisits := TransactionalStopVisits{model: model}
	stopVisits.resetCaches()
	return &stopVisits
}

func (manager *TransactionalStopVisits) resetCaches() {
	manager.saved = make(map[StopVisitId]*StopVisit)
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
