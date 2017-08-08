package core

import (
	"encoding/json"
	"fmt"
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

	Settings map[string]string `json:"Settings,omitempty"`

	collectManager    CollectManagerInterface
	broacasterManager BroadcastManagerInterface
	manager           Referentials
	model             *model.MemoryModel
	modelGuardian     *ModelGuardian
	partners          Partners
	startedAt         time.Time
	nextReloadAt      time.Time
	Tokens            []string `json:",omitempty"`
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
	SaveToDatabase() map[string]string
}

var referentials = NewMemoryReferentials()

type APIReferential struct {
	id       ReferentialId
	Slug     ReferentialSlug   `json:"Slug,omitempty"`
	Errors   Errors            `json:"Errors,omitempty"`
	Settings map[string]string `json:"Settings,omitempty"`
	Tokens   []string          `json:"Tokens,omitempty"`

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

	// if len(referential.Tokens) == 0 {
	// 	referential.Errors.Add("Tokens", ERROR_BLANK)
	// }
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

func (referential *Referential) MarshalJSON() ([]byte, error) {
	type Alias Referential
	aux := struct {
		Id           ReferentialId
		Slug         ReferentialSlug
		NextReloadAt *time.Time `json:",omitempty"`
		Partners     Partners   `json:",omitempty"`
		*Alias
	}{
		Id:    referential.id,
		Slug:  referential.slug,
		Alias: (*Alias)(referential),
	}

	if !referential.nextReloadAt.IsZero() {
		aux.NextReloadAt = &referential.nextReloadAt
	}
	if !referential.partners.IsEmpty() {
		aux.Partners = referential.partners
	}

	// if len(referential.Tokens) > 0 {
	// 	referentialMap["Tokens"] = referential.Tokens
	// }

	return json.Marshal(&aux)
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
		Tokens:   referential.Tokens,
	}
}

func (referential *Referential) SetDefinition(apiReferential *APIReferential) {
	initialReloadAt := referential.Setting(REFERENTIAL_SETTING_MODEL_RELOAD_AT)

	referential.slug = apiReferential.Slug
	referential.Settings = apiReferential.Settings
	referential.Tokens = apiReferential.Tokens

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
	referential.partners.CancelSubscriptions()
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
	referential.broacasterManager = NewBroadcastManager(referential)
	referential.model.SetBroadcasteChan(referential.broacasterManager.GetStopMonitoringBroadcastEventChan())

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
	referential.collectManager.HandleSituationUpdateEvent(model.NewSituationUpdateManager(referential))
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
		referential := manager.new()
		referential.id = ReferentialId(r.Referential_id)
		referential.slug = ReferentialSlug(r.Slug)

		if r.Settings.Valid && len(r.Settings.String) > 0 {
			if err = json.Unmarshal([]byte(r.Settings.String), &referential.Settings); err != nil {
				return err
			}
		}

		if r.Tokens.Valid && len(r.Tokens.String) > 0 {
			if err = json.Unmarshal([]byte(r.Tokens.String), &referential.Tokens); err != nil {
				return err
			}
		}

		manager.Save(referential)
		referential.Load()
	}

	logger.Log.Debugf("Loaded Referentials from database")
	return nil
}

func (manager *MemoryReferentials) SaveToDatabase() map[string]string {
	errors := make(map[string]string)

	// Truncate Table
	_, err := model.Database.Exec("truncate referentials;")
	if err != nil {
		errors["internal"] = fmt.Sprintf("Internal error: %v", err)
		return errors
	}

	// Insert referentials
	for _, referential := range manager.byId {
		dbReferential, err := manager.newDbReferential(referential)
		if err != nil {
			errors[string(referential.id)] = fmt.Sprintf("Error while saving referential: %v", err)
		}
		err = model.Database.Insert(dbReferential)
		if err != nil {
			errors[string(referential.id)] = fmt.Sprintf("Error while saving referential: %v", err)
		}
	}

	// Delete partners
	_, err = model.Database.Exec("delete from partners where referential_id not in (select referential_id from referentials);")
	if err != nil {
		errors["internal"] = fmt.Sprintf("Internal error: %v", err)
	}

	return errors
}

func (manager *MemoryReferentials) newDbReferential(referential *Referential) (*model.DatabaseReferential, error) {
	settings, err := json.Marshal(referential.Settings)
	if err != nil {
		return nil, err
	}
	tokens, err := json.Marshal(referential.Tokens)
	if err != nil {
		return nil, err
	}
	return &model.DatabaseReferential{
		ReferentialId: string(referential.id),
		Slug:          string(referential.slug),
		Settings:      string(settings),
		Tokens:        string(tokens),
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
