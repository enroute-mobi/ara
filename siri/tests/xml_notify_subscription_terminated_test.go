package siri_tests

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLNotifySubscriptionTerminated(t *testing.T) *sxml.XMLNotifySubscriptionTerminated {
	file, err := os.Open("testdata/notify-subscription-terminated-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLNotifySubscriptionTerminatedFromContent(content)
	return response
}

func Test_XMLNotifySubscriptionTerminated(t *testing.T) {
	response := getXMLNotifySubscriptionTerminated(t)

	if expected := "RELAIS"; response.ProducerRef() != expected {
		t.Errorf("Incorrect ProducerRef expected: %v\n got: %v", expected, response.ProducerRef())
	}

	if expected := "test-address"; response.Address() != expected {
		t.Errorf("Incorrect Address expected: %v\n got: %v", expected, response.Address())
	}

	if expected := "RELAIS:ResponseMessage::b3814ad3-d282-40ba-8895-68226045feb7:LOC"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Incorrect ResponseMessageIdentifier expected: %v\n got: %v", expected, response.ResponseMessageIdentifier())
	}

	if expected := "STIF:Message::811125622:LOC"; response.RequestMessageRef() != expected {
		t.Errorf("Incorrect RequestMessageRef expected: %v\n got: %v", expected, response.RequestMessageRef())
	}

	if expected := time.Date(2020, time.January, 30, 10, 0, 0, 0, time.UTC); response.ResponseTimestamp() != expected {
		t.Errorf("Incorrect ResponseTimestamp expected: %v\n got: %v", expected, response.ResponseTimestamp())
	}

	if expected := "STIF"; response.SubscriberRef() != expected {
		t.Errorf("Incorrect SubscriberRef expected: %v\n got: %v", expected, response.SubscriberRef())
	}

	if expected := "6ba7b814-9dad-11d1-0-00c04fd430c8"; response.SubscriptionRef() != expected {
		t.Errorf("Incorrect SubscriptionRef expected: %v\n got: %v", expected, response.SubscriptionRef())
	}
}
