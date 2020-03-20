package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
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
	if err := templates.ExecuteTemplate(&buffer, "check_status_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
