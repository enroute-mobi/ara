package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetFacilityMonitoring struct {
	XMLFacilityMonitoringRequest

	requestorRef string
}

type XMLFacilityMonitoringRequest struct {
	LightRequestXMLStructure

	facilityRef string
}

func NewXMLGetFacilityMonitoring(node xml.Node) *XMLGetFacilityMonitoring {
	xmlFacilityMonitoringRequest := &XMLGetFacilityMonitoring{}
	xmlFacilityMonitoringRequest.node = NewXMLNode(node)
	return xmlFacilityMonitoringRequest
}

func NewXMLGetFacilityMonitoringFromContent(content []byte) (*XMLGetFacilityMonitoring, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetFacilityMonitoring(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLGetFacilityMonitoring) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent(siri_attributes.RequestorRef)
	}
	return request.requestorRef
}

func (request *XMLFacilityMonitoringRequest) FacilityRef() string {
	if request.facilityRef == "" {
		request.facilityRef = request.findStringChildContent(siri_attributes.FacilityRef)
	}
	return request.facilityRef
}
