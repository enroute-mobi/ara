package model

import "encoding/json"

type LineId string

type LineAttributes struct {
	ObjectId ObjectID
	Name     string
}

type Line struct {
	ObjectIDConsumer
	model Model

	id LineId

	Name       string
	Attributes map[string]string
	References map[string]Reference
}

func NewLine(model Model) *Line {
	line := &Line{
		model:      model,
		Attributes: make(map[string]string),
		References: make(map[string]Reference),
	}

	line.objectids = make(ObjectIDs)
	return line
}

func (line *Line) Id() LineId {
	return line.id
}

func (line *Line) MarshalJSON() ([]byte, error) {
	lineMap := map[string]interface{}{
		"Id":         line.id,
		"Name":       line.Name,
		"Attributes": line.Attributes,
		"References": line.References,
	}
	if line.ObjectIDs() != nil {
		lineMap["ObjectIDs"] = line.ObjectIDs()
	}
	return json.Marshal(lineMap)
}

func (line *Line) UnmarshalJSON(data []byte) error {
	type Alias Line

	aux := &struct {
		ObjectIDs map[string]string
		Reference map[string]Reference
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

	return nil
}

func (line *Line) Attribute(key string) (string, bool) {
	value, present := line.Attributes[key]
	return value, present
}

func (line *Line) Reference(key string) (Reference, bool) {
	value, present := line.References[key]
	return value, present
}

func (line *Line) Save() (ok bool) {
	ok = line.model.Lines().Save(line)
	return
}

type MemoryLines struct {
	UUIDConsumer

	model Model

	byIdentifier map[LineId]*Line
	byObjectId   map[string]map[string]LineId
}

type Lines interface {
	UUIDInterface

	New() Line
	Find(id LineId) (Line, bool)
	FindByObjectId(objectid ObjectID) (Line, bool)
	FindAll() []Line
	Save(line *Line) bool
	Delete(line *Line) bool
}

func NewMemoryLines() *MemoryLines {
	return &MemoryLines{
		byIdentifier: make(map[LineId]*Line),
		byObjectId:   make(map[string]map[string]LineId),
	}
}

func (manager *MemoryLines) New() Line {
	line := NewLine(manager.model)
	return *line
}

func (manager *MemoryLines) Find(id LineId) (Line, bool) {
	line, ok := manager.byIdentifier[id]
	if ok {
		return *line, true
	} else {
		return Line{}, false
	}
}

func (manager *MemoryLines) FindByObjectId(objectid ObjectID) (Line, bool) {
	valueMap, ok := manager.byObjectId[objectid.Kind()]
	if !ok {
		return Line{}, false
	}
	id, ok := valueMap[objectid.Value()]
	if !ok {
		return Line{}, false
	}
	return *manager.byIdentifier[id], true
}

func (manager *MemoryLines) FindAll() (lines []Line) {
	if len(manager.byIdentifier) == 0 {
		return []Line{}
	}
	for _, line := range manager.byIdentifier {
		lines = append(lines, *line)
	}
	return
}

func (manager *MemoryLines) Save(line *Line) bool {
	if line.Id() == "" {
		line.id = LineId(manager.NewUUID())
	}
	line.model = manager.model
	manager.byIdentifier[line.Id()] = line
	for _, objectid := range line.ObjectIDs() {
		_, ok := manager.byObjectId[objectid.Kind()]
		if !ok {
			manager.byObjectId[objectid.Kind()] = make(map[string]LineId)
		}
		manager.byObjectId[objectid.Kind()][objectid.Value()] = line.Id()
	}
	return true
}

func (manager *MemoryLines) Delete(line *Line) bool {
	delete(manager.byIdentifier, line.Id())
	for _, objectid := range line.ObjectIDs() {
		valueMap := manager.byObjectId[objectid.Kind()]
		delete(valueMap, objectid.Value())
	}
	return true
}
