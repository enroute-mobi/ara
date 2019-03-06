package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLCheckStatusResponse struct {
	ResponseXMLStructureWithStatus

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

const checkStatusResponseTemplate = `<sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<CheckStatusAnswerInfo>
		<siri:ResponseTimestamp>{{.ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:ResponseTimestamp>
		<siri:ProducerRef>{{.ProducerRef}}</siri:ProducerRef>{{ if .Address }}
		<siri:Address>{{ .Address }}</siri:Address>{{ end }}
		<siri:ResponseMessageIdentifier>{{.ResponseMessageIdentifier}}</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>{{.RequestMessageRef}}</siri:RequestMessageRef>
	</CheckStatusAnswerInfo>
	<Answer>
		<siri:Status>{{.Status}}</siri:Status>{{ if not .Status }}
		<siri:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
			<siri:OtherError number="{{.ErrorNumber}}">{{ else }}
			<siri:{{.ErrorType}}>{{ end }}
				<siri:ErrorText>{{.ErrorText}}</siri:ErrorText>
			</siri:{{.ErrorType}}>
		</siri:ErrorCondition>{{ end }}
		<siri:ServiceStartedTime>{{.ServiceStartedTime.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:ServiceStartedTime>
	</Answer>
	<AnswerExtension/>
</sw:CheckStatusResponse>`

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
