package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

var lineReferentExtractor = func(instance ModelInstance) string { return string((instance.(*Line)).ReferentId) }

type memoryLines struct {
	IndexHandler
	memoryManager

	mutex        *sync.RWMutex
	byIdentifier map[LineId]*Line
}

type Lines interface {
	ModelManager[LineId, *Line]
	CodeHandler[*Line]
	Loadable

	FindFamily(LineId) []LineId
	FindFamilyFromCode(Code) []LineId
	FindCode(LineId, string) (Code, bool)
	CollectableLines(time.Time) []*Line
}

func NewMemoryLines() Lines {
	m := &memoryLines{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[LineId]*Line),
	}
	m.InitIndexes()
	m.AddIndex(ByReferent, OneToMany, lineReferentExtractor)

	return m
}

func (manager *memoryLines) New() *Line {
	return NewLine(manager.model)
}

func (manager *memoryLines) Find(id LineId) (*Line, bool) {
	manager.mutex.RLock()
	line, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return line.copy(), true
	}
	return &Line{}, false
}

func (manager *memoryLines) FindByReferentId(id LineId) (lines []*Line) {
	manager.mutex.RLock()

	ids, _ := manager.FindBy(ByReferent, string(id))

	for _, id := range ids {
		l := manager.byIdentifier[LineId(id)]
		lines = append(lines, l.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *memoryLines) FindByCode(code Code) (*Line, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.ByCode().Find(code)
	if ok {
		return manager.byIdentifier[LineId(id)].copy(), true
	}

	return &Line{}, false
}

func (manager *memoryLines) CodeExists(code Code) bool {
	manager.mutex.RLock()
	_, ok := manager.ByCode().Find(code)
	manager.mutex.RUnlock()

	return ok
}

func (manager *memoryLines) FindCode(id LineId, codespace string) (Code, bool) {
	manager.mutex.RLock()
	line, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return line.Code(codespace)
	}
	return Code{}, false
}

func (manager *memoryLines) FindAll() (lines []*Line) {
	manager.mutex.RLock()

	for _, line := range manager.byIdentifier {
		lines = append(lines, line.copy())
	}

	manager.mutex.RUnlock()
	return
}

// CollectableLines returns a copy of the lines due for collection at the given
// time. Only the due lines are copied, avoiding a deep copy of the whole
// collection on every guardian cycle.
func (manager *memoryLines) CollectableLines(now time.Time) (lines []*Line) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, line := range manager.byIdentifier {
		if !line.nextCollectAt.Before(now) {
			continue
		}
		lines = append(lines, line.copy())
	}
	return
}

func (manager *memoryLines) FindFamily(lineId LineId) (lineIds []LineId) {
	manager.mutex.RLock()

	lineIds = manager.findFamily(lineId)

	manager.mutex.RUnlock()

	return
}

func (manager *memoryLines) FindFamilyFromCode(code Code) (lineIds []LineId) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.ByCode().Find(code)
	if !ok {
		return
	}

	lineIds = manager.findFamily(LineId(id))

	return
}

func (manager *memoryLines) findFamily(lineId LineId) (lineIds []LineId) {
	lineIds = []LineId{lineId}

	ids, _ := manager.FindBy(ByReferent, string(lineId))
	for _, id := range ids {
		lineIds = append(lineIds, manager.findFamily(LineId(id))...)
	}

	return
}

func (manager *memoryLines) Save(line *Line) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if line.Id() == "" {
		line.id = LineId(manager.NewUUID())
	}

	line.model = manager.model
	manager.byIdentifier[line.Id()] = line
	manager.Index(line)

	return true
}

func (manager *memoryLines) Delete(line *Line) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, line.Id())
	manager.Deindex(string(line.id))

	return true
}

func (manager *memoryLines) Load(referentialSlug string) error {
	var selectLines []SelectLine
	modelDate := manager.model.Date()
	sqlQuery := fmt.Sprintf("select * from lines where referential_slug = '%s' and model_date = '%s'", referentialSlug, modelDate.String())
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
		if sl.ReferentId.Valid {
			line.ReferentId = LineId(sl.ReferentId.String)
		}

		if sl.CollectSituations.Valid {
			line.CollectSituations = sl.CollectSituations.Bool
		}

		if sl.Number.Valid {
			line.Number = sl.Number.String
		}
		if sl.RawAttributes.Valid && len(sl.RawAttributes.String) > 0 {
			if err = json.Unmarshal([]byte(sl.RawAttributes.String), &line.RawAttributes); err != nil {
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

		if sl.Codes.Valid && len(sl.Codes.String) > 0 {
			codeMap := make(map[string]string)
			if err = json.Unmarshal([]byte(sl.Codes.String), &codeMap); err != nil {
				return err
			}
			line.SetCodesFromMap(codeMap)
		}

		manager.Save(line)
	}
	return nil
}
