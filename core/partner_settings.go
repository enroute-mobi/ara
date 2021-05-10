package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/cache"
)

const (
	LOCAL_CREDENTIAL  = "local_credential"
	LOCAL_CREDENTIALS = "local_credentials"
	LOCAL_URL         = "local_url"

	PARTNER_MAX_RETRY = "partner.status.maximum_retry"

	REMOTE_CREDENTIAL            = "remote_credential"
	REMOTE_OBJECTID_KIND         = "remote_objectid_kind"
	VEHICLE_REMOTE_OBJECTID_KIND = "vehicle_remote_objectid_kind"
	REMOTE_URL                   = "remote_url"
	NOTIFICATIONS_REMOTE_URL     = "notifications.remote_url"
	SUBSCRIPTIONS_REMOTE_URL     = "subscriptions.remote_url"

	COLLECT_PRIORITY                 = "collect.priority"
	COLLECT_INCLUDE_LINES            = "collect.include_lines"
	COLLECT_INCLUDE_STOP_AREAS       = "collect.include_stop_areas"
	COLLECT_EXCLUDE_STOP_AREAS       = "collect.exclude_stop_areas"
	COLLECT_USE_DISCOVERED_SA        = "collect.use_discovered_stop_areas"
	COLLECT_SUBSCRIPTIONS_PERSISTENT = "collect.subscriptions.persistent"
	COLLECT_FILTER_GENERAL_MESSAGES  = "collect.filter_general_messages"

	DISCOVERY_INTERVAL = "discovery_interval"

	BROADCAST_SUBSCRIPTIONS_PERSISTENT         = "broadcast.subscriptions.persistent"
	BROADCAST_REWRITE_JOURNEY_PATTERN_REF      = "broadcast.rewrite_journey_pattern_ref"
	BROADCAST_NO_DESTINATIONREF_REWRITING_FROM = "broadcast.no_destinationref_rewriting_from"
	BROADCAST_NO_DATAFRAMEREF_REWRITING_FROM   = "broadcast.no_dataframeref_rewriting_from"
	BROADCAST_GZIP_GTFS                        = "broadcast.gzip_gtfs"
	BROADCAST_GTFS_CACHE_TIMEOUT               = "broadcast.gtfs.cache_timeout"

	IGNORE_STOP_WITHOUT_LINE        = "ignore_stop_without_line"
	GENEREAL_MESSAGE_REQUEST_2      = "generalMessageRequest.version2.2"
	SUBSCRIPTIONS_MAXIMUM_RESOURCES = "subscriptions.maximum_resources"

	LOGSTASH_LOG_DELIVERIES_IN_SM_COLLECT_NOTIFICATIONS = "logstash.log_deliveries_in_sm_collect_notifications"
	LOGSTASH_LOG_DELIVERIES_IN_SM_COLLECT_REQUESTS      = "logstash.log_deliveries_in_sm_collect_requests"

	CACHE_TIMEOUT = "cache_timeout"
)

type PartnerSettings struct {
	m *sync.RWMutex

	p *Partner

	s  map[string]string
	cs *CollectSettings
	g  map[string]*IdentifierGenerator
}

func NewPartnerSettings(p *Partner) PartnerSettings {
	return PartnerSettings{
		m: &sync.RWMutex{},
		p: p,
		s: make(map[string]string),
		g: make(map[string]*IdentifierGenerator),
	}
}

func (s *PartnerSettings) Setting(key string) string {
	s.m.RLock()
	defer s.m.RUnlock()
	return s.s[key]
}

// Should only be used in tests
func (s *PartnerSettings) SetSetting(k, v string) {
	s.m.Lock()
	s.s[k] = v
	s.reloadSettings()
	s.m.Unlock()
}

func (s *PartnerSettings) SettingsDefinition() (m map[string]string) {
	m = make(map[string]string)
	s.m.RLock()
	for k, v := range s.s {
		m[k] = v
	}
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) SetSettingsDefinition(m map[string]string) {
	if m == nil {
		return
	}
	s.m.Lock()
	for k, v := range m {
		s.s[k] = v
	}
	s.reloadSettings()
	s.m.Unlock()
	return
}

func (s *PartnerSettings) ToJson() ([]byte, error) {
	s.m.RLock()
	defer s.m.RUnlock()
	return json.Marshal(s.s)
}

func (s *PartnerSettings) Credentials() string {
	s.m.RLock()
	defer s.m.RUnlock()

	_, ok := s.s[LOCAL_CREDENTIAL]
	_, ok2 := s.s[LOCAL_CREDENTIALS]
	if !ok && !ok2 {
		return ""
	}
	return fmt.Sprintf("%v,%v", s.s[LOCAL_CREDENTIAL], s.s[LOCAL_CREDENTIALS])
}

func (s *PartnerSettings) RemoteObjectIDKind(connectorName string) string {
	s.m.RLock()
	defer s.m.RUnlock()

	if setting := s.s[fmt.Sprintf("%s.%s", connectorName, REMOTE_OBJECTID_KIND)]; setting != "" {
		return setting
	}
	return s.s[REMOTE_OBJECTID_KIND]
}

func (s *PartnerSettings) VehicleRemoteObjectIDKind(connectorName string) string {
	s.m.RLock()
	defer s.m.RUnlock()

	if setting := s.s[fmt.Sprintf("%s.%s", connectorName, VEHICLE_REMOTE_OBJECTID_KIND)]; setting != "" {
		return setting
	}
	return s.s[REMOTE_OBJECTID_KIND]
}

// Very specific for now, we'll refacto if we need to cache more
func (s *PartnerSettings) GtfsCacheTimeout() (t time.Duration) {
	s.m.RLock()
	t, _ = time.ParseDuration(s.s[BROADCAST_GTFS_CACHE_TIMEOUT])
	s.m.RUnlock()

	if t < cache.MIN_CACHE_LIFESPAN {
		t = cache.DEFAULT_CACHE_LIFESPAN
	}
	return
}

func (s *PartnerSettings) CacheTimeout(connectorName string) (t time.Duration) {
	s.m.RLock()
	t, _ = time.ParseDuration(s.s[fmt.Sprintf("%s.%s", connectorName, CACHE_TIMEOUT)])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) ProducerRef() (producerRef string) {
	s.m.RLock()
	producerRef = s.s[REMOTE_CREDENTIAL]
	s.m.RUnlock()

	if producerRef == "" {
		producerRef = "Ara"
	}
	return producerRef
}

// Ref Issue #4300
func (s *PartnerSettings) Address() string {
	// address := s.s("local_url")
	// if address == "" {
	// 	address = config.Config.DefaultAddress
	// }
	// return address
	s.m.RLock()
	defer s.m.RUnlock()
	return s.s[LOCAL_URL]
}

func (s *PartnerSettings) MaximumChechstatusRetry() (i int) {
	s.m.RLock()
	i, _ = strconv.Atoi(s.s[PARTNER_MAX_RETRY])
	if i < 0 {
		i = 0
	}
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) CollectPriority() (value int) {
	s.m.RLock()
	value, _ = strconv.Atoi(s.s[COLLECT_PRIORITY])
	s.m.RUnlock()
	return value
}

func (s *PartnerSettings) NoDestinationRefRewritingFrom() []string {
	s.m.RLock()
	defer s.m.RUnlock()

	return trimedSlice(s.s[BROADCAST_NO_DESTINATIONREF_REWRITING_FROM])
}

func (s *PartnerSettings) NoDataFrameRefRewritingFrom() []string {
	s.m.RLock()
	defer s.m.RUnlock()

	return trimedSlice(s.s[BROADCAST_NO_DATAFRAMEREF_REWRITING_FROM])
}

func (s *PartnerSettings) RewriteJourneyPatternRef() (r bool) {
	s.m.RLock()
	r, _ = strconv.ParseBool(s.s[BROADCAST_REWRITE_JOURNEY_PATTERN_REF])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) LogSubscriptionStopMonitoringDeliveries() (l bool) {
	s.m.RLock()
	l, _ = strconv.ParseBool(s.s[LOGSTASH_LOG_DELIVERIES_IN_SM_COLLECT_NOTIFICATIONS])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) LogRequestStopMonitoringDeliveries() (l bool) {
	s.m.RLock()
	l, _ = strconv.ParseBool(s.s[LOGSTASH_LOG_DELIVERIES_IN_SM_COLLECT_REQUESTS])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) GzipGtfs() (r bool) {
	s.m.RLock()
	r, _ = strconv.ParseBool(s.s[BROADCAST_GZIP_GTFS])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) DiscoveryInterval() (d time.Duration) {
	s.m.RLock()
	d, _ = time.ParseDuration(s.s[DISCOVERY_INTERVAL])
	s.m.RUnlock()
	if d == 0 {
		d = 1 * time.Hour
	}
	return -d
}

func (s *PartnerSettings) CollectSettings() *CollectSettings {
	if s.cs == nil {
		s.m.RLock()
		s.setCollectSettings()
		s.m.RUnlock()
	}

	return s.cs
}

// Warning, this method isn't threadsafe. Mutex must be handled before and after calling
func (s *PartnerSettings) setCollectSettings() {
	s.cs = &CollectSettings{
		UseDiscovered: s.s[COLLECT_USE_DISCOVERED_SA] != "",
		includedSA:    trimedSlice(s.s[COLLECT_INCLUDE_STOP_AREAS]),
		includedLines: trimedSlice(s.s[COLLECT_INCLUDE_LINES]),
		excludedSA:    trimedSlice(s.s[COLLECT_EXCLUDE_STOP_AREAS]),
	}
}

func trimedSlice(s string) (slc []string) {
	if s == "" {
		return
	}
	slc = strings.Split(s, ",")
	for i := range slc {
		slc[i] = strings.TrimSpace(slc[i])
	}
	return
}

func (s *PartnerSettings) IdentifierGenerator(generatorName string) *IdentifierGenerator {
	s.m.Lock()
	generator, ok := s.g[generatorName]
	if !ok {
		generator = NewIdentifierGenerator(s.idGeneratorFormat(generatorName), s.p.UUIDConsumer)
		s.g[generatorName] = generator
	}
	s.m.Unlock()
	return generator
}

func (s *PartnerSettings) idGeneratorFormat(generatorName string) (formatString string) {
	formatString = s.s[fmt.Sprintf("generators.%v", generatorName)]

	if formatString == "" {
		formatString = DefaultIdentifierGenerator(generatorName)
	}
	return
}

// Warning, this method isn't threadsafe. Mutex must be handled before and after calling
func (s *PartnerSettings) refreshGenerators() {
	s.g = make(map[string]*IdentifierGenerator)
}

// Warning, this method isn't threadsafe. Mutex must be handled before and after calling
func (s *PartnerSettings) reloadSettings() {
	s.setCollectSettings()
	s.refreshGenerators()
}
