package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/remote"

	"crypto/sha1"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/exp/maps"
)

type GtfsRequestCollectorFactory struct{}

func (factory *GtfsRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfLightRemoteCredentials()
}

func (factory *GtfsRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewGtfsRequestCollector(partner)
}

type GtfsRequestCollector struct {
	connector

	remoteCodeSpace string
	origin          string

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
	connector.remoteCodeSpace = connector.Partner().RemoteCodeSpace()
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
		} else if entity.GetAlert() != nil {
			connector.handleAlert(updateEvents, entity.GetAlert(), entity.GetId(), feed.GetHeader().GetTimestamp())
		}
	}

	// logging
	message.Lines = GetModelReferenceSlice(updateEvents.LineRefs)
	message.StopAreas = GetModelReferenceSlice(updateEvents.MonitoringRefs)

	// Broadcast all events
	connector.broadcastUpdateEvents(updateEvents)
	connector.Partner().GtfsStatus(OPERATIONNAL_STATUS_UP)
}

func (connector *GtfsRequestCollector) handleAlert(events *CollectUpdateEvents, a *gtfs.Alert, id string, timestamp uint64) {
	entities := a.GetInformedEntity()
	if len(entities) == 0 {
		logger.Log.Debugf("%d affects for this Alert, skipping message", len(entities))
		return
	}

	event := &model.SituationUpdateEvent{
		RecordedAt:  connector.Clock().Now(),
		VersionedAt: time.Unix(int64(timestamp), 0),
		Origin:      string(connector.Partner().Slug()),
		Progress:    model.SituationProgressPublished,
	}

	// Affects
	for _, entity := range entities {
		affect, collectedRefs, err := model.AffectFromProto(entity,
			connector.remoteCodeSpace,
			connector.Partner().Model(),
		)
		if err != nil {
			logger.Log.Debugf("cannot convert Proto entity: %v", err)
			continue
		}
		maps.Copy(events.MonitoringRefs, collectedRefs.MonitoringRefs)
		maps.Copy(events.LineRefs, collectedRefs.LineRefs)
		event.Affects = append(event.Affects, affect)
	}

	if len(event.Affects) == 0 {
		logger.Log.Debugf("%d affected line/stopArea found for this Alert, skipping message", len(event.Affects))
		return
	}

	// Version
	alert, err := json.Marshal(a)
	if err != nil {
		logger.Log.Debugf("Cannot Marshal gtfs Alert: %v", err)
		return
	}
	hasher := sha1.New()
	hasher.Write(alert)
	data := binary.BigEndian.Uint64(hasher.Sum(nil))
	version := int(data)
	if version < 0 {
		version = -version
	}
	event.Version = version

	// Code
	code := model.NewCode(connector.remoteCodeSpace, id)
	event.SituationCode = code

	// Internal tags
	event.InternalTags = append(event.InternalTags, connector.Partner().CollectSituationsInternalTags()...)

	// ValidityPeriods
	var validityPeriods []*model.TimeRange
	periods := a.GetActivePeriod()
	for _, period := range periods {
		var timePeriod model.TimeRange
		if err := timePeriod.FromProto(period); err != nil {
			logger.Log.Debugf("cannot convert Proto TimeRange: %v", err)
			continue
		}
		validityPeriods = append(validityPeriods, &timePeriod)
	}
	event.ValidityPeriods = validityPeriods

	// Summary
	var s model.SituationTranslatedString
	headerTexts := a.GetHeaderText().GetTranslation()
	if err := s.FromProto(headerTexts); err != nil {
		logger.Log.Debugf("cannot convert Proto HeaderText: %v", err)
	} else {
		event.Summary = &s
	}

	// Description
	var d model.SituationTranslatedString
	descriptionTexts := a.GetDescriptionText().GetTranslation()
	if err := d.FromProto(descriptionTexts); err != nil {
		logger.Log.Debugf("cannot convert Proto TranslatedString: %v", err)
	} else {
		event.Description = &d
	}

	// AlertCause
	var alertCause model.SituationAlertCause
	if err := alertCause.FromProto(a.GetCause()); err != nil {
		logger.Log.Debugf("error in alert cause: %v", err)
	} else {
		event.AlertCause = alertCause
	}

	// Severity
	var severity model.SituationSeverity
	if err := severity.FromProto(a.GetSeverityLevel()); err != nil {
		logger.Log.Debugf("error in severity: %v", err)
	} else {
		event.Severity = severity
	}

	// Condition
	var condition model.SituationCondition
	if err := condition.FromProto(a.GetEffect()); err != nil {
		logger.Log.Debugf("error in condition: %v", err)
	} else {
		consequence := &model.Consequence{
			Condition: condition,
		}
		event.Consequences = append(event.Consequences, consequence)
	}

	events.Situations = append(events.Situations, event)
}

func (connector *GtfsRequestCollector) handleTripUpdate(events *CollectUpdateEvents, t *gtfs.TripUpdate) {
	trip := t.GetTrip()
	if trip == nil {
		return
	}
	vjCode := connector.handleTrip(events, trip) // returns the vj code

	for _, stu := range t.GetStopTimeUpdate() {
		sid := stu.GetStopId()
		svid := fmt.Sprintf("%v-%v", vjCode.Value(), connector.handleStopSequence(stu))
		stopAreaCode := model.NewCode(connector.remoteCodeSpace, sid)

		if sid != "" {
			_, ok := events.StopAreas[sid]
			if !ok {
				// CollectedAlways is false by default
				event := &model.StopAreaUpdateEvent{
					Origin: connector.origin,
					Code:   stopAreaCode,
				}

				events.StopAreas[sid] = event
			}
		}

		_, ok := events.StopVisits[sid][svid]
		if !ok {
			stopVisitCode := model.NewCode(connector.remoteCodeSpace, svid)
			svEvent := &model.StopVisitUpdateEvent{
				Origin:             connector.origin,
				Code:               stopVisitCode,
				StopAreaCode:       stopAreaCode,
				VehicleJourneyCode: vjCode,
				PassageOrder:       connector.handleStopSequence(stu),
				Monitored:          true,
				RecordedAt:         connector.Clock().Now(),
				Schedules:          schedules.NewStopVisitSchedules(),
			}
			svEvent.Schedules.SetSchedule(
				schedules.Expected,
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
	vjCode := connector.handleTrip(events, trip, occupancy) // returns the vj code

	vid := v.GetVehicle().GetId()
	_, ok := events.Vehicles[vid]
	if !ok {
		vCode := model.NewCode(connector.remoteCodeSpace, vid)
		p := v.GetPosition()
		event := &model.VehicleUpdateEvent{
			Code:               vCode,
			StopAreaCode:       model.NewCode(connector.remoteCodeSpace, v.GetStopId()),
			VehicleJourneyCode: vjCode,
			Longitude:          float64(p.GetLongitude()),
			Latitude:           float64(p.GetLatitude()),
			Bearing:            float64(p.GetBearing()),
			Occupancy:          occupancyName(occupancy),
		}

		events.Vehicles[vid] = event
	}
}

// returns the vj code
func (connector *GtfsRequestCollector) handleTrip(events *CollectUpdateEvents, trip *gtfs.TripDescriptor, occupancy ...*gtfs.VehiclePosition_OccupancyStatus) model.Code {
	rid := trip.GetRouteId()
	tid := trip.GetTripId()
	lineCode := model.NewCode(connector.remoteCodeSpace, rid)
	vjCode := model.NewCode(connector.remoteCodeSpace, tid)

	_, ok := events.Lines[rid]
	if !ok {
		// CollectedAlways is false by default
		lineEvent := &model.LineUpdateEvent{
			Origin: connector.origin,
			Code:   lineCode,
		}

		events.Lines[rid] = lineEvent
	}

	_, ok = events.VehicleJourneys[tid]
	if !ok {
		vjEvent := &model.VehicleJourneyUpdateEvent{
			Origin:    connector.origin,
			Code:      vjCode,
			LineCode:  lineCode,
			Monitored: true,
		}
		if len(occupancy) != 0 {
			vjEvent.Occupancy = occupancyName(occupancy[0])
		}

		events.VehicleJourneys[tid] = vjEvent
	}

	return vjCode
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
	for _, e := range events.Situations {
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
		Type:      audit.GTFS_REQUEST,
		Protocol:  "gtfs",
		Direction: "sent",
		Partner:   string(connector.Partner().Slug()),
		Status:    "OK",
	}
}
