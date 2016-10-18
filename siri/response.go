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
	XMLStructure

	address                   string
	producerRef               string
	requestMessageRef         string
	responseMessageIdentifier string
	status                    bool
	errorType                 string
	errorNumber               int
	errorText                 string
	responseTimestamp         time.Time
	serviceStartedTime        time.Time
}

type SIRICheckStatusResponse struct {
	Address                   string
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string
	Status                    bool
	ErrorType                 string
	ErrorNumber               int
	ErrorText                 string
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
		<ns2:Address>{{.Address}}</ns2:Address>
		<ns2:ResponseMessageIdentifier>{{.ResponseMessageIdentifier}}</ns2:ResponseMessageIdentifier>
		<ns2:RequestMessageRef>{{.RequestMessageRef}}</ns2:RequestMessageRef>
	</CheckStatusAnswerInfo>
	<Answer>
		<ns2:Status>{{.Status}}</ns2:Status>{{ if not .Status }}
		<ns2:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
			<ns2:OtherError number="{{.ErrorNumber}}">{{ else }}
			<ns2:{{.ErrorType}}>{{ end }}
				<ns2:ErrorText>{{.ErrorText}}</ns2:ErrorText>
			</ns2:ServiceNotAvailableError>
		</ns2:ErrorCondition>{{ end }}
		<ns2:ServiceStartedTime>{{.ServiceStartedTime.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:ServiceStartedTime>
	</Answer>
	<AnswerExtension />
</ns7:CheckStatusResponse>`

func NewXMLCheckStatusResponse(node xml.Node) *XMLCheckStatusResponse {
	return &XMLCheckStatusResponse{XMLStructure: XMLStructure{node: node}}
}

func NewXMLCheckStatusResponseFromContent(content []byte) (*XMLCheckStatusResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLCheckStatusResponse(doc.Root().XmlNode)
	finalizer := func(response *XMLCheckStatusResponse) {
		doc.Free()
	}
	runtime.SetFinalizer(response, finalizer)
	return response, nil
}

func NewSIRICheckStatusResponse(
	address string,
	producerRef string,
	requestMessageRef string,
	responseMessageIdentifier string,
	status bool,
	errorType string,
	errorNumber int,
	errorText string,
	responseTimestamp time.Time,
	serviceStartedTime time.Time) *SIRICheckStatusResponse {
	return &SIRICheckStatusResponse{
		Address:                   address,
		ProducerRef:               producerRef,
		RequestMessageRef:         requestMessageRef,
		ResponseMessageIdentifier: responseMessageIdentifier,
		Status:             status,
		ErrorType:          errorType,
		ErrorNumber:        errorNumber,
		ErrorText:          errorText,
		ResponseTimestamp:  responseTimestamp,
		ServiceStartedTime: serviceStartedTime}
}

func (response *XMLCheckStatusResponse) Address() string {
	if response.address == "" {
		response.address = response.findStringChildContent("Address")
	}
	return response.address
}

func (response *XMLCheckStatusResponse) ProducerRef() string {
	if response.producerRef == "" {
		response.producerRef = response.findStringChildContent("ProducerRef")
	}
	return response.producerRef
}

func (response *XMLCheckStatusResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		response.requestMessageRef = response.findStringChildContent("RequestMessageRef")
	}
	return response.requestMessageRef
}

func (response *XMLCheckStatusResponse) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == "" {
		response.responseMessageIdentifier = response.findStringChildContent("ResponseMessageIdentifier")
	}
	return response.responseMessageIdentifier
}

func (response *XMLCheckStatusResponse) Status() bool {
	if !response.status {
		response.status = response.findBoolChildContent("Status")
	}
	return response.status
}

// TODO: See what to do if status is true
// we can't access directly the node, we search for errorText and get parent
// Gokogiri FirstChild() or LastChild()  doesn't work
func (response *XMLCheckStatusResponse) ErrorType() string {
	if !response.Status() && response.errorType == "" {
		node := response.findNode("ErrorText")
		response.errorType = node.Parent().Name()

		// Find errorText and errorNumber to avoir too much parsing
		response.errorText = strings.TrimSpace(node.Content())
		if response.errorType == "OtherError" {
			n, err := strconv.Atoi(node.Parent().Attr("number"))
			if err != nil {
				log.Fatal(err)
			}
			response.errorNumber = n
		}
	}
	return response.errorType
}

// TODO: See what to do if status is true
func (response *XMLCheckStatusResponse) ErrorNumber() int {
	if !response.Status() && response.ErrorType() == "OtherError" && response.errorNumber == 0 {
		node := response.findNode("ErrorText")
		n, err := strconv.Atoi(node.Parent().Attr("number"))
		if err != nil {
			log.Fatal(err)
		}
		response.errorNumber = n
	}
	return response.errorNumber
}

// TODO: See what to do if status is true
func (response *XMLCheckStatusResponse) ErrorText() string {
	if !response.Status() && response.errorText == "" {
		response.errorText = response.findStringChildContent("ErrorText")
	}
	return response.errorText
}

func (response *XMLCheckStatusResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		response.responseTimestamp = response.findTimeChildContent("ResponseTimestamp")
	}
	return response.responseTimestamp
}

func (response *XMLCheckStatusResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime.IsZero() {
		response.serviceStartedTime = response.findTimeChildContent("ServiceStartedTime")
	}
	return response.serviceStartedTime
}

// TODO : Handle errors
func (response *SIRICheckStatusResponse) BuildXML() string {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(SIRIResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		log.Fatal(err)
	}
	return buffer.String()
}
