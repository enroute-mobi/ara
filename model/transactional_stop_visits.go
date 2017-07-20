package model

import "sort"

type TransactionalStopVisits struct {
	UUIDConsumer
	ClockConsumer

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
	// stopVisit, ok := manager.saved[id]
	// if ok {
	// 	return *stopVisit, ok
	// }
	// stopVisit, ok = manager.deleted[id]
	// if ok {
	// 	return StopVisit{}, false
	// }

	// return manager.model.StopVisits().Find(id)
	stopVisit, ok := manager.collection()[id]
	if ok {
		return *stopVisit, true
	}
	return StopVisit{}, false
}

func (manager *TransactionalStopVisits) FindByObjectId(objectid ObjectID) (StopVisit, bool) {
	for _, stopVisit := range manager.collection() {
		stopVisitObjectId, _ := stopVisit.ObjectID(objectid.Kind())
		if stopVisitObjectId.Value() == objectid.Value() {
			return *stopVisit, true
		}
	}
	return StopVisit{}, false
}

func (manager *TransactionalStopVisits) collection() map[StopVisitId]*StopVisit {
	stopVisits := manager.model.StopVisits().collection()

	// Check Saved
	for id, stopVisit := range manager.saved {
		stopVisits[id] = stopVisit
	}

	// Check Deleted
	for id, _ := range manager.deleted {
		delete(stopVisits, id)
	}

	return stopVisits
}

func (manager *TransactionalStopVisits) FindByVehicleJourneyId(id VehicleJourneyId) (stopVisits []StopVisit) {
	// // Check saved StopVisits
	// for _, stopVisit := range manager.saved {
	// 	if stopVisit.VehicleJourneyId == id {
	// 		stopVisits = append(stopVisits, *stopVisit)
	// 	}
	// }

	// // Check model StopVisits
	// for _, modelStopVisit := range manager.model.StopVisits().FindByVehicleJourneyId(id) {
	// 	_, ok := manager.saved[modelStopVisit.Id()]
	// 	if !ok {
	// 		stopVisits = append(stopVisits, modelStopVisit)
	// 	}
	// }
	// return
	return FindStopVisitBy(manager.collection(), StopVisitSelectorByVehicleJourneyId(id))
}

// Temp
func (manager *TransactionalStopVisits) FindByStopAreaId(id StopAreaId) (stopVisits []StopVisit) {
	// // Check saved StopVisits
	// for _, stopVisit := range manager.saved {
	// 	if stopVisit.StopAreaId == id {
	// 		stopVisits = append(stopVisits, *stopVisit)
	// 	}
	// }

	// // Check model StopVisits
	// for _, modelStopVisit := range manager.model.StopVisits().FindByStopAreaId(id) {
	// 	_, ok := manager.saved[modelStopVisit.Id()]
	// 	if !ok {
	// 		stopVisits = append(stopVisits, modelStopVisit)
	// 	}
	// }
	// return
	return FindStopVisitBy(manager.collection(), StopVisitSelectorByStopAreaId(id))
}

func (manager *TransactionalStopVisits) FindFollowingByStopAreaId(id StopAreaId) (stopVisits []StopVisit) {
	// for _, stopVisit := range manager.saved {
	// 	if stopVisit.StopAreaId == id && stopVisit.ReferenceTime().After(manager.Clock().Now()) {
	// 		stopVisits = append(stopVisits, *stopVisit)
	// 	}
	// }
	// for _, modelStopVisit := range manager.model.StopVisits().FindFollowingByStopAreaId(id) {
	// 	_, saved := manager.saved[modelStopVisit.Id()]
	// 	if !saved {
	// 		stopVisits = append(stopVisits, modelStopVisit)
	// 	}
	// }
	// sort.Sort(ByTime(stopVisits))
	// return
	stopVisits = FindStopVisitBy(manager.collection(), StopVisitSelectorByStopAreaId(id), StopVisitSelectorFollowing(manager.Clock().Now()))
	sort.Sort(ByTime(stopVisits))
	return
}

func (manager *TransactionalStopVisits) FindFollowingByStopAreaIds(stopAreaIds []StopAreaId) (stopVisits []StopVisit) {
	// for _, stopAreaId := range stopAreaIds {
	// 	stopVisits = append(stopVisits, manager.FindFollowingByStopAreaId(stopAreaId)...)
	// }
	// sort.Sort(ByTime(stopVisits))
	// return
	stopVisits = FindStopVisitBy(manager.collection(), StopVisitSelectorByStopAreaIds(stopAreaIds), StopVisitSelectorFollowing(manager.Clock().Now()))
	sort.Sort(ByTime(stopVisits))
	return
}

func (manager *TransactionalStopVisits) FindAll() (stopVisits []StopVisit) {
	// stopVisits := []StopVisit{}
	// for _, stopVisit := range manager.saved {
	// 	stopVisits = append(stopVisits, *stopVisit)
	// }
	// savedStopVisits := manager.model.StopVisits().FindAll()
	// for _, stopVisit := range savedStopVisits {
	// 	_, ok := manager.saved[stopVisit.Id()]
	// 	if !ok {
	// 		stopVisits = append(stopVisits, stopVisit)
	// 	}
	// }
	// return stopVisits
	collection := manager.collection()
	if len(collection) == 0 {
		return []StopVisit{}
	}
	for _, stopVisit := range collection {
		stopVisits = append(stopVisits, *stopVisit)
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
