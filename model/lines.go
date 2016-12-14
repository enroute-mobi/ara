package model

import "encoding/json"

type LineId string

type Line struct {
	ObjectIDConsumer
	model Model

	id LineId
}

func NewLine(model Model) *Line {
	line := &Line{model: model}
	line.objectids = make(ObjectIDs)
	return line
}

func (line *Line) Id() LineId {
	return line.id
}

// WIP
func (line *Line) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Id": line.id,
	})
}

func (line *Line) UnmarshalJSON(data []byte) error {
	aux := &struct {
		ObjectIDs ObjectIDs
	}{
		ObjectIDs: make(ObjectIDs),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	line.ObjectIDConsumer.objectids = aux.ObjectIDs

	return nil
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
