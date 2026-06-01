package sxml

import (
	"strconv"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type ResponseXMLStructure struct {
	XMLStructure

	address                   *string
	producerRef               *string
	requestMessageRef         *string
	responseMessageIdentifier *string
	responseTimestamp         *time.Time
}

type ResponseXMLStructureWithStatus struct {
	DeliveryXMLStructure

	address                   *string
	producerRef               *string
	responseMessageIdentifier *string
}

type DeliveryXMLStructure struct {
	LightDeliveryXMLStructure

	requestMessageRef *string
}

type LightDeliveryXMLStructure struct {
	XMLStatus

	responseTimestamp *time.Time
}

type SubscriptionDeliveryXMLStructure struct {
	LightSubscriptionDeliveryXMLStructure

	requestMessageRef *string
}

type LightSubscriptionDeliveryXMLStructure struct {
	LightDeliveryXMLStructure

	subscriberRef   *string
	subscriptionRef *string
}

type XMLStatus struct {
	XMLStructure

	status           Bool
	errorType        *string
	errorNumber      Int
	errorText        *string
	errorDescription *string
}

func (response *ResponseXMLStructure) Address() string {
	if response.address == nil {
		s := response.findStringChildContent(siri_attributes.Address)
		response.address = &s
	}
	return *response.address
}

func (response *ResponseXMLStructure) ProducerRef() string {
	if response.producerRef == nil {
		s := response.findStringChildContent(siri_attributes.ProducerRef)
		response.producerRef = &s
	}
	return *response.producerRef
}

func (response *ResponseXMLStructure) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == nil {
		s := response.findStringChildContent(siri_attributes.ResponseMessageIdentifier)
		response.responseMessageIdentifier = &s
	}
	return *response.responseMessageIdentifier
}

func (response *ResponseXMLStructure) RequestMessageRef() string {
	if response.requestMessageRef == nil {
		s := response.findStringChildContent(siri_attributes.RequestMessageRef)
		response.requestMessageRef = &s
	}
	return *response.requestMessageRef
}

func (response *ResponseXMLStructure) ResponseTimestamp() time.Time {
	if response.responseTimestamp == nil {
		t := response.findTimeChildContent(siri_attributes.ResponseTimestamp)
		response.responseTimestamp = &t
	}
	return *response.responseTimestamp
}

func (response *ResponseXMLStructureWithStatus) Address() string {
	if response.address == nil {
		s := response.findStringChildContent(siri_attributes.Address)
		response.address = &s
	}
	return *response.address
}

func (response *ResponseXMLStructureWithStatus) ProducerRef() string {
	if response.producerRef == nil {
		s := response.findStringChildContent(siri_attributes.ProducerRef)
		response.producerRef = &s
	}
	return *response.producerRef
}

func (response *ResponseXMLStructureWithStatus) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == nil {
		s := response.findStringChildContent(siri_attributes.ResponseMessageIdentifier)
		response.responseMessageIdentifier = &s
	}
	return *response.responseMessageIdentifier
}

func (delivery *DeliveryXMLStructure) RequestMessageRef() string {
	if delivery.requestMessageRef == nil {
		s := delivery.findStringChildContent(siri_attributes.RequestMessageRef)
		delivery.requestMessageRef = &s
	}
	return *delivery.requestMessageRef
}

func (delivery *LightDeliveryXMLStructure) ResponseTimestamp() time.Time {
	if delivery.responseTimestamp == nil {
		t := delivery.findTimeChildContent(siri_attributes.ResponseTimestamp)
		delivery.responseTimestamp = &t
	}
	return *delivery.responseTimestamp
}

func (delivery *SubscriptionDeliveryXMLStructure) RequestMessageRef() string {
	if delivery.requestMessageRef == nil {
		s := delivery.findStringChildContent(siri_attributes.RequestMessageRef)
		delivery.requestMessageRef = &s
	}
	return *delivery.requestMessageRef
}

func (delivery *LightSubscriptionDeliveryXMLStructure) SubscriberRef() string {
	if delivery.subscriberRef == nil {
		s := delivery.findStringChildContent(siri_attributes.SubscriberRef)
		delivery.subscriberRef = &s
	}
	return *delivery.subscriberRef
}

func (delivery *LightSubscriptionDeliveryXMLStructure) SubscriptionRef() string {
	if delivery.subscriptionRef == nil {
		s := delivery.findStringChildContent(siri_attributes.SubscriptionRef)
		delivery.subscriptionRef = &s
	}
	return *delivery.subscriptionRef
}

func (response *XMLStatus) Status() bool {
	if !response.status.Defined {
		response.status.SetValue(response.findBoolChildContent(siri_attributes.Status))
	}
	return response.status.Value
}

func (response *XMLStatus) ErrorType() string {
	if !response.Status() && response.errorType == nil {
		node := response.findNode(siri_attributes.ErrorText)
		if node != nil {
			s := node.Parent().Name()
			response.errorType = &s
			// Find errorText and errorNumber to avoir too much parsing
			et := strings.TrimSpace(node.Content())
			response.errorText = &et
			if *response.errorType == siri_attributes.OtherError {
				n, err := strconv.Atoi(node.Parent().Attr(siri_attributes.Number))
				if err != nil {
					return ""
				}
				response.errorNumber.SetValue(n)
			}
		}
	}
	if response.errorType == nil {
		return ""
	}
	return *response.errorType
}

func (response *XMLStatus) ErrorNumber() int {
	if !response.Status() && response.ErrorType() == siri_attributes.OtherError && !response.errorNumber.Defined {
		node := response.findNode(siri_attributes.ErrorText)
		n, err := strconv.Atoi(node.Parent().Attr(siri_attributes.Number))
		if err != nil {
			return -1
		}
		response.errorNumber.SetValue(n)
	}
	return response.errorNumber.Value
}

func (response *XMLStatus) ErrorText() string {
	if !response.Status() && response.errorText == nil {
		s := response.findStringChildContent(siri_attributes.ErrorText)
		response.errorText = &s
	}
	if response.errorText == nil {
		return ""
	}
	return *response.errorText
}

func (response *XMLStatus) ErrorDescription() string {
	if !response.Status() && response.errorDescription == nil {
		s := response.findStringChildContent(siri_attributes.Description)
		response.errorDescription = &s
	}
	if response.errorDescription == nil {
		return ""
	}
	return *response.errorDescription
}
