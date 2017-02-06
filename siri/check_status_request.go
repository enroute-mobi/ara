package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/af83/edwig/logger"
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

const checkStatusRequestTemplate = `<ns7:CheckStatus xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
	<Request>
		<ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z"}}</ns2:RequestTimestamp>
		<ns2:RequestorRef>{{.RequestorRef}}</ns2:RequestorRef>
		<ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
	</Request>
	<RequestExtension/>
</ns7:CheckStatus>`

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
func (request *SIRICheckStatusRequest) BuildXML() string {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(checkStatusRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		logger.Log.Panicf("Error while using request template: %v", err)
	}
	return buffer.String()
}
