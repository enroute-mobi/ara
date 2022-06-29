package model

import (
	"time"
)

type VehicleUpdateEvent struct {
	Origin string

	ObjectId               ObjectID
	StopAreaObjectId       ObjectID
	VehicleJourneyObjectId ObjectID
	DriverRef              string
	Occupancy              string
	LinkDistance           float64
	Percentage             float64
	Longitude              float64
	Latitude               float64
	Bearing                float64
	ValidUntilTime         time.Time
	RecordedAt             time.Time
}

func NewVehicleUpdateEvent() *VehicleUpdateEvent {
	return &VehicleUpdateEvent{}
}

func (ue *VehicleUpdateEvent) EventKind() EventKind {
	return VEHICLE_EVENT
}
