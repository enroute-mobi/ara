package sxml

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSiriServiceRequest struct {
	RequestXMLStructure

	stopMonitoringRequests     []*XMLStopMonitoringRequest
	generalMessageRequests     []*XMLGeneralMessageRequest
	estimatedTimetableRequests []*XMLEstimatedTimetableRequest
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

func (request *XMLSiriServiceRequest) StopMonitoringRequests() []*XMLStopMonitoringRequest {
	if len(request.stopMonitoringRequests) == 0 {
		nodes := request.findNodes("StopMonitoringRequest")
		if nodes == nil {
			return request.stopMonitoringRequests
		}
		for _, stopMonitoringNode := range nodes {
			xmlStopMonitoringRequest := &XMLStopMonitoringRequest{}
			xmlStopMonitoringRequest.node = stopMonitoringNode
			request.stopMonitoringRequests = append(request.stopMonitoringRequests, xmlStopMonitoringRequest)
		}
	}
	return request.stopMonitoringRequests
}

func (request *XMLSiriServiceRequest) GeneralMessageRequests() []*XMLGeneralMessageRequest {
	if len(request.generalMessageRequests) == 0 {
		nodes := request.findNodes("GeneralMessageRequest")
		if nodes == nil {
			return request.generalMessageRequests
		}
		for _, generalMessageNode := range nodes {
			xmlGeneralMessageRequest := &XMLGeneralMessageRequest{}
			xmlGeneralMessageRequest.node = generalMessageNode
			request.generalMessageRequests = append(request.generalMessageRequests, xmlGeneralMessageRequest)
		}
	}
	return request.generalMessageRequests
}

func (request *XMLSiriServiceRequest) EstimatedTimetableRequests() []*XMLEstimatedTimetableRequest {
	if len(request.estimatedTimetableRequests) == 0 {
		nodes := request.findNodes("EstimatedTimetableRequest")
		if nodes == nil {
			return request.estimatedTimetableRequests
		}
		for _, estimatedTimetableNode := range nodes {
			xmlGetEstimatedTimetableRequest := &XMLEstimatedTimetableRequest{}
			xmlGetEstimatedTimetableRequest.node = estimatedTimetableNode
			request.estimatedTimetableRequests = append(request.estimatedTimetableRequests, xmlGetEstimatedTimetableRequest)
		}
	}
	return request.estimatedTimetableRequests
}
