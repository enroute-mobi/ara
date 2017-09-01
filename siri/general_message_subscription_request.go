package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri/xml"
)

type XMLGeneralMessageSubscriptionRequestEntry struct {
	XMLGeneralMessageRequest

	subscriberRef          string
	subscriptionIdentifier string

	initialTerminationTime time.Time
}

type SIRIGeneralMessageSubscriptionRequest struct {
	ConsumerAddress   string
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time

	Entries []*SIRIGeneralMessageSubscriptionRequestEntry
}

type SIRIGeneralMessageSubscriptionRequestEntry struct {
	SIRIGeneralMessageRequest

	SubscriberRef          string
	SubscriptionIdentifier string

	InitialTerminationTime time.Time
}

const generalMessageSubscriptionRequestTemplate = `<ws:Subscribe xmlns:ns2="http://www.siri.org.uk/siri"
											 xmlns:ns3="http://www.ifopt.org.uk/acsb"
											 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
											 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
											 xmlns:ns6="http://scma/siri"
											 xmlns:ns7="http://wsdl.siri.org.uk">
	<SubscriptionRequestInfo>
		<ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:RequestTimestamp>
		<ns2:RequestorRef>{{.RequestorRef}}</ns2:RequestorRef>
		<ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
		<ns2:ConsumerAddress>{{.ConsumerAddress}}</ns2:ConsumerAddress>
	</SubscriptionRequestInfo>
	<Request version="2.0:FR-IDF-2.4">{{ range .Entries }}
		<GeneralMessageSubscriptionRequest>
			<ns5:SubscriberRef>{{.SubscriberRef}}</ns5:SubscriberRef>
			<ns5:SubscriptionRef>{{.SubscriptionIdentifier}}</ns5:SubscriptionRef>
			<ns5:InitialTerminationTime>{{.InitialTerminationTime.Format "2006-01-02T15:04:05.000Z07:00"}}</ns5:InitialTerminationTime>
			<ns2:GeneralMessageRequest>
				{{ .BuildGeneralMessageRequestXML }}
			</ns2:GeneralMessageRequest>
		</GeneralMessageSubscriptionRequest>{{ end }}
	</Request>
	<RequestExtension />
</ws:Subscribe>`

func NewXMLGeneralMessageSubscriptionResponse(node xml.Node) *XMLGeneralMessageSubscriptionResponse {
	xmlGeneralMessageSubscriptionResponse := &XMLGeneralMessageSubscriptionResponse{}
	xmlGeneralMessageSubscriptionResponse.node = NewXMLNode(node)
	return xmlGeneralMessageSubscriptionResponse
}

func NewXMLGeneralMessageSubscriptionRequestEntry(node XMLNode) *XMLGeneralMessageSubscriptionRequestEntry {
	xmlGeneralMessageSubscriptionRequestEntry := &XMLGeneralMessageSubscriptionRequestEntry{}
	xmlGeneralMessageSubscriptionRequestEntry.node = node
	return xmlGeneralMessageSubscriptionRequestEntry
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

func (request *SIRIGeneralMessageSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(generalMessageSubscriptionRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
