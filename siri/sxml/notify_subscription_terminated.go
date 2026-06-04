package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLNotifySubscriptionTerminated struct {
	ResponseXMLStructure

	subscriptionRef *string
	subscriberRef   *string
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
	if delivery.subscriberRef == nil {
		s := delivery.findStringChildContent(siri_attributes.SubscriberRef)
		delivery.subscriberRef = &s
	}
	return *delivery.subscriberRef
}

func (delivery *XMLNotifySubscriptionTerminated) SubscriptionRef() string {
	if delivery.subscriptionRef == nil {
		s := delivery.findStringChildContent(siri_attributes.SubscriptionRef)
		delivery.subscriptionRef = &s
	}
	return *delivery.subscriptionRef
}
