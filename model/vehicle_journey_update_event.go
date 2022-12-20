package model

import (
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type VehicleJourneyUpdateEvent struct {
	references      *References
	SiriXML         *sxml.XMLMonitoredVehicleJourney
	attributes      Attributes
	ObjectId        ObjectID
	LineObjectId    ObjectID
	Direction       string
	DestinationName string
	DestinationRef  string
	DirectionType   string
	Occupancy       string
	OriginName      string
	ObjectidKind    string
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

	ue.references.SetObjectId("PlaceRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.PlaceRef()))
	ue.references.SetObjectId("JourneyPatternRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.JourneyPatternRef()))
	ue.references.SetObjectId("RouteRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.RouteRef()))
	return *ue.references
}
