package core

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	p "bitbucket.org/enroute-mobi/ara/core/partners"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func collectGtfs(t *testing.T, feed *gtfs.FeedMessage, fakeBroadcast bool) ([]model.UpdateEvent, *Partner) {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := proto.Marshal(feed)
		var buffer bytes.Buffer
		buffer.Write(data)
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.Write(buffer.Bytes())
	}))
	defer ts.Close()

	// Create a GtfsRequestCollector
	_, referential := newTestReferential(t)
	partners := referential.partners
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_url":        ts.URL,
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	gtfsCollector := NewGtfsRequestCollector(partner)
	gtfsCollector.SetClock(clock.NewFakeClock())

	if fakeBroadcast {
		fs := fakeBroadcaster{}
		gtfsCollector.SetSubscriber(fs.FakeBroadcaster)

		gtfsCollector.requestGtfs()
		time.Sleep(42 * time.Millisecond)

		return fs.Events, partner
	}

	gtfsCollector.Start()
	time.Sleep(42 * time.Millisecond)
	gtfsCollector.Clock().(clock.FakeClock).Advance(10 * time.Second)
	gtfsCollector.Stop()
	time.Sleep(42 * time.Millisecond)

	return nil, partner
}

func Test_GtfsCollectEvents(t *testing.T) {
	e := []*gtfs.FeedEntity{
		tripUpdate(),
		vehiclePosition(),
	}
	feed := newGtfsFeed(e)

	events, partner := collectGtfs(t, feed, true)

	if len(events) != 7 {
		t.Errorf("Should have 7 events after gtfs collect, got %v", len(events))
	}
	if partner.alternativeStatusCheck.LastCheck != clock.FAKE_CLOCK_INITIAL_DATE {
		t.Errorf("Partner alternative status time should be updated, got: %v", partner.alternativeStatusCheck.LastCheck)
	}
	if partner.alternativeStatusCheck.Status != p.OperationnalStatusUp {
		t.Errorf("Partner alternative status status should be updated, got: %v", partner.alternativeStatusCheck.Status)
	}
}

func Test_GtfsCollectEvents_SameEntities(t *testing.T) {
	e := []*gtfs.FeedEntity{
		tripUpdate(),
		tripUpdate(),
	}
	feed := newGtfsFeed(e)

	events, _ := collectGtfs(t, feed, true)

	if len(events) != 6 {
		t.Errorf("Should have 6 events after gtfs collect, got %v", len(events))
	}
}

func Test_GtfsCollect(t *testing.T) {
	e := []*gtfs.FeedEntity{
		tripUpdate(),
		vehiclePosition(),
	}
	feed := newGtfsFeed(e)

	_, partner := collectGtfs(t, feed, false)

	if c := len(partner.Referential().Model().StopAreas().FindAll()); c != 2 {
		t.Errorf("2 StopAreas should have been created, found %v", c)
	}
	if c := len(partner.Referential().Model().Lines().FindAll()); c != 1 {
		t.Errorf("1 Line should have been created, found %v", c)
	}
	if c := len(partner.Referential().Model().VehicleJourneys().FindAll()); c != 1 {
		t.Errorf("1 VehicleJourney should have been created, found %v", c)
	}
	if c := len(partner.Referential().Model().StopVisits().FindAll()); c != 2 {
		t.Errorf("2 StopVisits should have been created, found %v", c)
	}
	if c := len(partner.Referential().Model().Vehicles().FindAll()); c != 1 {
		t.Errorf("1 Vehicle should have been created, found %v", c)
	}
}

func Test_PartnerStatusDown(t *testing.T) {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	// Create a GtfsRequestCollector
	partners := newTestPartnerManager(t)
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_url":        ts.URL,
		"remote_code_space": "internal",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	gtfsCollector := NewGtfsRequestCollector(partner)
	gtfsCollector.SetClock(clock.NewFakeClock())
	gtfsCollector.requestGtfs()

	if partner.alternativeStatusCheck.LastCheck != clock.FAKE_CLOCK_INITIAL_DATE {
		t.Errorf("Partner alternative status time should be updated, got: %v", partner.alternativeStatusCheck.LastCheck)
	}
	if partner.alternativeStatusCheck.Status != p.OperationnalStatusDown {
		t.Errorf("Partner alternative status status should be updated, got: %v", partner.alternativeStatusCheck.Status)
	}
}

func Test_GtfsCollect_SkipExpiredStops(t *testing.T) {
	assert := assert.New(t)

	e := []*gtfs.FeedEntity{
		tripUpdateWithExpiredStops(),
	}
	feed := newGtfsFeed(e)

	_, partner := collectGtfs(t, feed, false)

	// The trip has 4 stops:
	// - 1 stop visit with both times after the persistence cutoff -> kept
	// - 1 stop visit departure-only after the persistence cutoff -> kept
	// - 1 stop visit with both times already older than the cutoff -> dropped
	// - 1 stop visit departure-only already older than the cutoff -> dropped
	// Only the two stops whose reference time survives persistence must produce a
	// StopVisit; the expired ones must be dropped rather than created and then
	// cleaned by the guardian on its next routine, every collect cycle.
	deleteBefore := clock.FAKE_CLOCK_INITIAL_DATE.Add(-3 * time.Hour)

	svs := partner.Referential().Model().StopVisits().FindAll()
	assert.Len(svs, 2, "only stops after the persistence cutoff should be created")

	for _, sv := range svs {
		assert.Falsef(sv.ReferenceTime().Before(deleteBefore),
			"created StopVisit reference time should survive persistence, got %v", sv.ReferenceTime())
	}
}

func Test_GtfsCollect_SkipFullyExpiredTrip(t *testing.T) {
	assert := assert.New(t)

	e := []*gtfs.FeedEntity{
		tripUpdateFullyExpired(),
	}
	feed := newGtfsFeed(e)

	_, partner := collectGtfs(t, feed, false)

	// Every stop of the trip is past the persistence window, so nothing is
	// collected — and crucially no empty VehicleJourney is left behind, since the
	// guardian only cleans VehicleJourneys whose StopVisits were deleted.
	assert.Empty(partner.Referential().Model().StopVisits().FindAll(), "no StopVisit should be created")
	assert.Empty(partner.Referential().Model().VehicleJourneys().FindAll(), "no empty VehicleJourney should be created")
}

func tripUpdateWithExpiredStops() *gtfs.FeedEntity {
	id := "id"
	tid := "tid"
	rid := "rid"

	// FAKE_CLOCK_INITIAL_DATE is 1984; persisted times are well after it,
	// expired times well before it.
	persisted := int64(1601875200)    // 2020, after the persistence cutoff
	persistedDep := int64(1601875610) // 2020
	expired := int64(315532800)       // 1980, before the persistence cutoff
	expiredDep := int64(315533000)    // 1980

	ss1 := uint32(1)
	sid1 := "sid1"
	ss2 := uint32(2)
	sid2 := "sid2"
	ss3 := uint32(3)
	sid3 := "sid3"
	ss4 := uint32(4)
	sid4 := "sid4"

	return &gtfs.FeedEntity{
		Id: &id,
		TripUpdate: &gtfs.TripUpdate{
			Trip: &gtfs.TripDescriptor{
				TripId:  &tid,
				RouteId: &rid,
			},
			StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{
				{ // regular stop after the cutoff, kept
					StopSequence: &ss1,
					StopId:       &sid1,
					Arrival:      &gtfs.TripUpdate_StopTimeEvent{Time: &persisted},
					Departure:    &gtfs.TripUpdate_StopTimeEvent{Time: &persisted},
				},
				{ // regular stop already past the persistence window, dropped
					StopSequence: &ss2,
					StopId:       &sid2,
					Arrival:      &gtfs.TripUpdate_StopTimeEvent{Time: &expired},
					Departure:    &gtfs.TripUpdate_StopTimeEvent{Time: &expiredDep},
				},
				{ // departure-only stop after the cutoff, kept
					StopSequence: &ss3,
					StopId:       &sid3,
					Departure:    &gtfs.TripUpdate_StopTimeEvent{Time: &persistedDep},
				},
				{ // departure-only stop already past the persistence window, dropped
					StopSequence: &ss4,
					StopId:       &sid4,
					Departure:    &gtfs.TripUpdate_StopTimeEvent{Time: &expiredDep},
				},
			},
		},
	}
}

func tripUpdateFullyExpired() *gtfs.FeedEntity {
	id := "id"
	tid := "tid"
	rid := "rid"

	expired := int64(315532800)    // 1980, before the persistence cutoff
	expiredDep := int64(315533000) // 1980

	ss1 := uint32(1)
	sid1 := "sid1"
	ss2 := uint32(2)
	sid2 := "sid2"

	return &gtfs.FeedEntity{
		Id: &id,
		TripUpdate: &gtfs.TripUpdate{
			Trip: &gtfs.TripDescriptor{
				TripId:  &tid,
				RouteId: &rid,
			},
			StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{
				{
					StopSequence: &ss1,
					StopId:       &sid1,
					Arrival:      &gtfs.TripUpdate_StopTimeEvent{Time: &expired},
					Departure:    &gtfs.TripUpdate_StopTimeEvent{Time: &expiredDep},
				},
				{
					StopSequence: &ss2,
					StopId:       &sid2,
					Arrival:      &gtfs.TripUpdate_StopTimeEvent{Time: &expired},
					Departure:    &gtfs.TripUpdate_StopTimeEvent{Time: &expiredDep},
				},
			},
		},
	}
}

func newGtfsFeed(e []*gtfs.FeedEntity) *gtfs.FeedMessage {
	v := "version"
	return &gtfs.FeedMessage{
		Header: &gtfs.FeedHeader{
			GtfsRealtimeVersion: &v,
		},
		Entity: e,
	}
}

func tripUpdate() *gtfs.FeedEntity {
	id := "id"
	tid := "tid"
	rid := "rid"
	ss1 := uint32(1)
	sid1 := "sid1"
	ss2 := uint32(2)
	sid2 := "sid2"
	at1 := int64(1601875200)
	dt1 := int64(1601875210)
	at2 := int64(1601875400)
	dt2 := int64(1601875410)

	return &gtfs.FeedEntity{
		Id: &id,
		TripUpdate: &gtfs.TripUpdate{
			Trip: &gtfs.TripDescriptor{
				TripId:  &tid,
				RouteId: &rid,
			},
			StopTimeUpdate: []*gtfs.TripUpdate_StopTimeUpdate{
				{
					StopSequence: &ss1,
					StopId:       &sid1,
					Arrival: &gtfs.TripUpdate_StopTimeEvent{
						Time: &at1,
					},
					Departure: &gtfs.TripUpdate_StopTimeEvent{
						Time: &dt1,
					},
				},
				{
					StopSequence: &ss2,
					StopId:       &sid2,
					Arrival: &gtfs.TripUpdate_StopTimeEvent{
						Time: &at2,
					},
					Departure: &gtfs.TripUpdate_StopTimeEvent{
						Time: &dt2,
					},
				},
			},
		},
	}
}

func vehiclePosition() *gtfs.FeedEntity {
	id := "id"
	tid := "tid"
	rid := "rid"
	vid := "vid"
	lat := float32(47.90258026123047)
	lon := float32(1.8717128038406372)
	bearing := float32(1.3)

	return &gtfs.FeedEntity{
		Id: &id,
		Vehicle: &gtfs.VehiclePosition{
			Trip: &gtfs.TripDescriptor{
				TripId:  &tid,
				RouteId: &rid,
			},
			Vehicle: &gtfs.VehicleDescriptor{
				Id: &vid,
			},
			Position: &gtfs.Position{
				Latitude:  &lat,
				Longitude: &lon,
				Bearing:   &bearing,
			},
		},
	}
}
