package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLStopPointsDiscoveryRequest(t *testing.T) *XMLStopPointsDiscoveryRequest {
	file, err := os.Open("testdata/stopdiscovery-request.xml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := NewXMLStopPointsDiscoveryRequestFromContent(content)
	return request
}

func Test_XMLStopPointsDiscoveryRequest_RequestorRef(t *testing.T) {
	request := getXMLStopPointsDiscoveryRequest(t)

	if expected := "STIF"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}
}

func Test_XMLStopPointsDiscoveryRequest_RequestTimestamp(t *testing.T) {
	request := getXMLStopPointsDiscoveryRequest(t)

	if expected := time.Date(2017, time.March, 03, 11, 28, 00, 359000000, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}
}
