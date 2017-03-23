package core

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopArea(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	mid := NewFormatMessageIdentifierGenerator("Edwig:Message::%s:LOC")
	mid.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SIRIPartner().SetMessageIdentifierGenerator(mid)
	connector.SetClock(model.NewFakeClock())

	objectid := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	file, err := os.Open("testdata/stopmonitoring-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLStopMonitoringRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, err := connector.RequestStopArea(request)
	if err != nil {
		t.Fatal(err)
	}

	if response.Address != "http://edwig" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://edwig", response.Address)
	}
	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.RequestMessageRef != "StopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
}

func Test_SIRIStopMonitoringRequestBroadcasterFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		Settings:       make(map[string]string),
		ConnectorTypes: []string{"siri-stop-monitoring-request-broadcaster"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have errors when local_credential and remote_objectid_kind aren't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_objectid_kind": "remote_objectid_kind",
		"local_credential":     "local_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when local_credential and remote_objectid_kind are set, got: %v", apiPartner.Errors)
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_LogXMLStopMonitoringRequest(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	file, err := os.Open("testdata/stopmonitoring-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLStopMonitoringRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	logXMLStopMonitoringRequest(logStashEvent, request)
	if logStashEvent["Connector"] != "StopMonitoringRequestBroadcaster" {
		t.Errorf("Wrong Connector logged:\n got: %v\n expected: StopMonitoringRequestBroadcaster", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["messageIdentifier"] != "StopMonitoring:Test:0" {
		t.Errorf("Wrong messageIdentifier logged:\n got: %v\n expected: StopMonitoring:Test:0", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["requestorRef"] != "NINOXE:default" {
		t.Errorf("Wrong requestorRef logged:\n got: %v\n expected: NINOXE:default", logStashEvent["requestorRef"])
	}
	if logStashEvent["monitoringRef"] != "NINOXE:StopPoint:SP:24:LOC" {
		t.Errorf("Wrong monitoringRef logged:\n got: %v\n expected: NINOXE:StopPoint:SP:24:LOC", logStashEvent["monitoringRef"])
	}
	if logStashEvent["requestTimestamp"] != "2016-09-22 07:54:52.977 +0000 UTC" {
		t.Errorf("Wrong requestTimestamp logged:\n got: %v\n expected: 2016-09-22 07:54:52.977 +0000 UTC", logStashEvent["requestTimestamp"])
	}
	if logStashEvent["requestXML"] != request.RawXML() {
		t.Errorf("Wrong requestXML logged:\n got: %v\n expected: %v", logStashEvent["requestXML"], request.RawXML())
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_LogSIRIStopMonitoringResponse(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	response := &siri.SIRIStopMonitoringResponse{
		Address:                   "edwig.edwig",
		ProducerRef:               "NINOXE:default",
		ResponseMessageIdentifier: "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26",
	}
	response.RequestMessageRef = "StopMonitoring:Test:0"
	response.Status = true
	response.ResponseTimestamp = time

	logSIRIStopMonitoringResponse(logStashEvent, response)

	if logStashEvent["address"] != "edwig.edwig" {
		t.Errorf("Wrong address logged:\n got: %v\n expected: http://appli.chouette.mobi/siri_france/siri", logStashEvent["address"])
	}
	if logStashEvent["producerRef"] != "NINOXE:default" {
		t.Errorf("Wrong producerRef logged:\n got: %v\n expected: NINOXE:default", logStashEvent["producerRef"])
	}
	if logStashEvent["requestMessageRef"] != "StopMonitoring:Test:0" {
		t.Errorf("Wrong requestMessageRef logged:\n got: %v\n expected: StopMonitoring:Test:0", logStashEvent["requestMessageRef"])
	}
	if logStashEvent["responseMessageIdentifier"] != "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26" {
		t.Errorf("Wrong responseMessageIdentifier logged:\n got: %v\n expected: fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26", logStashEvent["responseMessageIdentifier"])
	}
	if logStashEvent["responseTimestamp"] != "2009-11-10 23:00:00 +0000 UTC" {
		t.Errorf("Wrong responseTimestamp logged:\n got: %v\n expected: 2009-11-10 23:00:00 +0000 UTC", logStashEvent["responseTimestamp"])
	}
	xml, err := response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if logStashEvent["responseXML"] != xml {
		t.Errorf("Wrong responseXML logged:\n got: %v\n expected: %v", logStashEvent["responseXML"], xml)
	}
}
