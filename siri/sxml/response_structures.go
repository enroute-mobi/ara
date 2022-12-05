package sxml

import (
	"strconv"
	"strings"
	"time"
)

type ResponseXMLStructure struct {
	XMLStructure

	address                   string
	producerRef               string
	requestMessageRef         string
	responseMessageIdentifier string
	responseTimestamp         time.Time
}

type ResponseXMLStructureWithStatus struct {
	DeliveryXMLStructure

	address                   string
	producerRef               string
	responseMessageIdentifier string
}

type DeliveryXMLStructure struct {
	LightDeliveryXMLStructure

	requestMessageRef string
}

type LightDeliveryXMLStructure struct {
	XMLStatus

	responseTimestamp time.Time
}

type SubscriptionDeliveryXMLStructure struct {
	LightSubscriptionDeliveryXMLStructure

	requestMessageRef string
}

type LightSubscriptionDeliveryXMLStructure struct {
	LightDeliveryXMLStructure

	subscriberRef   string
	subscriptionRef string
}

type XMLStatus struct {
	XMLStructure

	status           Bool
	errorType        string
	errorNumber      Int
	errorText        string
	errorDescription string
}

func (response *ResponseXMLStructure) Address() string {
	if response.address == "" {
		response.address = response.findStringChildContent("Address")
	}
	return response.address
}

func (response *ResponseXMLStructure) ProducerRef() string {
	if response.producerRef == "" {
		response.producerRef = response.findStringChildContent("ProducerRef")
	}
	return response.producerRef
}

func (response *ResponseXMLStructure) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == "" {
		response.responseMessageIdentifier = response.findStringChildContent("ResponseMessageIdentifier")
	}
	return response.responseMessageIdentifier
}

func (response *ResponseXMLStructure) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *ResponseXMLStructure) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *ResponseXMLStructureWithStatus) Address() string {
	if response.address == "" {
		response.address = response.findStringChildContent("Address")
	}
	return response.address
}

func (response *ResponseXMLStructureWithStatus) ProducerRef() string {
	if response.producerRef == "" {
		response.producerRef = response.findStringChildContent("ProducerRef")
	}
	return response.producerRef
}

func (response *ResponseXMLStructureWithStatus) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == "" {
		response.responseMessageIdentifier = response.findStringChildContent("ResponseMessageIdentifier")
	}
	return response.responseMessageIdentifier
}

func (delivery *DeliveryXMLStructure) RequestMessageRef() string {
	if delivery.requestMessageRef == "" {
		delivery.requestMessageRef = delivery.findStringChildContent("RequestMessageRef")
	}
	return delivery.requestMessageRef
}

func (delivery *LightDeliveryXMLStructure) ResponseTimestamp() time.Time {
	if delivery.responseTimestamp.IsZero() {
		delivery.responseTimestamp = delivery.findTimeChildContent("ResponseTimestamp")
	}
	return delivery.responseTimestamp
}

func (delivery *SubscriptionDeliveryXMLStructure) RequestMessageRef() string {
	if delivery.requestMessageRef == "" {
		delivery.requestMessageRef = delivery.findStringChildContent("RequestMessageRef")
	}
	return delivery.requestMessageRef
}

func (delivery *LightSubscriptionDeliveryXMLStructure) SubscriberRef() string {
	if delivery.subscriberRef == "" {
		delivery.subscriberRef = delivery.findStringChildContent("SubscriberRef")
	}
	return delivery.subscriberRef
}

func (delivery *LightSubscriptionDeliveryXMLStructure) SubscriptionRef() string {
	if delivery.subscriptionRef == "" {
		delivery.subscriptionRef = delivery.findStringChildContent("SubscriptionRef")
	}
	return delivery.subscriptionRef
}

func (response *XMLStatus) Status() bool {
	if !response.status.Defined {
		response.status.SetValue(response.findBoolChildContent("Status"))
	}
	return response.status.Value
}

func (response *XMLStatus) ErrorType() string {
	if !response.Status() && response.errorType == "" {
		node := response.findNode("ErrorText")
		if node != nil {
			response.errorType = node.Parent().Name()
			// Find errorText and errorNumber to avoir too much parsing
			response.errorText = strings.TrimSpace(node.Content())
			if response.errorType == "OtherError" {
				n, err := strconv.Atoi(node.Parent().Attr("number"))
				if err != nil {
					return ""
				}
				response.errorNumber.SetValue(n)
			}
		}
	}
	return response.errorType
}

func (response *XMLStatus) ErrorNumber() int {
	if !response.Status() && response.ErrorType() == "OtherError" && !response.errorNumber.Defined {
		node := response.findNode("ErrorText")
		n, err := strconv.Atoi(node.Parent().Attr("number"))
		if err != nil {
			return -1
		}
		response.errorNumber.SetValue(n)
	}
	return response.errorNumber.Value
}

func (response *XMLStatus) ErrorText() string {
	if !response.Status() && response.errorText == "" {
		response.errorText = response.findStringChildContent("ErrorText")
	}
	return response.errorText
}

func (response *XMLStatus) ErrorDescription() string {
	if !response.Status() && response.errorDescription == "" {
		response.errorDescription = response.findStringChildContent("Description")
	}
	return response.errorDescription
}
