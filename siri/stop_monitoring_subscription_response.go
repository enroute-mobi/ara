package siri

import (
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

//WIP

type XMLStopMonitoringSubscriptionResponse struct {
	RequestXMLStructure

	address                string
	responderRef           string
	RequestMessageRef      string
	subscriberRef          string
	subscriptionIdentifier string
	initialTerminationTime time.Time
	status                 Bool
}

const stopMonitoringSubscriptionResponseTemplate = `<ns1:SubscribeResponse xmlns:ns1="http://wsdl.siri.org.uk">
    <SubscriptionAnswerInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:ResponseTimestamp>
        <ns5:Address>http://sqybus-siri:8080/ProfilSiriKidf2_4Producer-Sqybus/SiriServices</ns5:Address>
        <ns5:ResponderRef>SQYBUS</ns5:ResponderRef>
        <ns5:RequestMessageRef xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xsi:type="ns5:MessageRefStructure">Edwig:Message::6ba7b814-9dad-11d1-1-00c04fd430c8:LOC</ns5:RequestMessageRef>
    </SubscriptionAnswerInfo>
    <Answer xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
        <ns5:ResponseStatus>
            <ns5:ResponseTimestamp>2017-06-13T10:12:38.926+02:00</ns5:ResponseTimestamp>
            <ns5:RequestMessageRef>28679112-9dad-11d1-2-00c04fd430c8</ns5:RequestMessageRef>
            <ns5:SubscriberRef>RATPDEV:Concerto</ns5:SubscriberRef>
            <ns5:SubscriptionRef>Edwig:Subscription::6ba7b814-9dad-11d1-2-00c04fd430c8:LOC</ns5:SubscriptionRef>
            <ns5:Status>true</ns5:Status>
            <ns5:ValidUntil>2018-01-01T00:59:59.000+01:00</ns5:ValidUntil>
        </ns5:ResponseStatus>
        <ns5:ServiceStartedTime>2017-06-13T03:00:00.000+02:00</ns5:ServiceStartedTime>
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
