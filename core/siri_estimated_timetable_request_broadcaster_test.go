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

func Test_SIRIEstimatedTimetableBroadcaster_RequestStopAreaNoSelector(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"
	connector := NewSIRIEstimatedTimetableBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(model.NewObjectID("objectidKind", "stopArea1"))
	stopArea.Monitored = true
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetObjectID(model.NewObjectID("objectidKind", "stopArea2"))
	stopArea2.Monitored = true
	stopArea2.Save()

	line := referential.model.Lines().New()
	line.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:2:LOC"))
	line.Name = "lineName"
	line.Save()

	line2 := referential.model.Lines().New()
	line2.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:3:LOC"))
	line2.Name = "lineName2"
	line2.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vehicleJourney.SetObjectID(model.NewObjectID("objectidKind", "vehicleJourney"))
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	vehicleJourney2.SetObjectID(model.NewObjectID("objectidKind", "vehicleJourney2"))
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	vehicleJourney3 := referential.model.VehicleJourneys().New()
	vehicleJourney3.SetObjectID(model.NewObjectID("objectidKind", "vehicleJourney3"))
	vehicleJourney3.LineId = line2.Id()
	vehicleJourney3.Save()

	stopVisit := referential.model.StopVisits().New()
	stopVisit.SetObjectID(model.NewObjectID("objectidKind", "stopVisit"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.PassageOrder = 1
	stopVisit.ArrivalStatus = "onTime"
	stopVisit.Schedules.SetArrivalTime("aimed", connector.Clock().Now().Add(1*time.Minute))
	stopVisit.Schedules.SetArrivalTime("expected", connector.Clock().Now().Add(1*time.Minute))
	stopVisit.Save()

	stopVisit2 := referential.model.StopVisits().New()
	stopVisit2.SetObjectID(model.NewObjectID("objectidKind", "stopVisit2"))
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.StopAreaId = stopArea2.Id()
	stopVisit2.PassageOrder = 2
	stopVisit2.ArrivalStatus = "onTime"
	stopVisit2.Schedules.SetArrivalTime("aimed", connector.Clock().Now().Add(2*time.Minute))
	stopVisit2.Schedules.SetArrivalTime("expected", connector.Clock().Now().Add(2*time.Minute))
	stopVisit2.Save()

	stopVisit3 := referential.model.StopVisits().New()
	stopVisit3.SetObjectID(model.NewObjectID("objectidKind", "stopVisit3"))
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.StopAreaId = stopArea.Id()
	stopVisit3.PassageOrder = 1
	stopVisit3.ArrivalStatus = "onTime"
	stopVisit3.Schedules.SetArrivalTime("aimed", connector.Clock().Now().Add(1*time.Minute))
	stopVisit3.Schedules.SetArrivalTime("expected", connector.Clock().Now().Add(1*time.Minute))
	stopVisit3.Save()

	stopVisit4 := referential.model.StopVisits().New()
	stopVisit4.SetObjectID(model.NewObjectID("objectidKind", "stopVisit4"))
	stopVisit4.VehicleJourneyId = vehicleJourney3.Id()
	stopVisit4.StopAreaId = stopArea.Id()
	stopVisit4.PassageOrder = 1
	stopVisit4.ArrivalStatus = "onTime"
	stopVisit4.Schedules.SetArrivalTime("aimed", connector.Clock().Now().Add(1*time.Minute))
	stopVisit4.Schedules.SetArrivalTime("expected", connector.Clock().Now().Add(1*time.Minute))
	stopVisit4.Save()

	file, err := os.Open("testdata/estimated_timetable_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetEstimatedTimetableFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.RequestLine(request)

	if response.Address != "http://edwig" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://edwig", response.Address)
	}
	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.RequestMessageRef != "EstimatedTimetable:Test:0" {
		t.Errorf("Response has wrong requestMessageRef:\n got: %v\n expected: StopMonitoring:Test:0", response.RequestMessageRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC", response.ResponseMessageIdentifier)
	}
	if !response.ResponseTimestamp.Equal(connector.Clock().Now()) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}

	if len(response.EstimatedJourneyVersionFrames) != 2 {
		t.Fatalf("Response should have 2 EstimatedJourneyVersionFrames, got: %v", len(response.EstimatedJourneyVersionFrames))
	}
	// Test second Line because VehicleJourneys order is random
	if len(response.EstimatedJourneyVersionFrames[1].EstimatedVehicleJourneys) != 2 {
		t.Fatalf("Second EstimatedJourneyVersionFrame should have 2 EstimatedVehicleJourneys, got: %v", len(response.EstimatedJourneyVersionFrames[1].EstimatedVehicleJourneys))
	}

	firstLine := response.EstimatedJourneyVersionFrames[0]
	if firstLine.RecordedAtTime != connector.Clock().Now() {
		t.Errorf("Wrong EstimatedJourneyVersionFrames for first EstimatedJourneyVersionFrame:\n got: %v\n want: %v", firstLine.RecordedAtTime, connector.Clock().Now())
	}
	if len(firstLine.EstimatedVehicleJourneys) != 1 {
		t.Fatalf("First EstimatedJourneyVersionFrame should have 1 EstimatedVehicleJourneys, got: %v", len(firstLine.EstimatedVehicleJourneys))
	}
	firstVJ := firstLine.EstimatedVehicleJourneys[0]
	if firstVJ.LineRef != "NINOXE:Line:2:LOC" {
		t.Errorf("Wrong LineRef for first EstimatedVehicleJourney:\n got: %v\n want: NINOXE:Line:2:LOC", firstVJ.LineRef)
	}
	if firstVJ.DatedVehicleJourneyRef != "vehicleJourney" {
		t.Errorf("Wrong DatedVehicleJourneyRef for first EstimatedVehicleJourney:\n got: %v\n want: vehicleJourney", firstVJ.DatedVehicleJourneyRef)
	}
	if len(firstVJ.EstimatedCalls) != 2 {
		t.Fatalf("First EstimatedVehicleJourney should have 2 EstimatedCalls, got: %v", len(firstVJ.EstimatedCalls))
	}
	firstEC := firstVJ.EstimatedCalls[0]
	if firstEC.StopPointRef != "stopArea1" {
		t.Errorf("Wrong StopPointRef for first EstimatedCall:\n got: %v\n want: stopArea1", firstEC.StopPointRef)
	}
	if firstEC.Order != 1 {
		t.Errorf("Wrong Order for first EstimatedCall:\n got: %v\n want: 1", firstEC.Order)
	}
	if firstEC.ArrivalStatus != "onTime" {
		t.Errorf("Wrong ArrivalStatus for first EstimatedCall:\n got: %v\n want: onTime", firstEC.ArrivalStatus)
	}
	if !firstEC.AimedArrivalTime.Equal(connector.Clock().Now().Add(1 * time.Minute)) {
		t.Errorf("Wrong AimedArrivalTime for first EstimatedCall:\n got: %v\n want: %v", firstEC.AimedArrivalTime, connector.Clock().Now())
	}
	if !firstEC.ExpectedArrivalTime.Equal(connector.Clock().Now().Add(1 * time.Minute)) {
		t.Errorf("Wrong ExpectedArrivalTime for first EstimatedCall:\n got: %v\n want: %v", firstEC.ExpectedArrivalTime, connector.Clock().Now())
	}
}

func Test_SIRIEstimatedTimetableBroadcaster_RequestStopAreaWithReferent(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"
	connector := NewSIRIEstimatedTimetableBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(model.NewObjectID("objectidKind", "stopArea1"))
	stopArea.Monitored = true
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.SetObjectID(model.NewObjectID("wrongObjectidKind", "stopArea2"))
	stopArea2.ReferentId = stopArea.Id()
	stopArea2.Monitored = true
	stopArea2.Save()

	line := referential.model.Lines().New()
	line.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:2:LOC"))
	line.Name = "lineName"
	line.Save()

	line2 := referential.model.Lines().New()
	line2.SetObjectID(model.NewObjectID("objectidKind", "NINOXE:Line:3:LOC"))
	line2.Name = "lineName2"
	line2.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vehicleJourney.SetObjectID(model.NewObjectID("objectidKind", "vehicleJourney"))
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	vehicleJourney2.SetObjectID(model.NewObjectID("objectidKind", "vehicleJourney2"))
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit := referential.model.StopVisits().New()
	stopVisit.SetObjectID(model.NewObjectID("objectidKind", "stopVisit"))
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.PassageOrder = 1
	stopVisit.ArrivalStatus = "onTime"
	stopVisit.Schedules.SetArrivalTime("aimed", connector.Clock().Now().Add(1*time.Minute))
	stopVisit.Schedules.SetArrivalTime("expected", connector.Clock().Now().Add(1*time.Minute))
	stopVisit.Save()

	stopVisit3 := referential.model.StopVisits().New()
	stopVisit3.SetObjectID(model.NewObjectID("objectidKind", "stopVisit3"))
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.StopAreaId = stopArea2.Id()
	stopVisit3.PassageOrder = 1
	stopVisit3.ArrivalStatus = "onTime"
	stopVisit3.Schedules.SetArrivalTime("aimed", connector.Clock().Now().Add(1*time.Minute))
	stopVisit3.Schedules.SetArrivalTime("expected", connector.Clock().Now().Add(1*time.Minute))
	stopVisit3.Save()

	file, err := os.Open("testdata/estimated_timetable_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetEstimatedTimetableFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response := connector.RequestLine(request)

	if len(response.EstimatedJourneyVersionFrames) != 2 {
		t.Fatalf("Response should have 2 EstimatedJourneyVersionFrames, got: %v", len(response.EstimatedJourneyVersionFrames))
	}

	firstLine := response.EstimatedJourneyVersionFrames[0]
	if len(firstLine.EstimatedVehicleJourneys) != 1 {
		t.Fatalf("First EstimatedJourneyVersionFrame should have 1 EstimatedVehicleJourneys, got: %v", len(firstLine.EstimatedVehicleJourneys))
	}
	firstVJ := firstLine.EstimatedVehicleJourneys[0]
	if firstVJ.LineRef != "NINOXE:Line:2:LOC" {
		t.Errorf("Wrong LineRef for first EstimatedVehicleJourney:\n got: %v\n want: NINOXE:Line:2:LOC", firstVJ.LineRef)
	}
	if firstVJ.DatedVehicleJourneyRef != "vehicleJourney" {
		t.Errorf("Wrong DatedVehicleJourneyRef for first EstimatedVehicleJourney:\n got: %v\n want: vehicleJourney", firstVJ.DatedVehicleJourneyRef)
	}
	if len(firstVJ.EstimatedCalls) != 1 {
		t.Fatalf("First EstimatedVehicleJourney should have 1 EstimatedCalls, got: %v", len(firstVJ.EstimatedCalls))
	}
	firstEC := firstVJ.EstimatedCalls[0]
	if firstEC.StopPointRef != "stopArea1" {
		t.Errorf("Wrong StopPointRef for first EstimatedCall:\n got: %v\n want: stopArea1", firstEC.StopPointRef)
	}

	secondLine := response.EstimatedJourneyVersionFrames[1]
	if len(secondLine.EstimatedVehicleJourneys) != 1 {
		t.Fatalf("Second EstimatedJourneyVersionFrame should have 1 EstimatedVehicleJourneys, got: %v", len(secondLine.EstimatedVehicleJourneys))
	}
	secondVJ := secondLine.EstimatedVehicleJourneys[0]
	if secondVJ.LineRef != "NINOXE:Line:3:LOC" {
		t.Errorf("Wrong LineRef for second EstimatedVehicleJourney:\n got: %v\n want: NINOXE:Line:3:LOC", secondVJ.LineRef)
	}
	if secondVJ.DatedVehicleJourneyRef != "vehicleJourney2" {
		t.Errorf("Wrong DatedVehicleJourneyRef for second EstimatedVehicleJourney:\n got: %v\n want: vehicleJourney", secondVJ.DatedVehicleJourneyRef)
	}
	if len(secondVJ.EstimatedCalls) != 1 {
		t.Fatalf("Second EstimatedVehicleJourney should have 1 EstimatedCalls, got: %v", len(secondVJ.EstimatedCalls))
	}
	secondEC := secondVJ.EstimatedCalls[0]
	if secondEC.StopPointRef != "stopArea1" {
		t.Errorf("Wrong StopPointRef for second EstimatedCall:\n got: %v\n want: stopArea1", secondEC.StopPointRef)
	}
}

func Test_SIRIEstimatedTimetableBroadcasterFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		Settings:       make(map[string]string),
		ConnectorTypes: []string{"siri-estimated-timetable-request-broadcaster"},
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

func Test_SIRIEstimatedTimetableBroadcaster_LogXMLGetStopMonitoring(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	file, err := os.Open("testdata/estimated_timetable_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetEstimatedTimetableFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	logXMLEstimatedTimetableRequest(logStashEvent, &request.XMLEstimatedTimetableRequest)
	if logStashEvent["messageIdentifier"] != "EstimatedTimetable:Test:0" {
		t.Errorf("Wrong messageIdentifier logged:\n got: %v\n expected: EstimatedTimetable:Test:0", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["requestTimestamp"] != "2016-09-07 09:11:25.174 +0000 UTC" {
		t.Errorf("Wrong requestTimestamp logged:\n got: %v\n expected: 2016-09-07 09:11:25.174 +0000 UTC", logStashEvent["requestTimestamp"])
	}
	if logStashEvent["requestXML"] != request.RawXML() {
		t.Errorf("Wrong requestXML logged:\n got: %v\n expected: %v", logStashEvent["requestXML"], request.RawXML())
	}
}

func Test_SIRIEstimatedTimetableBroadcaster_LogSIRIStopMonitoringResponse(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	response := &siri.SIRIEstimatedTimeTableResponse{
		Address:                   "edwig.edwig",
		ProducerRef:               "NINOXE:default",
		ResponseMessageIdentifier: "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26",
	}
	response.RequestMessageRef = "EstimatedTimetable:Test:0"
	response.ResponseTimestamp = time
	response.Status = true

	logSIRIEstimatedTimetableDelivery(logStashEvent, response.SIRIEstimatedTimetableDelivery, nil, nil)
	logSIRIEstimatedTimetableResponse(logStashEvent, response)

	if logStashEvent["address"] != "edwig.edwig" {
		t.Errorf("Wrong address logged:\n got: %v\n expected: edwig.edwig", logStashEvent["address"])
	}
	if logStashEvent["producerRef"] != "NINOXE:default" {
		t.Errorf("Wrong producerRef logged:\n got: %v\n expected: NINOXE:default", logStashEvent["producerRef"])
	}
	if logStashEvent["requestMessageRef"] != "EstimatedTimetable:Test:0" {
		t.Errorf("Wrong requestMessageRef logged:\n got: %v\n expected: EstimatedTimetable:Test:0", logStashEvent["requestMessageRef"])
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

func Test_SIRIEstimatedTimetableBroadcaster_RemoteObjectIDKindPresent(t *testing.T) {
	partner := &Partner{}
	partner.Settings = make(map[string]string)

	partner.Settings["siri-estimated-timetable-request-broadcaster.remote_objectid_kind"] = "Kind1"
	partner.Settings["remote_objectid_kind"] = "Kind2"

	connector := NewSIRIEstimatedTimetableBroadcaster(partner)

	if connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER) != "Kind1" {
		t.Errorf("RemoteObjectIDKind should be egals to Kind1")
	}
}

func Test_SIRIEstimatedTimetableBroadcaster_RemoteObjectIDKindAbsent(t *testing.T) {
	partner := &Partner{}
	partner.Settings = make(map[string]string)

	partner.Settings["siri-estimated-timetable-request-broadcaster.remote_objectid_kind"] = ""
	partner.Settings["remote_objectid_kind"] = "Kind2"

	connector := NewSIRIEstimatedTimetableBroadcaster(partner)

	if connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER) != "Kind2" {
		t.Errorf("RemoteObjectIDKind should be egals to Kind2")
	}
}
