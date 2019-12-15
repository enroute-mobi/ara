package model

type VehicleUpdateEvent struct {
	ObjectId               ObjectID
	VehicleJourneyObjectId ObjectID
	Longitude              float64
	Latitude               float64
	Bearing                float64
}

func NewVehicleUpdateEvent() *VehicleUpdateEvent {
	return &VehicleUpdateEvent{}
}

func (ue *VehicleUpdateEvent) EventKind() EventKind {
	return VEHICLE_EVENT
}
