package sxml

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetVehicleMonitoring struct {
	XMLVehicleMonitoringRequest

	requestorRef string
}

type XMLVehicleMonitoringRequest struct {
	LightRequestXMLStructure

	lineRef           string
	messageIdentifier string
}

func NewXMLGetVehicleMonitoring(node xml.Node) *XMLGetVehicleMonitoring {
	xmlGetVehicleMonitoring := &XMLGetVehicleMonitoring{}
	xmlGetVehicleMonitoring.node = NewXMLNode(node)
	return xmlGetVehicleMonitoring
}

func NewXMLGetVehicleMonitoringFromContent(content []byte) (*XMLGetVehicleMonitoring, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetVehicleMonitoring(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLGetVehicleMonitoring) LineRef() string {
	if request.lineRef == "" {
		request.lineRef = request.findStringChildContent("LineRef")
	}
	return request.lineRef
}

func (request *XMLGetVehicleMonitoring) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *XMLGetVehicleMonitoring) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent("RequestorRef")
	}
	return request.requestorRef
}