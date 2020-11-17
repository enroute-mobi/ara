package model

import (
	"sort"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type TransactionalStopVisits struct {
	uuid.UUIDConsumer
	clock.ClockConsumer

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
		return *(stopVisit.copy()), ok
	}

	return manager.model.StopVisits().Find(id)
}

func (manager *TransactionalStopVisits) FindByObjectId(objectid ObjectID) (StopVisit, bool) {
	for _, stopVisit := range manager.saved {
		stopVisitObjectId, _ := stopVisit.ObjectID(objectid.Kind())
		if stopVisitObjectId.Value() == objectid.Value() {
			return *(stopVisit.copy()), true
		}
	}
	return manager.model.StopVisits().FindByObjectId(objectid)
}

func (manager *TransactionalStopVisits) FindByVehicleJourneyId(id VehicleJourneyId) (stopVisits []StopVisit) {
	// Check saved StopVisits
	for _, stopVisit := range manager.saved {
		if stopVisit.VehicleJourneyId == id {
			stopVisits = append(stopVisits, *(stopVisit.copy()))
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

func (manager *TransactionalStopVisits) FindFollowingByVehicleJourneyId(id VehicleJourneyId) (stopVisits []StopVisit) {
	for _, stopVisit := range manager.saved {
		if stopVisit.VehicleJourneyId == id && stopVisit.ReferenceTime().After(manager.Clock().Now()) {
			stopVisits = append(stopVisits, *(stopVisit.copy()))
		}
	}
	for _, modelStopVisit := range manager.model.StopVisits().FindFollowingByVehicleJourneyId(id) {
		_, saved := manager.saved[modelStopVisit.Id()]
		if !saved {
			stopVisits = append(stopVisits, modelStopVisit)
		}
	}
	sort.Sort(ByTime(stopVisits))
	return
}

// Temp
func (manager *TransactionalStopVisits) FindByStopAreaId(id StopAreaId) (stopVisits []StopVisit) {
	// Check saved StopVisits
	for _, stopVisit := range manager.saved {
		if stopVisit.StopAreaId == id {
			stopVisits = append(stopVisits, *(stopVisit.copy()))
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

func (manager *TransactionalStopVisits) FindFollowingByStopAreaId(id StopAreaId) (stopVisits []StopVisit) {
	for _, stopVisit := range manager.saved {
		if stopVisit.StopAreaId == id && stopVisit.ReferenceTime().After(manager.Clock().Now()) {
			stopVisits = append(stopVisits, *(stopVisit.copy()))
		}
	}
	for _, modelStopVisit := range manager.model.StopVisits().FindFollowingByStopAreaId(id) {
		_, saved := manager.saved[modelStopVisit.Id()]
		if !saved {
			stopVisits = append(stopVisits, modelStopVisit)
		}
	}
	sort.Sort(ByTime(stopVisits))
	return
}

func (manager *TransactionalStopVisits) FindFollowingByStopAreaIds(stopAreaIds []StopAreaId) (stopVisits []StopVisit) {
	for _, stopAreaId := range stopAreaIds {
		stopVisits = append(stopVisits, manager.FindFollowingByStopAreaId(stopAreaId)...)
	}
	sort.Sort(ByTime(stopVisits))
	return
}

func (manager *TransactionalStopVisits) FindAllAfter(t time.Time) (stopVisits []StopVisit) {
	for _, stopVisit := range manager.saved {
		if stopVisit.ReferenceTime().After(t) {
			stopVisits = append(stopVisits, *(stopVisit.copy()))
		}
	}
	savedStopVisits := manager.model.StopVisits().FindAll()
	for _, stopVisit := range savedStopVisits {
		_, ok := manager.saved[stopVisit.Id()]
		_, deleted := manager.deleted[stopVisit.Id()]
		if !ok && !deleted && stopVisit.ReferenceTime().After(t) {
			stopVisits = append(stopVisits, stopVisit)
		}
	}
	return
}

func (manager *TransactionalStopVisits) FindAll() (stopVisits []StopVisit) {
	for _, stopVisit := range manager.saved {
		stopVisits = append(stopVisits, *(stopVisit.copy()))
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
