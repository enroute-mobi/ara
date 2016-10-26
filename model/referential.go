package model

import "encoding/json"

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
	Save(stopArea *Referential) bool
	Delete(stopArea *Referential) bool
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

func (manager *MemoryReferentials) Find(id ReferentialId) (Referential, bool) {
	referential, ok := manager.byId[id]
	if ok {
		return *referential, true
	} else {
		return Referential{}, false
	}
}

func (manager *MemoryReferentials) Save(referential *Referential) bool {
	if referential.Id() == "" {
		referential.id = ReferentialId(manager.NewUUID())
	}
	referential.manager = manager
	manager.byId[referential.Id()] = referential
	return true
}

func (manager *MemoryReferentials) Delete(referential *Referential) bool {
	delete(manager.byId, referential.Id())
	return true
}
