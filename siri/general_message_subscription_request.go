package siri

import (
	"bytes"
	"text/template"
	"time"
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

const generalMessageSubscriptionRequestTemplate = `<sw:Subscribe xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri" xmlns:sws="http://wsdl.siri.org.uk/siri">
	<SubscriptionRequestInfo>
		<siri:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:RequestTimestamp>
		<siri:RequestorRef>{{.RequestorRef}}</siri:RequestorRef>
		<siri:MessageIdentifier>{{.MessageIdentifier}}</siri:MessageIdentifier>{{ if .ConsumerAddress }}
		<siri:ConsumerAddress>{{.ConsumerAddress}}</siri:ConsumerAddress>{{end}}
	</SubscriptionRequestInfo>
	<Request>{{ range .Entries }}
		<siri:GeneralMessageSubscriptionRequest>
			<siri:SubscriberRef>{{.SubscriberRef}}</siri:SubscriberRef>
			<siri:SubscriptionIdentifier>{{.SubscriptionIdentifier}}</siri:SubscriptionIdentifier>
			<siri:InitialTerminationTime>{{.InitialTerminationTime.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:InitialTerminationTime>
			<siri:GeneralMessageRequest version="2.0:FR-IDF-2.4">
				{{ .BuildGeneralMessageRequestXML }}
			</siri:GeneralMessageRequest>
		</siri:GeneralMessageSubscriptionRequest>{{ end }}
	</Request>
	<RequestExtension/>
</sw:Subscribe>`

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
		request.subscriptionIdentifier = request.findStringChildContent("SubscriptionIdentifier")
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
