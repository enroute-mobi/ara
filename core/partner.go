package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"sort"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/cache"
	e "bitbucket.org/enroute-mobi/ara/core/apierrs"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"golang.org/x/time/rate"
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
	uuid.UUIDInterface
	state.Startable
	state.Stopable

	New(PartnerSlug) *Partner
	Find(PartnerId) *Partner
	FindBySetting(string, string) (*Partner, bool)
	FindBySlug(PartnerSlug) (*Partner, bool)
	FindByCredential(string) (*Partner, bool)
	FindAllByCollectPriority() []*Partner
	FindAllWithConnector([]string) []*Partner
	FindAll() []*Partner
	Save(partner *Partner) bool
	Delete(partner *Partner) bool
	Model() model.Model
	Referential() *Referential
	UniqCredentials(PartnerId, string) bool
	IsEmpty() bool
	CancelSubscriptions()
	Load() error
	SaveToDatabase() (int, error)
}

type PartnerStatus struct {
	OperationnalStatus OperationnalStatus
	RetryCount         int
	ServiceStartedAt   time.Time
}

type PartnerStatusCheck struct {
	LastCheck time.Time
	Status    OperationnalStatus
}

type Partner struct {
	uuid.UUIDConsumer
	s.PartnerSettings

	mutex *sync.RWMutex

	manager             Partners
	subscriptionManager Subscriptions

	id            PartnerId
	slug          PartnerSlug
	Name          string `json:",omitempty"`
	PartnerStatus PartnerStatus

	ConnectorTypes []string

	connectors             map[string]Connector
	discoveredStopAreas    map[string]struct{}
	discoveredLines        map[string]struct{}
	startedAt              time.Time
	lastDiscovery          time.Time
	alternativeStatusCheck PartnerStatusCheck

	httpClient *remote.HTTPClient

	gtfsCache *cache.CacheTable

	limiters map[string]*rate.Limiter
}

type ByPriority []*Partner

func (a ByPriority) Len() int      { return len(a) }
func (a ByPriority) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByPriority) Less(i, j int) bool {
	return a[i].CollectPriority() > a[j].CollectPriority()
}

type PartnerManager struct {
	uuid.UUIDConsumer

	mutex *sync.RWMutex

	byId                  map[PartnerId]*Partner
	localCredentialsIndex *LocalCredentialsIndex
	guardian              *PartnersGuardian
	referential           *Referential
}

// Test Method
func NewPartner() *Partner {
	partner := &Partner{
		mutex:               &sync.RWMutex{},
		ConnectorTypes:      []string{},
		connectors:          make(map[string]Connector),
		discoveredStopAreas: make(map[string]struct{}),
		discoveredLines:     make(map[string]struct{}),
		PartnerStatus: PartnerStatus{
			OperationnalStatus: OPERATIONNAL_STATUS_UNKNOWN,
		},
		gtfsCache: cache.NewCacheTable(),
		limiters:  make(map[string]*rate.Limiter),
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator)
	partner.subscriptionManager = NewMemorySubscriptions(partner)

	return partner
}

func (partner *Partner) Referential() *Referential {
	return partner.manager.Referential()
}

func (partner *Partner) Subscriptions() Subscriptions {
	return partner.subscriptionManager
}

func (partner *Partner) StartedAt() time.Time {
	return partner.startedAt
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

func (partner *Partner) GtfsCache() *cache.CacheTable {
	return partner.gtfsCache
}

func (partner *Partner) HTTPClient() *remote.HTTPClient {
	opts := partner.HTTPClientOptions()
	if partner.httpClient == nil {
		logger.Log.Debugf("Create a new http client in partner %s to %s", partner.Name, opts.Urls.Url)
		partner.httpClient = remote.NewHTTPClient(opts)
	} else if partner.httpClient.HTTPClientUrls != opts.Urls {
		partner.httpClient.SetURLs(opts.Urls)
	}
	return partner.httpClient
}

func (partner *Partner) SIRIClient() *remote.SIRIClient {
	return partner.HTTPClient().SIRIClient()
}

func (partner *Partner) SIRILiteClient() *remote.SIRILiteClient {
	return partner.HTTPClient().SIRILiteClient()
}

func (partner *Partner) OperationnalStatus() OperationnalStatus {
	return partner.PartnerStatus.OperationnalStatus
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
		Settings map[string]string
	}{
		Id:       partner.id,
		Slug:     partner.slug,
		Settings: partner.SettingsDefinition(),
		Alias:    (*Alias)(partner),
	})
}

func (partner *Partner) Definition() *APIPartner {
	return &APIPartner{
		Id:             partner.id,
		Slug:           partner.slug,
		Name:           partner.Name,
		Settings:       partner.SettingsDefinition(),
		ConnectorTypes: partner.ConnectorTypes,
		factories:      make(map[string]ConnectorFactory),
		Errors:         e.NewErrors(),
		manager:        partner.manager,
	}
}

func (partner *Partner) Stop() {
	for _, connector := range partner.connectors {
		c, ok := connector.(state.Stopable)
		if ok {
			c.Stop()
		}
	}
	partner.CancelSubscriptions()
	partner.ResetCollectSettings()
	partner.gtfsCache.Clear()
	partner.httpClient = nil

	partner.limiters = make(map[string]*rate.Limiter)
}

func (partner *Partner) Start() {
	partner.startedAt = partner.manager.Referential().Clock().Now()
	partner.lastDiscovery = time.Time{}
	partner.alternativeStatusCheck.LastCheck = time.Time{}
	partner.alternativeStatusCheck.Status = OPERATIONNAL_STATUS_UNKNOWN

	for _, connector := range partner.connectors {
		c, ok := connector.(state.Startable)
		if ok {
			c.Start()
		}
	}

	to := partner.GtfsCacheTimeout()
	partner.gtfsCache.Add("trip-updates", to, nil)
	partner.gtfsCache.Add("vehicle-positions", to, nil)
	partner.gtfsCache.Add("trip-updates,vehicle-position", to, nil)

	partner.RefreshRateLimit()
}

func (partner *Partner) Allow(ip string) bool {
	limit := partner.RateLimit()
	if limit == 0 {
		return true
	}

	rl, ok := partner.limiters[ip]
	if !ok {
		logger.Log.Debugf("Create a new limiter for ip %v", ip)
		tick := time.Duration((1.0 / limit) * float64(time.Minute))
		burst := int(math.Min(limit, 10))
		limiter := rate.NewLimiter(rate.Every(tick), burst)
		partner.limiters[ip] = limiter
		rl = limiter
	}
	return rl.Allow()
}

func (partner *Partner) CanCollect(stopId string, lineIds map[string]struct{}) bool {
	if partner.CollectSettings().Empty() {
		return true
	}

	if partner.CollectSettings().UseDiscoveredSA || partner.CollectSettings().UseDiscoveredLines {
		return partner.checkDiscovered(stopId, lineIds) // Check excluded stops and lines
	}

	return partner.CollectSettings().CanCollectStop(stopId) && partner.CollectSettings().CanCollectLines(lineIds)
}

func (partner *Partner) CanCollectLine(lineId string) bool { // Used for vehicle collect
	if partner.CollectSettings().UseDiscoveredLines {
		return partner.checkDiscoveredLine(lineId)
	}

	return partner.CollectSettings().CanCollectLine(lineId)
}

func (partner *Partner) checkDiscovered(stopId string, lineIds map[string]struct{}) (ok bool) {
	// Return false if we exclude the stop or all the associated lines
	if partner.CollectSettings().ExcludeStop(stopId) || partner.CollectSettings().ExcludeAllLines(lineIds) {
		return
	}

	partner.mutex.RLock()
	defer partner.mutex.RUnlock()

	if partner.CollectSettings().UseDiscoveredLines {
		for l := range lineIds {
			_, ok = partner.discoveredLines[l]
			if ok {
				return
			}
		}
	}
	if partner.CollectSettings().UseDiscoveredSA {
		_, ok = partner.discoveredStopAreas[stopId]
	}

	return
}

func (partner *Partner) checkDiscoveredLine(lineId string) (ok bool) {
	if partner.CollectSettings().ExcludeLine(lineId) {
		return
	}

	partner.mutex.RLock()
	_, ok = partner.discoveredLines[lineId]
	partner.mutex.RUnlock()

	return
}

// APIPartner.Validate should be called for APIPartner factories to be set
func (partner *Partner) SetDefinition(apiPartner *APIPartner) {
	partner.id = apiPartner.Id
	partner.slug = apiPartner.Slug
	partner.Name = apiPartner.Name
	partner.SetSettingsDefinition(apiPartner.Settings)
	partner.ResetCollectSettings()
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
			factory := NewConnectorFactory(connectorType)
			if factory == nil {
				continue
			}
			partner.connectors[connectorType] = factory.CreateConnector(partner)
		}
	}
	partner.cleanConnectors()
}

// Delete from partner.Connectors connectors not in partner.ConnectorTypes
func (partner *Partner) cleanConnectors() {
	sort.Strings(partner.ConnectorTypes)

	for connectorType := range partner.connectors {
		if connectorType == SIRI_SUBSCRIPTION_REQUEST_DISPATCHER && partner.hasSubscribers() {
			continue
		}
		found := sort.SearchStrings(partner.ConnectorTypes, connectorType)
		if found == len(partner.ConnectorTypes) || partner.ConnectorTypes[found] != connectorType {
			delete(partner.connectors, connectorType)
		}
	}
}

func (partner *Partner) hasSubscribers() bool {
	_, ok := partner.connectors[SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER]
	if ok {
		return true
	}
	_, ok = partner.connectors[SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER]
	if ok {
		return true
	}
	_, ok = partner.connectors[SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER]
	if ok {
		return true
	}
	_, ok = partner.connectors[SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER]
	if ok {
		return true
	}

	_, ok = partner.connectors[SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER]
	if ok {
		return true
	}

	_, ok = partner.connectors[SIRI_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER]
	if ok {
		return true
	}

	return ok
}

func (partner *Partner) Connector(connectorType string) (Connector, bool) {
	connector, ok := partner.connectors[connectorType]
	return connector, ok
}

func (partner *Partner) HaveAtLeastOneConnector(connectorTypes []string) bool {
	for i := range connectorTypes {
		if _, present := partner.connectors[connectorTypes[i]]; present {
			return true
		}
	}
	return false
}

func (partner *Partner) CreateSubscriptionRequestDispatcher() {
	partner.connectors[SIRI_SUBSCRIPTION_REQUEST_DISPATCHER] = NewSIRISubscriptionRequestDispatcher(partner)
}

func (partner *Partner) CheckStatusClient() (csc CheckStatusClient) {
	c, ok := partner.connectors[SIRI_CHECK_STATUS_CLIENT_TYPE]
	if !ok {
		c = partner.connectors[TEST_CHECK_STATUS_CLIENT_TYPE]
	}
	csc, _ = c.(CheckStatusClient)
	return
}

func (partner *Partner) GtfsConnectors() (connectors []GtfsConnector, ok bool) {
	c, ok1 := partner.connectors[GTFS_RT_TRIP_UPDATES_BROADCASTER]
	if ok1 {
		connectors = append(connectors, c.(GtfsConnector))
	}
	c, ok2 := partner.connectors[GTFS_RT_VEHICLE_POSITIONS_BROADCASTER]
	if ok2 {
		connectors = append(connectors, c.(GtfsConnector))
	}
	ok = ok1 || ok2

	return
}

func (partner *Partner) GeneralMessageRequestCollector() GeneralMessageRequestCollector {
	client, ok := partner.connectors[SIRI_GENERAL_MESSAGE_REQUEST_COLLECTOR]
	if ok {
		return client.(GeneralMessageRequestCollector)
	}
	return nil
}

func (partner *Partner) GeneralMessageSubscriptionCollector() GeneralMessageSubscriptionCollector {
	client, ok := partner.connectors[SIRI_GENERAL_MESSAGE_SUBSCRIPTION_COLLECTOR]
	if ok {
		return client.(GeneralMessageSubscriptionCollector)
	}
	return nil
}

func (partner *Partner) StopMonitoringRequestCollector() (smrc StopMonitoringRequestCollector) {
	c, ok := partner.connectors[SIRI_STOP_MONITORING_REQUEST_COLLECTOR]
	if !ok {
		c = partner.connectors[TEST_STOP_MONITORING_REQUEST_COLLECTOR]
	}
	smrc, _ = c.(StopMonitoringRequestCollector)
	return
}

func (partner *Partner) LiteStopMonitoringRequestCollector() (lsmrc LiteStopMonitoringRequestCollector) {
	c, ok := partner.connectors[SIRI_LITE_STOP_MONITORING_REQUEST_COLLECTOR]
	if ok {
		return c.(LiteStopMonitoringRequestCollector)
	}

	return nil
}

func (partner *Partner) StopMonitoringSubscriptionCollector() StopMonitoringSubscriptionCollector {
	client, ok := partner.connectors[SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR]
	if ok {
		return client.(StopMonitoringSubscriptionCollector)
	}
	return nil
}

// func (partner *Partner) EstimatedTimetableRequestCollector() (smrc EstimatedTimetableRequestCollector) {
// 	client, ok := partner.connectors[SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_COLLECTOR]
// 	if ok {
// 		return client.(EstimatedTimetableRequestCollector)
// 	}
// 	return nil
// }

func (partner *Partner) EstimatedTimetableSubscriptionCollector() EstimatedTimetableSubscriptionCollector {
	client, ok := partner.connectors[SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_COLLECTOR]
	if ok {
		return client.(EstimatedTimetableSubscriptionCollector)
	}
	return nil
}

func (partner *Partner) VehicleMonitoringRequestCollector() VehicleMonitoringRequestCollector {
	client, ok := partner.connectors[SIRI_VEHICLE_MONITORING_REQUEST_COLLECTOR]
	if ok {
		return client.(VehicleMonitoringRequestCollector)
	}
	return nil
}

func (partner *Partner) VehicleMonitoringSubscriptionCollector() VehicleMonitoringSubscriptionCollector {
	client, ok := partner.connectors[SIRI_VEHICLE_MONITORING_SUBSCRIPTION_COLLECTOR]
	if ok {
		return client.(VehicleMonitoringSubscriptionCollector)
	}
	return nil
}

func (partner *Partner) hasPushCollector() (ok bool) {
	_, ok = partner.connectors[PUSH_COLLECTOR]
	return ok
}

func (partner *Partner) hasGtfsCollector() (ok bool) {
	_, ok = partner.connectors[GTFS_RT_REQUEST_COLLECTOR]
	return ok
}

func (partner *Partner) alternativeStatusConnector() string {
	if partner.hasPushCollector() {
		return "push"
	}
	if partner.hasGtfsCollector() {
		return "gtfs"
	}
	return "none"
}

func (partner *Partner) CheckStatus() (PartnerStatus, error) {
	logger.Log.Debugf("Check '%s' partner status", partner.slug)
	partnerStatus := PartnerStatus{}

	if partner.CheckStatusClient() == nil {
		switch partner.alternativeStatusConnector() {
		case "push":
			return partner.checkPushStatus()
		case "gtfs":
			return partner.checkGtfsStatus()
		default:
			logger.Log.Debugf("Can't define Status for partner %v", partner.slug)
			partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UNKNOWN
			return partnerStatus, errors.New("no way to define status")
		}
	}
	partnerStatus, err := partner.CheckStatusClient().Status()

	if err != nil {
		logger.Log.Printf("Error while checking %s partner status: %v", partner.Slug(), err)
	}
	logger.Log.Debugf("Partner %v status is %v", partner.slug, partnerStatus.OperationnalStatus)
	return partnerStatus, nil
}

func (partner *Partner) checkPushStatus() (partnerStatus PartnerStatus, _ error) {
	logger.Log.Debugf("Checking %v partner status with PushNotifications", partner.slug)
	if partner.alternativeStatusCheck.LastCheck.Before(partner.manager.Referential().Clock().Now().Add(-5 * time.Minute)) {
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_DOWN
	} else {
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_UP
	}
	logger.Log.Debugf("Partner %v status is %v", partner.slug, partnerStatus.OperationnalStatus)
	return partnerStatus, nil
}

func (partner *Partner) checkGtfsStatus() (partnerStatus PartnerStatus, _ error) {
	logger.Log.Debugf("Checking %v partner status with Gtfs collect", partner.slug)
	if partner.alternativeStatusCheck.LastCheck.Before(partner.manager.Referential().Clock().Now().Add(-5 * time.Minute)) {
		partnerStatus.OperationnalStatus = OPERATIONNAL_STATUS_DOWN
	} else {
		partnerStatus.OperationnalStatus = partner.alternativeStatusCheck.Status
	}
	logger.Log.Debugf("Partner %v status is %v", partner.slug, partnerStatus.OperationnalStatus)
	return partnerStatus, nil
}

func (partner *Partner) Model() model.Model {
	return partner.manager.Model()
}

func (partner *Partner) CancelSubscriptions() {
	partner.subscriptionManager.CancelSubscriptions()
}

func (partner *Partner) CancelBroadcastSubscriptions() {
	partner.subscriptionManager.CancelBroadcastSubscriptions()
}

func (partner *Partner) LastDiscovery() time.Time {
	return partner.lastDiscovery
}

func (partner *Partner) Discover() {
	partner.lastDiscovery = partner.manager.Referential().Clock().Now()
	partner.lineDiscovery()
	partner.stopDiscovery()
}

func (partner *Partner) stopDiscovery() {
	logger.Log.Debugf("StopDiscovery for partner '%s'", partner.slug)

	c, ok := partner.connectors[SIRI_STOP_POINTS_DISCOVERY_REQUEST_COLLECTOR]
	if !ok {
		logger.Log.Debugf("No SiriStopPointsDiscoveryRequestCollector found for partner '%s'", partner.slug)
		return
	}

	c.(StopPointsDiscoveryRequestCollector).RequestStopPoints()
}

func (partner *Partner) lineDiscovery() {
	logger.Log.Debugf("LineDiscovery for partner '%s'", partner.slug)

	c, ok := partner.connectors[SIRI_LINES_DISCOVERY_REQUEST_COLLECTOR]
	if !ok {
		logger.Log.Debugf("No SiriLinesDiscoveryRequestCollector found for partner '%s'", partner.slug)
		return
	}

	c.(LinesDiscoveryRequestCollector).RequestLines()
}

func (partner *Partner) Pushed() {
	partner.alternativeStatusCheck.LastCheck = partner.manager.Referential().Clock().Now()
}

func (partner *Partner) GtfsStatus(s OperationnalStatus) {
	partner.alternativeStatusCheck.LastCheck = partner.manager.Referential().Clock().Now()
	partner.alternativeStatusCheck.Status = s
}

func (partner *Partner) RegisterDiscoveredStopAreas(stops []string) {
	if !partner.CollectSettings().UseDiscoveredSA {
		return
	}

	partner.mutex.Lock()

	for i := range stops {
		partner.discoveredStopAreas[stops[i]] = struct{}{}
	}

	partner.mutex.Unlock()
}

func (partner *Partner) RegisterDiscoveredLines(lines []string) {
	if !partner.CollectSettings().UseDiscoveredLines {
		return
	}

	partner.mutex.Lock()

	for i := range lines {
		partner.discoveredLines[lines[i]] = struct{}{}
	}

	partner.mutex.Unlock()
}

func NewPartnerManager(referential *Referential) *PartnerManager {
	manager := &PartnerManager{
		mutex:                 &sync.RWMutex{},
		byId:                  make(map[PartnerId]*Partner),
		localCredentialsIndex: NewLocalCredentialsIndex(),
		referential:           referential,
	}
	manager.guardian = NewPartnersGuardian(referential)
	return manager
}

func (manager *PartnerManager) Guardian() *PartnersGuardian {
	return manager.guardian
}

func (manager *PartnerManager) UniqCredentials(modelId PartnerId, localCredentials string) bool {
	return manager.localCredentialsIndex.UniqCredentials(modelId, localCredentials)
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
		mutex:               &sync.RWMutex{},
		slug:                slug,
		manager:             manager,
		connectors:          make(map[string]Connector),
		discoveredStopAreas: make(map[string]struct{}),
		discoveredLines:     make(map[string]struct{}),
		PartnerStatus: PartnerStatus{
			OperationnalStatus: OPERATIONNAL_STATUS_UNKNOWN,
		},
		ConnectorTypes: []string{},
		gtfsCache:      cache.NewCacheTable(),
		limiters:       make(map[string]*rate.Limiter),
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator)
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	return partner
}

func (manager *PartnerManager) MarshalJSON() ([]byte, error) {
	partnersId := []PartnerId{}
	for id := range manager.byId {
		partnersId = append(partnersId, id)
	}
	return json.Marshal(partnersId)
}

func (manager *PartnerManager) Find(id PartnerId) *Partner {
	manager.mutex.RLock()
	partner := manager.byId[id]
	manager.mutex.RUnlock()
	return partner
}

func (manager *PartnerManager) FindBySetting(setting, value string) (*Partner, bool) {
	manager.mutex.RLock()
	for _, partner := range manager.byId {
		if partner.Setting(setting) == value {
			manager.mutex.RUnlock()
			return partner, true
		}
	}

	manager.mutex.RUnlock()
	return nil, false
}

func (manager *PartnerManager) FindBySlug(slug PartnerSlug) (*Partner, bool) {
	manager.mutex.RLock()
	for _, partner := range manager.byId {
		if partner.slug == slug {
			manager.mutex.RUnlock()
			return partner, true
		}
	}

	manager.mutex.RUnlock()
	return nil, false
}

func (manager *PartnerManager) FindByCredential(c string) (*Partner, bool) {
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()

	id, ok := manager.localCredentialsIndex.Find(c)
	if !ok {
		return nil, false
	}
	partner, ok := manager.byId[id]
	return partner, ok
}

func (manager *PartnerManager) FindAll() (partners []*Partner) {
	manager.mutex.RLock()
	for _, partner := range manager.byId {
		partners = append(partners, partner)
	}
	manager.mutex.RUnlock()
	return
}

func (manager *PartnerManager) FindAllByCollectPriority() []*Partner {
	partners := manager.FindAll()
	sort.Sort(ByPriority(partners))
	return partners
}

func (manager *PartnerManager) FindAllWithConnector(connectorTypes []string) (partners []*Partner) {
	manager.mutex.RLock()
	for _, partner := range manager.byId {
		if partner.HaveAtLeastOneConnector(connectorTypes) {
			partners = append(partners, partner)
		}
	}
	manager.mutex.RUnlock()
	return
}

func (manager *PartnerManager) Save(partner *Partner) bool {
	if partner.id == "" {
		partner.id = PartnerId(manager.NewUUID())
	}
	partner.manager = manager

	manager.mutex.Lock()
	manager.byId[partner.id] = partner
	manager.localCredentialsIndex.Index(partner.id, partner.Credentials())
	manager.mutex.Unlock()

	return true
}

func (manager *PartnerManager) Delete(partner *Partner) bool {
	manager.mutex.Lock()
	manager.localCredentialsIndex.Delete(partner.id)
	delete(manager.byId, partner.id)
	manager.mutex.Unlock()

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
	manager.mutex.Lock()
	for _, partner := range manager.byId {
		partner.CancelSubscriptions()
	}
	manager.mutex.Unlock()
}

func (manager *PartnerManager) Load() error {
	selectPartners := []model.SelectPartner{}
	sqlQuery := fmt.Sprintf("select * from partners where referential_id = '%s'", manager.referential.Id())
	_, err := model.Database.Select(&selectPartners, sqlQuery)
	if err != nil {
		return err
	}
	for _, p := range selectPartners {
		partner := manager.New(PartnerSlug(p.Slug))
		partner.id = PartnerId(p.Id)
		if p.Name.Valid {
			partner.Name = p.Name.String
		}

		if p.Settings.Valid && len(p.Settings.String) > 0 {
			m := make(map[string]string)
			if err = json.Unmarshal([]byte(p.Settings.String), &m); err != nil {
				return err
			}
			partner.SetSettingsDefinition(m)
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

func (manager *PartnerManager) SaveToDatabase() (int, error) {
	// Check presence of Referential
	selectReferentials := []model.SelectReferential{}
	sqlQuery := fmt.Sprintf("select * from referentials where referential_id = '%s'", manager.referential.Id())
	_, err := model.Database.Select(&selectReferentials, sqlQuery)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}
	if len(selectReferentials) == 0 {
		return http.StatusNotAcceptable, errors.New("can't save Partners without Referential in Database")
	}

	// Begin transaction
	tx, err := model.Database.Begin()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	// Delete partners
	sqlQuery = fmt.Sprintf("delete from partners where referential_id = '%s';", manager.referential.Id())
	_, err = tx.Exec(sqlQuery)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	// Insert partners
	manager.mutex.RLock()
	defer manager.mutex.RUnlock()
	for _, partner := range manager.byId {
		dbPartner, err := manager.newDbPartner(partner)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, fmt.Errorf("internal error: %v", err)
		}
		err = tx.Insert(dbPartner)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, fmt.Errorf("internal error: %v", err)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("database error: %v", err)
	}

	return http.StatusOK, nil
}

func (manager *PartnerManager) newDbPartner(partner *Partner) (*model.DatabasePartner, error) {
	settings, err := partner.PartnerSettings.ToJson()
	if err != nil {
		return nil, err
	}
	connectors, err := json.Marshal(partner.ConnectorTypes)
	if err != nil {
		return nil, err
	}
	return &model.DatabasePartner{
		Id:             string(partner.id),
		ReferentialId:  string(manager.referential.id),
		Slug:           string(partner.slug),
		Name:           partner.Name,
		Settings:       string(settings),
		ConnectorTypes: string(connectors),
	}, nil
}
