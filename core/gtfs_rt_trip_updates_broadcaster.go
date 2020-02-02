package core

import (
	"bitbucket.org/enroute-mobi/edwig/model"
	"github.com/MobilityData/gtfs-realtime-bindings/golang/gtfs"
)

type TripUpdatesBroadcaster struct {
	model.ClockConsumer

	BaseConnector
}

type TripUpdatesBroadcasterFactory struct{}

func (factory *TripUpdatesBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTripUpdatesBroadcaster(partner)
}

func (factory *TripUpdatesBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	return ok
}

func NewTripUpdatesBroadcaster(partner *Partner) *TripUpdatesBroadcaster {
	connector := &TripUpdatesBroadcaster{}
	connector.partner = partner

	return connector
}

func (tub *TripUpdatesBroadcaster) HandleGtfs(feed *gtfs.FeedMessage) {
}
