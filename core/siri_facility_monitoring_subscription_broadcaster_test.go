package core

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FacilityMonitoringBroadcaster_Create_Events(t *testing.T) {
	assert := assert.New(t)

	clock.SetDefaultClock(clock.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewTestMemoryModel()

	referential.model.SetBroadcastFMChan(referential.broacasterManager.GetFacilityBroadcastEventChan())
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{TEST_FACILITY_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(TEST_FACILITY_MONITORING_SUBSCRIPTION_BROADCASTER)

	facility := referential.Model().Facilities().New()
	code := model.NewCode("internal", string(facility.Id()))
	facility.SetCode(code)
	facility.Save()

	reference := model.Reference{
		Code: &code,
		Type: "Facility",
	}

	subs := partner.Subscriptions().New(FacilityMonitoringBroadcast)
	subs.Save()
	subs.CreateAndAddNewResource(reference)
	subs.SetExternalId("externalId")
	subs.Save()

	time.Sleep(5 * time.Millisecond) // Wait for the goRoutine to start ...

	time.Sleep(5 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	assert.Len(connector.(*TestFMSubscriptionBroadcaster).events, 1, "1 event should have been generated")

	event := connector.(*TestFMSubscriptionBroadcaster).events[0]
	assert.Equal(event.ModelId, string(facility.Id()))
	assert.Equal(event.ModelType, "Facility")

}

func Test_checkFacilities(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Test Setup
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"local_url":         "http://ara",
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := newSIRIFacilityMonitoringSubscriptionBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())

	facility := referential.model.Facilities().New()
	facility.SetCode(model.NewCode("codeSpace", "NINOXE:Facility:1"))
	facility.Save()

	facility2 := referential.model.Facilities().New()
	facility2.SetCode(model.NewCode("codeSpace", "NINOXE:Facility:2"))
	facility2.Save()

	facility3 := referential.model.Facilities().New()
	facility3.SetCode(model.NewCode("AnotherCodeSpace", "NINOXE:Facility:3"))
	facility3.Save()

	// test request for subscription to all Facilities having the same remote_code_space
	request := []byte(`
    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:FacilityMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>subscription-1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:FacilityMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
              </siri:FacilityMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
            </siri:FacilityMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>`)

	subs, err := sxml.NewXMLSubscriptionRequestFromContent(request)
	require.NoError(err)

	fm := subs.XMLSubscriptionFMEntries()
	facilitys, unknownFacilities := connector.checkFacilities(fm[0])

	assert.Equal(len(facilitys), 2)
	assert.Equal(len(unknownFacilities), 0)

	// test subscription to a Facility not having the same remote_code_space
	request1 := []byte(`    <?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:FacilityMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>subscription-1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:FacilityMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
                <siri:FacilityRef>WRONG</siri:FacilityRef>
              </siri:FacilityMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
            </siri:FacilityMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>`)

	subs1, err1 := sxml.NewXMLSubscriptionRequestFromContent(request1)
	require.NoError(err1)

	fm1 := subs1.XMLSubscriptionFMEntries()
	facilitys1, unknownFacilities1 := connector.checkFacilities(fm1[0])

	assert.Equal(len(facilitys1), 0)
	assert.Equal(len(unknownFacilities1), 1)

	// test subscription to multiple Facilities with both remote_code_space from partner and unknown remote_code_space
	request2 := []byte(`<?xml version='1.0' encoding='utf-8'?>
    <S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
      <S:Body>
        <ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
          <SubscriptionRequestInfo>
            <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
            <siri:RequestorRef>NINOXE:default</siri:RequestorRef>
            <siri:MessageIdentifier>6ba7b814-9dad-11d1-7-00c04fd430c8</siri:MessageIdentifier>
          </SubscriptionRequestInfo>
          <Request>
            <siri:FacilityMonitoringSubscriptionRequest>
              <siri:SubscriberRef>test</siri:SubscriberRef>
              <siri:SubscriptionIdentifier>subscription-1</siri:SubscriptionIdentifier>
              <siri:InitialTerminationTime>2017-01-03T12:03:00.000Z</siri:InitialTerminationTime>
              <siri:FacilityMonitoringRequest version="2.0:FR-IDF-2.4">
                <siri:RequestTimestamp>2017-01-01T12:03:00.000Z</siri:RequestTimestamp>
                <siri:MessageIdentifier>6ba7b814-9dad-11d1-6-00c04fd430c8</siri:MessageIdentifier>
                <siri:FacilityRef>NINOXE:Facility:1</siri:FacilityRef>
                <siri:FacilityRef>NINOXE:Facility:3</siri:FacilityRef>
              </siri:FacilityMonitoringRequest>
              <siri:IncrementalUpdates>true</siri:IncrementalUpdates>
            </siri:FacilityMonitoringSubscriptionRequest>
          </Request>
          <RequestExtension />
        </ws:Subscribe>
      </S:Body>`)

	subs2, err2 := sxml.NewXMLSubscriptionRequestFromContent(request2)
	require.NoError(err2)

	fm2 := subs2.XMLSubscriptionFMEntries()
	facilitys2, unknownFacilities2 := connector.checkFacilities(fm2[0])

	assert.Equal(len(facilitys2), 1)
	assert.Equal(len(unknownFacilities2), 1)
}
