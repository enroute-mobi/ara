package core

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
)

func Test_VehiclePositionBroadcaster_HandleGtfs(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	connector := NewVehiclePositionBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.model.Lines().New()
	lId := model.NewObjectID("objectidKind", "lId")
	line.SetObjectID(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewObjectID("objectidKind", "vjId")
	vehicleJourney.SetObjectID(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewObjectID("objectidKind", "svId1")
	stopVisit.SetObjectID(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewObjectID("objectidKind", "vId")
	vehicle.SetObjectID(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewObjectID("objectidKind", "vId2")
	vehicle2.SetObjectID(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 2 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 2", l)
	}
	var entity *gtfs.FeedEntity
	if gtfsFeed.Entity[0].GetId() == "vehicle:vId" && gtfsFeed.Entity[1].GetId() == "vehicle:vId2" {
		entity = gtfsFeed.Entity[0]
	} else if gtfsFeed.Entity[1].GetId() == "vehicle:vId" && gtfsFeed.Entity[0].GetId() == "vehicle:vId2" {
		entity = gtfsFeed.Entity[1]
	} else {
		t.Fatalf("Incorrect Ids for gtfs entities:\n got: %v and %v\n want vehicle:vId and vehicle:vId2", gtfsFeed.Entity[0].GetId(), gtfsFeed.Entity[1].GetId())
	}

	if r := "vId"; entity.GetVehicle().GetVehicle().GetId() != r {
		t.Errorf("Incorrect Vehicle Id:\n got: %v\n want: %v", entity.GetVehicle().GetVehicle().GetId(), r)
	}

	trip := entity.GetVehicle().GetTrip()
	if r := "vjId"; trip.GetTripId() != r {
		t.Errorf("Incorrect TripId for entity TripUpdate:\n got: %v\n want: %v", trip.GetTripId(), r)
	}
	if r := "lId"; trip.GetRouteId() != r {
		t.Errorf("Incorrect RouteId for entity TripUpdate:\n got: %v\n want: %v", trip.GetRouteId(), r)
	}
	if r := connector.Clock().Now().Add(10 * time.Minute).Format("15:04:05"); trip.GetStartTime() != r {
		t.Errorf("Incorrect StartTime for entity TripUpdate:\n got: %v\n want: %v", trip.GetStartTime(), r)
	}

	position := entity.GetVehicle().GetPosition()
	if r := float32(1.23456); position.GetLongitude() != r {
		t.Errorf("Incorrect Longitude for entity PositionUpdate:\n got: %v\n want: %v", position.GetLongitude(), r)
	}
	if r := float32(2.34567); position.GetLatitude() != r {
		t.Errorf("Incorrect Latitude for entity PositionUpdate:\n got: %v\n want: %v", position.GetLatitude(), r)
	}
	if r := float32(1.2); position.GetBearing() != r {
		t.Errorf("Incorrect Bearing for entity PositionUpdate:\n got: %v\n want: %v", position.GetBearing(), r)
	}
}

func Test_VehiclePositionBroadcaster_HandleGtfs_WrongLineId(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	connector := NewVehiclePositionBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.model.Lines().New()
	lId := model.NewObjectID("WRONG_KIND", "lId")
	line.SetObjectID(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewObjectID("objectidKind", "vjId")
	vehicleJourney.SetObjectID(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewObjectID("objectidKind", "svId1")
	stopVisit.SetObjectID(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewObjectID("objectidKind", "vId")
	vehicle.SetObjectID(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewObjectID("objectidKind", "vId2")
	vehicle2.SetObjectID(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_VehiclePositionBroadcaster_HandleGtfs_WrongVJId(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	connector := NewVehiclePositionBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.model.Lines().New()
	lId := model.NewObjectID("objectidKind", "lId")
	line.SetObjectID(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewObjectID("WRONG_ID", "vjId")
	vehicleJourney.SetObjectID(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewObjectID("objectidKind", "svId1")
	stopVisit.SetObjectID(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewObjectID("objectidKind", "vId")
	vehicle.SetObjectID(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewObjectID("objectidKind", "vId2")
	vehicle2.SetObjectID(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_VehiclePositionBroadcaster_HandleGtfs_WrongVehicleId(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	connector := NewVehiclePositionBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.model.Lines().New()
	lId := model.NewObjectID("objectidKind", "lId")
	line.SetObjectID(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewObjectID("objectidKind", "vjId")
	vehicleJourney.SetObjectID(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewObjectID("objectidKind", "svId1")
	stopVisit.SetObjectID(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewObjectID("WRONG_ID", "vId")
	vehicle.SetObjectID(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewObjectID("objectidKind", "vId2")
	vehicle2.SetObjectID(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

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
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["gtfs-rt-vehicle-positions-broadcaster.vehicle_remote_objectid_kind"] = "WRONG_ID"
	connector := NewVehiclePositionBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.model.Lines().New()
	lId := model.NewObjectID("objectidKind", "lId")
	line.SetObjectID(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewObjectID("objectidKind", "vjId")
	vehicleJourney.SetObjectID(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewObjectID("objectidKind", "svId1")
	stopVisit.SetObjectID(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewObjectID("WRONG_ID", "vId")
	vehicle.SetObjectID(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	vehicle2 := referential.model.Vehicles().New()
	vId2 := model.NewObjectID("objectidKind", "vId2")
	vehicle2.SetObjectID(vId2)
	vehicle2.VehicleJourneyId = vehicleJourney.Id()
	vehicle2.LineId = line.Id()
	vehicle2.Longitude = 3.45678
	vehicle2.Latitude = 4.56789
	vehicle2.Bearing = 2.3
	vehicle2.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

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
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.reference_identifier"] = "%{type}:%{objectid}:LOC"
	connector := NewVehiclePositionBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.model.Lines().New()
	lId := model.NewObjectID("objectidKind", "lId")
	line.SetObjectID(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewObjectID("objectidKind", "vjId")
	vehicleJourney.SetObjectID(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewObjectID("objectidKind", "svId1")
	stopVisit.SetObjectID(svId1)
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewObjectID("objectidKind", "vId")
	vehicle.SetObjectID(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.LineId = line.Id()
	vehicle.Longitude = 1.23456
	vehicle.Latitude = 2.34567
	vehicle.Bearing = 1.2
	vehicle.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 1 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 1", l)
	}
	entity := gtfsFeed.Entity[0]

	trip := entity.GetVehicle().GetTrip()
	if r := "VehicleJourney:vjId:LOC"; trip.GetTripId() != r {
		t.Errorf("Incorrect TripId for entity TripUpdate:\n got: %v\n want: %v", trip.GetTripId(), r)
	}
	if r := "Line:lId:LOC"; trip.GetRouteId() != r {
		t.Errorf("Incorrect RouteId for entity TripUpdate:\n got: %v\n want: %v", trip.GetRouteId(), r)
	}
}
