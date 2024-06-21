package sxml

import (
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

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
		request.subscriptionRef = request.findStringChildContent(siri_attributes.SubscriptionRef)
	}
	return request.subscriptionRef
}

func (request *XMLDeleteSubscriptionRequest) CancelAll() bool {
	if !request.cancelAll.Defined {
		request.cancelAll.SetValue(request.containSelfClosing("All"))
	}
	return request.cancelAll.Value
}
