package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/remote"
)

type GtfsRequestCollectorFactory struct{}

func (factory *GtfsRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLightRemoteCredentials()
}

func (factory *GtfsRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewGtfsRequestCollector(partner)
}

type GtfsRequestCollector struct {
	clock.ClockConsumer

	connector

	remoteObjectidKind string
	origin             string

	ttl        time.Duration
	subscriber UpdateSubscriber
	stop       chan struct{}
}

func NewGtfsRequestCollector(partner *Partner) *GtfsRequestCollector {
	connector := &GtfsRequestCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.subscriber = manager.BroadcastUpdateEvent

	return connector
}

func (connector *GtfsRequestCollector) Start() {
	logger.Log.Debugf("Start GtfsRequestCollector")

	connector.ttl = connector.Partner().GtfsTTL()
	connector.remoteObjectidKind = connector.Partner().RemoteObjectIDKind()
	connector.origin = string(connector.Partner().Slug())
	connector.stop = make(chan struct{})
	go connector.run()
}

func (connector *GtfsRequestCollector) run() {
	c := connector.Clock().After(5 * time.Second)

	for {
		select {
		case <-connector.stop:
			logger.Log.Debugf("gtfs collector routine stop")

			return
		case <-c:
			logger.Log.Debugf("gtfs collector routine routine")

			connector.requestGtfs()

			c = connector.Clock().After(connector.ttl)
		}
	}
}

func (connector *GtfsRequestCollector) Stop() {
	if connector.stop != nil {
		close(connector.stop)
	}
}

func (connector *GtfsRequestCollector) requestGtfs() {
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	feed, err := connector.Partner().HTTPClient().GTFSRequest()
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		s := operationnalStatusFromError(err)
		connector.Partner().GtfsStatus(s)
		message.Status = "Error"
		message.ErrorDetails = fmt.Sprintf("Error while making a GTFS Request: %v", err)
		return
	}

	updateEvents := NewCollectUpdateEvents()

	for _, entity := range feed.GetEntity() {
		if entity.GetTripUpdate() != nil {
			connector.handleTripUpdate(updateEvents, entity.GetTripUpdate())
		} else if entity.GetVehicle() != nil {
			connector.handleVehicle(updateEvents, entity.GetVehicle())
		}
	}

	// Broadcast all events
	connector.broadcastUpdateEvents(updateEvents)
	connector.Partner().GtfsStatus(OPERATIONNAL_STATUS_UP)
}

func (connector *GtfsRequestCollector) handleTripUpdate(events *CollectUpdateEvents, t *gtfs.TripUpdate) {
	trip := t.GetTrip()
	if trip == nil {
		return
	}
	vjObjectId := connector.handleTrip(events, trip) // returns the vj objectid

	for _, stu := range t.GetStopTimeUpdate() {
		sid := stu.GetStopId()
		svid := fmt.Sprintf("%v-%v", vjObjectId.Value(), stu.GetStopSequence())
		stopAreaObjectId := model.NewObjectID(connector.remoteObjectidKind, sid)

		if sid != "" {
			_, ok := events.StopAreas[sid]
			if !ok {
				// CollectedAlways is false by default
				event := &model.StopAreaUpdateEvent{
					Origin:   connector.origin,
					ObjectId: stopAreaObjectId,
				}

				events.StopAreas[sid] = event
			}
		}

		_, ok := events.StopVisits[sid][svid]
		if !ok {
			stopVisitObjectId := model.NewObjectID(connector.remoteObjectidKind, svid)
			svEvent := &model.StopVisitUpdateEvent{
				Origin:                 connector.origin,
				ObjectId:               stopVisitObjectId,
				StopAreaObjectId:       stopAreaObjectId,
				VehicleJourneyObjectId: vjObjectId,
				PassageOrder:           int(stu.GetStopSequence()),
				Monitored:              true,
				RecordedAt:             connector.Clock().Now(),
				Schedules:              model.NewStopVisitSchedules(),
			}
			svEvent.Schedules.SetSchedule(
				model.STOP_VISIT_SCHEDULE_EXPECTED,
				time.Unix(stu.GetDeparture().GetTime(), 0),
				time.Unix(stu.GetArrival().GetTime(), 0))

			if events.StopVisits[sid] == nil {
				events.StopVisits[sid] = make(map[string]*model.StopVisitUpdateEvent)
			}
			events.StopVisits[sid][svid] = svEvent
		}
	}
}

func (connector *GtfsRequestCollector) handleVehicle(events *CollectUpdateEvents, v *gtfs.VehiclePosition) {
	trip := v.GetTrip()
	if trip == nil || v.GetVehicle() == nil {
		return
	}
	vjObjectId := connector.handleTrip(events, trip) // returns the vj objectid

	vid := v.GetVehicle().GetId()
	_, ok := events.Vehicles[vid]
	if !ok {
		vObjectId := model.NewObjectID(connector.remoteObjectidKind, vid)
		p := v.GetPosition()
		event := &model.VehicleUpdateEvent{
			ObjectId:               vObjectId,
			VehicleJourneyObjectId: vjObjectId,
			Longitude:              float64(p.GetLongitude()),
			Latitude:               float64(p.GetLatitude()),
			Bearing:                float64(p.GetBearing()),
		}
		event.Attributes().Set("Occupancy", v.GetOccupancyStatus().String())

		events.Vehicles[vid] = event
	}
}

// returns the vj objectid
func (connector *GtfsRequestCollector) handleTrip(events *CollectUpdateEvents, trip *gtfs.TripDescriptor) model.ObjectID {
	rid := trip.GetRouteId()
	tid := trip.GetTripId()
	lineObjectId := model.NewObjectID(connector.remoteObjectidKind, rid)
	vjObjectId := model.NewObjectID(connector.remoteObjectidKind, tid)

	_, ok := events.Lines[rid]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin:   connector.origin,
			ObjectId: lineObjectId,
		}

		events.Lines[rid] = lineEvent
	}

	_, ok = events.VehicleJourneys[tid]
	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:       connector.origin,
			ObjectId:     vjObjectId,
			LineObjectId: lineObjectId,
			Monitored:    true,
		}

		events.VehicleJourneys[tid] = vjEvent
	}

	return vjObjectId
}

func (connector *GtfsRequestCollector) SetSubscriber(s UpdateSubscriber) {
	connector.subscriber = s
}

func (connector *GtfsRequestCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {
	if connector.subscriber == nil {
		return
	}
	for _, e := range events.StopAreas {
		connector.subscriber(e)
	}
	for _, e := range events.Lines {
		connector.subscriber(e)
	}
	for _, e := range events.VehicleJourneys {
		connector.subscriber(e)
	}
	for _, es := range events.StopVisits { // Stopvisits are map[MonitoringRef]map[ItemIdentifier]event
		for _, e := range es {
			connector.subscriber(e)
		}
	}
	for _, e := range events.Vehicles {
		connector.subscriber(e)
	}
}

func operationnalStatusFromError(err error) OperationnalStatus {
	if _, ok := err.(remote.GtfsError); ok {
		return OPERATIONNAL_STATUS_DOWN
	}
	return OPERATIONNAL_STATUS_UNKNOWN
}

func (connector *GtfsRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "GtfsRequest",
		Protocol:  "gtfs",
		Direction: "sent",
		Partner:   string(connector.Partner().Slug()),
		Status:    "OK",
	}
}
