package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLDeleteSubscriptionResponse struct {
	XMLStructure

	responderRef      string
	requestMessageRef string

	responseTimestamp time.Time

	responseStatus []*XMLTerminationResponseStatus
}

type XMLTerminationResponseStatus struct {
	LightSubscriptionDeliveryXMLStructure
}

func NewXMLDeleteSubscriptionResponse(node xml.Node) *XMLDeleteSubscriptionResponse {
	xmlDeleteSubscriptionResponse := &XMLDeleteSubscriptionResponse{}
	xmlDeleteSubscriptionResponse.node = NewXMLNode(node)
	return xmlDeleteSubscriptionResponse
}

func NewXMLDeleteSubscriptionResponseFromContent(content []byte) (*XMLDeleteSubscriptionResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLDeleteSubscriptionResponse(doc.Root().XmlNode)
	return request, nil
}

func NewXMLTerminationResponseStatus(node XMLNode) *XMLTerminationResponseStatus {
	responseStatus := &XMLTerminationResponseStatus{}
	responseStatus.node = node
	return responseStatus
}

func (response *XMLDeleteSubscriptionResponse) ResponderRef() string {
	if response.responderRef == "" {
		response.responderRef = response.findStringChildContent(siri_attributes.ResponderRef)
	}
	return response.responderRef
}

func (response *XMLDeleteSubscriptionResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent(siri_attributes.RequestMessageRef)
	}
	return response.requestMessageRef
}

func (response *XMLDeleteSubscriptionResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent(siri_attributes.ResponseTimestamp)
	}
	return response.responseTimestamp
}

func (response *XMLDeleteSubscriptionResponse) ResponseStatus() []*XMLTerminationResponseStatus {
	if len(response.responseStatus) == 0 {
		nodes := response.findNodes(siri_attributes.TerminationResponseStatus)
		if nodes == nil {
			return response.responseStatus
		}
		for _, responseStatus := range nodes {
			response.responseStatus = append(response.responseStatus, NewXMLTerminationResponseStatus(responseStatus))
		}
	}
	return response.responseStatus
}
