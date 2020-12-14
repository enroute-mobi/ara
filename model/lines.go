package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

type LineId ModelId

type LineAttributes struct {
	ObjectId ObjectID
	Name     string
}

type Line struct {
	ObjectIDConsumer
	model  Model
	origin string

	id LineId

	CollectGeneralMessages bool
	nextCollectAt          time.Time
	collectedAt            time.Time

	Name       string `json:",omitempty"`
	Attributes Attributes
	References References
}

func NewLine(model Model) *Line {
	line := &Line{
		model:      model,
		Attributes: NewAttributes(),
		References: NewReferences(),
	}

	line.objectids = make(ObjectIDs)
	return line
}

func (line *Line) modelId() ModelId {
	return ModelId(line.id)
}

func (line *Line) Id() LineId {
	return line.id
}

func (line *Line) NextCollectAt() time.Time {
	return line.nextCollectAt
}

func (line *Line) NextCollect(collectTime time.Time) {
	line.nextCollectAt = collectTime
}

func (line *Line) CollectedAt() time.Time {
	return line.collectedAt
}

func (line *Line) Updated(updateTime time.Time) {
	line.collectedAt = updateTime
}

func (line *Line) Origin() string {
	return line.origin
}

func (line *Line) SetOrigin(origin string) {
	line.origin = origin
}

func (line *Line) MarshalJSON() ([]byte, error) {
	type Alias Line
	aux := struct {
		Id            LineId
		ObjectIDs     ObjectIDs            `json:",omitempty"`
		NextCollectAt *time.Time           `json:",omitempty"`
		CollectedAt   *time.Time           `json:",omitempty"`
		Attributes    Attributes           `json:",omitempty"`
		References    map[string]Reference `json:",omitempty"`
		*Alias
	}{
		Id:    line.id,
		Alias: (*Alias)(line),
	}

	if !line.ObjectIDs().Empty() {
		aux.ObjectIDs = line.ObjectIDs()
	}
	if !line.nextCollectAt.IsZero() {
		aux.NextCollectAt = &line.nextCollectAt
	}
	if !line.collectedAt.IsZero() {
		aux.CollectedAt = &line.collectedAt
	}
	if !line.Attributes.IsEmpty() {
		aux.Attributes = line.Attributes
	}

	if !line.References.IsEmpty() {
		aux.References = line.References.GetReferences()
	}

	return json.Marshal(&aux)
}

func (line *Line) UnmarshalJSON(data []byte) error {
	type Alias Line

	aux := &struct {
		ObjectIDs  map[string]string
		References map[string]Reference
		*Alias
	}{
		Alias: (*Alias)(line),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ObjectIDs != nil {
		line.ObjectIDConsumer.objectids = NewObjectIDsFromMap(aux.ObjectIDs)
	}

	if aux.References != nil {
		line.References.SetReferences(aux.References)
	}
	return nil
}

func (line *Line) Attribute(key string) (string, bool) {
	value, present := line.Attributes[key]
	return value, present
}

func (line *Line) Reference(key string) (Reference, bool) {
	value, present := line.References.Get(key)
	return value, present
}

func (line *Line) Save() (ok bool) {
	ok = line.model.Lines().Save(line)
	return
}

type MemoryLines struct {
	uuid.UUIDConsumer

	model Model

	mutex        *sync.RWMutex
	byIdentifier map[LineId]*Line
	byObjectId   *ObjectIdIndex
}

type Lines interface {
	uuid.UUIDInterface

	New() Line
	Find(id LineId) (Line, bool)
	FindByObjectId(objectid ObjectID) (Line, bool)
	FindAll() []Line
	Save(line *Line) bool
	Delete(line *Line) bool
}

func NewMemoryLines() *MemoryLines {
	return &MemoryLines{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[LineId]*Line),
		byObjectId:   NewObjectIdIndex(),
	}
}

func (manager *MemoryLines) Clone(model Model) *MemoryLines {
	clone := NewMemoryLines()
	clone.model = model

	for _, line := range manager.byIdentifier {
		cloneLine := *line
		clone.Save(&cloneLine)
	}

	return clone
}

func (manager *MemoryLines) New() Line {
	line := NewLine(manager.model)
	return *line
}

func (manager *MemoryLines) Find(id LineId) (Line, bool) {
	manager.mutex.RLock()

	line, ok := manager.byIdentifier[id]
	if ok {
		manager.mutex.RUnlock()
		return *line, true
	} else {
		manager.mutex.RUnlock()
		return Line{}, false
	}
}

func (manager *MemoryLines) FindByObjectId(objectid ObjectID) (Line, bool) {
	manager.mutex.RLock()

	id, ok := manager.byObjectId.Find(objectid)
	if ok {
		manager.mutex.RUnlock()
		return *manager.byIdentifier[LineId(id)], true
	}

	manager.mutex.RUnlock()
	return Line{}, false
}

func (manager *MemoryLines) FindAll() (lines []Line) {
	manager.mutex.RLock()

	for _, line := range manager.byIdentifier {
		lines = append(lines, *line)
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryLines) Save(line *Line) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if line.Id() == "" {
		line.id = LineId(manager.NewUUID())
	}

	line.model = manager.model
	manager.byIdentifier[line.Id()] = line
	manager.byObjectId.Index(line)

	return true
}

func (manager *MemoryLines) Delete(line *Line) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, line.Id())
	manager.byObjectId.Delete(ModelId(line.id))

	return true
}

func (manager *MemoryLines) Load(referentialSlug string) error {
	var selectLines []SelectLine
	modelName := manager.model.Date()
	sqlQuery := fmt.Sprintf("select * from lines where referential_slug = '%s' and model_name = '%s'", referentialSlug, modelName.String())
	_, err := Database.Select(&selectLines, sqlQuery)
	if err != nil {
		return err
	}
	for _, sl := range selectLines {
		line := manager.New()
		line.id = LineId(sl.Id)
		if sl.Name.Valid {
			line.Name = sl.Name.String
		}
		if sl.CollectGeneralMessages.Valid {
			line.CollectGeneralMessages = sl.CollectGeneralMessages.Bool
		}

		if sl.Attributes.Valid && len(sl.Attributes.String) > 0 {
			if err = json.Unmarshal([]byte(sl.Attributes.String), &line.Attributes); err != nil {
				return err
			}
		}

		if sl.References.Valid && len(sl.References.String) > 0 {
			references := make(map[string]Reference)
			if err = json.Unmarshal([]byte(sl.References.String), &references); err != nil {
				return err
			}
			line.References.SetReferences(references)
		}

		if sl.ObjectIDs.Valid && len(sl.ObjectIDs.String) > 0 {
			objectIdMap := make(map[string]string)
			if err = json.Unmarshal([]byte(sl.ObjectIDs.String), &objectIdMap); err != nil {
				return err
			}
			line.objectids = NewObjectIDsFromMap(objectIdMap)
		}

		manager.Save(&line)
	}
	return nil
}
