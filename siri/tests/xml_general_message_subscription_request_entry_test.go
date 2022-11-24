package siri_tests

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLGeneralMessageSubscriptionRequest(t *testing.T) *sxml.XMLSubscriptionRequest {
	file, err := os.Open("testdata/generalmessagesubscription-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := sxml.NewXMLSubscriptionRequestFromContent(content)
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
