package siri

import (
	"bytes"
	"html/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGeneralMessageRequest struct {
	RequestXMLStructure
}

type SIRIGeneralMessageRequest struct {
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time
}

const generalMessageRequestTemplate = `<ns7:GetGeneralMessage xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://wsdl.siri.org.uk/siri" xmlns:ns7="http://wsdl.siri.org.uk">
      <ServiceRequestInfo>
        <ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z"}}</ns2:RequestTimestamp>
        <ns2:RequestorRef>{{.RequestorRef}}</ns2:RequestorRef>
        <ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
      </ServiceRequestInfo>
      <Request version="2.0:FR-IDF-2.4">
        <ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z"}}/ns2:RequestTimestamp>
        <ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
      </Request>
      <RequestExtension/>
</ns7:GetGeneralMessage>`

func NewXMLGeneralMessageRequest(node xml.Node) *XMLGeneralMessageRequest {
	xmlGeneralMessageRequest := &XMLGeneralMessageRequest{}
	xmlGeneralMessageRequest.node = NewXMLNode(node)
	return xmlGeneralMessageRequest
}

func NewXMLGeneralMessageRequestFromContent(content []byte) (*XMLGeneralMessageRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGeneralMessageRequest(doc.Root().XmlNode)
	return request, nil
}

func NewSIRIGeneralMessageRequest(
	messageIdentifier,
	requestorRef string,
	requestTimestamp time.Time) *SIRIGeneralMessageRequest {
	return &SIRIGeneralMessageRequest{
		MessageIdentifier: messageIdentifier,
		RequestorRef:      requestorRef,
		RequestTimestamp:  requestTimestamp,
	}
}

func (request *SIRIGeneralMessageRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(generalMessageRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
