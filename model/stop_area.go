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

	id         StopAreaId
	ParentId   StopAreaId `json:",omitempty"`
	ReferentId StopAreaId `json:",omitempty"`

	nextCollectAt          time.Time
	collectedAt            time.Time
	CollectedUntil         time.Time
	CollectedAlways        bool
	CollectGeneralMessages bool

	Monitored bool
	Origins   StopAreaOrigins

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
		Origins:         NewStopAreaOrigins(),
		Attributes:      NewAttributes(),
		References:      NewReferences(),
		CollectedAlways: true,
	}
	stopArea.objectids = make(ObjectIDs)
	return stopArea
}

func (stopArea *StopArea) copy() *StopArea {
	s := *stopArea
	s.Origins = *(stopArea.Origins.Copy())
	return &s
}

func (stopArea *StopArea) Id() StopAreaId {
	return stopArea.id
}

func (stopArea *StopArea) NextCollectAt() time.Time {
	return stopArea.nextCollectAt
}

func (stopArea *StopArea) NextCollect(collectTime time.Time) {
	stopArea.nextCollectAt = collectTime
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
	if !stopArea.nextCollectAt.IsZero() {
		aux.NextCollectAt = &stopArea.nextCollectAt
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
	if stopArea.ParentId == "" {
		return StopArea{}, false
	}
	return stopArea.model.StopAreas().Find(stopArea.ParentId)
}

func (stopArea *StopArea) Referent() (StopArea, bool) {
	if stopArea.ReferentId == "" {
		return StopArea{}, false
	}
	return stopArea.model.StopAreas().Find(stopArea.ReferentId)
}

func (stopArea *StopArea) ReferentOrSelfObjectId(objectIDKind string) (ObjectID, bool) {
	ref, ok := stopArea.Referent()
	if ok {
		objectID, ok := ref.ObjectID(objectIDKind)
		if ok {
			return objectID, true
		}
	}
	objectID, ok := stopArea.ObjectID(objectIDKind)
	if ok {
		return objectID, true
	}
	return ObjectID{}, false
}

func (stopArea *StopArea) Save() (ok bool) {
	ok = stopArea.model.StopAreas().Save(stopArea)
	return
}

type MemoryStopAreas struct {
	UUIDConsumer

	model *MemoryModel

	mutex          *sync.RWMutex
	byIdentifier   map[StopAreaId]*StopArea
	broadcastEvent func(event StopMonitoringBroadcastEvent)
}

type StopAreas interface {
	UUIDInterface

	New() StopArea
	Find(id StopAreaId) (StopArea, bool)
	FindByObjectId(objectid ObjectID) (StopArea, bool)
	FindByLineId(id LineId) []StopArea
	FindByOrigin(origin string) []StopAreaId
	FindAll() []StopArea
	FindFamily(id StopAreaId) []StopAreaId
	FindAscendants(id StopAreaId) (stopAreaIds []StopAreaId)
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
		return *(stopArea.copy()), true
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
			return *(stopArea.copy()), true
		}
	}
	return StopArea{}, false
}

func (manager *MemoryStopAreas) FindByLineId(id LineId) (stopAreas []StopArea) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, stopArea := range manager.byIdentifier {
		if stopArea.LineIds.Contains(id) {
			stopAreas = append(stopAreas, *(stopArea.copy()))
		}
	}
	return
}

func (manager *MemoryStopAreas) FindByOrigin(origin string) (stopAreas []StopAreaId) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, stopArea := range manager.byIdentifier {
		if _, ok := stopArea.Origins.Origin(origin); ok {
			stopAreas = append(stopAreas, stopArea.Id())
		}
	}
	return
}

func (manager *MemoryStopAreas) FindAll() (stopAreas []StopArea) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if len(manager.byIdentifier) == 0 {
		return []StopArea{}
	}
	for _, stopArea := range manager.byIdentifier {
		stopAreas = append(stopAreas, *(stopArea.copy()))
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

	event := StopMonitoringBroadcastEvent{
		ModelId:   string(stopArea.id),
		ModelType: "StopArea",
	}

	if manager.broadcastEvent != nil {
		manager.broadcastEvent(event)
	}

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
		if stopArea.ParentId == stopAreaId || stopArea.ReferentId == stopAreaId {
			stopAreaIds = append(stopAreaIds, manager.FindFamily(stopArea.id)...)
		}
	}
	return
}

func (manager *MemoryStopAreas) FindAscendants(stopAreaId StopAreaId) (stopAreaIds []StopAreaId) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	sa, ok := manager.Find(stopAreaId)
	if !ok {
		return
	}
	stopAreaIds = []StopAreaId{stopAreaId}
	if sa.ParentId != "" {
		stopAreaIds = append(stopAreaIds, manager.FindAscendants(sa.ParentId)...)
	}
	if sa.ReferentId != "" {
		stopAreaIds = append(stopAreaIds, manager.FindAscendants(sa.ReferentId)...)
	}

	return
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
		if sa.ReferentId.Valid {
			stopArea.ReferentId = StopAreaId(sa.ReferentId.String)
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
