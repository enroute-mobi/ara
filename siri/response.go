package siri

import (
	"bytes"
	"html/template"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLCheckStatusResponse struct {
	node *xml.XmlNode

	producerRef               string
	requestMessageRef         string
	responseMessageIdentifier string
	status                    bool
	responseTimestamp         time.Time
	serviceStartedTime        time.Time
}

type SIRICheckStatusResponse struct {
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string
	Status                    bool
	ResponseTimestamp         time.Time
	ServiceStartedTime        time.Time
}

const SIRIResponseTemplate = `<ns7:CheckStatusResponse xmlns:ns2="http://www.siri.org.uk/siri"
												 xmlns:ns3="http://www.ifopt.org.uk/acsb"
												 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
												 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
												 xmlns:ns6="http://scma/siri"
												 xmlns:ns7="http://wsdl.siri.org.uk"
												 xmlns:ns8="http://wsdl.siri.org.uk/siri">
	<CheckStatusAnswerInfo>
		<ns2:ResponseTimestamp>{{.ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:ResponseTimestamp>
		<ns2:ProducerRef>{{.ProducerRef}}</ns2:ProducerRef>
		<ns2:Address>http://appli.chouette.mobi/siri_france/siri</ns2:Address>
		<ns2:ResponseMessageIdentifier>{{.ResponseMessageIdentifier}}</ns2:ResponseMessageIdentifier>
		<ns2:RequestMessageRef>{{.RequestMessageRef}}</ns2:RequestMessageRef>
	</CheckStatusAnswerInfo>
	<Answer>
		<ns2:Status>{{.Status}}</ns2:Status>
		<ns2:ServiceStartedTime>{{.ServiceStartedTime.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:ServiceStartedTime>
	</Answer>
	<AnswerExtension />
</ns7:CheckStatusResponse>`

func NewXMLCheckStatusResponse(node *xml.XmlNode) *XMLCheckStatusResponse {
	return &XMLCheckStatusResponse{node: node}
}

func NewXMLCheckStatusResponseFromContent(content []byte) *XMLCheckStatusResponse {
	doc, _ := gokogiri.ParseXml(content)
	request := NewXMLCheckStatusResponse(doc.Root().XmlNode)
	finalizer := func(request *XMLCheckStatusResponse) {
		doc.Free()
	}
	runtime.SetFinalizer(request, finalizer)
	return request
}

func NewSIRICheckStatusResponse(
	producerRef string,
	requestMessageRef string,
	responseMessageIdentifier string,
	status bool,
	responseTimestamp time.Time,
	serviceStartedTime time.Time) *SIRICheckStatusResponse {
	return &SIRICheckStatusResponse{
		ProducerRef:               producerRef,
		RequestMessageRef:         requestMessageRef,
		ResponseMessageIdentifier: responseMessageIdentifier,
		Status:             status,
		ResponseTimestamp:  responseTimestamp,
		ServiceStartedTime: serviceStartedTime}
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) ProducerRef() string {
	if request.producerRef == "" {
		nodes, _ := request.node.Search("//*[local-name()='ProducerRef']")
		request.producerRef = strings.TrimSpace(nodes[0].Content())
	}
	return request.producerRef
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) RequestMessageRef() string {
	if request.requestMessageRef == "" {
		nodes, _ := request.node.Search("//*[local-name()='RequestMessageRef']")
		request.requestMessageRef = strings.TrimSpace(nodes[0].Content())
	}
	return request.requestMessageRef
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) ResponseMessageIdentifier() string {
	if request.responseMessageIdentifier == "" {
		nodes, _ := request.node.Search("//*[local-name()='ResponseMessageIdentifier']")
		request.responseMessageIdentifier = strings.TrimSpace(nodes[0].Content())
	}
	return request.responseMessageIdentifier
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) Status() bool {
	if !request.status {
		nodes, _ := request.node.Search("//*[local-name()='Status']")
		s, _ := strconv.ParseBool(nodes[0].Content())
		request.status = s
	}
	return request.status
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) ResponseTimestamp() time.Time {
	if request.responseTimestamp.IsZero() {
		nodes, _ := request.node.Search("//*[local-name()='ResponseTimestamp']")
		t, _ := time.Parse("2006-01-02T15:04:05.000Z07:00", strings.TrimSpace(nodes[0].Content()))
		request.responseTimestamp = t
	}
	return request.responseTimestamp
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) ServiceStartedTime() time.Time {
	if request.serviceStartedTime.IsZero() {
		nodes, _ := request.node.Search("//*[local-name()='ServiceStartedTime']")
		t, _ := time.Parse("2006-01-02T15:04:05.000Z07:00", strings.TrimSpace(nodes[0].Content()))
		request.serviceStartedTime = t
	}
	return request.serviceStartedTime
}

// TODO : Handle errors
func (request *SIRICheckStatusResponse) BuildXML() string {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(SIRIResponseTemplate))
	if err := siriResponse.Execute(&buffer, request); err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}
