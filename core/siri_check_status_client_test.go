package core

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/siri"
)

func prepare_siriCheckStatusClient(t *testing.T, responseFilePath string) OperationnalStatus {
	audit.SetCurrentLogstash(audit.NewFakeLogStash())
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}
		file, err := os.Open(responseFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))
	defer ts.Close()

	// Create a CheckStatusClient
	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"remote_url": ts.URL,
		},
	}
	checkStatusClient := NewSIRICheckStatusClient(partner)

	status, err := checkStatusClient.Status()

	if err != nil {
		t.Fatal(err)
	}

	return status
}

func testLogStash(t *testing.T) {
	events := audit.CurrentLogStash().(*audit.FakeLogStash).Events()
	if len(events) != 1 {
		t.Errorf("Logstash should have recieved an event, got: %v", events)
	}
	if len(events[0]) != 13 {
		t.Errorf("LogstashEvent should have 13 values, got: %v", events[0])
	}
}

func Test_SIRICheckStatusClient_Status_OK(t *testing.T) {
	status := prepare_siriCheckStatusClient(t, "testdata/checkstatus-response-soap.xml")
	if status != OPERATIONNAL_STATUS_UP {
		t.Errorf("Wrong status found:\n got: %v\n expected: up", status)
	}
	testLogStash(t)
}

func Test_SIRICheckStatusClient_Status_KO(t *testing.T) {
	status := prepare_siriCheckStatusClient(t, "testdata/checkstatus-negative-response-soap.xml")
	if status != OPERATIONNAL_STATUS_DOWN {
		t.Errorf("Wrong status found:\n got: %v\n expected: down", status)
	}
	testLogStash(t)
}

func Test_SIRICheckStatusClientFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		Settings:       make(map[string]string),
		ConnectorTypes: []string{"siri-check-status-client"},
		connectors:     make(map[string]Connector),
	}
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have two errors when remote_url isn't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_url":        "remote_url",
		"remote_credential": "remote_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when remote_url is set, got: %v", apiPartner.Errors)
	}
}

func Test_SIRICheckStatusClient_LogRequest(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)
	time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      "Edwig",
		RequestTimestamp:  time,
		MessageIdentifier: "0000-0000-0000-0000",
	}
	logRequest(logStashEvent, request)
	if logStashEvent["messageIdentifier"] != "0000-0000-0000-0000" {
		t.Errorf("Wrong messageIdentifier logged:\n got: %v\n expected: 0000-0000-0000-0000", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["requestorRef"] != "Edwig" {
		t.Errorf("Wrong requestorRef logged:\n got: %v\n expected: Edwig", logStashEvent["requestorRef"])
	}
	if expected := time.String(); logStashEvent["requestTimestamp"] != expected {
		t.Errorf("Wrong requestTimestamp logged:\n got: %v\n expected: %v", logStashEvent["requestTimestamp"], expected)
	}
	if logStashEvent["requestXML"] != request.BuildXML() {
		t.Errorf("Wrong requestXML logged:\n got: %v\n expected: %v", logStashEvent["requestXML"], request.BuildXML())
	}
}

func Test_SIRICheckStatusClient_LogResponse(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	file, err := os.Open("testdata/checkstatus-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	response, err := siri.NewXMLCheckStatusResponseFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	logResponse(logStashEvent, response)

	if logStashEvent["address"] != "http://appli.chouette.mobi/siri_france/siri" {
		t.Errorf("Wrong address logged:\n got: %v\n expected: http://appli.chouette.mobi/siri_france/siri", logStashEvent["address"])
	}
	if logStashEvent["producerRef"] != "NINOXE:default" {
		t.Errorf("Wrong producerRef logged:\n got: %v\n expected: NINOXE:default", logStashEvent["producerRef"])
	}
	if logStashEvent["requestMessageRef"] != "CheckStatus:Test:0" {
		t.Errorf("Wrong requestMessageRef logged:\n got: %v\n expected: CheckStatus:Test:0", logStashEvent["requestMessageRef"])
	}
	if logStashEvent["responseMessageIdentifier"] != "c464f588-5128-46c8-ac3f-8b8a465692ab" {
		t.Errorf("Wrong responseMessageIdentifier logged:\n got: %v\n expected: c464f588-5128-46c8-ac3f-8b8a465692ab", logStashEvent["responseMessageIdentifier"])
	}
	if logStashEvent["status"] != "true" {
		t.Errorf("Wrong status logged:\n got: %v\n expected: true", logStashEvent["status"])
	}
	if logStashEvent["responseTimestamp"] != "2016-09-22 07:58:34 +0200 CEST" {
		t.Errorf("Wrong responseTimestamp logged:\n got: %v\n expected: 2016-09-22 07:58:34 +0200 CEST", logStashEvent["responseTimestamp"])
	}
	if logStashEvent["serviceStartedTime"] != "2016-09-22 03:30:32 +0200 CEST" {
		t.Errorf("Wrong serviceStartedTime logged:\n got: %v\n expected: 2016-09-22 03:30:32 +0200 CEST", logStashEvent["serviceStartedTime"])
	}
	if logStashEvent["responseXML"] != response.RawXML() {
		t.Errorf("Wrong responseXML logged:\n got: %v\n expected: %v", logStashEvent["responseXML"], response.RawXML())
	}
}
