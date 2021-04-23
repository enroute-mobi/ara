package core

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
	"github.com/golang/protobuf/proto"
)

func collectGtfs(t *testing.T, feed *gtfs.FeedMessage, fakeBroadcast bool) ([]model.UpdateEvent, *Partner) {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data, _ := proto.Marshal(feed)
		var buffer bytes.Buffer
		buffer.Write(data)
		// data := d.([]byte)
		w.Header().Set("Content-Type", "application/x-protobuf")
		w.Write(buffer.Bytes())
	}))
	defer ts.Close()

	// Create a GtfsRequestCollector
	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	// referential.collectManager = NewTestCollectManager()
	referentials.Save(referential)
	partners := referential.partners
	partner := partners.New("slug")
	partner.SetSettingsDefinition(map[string]string{
		"remote_url":           ts.URL,
		"remote_objectid_kind": "test_kind",
	})
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
	if partner.alternativeStatusCheck.Status != OPERATIONNAL_STATUS_UP {
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
	partners := createTestPartnerManager()
	partner := partners.New("slug")
	partner.SetSettingsDefinition(map[string]string{
		"remote_url":           ts.URL,
		"remote_objectid_kind": "test_kind",
	})
	partners.Save(partner)

	gtfsCollector := NewGtfsRequestCollector(partner)
	gtfsCollector.SetClock(clock.NewFakeClock())
	gtfsCollector.requestGtfs()

	if partner.alternativeStatusCheck.LastCheck != clock.FAKE_CLOCK_INITIAL_DATE {
		t.Errorf("Partner alternative status time should be updated, got: %v", partner.alternativeStatusCheck.LastCheck)
	}
	if partner.alternativeStatusCheck.Status != OPERATIONNAL_STATUS_DOWN {
		t.Errorf("Partner alternative status status should be updated, got: %v", partner.alternativeStatusCheck.Status)
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
				&gtfs.TripUpdate_StopTimeUpdate{
					StopSequence: &ss1,
					StopId:       &sid1,
					Arrival: &gtfs.TripUpdate_StopTimeEvent{
						Time: &at1,
					},
					Departure: &gtfs.TripUpdate_StopTimeEvent{
						Time: &dt1,
					},
				},
				&gtfs.TripUpdate_StopTimeUpdate{
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
