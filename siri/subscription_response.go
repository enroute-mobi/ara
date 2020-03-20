package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLSubscriptionResponse struct {
	XMLStructure

	address           string
	requestMessageRef string
	responderRef      string

	responseTimestamp  time.Time
	serviceStartedTime time.Time

	responseStatus []*XMLResponseStatus
}

type XMLResponseStatus struct {
	SubscriptionDeliveryXMLStructure

	validUntil time.Time
}

type SIRISubscriptionResponse struct {
	Address           string
	ResponderRef      string
	RequestMessageRef string

	ResponseTimestamp  time.Time
	ServiceStartedTime time.Time

	ResponseStatus []SIRIResponseStatus
}

type SIRIResponseStatus struct {
	RequestMessageRef string
	SubscriberRef     string
	SubscriptionRef   string

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	ResponseTimestamp time.Time
	ValidUntil        time.Time
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
		nodes := response.findNodes("ResponseStatus")
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
		response.address = response.findStringChildContent("Address")
	}
	return response.address
}

func (response *XMLSubscriptionResponse) ResponderRef() string {
	if response.responderRef == "" {
		response.responderRef = response.findStringChildContent("ResponderRef")
	}
	return response.responderRef
}

func (response *XMLSubscriptionResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *XMLSubscriptionResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime.IsZero() {
		response.serviceStartedTime = response.findTimeChildContent("ServiceStartedTime")
	}
	return response.serviceStartedTime
}

func (response *XMLSubscriptionResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *XMLResponseStatus) ValidUntil() time.Time {
	if response.validUntil.IsZero() {
		response.validUntil = response.findTimeChildContent("ValidUntil")
	}
	return response.validUntil
}

func (response *SIRISubscriptionResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "subscription_response.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
