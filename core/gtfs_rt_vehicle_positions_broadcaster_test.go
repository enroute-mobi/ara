package core

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_VehiclePositionBroadcaster_HandleGtfs(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewVehiclePositionBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("codeSpace", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("codeSpace", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.DirectionType = model.VEHICLE_DIRECTION_INBOUND
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("codeSpace", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewCode("codeSpace", "vId")
	vehicle.SetCode(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewCode("codeSpace", "vId2")
	vehicle2.SetCode(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)
	assert.Len(gtfsFeed.Entity, 2)

	var entity *gtfs.FeedEntity

	if gtfsFeed.Entity[0].GetId() == "vehicle:vId" && gtfsFeed.Entity[1].GetId() == "vehicle:vId2" {
		entity = gtfsFeed.Entity[0]
	} else if gtfsFeed.Entity[1].GetId() == "vehicle:vId" && gtfsFeed.Entity[0].GetId() == "vehicle:vId2" {
		entity = gtfsFeed.Entity[1]
	} else {
		t.Fatalf("Incorrect Ids for gtfs entities:\n got: %v and %v\n want vehicle:vId and vehicle:vId2", gtfsFeed.Entity[0].GetId(), gtfsFeed.Entity[1].GetId())
	}

	assert.Equal("vId", entity.GetVehicle().GetVehicle().GetId())

	trip := entity.GetVehicle().GetTrip()
	assert.Equal("vjId", trip.GetTripId())
	assert.Equal("lId", trip.GetRouteId())
	assert.Equal(uint32(1), trip.GetDirectionId())

	position := entity.GetVehicle().GetPosition()
	assert.Equal(float32(1.23456), position.GetLongitude())
	assert.Equal(float32(2.34567), position.GetLatitude())
	assert.Equal(float32(1.2), position.GetBearing())
}

func Test_VehiclePositionBroadcaster_HandleGtfs_WrongLineId(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewVehiclePositionBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("WRONG_KIND", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("codeSpace", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("codeSpace", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewCode("codeSpace", "vId")
	vehicle.SetCode(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewCode("codeSpace", "vId2")
	vehicle2.SetCode(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_VehiclePositionBroadcaster_HandleGtfs_WrongVJId(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewVehiclePositionBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("codeSpace", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("WRONG_ID", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("codeSpace", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewCode("codeSpace", "vId")
	vehicle.SetCode(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewCode("codeSpace", "vId2")
	vehicle2.SetCode(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_VehiclePositionBroadcaster_HandleGtfs_WrongVehicleId(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewVehiclePositionBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("codeSpace", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("codeSpace", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("codeSpace", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewCode("WRONG_ID", "vId")
	vehicle.SetCode(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewCode("codeSpace", "vId2")
	vehicle2.SetCode(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 1 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 1", l)
	}
	if gtfsFeed.Entity[0].GetId() != "vehicle:vId2" {
		t.Errorf("Response have the wrong Vehicle ID:\n got: %v\n want: vehicle:vId", gtfsFeed.Entity[0].GetId())
	}
}

func Test_VehiclePositionBroadcaster_HandleGtfs_WrongVehicleIdWithSetting(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"remote_code_space": "codeSpace",
		"gtfs-rt-vehicle-positions-broadcaster.vehicle_remote_code_space": "WRONG_ID",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewVehiclePositionBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("codeSpace", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("codeSpace", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("codeSpace", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewCode("WRONG_ID", "vId")
	vehicle.SetCode(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewCode("codeSpace", "vId2")
	vehicle2.SetCode(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 1 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 1", l)
	}
	if gtfsFeed.Entity[0].GetId() != "vehicle:vId" {
		t.Errorf("Response have the wrong Vehicle ID:\n got: %v\n want: vehicle:vId", gtfsFeed.Entity[0].GetId())
	}
}

func Test_VehiclePositionBroadcaster_HandleGtfs_Generators(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewVehiclePositionBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("codeSpace", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("codeSpace", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("codeSpace", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewCode("codeSpace", "vId")
	vehicle.SetCode(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 1 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 1", l)
	}
	entity := gtfsFeed.Entity[0]

	trip := entity.GetVehicle().GetTrip()
	if r := "vjId"; trip.GetTripId() != r {
		t.Errorf("Incorrect TripId for entity TripUpdate:\n got: %v\n want: %v", trip.GetTripId(), r)
	}
	if r := "lId"; trip.GetRouteId() != r {
		t.Errorf("Incorrect RouteId for entity TripUpdate:\n got: %v\n want: %v", trip.GetRouteId(), r)
	}
}
