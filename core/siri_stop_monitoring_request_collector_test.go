package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/model"
	"github.com/jonboulle/clockwork"
)

func prepare_SIRIStopMonitoringRequestCollector(t *testing.T, responseFilePath string) *model.StopAreaUpdateEvent {
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
	partners := NewPartnerManager(model.NewMemoryModel())
	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"remote_url":           ts.URL,
			"remote_objectid_kind": "test kind",
		},
	}
	partners.Save(partner)

	// Create StopArea with ObjectId
	stopArea := partners.Model().StopAreas().New()
	objectid := model.NewObjectID("test kind", "test value")
	stopArea.SetObjectID(objectid)
	partners.Model().StopAreas().Save(&stopArea)

	siriStopMonitoringRequestCollector := NewSIRIStopMonitoringRequestCollector(partner)
	siriStopMonitoringRequestCollector.SetClock(clockwork.NewFakeClock())
	stopAreaUpdateRequest := NewStopAreaUpdateRequest(stopArea.Id())
	stopAreaUpdateEvent, err := siriStopMonitoringRequestCollector.RequestStopAreaUpdate(stopAreaUpdateRequest)
	if err != nil {
		t.Fatal(err)
	}

	return stopAreaUpdateEvent
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
	if expected := time.Date(1984, time.April, 4, 0, 0, 0, 0, time.UTC); stopVisitEvent.Created_at != expected {
		t.Errorf("Wrong Created_At for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.Created_at)
	}
	if expected := model.STOP_VISIT_ARRIVAL_ARRIVED; stopVisitEvent.ArrivalStatuts != expected {
		t.Errorf("Wrong ArrivalStatuts for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.ArrivalStatuts)
	}
	if expected := model.STOP_VISIT_DEPARTURE_UNDEFINED; stopVisitEvent.DepartureStatus != expected {
		t.Errorf("Wrong DepartureStatuts for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.DepartureStatus)
	}
	if expected := "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3"; stopVisitEvent.Stop_visit_objectid.Value() != expected {
		t.Errorf("Wrong ObjectID for stopVisitEvent:\n expected: %v\n got: %v", expected, stopVisitEvent.Stop_visit_objectid.Value())
	}
	// Aimed schedule
	schedule := stopVisitEvent.Schedules[model.STOP_VISIT_SCHEDULE_AIMED]
	if !schedule.DepartureTime().IsZero() {
		t.Errorf("AimedDepartureTime for stopVisitEvent should be zero, got: %v", schedule.DepartureTime())
	}
	if expected := time.Date(2016, time.September, 22, 5, 54, 0, 0, time.UTC); !schedule.ArrivalTime().Equal(expected) {
		t.Errorf("Wrong AimedArrivalTime for stopVisitEvent:\n expected: %v\n got: %v", expected, schedule.ArrivalTime())
	}
	// Expected schedule
	schedule = stopVisitEvent.Schedules[model.STOP_VISIT_SCHEDULE_EXPECTED]
	if !schedule.DepartureTime().IsZero() {
		t.Errorf("ExpectedDepartureTime for stopVisitEvent should be zero, got: %v", schedule.DepartureTime())
	}
	if !schedule.ArrivalTime().IsZero() {
		t.Errorf("ExpectedArrivalTime for stopVisitEvent should be zero, got: %v", schedule.ArrivalTime())
	}
	// Actual schedule
	schedule = stopVisitEvent.Schedules[model.STOP_VISIT_SCHEDULE_ACTUAL]
	if !schedule.DepartureTime().IsZero() {
		t.Errorf("ActualDepartureTime for stopVisitEvent should be zero, got: %v", schedule.DepartureTime())
	}
	if expected := time.Date(2016, time.September, 22, 5, 54, 0, 0, time.UTC); !schedule.ArrivalTime().Equal(expected) {
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
