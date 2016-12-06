package model

import (
	"encoding/json"
	"time"
)

type StopAreaId string

type StopArea struct {
	model Model

	id          StopAreaId
	objectids   ObjectIDs
	requestedAt time.Time
	updatedAt   time.Time

	Name string
	// ...
}

func NewStopArea(model Model) *StopArea {
	return &StopArea{model: model, objectids: make(ObjectIDs)}
}

func (stopArea *StopArea) Id() StopAreaId {
	return stopArea.id
}

func (stopArea *StopArea) RequestedAt() time.Time {
	return stopArea.requestedAt
}

func (stopArea *StopArea) Requested(requestTime time.Time) {
	stopArea.requestedAt = requestTime
}

func (stopArea *StopArea) UpdatedAt() time.Time {
	return stopArea.updatedAt
}

func (stopArea *StopArea) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Id":   stopArea.id,
		"Name": stopArea.Name,
	})
}

func (stopArea *StopArea) Save() (ok bool) {
	ok = stopArea.model.StopAreas().Save(stopArea)
	return
}

func (stopArea *StopArea) ObjectID(kind string) (ObjectID, bool) {
	objectid, ok := stopArea.objectids[kind]
	if ok {
		return objectid, true
	}
	return ObjectID{}, false
}

func (stopArea *StopArea) SetObjectID(objectid ObjectID) {
	stopArea.objectids[objectid.Kind()] = objectid
}

func (stopArea *StopArea) ObjectIDs() (objectidArray []ObjectID) {
	for _, objectid := range stopArea.objectids {
		objectidArray = append(objectidArray, objectid)
	}
	return
}

type MemoryStopAreas struct {
	UUIDConsumer

	model Model

	byIdentifier map[StopAreaId]*StopArea
}

type StopAreas interface {
	UUIDInterface

	New() StopArea
	Find(id StopAreaId) (StopArea, bool)
	FindAll() []StopArea
	Save(stopArea *StopArea) bool
	Delete(stopArea *StopArea) bool
}

func NewMemoryStopAreas() *MemoryStopAreas {
	return &MemoryStopAreas{
		byIdentifier: make(map[StopAreaId]*StopArea),
	}
}

func (manager *MemoryStopAreas) New() StopArea {
	return StopArea{model: manager.model, objectids: make(ObjectIDs)}
}

func (manager *MemoryStopAreas) Find(id StopAreaId) (StopArea, bool) {
	stopArea, ok := manager.byIdentifier[id]
	if ok {
		return *stopArea, true
	} else {
		return StopArea{}, false
	}
}

func (manager *MemoryStopAreas) FindAll() (stopAreas []StopArea) {
	for _, stopArea := range manager.byIdentifier {
		stopAreas = append(stopAreas, *stopArea)
	}
	return
}

func (manager *MemoryStopAreas) Save(stopArea *StopArea) bool {
	if stopArea.Id() == "" {
		stopArea.id = StopAreaId(manager.NewUUID())
	}
	stopArea.model = manager.model
	manager.byIdentifier[stopArea.Id()] = stopArea
	return true
}

func (manager *MemoryStopAreas) Delete(stopArea *StopArea) bool {
	delete(manager.byIdentifier, stopArea.Id())
	return true
}
