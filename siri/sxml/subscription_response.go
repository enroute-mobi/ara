package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSubscriptionResponse struct {
	XMLStructure

	address                   *string
	requestMessageRef         *string
	responderRef              *string
	responseMessageIdentifier *string

	responseTimestamp  *time.Time
	serviceStartedTime *time.Time

	responseStatus []*XMLResponseStatus
}

type XMLResponseStatus struct {
	SubscriptionDeliveryXMLStructure

	validUntil *time.Time
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
	if response.address == nil {
		s := response.findStringChildContent(siri_attributes.Address)
		response.address = &s
	}
	return *response.address
}

func (response *XMLSubscriptionResponse) ResponderRef() string {
	if response.responderRef == nil {
		s := response.findStringChildContent(siri_attributes.ResponderRef)
		response.responderRef = &s
	}
	return *response.responderRef
}

func (response *XMLSubscriptionResponse) RequestMessageRef() string {
	if response.requestMessageRef == nil {
		s := response.findStringChildContent(siri_attributes.RequestMessageRef)
		response.requestMessageRef = &s
	}
	return *response.requestMessageRef
}

func (response *XMLSubscriptionResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime == nil {
		t := response.findTimeChildContent(siri_attributes.ServiceStartedTime)
		response.serviceStartedTime = &t
	}
	return *response.serviceStartedTime
}

func (response *XMLSubscriptionResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp == nil {
		t := response.findTimeChildContent(siri_attributes.ResponseTimestamp)
		response.responseTimestamp = &t
	}
	return *response.responseTimestamp
}

func (response *XMLSubscriptionResponse) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == nil {
		s := response.findStringChildContent(siri_attributes.ResponseMessageIdentifier)
		response.responseMessageIdentifier = &s
	}
	return *response.responseMessageIdentifier
}

func (response *XMLResponseStatus) ValidUntil() time.Time {
	if response.validUntil == nil {
		t := response.findTimeChildContent(siri_attributes.ValidUntil)
		response.validUntil = &t
	}
	return *response.validUntil
}
