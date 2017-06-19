package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringRequest struct {
	XMLStopMonitoringSubRequest

	requestorRef string
}

type XMLStopMonitoringSubRequest struct {
	XMLStructure

	messageIdentifier string
	monitoringRef     string
	stopVisitTypes    string
	lineRef           string

	maximumStopVisits int

	previewInterval time.Duration

	startTime        time.Time
	requestTimestamp time.Time
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

func (request *XMLStopMonitoringRequest) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent("RequestorRef")
	}
	return request.requestorRef
}

func (request *SIRIStopMonitoringRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(stopMonitoringRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (request *XMLStopMonitoringSubRequest) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *XMLStopMonitoringSubRequest) MonitoringRef() string {
	if request.monitoringRef == "" {
		request.monitoringRef = request.findStringChildContent("MonitoringRef")
	}
	return request.monitoringRef
}

func (request *XMLStopMonitoringSubRequest) StopVisitTypes() string {
	if request.stopVisitTypes == "" {
		request.stopVisitTypes = request.findStringChildContent("StopVisitTypes")
	}
	return request.stopVisitTypes
}

func (request *XMLStopMonitoringSubRequest) LineRef() string {
	if request.lineRef == "" {
		request.lineRef = request.findStringChildContent("LineRef")
	}
	return request.lineRef
}

func (request *XMLStopMonitoringSubRequest) MaximumStopVisits() int {
	if request.maximumStopVisits == 0 {
		request.maximumStopVisits = request.findIntChildContent("MaximumStopVisits")
	}
	return request.maximumStopVisits
}

func (request *XMLStopMonitoringSubRequest) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent("RequestTimestamp")
	}
	return request.requestTimestamp
}

func (request *XMLStopMonitoringSubRequest) PreviewInterval() time.Duration {
	if request.previewInterval == 0 {
		request.previewInterval = request.findDurationChildContent("PreviewInterval")
	}
	return request.previewInterval
}

func (request *XMLStopMonitoringSubRequest) StartTime() time.Time {
	if request.startTime.IsZero() {
		request.startTime = request.findTimeChildContent("StartTime")
	}
	return request.startTime
}
