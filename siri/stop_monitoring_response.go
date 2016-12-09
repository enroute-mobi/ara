package siri

import (
	"runtime"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringResponse struct {
	ResponseXMLStructure

	monitoredStopVisits []*XMLMonitoredStopVisit
}

type XMLMonitoredStopVisit struct {
	XMLStructure

	itemIdentifier  string
	departureStatus string
	arrivalStatus   string

	aimedArrivalTime    time.Time
	expectedArrivalTime time.Time
	actualArrivalTime   time.Time

	aimedDepartureTime    time.Time
	expectedDepartureTime time.Time
	actualDepartureTime   time.Time
}

func NewXMLStopMonitoringResponse(node xml.Node) *XMLStopMonitoringResponse {
	xmlStopMonitoringResponse := &XMLStopMonitoringResponse{}
	xmlStopMonitoringResponse.node = node
	return xmlStopMonitoringResponse
}

func NewXMLStopMonitoringResponseFromContent(content []byte) (*XMLStopMonitoringResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLStopMonitoringResponse(doc.Root().XmlNode)
	finalizer := func(response *XMLStopMonitoringResponse) {
		doc.Free()
	}
	runtime.SetFinalizer(response, finalizer)
	return response, nil
}

func (response *XMLStopMonitoringResponse) XMLMonitoredStopVisits() []*XMLMonitoredStopVisit {
	if len(response.monitoredStopVisits) == 0 {
		nodes := response.findNodes("MonitoredStopVisit")
		if nodes == nil {
			return response.monitoredStopVisits
		}
		for _, stopVisitNode := range nodes {
			response.monitoredStopVisits = append(response.monitoredStopVisits, NewXMLMonitoredStopVisitFromContent(stopVisitNode))
		}
	}
	return response.monitoredStopVisits
}

func NewXMLMonitoredStopVisitFromContent(node xml.Node) *XMLMonitoredStopVisit {
	stopVisit := &XMLMonitoredStopVisit{}
	stopVisit.node = node
	return stopVisit
}

func (visit *XMLMonitoredStopVisit) ItemIdentifier() string {
	if visit.itemIdentifier == "" {
		visit.itemIdentifier = visit.findStringChildContent("ItemIdentifier")
	}
	return visit.itemIdentifier
}

func (visit *XMLMonitoredStopVisit) DepartureStatus() string {
	if visit.departureStatus == "" {
		visit.departureStatus = visit.findStringChildContent("DepartureStatus")
	}
	return visit.departureStatus
}

func (visit *XMLMonitoredStopVisit) ArrivalStatus() string {
	if visit.arrivalStatus == "" {
		visit.arrivalStatus = visit.findStringChildContent("ArrivalStatus")
	}
	return visit.arrivalStatus
}

func (visit *XMLMonitoredStopVisit) AimedArrivalTime() time.Time {
	if visit.aimedArrivalTime.IsZero() {
		visit.aimedArrivalTime = visit.findTimeChildContent("AimedArrivalTime")
	}
	return visit.aimedArrivalTime
}

func (visit *XMLMonitoredStopVisit) ExpectedArrivalTime() time.Time {
	if visit.expectedArrivalTime.IsZero() {
		visit.expectedArrivalTime = visit.findTimeChildContent("ExpectedArrivalTime")
	}
	return visit.expectedArrivalTime
}

func (visit *XMLMonitoredStopVisit) ActualArrivalTime() time.Time {
	if visit.actualArrivalTime.IsZero() {
		visit.actualArrivalTime = visit.findTimeChildContent("ActualArrivalTime")
	}
	return visit.actualArrivalTime
}

func (visit *XMLMonitoredStopVisit) AimedDepartureTime() time.Time {
	if visit.aimedDepartureTime.IsZero() {
		visit.aimedDepartureTime = visit.findTimeChildContent("AimedDepartureTime")
	}
	return visit.aimedDepartureTime
}

func (visit *XMLMonitoredStopVisit) ExpectedDepartureTime() time.Time {
	if visit.expectedDepartureTime.IsZero() {
		visit.expectedDepartureTime = visit.findTimeChildContent("ExpectedDepartureTime")
	}
	return visit.expectedDepartureTime
}

func (visit *XMLMonitoredStopVisit) ActualDepartureTime() time.Time {
	if visit.actualDepartureTime.IsZero() {
		visit.actualDepartureTime = visit.findTimeChildContent("ActualDepartureTime")
	}
	return visit.actualDepartureTime
}
