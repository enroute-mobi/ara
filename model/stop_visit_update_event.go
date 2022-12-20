package model

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type StopVisitUpdateEvent struct {
	RecordedAt             time.Time
	Schedules              *StopVisitSchedules
	attributes             Attributes
	SiriXML                *sxml.XMLMonitoredStopVisit
	references             *References
	VehicleJourneyObjectId ObjectID
	StopAreaObjectId       ObjectID
	ObjectId               ObjectID
	ObjectidKind           string
	DepartureStatus        StopVisitDepartureStatus
	ArrivalStatus          StopVisitArrivalStatus
	DataFrameRef           string
	Origin                 string
	PassageOrder           int
	Monitored              bool
	VehicleAtStop          bool
}

func NewStopVisitUpdateEvent() *StopVisitUpdateEvent {
	return &StopVisitUpdateEvent{
		Schedules: NewStopVisitSchedules(),
	}
}

func (ue *StopVisitUpdateEvent) EventKind() EventKind {
	return STOP_VISIT_EVENT
}

func (ue *StopVisitUpdateEvent) Attributes() Attributes {
	if ue.attributes != nil {
		return ue.attributes
	}
	ue.attributes = NewAttributes()

	if ue.SiriXML == nil {
		return ue.attributes
	}

	ue.attributes.Set("Delay", ue.SiriXML.Delay())
	ue.attributes.Set("ActualQuayName", ue.SiriXML.ActualQuayName())
	ue.attributes.Set("AimedHeadwayInterval", ue.SiriXML.AimedHeadwayInterval())
	ue.attributes.Set("ArrivalPlatformName", ue.SiriXML.ArrivalPlatformName())
	ue.attributes.Set("ArrivalProximyTest", ue.SiriXML.ArrivalProximyTest())
	ue.attributes.Set("DepartureBoardingActivity", ue.SiriXML.DepartureBoardingActivity())
	ue.attributes.Set("DeparturePlatformName", ue.SiriXML.DeparturePlatformName())
	ue.attributes.Set("DestinationDisplay", ue.SiriXML.DestinationDisplay())
	ue.attributes.Set("DistanceFromStop", ue.SiriXML.DistanceFromStop())
	ue.attributes.Set("ExpectedHeadwayInterval", ue.SiriXML.ExpectedHeadwayInterval())
	ue.attributes.Set("NumberOfStopsAway", ue.SiriXML.NumberOfStopsAway())
	ue.attributes.Set("PlatformTraversal", ue.SiriXML.PlatformTraversal())

	return ue.attributes
}

func (ue *StopVisitUpdateEvent) References() References {
	if ue.references != nil {
		return *ue.references
	}
	refs := NewReferences()
	ue.references = &refs

	if ue.SiriXML == nil {
		return *ue.references
	}

	ue.references.SetObjectId("OperatorRef", NewObjectID(ue.ObjectidKind, ue.SiriXML.OperatorRef()))

	return *ue.references
}
