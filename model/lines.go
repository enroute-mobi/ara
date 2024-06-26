package model

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

type LineId ModelId

type Line struct {
	Collectable
	model      Model
	References References
	CodeConsumer
	Attributes        Attributes
	id                LineId
	ReferentId        LineId `json:",omitempty"`
	Name              string `json:",omitempty"`
	Number            string `json:",omitempty"`
	origin            string
	CollectSituations bool
}

func NewLine(model Model) *Line {
	line := &Line{
		model:      model,
		Attributes: NewAttributes(),
		References: NewReferences(),
	}

	line.codes = make(Codes)
	return line
}

func (line *Line) modelId() ModelId {
	return ModelId(line.id)
}

func (line *Line) copy() *Line {
	l := *line
	l.Attributes = line.Attributes.Copy()
	l.References = line.References.Copy()
	return &l
}

func (line *Line) Id() LineId {
	return line.id
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
		*Alias
		Codes         Codes                `json:",omitempty"`
		NextCollectAt *time.Time           `json:",omitempty"`
		CollectedAt   *time.Time           `json:",omitempty"`
		Attributes    Attributes           `json:",omitempty"`
		References    map[string]Reference `json:",omitempty"`
		Id            LineId
	}{
		Id:    line.id,
		Alias: (*Alias)(line),
	}

	if !line.Codes().Empty() {
		aux.Codes = line.Codes()
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
		Codes      map[string]string
		References map[string]Reference
		*Alias
	}{
		Alias: (*Alias)(line),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Codes != nil {
		line.CodeConsumer.codes = NewCodesFromMap(aux.Codes)
	}

	if aux.References != nil {
		line.References.SetReferences(aux.References)
	}
	return nil
}

func (line *Line) Referent() (*Line, bool) {
	return line.model.Lines().Find(line.ReferentId)
}

func (line *Line) ReferentOrSelfCode(codeSpace string) (Code, bool) {
	ref, ok := line.Referent()
	if ok {
		code, ok := ref.Code(codeSpace)
		if ok {
			return code, true
		}
	}
	code, ok := line.Code(codeSpace)
	if ok {
		return code, true
	}
	return Code{}, false
}

/*
Returns true if we need to send the Line in a Discovery

We only send the Line if it has no referent with a correct codeSpace.
If that's the case, we'll send the Referent instead
*/
func (line *Line) DiscoveryCode(codeSpace string) (Code, bool) {
	ref, ok := line.Referent()
	if ok {
		_, ok := ref.Code(codeSpace)
		if ok {
			return Code{}, false
		}
	}
	code, ok := line.Code(codeSpace)
	if ok {
		return code, true
	}
	return Code{}, false
}

func (line *Line) Attribute(key string) (string, bool) {
	value, present := line.Attributes[key]
	return value, present
}

func (line *Line) Reference(key string) (Reference, bool) {
	value, present := line.References.Get(key)
	return value, present
}

func (line *Line) Save() bool {
	return line.model.Lines().Save(line)
}

type MemoryLines struct {
	uuid.UUIDConsumer

	model Model

	mutex        *sync.RWMutex
	byIdentifier map[LineId]*Line
	byReferent   *Index
	byCode       *CodeIndex
}

type Lines interface {
	uuid.UUIDInterface

	New() *Line
	Find(LineId) (*Line, bool)
	FindByCode(Code) (*Line, bool)
	FindAll() []*Line
	FindFamily(LineId) []LineId
	FindFamilyFromCode(Code) []LineId
	Save(*Line) bool
	Delete(*Line) bool
}

func NewMemoryLines() *MemoryLines {
	referentExtractor := func(instance ModelInstance) ModelId { return ModelId((instance.(*Line)).ReferentId) }

	return &MemoryLines{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[LineId]*Line),
		byReferent:   NewIndex(referentExtractor),
		byCode:       NewCodeIndex(),
	}
}

func (manager *MemoryLines) New() *Line {
	return NewLine(manager.model)
}

func (manager *MemoryLines) Find(id LineId) (*Line, bool) {
	manager.mutex.RLock()
	line, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return line.copy(), true
	}
	return &Line{}, false
}

func (manager *MemoryLines) FindByReferentId(id LineId) (lines []*Line) {
	manager.mutex.RLock()

	ids, _ := manager.byReferent.Find(ModelId(id))

	for _, id := range ids {
		l := manager.byIdentifier[LineId(id)]
		lines = append(lines, l.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryLines) FindByCode(code Code) (*Line, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byCode.Find(code)
	if ok {
		return manager.byIdentifier[LineId(id)].copy(), true
	}

	return &Line{}, false
}

func (manager *MemoryLines) FindAll() (lines []*Line) {
	manager.mutex.RLock()

	for _, line := range manager.byIdentifier {
		lines = append(lines, line.copy())
	}

	manager.mutex.RUnlock()
	return
}

func (manager *MemoryLines) FindFamily(lineId LineId) (lineIds []LineId) {
	manager.mutex.RLock()

	lineIds = manager.findFamily(lineId)

	manager.mutex.RUnlock()

	return
}

func (manager *MemoryLines) FindFamilyFromCode(code Code) (lineIds []LineId) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.byCode.Find(code)
	if !ok {
		return
	}

	lineIds = manager.findFamily(LineId(id))

	return
}

func (manager *MemoryLines) findFamily(lineId LineId) (lineIds []LineId) {
	lineIds = []LineId{lineId}

	ids, _ := manager.byReferent.Find(ModelId(lineId))
	for _, id := range ids {
		lineIds = append(lineIds, manager.findFamily(LineId(id))...)
	}

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
	manager.byReferent.Index(line)
	manager.byCode.Index(line)

	return true
}

func (manager *MemoryLines) Delete(line *Line) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, line.Id())
	manager.byReferent.Delete(ModelId(line.id))
	manager.byCode.Delete(ModelId(line.id))

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
		if sl.ReferentId.Valid {
			line.ReferentId = LineId(sl.ReferentId.String)
		}

		if sl.CollectSituations.Valid {
			line.CollectSituations = sl.CollectSituations.Bool
		}

		if sl.Number.Valid {
			line.Number = sl.Number.String
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

		if sl.Codes.Valid && len(sl.Codes.String) > 0 {
			codeMap := make(map[string]string)
			if err = json.Unmarshal([]byte(sl.Codes.String), &codeMap); err != nil {
				return err
			}
			line.codes = NewCodesFromMap(codeMap)
		}

		manager.Save(line)
	}
	return nil
}
