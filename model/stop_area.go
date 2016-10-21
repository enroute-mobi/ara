package model

import "encoding/json"

type StopAreaId string

type StopArea struct {
	id   StopAreaId
	Name string
	// ...
}

func (stopArea *StopArea) Id() StopAreaId {
	return stopArea.id
}

func (stopArea *StopArea) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Id":   stopArea.id,
		"Name": stopArea.Name,
	})
}

type MemoryStopAreas struct {
	UUIDConsumer

	byIdentifier map[StopAreaId]*StopArea
}

func NewMemoryStopAreas() *MemoryStopAreas {
	return &MemoryStopAreas{
		byIdentifier: make(map[StopAreaId]*StopArea),
	}
}

func (manager *MemoryStopAreas) New() StopArea {
	return StopArea{}
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
	manager.byIdentifier[stopArea.Id()] = stopArea
	return true
}

func (manager *MemoryStopAreas) Delete(stopArea *StopArea) bool {
	delete(manager.byIdentifier, stopArea.Id())
	return true
}
