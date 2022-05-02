package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	ps "bitbucket.org/enroute-mobi/ara/core/psettings"
	"bitbucket.org/enroute-mobi/ara/model"
)

type fakeBroadcaster struct {
	Events []model.UpdateEvent
}

func (fb *fakeBroadcaster) FakeBroadcaster(event model.UpdateEvent) {
	fb.Events = append(fb.Events, event)
}

func prepare_SIRIStopMonitoringRequestCollector(t *testing.T, responseFilePath string) []model.UpdateEvent {
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
	partner := partners.New("slug")
	partner.SetSettingsDefinition(map[string]string{
		"remote_url":           ts.URL,
		"remote_objectid_kind": "test kind",
	})
	partners.Save(partner)

	// Create StopArea with ObjectId
	stopArea := partners.Model().StopAreas().New()
	objectid := model.NewObjectID("test kind", "test value")
	stopArea.SetObjectID(objectid)
	partners.Model().StopAreas().Save(&stopArea)

	siriStopMonitoringRequestCollector := NewSIRIStopMonitoringRequestCollector(partner)

	fs := fakeBroadcaster{}
	siriStopMonitoringRequestCollector.SetUpdateSubscriber(fs.FakeBroadcaster)
	siriStopMonitoringRequestCollector.SetClock(clock.NewFakeClock())
	stopAreaUpdateRequest := NewStopAreaUpdateRequest(stopArea.Id())
	siriStopMonitoringRequestCollector.RequestStopAreaUpdate(stopAreaUpdateRequest)

	time.Sleep(42 * time.Millisecond)

	return fs.Events
}

func Test_SIRIStopMonitoringRequestCollector_RequestStopAreaUpdate(t *testing.T) {
	updateEvents := prepare_SIRIStopMonitoringRequestCollector(t, "testdata/stopmonitoring-response-soap.xml")

	// 2 stops 1 Line 2 VehicleJourneys 2 StopVisits
	if len(updateEvents) != 7 {
		t.Fatalf("Should have 7 update events, got %v", len(updateEvents))
	}

	stopVisitEvent := findSVEvent(updateEvents, "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3")
	if stopVisitEvent == nil {
		t.Fatal("Cannot find StopVisit event")
	}

	// Date is time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC) with fake clock
	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T07:56:53+02:00"); !stopVisitEvent.RecordedAt.Equal(expected) {
		t.Errorf("Wrong RecorderAt for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.RecordedAt)
	}
	if expected := model.STOP_VISIT_ARRIVAL_ARRIVED; stopVisitEvent.ArrivalStatus != expected {
		t.Errorf("Wrong ArrivalStatuts for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.ArrivalStatus)
	}
	if expected := model.STOP_VISIT_DEPARTURE_UNDEFINED; stopVisitEvent.DepartureStatus != expected {
		t.Errorf("Wrong DepartureStatuts for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.DepartureStatus)
	}
	if expected := "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3"; stopVisitEvent.ObjectId.Value() != expected {
		t.Errorf("Wrong ObjectID for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.ObjectId.Value())
	}
	// Aimed schedule
	schedule := stopVisitEvent.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED)
	if !schedule.DepartureTime().IsZero() {
		t.Errorf("AimedDepartureTime for stopVisitEvent should be zero, got: %v", schedule.DepartureTime())
	}
	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T07:54:00+02:00"); !schedule.ArrivalTime().Equal(expected) {
		t.Errorf("Wrong AimedArrivalTime for stopVisitEvent:\n expected: %v\n got: %v", expected, schedule.ArrivalTime())
	}
	// Expected schedule
	schedule = stopVisitEvent.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED)
	if !schedule.DepartureTime().IsZero() || !schedule.ArrivalTime().IsZero() {
		t.Errorf("Expected Schedule shouldn't be created, got: %v", stopVisitEvent.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED))
	}
	// Actual schedule
	schedule = stopVisitEvent.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL)
	if !schedule.DepartureTime().IsZero() {
		t.Errorf("ActualDepartureTime for stopVisitEvent should be zero, got: %v", schedule.DepartureTime())
	}
	if expected, _ := time.Parse(time.RFC3339, "2016-09-22T07:54:00+02:00"); !schedule.ArrivalTime().Equal(expected) {
		t.Errorf("Wrong ActualArrivalTime for stopVisitEvent:\n expected: %v\n got: %v", expected, schedule.ArrivalTime())
	}
}

func findSVEvent(events []model.UpdateEvent, ref string) *model.StopVisitUpdateEvent {
	for _, e := range events {
		svEvent, ok := e.(*model.StopVisitUpdateEvent)
		if !ok {
			continue
		}
		if svEvent.ObjectId.Value() == ref {
			return svEvent
		}
	}
	return nil
}

func Test_SIRIStopMonitoringRequestCollector_RequestStopAreaUpdate_MultipleDeliveries(t *testing.T) {
	updateEvents := prepare_SIRIStopMonitoringRequestCollector(t, "testdata/stopmonitoring-response-double-delivery-soap.xml")
	// 2 StopAreas 1 Line 2 VehicleJourneys 2 StopVisits
	if len(updateEvents) != 7 {
		t.Errorf("Should have 7 update events, got %v", len(updateEvents))
	}
}

// Test Factory Validate
func Test_SIRIStopMonitoringRequestCollectorFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		ConnectorTypes: []string{"siri-stop-monitoring-request-collector"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = ps.NewPartnerSettings(partner.UUIDGenerator)
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
