package model

import (
	"encoding/json"

	"github.com/af83/edwig/logger"
)

type OperationnalStatus int

const (
	OPERATIONNAL_STATUS_UNKNOWN OperationnalStatus = iota
	OPERATIONNAL_STATUS_UP
	OPERATIONNAL_STATUS_DOWN
)

type PartnerId string

type Partners interface {
	UUIDInterface
	Startable

	New() *Partner
	Find(id PartnerId) *Partner
	FindAll() []*Partner
	Save(partner *Partner) bool
	Delete(partner *Partner) bool
}

type Partner struct {
	id                 PartnerId
	Name               string
	Settings           map[string]string
	operationnalStatus OperationnalStatus

	checkStatusClient CheckStatusClient

	manager Partners
}

type PartnerManager struct {
	UUIDConsumer

	byId     map[PartnerId]*Partner
	guardian *PartnersGuardian
}

func (partner *Partner) Id() PartnerId {
	return partner.id
}

func (partner *Partner) Setting(key string) string {
	return partner.Settings[key]
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
		"Name": partner.Name,
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
	logger.Log.Debugf("Check '%s' partner status", partner.Name)

	partner.operationnalStatus, _ = partner.CheckStatusClient().Status()
}

func NewPartnerManager() *PartnerManager {
	manager := &PartnerManager{
		byId: make(map[PartnerId]*Partner),
	}
	manager.guardian = NewPartnersGuardian(manager)
	return manager
}

func (manager *PartnerManager) Start() {
	manager.guardian.Start()
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
	return partners
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
