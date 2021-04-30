package core

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/siri"
)

func prepare_siriCheckStatusClient(t *testing.T, responseFilePath string) PartnerStatus {
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
	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	partner := referential.Partners().New("slug")
	partner.SetSetting(REMOTE_URL, ts.URL)
	checkStatusClient := NewSIRICheckStatusClient(partner)

	partnerStatus, err := checkStatusClient.Status()
	if err != nil {
		t.Fatal(err)
	}

	return partnerStatus
}

func Test_SIRICheckStatusClient_Status_OK(t *testing.T) {
	partnerStatus := prepare_siriCheckStatusClient(t, "testdata/checkstatus-response-soap.xml")
	if partnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
		t.Errorf("Wrong status found:\n got: %v\n expected: up", partnerStatus.OperationnalStatus)
	}
	events := audit.CurrentLogStash().(*audit.FakeLogStash).Events()
	if len(events) != 1 {
		t.Errorf("Logstash should have recieved an event, got: %v", events)
	}
	if len(events[0]) != 17 {
		t.Errorf("LogstashEvent should have 16 values, got %v:\n%v", len(events[0]), events[0])
	}
}

func Test_SIRICheckStatusClient_Status_KO(t *testing.T) {
	partnerStatus := prepare_siriCheckStatusClient(t, "testdata/checkstatus-negative-response-soap.xml")
	if partnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_DOWN {
		t.Errorf("Wrong status found:\n got: %v\n expected: down", partnerStatus.OperationnalStatus)
	}
	events := audit.CurrentLogStash().(*audit.FakeLogStash).Events()
	if len(events) != 1 {
		t.Errorf("Logstash should have recieved an event, got: %v", events)
	}
	if len(events[0]) != 20 {
		t.Errorf("LogstashEvent should have 19 values, got %v:\n%v", len(events[0]), events[0])
	}
}

func Test_SIRICheckStatusClientFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		ConnectorTypes: []string{"siri-check-status-client"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = NewPartnerSettings(partner)
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

func Test_SIRICheckStatusClient_LogCheckStatusRequest(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)
	time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      "Ara",
		RequestTimestamp:  time,
		MessageIdentifier: "0000-0000-0000-0000",
	}
	logSIRICheckStatusRequest(logStashEvent, &audit.BigQueryMessage{}, request)
	if logStashEvent["messageIdentifier"] != "0000-0000-0000-0000" {
		t.Errorf("Wrong messageIdentifier logged:\n got: %v\n expected: 0000-0000-0000-0000", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["requestorRef"] != "Ara" {
		t.Errorf("Wrong requestorRef logged:\n got: %v\n expected: Ara", logStashEvent["requestorRef"])
	}
	if expected := time.String(); logStashEvent["requestTimestamp"] != expected {
		t.Errorf("Wrong requestTimestamp logged:\n got: %v\n expected: %v", logStashEvent["requestTimestamp"], expected)
	}
	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if logStashEvent["requestXML"] != xml {
		t.Errorf("Wrong requestXML logged:\n got: %v\n expected: %v", logStashEvent["requestXML"], xml)
	}
}

func Test_SIRICheckStatusClient_LogCheckStatusResponse(t *testing.T) {
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

	logXMLCheckStatusResponse(logStashEvent, &audit.BigQueryMessage{}, response)

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
