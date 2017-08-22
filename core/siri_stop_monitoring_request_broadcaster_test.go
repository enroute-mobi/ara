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

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaNoSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.SIRIPartner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	objectid := model.NewObjectID("objectidKind", "modelOperatorRef")

	operator := referential.Model().Operators().New()
	operator.SetObjectID(objectid)

	operator.Save()

	objectid = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	stopVisit := referential.model.StopVisits().New()
	stopVisitRef := model.Reference{}
	obj := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:25:LOC")
	stopVisitRef.ObjectId = &obj
	stopVisitRef.Id = string(operator.Id())

	stopVisit.SetObjectID(obj)
	stopVisit.References.Set("OperatorRef", stopVisitRef)
	stopVisit.StopAreaId = stopArea.Id()

	sVSchedule := model.StopVisitSchedule{}

	sVSchedule.SetArrivalTime(connector.Clock().Now().Add(10 * time.Minute))
	stopVisit.Schedules["actual"] = &sVSchedule
	stopVisit.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:26:LOC")
	vehicleJourney.SetObjectID(obj)
	vehicleJourney.Save()

	stopVisit.VehicleJourneyId = vehicleJourney.Id()

	line := referential.model.Lines().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:27:LOC")
	line.SetObjectID(obj)
	line.Save()

	vehicleJourney.LineId = line.Id()

	stopVisit2 := referential.model.StopVisits().New()
	stopVisitRef2 := model.Reference{}
	obj2 := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:28:LOC")
	stopVisitRef2.ObjectId = &obj2

	stopVisit2.SetObjectID(obj2)
	stopVisit2.References.Set("OperatorRef", stopVisitRef)
	stopVisit2.StopAreaId = stopArea.Id()

	sVSchedule2 := model.StopVisitSchedule{}

	sVSchedule2.SetArrivalTime(connector.Clock().Now().Add(10 * time.Minute))
	stopVisit2.Schedules["actual"] = &sVSchedule2
	stopVisit2.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	obj2 = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:29:LOC")
	vehicleJourney2.SetObjectID(obj2)
	vehicleJourney2.Save()

	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()

	line2 := referential.model.Lines().New()
	obj = model.NewObjectID("WrongKind", "NINOXE:StopPoint:SP:30:LOC")
	line2.SetObjectID(obj)
	line2.Save()

	vehicleJourney2.LineId = line2.Id()

	file, err := os.Open("testdata/stopmonitoring-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetStopMonitoringFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.RequestStopArea(request)

	if response.Address != "http://edwig" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://edwig", response.Address)
	}
	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.RequestMessageRef != "StopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}

	if len(response.MonitoredStopVisits) != 1 {
		t.Fatalf("Response.MonitoredStopVisits should be 1 is %v", len(response.MonitoredStopVisits))
	}

	if response.MonitoredStopVisits[0].References["StopVisitReferences"]["OperatorRef"].ObjectId.Value() != "modelOperatorRef" {
		t.Errorf("OperatorRef should be modelOperatorRef, got: %v", response.MonitoredStopVisits[0].References["StopVisitReferences"]["OperatorRef"].ObjectId.Value())
	}

}

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaLineSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.SIRIPartner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	line := referential.model.Lines().New()
	obj := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:27:LOC")
	line.SetObjectID(obj)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:26:LOC")
	vehicleJourney.SetObjectID(obj)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:25:LOC")
	stopVisit.SetObjectID(obj)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	sVSchedule := model.StopVisitSchedule{}
	sVSchedule.SetArrivalTime(connector.Clock().Now().Add(10 * time.Minute))
	stopVisit.Schedules["actual"] = &sVSchedule
	stopVisit.Save()

	line2 := referential.model.Lines().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:7:LOC")
	line2.SetObjectID(obj)
	line.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:6:LOC")
	vehicleJourney2.SetObjectID(obj)
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit2 := referential.model.StopVisits().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:5:LOC")
	stopVisit2.SetObjectID(obj)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()
	sVSchedule2 := model.StopVisitSchedule{}
	sVSchedule2.SetArrivalTime(connector.Clock().Now().Add(15 * time.Minute))
	stopVisit.Schedules["actual"] = &sVSchedule2
	stopVisit2.Save()

	file, err := os.Open("testdata/stopmonitoring-request-line-selector-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetStopMonitoringFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.RequestStopArea(request)

	if len(response.MonitoredStopVisits) != 1 {
		t.Fatalf("Response.MonitoredStopVisits should be 1 is %v", len(response.MonitoredStopVisits))
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaTimeSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.SIRIPartner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	line := referential.model.Lines().New()
	obj := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:27:LOC")
	line.SetObjectID(obj)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:26:LOC")
	vehicleJourney.SetObjectID(obj)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	startTime, _ := time.Parse(time.RFC3339, "2016-09-22T07:54:52.977Z")

	stopVisit := referential.model.StopVisits().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:25:LOC")
	stopVisit.SetObjectID(obj)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	sVSchedule := model.StopVisitSchedule{}
	sVSchedule.SetArrivalTime(startTime.Add(2 * time.Minute))
	stopVisit.Schedules["actual"] = &sVSchedule
	stopVisit.Save()

	stopVisit2 := referential.model.StopVisits().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:5:LOC")
	stopVisit2.SetObjectID(obj)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	sVSchedule2 := model.StopVisitSchedule{}
	sVSchedule2.SetArrivalTime(startTime.Add(10 * time.Minute))
	stopVisit.Schedules["actual"] = &sVSchedule2
	stopVisit2.Save()

	file, err := os.Open("testdata/stopmonitoring-request-time-selector-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetStopMonitoringFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.RequestStopArea(request)

	if len(response.MonitoredStopVisits) != 1 {
		t.Fatalf("Response.MonitoredStopVisits should be 1 is %v", len(response.MonitoredStopVisits))
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaNotFound(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.SIRIPartner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	file, err := os.Open("testdata/stopmonitoring-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetStopMonitoringFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.RequestStopArea(request)

	if response.Address != "http://edwig" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://edwig", response.Address)
	}
	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.RequestMessageRef != "StopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
	if response.Status {
		t.Errorf("Response has wrong status:\n got: %v\n expected: false", response.Status)
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

func Test_SIRIStopMonitoringRequestBroadcaster_LogXMLGetStopMonitoring(t *testing.T) {
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
	request, err := siri.NewXMLGetStopMonitoringFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	logXMLStopMonitoringRequest(logStashEvent, &request.XMLStopMonitoringRequest)
	if logStashEvent["messageIdentifier"] != "StopMonitoring:Test:0" {
		t.Errorf("Wrong messageIdentifier logged:\n got: %v\n expected: StopMonitoring:Test:0", logStashEvent["messageIdentifier"])
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

	logSIRIStopMonitoringDelivery(logStashEvent, response.SIRIStopMonitoringDelivery)
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

func Test_SIRIStopMonitoringRequestBroadcaster_RemoteObjectIDKindPresent(t *testing.T) {
	partner := &Partner{}
	partner.Settings = make(map[string]string)

	partner.Settings["siri-stop-monitoring-request-broadcaster.remote_objectid_kind"] = "Kind1"
	partner.Settings["remote_objectid_kind"] = "Kind2"

	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)

	if connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_REQUEST_BROADCASTER) != "Kind1" {
		t.Errorf("RemoteObjectIDKind should be egals to Kind1")
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RemoteObjectIDKindAbsent(t *testing.T) {
	partner := &Partner{}
	partner.Settings = make(map[string]string)

	partner.Settings["siri-stop-monitoring-request-broadcaster.remote_objectid_kind"] = ""
	partner.Settings["remote_objectid_kind"] = "Kind2"

	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)

	if connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_REQUEST_BROADCASTER) != "Kind2" {
		t.Errorf("RemoteObjectIDKind should be egals to Kind2")
	}
}
