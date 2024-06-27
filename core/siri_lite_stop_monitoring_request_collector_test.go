package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
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

	settings := map[string]string{
		"remote_url":        ts.URL,
		"remote_code_space": "test kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	// Create StopArea with Code
	stopArea := partners.Model().StopAreas().New()
	code := model.NewCode("test kind", "test value")
	stopArea.SetCode(code)
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

	var eventCodeSpaces []model.EventKind
	for i := range updateEvents {
		eventCodeSpaces = append(eventCodeSpaces, updateEvents[i].EventKind())
	}

	expectedEventKinds := []model.EventKind{
		model.STOP_AREA_EVENT,
		model.LINE_EVENT,
		model.VEHICLE_JOURNEY_EVENT,
		model.STOP_VISIT_EVENT}

	assert.ElementsMatch(expectedEventKinds, eventCodeSpaces,
		"Should have 1 StopArea, 1 Line, 1 VehicleJourney, 1 StopVisit")

	stopVisitEvent := findSVEvent(updateEvents, "SNCF_ACCES_CLOUD:Item::41178_133528:LOC")
	assert.NotNil(stopVisitEvent)

	expectedRecordedAt, _ := time.Parse(time.RFC3339, "2023-06-02T01:07:19.892Z")
	assert.Equal(expectedRecordedAt, stopVisitEvent.RecordedAt)

	// Arrival & departure status
	assert.Equal(model.STOP_VISIT_ARRIVAL_ONTIME, stopVisitEvent.ArrivalStatus)
	assert.Equal(model.STOP_VISIT_DEPARTURE_ONTIME, stopVisitEvent.DepartureStatus)

	// No actual Schedules
	actualSchedule := stopVisitEvent.Schedules.Schedule(schedules.Actual)
	assert.Zero(actualSchedule.ArrivalTime())
	assert.Zero(actualSchedule.DepartureTime())

	// Existing expected schedules
	expectedSchedule := stopVisitEvent.Schedules.Schedule(schedules.Expected)
	expectedDeparture, _ := time.Parse(time.RFC3339, "2023-06-02T08:47:40.000Z")
	assert.Equal(expectedDeparture, expectedSchedule.DepartureTime(), "expected Departure")

	expectedArrival, _ := time.Parse(time.RFC3339, "2023-06-02T08:46:40.000Z")
	assert.Equal(expectedArrival, expectedSchedule.ArrivalTime(), "expected Arrival")

	// Existing aimed schedules
	aimedSchedule := stopVisitEvent.Schedules.Schedule(schedules.Aimed)
	aimedDeparture, _ := time.Parse(time.RFC3339, "2023-06-02T08:47:40.000Z")
	assert.Equal(aimedDeparture, aimedSchedule.DepartureTime())

	aimedArrival, _ := time.Parse(time.RFC3339, "2023-06-02T08:46:40.000Z")
	assert.Equal(aimedArrival, aimedSchedule.ArrivalTime())

}
