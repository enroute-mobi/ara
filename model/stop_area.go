package model

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopAreaId ModelId

type StopArea struct {
	Collectable
	CollectedUntil time.Time
	model          Model
	References     References
	CodeConsumer
	Origins           *StopAreaOrigins
	Attributes        Attributes
	ReferentId        StopAreaId `json:",omitempty"`
	id                StopAreaId
	ParentId          StopAreaId `json:",omitempty"`
	Name              string
	LineIds           StopAreaLineIds `json:"Lines,omitempty"`
	Latitude          float64         `json:",omitempty"`
	Longitude         float64         `json:",omitempty"`
	CollectChildren   bool
	CollectSituations bool
	CollectedAlways   bool
	Monitored         bool
}

func NewStopArea(model Model) *StopArea {
	stopArea := &StopArea{
		model:           model,
		Origins:         NewStopAreaOrigins(),
		Attributes:      NewAttributes(),
		References:      NewReferences(),
		CollectedAlways: true,
	}
	stopArea.codes = make(Codes)
	return stopArea
}

func (stopArea *StopArea) modelId() ModelId {
	return ModelId(stopArea.id)
}

func (stopArea *StopArea) copy() *StopArea {
	return &StopArea{
		Collectable: Collectable{
			nextCollectAt: stopArea.nextCollectAt,
			collectedAt:   stopArea.collectedAt,
		},
		CodeConsumer:      stopArea.CodeConsumer.Copy(),
		model:             stopArea.model,
		id:                stopArea.id,
		ParentId:          stopArea.ParentId,
		ReferentId:        stopArea.ReferentId,
		CollectedUntil:    stopArea.CollectedUntil,
		CollectedAlways:   stopArea.CollectedAlways,
		CollectSituations: stopArea.CollectSituations,
		Monitored:         stopArea.Monitored,
		Origins:           stopArea.Origins.Copy(),
		Name:              stopArea.Name,
		LineIds:           stopArea.LineIds.Copy(),
		CollectChildren:   stopArea.CollectChildren,
		Attributes:        stopArea.Attributes.Copy(),
		References:        stopArea.References.Copy(),
		Longitude:         stopArea.Longitude,
		Latitude:          stopArea.Latitude,
	}
}

func (stopArea *StopArea) Id() StopAreaId {
	return stopArea.id
}

func (stopArea *StopArea) MarshalJSON() ([]byte, error) {
	type Alias StopArea
	aux := struct {
		References     map[string]Reference `json:",omitempty"`
		Codes          Codes                `json:",omitempty"`
		NextCollectAt  *time.Time           `json:",omitempty"`
		CollectedAt    *time.Time           `json:",omitempty"`
		CollectedUntil *time.Time           `json:",omitempty"`
		Attributes     Attributes           `json:",omitempty"`
		*Alias
		Id StopAreaId
	}{
		Id:    stopArea.id,
		Alias: (*Alias)(stopArea),
	}

	if !stopArea.Codes().Empty() {
		aux.Codes = stopArea.Codes()
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
		Codes      map[string]string
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

	if aux.Codes != nil {
		stopArea.CodeConsumer.codes = NewCodesFromMap(aux.Codes)
	}

	if aux.References != nil {
		stopArea.References.SetReferences(aux.References)
	}

	if aux.Origins != nil {
		stopArea.Origins.SetOriginsFromMap(aux.Origins)
	}

	return nil
}

func (stopArea *StopArea) Attribute(key string) (value string, present bool) {
	value, present = stopArea.Attributes[key]
	return
}

func (stopArea *StopArea) Reference(key string) (Reference, bool) {
	value, present := stopArea.References.Get(key)
	return value, present
}

func (stopArea *StopArea) Lines() []*Line {
	var lines []*Line
	h := make(map[LineId]*Line)
	for _, lineId := range stopArea.LineIds {
		foundLine, ok := stopArea.model.Lines().Find(lineId)
		if ok {
			h[foundLine.Id()] = foundLine
		}
	}

	particulars := stopArea.model.StopAreas().FindByReferentId(stopArea.Id())
	for i := range particulars {
		ls := particulars[i].Lines()
		for j := range ls {
			h[ls[j].Id()] = ls[j]
		}
	}

	for _, v := range h {
		lines = append(lines, v)
	}

	return lines
}

func (stopArea *StopArea) Parent() (*StopArea, bool) {
	return stopArea.model.StopAreas().Find(stopArea.ParentId)
}

func (stopArea *StopArea) Referent() (*StopArea, bool) {
	return stopArea.model.StopAreas().Find(stopArea.ReferentId)
}

func (stopArea *StopArea) ReferentOrSelfCode(codeSpace string) (Code, bool) {
	ref, ok := stopArea.Referent()
	if ok {
		code, ok := ref.Code(codeSpace)
		if ok {
			return code, true
		}
	}
	code, ok := stopArea.Code(codeSpace)
	if ok {
		return code, true
	}
	return Code{}, false
}

func (stopArea *StopArea) SetPartnerStatus(partner string, status bool) {
	stopArea.Origins.SetPartnerStatus(partner, status)
	stopArea.Monitored = stopArea.Origins.Monitored()
}

func (stopArea *StopArea) Save() bool {
	return stopArea.model.StopAreas().Save(stopArea)
}

type MemoryStopAreas struct {
	uuid.UUIDConsumer

	model *MemoryModel

	mutex        *sync.RWMutex
	byIdentifier map[StopAreaId]*StopArea
	byParent     *Index
	byReferent   *Index
	byCode       *CodeIndex

	broadcastEvent func(event StopMonitoringBroadcastEvent)
}

type StopAreas interface {
	uuid.UUIDInterface

	New() *StopArea
	Find(StopAreaId) (*StopArea, bool)
	FindByCode(Code) (*StopArea, bool)
	FindByLineId(LineId) []*StopArea
	FindByOrigin(string) []StopAreaId
	FindAll() []*StopArea
	FindAllValues() []StopArea
	FindFamily(StopAreaId) []StopAreaId
	FindByReferentId(StopAreaId) []*StopArea
	FindAscendants(StopAreaId) []*StopArea
	FindAscendantsWithCodeSpace(StopAreaId, string) []Code
	Save(*StopArea) bool
	Delete(*StopArea) bool
}

func NewMemoryStopAreas() *MemoryStopAreas {
	referentExtractor := func(instance ModelInstance) ModelId { return ModelId((instance.(*StopArea)).ReferentId) }
	parentExtractor := func(instance ModelInstance) ModelId { return ModelId((instance.(*StopArea)).ParentId) }

	return &MemoryStopAreas{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[StopAreaId]*StopArea),
		byReferent:   NewIndex(referentExtractor),
		byParent:     NewIndex(parentExtractor),
		byCode:       NewCodeIndex(),
	}
}

func (manager *MemoryStopAreas) New() *StopArea {
	return NewStopArea(manager.model)
}

func (manager *MemoryStopAreas) Find(id StopAreaId) (*StopArea, bool) {
	manager.mutex.RLock()
	stopArea, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return stopArea.copy(), true
	}
	return &StopArea{}, false
}

func (manager *MemoryStopAreas) FindByCode(code Code) (*StopArea, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byCode.Find(code)
	if ok {
		return manager.byIdentifier[StopAreaId(id)].copy(), true
	}

	return &StopArea{}, false
}

func (manager *MemoryStopAreas) FindByLineId(id LineId) (stopAreas []*StopArea) {
	manager.mutex.RLock()

	for _, stopArea := range manager.byIdentifier {
		if stopArea.LineIds.Contains(id) {
			stopAreas = append(stopAreas, stopArea.copy())
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

func (manager *MemoryStopAreas) FindByReferentId(id StopAreaId) (stopAreas []*StopArea) {
	manager.mutex.RLock()

	ids, _ := manager.byReferent.Find(ModelId(id))

	for _, id := range ids {
		sa := manager.byIdentifier[StopAreaId(id)]
		stopAreas = append(stopAreas, sa.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopAreas) FindByParentId(id StopAreaId) (stopAreas []*StopArea) {
	manager.mutex.RLock()

	ids, _ := manager.byParent.Find(ModelId(id))

	for _, id := range ids {
		sa := manager.byIdentifier[StopAreaId(id)]
		stopAreas = append(stopAreas, sa.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopAreas) FindAllValues() (stopAreas []StopArea) {
	manager.mutex.RLock()

	for _, stopArea := range manager.byIdentifier {
		stopAreas = append(stopAreas, *stopArea.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryStopAreas) FindAll() (stopAreas []*StopArea) {
	manager.mutex.RLock()

	for _, stopArea := range manager.byIdentifier {
		stopAreas = append(stopAreas, stopArea.copy())
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
	manager.byReferent.Index(stopArea)
	manager.byParent.Index(stopArea)
	manager.byCode.Index(stopArea)

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
	manager.byReferent.Delete(ModelId(stopArea.id))
	manager.byParent.Delete(ModelId(stopArea.id))
	manager.byCode.Delete(ModelId(stopArea.id))

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

	ids, _ := manager.byParent.Find(ModelId(stopAreaId))
	for _, id := range ids {
		stopAreaIds = append(stopAreaIds, manager.findFamily(StopAreaId(id))...)
	}
	ids, _ = manager.byReferent.Find(ModelId(stopAreaId))
	for _, id := range ids {
		stopAreaIds = append(stopAreaIds, manager.findFamily(StopAreaId(id))...)
	}

	return
}

func (manager *MemoryStopAreas) FindAscendants(stopAreaId StopAreaId) (stopAreas []*StopArea) {
	manager.mutex.RLock()

	var count uint8

	stopAreas = manager.findAscendants(stopAreaId, count)

	manager.mutex.RUnlock()

	return
}

func (manager *MemoryStopAreas) findAscendants(stopAreaId StopAreaId, count uint8) (stopAreas []*StopArea) {
	if count >= 20 {
		logger.Log.Printf("Loop in StopAreas when finding Ascendants: %v", stopAreaId)
		return
	}
	count++

	sa, ok := manager.byIdentifier[stopAreaId]
	if !ok {
		return
	}
	stopAreas = []*StopArea{sa.copy()}
	if sa.ParentId != "" {
		stopAreas = append(stopAreas, manager.findAscendants(sa.ParentId, count)...)
	}
	if sa.ReferentId != "" {
		stopAreas = append(stopAreas, manager.findAscendants(sa.ReferentId, count)...)
	}

	return
}

func (manager *MemoryStopAreas) FindAscendantsWithCodeSpace(stopAreaId StopAreaId, kind string) (stopAreaCodes []Code) {
	manager.mutex.RLock()

	stopAreaCodes = manager.findAscendantsWithCodeSpace(stopAreaId, kind)

	manager.mutex.RUnlock()

	return
}

func (manager *MemoryStopAreas) findAscendantsWithCodeSpace(stopAreaId StopAreaId, kind string) (stopAreaCodes []Code) {
	sa, ok := manager.byIdentifier[stopAreaId]
	if !ok {
		return
	}

	id, ok := sa.Code(kind)
	if ok {
		stopAreaCodes = []Code{id}
	}

	if sa.ParentId != "" {
		stopAreaCodes = append(stopAreaCodes, manager.findAscendantsWithCodeSpace(sa.ParentId, kind)...)
	}
	if sa.ReferentId != "" {
		stopAreaCodes = append(stopAreaCodes, manager.findAscendantsWithCodeSpace(sa.ReferentId, kind)...)
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
		if sa.CollectSituations.Valid {
			stopArea.CollectSituations = sa.CollectSituations.Bool
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

		if sa.Codes.Valid && len(sa.Codes.String) > 0 {
			codeMap := make(map[string]string)
			if err = json.Unmarshal([]byte(sa.Codes.String), &codeMap); err != nil {
				return err
			}
			stopArea.codes = NewCodesFromMap(codeMap)
		}

		manager.Save(stopArea)
	}
	return nil
}
