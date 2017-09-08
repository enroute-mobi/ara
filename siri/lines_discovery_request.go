package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLLinesDiscoveryRequest struct {
	RequestXMLStructure
}

func NewXMLLinesDiscoveryRequest(node xml.Node) *XMLLinesDiscoveryRequest {
	xmlLinesDiscoveryRequest := &XMLLinesDiscoveryRequest{}
	xmlLinesDiscoveryRequest.node = NewXMLNode(node)
	return xmlLinesDiscoveryRequest
}

func NewXMLLinesDiscoveryRequestFromContent(content []byte) (*XMLLinesDiscoveryRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLLinesDiscoveryRequest(doc.Root().XmlNode)
	return request, nil
}
