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

	manager Referentials
	model   Model
}

type Referentials interface {
	New(slug ReferentialSlug) Referential
	Find(id ReferentialId) (Referential, bool)
	FindBySlug(slug ReferentialSlug) (Referential, bool)
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

func (referential *Referential) Save() (ok bool) {
	ok = referential.manager.Save(referential)
	return
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

func (manager *MemoryReferentials) New(slug ReferentialSlug) Referential {
	model := NewMemoryModel()
	return Referential{slug: slug, manager: manager, model: model}
}

func (manager *MemoryReferentials) NewWithId(slug ReferentialSlug, id ReferentialId) Referential {
	model := NewMemoryModel()
	return Referential{id: id, slug: slug, manager: manager, model: model}
}

func (manager *MemoryReferentials) Find(id ReferentialId) (Referential, bool) {
	referential, ok := manager.byId[id]
	if ok {
		return *referential, true
	} else {
		return Referential{}, false
	}
}

func (manager *MemoryReferentials) FindBySlug(slug ReferentialSlug) (Referential, bool) {
	for _, referential := range manager.byId {
		if referential.slug == slug {
			return *referential, true
		}
	}
	return Referential{}, false
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
		referential := manager.NewWithId(ReferentialSlug(r.Slug), ReferentialId(r.Referential_id))
		manager.Save(&referential)
	}

	logger.Log.Debugf("Loaded Referentials from database")
	return nil
}
