package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/edwig/logger"
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
	if err := templates.ExecuteTemplate(&buffer, "stop_points_discovery_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
