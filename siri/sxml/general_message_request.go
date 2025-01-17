package sxml

import (
	"strings"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
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
	infoChannelRef []string
	lineRef        []string
	stopPointRef   []string
	destinationRef []string
	routeRef       []string
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
		request.requestorRef = request.findStringChildContent(siri_attributes.RequestorRef)
	}
	return request.requestorRef
}

func (request *XMLGeneralMessageRequest) InfoChannelRef() []string {
	if len(request.infoChannelRef) == 0 {
		nodes := request.findNodes(siri_attributes.InfoChannelRef)
		for _, infoChannelRef := range nodes {
			request.infoChannelRef = append(request.infoChannelRef, strings.TrimSpace(infoChannelRef.NativeNode().Content()))
		}
	}
	return request.infoChannelRef
}

func (request *XMLGeneralMessageRequest) RouteRef() []string {
	if len(request.routeRef) == 0 {
		nodes := request.findNodes(siri_attributes.RouteRef)
		for _, routeRef := range nodes {
			request.routeRef = append(request.routeRef, strings.TrimSpace(routeRef.NativeNode().Content()))
		}
	}
	return request.routeRef
}

func (request *XMLGeneralMessageRequest) DestinationRef() []string {
	if len(request.destinationRef) == 0 {
		nodes := request.findNodes(siri_attributes.DestinationRef)
		for _, destinationRef := range nodes {
			request.destinationRef = append(request.destinationRef, strings.TrimSpace(destinationRef.NativeNode().Content()))
		}
	}
	return request.destinationRef
}

func (request *XMLGeneralMessageRequest) StopPointRef() []string {
	if len(request.stopPointRef) == 0 {
		nodes := request.findNodes(siri_attributes.StopPointRef)
		for _, stopPointRef := range nodes {
			request.stopPointRef = append(request.stopPointRef, strings.TrimSpace(stopPointRef.NativeNode().Content()))
		}
	}
	return request.stopPointRef
}

func (request *XMLGeneralMessageRequest) LineRef() []string {
	if len(request.lineRef) == 0 {
		nodes := request.findNodes(siri_attributes.LineRef)
		for _, lineRef := range nodes {
			request.lineRef = append(request.lineRef, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return request.lineRef
}
