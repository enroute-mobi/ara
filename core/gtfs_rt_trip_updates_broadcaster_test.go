package core

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_TripUpdatesBroadcaster_HandleGtfs(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	connector := NewTripUpdatesBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	saId := model.NewObjectID("objectidKind", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(saId)
	stopArea.Save()

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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	line2 := referential.model.Lines().New()
	iId2 := model.NewObjectID("objectidKind", "lId2")
	line2.SetObjectID(iId2)
	line2.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	vjId2 := model.NewObjectID("objectidKind", "vjId2")
	vehicleJourney2.SetObjectID(vjId2)
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit2 := referential.model.StopVisits().New()
	svId2 := model.NewObjectID("objectidKind", "svId2")
	stopVisit2.SetObjectID(svId2)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit2.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit2.Save()

	stopVisit3 := referential.model.StopVisits().New()
	svId3 := model.NewObjectID("objectidKind", "svId3")
	stopVisit3.SetObjectID(svId3)
	stopVisit3.StopAreaId = stopArea.Id()
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit3.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 2 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 2", l)
	}
	var entity *gtfs.FeedEntity
	if len(gtfsFeed.Entity[0].TripUpdate.StopTimeUpdate) == 1 && len(gtfsFeed.Entity[1].TripUpdate.StopTimeUpdate) == 2 {
		entity = gtfsFeed.Entity[0]
	} else if len(gtfsFeed.Entity[0].TripUpdate.StopTimeUpdate) == 2 && len(gtfsFeed.Entity[1].TripUpdate.StopTimeUpdate) == 1 {
		entity = gtfsFeed.Entity[1]
	} else {
		t.Fatalf("Incorrect number of StopTimeUpdates in gtfs entities:\n got: %v and %v\n want 1 and 2", len(gtfsFeed.Entity[0].TripUpdate.StopTimeUpdate), len(gtfsFeed.Entity[1].TripUpdate.StopTimeUpdate))
	}

	if r := "trip:vjId"; entity.GetId() != r {
		t.Errorf("Response first Feed entity have incorrect Id:\n got: %v\n want: %v", entity.GetId(), r)
	}
	tripUpdate := entity.TripUpdate
	if r := "vjId"; tripUpdate.Trip.GetTripId() != r {
		t.Errorf("Response first Trip Update have incorrect TripId:\n got: %v\n want: %v", tripUpdate.Trip.GetTripId(), r)
	}
	if r := "lId"; tripUpdate.Trip.GetRouteId() != r {
		t.Errorf("Response first Trip Update have incorrect RouteId:\n got: %v\n want: %v", tripUpdate.Trip.GetRouteId(), r)
	}
	// ARA-829
	// if r := "02:10:00"; tripUpdate.Trip.GetStartTime() != r {
	// 	t.Errorf("Response first Trip Update have incorrect StartTime:\n got: %v\n want: %v", tripUpdate.Trip.GetStartTime(), r)
	// }
	if l := len(tripUpdate.StopTimeUpdate); l != 1 {
		t.Errorf("Response first Trip Update have incorrect number of StopTimeUpdate:\n got: %v\n want: 1", l)
	}
	stopTimeUpdate := tripUpdate.StopTimeUpdate[0]
	if r := uint32(1); stopTimeUpdate.GetStopSequence() != r {
		t.Errorf("Incorrect StopSequence in StopTimeUpdate:\n got: %v\n want: %v", stopTimeUpdate.GetStopSequence(), r)
	}
	if r := "saId"; stopTimeUpdate.GetStopId() != r {
		t.Errorf("Incorrect StopId in StopTimeUpdate:\n got: %v\n want: %v", stopTimeUpdate.GetStopId(), r)
	}
	if r := int64(referential.Clock().Now().Add(10 * time.Minute).Unix()); stopTimeUpdate.Departure.GetTime() != r {
		t.Errorf("Incorrect Departure.Time in StopTimeUpdate:\n got: %v\n want: %v", stopTimeUpdate.GetDeparture().Time, r)
	}
	if r := int64(0); stopTimeUpdate.Arrival.GetTime() != r {
		t.Errorf("Incorrect Arrival.Time in StopTimeUpdate:\n got: %v\n want: %v", stopTimeUpdate.GetArrival().Time, r)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongStopIdKind(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	connector := NewTripUpdatesBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	saId := model.NewObjectID("WRONG_ID", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(saId)
	stopArea.Save()

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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Errorf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongLineIdKind(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	connector := NewTripUpdatesBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	saId := model.NewObjectID("objectidKind", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewObjectID("WRONG_ID", "lId")
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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Errorf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongVJIdKind(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	connector := NewTripUpdatesBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	saId := model.NewObjectID("objectidKind", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(saId)
	stopArea.Save()

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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Errorf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongSVIdKind(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	connector := NewTripUpdatesBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	saId := model.NewObjectID("objectidKind", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(saId)
	stopArea.Save()

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
	svId1 := model.NewObjectID("WRONG_ID", "svId1")
	stopVisit.SetObjectID(svId1)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 1 {
		t.Errorf("Response have incorrect number of entities:\n got: %v\n want: 1", l)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_Generators(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("remote_objectid_kind", "objectidKind")
	connector := NewTripUpdatesBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	saId := model.NewObjectID("objectidKind", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(saId)
	stopArea.Save()

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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	l := partner.NewLogStashEvent()
	connector.HandleGtfs(gtfsFeed, l)

	if l := len(gtfsFeed.Entity); l != 1 {
		t.Fatalf("Response have incorrect number of entities:\n got: %v\n want: 1", l)
	}
	entity := gtfsFeed.Entity[0]

	if r := "trip:vjId"; entity.GetId() != r {
		t.Errorf("Response first Feed entity have incorrect Id:\n got: %v\n want: %v", entity.GetId(), r)
	}
	tripUpdate := entity.TripUpdate
	if r := "vjId"; tripUpdate.Trip.GetTripId() != r {
		t.Errorf("Response first Trip Update have incorrect TripId:\n got: %v\n want: %v", tripUpdate.Trip.GetTripId(), r)
	}
	if r := "lId"; tripUpdate.Trip.GetRouteId() != r {
		t.Errorf("Response first Trip Update have incorrect RouteId:\n got: %v\n want: %v", tripUpdate.Trip.GetRouteId(), r)
	}
	stopTimeUpdate := tripUpdate.StopTimeUpdate[0]
	if r := "saId"; stopTimeUpdate.GetStopId() != r {
		t.Errorf("Incorrect StopId in StopTimeUpdate:\n got: %v\n want: %v", stopTimeUpdate.GetStopId(), r)
	}
}
