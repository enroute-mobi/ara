package core

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type fakeBroadcaster struct {
	Events []*model.StopAreaUpdateEvent
}

func (fb *fakeBroadcaster) FakeBroadcaster(event *model.StopAreaUpdateEvent) {
	fb.Events = append(fb.Events, event)
}

func prepare_SIRIStopMonitoringRequestCollector(t *testing.T, responseFilePath string) *model.StopAreaUpdateEvent {
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

	// Create a SIRIStopMonitoringRequestCollector
	partners := createTestPartnerManager()
	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"remote_url":           ts.URL,
			"remote_objectid_kind": "test kind",
		},
		manager: partners,
	}
	partners.Save(partner)

	// Create StopArea with ObjectId
	stopArea := partners.Model().StopAreas().New()
	objectid := model.NewObjectID("test kind", "test value")
	stopArea.SetObjectID(objectid)
	partners.Model().StopAreas().Save(&stopArea)

	siriStopMonitoringRequestCollector := NewSIRIStopMonitoringRequestCollector(partner)

	fs := fakeBroadcaster{}
	siriStopMonitoringRequestCollector.SetStopAreaUpdateSubscriber(fs.FakeBroadcaster)
	siriStopMonitoringRequestCollector.SetClock(model.NewFakeClock())
	stopAreaUpdateRequest := NewStopAreaUpdateRequest(stopArea.Id())
	siriStopMonitoringRequestCollector.RequestStopAreaUpdate(stopAreaUpdateRequest)

	time.Sleep(42 * time.Millisecond)
	return fs.Events[0]
}

func testStopMonitoringLogStash(t *testing.T) {
	events := audit.CurrentLogStash().(*audit.FakeLogStash).Events()
	if len(events) != 1 {
		t.Errorf("Logstash should have recieved an event, got: %v", events)
	}
	if len(events[0]) != 12 {
		t.Errorf("LogstashEvent should have 12 values, got: %v", events[0])
	}
}

// WIP
func Test_SIRIStopMonitoringRequestCollector_RequestStopAreaUpdate(t *testing.T) {
	stopAreaUpdateEvent := prepare_SIRIStopMonitoringRequestCollector(t, "testdata/stopmonitoring-response-soap.xml")
	if stopAreaUpdateEvent == nil {
		t.Error("RequestStopAreaUpdate should not return nil")
	}
	if len(stopAreaUpdateEvent.StopVisitUpdateEvents) != 2 {
		t.Errorf("RequestStopAreaUpdate should have 2 StopVisitUpdateEvents, got: %v", len(stopAreaUpdateEvent.StopVisitUpdateEvents))
	}
	stopVisitEvent := stopAreaUpdateEvent.StopVisitUpdateEvents[0]

	// Date is time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC) with fake clock
	if expected := model.FAKE_CLOCK_INITIAL_DATE; stopVisitEvent.Created_at != expected {
		t.Errorf("Wrong Created_At for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.Created_at)
	}
	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T07:56:53+02:00"); !stopVisitEvent.RecordedAt.Equal(expected) {
		t.Errorf("Wrong RecorderAt for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.RecordedAt)
	}
	if expected := model.STOP_VISIT_ARRIVAL_ARRIVED; stopVisitEvent.ArrivalStatuts != expected {
		t.Errorf("Wrong ArrivalStatuts for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.ArrivalStatuts)
	}
	if expected := model.STOP_VISIT_DEPARTURE_UNDEFINED; stopVisitEvent.DepartureStatus != expected {
		t.Errorf("Wrong DepartureStatuts for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.DepartureStatus)
	}
	if expected := "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3"; stopVisitEvent.StopVisitObjectid.Value() != expected {
		t.Errorf("Wrong ObjectID for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.StopVisitObjectid.Value())
	}
	// Aimed schedule
	schedule := stopVisitEvent.Schedules[model.STOP_VISIT_SCHEDULE_AIMED]
	if !schedule.DepartureTime().IsZero() {
		t.Errorf("AimedDepartureTime for stopVisitEvent should be zero, got: %v", schedule.DepartureTime())
	}
	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T07:54:00+02:00"); !schedule.ArrivalTime().Equal(expected) {
		t.Errorf("Wrong AimedArrivalTime for stopVisitEvent:\n expected: %v\n got: %v", expected, schedule.ArrivalTime())
	}
	// Expected schedule
	_, ok := stopVisitEvent.Schedules[model.STOP_VISIT_SCHEDULE_EXPECTED]
	if ok {
		t.Errorf("Expected Schedule shouldn't be created, got: %v", stopVisitEvent.Schedules[model.STOP_VISIT_SCHEDULE_EXPECTED])
	}
	// Actual schedule
	schedule = stopVisitEvent.Schedules[model.STOP_VISIT_SCHEDULE_ACTUAL]
	if !schedule.DepartureTime().IsZero() {
		t.Errorf("ActualDepartureTime for stopVisitEvent should be zero, got: %v", schedule.DepartureTime())
	}
	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T07:54:00+02:00"); !schedule.ArrivalTime().Equal(expected) {
		t.Errorf("Wrong ActualArrivalTime for stopVisitEvent:\n expected: %v\n got: %v", expected, schedule.ArrivalTime())
	}
}

// Test Factory Validate
func Test_SIRIStopMonitoringRequestCollectorFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		Settings:       make(map[string]string),
		ConnectorTypes: []string{"siri-stop-monitoring-request-collector"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have three errors when remote_url and remote_objectid_kind aren't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_url":           "remote_url",
		"remote_objectid_kind": "remote_objectid_kind",
		"remote_credential":    "remote_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when remote_url is set, got: %v", apiPartner.Errors)
	}
}

func Test_SIRIStopMonitoringRequestCollector_LogStopMonitoringRequest(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)
	time := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &siri.SIRIStopMonitoringRequest{
		MessageIdentifier: "0000-0000-0000-0000",
		MonitoringRef:     "test",
		RequestorRef:      "Edwig",
		RequestTimestamp:  time,
	}

	logSIRIStopMonitoringRequest(logStashEvent, request)
	if logStashEvent["Connector"] != "StopMonitoringRequestCollector" {
		t.Errorf("Wrong Connector logged:\n got: %v\n expected: StopMonitoringRequestCollector", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["messageIdentifier"] != "0000-0000-0000-0000" {
		t.Errorf("Wrong messageIdentifier logged:\n got: %v\n expected: 0000-0000-0000-0000", logStashEvent["messageIdentifier"])
	}
	if logStashEvent["requestorRef"] != "Edwig" {
		t.Errorf("Wrong requestorRef logged:\n got: %v\n expected: Edwig", logStashEvent["requestorRef"])
	}
	if logStashEvent["monitoringRef"] != "test" {
		t.Errorf("Wrong monitoringRef logged:\n got: %v\n expected: test", logStashEvent["monitoringRef"])
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

func Test_SIRIStopMonitoringRequestCollector_LogStopMonitoringResponse(t *testing.T) {
	logStashEvent := make(audit.LogStashEvent)

	file, err := os.Open("testdata/stopmonitoring-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	response, err := siri.NewXMLStopMonitoringResponseFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	logXMLStopMonitoringResponse(logStashEvent, response)

	if logStashEvent["address"] != "http://appli.chouette.mobi/siri_france/siri" {
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
	if logStashEvent["responseTimestamp"] != "2016-09-22 08:01:20.227 +0200 CEST" {
		t.Errorf("Wrong responseTimestamp logged:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", logStashEvent["responseTimestamp"])
	}
	if logStashEvent["responseXML"] != response.RawXML() {
		t.Errorf("Wrong responseXML logged:\n got: %v\n expected: %v", logStashEvent["responseXML"], response.RawXML())
	}
}
