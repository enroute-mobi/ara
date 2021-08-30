package core

import (
	"fmt"
	"sort"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/cache"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
)

const (
	PAST_STOP_VISITS_MAX_TIME = -2 * time.Minute
)

type TripUpdatesBroadcaster struct {
	clock.ClockConsumer

	connector

	cache *cache.CachedItem
}

type TripUpdatesBroadcasterFactory struct{}

func (factory *TripUpdatesBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTripUpdatesBroadcaster(partner)
}

func (factory *TripUpdatesBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
}

func NewTripUpdatesBroadcaster(partner *Partner) *TripUpdatesBroadcaster {
	connector := &TripUpdatesBroadcaster{}
	connector.partner = partner
	connector.cache = cache.NewCachedItem("TripUpdates", partner.CacheTimeout(GTFS_RT_TRIP_UPDATES_BROADCASTER), nil, func(...interface{}) (interface{}, error) { return connector.handleGtfs() })

	return connector
}

func (connector *TripUpdatesBroadcaster) HandleGtfs(feed *gtfs.FeedMessage, logStashEvent audit.LogStashEvent) {
	entities, _ := connector.cache.Value()
	feedEntities := entities.([]*gtfs.FeedEntity)

	for i := range feedEntities {
		feed.Entity = append(feed.Entity, feedEntities[i])
	}
	logStashEvent["trip_update_quantity"] = strconv.Itoa(len(feedEntities))
}

func (connector *TripUpdatesBroadcaster) handleGtfs() (entities []*gtfs.FeedEntity, err error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	stopVisits := tx.Model().StopVisits().FindAllAfter(connector.Clock().Now().Add(PAST_STOP_VISITS_MAX_TIME))
	linesObjectId := make(map[model.VehicleJourneyId]model.ObjectID)
	feedEntities := make(map[model.VehicleJourneyId]*gtfs.FeedEntity)

	objectidKind := connector.partner.RemoteObjectIDKind(GTFS_RT_TRIP_UPDATES_BROADCASTER)

	for i := range stopVisits {
		sa, ok := tx.Model().StopAreas().Find(stopVisits[i].StopAreaId)
		if !ok { // Should never happen
			logger.Log.Debugf("Can't find StopArea %v of StopVisit %v", stopVisits[i].StopAreaId, stopVisits[i].Id())
			continue
		}
		saId, ok := sa.ObjectID(objectidKind)
		if !ok {
			continue
		}

		feedEntity, ok := feedEntities[stopVisits[i].VehicleJourneyId]
		// If we don't already have a tripUpdate with the VehicleJourney we create one
		if !ok {
			// Fetch all needed models and objectids
			vj, ok := tx.Model().VehicleJourneys().Find(stopVisits[i].VehicleJourneyId)
			if !ok {
				continue
			}
			vjId, ok := vj.ObjectID(objectidKind)
			if !ok {
				continue
			}

			var routeId string
			lineObjectid, ok := linesObjectId[vj.Id()]
			if !ok {
				l, ok := tx.Model().Lines().Find(vj.LineId)
				if !ok {
					continue
				}
				lineObjectid, ok = l.ObjectID(objectidKind)
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
		stopSequence := uint32(stopVisits[i].PassageOrder)
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
