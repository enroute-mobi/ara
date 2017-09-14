package core

import (
	"testing"
	"time"

	"github.com/af83/edwig/model"
)

func Test_EstimatedTimeTableBroadcaster_Create_Events(t *testing.T) {
	model.SetDefaultClock(model.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	referential.model.SetBroadcastSMChan(referential.broacasterManager.GetStopMonitoringBroadcastEventChan())
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")
	partner.Settings["remote_objectid_kind"] = "internal"
	partner.ConnectorTypes = []string{TEST_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(TEST_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)

	line := referential.Model().Lines().New()
	line.Save()

	objectid := model.NewObjectID("internal", string(line.Id()))
	line.SetObjectID(objectid)

	reference := model.Reference{
		ObjectId: &objectid,
		Type:     "Line",
	}

	vj := referential.Model().VehicleJourneys().New()
	vj.LineId = line.Id()
	vj.Save()

	sv := referential.Model().StopVisits().New()
	sv.VehicleJourneyId = vj.Id()

	subs := partner.Subscriptions().New("EstimatedTimeTable")
	subs.Save()
	subs.CreateAddNewResource(reference)
	subs.SetExternalId("externalId")
	subs.Save()

	stopVisit := referential.Model().StopVisits().New()

	time.Sleep(5 * time.Millisecond) // Wait for the goRoutine to start ...
	stopVisit.Save()

	time.Sleep(5 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work
	if len(connector.(*TestETTSubscriptionBroadcaster).events) != 1 {
		t.Error("1 event should have been generated got: ", len(connector.(*TestETTSubscriptionBroadcaster).events))
	}
}
