package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetSituationExchange struct {
	XMLSituationExchangeRequest

	requestorRef string
}

type XMLSituationExchangeRequest struct {
	LightRequestXMLStructure

	lineRefs      []string
	stopPointRefs []string
}

func NewXMLGetSituationExchange(node xml.Node) *XMLGetSituationExchange {
	xmlGetSituationExchange := &XMLGetSituationExchange{}
	xmlGetSituationExchange.node = NewXMLNode(node)
	return xmlGetSituationExchange
}

func NewXMLGetSituationExchangeFromContent(content []byte) (*XMLGetSituationExchange, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetSituationExchange(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLGetSituationExchange) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent(siri_attributes.RequestorRef)
	}
	return request.requestorRef
}
