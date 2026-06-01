package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLEstimatedJourneyVersionFrame struct {
	XMLStructure

	recordedAt *time.Time

	estimatedVehicleJourneys []*XMLEstimatedVehicleJourney
}

type XMLEstimatedVehicleJourney struct {
	XMLStructure

	cancellation           Bool
	datedVehicleJourneyRef *string
	destinationRef         *string
	destinationName        *string
	originName             *string
	directionRef           *string
	operatorRef            *string
	originRef              *string
	lineRef                *string

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
	if ejvf.recordedAt == nil {
		t := ejvf.findTimeChildContent(siri_attributes.RecordedAtTime)
		ejvf.recordedAt = &t
	}
	return *ejvf.recordedAt
}

func (ejvf *XMLEstimatedJourneyVersionFrame) EstimatedVehicleJourneys() []*XMLEstimatedVehicleJourney {
	if ejvf.estimatedVehicleJourneys == nil {
		estimatedVehicleJourneys := []*XMLEstimatedVehicleJourney{}
		nodes := ejvf.findNodes(siri_attributes.EstimatedVehicleJourney)
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
		nodes := evj.findNodes(siri_attributes.EstimatedCall)
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
		nodes := evj.findNodes(siri_attributes.RecordedCall)
		for _, node := range nodes {
			recordedCalls = append(recordedCalls, NewXMLCall(node))
		}
		evj.recordedCalls = recordedCalls
	}
	return evj.recordedCalls
}

func (evj *XMLEstimatedVehicleJourney) LineRef() string {
	if evj.lineRef == nil {
		s := evj.findStringChildContent(siri_attributes.LineRef)
		evj.lineRef = &s
	}
	return *evj.lineRef
}

func (evj *XMLEstimatedVehicleJourney) Cancellation() bool {
	if !evj.cancellation.Defined {
		evj.cancellation.SetValue(evj.findBoolChildContent(siri_attributes.Cancellation))
	}
	return evj.cancellation.Value
}

func (evj *XMLEstimatedVehicleJourney) DirectionRef() string {
	if evj.directionRef == nil {
		s := evj.findStringChildContent(siri_attributes.DirectionRef)
		evj.directionRef = &s
	}
	return *evj.directionRef
}

func (evj *XMLEstimatedVehicleJourney) OperatorRef() string {
	if evj.operatorRef == nil {
		s := evj.findStringChildContent(siri_attributes.OperatorRef)
		evj.operatorRef = &s
	}
	return *evj.operatorRef
}

func (evj *XMLEstimatedVehicleJourney) DatedVehicleJourneyRef() string {
	if evj.datedVehicleJourneyRef == nil {
		s := evj.findStringChildContent(siri_attributes.DatedVehicleJourneyRef)
		evj.datedVehicleJourneyRef = &s
	}
	return *evj.datedVehicleJourneyRef
}

func (evj *XMLEstimatedVehicleJourney) OriginRef() string {
	if evj.originRef == nil {
		s := evj.findStringChildContent(siri_attributes.OriginRef)
		evj.originRef = &s
	}
	return *evj.originRef
}

func (evj *XMLEstimatedVehicleJourney) DestinationRef() string {
	if evj.destinationRef == nil {
		s := evj.findStringChildContent(siri_attributes.DestinationRef)
		evj.destinationRef = &s
	}
	return *evj.destinationRef
}

func (evj *XMLEstimatedVehicleJourney) DestinationName() string {
	if evj.destinationName == nil {
		s := evj.findStringChildContent(siri_attributes.DestinationName)
		evj.destinationName = &s
	}
	return *evj.destinationName
}

func (evj *XMLEstimatedVehicleJourney) OriginName() string {
	if evj.originName == nil {
		s := evj.findStringChildContent(siri_attributes.OriginName)
		evj.originName = &s
	}
	return *evj.originName
}
