package siri

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetGeneralMessage struct {
	XMLGeneralMessageRequest

	requestorRef string
}

type XMLGeneralMessageRequest struct {
	XMLStructure

	messageIdentifier string

	requestTimestamp time.Time

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

	// Filters are not used by Edwig for now, we always request all GM
	InfoChannelRef []string

	LineRef           []string
	StopPointRef      []string
	JourneyPatternRef []string
	DestinationRef    []string
	RouteRef          []string
	GroupOfLinesRef   []string
}

const getGeneralMessageRequestTemplate = `<sw:GetGeneralMessage xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri" xmlns:sws="http://wsdl.siri.org.uk/siri">
	<ServiceRequestInfo>
		<siri:RequestTimestamp>{{ .RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:RequestTimestamp>
		<siri:RequestorRef>{{ .RequestorRef }}</siri:RequestorRef>
		<siri:MessageIdentifier>{{ .MessageIdentifier }}</siri:MessageIdentifier>
	</ServiceRequestInfo>
	<Request version="2.0:FR-IDF-2.4">
		{{ .BuildGeneralMessageRequestXML }}
	</Request>
	<RequestExtension/>
</sw:GetGeneralMessage>`

const generalMessageRequestTemplate = `<siri:RequestTimestamp>{{ .RequestTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:RequestTimestamp>
		<siri:MessageIdentifier>{{ .MessageIdentifier }}</siri:MessageIdentifier>{{ range .InfoChannelRef }}
		<siri:InfoChannelRef>{{ . }}</siri:InfoChannelRef>{{ end }}
		<siri:Extensions>{{ if .XsdInWsdl }}
			<siri:IDFGeneralMessageRequestFilter>{{ else }}
			<sws:IDFGeneralMessageRequestFilter>{{ end }}{{ range .LineRef }}
				<siri:LineRef>{{ . }}</siri:LineRef>{{ end }}{{ range .StopPointRef }}
				<siri:StopPointRef>{{ . }}</siri:StopPointRef>{{ end }}{{ range .JourneyPatternRef }}
				<siri:JourneyPatternRef>{{ . }}</siri:JourneyPatternRef>{{ end }}{{ range .DestinationRef }}
				<siri:DestinationRef>{{ . }}</siri:DestinationRef>{{ end }}{{ range .RouteRef }}
				<siri:RouteRef>{{ . }}</siri:RouteRef>{{ end }}{{ range .GroupOfLinesRef }}
				<siri:GroupOfLinesRef>{{ . }}</siri:GroupOfLinesRef>{{ end }}{{ if .XsdInWsdl }}
			</siri:IDFGeneralMessageRequestFilter>{{ else }}
			</sws:IDFGeneralMessageRequestFilter>{{ end }}
		</siri:Extensions>`

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

func (request *XMLGeneralMessageRequest) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *XMLGeneralMessageRequest) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent("RequestTimestamp")
	}
	return request.requestTimestamp
}

func (request *XMLGeneralMessageRequest) InfoChannelRef() []string {
	if len(request.infoChannelRef) == 0 {
		nodes := request.findNodes("InfoChannelRef")
		if nodes != nil {
			for _, infoChannelRef := range nodes {
				request.infoChannelRef = append(request.infoChannelRef, strings.TrimSpace(infoChannelRef.NativeNode().Content()))
			}
		}
	}
	return request.infoChannelRef
}

func (request *XMLGeneralMessageRequest) GroupOfLinesRef() []string {
	if len(request.groupOfLinesRef) == 0 {
		nodes := request.findNodes("GroupOfLinesRef")
		if nodes != nil {
			for _, groupOfLinesRef := range nodes {
				request.groupOfLinesRef = append(request.groupOfLinesRef, strings.TrimSpace(groupOfLinesRef.NativeNode().Content()))
			}
		}
	}
	return request.groupOfLinesRef
}

func (request *XMLGeneralMessageRequest) RouteRef() []string {
	if len(request.routeRef) == 0 {
		nodes := request.findNodes("RouteRef")
		if nodes != nil {
			for _, routeRef := range nodes {
				request.routeRef = append(request.routeRef, strings.TrimSpace(routeRef.NativeNode().Content()))
			}
		}
	}
	return request.routeRef
}

func (request *XMLGeneralMessageRequest) DestinationRef() []string {
	if len(request.destinationRef) == 0 {
		nodes := request.findNodes("DestinationRef")
		if nodes != nil {
			for _, destinationRef := range nodes {
				request.destinationRef = append(request.destinationRef, strings.TrimSpace(destinationRef.NativeNode().Content()))
			}
		}
	}
	return request.destinationRef
}

func (request *XMLGeneralMessageRequest) JourneyPatternRef() []string {
	if len(request.journeyPatternRef) == 0 {
		nodes := request.findNodes("JourneyPatternRef")
		if nodes != nil {
			for _, journeyPatternRef := range nodes {
				request.journeyPatternRef = append(request.journeyPatternRef, strings.TrimSpace(journeyPatternRef.NativeNode().Content()))
			}
		}
	}
	return request.journeyPatternRef
}

func (request *XMLGeneralMessageRequest) StopPointRef() []string {
	if len(request.stopPointRef) == 0 {
		nodes := request.findNodes("StopPointRef")
		if nodes != nil {
			for _, stopPointRef := range nodes {
				request.stopPointRef = append(request.stopPointRef, strings.TrimSpace(stopPointRef.NativeNode().Content()))
			}
		}
	}
	return request.stopPointRef
}

func (request *XMLGeneralMessageRequest) LineRef() []string {
	if len(request.lineRef) == 0 {
		nodes := request.findNodes("LineRef")
		if nodes != nil {
			for _, lineRef := range nodes {
				request.lineRef = append(request.lineRef, strings.TrimSpace(lineRef.NativeNode().Content()))
			}
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

func (request *SIRIGetGeneralMessageRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("getGeneralMessageRequest").Parse(getGeneralMessageRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIGeneralMessageRequest) BuildGeneralMessageRequestXML() (string, error) {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("generalMessageRequest").Parse(generalMessageRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
