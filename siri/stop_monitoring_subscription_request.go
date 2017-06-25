package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringSubscriptionRequest struct {
	RequestXMLStructure

	consumerAddress string

	entries []*XMLStopMonitoringSubscriptionRequestEntry
}

type XMLStopMonitoringSubscriptionRequestEntry struct {
	XMLStructure

	messageIdentifier      string
	monitoringRef          string
	stopVisitTypes         string
	subscriberRef          string
	subscriptionIdentifier string

	initialTerminationTime time.Time
	requestTimestamp       time.Time
}

type SIRIStopMonitoringSubscriptionRequest struct {
	ConsumerAddress   string
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time

	Entries []*SIRIStopMonitoringSubscriptionRequestEntry
}

type SIRIStopMonitoringSubscriptionRequestEntry struct {
	MessageIdentifier      string
	MonitoringRef          string
	SubscriberRef          string
	SubscriptionIdentifier string

	InitialTerminationTime time.Time
	RequestTimestamp       time.Time
}

const stopMonitoringSubscriptionRequestTemplate = `<ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<SubscriptionRequestInfo>
		<siri:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:RequestTimestamp>
		<siri:RequestorRef>{{.RequestorRef}}</siri:RequestorRef>
		<siri:MessageIdentifier>{{.MessageIdentifier}}</siri:MessageIdentifier>{{ if .ConsumerAddress }}
		<siri:ConsumerAddress>{{.ConsumerAddress}}</siri:ConsumerAddress>{{end}}
  </SubscriptionRequestInfo>
	<Request>{{ range .Entries }}
		<StopMonitoringSubscriptionRequest>
			<SubscriberRef>{{.SubscriberRef}}</SubscriberRef>
			<SubscriptionIdentifier>{{.SubscriptionIdentifier}}</SubscriptionIdentifier>
			<InitialTerminationTime>{{.InitialTerminationTime.Format "2006-01-02T15:04:05.000Z07:00"}}</InitialTerminationTime>
			<StopMonitoringRequest>
				<MessageIdentifier>{{.MessageIdentifier}}</MessageIdentifier>
				<RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</RequestTimestamp>
				<MonitoringRef>{{.MonitoringRef}}</MonitoringRef>
				<StopVisitTypes>all</StopVisitTypes>
      </StopMonitoringRequest>
      <IncrementalUpdates>true</IncrementalUpdates>
    	<ChangeBeforeUpdates>PT1M</ChangeBeforeUpdates>
    </StopMonitoringSubscriptionRequest>{{end}}
	</Request>
  <RequestExtension />
</ws:Subscribe>`

func NewXMLStopMonitoringSubscriptionRequest(node xml.Node) *XMLStopMonitoringSubscriptionRequest {
	xmlStopMonitoringSubscriptionRequest := &XMLStopMonitoringSubscriptionRequest{}
	xmlStopMonitoringSubscriptionRequest.node = NewXMLNode(node)
	return xmlStopMonitoringSubscriptionRequest
}

func NewXMLStopMonitoringSubscriptionRequestEntry(node XMLNode) *XMLStopMonitoringSubscriptionRequestEntry {
	xmlStopMonitoringSubscriptionRequestEntry := &XMLStopMonitoringSubscriptionRequestEntry{}
	xmlStopMonitoringSubscriptionRequestEntry.node = node
	return xmlStopMonitoringSubscriptionRequestEntry
}

func NewXMLStopMonitoringSubscriptionRequestFromContent(content []byte) (*XMLStopMonitoringSubscriptionRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLStopMonitoringSubscriptionRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLStopMonitoringSubscriptionRequest) XMLSubscriptionEntries() []*XMLStopMonitoringSubscriptionRequestEntry {
	if len(request.entries) != 0 {
		return request.entries
	}
	nodes := request.findNodes("StopMonitoringSubscriptionRequest")
	if nodes != nil {
		for _, stopMonitoring := range nodes {
			request.entries = append(request.entries, NewXMLStopMonitoringSubscriptionRequestEntry(stopMonitoring))
		}
	}
	return request.entries
}

func (request *XMLStopMonitoringSubscriptionRequest) ConsumerAddress() string {
	if request.consumerAddress == "" {
		request.consumerAddress = request.findStringChildContent("ConsumerAddress")
	}
	return request.consumerAddress
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

func (request *XMLStopMonitoringSubscriptionRequestEntry) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) MonitoringRef() string {
	if request.monitoringRef == "" {
		request.monitoringRef = request.findStringChildContent("MonitoringRef")
	}
	return request.monitoringRef
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent("RequestTimestamp")
	}
	return request.requestTimestamp
}

func (request *SIRIStopMonitoringSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(stopMonitoringSubscriptionRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
