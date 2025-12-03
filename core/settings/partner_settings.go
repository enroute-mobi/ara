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

	REMOTE_CREDENTIAL                 = "remote_credential"
	REMOTE_CODE_SPACE                 = "remote_code_space"
	VEHICLE_REMOTE_CODE_SPACE         = "vehicle_remote_code_space"
	VEHICLE_JOURNEY_REMOTE_CODE_SPACE = "vehicle_journey_remote_code_space"
	REMOTE_URL                        = "remote_url"
	NOTIFICATIONS_REMOTE_URL          = "notifications.remote_url"
	SUBSCRIPTIONS_REMOTE_URL          = "subscriptions.remote_url"

	COLLECT_DEFAULT_SRS_NAME                                              = "collect.default_srs_name"
	COLLECT_EXCLUDE_FACILITIES                                            = "collect.exclude_facilities"
	COLLECT_EXCLUDE_LINES                                                 = "collect.exclude_lines"
	COLLECT_EXCLUDE_LINE_GROUPS                                           = "collect.exclude_line_groups"
	COLLECT_EXCLUDE_STOP_AREAS                                            = "collect.exclude_stop_areas"
	COLLECT_EXCLUDE_STOP_AREA_GROUPS                                      = "collect.exclude_stop_area_groups"
	COLLECT_FILTER_GENERAL_MESSAGES                                       = "collect.filter_general_messages" // Kept for retro compatibility
	COLLECT_FILTER_SITUATIONS                                             = "collect.filter_situations"
	COLLECT_GTFS_TTL                                                      = "collect.gtfs.ttl"
	COLLECT_INCLUDE_FACILITIES                                            = "collect.include_facilities"
	COLLECT_INCLUDE_LINES                                                 = "collect.include_lines"
	COLLECT_INCLUDE_LINE_GROUPS                                           = "collect.include_line_groups"
	COLLECT_INCLUDE_STOP_AREAS                                            = "collect.include_stop_areas"
	COLLECT_INCLUDE_STOP_AREA_GROUPS                                      = "collect.include_stop_area_groups"
	COLLECT_PERSISTENT                                                    = "collect.persistent"
	COLLECT_PRIORITY                                                      = "collect.priority"
	COLLECT_SITUATIONS_INTERNAL_TAGS                                      = "collect.situations.internal_tags"
	COLLECT_SUBSCRIPTIONS_PERSISTENT                                      = "collect.subscriptions.persistent"
	COLLECT_USE_DISCOVERED_LINES                                          = "collect.use_discovered_lines"
	COLLECT_USE_DISCOVERED_SA                                             = "collect.use_discovered_stop_areas"
	COLLECT_SIRI_STOP_MONITORING_MAXIMUM_SUBSCRIPTION_PER_REQUEST         = "collect.siri.stop_monitoring.maximum_subscriptions_per_request"
	COLLECT_DEFAULT_SIRI_STOP_MONITORING_MAXIMUM_SUBSCRIPTION_PER_REQUEST = 1000

	DISCOVERY_INTERVAL = "discovery_interval"

	BROADCAST_GTFS_CACHE_TIMEOUT                          = "broadcast.gtfs.cache_timeout"
	BROADCAST_GZIP_GTFS                                   = "broadcast.gzip_gtfs"
	BROADCAST_NO_DATAFRAMEREF_REWRITING_FROM              = "broadcast.no_dataframeref_rewriting_from"
	BROADCAST_NO_DESTINATIONREF_REWRITING_FROM            = "broadcast.no_destinationref_rewriting_from"
	BROADCAST_PREFER_REFERENT_STOP_AREAS                  = "broadcast.prefer_referent_stop_areas"
	BROADCAST_RECORDED_CALLS_DURATION                     = "broadcast.recorded_calls.duration"
	BROADCAST_REWRITE_JOURNEY_PATTERN_REF                 = "broadcast.rewrite_journey_pattern_ref"
	BROADCAST_SIRI_IGNORE_TERMINATE_SUBSCRIPTION_REQUESTS = "broadcast.siri.ignore_terminate_subscription_requests"
	BROADCAST_SIRI_SM_MULTIPLE_SUBSCRIPTIONS              = "broadcast.siri.stop_monitoring.multiple_subscriptions"
	BROADCAST_SIRI_SM_MAXIMUM_RESOURCES_PER_DELIVERY      = "broadcast.siri.stop_monitoring.maximum_resources_per_delivery"
	BROADCAST_DEFAULT_SM_MAXIMUM_RESOURCES_PER_DELIVERY   = 50
	BROADCAST_SITUATIONS_INTERNAL_TAGS                    = "broadcast.situations.internal_tags"
	BROADCAST_SITUATIONS_TTL                              = "broadcast.situations.time_to_live"
	BROADCAST_SUBSCRIPTIONS_PERSISTENT                    = "broadcast.subscriptions.persistent"

	IGNORE_STOP_WITHOUT_LINE        = "ignore_stop_without_line"
	GENERAL_MESSAGE_REQUEST_2_2     = "generalMessageRequest.version2.2"
	SUBSCRIPTIONS_MAXIMUM_RESOURCES = "subscriptions.maximum_resources"

	CACHE_TIMEOUT = "cache_timeout"

	OAUTH_CLIENT_ID     = "remote_authentication.oauth.client_id"
	OAUTH_CLIENT_SECRET = "remote_authentication.oauth.client_secret"
	OAUTH_TOKEN_URL     = "remote_authentication.oauth.token_url"
	OAUTH_SCOPES        = "remote_authentication.oauth.scopes"

	SIRI_ENVELOPE                            = "siri.envelope"
	SIRI_LINE_PUBLISHED_NAME                 = "siri.line.published_name"
	SIRI_DIRECTION_TYPE                      = "siri.direction_type"
	SIRI_PASSAGE_ORDER                       = "siri.passage_order"
	SIRI_CREDENTIAL_HEADER                   = "siri.credential.header"
	SIRI_SOAP_EMPTY_RESPONSE_ON_NOTIFICATION = "siri.soap.empty_response_on_notification"
	DEFAULT_GTFS_TTL                         = 30 * time.Second
	DEFAULT_SITUATIONS_TTL                   = 1 * time.Hour

	HTTP_CUSTOM_HEADERS = "http.custom_headers"

	SORT_PAYLOAD_FOR_TEST = "sort_payload_for_test"

	GRAPHQL_MUTABLE_ATTRIBUTES = "graphql.mutable_attributes"
)

type PartnerSettings struct {
	idgen.IdentifierGenerator

	ug func() uuid.UUIDGenerator

	collectSettings *CollectSettings

	credentials                                             string
	rateLimit                                               float64
	gtfsTTL                                                 time.Duration
	gtfsCacheTimeout                                        time.Duration
	httpCustomHeaders                                       []string
	siriCredentialHeader                                    string
	siriEnvelopeType                                        string
	siriSoapEmptyResponseOnNotification                     bool
	httpClientOAuth                                         *remote.HTTPClientOAuth
	recordedCallsDuration                                   time.Duration
	producerRef                                             string
	address                                                 string
	linePublishedName                                       string
	passageOrder                                            string
	envelopeType                                            string
	collectPriority                                         int
	collectSituationsInternalTags                           []string
	broadcastSituationsInternalTags                         []string
	broadcastSituationsTTL                                  time.Duration
	defaultSRSName                                          string
	noDestinationRefRewritingFrom                           []string
	noDataFrameRefRewritingFrom                             []string
	rewriteJourneyPatternRef                                bool
	preferReferentStopArea                                  bool
	gzipGtfs                                                bool
	generalMessageRequestVersion22                          bool
	collectFilteredSituations                               bool
	collectSiriStopMonitoringMaximumSubscriptionsPerRequest int
	ignoreStopWithoutLine                                   bool
	smMultipleDeliveriesPerNotify                           bool
	smMaxStopVisitPerDelivery                               int
	discoveryInterval                                       time.Duration
	cacheTimeouts                                           sync.Map
	siriDirectionTypeInbound                                string
	siriDirectionTypeOutbound                               string

	maximumCheckstatusRetry          int
	subscriptionMaximumResources     int
	persistentCollect                bool
	persistentBroadcastSubscriptions bool

	remoteCodeSpaces sync.Map
	remoteCodeSpace  string

	httpClientOptions remote.HTTPClientOptions

	sortPayloadForTest                  bool
	ignoreTerminateSubscriptionsRequest bool

	vehicleRemoteCodeSpaces            []string
	vehicleRemoteCodeSpacesByConnector sync.Map

	vehicleJourneyRemoteCodeSpaces            []string
	vehicleJourneyRemoteCodeSpacesByConnector sync.Map

	graphQLMutableAttributes sync.Map

	// ! Never use these values outside SettingsDefinition()
	originalSettings map[string]string
}

func NewEmptyPartnerSettings(ug func() uuid.UUIDGenerator) *PartnerSettings {
	partnerSettings := &PartnerSettings{
		ug: ug,
	}
	partnerSettings.parseSettings(map[string]string{}, nil)
	return partnerSettings
}

func NewPartnerSettings(generator func() uuid.UUIDGenerator, settings map[string]string, resolvers ...func(string, string) ([]string, bool)) *PartnerSettings {
	partnerSettings := &PartnerSettings{
		ug: generator,
	}

	partnerSettings.parseSettings(settings, resolvers)

	return partnerSettings
}

func (s *PartnerSettings) parseSettings(settings map[string]string, resolvers []func(string, string) ([]string, bool)) {
	s.setRemoteCodeSpaces(settings)
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
	s.setCollectSettings(settings, resolvers)
	s.setHttpCustomHeaders(settings)
	s.setSiriCredentialHeader(settings)
	s.setSiriEnvelopeType(settings)
	s.setSiriSoapEmptyResponseOnNotification(settings)
	s.setSIRIPassageOrder(settings)
	s.setMaximumChechstatusRetry(settings)
	s.setSubscriptionMaximumResources(settings)
	s.setCollectPriority(settings)
	s.setCollectSituationsInternalTags(settings)
	s.setSiriStopMonitoringMaximumSubscriptionsPerRequest(settings)
	s.setBroadcastSituationsInternalTags(settings)
	s.setSituationsTTL(settings)
	s.setDefaultSRSName(settings)
	s.setNoDestinationRefRewritingFrom(settings)
	s.setNoDataFrameRefRewritingFrom(settings)
	s.setGzipGtfs(settings)
	s.setGeneralMessageRequestVersion22(settings)
	s.setPersistentCollect(settings)
	s.setPersistentBroadcastSubscriptions(settings)
	s.setPreferReferentStopArea(settings)
	s.setRewriteJourneyPatternRef(settings)
	s.setCollectFilteredSituations(settings)
	s.setIgnoreStopWithoutLine(settings)
	s.setDiscoveryInterval(settings)
	s.setCacheTimeouts(settings)
	s.setSortPayloadForTest(settings)
	s.setSmMultipleDeliveriesPerNotify(settings)
	s.setMaxStopVisitPerDelivery(settings)

	s.setVehicleRemoteCodeSpaceWithFallback(settings)
	s.setVehicleJourneyRemoteCodeSpaceWithFallback(settings)

	s.setIgnoreTerminateSubscriptionsRequest(settings)
	s.setIdentifierGenerator(settings)

	s.setHttpClientOAuth(settings)

	// depends on other settings
	s.setHTTPClientOptions(settings)

	s.setGraphQLMutableAttributes(settings)

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

func (s *PartnerSettings) setRemoteCodeSpaces(settings map[string]string) {
	r, _ := regexp.Compile(`(.+)\.remote_code_space`)

	// xxxx.remote_code_space = dummy -> xxxx = dummy
	for key, value := range settings {
		if len(value) == 0 {
			continue
		}

		matches := r.FindStringSubmatch(key)
		if len(matches) == 0 {
			continue
		}

		connectorName := matches[1]
		s.remoteCodeSpaces.Store(connectorName, value)
	}

	s.remoteCodeSpace = settings[REMOTE_CODE_SPACE]
}

func (s *PartnerSettings) RemoteCodeSpace(optionalConnectorName ...string) string {
	if len(optionalConnectorName) == 1 {
		connectorName := optionalConnectorName[0]

		value, ok := s.remoteCodeSpaces.Load(connectorName)
		if ok {
			return value.(string)
		}
	}

	return s.remoteCodeSpace
}

func (s *PartnerSettings) setVehicleJourneyRemoteCodeSpaceWithFallback(settings map[string]string) {
	// xxxx.vehicle_journey_remote_code_space
	r, _ := regexp.Compile(fmt.Sprintf("(.+)\\.%s", VEHICLE_JOURNEY_REMOTE_CODE_SPACE))

	s.vehicleJourneyRemoteCodeSpaces = trimedSlice(settings[VEHICLE_JOURNEY_REMOTE_CODE_SPACE])

	// xxxx.vehicle_journey_remote_code_space = dummy -> xxxx = dummy
	for key, value := range settings {
		if len(value) == 0 {
			continue
		}

		matches := r.FindStringSubmatch(key)
		if len(matches) == 0 {
			continue
		}

		var connectorRemoteCodeSpaces []string
		connectorRemoteCodeSpaces = append(connectorRemoteCodeSpaces, trimedSlice(value)...)
		connectorRemoteCodeSpaces = append(connectorRemoteCodeSpaces, s.vehicleJourneyRemoteCodeSpaces...)

		// "a,b,c" -> [a,b,c]

		connectorName := matches[1]
		s.vehicleJourneyRemoteCodeSpacesByConnector.Store(connectorName, connectorRemoteCodeSpaces)
	}
}

func (s *PartnerSettings) VehicleJourneyRemoteCodeSpaceWithFallback(connectorName string) []string {
	value, ok := s.vehicleJourneyRemoteCodeSpacesByConnector.Load(connectorName)
	if ok {
		return value.([]string)
	}

	if len(s.vehicleJourneyRemoteCodeSpaces) > 0 {
		return s.vehicleJourneyRemoteCodeSpaces
	}

	return []string{s.RemoteCodeSpace(connectorName)}
}

func (s *PartnerSettings) setVehicleRemoteCodeSpaceWithFallback(settings map[string]string) {
	// xxxx.vehicle_journey_remote_code_space
	r, _ := regexp.Compile(fmt.Sprintf("(.+)\\.%s", VEHICLE_REMOTE_CODE_SPACE))

	s.vehicleRemoteCodeSpaces = trimedSlice(settings[VEHICLE_REMOTE_CODE_SPACE])

	// xxxx.vehicle_journey_remote_code_space = dummy -> xxxx = dummy
	for key, value := range settings {
		if len(value) == 0 {
			continue
		}

		matches := r.FindStringSubmatch(key)
		if len(matches) == 0 {
			continue
		}

		var connectorRemoteCodeSpaces []string
		connectorRemoteCodeSpaces = append(connectorRemoteCodeSpaces, trimedSlice(value)...)
		connectorRemoteCodeSpaces = append(connectorRemoteCodeSpaces, s.vehicleRemoteCodeSpaces...)

		// "a,b,c" -> [a,b,c]

		connectorName := matches[1]
		s.vehicleRemoteCodeSpacesByConnector.Store(connectorName, connectorRemoteCodeSpaces)
	}
}

func (s *PartnerSettings) VehicleRemoteCodeSpaceWithFallback(connectorName string) []string {
	value, ok := s.vehicleRemoteCodeSpacesByConnector.Load(connectorName)
	if ok {
		return value.([]string)
	}

	if len(s.vehicleRemoteCodeSpaces) > 0 {
		return s.vehicleRemoteCodeSpaces
	}

	return []string{s.RemoteCodeSpace(connectorName)}
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
	r, _ := regexp.Compile(`(.+)\.cache_timeout`)

	// xxxx.cache_timeout = dummy -> xxxx = dummy
	for key, value := range settings {
		if len(value) == 0 {
			continue
		}

		matches := r.FindStringSubmatch(key)
		if len(matches) == 0 {
			continue
		}

		connectorName := matches[1]
		s.cacheTimeouts.Store(connectorName, value)
	}
}

func (s *PartnerSettings) CacheTimeout(connectorName string) (t time.Duration) {
	value, ok := s.cacheTimeouts.Load(connectorName)
	if ok {
		return value.(time.Duration)
	}

	return 0
}

func (s *PartnerSettings) setSortPayloadForTest(settings map[string]string) {
	sortPayload, _ := strconv.ParseBool(settings[SORT_PAYLOAD_FOR_TEST])
	s.sortPayloadForTest = sortPayload
}
func (s *PartnerSettings) SortPaylodForTest() bool {
	return s.sortPayloadForTest
}

func (s *PartnerSettings) setSmMultipleDeliveriesPerNotify(settings map[string]string) {
	m, _ := strconv.ParseBool(settings[BROADCAST_SIRI_SM_MULTIPLE_SUBSCRIPTIONS])
	s.smMultipleDeliveriesPerNotify = m
}

func (s *PartnerSettings) SmMultipleDeliveriesPerNotify() bool {
	return s.smMultipleDeliveriesPerNotify
}

func (s *PartnerSettings) setSiriStopMonitoringMaximumSubscriptionsPerRequest(settings map[string]string) {
	max, _ := strconv.Atoi(settings[COLLECT_SIRI_STOP_MONITORING_MAXIMUM_SUBSCRIPTION_PER_REQUEST])
	if max > COLLECT_DEFAULT_SIRI_STOP_MONITORING_MAXIMUM_SUBSCRIPTION_PER_REQUEST || max == 0 {
		max = COLLECT_DEFAULT_SIRI_STOP_MONITORING_MAXIMUM_SUBSCRIPTION_PER_REQUEST
	}
	s.collectSiriStopMonitoringMaximumSubscriptionsPerRequest = max
}

func (s *PartnerSettings) StopMonitoringMaxSubscriptionPerRequest() int {
	return s.collectSiriStopMonitoringMaximumSubscriptionsPerRequest
}

func (s *PartnerSettings) setMaxStopVisitPerDelivery(settings map[string]string) {
	max, _ := strconv.Atoi(settings[BROADCAST_SIRI_SM_MAXIMUM_RESOURCES_PER_DELIVERY])
	if max > BROADCAST_DEFAULT_SM_MAXIMUM_RESOURCES_PER_DELIVERY {
		max = BROADCAST_DEFAULT_SM_MAXIMUM_RESOURCES_PER_DELIVERY
	}
	s.smMaxStopVisitPerDelivery = max
}

func (s *PartnerSettings) MaxStopVisitPerDelivery() int {
	return s.smMaxStopVisitPerDelivery
}

func (s *PartnerSettings) setIgnoreTerminateSubscriptionsRequest(settings map[string]string) {
	ignore, _ := strconv.ParseBool(settings[BROADCAST_SIRI_IGNORE_TERMINATE_SUBSCRIPTION_REQUESTS])
	s.ignoreTerminateSubscriptionsRequest = ignore
}

func (s *PartnerSettings) IgnoreTerminateSubscriptionsRequest() bool {
	return s.ignoreTerminateSubscriptionsRequest
}

func (s *PartnerSettings) setProducerRef(settings map[string]string) {
	producerRef := settings[REMOTE_CREDENTIAL]
	if producerRef == "" {
		producerRef = "Ara"
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
		s.siriDirectionTypeInbound = ""
		s.siriDirectionTypeOutbound = ""
		return
	}

	s.siriDirectionTypeInbound = directions[0]
	s.siriDirectionTypeOutbound = directions[1]
}

func (s *PartnerSettings) SIRIDirectionType() (string, string, bool) {
	valid := len(s.siriDirectionTypeOutbound) > 0 && len(s.siriDirectionTypeInbound) > 0
	return s.siriDirectionTypeInbound, s.siriDirectionTypeOutbound, !valid
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

func (s *PartnerSettings) setSiriSoapEmptyResponseOnNotification(settings map[string]string) {
	s.siriSoapEmptyResponseOnNotification, _ = strconv.ParseBool(settings[SIRI_SOAP_EMPTY_RESPONSE_ON_NOTIFICATION])
}

func (s *PartnerSettings) SiriSoapEmptyResponseOnNotification() bool {
	return s.siriSoapEmptyResponseOnNotification
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

func (s *PartnerSettings) setCollectSituationsInternalTags(settings map[string]string) {
	values := trimedSlice(settings[COLLECT_SITUATIONS_INTERNAL_TAGS])
	s.collectSituationsInternalTags = values
}

func (s *PartnerSettings) CollectSituationsInternalTags() []string {
	return s.collectSituationsInternalTags
}

func (s *PartnerSettings) setBroadcastSituationsInternalTags(settings map[string]string) {
	values := trimedSlice(settings[BROADCAST_SITUATIONS_INTERNAL_TAGS])
	s.broadcastSituationsInternalTags = values
}

func (s *PartnerSettings) BroadcastSituationsInternalTags() []string {
	return s.broadcastSituationsInternalTags
}

func (s *PartnerSettings) setSituationsTTL(settings map[string]string) {
	duration := DEFAULT_SITUATIONS_TTL

	if settings[BROADCAST_SITUATIONS_TTL] != "" {
		duration, _ = time.ParseDuration(settings[BROADCAST_SITUATIONS_TTL])
	}

	s.broadcastSituationsTTL = duration
}

func (s *PartnerSettings) SituationsTTL() (t time.Duration) {
	return s.broadcastSituationsTTL
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

func (s *PartnerSettings) setPreferReferentStopArea(settings map[string]string) {
	preferReferent, _ := strconv.ParseBool(settings[BROADCAST_PREFER_REFERENT_STOP_AREAS])
	s.preferReferentStopArea = preferReferent
}

func (s *PartnerSettings) PreferReferentStopArea() bool {
	return s.preferReferentStopArea
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

func (s *PartnerSettings) setCollectFilteredSituations(settings map[string]string) {
	cgm, _ := strconv.ParseBool(settings[COLLECT_FILTER_GENERAL_MESSAGES])
	cs, _ := strconv.ParseBool(settings[COLLECT_FILTER_SITUATIONS])
	s.collectFilteredSituations = cgm || cs
}

func (s *PartnerSettings) CollectFilteredSituations() bool {
	return s.collectFilteredSituations
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

func (s *PartnerSettings) setCollectSettings(settings map[string]string, resolvers []func(string, string) ([]string, bool)) {
	s.collectSettings = &CollectSettings{
		UseDiscoveredSA:    settings[COLLECT_USE_DISCOVERED_SA] != "",
		UseDiscoveredLines: settings[COLLECT_USE_DISCOVERED_LINES] != "",
		includedSA:         toMap(settings[COLLECT_INCLUDE_STOP_AREAS]),
		excludedSA:         toMap(settings[COLLECT_EXCLUDE_STOP_AREAS]),
		includedLines:      toMap(settings[COLLECT_INCLUDE_LINES]),
		excludedLines:      toMap(settings[COLLECT_EXCLUDE_LINES]),
		includedFacilities: toMap(settings[COLLECT_INCLUDE_FACILITIES]),
		excludedFacilities: toMap(settings[COLLECT_EXCLUDE_FACILITIES]),
	}

	remoteCodeSpace := s.RemoteCodeSpace()
	if len(resolvers) != 0 {
		stopAreaResolver := resolvers[0]

		for shortName := range toMap(settings[COLLECT_INCLUDE_STOP_AREA_GROUPS]) {
			stopAreaValues, ok := stopAreaResolver(shortName, remoteCodeSpace)
			if !ok {
				continue
			}

			for _, stopAreaValue := range stopAreaValues {
				s.collectSettings.includedSA[stopAreaValue] = struct{}{}
			}
		}

		for shortName := range toMap(settings[COLLECT_EXCLUDE_STOP_AREA_GROUPS]) {
			stopAreaValues, ok := stopAreaResolver(shortName, remoteCodeSpace)
			if !ok {
				continue
			}

			for _, stopAreaValue := range stopAreaValues {
				s.collectSettings.excludedSA[stopAreaValue] = struct{}{}
			}
		}
	}

	if len(resolvers) > 1 {
		lineResolver := resolvers[1]

		for shortName := range toMap(settings[COLLECT_INCLUDE_LINE_GROUPS]) {
			lineValues, ok := lineResolver(shortName, remoteCodeSpace)
			if !ok {
				continue
			}

			for _, lineValue := range lineValues {
				s.collectSettings.includedLines[lineValue] = struct{}{}
			}
		}

		for shortName := range toMap(settings[COLLECT_EXCLUDE_LINE_GROUPS]) {
			lineValues, ok := lineResolver(shortName, remoteCodeSpace)
			if !ok {
				continue
			}

			for _, lineValue := range lineValues {
				s.collectSettings.excludedLines[lineValue] = struct{}{}
			}
		}
	}
}

func (s *PartnerSettings) CollectSettings() *CollectSettings {
	return s.collectSettings
}

func (s *PartnerSettings) setHttpCustomHeaders(settings map[string]string) {
	headers := trimedSlice(settings[HTTP_CUSTOM_HEADERS])

	s.httpCustomHeaders = headers
}

func (s *PartnerSettings) HttpCustomHeaders() []string {
	return s.httpCustomHeaders
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
		CustomHeaders:    s.HttpCustomHeaders(),
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

func (s *PartnerSettings) setGraphQLMutableAttributes(settings map[string]string) {
	for key := range toMap(settings[GRAPHQL_MUTABLE_ATTRIBUTES]) {
		s.graphQLMutableAttributes.Store(key, struct{}{})
	}
}

func (s *PartnerSettings) IsMutable(attribute string) (ok bool) {
	_, ok = s.graphQLMutableAttributes.Load(attribute)
	return
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
	scopes, scopesFound := settings[OAUTH_SCOPES]

	if clientIdFound && clientSecretFound && tokenURLFound {
		s.httpClientOAuth = &remote.HTTPClientOAuth{
			ClientID:     clientId,
			ClientSecret: clientSecret,
			TokenURL:     tokenURL,
		}
		if scopesFound {
			s.httpClientOAuth.Scopes = strings.Split(scopes, ",")
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

func (s *PartnerSettings) setIdentifierGenerator(settings map[string]string) {
	s.IdentifierGenerator = idgen.NewIdentifierGenerator(settings, s.ug())
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
