package sxml

import (
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetStopMonitoring struct {
	XMLStopMonitoringRequest

	requestorRef string
}

type XMLStopMonitoringRequest struct {
	LightXMLStopMonitoringRequest

	previewInterval time.Duration
	startTime       time.Time
}

type LightXMLStopMonitoringRequest struct {
	LightRequestXMLStructure

	monitoringRef     string
	stopVisitTypes    string
	lineRef           string
	maximumStopVisits Int
}

func NewXMLGetStopMonitoring(node xml.Node) *XMLGetStopMonitoring {
	xmlStopMonitoringRequest := &XMLGetStopMonitoring{}
	xmlStopMonitoringRequest.node = NewXMLNode(node)
	return xmlStopMonitoringRequest
}

func NewXMLGetStopMonitoringFromContent(content []byte) (*XMLGetStopMonitoring, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetStopMonitoring(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLGetStopMonitoring) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent("RequestorRef")
	}
	return request.requestorRef
}

func (request *LightXMLStopMonitoringRequest) MonitoringRef() string {
	if request.monitoringRef == "" {
		request.monitoringRef = request.findStringChildContent("MonitoringRef")
	}
	return request.monitoringRef
}

func (request *LightXMLStopMonitoringRequest) StopVisitTypes() string {
	if request.stopVisitTypes == "" {
		request.stopVisitTypes = request.findStringChildContent("StopVisitTypes")
	}
	return request.stopVisitTypes
}

func (request *LightXMLStopMonitoringRequest) LineRef() string {
	if request.lineRef == "" {
		request.lineRef = request.findStringChildContent("LineRef")
	}
	return request.lineRef
}

func (request *LightXMLStopMonitoringRequest) MaximumStopVisits() int {
	if !request.maximumStopVisits.Defined {
		request.maximumStopVisits.SetValue(request.findIntChildContent("MaximumStopVisits"))
	}
	return request.maximumStopVisits.Value
}

func (request *XMLStopMonitoringRequest) PreviewInterval() time.Duration {
	if request.previewInterval == 0 {
		request.previewInterval = request.findDurationChildContent("PreviewInterval")
	}
	return request.previewInterval
}

func (request *XMLStopMonitoringRequest) StartTime() time.Time {
	if request.startTime.IsZero() {
		request.startTime = request.findTimeChildContent("StartTime")
	}
	return request.startTime
}
