package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringSubscriptionTerminatedResponse struct {
	SubscriptionDeliveryXMLStructure

	producerRef string
}

func NewXMLStopMonitoringSubscriptionTerminatedResponse(node xml.Node) *XMLStopMonitoringSubscriptionTerminatedResponse {
	xmlStopMonitoringSubscriptionTerminatedResponse := &XMLStopMonitoringSubscriptionTerminatedResponse{}
	xmlStopMonitoringSubscriptionTerminatedResponse.node = NewXMLNode(node)
	return xmlStopMonitoringSubscriptionTerminatedResponse
}

func NewXMLStopMonitoringSubscriptionTerminatedResponseFromContent(content []byte) (*XMLStopMonitoringSubscriptionTerminatedResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLStopMonitoringSubscriptionTerminatedResponse(doc.Root().XmlNode)
	return request, nil
}

func (response *XMLStopMonitoringSubscriptionTerminatedResponse) ProducerRef() string {
	if response.producerRef == "" {
		response.producerRef = response.findStringChildContent("ProducerRef")
	}
	return response.producerRef
}
