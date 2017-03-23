package siri

import (
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSiriServiceRequest struct {
	RequestXMLStructure

	stopMonitoringRequests []*XMLSiriServiceStopMonitoringRequest
}

type XMLSiriServiceStopMonitoringRequest struct {
	XMLStructure

	messageIdentifier string
	monitoringRef     string
	stopVisitTypes    string

	requestTimestamp time.Time
}

func NewXMLSiriServiceRequest(node xml.Node) *XMLSiriServiceRequest {
	xmlSiriServiceRequest := &XMLSiriServiceRequest{}
	xmlSiriServiceRequest.node = NewXMLNode(node)
	return xmlSiriServiceRequest
}

func NewXMLSiriServiceRequestFromContent(content []byte) (*XMLSiriServiceRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLSiriServiceRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLSiriServiceRequest) StopMonitoringRequests() []*XMLSiriServiceStopMonitoringRequest {
	if len(request.stopMonitoringRequests) == 0 {
		nodes := request.findNodes("StopMonitoringRequest")
		if nodes == nil {
			return request.stopMonitoringRequests
		}
		for _, stopMonitoringNode := range nodes {
			xmlStopMonitoringRequest := &XMLSiriServiceStopMonitoringRequest{}
			xmlStopMonitoringRequest.node = stopMonitoringNode
			request.stopMonitoringRequests = append(request.stopMonitoringRequests, xmlStopMonitoringRequest)
		}
	}
	return request.stopMonitoringRequests
}

func (request *XMLSiriServiceStopMonitoringRequest) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *XMLSiriServiceStopMonitoringRequest) MonitoringRef() string {
	if request.monitoringRef == "" {
		request.monitoringRef = request.findStringChildContent("MonitoringRef")
	}
	return request.monitoringRef
}

func (request *XMLSiriServiceStopMonitoringRequest) StopVisitTypes() string {
	if request.stopVisitTypes == "" {
		request.stopVisitTypes = request.findStringChildContent("StopVisitTypes")
	}
	return request.stopVisitTypes
}

func (request *XMLSiriServiceStopMonitoringRequest) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent("RequestTimestamp")
	}
	return request.requestTimestamp
}
