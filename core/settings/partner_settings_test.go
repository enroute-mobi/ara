package settings

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/cache"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_Credentials_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	assert.Equal("", partnerSettings.Credentials())
}

func Test_Credentials_Single_Credential(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		LOCAL_CREDENTIAL: "credential",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal("credential,", partnerSettings.Credentials())
}

func Test_Credentials_Multiple_Credentials(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		LOCAL_CREDENTIAL:  "credential",
		LOCAL_CREDENTIALS: "another_credential",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal("credential,another_credential", partnerSettings.Credentials())
}

func Test_Credentials_Multiple_Credentials_Without_Local_Credential(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		LOCAL_CREDENTIALS: "another_credential",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(",another_credential", partnerSettings.Credentials())
}

func Test_RateLimit(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		RATE_LIMIT_PER_IP: "100",
	}

	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal(float64(100), partnerSettings.RateLimit())
}

func Test_RemoteCodeSpace_Without_Connector(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		REMOTE_CODE_SPACE: "external",
	}

	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal("external", partnerSettings.RemoteCodeSpace())
}

func Test_RemoteCodeSpace_With_Connector(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"connector_name.remote_code_space":  "external",
		"connector_name1.remote_code_space": "another",
	}

	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal("external", partnerSettings.RemoteCodeSpace("connector_name"))
	assert.Equal("another", partnerSettings.RemoteCodeSpace("connector_name1"))
}

func Test_VehicleJourneyRemoteCodeSpaceWithFallback(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		VEHICLE_JOURNEY_REMOTE_CODE_SPACE: "external",
	}

	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal([]string{"external"}, partnerSettings.VehicleJourneyRemoteCodeSpaceWithFallback((VEHICLE_JOURNEY_REMOTE_CODE_SPACE)))
}

func Test_VehicleJourneyRemoteCodeSpaceWithFallback_With_Multiple_Connectors(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"connector_name.vehicle_journey_remote_code_space":  "external",
		"connector_name1.vehicle_journey_remote_code_space": "another",
	}

	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal([]string{"external"}, partnerSettings.VehicleJourneyRemoteCodeSpaceWithFallback("connector_name"))
	assert.Equal([]string{"another"}, partnerSettings.VehicleJourneyRemoteCodeSpaceWithFallback("connector_name1"))
}

func Test_VehicleRemoteCodeSpaceWithFallback(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		VEHICLE_REMOTE_CODE_SPACE: "external",
	}

	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal([]string{"external"}, partnerSettings.VehicleRemoteCodeSpaceWithFallback((VEHICLE_REMOTE_CODE_SPACE)))
}

func Test_VehicleRemoteCodeSpaceWithFallback_With_Multiple_Connectors(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"connector_name.vehicle_remote_code_space":  "external",
		"connector_name1.vehicle_remote_code_space": "another",
	}

	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal([]string{"external"}, partnerSettings.VehicleRemoteCodeSpaceWithFallback("connector_name"))
	assert.Equal([]string{"another"}, partnerSettings.VehicleRemoteCodeSpaceWithFallback("connector_name1"))
}

func Test_GtfsTTLBelow30Seconds(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		COLLECT_GTFS_TTL: "29s",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(time.Duration(30000000000), partnerSettings.GtfsTTL(), "Should set GtfsTTL at default 30 minutes")
}

func Test_GtfsTTL_Above_30_Seconds(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		COLLECT_GTFS_TTL: "31s",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(time.Duration(31000000000), partnerSettings.GtfsTTL())
}

func Test_RecordedCallsDurations(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_RECORDED_CALLS_DURATION: "1h",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(time.Duration(3600000000000), partnerSettings.RecordedCallsDuration())
}

func Test_GtfsCacheTimeout_Below_Min_LifeSpan(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_GTFS_CACHE_TIMEOUT: "1s",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(cache.DEFAULT_CACHE_LIFESPAN, partnerSettings.GtfsCacheTimeout(), "Should set GtfsCacheTimeout to DEFAULT_CACHE_LIFESPAN")
}

func Test_GtfsCacheTimeout_Above_MinLifeSpan(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_GTFS_CACHE_TIMEOUT: "30s",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(time.Duration(30000000000), partnerSettings.GtfsCacheTimeout())
}

func Test_SortPaylodForTest(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		SORT_PAYLOAD_FOR_TEST: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.SortPaylodForTest())
}

func Test_IgnoreTerminateSubscriptionsRequest(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_SIRI_IGNORE_TERMINATE_SUBSCRIPTION_REQUESTS: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.IgnoreTerminateSubscriptionsRequest())
}

func Test_ProducerRef_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	assert.Equal("Ara", partnerSettings.ProducerRef())
}

func Test_ProducerRef_And_RequestoRef(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		REMOTE_CREDENTIAL: "credential",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal("credential", partnerSettings.ProducerRef())
	assert.Equal("credential", partnerSettings.RequestorRef())
}

func Test_Address(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		LOCAL_URL: "https://url.com",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal("https://url.com", partnerSettings.Address())
}

func Test_SIRIDirectionType_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	inboundValue, outboundValue, err := partnerSettings.SIRIDirectionType()
	assert.Equal("", inboundValue)
	assert.Equal("", outboundValue)
	assert.True(err)
}

func Test_SIRIDirectionType(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		SIRI_DIRECTION_TYPE: "A,B",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	inboundValue, outboundValue, err := partnerSettings.SIRIDirectionType()
	assert.Equal("A", inboundValue)
	assert.Equal("B", outboundValue)
	assert.False(err)
}

func Test_SIRILinePublishedName(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		SIRI_LINE_PUBLISHED_NAME: "name",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal("name", partnerSettings.SIRILinePublishedName())
}

func Test_SIRIEnvelopeType_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	assert.Equal("soap", partnerSettings.SIRIEnvelopeType())
}

func Test_SIRIEnvelopeType(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		SIRI_ENVELOPE: "envelope",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal("envelope", partnerSettings.SIRIEnvelopeType())
}

func Test_MaximumChechstatusRetry_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	assert.Equal(0, partnerSettings.MaximumChechstatusRetry())
}

func Test_MaximumChechstatusRetry_Below_Zero(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		PARTNER_MAX_RETRY: "-20",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(0, partnerSettings.MaximumChechstatusRetry())
}

func Test_MaximumChechstatusRetry(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		PARTNER_MAX_RETRY: "20",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(20, partnerSettings.MaximumChechstatusRetry())
}

func Test_SubscriptionMaximumResources(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		SUBSCRIPTIONS_MAXIMUM_RESOURCES: "5",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(5, partnerSettings.SubscriptionMaximumResources())
}

func Test_CollectPriority(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		COLLECT_PRIORITY: "5",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(5, partnerSettings.CollectPriority())
}

func Test_DefaultSRSName_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	assert.Equal("EPSG:2154", partnerSettings.DefaultSRSName())
}

func Test_DefaultSRSName(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		COLLECT_DEFAULT_SRS_NAME: "dummy",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal("dummy", partnerSettings.DefaultSRSName())
}

func Test_NoDestinationRefRewritingFrom(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_NO_DESTINATIONREF_REWRITING_FROM: "dummy",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal([]string{"dummy"}, partnerSettings.NoDestinationRefRewritingFrom())
}

func Test_NoDataFrameRefRewritingFrom(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_NO_DATAFRAMEREF_REWRITING_FROM: "dummy",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal([]string{"dummy"}, partnerSettings.NoDataFrameRefRewritingFrom())
}

func Test_RewriteJourneyPatternRef(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_REWRITE_JOURNEY_PATTERN_REF: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.RewriteJourneyPatternRef())
}

func Test_GzipGtfs(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_GZIP_GTFS: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.GzipGtfs())
}

func Test_GeneralMessageRequestVersion22(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		GENERAL_MESSAGE_REQUEST_2_2: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.GeneralMessageRequestVersion22())
}

func Test_PersistentCollect_With_Collect_Subscriptions_Persistent(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		COLLECT_SUBSCRIPTIONS_PERSISTENT: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.PersistentCollect())
}

func Test_PersistentCollect_With_Collect_Persistent(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		COLLECT_PERSISTENT: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.PersistentCollect())
}

func Test_PersistentBroadcastSubscriptions(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		BROADCAST_SUBSCRIPTIONS_PERSISTENT: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.PersistentBroadcastSubscriptions())
}

func Test_CollectFilteredGeneralMessages(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		COLLECT_FILTER_GENERAL_MESSAGES: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.CollectFilteredSituations())

	settings = map[string]string{
		COLLECT_FILTER_SITUATIONS: "true",
	}
	partnerSettings = NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.CollectFilteredSituations())
}

func Test_IgnoreStopWithoutLine(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		IGNORE_STOP_WITHOUT_LINE: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.IgnoreStopWithoutLine())
}

func Test_DiscoveryInterval_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	assert.Equal(time.Duration(-3600000000000), partnerSettings.DiscoveryInterval(), "should be set to -1 hour")
}

func Test_DiscoveryInterval(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		DISCOVERY_INTERVAL: "20m",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.Equal(time.Duration(-1200000000000), partnerSettings.DiscoveryInterval())
}

func Test_MessageIdentifierGenerator_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)
	assert.Equal("%{uuid}", partnerSettings.IdentifierGenerator.FormatString("Message"))
}

func Test_MessageIdentifierGenerator(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"generators.message_identifier": "message:%{uuid}:LOC",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal("message:%{uuid}:LOC", partnerSettings.IdentifierGenerator.FormatString("Message"))
}

func Test_ResponseMessageIdentifierGenerator_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)
	assert.Equal("%{uuid}", partnerSettings.IdentifierGenerator.FormatString("ResponseMessage"))
}

func Test_ResponseMessageIdentifierGenerator(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"generators.response_message_identifier": "response:message:%{uuid}:LOC",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal("response:message:%{uuid}:LOC", partnerSettings.IdentifierGenerator.FormatString("ResponseMessage"))
}

func Test_DataFrameIdentifierGenerator_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)
	assert.Equal("%{id}", partnerSettings.IdentifierGenerator.FormatString("DataFrame"))
}

func Test_DataFrameIdentifierGenerator(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"generators.data_frame_identifier": "dataframe:response:%{id}:LOC",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal("dataframe:response:%{id}:LOC", partnerSettings.IdentifierGenerator.FormatString("DataFrame"))
}

func Test_ReferenceIdentifierGenerator_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)
	assert.Equal("%{type}:%{id}", partnerSettings.IdentifierGenerator.FormatString("Reference"))
}

func Test_ReferenceIdentifierGenerator(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"generators.reference_identifier": "reference:%{type}:value:%{id}",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal("reference:%{type}:value:%{id}", partnerSettings.IdentifierGenerator.FormatString("Reference"))
}

func Test_ReferenceStopAreaIdentIdentifierGenerator_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)
	assert.Equal("%{id}", partnerSettings.IdentifierGenerator.FormatString("StopArea"))
}

func Test_ReferenceStopAreaIdentifierGenerator(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"generators.reference_stop_area_identifier": "stoparea:%{id}",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal("stoparea:%{id}", partnerSettings.IdentifierGenerator.FormatString("StopArea"))
}

func Test_SubscriptionIdentifierGenerator_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)
	assert.Equal("%{id}", partnerSettings.IdentifierGenerator.FormatString("Subscription"))
}

func Test_SubscriptionIdentifierGenerator(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"generators.subscription_identifier": "subscription:%{id}",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	assert.Equal("subscription:%{id}", partnerSettings.IdentifierGenerator.FormatString("Subscription"))
}

func Test_CollectSettings_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	expected := &CollectSettings{
		UseDiscoveredSA:    false,
		UseDiscoveredLines: false,
		includedSA:         collection{},
		excludedSA:         collection{},
		includedLines:      collection{},
		excludedLines:      collection{},
	}
	assert.Equal(expected, partnerSettings.CollectSettings())
}

func Test_CollectSettings(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"collect.use_discovered_stop_areas": "true",
		"collect.use_discovered_lines":      "true",
		"collect.include_stop_areas":        "STIF:StopPoint:Q:46647:,STIF:StopPoint:Q:555:",
		"collect.exclude_stop_areas":        "STIF:StopPoint:X:,STIF:StopPoint:J:",
		"collect.include_lines":             "Line:B:Tram,Line:A:Metro",
		"collect.exclude_lines":             "Line:C:Bus",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	expected := &CollectSettings{
		UseDiscoveredSA:    true,
		UseDiscoveredLines: true,
		includedSA:         collection{"STIF:StopPoint:Q:46647:": struct{}{}, "STIF:StopPoint:Q:555:": struct{}{}},
		excludedSA:         collection{"STIF:StopPoint:X:": struct{}{}, "STIF:StopPoint:J:": struct{}{}},
		includedLines:      collection{"Line:B:Tram": struct{}{}, "Line:A:Metro": struct{}{}},
		excludedLines:      collection{"Line:C:Bus": struct{}{}},
	}

	assert.Equal(expected, partnerSettings.CollectSettings())
}

func Test_HTTPClientOAuth_Empty(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)
	assert.Nil(partnerSettings.HTTPClientOAuth())
}

func Test_HTTPClientOAuth(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"remote_authentication.oauth.client_id":     "client_id",
		"remote_authentication.oauth.client_secret": "client_secret",
		"remote_authentication.oauth.token_url":     "https://token-url.com",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	expected := &remote.HTTPClientOAuth{
		ClientID:     "client_id",
		ClientSecret: "client_secret",
		TokenURL:     "https://token-url.com",
	}

	assert.EqualValues(expected, partnerSettings.HTTPClientOAuth())
}

func Test_HTTPClientOptions_Without_OAuth(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		"remote_url":               "remote_url",
		"notifications.remote_url": "notification-remote-url",
		"subscriptions.remote_url": "subscription-remote-url",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)
	expected := remote.HTTPClientOptions{
		SiriEnvelopeType: "soap",
		OAuth:            nil,
		SiriCredential: remote.SiriCredentialHeader{
			CredentialHeader: "X-SIRI-Requestor",
			Value:            "Ara",
		},
		Urls: remote.HTTPClientUrls{
			Url:              "remote_url",
			SubscriptionsUrl: "subscription-remote-url",
			NotificationsUrl: "notification-remote-url",
		},
	}
	assert.Equal(expected, partnerSettings.HTTPClientOptions())
}

func Test_SiriSoapEmptyResponseOnNotification_true(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		SIRI_SOAP_EMPTY_RESPONSE_ON_NOTIFICATION: "true",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.True(partnerSettings.SiriSoapEmptyResponseOnNotification())
}

func Test_SiriSoapEmptyResponseOnNotification_undefined(t *testing.T) {
	assert := assert.New(t)

	partnerSettings := NewEmptyPartnerSettings(uuid.DefaultUUIDGenerator)

	assert.False(partnerSettings.SiriSoapEmptyResponseOnNotification())
}

func Test_SiriSoapEmptyResponseOnNotification_false(t *testing.T) {
	assert := assert.New(t)

	settings := map[string]string{
		SIRI_SOAP_EMPTY_RESPONSE_ON_NOTIFICATION: "false",
	}
	partnerSettings := NewPartnerSettings(uuid.DefaultUUIDGenerator, settings)

	assert.False(partnerSettings.SiriSoapEmptyResponseOnNotification())
}
