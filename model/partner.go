package model

import (
	"encoding/json"
	"sort"

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

	connectors map[string]Connector

	manager Partners
}

type APIPartner struct {
	Id             PartnerId `json:"Id,omitempty"`
	Slug           PartnerSlug
	Settings       map[string]string `json:"Settings,omitempty"`
	ConnectorTypes []string          `json:"ConnectorTypes,omitempty"`
	Errors         []string          `json:"Errors,omitempty"`

	factories map[string]ConnectorFactory
}

type PartnerManager struct {
	UUIDConsumer

	byId     map[PartnerId]*Partner
	guardian *PartnersGuardian
}

func (partner *APIPartner) Validate() bool {
	partner.setFactories()
	valid := true
	for _, factory := range partner.factories {
		if !factory.Validate(partner) {
			valid = false
		}
	}
	return valid
}

func (partner *APIPartner) setFactories() {
	for _, connectorType := range partner.ConnectorTypes {
		factory := NewConnectorFactory(connectorType)
		if factory != nil {
			partner.factories[connectorType] = factory
		}
	}
}

func (partner *APIPartner) isConnectorDefined(expected string) bool {
	for _, connectorType := range partner.ConnectorTypes {
		if connectorType == expected {
			return true
		}
	}
	return false
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
	return partner.manager.Save(partner)
}

func (partner *Partner) MarshalJSON() ([]byte, error) {
	return json.Marshal(APIPartner{
		Id:             partner.id,
		Slug:           partner.slug,
		Settings:       partner.Settings,
		ConnectorTypes: partner.ConnectorTypes,
	})
}

func (partner *Partner) Definition() *APIPartner {
	return &APIPartner{
		Id:             partner.id,
		Slug:           partner.slug,
		Settings:       partner.Settings,
		ConnectorTypes: partner.ConnectorTypes,
		factories:      make(map[string]ConnectorFactory),
	}
}

// APIPartner.Validate should be called for APIPartner factories to be set
func (partner *Partner) SetDefinition(apiPartner *APIPartner) {
	partner.id = apiPartner.Id
	partner.slug = apiPartner.Slug
	partner.Settings = apiPartner.Settings
	partner.ConnectorTypes = apiPartner.ConnectorTypes

	for id, factory := range apiPartner.factories {
		if _, ok := partner.connectors[id]; !ok {
			partner.connectors[id] = factory.CreateConnector(partner)
		}
	}
	partner.cleanConnectors()
}

// Test method, refresh Connector instances according to connector type list without validation
func (partner *Partner) RefreshConnectors() {
	logger.Log.Debugf("Initialize Connectors %#v for %s", partner.ConnectorTypes, partner.slug)

	for _, connectorType := range partner.ConnectorTypes {
		if _, ok := partner.connectors[connectorType]; !ok {
			partner.connectors[connectorType] = NewConnectorFactory(connectorType).CreateConnector(partner)
		}
	}
	partner.cleanConnectors()
}

// Delete from partner.Connectors connectors not in partner.ConnectorTypes
func (partner *Partner) cleanConnectors() {
	sort.Strings(partner.ConnectorTypes)

	for connector, _ := range partner.connectors {
		found := sort.SearchStrings(partner.ConnectorTypes, connector)
		if found == len(partner.ConnectorTypes) {
			delete(partner.connectors, connector)
		}
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
	if partner.isConnectorDefined(SIRI_CHECK_STATUS_CLIENT_TYPE) {
		return partner.connectors[SIRI_CHECK_STATUS_CLIENT_TYPE].(CheckStatusClient)
	} else if partner.isConnectorDefined(TEST_CHECK_STATUS_CLIENT_TYPE) {
		return partner.connectors[TEST_CHECK_STATUS_CLIENT_TYPE].(CheckStatusClient)
	} else {
		return nil
	}
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

func (manager *PartnerManager) Guardian() *PartnersGuardian {
	return manager.guardian
}

func (manager *PartnerManager) Start() {
	manager.guardian.Start()
}

func (manager *PartnerManager) Stop() {
	manager.guardian.Stop()
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
