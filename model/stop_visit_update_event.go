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
	// attributes is a fallback for non-SIRI builders; when SiriXML is set, RawAttributes() delegates to the XML cache.
	attributes         RawAttributes
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

func (ue *StopVisitUpdateEvent) RawAttributes() RawAttributes {
	if ue.SiriXML != nil {
		return RawAttributes(ue.SiriXML.RawAttributes())
	}
	if ue.attributes == nil {
		ue.attributes = NewRawAttributes()
	}
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
