package core

import (
	"fmt"
	"sort"
	"time"

	"bitbucket.org/enroute-mobi/ara/cache"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/state"
)

const (
	PAST_STOP_VISITS_MAX_TIME = -2 * time.Minute
)

type TripUpdatesBroadcaster struct {
	state.Startable
	connector

	vjRemoteObjectidKinds []string
	cache                 *cache.CachedItem
}

type TripUpdatesBroadcasterFactory struct{}

func (factory *TripUpdatesBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTripUpdatesBroadcaster(partner)
}

func (factory *TripUpdatesBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
}

func NewTripUpdatesBroadcaster(partner *Partner) *TripUpdatesBroadcaster {
	connector := &TripUpdatesBroadcaster{}
	connector.partner = partner

	return connector
}

func (connector *TripUpdatesBroadcaster) Start() {
	connector.remoteObjectidKind = connector.partner.RemoteObjectIDKind(GTFS_RT_TRIP_UPDATES_BROADCASTER)
	connector.vjRemoteObjectidKinds = connector.partner.VehicleJourneyRemoteObjectIDKindWithFallback(GTFS_RT_TRIP_UPDATES_BROADCASTER)
	connector.cache = cache.NewCachedItem("TripUpdates", connector.partner.CacheTimeout(GTFS_RT_TRIP_UPDATES_BROADCASTER), nil, func(...interface{}) (interface{}, error) { return connector.handleGtfs() })
}

func (connector *TripUpdatesBroadcaster) HandleGtfs(feed *gtfs.FeedMessage) {
	entities, _ := connector.cache.Value()
	feedEntities := entities.([]*gtfs.FeedEntity)

	feed.Entity = append(feed.Entity, feedEntities...)
}

func (connector *TripUpdatesBroadcaster) handleGtfs() (entities []*gtfs.FeedEntity, err error) {
	stopVisits := connector.partner.Model().StopVisits().FindAllAfter(connector.Clock().Now().Add(PAST_STOP_VISITS_MAX_TIME))
	linesObjectId := make(map[model.VehicleJourneyId]model.ObjectID)
	feedEntities := make(map[model.VehicleJourneyId]*gtfs.FeedEntity)

	for i := range stopVisits {
		sa, ok := connector.partner.Model().StopAreas().Find(stopVisits[i].StopAreaId)
		if !ok { // Should never happen
			logger.Log.Debugf("Can't find StopArea %v of StopVisit %v", stopVisits[i].StopAreaId, stopVisits[i].Id())
			continue
		}
		saId, ok := sa.ObjectID(connector.remoteObjectidKind)
		if !ok {
			continue
		}

		feedEntity, ok := feedEntities[stopVisits[i].VehicleJourneyId]
		// If we don't already have a tripUpdate with the VehicleJourney we create one
		if !ok {
			// Fetch all needed models and objectids
			vj, ok := connector.partner.Model().VehicleJourneys().Find(stopVisits[i].VehicleJourneyId)
			if !ok {
				continue
			}
			vjId, ok := vj.ObjectIDWithFallback(connector.vjRemoteObjectidKinds)
			if !ok {
				continue
			}

			var routeId string
			lineObjectid, ok := linesObjectId[vj.Id()]
			if !ok {
				l, ok := connector.partner.Model().Lines().Find(vj.LineId)
				if !ok {
					continue
				}
				lineObjectid, ok = l.ObjectID(connector.remoteObjectidKind)
				if !ok {
					continue
				}
				linesObjectId[stopVisits[i].VehicleJourneyId] = lineObjectid
			}
			routeId = lineObjectid.Value()
			tripId := vjId.Value()
			// Fill the tripDescriptor
			tripDescriptor := &gtfs.TripDescriptor{
				TripId:  &tripId,
				RouteId: &routeId,
			}

			// Fill the FeedEntity
			newId := fmt.Sprintf("trip:%v", vjId.Value())
			feedEntity = &gtfs.FeedEntity{
				Id:         &newId,
				TripUpdate: &gtfs.TripUpdate{Trip: tripDescriptor},
			}

			feedEntities[stopVisits[i].VehicleJourneyId] = feedEntity
		}

		stopId := saId.Value()
		stopSequence := connector.gtfsStopSequence(stopVisits[i].PassageOrder)
		arrival := &gtfs.TripUpdate_StopTimeEvent{}
		departure := &gtfs.TripUpdate_StopTimeEvent{}

		if a := stopVisits[i].ReferenceArrivalTime(); !a.IsZero() {
			arrivalTime := int64(a.Unix())
			arrival.Time = &arrivalTime
		}
		if d := stopVisits[i].ReferenceDepartureTime(); !d.IsZero() {
			departureTime := int64(d.Unix())
			departure.Time = &departureTime
		}

		stopTimeUpdate := &gtfs.TripUpdate_StopTimeUpdate{
			StopSequence: &stopSequence,
			StopId:       &stopId,
			Arrival:      arrival,
			Departure:    departure,
		}

		if stopVisits[i].DepartureStatus == model.STOP_VISIT_DEPARTURE_CANCELLED {
			skipped := gtfs.TripUpdate_StopTimeUpdate_SKIPPED
			stopTimeUpdate.ScheduleRelationship = &skipped
		}

		feedEntity.TripUpdate.StopTimeUpdate = append(feedEntity.TripUpdate.StopTimeUpdate, stopTimeUpdate)
	}

	for _, entity := range feedEntities {
		if len(entity.TripUpdate.StopTimeUpdate) == 0 {
			continue
		}
		sort.Slice(entity.TripUpdate.StopTimeUpdate, func(i, j int) bool {
			return *entity.TripUpdate.StopTimeUpdate[i].StopSequence < *entity.TripUpdate.StopTimeUpdate[j].StopSequence
		})
		// ARA-829
		// if entity.TripUpdate.StopTimeUpdate[0].Departure.Time != nil {
		// 	startTime := time.Unix(*entity.TripUpdate.StopTimeUpdate[0].Departure.Time, 0).Format("15:04:05")
		// 	entity.TripUpdate.Trip.StartTime = &startTime
		// }
		entities = append(entities, entity)
	}
	return
}

func (connector *TripUpdatesBroadcaster) gtfsStopSequence(stopSequence int) uint32 {
	return uint32(stopSequence - 1)
}
