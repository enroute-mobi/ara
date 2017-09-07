package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/af83/edwig/model"
)

func Test_EstimatedTimeTableBroadcaster_Receive_Notify(t *testing.T) {
	fakeClock := model.NewFakeClock()
	model.SetDefaultClock(fakeClock)

	// Create a test http server

	response := []byte{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response, _ = ioutil.ReadAll(r.Body)
		w.Header().Add("Content-Type", "text/xml")
	}))
	defer ts.Close()

	// Create a test http server
	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.SetClock(fakeClock)
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "internal"
	partner.Settings["remote_credential"] = "external"
	partner.Settings["remote_url"] = ts.URL

	partner.ConnectorTypes = []string{SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)

	connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).SetClock(fakeClock)
	connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).estimatedTimeTableBroadcaster = NewFakeEstimatedTimeTableBroadcaster(connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster))

	line := referential.Model().Lines().New()
	line.Save()

	objectid := model.NewObjectID("internal", string(line.Id()))
	line.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Id:       string(line.Id()),
		Type:     "Line",
	}

	subscription := partner.Subscriptions().New("sub")
	subscription.SetExternalId("externalId")
	subscription.CreateAddNewResource(reference)
	subscription.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.SetObjectID(objectid)
	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	connector.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).estimatedTimeTableBroadcaster.Start()

	if response == nil {
		t.Errorf("Should have received a notify")
	}
}
