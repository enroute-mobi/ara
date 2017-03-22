package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLStopDiscoveryRequest(t *testing.T) *XMLStopDiscoveryRequest {
	file, err := os.Open("testdata/stopdiscovery-request.xml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := NewXMLStopDiscoveryRequestFromContent(content)
	return request
}

func Test_XMLStopDiscoveryRequest_RequestorRef(t *testing.T) {
	request := getXMLStopDiscoveryRequest(t)

	if expected := "STIF"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}
}

func Test_XMLStopDiscoveryRequest_RequestTimestamp(t *testing.T) {
	request := getXMLStopDiscoveryRequest(t)

	if expected := time.Date(2017, time.March, 03, 11, 28, 00, 359000000, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}
}
