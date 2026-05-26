package core

import (
	"fmt"
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

	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("internal", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("internal", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.DirectionType = model.VEHICLE_DIRECTION_OUTBOUND
	vehicleJourney.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewCode("internal", "vId")
	vehicle.SetCode(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("internal", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 12
	stopVisit.Save()

	line2 := referential.model.Lines().New()
	iId2 := model.NewCode("internal", "lId2")
	line2.SetCode(iId2)
	line2.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	vjId2 := model.NewCode("internal", "vjId2")
	vehicleJourney2.SetCode(vjId2)
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit2 := referential.model.StopVisits().New()
	svId2 := model.NewCode("internal", "svId2")
	stopVisit2.SetCode(svId2)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit2.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit2.PassageOrder = 34
	stopVisit2.Save()

	stopVisit3 := referential.model.StopVisits().New()
	svId3 := model.NewCode("internal", "svId3")
	stopVisit3.SetCode(svId3)
	stopVisit3.StopAreaId = stopArea.Id()
	stopVisit3.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit3.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit3.PassageOrder = 56
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
	assert.Equal(uint32(0), tripUpdate.Trip.GetDirectionId())
	assert.Len(tripUpdate.StopTimeUpdate, 1)

	vehicleDescriptor := tripUpdate.Vehicle
	assert.NotNil(vehicleDescriptor)
	assert.Equal("vId", vehicleDescriptor.GetId())

	stopTimeUpdate := tripUpdate.StopTimeUpdate[0]
	assert.Equal(uint32(0), stopTimeUpdate.GetStopSequence())
	assert.Equal("saId", stopTimeUpdate.GetStopId())
	assert.Equal(connector.Clock().Now().Add(10*time.Minute).Unix(), stopTimeUpdate.Departure.GetTime())
	assert.Equal(int64(0), stopTimeUpdate.Arrival.GetTime())
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongVehicleIdWithSetting(t *testing.T) {
	assert := assert.New(t)

	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
		"gtfs-rt-trip-updates-broadcaster.vehicle_remote_code_space": "external",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("internal", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("internal", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.DirectionType = model.VEHICLE_DIRECTION_OUTBOUND
	vehicleJourney.Save()

	vehicle := referential.model.Vehicles().New()
	vId := model.NewCode("internal", "vId")
	vehicle.SetCode(vId)
	vehicle.VehicleJourneyId = vehicleJourney.Id()
	vehicle.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("internal", "svId1")
	stopVisit.SetCode(svId1)
	stopVisit.StopAreaId = stopArea.Id()
	stopVisit.VehicleJourneyId = vehicleJourney.Id()
	stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit.PassageOrder = 1
	stopVisit.Save()

	line2 := referential.model.Lines().New()
	iId2 := model.NewCode("internal", "lId2")
	line2.SetCode(iId2)
	line2.Save()

	vehicleJourney2 := referential.model.VehicleJourneys().New()
	vjId2 := model.NewCode("internal", "vjId2")
	vehicleJourney2.SetCode(vjId2)
	vehicleJourney2.LineId = line2.Id()
	vehicleJourney2.Save()

	stopVisit2 := referential.model.StopVisits().New()
	svId2 := model.NewCode("internal", "svId2")
	stopVisit2.SetCode(svId2)
	stopVisit2.StopAreaId = stopArea.Id()
	stopVisit2.VehicleJourneyId = vehicleJourney2.Id()
	stopVisit2.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(10*time.Minute))
	stopVisit2.Save()

	stopVisit3 := referential.model.StopVisits().New()
	svId3 := model.NewCode("internal", "svId3")
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
	assert.Equal(uint32(0), tripUpdate.Trip.GetDirectionId())
	assert.Len(tripUpdate.StopTimeUpdate, 1)

	vehicleDescriptor := tripUpdate.Vehicle
	assert.Equal(vehicleDescriptor, &gtfs.VehicleDescriptor{}, "vehicleDescriptor should be empty with wrong code_space for vehicle")
}

func Test_TripUpdatesBroadcaster_HandleGtfs_WrongStopIdCodeSpace(t *testing.T) {
	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("external", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("internal", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("internal", "svId1")
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
	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("internal", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("external", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("internal", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("internal", "svId1")
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
	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("internal", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("external", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("internal", "svId1")
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
	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("internal", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("internal", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("external", "svId1")
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
	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	saId := model.NewCode("internal", "saId")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(saId)
	stopArea.Save()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("internal", "vjId")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	stopVisit := referential.model.StopVisits().New()
	svId1 := model.NewCode("internal", "svId1")
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

func Test_RewriteStopSequence(t *testing.T) {
	assert := assert.New(t)

	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("internal", "value")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.Save()

	// stopAreas & stopVisits
	for j := 0; j < 5; j++ {
		saId := model.NewCode("internal", fmt.Sprintf("saId%d", j))
		stopArea := referential.Model().StopAreas().New()
		stopArea.SetCode(saId)
		stopArea.Save()

		stopVisit := referential.model.StopVisits().New()
		svId1 := model.NewCode("internal", fmt.Sprintf("svId%d", j))
		stopVisit.SetCode(svId1)
		stopVisit.StopAreaId = stopArea.Id()
		stopVisit.VehicleJourneyId = vehicleJourney.Id()
		stopVisit.Schedules.SetDepartureTime("actual", connector.Clock().Now().Add(time.Duration(float64(j*1e9))+10*time.Minute))
		stopVisit.PassageOrder = j * 3
		stopVisit.Save()
	}

	stopVisits := connector.partner.Model().StopVisits().FindAll()
	var stopVisitPassageOrders []int
	for i := range stopVisits {
		stopVisitPassageOrders = append(stopVisitPassageOrders, stopVisits[i].PassageOrder)
	}
	assert.ElementsMatch([]int{0, 3, 6, 9, 12}, stopVisitPassageOrders)

	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)
	assert.Len(gtfsFeed.Entity, 1)

	entity := gtfsFeed.Entity[0]

	var gtfsStopSequences []int
	for i := range entity.TripUpdate.StopTimeUpdate {
		gtfsStopSequences = append(gtfsStopSequences, int(entity.TripUpdate.StopTimeUpdate[i].GetStopSequence()))
	}
	assert.ElementsMatch([]int{0, 1, 2, 3, 4}, gtfsStopSequences)
}

func Test_RewriteStopSequenceWithPassedStopVisits(t *testing.T) {
	assert := assert.New(t)

	// Setup
	_, referential := newTestReferential(t)
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewTripUpdatesBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.model.Lines().New()
	lId := model.NewCode("internal", "lId")
	line.SetCode(lId)
	line.Save()

	vehicleJourney := referential.model.VehicleJourneys().New()
	vjId := model.NewCode("internal", "value")
	vehicleJourney.SetCode(vjId)
	vehicleJourney.LineId = line.Id()
	vehicleJourney.AimedStopVisitCount = 20
	vehicleJourney.Save()

	// stopAreas & stopVisits
	for j := 0; j < 5; j++ {
		saId := model.NewCode("internal", fmt.Sprintf("saId%d", j))
		stopArea := referential.Model().StopAreas().New()
		stopArea.SetCode(saId)
		stopArea.Save()

		stopVisit := referential.model.StopVisits().New()
		svId1 := model.NewCode("internal", fmt.Sprintf("svId%d", j))
		stopVisit.SetCode(svId1)
		stopVisit.StopAreaId = stopArea.Id()
		stopVisit.VehicleJourneyId = vehicleJourney.Id()
		delta := time.Duration(float64(j) * float64(time.Minute))
		base := connector.Clock().Now().Add(-4 * time.Minute)
		stopVisit.Schedules.SetDepartureTime("actual", base.Add(delta))
		stopVisit.PassageOrder = j * 3
		stopVisit.Save()
	}

	// verify Setup
	stopVisits := connector.partner.Model().StopVisits().FindAll()
	var stopVisitPassageOrders []int
	for i := range stopVisits {
		stopVisitPassageOrders = append(stopVisitPassageOrders, stopVisits[i].PassageOrder)
	}
	assert.ElementsMatch([]int{0, 3, 6, 9, 12}, stopVisitPassageOrders)

	// verify Broadcastble stopVisits
	referenceTime := connector.Clock().Now().Add(PAST_STOP_VISITS_MAX_TIME)
	broadcastableStopVisits := connector.partner.Model().StopVisits().FindByVehicleJourneyIdAfterUnsorted(vehicleJourney.Id(), referenceTime)
	assert.Len(broadcastableStopVisits, 2)

	var broadcastableStopVisitsPassageOrder []int
	for i := range broadcastableStopVisits {
		broadcastableStopVisitsPassageOrder = append(broadcastableStopVisitsPassageOrder, broadcastableStopVisits[i].PassageOrder)
	}
	assert.ElementsMatch([]int{9, 12}, broadcastableStopVisitsPassageOrder)

	// GTFS feed
	gtfsFeed := &gtfs.FeedMessage{}

	connector.HandleGtfs(gtfsFeed)
	assert.Len(gtfsFeed.Entity, 1)

	entity := gtfsFeed.Entity[0]

	var gtfsStopSequences []int
	for i := range entity.TripUpdate.StopTimeUpdate {
		gtfsStopSequences = append(gtfsStopSequences, int(entity.TripUpdate.StopTimeUpdate[i].GetStopSequence()))
	}
	// vehicleJourney has 20 AimedStopVisits
	// vehicleJourney has 5 stopVisits in memory
	//
	// Only 2 can be broadcasted because of their reference time
	// PassageOrder are 9 and 12
	//
	// The passageOrder offset should be 20 - 2 = 18
	// the broadcast should start at passageOrder 18
	assert.ElementsMatch([]int{18, 19}, gtfsStopSequences)
}
