package model

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type VehicleJourneyUpdateEvent struct {
	references            *References
	// attributes is a fallback for non-SIRI builders; when SiriXML is set, RawAttributes() delegates to the XML cache.
	attributes            RawAttributes
	SiriXML               *sxml.XMLMonitoredVehicleJourney
	Cancellation          bool
	Code                  Code
	LineCode              Code
	Direction             string
	DestinationName       string
	DestinationRef        string
	DirectionType         string
	Occupancy             string
	OriginName            string
	CodeSpace             string
	OriginRef             string
	Origin                string
	Monitored             bool
	FromVehicleMonitoring bool
}

func NewVehicleJourneyUpdateEvent() *VehicleJourneyUpdateEvent {
	return &VehicleJourneyUpdateEvent{}
}

func (ue *VehicleJourneyUpdateEvent) EventKind() EventKind {
	return VEHICLE_JOURNEY_EVENT
}

func (ue *VehicleJourneyUpdateEvent) RawAttributes() RawAttributes {
	if ue.SiriXML != nil {
		return RawAttributes(ue.SiriXML.RawAttributes())
	}
	if ue.attributes == nil {
		ue.attributes = NewRawAttributes()
	}
	return ue.attributes
}

func (ue *VehicleJourneyUpdateEvent) References() References {
	if ue.references != nil {
		return *ue.references
	}
	refs := NewReferences()
	ue.references = &refs

	if ue.SiriXML == nil {
		return *ue.references
	}

	ue.references.SetCode(siri_attributes.PlaceRef, NewCode(ue.CodeSpace, ue.SiriXML.PlaceRef()))
	ue.references.SetCode(siri_attributes.JourneyPatternRef, NewCode(ue.CodeSpace, ue.SiriXML.JourneyPatternRef()))
	ue.references.SetCode(siri_attributes.RouteRef, NewCode(ue.CodeSpace, ue.SiriXML.RouteRef()))
	return *ue.references
}
