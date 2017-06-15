package siri

import (
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

//WIP

type XMLStopMonitoringSubscriptionResponse struct {
	ResponseXMLStructure

	address                string
	responderRef           string
	requestMessageRef      string
	subscriberRef          string
	subscriptionIdentifier string
	validUntil             time.Time
	serviceStartedTime     time.Time
}

type SIRIStopMonitoringSubscriptionResponse struct {
	Address            string
	ResponderRef       string
	RequestMessageRef  string
	SubscriberRef      string
	SubscriptionRef    string
	ResponseTimestamp  time.Time
	Status             bool
	ValidUntil         time.Time
	ServiceStartedTime time.Time
}

const stopMonitoringSubscriptionResponseTemplate = `<ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
    <SubscriptionAnswerInfo
			xmlns:ns2="http://www.ifopt.org.uk/acsb"
			xmlns:ns3="http://www.ifopt.org.uk/ifopt"
			xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
			xmlns:ns5="http://www.siri.org.uk/siri"
			xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:ResponseTimestamp>
        <ns5:Address>{{.Address}}</ns5:Address>
        <ns5:ResponderRef>{{.ResponderRef}}</ns5:ResponderRef>
        <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">{{.RequestMessageRef}}</ns5:RequestMessageRef>
    </SubscriptionAnswerInfo>
    <Answer
			xmlns:ns2="http://www.ifopt.org.uk/acsb"
			xmlns:ns3="http://www.ifopt.org.uk/ifopt"
			xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
			xmlns:ns5="http://www.siri.org.uk/siri"
			xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseStatus>
            <ns5:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:ResponseTimestamp>
            <ns5:RequestMessageRef>{{.RequestMessageRef}}</ns5:RequestMessageRef>
            <ns5:SubscriberRef>{{.SubscriberRef}}</ns5:SubscriberRef>
            <ns5:SubscriptionRef>{{.SubscriptionIdentifier}}</ns5:SubscriptionRef>
            <ns5:Status>{{.Status}}</ns5:Status>
            <ns5:ValidUntil>{{.ValidUntil}}</ns5:ValidUntil>
        </ns5:ResponseStatus>
        <ns5:ServiceStartedTime>{{.ServiceStartedTime}}</ns5:ServiceStartedTime>
    </Answer>
    <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
</ns1:SubscribeResponse>`

func NewXMLStopMonitoringSubscriptionResponse(node xml.Node) *XMLStopMonitoringSubscriptionResponse {
	xmlStopMonitoringSubscriptionResponse := &XMLStopMonitoringSubscriptionResponse{}
	xmlStopMonitoringSubscriptionResponse.node = NewXMLNode(node)
	return xmlStopMonitoringSubscriptionResponse
}

func NewXMLStopMonitoringSubscriptionResponseFromContent(content []byte) (*XMLStopMonitoringSubscriptionResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLStopMonitoringSubscriptionResponse(doc.Root().XmlNode)
	return response, nil
}

func (request *XMLStopMonitoringSubscriptionResponse) ResponderRef() string {
	if request.responderRef == "" {
		request.responderRef = request.findStringChildContent("ResponderRef")
	}
	return request.responderRef
}

func (request *XMLStopMonitoringSubscriptionResponse) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLStopMonitoringSubscriptionResponse) SubscriptionRef() string {
	if request.subscriptionIdentifier == "" {
		request.subscriptionIdentifier = request.findStringChildContent("SubscriptionRef")
	}
	return request.subscriptionIdentifier
}

func (request *XMLStopMonitoringSubscriptionResponse) RequestMessageRef() string {
	if request.requestMessageRef == "" {
		request.requestMessageRef = request.findStringChildContent("RequestMessageRef")
	}
	return request.requestMessageRef
}

func (request *XMLStopMonitoringSubscriptionResponse) ValidUntil() time.Time {
	if request.validUntil.IsZero() {
		request.validUntil = request.findTimeChildContent("ValidUntil")
	}
	return request.validUntil
}

func (request *XMLStopMonitoringSubscriptionResponse) ServiceStartedTime() time.Time {
	if request.serviceStartedTime.IsZero() {
		request.serviceStartedTime = request.findTimeChildContent("ServiceStartedTime")
	}
	return request.serviceStartedTime
}
