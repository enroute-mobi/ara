package siri

import (
	"bytes"
	"html/template"
	"time"
)

type SIRISubscriptionTerminated struct {
	ResponseMessageIdentifier string
	RequestMessageRef         string

	ProducerRef     string
	SubscriberRef   string
	SubscriptionRef string

	ResponseTimestamp time.Time
}

const subscriptionTerminatedTemplate = `<ns1:SubscriptionTerminatedNotification xmlns:ns1="http://wsdl.siri.org.uk">
      <SubscriptionTerminatedDelivery
       xmlns:ns2="http://www.ifopt.org.uk/acsb"
       xmlns:ns3="http://www.ifopt.org.uk/ifopt"
       xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
       xmlns:ns5="http://www.siri.org.uk/siri"
       xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>{{.ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns5:ResponseTimestamp>
        <ns5:ResponseMessageIdentifier>{{.ResponseMessageIdentifier}}</ns5:ResponseMessageIdentifier>
        <ns5:RequestMessageRef>{{.RequestMessageRef}}</ns5:RequestMessageRef>
        <ns5:ProducerRef>{{.ProducerRef}}</ns5:ProducerRef>
      </SubscriptionTerminatedDelivery>
      <Answer>
        <ns5:ResponseStatus>
          <ns5:SubscriberRef>{{.SubscriberRef}}</ns5:SubscriberRef>
          <ns5:SubscriptionRef>{{.SubscriptionRef}}</ns5:SubscriptionRef>
        </ns5:ResponseStatus>
      </Answer>
      <AnswerExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
    </ns1:SubscriptionTerminatedNotification>`

func (request *SIRISubscriptionTerminated) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("subscriptionTerminated").Parse(subscriptionTerminatedTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
