package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
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

func (request *SIRIDeleteSubscriptionRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("delete_subscription_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
