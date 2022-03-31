package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
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

type SIRIDeleteSubscriptionResponse struct {
	ResponderRef      string
	RequestMessageRef string
	ResponseTimestamp time.Time

	ResponseStatus []*SIRITerminationResponseStatus
}

type SIRITerminationResponseStatus struct {
	SubscriberRef     string
	SubscriptionRef   string
	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber string
	ErrorText   string
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
		response.responderRef = response.findStringChildContent("ResponderRef")
	}
	return response.responderRef
}

func (response *XMLDeleteSubscriptionResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *XMLDeleteSubscriptionResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *XMLDeleteSubscriptionResponse) ResponseStatus() []*XMLTerminationResponseStatus {
	if len(response.responseStatus) == 0 {
		nodes := response.findNodes("TerminationResponseStatus")
		if nodes == nil {
			return response.responseStatus
		}
		for _, responseStatus := range nodes {
			response.responseStatus = append(response.responseStatus, NewXMLTerminationResponseStatus(responseStatus))
		}
	}
	return response.responseStatus
}

func (notify *SIRIDeleteSubscriptionResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	} else {
		envType = ""
	}

	templateName = fmt.Sprintf("delete_subscription_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
