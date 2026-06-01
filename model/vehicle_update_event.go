package model

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type VehicleUpdateEvent struct {
	// attributes is a fallback for non-SIRI builders; when SiriXML is set, RawAttributes() delegates to the XML cache.
	attributes         RawAttributes
	SiriXML            *sxml.XMLVehicleActivity
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
	if ue.SiriXML != nil {
		return RawAttributes(ue.SiriXML.RawAttributes())
	}
	if ue.attributes == nil {
		ue.attributes = NewRawAttributes()
	}
	return ue.attributes
}
