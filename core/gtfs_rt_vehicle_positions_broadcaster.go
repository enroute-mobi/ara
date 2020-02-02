package core

import (
	"bitbucket.org/enroute-mobi/edwig/model"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
)

type VehiclePositionBroadcaster struct {
	model.ClockConsumer

	BaseConnector
}

type VehiclePositionBroadcasterFactory struct{}

func (factory *VehiclePositionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewVehiclePositionBroadcaster(partner)
}

func (factory *VehiclePositionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	return ok
}

func NewVehiclePositionBroadcaster(partner *Partner) *VehiclePositionBroadcaster {
	connector := &VehiclePositionBroadcaster{}
	connector.partner = partner

	return connector
}

func (vpb *VehiclePositionBroadcaster) HandleGtfs(feed *gtfs.FeedMessage) {
}
