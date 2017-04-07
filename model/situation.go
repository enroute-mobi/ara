package model

import (
	"encoding/json"
	"time"
)

type SituationId string

type Message struct {
	Content             string `xml:"MessageText,chardata"`
	Type                string `xml:"MessageType,chardata"`
	NumberOfLines       int    `xml:"NumberOfLines,attr"`
	NumberOfCharPerLine int    `xml:"NumberOfCharPerLine,attr"`
}

type Situation struct {
	ObjectIDConsumer

	model Model

	id SituationId

	References References
	Messages   []*Message

	RecordedAt time.Time
	ValidUntil time.Time
	Format     string
	Channel    string
	Version    int64
}

func NewSituation(model Model) *Situation {
	situation := &Situation{
		model:      model,
		References: NewReferences(),
	}

	situation.objectids = make(ObjectIDs)
	return situation
}

func (situation *Situation) Id() SituationId {
	return situation.id
}

func (situation *Situation) Save() (ok bool) {
	ok = situation.model.Situations().Save(situation)
	return
}

func (situation *Situation) UnmarshalJSON(data []byte) error {
	type Alias Situation

	aux := &struct {
		ObjectIDs map[string]string
		*Alias
	}{
		Alias: (*Alias)(situation),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.ObjectIDs != nil {
		situation.ObjectIDConsumer.objectids = NewObjectIDsFromMap(aux.ObjectIDs)
	}
	return nil
}

func (situation *Situation) fillSituation(situationMap map[string]interface{}) {
	if situation.id != "" {
		situationMap["Id"] = situation.id
	}

	if len(situation.Messages) != 0 {
		situationMap["Message"] = situation.Messages
	}

	if !situation.References.IsEmpty() {
		situationMap["References"] = situation.References
	}

}

func (situation *Situation) MarshalJSON() ([]byte, error) {
	situationMap := make(map[string]interface{})

	if !situation.ObjectIDs().Empty() {
		situationMap["ObjectIDs"] = situation.ObjectIDs()
	}

	situation.fillSituation(situationMap)
	return json.Marshal(situationMap)
}

type MemorySituations struct {
	UUIDConsumer

	model *MemoryModel

	byIdentifier map[SituationId]*Situation
}

type Situations interface {
	UUIDInterface

	New() Situation
	Find(id SituationId) (Situation, bool)
	FindByObjectId(objectid ObjectID) (Situation, bool)
	FindAll() []Situation
	Save(situation *Situation) bool
	Delete(situation *Situation) bool
}

func NewMemorySituations() *MemorySituations {
	return &MemorySituations{
		byIdentifier: make(map[SituationId]*Situation),
	}
}

func (manager *MemorySituations) New() Situation {
	situation := NewSituation(manager.model)
	return *situation
}

func (manager *MemorySituations) Find(id SituationId) (Situation, bool) {
	situation, ok := manager.byIdentifier[id]
	if ok {
		return *situation, true
	} else {
		return Situation{}, false
	}
}

func (manager *MemorySituations) FindAll() (situations []Situation) {
	if len(manager.byIdentifier) == 0 {
		return []Situation{}
	}
	for _, situation := range manager.byIdentifier {
		situations = append(situations, *situation)
	}
	return
}

func (manager *MemorySituations) FindByObjectId(objectid ObjectID) (Situation, bool) {
	for _, situation := range manager.byIdentifier {
		situationObjectId, _ := situation.ObjectID(objectid.Kind())
		if situationObjectId.Value() == objectid.Value() {
			return *situation, true
		}
	}
	return Situation{}, false
}

func (manager *MemorySituations) Save(situation *Situation) bool {
	if situation.Id() == "" {
		situation.id = SituationId(manager.NewUUID())
	}
	situation.model = manager.model
	manager.byIdentifier[situation.Id()] = situation
	return true
}

func (manager *MemorySituations) Delete(situation *Situation) bool {
	delete(manager.byIdentifier, situation.Id())
	return true
}