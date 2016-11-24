package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getStopMonitoringResponseBody(t *testing.T) []byte {
	file, err := os.Open("testdata/stopmonitoring-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	return content
}

func Test_XMLStopMonitoringRequest_Address(t *testing.T) {
	content := getStopMonitoringResponseBody(t)

	response, _ := NewXMLStopMonitoringResponseFromContent(content)
	if expected := "http://appli.chouette.mobi/siri_france/siri"; response.Address() != expected {
		t.Errorf("Wrong Address:\n got: %v\nwant: %v", response.Address(), expected)
	}
}

func Test_XMLStopMonitoringRequest_ProducerRef(t *testing.T) {
	content := getStopMonitoringResponseBody(t)

	response, _ := NewXMLStopMonitoringResponseFromContent(content)
	if expected := "NINOXE:default"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef:\n got: %v\nwant: %v", response.ProducerRef(), expected)
	}
}

func Test_XMLStopMonitoringRequest_RequestMessageRef(t *testing.T) {
	content := getStopMonitoringResponseBody(t)

	response, _ := NewXMLStopMonitoringResponseFromContent(content)
	if expected := "StopMonitoring:Test:0"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef:\n got: %v\nwant: %v", response.RequestMessageRef(), expected)
	}
}

func Test_XMLStopMonitoringRequest_ResponseMessageIdentifier(t *testing.T) {
	content := getStopMonitoringResponseBody(t)

	response, _ := NewXMLStopMonitoringResponseFromContent(content)
	if expected := "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Wrong ResponseMessageIdentifier:\n got: %v\nwant: %v", response.ResponseMessageIdentifier(), expected)
	}
}

func Test_XMLStopMonitoringRequest_ResponseTimestamp(t *testing.T) {
	content := getStopMonitoringResponseBody(t)

	response, _ := NewXMLStopMonitoringResponseFromContent(content)

	if expected := time.Date(2016, time.September, 22, 6, 01, 20, 227000000, time.UTC); !response.ResponseTimestamp().Equal(expected) {
		t.Errorf("Wrong ResponseTimestamp:\n got: %v\nwant: %v", response.ResponseTimestamp(), expected)
	}
}
