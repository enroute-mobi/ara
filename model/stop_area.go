package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type StopAreaId string

type StopAreaAttributes struct {
	ObjectId        ObjectID
	Name            string
	CollectedAlways bool
}

type StopArea struct {
	ObjectIDConsumer
	model Model

	id              StopAreaId
	requestedAt     time.Time
	collectedAt     time.Time
	CollectedUntil  time.Time
	CollectedAlways bool

	Name       string
	Attributes Attributes
	References References
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

func (stopArea *StopArea) RequestedAt() time.Time {
	return stopArea.requestedAt
}

func (stopArea *StopArea) Requested(requestTime time.Time) {
	stopArea.requestedAt = requestTime
}

func (stopArea *StopArea) CollectedAt() time.Time {
	return stopArea.collectedAt
}

func (stopArea *StopArea) Updated(updateTime time.Time) {
	stopArea.collectedAt = updateTime
}

func (stopArea *StopArea) FillStopArea(stopAreaMap map[string]interface{}) {
	if stopArea.id != "" {
		stopAreaMap["Id"] = stopArea.id
	}

	if stopArea.Name != "" {
		stopAreaMap["Name"] = stopArea.Name
	}

	if !stopArea.Attributes.IsEmpty() {
		stopAreaMap["Attributes"] = stopArea.Attributes
	}

	if !stopArea.References.IsEmpty() {
		stopAreaMap["References"] = stopArea.References
	}

	if !stopArea.requestedAt.IsZero() {
		stopAreaMap["RequestedAt"] = stopArea.requestedAt
	}
	if !stopArea.collectedAt.IsZero() {
		stopAreaMap["CollectedAt"] = stopArea.collectedAt
	}
	if !stopArea.ObjectIDs().Empty() {
		stopAreaMap["ObjectIDs"] = stopArea.ObjectIDs()
	}
	if stopAreaMap["CollectedAlways"] == false {
		stopAreaMap["CollectedUntil"] = stopArea.CollectedUntil
	}
	stopAreaMap["CollectedAlways"] = stopArea.CollectedAlways
}

func (stopArea *StopArea) MarshalJSON() ([]byte, error) {
	stopAreaMap := make(map[string]interface{})

	stopArea.FillStopArea(stopAreaMap)

	return json.Marshal(stopAreaMap)
}

func (stopArea *StopArea) Attribute(key string) (string, bool) {
	value, present := stopArea.Attributes[key]
	return value, present
}

func (stopArea *StopArea) Reference(key string) (Reference, bool) {
	value, present := stopArea.References[key]
	return value, present
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

	return nil
}

func (stopArea *StopArea) Save() (ok bool) {
	ok = stopArea.model.StopAreas().Save(stopArea)
	return
}

type MemoryStopAreas struct {
	UUIDConsumer

	model *MemoryModel

	byIdentifier map[StopAreaId]*StopArea
}

type StopAreas interface {
	UUIDInterface

	New() StopArea
	Find(id StopAreaId) (StopArea, bool)
	FindByObjectId(objectid ObjectID) (StopArea, bool)
	FindAll() []StopArea
	Save(stopArea *StopArea) bool
	Delete(stopArea *StopArea) bool
}

func NewMemoryStopAreas() *MemoryStopAreas {
	return &MemoryStopAreas{
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
	stopArea, ok := manager.byIdentifier[id]
	if ok {
		return *stopArea, true
	} else {
		return StopArea{}, false
	}
}

func (manager *MemoryStopAreas) FindByObjectId(objectid ObjectID) (StopArea, bool) {
	for _, stopArea := range manager.byIdentifier {
		stopAreaObjectId, _ := stopArea.ObjectID(objectid.Kind())
		if stopAreaObjectId.Value() == objectid.Value() {
			return *stopArea, true
		}
	}
	return StopArea{}, false
}

func (manager *MemoryStopAreas) FindAll() (stopAreas []StopArea) {
	if len(manager.byIdentifier) == 0 {
		return []StopArea{}
	}
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

func (manager *MemoryStopAreas) Load(referentialId string) error {
	var selectStopAreas []struct {
		Id              string
		ReferentialId   string `db:"referential_id"`
		Name            string
		ObjectIDs       string `db:"object_ids"`
		Attributes      string
		References      string    `db:"siri_references"`
		RequestedAt     time.Time `db:"requested_at"`
		CollectedAt     time.Time `db:"collected_at"`
		CollectedUntil  time.Time `db:"collected_until"`
		CollectedAlways bool      `db:"collected_always"`
	}
	sqlQuery := fmt.Sprintf("select * from stop_areas where referential_id = '%s'", referentialId)
	_, err := Database.Select(&selectStopAreas, sqlQuery)
	if err != nil {
		return err
	}
	for _, sa := range selectStopAreas {
		stopArea := manager.New()
		stopArea.id = StopAreaId(sa.Id)
		stopArea.Name = sa.Name
		stopArea.requestedAt = sa.RequestedAt
		stopArea.collectedAt = sa.CollectedAt
		stopArea.CollectedAlways = sa.CollectedAlways
		stopArea.CollectedUntil = sa.CollectedUntil

		if err = json.Unmarshal([]byte(sa.Attributes), &stopArea.Attributes); err != nil {
			return err
		}

		if err = json.Unmarshal([]byte(sa.References), &stopArea.References); err != nil {
			return err
		}

		objectIdMap := make(map[string]string)
		if err = json.Unmarshal([]byte(sa.ObjectIDs), &objectIdMap); err != nil {
			return err
		}
		stopArea.objectids = NewObjectIDsFromMap(objectIdMap)

		manager.Save(&stopArea)
	}
	return nil
}
