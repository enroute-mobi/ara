package model

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type VehicleUpdateEvent struct {
	SiriXML            *sxml.XMLVehicleActivity
	attributes         RawAttributes
	ValidUntilTime     time.Time
	RecordedAt         time.Time
	Code               Code
	StopAreaCode       Code
	VehicleJourneyCode Code
	Occupancy          string
	DriverRef          string
	Origin             string
	Longitude          float64
	Latitude           float64
	Bearing            float64
	Percentage         float64
	LinkDistance       float64
	NextStopPointOrder int
}

func NewVehicleUpdateEvent() *VehicleUpdateEvent {
	return &VehicleUpdateEvent{}
}

func (ue *VehicleUpdateEvent) EventKind() EventKind {
	return VEHICLE_EVENT
}

func (ue *VehicleUpdateEvent) RawAttributes() RawAttributes {
	if ue.attributes != nil {
		return ue.attributes
	}
	ue.attributes = NewRawAttributes()

	if ue.SiriXML == nil {
		return ue.attributes
	}

	ue.attributes.Set(siri_attributes.VehicleActivityNote, ue.SiriXML.VehicleActivityNote())

	return ue.attributes
}
