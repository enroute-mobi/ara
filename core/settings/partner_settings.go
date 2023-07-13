package settings

import (
	"fmt"
	"regexp"
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
	RATE_LIMIT_PER_IP = "rate_limit_per_ip"

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
	COLLECT_USE_DISCOVERED_LINES     = "collect.use_discovered_lines"
	COLLECT_SUBSCRIPTIONS_PERSISTENT = "collect.subscriptions.persistent"
	COLLECT_PERSISTENT               = "collect.persistent"
	COLLECT_FILTER_GENERAL_MESSAGES  = "collect.filter_general_messages"
	COLLECT_GTFS_TTL                 = "collect.gtfs.ttl"
	COLLECT_DEFAULT_SRS_NAME         = "collect.default_srs_name"

	DISCOVERY_INTERVAL = "discovery_interval"

	BROADCAST_SUBSCRIPTIONS_PERSISTENT         = "broadcast.subscriptions.persistent"
	BROADCAST_RECORDED_CALLS_DURATION          = "broadcast.recorded_calls.duration"
	BROADCAST_REWRITE_JOURNEY_PATTERN_REF      = "broadcast.rewrite_journey_pattern_ref"
	BROADCAST_NO_DESTINATIONREF_REWRITING_FROM = "broadcast.no_destinationref_rewriting_from"
	BROADCAST_NO_DATAFRAMEREF_REWRITING_FROM   = "broadcast.no_dataframeref_rewriting_from"
	BROADCAST_GZIP_GTFS                        = "broadcast.gzip_gtfs"
	BROADCAST_GTFS_CACHE_TIMEOUT               = "broadcast.gtfs.cache_timeout"

	IGNORE_STOP_WITHOUT_LINE        = "ignore_stop_without_line"
	GENERAL_MESSAGE_REQUEST_2_2     = "generalMessageRequest.version2.2"
	SUBSCRIPTIONS_MAXIMUM_RESOURCES = "subscriptions.maximum_resources"

	CACHE_TIMEOUT = "cache_timeout"

	OAUTH_CLIENT_ID     = "remote_authentication.oauth.client_id"
	OAUTH_CLIENT_SECRET = "remote_authentication.oauth.client_secret"
	OAUTH_TOKEN_URL     = "remote_authentication.oauth.token_url"

	SIRI_ENVELOPE                                         = "siri.envelope"
	SIRI_LINE_PUBLISHED_NAME                              = "siri.line.published_name"
	SIRI_DIRECTION_TYPE                                   = "siri.direction_type"
	SIRI_PASSAGE_ORDER                                    = "siri.passage_order"
	SIRI_CREDENTIAL_HEADER                                = "siri.credential.header"
	DEFAULT_GTFS_TTL                                      = 30 * time.Second
	BROADCAST_SIRI_IGNORE_TERMINATE_SUBSCRIPTION_REQUESTS = "broadcast.siri.ignore_terminate_subscription_requests"

	SORT_PAYLOAD_FOR_TEST = "sort_payload_for_test"
)

type PartnerSettings struct {
	ug func() uuid.UUIDGenerator

	collectSettings *CollectSettings

	credentials                    string
	rateLimit                      float64
	gtfsTTL                        time.Duration
	gtfsCacheTimeout               time.Duration
	siriCredentialHeader           string
	siriEnvelopeType               string
	httpClientOAuth                *remote.HTTPClientOAuth
	recordedCallsDuration          time.Duration
	producerRef                    string
	address                        string
	linePublishedName              string
	passageOrder                   string
	envelopeType                   string
	collectPriority                int
	defaultSRSName                 string
	noDestinationRefRewritingFrom  []string
	noDataFrameRefRewritingFrom    []string
	rewriteJourneyPatternRef       bool
	gzipGtfs                       bool
	generalMessageRequestVersion22 bool
	collectFilteredGeneralMessages bool
	ignoreStopWithoutLine          bool
	discoveryInterval              time.Duration
	cacheTimeouts                  sync.Map
	cacheTimeout                   time.Duration
	siriDirectionTypeInbound       string
	siriDirectionTypeOutbound      string

	maximumCheckstatusRetry          int
	subscriptionMaximumResources     int
	persistentCollect                bool
	persistentBroadcastSubscriptions bool

	messageIdentifierGenerator           *idgen.IdentifierGenerator
	responseMessageIdentifierGenerator   *idgen.IdentifierGenerator
	dataFrameIdentifierGenerator         *idgen.IdentifierGenerator
	referenceIdentifierGenerator         *idgen.IdentifierGenerator
	referenceStopAreaIdentifierGenerator *idgen.IdentifierGenerator
	subscriptionIdentifierGenerator      *idgen.IdentifierGenerator

	remoteObjectIDKinds sync.Map
	remoteObjectIDKind  string

	httpClientOptions remote.HTTPClientOptions

	sortPayloadForTest                  bool
	ignoreTerminateSubscriptionsRequest bool

	vehicleRemoteObjectIDKinds            []string
	vehicleRemoteObjectIDKindsByConnector sync.Map

	vehicleJourneyRemoteObjectIDKinds            []string
	vehicleJourneyRemoteObjectIDKindsByConnector sync.Map

	// ! Never use these values outside SettingsDefinition()
	originalSettings map[string]string
}

func NewEmptyPartnerSettings(ug func() uuid.UUIDGenerator) PartnerSettings {
	partnerSettings := PartnerSettings{
		ug: ug,
	}
	partnerSettings.parseSettings(map[string]string{})
	return partnerSettings
}

func NewPartnerSettings(generator func() uuid.UUIDGenerator, settings map[string]string) PartnerSettings {
	partnerSettings := PartnerSettings{
		ug: generator,
	}

	partnerSettings.parseSettings(settings)

	return partnerSettings
}

func (s *PartnerSettings) parseSettings(settings map[string]string) {
	s.setRemoteObjectIDKinds(settings)
	s.setCredentials(settings)
	s.setRateLimit(settings)
	s.setGtfsTTL(settings)
	s.setGtfsCacheTimeout(settings)
	s.setRecordedCallsDuration(settings)

	s.setSIRIEnvelopeType(settings)
	s.setAddress(settings)
	s.setProducerRef(settings)
	s.setSIRILinePublishedName(settings)
	s.setSIRIDirectionType(settings)
	s.setCollectSettings(settings)
	s.setSiriCredentialHeader(settings)
	s.setSiriEnvelopeType(settings)
	s.setSIRIPassageOrder(settings)
	s.setMaximumChechstatusRetry(settings)
	s.setSubscriptionMaximumResources(settings)
	s.setCollectPriority(settings)
	s.setDefaultSRSName(settings)
	s.setNoDestinationRefRewritingFrom(settings)
	s.setNoDataFrameRefRewritingFrom(settings)
	s.setGzipGtfs(settings)
	s.setGeneralMessageRequestVersion22(settings)
	s.setPersistentCollect(settings)
	s.setPersistentBroadcastSubscriptions(settings)
	s.setRewriteJourneyPatternRef(settings)
	s.setCollectFilteredGeneralMessages(settings)
	s.setIgnoreStopWithoutLine(settings)
	s.setDiscoveryInterval(settings)
	s.setCacheTimeouts(settings)
	s.setSortPayloadForTest(settings)

	s.setVehicleRemoteObjectIDKindWithFallback(settings)

	s.setIgnoreTerminateSubscriptionsRequest(settings)
	s.setIdentifierGenerators(settings)

	// depends on other settings
	s.setHTTPClientOptions(settings)

	s.originalSettings = settings
}

func (s *PartnerSettings) SettingsDefinition() map[string]string {
	return s.originalSettings
}

func (s *PartnerSettings) setCredentials(settings map[string]string) {
	_, ok := settings[LOCAL_CREDENTIAL]
	_, ok2 := settings[LOCAL_CREDENTIALS]
	if !ok && !ok2 {
		s.credentials = ""
	} else {
		s.credentials = fmt.Sprintf("%v,%v", settings[LOCAL_CREDENTIAL], settings[LOCAL_CREDENTIALS])
	}
}

func (s *PartnerSettings) Credentials() string {
	return s.credentials
}

func (s *PartnerSettings) setRateLimit(settings map[string]string) {
	value, _ := strconv.Atoi(settings[RATE_LIMIT_PER_IP])
	if value < 0 {
		value = 0
	}
	s.rateLimit = float64(value)
}

func (s *PartnerSettings) RateLimit() float64 {
	return s.rateLimit
}

func (s *PartnerSettings) setRemoteObjectIDKinds(settings map[string]string) {
	r, _ := regexp.Compile("(.+)\\.remote_objectid_kind")

	// xxxx.remote_objectid_kind = dummy -> xxxx = dummy
	for key, value := range settings {
		matches := r.FindStringSubmatch(key)
		if len(matches) != 1 {
			break
		}

		connectorName := matches[0]
		s.remoteObjectIDKinds.Store(connectorName, value)
	}

	s.remoteObjectIDKind = settings[REMOTE_OBJECTID_KIND]
}

func (s *PartnerSettings) RemoteObjectIDKind(optionalConnectorName ...string) string {
	if len(optionalConnectorName) == 1 {
		connectorName := optionalConnectorName[0]

		value, ok := s.remoteObjectIDKinds.Load(connectorName)
		if ok {
			return value.(string)
		}
	}

	return s.remoteObjectIDKind
}

func (s *PartnerSettings) setVehicleJourneyRemoteObjectIDKindWithFallback(settings map[string]string) {
	// xxxx.vehicle_journey_remote_objectid_kind
	r, _ := regexp.Compile(fmt.Sprintf("(.+)\\.%s", VEHICLE_JOURNEY_REMOTE_OBJECTID_KIND))

	s.vehicleJourneyRemoteObjectIDKinds = trimedSlice(settings[VEHICLE_JOURNEY_REMOTE_OBJECTID_KIND])

	// xxxx.vehicle_journey_remote_objectid_kind = dummy -> xxxx = dummy
	for key, value := range settings {
		matches := r.FindStringSubmatch(key)
		if len(matches) != 1 {
			break
		}

		var connectorRemoteObjectIDKinds []string
		connectorRemoteObjectIDKinds = append(connectorRemoteObjectIDKinds, trimedSlice(value)...)
		connectorRemoteObjectIDKinds = append(connectorRemoteObjectIDKinds, s.vehicleJourneyRemoteObjectIDKinds...)

		// "a,b,c" -> [a,b,c]

		connectorName := matches[0]
		s.vehicleJourneyRemoteObjectIDKindsByConnector.Store(connectorName, connectorRemoteObjectIDKinds)
	}
}

func (s *PartnerSettings) VehicleJourneyRemoteObjectIDKindWithFallback(connectorName string) []string {
	value, ok := s.vehicleJourneyRemoteObjectIDKindsByConnector.Load(connectorName)
	if ok {
		return value.([]string)
	}

	if len(s.vehicleJourneyRemoteObjectIDKinds) > 0 {
		return s.vehicleJourneyRemoteObjectIDKinds
	}

	return []string{s.RemoteObjectIDKind(connectorName)}
}

func (s *PartnerSettings) setVehicleRemoteObjectIDKindWithFallback(settings map[string]string) {
	// xxxx.vehicle_journey_remote_objectid_kind
	r, _ := regexp.Compile(fmt.Sprintf("(.+)\\.%s", VEHICLE_REMOTE_OBJECTID_KIND))

	s.vehicleRemoteObjectIDKinds = trimedSlice(settings[VEHICLE_REMOTE_OBJECTID_KIND])

	// xxxx.vehicle_journey_remote_objectid_kind = dummy -> xxxx = dummy
	for key, value := range settings {
		matches := r.FindStringSubmatch(key)
		if len(matches) != 1 {
			break
		}

		var connectorRemoteObjectIDKinds []string
		connectorRemoteObjectIDKinds = append(connectorRemoteObjectIDKinds, trimedSlice(value)...)
		connectorRemoteObjectIDKinds = append(connectorRemoteObjectIDKinds, s.vehicleRemoteObjectIDKinds...)

		// "a,b,c" -> [a,b,c]

		connectorName := matches[0]
		s.vehicleRemoteObjectIDKindsByConnector.Store(connectorName, connectorRemoteObjectIDKinds)
	}
}

func (s *PartnerSettings) VehicleRemoteObjectIDKindWithFallback(connectorName string) []string {
	value, ok := s.vehicleRemoteObjectIDKindsByConnector.Load(connectorName)
	if ok {
		return value.([]string)
	}

	if len(s.vehicleRemoteObjectIDKinds) > 0 {
		return s.vehicleRemoteObjectIDKinds
	}

	return []string{s.RemoteObjectIDKind(connectorName)}
}

func (s *PartnerSettings) setGtfsTTL(settings map[string]string) {
	duration, _ := time.ParseDuration(settings[COLLECT_GTFS_TTL])
	if duration < DEFAULT_GTFS_TTL {
		duration = DEFAULT_GTFS_TTL
	}
	s.gtfsTTL = duration
}

func (s *PartnerSettings) GtfsTTL() (t time.Duration) {
	return s.gtfsTTL
}

func (s *PartnerSettings) setRecordedCallsDuration(settings map[string]string) {
	duration, _ := time.ParseDuration(settings[BROADCAST_RECORDED_CALLS_DURATION])
	s.recordedCallsDuration = duration
}
func (s *PartnerSettings) RecordedCallsDuration() (t time.Duration) {
	return s.recordedCallsDuration
}

func (s *PartnerSettings) setGtfsCacheTimeout(settings map[string]string) {
	duration, _ := time.ParseDuration(settings[BROADCAST_GTFS_CACHE_TIMEOUT])

	if duration < cache.MIN_CACHE_LIFESPAN {
		duration = cache.DEFAULT_CACHE_LIFESPAN
	}

	s.gtfsCacheTimeout = duration
}

// Very specific for now, we'll refacto if we need to cache more
func (s *PartnerSettings) GtfsCacheTimeout() (t time.Duration) {
	return s.gtfsCacheTimeout
}

func (s *PartnerSettings) setCacheTimeouts(settings map[string]string) {
	r, _ := regexp.Compile("(.+)\\.cache_timeout")

	// xxxx.cache_timeout = dummy -> xxxx = dummy
	for key, value := range settings {
		matches := r.FindStringSubmatch(key)
		if len(matches) != 1 {
			break
		}

		connectorName := matches[0]
		s.cacheTimeouts.Store(connectorName, value)
	}
	return
}

func (s *PartnerSettings) CacheTimeout(connectorName string) (t time.Duration) {
	value, ok := s.cacheTimeouts.Load(connectorName)
	if ok {
		return value.(time.Duration)
	}

	emptyDuration, _ := time.ParseDuration("")
	return emptyDuration
}

func (s *PartnerSettings) setSortPayloadForTest(settings map[string]string) {
	sortPayload, _ := strconv.ParseBool(settings[SORT_PAYLOAD_FOR_TEST])
	s.sortPayloadForTest = sortPayload
}
func (s *PartnerSettings) SortPaylodForTest() bool {
	return s.sortPayloadForTest
}

func (s *PartnerSettings) setIgnoreTerminateSubscriptionsRequest(settings map[string]string) {
	ignore, _ := strconv.ParseBool(settings[BROADCAST_SIRI_IGNORE_TERMINATE_SUBSCRIPTION_REQUESTS])
	s.ignoreTerminateSubscriptionsRequest = ignore
	return
}

func (s *PartnerSettings) IgnoreTerminateSubscriptionsRequest() bool {
	return s.ignoreTerminateSubscriptionsRequest
}

func (s *PartnerSettings) setProducerRef(settings map[string]string) {
	producerRef := settings[REMOTE_CREDENTIAL]
	if producerRef == "" {
		s.producerRef = "Ara"
	}
	s.producerRef = producerRef
}

func (s *PartnerSettings) ProducerRef() string {
	return s.producerRef
}

func (s *PartnerSettings) RequestorRef() string {
	return s.ProducerRef()
}

func (s *PartnerSettings) setAddress(settings map[string]string) {
	address := settings[LOCAL_URL]
	s.address = address
}

func (s *PartnerSettings) Address() string {
	return s.address
}

func (s *PartnerSettings) setSIRIDirectionType(settings map[string]string) {
	directions := strings.Split(settings[SIRI_DIRECTION_TYPE], ",")
	// ensure the correctness of the setting
	if len(directions) != 2 {
		s.siriDirectionTypeOutbound = "Aller"
		s.siriDirectionTypeInbound = "Retour"
		return
	}

	s.siriDirectionTypeOutbound = directions[0]
	s.siriDirectionTypeInbound = directions[1]
}

func (s *PartnerSettings) SIRIDirectionType() (string, string, bool) {
	return s.siriDirectionTypeInbound, s.siriDirectionTypeOutbound, true
}

func (s *PartnerSettings) setSIRILinePublishedName(settings map[string]string) {
	lineName := settings[SIRI_LINE_PUBLISHED_NAME]
	s.linePublishedName = lineName
}
func (s *PartnerSettings) SIRILinePublishedName() string {
	return s.linePublishedName
}

func (s *PartnerSettings) setSIRIPassageOrder(settings map[string]string) {
	passageOrder := settings[SIRI_PASSAGE_ORDER]
	s.passageOrder = passageOrder
}

func (s *PartnerSettings) SIRIPassageOrder() string {
	return s.passageOrder
}

func (s *PartnerSettings) setSIRIEnvelopeType(settings map[string]string) {
	envelopeType := settings[SIRI_ENVELOPE]
	if envelopeType == "" {
		envelopeType = "soap"
	}
	s.envelopeType = envelopeType
}

func (s *PartnerSettings) SIRIEnvelopeType() string {
	return s.envelopeType
}

func (s *PartnerSettings) setMaximumChechstatusRetry(settings map[string]string) {
	maxRetry, _ := strconv.Atoi(settings[PARTNER_MAX_RETRY])
	if maxRetry < 0 {
		maxRetry = 0
	}
	s.maximumCheckstatusRetry = maxRetry
}

func (s *PartnerSettings) MaximumChechstatusRetry() int {
	return s.maximumCheckstatusRetry
}

func (s *PartnerSettings) setSubscriptionMaximumResources(settings map[string]string) {
	maxResources, _ := strconv.Atoi(settings[SUBSCRIPTIONS_MAXIMUM_RESOURCES])
	s.subscriptionMaximumResources = maxResources
}

func (s *PartnerSettings) SubscriptionMaximumResources() int {
	return s.subscriptionMaximumResources
}

func (s *PartnerSettings) setCollectPriority(settings map[string]string) {
	collectPriority, _ := strconv.Atoi(settings[COLLECT_PRIORITY])
	s.collectPriority = collectPriority
}
func (s *PartnerSettings) CollectPriority() int {
	return s.collectPriority
}

func (s *PartnerSettings) setDefaultSRSName(settings map[string]string) {
	var srsName string
	if settings[COLLECT_DEFAULT_SRS_NAME] == "" {
		srsName = "EPSG:2154"
	} else {
		srsName = settings[COLLECT_DEFAULT_SRS_NAME]
	}
	s.defaultSRSName = srsName
}
func (s *PartnerSettings) DefaultSRSName() string {
	return s.defaultSRSName
}

func (s *PartnerSettings) setNoDestinationRefRewritingFrom(settings map[string]string) {
	values := trimedSlice(settings[BROADCAST_NO_DESTINATIONREF_REWRITING_FROM])
	s.noDestinationRefRewritingFrom = values
}
func (s *PartnerSettings) NoDestinationRefRewritingFrom() []string {
	return s.noDestinationRefRewritingFrom
}

func (s *PartnerSettings) setNoDataFrameRefRewritingFrom(settings map[string]string) {
	values := trimedSlice(settings[BROADCAST_NO_DATAFRAMEREF_REWRITING_FROM])
	s.noDataFrameRefRewritingFrom = values
}

func (s *PartnerSettings) NoDataFrameRefRewritingFrom() []string {
	return s.noDataFrameRefRewritingFrom
}

func (s *PartnerSettings) setRewriteJourneyPatternRef(settings map[string]string) {
	rewrite, _ := strconv.ParseBool(settings[BROADCAST_REWRITE_JOURNEY_PATTERN_REF])
	s.rewriteJourneyPatternRef = rewrite
}
func (s *PartnerSettings) RewriteJourneyPatternRef() bool {
	return s.rewriteJourneyPatternRef
}

func (s *PartnerSettings) setGzipGtfs(settings map[string]string) {
	gzipGtfs, _ := strconv.ParseBool(settings[BROADCAST_GZIP_GTFS])
	s.gzipGtfs = gzipGtfs
}
func (s *PartnerSettings) GzipGtfs() bool {
	return s.gzipGtfs
}

func (s *PartnerSettings) setGeneralMessageRequestVersion22(settings map[string]string) {
	version22, _ := strconv.ParseBool(settings[GENERAL_MESSAGE_REQUEST_2_2])
	s.generalMessageRequestVersion22 = version22
}
func (s *PartnerSettings) GeneralMessageRequestVersion22() bool {
	return s.generalMessageRequestVersion22
}

func (s *PartnerSettings) setPersistentCollect(settings map[string]string) {
	subscription, _ := strconv.ParseBool(settings[COLLECT_SUBSCRIPTIONS_PERSISTENT])
	collect, _ := strconv.ParseBool(settings[COLLECT_PERSISTENT])

	peristent := subscription || collect
	s.persistentCollect = peristent
}

func (s *PartnerSettings) PersistentCollect() bool {
	return s.persistentCollect
}

func (s *PartnerSettings) setPersistentBroadcastSubscriptions(settings map[string]string) {
	persistent, _ := strconv.ParseBool(settings[BROADCAST_SUBSCRIPTIONS_PERSISTENT])
	s.persistentBroadcastSubscriptions = persistent
}

func (s *PartnerSettings) PersistentBroadcastSubscriptions() bool {
	return s.persistentBroadcastSubscriptions
}

func (s *PartnerSettings) setCollectFilteredGeneralMessages(settings map[string]string) {
	collect, _ := strconv.ParseBool(settings[COLLECT_FILTER_GENERAL_MESSAGES])
	s.collectFilteredGeneralMessages = collect
}

func (s *PartnerSettings) CollectFilteredGeneralMessages() bool {
	return s.collectFilteredGeneralMessages
}

func (s *PartnerSettings) setIgnoreStopWithoutLine(settings map[string]string) {
	s.ignoreStopWithoutLine = settings[IGNORE_STOP_WITHOUT_LINE] != "false"
}

func (s *PartnerSettings) IgnoreStopWithoutLine() bool {
	return s.ignoreStopWithoutLine
}

func (s *PartnerSettings) setDiscoveryInterval(settings map[string]string) {
	interval, _ := time.ParseDuration(settings[DISCOVERY_INTERVAL])
	if interval == 0 {
		interval = 1 * time.Hour
	}

	s.discoveryInterval = -interval
}

func (s *PartnerSettings) DiscoveryInterval() time.Duration {
	return s.discoveryInterval
}

func (s *PartnerSettings) setCollectSettings(settings map[string]string) {
	s.collectSettings = &CollectSettings{
		UseDiscoveredSA:    settings[COLLECT_USE_DISCOVERED_SA] != "",
		UseDiscoveredLines: settings[COLLECT_USE_DISCOVERED_LINES] != "",
		includedSA:         toMap(settings[COLLECT_INCLUDE_STOP_AREAS]),
		excludedSA:         toMap(settings[COLLECT_EXCLUDE_STOP_AREAS]),
		includedLines:      toMap(settings[COLLECT_INCLUDE_LINES]),
		excludedLines:      toMap(settings[COLLECT_EXCLUDE_LINES]),
	}
}

func (s *PartnerSettings) CollectSettings() *CollectSettings {
	return s.collectSettings
}

func (s *PartnerSettings) setSiriCredentialHeader(settings map[string]string) {
	header := settings[SIRI_CREDENTIAL_HEADER]

	if header == "" {
		header = "X-SIRI-Requestor"
	}

	s.siriCredentialHeader = header
}

func (s *PartnerSettings) SiriCredentialHeader() string {
	return s.siriCredentialHeader
}

func (s *PartnerSettings) setHTTPClientOptions(settings map[string]string) {
	credential := remote.SiriCredentialHeader{
		CredentialHeader: s.SiriCredentialHeader(),
		Value:            s.RequestorRef(),
	}

	s.httpClientOptions = remote.HTTPClientOptions{
		SiriEnvelopeType: s.SiriEnvelopeType(),
		OAuth:            s.HTTPClientOAuth(),
		SiriCredential:   credential,
		Urls: remote.HTTPClientUrls{
			Url:              settings[REMOTE_URL],
			SubscriptionsUrl: settings[SUBSCRIPTIONS_REMOTE_URL],
			NotificationsUrl: settings[NOTIFICATIONS_REMOTE_URL],
		},
	}
}

func (s *PartnerSettings) HTTPClientOptions() remote.HTTPClientOptions {
	return s.httpClientOptions
}

func (s *PartnerSettings) setSiriEnvelopeType(settings map[string]string) {
	value := settings[SIRI_ENVELOPE]
	if value == "" {
		value = "soap"
	}

	s.siriEnvelopeType = value
}

func (s *PartnerSettings) SiriEnvelopeType() string {
	return s.siriEnvelopeType
}

func (s *PartnerSettings) setHttpClientOAuth(settings map[string]string) {
	clientId, clientIdFound := settings[OAUTH_CLIENT_ID]
	clientSecret, clientSecretFound := settings[OAUTH_CLIENT_SECRET]
	tokenURL, tokenURLFound := settings[OAUTH_TOKEN_URL]

	if clientIdFound && clientSecretFound && tokenURLFound {
		s.httpClientOAuth = &remote.HTTPClientOAuth{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			TokenURL:     tokenURL,
		}
	} else {
		s.httpClientOAuth = nil
	}
}

func (s *PartnerSettings) HTTPClientOAuth() *remote.HTTPClientOAuth {
	return s.httpClientOAuth
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

func (s *PartnerSettings) setIdentifierGenerators(settings map[string]string) {
	s.messageIdentifierGenerator = s.createIdentifierGenerator(settings, idgen.MESSAGE_IDENTIFIER)
	s.responseMessageIdentifierGenerator = s.createIdentifierGenerator(settings, idgen.RESPONSE_MESSAGE_IDENTIFIER)
	s.dataFrameIdentifierGenerator = s.createIdentifierGenerator(settings, idgen.DATA_FRAME_IDENTIFIER)
	s.referenceIdentifierGenerator = s.createIdentifierGenerator(settings, idgen.REFERENCE_IDENTIFIER)
	s.referenceStopAreaIdentifierGenerator = s.createIdentifierGenerator(settings, idgen.REFERENCE_STOP_AREA_IDENTIFIER)
	s.subscriptionIdentifierGenerator = s.createIdentifierGenerator(settings, idgen.SUBSCRIPTION_IDENTIFIER)
}

func (s *PartnerSettings) MessageIdentifierGenerator() *idgen.IdentifierGenerator {
	return s.messageIdentifierGenerator
}

func (s *PartnerSettings) ResponseMessageIdentifierGenerator() *idgen.IdentifierGenerator {
	return s.responseMessageIdentifierGenerator
}

func (s *PartnerSettings) DataFrameIdentifierGenerator() *idgen.IdentifierGenerator {
	return s.dataFrameIdentifierGenerator
}

func (s *PartnerSettings) ReferenceIdentifierGenerator() *idgen.IdentifierGenerator {
	return s.referenceIdentifierGenerator
}

func (s *PartnerSettings) ReferenceStopAreaIdentifierGenerator() *idgen.IdentifierGenerator {
	return s.referenceStopAreaIdentifierGenerator
}

func (s *PartnerSettings) SubscriptionIdentifierGenerator() *idgen.IdentifierGenerator {
	return s.subscriptionIdentifierGenerator
}

func (s *PartnerSettings) createIdentifierGenerator(settings map[string]string, generatorName string) *idgen.IdentifierGenerator {
	format := settings[fmt.Sprintf("generators.%v", generatorName)]

	if format == "" {
		format = idgen.DefaultIdentifierGenerator(generatorName)
	}

	return idgen.NewIdentifierGenerator(format, s.ug())
}

func (s *PartnerSettings) NewMessageIdentifier() string {
	return s.messageIdentifierGenerator.NewMessageIdentifier()
}

func (s *PartnerSettings) NewResponseMessageIdentifier() string {
	return s.responseMessageIdentifierGenerator.NewMessageIdentifier()
}

// func (s *PartnerSettings) idGeneratorFormat(generatorName string) (formatString string) {
// 	formatString = s.s[fmt.Sprintf("generators.%v", generatorName)]

// 	if formatString == "" {
// 		formatString = idgen.DefaultIdentifierGenerator(generatorName)
// 	}
// 	return
// }

// // Warning, this method isn't threadsafe. Mutex must be handled before and after calling
// func (s *PartnerSettings) refreshGenerators() {
// 	for k := range s.g {
// 		delete(s.g, k)
// 	}
// }

// // Warning, this method isn't threadsafe. Mutex must be handled before and after calling
// func (s *PartnerSettings) reloadSettings() {
// 	s.SetCollectSettings()
// 	s.refreshGenerators()
// }
