package model

import (
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type VehicleJourneyUpdateEvent struct {
	Origin string

	ObjectId        ObjectID
	LineObjectId    ObjectID
	OriginRef       string
	OriginName      string
	DestinationRef  string
	DestinationName string
	Direction       string
	Occupancy       string
	Monitored       bool

	ObjectidKind string
	SiriXML      *sxml.XMLMonitoredVehicleJourney
	attributes   Attributes
	references   *References
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

	ue.attributes.Set("Delay", ue.SiriXML.Delay())
	ue.attributes.Set("Bearing", ue.SiriXML.Bearing())
	ue.attributes.Set("InPanic", ue.SiriXML.InPanic())
	ue.attributes.Set("InCongestion", ue.SiriXML.InCongestion())
	ue.attributes.Set("SituationRef", ue.SiriXML.SituationRef())
	ue.attributes.Set("DirectionName", ue.SiriXML.DirectionName())
	ue.attributes.Set("FirstOrLastJourney", ue.SiriXML.FirstOrLastJourney())
	ue.attributes.Set("HeadwayService", ue.SiriXML.HeadwayService())
	ue.attributes.Set("JourneyNote", ue.SiriXML.JourneyNote())
	ue.attributes.Set("JourneyPatternName", ue.SiriXML.JourneyPatternName())
	ue.attributes.Set("MonitoringError", ue.SiriXML.MonitoringError())
	ue.attributes.Set("OriginAimedDepartureTime", ue.SiriXML.OriginAimedDepartureTime())
	ue.attributes.Set("DestinationAimedArrivalTime", ue.SiriXML.DestinationAimedArrivalTime())
	ue.attributes.Set("ProductCategoryRef", ue.SiriXML.ProductCategoryRef())
	ue.attributes.Set("ServiceFeatureRef", ue.SiriXML.ServiceFeatureRef())
	ue.attributes.Set("TrainNumberRef", ue.SiriXML.TrainNumberRef())
	ue.attributes.Set("VehicleFeature", ue.SiriXML.VehicleFeature())
	ue.attributes.Set("VehicleMode", ue.SiriXML.VehicleMode())
	ue.attributes.Set("ViaPlaceName", ue.SiriXML.ViaPlaceName())
	ue.attributes.Set("VehicleJourneyName", ue.SiriXML.VehicleJourneyName())
	ue.attributes.Set("DirectionRef", ue.SiriXML.DirectionRef())
	ue.attributes.Set("DestinationName", ue.SiriXML.DestinationName())
	ue.attributes.Set("OriginName", ue.SiriXML.OriginName())

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

	ue.references.SetObjectId("PlaceRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.PlaceRef()))
	ue.references.SetObjectId("JourneyPatternRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.JourneyPatternRef()))
	ue.references.SetObjectId("RouteRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.RouteRef()))
	ue.references.SetObjectId("DestinationRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.DestinationRef()))
	ue.references.SetObjectId("OriginRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.OriginRef()))
	return *ue.references
}
