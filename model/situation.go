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

const (
	SituationReportTypeGeneral  ReportType    = "general"
	SituationReportTypeIncident ReportType    = "incident"
	SituationTypeLine           SituationType = "Line"
	SituationTypeStopArea       SituationType = "StopArea"
)

type ReportType string
type SituationType string

type SituationTranslatedString struct {
	DefaultValue string            `json:",omitempty"`
	Translations map[string]string `json:",omitempty"`
}

type Situation struct {
	model Model
	CodeConsumer
	id     SituationId
	Origin string

	RecordedAt time.Time
	Version    int `json:",omitempty"`

	VersionedAt     time.Time
	ValidityPeriods []*TimeRange

	Keywords   []string   `json:",omitempty"`
	ReportType ReportType `json:",omitempty"`

	ProducerRef string `json:",omitempty"`
	Format      string `json:",omitempty"`

	ParticipantRef string                     `json:",omitempty"`
	Summary        *SituationTranslatedString `json:",omitempty"`
	Description    *SituationTranslatedString `json:",omitempty"`

	Affects []Affect `json:",omitempty"`
}

// SubTypes of Affect
type Affect interface {
	GetType() SituationType
	GetId() ModelId
}

type AffectedStopArea struct {
	StopAreaId StopAreaId `json:",omitempty"`
}

func (a AffectedStopArea) GetId() ModelId {
	return ModelId(a.StopAreaId)
}

func (a AffectedStopArea) GetType() SituationType {
	return SituationTypeStopArea
}

func NewAffectedStopArea() *AffectedStopArea {
	return &AffectedStopArea{}
}

type AffectedLine struct {
	LineId               LineId                 `json:",omitempty"`
	AffectedDestinations []*AffectedDestination `json:",omitempty"`
	AffectedSections     []*AffectedSection     `json:",omitempty"`
	AffectedRoutes       []*AffectedRoute       `json:",omitempty"`
}

type AffectedDestination struct {
	StopAreaId StopAreaId
}

type AffectedSection struct {
	FirstStop StopAreaId
	LastStop  StopAreaId
}

type AffectedRoute struct {
	RouteRef string
}

func (a AffectedLine) GetId() ModelId {
	return ModelId(a.LineId)
}

func (a AffectedLine) GetType() SituationType {
	return SituationTypeLine
}

func NewAffectedLine() *AffectedLine {
	return &AffectedLine{}
}

type TimeRange struct {
	StartTime time.Time `json:",omitempty"`
	EndTime   time.Time `json:",omitempty"`
}

func NewSituation(model Model) *Situation {
	situation := &Situation{
		model: model,
	}

	situation.codes = make(Codes)
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
		Codes   map[string]string
		Affects []json.RawMessage
		*Alias
	}{
		Alias: (*Alias)(situation),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}
	if aux.Affects != nil {
		for _, v := range aux.Affects {
			var affectedSubtype = struct {
				Type SituationType
			}{}
			err = json.Unmarshal(v, &affectedSubtype)
			if err != nil {
				return err
			}
			switch affectedSubtype.Type {
			case "StopArea":
				a := NewAffectedStopArea()
				json.Unmarshal(v, a)
				situation.Affects = append(situation.Affects, a)
			case "Line":
				l := NewAffectedLine()
				json.Unmarshal(v, l)
				situation.Affects = append(situation.Affects, l)
			}
		}
	}
	if aux.Codes != nil {
		situation.CodeConsumer.codes = NewCodesFromMap(aux.Codes)
	}
	return nil
}

func (affect AffectedStopArea) MarshalJSON() ([]byte, error) {
	type Alias AffectedStopArea
	aux := struct {
		Type SituationType
		Alias
	}{
		Type:  SituationTypeStopArea,
		Alias: (Alias)(affect),
	}

	return json.Marshal(&aux)
}

func (affect AffectedLine) MarshalJSON() ([]byte, error) {
	type Alias AffectedLine
	aux := struct {
		Type SituationType
		Alias
	}{
		Type:  SituationTypeLine,
		Alias: (Alias)(affect),
	}

	return json.Marshal(&aux)
}

func (situation *Situation) MarshalJSON() ([]byte, error) {
	type Alias Situation
	aux := struct {
		Codes       Codes      `json:",omitempty"`
		RecordedAt  *time.Time `json:",omitempty"`
		VersionedAt *time.Time `json:",omitempty"`
		*Alias
		Id      SituationId
		Affects []Affect `json:",omitempty"`
	}{
		Id:    situation.id,
		Alias: (*Alias)(situation),
	}

	if !situation.Codes().Empty() {
		aux.Codes = situation.Codes()
	}

	if !situation.RecordedAt.IsZero() {
		aux.RecordedAt = &situation.RecordedAt
	}

	if !situation.VersionedAt.IsZero() {
		aux.VersionedAt = &situation.VersionedAt
	}

	if len(situation.Affects) != 0 {
		aux.Affects = situation.Affects
	}

	return json.Marshal(&aux)
}

func (situation *Situation) GMValidUntil() time.Time {
	if len(situation.ValidityPeriods) == 0 {
		return time.Time{}
	}
	return situation.ValidityPeriods[0].EndTime
}

func (situation *Situation) GetGMChannel() (string, bool) {
	switch {
	case situation.containsKeyword("Perturbation"):
		return "Perturbation", true
	case situation.containsKeyword("Information"):
		return "Information", true
	case situation.containsKeyword("Commercial"):
		return "Commercial", true
	default:
		return "", false
	}
}

func (situation *Situation) containsKeyword(str string) bool {
	if len(situation.Keywords) == 0 {
		return false
	}
	for _, v := range situation.Keywords {
		if v == str {
			return true
		}
	}
	return false
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
	FindByCode(code Code) (Situation, bool)
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

func (manager *MemorySituations) FindByCode(code Code) (Situation, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	for _, situation := range manager.byIdentifier {
		situationCode, _ := situation.Code(code.CodeSpace())
		if situationCode.Value() == code.Value() {
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
