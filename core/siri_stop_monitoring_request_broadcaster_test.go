package core

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	ps "bitbucket.org/enroute-mobi/ara/core/partner_settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaNoSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("local_url", "http://ara")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	partner.SetSetting("generators.response_message_identifier", "Ara:ResponseMessage::%{uuid}:LOC")
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

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
	obj1 := model.NewObjectID("objectidKind", "modelOperatorRef")
	stopVisitRef.ObjectId = &obj1

	stopVisit.SetObjectID(obj1)
	stopVisit.References.Set("OperatorRef", stopVisitRef)
	stopVisit.StopAreaId = stopArea.Id()

	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:26:LOC")
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

	stopVisit2.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
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

	response := connector.RequestStopArea(request, &audit.BigQueryMessage{})

	if response.Address != "http://ara" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://ara", response.Address)
	}
	if response.ProducerRef != "Ara" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Ara", response.ProducerRef)
	}
	if response.RequestMessageRef != "StopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}

	if len(response.MonitoredStopVisits) != 1 {
		t.Fatalf("Response.MonitoredStopVisits should be 1 is %v", len(response.MonitoredStopVisits))
	}

	if response.MonitoredStopVisits[0].References["StopVisitReferences"]["OperatorRef"] != "modelOperatorRef" {
		t.Errorf("OperatorRef should be modelOperatorRef, got: %v", response.MonitoredStopVisits[0].References["StopVisitReferences"]["OperatorRef"])
	}

}

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopWithReferent(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("local_url", "http://ara")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	objectid := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	objectid2 := model.NewObjectID("wrongObjectidKind", "NINOXE:StopPoint:SP:20:LOC")
	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetObjectID(objectid2)
	stopArea2.ReferentId = stopArea.Id()
	stopArea2.Save()

	stopVisit := referential.model.StopVisits().New()
	stopVisitRef := model.Reference{}
	obj1 := model.NewObjectID("objectidKind", "modelOperatorRef")
	stopVisitRef.ObjectId = &obj1

	stopVisit.SetObjectID(obj1)
	stopVisit.References.Set("OperatorRef", stopVisitRef)
	stopVisit.StopAreaId = stopArea.Id()

	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:26:LOC")
	vehicleJourney.SetObjectID(obj)
	vehicleJourney.Save()

	stopVisit.VehicleJourneyId = vehicleJourney.Id()

	line := referential.model.Lines().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:27:LOC")
	line.SetObjectID(obj)
	line.Save()

	vehicleJourney.LineId = line.Id()

	stopVisit2 := referential.model.StopVisits().New()
	obj2 := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:28:LOC")
	stopVisit2.SetObjectID(obj2)
	stopVisit2.References.Set("OperatorRef", stopVisitRef)
	stopVisit2.StopAreaId = stopArea2.Id()

	stopVisit2.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit2.Save()

	stopVisit2.VehicleJourneyId = vehicleJourney.Id()

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

	response := connector.RequestStopArea(request, &audit.BigQueryMessage{})

	if response.Address != "http://ara" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://ara", response.Address)
	}
	if response.ProducerRef != "Ara" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Ara", response.ProducerRef)
	}
	if response.RequestMessageRef != "StopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}

	if len(response.MonitoredStopVisits) != 2 {
		t.Fatalf("Response.MonitoredStopVisits should be 2 is %v", len(response.MonitoredStopVisits))
	}
	if response.MonitoredStopVisits[0].StopPointRef != objectid.Value() || response.MonitoredStopVisits[1].StopPointRef != objectid.Value() {
		t.Errorf("Both MonitoredStopVisits should have the same StopPointRef, got %v and %v", response.MonitoredStopVisits[0].StopPointRef, response.MonitoredStopVisits[1].StopPointRef)
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaLineSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("local_url", "http://ara")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	partner.SetSetting("generators.response_message_identifier", "Ara:ResponseMessage::%{uuid}:LOC")
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

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
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
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
	stopVisit2.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(15*time.Minute))
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

	response := connector.RequestStopArea(request, &audit.BigQueryMessage{})

	if len(response.MonitoredStopVisits) != 1 {
		t.Fatalf("Response.MonitoredStopVisits should be 1 is %v", len(response.MonitoredStopVisits))
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaTimeSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("local_url", "http://ara")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	partner.SetSetting("generators.response_message_identifier", "Ara:ResponseMessage::%{uuid}:LOC")
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

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
	stopVisit.Schedules.SetArrivalTime("actual", startTime.Add(2*time.Minute))
	stopVisit.Save()

	stopVisit2 := referential.model.StopVisits().New()
	obj = model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:5:LOC")
	stopVisit2.SetObjectID(obj)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.Schedules.SetArrivalTime("actual", startTime.Add(10*time.Minute))
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

	response := connector.RequestStopArea(request, &audit.BigQueryMessage{})

	if len(response.MonitoredStopVisits) != 1 {
		t.Fatalf("Response.MonitoredStopVisits should be 1 is %v", len(response.MonitoredStopVisits))
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaNotFound(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("local_url", "http://ara")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	partner.SetSetting("generators.response_message_identifier", "Ara:ResponseMessage::%{uuid}:LOC")
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

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

	response := connector.RequestStopArea(request, &audit.BigQueryMessage{})

	if response.Address != "http://ara" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://ara", response.Address)
	}
	if response.ProducerRef != "Ara" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Ara", response.ProducerRef)
	}
	if response.RequestMessageRef != "StopMonitoring:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
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
		ConnectorTypes: []string{"siri-stop-monitoring-request-broadcaster"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = ps.NewPartnerSettings(partner.UUIDGenerator)
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
		Address:                   "ara.ara",
		ProducerRef:               "NINOXE:default",
		ResponseMessageIdentifier: "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26",
	}
	response.RequestMessageRef = "StopMonitoring:Test:0"
	response.Status = true
	response.ResponseTimestamp = time

	logSIRIStopMonitoringDelivery(logStashEvent, response.SIRIStopMonitoringDelivery)
	logSIRIStopMonitoringResponse(logStashEvent, response)

	if logStashEvent["address"] != "ara.ara" {
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
	partner := NewPartner()

	partner.SetSetting("siri-stop-monitoring-request-broadcaster.remote_objectid_kind", "Kind1")
	partner.SetSetting("remote_objectid_kind", "Kind2")

	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)

	if connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_REQUEST_BROADCASTER) != "Kind1" {
		t.Errorf("RemoteObjectIDKind should be egals to Kind1")
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RemoteObjectIDKindAbsent(t *testing.T) {
	partner := NewPartner()

	partner.SetSetting("siri-stop-monitoring-request-broadcaster.remote_objectid_kind", "")
	partner.SetSetting("remote_objectid_kind", "Kind2")

	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)

	if connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_REQUEST_BROADCASTER) != "Kind2" {
		t.Errorf("RemoteObjectIDKind should be egals to Kind2")
	}
}
