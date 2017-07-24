package core

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type OperationnalStatus string

const (
	OPERATIONNAL_STATUS_UNKNOWN OperationnalStatus = "unknown"
	OPERATIONNAL_STATUS_UP      OperationnalStatus = "up"
	OPERATIONNAL_STATUS_DOWN    OperationnalStatus = "down"
)

type PartnerId string
type PartnerSlug string

type Partners interface {
	model.UUIDInterface
	model.Startable
	model.Stopable

	New(slug PartnerSlug) *Partner
	Find(id PartnerId) *Partner
	FindByLocalCredential(credential string) (*Partner, bool)
	FindBySlug(slug PartnerSlug) (*Partner, bool)
	FindAllByCollectPriority() []*Partner
	FindAll() []*Partner
	Save(partner *Partner) bool
	Delete(partner *Partner) bool
	Model() model.Model
	Referential() *Referential
	IsEmpty() bool
	CancelSubscriptions()
	Load() error
}

type PartnerStatus struct {
	OperationnalStatus OperationnalStatus
	ServiceStartedAt   time.Time
}

type Partner struct {
	id            PartnerId
	slug          PartnerSlug
	PartnerStatus PartnerStatus

	ConnectorTypes []string
	Settings       map[string]string

	connectors          map[string]Connector
	context             Context
	subscriptionManager Subscriptions
	manager             Partners
}

type ByPriority []*Partner

func (a ByPriority) Len() int      { return len(a) }
func (a ByPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool {
	first, _ := strconv.Atoi(a[i].Settings["collect.priority"])
	second, _ := strconv.Atoi(a[j].Settings["collect.priority"])
	return first > second
}

func (manager *PartnerManager) FindAllByCollectPriority() []*Partner {
	partners := manager.FindAll()
	sort.Sort(ByPriority(partners))
	return partners
}

type APIPartner struct {
	Id             PartnerId `json:"Id,omitempty"`
	Slug           PartnerSlug
	Settings       map[string]string `json:"Settings,omitempty"`
	ConnectorTypes []string          `json:"ConnectorTypes,omitempty"`
	Errors         Errors            `json:"Errors,omitempty"`

	factories map[string]ConnectorFactory
	manager   Partners
}

type PartnerManager struct {
	model.UUIDConsumer

	byId        map[PartnerId]*Partner
	guardian    *PartnersGuardian
	referential *Referential
}

func (partner *APIPartner) Validate() bool {
	partner.Errors = NewErrors()

	// Check if slug is non null
	if partner.Slug == "" {
		partner.Errors.Add("Slug", ERROR_BLANK)
	}

	// Check factories
	partner.setFactories()
	for _, factory := range partner.factories {
		factory.Validate(partner)
	}

	// Check local_credential and Slug uniqueness
	credentials, ok := partner.Settings["local_credential"]
	for _, existingPartner := range partner.manager.FindAll() {
		if existingPartner.id != partner.Id {
			if partner.Slug == existingPartner.slug {
				partner.Errors.Add("Slug", ERROR_UNIQUE)
			}
			if ok && credentials == existingPartner.Settings["local_credential"] {
				partner.Errors.Add("Settings[\"local_credential\"]", ERROR_UNIQUE)
			}
		}
	}

	return len(partner.Errors) == 0
}

func (partner *APIPartner) setFactories() {
	for _, connectorType := range partner.ConnectorTypes {
		factory := NewConnectorFactory(connectorType)
		if factory != nil {
			partner.factories[connectorType] = factory
		}
	}
}

func (partner *APIPartner) IsSettingDefined(setting string) (ok bool) {
	_, ok = partner.Settings[setting]
	return
}

func (partner *APIPartner) ValidatePresenceOfSetting(setting string) bool {
	if !partner.IsSettingDefined(setting) {
		partner.Errors.Add(fmt.Sprintf("Setting %s", setting), ERROR_BLANK)
		return false
	}
	return true
}

func (partner *APIPartner) ValidatePresenceOfConnector(connector string) bool {
	for _, listedConnector := range partner.ConnectorTypes {
		if listedConnector == connector {
			return true
		}
	}
	partner.Errors.Add(fmt.Sprintf("Connector %s", connector), ERROR_BLANK)
	return false
}

func (partner *APIPartner) UnmarshalJSON(data []byte) error {
	type Alias APIPartner
	aux := &struct {
		Settings map[string]string
		*Alias
	}{
		Alias: (*Alias)(partner),
	}
	err := json.Unmarshal(data, aux)
	if err != nil {
		return err
	}

	if aux.Settings != nil {
		partner.Settings = aux.Settings
	}

	return nil
}

func NewPartner() *Partner {
	partner := &Partner{
		Settings:       make(map[string]string),
		ConnectorTypes: []string{},
		connectors:     make(map[string]Connector),
		context:        make(Context),
		PartnerStatus: PartnerStatus{
			OperationnalStatus: OPERATIONNAL_STATUS_UNKNOWN,
		},
	}
	partner.subscriptionManager = NewMemorySubscriptions(partner)

	return partner
}

func (partner *Partner) Referential() *Referential {
	return partner.manager.Referential()
}

func (partner *Partner) Subscriptions() Subscriptions {
	return partner.subscriptionManager
}

func (partner *Partner) Id() PartnerId {
	return partner.id
}

func (partner *Partner) Slug() PartnerSlug {
	return partner.slug
}

func (partner *Partner) SetSlug(s PartnerSlug) {
	partner.slug = s
}

func (partner *Partner) Setting(key string) string {
	return partner.Settings[key]
}

func (partner *Partner) OperationnalStatus() OperationnalStatus {
	return partner.PartnerStatus.OperationnalStatus
}

func (partner *Partner) Context() *Context {
	return &partner.context
}

func (partner *Partner) Save() (ok bool) {
	return partner.manager.Save(partner)
}

func (partner *Partner) MarshalJSON() ([]byte, error) {
	type Alias Partner
	return json.Marshal(&struct {
		Id   PartnerId
		Slug PartnerSlug
		*Alias
	}{
		Id:    partner.id,
		Slug:  partner.slug,
		Alias: (*Alias)(partner),
	})
}

func (partner *Partner) Definition() *APIPartner {
	return &APIPartner{
		Id:             partner.id,
		Slug:           partner.slug,
		Settings:       partner.Settings,
		ConnectorTypes: partner.ConnectorTypes,
		factories:      make(map[string]ConnectorFactory),
		Errors:         NewErrors(),
		manager:        partner.manager,
	}
}

func (partner *Partner) Stop() {
	for _, connector := range partner.connectors {
		c, ok := connector.(model.Stopable)
		if ok {
			c.Stop()
		}
	}
}

func (partner *Partner) Start() {
	for _, connector := range partner.connectors {
		c, ok := connector.(model.Startable)
		if ok {
			c.Start()
		}
	}
}

func (partner *Partner) CollectPriority() int {
	value, _ := strconv.Atoi(partner.Settings["collect.priority"])
	return value
}

func (partner *Partner) CanCollect(stopAreaObjectId model.ObjectID) bool {
	if partner.Settings["collect.include_stop_areas"] == "" {
		return true
	}
	stopAreas := strings.Split(partner.Settings["collect.include_stop_areas"], ",")
	for _, stopArea := range stopAreas {
		if stopArea == stopAreaObjectId.Value() {
			return true
		}
	}
	return false
}

// APIPartner.Validate should be called for APIPartner factories to be set
func (partner *Partner) SetDefinition(apiPartner *APIPartner) {
	partner.id = apiPartner.Id
	partner.slug = apiPartner.Slug
	partner.Settings = apiPartner.Settings
	partner.ConnectorTypes = apiPartner.ConnectorTypes
	partner.PartnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UNKNOWN

	for id, factory := range apiPartner.factories {
		if _, ok := partner.connectors[id]; !ok {
			logger.Log.Debugf("Create connector %v for partner %v", id, partner.slug)
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

	for connectorType, _ := range partner.connectors {
		found := sort.SearchStrings(partner.ConnectorTypes, connectorType)
		if found == len(partner.ConnectorTypes) || partner.ConnectorTypes[found] != connectorType {
			delete(partner.connectors, connectorType)
		}
	}
}

func (partner *Partner) Connector(connectorType string) (Connector, bool) {
	connector, ok := partner.connectors[connectorType]
	return connector, ok
}

func (partner *Partner) CheckStatusClient() CheckStatusClient {
	// WIP
	client, ok := partner.connectors[SIRI_CHECK_STATUS_CLIENT_TYPE]
	if ok {
		return client.(CheckStatusClient)
	}
	client, ok = partner.connectors[TEST_CHECK_STATUS_CLIENT_TYPE]
	if ok {
		return client.(CheckStatusClient)
	}
	return nil
}

func (partner *Partner) GeneralMessageRequestCollector() GeneralMessageRequestCollector {
	// WIP
	client, ok := partner.connectors[SIRI_GENERAL_MESSAGE_REQUEST_COLLECTOR]
	if ok {
		return client.(GeneralMessageRequestCollector)
	}
	return nil
}

func (partner *Partner) GeneralMessageSubscriptionCollector() GeneralMessageSubscriptionCollector {
	// WIP
	client, ok := partner.connectors[SIRI_GENERAL_MESSAGE_DELIVERIES_RESPONSE_COLLECTOR]
	if ok {
		return client.(GeneralMessageSubscriptionCollector)
	}
	return nil
}

func (partner *Partner) StopMonitoringSubscriptionRequestCollector() StopMonitoringSubscriptionCollector {
	// WIP
	client, ok := partner.connectors[SIRI_STOP_MONITORING_DELIVERIES_RESPONSE_COLLECTOR]
	if ok {
		return client.(StopMonitoringSubscriptionCollector)
	}
	return nil
}

func (partner *Partner) StopMonitoringRequestCollector() StopMonitoringRequestCollector {
	// WIP
	client, ok := partner.connectors[SIRI_STOP_MONITORING_REQUEST_COLLECTOR]
	if ok {
		return client.(StopMonitoringRequestCollector)
	}
	client, ok = partner.connectors[TEST_STOP_MONITORING_REQUEST_COLLECTOR]
	if ok {
		return client.(StopMonitoringRequestCollector)
	}
	return nil
}

func (partner *Partner) CheckStatus() (PartnerStatus, error) {
	logger.Log.Debugf("Check '%s' partner status", partner.slug)
	partnerStatus := PartnerStatus{}

	if partner.CheckStatusClient() == nil {
		logger.Log.Debugf("No CheckStatusClient connector")
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UNKNOWN
		return partnerStatus, errors.New("No CheckStatusClient connector")
	}
	partnerStatus, err := partner.CheckStatusClient().Status()

	if err != nil {
		logger.Log.Printf("Error while checking status: %v", err)
	}
	logger.Log.Debugf("Partner status is %v", partner.PartnerStatus.OperationnalStatus)
	return partnerStatus, nil
}

func (partner *Partner) Model() model.Model {
	return partner.manager.Model()
}

func (partner *Partner) DeleteSubscription(sub Subscription) {
	partner.Subscriptions().Delete(&sub)
}

func (partner *Partner) CancelSubscriptions() {
	subscriptions := partner.subscriptionManager.FindAll()
	for _, sub := range subscriptions {
		logger.Log.Printf("Deleting subscription id %v and kind %v", sub.Id(), sub.Kind())
		partner.DeleteSubscription(sub)
	}
}

func NewPartnerManager(referential *Referential) *PartnerManager {
	manager := &PartnerManager{
		byId:        make(map[PartnerId]*Partner),
		referential: referential,
	}
	manager.guardian = NewPartnersGuardian(manager)
	return manager
}

func (manager *PartnerManager) Guardian() *PartnersGuardian {
	return manager.guardian
}

func (manager *PartnerManager) Start() {
	manager.guardian.Start()
	for _, partner := range manager.byId {
		partner.Start()
	}
}

func (manager *PartnerManager) Stop() {
	manager.guardian.Stop()
	for _, partner := range manager.byId {
		partner.Stop()
	}
}

func (manager *PartnerManager) New(slug PartnerSlug) *Partner {
	partner := &Partner{
		slug:       slug,
		manager:    manager,
		Settings:   make(map[string]string),
		connectors: make(map[string]Connector),
		context:    make(Context),
		PartnerStatus: PartnerStatus{
			OperationnalStatus: OPERATIONNAL_STATUS_UNKNOWN,
		},
		ConnectorTypes: []string{},
	}
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	return partner
}

func (manager *PartnerManager) MarshalJSON() ([]byte, error) {
	partnersId := []PartnerId{}
	for id, _ := range manager.byId {
		partnersId = append(partnersId, id)
	}
	return json.Marshal(partnersId)
}

func (manager *PartnerManager) Find(id PartnerId) *Partner {
	partner, _ := manager.byId[id]
	return partner
}

func (manager *PartnerManager) FindByLocalCredential(credential string) (*Partner, bool) {
	for _, partner := range manager.byId {
		if partner.Setting("local_credential") == credential {
			return partner, true
		}
	}
	return nil, false
}

func (manager *PartnerManager) FindBySlug(slug PartnerSlug) (*Partner, bool) {
	for _, partner := range manager.byId {
		if partner.slug == slug {
			return partner, true
		}
	}
	return nil, false
}

func (manager *PartnerManager) FindAll() (partners []*Partner) {
	if len(manager.byId) == 0 {
		return []*Partner{}
	}
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

func (manager *PartnerManager) Model() model.Model {
	return manager.referential.Model()
}

func (manager *PartnerManager) IsEmpty() bool {
	return len(manager.byId) == 0
}

func (manager *PartnerManager) Referential() *Referential {
	return manager.referential
}

func (manager *PartnerManager) CancelSubscriptions() {
	for _, partner := range manager.byId {
		partner.CancelSubscriptions()
	}
}

func (manager *PartnerManager) Load() error {
	var selectPartners []struct {
		Id             string
		ReferentialId  string `db:"referential_id"`
		Slug           string
		Settings       sql.NullString
		ConnectorTypes sql.NullString `db:"connector_types"`
	}
	sqlQuery := fmt.Sprintf("select * from partners where referential_id = '%s'", manager.referential.Id())
	_, err := model.Database.Select(&selectPartners, sqlQuery)
	if err != nil {
		return err
	}
	for _, p := range selectPartners {
		partner := manager.New(PartnerSlug(p.Slug))
		partner.id = PartnerId(p.Id)

		if p.Settings.Valid && len(p.Settings.String) > 0 {
			if err = json.Unmarshal([]byte(p.Settings.String), &partner.Settings); err != nil {
				return err
			}
		}

		if p.ConnectorTypes.Valid && len(p.ConnectorTypes.String) > 0 {
			if err = json.Unmarshal([]byte(p.ConnectorTypes.String), &partner.ConnectorTypes); err != nil {
				return err
			}
		}

		partner.RefreshConnectors()
		manager.Save(partner)
	}
	return nil
}
