package sxml

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifySituationExchange struct {
	ResponseXMLStructure

	deliveries []*XMLSituationExchangeDelivery
}

func NewXMLNotifySituationExchange(node xml.Node) *XMLNotifySituationExchange {
	xmlSituationExchangeResponse := &XMLNotifySituationExchange{}
	xmlSituationExchangeResponse.node = NewXMLNode(node)
	return xmlSituationExchangeResponse
}

func NewXMLNotifySituationExchangeFromContent(content []byte) (*XMLNotifySituationExchange, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLNotifySituationExchange(doc.Root().XmlNode)
	return response, nil
}

func (notify *XMLNotifySituationExchange) SituationExchangesDeliveries() []*XMLSituationExchangeDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLSituationExchangeDelivery{}
		nodes := notify.findNodes("SituationExchangeDelivery")
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLSituationExchangeDelivery(node))
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}
