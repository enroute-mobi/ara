package siri

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetGeneralMessage struct {
	XMLGeneralMessageRequest

	requestorRef string
}

type XMLGeneralMessageRequest struct {
	LightRequestXMLStructure

	// Filters
	infoChannelRef    []string
	lineRef           []string
	stopPointRef      []string
	journeyPatternRef []string
	destinationRef    []string
	routeRef          []string
	groupOfLinesRef   []string
}

type SIRIGetGeneralMessageRequest struct {
	SIRIGeneralMessageRequest

	RequestorRef string
}

type SIRIGeneralMessageRequest struct {
	XsdInWsdl bool

	MessageIdentifier string

	RequestTimestamp time.Time

	InfoChannelRef []string

	LineRef           []string
	StopPointRef      []string
	JourneyPatternRef []string
	DestinationRef    []string
	RouteRef          []string
	GroupOfLinesRef   []string
}

func NewXMLGetGeneralMessage(node xml.Node) *XMLGetGeneralMessage {
	xmlGeneralMessageRequest := &XMLGetGeneralMessage{}
	xmlGeneralMessageRequest.node = NewXMLNode(node)
	return xmlGeneralMessageRequest
}

func NewXMLGetGeneralMessageFromContent(content []byte) (*XMLGetGeneralMessage, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetGeneralMessage(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLGetGeneralMessage) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent("RequestorRef")
	}
	return request.requestorRef
}

func (request *XMLGeneralMessageRequest) InfoChannelRef() []string {
	if len(request.infoChannelRef) == 0 {
		nodes := request.findNodes("InfoChannelRef")
		for _, infoChannelRef := range nodes {
			request.infoChannelRef = append(request.infoChannelRef, strings.TrimSpace(infoChannelRef.NativeNode().Content()))
		}
	}
	return request.infoChannelRef
}

func (request *XMLGeneralMessageRequest) GroupOfLinesRef() []string {
	if len(request.groupOfLinesRef) == 0 {
		nodes := request.findNodes("GroupOfLinesRef")
		for _, groupOfLinesRef := range nodes {
			request.groupOfLinesRef = append(request.groupOfLinesRef, strings.TrimSpace(groupOfLinesRef.NativeNode().Content()))
		}
	}
	return request.groupOfLinesRef
}

func (request *XMLGeneralMessageRequest) RouteRef() []string {
	if len(request.routeRef) == 0 {
		nodes := request.findNodes("RouteRef")
		for _, routeRef := range nodes {
			request.routeRef = append(request.routeRef, strings.TrimSpace(routeRef.NativeNode().Content()))
		}
	}
	return request.routeRef
}

func (request *XMLGeneralMessageRequest) DestinationRef() []string {
	if len(request.destinationRef) == 0 {
		nodes := request.findNodes("DestinationRef")
		for _, destinationRef := range nodes {
			request.destinationRef = append(request.destinationRef, strings.TrimSpace(destinationRef.NativeNode().Content()))
		}
	}
	return request.destinationRef
}

func (request *XMLGeneralMessageRequest) JourneyPatternRef() []string {
	if len(request.journeyPatternRef) == 0 {
		nodes := request.findNodes("JourneyPatternRef")
		for _, journeyPatternRef := range nodes {
			request.journeyPatternRef = append(request.journeyPatternRef, strings.TrimSpace(journeyPatternRef.NativeNode().Content()))
		}
	}
	return request.journeyPatternRef
}

func (request *XMLGeneralMessageRequest) StopPointRef() []string {
	if len(request.stopPointRef) == 0 {
		nodes := request.findNodes("StopPointRef")
		for _, stopPointRef := range nodes {
			request.stopPointRef = append(request.stopPointRef, strings.TrimSpace(stopPointRef.NativeNode().Content()))
		}
	}
	return request.stopPointRef
}

func (request *XMLGeneralMessageRequest) LineRef() []string {
	if len(request.lineRef) == 0 {
		nodes := request.findNodes("LineRef")
		for _, lineRef := range nodes {
			request.lineRef = append(request.lineRef, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return request.lineRef
}

func NewSIRIGeneralMessageRequest(
	messageIdentifier,
	requestorRef string,
	requestTimestamp time.Time) *SIRIGetGeneralMessageRequest {
	request := &SIRIGetGeneralMessageRequest{
		RequestorRef: requestorRef,
	}
	request.MessageIdentifier = messageIdentifier
	request.RequestTimestamp = requestTimestamp
	return request
}

func (request *SIRIGetGeneralMessageRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("get_general_message_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIGeneralMessageRequest) BuildGeneralMessageRequestXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "general_message_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
