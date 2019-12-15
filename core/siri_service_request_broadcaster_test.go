package core

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/edwig/audit"
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

func Test_SIRISiriServiceRequestBroadcaster_NoConnectors(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.ConnectorTypes = []string{SIRI_SERVICE_REQUEST_BROADCASTER}
	partner.RefreshConnectors()
	c, _ := partner.Connector(SIRI_SERVICE_REQUEST_BROADCASTER)
	connector := c.(*SIRIServiceRequestBroadcaster)

	file, err := os.Open("testdata/siri-service-multiple-request-soap.xml")
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
	if response == nil {
		t.Fatalf("HandleRequests should return a response")
	}
	if len(response.StopMonitoringDeliveries) != 1 {
		t.Fatal("Response should have 1 StopMonitoring delivery")
	}
	if len(response.GeneralMessageDeliveries) != 1 {
		t.Fatal("Response should have 1 GeneralMessage delivery")
	}
	if len(response.EstimatedTimetableDeliveries) != 1 {
		t.Fatal("Response should have 1 EstimatedTimetable delivery")
	}

	if response.StopMonitoringDeliveries[0].Status {
		t.Error("Response status should be false, got true")
	}
	if response.StopMonitoringDeliveries[0].ErrorType != "CapabilityNotSupportedError" {
		t.Errorf("Response Errortype should be CapabilityNotSupportedError, got: %v", response.StopMonitoringDeliveries[0].ErrorType)
	}
	expected := "Can't find a StopMonitoringRequestBroadcaster connector"
	if response.StopMonitoringDeliveries[0].ErrorText != expected {
		t.Errorf("Wrong response Errortype:\n got: %v\n want: %v", response.StopMonitoringDeliveries[0].ErrorText, expected)
	}

	if response.GeneralMessageDeliveries[0].Status {
		t.Error("Response status should be false, got true")
	}
	if response.GeneralMessageDeliveries[0].ErrorType != "CapabilityNotSupportedError" {
		t.Errorf("Response Errortype should be CapabilityNotSupportedError, got: %v", response.GeneralMessageDeliveries[0].ErrorType)
	}
	expected = "Can't find a GeneralMessageRequestBroadcaster connector"
	if response.GeneralMessageDeliveries[0].ErrorText != expected {
		t.Errorf("Wrong response Errortype:\n got: %v\n want: %v", response.GeneralMessageDeliveries[0].ErrorText, expected)
	}

	if response.EstimatedTimetableDeliveries[0].Status {
		t.Error("Response status should be false, got true")
	}
	if response.EstimatedTimetableDeliveries[0].ErrorType != "CapabilityNotSupportedError" {
		t.Errorf("Response Errortype should be CapabilityNotSupportedError, got: %v", response.EstimatedTimetableDeliveries[0].ErrorType)
	}
	expected = "Can't find a EstimatedTimetableBroadcaster connector"
	if response.EstimatedTimetableDeliveries[0].ErrorText != expected {
		t.Errorf("Wrong response Errortype:\n got: %v\n want: %v", response.EstimatedTimetableDeliveries[0].ErrorText, expected)
	}
}

func Test_SIRISiriServiceRequestBroadcaster_HandleRequests(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"
	partner.ConnectorTypes = []string{
		SIRI_SERVICE_REQUEST_BROADCASTER,
		SIRI_STOP_MONITORING_REQUEST_BROADCASTER,
		SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER,
		SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER,
	}
	partner.RefreshConnectors()
	c, _ := partner.Connector(SIRI_SERVICE_REQUEST_BROADCASTER)
	connector := c.(*SIRIServiceRequestBroadcaster)
	connector.SIRIPartner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	objectid := model.NewObjectID("objectidKind", "boaarle")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Monitored = true
	stopArea.Save()

	line := referential.model.Lines().New()
	line.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:3:LOC"))
	line.Name = "lineName"
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vehicleJourney.SetObjectID(model.NewObjectID("objectidKind", "vehicleJourney"))
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	stopVisit.SetObjectID(model.NewObjectID("objectidKind", "stopVisit"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.PassageOrder = 1
	stopVisit.ArrivalStatus = "onTime"
	stopVisit.Schedules.SetArrivalTime("aimed", connector.Clock().Now().Add(1*time.Minute))
	stopVisit.Schedules.SetArrivalTime("expected", connector.Clock().Now().Add(1*time.Minute))
	stopVisit.Save()

	file, err := os.Open("testdata/siri-service-multiple-request-soap.xml")
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

	if response == nil {
		t.Fatalf("HandleRequests should return a response")
	}
	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.RequestMessageRef != "GetSIRIStopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
	if !response.Status {
		t.Errorf("Response has wrong status:\n got: %v\n expected: true", response.Status)
	}

	if len(response.StopMonitoringDeliveries) != 1 {
		t.Fatal("Response should have 1 StopMonitoring delivery")
	}
	if !response.StopMonitoringDeliveries[0].Status {
		xml, err := response.StopMonitoringDeliveries[0].BuildStopMonitoringDeliveryXML()
		if err != nil {
			t.Fatalf("Error whild building xml: %v", err)
		}
		t.Errorf("StopMonitoring delivery should have status true: %v", xml)
	}
	if len(response.GeneralMessageDeliveries) != 1 {
		t.Fatal("Response should have 1 GeneralMessage delivery")
	}
	if !response.GeneralMessageDeliveries[0].Status {
		xml, err := response.GeneralMessageDeliveries[0].BuildGeneralMessageDeliveryXML()
		if err != nil {
			t.Fatalf("Error whild building xml: %v", err)
		}
		t.Errorf("GeneralMessage delivery should have status true: %v", xml)
	}
	if len(response.EstimatedTimetableDeliveries) != 1 {
		t.Fatal("Response should have 1 EstimatedTimetable delivery")
	}
	if !response.EstimatedTimetableDeliveries[0].Status {
		xml, err := response.EstimatedTimetableDeliveries[0].BuildEstimatedTimetableDeliveryXML()
		if err != nil {
			t.Fatalf("Error whild building xml: %v", err)
		}
		t.Errorf("EstimatedTimetable delivery should have status true: %v", xml)
	}
}

func Test_SIRISiriServiceRequestBroadcaster_HandleRequestsStopAreaNotFound(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"
	partner.ConnectorTypes = []string{SIRI_SERVICE_REQUEST_BROADCASTER, SIRI_STOP_MONITORING_REQUEST_BROADCASTER}
	partner.RefreshConnectors()
	c, _ := partner.Connector(SIRI_SERVICE_REQUEST_BROADCASTER)
	connector := c.(*SIRIServiceRequestBroadcaster)
	connector.SIRIPartner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	file, err := os.Open("testdata/siri-service-smrequest-soap.xml")
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
	if response.ResponseMessageIdentifier != "Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
	if response.StopMonitoringDeliveries[0].Status || response.StopMonitoringDeliveries[1].Status {
		t.Errorf("Response deliveries have wrong status:\n got: %v %v\n expected: false", response.StopMonitoringDeliveries[0].Status, response.StopMonitoringDeliveries[1].Status)
	}
	if len(response.StopMonitoringDeliveries) != 2 {
		t.Errorf("Response has the wrong number of deliveries:\n got: %v\n expected: 2", len(response.StopMonitoringDeliveries))
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

	file, err := os.Open("testdata/siri-service-smrequest-soap.xml")
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
	if logStashEvent["requestorRef"] != "RATPDev:Concerto" {
		t.Errorf("Wrong requestorRef logged:\n got: %v\n expected: RATPDev:Concerto", logStashEvent["requestorRef"])
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
