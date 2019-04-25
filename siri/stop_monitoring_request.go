package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetStopMonitoring struct {
	XMLStopMonitoringRequest

	requestorRef string
}

type XMLStopMonitoringRequest struct {
	LightXMLStopMonitoringRequest

	previewInterval time.Duration
	startTime       time.Time
}

type LightXMLStopMonitoringRequest struct {
	LightRequestXMLStructure

	monitoringRef     string
	stopVisitTypes    string
	lineRef           string
	maximumStopVisits int
}

type SIRIGetStopMonitoringRequest struct {
	SIRIStopMonitoringRequest

	RequestorRef string
}

type SIRIStopMonitoringRequest struct {
	MessageIdentifier string
	MonitoringRef     string
	StopVisitTypes    string

	RequestTimestamp time.Time
}

const getStopMonitoringRequestTemplate = `<sw:GetStopMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceRequestInfo>
		<siri:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:RequestTimestamp>
		<siri:RequestorRef>{{.RequestorRef}}</siri:RequestorRef>
		<siri:MessageIdentifier>{{.MessageIdentifier}}</siri:MessageIdentifier>
	</ServiceRequestInfo>
	<Request version="2.0:FR-IDF-2.4">
		{{ .BuildStopMonitoringRequestXML }}
	</Request>
	<RequestExtension />
</sw:GetStopMonitoring>`

const stopMonitoringRequestTemplate = `<siri:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:RequestTimestamp>
		<siri:MessageIdentifier>{{.MessageIdentifier}}</siri:MessageIdentifier>
		<siri:MonitoringRef>{{.MonitoringRef}}</siri:MonitoringRef>{{ if .StopVisitTypes }}
		<siri:StopVisitTypes>{{.StopVisitTypes}}</siri:StopVisitTypes>{{ end }}`

func NewXMLGetStopMonitoring(node xml.Node) *XMLGetStopMonitoring {
	xmlStopMonitoringRequest := &XMLGetStopMonitoring{}
	xmlStopMonitoringRequest.node = NewXMLNode(node)
	return xmlStopMonitoringRequest
}

func NewXMLGetStopMonitoringFromContent(content []byte) (*XMLGetStopMonitoring, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetStopMonitoring(doc.Root().XmlNode)
	return request, nil
}

func NewSIRIGetStopMonitoringRequest(
	messageIdentifier,
	monitoringRef,
	requestorRef string,
	requestTimestamp time.Time) *SIRIGetStopMonitoringRequest {
	request := &SIRIGetStopMonitoringRequest{
		RequestorRef: requestorRef,
	}
	request.MessageIdentifier = messageIdentifier
	request.MonitoringRef = monitoringRef
	request.RequestTimestamp = requestTimestamp
	return request
}

func (request *XMLGetStopMonitoring) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent("RequestorRef")
	}
	return request.requestorRef
}

func (request *SIRIGetStopMonitoringRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(getStopMonitoringRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIStopMonitoringRequest) BuildStopMonitoringRequestXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(stopMonitoringRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (request *LightXMLStopMonitoringRequest) MonitoringRef() string {
	if request.monitoringRef == "" {
		request.monitoringRef = request.findStringChildContent("MonitoringRef")
	}
	return request.monitoringRef
}

func (request *LightXMLStopMonitoringRequest) StopVisitTypes() string {
	if request.stopVisitTypes == "" {
		request.stopVisitTypes = request.findStringChildContent("StopVisitTypes")
	}
	return request.stopVisitTypes
}

func (request *LightXMLStopMonitoringRequest) LineRef() string {
	if request.lineRef == "" {
		request.lineRef = request.findStringChildContent("LineRef")
	}
	return request.lineRef
}

func (request *LightXMLStopMonitoringRequest) MaximumStopVisits() int {
	if request.maximumStopVisits == 0 {
		request.maximumStopVisits = request.findIntChildContent("MaximumStopVisits")
	}
	return request.maximumStopVisits
}

func (request *XMLStopMonitoringRequest) PreviewInterval() time.Duration {
	if request.previewInterval == 0 {
		request.previewInterval = request.findDurationChildContent("PreviewInterval")
	}
	return request.previewInterval
}

func (request *XMLStopMonitoringRequest) StartTime() time.Time {
	if request.startTime.IsZero() {
		request.startTime = request.findTimeChildContent("StartTime")
	}
	return request.startTime
}
