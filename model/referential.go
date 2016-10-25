package model

import "encoding/json"

type ReferentialSlug string

type Referential struct {
	slug ReferentialSlug

	manager Referentials
	model   Model
}

type Referentials interface {
	New(slug ReferentialSlug) Referential
	Find(slug ReferentialSlug) (Referential, bool)
	Save(stopArea *Referential) bool
	Delete(stopArea *Referential) bool
}

var referentials = NewMemoryReferentials()

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
		"Slug": referential.slug,
	})
}

type MemoryReferentials struct {
	bySlug map[ReferentialSlug]*Referential
}

func NewMemoryReferentials() *MemoryReferentials {
	return &MemoryReferentials{
		bySlug: make(map[ReferentialSlug]*Referential),
	}
}

func CurrentReferentials() Referentials {
	return referentials
}

func (manager *MemoryReferentials) New(slug ReferentialSlug) Referential {
	model := NewMemoryModel()
	return Referential{slug: slug, manager: manager, model: model}
}

func (manager *MemoryReferentials) Find(slug ReferentialSlug) (Referential, bool) {
	referential, ok := manager.bySlug[slug]
	if ok {
		return *referential, true
	} else {
		return Referential{}, false
	}
}

func (manager *MemoryReferentials) Save(referential *Referential) bool {
	if referential.Slug() == "" {
		return false
	}
	referential.manager = manager
	manager.bySlug[referential.Slug()] = referential
	return true
}

func (manager *MemoryReferentials) Delete(referential *Referential) bool {
	delete(manager.bySlug, referential.Slug())
	return true
}
