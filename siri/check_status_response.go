package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLCheckStatusResponse struct {
	ResponseXMLStructure

	serviceStartedTime time.Time
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

const checkStatusResponseTemplate = `<ns7:CheckStatusResponse xmlns:ns2="http://www.siri.org.uk/siri"
												 xmlns:ns3="http://www.ifopt.org.uk/acsb"
												 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
												 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
												 xmlns:ns6="http://scma/siri"
												 xmlns:ns7="http://wsdl.siri.org.uk"
												 xmlns:ns8="http://wsdl.siri.org.uk/siri">
	<CheckStatusAnswerInfo>
		<ns2:ResponseTimestamp>{{.ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns2:ResponseTimestamp>
		<ns2:ProducerRef>{{.ProducerRef}}</ns2:ProducerRef>{{ if .Address }}
		<ns2:Address>{{ .Address }}</ns2:Address>{{ end }}
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
	xmlCheckStatusResponse := &XMLCheckStatusResponse{}
	xmlCheckStatusResponse.node = NewXMLNode(node)
	return xmlCheckStatusResponse
}

func NewXMLCheckStatusResponseFromContent(content []byte) (*XMLCheckStatusResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLCheckStatusResponse(doc.Root().XmlNode)
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

func (response *XMLCheckStatusResponse) ServiceStartedTime() time.Time {
	if response.serviceStartedTime.IsZero() {
		response.serviceStartedTime = response.findTimeChildContent("ServiceStartedTime")
	}
	return response.serviceStartedTime
}

// TODO : Handle errors
func (response *SIRICheckStatusResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(checkStatusResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
