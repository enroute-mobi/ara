package model

import "time"

type VehicleUpdateEvent struct {
	Origin string

	ObjectId               ObjectID
	StopAreaObjectId       ObjectID
	VehicleJourneyObjectId ObjectID
	DriverRef              string
	LinkDistance           float64
	Percentage             float64
	Longitude              float64
	Latitude               float64
	Bearing                float64
	ValidUntilTime         time.Time
	RecordedAt             time.Time

	attributes Attributes
}

func NewVehicleUpdateEvent() *VehicleUpdateEvent {
	return &VehicleUpdateEvent{
		attributes: NewAttributes(),
	}
}

func (ue *VehicleUpdateEvent) EventKind() EventKind {
	return VEHICLE_EVENT
}

func (ue *VehicleUpdateEvent) Attributes() Attributes {
	if ue.attributes == nil {
		ue.attributes = NewAttributes()
	}
	return ue.attributes
}
