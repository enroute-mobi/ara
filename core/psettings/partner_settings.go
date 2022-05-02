package psettings

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"bitbucket.org/enroute-mobi/ara/cache"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

const (
	LOCAL_CREDENTIAL  = "local_credential"
	LOCAL_CREDENTIALS = "local_credentials"
	LOCAL_URL         = "local_url"

	PARTNER_MAX_RETRY = "partner.status.maximum_retry"

	REMOTE_CREDENTIAL                    = "remote_credential"
	REMOTE_OBJECTID_KIND                 = "remote_objectid_kind"
	VEHICLE_REMOTE_OBJECTID_KIND         = "vehicle_remote_objectid_kind"
	VEHICLE_JOURNEY_REMOTE_OBJECTID_KIND = "vehicle_journey_remote_objectid_kind"
	REMOTE_URL                           = "remote_url"
	NOTIFICATIONS_REMOTE_URL             = "notifications.remote_url"
	SUBSCRIPTIONS_REMOTE_URL             = "subscriptions.remote_url"

	COLLECT_PRIORITY                 = "collect.priority"
	COLLECT_INCLUDE_LINES            = "collect.include_lines"
	COLLECT_EXCLUDE_LINES            = "collect.exclude_lines"
	COLLECT_INCLUDE_STOP_AREAS       = "collect.include_stop_areas"
	COLLECT_EXCLUDE_STOP_AREAS       = "collect.exclude_stop_areas"
	COLLECT_USE_DISCOVERED_SA        = "collect.use_discovered_stop_areas"
	COLLECT_SUBSCRIPTIONS_PERSISTENT = "collect.subscriptions.persistent"
	COLLECT_FILTER_GENERAL_MESSAGES  = "collect.filter_general_messages"
	COLLECT_GTFS_TTL                 = "collect.gtfs.ttl"
	COLLECT_DEFAULT_SRS_NAME         = "collect.default_srs_name"

	DISCOVERY_INTERVAL = "discovery_interval"

	BROADCAST_SUBSCRIPTIONS_PERSISTENT         = "broadcast.subscriptions.persistent"
	BROADCAST_REWRITE_JOURNEY_PATTERN_REF      = "broadcast.rewrite_journey_pattern_ref"
	BROADCAST_NO_DESTINATIONREF_REWRITING_FROM = "broadcast.no_destinationref_rewriting_from"
	BROADCAST_NO_DATAFRAMEREF_REWRITING_FROM   = "broadcast.no_dataframeref_rewriting_from"
	SEND_PRODUCER_UNAVAILABLE_ERROR            = "broadcast.send_producer_unavailable_error"
	BROADCAST_GZIP_GTFS                        = "broadcast.gzip_gtfs"
	BROADCAST_GTFS_CACHE_TIMEOUT               = "broadcast.gtfs.cache_timeout"

	IGNORE_STOP_WITHOUT_LINE        = "ignore_stop_without_line"
	GENEREAL_MESSAGE_REQUEST_2_2    = "generalMessageRequest.version2.2"
	SUBSCRIPTIONS_MAXIMUM_RESOURCES = "subscriptions.maximum_resources"

	CACHE_TIMEOUT = "cache_timeout"

	OAUTH_CLIENT_ID     = "remote_authentication.oauth.client_id"
	OAUTH_CLIENT_SECRET = "remote_authentication.oauth.client_secret"
	OAUTH_TOKEN_URL     = "remote_authentication.oauth.token_url"

	SIRI_ENVELOPE = "siri.envelope"

	DEFAULT_GTFS_TTL = 30 * time.Second
)

type PartnerSettings struct {
	m *sync.RWMutex

	ug func() uuid.UUIDGenerator

	s  map[string]string
	cs *CollectSettings
	g  map[string]*idgen.IdentifierGenerator
}

func NewPartnerSettings(ug func() uuid.UUIDGenerator) PartnerSettings {
	return PartnerSettings{
		m:  &sync.RWMutex{},
		ug: ug,
		s:  make(map[string]string),
		g:  make(map[string]*idgen.IdentifierGenerator),
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
	s.s = make(map[string]string)
	for k, v := range m {
		s.s[k] = v
	}
	s.reloadSettings()
	s.m.Unlock()
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

func (s *PartnerSettings) RemoteObjectIDKind(connectorName ...string) string {
	var cn string
	if len(connectorName) != 0 {
		cn = connectorName[0]
	}

	s.m.RLock()
	defer s.m.RUnlock()

	if setting := s.s[fmt.Sprintf("%s.%s", cn, REMOTE_OBJECTID_KIND)]; setting != "" {
		return setting
	}
	return s.s[REMOTE_OBJECTID_KIND]
}

func (s *PartnerSettings) VehicleRemoteObjectIDKindWithFallback(connectorName ...string) []string {
	return s.remoteObjectIDKindWithFallback(VEHICLE_REMOTE_OBJECTID_KIND, connectorName...)
}

func (s *PartnerSettings) VehicleJourneyRemoteObjectIDKindWithFallback(connectorName ...string) []string {
	return s.remoteObjectIDKindWithFallback(VEHICLE_JOURNEY_REMOTE_OBJECTID_KIND, connectorName...)
}

func (s *PartnerSettings) remoteObjectIDKindWithFallback(settingName string, connectorName ...string) (k []string) {
	var cn string
	if len(connectorName) != 0 {
		cn = connectorName[0]
	}

	s.m.RLock()

	if setting := s.s[fmt.Sprintf("%s.%s", cn, settingName)]; setting != "" {
		k = append(k, trimedSlice(setting)...)
	}
	if setting := s.s[settingName]; setting != "" {
		k = append(k, trimedSlice(setting)...)
	}

	if len(k) == 0 {
		if setting := s.s[fmt.Sprintf("%s.%s", cn, REMOTE_OBJECTID_KIND)]; setting != "" {
			k = append(k, setting)
		} else {
			k = append(k, s.s[REMOTE_OBJECTID_KIND])
		}
	}

	s.m.RUnlock()
	return
}

func (s *PartnerSettings) GtfsTTL() (t time.Duration) {
	s.m.RLock()
	t, _ = time.ParseDuration(s.s[COLLECT_GTFS_TTL])
	s.m.RUnlock()
	if t < DEFAULT_GTFS_TTL {
		t = DEFAULT_GTFS_TTL
	}

	return
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

func (s *PartnerSettings) RequestorRef() string {
	return s.ProducerRef()
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

func (s *PartnerSettings) SIRIEnvelopeType() (set string) {
	s.m.RLock()
	set = s.s[SIRI_ENVELOPE]
	s.m.RUnlock()

	if set == "" {
		set = "soap"
	}
	return
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

func (s *PartnerSettings) SubscriptionMaximumResources() (i int) {
	s.m.RLock()
	i, _ = strconv.Atoi(s.s[SUBSCRIPTIONS_MAXIMUM_RESOURCES])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) CollectPriority() (value int) {
	s.m.RLock()
	value, _ = strconv.Atoi(s.s[COLLECT_PRIORITY])
	s.m.RUnlock()
	return value
}

func (s *PartnerSettings) DefaultSRSName() (srsName string) {
	s.m.RLock()
	if s.s[COLLECT_DEFAULT_SRS_NAME] == "" {
		srsName = "EPSG:2154"
	} else {
		srsName = s.s[COLLECT_DEFAULT_SRS_NAME]
	}
	s.m.RUnlock()
	return srsName
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

func (s *PartnerSettings) GzipGtfs() (r bool) {
	s.m.RLock()
	r, _ = strconv.ParseBool(s.s[BROADCAST_GZIP_GTFS])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) GeneralMessageRequestVersion22() (r bool) {
	s.m.RLock()
	r, _ = strconv.ParseBool(s.s[GENEREAL_MESSAGE_REQUEST_2_2])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) PersistentCollectSubscriptions() (r bool) {
	s.m.RLock()
	r, _ = strconv.ParseBool(s.s[COLLECT_SUBSCRIPTIONS_PERSISTENT])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) PersistentBroadcastSubscriptions() (r bool) {
	s.m.RLock()
	r, _ = strconv.ParseBool(s.s[BROADCAST_SUBSCRIPTIONS_PERSISTENT])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) CollectFilteredGeneralMessages() (r bool) {
	s.m.RLock()
	r, _ = strconv.ParseBool(s.s[COLLECT_FILTER_GENERAL_MESSAGES])
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) IgnoreStopWithoutLine() (r bool) {
	s.m.RLock()
	r = s.s[IGNORE_STOP_WITHOUT_LINE] != "false"
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) SendProducerUnavailableError() (r bool) {
	s.m.RLock()
	r, _ = strconv.ParseBool(s.s[SEND_PRODUCER_UNAVAILABLE_ERROR])
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
		s.SetCollectSettings()
		s.m.RUnlock()
	}

	return s.cs
}

// Warning, this method isn't threadsafe. Mutex must be handled before and after calling
func (s *PartnerSettings) SetCollectSettings() {
	s.cs = &CollectSettings{
		UseDiscovered: s.s[COLLECT_USE_DISCOVERED_SA] != "",
		includedSA:    toMap(s.s[COLLECT_INCLUDE_STOP_AREAS]),
		excludedSA:    toMap(s.s[COLLECT_EXCLUDE_STOP_AREAS]),
		includedLines: toMap(s.s[COLLECT_INCLUDE_LINES]),
		excludedLines: toMap(s.s[COLLECT_EXCLUDE_LINES]),
	}
}

func (s *PartnerSettings) HTTPClientOptions() (opts remote.HTTPClientOptions) {
	s.m.RLock()
	opts = remote.HTTPClientOptions{
		SiriEnvelopeType: s.siriEnvelopeType(),
		OAuth:            s.httpClientOAuth(),
		Urls: remote.HTTPClientUrls{
			Url:              s.s[REMOTE_URL],
			SubscriptionsUrl: s.s[SUBSCRIPTIONS_REMOTE_URL],
			NotificationsUrl: s.s[NOTIFICATIONS_REMOTE_URL],
		},
	}
	s.m.RUnlock()
	return
}

func (s *PartnerSettings) siriEnvelopeType() (set string) {
	set = s.s[SIRI_ENVELOPE]
	if set == "" {
		set = "soap"
	}

	return set
}

// Warning, this method isn't threadsafe. Mutex must be handled before and after calling
func (s *PartnerSettings) httpClientOAuth() (opts *remote.HTTPClientOAuth) {
	cid, ok1 := s.s[OAUTH_CLIENT_ID]
	cs, ok2 := s.s[OAUTH_CLIENT_SECRET]
	t, ok3 := s.s[OAUTH_TOKEN_URL]
	if ok1 && ok2 && ok3 {
		opts = &remote.HTTPClientOAuth{
			ClientID:     cid,
			ClientSecret: cs,
			TokenURL:     t,
		}
	}
	return
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

func toMap(s string) (m map[string]struct{}) {
	m = make(map[string]struct{})
	if s == "" {
		return
	}
	t := strings.Split(s, ",")
	for i := range t {
		m[strings.TrimSpace(t[i])] = struct{}{}
	}
	return
}

func (s *PartnerSettings) IdentifierGenerator(generatorName string) *idgen.IdentifierGenerator {
	s.m.Lock()
	generator, ok := s.g[generatorName]
	if !ok {
		generator = idgen.NewIdentifierGenerator(s.idGeneratorFormat(generatorName), s.ug())
		s.g[generatorName] = generator
	}
	s.m.Unlock()
	return generator
}

func (s *PartnerSettings) NewMessageIdentifier() string {
	return s.IdentifierGenerator(idgen.MESSAGE_IDENTIFIER).NewMessageIdentifier()
}

func (s *PartnerSettings) NewResponseMessageIdentifier() string {
	return s.IdentifierGenerator(idgen.RESPONSE_MESSAGE_IDENTIFIER).NewMessageIdentifier()
}

func (s *PartnerSettings) idGeneratorFormat(generatorName string) (formatString string) {
	formatString = s.s[fmt.Sprintf("generators.%v", generatorName)]

	if formatString == "" {
		formatString = idgen.DefaultIdentifierGenerator(generatorName)
	}
	return
}

// Warning, this method isn't threadsafe. Mutex must be handled before and after calling
func (s *PartnerSettings) refreshGenerators() {
	s.g = make(map[string]*idgen.IdentifierGenerator)
}

// Warning, this method isn't threadsafe. Mutex must be handled before and after calling
func (s *PartnerSettings) reloadSettings() {
	s.SetCollectSettings()
	s.refreshGenerators()
}
