package model

import (
	"encoding/json"
	"fmt"
	"sync"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopAreaGroupId ModelId

type StopAreaGroup struct {
	model       Model
	id          StopAreaGroupId `json:",omitempty"`
	Name        string          `json:",omitempty"`
	ShortName   string          `json:",omitempty"`
	StopAreaIds []StopAreaId    `json:",omitempty"`
}

func NewStopAreaGroup(model Model) *StopAreaGroup {
	stopAreaGroup := &StopAreaGroup{
		model: model,
	}

	return stopAreaGroup
}

func (stopAreaGroup *StopAreaGroup) copy() *StopAreaGroup {
	l := *stopAreaGroup
	return &l
}

func (stopAreaGroup *StopAreaGroup) Id() StopAreaGroupId {
	return stopAreaGroup.id
}
func (stopAreaGroup *StopAreaGroup) Save() bool {
	return stopAreaGroup.model.StopAreaGroups().Save(stopAreaGroup)
}

type MemoryStopAreaGroups struct {
	uuid.UUIDConsumer

	model *MemoryModel
	mutex *sync.RWMutex

	byIdentifier map[StopAreaGroupId]*StopAreaGroup
	byShortName  map[string]*StopAreaGroup
}

type StopAreaGroups interface {
	uuid.UUIDInterface

	New() *StopAreaGroup

	Find(StopAreaGroupId) (*StopAreaGroup, bool)
	FindByShortName(string) (*StopAreaGroup, bool)
	FindAll() []*StopAreaGroup
	Save(*StopAreaGroup) bool
	Delete(*StopAreaGroup) bool
}

func NewMemoryStopAreaGroups() *MemoryStopAreaGroups {
	return &MemoryStopAreaGroups{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[StopAreaGroupId]*StopAreaGroup),
		byShortName:  make(map[string]*StopAreaGroup),
	}
}

func (manager *MemoryStopAreaGroups) Find(id StopAreaGroupId) (*StopAreaGroup, bool) {
	manager.mutex.RLock()
	stopAreaGroup, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return stopAreaGroup.copy(), true
	}
	return &StopAreaGroup{}, false
}

func (manager *MemoryStopAreaGroups) FindByShortName(shortName string) (*StopAreaGroup, bool) {
	manager.mutex.RLock()
	stopAreaGroup, ok := manager.byShortName[shortName]
	manager.mutex.RUnlock()

	if ok {
		return stopAreaGroup.copy(), true
	}
	return &StopAreaGroup{}, false
}

func (manager *MemoryStopAreaGroups) New() *StopAreaGroup {
	return NewStopAreaGroup(manager.model)
}

func (manager *MemoryStopAreaGroups) FindAll() (stopAreaGroups []*StopAreaGroup) {
	manager.mutex.RLock()

	for _, stopAreaGroup := range manager.byIdentifier {
		stopAreaGroups = append(stopAreaGroups, stopAreaGroup.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopAreaGroups) Save(stopAreaGroup *StopAreaGroup) bool {
	manager.mutex.Lock()

	if stopAreaGroup.Id() == "" {
		stopAreaGroup.id = StopAreaGroupId(manager.NewUUID())
	}
	stopAreaGroup.model = manager.model
	manager.byIdentifier[stopAreaGroup.Id()] = stopAreaGroup
	manager.byShortName[stopAreaGroup.ShortName] = stopAreaGroup

	manager.mutex.Unlock()
	return true
}

func (manager *MemoryStopAreaGroups) Delete(stopAreaGroup *StopAreaGroup) bool {
	manager.mutex.Lock()

	delete(manager.byIdentifier, stopAreaGroup.Id())
	delete(manager.byShortName, stopAreaGroup.ShortName)

	manager.mutex.Unlock()
	return true
}

func (stopAreaGroup *StopAreaGroup) MarshalJSON() ([]byte, error) {
	type Alias StopAreaGroup

	aux := struct {
		*Alias
		Id StopAreaGroupId
	}{
		Id:    stopAreaGroup.id,
		Alias: (*Alias)(stopAreaGroup),
	}

	return json.Marshal(&aux)
}

func (manager *MemoryStopAreaGroups) Load(referentialSlug string) error {
	var selectStopAreaGroups []SelectStopAreaGroup
	modelName := manager.model.Date()
	sqlQuery := fmt.Sprintf("select * from stop_area_groups where referential_slug = '%s' and model_name = '%s'", referentialSlug, modelName.String())
	_, err := Database.Select(&selectStopAreaGroups, sqlQuery)
	if err != nil {
		return err
	}

	for _, sag := range selectStopAreaGroups {
		stopAreaGroup := manager.New()
		stopAreaGroup.id = StopAreaGroupId(sag.Id)

		if sag.Name.Valid {
			stopAreaGroup.Name = sag.Name.String
		}

		if sag.ShortName.Valid {
			stopAreaGroup.ShortName = sag.ShortName.String
		}

		if sag.StopAreaIds.Valid && len(sag.StopAreaIds.String) > 0 {
			var stopAreaIds []StopAreaId
			if err = json.Unmarshal([]byte(sag.StopAreaIds.String), &stopAreaIds); err != nil {
				return err
			}
			for i := range stopAreaIds {
				stopAreaGroup.StopAreaIds = append(stopAreaGroup.StopAreaIds, StopAreaId(stopAreaIds[i]))
			}
		}

		manager.Save(stopAreaGroup)
	}
	return nil
}
