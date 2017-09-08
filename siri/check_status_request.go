package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLCheckStatusRequest struct {
	RequestXMLStructure
}

type SIRICheckStatusRequest struct {
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time
}

const checkStatusRequestTemplate = `<sw:CheckStatus xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Request>
		<siri:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:RequestTimestamp>
		<siri:RequestorRef>{{.RequestorRef}}</siri:RequestorRef>
		<siri:MessageIdentifier>{{.MessageIdentifier}}</siri:MessageIdentifier>
	</Request>
</sw:CheckStatus>`

func NewXMLCheckStatusRequest(node xml.Node) *XMLCheckStatusRequest {
	xmlCheckStatusRequest := &XMLCheckStatusRequest{}
	xmlCheckStatusRequest.node = NewXMLNode(node)
	return xmlCheckStatusRequest
}

func NewXMLCheckStatusRequestFromContent(content []byte) (*XMLCheckStatusRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLCheckStatusRequest(doc.Root().XmlNode)
	return request, nil
}

func NewSIRICheckStatusRequest(
	RequestorRef string,
	RequestTimestamp time.Time,
	MessageIdentifier string) *SIRICheckStatusRequest {
	return &SIRICheckStatusRequest{
		RequestorRef:      RequestorRef,
		RequestTimestamp:  RequestTimestamp,
		MessageIdentifier: MessageIdentifier,
	}
}

// TODO : Handle errors
func (request *SIRICheckStatusRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(checkStatusRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
