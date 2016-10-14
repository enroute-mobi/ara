package model

type StopAreaId int64

type StopArea struct {
	id StopAreaId

	// Name string
	// ...
}

func (stopArea *StopArea) Id() StopAreaId {
	return stopArea.id
}

type MemoryStopAreas struct {
	byIdentifier   map[StopAreaId]*StopArea
	lastIdentifier StopAreaId
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

func (manager *MemoryStopAreas) Save(stopArea *StopArea) bool {
	if stopArea.Id() == 0 {
		stopArea.id = manager.lastIdentifier + 1
		manager.lastIdentifier = stopArea.id
	}
	manager.byIdentifier[stopArea.Id()] = stopArea
	return true
}

func (manager *MemoryStopAreas) Delete(stopArea *StopArea) bool {
	delete(manager.byIdentifier, stopArea.Id())
	return true
}
