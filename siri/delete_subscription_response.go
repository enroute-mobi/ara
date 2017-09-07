package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLDeleteSubscriptionResponse struct {
	ResponseXMLStructure

	responderRef    string
	subscriberRef   string
	subscriptionRef string
}

type SIRIDeleteSubscriptionResponse struct {
	ResponseTimestamp time.Time
	ResponderRef      string
	Status            bool

	SubscriberRef   string
	SubscriptionRef string
}

const deleteSubscriptionResponseTemplate = `<ns3:DeleteSubscriptionResponse version="2.0:FR-IDF-2.4">
  <ns5:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:ResponseTimestamp>
  <ns3:ResponderRef>{{ .ResponderRef }}</ns3:ResponderRef>
  <ns5:RequestMessageRef>{{ .RequestMessageRef }}</ns5:RequestMessageRef>
  <ns3:TerminationResponseStatus>
	  <ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
	  <ns3:SubscriberRef>{{ .SubscriberRef }}</ns3:SubscriberRef>
	  <ns3:SubscriptionRef>{{ .SubscriptionRef }}</ns3:SubscriptionRef>
	  <ns3:Status>{{ .Status }}</ns3:Status>{{ if not .Status }}
	  <ns3:ErrorCondition>
	    <ns3:UnknownSubscriptionError/>
	  </ns3:ErrorCondition>{{ end }}
  </ns3:TerminationResponseStatus>
</ns3:DeleteSubscriptionResponse>`

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

func (response *XMLDeleteSubscriptionResponse) SubscriptionRef() string {
	if response.subscriptionRef == "" {
		response.subscriptionRef = response.findStringChildContent("SubscriptionRef")
	}
	return response.subscriptionRef
}

func (response *XMLDeleteSubscriptionResponse) ResponderRef() string {
	if response.responderRef == "" {
		response.responderRef = response.findStringChildContent("ResponderRef")
	}
	return response.responderRef
}

func (response *XMLDeleteSubscriptionResponse) SubscriberRef() string {
	if response.subscriptionRef == "" {
		response.subscriptionRef = response.findStringChildContent("SubscriberRef")
	}
	return response.subscriptionRef
}

func (notify *SIRIDeleteSubscriptionResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var deleteSubscriptionResponse = template.Must(template.New("deleteSubscriptionResponseTemplate").Parse(deleteSubscriptionResponseTemplate))
	if err := deleteSubscriptionResponse.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
