package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_SIRISiriServiceRequestBroadcaster_HandleRequests(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	connector := NewSIRIServiceRequestBroadcaster(partner)
	mid := NewFormatMessageIdentifierGenerator("Edwig:Message::%s:LOC")
	mid.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SIRIPartner().SetMessageIdentifierGenerator(mid)
	connector.SetClock(model.NewFakeClock())

	objectid := model.NewObjectID("objectidKind", "boaarle")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	objectid2 := model.NewObjectID("objectidKind", "cladebr")
	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetObjectID(objectid2)
	stopArea2.Save()

	file, err := os.Open("testdata/siri-service-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLSiriServiceRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.HandleRequests(request)

	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.RequestMessageRef != "GetSIRIStopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
	if !response.Status {
		fmt.Println(referential.Model().StopAreas().FindAll())
		t.Errorf("Response has wrong status:\n got: %v\n expected: true", response.Status)
	}
	if len(response.Deliveries) != 2 {
		t.Errorf("Response has the wrong number of deliveries:\n got: %v\n expected: 2", len(response.Deliveries))
	}
}

func Test_SIRISiriServiceRequestBroadcaster_HandleRequestsNotFound(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	connector := NewSIRIServiceRequestBroadcaster(partner)
	mid := NewFormatMessageIdentifierGenerator("Edwig:Message::%s:LOC")
	mid.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SIRIPartner().SetMessageIdentifierGenerator(mid)
	connector.SetClock(model.NewFakeClock())

	file, err := os.Open("testdata/siri-service-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLSiriServiceRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.HandleRequests(request)

	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.RequestMessageRef != "GetSIRIStopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
	if response.Deliveries[0].Status || response.Deliveries[1].Status {
		t.Errorf("Response deliveries have wrong status:\n got: %v %v\n expected: false", response.Deliveries[0].Status, response.Deliveries[1].Status)
	}
	if len(response.Deliveries) != 2 {
		t.Errorf("Response has the wrong number of deliveries:\n got: %v\n expected: 2", len(response.Deliveries))
	}
}

func Test_SIRIServiceRequestBroadcasterFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		Settings:       make(map[string]string),
		ConnectorTypes: []string{"siri-service-request-broadcaster"},
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

func Test_SIRIServiceRequestBroadcaster_LogXMLSiriServiceRequest(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	file, err := os.Open("testdata/siri-service-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLSiriServiceRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	logXMLSiriServiceRequest(logStashEvent, request)
	if logStashEvent["messageIdentifier"] != "GetSIRIStopMonitoring:Test:0" {
		t.Errorf("Wrong messageIdentifier logged:\n got: %v\n expected: GetSIRIStopMonitoring:Test:0", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["requestorRef"] != "RATPDEV:Concerto" {
		t.Errorf("Wrong requestorRef logged:\n got: %v\n expected: RATPDEV:Concerto", logStashEvent["requestorRef"])
	}
	if logStashEvent["requestTimestamp"] != "2001-12-17 09:30:47 +0000 UTC" {
		t.Errorf("Wrong requestTimestamp logged:\n got: %v\n expected: 2001-12-17 09:30:47 +0000 UTC", logStashEvent["requestTimestamp"])
	}
	if logStashEvent["requestXML"] != request.RawXML() {
		t.Errorf("Wrong requestXML logged:\n got: %v\n expected: %v", logStashEvent["requestXML"], request.RawXML())
	}
}

func Test_SIRIServiceRequestBroadcaster_LogSIRIServiceResponse(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	response := &siri.SIRIServiceResponse{
		ProducerRef:               "NINOXE:default",
		ResponseMessageIdentifier: "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26",
	}
	response.RequestMessageRef = "StopMonitoring:Test:0"
	response.Status = true
	response.ResponseTimestamp = time

	logSIRIServiceResponse(logStashEvent, response)

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
		t.Errorf("Wrong responseTimestamp logged:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", logStashEvent["responseTimestamp"])
	}
	if logStashEvent["status"] != "true" {
		t.Errorf("Wrong status logged:\n got: %v\n expected: true", logStashEvent["true"])
	}
	xml, err := response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if logStashEvent["responseXML"] != xml {
		t.Errorf("Wrong responseXML logged:\n got: %v\n expected: %v", logStashEvent["responseXML"], xml)
	}
}
