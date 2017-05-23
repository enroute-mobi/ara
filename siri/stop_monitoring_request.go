package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringRequest struct {
	RequestXMLStructure

	monitoringRef string
}

type SIRIStopMonitoringRequest struct {
	MessageIdentifier string
	MonitoringRef     string
	RequestorRef      string
	RequestTimestamp  time.Time
}

const stopMonitoringRequestTemplate = `<ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
													 xmlns:ns3="http://www.ifopt.org.uk/acsb"
													 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
													 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
													 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
	<ServiceRequestInfo>
		<ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:RequestTimestamp>
		<ns2:RequestorRef>{{.RequestorRef}}</ns2:RequestorRef>
		<ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
	</ServiceRequestInfo>
	<Request version="2.0:FR-IDF-2.4">
		<ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:RequestTimestamp>
		<ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
		<ns2:MonitoringRef>{{.MonitoringRef}}</ns2:MonitoringRef>
		<ns2:StopVisitTypes>all</ns2:StopVisitTypes>
	</Request>
	<RequestExtension />
</ns7:GetStopMonitoring>`

func NewXMLStopMonitoringRequest(node xml.Node) *XMLStopMonitoringRequest {
	xmlStopMonitoringRequest := &XMLStopMonitoringRequest{}
	xmlStopMonitoringRequest.node = NewXMLNode(node)
	return xmlStopMonitoringRequest
}

func NewXMLStopMonitoringRequestFromContent(content []byte) (*XMLStopMonitoringRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLStopMonitoringRequest(doc.Root().XmlNode)
	return request, nil
}

func NewSIRIStopMonitoringRequest(
	messageIdentifier,
	monitoringRef,
	requestorRef string,
	requestTimestamp time.Time) *SIRIStopMonitoringRequest {
	return &SIRIStopMonitoringRequest{
		MessageIdentifier: messageIdentifier,
		MonitoringRef:     monitoringRef,
		RequestorRef:      requestorRef,
		RequestTimestamp:  requestTimestamp,
	}
}

func (request *XMLStopMonitoringRequest) MonitoringRef() string {
	if request.monitoringRef == "" {
		request.monitoringRef = request.findStringChildContent("MonitoringRef")
	}
	return request.monitoringRef
}

func (request *SIRIStopMonitoringRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(stopMonitoringRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
