package siri_tests

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLSubscriptionResponse(t *testing.T) *sxml.XMLSubscriptionResponse {
	file, err := os.Open("testdata/stopmonitoringsubscription-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLSubscriptionResponseFromContent(content)
	return response
}

func Test_XMLSubscriptionResponse(t *testing.T) {
	response := getXMLSubscriptionResponse(t)
	responseStatus := response.ResponseStatus()[0]

	if expected := "28679112-9dad-11d1-2-00c04fd430c8"; response.RequestMessageRef() != expected {
		t.Errorf("Incorrect RequestMessageRef expected: %v\n got: %v", expected, response.RequestMessageRef())
	}

	if expected := "SQYBUS"; response.ResponderRef() != expected {
		t.Errorf("Incorrect ResponderRef expected: %v\n got: %v", expected, response.ResponderRef())
	}

	if expected := "RATPDEV:Concerto"; responseStatus.SubscriberRef() != expected {
		t.Errorf("Incorrect SubscriberRef expected: %v\n got: %v", expected, responseStatus.SubscriberRef())
	}

	if expected := "Ara:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC"; responseStatus.SubscriptionRef() != expected {
		t.Errorf("Incorrect SubscriptionRef expected: %v\n got: %v", expected, responseStatus.SubscriptionRef())
	}

	if expected := time.Date(2016, time.September, 22, 6, 01, 20, 227000000, time.UTC); !response.ResponseTimestamp().Equal(expected) {
		t.Errorf("Incorrect ResponseTimestamp expected: %v\n got: %v", expected, response.ResponseTimestamp())
	}

	if expected := time.Date(2016, time.September, 22, 6, 01, 20, 227000000, time.UTC); !response.ServiceStartedTime().Equal(expected) {
		t.Errorf("Incorrect ValidUntil expected: %v\n got: %v", expected, response.ServiceStartedTime())
	}

	if expected := time.Date(2016, time.September, 22, 6, 01, 20, 227000000, time.UTC); !responseStatus.ValidUntil().Equal(expected) {
		t.Errorf("Incorrect ValidUntil expected: %v\n got: %v", expected, responseStatus.ValidUntil())
	}

	if expected := true; responseStatus.Status() != expected {
		t.Errorf("Incorrect Status expected: %v\n got: %v", expected, responseStatus.Status())
	}
}
