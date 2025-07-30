package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/cache"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/state"
)

type VehiclePositionBroadcaster struct {
	state.Startable
	connector

	vjRemoteCodeSpaces      []string
	vehicleRemoteCodeSpaces []string
	cache                   *cache.CachedItem
}

type VehiclePositionBroadcasterFactory struct{}

func (factory *VehiclePositionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewVehiclePositionBroadcaster(partner)
}

func (factory *VehiclePositionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
}

func NewVehiclePositionBroadcaster(partner *Partner) *VehiclePositionBroadcaster {
	connector := &VehiclePositionBroadcaster{}
	connector.partner = partner

	return connector
}

func (connector *VehiclePositionBroadcaster) Start() {
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)
	connector.vjRemoteCodeSpaces = connector.partner.VehicleJourneyRemoteCodeSpaceWithFallback(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)
	connector.vehicleRemoteCodeSpaces = connector.partner.VehicleRemoteCodeSpaceWithFallback(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)
	connector.cache = cache.NewCachedItem("VehiclePositions", connector.partner.CacheTimeout(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER), nil, func(...interface{}) (interface{}, error) { return connector.handleGtfs() })
}

func (connector *VehiclePositionBroadcaster) HandleGtfs(feed *gtfs.FeedMessage) {
	entities, _ := connector.cache.Value()
	feedEntities := entities.([]*gtfs.FeedEntity)

	feed.Entity = append(feed.Entity, feedEntities...)

}

func (connector *VehiclePositionBroadcaster) handleGtfs() (entities []*gtfs.FeedEntity, err error) {
	vehicles := connector.partner.Model().Vehicles().FindAll()
	linesCode := make(map[model.VehicleJourneyId]model.Code)
	trips := make(map[model.VehicleJourneyId]*gtfs.TripDescriptor)

	for i := range vehicles {
		vehicleId, ok := vehicles[i].CodeWithFallback(connector.vehicleRemoteCodeSpaces)
		if !ok {
			continue
		}

		trip, ok := trips[vehicles[i].VehicleJourneyId]
		// If we don't already have a tripUpdate with the VehicleJourney we create one
		if !ok {
			// Fetch all needed models and codes
			vj, ok := connector.partner.Model().VehicleJourneys().Find(vehicles[i].VehicleJourneyId)
			if !ok {
				continue
			}
			vjId, ok := vj.CodeWithFallback(connector.vjRemoteCodeSpaces)
			if !ok {
				continue
			}

			var routeId string
			lineCode, ok := linesCode[vj.Id()]
			if !ok {
				l, ok := connector.partner.Model().Lines().Find(vj.LineId)
				if !ok {
					continue
				}
				lineCode, ok = l.Code(connector.remoteCodeSpace)
				if !ok {
					continue
				}
				linesCode[vehicles[i].VehicleJourneyId] = lineCode
			}
			routeId = lineCode.Value()

			// Fill the tripDescriptor
			tripId := vjId.Value()
			trip = &gtfs.TripDescriptor{
				TripId:  &tripId,
				RouteId: &routeId,
			}

			if directionId := vj.GtfsDirectionId(); directionId != nil {
				trip.DirectionId = directionId
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
		feedEntity := &gtfs.FeedEntity{
			Id: &newId,
			Vehicle: &gtfs.VehiclePosition{
				Trip:            trip,
				Vehicle:         &gtfs.VehicleDescriptor{Id: &vId},
				OccupancyStatus: occupancyCode(vehicles[i].Occupancy),
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
			saId, ok := sa.Code(connector.remoteCodeSpace)
			if ok {
				id := saId.Value()
				feedEntity.Vehicle.StopId = &id
			}
		}

		entities = append(entities, feedEntity)
	}
	return
}

func occupancyCode(occupancy string) *gtfs.VehiclePosition_OccupancyStatus {
	var o gtfs.VehiclePosition_OccupancyStatus
	switch occupancy {
	case model.Unknown:
		o = gtfs.VehiclePosition_NO_DATA_AVAILABLE
	case model.Empty:
		o = gtfs.VehiclePosition_EMPTY
	case model.ManySeatsAvailable:
		o = gtfs.VehiclePosition_MANY_SEATS_AVAILABLE
	case model.FewSeatsAvailable:
		o = gtfs.VehiclePosition_FEW_SEATS_AVAILABLE
	case model.StandingRoomOnly:
		o = gtfs.VehiclePosition_STANDING_ROOM_ONLY
	case model.CrushedStandingRoomOnly:
		o = gtfs.VehiclePosition_CRUSHED_STANDING_ROOM_ONLY
	case model.Full:
		o = gtfs.VehiclePosition_FULL
	case model.NotAcceptingPassengers:
		o = gtfs.VehiclePosition_NOT_ACCEPTING_PASSENGERS
	default:
		return nil
	}
	return &o
}
