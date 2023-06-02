package core

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/config"
	e "bitbucket.org/enroute-mobi/ara/core/apierrs"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type ReferentialId string
type ReferentialSlug string

// Validation
var slugRegexp = regexp.MustCompile(`^[a-z0-9-_]+$`)

type Referential struct {
	clock.ClockConsumer
	s.ReferentialSettings

	id   ReferentialId
	slug ReferentialSlug
	Name string `json:",omitempty"`

	OrganisationId string `json:",omitempty"`

	collectManager    CollectManagerInterface
	broacasterManager BroadcastManagerInterface
	manager           Referentials
	model             *model.MemoryModel
	modelGuardian     *ModelGuardian
	partners          Partners
	startedAt         time.Time
	nextReloadAt      time.Time
	Tokens            []string `json:",omitempty"`
	ImportTokens      []string `json:",omitempty"`
}

type Referentials interface {
	state.Startable

	New(ReferentialSlug) *Referential
	Find(ReferentialId) *Referential
	FindBySlug(ReferentialSlug) *Referential
	FindAll() []*Referential
	Save(*Referential) bool
	Delete(*Referential) bool
	Load() error
	SaveToDatabase() (int, error)
}

var referentials = NewMemoryReferentials()

type APIReferential struct {
	id             ReferentialId
	OrganisationId string            `json:",omitempty"`
	Slug           ReferentialSlug   `json:"Slug,omitempty"`
	Name           string            `json:",omitempty"`
	Errors         e.Errors          `json:"Errors,omitempty"`
	Settings       map[string]string `json:"Settings,omitempty"`
	Tokens         []string          `json:"Tokens,omitempty"`
	ImportTokens   []string          `json:"ImportTokens,omitempty"`

	manager Referentials
}

func (referential *APIReferential) Id() ReferentialId {
	return referential.id
}

func (referential *APIReferential) Validate() bool {
	referential.Errors = e.NewErrors()

	if referential.Slug == "" {
		referential.Errors.Add("Slug", e.ERROR_BLANK)
	} else if !slugRegexp.MatchString(string(referential.Slug)) {
		referential.Errors.Add("Slug", e.ERROR_SLUG_FORMAT)
	}

	// Check Slug uniqueness
	for _, existingReferential := range referential.manager.FindAll() {
		if existingReferential.id != referential.Id() {
			if referential.Slug == existingReferential.slug {
				referential.Errors.Add("Slug", e.ERROR_UNIQUE)
			}
		}
	}

	return len(referential.Errors) == 0
}

func (referential *Referential) Id() ReferentialId {
	return referential.id
}

func (referential *Referential) Slug() ReferentialSlug {
	return referential.slug
}

func (referential *Referential) StartedAt() time.Time {
	return referential.startedAt
}

// WIP: Interface ?
func (referential *Referential) CollectManager() CollectManagerInterface {
	return referential.collectManager
}

func (referential *Referential) Model() model.Model {
	return referential.model
}

func (referential *Referential) ModelGuardian() *ModelGuardian {
	return referential.modelGuardian
}

func (referential *Referential) Partners() Partners {
	return referential.partners
}

func (referential *Referential) DatabaseOrganisationId() sql.NullString {
	if referential.OrganisationId == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: referential.OrganisationId,
		Valid:  true,
	}
}

func (referential *Referential) Start() {
	referential.startedAt = referential.Clock().Now()

	// Configure BigQuery
	if config.Config.ValidBQConfig() {
		dataset := fmt.Sprintf("%v_%v", config.Config.BigQueryDatasetPrefix, referential.slug)
		audit.SetCurrentBigQuery(string(referential.slug), audit.NewBigQuery(dataset))
		audit.CurrentBigQuery(string(referential.slug)).Start()
	}

	referential.partners.Start()
	referential.modelGuardian.Start()

	referential.broacasterManager = NewBroadcastManager(referential)
	referential.model.SetBroadcastSMChan(referential.broacasterManager.GetStopMonitoringBroadcastEventChan())
	referential.model.SetBroadcastGMChan(referential.broacasterManager.GetGeneralMessageBroadcastEventChan())
	referential.model.SetBroadcastVeChan(referential.broacasterManager.GetVehicleBroadcastEventChan())

	referential.broacasterManager.Start()

}

func (referential *Referential) Stop() {
	referential.partners.Stop()
	referential.modelGuardian.Stop()
	referential.broacasterManager.Stop()
	audit.CurrentBigQuery(string(referential.slug)).Stop()
}

func (referential *Referential) Save() (ok bool) {
	ok = referential.manager.Save(referential)
	return
}

func (referential *Referential) MarshalJSON() ([]byte, error) {
	type Alias Referential
	aux := struct {
		Id           ReferentialId
		Slug         ReferentialSlug
		NextReloadAt *time.Time `json:",omitempty"`
		Partners     Partners   `json:",omitempty"`
		Settings     map[string]string
		*Alias
	}{
		Id:       referential.id,
		Slug:     referential.slug,
		Settings: referential.SettingsDefinition(),
		Alias:    (*Alias)(referential),
	}

	if !referential.nextReloadAt.IsZero() {
		aux.NextReloadAt = &referential.nextReloadAt
	}
	if !referential.partners.IsEmpty() {
		aux.Partners = referential.partners
	}

	return json.Marshal(&aux)
}

func (referential *Referential) Definition() *APIReferential {
	return &APIReferential{
		id:             referential.id,
		OrganisationId: referential.OrganisationId,
		Slug:           referential.slug,
		Name:           referential.Name,
		Settings:       referential.SettingsDefinition(),
		Errors:         e.NewErrors(),
		manager:        referential.manager,
		Tokens:         referential.Tokens,
		ImportTokens:   referential.ImportTokens,
	}
}

func (referential *Referential) SetDefinition(apiReferential *APIReferential) {
	initialReloadAt := referential.Setting(s.MODEL_RELOAD_AT)

	referential.OrganisationId = apiReferential.OrganisationId
	referential.slug = apiReferential.Slug
	referential.Name = apiReferential.Name
	referential.SetSettingsDefinition(apiReferential.Settings)
	referential.Tokens = apiReferential.Tokens
	referential.ImportTokens = apiReferential.ImportTokens

	if initialReloadAt != referential.Setting(s.MODEL_RELOAD_AT) {
		referential.setNextReloadAt()
	}
}

func (referential *Referential) NextReloadAt() time.Time {
	return referential.nextReloadAt
}

func (referential *Referential) ReloadModel() {
	logger.Log.Printf("Reset Model for referential %v", referential.slug)
	referential.Stop()
	referential.model = referential.model.Reload(string(referential.Slug()))
	referential.setNextReloadAt()
	referential.Start()
}

func (referential *Referential) setNextReloadAt() {
	hour, minute := referential.NextReloadAtSetting()

	now := referential.Clock().Now()
	day := now.Day()

	if now.Hour() > hour || (now.Hour() == hour && now.Minute() >= minute) {
		day += 1
	}

	referential.nextReloadAt = time.Date(now.Year(), now.Month(), day, hour, minute, 0, 0, now.Location())
	logger.Log.Printf("Next reload at: %v", referential.nextReloadAt)
}

func (referential *Referential) Load() {
	referential.Partners().Load()
	referential.model.Load(string(referential.slug))
}

type MemoryReferentials struct {
	uuid.UUIDConsumer

	byId map[ReferentialId]*Referential
}

func NewMemoryReferentials() *MemoryReferentials {
	return &MemoryReferentials{
		byId: make(map[ReferentialId]*Referential),
	}
}

func CurrentReferentials() Referentials {
	return referentials
}

func (manager *MemoryReferentials) New(slug ReferentialSlug) *Referential {
	model := model.NewMemoryModel(string(slug))

	referential := &Referential{
		ReferentialSettings: s.NewReferentialSettings(),
		manager:             manager,
		model:               model,
		slug:                slug,
	}

	referential.partners = NewPartnerManager(referential)
	referential.collectManager = NewCollectManager(referential)
	referential.broacasterManager = NewBroadcastManager(referential)

	referential.model.SetBroadcastSMChan(referential.broacasterManager.GetStopMonitoringBroadcastEventChan())
	referential.model.SetBroadcastGMChan(referential.broacasterManager.GetGeneralMessageBroadcastEventChan())

	referential.modelGuardian = NewModelGuardian(referential)
	referential.setNextReloadAt()

	return referential
}

func (manager *MemoryReferentials) Find(id ReferentialId) *Referential {
	referential := manager.byId[id]
	return referential
}

func (manager *MemoryReferentials) FindBySlug(slug ReferentialSlug) *Referential {
	for _, referential := range manager.byId {
		if referential.slug == slug {
			return referential
		}
	}
	return nil
}

func (manager *MemoryReferentials) FindAll() (referentials []*Referential) {
	if len(manager.byId) == 0 {
		return []*Referential{}
	}
	for _, referential := range manager.byId {
		referentials = append(referentials, referential)
	}
	return
}

func (manager *MemoryReferentials) Save(referential *Referential) bool {
	if referential.id == "" {
		referential.id = ReferentialId(manager.NewUUID())
	}
	referential.manager = manager
	referential.model.SetReferential(string(referential.slug))
	referential.collectManager.HandleSituationUpdateEvent(model.NewSituationUpdateManager(referential.Model()))
	referential.collectManager.HandleUpdateEvent(model.NewUpdateManager(referential.Model()))
	manager.byId[referential.id] = referential
	return true
}

func (manager *MemoryReferentials) Delete(referential *Referential) bool {
	delete(manager.byId, referential.id)
	return true
}

func (manager *MemoryReferentials) Load() error {
	selectReferentials := []model.SelectReferential{}
	_, err := model.Database.Select(&selectReferentials, "select * from referentials")
	if err != nil {
		return err
	}

	for _, r := range selectReferentials {
		referential := manager.New(ReferentialSlug(r.Slug))
		referential.id = ReferentialId(r.ReferentialId)

		if r.Name.Valid {
			referential.Name = r.Name.String
		}

		if r.OrganisationId.Valid {
			referential.OrganisationId = r.OrganisationId.String
		}

		if r.Settings.Valid && len(r.Settings.String) > 0 {
			m := make(map[string]string)
			if err = json.Unmarshal([]byte(r.Settings.String), &m); err != nil {
				return err
			}
			referential.SetSettingsDefinition(m)
		}

		if r.Tokens.Valid && len(r.Tokens.String) > 0 {
			if err = json.Unmarshal([]byte(r.Tokens.String), &referential.Tokens); err != nil {
				return err
			}
		}

		if r.ImportTokens.Valid && len(r.ImportTokens.String) > 0 {
			if err = json.Unmarshal([]byte(r.ImportTokens.String), &referential.ImportTokens); err != nil {
				return err
			}
		}

		referential.setNextReloadAt()
		manager.Save(referential)
		referential.Load()
	}

	logger.Log.Debugf("Loaded Referentials from database")
	return nil
}

func (manager *MemoryReferentials) SaveToDatabase() (int, error) {
	// Begin transaction
	tx, err := model.Database.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	// Truncate Table
	_, err = tx.Exec("truncate referentials;")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	// Insert referentials
	for _, referential := range manager.byId {
		dbReferential, err := manager.newDbReferential(referential)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, fmt.Errorf("internal error: %v", err)
		}
		err = tx.Insert(dbReferential)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
		}
	}

	// Delete partners
	_, err = tx.Exec("delete from partners where referential_id not in (select referential_id from referentials);")
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	return http.StatusOK, nil
}

func (manager *MemoryReferentials) newDbReferential(referential *Referential) (*model.DatabaseReferential, error) {
	settings, err := referential.ReferentialSettings.ToJson()
	if err != nil {
		return nil, err
	}
	tokens, err := json.Marshal(referential.Tokens)
	if err != nil {
		return nil, err
	}
	importTokens, err := json.Marshal(referential.ImportTokens)
	if err != nil {
		return nil, err
	}
	return &model.DatabaseReferential{
		ReferentialId:  string(referential.id),
		OrganisationId: referential.DatabaseOrganisationId(),
		Slug:           string(referential.slug),
		Name:           referential.Name,
		Settings:       string(settings),
		Tokens:         string(tokens),
		ImportTokens:   string(importTokens),
	}, nil
}

func (manager *MemoryReferentials) Start() {
	for _, referential := range manager.byId {
		referential.Start()
	}
}

type ReferentialsConsumer struct {
	referentials Referentials
}

func (consumer *ReferentialsConsumer) SetReferentials(referentials Referentials) {
	consumer.referentials = referentials
}

func (consumer *ReferentialsConsumer) CurrentReferentials() Referentials {
	if consumer.referentials == nil {
		consumer.referentials = CurrentReferentials()
	}
	return consumer.referentials
}
