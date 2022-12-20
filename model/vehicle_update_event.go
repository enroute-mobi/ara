package model

import (
	"time"
)

type VehicleUpdateEvent struct {
	ValidUntilTime         time.Time
	RecordedAt             time.Time
	ObjectId               ObjectID
	StopAreaObjectId       ObjectID
	VehicleJourneyObjectId ObjectID
	DriverRef              string
	Occupancy              string
	Origin                 string
	Longitude              float64
	Latitude               float64
	Bearing                float64
	Percentage             float64
	LinkDistance           float64
}

func NewVehicleUpdateEvent() *VehicleUpdateEvent {
	return &VehicleUpdateEvent{}
}

func (ue *VehicleUpdateEvent) EventKind() EventKind {
	return VEHICLE_EVENT
}
