package siri_tests

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLStopMonitoringSubscriptionRequest(t *testing.T) *sxml.XMLSubscriptionRequest {
	file, err := os.Open("testdata/stopmonitoringsubscription-request-with-lineref-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := sxml.NewXMLSubscriptionRequestFromContent(content)
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

	if expected := "Ara:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"; entry.SubscriptionIdentifier() != expected {
		t.Errorf("Wrong SubscriptionIdentifier:\n got: %v\nwant: %v", entry.SubscriptionIdentifier(), expected)
	}

	if expected := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC); entry.InitialTerminationTime() != expected {
		t.Errorf("Wrong InitialTerminationTime:\n got: %v\nwant: %v", entry.InitialTerminationTime(), expected)
	}

	if expected := "STIF:Line::C00064:"; entry.LineRef() != expected {
		t.Errorf("Wrong LineRef:\n got: %v\nwant: %v", entry.LineRef(), expected)
	}
}
