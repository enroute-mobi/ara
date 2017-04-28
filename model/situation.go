package model

import (
	"encoding/json"
	"time"
)

type SituationId string

type Message struct {
	Content             string `xml:"MessageText" json:"MessageText,omitempty"`
	Type                string `xml:"MessageType" json:"MessageType,omitempty"`
	NumberOfLines       int    `xml:"NumberOfLines" json:",omitempty"`
	NumberOfCharPerLine int    `xml:"NumberOfCharPerLine" json:",omitempty"`
}

type Situation struct {
	ObjectIDConsumer

	model Model

	id SituationId

	References References
	Messages   []*Message

	RecordedAt  time.Time
	ValidUntil  time.Time
	Format      string `json:",omitempty"`
	Channel     string `json:",omitempty"`
	ProducerRef string `json:",omitempty"`
	Version     int64  `json:",omitempty"`
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

func (situation *Situation) MarshalJSON() ([]byte, error) {
	type Alias Situation
	aux := struct {
		Id         SituationId
		ObjectIDs  ObjectIDs  `json:",omitempty"`
		RecordedAt *time.Time `json:",omitempty"`
		ValidUntil *time.Time `json:",omitempty"`
		Messages   []*Message `json:",omitempty"`
		References References `json:",omitempty"`
		*Alias
	}{
		Id:    situation.id,
		Alias: (*Alias)(situation),
	}

	if !situation.ObjectIDs().Empty() {
		aux.ObjectIDs = situation.ObjectIDs()
	}
	if len(situation.Messages) != 0 {
		aux.Messages = situation.Messages
	}
	if !situation.References.IsEmpty() {
		aux.References = situation.References
	}
	if !situation.RecordedAt.IsZero() {
		aux.RecordedAt = &situation.RecordedAt
	}
	if !situation.ValidUntil.IsZero() {
		aux.ValidUntil = &situation.ValidUntil
	}

	return json.Marshal(&aux)
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
