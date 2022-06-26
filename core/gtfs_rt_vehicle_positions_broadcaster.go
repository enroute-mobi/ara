package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/cache"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/model"
)

type VehiclePositionBroadcaster struct {
	connector

	vjRemoteObjectidKinds      []string
	vehicleRemoteObjectidKinds []string
	cache                      *cache.CachedItem
}

type VehiclePositionBroadcasterFactory struct{}

func (factory *VehiclePositionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewVehiclePositionBroadcaster(partner)
}

func (factory *VehiclePositionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
}

func NewVehiclePositionBroadcaster(partner *Partner) *VehiclePositionBroadcaster {
	connector := &VehiclePositionBroadcaster{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)
	connector.vjRemoteObjectidKinds = partner.VehicleJourneyRemoteObjectIDKindWithFallback(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)
	connector.vehicleRemoteObjectidKinds = partner.VehicleRemoteObjectIDKindWithFallback(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)
	connector.partner = partner
	connector.cache = cache.NewCachedItem("VehiclePositions", partner.CacheTimeout(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER), nil, func(...interface{}) (interface{}, error) { return connector.handleGtfs() })

	return connector
}

func (connector *VehiclePositionBroadcaster) HandleGtfs(feed *gtfs.FeedMessage) {
	entities, _ := connector.cache.Value()
	feedEntities := entities.([]*gtfs.FeedEntity)

	for i := range feedEntities {
		feed.Entity = append(feed.Entity, feedEntities[i])
	}
}

func (connector *VehiclePositionBroadcaster) handleGtfs() (entities []*gtfs.FeedEntity, err error) {
	vehicles := connector.partner.Model().Vehicles().FindAll()
	linesObjectId := make(map[model.VehicleJourneyId]model.ObjectID)
	trips := make(map[model.VehicleJourneyId]*gtfs.TripDescriptor)

	for i := range vehicles {
		vehicleId, ok := vehicles[i].ObjectIDWithFallback(connector.vehicleRemoteObjectidKinds)
		if !ok {
			continue
		}

		trip, ok := trips[vehicles[i].VehicleJourneyId]
		// If we don't already have a tripUpdate with the VehicleJourney we create one
		if !ok {
			// Fetch all needed models and objectids
			vj, ok := connector.partner.Model().VehicleJourneys().Find(vehicles[i].VehicleJourneyId)
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
				linesObjectId[vehicles[i].VehicleJourneyId] = lineObjectid
			}
			routeId = lineObjectid.Value()

			// Fill the tripDescriptor
			tripId := vjId.Value()
			trip = &gtfs.TripDescriptor{
				TripId:  &tripId,
				RouteId: &routeId,
			}

			// ARA-874
			// // That part is really not optimized, we could cut it if we have performance problems as StartTime is optionnal
			// stopVisits := connector.partner.Model().StopVisits().FindByVehicleJourneyId(vj.Id())
			// if len(stopVisits) > 0 {
			// 	sort.Slice(stopVisits, func(i, j int) bool {
			// 		return stopVisits[i].PassageOrder < stopVisits[j].PassageOrder
			// 	})
			// 	startTime := stopVisits[0].ReferenceDepartureTime()
			// 	if !startTime.IsZero() {
			// 		t := startTime.Format("15:04:05")
			// 		trip.StartTime = &t
			// 	}
			// }

			trips[vehicles[i].VehicleJourneyId] = trip
		}

		vId := vehicleId.Value()
		newId := fmt.Sprintf("vehicle:%v", vId)
		lat := float32(vehicles[i].Latitude)
		lon := float32(vehicles[i].Longitude)
		bearing := float32(vehicles[i].Bearing)
		timestamp := uint64(vehicles[i].RecordedAtTime.Unix())
		occupancy := gtfs.VehiclePosition_OccupancyStatus(vehicles[i].Occupancy)
		feedEntity := &gtfs.FeedEntity{
			Id: &newId,
			Vehicle: &gtfs.VehiclePosition{
				Trip:            trip,
				Vehicle:         &gtfs.VehicleDescriptor{Id: &vId},
				OccupancyStatus: &occupancy,
				Position: &gtfs.Position{
					Latitude:  &lat,
					Longitude: &lon,
					Bearing:   &bearing,
				},
				Timestamp: &timestamp,
			},
		}

		// Fill StopId, but as it's uptionnal we just fill it if we can find it
		sa, ok := connector.partner.Model().StopAreas().Find(vehicles[i].StopAreaId)
		if ok {
			saId, ok := sa.ObjectID(connector.remoteObjectidKind)
			if ok {
				id := saId.Value()
				feedEntity.Vehicle.StopId = &id
			}
		}

		entities = append(entities, feedEntity)
	}
	return
}
