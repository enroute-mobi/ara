package siri

import (
	"bytes"
	"text/template"
	"time"
)

type XMLStopMonitoringSubscriptionRequestEntry struct {
	LightXMLStopMonitoringRequest

	subscriberRef          string
	subscriptionIdentifier string
	initialTerminationTime time.Time
}

type SIRIStopMonitoringSubscriptionRequest struct {
	ConsumerAddress   string
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time

	Entries []*SIRIStopMonitoringSubscriptionRequestEntry
}

type SIRIStopMonitoringSubscriptionRequestEntry struct {
	SIRIStopMonitoringRequest

	SubscriberRef          string
	SubscriptionIdentifier string

	InitialTerminationTime time.Time
}

const stopMonitoringSubscriptionRequestTemplate = `<ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<SubscriptionRequestInfo>
		<siri:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:RequestTimestamp>
		<siri:RequestorRef>{{.RequestorRef}}</siri:RequestorRef>
		<siri:MessageIdentifier>{{.MessageIdentifier}}</siri:MessageIdentifier>{{ if .ConsumerAddress }}
		<siri:ConsumerAddress>{{.ConsumerAddress}}</siri:ConsumerAddress>{{end}}
	</SubscriptionRequestInfo>
	<Request>{{ range .Entries }}
		<siri:StopMonitoringSubscriptionRequest>
			<siri:SubscriberRef>{{.SubscriberRef}}</siri:SubscriberRef>
			<siri:SubscriptionIdentifier>{{.SubscriptionIdentifier}}</siri:SubscriptionIdentifier>
			<siri:InitialTerminationTime>{{.InitialTerminationTime.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:InitialTerminationTime>
			<siri:StopMonitoringRequest version="2.0:FR-IDF-2.4">
				{{ .BuildStopMonitoringRequestXML }}
			</siri:StopMonitoringRequest>
			<siri:IncrementalUpdates>true</siri:IncrementalUpdates>
			<siri:ChangeBeforeUpdates>PT1M</siri:ChangeBeforeUpdates>
		</siri:StopMonitoringSubscriptionRequest>{{end}}
	</Request>
	<RequestExtension />
</ws:Subscribe>`

func NewXMLStopMonitoringSubscriptionRequestEntry(node XMLNode) *XMLStopMonitoringSubscriptionRequestEntry {
	xmlStopMonitoringSubscriptionRequestEntry := &XMLStopMonitoringSubscriptionRequestEntry{}
	xmlStopMonitoringSubscriptionRequestEntry.node = node
	return xmlStopMonitoringSubscriptionRequestEntry
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionIdentifier == "" {
		request.subscriptionIdentifier = request.findStringChildContent("SubscriptionIdentifier")
	}
	return request.subscriptionIdentifier
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent("InitialTerminationTime")
	}
	return request.initialTerminationTime
}

func (request *SIRIStopMonitoringSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(stopMonitoringSubscriptionRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
