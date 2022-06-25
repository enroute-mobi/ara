package sxml

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLCheckStatusRequest struct {
	RequestXMLStructure
}

func NewXMLCheckStatusRequest(node xml.Node) *XMLCheckStatusRequest {
	xmlCheckStatusRequest := &XMLCheckStatusRequest{}
	xmlCheckStatusRequest.node = NewXMLNode(node)
	return xmlCheckStatusRequest
}

func NewXMLCheckStatusRequestFromContent(content []byte) (*XMLCheckStatusRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLCheckStatusRequest(doc.Root().XmlNode)
	return request, nil
}
