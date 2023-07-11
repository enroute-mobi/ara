package core

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_VehicleMonitoringBroadcaster_Create_Events(t *testing.T) {
	assert := assert.New(t)
	clock.SetDefaultClock(clock.NewFakeClock())

	referentials := NewMemoryReferentials()
	referential := referentials.New("Un Referential Plutot Cool")
	referential.model = model.NewMemoryModel()

	referential.model.SetBroadcastVeChan(referential.broacasterManager.GetVehicleBroadcastEventChan())
	referential.broacasterManager.Start()
	defer referential.broacasterManager.Stop()

	partner := referential.Partners().New("Un Partner tout autant cool")

	settings := map[string]string{
		"remote_objectid_kind": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	partner.ConnectorTypes = []string{TEST_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER}
	partner.RefreshConnectors()
	referential.Partners().Save(partner)

	connector, _ := partner.Connector(TEST_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER)

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

	subs := partner.Subscriptions().New("VehicleJourneyBroadcast")
	subs.Save()
	subs.CreateAndAddNewResource(reference)
	subs.SetExternalId("externalId")
	subs.Save()

	vehicle := referential.Model().Vehicles().New()
	vehicle.LineId = line.Id()
	vehicle.VehicleJourneyId = vj.Id()

	time.Sleep(5 * time.Millisecond) // Wait for the goRoutine to start ...
	vehicle.Save()

	time.Sleep(5 * time.Millisecond) // Wait for the Broadcaster and Connector to finish their work

	events := connector.(*TestVMSubscriptionBroadcaster).events
	assert.Equal(len(events), 1, "1 event should be generated")
	assert.Equal(events[0].ModelType, "Vehicle", "event should be of type Vehicle")
}
