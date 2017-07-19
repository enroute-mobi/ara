package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGeneralMessageSubscriptionRequest struct {
	RequestXMLStructure

	consumerAddress string

	entries []*XMLGeneralMessageSubscriptionRequestEntry
}

type XMLGeneralMessageSubscriptionRequestEntry struct {
	XMLStructure

	messageIdentifier      string
	subscriberRef          string
	subscriptionIdentifier string

	initialTerminationTime time.Time
	requestTimestamp       time.Time
}

type SIRIGeneralMessageSubscriptionRequest struct {
	ConsumerAddress   string
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time

	Entry *SIRIGeneralMessageSubscriptionRequestEntry
}

type SIRIGeneralMessageSubscriptionRequestEntry struct {
	MessageIdentifier      string
	MonitoringRef          string
	SubscriberRef          string
	SubscriptionIdentifier string

	InitialTerminationTime time.Time
	RequestTimestamp       time.Time
}

const generalMessageSubscriptionRequestTemplate = `<ws:Subscribe xmlns:ns2="http://www.siri.org.uk/siri"
                       xmlns:ns3="http://www.ifopt.org.uk/acsb"
                       xmlns:ns4="http://www.ifopt.org.uk/ifopt"
                       xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
                       xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
  <SubscriptionRequestInfo>
    <ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:RequestTimestamp>
    <ns2:RequestorRef>{{.RequestorRef}}</ns2:RequestorRef>
    <ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
    <ns2:ConsumerAddress>{{.ConsumerAddress}}</ns2:ConsumerAddress>
  </SubscriptionRequestInfo>
  <Request version="2.0:FR-IDF-2.4">
    <GeneralMessageSubscriptionRequest>
      <ns2:MessageIdentifier>{{.Entry.MessageIdentifier}}</ns2:MessageIdentifier>
      <ns2:RequestTimestamp>{{.Entry.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:RequestTimestamp>
      <ns5:SubscriberRef>{{.Entry.SubscriberRef}}</ns5:SubscriberRef>
      <ns5:SubscriptionRef>{{.Entry.SubscriptionIdentifier}}</ns5:SubscriptionRef>
      <ns5:InitialTerminationTime>{{.Entry.InitialTerminationTime.Format "2006-01-02T15:04:05.000Z07:00"}}</ns5:InitialTerminationTime>
    </GeneralMessageSubscriptionRequest>
  </Request>
  <RequestExtension />
</ws:Subscribe>`

func NewXMLGeneralMessageSubscriptionResponse(node xml.Node) *XMLGeneralMessageSubscriptionResponse {
	xmlGeneralMessageSubscriptionResponse := &XMLGeneralMessageSubscriptionResponse{}
	xmlGeneralMessageSubscriptionResponse.node = NewXMLNode(node)
	return xmlGeneralMessageSubscriptionResponse
}

func NewXMLGeneralMessageSubscriptionRequest(node xml.Node) *XMLGeneralMessageSubscriptionRequest {
	xmlGeneralMessageSubscriptionRequest := &XMLGeneralMessageSubscriptionRequest{}
	xmlGeneralMessageSubscriptionRequest.node = NewXMLNode(node)
	return xmlGeneralMessageSubscriptionRequest
}

func NewXMLGeneralMessageSubscriptionRequestEntry(node XMLNode) *XMLGeneralMessageSubscriptionRequestEntry {
	xmlGeneralMessageSubscriptionRequestEntry := &XMLGeneralMessageSubscriptionRequestEntry{}
	xmlGeneralMessageSubscriptionRequestEntry.node = node
	return xmlGeneralMessageSubscriptionRequestEntry
}

func NewXMLGeneralMessageSubscriptionRequestFromContent(content []byte) (*XMLGeneralMessageSubscriptionRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGeneralMessageSubscriptionRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLGeneralMessageSubscriptionRequest) XMLSubscriptionEntries() []*XMLGeneralMessageSubscriptionRequestEntry {
	if len(request.entries) != 0 {
		return request.entries
	}
	nodes := request.findNodes("GeneralMessageSubscriptionRequest")
	if nodes != nil {
		for _, generalMessage := range nodes {
			request.entries = append(request.entries, NewXMLGeneralMessageSubscriptionRequestEntry(generalMessage))
		}
	}
	return request.entries
}

func (request *XMLGeneralMessageSubscriptionRequest) ConsumerAddress() string {
	if request.consumerAddress == "" {
		request.consumerAddress = request.findStringChildContent("ConsumerAddress")
	}
	return request.consumerAddress
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionIdentifier == "" {
		request.subscriptionIdentifier = request.findStringChildContent("SubscriptionRef")
	}
	return request.subscriptionIdentifier
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent("InitialTerminationTime")
	}
	return request.initialTerminationTime
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent("RequestTimestamp")
	}
	return request.requestTimestamp
}

func (request *SIRIGeneralMessageSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(generalMessageSubscriptionRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
