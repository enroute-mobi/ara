package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
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
		svid := fmt.Sprintf("%v-%v", vjObjectId.Value(), connector.handleStopSequence(stu))
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
				PassageOrder:           connector.handleStopSequence(stu),
				Monitored:              true,
				RecordedAt:             connector.Clock().Now(),
				Schedules:              model.NewStopVisitSchedules(),
			}
			svEvent.Schedules.SetSchedule(
				model.STOP_VISIT_SCHEDULE_EXPECTED,
				time.Unix(stu.GetDeparture().GetTime(), 0),
				time.Unix(stu.GetArrival().GetTime(), 0))

			if connector.hasSkippedScheduleRelationship(stu) {
				svEvent.DepartureStatus = model.STOP_VISIT_DEPARTURE_CANCELLED
				svEvent.ArrivalStatus = model.STOP_VISIT_ARRIVAL_CANCELLED
			}

			if events.StopVisits[sid] == nil {
				events.StopVisits[sid] = make(map[string]*model.StopVisitUpdateEvent)
			}
			events.StopVisits[sid][svid] = svEvent
		}
	}
}

func (connector *GtfsRequestCollector) hasSkippedScheduleRelationship(stu *gtfs.TripUpdate_StopTimeUpdate) bool {
	return stu.GetScheduleRelationship() == gtfs.TripUpdate_StopTimeUpdate_SKIPPED
}

func (connector *GtfsRequestCollector) handleStopSequence(st *gtfs.TripUpdate_StopTimeUpdate) int {
	return int(st.GetStopSequence() + uint32(1))
}

func (connector *GtfsRequestCollector) handleVehicle(events *CollectUpdateEvents, v *gtfs.VehiclePosition) {
	trip := v.GetTrip()
	if trip == nil || v.GetVehicle() == nil {
		return
	}
	occupancy := v.OccupancyStatus
	vjObjectId := connector.handleTrip(events, trip, occupancy) // returns the vj objectid

	vid := v.GetVehicle().GetId()
	_, ok := events.Vehicles[vid]
	if !ok {
		vObjectId := model.NewObjectID(connector.remoteObjectidKind, vid)
		p := v.GetPosition()
		event := &model.VehicleUpdateEvent{
			ObjectId:               vObjectId,
			StopAreaObjectId:       model.NewObjectID(connector.remoteObjectidKind, v.GetStopId()),
			VehicleJourneyObjectId: vjObjectId,
			Longitude:              float64(p.GetLongitude()),
			Latitude:               float64(p.GetLatitude()),
			Bearing:                float64(p.GetBearing()),
			Occupancy:              occupancyName(occupancy),
			OriginFromGtfsRT:       true,
		}

		events.Vehicles[vid] = event
	}
}

// returns the vj objectid
func (connector *GtfsRequestCollector) handleTrip(events *CollectUpdateEvents, trip *gtfs.TripDescriptor, occupancy ...*gtfs.VehiclePosition_OccupancyStatus) model.ObjectID {
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
		if len(occupancy) != 0 {
			vjEvent.Occupancy = occupancyName(occupancy[0])
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

func occupancyName(occupancy *gtfs.VehiclePosition_OccupancyStatus) string {
	if occupancy == nil {
		return model.Undefined
	}
	switch *occupancy {
	case gtfs.VehiclePosition_NO_DATA_AVAILABLE:
		return model.Undefined
	case gtfs.VehiclePosition_EMPTY:
		return model.Empty
	case gtfs.VehiclePosition_MANY_SEATS_AVAILABLE:
		return model.ManySeatsAvailable
	case gtfs.VehiclePosition_FEW_SEATS_AVAILABLE:
		return model.FewSeatsAvailable
	case gtfs.VehiclePosition_STANDING_ROOM_ONLY:
		return model.StandingRoomOnly
	case gtfs.VehiclePosition_CRUSHED_STANDING_ROOM_ONLY:
		return model.CrushedStandingRoomOnly
	case gtfs.VehiclePosition_FULL:
		return model.Full
	case gtfs.VehiclePosition_NOT_ACCEPTING_PASSENGERS:
		return model.NotAcceptingPassengers
	// case gtfs.VehiclePosition_NOT_BOARDABLE:
	// 	return model.Unknown
	default:
		return model.Unknown
	}
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
