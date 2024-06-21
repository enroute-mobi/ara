package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

// MonitoredCall, EstimatedCall, RecordedCall
type XMLCall struct {
	XMLStructure

	stopPointRef       string
	stopPointName      string
	destinationDisplay string
	arrivalStatus      string
	departureStatus    string

	order Int

	vehicleAtStop Bool

	aimedArrivalTime    time.Time
	expectedArrivalTime time.Time
	actualArrivalTime   time.Time

	aimedDepartureTime    time.Time
	expectedDepartureTime time.Time
	actualDepartureTime   time.Time
}

func NewXMLCall(node XMLNode) *XMLCall {
	call := &XMLCall{}
	call.node = node
	return call
}

func (c *XMLCall) StopPointRef() string {
	if c.stopPointRef == "" {
		c.stopPointRef = c.findStringChildContent(siri_attributes.StopPointRef)
	}
	return c.stopPointRef
}

func (c *XMLCall) StopPointName() string {
	if c.stopPointName == "" {
		c.stopPointName = c.findStringChildContent(siri_attributes.StopPointName)
	}
	return c.stopPointName
}

func (c *XMLCall) DestinationDisplay() string {
	if c.destinationDisplay == "" {
		c.destinationDisplay = c.findStringChildContent(siri_attributes.DestinationDisplay)
	}
	return c.destinationDisplay
}

func (c *XMLCall) ArrivalStatus() string {
	if c.arrivalStatus == "" {
		c.arrivalStatus = c.findStringChildContent(siri_attributes.ArrivalStatus)
	}
	return c.arrivalStatus
}

func (c *XMLCall) DepartureStatus() string {
	if c.departureStatus == "" {
		c.departureStatus = c.findStringChildContent(siri_attributes.DepartureStatus)
	}
	return c.departureStatus
}

func (c *XMLCall) VehicleAtStop() bool {
	if !c.vehicleAtStop.Defined {
		c.vehicleAtStop.SetValue(c.findBoolChildContent(siri_attributes.VehicleAtStop))
	}
	return c.vehicleAtStop.Value
}

func (c *XMLCall) Order() int {
	if !c.order.Defined {
		if c.findNode(siri_attributes.Order) != nil {
			c.order.SetValue(c.findIntChildContent(siri_attributes.Order))

		} else {
			c.order.SetValue(c.findIntChildContent(siri_attributes.VisitNumber))
		}
	}

	return c.order.Value
}

func (c *XMLCall) AimedArrivalTime() time.Time {
	if c.aimedArrivalTime.IsZero() {
		c.aimedArrivalTime = c.findTimeChildContent(siri_attributes.AimedArrivalTime)
	}
	return c.aimedArrivalTime
}

func (c *XMLCall) ExpectedArrivalTime() time.Time {
	if c.expectedArrivalTime.IsZero() {
		c.expectedArrivalTime = c.findTimeChildContent(siri_attributes.ExpectedArrivalTime)
	}
	return c.expectedArrivalTime
}

func (c *XMLCall) ActualArrivalTime() time.Time {
	if c.actualArrivalTime.IsZero() {
		c.actualArrivalTime = c.findTimeChildContent(siri_attributes.ActualArrivalTime)
	}
	return c.actualArrivalTime
}

func (c *XMLCall) AimedDepartureTime() time.Time {
	if c.aimedDepartureTime.IsZero() {
		c.aimedDepartureTime = c.findTimeChildContent(siri_attributes.AimedDepartureTime)
	}
	return c.aimedDepartureTime
}

func (c *XMLCall) ExpectedDepartureTime() time.Time {
	if c.expectedDepartureTime.IsZero() {
		c.expectedDepartureTime = c.findTimeChildContent(siri_attributes.ExpectedDepartureTime)
	}
	return c.expectedDepartureTime
}

func (c *XMLCall) ActualDepartureTime() time.Time {
	if c.actualDepartureTime.IsZero() {
		c.actualDepartureTime = c.findTimeChildContent(siri_attributes.ActualDepartureTime)
	}
	return c.actualDepartureTime
}
