package model

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri"
)

type StopVisitDepartureStatus string

const (
	STOP_VISIT_DEPARTURE_ONTIME    StopVisitDepartureStatus = "onTime"
	STOP_VISIT_DEPARTURE_EARLY     StopVisitDepartureStatus = "early"
	STOP_VISIT_DEPARTURE_DELAYED   StopVisitDepartureStatus = "delayed"
	STOP_VISIT_DEPARTURE_CANCELLED StopVisitDepartureStatus = "cancelled"
	STOP_VISIT_DEPARTURE_NOREPORT  StopVisitDepartureStatus = "noreport"
	STOP_VISIT_DEPARTURE_DEPARTED  StopVisitDepartureStatus = "departed"
	STOP_VISIT_DEPARTURE_UNDEFINED StopVisitDepartureStatus = ""
)

type StopVisitArrivalStatus string

const (
	STOP_VISIT_ARRIVAL_ARRIVED      StopVisitArrivalStatus = "arrived"
	STOP_VISIT_ARRIVAL_ONTIME       StopVisitArrivalStatus = "onTime"
	STOP_VISIT_ARRIVAL_EARLY        StopVisitArrivalStatus = "early"
	STOP_VISIT_ARRIVAL_DELAYED      StopVisitArrivalStatus = "delayed"
	STOP_VISIT_ARRIVAL_CANCELLED    StopVisitArrivalStatus = "cancelled"
	STOP_VISIT_ARRIVAL_NOREPORT     StopVisitArrivalStatus = "noreport"
	STOP_VISIT_ARRIVAL_MISSED       StopVisitArrivalStatus = "missed"
	STOP_VISIT_ARRIVAL_NOT_EXPECTED StopVisitArrivalStatus = "notExpected"
	STOP_VISIT_ARRIVAL_UNDEFINED    StopVisitArrivalStatus = ""
)

type StopVisitUpdateEvent struct {
	Origin string

	ObjectId               ObjectID
	StopAreaObjectId       ObjectID
	VehicleJourneyObjectId ObjectID

	DataFrameRef    string
	PassageOrder    int
	Monitored       bool
	VehicleAtStop   bool
	Schedules       StopVisitSchedules
	DepartureStatus StopVisitDepartureStatus
	ArrivalStatus   StopVisitArrivalStatus
	RecordedAt      time.Time

	ObjectidKind string
	SiriXML      *siri.XMLMonitoredStopVisit
	attributes   Attributes
	references   *References
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
