package model

type VehicleJourneyUpdateEvent struct {
	Origin string

	ObjectId        ObjectID
	LineObjectId    ObjectID
	OriginRef       string
	OriginName      string
	DestinationRef  string
	DestinationName string
	Direction       string
}

func NewVehicleJourneyUpdateEvent() *VehicleJourneyUpdateEvent {
	return &VehicleJourneyUpdateEvent{}
}

func (ue *VehicleJourneyUpdateEvent) EventKind() EventKind {
	return VEHICLE_JOURNEY_EVENT
}
