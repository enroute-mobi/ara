package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifyGeneralMessage struct {
	ResponseXMLStructure

	deliveries []*XMLGeneralMessageDelivery
}

func (notify *XMLNotifyGeneralMessage) GeneralMessagesDeliveries() []*XMLGeneralMessageDelivery {
	if notify.deliveries == nil {
		deliveries := []*XMLGeneralMessageDelivery{}
		nodes := notify.findNodes("GeneralMessageDelivery")
		if nodes != nil {
			for _, node := range nodes {
				deliveries = append(deliveries, NewXMLGeneralMessageDelivery(node))
			}
		}
		notify.deliveries = deliveries
	}
	return notify.deliveries
}

func (delivery *XMLGeneralMessageDelivery) SubscriptionRef() string {
	if delivery.subscriptionRef == "" {
		delivery.subscriptionRef = delivery.findStringChildContent("SubscriptionRef")
	}
	return delivery.subscriptionRef
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
