package model

import "encoding/json"

type PartnerId string

type Partners interface {
	New(name string) Partner
	Find(id PartnerId) (Partner, bool)
	FindByName(name string) (Partner, bool)
	Save(partner *Partner) bool
	Delete(partner *Partner) bool
}

type Partner struct {
	id   PartnerId
	name string

	manager Partners
}

type PartnerManager struct {
	UUIDConsumer

	byId map[PartnerId]*Partner
}

func (partner *Partner) Id() PartnerId {
	return partner.id
}

func (partner *Partner) Name() string {
	return partner.name
}

func (partner *Partner) Save() (ok bool) {
	ok = partner.manager.Save(partner)
	return
}

func (partner *Partner) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Id":   partner.id,
		"Name": partner.name,
	})
}

func NewPartnerManager() *PartnerManager {
	return &PartnerManager{
		byId: make(map[PartnerId]*Partner),
	}
}

func (manager *PartnerManager) New(name string) Partner {
	return Partner{name: name, manager: manager}
}

func (manager *PartnerManager) Find(id PartnerId) (Partner, bool) {
	partner, ok := manager.byId[id]
	if ok {
		return *partner, true
	}
	return Partner{}, false
}

func (manager *PartnerManager) FindByName(name string) (Partner, bool) {
	for _, partner := range manager.byId {
		if partner.name == name {
			return *partner, true
		}
	}
	return Partner{}, false
}

func (manager *PartnerManager) Save(partner *Partner) bool {
	if partner.id == "" {
		partner.id = PartnerId(manager.NewUUID())
	}
	partner.manager = manager
	manager.byId[partner.id] = partner
	return true
}

func (manager *PartnerManager) Delete(partner *Partner) bool {
	delete(manager.byId, partner.id)
	return true
}
