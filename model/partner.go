package model

import "encoding/json"

type OperationnalStatus int

const (
	OPERATIONNAL_STATUS_UNKNOWN OperationnalStatus = iota
	OPERATIONNAL_STATUS_UP
	OPERATIONNAL_STATUS_DOWN
)

type PartnerId string

type Partners interface {
	New() *Partner
	Find(id PartnerId) *Partner
	FindAll() []*Partner
	Save(partner *Partner) bool
	Delete(partner *Partner) bool
}

type Partner struct {
	id                 PartnerId
	name               string
	operationnalStatus OperationnalStatus

	checkStatusClient CheckStatusClient

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

func (partner *Partner) OperationnalStatus() OperationnalStatus {
	return partner.operationnalStatus
}

func (partner *Partner) Save() (ok bool) {
	return partner.manager.Save(partner)
}

func (partner *Partner) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"Id":   partner.id,
		"Name": partner.name,
	})
}

// Refresh Connector instances according to connector type list
func (partner *Partner) RefreshConnectors() {
	// WIP
	if partner.checkStatusClient != nil {
		siriPartner := NewSIRIPartner(partner)
		partner.checkStatusClient = NewSIRICheckStatusClient(siriPartner)
	}
}

func (partner *Partner) CheckStatusClient() CheckStatusClient {
	// WIP
	return partner.checkStatusClient
}

func (partner *Partner) CheckStatus() {
	partner.operationnalStatus, _ = partner.CheckStatusClient().Status()
}

func NewPartnerManager() *PartnerManager {
	return &PartnerManager{
		byId: make(map[PartnerId]*Partner),
	}
}

func (manager *PartnerManager) New() *Partner {
	return &Partner{manager: manager}
}

func (manager *PartnerManager) Find(id PartnerId) *Partner {
	partner, _ := manager.byId[id]
	return partner
}

func (manager *PartnerManager) FindAll() (partners []*Partner) {
	for _, partner := range manager.byId {
		partners = append(partners, partner)
	}
	return
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
