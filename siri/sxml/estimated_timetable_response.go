package sxml

import (
	"time"
)

type XMLEstimatedJourneyVersionFrame struct {
	XMLStructure

	lineRef                string
	directionRef           string
	operatorRef            string
	datedVehicleJourneyRef string
	originRef              string
	destinationRef         string

	recordedAt time.Time

	estimatedCalls []*XMLCall
	recordedCalls  []*XMLCall
}

func NewXMLEstimatedJourneyVersionFrame(node XMLNode) *XMLEstimatedJourneyVersionFrame {
	estimatedJourney := &XMLEstimatedJourneyVersionFrame{}
	estimatedJourney.node = node
	return estimatedJourney
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) EstimatedCalls() []*XMLCall {
	if estimatedJourney.estimatedCalls == nil {
		estimatedCalls := []*XMLCall{}
		nodes := estimatedJourney.findNodes("EstimatedCall")
		for _, node := range nodes {
			estimatedCalls = append(estimatedCalls, NewXMLCall(node))
		}
		estimatedJourney.estimatedCalls = estimatedCalls
	}
	return estimatedJourney.estimatedCalls
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) RecordedCalls() []*XMLCall {
	if estimatedJourney.recordedCalls == nil {
		recordedCalls := []*XMLCall{}
		nodes := estimatedJourney.findNodes("RecordedCall")
		for _, node := range nodes {
			recordedCalls = append(recordedCalls, NewXMLCall(node))
		}
		estimatedJourney.recordedCalls = recordedCalls
	}
	return estimatedJourney.recordedCalls
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) LineRef() string {
	if estimatedJourney.lineRef == "" {
		estimatedJourney.lineRef = estimatedJourney.findStringChildContent("LineRef")
	}
	return estimatedJourney.lineRef
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) DirectionRef() string {
	if estimatedJourney.directionRef == "" {
		estimatedJourney.directionRef = estimatedJourney.findStringChildContent("DirectionRef")
	}
	return estimatedJourney.directionRef
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) OperatorRef() string {
	if estimatedJourney.operatorRef == "" {
		estimatedJourney.operatorRef = estimatedJourney.findStringChildContent("OperatorRef")
	}
	return estimatedJourney.operatorRef
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) DatedVehicleJourneyRef() string {
	if estimatedJourney.datedVehicleJourneyRef == "" {
		estimatedJourney.datedVehicleJourneyRef = estimatedJourney.findStringChildContent("DatedVehicleJourneyRef")
	}
	return estimatedJourney.datedVehicleJourneyRef
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) OriginRef() string {
	if estimatedJourney.originRef == "" {
		estimatedJourney.originRef = estimatedJourney.findStringChildContent("OriginRef")
	}
	return estimatedJourney.originRef
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) DestinationRef() string {
	if estimatedJourney.destinationRef == "" {
		estimatedJourney.destinationRef = estimatedJourney.findStringChildContent("DestinationRef")
	}
	return estimatedJourney.destinationRef
}

func (estimatedJourney *XMLEstimatedJourneyVersionFrame) RecordedAt() time.Time {
	if estimatedJourney.recordedAt.IsZero() {
		estimatedJourney.recordedAt = estimatedJourney.findTimeChildContent("RecordedAtTime")
	}
	return estimatedJourney.recordedAt
}
