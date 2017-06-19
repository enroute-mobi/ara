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

	monitoringRef          string
	subscriberRef          string
	subscriptionIdentifier string
	consumerAddress        string
	initialTerminationTime time.Time
}

type SIRIStopMonitoringSubscriptionRequest struct {
	MessageIdentifier string
	MonitoringRef     string
	RequestorRef      string
	RequestTimestamp  time.Time

	SubscriberRef          string
	SubscriptionIdentifier string
	InitialTerminationTime time.Time
	ConsumerAddress        string
}

const StopMonitoringSubscriptionRequestTemplate = `<ws:Subscribe xmlns:ws="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<SubscriptionRequestInfo>
		<siri:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:RequestTimestamp>
		<siri:RequestorRef>{{.RequestorRef}}</siri:RequestorRef>
		<siri:MessageIdentifier>{{.MessageIdentifier}}</siri:MessageIdentifier>
		<siri:ConsumerAddress>https://edwig-staging.af83.io/test/siri</siri:ConsumerAddress>
  </SubscriptionRequestInfo>
	<Request>
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
    </StopMonitoringSubscriptionRequest>
	</Request>
  <RequestExtension />
</ws:Subscribe>`

func NewXMLStopMonitoringSubscriptionRequest(node xml.Node) *XMLStopMonitoringSubscriptionRequest {
	xmlStopMonitoringSubscriptionRequest := &XMLStopMonitoringSubscriptionRequest{}
	xmlStopMonitoringSubscriptionRequest.node = NewXMLNode(node)
	return xmlStopMonitoringSubscriptionRequest
}

func NewXMLStopMonitoringSubscriptionRequestFromContent(content []byte) (*XMLStopMonitoringSubscriptionRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLStopMonitoringSubscriptionRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *SIRIStopMonitoringSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(StopMonitoringSubscriptionRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (request *XMLStopMonitoringSubscriptionRequest) MonitoringRef() string {
	if request.monitoringRef == "" {
		request.monitoringRef = request.findStringChildContent("MonitoringRef")
	}
	return request.monitoringRef
}

func (request *XMLStopMonitoringSubscriptionRequest) ConsumerAddress() string {
	if request.consumerAddress == "" {
		request.consumerAddress = request.findStringChildContent("ConsumerAddress")
	}
	return request.consumerAddress
}

func (request *XMLStopMonitoringSubscriptionRequest) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLStopMonitoringSubscriptionRequest) SubscriptionIdentifier() string {
	if request.subscriptionIdentifier == "" {
		request.subscriptionIdentifier = request.findStringChildContent("SubscriptionIdentifier")
	}
	return request.subscriptionIdentifier
}

func (request *XMLStopMonitoringSubscriptionRequest) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent("InitialTerminationTime")
	}
	return request.initialTerminationTime
}
