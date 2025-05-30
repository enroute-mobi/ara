package model

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type StopVisitUpdateEvent struct {
	RecordedAt         time.Time
	Schedules          *schedules.StopVisitSchedules
	attributes         Attributes
	SiriXML            *sxml.XMLMonitoredStopVisit
	references         *References
	VehicleJourneyCode Code
	StopAreaCode       Code
	Code               Code
	CodeSpace          string
	DepartureStatus    StopVisitDepartureStatus
	ArrivalStatus      StopVisitArrivalStatus
	DataFrameRef       string
	Origin             string
	PassageOrder       int
	Monitored          bool
	VehicleAtStop      bool
}

func NewStopVisitUpdateEvent() *StopVisitUpdateEvent {
	return &StopVisitUpdateEvent{
		Schedules: schedules.NewStopVisitSchedules(),
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

	ue.attributes.Set(siri_attributes.Delay, ue.SiriXML.Delay())
	ue.attributes.Set(siri_attributes.ActualQuayName, ue.SiriXML.ActualQuayName())
	ue.attributes.Set(siri_attributes.AimedHeadwayInterval, ue.SiriXML.AimedHeadwayInterval())
	ue.attributes.Set(siri_attributes.ArrivalPlatformName, ue.SiriXML.ArrivalPlatformName())
	ue.attributes.Set(siri_attributes.ArrivalProximityText, ue.SiriXML.ArrivalProximityText())
	ue.attributes.Set(siri_attributes.DepartureBoardingActivity, ue.SiriXML.DepartureBoardingActivity())
	ue.attributes.Set(siri_attributes.DeparturePlatformName, ue.SiriXML.DeparturePlatformName())
	ue.attributes.Set(siri_attributes.DestinationDisplay, ue.SiriXML.DestinationDisplay())
	ue.attributes.Set(siri_attributes.DistanceFromStop, ue.SiriXML.DistanceFromStop())
	ue.attributes.Set(siri_attributes.ExpectedHeadwayInterval, ue.SiriXML.ExpectedHeadwayInterval())
	ue.attributes.Set(siri_attributes.NumberOfStopsAway, ue.SiriXML.NumberOfStopsAway())
	ue.attributes.Set(siri_attributes.PlatformTraversal, ue.SiriXML.PlatformTraversal())

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

	ue.references.SetCode(siri_attributes.OperatorRef, NewCode(ue.CodeSpace, ue.SiriXML.OperatorRef()))

	return *ue.references
}
