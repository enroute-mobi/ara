package sxml

import (
	"time"
)

type XMLEstimatedJourneyVersionFrame struct {
	XMLStructure

	recordedAt time.Time

	estimatedVehicleJourneys []*XMLEstimatedVehicleJourney
}

type XMLEstimatedVehicleJourney struct {
	XMLStructure

	lineRef                string
	directionRef           string
	operatorRef            string
	datedVehicleJourneyRef string
	originRef              string
	destinationRef         string

	estimatedCalls []*XMLCall
	recordedCalls  []*XMLCall
}

func NewXMLEstimatedJourneyVersionFrame(node XMLNode) *XMLEstimatedJourneyVersionFrame {
	ejvf := &XMLEstimatedJourneyVersionFrame{}
	ejvf.node = node
	return ejvf
}

func NewXMLEstimatedVehicleJourney(node XMLNode) *XMLEstimatedVehicleJourney {
	evj := &XMLEstimatedVehicleJourney{}
	evj.node = node
	return evj
}

func (ejvf *XMLEstimatedJourneyVersionFrame) RecordedAt() time.Time {
	if ejvf.recordedAt.IsZero() {
		ejvf.recordedAt = ejvf.findTimeChildContent("RecordedAtTime")
	}
	return ejvf.recordedAt
}

func (ejvf *XMLEstimatedJourneyVersionFrame) EstimatedVehicleJourneys() []*XMLEstimatedVehicleJourney {
	if ejvf.estimatedVehicleJourneys == nil {
		estimatedVehicleJourneys := []*XMLEstimatedVehicleJourney{}
		nodes := ejvf.findNodes("EstimatedVehicleJourney")
		for _, node := range nodes {
			estimatedVehicleJourneys = append(estimatedVehicleJourneys, NewXMLEstimatedVehicleJourney(node))
		}
		ejvf.estimatedVehicleJourneys = estimatedVehicleJourneys
	}
	return ejvf.estimatedVehicleJourneys
}

func (evj *XMLEstimatedVehicleJourney) EstimatedCalls() []*XMLCall {
	if evj.estimatedCalls == nil {
		estimatedCalls := []*XMLCall{}
		nodes := evj.findNodes("EstimatedCall")
		for _, node := range nodes {
			estimatedCalls = append(estimatedCalls, NewXMLCall(node))
		}
		evj.estimatedCalls = estimatedCalls
	}
	return evj.estimatedCalls
}

func (evj *XMLEstimatedVehicleJourney) RecordedCalls() []*XMLCall {
	if evj.recordedCalls == nil {
		recordedCalls := []*XMLCall{}
		nodes := evj.findNodes("RecordedCall")
		for _, node := range nodes {
			recordedCalls = append(recordedCalls, NewXMLCall(node))
		}
		evj.recordedCalls = recordedCalls
	}
	return evj.recordedCalls
}

func (evj *XMLEstimatedVehicleJourney) LineRef() string {
	if evj.lineRef == "" {
		evj.lineRef = evj.findStringChildContent("LineRef")
	}
	return evj.lineRef
}

func (evj *XMLEstimatedVehicleJourney) DirectionRef() string {
	if evj.directionRef == "" {
		evj.directionRef = evj.findStringChildContent("DirectionRef")
	}
	return evj.directionRef
}

func (evj *XMLEstimatedVehicleJourney) OperatorRef() string {
	if evj.operatorRef == "" {
		evj.operatorRef = evj.findStringChildContent("OperatorRef")
	}
	return evj.operatorRef
}

func (evj *XMLEstimatedVehicleJourney) DatedVehicleJourneyRef() string {
	if evj.datedVehicleJourneyRef == "" {
		evj.datedVehicleJourneyRef = evj.findStringChildContent("DatedVehicleJourneyRef")
	}
	return evj.datedVehicleJourneyRef
}

func (evj *XMLEstimatedVehicleJourney) OriginRef() string {
	if evj.originRef == "" {
		evj.originRef = evj.findStringChildContent("OriginRef")
	}
	return evj.originRef
}

func (evj *XMLEstimatedVehicleJourney) DestinationRef() string {
	if evj.destinationRef == "" {
		evj.destinationRef = evj.findStringChildContent("DestinationRef")
	}
	return evj.destinationRef
}
