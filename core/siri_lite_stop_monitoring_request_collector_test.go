package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func prepare_SIRILiteStopMonitoringRequestCollector(t *testing.T, responseFilePath string) []model.UpdateEvent {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
	partners.Model().StopAreas().Save(stopArea)

	siriLiteStopMonitoringRequestCollector := NewSIRILiteStopMonitoringRequestCollector(partner)

	fs := fakeBroadcaster{}
	siriLiteStopMonitoringRequestCollector.SetUpdateSubscriber(fs.FakeBroadcaster)
	siriLiteStopMonitoringRequestCollector.SetClock(clock.NewFakeClock())
	stopAreaUpdateRequest := NewStopAreaUpdateRequest(stopArea.Id())
	siriLiteStopMonitoringRequestCollector.RequestStopAreaUpdate(stopAreaUpdateRequest)

	time.Sleep(42 * time.Millisecond)

	return fs.Events
}

func Test_SIRILiteStopMonitoringRequestCollector_RequestStopAreaUpdate(t *testing.T) {
	assert := assert.New(t)

	// Get update events
	updateEvents := prepare_SIRILiteStopMonitoringRequestCollector(t, "testdata/stopmonitoring-lite-delivery.json")
	assert.NotEmpty(updateEvents)
	assert.Equal(4, len(updateEvents), "4 update events")

	var eventKinds []model.EventKind
	for i := range updateEvents {
		eventKinds = append(eventKinds, updateEvents[i].EventKind())
	}
	// 1 StopArea, 1 Line, 1 VehicleJourney, 1 StopVisit
	assert.ElementsMatch([]model.EventKind{0, 2, 3, 4}, eventKinds)

	stopVisitEvent := findSVEvent(updateEvents, "SNCF_ACCES_CLOUD:Item::41178_133528:LOC")
	assert.NotNil(stopVisitEvent)

	expectedRecordedAt, _ := time.Parse(time.RFC3339, "2023-06-02T01:07:19.892Z")
	assert.Equal(expectedRecordedAt, stopVisitEvent.RecordedAt)

	// Arrival & departure status
	assert.Equal(model.STOP_VISIT_ARRIVAL_ONTIME, stopVisitEvent.ArrivalStatus)
	assert.Equal(model.STOP_VISIT_DEPARTURE_ONTIME, stopVisitEvent.DepartureStatus)

	// No actual Schedules
	actualSchedule := stopVisitEvent.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_ACTUAL)
	assert.Zero(actualSchedule.ArrivalTime())
	assert.Zero(actualSchedule.DepartureTime())

	// Existing expected schedules
	expectedSchedule := stopVisitEvent.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_EXPECTED)
	expectedDeparture, _ := time.Parse(time.RFC3339, "2023-06-02T08:47:40.000Z")
	assert.Equal(expectedDeparture, expectedSchedule.DepartureTime(), "expected Departure")

	expectedArrival, _ := time.Parse(time.RFC3339, "2023-06-02T08:46:40.000Z")
	assert.Equal(expectedArrival, expectedSchedule.ArrivalTime(), "expected Arrival")

	// Existing aimed schedules
	aimedSchedule := stopVisitEvent.Schedules.Schedule(model.STOP_VISIT_SCHEDULE_AIMED)
	aimedDeparture, _ := time.Parse(time.RFC3339, "2023-06-02T08:47:40.000Z")
	assert.Equal(aimedDeparture, aimedSchedule.DepartureTime())

	aimedArrival, _ := time.Parse(time.RFC3339, "2023-06-02T08:46:40.000Z")
	assert.Equal(aimedArrival, aimedSchedule.ArrivalTime())

}
