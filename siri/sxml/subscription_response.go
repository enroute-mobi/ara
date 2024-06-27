package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSubscriptionResponse struct {
	XMLStructure

	address                   string
	requestMessageRef         string
	responderRef              string
	responseMessageIdentifier string

	responseTimestamp  time.Time
	serviceStartedTime time.Time

	responseStatus []*XMLResponseStatus
}

type XMLResponseStatus struct {
	SubscriptionDeliveryXMLStructure

	validUntil time.Time
}

func NewXMLSubscriptionResponse(node xml.Node) *XMLSubscriptionResponse {
	xmlStopMonitoringSubscriptionResponse := &XMLSubscriptionResponse{}
	xmlStopMonitoringSubscriptionResponse.node = NewXMLNode(node)
	return xmlStopMonitoringSubscriptionResponse
}

func NewXMLSubscriptionResponseFromContent(content []byte) (*XMLSubscriptionResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLSubscriptionResponse(doc.Root().XmlNode)
	return response, nil
}

func (response *XMLSubscriptionResponse) ResponseStatus() []*XMLResponseStatus {
	if len(response.responseStatus) == 0 {
		nodes := response.findNodes(siri_attributes.ResponseStatus)
		if nodes == nil {
			return response.responseStatus
		}
		for _, responseStatusNode := range nodes {
			xmlResponseStatus := &XMLResponseStatus{}
			xmlResponseStatus.node = responseStatusNode
			response.responseStatus = append(response.responseStatus, xmlResponseStatus)
		}
	}
	return response.responseStatus
}

func (response *XMLSubscriptionResponse) Address() string {
	if response.address == "" {
		response.address = response.findStringChildContent(siri_attributes.Address)
	}
	return response.address
}

func (response *XMLSubscriptionResponse) ResponderRef() string {
	if response.responderRef == "" {
		response.responderRef = response.findStringChildContent(siri_attributes.ResponderRef)
	}
	return response.responderRef
}

func (response *XMLSubscriptionResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent(siri_attributes.RequestMessageRef)
	}
	return response.requestMessageRef
}

func (response *XMLSubscriptionResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime.IsZero() {
		response.serviceStartedTime = response.findTimeChildContent(siri_attributes.ServiceStartedTime)
	}
	return response.serviceStartedTime
}

func (response *XMLSubscriptionResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent(siri_attributes.ResponseTimestamp)
	}
	return response.responseTimestamp
}

func (response *XMLSubscriptionResponse) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == "" {
		response.responseMessageIdentifier = response.findStringChildContent(siri_attributes.ResponseMessageIdentifier)
	}
	return response.responseMessageIdentifier
}

func (response *XMLResponseStatus) ValidUntil() time.Time {
	if response.validUntil.IsZero() {
		response.validUntil = response.findTimeChildContent(siri_attributes.ValidUntil)
	}
	return response.validUntil
}
