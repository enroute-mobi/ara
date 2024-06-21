package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyGeneralMessage struct {
	ResponseXMLStructure

	deliveries []*XMLGeneralMessageDelivery
}

type XMLGeneralMessageDelivery struct {
	SubscriptionDeliveryXMLStructure

	xmlGeneralMessages              []*XMLGeneralMessage
	xmlGeneralMessagesCancellations []*XMLGeneralMessageCancellation
}

func NewXMLNotifyGeneralMessage(node xml.Node) *XMLNotifyGeneralMessage {
	xmlGeneralMessageResponse := &XMLNotifyGeneralMessage{}
	xmlGeneralMessageResponse.node = NewXMLNode(node)
	return xmlGeneralMessageResponse
}

func NewXMLNotifyGeneralMessageFromContent(content []byte) (*XMLNotifyGeneralMessage, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLNotifyGeneralMessage(doc.Root().XmlNode)
	return response, nil
}

func NewXMLGeneralMessageDelivery(node XMLNode) *XMLGeneralMessageDelivery {
	delivery := &XMLGeneralMessageDelivery{}
	delivery.node = node
	return delivery
}

func (notify *XMLNotifyGeneralMessage) GeneralMessagesDeliveries() []*XMLGeneralMessageDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLGeneralMessageDelivery{}
		nodes := notify.findNodes(siri_attributes.GeneralMessageDelivery)
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLGeneralMessageDelivery(node))
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}

func (delivery *XMLGeneralMessageDelivery) XMLGeneralMessages() []*XMLGeneralMessage {
	if delivery.xmlGeneralMessages == nil {
		nodes := delivery.findNodes(siri_attributes.GeneralMessage)
		for _, node := range nodes {
			delivery.xmlGeneralMessages = append(delivery.xmlGeneralMessages, NewXMLGeneralMessage(node))
		}
	}
	return delivery.xmlGeneralMessages
}

func (delivery *XMLGeneralMessageDelivery) XMLGeneralMessagesCancellations() []*XMLGeneralMessageCancellation {
	if delivery.xmlGeneralMessagesCancellations == nil {
		cancellations := []*XMLGeneralMessageCancellation{}
		nodes := delivery.findNodes(siri_attributes.GeneralMessageCancellation)
		for _, node := range nodes {
			cancellations = append(cancellations, NewXMLCancelledGeneralMessage(node))
		}
		delivery.xmlGeneralMessagesCancellations = cancellations
	}
	return delivery.xmlGeneralMessagesCancellations
}
