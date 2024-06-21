package model

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type VehicleJourneyUpdateEvent struct {
	references      *References
	SiriXML         *sxml.XMLMonitoredVehicleJourney
	attributes      Attributes
	Code            Code
	LineCode        Code
	Direction       string
	DestinationName string
	DestinationRef  string
	DirectionType   string
	Occupancy       string
	OriginName      string
	CodeSpace       string
	OriginRef       string
	Origin          string
	Monitored       bool
}

func NewVehicleJourneyUpdateEvent() *VehicleJourneyUpdateEvent {
	return &VehicleJourneyUpdateEvent{}
}

func (ue *VehicleJourneyUpdateEvent) EventKind() EventKind {
	return VEHICLE_JOURNEY_EVENT
}

func (ue *VehicleJourneyUpdateEvent) Attributes() Attributes {
	if ue.attributes != nil {
		return ue.attributes
	}
	ue.attributes = NewAttributes()

	if ue.SiriXML == nil {
		return ue.attributes
	}

	ue.attributes.Set(siri_attributes.Delay, ue.SiriXML.Delay())
	ue.attributes.Set(siri_attributes.Bearing, ue.SiriXML.Bearing())
	ue.attributes.Set(siri_attributes.InPanic, ue.SiriXML.InPanic())
	ue.attributes.Set(siri_attributes.InCongestion, ue.SiriXML.InCongestion())
	ue.attributes.Set(siri_attributes.SituationRef, ue.SiriXML.SituationRef())
	ue.attributes.Set(siri_attributes.DirectionName, ue.SiriXML.DirectionName())
	ue.attributes.Set(siri_attributes.FirstOrLastJourney, ue.SiriXML.FirstOrLastJourney())
	ue.attributes.Set(siri_attributes.HeadwayService, ue.SiriXML.HeadwayService())
	ue.attributes.Set(siri_attributes.JourneyNote, ue.SiriXML.JourneyNote())
	ue.attributes.Set(siri_attributes.JourneyPatternName, ue.SiriXML.JourneyPatternName())
	ue.attributes.Set(siri_attributes.MonitoringError, ue.SiriXML.MonitoringError())
	ue.attributes.Set(siri_attributes.OriginAimedDepartureTime, ue.SiriXML.OriginAimedDepartureTime())
	ue.attributes.Set(siri_attributes.DestinationAimedArrivalTime, ue.SiriXML.DestinationAimedArrivalTime())
	ue.attributes.Set(siri_attributes.ProductCategoryRef, ue.SiriXML.ProductCategoryRef())
	ue.attributes.Set(siri_attributes.ServiceFeatureRef, ue.SiriXML.ServiceFeatureRef())
	ue.attributes.Set(siri_attributes.TrainNumberRef, ue.SiriXML.TrainNumberRef())
	ue.attributes.Set(siri_attributes.VehicleFeatureRef, ue.SiriXML.VehicleFeatureRef())
	ue.attributes.Set(siri_attributes.VehicleMode, ue.SiriXML.VehicleMode())
	ue.attributes.Set(siri_attributes.ViaPlaceName, ue.SiriXML.ViaPlaceName())
	ue.attributes.Set(siri_attributes.VehicleJourneyName, ue.SiriXML.VehicleJourneyName())

	// filter empty attributes
	for k := range ue.attributes {
		if ue.attributes[k] == "" {
			delete(ue.attributes, k)
		}
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
