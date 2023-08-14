package model

import (
	"encoding/json"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/uuid"
)

type SituationId string

type Message struct {
	Content             string `json:"MessageText,omitempty"`
	Type                string `json:"MessageType,omitempty"`
	NumberOfLines       int    `json:",omitempty"`
	NumberOfCharPerLine int    `json:",omitempty"`
}

type Situation struct {
	model Model
	ObjectIDConsumer
	id     SituationId
	Origin string

	RecordedAt time.Time
	Version    int `json:",omitempty"`

	ValidityPeriods []*TimeRange `json:",omitempty"`

	ProducerRef string `json:",omitempty"`
	Channel     string `json:",omitempty"`
	Format      string `json:",omitempty"`

	Messages     []*Message
	LineSections []*References
	References   []*Reference
}

type TimeRange struct {
	StartTime time.Time `json:",omitempty"`
	EndTime   time.Time `json:",omitempty"`
}

func NewSituation(model Model) *Situation {
	situation := &Situation{
		model: model,
	}

	situation.objectids = make(ObjectIDs)
	return situation
}

func (situation *Situation) Id() SituationId {
	return situation.id
}

func (situation *Situation) FindReferenceByObjectId(obj *ObjectID) (*Reference, bool) {
	for _, ref := range situation.References {
		if ref.ObjectId.String() == obj.String() {
			return ref, true
		}
	}

	return &Reference{}, false
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
		ObjectIDs  ObjectIDs  `json:",omitempty"`
		RecordedAt *time.Time `json:",omitempty"`
		*Alias
		Id              SituationId
		ValidityPeriods []*TimeRange  `json:",omitempty"`
		Messages        []*Message    `json:",omitempty"`
		References      []*Reference  `json:",omitempty"`
		LineSections    []*References `json:",omitempty"`
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
	if len(situation.References) != 0 {
		aux.References = situation.References
	}
	if len(situation.LineSections) != 0 {
		aux.LineSections = situation.LineSections
	}
	if !situation.RecordedAt.IsZero() {
		aux.RecordedAt = &situation.RecordedAt
	}
	if len(situation.ValidityPeriods) != 0 {
		aux.ValidityPeriods = situation.ValidityPeriods
	}

	return json.Marshal(&aux)
}

func (situation *Situation) GMValidUntil() time.Time {
	if len(situation.ValidityPeriods) == 0 {
		return time.Time{}
	}
	return situation.ValidityPeriods[0].EndTime
}

type MemorySituations struct {
	uuid.UUIDConsumer

	model *MemoryModel

	mutex          *sync.RWMutex
	broadcastEvent func(event GeneralMessageBroadcastEvent)
	byIdentifier   map[SituationId]*Situation
}

type Situations interface {
	uuid.UUIDInterface

	New() Situation
	Find(id SituationId) (Situation, bool)
	FindByObjectId(objectid ObjectID) (Situation, bool)
	FindAll() []Situation
	Save(situation *Situation) bool
	Delete(situation *Situation) bool
}

func NewMemorySituations() *MemorySituations {
	return &MemorySituations{
		mutex:        &sync.RWMutex{},
		byIdentifier: make(map[SituationId]*Situation),
	}
}

func (manager *MemorySituations) New() Situation {
	situation := NewSituation(manager.model)
	return *situation
}

func (manager *MemorySituations) Find(id SituationId) (Situation, bool) {
	manager.mutex.RLock()
	situation, ok := manager.byIdentifier[id]
	manager.mutex.RUnlock()

	if ok {
		return *situation, true
	}
	return Situation{}, false
}

func (manager *MemorySituations) FindAll() (situations []Situation) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	if len(manager.byIdentifier) == 0 {
		return []Situation{}
	}
	for _, situation := range manager.byIdentifier {
		situations = append(situations, *situation)
	}
	return
}

func (manager *MemorySituations) FindByObjectId(objectid ObjectID) (Situation, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, situation := range manager.byIdentifier {
		situationObjectId, _ := situation.ObjectID(objectid.Kind())
		if situationObjectId.Value() == objectid.Value() {
			return *situation, true
		}
	}
	return Situation{}, false
}

func (manager *MemorySituations) Save(situation *Situation) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	if situation.Id() == "" {
		situation.id = SituationId(manager.NewUUID())
	}
	situation.model = manager.model
	manager.byIdentifier[situation.Id()] = situation

	event := GeneralMessageBroadcastEvent{
		SituationId: situation.id,
	}

	if manager.broadcastEvent != nil {
		manager.broadcastEvent(event)
	}
	return true
}

func (manager *MemorySituations) Delete(situation *Situation) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, situation.Id())
	return true
}
