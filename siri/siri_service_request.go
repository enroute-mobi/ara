package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSiriServiceRequest struct {
	RequestXMLStructure

	stopMonitoringRequests []*XMLStopMonitoringSubRequest
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

func (request *XMLSiriServiceRequest) StopMonitoringRequests() []*XMLStopMonitoringSubRequest {
	if len(request.stopMonitoringRequests) == 0 {
		nodes := request.findNodes("StopMonitoringRequest")
		if nodes == nil {
			return request.stopMonitoringRequests
		}
		for _, stopMonitoringNode := range nodes {
			xmlStopMonitoringRequest := &XMLStopMonitoringSubRequest{}
			xmlStopMonitoringRequest.node = stopMonitoringNode
			request.stopMonitoringRequests = append(request.stopMonitoringRequests, xmlStopMonitoringRequest)
		}
	}
	return request.stopMonitoringRequests
}
