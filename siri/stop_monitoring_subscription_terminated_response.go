package siri

import (
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringSubscriptionTerminatedResponse struct {
	ResponseXMLStructure

	subscriptionTerminateds []*XMLSubscriptionTerminated
}

type XMLSubscriptionTerminated struct {
	ResponseXMLStructure

	subscriberRef   string
	subscriptionRef string
}

func NewXMLStopMonitoringSubscriptionTerminatedResponse(node xml.Node) *XMLStopMonitoringSubscriptionTerminatedResponse {
	xmlStopMonitoringSubscriptionTerminatedResponse := &XMLStopMonitoringSubscriptionTerminatedResponse{}
	xmlStopMonitoringSubscriptionTerminatedResponse.node = NewXMLNode(node)
	return xmlStopMonitoringSubscriptionTerminatedResponse
}

func NewXMLSubscriptionTerminated(node XMLNode) *XMLSubscriptionTerminated {
	xmlSubscriptionTerminated := &XMLSubscriptionTerminated{}
	xmlSubscriptionTerminated.node = node
	return xmlSubscriptionTerminated
}

func NewXMLStopMonitoringSubscriptionTerminatedResponseFromContent(content []byte) (*XMLStopMonitoringSubscriptionTerminatedResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLStopMonitoringSubscriptionTerminatedResponse(doc.Root().XmlNode)
	return request, nil
}

func (response *XMLStopMonitoringSubscriptionTerminatedResponse) XMLSubscriptionTerminateds() []*XMLSubscriptionTerminated {
	if len(response.subscriptionTerminateds) == 0 {
		nodes := response.findNodes("SubscriptionTerminatedNotification")
		if nodes == nil {
			return response.subscriptionTerminateds
		}
		for _, subscriptionTerminatedNode := range nodes {
			response.subscriptionTerminateds = append(response.subscriptionTerminateds, NewXMLSubscriptionTerminated(subscriptionTerminatedNode))
		}
	}
	return response.subscriptionTerminateds
}

func (sub *XMLSubscriptionTerminated) SubscriberRef() string {
	if sub.subscriberRef == "" {
		sub.subscriberRef = sub.findStringChildContent("SubscriberRef")
	}
	return sub.subscriberRef
}

func (sub *XMLSubscriptionTerminated) SubscriptionRef() string {
	if sub.subscriptionRef == "" {
		sub.subscriptionRef = sub.findStringChildContent("SubscriptionRef")
	}
	return sub.subscriptionRef
}
