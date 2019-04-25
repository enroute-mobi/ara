package siri

import (
	"io/ioutil"
	"os"
	"testing"
)

func getXMLStopMonitoringSubscriptionTerminatedNotification(t *testing.T) *XMLStopMonitoringSubscriptionTerminatedResponse {
	file, err := os.Open("testdata/subscription_terminated_notification-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := NewXMLStopMonitoringSubscriptionTerminatedResponseFromContent(content)
	return response
}

func Test_XMLStopMonitoringSubscriptionTerminatedResponse(t *testing.T) {
	response := getXMLStopMonitoringSubscriptionTerminatedNotification(t)

	if expected := "KUBRICK"; response.ProducerRef() != expected {
		t.Errorf("Incorrect ProducerRef expected: %v\n got: %v", expected, response.ProducerRef())
	}

	if expected := "NADER"; response.SubscriberRef() != expected {
		t.Errorf("Incorrect SubscriberRef expected: %v\n got: %v", expected, response.SubscriberRef())
	}

	if expected := "6ba7b814-9dad-11d1-0-00c04fd430c8"; response.SubscriptionRef() != expected {
		t.Errorf("Incorrect SubscriptionRef expected: %v\n got: %v", expected, response.SubscriptionRef())
	}

	if expected := "Weekley restart"; response.ErrorDescription() != expected {
		t.Errorf("Incorrect ErrorNumber ErrorDescription: %v\n got: %v", expected, response.ErrorDescription())
	}
}
