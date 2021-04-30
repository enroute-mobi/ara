package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/cache"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
)

type VehiclePositionBroadcaster struct {
	clock.ClockConsumer

	BaseConnector

	cache *cache.CachedItem

	referenceGenerator *IdentifierGenerator
}

type VehiclePositionBroadcasterFactory struct{}

func (factory *VehiclePositionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewVehiclePositionBroadcaster(partner)
}

func (factory *VehiclePositionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
}

func NewVehiclePositionBroadcaster(partner *Partner) *VehiclePositionBroadcaster {
	connector := &VehiclePositionBroadcaster{}
	connector.partner = partner
	connector.referenceGenerator = partner.IdentifierGeneratorWithDefault("reference_identifier", "%{objectid}")
	connector.cache = cache.NewCachedItem("VehiclePositions", partner.CacheTimeout(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER), nil, func(...interface{}) (interface{}, error) { return connector.handleGtfs() })

	return connector
}

func (connector *VehiclePositionBroadcaster) HandleGtfs(feed *gtfs.FeedMessage, logStashEvent audit.LogStashEvent) {
	entities, _ := connector.cache.Value()
	feedEntities := entities.([]*gtfs.FeedEntity)

	for i := range feedEntities {
		feed.Entity = append(feed.Entity, feedEntities[i])
	}
	logStashEvent["vehicle_position_quantity"] = strconv.Itoa(len(feedEntities))
}

func (connector *VehiclePositionBroadcaster) handleGtfs() (entities []*gtfs.FeedEntity, err error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	vehicles := tx.Model().Vehicles().FindAll()
	linesObjectId := make(map[model.VehicleJourneyId]model.ObjectID)
	trips := make(map[model.VehicleJourneyId]*gtfs.TripDescriptor)

	objectidKind := connector.partner.RemoteObjectIDKind(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)
	vehicleObjectidKind := connector.partner.VehicleRemoteObjectIDKind(GTFS_RT_VEHICLE_POSITIONS_BROADCASTER)

	for i := range vehicles {
		vehicleId, ok := vehicles[i].ObjectID(vehicleObjectidKind)
		if !ok {
			continue
		}

		trip, ok := trips[vehicles[i].VehicleJourneyId]
		// If we don't already have a tripUpdate with the VehicleJourney we create one
		if !ok {
			// Fetch all needed models and objectids
			vj, ok := tx.Model().VehicleJourneys().Find(vehicles[i].VehicleJourneyId)
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
				linesObjectId[vehicles[i].VehicleJourneyId] = lineObjectid
			}
			routeId = connector.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "Line", ObjectId: lineObjectid.Value()})

			// Fill the tripDescriptor
			tripId := connector.referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "VehicleJourney", ObjectId: vjId.Value()})
			trip = &gtfs.TripDescriptor{
				TripId:  &tripId,
				RouteId: &routeId,
			}

			// ARA-874
			// // That part is really not optimized, we could cut it if we have performance problems as StartTime is optionnal
			// stopVisits := tx.Model().StopVisits().FindByVehicleJourneyId(vj.Id())
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
		feedEntity := &gtfs.FeedEntity{
			Id: &newId,
			Vehicle: &gtfs.VehiclePosition{
				Trip:    trip,
				Vehicle: &gtfs.VehicleDescriptor{Id: &vId},
				Position: &gtfs.Position{
					Latitude:  &lat,
					Longitude: &lon,
					Bearing:   &bearing,
				},
			},
		}

		entities = append(entities, feedEntity)
	}
	return
}
