package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

// MonitoredCall, EstimatedCall, RecordedCall
type XMLCall struct {
	XMLStructure

	stopPointRef       *string
	stopPointName      *string
	destinationDisplay *string
	arrivalStatus      *string
	departureStatus    *string

	order Int

	vehicleAtStop Bool

	aimedArrivalTime    *time.Time
	expectedArrivalTime *time.Time
	actualArrivalTime   *time.Time

	aimedDepartureTime    *time.Time
	expectedDepartureTime *time.Time
	actualDepartureTime   *time.Time
}

func NewXMLCall(node XMLNode) *XMLCall {
	call := &XMLCall{}
	call.node = node
	return call
}

func (c *XMLCall) StopPointRef() string {
	if c.stopPointRef == nil {
		s := c.findStringChildContent(siri_attributes.StopPointRef)
		c.stopPointRef = &s
	}
	return *c.stopPointRef
}

func (c *XMLCall) StopPointName() string {
	if c.stopPointName == nil {
		s := c.findStringChildContent(siri_attributes.StopPointName)
		c.stopPointName = &s
	}
	return *c.stopPointName
}

func (c *XMLCall) DestinationDisplay() string {
	if c.destinationDisplay == nil {
		s := c.findStringChildContent(siri_attributes.DestinationDisplay)
		c.destinationDisplay = &s
	}
	return *c.destinationDisplay
}

func (c *XMLCall) ArrivalStatus() string {
	if c.arrivalStatus == nil {
		s := c.findStringChildContent(siri_attributes.ArrivalStatus)
		c.arrivalStatus = &s
	}
	return *c.arrivalStatus
}

func (c *XMLCall) DepartureStatus() string {
	if c.departureStatus == nil {
		s := c.findStringChildContent(siri_attributes.DepartureStatus)
		c.departureStatus = &s
	}
	return *c.departureStatus
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
	if c.aimedArrivalTime == nil {
		t := c.findTimeChildContent(siri_attributes.AimedArrivalTime)
		c.aimedArrivalTime = &t
	}
	return *c.aimedArrivalTime
}

func (c *XMLCall) ExpectedArrivalTime() time.Time {
	if c.expectedArrivalTime == nil {
		t := c.findTimeChildContent(siri_attributes.ExpectedArrivalTime)
		c.expectedArrivalTime = &t
	}
	return *c.expectedArrivalTime
}

func (c *XMLCall) ActualArrivalTime() time.Time {
	if c.actualArrivalTime == nil {
		t := c.findTimeChildContent(siri_attributes.ActualArrivalTime)
		c.actualArrivalTime = &t
	}
	return *c.actualArrivalTime
}

func (c *XMLCall) AimedDepartureTime() time.Time {
	if c.aimedDepartureTime == nil {
		t := c.findTimeChildContent(siri_attributes.AimedDepartureTime)
		c.aimedDepartureTime = &t
	}
	return *c.aimedDepartureTime
}

func (c *XMLCall) ExpectedDepartureTime() time.Time {
	if c.expectedDepartureTime == nil {
		t := c.findTimeChildContent(siri_attributes.ExpectedDepartureTime)
		c.expectedDepartureTime = &t
	}
	return *c.expectedDepartureTime
}

func (c *XMLCall) ActualDepartureTime() time.Time {
	if c.actualDepartureTime == nil {
		t := c.findTimeChildContent(siri_attributes.ActualDepartureTime)
		c.actualDepartureTime = &t
	}
	return *c.actualDepartureTime
}
