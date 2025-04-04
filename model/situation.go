package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"sync"
	"time"

	e "bitbucket.org/enroute-mobi/ara/core/apierrs"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/sym01/htmlsanitizer"
)

type SituationId ModelId

const (
	SituationReportTypeGeneral  ReportType    = "general"
	SituationReportTypeIncident ReportType    = "incident"
	SituationTypeLine           SituationType = "Line"
	SituationTypeStopArea       SituationType = "StopArea"
	SituationTypeAllLines       SituationType = "AllLines"
)

type ReportType string
type SituationType string

type TranslatedString struct {
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

	VersionedAt        time.Time
	ValidityPeriods    []*TimeRange
	PublicationWindows []*TimeRange

	Progress       SituationProgress   `json:",omitempty"`
	Severity       SituationSeverity   `json:",omitempty"`
	Keywords       []string            `json:",omitempty"`
	ReportType     ReportType          `json:",omitempty"`
	AlertCause     SituationAlertCause `json:",omitempty"`
	Reality        SituationReality    `json:",omitempty"`
	ProducerRef    string              `json:",omitempty"`
	Format         string              `json:",omitempty"`
	InternalTags   []string            `json:",omitempty"`
	ParticipantRef string              `json:",omitempty"`
	Summary        *TranslatedString   `json:",omitempty"`
	Description    *TranslatedString   `json:",omitempty"`

	Affects                 Affects                   `json:",omitempty"`
	Consequences            []*Consequence            `json:",omitempty"`
	PublishToWebActions     []*PublishToWebAction     `json:",omitempty"`
	PublishToMobileActions  []*PublishToMobileAction  `json:",omitempty"`
	PublishToDisplayActions []*PublishToDisplayAction `json:",omitempty"`
	InfoLinks               []*InfoLink               `json:",omitempty"`
}

type ActionData struct {
	Name       string            `json:",omitempty"`
	ActionType string            `json:",omitempty"`
	Value      string            `json:",omitempty"`
	Prompt     *TranslatedString `json:",omitempty"`
}

type ActionCommon struct {
	Name               string                `json:",omitempty"`
	ActionType         string                `json:",omitempty"`
	Value              string                `json:",omitempty"`
	Prompt             *TranslatedString     `json:",omitempty"`
	ScopeType          SituationScopeType    `json:",omitempty"`
	Affects            Affects               `json:",omitempty"`
	ActionStatus       SituationActionStatus `json:",omitempty"`
	Description        *TranslatedString     `json:",omitempty"`
	PublicationWindows []*TimeRange          `json:",omitempty"`
}

type PublishToWebAction struct {
	ActionCommon

	Incidents      *bool    `json:",omitempty"`
	HomePage       *bool    `json:",omitempty"`
	Ticker         *bool    `json:",omitempty"`
	SocialNetworks []string `json:",omitempty"`
}

type PublishToMobileAction struct {
	ActionCommon

	Incidents *bool `json:",omitempty"`
	HomePage  *bool `json:",omitempty"`
}

type PublishToDisplayAction struct {
	ActionCommon

	OnPlace *bool `json:",omitempty"`
	OnBoard *bool `json:",omitempty"`
}

type InfoLink struct {
	Uri         string
	Label       string
	ImageRef    string
	LinkContent string
}
type Consequence struct {
	Periods   []*TimeRange       `json:",omitempty"`
	Condition SituationCondition `json:",omitempty"`
	Severity  SituationSeverity  `json:",omitempty"`
	Affects   Affects            `json:",omitempty"`
	Blocking  *Blocking          `json:",omitempty"`
}

type Blocking struct {
	JourneyPlanner bool
	RealTime       bool
}

// SubTypes of Affect
type Affect interface {
	GetType() SituationType
	GetId() ModelId
}

type Affects []Affect

type AffectedStopArea struct {
	StopAreaId StopAreaId `json:",omitempty"`
	LineIds    []LineId   `json:",omitempty"`
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

type AffectedAllLines struct{}

func (a AffectedAllLines) GetId() ModelId {
	return ModelId("")
}

func (a AffectedAllLines) GetType() SituationType {
	return SituationTypeAllLines
}

func NewAffectedAllLines() *AffectedAllLines {
	return &AffectedAllLines{}
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
	RouteRef    string       `json:",omitempty"`
	StopAreaIds []StopAreaId `json:",omitempty"`
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
	StartTime time.Time
	EndTime   time.Time
}

func (s *Situation) BroadcastPeriod() *TimeRange {
	var possibleMin []time.Time
	var possibleMax []time.Time
	for _, validityPeriod := range s.ValidityPeriods {
		possibleMin = append(possibleMin, validityPeriod.StartTime)
		possibleMax = append(possibleMax, validityPeriod.EndTime)
	}

	for _, publicationWindow := range s.PublicationWindows {
		possibleMin = append(possibleMin, publicationWindow.StartTime)
		possibleMax = append(possibleMax, publicationWindow.EndTime)
	}

	broadcastPeriod := &TimeRange{}
	broadcastPeriod.StartTime = slices.MinFunc(possibleMin, time.Time.Compare)

	for i := range possibleMax {
		if possibleMax[i].IsZero() {
			return broadcastPeriod
		}
	}
	broadcastPeriod.EndTime = slices.MaxFunc(possibleMax, time.Time.Compare)

	return broadcastPeriod
}

func (t *TimeRange) Overlaps(other *TimeRange) bool {
	if other.EndTime.IsZero() && t.EndTime.IsZero() {
		return true
	}

	if other.EndTime.IsZero() {
		return other.StartTime.Before(t.EndTime)
	}

	if t.EndTime.IsZero() {
		return t.StartTime.Before(other.EndTime)
	}

	return t.EndTime.After(other.StartTime) && t.StartTime.Before(other.EndTime)
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

func (affects *Affects) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}

	*affects = make(Affects, len(raw))
	for i, v := range raw {
		var affectedSubtype = struct {
			Type SituationType
		}{}
		err = json.Unmarshal(v, &affectedSubtype)
		if err != nil {
			return err
		}
		switch affectedSubtype.Type {
		case SituationTypeStopArea:
			a := NewAffectedStopArea()
			json.Unmarshal(v, a)
			(*affects)[i] = a
		case SituationTypeLine:
			l := NewAffectedLine()
			json.Unmarshal(v, l)
			(*affects)[i] = l
		case SituationTypeAllLines:
			all := NewAffectedAllLines()
			json.Unmarshal(v, all)
			(*affects)[i] = all
		}
	}
	return nil
}

func (apiSituation *APISituation) MarshalJSON() ([]byte, error) {
	type Alias APISituation
	aux := struct {
		*Alias
		Codes  Codes    `json:",omitempty"`
		Errors e.Errors `json:"Errors,omitempty"`
	}{
		Alias: (*Alias)(apiSituation),
	}

	if !apiSituation.Codes().Empty() {
		aux.Codes = apiSituation.Codes()
	}

	if !apiSituation.Errors.Empty() {
		aux.Errors = apiSituation.Errors
	}
	return json.Marshal(&aux)
}

func (situation *APISituation) UnmarshalJSON(data []byte) error {
	type Alias APISituation

	aux := &struct {
		Codes map[string]string
		*Alias
		Summary     *TranslatedString
		Description *TranslatedString
	}{
		Alias: (*Alias)(situation),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Codes != nil {
		situation.CodeConsumer.codes = NewCodesFromMap(aux.Codes)
	}

	situation.Summary = aux.Summary
	situation.Description = aux.Description

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

func (affect AffectedAllLines) MarshalJSON() ([]byte, error) {
		type Alias AffectedAllLines
	aux := struct {
		Type SituationType
		Alias
	}{
		Type:  SituationTypeAllLines,
		Alias: (Alias)(affect),
	}

	return json.Marshal(&aux)
}

func (t *TimeRange) UnmarshalJSON(data []byte) error {
	aux := &struct {
		StartTime *time.Time
		EndTime   *time.Time
	}{}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.StartTime == nil {
		t.StartTime = time.Time{}
	} else {
		t.StartTime = *aux.StartTime
	}

	if aux.EndTime == nil {
		t.EndTime = time.Time{}
	} else {
		t.EndTime = *aux.EndTime
	}

	return nil
}

func (t *TimeRange) MarshalJSON() ([]byte, error) {
	type Alias TimeRange
	aux := struct {
		StartTime *time.Time `json:",omitempty"`
		EndTime   *time.Time `json:",omitempty"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	if !t.StartTime.IsZero() {
		aux.StartTime = &t.StartTime
	}

	if !t.EndTime.IsZero() {
		aux.EndTime = &t.EndTime
	}

	return json.Marshal(&aux)
}

func (situation *Situation) modelId() ModelId {
	return ModelId(situation.id)
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

type APISituation struct {
	Id     SituationId `json:",omitempty"`
	Origin string      `json:",omitempty"`
	CodeConsumer

	CodeSpace             string    `json:",omitempty"`
	SituationNumber       string    `json:",omitempty"`
	ExistingSituationCode bool      `json:"-"`
	RecordedAt            time.Time `json:",omitempty"`
	Version               int       `json:",omitempty"`

	VersionedAt        time.Time    `json:",omitempty"`
	ValidityPeriods    []*TimeRange `json:",omitempty"`
	PublicationWindows []*TimeRange `json:",omitempty"`

	Progress       SituationProgress   `json:",omitempty"`
	Severity       SituationSeverity   `json:",omitempty"`
	Reality        SituationReality    `json:",omitempty"`
	Keywords       []string            `json:",omitempty"`
	ReportType     ReportType          `json:",omitempty"`
	AlertCause     SituationAlertCause `json:",omitempty"`
	ProducerRef    string              `json:",omitempty"`
	Format         string              `json:",omitempty"`
	InternalTags   []string            `json:",omitempty"`
	ParticipantRef string              `json:",omitempty"`
	Summary        *TranslatedString   `json:",omitempty"`
	Description    *TranslatedString   `json:",omitempty"`

	Affects                 Affects                   `json:",omitempty"`
	Consequences            []*Consequence            `json:",omitempty"`
	PublishToWebActions     []*PublishToWebAction     `json:",omitempty"`
	PublishToMobileActions  []*PublishToMobileAction  `json:",omitempty"`
	PublishToDisplayActions []*PublishToDisplayAction `json:",omitempty"`
	InfoLinks               []*InfoLink               `json:",omitempty"`

	Errors e.Errors `json:"Errors,omitempty"`

	IgnoreValidation bool `json:",omitempty"`
}

func (apiSituation *APISituation) Validate() bool {
	if apiSituation.CodeSpace == "" {
		apiSituation.Errors.Add("CodeSpace", e.ERROR_BLANK)
	}

	if apiSituation.SituationNumber == "" {
		apiSituation.Errors.Add("SituationNumber", e.ERROR_BLANK)
	}

	if apiSituation.ExistingSituationCode {
		apiSituation.Errors.Add("SituationNumber", e.ERROR_UNIQUE)
	}

	if apiSituation.Version == 0 && apiSituation.Id == "" {
		apiSituation.Errors.Add("Version", e.ERROR_BLANK)
	}

	if apiSituation.Summary == nil {
		apiSituation.Errors.Add("Summary", e.ERROR_BLANK)
	}

	if apiSituation.Summary != nil {
		if apiSituation.Summary.DefaultValue == "" {
			if apiSituation.Summary.Translations == nil {
				apiSituation.Errors.Add("Summary", e.ERROR_BLANK)
			}
		}
	}

	sanitizer := htmlsanitizer.NewHTMLSanitizer()
	if apiSituation.Summary != nil {
		sanitizedSummary, err := sanitizer.Sanitize([]byte(apiSituation.Summary.DefaultValue))
		if err != nil {
			apiSituation.Errors.Add("Summary", fmt.Sprintf("%s: %v", e.ERROR_FORMAT, err))
		} else {
			apiSituation.Summary.DefaultValue = string(sanitizedSummary)
		}
	}

	if apiSituation.Description != nil {
		sanitizedDescription, err := sanitizer.Sanitize([]byte(apiSituation.Description.DefaultValue))
		if err != nil {
			apiSituation.Errors.Add("Description", fmt.Sprintf("%s: %v", e.ERROR_FORMAT, err))
		} else {
			apiSituation.Description.DefaultValue = string(sanitizedDescription)
		}
	}

	if len(apiSituation.ValidityPeriods) == 0 {
		apiSituation.Errors.Add("ValidityPeriods", e.ERROR_BLANK)
	}

	for _, period := range apiSituation.ValidityPeriods {
		if period.StartTime.IsZero() {
			apiSituation.Errors.Add("ValidityPeriods", e.ERROR_BLANK)
			break
		}
	}

	if len(apiSituation.Affects) == 0 {
		apiSituation.Errors.Add("Affects", e.ERROR_BLANK)
	}

	if len(apiSituation.InfoLinks) != 0 {
		for i := range apiSituation.InfoLinks {
			if apiSituation.InfoLinks[i].Uri == "" {
				apiSituation.Errors.Add("InfoLinks", e.ERROR_BLANK)
				break
			}
		}
	}

	return len(apiSituation.Errors) == 0
}

func (situation *Situation) Definition() *APISituation {
	apiSituation := &APISituation{
		Id:                      situation.Id(),
		Affects:                 []Affect{},
		AlertCause:              situation.AlertCause,
		Consequences:            []*Consequence{},
		Description:             situation.Description,
		Errors:                  e.NewErrors(),
		Format:                  situation.Format,
		InternalTags:            situation.InternalTags,
		InfoLinks:               situation.InfoLinks,
		Keywords:                situation.Keywords,
		Origin:                  situation.Origin,
		ParticipantRef:          situation.ParticipantRef,
		PublishToWebActions:     []*PublishToWebAction{},
		PublishToMobileActions:  []*PublishToMobileAction{},
		PublishToDisplayActions: []*PublishToDisplayAction{},
		ProducerRef:             situation.ProducerRef,
		Progress:                situation.Progress,
		PublicationWindows:      situation.PublicationWindows,
		Reality:                 situation.Reality,
		RecordedAt:              situation.RecordedAt,
		ReportType:              situation.ReportType,
		Severity:                situation.Severity,
		Summary:                 situation.Summary,
		ValidityPeriods:         situation.ValidityPeriods,
		Version:                 situation.Version,
		VersionedAt:             situation.VersionedAt,
		IgnoreValidation:        false,
	}

	apiSituation.codes = make(Codes)
	return apiSituation
}

func (situation *Situation) SetDefinition(apiSituation *APISituation) {
	situation.Affects = apiSituation.Affects
	situation.AlertCause = apiSituation.AlertCause
	situation.Consequences = apiSituation.Consequences
	situation.Description = apiSituation.Description
	situation.Format = apiSituation.Format
	situation.InternalTags = apiSituation.InternalTags
	situation.InfoLinks = apiSituation.InfoLinks
	situation.Keywords = apiSituation.Keywords
	situation.Origin = apiSituation.Origin
	situation.ParticipantRef = apiSituation.ParticipantRef
	situation.PublishToWebActions = apiSituation.PublishToWebActions
	situation.PublishToMobileActions = apiSituation.PublishToMobileActions
	situation.PublishToDisplayActions = apiSituation.PublishToDisplayActions
	situation.ProducerRef = apiSituation.ProducerRef
	situation.Progress = apiSituation.Progress
	situation.PublicationWindows = apiSituation.PublicationWindows
	situation.Reality = apiSituation.Reality
	situation.RecordedAt = apiSituation.RecordedAt
	situation.ReportType = apiSituation.ReportType
	situation.Severity = apiSituation.Severity
	situation.Summary = apiSituation.Summary
	situation.ValidityPeriods = apiSituation.ValidityPeriods
	situation.Version = apiSituation.Version
	situation.VersionedAt = apiSituation.VersionedAt

	if apiSituation.codes.Empty() {
		if apiSituation.CodeSpace != "" && apiSituation.SituationNumber != "" {
			situation.codes = make(Codes)
			code := NewCode(apiSituation.CodeSpace, apiSituation.SituationNumber)
			situation.SetCode(code)
		}
	} else {
		// keep cucumber scenarios compatibility with API
		situation.codes = apiSituation.codes
	}
}

type MemorySituations struct {
	uuid.UUIDConsumer

	model *MemoryModel

	mutex            *sync.RWMutex
	GMbroadcastEvent func(event SituationBroadcastEvent)
	SXbroadcastEvent func(event SituationBroadcastEvent)
	byIdentifier     map[SituationId]*Situation
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

	event := SituationBroadcastEvent{
		SituationId: situation.id,
	}

	if manager.GMbroadcastEvent != nil {
		manager.GMbroadcastEvent(event)
	}

	if manager.SXbroadcastEvent != nil {
		manager.SXbroadcastEvent(event)
	}
	return true
}

func (manager *MemorySituations) Delete(situation *Situation) bool {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	delete(manager.byIdentifier, situation.Id())
	return true
}

func NewTranslatedStringFromMap(translations map[string]string) *TranslatedString {
	ts := &TranslatedString{
		Translations: make(map[string]string),
	}

	for lang, text := range translations {
		if lang == "" {
			ts.DefaultValue = text
			continue
		}
		ts.Translations[lang] = text
	}

	return ts
}

func NewTranslatedStringFromProto(value []*gtfs.TranslatedString_Translation) *TranslatedString {
	ts := &TranslatedString{
		Translations: make(map[string]string),
	}

	for _, translation := range value {
		if translation.GetLanguage() == "" {
			ts.DefaultValue = translation.GetText()
			continue
		}

		ts.Translations[translation.GetLanguage()] = translation.GetText()
	}

	return ts
}

func (ts *TranslatedString) ToProto(dest *gtfs.TranslatedString) {
	translations := []*gtfs.TranslatedString_Translation{}
	if ts.DefaultValue != "" {
		var emptyLanguage string
		gtfsTranslation := &gtfs.TranslatedString_Translation{
			Language: &emptyLanguage,
			Text:     &ts.DefaultValue,
		}

		translations = append(translations, gtfsTranslation)
	}

	for lang, text := range ts.Translations {
		gtfsTranslation := &gtfs.TranslatedString_Translation{
			Language: &lang,
			Text:     &text,
		}

		translations = append(translations, gtfsTranslation)
	}

	dest.Translation = translations
}

func (t *TimeRange) FromProto(value interface{}) error {
	var timeRange TimeRange
	switch v := value.(type) {
	case *gtfs.TimeRange:
		if v.GetStart() == 0 {
			return errors.New("gtfs.Timerange missing Start")
		}
		timeRange.StartTime = time.Unix(int64(v.GetStart()), 0)

		if v.GetEnd() != 0 {
			timeRange.EndTime = time.Unix(int64(v.GetEnd()), 0)
		}

	default:
		return fmt.Errorf("unsupported value %T", value)
	}

	*t = timeRange
	return nil
}

func (t *TimeRange) ToProto(dest interface{}) error {
	if t == nil {
		return errors.New("nil TimeRange")
	}

	switch v := dest.(type) {
	case *gtfs.TimeRange:
		if start := t.StartTime; !start.IsZero() {
			startTime := uint64(start.Unix())
			v.Start = &startTime
		}
		if end := t.EndTime; !end.IsZero() {
			endTime := uint64(end.Unix())
			v.End = &endTime
		}
	default:
		return fmt.Errorf("unsupported value %T", dest)
	}
	return nil
}

type AffectRefs struct {
	MonitoringRefs map[string]struct{}
	LineRefs       map[string]struct{}
}

func AffectFromProto(value interface{}, remoteCodeSpace string, m Model) (Affect, *AffectRefs, error) {
	collectedRefs := &AffectRefs{
		MonitoringRefs: make(map[string]struct{}),
		LineRefs:       make(map[string]struct{}),
	}

	switch v := value.(type) {
	case *gtfs.EntitySelector:
		lineId := v.GetRouteId()
		stopId := v.GetStopId()

		if stopId != "" {
			stopCode := NewCode(remoteCodeSpace, stopId)
			stopArea, ok := m.StopAreas().FindByCode(stopCode)
			if !ok {
				return nil, nil, fmt.Errorf("unknow stopId: %v", stopId)
			}
			affect := NewAffectedStopArea()
			affect.StopAreaId = stopArea.Id()
			collectedRefs.MonitoringRefs[stopId] = struct{}{}
			if lineId != "" {
				lineCode := NewCode(remoteCodeSpace, lineId)
				line, ok := m.Lines().FindByCode(lineCode)
				if ok {
					affect.LineIds = append(affect.LineIds, line.Id())
					collectedRefs.LineRefs[lineId] = struct{}{}
				}
			}
			return affect, collectedRefs, nil
		}

		if lineId != "" {
			lineCode := NewCode(remoteCodeSpace, lineId)
			line, ok := m.Lines().FindByCode(lineCode)
			if !ok {
				return nil, nil, fmt.Errorf("unknow lineId: %v", lineId)
			}
			affect := NewAffectedLine()
			affect.LineId = line.Id()
			collectedRefs.LineRefs[lineId] = struct{}{}
			return affect, collectedRefs, nil
		}
	default:
		return nil, nil, fmt.Errorf("invalide type: %T", value)
	}
	return nil, nil, errors.New("cannot find line/stopArea model from gtfs ")
}

func AffectToProto(a Affect, remoteCodeSpace string, m Model) ([]*gtfs.EntitySelector, *AffectRefs, error) {
	collectedRefs := &AffectRefs{
		MonitoringRefs: make(map[string]struct{}),
		LineRefs:       make(map[string]struct{}),
	}
	var entities []*gtfs.EntitySelector

	switch v := a.(type) {
	case *AffectedLine:
		line, ok := m.Lines().Find(v.LineId)
		if !ok {
			return nil, nil, fmt.Errorf("unknown lineId: %v", v.LineId)
		}

		lineCode, ok := line.ReferentOrSelfCode(remoteCodeSpace)
		if !ok {
			return nil, nil, fmt.Errorf("lineId %v does not have right codeSpace %v", v.LineId, remoteCodeSpace)
		}

		var routeId *string
		value := lineCode.Value()
		routeId = &value

		collectedRefs.LineRefs[lineCode.Value()] = struct{}{}
		entities = append(entities, &gtfs.EntitySelector{RouteId: routeId})

		return entities, collectedRefs, nil
	case *AffectedStopArea:
		sa, ok := m.StopAreas().Find(v.StopAreaId)
		if !ok {
			return nil, nil, fmt.Errorf("unknown stopAreaId: %v", v.StopAreaId)
		}

		saCode, ok := sa.ReferentOrSelfCode(remoteCodeSpace)
		if !ok {
			return nil, nil, fmt.Errorf("stopId %v does not have right codeSpace %v", v.StopAreaId, remoteCodeSpace)
		}

		for i := range v.LineIds {
			line, ok := m.Lines().Find(v.LineIds[i])
			if !ok {
				logger.Log.Debugf("unknown line id: %v", v.LineIds[i])
				continue
			}
			lineCode, ok := line.ReferentOrSelfCode(remoteCodeSpace)
			if !ok {
				logger.Log.Debugf("line id: %v does not have right codeSpace %v", v.LineIds[i], remoteCodeSpace)
				continue
			}
			var stopId *string
			saValue := saCode.Value()
			stopId = &saValue

			var lineId *string
			lineValue := lineCode.Value()
			lineId = &lineValue

			e := &gtfs.EntitySelector{StopId: stopId, RouteId: lineId}
			collectedRefs.LineRefs[lineCode.Value()] = struct{}{}
			collectedRefs.MonitoringRefs[*stopId] = struct{}{}
			entities = append(entities, e)
		}
		var stopId *string
		value := saCode.Value()
		stopId = &value

		collectedRefs.MonitoringRefs[saCode.Value()] = struct{}{}
		entities = append(entities, &gtfs.EntitySelector{StopId: stopId})
		return entities, collectedRefs, nil
	}
	return nil, nil, fmt.Errorf("unsupported value: %T", a)
}
