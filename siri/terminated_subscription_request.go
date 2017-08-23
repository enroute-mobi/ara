package siri

import (
	"bytes"
	"html/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type SIRITerminatedSubscriptionRequest struct {
	RequestorRef     string
	RequestTimestamp time.Time

	SubscriptionRef string
}

type XMLTerminatedSubscriptionRequest struct {
	RequestXMLStructure

	subscriptionRef string
}

const terminateSubscriptionRequestTemplate = `<ns1:TerminateSubscriptionRequest xmlns:ns1="http://wsdl.siri.org.uk">
  <ServiceRequestInfo
   xmlns:ns2="http://www.ifopt.org.uk/acsb"
   xmlns:ns3="http://www.ifopt.org.uk/ifopt"
   xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0"
   xmlns:ns5="http://www.siri.org.uk/siri"
   xmlns:ns6="http://wsdl.siri.org.uk/siri">
    <ns5:RequestTimestamp>{{ .RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:RequestTimestamp>
    <ns5:RequestorRef>{{.RequestorRef}}</ns5:RequestorRef>
  </ServiceRequestInfo>
  <Request version="2.0:FR-IDF-2.4">
    <ns5:SubscriptionRef>{{.SubscriptionRef}}</ns5:SubscriptionRef>
  </Request>
  <RequestExtension/>
</ns1:TerminateSubscriptionRequest>`

func NewXMLTerminatedSubscriptionRequest(node xml.Node) *XMLTerminatedSubscriptionRequest {
	xmlTerminatedSubscriptionRequest := &XMLTerminatedSubscriptionRequest{}
	xmlTerminatedSubscriptionRequest.node = NewXMLNode(node)
	return xmlTerminatedSubscriptionRequest
}

func NewXMLTerminatedSubscriptionRequestFromContent(content []byte) (*XMLTerminatedSubscriptionRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLTerminatedSubscriptionRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLTerminatedSubscriptionRequest) SubscriptionRef() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent("SubscriptionRef")
	}
	return request.subscriptionRef
}

func (request *SIRITerminatedSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var terminatedSubscriptionRequest = template.Must(template.New("terminateSubscriptionRequestTemplate").Parse(terminateSubscriptionRequestTemplate))
	if err := terminatedSubscriptionRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
