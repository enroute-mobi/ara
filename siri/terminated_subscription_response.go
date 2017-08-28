package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLTerminatedSubscriptionResponse struct {
	ResponseXMLStructure

	subscriberRef   string
	subscriptionRef string
}

type SIRITerminatedSubscriptionResponse struct {
	ResponseTimestamp time.Time
	ResponderRef      string
	Status            bool

	SubscriberRef   string
	SubscriptionRef string
}

const terminatedSubscriptionResponseTemplate = `<ns3:TerminateSubscriptionResponse version="2.0:FR-IDF-2.4">
  <ns3:ResponderRef>{{.ResponderRef}}</ns3:ResponderRef>
  <ns3:TerminationResponseStatus>
  <ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
  <ns3:SubscriberRef>{{.SubscriberRef}}</ns3:SubscriberRef>
  <ns3:SubscriptionRef>{{.SubscriptionRef}}</ns3:SubscriptionRef>
  <ns3:Status>{{.Status}}</ns3:Status>{{if not .Status}}
  <ns3:ErrorCondition>
    <ns3:UnknownSubscriptionError/>
  </ns3: ErrorCondition>{{end}}
</ns3:TerminateSubscriptionResponse>`

func NewXMLTerminatedSubscriptionResponse(node xml.Node) *XMLTerminatedSubscriptionResponse {
	xmlTerminatedSubscriptionResponse := &XMLTerminatedSubscriptionResponse{}
	xmlTerminatedSubscriptionResponse.node = NewXMLNode(node)
	return xmlTerminatedSubscriptionResponse
}

func NewXMLTerminatedSubscriptionResponseFromContent(content []byte) (*XMLTerminatedSubscriptionResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLTerminatedSubscriptionResponse(doc.Root().XmlNode)
	return request, nil
}

func (response *XMLTerminatedSubscriptionResponse) SubscriptionRef() string {
	if response.subscriptionRef == "" {
		response.subscriptionRef = response.findStringChildContent("SubscriptionRef")
	}
	return response.subscriptionRef
}

func (response *XMLTerminatedSubscriptionResponse) SubscriberRef() string {
	if response.subscriptionRef == "" {
		response.subscriptionRef = response.findStringChildContent("SubscriberRef")
	}
	return response.subscriptionRef
}

func (notify *SIRITerminatedSubscriptionResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var terminatedSubscriptionResponse = template.Must(template.New("terminatedSubscriptionResponseTemplate").Parse(terminatedSubscriptionResponseTemplate))
	if err := terminatedSubscriptionResponse.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
