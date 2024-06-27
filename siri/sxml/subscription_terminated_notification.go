package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSubscriptionTerminatedNotification struct {
	SubscriptionDeliveryXMLStructure

	producerRef string
}

func NewXMLSubscriptionTerminatedNotification(node xml.Node) *XMLSubscriptionTerminatedNotification {
	xmlSubscriptionTerminatedNotification := &XMLSubscriptionTerminatedNotification{}
	xmlSubscriptionTerminatedNotification.node = NewXMLNode(node)
	return xmlSubscriptionTerminatedNotification
}

func NewXMLSubscriptionTerminatedNotificationFromContent(content []byte) (*XMLSubscriptionTerminatedNotification, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLSubscriptionTerminatedNotification(doc.Root().XmlNode)
	return request, nil
}

func (response *XMLSubscriptionTerminatedNotification) ProducerRef() string {
	if response.producerRef == "" {
		response.producerRef = response.findStringChildContent(siri_attributes.ProducerRef)
	}
	return response.producerRef
}
