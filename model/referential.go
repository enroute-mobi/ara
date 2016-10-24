package model

type ReferentialSlug string

type Referential struct {
	slug  ReferentialSlug
	model *MemoryModel
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

type MemoryReferentials struct {
	model *MemoryModel

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
	return Referential{}
}
func (manager *MemoryReferentials) Find(slug ReferentialSlug) (Referential, bool) {
	return Referential{}, false
}
func (manager *MemoryReferentials) Save(stopArea *Referential) bool {
	return false
}
func (manager *MemoryReferentials) Delete(stopArea *Referential) bool {
	return false
}
