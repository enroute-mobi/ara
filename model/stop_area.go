package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopAreaId ModelId

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

	Longitude float64 `json:",omitempty"`
	Latitude  float64 `json:",omitempty"`
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

func (stopArea *StopArea) modelId() ModelId {
	return ModelId(stopArea.id)
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
		Origins    map[string]bool
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

	if aux.Origins != nil {
		stopArea.Origins.SetOriginsFromMap(aux.Origins)
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

func (stopArea *StopArea) Referent() (StopArea, bool) {
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

func (stopArea *StopArea) SetPartnerStatus(partner string, status bool) {
	stopArea.Origins.SetPartnerStatus(partner, status)
	stopArea.Monitored = stopArea.Origins.Monitored()
}

func (stopArea *StopArea) Save() (ok bool) {
	ok = stopArea.model.StopAreas().Save(stopArea)
	return
}

type MemoryStopAreas struct {
	uuid.UUIDConsumer

	model *MemoryModel

	mutex        *sync.RWMutex
	byIdentifier map[StopAreaId]*StopArea
	byObjectId   *ObjectIdIndex

	broadcastEvent func(event StopMonitoringBroadcastEvent)
}

type StopAreas interface {
	uuid.UUIDInterface

	New() StopArea
	Find(id StopAreaId) (StopArea, bool)
	FindByObjectId(objectid ObjectID) (StopArea, bool)
	FindByLineId(id LineId) []StopArea
	FindByOrigin(origin string) []StopAreaId
	FindAll() []StopArea
	FindFamily(id StopAreaId) []StopAreaId
	FindAscendants(id StopAreaId) (stopAreas []StopArea)
	FindAscendantsWithObjectIdKind(stopAreaId StopAreaId, kind string) (stopAreaIds []ObjectID)
	Save(stopArea *StopArea) bool
	Delete(stopArea *StopArea) bool
}

func NewMemoryStopAreas() *MemoryStopAreas {
	return &MemoryStopAreas{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[StopAreaId]*StopArea),
		byObjectId:   NewObjectIdIndex(),
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

	stopArea, ok := manager.byIdentifier[id]
	if ok {
		manager.mutex.RUnlock()
		return *(stopArea.copy()), true
	} else {
		manager.mutex.RUnlock()
		return StopArea{}, false
	}
}

func (manager *MemoryStopAreas) FindByObjectId(objectid ObjectID) (StopArea, bool) {
	manager.mutex.RLock()

	id, ok := manager.byObjectId.Find(objectid)
	if ok {
		manager.mutex.RUnlock()
		return *manager.byIdentifier[StopAreaId(id)], true
	}

	manager.mutex.RUnlock()
	return StopArea{}, false
}

func (manager *MemoryStopAreas) FindByLineId(id LineId) (stopAreas []StopArea) {
	manager.mutex.RLock()

	for _, stopArea := range manager.byIdentifier {
		if stopArea.LineIds.Contains(id) {
			stopAreas = append(stopAreas, *(stopArea.copy()))
		}
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopAreas) FindByOrigin(origin string) (stopAreas []StopAreaId) {
	manager.mutex.RLock()

	for _, stopArea := range manager.byIdentifier {
		if _, ok := stopArea.Origins.Origin(origin); ok {
			stopAreas = append(stopAreas, stopArea.Id())
		}
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopAreas) FindAll() (stopAreas []StopArea) {
	manager.mutex.RLock()

	for _, stopArea := range manager.byIdentifier {
		stopAreas = append(stopAreas, *(stopArea.copy()))
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopAreas) Save(stopArea *StopArea) bool {
	if stopArea.Id() == "" {
		stopArea.id = StopAreaId(manager.NewUUID())
	}

	manager.mutex.Lock()

	stopArea.model = manager.model
	manager.byIdentifier[stopArea.Id()] = stopArea
	manager.byObjectId.Index(stopArea)

	manager.mutex.Unlock()

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
	manager.byObjectId.Delete(ModelId(stopArea.id))

	return true
}

func (manager *MemoryStopAreas) FindFamily(stopAreaId StopAreaId) (stopAreaIds []StopAreaId) {
	manager.mutex.RLock()

	stopAreaIds = manager.findFamily(stopAreaId)

	manager.mutex.RUnlock()

	return
}

func (manager *MemoryStopAreas) findFamily(stopAreaId StopAreaId) (stopAreaIds []StopAreaId) {
	stopAreaIds = []StopAreaId{stopAreaId}
	for _, stopArea := range manager.byIdentifier {
		if stopArea.ParentId == stopAreaId || stopArea.ReferentId == stopAreaId {
			stopAreaIds = append(stopAreaIds, manager.findFamily(stopArea.id)...)
		}
	}
	return
}

func (manager *MemoryStopAreas) FindAscendants(stopAreaId StopAreaId) (stopAreas []StopArea) {
	manager.mutex.RLock()

	stopAreas = manager.findAscendants(stopAreaId)

	manager.mutex.RUnlock()

	return
}

func (manager *MemoryStopAreas) findAscendants(stopAreaId StopAreaId) (stopAreas []StopArea) {
	sa, ok := manager.byIdentifier[stopAreaId]
	if !ok {
		return
	}
	stopAreas = []StopArea{*(sa.copy())}
	if sa.ParentId != "" {
		stopAreas = append(stopAreas, manager.findAscendants(sa.ParentId)...)
	}
	if sa.ReferentId != "" {
		stopAreas = append(stopAreas, manager.findAscendants(sa.ReferentId)...)
	}

	return
}

func (manager *MemoryStopAreas) FindAscendantsWithObjectIdKind(stopAreaId StopAreaId, kind string) (stopAreaObjectIds []ObjectID) {
	manager.mutex.RLock()

	stopAreaObjectIds = manager.findAscendantsWithObjectIdKind(stopAreaId, kind)

	manager.mutex.RUnlock()

	return
}

func (manager *MemoryStopAreas) findAscendantsWithObjectIdKind(stopAreaId StopAreaId, kind string) (stopAreaObjectIds []ObjectID) {
	sa, ok := manager.byIdentifier[stopAreaId]
	if !ok {
		return
	}

	id, ok := sa.ObjectID(kind)
	if ok {
		stopAreaObjectIds = []ObjectID{id}
	}

	if sa.ParentId != "" {
		stopAreaObjectIds = append(stopAreaObjectIds, manager.findAscendantsWithObjectIdKind(sa.ParentId, kind)...)
	}
	if sa.ReferentId != "" {
		stopAreaObjectIds = append(stopAreaObjectIds, manager.findAscendantsWithObjectIdKind(sa.ReferentId, kind)...)
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
		if stopArea.CollectedAlways { // To prevent too much spam when initializing
			rand_duration := time.Duration(rand.Intn(30)) * time.Second
			stopArea.NextCollect(clock.DefaultClock().Now().Add(rand_duration))
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
			references := make(map[string]Reference)
			if err = json.Unmarshal([]byte(sa.References.String), &references); err != nil {
				return err
			}
			stopArea.References.SetReferences(references)
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
