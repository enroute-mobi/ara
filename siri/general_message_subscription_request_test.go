package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLGeneralMessageSubscriptionRequest(t *testing.T) *XMLSubscriptionRequest {
	file, err := os.Open("../core/testdata/generalmessagesubscription-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := NewXMLSubscriptionRequestFromContent(content)
	return request
}

func Test_XMLGeneralMessageSubscriptionRequest(t *testing.T) {
	request := getXMLGeneralMessageSubscriptionRequest(t)
	entry := request.XMLSubscriptionGMEntries()[0]

	if len(request.XMLSubscriptionGMEntries()) != 1 {
		t.Errorf("Wrong len should be 1 got %v", len(request.XMLSubscriptionGMEntries()))
	}

	if expected := "RATPDEV:Concerto"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}

	if expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}

	if expected := "NINOXE:default"; entry.SubscriberRef() != expected {
		t.Errorf("Wrong SubscriberRef:\n got: %v\nwant: %v", entry.SubscriberRef(), expected)
	}

	if expected := "28679112-9dad-11d1-2-00c04fd430c8"; request.MessageIdentifier() != expected {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", request.MessageIdentifier(), expected)
	}

	if expected := "NINOXE:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"; entry.SubscriptionIdentifier() != expected {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", entry.SubscriptionIdentifier(), expected)
	}

	if expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC); entry.InitialTerminationTime() != expected {
		t.Errorf("Wrong InitialTerminationTime:\n got: %v\nwant: %v", entry.InitialTerminationTime(), expected)
	}
}

func Test_SIRIGeneralMessageSubscriptionRequest_BuildXML(t *testing.T) {

	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &SIRIGeneralMessageSubscriptionRequest{
		ConsumerAddress:   "https://ara-staging.af83.io/test/siri",
		MessageIdentifier: "test",
		RequestorRef:      "test",
		RequestTimestamp:  date,
	}

	entry := &SIRIGeneralMessageSubscriptionRequestEntry{
		SubscriberRef:          "SubscriberRef",
		SubscriptionIdentifier: "SubscriptionIdentifier",
		InitialTerminationTime: date,
	}
	entry.MessageIdentifier = "test"
	entry.RequestTimestamp = date

	request.Entries = []*SIRIGeneralMessageSubscriptionRequestEntry{entry}
	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}

	smsr, err := NewXMLSubscriptionRequestFromContent([]byte(xml))
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
