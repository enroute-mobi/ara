package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type StopAreaId string

type StopAreaAttributes struct {
	ObjectId        ObjectID
	ParentObjectId  ObjectID
	Name            string
	CollectedAlways bool
}

type StopArea struct {
	ObjectIDConsumer
	model Model

	id       StopAreaId
	ParentId StopAreaId `json:",omitempty"`

	NextCollectAt          time.Time
	collectedAt            time.Time
	CollectedUntil         time.Time
	CollectedAlways        bool
	CollectGeneralMessages bool

	Monitored bool

	Name            string
	LineIds         StopAreaLineIds `json:"Lines,omitempty"`
	CollectChildren bool
	Attributes      Attributes
	References      References
	// ...
}

func NewStopArea(model Model) *StopArea {
	stopArea := &StopArea{
		model:           model,
		Attributes:      NewAttributes(),
		References:      NewReferences(),
		CollectedAlways: true,
	}
	stopArea.objectids = make(ObjectIDs)
	return stopArea
}

func (stopArea *StopArea) Id() StopAreaId {
	return stopArea.id
}

func (stopArea *StopArea) NextCollect(collectTime time.Time) {
	stopArea.NextCollectAt = collectTime
}

func (stopArea *StopArea) CollectedAt() time.Time {
	return stopArea.collectedAt
}

func (stopArea *StopArea) Updated(updateTime time.Time) {
	stopArea.collectedAt = updateTime
}

func (stopArea *StopArea) MarshalJSON() ([]byte, error) {
	type Alias StopArea
	aux := struct {
		Id             StopAreaId
		ObjectIDs      ObjectIDs            `json:",omitempty"`
		NextCollectAt  *time.Time           `json:",omitempty"`
		CollectedAt    *time.Time           `json:",omitempty"`
		CollectedUntil *time.Time           `json:",omitempty"`
		Attributes     Attributes           `json:",omitempty"`
		References     map[string]Reference `json:",omitempty"`
		*Alias
	}{
		Id:    stopArea.id,
		Alias: (*Alias)(stopArea),
	}

	if !stopArea.ObjectIDs().Empty() {
		aux.ObjectIDs = stopArea.ObjectIDs()
	}
	if !stopArea.Attributes.IsEmpty() {
		aux.Attributes = stopArea.Attributes
	}
	if !stopArea.References.IsEmpty() {
		aux.References = stopArea.References.GetReferences()
	}
	if !stopArea.NextCollectAt.IsZero() {
		aux.NextCollectAt = &stopArea.NextCollectAt
	}
	if !stopArea.collectedAt.IsZero() {
		aux.CollectedAt = &stopArea.collectedAt
	}
	if !stopArea.CollectedAlways && !stopArea.CollectedUntil.IsZero() {
		aux.CollectedUntil = &stopArea.CollectedUntil
	}

	return json.Marshal(&aux)
}

func (stopArea *StopArea) UnmarshalJSON(data []byte) error {
	type Alias StopArea
	aux := &struct {
		ObjectIDs  map[string]string
		References map[string]Reference
		*Alias
	}{
		Alias: (*Alias)(stopArea),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ObjectIDs != nil {
		stopArea.ObjectIDConsumer.objectids = NewObjectIDsFromMap(aux.ObjectIDs)
	}

	if aux.References != nil {
		stopArea.References.SetReferences(aux.References)
	}

	return nil
}

func (stopArea *StopArea) Attribute(key string) (string, bool) {
	value, present := stopArea.Attributes[key]
	return value, present
}

func (stopArea *StopArea) Reference(key string) (Reference, bool) {
	value, present := stopArea.References.Get(key)
	return value, present
}

func (stopArea *StopArea) Lines() (lines []Line) {
	for _, lineId := range stopArea.LineIds {
		foundLine, ok := stopArea.model.Lines().Find(lineId)
		if ok {
			lines = append(lines, foundLine)
		}
	}
	return
}

func (stopArea *StopArea) Parent() (StopArea, bool) {
	return stopArea.model.StopAreas().Find(stopArea.ParentId)
}

func (stopArea *StopArea) Save() (ok bool) {
	ok = stopArea.model.StopAreas().Save(stopArea)
	return
}

type MemoryStopAreas struct {
	UUIDConsumer

	model *MemoryModel

	mutex        *sync.RWMutex
	byIdentifier map[StopAreaId]*StopArea
}

type StopAreas interface {
	UUIDInterface

	New() StopArea
	Find(id StopAreaId) (StopArea, bool)
	FindByObjectId(objectid ObjectID) (StopArea, bool)
	FindAll() []StopArea
	FindFamily(stopAreaId StopAreaId) []StopAreaId
	Save(stopArea *StopArea) bool
	Delete(stopArea *StopArea) bool
}

func NewMemoryStopAreas() *MemoryStopAreas {
	return &MemoryStopAreas{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[StopAreaId]*StopArea),
	}
}

func (manager *MemoryStopAreas) Clone(model *MemoryModel) *MemoryStopAreas {
	clone := NewMemoryStopAreas()
	clone.model = model

	for _, stopArea := range manager.byIdentifier {
		cloneStopArea := *stopArea
		cloneStopArea.id = StopAreaId("")
		clone.Save(&cloneStopArea)
	}

	return clone
}

func (manager *MemoryStopAreas) New() StopArea {
	stopArea := NewStopArea(manager.model)
	return *stopArea
}

func (manager *MemoryStopAreas) Find(id StopAreaId) (StopArea, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	stopArea, ok := manager.byIdentifier[id]
	if ok {
		return *stopArea, true
	} else {
		return StopArea{}, false
	}
}

func (manager *MemoryStopAreas) FindByObjectId(objectid ObjectID) (StopArea, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, stopArea := range manager.byIdentifier {
		stopAreaObjectId, _ := stopArea.ObjectID(objectid.Kind())
		if stopAreaObjectId.Value() == objectid.Value() {
			return *stopArea, true
		}
	}
	return StopArea{}, false
}

func (manager *MemoryStopAreas) FindAll() (stopAreas []StopArea) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if len(manager.byIdentifier) == 0 {
		return []StopArea{}
	}
	for _, stopArea := range manager.byIdentifier {
		stopAreas = append(stopAreas, *stopArea)
	}
	return
}

func (manager *MemoryStopAreas) Save(stopArea *StopArea) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if stopArea.Id() == "" {
		stopArea.id = StopAreaId(manager.NewUUID())
	}

	stopArea.model = manager.model
	manager.byIdentifier[stopArea.Id()] = stopArea
	return true
}

func (manager *MemoryStopAreas) Delete(stopArea *StopArea) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, stopArea.Id())
	return true
}

func (manager *MemoryStopAreas) FindFamily(stopAreaId StopAreaId) (stopAreaIds []StopAreaId) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	stopAreaIds = []StopAreaId{stopAreaId}
	for _, stopArea := range manager.byIdentifier {
		if stopArea.ParentId == stopAreaId {
			stopAreaIds = append(stopAreaIds, manager.FindFamily(stopArea.id)...)
		}
	}
	return stopAreaIds
}

func (manager *MemoryStopAreas) Load(referentialSlug string) error {
	var selectStopAreas []SelectStopArea
	modelName := manager.model.Date()

	sqlQuery := fmt.Sprintf("select * from stop_areas where referential_slug = '%s' and model_name = '%s'", referentialSlug, modelName.String())
	_, err := Database.Select(&selectStopAreas, sqlQuery)
	if err != nil {
		return err
	}

	for _, sa := range selectStopAreas {
		stopArea := manager.New()
		stopArea.id = StopAreaId(sa.Id)
		if sa.Name.Valid {
			stopArea.Name = sa.Name.String
		}
		if sa.ParentId.Valid {
			stopArea.ParentId = StopAreaId(sa.ParentId.String)
		}
		if sa.CollectedAlways.Valid {
			stopArea.CollectedAlways = sa.CollectedAlways.Bool
		}
		if sa.CollectChildren.Valid {
			stopArea.CollectChildren = sa.CollectChildren.Bool
		}
		if sa.CollectGeneralMessages.Valid {
			stopArea.CollectGeneralMessages = sa.CollectGeneralMessages.Bool
		}

		if sa.LineIds.Valid && len(sa.LineIds.String) > 0 {
			var lineIds []string
			if err = json.Unmarshal([]byte(sa.LineIds.String), &lineIds); err != nil {
				return err
			}
			for i := range lineIds {
				stopArea.LineIds = append(stopArea.LineIds, LineId(lineIds[i]))
			}
		}

		if sa.Attributes.Valid && len(sa.Attributes.String) > 0 {
			if err = json.Unmarshal([]byte(sa.Attributes.String), &stopArea.Attributes); err != nil {
				return err
			}
		}

		if sa.References.Valid && len(sa.References.String) > 0 {
			if err = json.Unmarshal([]byte(sa.References.String), &stopArea.References); err != nil {
				return err
			}
		}

		if sa.ObjectIDs.Valid && len(sa.ObjectIDs.String) > 0 {
			objectIdMap := make(map[string]string)
			if err = json.Unmarshal([]byte(sa.ObjectIDs.String), &objectIdMap); err != nil {
				return err
			}
			stopArea.objectids = NewObjectIDsFromMap(objectIdMap)
		}

		manager.Save(&stopArea)
	}
	return nil
}
