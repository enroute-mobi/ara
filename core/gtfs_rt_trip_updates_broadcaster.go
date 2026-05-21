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

	vjRemoteCodeSpaces      []string
	vehicleRemoteCodeSpaces []string
	cache                   *cache.CachedItem
}

type TripUpdatesBroadcasterFactory struct{}

func (factory *TripUpdatesBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTripUpdatesBroadcaster(partner)
}

func (factory *TripUpdatesBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
}

func NewTripUpdatesBroadcaster(partner *Partner) *TripUpdatesBroadcaster {
	connector := &TripUpdatesBroadcaster{}
	connector.partner = partner

	return connector
}

func (connector *TripUpdatesBroadcaster) Start() {
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace(GTFS_RT_TRIP_UPDATES_BROADCASTER)
	connector.vjRemoteCodeSpaces = connector.partner.VehicleJourneyRemoteCodeSpaceWithFallback(GTFS_RT_TRIP_UPDATES_BROADCASTER)
	connector.vehicleRemoteCodeSpaces = connector.partner.VehicleRemoteCodeSpaceWithFallback(GTFS_RT_TRIP_UPDATES_BROADCASTER)
	connector.cache = cache.NewCachedItem("TripUpdates", connector.partner.CacheTimeout(GTFS_RT_TRIP_UPDATES_BROADCASTER), nil, func(...any) (any, error) { return connector.handleGtfs() })
}

func (connector *TripUpdatesBroadcaster) HandleGtfs(feed *gtfs.FeedMessage) {
	entities, _ := connector.cache.Value()
	feedEntities := entities.([]*gtfs.FeedEntity)

	feed.Entity = append(feed.Entity, feedEntities...)
}

func (connector *TripUpdatesBroadcaster) handleGtfs() (entities []*gtfs.FeedEntity, err error) {
	gtfsStopSequenceOffset := uint32(0)
	if connector.partner.GtfsEnforceStopSequence() {
		gtfsStopSequenceOffset += uint32(1)
	}

	referenceTime := connector.Clock().Now().Add(PAST_STOP_VISITS_MAX_TIME)

	vehicleJourneys := connector.partner.Model().VehicleJourneys().FindAll()
	for i := range vehicleJourneys {
		vjId, ok := vehicleJourneys[i].CodeWithFallback(connector.vjRemoteCodeSpaces)
		if !ok {
			continue
		}

		var routeId string
		l, ok := connector.partner.Model().Lines().Find(vehicleJourneys[i].LineId)
		if !ok {
			continue
		}
		lineCode, ok := l.Code(connector.remoteCodeSpace)
		if !ok {
			continue
		}

		routeId = lineCode.Value()
		tripId := vjId.Value()

		// Fill the tripDescriptor
		tripDescriptor := &gtfs.TripDescriptor{
			TripId:  &tripId,
			RouteId: &routeId,
		}

		if directionId := vehicleJourneys[i].GtfsDirectionId(); directionId != nil {
			tripDescriptor.DirectionId = directionId
		}

		if vehicleJourneys[i].IsCancelled() {
			cancelled := gtfs.TripDescriptor_CANCELED
			tripDescriptor.ScheduleRelationship = &cancelled
		}

		tu := &gtfs.TripUpdate{Trip: tripDescriptor}

		// Fetch the Vehicle Informations
		v := vehicleJourneys[i].Vehicle()
		if v != nil {
			vd := &gtfs.VehicleDescriptor{}

			vehicleId, ok := v.CodeWithFallback(connector.vehicleRemoteCodeSpaces)
			if ok {
				vehicleId := vehicleId.Value()
				vd.Id = &vehicleId
			}

			// The other GTFS fields are a label, the licence plate,
			// and wheelchair accessible, but we don't have anything to fill these

			tu.Vehicle = vd
		}

		// Fill the FeedEntity
		newId := fmt.Sprintf("trip:%v", vjId.Value())
		feedEntity := &gtfs.FeedEntity{
			Id:         &newId,
			TripUpdate: tu,
		}

		stopVisits := connector.partner.Model().StopVisits().FindByVehicleJourneyIdAfterUnsorted(vehicleJourneys[i].Id(), referenceTime)
		sort.Slice(stopVisits, func(i, j int) bool {
			return stopVisits[i].PassageOrder < stopVisits[j].PassageOrder
		})

		var passageOrderOffset int
		if vehicleJourneys[i].AimedStopVisitCount != 0 && vehicleJourneys[i].AimedStopVisitCount > len(stopVisits) {
			passageOrderOffset = vehicleJourneys[i].AimedStopVisitCount - len(stopVisits)
		}

		for i := range stopVisits {
			sa, ok := connector.partner.Model().StopAreas().Find(stopVisits[i].StopAreaId)
			if !ok { // Should never happen
				logger.Log.Debugf("Can't find StopArea %v of StopVisit %v", stopVisits[i].StopAreaId, stopVisits[i].Id())
				continue
			}
			saId, ok := sa.Code(connector.remoteCodeSpace)
			if !ok {
				continue
			}

			stopId := saId.Value()

			// rewrite stopSequence
			stopSequence := uint32(i) + uint32(passageOrderOffset) + gtfsStopSequenceOffset

			stopTimeUpdate := &gtfs.TripUpdate_StopTimeUpdate{
				StopSequence: &stopSequence,
				StopId:       &stopId,
			}

			arrival := &gtfs.TripUpdate_StopTimeEvent{}
			departure := &gtfs.TripUpdate_StopTimeEvent{}

			if a := stopVisits[i].ReferenceArrivalTime(); !a.IsZero() {
				arrivalTime := int64(a.Unix())
				arrival.Time = &arrivalTime
				stopTimeUpdate.Arrival = arrival

			}
			if d := stopVisits[i].ReferenceDepartureTime(); !d.IsZero() {
				departureTime := int64(d.Unix())
				departure.Time = &departureTime
				stopTimeUpdate.Departure = departure
			}

			if stopVisits[i].DepartureStatus == model.STOP_VISIT_DEPARTURE_CANCELLED {
				skipped := gtfs.TripUpdate_StopTimeUpdate_SKIPPED
				stopTimeUpdate.ScheduleRelationship = &skipped
			}
			feedEntity.TripUpdate.StopTimeUpdate = append(feedEntity.TripUpdate.StopTimeUpdate, stopTimeUpdate)
		}

		if len(feedEntity.TripUpdate.StopTimeUpdate) != 0 {
			entities = append(entities, feedEntity)
		}
	}
	return
}
