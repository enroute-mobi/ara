package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLStopMonitoringSubscriptionResponse(t *testing.T) *XMLStopMonitoringSubscriptionResponse {
	file, err := os.Open("testdata/stopmonitoringsubscription-response-soap.xml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := NewXMLStopMonitoringSubscriptionResponseFromContent(content)
	return response
}

func Test_XMLStopMonitoringSubscriptionResponse(t *testing.T) {
	response := getXMLStopMonitoringSubscriptionResponse(t)

	if expected := "28679112-9dad-11d1-2-00c04fd430c8"; response.RequestMessageRef() != expected {
		t.Errorf("Incorrect RequestMessageRef expected: %v\n got: %v", expected, response.RequestMessageRef())
	}

	if expected := "SQYBUS"; response.ResponderRef() != expected {
		t.Errorf("Incorrect ResponderRef expected: %v\n got: %v", expected, response.ResponderRef())
	}

	if expected := "RATPDEV:Concerto"; response.SubscriberRef() != expected {
		t.Errorf("Incorrect SubscriberRef expected: %v\n got: %v", expected, response.SubscriberRef())
	}

	if expected := "Edwig:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"; response.SubscriptionRef() != expected {
		t.Errorf("Incorrect SubscriptionRef expected: %v\n got: %v", expected, response.SubscriptionRef())
	}

	if expected := time.Date(2016, time.September, 22, 6, 01, 20, 227000000, time.UTC); !response.ResponseTimestamp().Equal(expected) {
		t.Errorf("Incorrect ResponseTimestamp expected: %v\n got: %v", expected, response.ResponseTimestamp())
	}

	if expected := time.Date(2016, time.September, 22, 6, 01, 20, 227000000, time.UTC); !response.ServiceStartedTime().Equal(expected) {
		t.Errorf("Incorrect ValidUntil expected: %v\n got: %v", expected, response.ServiceStartedTime())
	}

	if expected := time.Date(2016, time.September, 22, 6, 01, 20, 227000000, time.UTC); !response.ValidUntil().Equal(expected) {
		t.Errorf("Incorrect ValidUntil expected: %v\n got: %v", expected, response.ValidUntil())
	}

	if expected := true; response.Status() != expected {
		t.Errorf("Incorrect Status expected: %v\n got: %v", expected, response.Status())
	}
}
