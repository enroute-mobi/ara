package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/edwig/logger"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type SIRIDeleteSubscriptionRequest struct {
	RequestorRef     string
	RequestTimestamp time.Time

	MessageIdentifier string
	SubscriptionRef   string
	CancelAll         bool
}

type XMLDeleteSubscriptionRequest struct {
	RequestXMLStructure

	cancelAll       Bool
	subscriptionRef string
}

func NewXMLDeleteSubscriptionRequest(node xml.Node) *XMLDeleteSubscriptionRequest {
	xmlDeleteSubscriptionRequest := &XMLDeleteSubscriptionRequest{}
	xmlDeleteSubscriptionRequest.node = NewXMLNode(node)
	return xmlDeleteSubscriptionRequest
}

func NewXMLDeleteSubscriptionRequestFromContent(content []byte) (*XMLDeleteSubscriptionRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLDeleteSubscriptionRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLDeleteSubscriptionRequest) SubscriptionRef() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent("SubscriptionRef")
	}
	return request.subscriptionRef
}

func (request *XMLDeleteSubscriptionRequest) CancelAll() bool {
	if !request.cancelAll.Defined {
		request.cancelAll.SetValue(request.containSelfClosing("All"))
	}
	return request.cancelAll.Value
}

func (request *SIRIDeleteSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "delete_subscription_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
