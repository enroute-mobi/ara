package model

import (
	"encoding/json"

	"github.com/af83/edwig/logger"
)

type ReferentialId string
type ReferentialSlug string

type Referential struct {
	id   ReferentialId
	slug ReferentialSlug

	manager  Referentials
	model    Model
	partners Partners
}

type Referentials interface {
	New(slug ReferentialSlug) *Referential
	Find(id ReferentialId) *Referential
	FindBySlug(slug ReferentialSlug) *Referential
	Save(stopArea *Referential) bool
	Delete(stopArea *Referential) bool
	Load() error
}

var referentials = NewMemoryReferentials()

func (referential *Referential) Id() ReferentialId {
	return referential.id
}

func (referential *Referential) Slug() ReferentialSlug {
	return referential.slug
}

func (referential *Referential) Model() Model {
	return referential.model
}

func (referential *Referential) Partners() Partners {
	return referential.partners
}

func (referential *Referential) Start() {
	referential.partners.Start()
}

func (referential *Referential) Save() (ok bool) {
	ok = referential.manager.Save(referential)
	return
}

func (referential *Referential) NewTransaction() *Transaction {
	return NewTransaction(referential.model)
}

func (referential *Referential) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Id":   referential.id,
		"Slug": referential.slug,
	})
}

type MemoryReferentials struct {
	UUIDConsumer

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
	return &Referential{
		manager:  manager,
		model:    NewMemoryModel(),
		partners: NewPartnerManager(),
	}
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

func (manager *MemoryReferentials) Save(referential *Referential) bool {
	if referential.id == "" {
		referential.id = ReferentialId(manager.NewUUID())
	}
	referential.manager = manager
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
	}
	_, err := Database.Select(&selectReferentials, "select * from referentials")
	if err != nil {
		return err
	}

	for _, r := range selectReferentials {
		referential := manager.new()
		referential.id = ReferentialId(r.Referential_id)
		referential.slug = ReferentialSlug(r.Slug)
		manager.Save(referential)
	}

	logger.Log.Debugf("Loaded Referentials from database")
	return nil
}
