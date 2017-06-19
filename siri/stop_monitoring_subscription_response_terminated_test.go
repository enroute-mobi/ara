package siri

import (
	"io/ioutil"
	"os"
	"testing"
)

func getXMLStopMonitoringSubscriptionTerminatedNotification(t *testing.T) *XMLStopMonitoringSubscriptionTerminatedResponse {
	file, err := os.Open("testdata/subscription_terminated_notification-soap.xml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := NewXMLStopMonitoringSubscriptionTerminatedResponseFromContent(content)
	return response
}

func Test_XMLStopMonitoringSubscriptionTerminatedResponse(t *testing.T) {
	response := getXMLStopMonitoringSubscriptionTerminatedNotification(t)
	subscriptionTerminated := response.XMLSubscriptionTerminateds()[0]

	if expected := "KUBRICK"; subscriptionTerminated.ProducerRef() != expected {
		t.Errorf("Incorrect ProducerRef expected: %v\n got: %v", expected, subscriptionTerminated.ProducerRef())
	}

	if expected := "NADER"; subscriptionTerminated.SubscriberRef() != expected {
		t.Errorf("Incorrect SubscriberRef expected: %v\n got: %v", expected, subscriptionTerminated.SubscriberRef())
	}

	if expected := "6ba7b814-9dad-11d1-0-00c04fd430c8"; subscriptionTerminated.SubscriptionRef() != expected {
		t.Errorf("Incorrect SubscriptionRef expected: %v\n got: %v", expected, subscriptionTerminated.SubscriptionRef())
	}

	if expected := "Weekley restart"; subscriptionTerminated.ErrorDescription() != expected {
		t.Errorf("Incorrect ErrorNumber ErrorDescription: %v\n got: %v", expected, subscriptionTerminated.ErrorDescription())
	}
}
