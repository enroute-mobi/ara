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

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaNoSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"local_url":                              "http://ara",
		"remote_code_space":                   "codeSpace",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	code := model.NewCode("codeSpace", "modelOperatorRef")

	operator := referential.Model().Operators().New()
	operator.SetCode(code)

	operator.Save()

	code = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:24:LOC")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	stopVisit := referential.model.StopVisits().New()
	stopVisitRef := model.Reference{}
	obj1 := model.NewCode("codeSpace", "modelOperatorRef")
	stopVisitRef.Code = &obj1

	stopVisit.SetCode(obj1)
	stopVisit.References.Set("OperatorRef", stopVisitRef)
	stopVisit.StopAreaId = stopArea.Id()

	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:26:LOC")
	vehicleJourney.SetCode(obj)
	vehicleJourney.Save()

	stopVisit.VehicleJourneyId = vehicleJourney.Id()

	line := referential.model.Lines().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:27:LOC")
	line.SetCode(obj)
	line.Save()

	vehicleJourney.LineId = line.Id()

	stopVisit2 := referential.model.StopVisits().New()
	stopVisitRef2 := model.Reference{}
	obj2 := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:28:LOC")
	stopVisitRef2.Code = &obj2

	stopVisit2.SetCode(obj2)
	stopVisit2.References.Set("OperatorRef", stopVisitRef)
	stopVisit2.StopAreaId = stopArea.Id()

	stopVisit2.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit2.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	obj2 = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:29:LOC")
	vehicleJourney2.SetCode(obj2)
	vehicleJourney2.Save()

	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()

	line2 := referential.model.Lines().New()
	obj = model.NewCode("WrongCodeSpace", "NINOXE:StopPoint:SP:30:LOC")
	line2.SetCode(obj)
	line2.Save()

	vehicleJourney2.LineId = line2.Id()

	file, err := os.Open("testdata/stopmonitoring-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetStopMonitoringFromContent(content)
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
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"local_url":            "http://ara",
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	code := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:24:LOC")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	code2 := model.NewCode("wrongCodeSpace", "NINOXE:StopPoint:SP:20:LOC")
	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetCode(code2)
	stopArea2.ReferentId = stopArea.Id()
	stopArea2.Save()

	stopVisit := referential.model.StopVisits().New()
	stopVisitRef := model.Reference{}
	obj1 := model.NewCode("codeSpace", "modelOperatorRef")
	stopVisitRef.Code = &obj1

	stopVisit.SetCode(obj1)
	stopVisit.References.Set("OperatorRef", stopVisitRef)
	stopVisit.StopAreaId = stopArea.Id()

	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:26:LOC")
	vehicleJourney.SetCode(obj)
	vehicleJourney.Save()

	stopVisit.VehicleJourneyId = vehicleJourney.Id()

	line := referential.model.Lines().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:27:LOC")
	line.SetCode(obj)
	line.Save()

	vehicleJourney.LineId = line.Id()

	stopVisit2 := referential.model.StopVisits().New()
	obj2 := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:28:LOC")
	stopVisit2.SetCode(obj2)
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
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetStopMonitoringFromContent(content)
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
	if response.MonitoredStopVisits[0].StopPointRef != code.Value() || response.MonitoredStopVisits[1].StopPointRef != code.Value() {
		t.Errorf("Both MonitoredStopVisits should have the same StopPointRef, got %v and %v", response.MonitoredStopVisits[0].StopPointRef, response.MonitoredStopVisits[1].StopPointRef)
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RequestStopAreaLineSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"local_url":                              "http://ara",
		"remote_code_space":                   "codeSpace",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	stopArea := referential.Model().StopAreas().New()
	code := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:24:LOC")
	stopArea.SetCode(code)
	stopArea.Save()

	line := referential.model.Lines().New()
	obj := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:27:LOC")
	line.SetCode(obj)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:26:LOC")
	vehicleJourney.SetCode(obj)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:25:LOC")
	stopVisit.SetCode(obj)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.Save()

	line2 := referential.model.Lines().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:7:LOC")
	line2.SetCode(obj)
	line.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:6:LOC")
	vehicleJourney2.SetCode(obj)
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit2 := referential.model.StopVisits().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:5:LOC")
	stopVisit2.SetCode(obj)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit2.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(15*time.Minute))
	stopVisit2.Save()

	file, err := os.Open("testdata/stopmonitoring-request-line-selector-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetStopMonitoringFromContent(content)
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
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"local_url":                              "http://ara",
		"remote_code_space":                   "codeSpace",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	stopArea := referential.Model().StopAreas().New()
	code := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:24:LOC")
	stopArea.SetCode(code)
	stopArea.Save()

	line := referential.model.Lines().New()
	obj := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:27:LOC")
	line.SetCode(obj)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:26:LOC")
	vehicleJourney.SetCode(obj)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	startTime, _ := time.Parse(time.RFC3339, "2016-09-22T07:54:52.977Z")

	stopVisit := referential.model.StopVisits().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:25:LOC")
	stopVisit.SetCode(obj)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", startTime.Add(2*time.Minute))
	stopVisit.Save()

	stopVisit2 := referential.model.StopVisits().New()
	obj = model.NewCode("codeSpace", "NINOXE:StopPoint:SP:5:LOC")
	stopVisit2.SetCode(obj)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.Schedules.SetArrivalTime("actual", startTime.Add(10*time.Minute))
	stopVisit2.Save()

	file, err := os.Open("testdata/stopmonitoring-request-time-selector-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetStopMonitoringFromContent(content)
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
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"local_url":                              "http://ara",
		"remote_code_space":                   "codeSpace",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)

	connector.SetClock(clock.NewFakeClock())

	file, err := os.Open("testdata/stopmonitoring-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetStopMonitoringFromContent(content)
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
		t.Errorf("Response has wrong responseMessageIdentifier:\n got: %v\n expected: Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
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
	partner.PartnerSettings = s.NewEmptyPartnerSettings(partner.UUIDGenerator)
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have errors when local_credential and remote_code_space aren't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_code_space": "remote_code_space",
		"local_credential":     "local_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when local_credential and remote_code_space are set, got: %v", apiPartner.Errors)
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RemoteCodeSpacePresent(t *testing.T) {
	partner := NewPartner()

	settings := map[string]string{
		"siri-stop-monitoring-request-broadcaster.remote_code_space": "CodeSpace1",
		"remote_code_space": "CodeSpace2",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)

	if connector.partner.RemoteCodeSpace(SIRI_STOP_MONITORING_REQUEST_BROADCASTER) != "CodeSpace1" {
		t.Errorf("RemoteCodeSpace should be equal to CodeSpace1")
	}
}

func Test_SIRIStopMonitoringRequestBroadcaster_RemoteCodeSpaceAbsent(t *testing.T) {
	partner := NewPartner()

	settings := map[string]string{
		"siri-stop-monitoring-request-broadcaster.remote_code_space": "",
		"remote_code_space": "CodeSpace2",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIStopMonitoringRequestBroadcaster(partner)

	if connector.partner.RemoteCodeSpace(SIRI_STOP_MONITORING_REQUEST_BROADCASTER) != "CodeSpace2" {
		t.Errorf("RemoteCodeSpace should be equal to CodeSpace2")
	}
}
