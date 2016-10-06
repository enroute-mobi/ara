package siri

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/api"
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
	api.UUIDConsumer
	api.ClockConsumer

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

// TODO : Handle errors
func (response *XMLCheckStatusResponse) Address() string {
	if response.address == "" {
		nodes, err := response.node.Search("//*[local-name()='Address']")
		if err != nil {
			log.Fatal(err)
		}
		response.address = strings.TrimSpace(nodes[0].Content())
	}
	return response.address
}

// TODO : Handle errors
func (response *XMLCheckStatusResponse) ProducerRef() string {
	if response.producerRef == "" {
		nodes, err := response.node.Search("//*[local-name()='ProducerRef']")
		if err != nil {
			log.Fatal(err)
		}
		response.producerRef = strings.TrimSpace(nodes[0].Content())
	}
	return response.producerRef
}

// TODO : Handle errors
func (response *XMLCheckStatusResponse) RequestMessageRef() string {
	if response.requestMessageRef == "" {
		nodes, err := response.node.Search("//*[local-name()='RequestMessageRef']")
		if err != nil {
			log.Fatal(err)
		}
		response.requestMessageRef = strings.TrimSpace(nodes[0].Content())
	}
	return response.requestMessageRef
}

// TODO : Handle errors
func (response *XMLCheckStatusResponse) ResponseMessageIdentifier() string {
	if response.responseMessageIdentifier == "" {
		nodes, err := response.node.Search("//*[local-name()='ResponseMessageIdentifier']")
		if err != nil {
			log.Fatal(err)
		}
		response.responseMessageIdentifier = strings.TrimSpace(nodes[0].Content())
	}
	return response.responseMessageIdentifier
}

// TODO : Handle errors
func (response *XMLCheckStatusResponse) Status() bool {
	if !response.status {
		nodes, err := response.node.Search("//*[local-name()='Status']")
		if err != nil {
			log.Fatal(err)
		}
		s, err := strconv.ParseBool(strings.TrimSpace(nodes[0].Content()))
		if err != nil {
			log.Fatal(err)
		}
		response.status = s
	}
	return response.status
}

// TODO : Handle errors and see what to do if status is true
func (response *XMLCheckStatusResponse) ErrorType() string {
	if !response.Status() && response.errorType == "" {
		nodes, err := response.node.Search("//*[local-name()='ErrorText']")
		if err != nil {
			log.Fatal(err)
		}
		response.errorType = nodes[0].Parent().Name()
	}
	return response.errorType
}

// TODO : Handle errors and see what to do if status is true
func (response *XMLCheckStatusResponse) ErrorNumber() int {
	if !response.Status() && response.ErrorType() == "OtherError" && response.errorNumber == 0 {
		nodes, err := response.node.Search("//*[local-name()='ErrorText']")
		if err != nil {
			log.Fatal(err)
		}
		n, err := strconv.Atoi(nodes[0].Parent().Attr("number"))
		if err != nil {
			log.Fatal(err)
		}
		response.errorNumber = n
	}
	return response.errorNumber
}

// TODO : Handle errors and see what to do if status is true
func (response *XMLCheckStatusResponse) ErrorText() string {
	if !response.Status() && response.errorText == "" {
		nodes, err := response.node.Search("//*[local-name()='ErrorText']")
		if err != nil {
			log.Fatal(err)
		}
		response.errorText = strings.TrimSpace(nodes[0].Content())
	}
	return response.errorText
}

// TODO : Handle errors
func (response *XMLCheckStatusResponse) ResponseTimestamp() time.Time {
	if response.responseTimestamp.IsZero() {
		nodes, err := response.node.Search("//*[local-name()='ResponseTimestamp']")
		if err != nil {
			log.Fatal(err)
		}
		t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", strings.TrimSpace(nodes[0].Content()))
		if err != nil {
			log.Fatal(err)
		}
		response.responseTimestamp = t
	}
	return response.responseTimestamp
}

// TODO : Handle errors
func (response *XMLCheckStatusResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime.IsZero() {
		nodes, err := response.node.Search("//*[local-name()='ServiceStartedTime']")
		if err != nil {
			log.Fatal(err)
		}
		t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", strings.TrimSpace(nodes[0].Content()))
		if err != nil {
			log.Fatal(err)
		}
		response.serviceStartedTime = t
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

func (response *SIRICheckStatusResponse) GenerateMessageIdentifier() {
	response.ResponseMessageIdentifier = fmt.Sprintf("Edwig:ResponseMessage::%s:LOC", response.NewUUID())
}

func (response *SIRICheckStatusResponse) SetResponseTimestamp() {
	response.ResponseTimestamp = response.Clock().Now()
}
