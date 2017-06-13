package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLStopMonitoringSubscriptionRequest(t *testing.T) *XMLStopMonitoringSubscriptionRequest {
	file, err := os.Open("testdata/stopmonitoringsubscription-request.soap.xml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := NewXMLStopMonitoringSubscriptionRequestFromContent(content)
	return request
}

func Test_XMLStopMonitoringSubscriptionRequest(t *testing.T) {
	request := getXMLStopMonitoringSubscriptionRequest(t)

	if expected := "RATPDEV:Concerto"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}

	if expected := time.Date(2017, time.January, 01, 12, 0, 0, 0, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}

	if expected := "coicogn2"; request.MonitoringRef() != expected {
		t.Errorf("Wrong MonitoringRef:\n got: %v\nwant: %v", request.MonitoringRef(), expected)
	}

	if expected := "RATPDEV:Concerto"; request.SubscriberRef() != expected {
		t.Errorf("Wrong SubscriberRef:\n got: %v\nwant: %v", request.SubscriberRef(), expected)
	}

	if expected := "28679112-9dad-11d1-2-00c04fd430c8"; request.MessageIdentifier() != expected {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", request.MessageIdentifier(), expected)
	}

	if expected := "Edwig:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"; request.SubscriptionIdentifier() != expected {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", request.SubscriptionIdentifier(), expected)
	}

	if expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC); request.InitialTerminationTime() != expected {
		t.Errorf("Wrong InitialTerminationTime:\n got: %v\nwant: %v", request.InitialTerminationTime(), expected)
	}
}

func Test_SIRIStopMonitoringSubscriptionRequest_BuildXML(t *testing.T) {

	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &SIRIStopMonitoringSubscriptionRequest{
		MessageIdentifier:      "test",
		MonitoringRef:          "test",
		RequestorRef:           "test",
		RequestTimestamp:       date,
		ConsumerAddress:        "https://edwig-staging.af83.io/test/siri",
		SubscriberRef:          "SubscriberRef",
		SubscriptionIdentifier: "SubscriptionIdentifier",
		InitialTerminationTime: date,
	}

	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}

	smsr, err := NewXMLStopMonitoringSubscriptionRequestFromContent([]byte(xml))

	if smsr.RequestorRef() != request.RequestorRef {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", smsr.RequestorRef(), request.RequestorRef)
	}

	if smsr.MessageIdentifier() != request.MessageIdentifier {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", smsr.MessageIdentifier(), request.MessageIdentifier)
	}

	if smsr.ConsumerAddress() != request.ConsumerAddress {
		t.Errorf("Wrong ConsumerAddress:\n got: %v\nwant: %v", smsr.ConsumerAddress(), request.ConsumerAddress)
	}

	if smsr.MonitoringRef() != request.MonitoringRef {
		t.Errorf("Wrong MonitoringRef:\n got: %v\nwant: %v", smsr.MonitoringRef(), request.MonitoringRef)
	}

	if smsr.RequestTimestamp() != request.RequestTimestamp {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", smsr.RequestTimestamp(), request.RequestTimestamp)
	}

	if smsr.SubscriberRef() != request.SubscriberRef {
		t.Errorf("Wrong SubscriberRef:\n got: %v\nwant: %v", smsr.SubscriberRef(), request.SubscriberRef)
	}

	if smsr.SubscriptionIdentifier() != request.SubscriptionIdentifier {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", smsr.SubscriptionIdentifier(), request.SubscriptionIdentifier)
	}

	if smsr.InitialTerminationTime() != request.InitialTerminationTime {
		t.Errorf("Wrong InitialTerminationTime:\n got: %v\nwant: %v", smsr.InitialTerminationTime(), request.InitialTerminationTime)
	}
}
