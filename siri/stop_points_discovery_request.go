package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopPointsDiscoveryRequest struct {
	RequestXMLStructure
}

type SIRIStopPointsDiscoveryRequest struct {
	MessageIdentifier string
	RequestorRef      string

	RequestTimestamp time.Time
}

const stopPointsDiscoveryRequestTemplate = `<sw:StopPointsDiscovery xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Request>
		<siri:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:RequestTimestamp>
		<siri:RequestorRef>{{.RequestorRef}}</siri:RequestorRef>
		<siri:MessageIdentifier>{{.MessageIdentifier}}</siri:MessageIdentifier>
	</Request>
	<RequestExtension />
</sw:StopPointsDiscovery>`

func NewXMLStopPointsDiscoveryRequest(node xml.Node) *XMLStopPointsDiscoveryRequest {
	xmlStopDiscoveryRequest := &XMLStopPointsDiscoveryRequest{}
	xmlStopDiscoveryRequest.node = NewXMLNode(node)
	return xmlStopDiscoveryRequest
}

func NewXMLStopPointsDiscoveryRequestFromContent(content []byte) (*XMLStopPointsDiscoveryRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLStopPointsDiscoveryRequest(doc.Root().XmlNode)
	return request, nil
}

func NewSIRIStopPointsDiscoveryRequest(messageIdentifier, requestorRef string, requestTimestamp time.Time) *SIRIStopPointsDiscoveryRequest {
	return &SIRIStopPointsDiscoveryRequest{
		MessageIdentifier: messageIdentifier,
		RequestorRef:      requestorRef,
		RequestTimestamp:  requestTimestamp,
	}
}

func (request *SIRIStopPointsDiscoveryRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(stopPointsDiscoveryRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
