package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifySubscriptionTerminated struct {
	ResponseXMLStructure

	subscriptionRef string
	subscriberRef   string
}

func NewXMLNotifySubscriptionTerminated(node xml.Node) *XMLNotifySubscriptionTerminated {
	xmlDeleteSubscriptionRequest := &XMLNotifySubscriptionTerminated{}
	xmlDeleteSubscriptionRequest.node = NewXMLNode(node)
	return xmlDeleteSubscriptionRequest
}

func NewXMLNotifySubscriptionTerminatedFromContent(content []byte) (*XMLNotifySubscriptionTerminated, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLNotifySubscriptionTerminated(doc.Root().XmlNode)
	return request, nil
}

func (delivery *XMLNotifySubscriptionTerminated) SubscriberRef() string {
	if delivery.subscriberRef == "" {
		delivery.subscriberRef = delivery.findStringChildContent("SubscriberRef")
	}
	return delivery.subscriberRef
}

func (delivery *XMLNotifySubscriptionTerminated) SubscriptionRef() string {
	if delivery.subscriptionRef == "" {
		delivery.subscriptionRef = delivery.findStringChildContent("SubscriptionRef")
	}
	return delivery.subscriptionRef
}
