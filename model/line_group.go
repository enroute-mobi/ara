package model

import (
	"encoding/json"
	"sync"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

type LineGroupId ModelId

type LineGroup struct {
	model     Model
	id        LineGroupId `json:",omitempty"`
	Name      string      `json:",omitempty"`
	ShortName string      `json:",omitempty"`
	LineIds   []LineId    `json:",omitempty"`
}

func NewLineGroup(model Model) *LineGroup {
	lineGroup := &LineGroup{
		model: model,
	}

	return lineGroup
}

func (lineGroup *LineGroup) copy() *LineGroup {
	l := *lineGroup
	return &l
}

func (lineGroup *LineGroup) Id() LineGroupId {
	return lineGroup.id
}
func (lineGroup *LineGroup) Save() bool {
	return lineGroup.model.LineGroups().Save(lineGroup)
}

type MemoryLineGroups struct {
	uuid.UUIDConsumer

	model *MemoryModel
	mutex *sync.RWMutex

	byIdentifier map[LineGroupId]*LineGroup
}

type LineGroups interface {
	uuid.UUIDInterface

	New() *LineGroup

	Find(LineGroupId) (*LineGroup, bool)
	FindAll() []*LineGroup
	Save(*LineGroup) bool
	Delete(*LineGroup) bool
}

func NewMemoryLineGroups() *MemoryLineGroups {
	return &MemoryLineGroups{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[LineGroupId]*LineGroup),
	}
}

func (manager *MemoryLineGroups) Find(id LineGroupId) (*LineGroup, bool) {
	manager.mutex.RLock()
	lineGroup, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return lineGroup.copy(), true
	}
	return &LineGroup{}, false
}

func (manager *MemoryLineGroups) New() *LineGroup {
	return NewLineGroup(manager.model)
}

func (manager *MemoryLineGroups) FindAll() (lineGroups []*LineGroup) {
	manager.mutex.RLock()

	for _, lineGroup := range manager.byIdentifier {
		lineGroups = append(lineGroups, lineGroup.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryLineGroups) Save(lineGroup *LineGroup) bool {
	manager.mutex.Lock()

	if lineGroup.Id() == "" {
		lineGroup.id = LineGroupId(manager.NewUUID())
	}
	lineGroup.model = manager.model
	manager.byIdentifier[lineGroup.Id()] = lineGroup

	manager.mutex.Unlock()
	return true
}

func (manager *MemoryLineGroups) Delete(lineGroup *LineGroup) bool {
	manager.mutex.Lock()

	delete(manager.byIdentifier, lineGroup.Id())

	manager.mutex.Unlock()
	return true
}

func (lineGroup *LineGroup) MarshalJSON() ([]byte, error) {
	type Alias LineGroup

	aux := struct {
		*Alias
		Id LineGroupId
	}{
		Id:    lineGroup.id,
		Alias: (*Alias)(lineGroup),
	}

	return json.Marshal(&aux)
}
