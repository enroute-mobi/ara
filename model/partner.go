package model

import (
	"encoding/json"
	"fmt"

	"github.com/af83/edwig/logger"
)

type OperationnalStatus int

const (
	OPERATIONNAL_STATUS_UNKNOWN OperationnalStatus = iota
	OPERATIONNAL_STATUS_UP
	OPERATIONNAL_STATUS_DOWN
)

type PartnerId string
type PartnerSlug string

type Partners interface {
	UUIDInterface
	Startable

	New(slug PartnerSlug) *Partner
	Find(id PartnerId) *Partner
	FindAll() []*Partner
	Save(partner *Partner) bool
	Delete(partner *Partner) bool
}

type Partner struct {
	id                 PartnerId
	slug               PartnerSlug
	Settings           map[string]string
	ConnectorTypes     []string
	operationnalStatus OperationnalStatus

	// WIP
	checkStatusClient CheckStatusClient

	manager Partners
}

type APIPartner struct {
	Id             PartnerId
	Slug           PartnerSlug
	Settings       map[string]string `json:"Settings,omitempty"`
	ConnectorTypes []string          `json:"ConnectorTypes,omitempty"`
}

type PartnerManager struct {
	UUIDConsumer

	byId     map[PartnerId]*Partner
	guardian *PartnersGuardian
}

// WIP
func (partner *APIPartner) Validate() error {
	return nil
}

func (partner *Partner) Id() PartnerId {
	return partner.id
}

func (partner *Partner) Slug() PartnerSlug {
	return partner.slug
}

func (partner *Partner) Setting(key string) string {
	return partner.Settings[key]
}

func (partner *Partner) OperationnalStatus() OperationnalStatus {
	return partner.operationnalStatus
}

func (partner *Partner) Save() (ok bool) {
	// WIP
	partner.RefreshConnectors()

	return partner.manager.Save(partner)
}

func (partner *Partner) MarshalJSON() ([]byte, error) {
	return json.Marshal(APIPartner{
		partner.id,
		partner.slug,
		partner.Settings,
		partner.ConnectorTypes,
	})
}

func (partner *Partner) UnmarshalJSON(b []byte) error {
	var apiPartner APIPartner
	if err := json.Unmarshal(b, &apiPartner); err != nil {
		return fmt.Errorf("Can't parse JSON")
	}
	if err := apiPartner.Validate(); err != nil {
		return fmt.Errorf("Invalid Partner")
	}

	if apiPartner.Id != "" {
		partner.id = apiPartner.Id
	}
	if apiPartner.Slug != "" {
		partner.slug = apiPartner.Slug
	}
	if len(apiPartner.Settings) > 0 {
		partner.Settings = apiPartner.Settings
	}
	if len(apiPartner.ConnectorTypes) > 0 {
		partner.ConnectorTypes = apiPartner.ConnectorTypes
	}

	return nil
}

// Refresh Connector instances according to connector type list
func (partner *Partner) RefreshConnectors() {
	// WIP
	logger.Log.Debugf("Initialize Connectors %#v for %s", partner.ConnectorTypes, partner.slug)

	if partner.isConnectorDefined(SIRI_CHECK_STATUS_CLIENT_TYPE) {
		if partner.checkStatusClient == nil {
			siriPartner := NewSIRIPartner(partner)
			partner.checkStatusClient = NewSIRICheckStatusClient(siriPartner)
		}
	} else {
		partner.checkStatusClient = nil
	}
}

func (partner *Partner) isConnectorDefined(expected string) bool {
	for _, connectorType := range partner.ConnectorTypes {
		if connectorType == expected {
			return true
		}
	}
	return false
}

func (partner *Partner) CheckStatusClient() CheckStatusClient {
	// WIP
	return partner.checkStatusClient
}

func (partner *Partner) CheckStatus() {
	logger.Log.Debugf("Check '%s' partner status", partner.slug)

	if partner.CheckStatusClient() == nil {
		logger.Log.Debugf("No CheckStatusClient connector")
		return
	}

	status, err := partner.CheckStatusClient().Status()
	if err != nil {
		logger.Log.Printf("Error while checking status: %v", err)
	}

	partner.operationnalStatus = status
	logger.Log.Debugf("Partner status is %v", partner.operationnalStatus)
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

func (manager *PartnerManager) New(slug PartnerSlug) *Partner {
	return &Partner{slug: slug, manager: manager}
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
