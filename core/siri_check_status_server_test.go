package core

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	ps "bitbucket.org/enroute-mobi/ara/core/partner_settings"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_SIRICheckStatusServer_CheckStatus(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("local_url", "http://ara")
	partner.SetSetting("generators.response_message_identifier", "Ara:ResponseMessage::%{uuid}:LOC")
	partner.Save()
	referential.Start()
	referential.Stop()
	connector := NewSIRICheckStatusServer(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	file, err := os.Open("testdata/checkstatus_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLCheckStatusRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, err := connector.CheckStatus(request, &audit.BigQueryMessage{})
	if err != nil {
		t.Fatal(err)
	}

	time := clock.DefaultClock().Now()
	if response.Address != "http://ara" {
		t.Errorf("Wrong Address in response:\n got: %v\n want: http://ara", response.Address)
	}
	if response.ProducerRef != "Ara" {
		t.Errorf("Wrong ProducerRef in response:\n got: %v\n want: Ara", response.ProducerRef)
	}
	if response.RequestMessageRef != "CheckStatus:Test:0" {
		t.Errorf("Wrong RequestMessageRef in response:\n got: %v\n want: CheckStatus:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Wrong ResponseMessageIdentifier in response:\n got: %v\n want: Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	if !response.Status {
		t.Errorf("Wrong Status in response:\n got: %v\n want: true", response.Status)
	}
	if response.ResponseTimestamp != time {
		t.Errorf("Wrong Address in response:\n got: %v\n want: %v", response.ResponseTimestamp, time)
	}
	if response.ServiceStartedTime != time {
		t.Errorf("Wrong ServiceStartedTime in response:\n got: %v\n want: %v", response.ServiceStartedTime, time)
	}
}

func Test_SIRICheckStatusServerFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		ConnectorTypes: []string{"siri-check-status-server"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = ps.NewPartnerSettings(partner.UUIDGenerator)
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have an error when local_credential isn't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"local_credential": "local_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when local_credential is set, got: %v", apiPartner.Errors)
	}
}

func Test_SIRICheckStatusServer_LogCheckStatusRequest(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	file, err := os.Open("testdata/checkstatus_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLCheckStatusRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	logXMLCheckStatusRequest(logStashEvent, request)

	if logStashEvent["requestorRef"] != "NINOXE:default" {
		t.Errorf("Wrong requestorRef logged:\n got: %v\n expected: NINOXE:default", logStashEvent["requestorRef"])
	}
	if logStashEvent["messageIdentifier"] != "CheckStatus:Test:0" {
		t.Errorf("Wrong messageIdentifier logged:\n got: %v\n expected: CheckStatus:Test:0", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["requestTimestamp"] != "2016-09-07 09:11:25.174 +0000 UTC" {
		t.Errorf("Wrong requestTimestamp logged:\n got: %v\n expected: 2016-09-22 07:58:34 +0200 CEST", logStashEvent["requestTimestamp"])
	}
}

func Test_SIRICheckStatusServer_LogCheckStatusResponse(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)
	time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	response := &siri.SIRICheckStatusResponse{
		Status:             true,
		ServiceStartedTime: time,
	}
	response.Address = "Address"
	response.ProducerRef = "ProducerRef"
	response.RequestMessageRef = "RequestMessageRef"
	response.ResponseMessageIdentifier = "ResponseMessageIdentifier"
	response.ResponseTimestamp = time

	logSIRICheckStatusResponse(logStashEvent, response)

	if logStashEvent["address"] != "Address" {
		t.Errorf("Wrong Address logged:\n got: %v\n expected: Address", logStashEvent["address"])
	}
	if logStashEvent["producerRef"] != "ProducerRef" {
		t.Errorf("Wrong ProducerRef logged:\n got: %v\n expected: ProducerRef", logStashEvent["producerRef"])
	}
	if logStashEvent["requestMessageRef"] != "RequestMessageRef" {
		t.Errorf("Wrong RequestMessageRef logged:\n got: %v\n expected: RequestMessageRef", logStashEvent["requestMessageRef"])
	}
	if logStashEvent["responseMessageIdentifier"] != "ResponseMessageIdentifier" {
		t.Errorf("Wrong ResponseMessageIdentifier logged:\n got: %v\n expected: ResponseMessageIdentifier", logStashEvent["responseMessageIdentifier"])
	}
	if logStashEvent["status"] != "true" {
		t.Errorf("Wrong Status logged:\n got: %v\n expected: true", logStashEvent["status"])
	}
	if expected := time.String(); logStashEvent["responseTimestamp"] != expected {
		t.Errorf("Wrong ResponseTimestamp logged:\n got: %v\n expected: %v", logStashEvent["responseTimestamp"], expected)
	}
	if expected := time.String(); logStashEvent["serviceStartedTime"] != expected {
		t.Errorf("Wrong ServiceStartedTime logged:\n got: %v\n expected: %v", logStashEvent["serviceStartedTime"], expected)
	}
	xml, err := response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if logStashEvent["responseXML"] != xml {
		t.Errorf("Wrong responseXML logged:\n got: %v\n expected: %v", logStashEvent["responseXML"], xml)
	}
}
