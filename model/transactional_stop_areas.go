package model

import "bitbucket.org/enroute-mobi/ara/uuid"

type TransactionalStopAreas struct {
	uuid.UUIDConsumer

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
		return *(stopArea.copy()), ok
	}

	return manager.model.StopAreas().Find(id)
}

func (manager *TransactionalStopAreas) FindByObjectId(objectid ObjectID) (StopArea, bool) {
	for _, stopArea := range manager.saved {
		stopAreaObjectId, _ := stopArea.ObjectID(objectid.Kind())
		if stopAreaObjectId.Value() == objectid.Value() {
			return *(stopArea.copy()), true
		}
	}
	return manager.model.StopAreas().FindByObjectId(objectid)
}

// Temp
func (manager *TransactionalStopAreas) FindByLineId(id LineId) (stopAreas []StopArea) {
	// Check saved StopAreas
	for _, stopArea := range manager.saved {
		if stopArea.LineIds.Contains(id) {
			stopAreas = append(stopAreas, *(stopArea.copy()))
		}
	}

	// Check model StopAreas
	for _, modelStopArea := range manager.model.StopAreas().FindByLineId(id) {
		_, ok := manager.saved[modelStopArea.Id()]
		if !ok {
			stopAreas = append(stopAreas, modelStopArea)
		}
	}
	return
}

func (manager *TransactionalStopAreas) FindAll() []StopArea {
	stopAreas := []StopArea{}
	for _, savedStopArea := range manager.saved {
		stopAreas = append(stopAreas, *(savedStopArea.copy()))
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

func (manager *TransactionalStopAreas) FindByOrigin(origin string) (stopAreaIds []StopAreaId) {
	for _, stopAreaId := range manager.model.StopAreas().FindByOrigin(origin) {
		_, ok := manager.deleted[stopAreaId]
		if !ok {
			stopAreaIds = append(stopAreaIds, stopAreaId)
		}
	}
	return
}

func (manager *TransactionalStopAreas) FindFamily(stopAreaId StopAreaId) (stopAreaIds []StopAreaId) {
	for _, stopAreaId := range manager.model.StopAreas().FindFamily(stopAreaId) {
		_, ok := manager.deleted[stopAreaId]
		if !ok {
			stopAreaIds = append(stopAreaIds, stopAreaId)
		}
	}
	return
}

func (manager *TransactionalStopAreas) FindAscendants(stopAreaId StopAreaId) (stopAreas []StopArea) {
	for _, stopArea := range manager.model.StopAreas().FindAscendants(stopAreaId) {
		_, deleted := manager.deleted[stopArea.Id()]
		if deleted {
			continue
		}
		savedStopArea, saved := manager.saved[stopArea.Id()]
		if saved {
			stopAreas = append(stopAreas, *(savedStopArea.copy()))
		} else {
			stopAreas = append(stopAreas, stopArea)
		}

	}
	return
}

func (manager *TransactionalStopAreas) FindAscendantsWithObjectIdKind(stopAreaId StopAreaId, kind string) (stopAreaObjectIds []ObjectID) {
	for _, stopAreaObjectId := range manager.model.StopAreas().FindAscendantsWithObjectIdKind(stopAreaId, kind) {
		if !manager.findByObjectIdInDeleted(stopAreaObjectId) {
			stopAreaObjectIds = append(stopAreaObjectIds, stopAreaObjectId)
		}
	}
	return
}

func (manager *TransactionalStopAreas) findByObjectIdInDeleted(objectid ObjectID) bool {
	for _, sa := range manager.deleted {
		if id, ok := sa.ObjectID(objectid.Kind()); ok && id.Value() == objectid.Value() {
			return true
		}
	}
	return false
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
