package core

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_SIRISiriServiceRequestBroadcaster_NoConnectors(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")

	settings := map[string]string{"remote_objectid_kind": "objectidKind"}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{SIRI_SERVICE_REQUEST_BROADCASTER}
	partner.RefreshConnectors()
	c, _ := partner.Connector(SIRI_SERVICE_REQUEST_BROADCASTER)
	connector := c.(*SIRIServiceRequestBroadcaster)

	file, err := os.Open("testdata/siri-service-multiple-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLSiriServiceRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.HandleRequests(request, &audit.BigQueryMessage{})
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
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"remote_objectid_kind":                   "objectidKind",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{
		SIRI_SERVICE_REQUEST_BROADCASTER,
		SIRI_STOP_MONITORING_REQUEST_BROADCASTER,
		SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER,
		SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER,
	}
	partner.RefreshConnectors()
	c, _ := partner.Connector(SIRI_SERVICE_REQUEST_BROADCASTER)
	connector := c.(*SIRIServiceRequestBroadcaster)
	connector.SetClock(clock.NewFakeClock())

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
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLSiriServiceRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.HandleRequests(request, &audit.BigQueryMessage{})

	if response == nil {
		t.Fatalf("HandleRequests should return a response")
	}
	if response.ProducerRef != "Ara" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Ara", response.ProducerRef)
	}
	if response.RequestMessageRef != "GetSIRIStopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
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
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"remote_objectid_kind":                   "objectidKind",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.ConnectorTypes = []string{SIRI_SERVICE_REQUEST_BROADCASTER, SIRI_STOP_MONITORING_REQUEST_BROADCASTER}
	partner.RefreshConnectors()
	c, _ := partner.Connector(SIRI_SERVICE_REQUEST_BROADCASTER)
	connector := c.(*SIRIServiceRequestBroadcaster)
	connector.SetClock(clock.NewFakeClock())

	file, err := os.Open("testdata/siri-service-smrequest-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLSiriServiceRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.HandleRequests(request, &audit.BigQueryMessage{})

	if response.ProducerRef != "Ara" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Ara", response.ProducerRef)
	}
	if response.RequestMessageRef != "GetSIRIStopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
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
		ConnectorTypes: []string{"siri-service-request-broadcaster"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = s.NewEmptyPartnerSettings(partner.UUIDGenerator)
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
