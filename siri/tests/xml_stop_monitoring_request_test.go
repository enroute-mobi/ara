package siri_tests

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLGetStopMonitoring(t *testing.T) *sxml.XMLGetStopMonitoring {
	file, err := os.Open("testdata/stopmonitoring-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := sxml.NewXMLGetStopMonitoringFromContent(content)
	return request
}

func Test_XMLGetStopMonitoring_RequestorRef(t *testing.T) {
	request := getXMLGetStopMonitoring(t)

	if expected := "NINOXE:default"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}
}

func Test_XMLGetStopMonitoring_RequestTimestamp(t *testing.T) {
	request := getXMLGetStopMonitoring(t)

	if expected := time.Date(2016, time.September, 22, 7, 54, 52, 977000000, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}
}

func Test_XMLGetStopMonitoring_MessageIdentifier(t *testing.T) {
	request := getXMLGetStopMonitoring(t)

	if expected := "StopMonitoring:Test:0"; request.MessageIdentifier() != expected {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", request.MessageIdentifier(), expected)
	}
}

func Test_XMLGetStopMonitoring_MonitoringRef(t *testing.T) {
	request := getXMLGetStopMonitoring(t)

	if expected := "NINOXE:StopPoint:SP:24:LOC"; request.MonitoringRef() != expected {
		t.Errorf("Wrong MonitoringRef:\n got: %v\nwant: %v", request.MonitoringRef(), expected)
	}
}
