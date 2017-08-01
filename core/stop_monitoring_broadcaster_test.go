package core

import (
	"testing"
	"time"

	"github.com/af83/edwig/model"
)

func Test_StopMonitoringBroadcaster_Create_Events(t *testing.T) {
	model.SetDefaultClock(model.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	referential.brocasterManager = NewBroadcasterManager(referential)
	referential.brocasterManager.Run()
	referential.Model().StopVisits().(*model.MemoryStopVisits).BroadcastEventChan = referential.brocasterManager.StopVisitBroadcastEvent()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "internal"
	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER_TEST}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER_TEST)

	stopArea := referential.Model().StopAreas().New()
	stopArea.Save()

	objectid := model.NewObjectID("internal", string(stopArea.Id()))
	stopArea.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Id:       string(stopArea.Id()),
		Type:     "StopArea",
	}

	subscriptionRessource, _ := partner.Subscriptions().FindOrCreateByKind("stopMonitoring")
	subscriptionRessource.CreateAddNewResource(reference)

	stopVisit := referential.Model().StopVisits().New()
	stopVisit.StopAreaId = stopArea.Id()
	time.Sleep(5 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()

	time.Sleep(5 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	if len(connector.(*TestStopMonitoringSubscriptionBroadcaster).events) != 1 {
		t.Error("1 event should have been generated got: ", len(connector.(*TestStopMonitoringSubscriptionBroadcaster).events))
	}
}

// func Test_StopMonitoringBroadcaster_Receive_Notify(t *testing.T) {
// 	//fakeClock := model.NewFakeClock()
//
// 	// Create a test http server
// 	response := []byte{}
// 	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		if r.ContentLength <= 0 {
// 			t.Errorf("Notify ContentLength should be zero")
// 		}
// 		r.Body.Read(response)
// 	}))
// 	defer ts.Close()
//
// 	// Create a test http server
// 	referentials := NewMemoryReferentials()
// 	referential := referentials.New("Un Referential Plutot Cool")
// 	referential.model = model.NewMemoryModel()
//
// 	referential.brocasterManager = NewBroadcasterManager(referential)
// 	referential.brocasterManager.Run()
// 	referential.Model().StopVisits().(*model.MemoryStopVisits).BroadcastEventChan = referential.brocasterManager.StopVisitBroadcastEvent()
//
// 	partner := referential.Partners().New("Un Partner tout autant cool")
// 	partner.Settings["remote_objectid_kind"] = "internal"
// 	partner.Settings["remote_url"] = ts.URL
//
// 	partner.ConnectorTypes = []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
// 	partner.RefreshConnectors()
// 	referential.Partners().Save(partner)
//
// 	connector, _ := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
// 	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster = NewFakeStopMonitoringBroadcaster(connector.(*SIRIStopMonitoringSubscriptionBroadcaster))
//
// 	stopArea := referential.Model().StopAreas().New()
// 	stopArea.Save()
//
// 	objectid := model.NewObjectID("internal", string(stopArea.Id()))
// 	stopArea.SetObjectID(objectid)
//
// 	reference := model.Reference{
// 		ObjectId: &objectid,
// 		Id:       string(stopArea.Id()),
// 		Type:     "StopArea",
// 	}
//
// 	subscriptionRessource, _ := partner.Subscriptions().FindOrCreateByKind("stopMonitoring")
// 	subscriptionRessource.CreateAddNewResource(reference)
//
// 	stopVisit := referential.Model().StopVisits().New()
// 	stopVisit.StopAreaId = stopArea.Id()
// 	time.Sleep(10 * time.Millisecond) // Wait for the goRoutine to start ...
// 	stopVisit.Save()
//
// 	time.Sleep(100 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
// 	connector.(*SIRIStopMonitoringSubscriptionBroadcaster).stopMonitoringBroadcaster.Run()
// 	time.Sleep(100 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
//
// 	if len(response) == 0 {
// 		t.Errorf("reponse == %v", response)
// 	}
// }