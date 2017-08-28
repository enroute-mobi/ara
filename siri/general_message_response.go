package siri

import (
	"bytes"
	"strings"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGeneralMessageResponse struct {
	ResponseXMLStructure

	xmlGeneralMessages []*XMLGeneralMessage
}

type XMLGeneralMessageCancellation struct {
	XMLStructure

	infoMessageIdentifier string
}

type XMLGeneralMessage struct {
	XMLStructure

	itemIdentifier        string
	infoMessageIdentifier string
	infoChannelRef        string
	formatRef             string

	infoMessageVersion int

	recordedAtTime time.Time
	validUntilTime time.Time

	content interface{}
}

type IDFGeneralMessageStructure struct {
	XMLStructure

	lineRef           []string
	stopPointRef      []string
	journeyPatternRef []string
	destinationRef    []string
	routeRef          []string
	format            []string
	groupOfLinesRef   []string

	lineSections []*IDFLineSectionStructure
	messages     []*XMLMessage
}

type XMLMessage struct {
	XMLStructure

	messageText         string
	messageType         string
	numberOfLines       int
	numberOfCharPerLine int
}

type IDFLineSectionStructure struct {
	XMLStructure

	firstStop string
	lastStop  string
	lineRef   string
}

type SIRIGeneralMessageResponse struct {
	SIRIGeneralMessageDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
}

type SIRIGeneralMessageDelivery struct {
	RequestMessageRef string
	Status            bool
	ErrorType         string
	ErrorNumber       int
	ErrorText         string
	ResponseTimestamp time.Time

	GeneralMessages []*SIRIGeneralMessage
}

type SIRIGeneralMessage struct {
	RecordedAtTime        time.Time
	ValidUntilTime        time.Time
	ItemIdentifier        string
	InfoMessageIdentifier string
	FormatRef             string
	InfoMessageVersion    int
	InfoChannelRef        string

	LineRef           []string
	StopPointRef      []string
	JourneyPatternRef []string
	DestinationRef    []string
	RouteRef          []string
	GroupOfLinesRef   []string

	LineSections []*SIRILineSection
	Messages     []*SIRIMessage
}

type SIRILineSection struct {
	FirstStop string
	LastStop  string
	LineRef   string
}

type SIRIMessage struct {
	Content             string
	Type                string
	NumberOfLines       int
	NumberOfCharPerLine int
}

const generalMessageResponseTemplate = `<ns8:GetGeneralMessageResponse xmlns:ns3="http://www.siri.org.uk/siri"
															 xmlns:ns4="http://www.ifopt.org.uk/acsb"
															 xmlns:ns5="http://www.ifopt.org.uk/ifopt"
															 xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
															 xmlns:ns7="http://scma/siri"
															 xmlns:ns8="http://wsdl.siri.org.uk"
															 xmlns:ns9="http://wsdl.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
		<ns3:ProducerRef>{{ .ProducerRef }}</ns3:ProducerRef>{{ if .Address }}
		<ns3:Address>{{ .Address }}</ns3:Address>{{ end }}
		<ns3:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		{{ .BuildGeneralMessageDeliveryXML }}
	</Answer>
	<AnswerExtension/>
</ns8:GetGeneralMessageResponse>`

const generalMessageDeliveryTemplate = `<ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
			<ns3:Status>{{.Status}}</ns3:Status>{{ if not .Status }}
			<ns3:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<ns3:OtherError number="{{.ErrorNumber}}">{{ else }}
				<ns3:{{.ErrorType}}>{{ end }}
					<ns3:ErrorText>{{.ErrorText}}</ns3:ErrorText>
				</ns3:{{.ErrorType}}>
			</ns3:ErrorCondition>{{ else }}{{range .GeneralMessages}}
			{{ .BuildGeneralMessageXML }}{{end}}{{end}}
		</ns3:GeneralMessageDelivery>`

const generalMessageTemplate = `<ns3:GeneralMessage>
				<ns3:formatRef>{{ .FormatRef }}</ns3:formatRef>
				<ns3:RecordedAtTime>{{ .RecordedAtTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:RecordedAtTime>
				<ns3:ItemIdentifier>{{ .ItemIdentifier }}</ns3:ItemIdentifier>
				<ns3:InfoMessageIdentifier>{{ .InfoMessageIdentifier }}</ns3:InfoMessageIdentifier>
				<ns3:InfoMessageVersion>{{ .InfoMessageVersion }}</ns3:InfoMessageVersion>
				<ns3:InfoChannelRef>{{ .InfoChannelRef }}</ns3:InfoChannelRef>
				<ns3:ValidUntilTime>{{ .ValidUntilTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ValidUntilTime>
				<ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
				xsi:type="ns9:IDFLineSectionStructure">{{range .LineRef }}
					<ns3:LineRef>{{ . }}</ns3:LineRef>{{end}}{{range .StopPointRef }}
					<ns3:StopPointRef>{{ . }}</ns3:StopPointRef>{{end}}{{range .JourneyPatternRef }}
					<ns3:JourneyPatternRef>{{ . }}</ns3:JourneyPatternRef>{{end}}{{range .DestinationRef }}
					<ns3:DestinationRef>{{ . }}</ns3:DestinationRef>{{end}}{{range .RouteRef }}
					<ns3:RouteRef>{{ . }}</ns3:RouteRef>{{end}}{{range .GroupOfLinesRef }}
					<ns3:GroupOfLinesRef>{{ . }}</ns3:GroupOfLinesRef>{{end}}{{ range .LineSections }}
					<ns3:LineSection>{{ if .FirstStop }}
						<ns3:FirstStop>{{ .FirstStop }}</ns3:FirstStop>{{end}}{{if .LastStop }}
						<ns3:LastStop>{{ .LastStop }}</ns3:LastStop>{{end}}{{if .LineRef }}
						<ns3:LineRef>{{ .LineRef }}</ns3:LineRef>{{end}}
					</ns3:LineSection>{{end}}{{range .Messages}}
					<ns3:Message>{{if .Type}}
						<ns3:MessageType>{{ .Type }}</ns3:MessageType>{{end}}{{if .Content }}
						<ns3:MessageText>{{ .Content }}</ns3:MessageText>{{end}}{{if .NumberOfLines }}
						<ns3:NumberOfLines>{{ .NumberOfLines }}</ns3:NumberOfLines>{{end}}{{if .NumberOfCharPerLine }}
						<ns3:NumberOfCharPerLine>{{ .NumberOfCharPerLine }}</ns3:NumberOfCharPerLine>{{end}}
					</ns3:Message>{{end}}
				</ns3:Content>
			</ns3:GeneralMessage>`

func NewXMLGeneralMessageResponseFromContent(content []byte) (*XMLGeneralMessageResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLGeneralMessageResponse(doc.Root().XmlNode)
	return response, nil
}

func NewXMLGeneralMessageResponse(node xml.Node) *XMLGeneralMessageResponse {
	xmlGeneralMessageResponse := &XMLGeneralMessageResponse{}
	xmlGeneralMessageResponse.node = NewXMLNode(node)
	return xmlGeneralMessageResponse
}

func NewXMLCancelledGeneralMessage(node XMLNode) *XMLGeneralMessageCancellation {
	cancelledGeneralMessage := &XMLGeneralMessageCancellation{}
	cancelledGeneralMessage.node = node
	return cancelledGeneralMessage
}

func NewXMLGeneralMessage(node XMLNode) *XMLGeneralMessage {
	generalMessage := &XMLGeneralMessage{}
	generalMessage.node = node
	return generalMessage
}

func NewXMLLineSection(node XMLNode) *IDFLineSectionStructure {
	lineSection := &IDFLineSectionStructure{}
	lineSection.node = node
	return lineSection
}

func NewXMLMessage(node XMLNode) *XMLMessage {
	message := &XMLMessage{}
	message.node = node
	return message
}

func (response *XMLGeneralMessageResponse) XMLGeneralMessages() []*XMLGeneralMessage {
	if len(response.xmlGeneralMessages) == 0 {
		nodes := response.findNodes("GeneralMessage")
		if nodes == nil {
			return response.xmlGeneralMessages
		}
		for _, generalMessage := range nodes {
			response.xmlGeneralMessages = append(response.xmlGeneralMessages, NewXMLGeneralMessage(generalMessage))
		}
	}
	return response.xmlGeneralMessages
}

func (visit *XMLGeneralMessageCancellation) InfoMessageIdentifier() string {
	if visit.infoMessageIdentifier == "" {
		visit.infoMessageIdentifier = visit.findStringChildContent("InfoMessageIdentifier")
	}
	return visit.infoMessageIdentifier
}

func (visit *XMLGeneralMessage) RecordedAtTime() time.Time {
	if visit.recordedAtTime.IsZero() {
		visit.recordedAtTime = visit.findTimeChildContent("RecordedAtTime")
	}
	return visit.recordedAtTime
}

func (visit *XMLGeneralMessage) ValidUntilTime() time.Time {
	if visit.validUntilTime.IsZero() {
		visit.validUntilTime = visit.findTimeChildContent("ValidUntilTime")
	}
	return visit.validUntilTime
}

func (visit *XMLGeneralMessage) ItemIdentifier() string {
	if visit.itemIdentifier == "" {
		visit.itemIdentifier = visit.findStringChildContent("ItemIdentifier")
	}
	return visit.itemIdentifier
}

func (visit *XMLGeneralMessage) InfoMessageIdentifier() string {
	if visit.infoMessageIdentifier == "" {
		visit.infoMessageIdentifier = visit.findStringChildContent("InfoMessageIdentifier")
	}
	return visit.infoMessageIdentifier
}

func (visit *XMLGeneralMessage) InfoMessageVersion() int {
	if visit.infoMessageVersion == 0 {
		visit.infoMessageVersion = visit.findIntChildContent("InfoMessageVersion")
	}
	return visit.infoMessageVersion
}

func (visit *XMLGeneralMessage) InfoChannelRef() string {
	if visit.infoChannelRef == "" {
		visit.infoChannelRef = visit.findStringChildContent("InfoChannelRef")
	}
	return visit.infoChannelRef
}

func (visit *XMLGeneralMessage) FormatRef() string {
	if visit.formatRef == "" {
		visit.formatRef = visit.findStringChildContent("formatRef")
	}
	return visit.formatRef
}

func (visit *XMLGeneralMessage) createNewContent() IDFGeneralMessageStructure {
	content := IDFGeneralMessageStructure{}
	content.node = NewXMLNode(visit.findNode("Content"))
	return content
}

func (visit *XMLGeneralMessage) Content() interface{} {
	if visit.content != nil {
		return visit.content
	}
	visit.content = visit.createNewContent()
	return visit.content
}

func (visit *IDFGeneralMessageStructure) GroupOfLinesRef() []string {
	if len(visit.groupOfLinesRef) == 0 {
		nodes := visit.findNodes("GroupOfLinesRef")
		if nodes != nil {
			for _, groupOfLinesRef := range nodes {
				visit.groupOfLinesRef = append(visit.groupOfLinesRef, strings.TrimSpace(groupOfLinesRef.NativeNode().Content()))
			}
		}
	}
	return visit.groupOfLinesRef
}

func (visit *IDFGeneralMessageStructure) RouteRef() []string {
	if len(visit.routeRef) == 0 {
		nodes := visit.findNodes("RouteRef")
		if nodes != nil {
			for _, routeRef := range nodes {
				visit.routeRef = append(visit.routeRef, strings.TrimSpace(routeRef.NativeNode().Content()))
			}
		}
	}
	return visit.routeRef
}

func (visit *IDFGeneralMessageStructure) DestinationRef() []string {
	if len(visit.destinationRef) == 0 {
		nodes := visit.findNodes("DestinationRef")
		if nodes != nil {
			for _, destinationRef := range nodes {
				visit.destinationRef = append(visit.destinationRef, strings.TrimSpace(destinationRef.NativeNode().Content()))
			}
		}
	}
	return visit.destinationRef
}

func (visit *IDFGeneralMessageStructure) JourneyPatternRef() []string {
	if len(visit.journeyPatternRef) == 0 {
		nodes := visit.findNodes("JourneyPatternRef")
		if nodes != nil {
			for _, journeyPatternRef := range nodes {
				visit.journeyPatternRef = append(visit.journeyPatternRef, strings.TrimSpace(journeyPatternRef.NativeNode().Content()))
			}
		}
	}
	return visit.journeyPatternRef
}

func (visit *IDFGeneralMessageStructure) StopPointRef() []string {
	if len(visit.stopPointRef) == 0 {
		nodes := visit.findNodes("StopPointRef")
		if nodes != nil {
			for _, stopPointRef := range nodes {
				visit.stopPointRef = append(visit.stopPointRef, strings.TrimSpace(stopPointRef.NativeNode().Content()))
			}
		}
	}
	return visit.stopPointRef
}

func (visit *IDFGeneralMessageStructure) LineRef() []string {
	if len(visit.lineRef) == 0 {
		nodes := visit.findNodes("LineRef")
		if nodes != nil {
			for _, lineRef := range nodes {
				visit.lineRef = append(visit.lineRef, strings.TrimSpace(lineRef.NativeNode().Content()))
			}
		}
	}
	return visit.lineRef
}

func (visit *IDFGeneralMessageStructure) LineSections() []*IDFLineSectionStructure {
	if len(visit.lineSections) == 0 {
		nodes := visit.findNodes("LineSection")
		if nodes != nil {
			for _, lineNode := range nodes {
				visit.lineSections = append(visit.lineSections, NewXMLLineSection(lineNode))
			}
		}
	}
	return visit.lineSections
}

func (visit *IDFGeneralMessageStructure) Messages() []*XMLMessage {
	if len(visit.messages) == 0 {
		nodes := visit.findNodes("Message")
		if nodes != nil {
			for _, messageNode := range nodes {
				visit.messages = append(visit.messages, NewXMLMessage(messageNode))
			}
		}
	}
	return visit.messages
}

func (visit *IDFLineSectionStructure) FirstStop() string {
	if visit.firstStop == "" {
		visit.firstStop = visit.findStringChildContent("FirstStop")
	}
	return visit.firstStop
}

func (visit *IDFLineSectionStructure) LastStop() string {
	if visit.lastStop == "" {
		visit.lastStop = visit.findStringChildContent("LastStop")
	}
	return visit.lastStop
}

func (visit *IDFLineSectionStructure) LineRef() string {
	if visit.lineRef == "" {
		visit.lineRef = visit.findStringChildContent("LineRef")
	}
	return visit.lineRef
}

func (message *XMLMessage) MessageText() string {
	if message.messageText == "" {
		message.messageText = message.findStringChildContent("MessageText")
	}
	return message.messageText
}

func (message *XMLMessage) MessageType() string {
	if message.messageType == "" {
		message.messageType = message.findStringChildContent("MessageType")
	}
	return message.messageType
}

func (message *XMLMessage) NumberOfLines() int {
	if message.numberOfLines == 0 {
		message.numberOfLines = message.findIntChildContent("NumberOfLines")
	}
	return message.numberOfLines
}

func (message *XMLMessage) NumberOfCharPerLine() int {
	if message.numberOfCharPerLine == 0 {
		message.numberOfCharPerLine = message.findIntChildContent("NumberOfCharPerLine")
	}
	return message.numberOfCharPerLine
}

func (response *SIRIGeneralMessageResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var generalMessage = template.Must(template.New("generalMessageResponse").Parse(generalMessageResponseTemplate))
	if err := generalMessage.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIGeneralMessageDelivery) BuildGeneralMessageDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	var generalMessageDelivery = template.Must(template.New("generalMessageDelivery").Parse(generalMessageDeliveryTemplate))
	if err := generalMessageDelivery.Execute(&buffer, delivery); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (message *SIRIGeneralMessage) BuildGeneralMessageXML() (string, error) {
	var buffer bytes.Buffer
	var generalMessage = template.Must(template.New("generalMessage").Parse(generalMessageTemplate))
	if err := generalMessage.Execute(&buffer, message); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
