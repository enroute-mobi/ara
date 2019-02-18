package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLStopMonitoringSubscriptionRequest(t *testing.T) *XMLSubscriptionRequest {
	file, err := os.Open("../core/testdata/stopmonitoringsubscription-request-with-lineref-soap.xml")
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

func Test_XMLStopMonitoringSubscriptionRequest(t *testing.T) {
	request := getXMLStopMonitoringSubscriptionRequest(t)
	entry := request.XMLSubscriptionSMEntries()[0]

	if expected := "RATPDEV:Concerto"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}

	if expected := time.Date(2017, time.January, 01, 12, 0, 0, 0, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}

	if expected := "RATPDEV:Concerto"; entry.SubscriberRef() != expected {
		t.Errorf("Wrong SubscriberRef:\n got: %v\nwant: %v", entry.SubscriberRef(), expected)
	}

	if expected := "28679112-9dad-11d1-2-00c04fd430c8"; request.MessageIdentifier() != expected {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", request.MessageIdentifier(), expected)
	}

	if expected := "Edwig:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"; entry.SubscriptionIdentifier() != expected {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", entry.SubscriptionIdentifier(), expected)
	}

	if expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC); entry.InitialTerminationTime() != expected {
		t.Errorf("Wrong InitialTerminationTime:\n got: %v\nwant: %v", entry.InitialTerminationTime(), expected)
	}

	if expected := "STIF:Line::C00064:"; entry.LineRef() != expected {
		t.Errorf("Wrong LineRef:\n got: %v\nwant: %v", entry.LineRef(), expected)
	}
}

func Test_SIRIStopMonitoringSubscriptionRequest_BuildXML(t *testing.T) {

	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &SIRIStopMonitoringSubscriptionRequest{
		ConsumerAddress:   "https://edwig-staging.af83.io/test/siri",
		MessageIdentifier: "test",
		RequestorRef:      "test",
		RequestTimestamp:  date,
	}

	entry := &SIRIStopMonitoringSubscriptionRequestEntry{
		SubscriberRef:          "SubscriberRef",
		SubscriptionIdentifier: "SubscriptionIdentifier",
		InitialTerminationTime: date,
	}
	entry.MessageIdentifier = "test"
	entry.MonitoringRef = "MonitoringRef"
	entry.RequestTimestamp = date

	request.Entries = append(request.Entries, entry)
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

	xse := smsr.XMLSubscriptionSMEntries()
	if len(xse) != 1 {
		t.Errorf("Wrong number of subscriptions entries :\n got: %v\nwant: %v", len(xse), 1)
	}

	xmlEntry := xse[0]

	if xmlEntry.MessageIdentifier() != entry.MessageIdentifier {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", xmlEntry.MessageIdentifier(), entry.MessageIdentifier)
	}

	if xmlEntry.MonitoringRef() != entry.MonitoringRef {
		t.Errorf("Wrong MonitoringRef:\n got: %v\nwant: %v", xmlEntry.MonitoringRef(), entry.MonitoringRef)
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
