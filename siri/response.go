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

func NewXMLCheckStatusResponse(node *xml.XmlNode) *XMLCheckStatusResponse {
	return &XMLCheckStatusResponse{node: node}
}

func NewXMLCheckStatusResponseFromContent(content []byte) (*XMLCheckStatusResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLCheckStatusResponse(doc.Root().XmlNode)
	finalizer := func(request *XMLCheckStatusResponse) {
		doc.Free()
	}
	runtime.SetFinalizer(request, finalizer)
	return request, nil
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

// TODO : Handle errors
func (request *XMLCheckStatusResponse) Address() string {
	if request.address == "" {
		nodes, err := request.node.Search("//*[local-name()='Address']")
		if err != nil {
			log.Fatal(err)
		}
		request.address = strings.TrimSpace(nodes[0].Content())
	}
	return request.address
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) ProducerRef() string {
	if request.producerRef == "" {
		nodes, err := request.node.Search("//*[local-name()='ProducerRef']")
		if err != nil {
			log.Fatal(err)
		}
		request.producerRef = strings.TrimSpace(nodes[0].Content())
	}
	return request.producerRef
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) RequestMessageRef() string {
	if request.requestMessageRef == "" {
		nodes, err := request.node.Search("//*[local-name()='RequestMessageRef']")
		if err != nil {
			log.Fatal(err)
		}
		request.requestMessageRef = strings.TrimSpace(nodes[0].Content())
	}
	return request.requestMessageRef
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) ResponseMessageIdentifier() string {
	if request.responseMessageIdentifier == "" {
		nodes, err := request.node.Search("//*[local-name()='ResponseMessageIdentifier']")
		if err != nil {
			log.Fatal(err)
		}
		request.responseMessageIdentifier = strings.TrimSpace(nodes[0].Content())
	}
	return request.responseMessageIdentifier
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) Status() bool {
	if !request.status {
		nodes, err := request.node.Search("//*[local-name()='Status']")
		if err != nil {
			log.Fatal(err)
		}
		s, err := strconv.ParseBool(strings.TrimSpace(nodes[0].Content()))
		if err != nil {
			log.Fatal(err)
		}
		request.status = s
	}
	return request.status
}

// TODO : Handle errors and see what to do if status is true
func (request *XMLCheckStatusResponse) ErrorType() string {
	if !request.Status() && request.errorType == "" {
		nodes, err := request.node.Search("//*[local-name()='ErrorText']")
		if err != nil {
			log.Fatal(err)
		}
		request.errorType = nodes[0].Parent().Name()
	}
	return request.errorType
}

// TODO : Handle errors and see what to do if status is true
func (request *XMLCheckStatusResponse) ErrorNumber() int {
	if !request.Status() && request.ErrorType() == "OtherError" && request.errorNumber == 0 {
		nodes, err := request.node.Search("//*[local-name()='ErrorText']")
		if err != nil {
			log.Fatal(err)
		}
		n, err := strconv.Atoi(nodes[0].Parent().Attr("number"))
		if err != nil {
			log.Fatal(err)
		}
		request.errorNumber = n
	}
	return request.errorNumber
}

// TODO : Handle errors and see what to do if status is true
func (request *XMLCheckStatusResponse) ErrorText() string {
	if !request.Status() && request.errorText == "" {
		nodes, err := request.node.Search("//*[local-name()='ErrorText']")
		if err != nil {
			log.Fatal(err)
		}
		request.errorText = strings.TrimSpace(nodes[0].Content())
	}
	return request.errorText
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) ResponseTimestamp() time.Time {
	if request.responseTimestamp.IsZero() {
		nodes, err := request.node.Search("//*[local-name()='ResponseTimestamp']")
		if err != nil {
			log.Fatal(err)
		}
		t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", strings.TrimSpace(nodes[0].Content()))
		if err != nil {
			log.Fatal(err)
		}
		request.responseTimestamp = t
	}
	return request.responseTimestamp
}

// TODO : Handle errors
func (request *XMLCheckStatusResponse) ServiceStartedTime() time.Time {
	if request.serviceStartedTime.IsZero() {
		nodes, err := request.node.Search("//*[local-name()='ServiceStartedTime']")
		if err != nil {
			log.Fatal(err)
		}
		t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", strings.TrimSpace(nodes[0].Content()))
		if err != nil {
			log.Fatal(err)
		}
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
