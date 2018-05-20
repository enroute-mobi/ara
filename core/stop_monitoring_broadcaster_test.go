package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_StopMonitoringBroadcaster_Create_Events(t *testing.T) {
	model.SetDefaultClock(model.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	referential.model.SetBroadcastSMChan(referential.broacasterManager.GetStopMonitoringBroadcastEventChan())
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "internal"
	partner.ConnectorTypes = []string{TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	objectid := model.NewObjectID("internal", string(stopArea.Id()))
	stopArea.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Type:     "StopArea",
	}

	subs := partner.Subscriptions().New("kind")
	subs.Save()
	subs.CreateAddNewResource(reference)
	subs.SetExternalId("externalId")
	subs.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	time.Sleep(5 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()

	time.Sleep(5 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	if len(connector.(*TestStopMonitoringSubscriptionBroadcaster).events) != 2 {
		t.Error("2 events should have been generated got: ", len(connector.(*TestStopMonitoringSubscriptionBroadcaster).events))
	}
}

func Test_StopMonitoringBroadcaster_HandleStopMonitoringBroadcastWithReferent(t *testing.T) {
	model.SetDefaultClock(model.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.Save()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "internal"
	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	stopArea := referential.Model().StopAreas().New()
	objectid := model.NewObjectID("internal", string(stopArea.Id()))
	stopArea.SetObjectID(objectid)
	stopArea.Save()

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.ReferentId = stopArea.Id()
	stopArea2.Save()

	reference := model.Reference{
		ObjectId: &objectid,
		Type:     "StopArea",
	}

	subs := partner.Subscriptions().New("StopMonitoringBroadcast")
	subs.CreateAddNewResource(reference)
	subs.SetExternalId("externalId")
	subs.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea2.Id()
	stopVisit.Save()

	event := &model.StopMonitoringBroadcastEvent{
		ModelId:   string(stopVisit.Id()),
		ModelType: "StopVisit",
	}

	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)

	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleStopMonitoringBroadcastEvent(event)
	if len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast) != 1 {
		t.Error("1 events should have been generated got: ", len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast))
	}
}

func Test_StopMonitoringBroadcaster_Receive_Notify(t *testing.T) {
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

	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).SetClock(fakeClock)
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster = NewFakeStopMonitoringBroadcaster(connector.(*SIRIStopMonitoringSubscriptionBroadcaster))

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	objectid := model.NewObjectID("internal", string(stopArea.Id()))
	stopArea.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Type:     "StopArea",
	}

	stopArea2 := referential.Model().StopAreas().New()
	stopArea2.Save()

	objectid2 := model.NewObjectID("internal", string(stopArea2.Id()))
	stopArea2.SetObjectID(objectid2)

	reference2 := model.Reference{
		ObjectId: &objectid2,
		Type:     "StopArea",
	}

	subscription := partner.Subscriptions().New("StopMonitoringBroadcast")
	subscription.SetExternalId("externalId")
	subscription.CreateAddNewResource(reference)
	subscription.CreateAddNewResource(reference2)
	subscription.subscriptionOptions["ChangeBeforeUpdates"] = "PT4M"
	subscription.subscriptionOptions["MaximumStopVisits"] = "1"
	subscription.Save()

	line := referential.Model().Lines().New()
	line.SetObjectID(objectid)
	line.Save()

	vehicleJourney := referential.Model().VehicleJourneys().New()
	vehicleJourney.LineId = line.Id()
	vehicleJourney.SetObjectID(objectid)
	vehicleJourney.Save()

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.SetObjectID(objectid)
	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	stopVisit2 := referential.Model().StopVisits().New()
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney.Id()
	stopVisit2.SetObjectID(model.NewObjectID("internal", string(stopArea.Id())))
	stopVisit2.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))

	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()
	time.Sleep(10 * time.Millisecond)
	stopVisit2.Save()

	time.Sleep(10 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster.Start()

	notify, _ := siri.NewXMLNotifyStopMonitoringFromContent(response)
	delivery := notify.StopMonitoringDeliveries()

	if len(delivery) != 1 {
		t.Errorf("Should have received 1 delivery but got == %v", len(delivery))
	}

	if delivery[0].SubscriberRef() != "external" {
		t.Errorf("SubscriberRef should be external but got == %v", delivery[0].SubscriberRef())
	}

	if delivery[0].SubscriptionRef() != "externalId" {
		t.Errorf("SubscriptionRef should be externalId but got == %v", delivery[0].SubscriptionRef())
	}

	sv := delivery[0].XMLMonitoredStopVisits()

	if len(sv) != 1 {
		t.Errorf("Should have received 1 StopVisit but got == %v", len(sv))
	}

	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(1*time.Minute))
	stopVisit.Save()

	time.Sleep(10 * time.Millisecond)
	if len := len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast); len != 0 {
		t.Errorf("No stopVisit should need to be broadcasted %v", len)
	}

	stopVisit.Schedules.SetArrivalTime("actual", referential.Clock().Now().Add(10*time.Minute))
	stopVisit.Save()

	time.Sleep(10 * time.Millisecond)
	if len := len(connector.(*SIRIStopMonitoringSubscriptionBroadcaster).toBroadcast); len != 1 {
		t.Errorf("1 stopVisit should need to be broadcasted %v", len)
	}
}
