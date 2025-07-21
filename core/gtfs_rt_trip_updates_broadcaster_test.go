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

func Test_TripUpdatesBroadcaster_HandleGtfs(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("codeSpace", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("codeSpace", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("codeSpace", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.DirectionType = model.VEHICLE_DIRECTION_OUTBOUND
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("codeSpace", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	line2 := referential.model.Lines().New()
	iId2 := model.NewCode("codeSpace", "lId2")
	line2.SetCode(iId2)
	line2.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	vjId2 := model.NewCode("codeSpace", "vjId2")
	vehicleJourney2.SetCode(vjId2)
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit2 := referential.model.StopVisits().New()
	svId2 := model.NewCode("codeSpace", "svId2")
	stopVisit2.SetCode(svId2)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit2.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit2.Save()

	stopVisit3 := referential.model.StopVisits().New()
	svId3 := model.NewCode("codeSpace", "svId3")
	stopVisit3.SetCode(svId3)
	stopVisit3.StopAreaId = stopArea.Id()
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit3.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)
	assert.Len(gtfsFeed.Entity, 2)

	var entity *gtfs.FeedEntity
	if len(gtfsFeed.Entity[0].TripUpdate.StopTimeUpdate) == 1 && len(gtfsFeed.Entity[1].TripUpdate.StopTimeUpdate) == 2 {
		entity = gtfsFeed.Entity[0]
	} else if len(gtfsFeed.Entity[0].TripUpdate.StopTimeUpdate) == 2 && len(gtfsFeed.Entity[1].TripUpdate.StopTimeUpdate) == 1 {
		entity = gtfsFeed.Entity[1]
	} else {
		t.Fatalf("Incorrect number of StopTimeUpdates in gtfs entities:\n got: %v and %v\n want 1 and 2", len(gtfsFeed.Entity[0].TripUpdate.StopTimeUpdate), len(gtfsFeed.Entity[1].TripUpdate.StopTimeUpdate))
	}

	assert.Equal("trip:vjId", entity.GetId())

	tripUpdate := entity.TripUpdate
	assert.Equal("vjId", tripUpdate.Trip.GetTripId())
	assert.Equal("lId", tripUpdate.Trip.GetRouteId())
	assert.Equal(uint32(1), tripUpdate.Trip.GetDirectionId())
	assert.Len(tripUpdate.StopTimeUpdate, 1)

	stopTimeUpdate := tripUpdate.StopTimeUpdate[0]
	assert.Equal(uint32(0), stopTimeUpdate.GetStopSequence())
	assert.Equal("saId", stopTimeUpdate.GetStopId())
	assert.Equal(referential.Clock().Now().Add(10*time.Minute).Unix(), stopTimeUpdate.Departure.GetTime())
	assert.Equal(int64(0), stopTimeUpdate.Arrival.GetTime())
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongStopIdCodeSpace(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("WRONG_ID", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Errorf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongLineIdCodeSpace(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("codeSpace", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("WRONG_ID", "lId")
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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Errorf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongVJIdCodeSpace(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("codeSpace", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 0 {
		t.Errorf("Response have incorrect number of entities:\n got: %v\n want: 0", l)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongSVIdCodeSpace(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("codeSpace", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

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
	svId1 := model.NewCode("WRONG_ID", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetArrivalTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

	if l := len(gtfsFeed.Entity); l != 1 {
		t.Errorf("Response have incorrect number of entities:\n got: %v\n want: 1", l)
	}
}

func Test_TripUpdatesBroadcaster_HandleGtfs_Generators(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("codeSpace", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

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
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)

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
