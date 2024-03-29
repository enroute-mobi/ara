package core

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_PTT_checklines(t *testing.T) {
	assert := assert.New(t)

	// Test Setup
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"local_url":            "http://ara",
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := newSIRIProductionTimetableSubscriptionBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())

	line := referential.model.Lines().New()
	line.SetCode(model.NewCode("codeSpace", "NINOXE:Line:2:LOC"))
	line.Name = "lineName"
	line.Save()

	line2 := referential.model.Lines().New()
	line2.SetCode(model.NewCode("codeSpace", "NINOXE:Line:3:LOC"))
	line2.Name = "lineName2"
	line2.Save()

	line3 := referential.model.Lines().New()
	line3.SetCode(model.NewCode("AnotherCodeSpace", "NINOXE:Line:A:BUS"))
	line3.Name = "lineName3"
	line3.Save()

	// test request for subscription to all Lines having the same remote_code_space
	request := []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>" +
		"<Siri xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" version=\"2.0\" xmlns=\"http://www.siri.org.uk/siri\">" +
		"  <SubscriptionRequest>" +
		"      <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>" +
		"      <RequestorRef>NINOXE:default</RequestorRef>" +
		"      <ProductionTimetableSubscriptionRequest>" +
		"         <SubscriptionIdentifier>test1</SubscriptionIdentifier>" +
		"         <InitialTerminationTime>2022-02-10T02:00:00Z</InitialTerminationTime>" +
		"        <ProductionTimetableRequest>" +
		"            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>" +
		"            <PreviewInterval>PT3H0S</PreviewInterval>" +
		"         </ProductionTimetableRequest>" +
		"        <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>" +
		"      </ProductionTimetableSubscriptionRequest>" +
		"   </SubscriptionRequest>" +
		"</Siri>")

	subs, err := sxml.NewXMLSubscriptionRequestFromContent(request)
	if err != nil {
		t.Errorf("cannot parse xml: %s", err)
	}

	ptt := subs.XMLSubscriptionPTTEntries()
	lines, unknownLines := connector.checkLines(ptt[0])

	assert.Equal(len(lines), 2)
	assert.Equal(len(unknownLines), 0)

	// test subscription to a Line not having the same remote_code_space
	request1 := []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>" +
		"<Siri xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" version=\"2.0\" xmlns=\"http://www.siri.org.uk/siri\">" +
		"  <SubscriptionRequest>" +
		"      <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>" +
		"      <RequestorRef>NINOXE:default</RequestorRef>" +
		"      <ProductionTimetableSubscriptionRequest>" +
		"         <SubscriptionIdentifier>test1</SubscriptionIdentifier>" +
		"         <InitialTerminationTime>2022-02-10T02:00:00Z</InitialTerminationTime>" +
		"        <ProductionTimetableRequest>" +
		"            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>" +
		"            <Line><LineDirection><LineRef>NINOXE:Line:A:BUS</LineRef></LineDirection></Line>" +
		"            <PreviewInterval>PT3H0S</PreviewInterval>" +
		"         </ProductionTimetableRequest>" +
		"        <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>" +
		"      </ProductionTimetableSubscriptionRequest>" +
		"   </SubscriptionRequest>" +
		"</Siri>")

	subs1, err1 := sxml.NewXMLSubscriptionRequestFromContent(request1)
	if err1 != nil {
		t.Errorf("cannot parse xml: %s", err1)
	}
	ptt1 := subs1.XMLSubscriptionPTTEntries()
	lines1, unknownLines1 := connector.checkLines(ptt1[0])

	assert.Equal(len(lines1), 0)
	assert.Equal(len(unknownLines1), 1)

	// test subscription to multiple Lines with both remote_code_space from partner and unknown remote_code_space
	request2 := []byte("<?xml version=\"1.0\" encoding=\"utf-8\"?>" +
		"<Siri xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\" xmlns:xsi=\"http://www.w3.org/2001/XMLSchema-instance\" version=\"2.0\" xmlns=\"http://www.siri.org.uk/siri\">" +
		"  <SubscriptionRequest>" +
		"      <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>" +
		"      <RequestorRef>NINOXE:default</RequestorRef>" +
		"      <ProductionTimetableSubscriptionRequest>" +
		"         <SubscriptionIdentifier>test1</SubscriptionIdentifier>" +
		"         <InitialTerminationTime>2022-02-10T02:00:00Z</InitialTerminationTime>" +
		"        <ProductionTimetableRequest>" +
		"            <RequestTimestamp>2017-01-01T12:01:00.000Z</RequestTimestamp>" +
		"            <Line><LineDirection><LineRef>NINOXE:Line:A:BUS</LineRef></LineDirection></Line>" +
		"            <Line><LineDirection><LineRef>NINOXE:Line:3:LOC</LineRef></LineDirection></Line>" +
		"            <PreviewInterval>PT3H0S</PreviewInterval>" +
		"         </ProductionTimetableRequest>" +
		"        <ChangeBeforeUpdates>PT30S</ChangeBeforeUpdates>" +
		"      </ProductionTimetableSubscriptionRequest>" +
		"   </SubscriptionRequest>" +
		"</Siri>")

	subs2, err2 := sxml.NewXMLSubscriptionRequestFromContent(request2)
	if err2 != nil {
		t.Errorf("cannot parse xml: %s", err2)
	}
	ptt2 := subs2.XMLSubscriptionPTTEntries()
	lines2, unknownLines2 := connector.checkLines(ptt2[0])

	assert.Equal(len(lines2), 1)
	assert.Equal(len(unknownLines2), 1)
}
