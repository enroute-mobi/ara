package sxml

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopPointsDiscoveryRequest struct {
	RequestXMLStructure
}

func NewXMLStopPointsDiscoveryRequest(node xml.Node) *XMLStopPointsDiscoveryRequest {
	xmlStopDiscoveryRequest := &XMLStopPointsDiscoveryRequest{}
	xmlStopDiscoveryRequest.node = NewXMLNode(node)
	return xmlStopDiscoveryRequest
}

func NewXMLStopPointsDiscoveryRequestFromContent(content []byte) (*XMLStopPointsDiscoveryRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLStopPointsDiscoveryRequest(doc.Root().XmlNode)
	return request, nil
}
