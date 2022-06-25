package siri_tests

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func Test_SIRIGeneralMessageSubscriptionRequest_BuildXML(t *testing.T) {

	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &siri.SIRIGeneralMessageSubscriptionRequest{
		ConsumerAddress:   "https://ara-staging.af83.io/test/siri",
		MessageIdentifier: "test",
		RequestorRef:      "test",
		RequestTimestamp:  date,
	}

	entry := &siri.SIRIGeneralMessageSubscriptionRequestEntry{
		SubscriberRef:          "SubscriberRef",
		SubscriptionIdentifier: "SubscriptionIdentifier",
		InitialTerminationTime: date,
	}
	entry.MessageIdentifier = "test"
	entry.RequestTimestamp = date

	request.Entries = []*siri.SIRIGeneralMessageSubscriptionRequestEntry{entry}
	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}

	smsr, err := sxml.NewXMLSubscriptionRequestFromContent([]byte(xml))
	if err != nil {
		t.Fatal(err)
	}

	if smsr.RequestorRef() != request.RequestorRef {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", smsr.RequestorRef(), request.RequestorRef)
	}

	if smsr.MessageIdentifier() != request.MessageIdentifier {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", smsr.MessageIdentifier(), request.MessageIdentifier)
	}

	if smsr.RequestTimestamp() != request.RequestTimestamp {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", smsr.RequestTimestamp(), request.RequestTimestamp)
	}

	if smsr.ConsumerAddress() != request.ConsumerAddress {
		t.Errorf("Wrong ConsumerAddress:\n got: %v\nwant: %v", smsr.ConsumerAddress(), request.ConsumerAddress)
	}

	xse := smsr.XMLSubscriptionGMEntries()
	if len(xse) != 1 {
		t.Errorf("Wrong number of subscriptions entries :\n got: %v\nwant: %v", len(xse), 1)
	}

	xmlEntry := xse[0]

	if xmlEntry.MessageIdentifier() != entry.MessageIdentifier {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", xmlEntry.MessageIdentifier(), entry.MessageIdentifier)
	}

	if xmlEntry.RequestTimestamp() != entry.RequestTimestamp {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", xmlEntry.RequestTimestamp(), entry.RequestTimestamp)
	}

	if xmlEntry.SubscriberRef() != entry.SubscriberRef {
		t.Errorf("Wrong SubscriberRef:\n got: %v\nwant: %v", xmlEntry.SubscriberRef(), entry.SubscriberRef)
	}

	if xmlEntry.SubscriptionIdentifier() != entry.SubscriptionIdentifier {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", xmlEntry.SubscriptionIdentifier(), entry.SubscriptionIdentifier)
	}

	if xmlEntry.InitialTerminationTime() != entry.InitialTerminationTime {
		t.Errorf("Wrong InitialTerminationTime:\n got: %v\nwant: %v", xmlEntry.InitialTerminationTime(), entry.InitialTerminationTime)
	}
}
