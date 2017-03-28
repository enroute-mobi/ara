package core

import (
	"database/sql"
	"encoding/json"
	"strconv"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type ReferentialId string
type ReferentialSlug string

const (
	REFERENTIAL_SETTING_MODEL_RELOAD_AT = "model.reload_at"
)

type Referential struct {
	model.ClockConsumer

	id   ReferentialId
	slug ReferentialSlug

	Settings map[string]string

	collectManager CollectManagerInterface
	manager        Referentials
	model          *model.MemoryModel
	modelGuardian  *ModelGuardian
	partners       Partners
	startedAt      time.Time
	nextReloadAt   time.Time
}

type Referentials interface {
	model.Startable

	New(slug ReferentialSlug) *Referential
	Find(id ReferentialId) *Referential
	FindBySlug(slug ReferentialSlug) *Referential
	FindAll() []*Referential
	Save(stopArea *Referential) bool
	Delete(stopArea *Referential) bool
	Load() error
}

var referentials = NewMemoryReferentials()

type APIReferential struct {
	id       ReferentialId
	Slug     ReferentialSlug   `json:"Slug,omitempty"`
	Errors   Errors            `json:"Errors,omitempty"`
	Settings map[string]string `json:"Settings,omitempty"`

	manager Referentials
}

func (referential *APIReferential) Id() ReferentialId {
	return referential.id
}

func (referential *APIReferential) Validate() bool {
	referential.Errors = NewErrors()

	if referential.Slug == "" {
		referential.Errors.Add("Slug", ERROR_BLANK)
	}

	// Check Slug uniqueness
	for _, existingReferential := range referential.manager.FindAll() {
		if existingReferential.id != referential.Id() {
			if referential.Slug == existingReferential.slug {
				referential.Errors.Add("Slug", ERROR_UNIQUE)
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

func (referential *Referential) Setting(key string) string {
	return referential.Settings[key]
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

func (referential *Referential) Start() {
	referential.startedAt = referential.Clock().Now()

	referential.partners.Start()
	referential.modelGuardian.Start()
}

func (referential *Referential) Stop() {
	referential.partners.Stop()
	referential.modelGuardian.Stop()
}

func (referential *Referential) Save() (ok bool) {
	ok = referential.manager.Save(referential)
	return
}

func (referential *Referential) NewTransaction() *model.Transaction {
	return model.NewTransaction(referential.model)
}

func (referential *Referential) FillReferential(referentialMap map[string]interface{}) {
	if referential.id != "" {
		referentialMap["Id"] = referential.id
	}

	if referential.slug != "" {
		referentialMap["Slug"] = referential.slug
	}

	if len(referential.Settings) > 0 {
		referentialMap["Settings"] = referential.Settings
	}

	if !referential.nextReloadAt.IsZero() {
		referentialMap["NextReloadAt"] = referential.nextReloadAt
	}

	if !referential.partners.IsEmpty() {
		referentialMap["Partners"] = referential.partners
	}
}

func (referential *Referential) MarshalJSON() ([]byte, error) {

	referentialMap := make(map[string]interface{})

	referential.FillReferential(referentialMap)

	return json.Marshal(referentialMap)
}

func (referential *Referential) Definition() *APIReferential {
	settings := map[string]string{}
	for k, v := range referential.Settings {
		settings[k] = v
	}

	return &APIReferential{
		id:       referential.id,
		Slug:     referential.slug,
		Settings: settings,
		Errors:   NewErrors(),
		manager:  referential.manager,
	}
}

func (referential *Referential) SetDefinition(apiReferential *APIReferential) {
	initialReloadAt := referential.Setting(REFERENTIAL_SETTING_MODEL_RELOAD_AT)

	referential.slug = apiReferential.Slug
	referential.Settings = apiReferential.Settings

	if initialReloadAt != referential.Setting(REFERENTIAL_SETTING_MODEL_RELOAD_AT) {
		referential.setNextReloadAt()
	}
}

func (referential *Referential) NextReloadAt() time.Time {
	return referential.nextReloadAt
}

func (referential *Referential) ReloadModel() {
	logger.Log.Printf("Reset Model")
	referential.model = referential.model.Clone()
	referential.setNextReloadAt()
}

func (referential *Referential) setNextReloadAt() {
	reloadHour := referential.Setting(REFERENTIAL_SETTING_MODEL_RELOAD_AT)
	hour, minute := 4, 0

	if len(reloadHour) == 5 {
		hour, _ = strconv.Atoi(reloadHour[0:2])
		minute, _ = strconv.Atoi(reloadHour[3:5])
	}
	now := referential.Clock().Now()

	day := now.Day()
	if now.Hour() >= hour && now.Minute() >= minute {
		day += 1
	}

	referential.nextReloadAt = time.Date(now.Year(), now.Month(), day, hour, minute, 0, 0, now.Location())
	logger.Log.Printf("Next reload at: %v", referential.nextReloadAt)
}

func (referential *Referential) Load() {
	referential.Partners().Load()
	referential.model.Load(string(referential.id))
}

type MemoryReferentials struct {
	model.UUIDConsumer

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
	referential := manager.new()
	referential.slug = slug
	return referential
}

func (manager *MemoryReferentials) new() *Referential {
	model := model.NewMemoryModel()

	referential := &Referential{
		manager:  manager,
		model:    model,
		Settings: make(map[string]string),
	}

	referential.partners = NewPartnerManager(referential)
	referential.collectManager = NewCollectManager(referential)

	referential.modelGuardian = NewModelGuardian(referential)
	referential.setNextReloadAt()

	return referential
}

func (manager *MemoryReferentials) Find(id ReferentialId) *Referential {
	referential, _ := manager.byId[id]
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
	referential.collectManager.HandleStopAreaUpdateEvent(model.NewStopAreaUpdateManager(referential))
	manager.byId[referential.id] = referential
	return true
}

func (manager *MemoryReferentials) Delete(referential *Referential) bool {
	delete(manager.byId, referential.id)
	return true
}

func (manager *MemoryReferentials) Load() error {
	var selectReferentials []struct {
		Referential_id string
		Slug           string
		Settings       sql.NullString
	}

	_, err := model.Database.Select(&selectReferentials, "select * from referentials")
	if err != nil {
		return err
	}

	for _, r := range selectReferentials {
		referential := manager.new()
		referential.id = ReferentialId(r.Referential_id)
		referential.slug = ReferentialSlug(r.Slug)

		if r.Settings.Valid && len(r.Settings.String) > 0 {
			if err = json.Unmarshal([]byte(r.Settings.String), &referential.Settings); err != nil {
				return err
			}
		}

		manager.Save(referential)
		referential.Load()
	}

	logger.Log.Debugf("Loaded Referentials from database")
	return nil
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
