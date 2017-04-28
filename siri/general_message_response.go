package siri

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGeneralMessageResponse struct {
	ResponseXMLStructure

	xmlGeneralMessages []*XMLGeneralMessage
}

type XMLGeneralMessage struct {
	XMLStructure

	recordedAtTime        time.Time
	validUntilTime        time.Time
	itemIdentifier        string
	infoMessageIdentifier string
	infoChannelRef        string
	format                string
	infoMessageVersion    int
	numberOfLines         int
	numberOfCharPerLine   int
	content               interface{}
}

type IDFGeneralMessageStructure struct {
	XMLStructure

	messages          []*XMLMessage
	lineRef           string
	stopPointRef      string
	journeyPatternRef string
	destinationRef    string
	routeRef          string
	format            string
	groupOfLinesRef   string
	lineSection       IDFLineSectionStructure
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
	InfoMessageVersion    int64
	InfoChannelRef        string

	LineRefContent    string
	StopPointRef      string
	JourneyPatternRef string
	DestinationRef    string
	RouteRef          string
	GroupOfLinesRef   string

	Messages []*SIRIMessage

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

const generalMessageTemplate = `<ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
				  <ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
					<ns3:Status>{{.Status}}</ns3:Status>{{range .GeneralMessages}}
					<ns3:GeneralMessage formatRef="FRANCE">
						<ns3:RecordedAtTime>{{ .RecordedAtTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:RecordedAtTime>
						<ns3:InfoMessageIdentifier>{{.InfoMessageIdentifier}}</ns3:InfoMessageIdentifier>
						<ns3:InfoMessageVersion>{{.InfoMessageVersion}}</ns3:InfoMessageVersion>
						<ns3:InfoChannelRef>{{.InfoChannelRef}}</ns3:InfoChannelRef>
						<ns3:ValidUntilTime>{{ .ValidUntilTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ValidUntilTime>
						<ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
						xsi:type="ns9:IDFLineSectionStructure">{{range .Messages}}
							<Message>{{if .Type}}
								<MessageType>{{ .Type }}</MessageType>{{end}}{{if .Content}}
								<MessageText>{{ .Content }}</MessageText>{{end}}{{if .NumberOfLines}}
								<NumberOfLines>{{ .NumberOfLines }}</NumberOfLines>{{end}}{{if .NumberOfCharPerLine}}
								<NumberOfCharPerLine>{{ .NumberOfCharPerLine }}</NumberOfCharPerLine>{{end}}
							</Message>{{end}}{{ if or .FirstStop .LastStop .LineRef}}
							<LineSection>{{ if .FirstStop }}
								<FirstStop>{{.FirstStop}}</FirstStop>{{end}}{{if .LastStop}}
							  <LastStop>{{.LastStop}}</LastStop>{{end}}{{if .LineRef}}
							  <LineRef>{{.LineRef}}</LineRef>{{end}}
							</LineSection>{{end}}
						</ns3:Content>
					</ns3:GeneralMessage>{{end}}
				</ns3:GeneralMessageDelivery>`

const generalMessageDeliveryTemplate = `<ns8:GetGeneralMessageResponse xmlns:ns3="http://www.siri.org.uk/siri"
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
		{{ .BuildGeneralMessageXML }}
	</Answer>
	<AnswerExtension/>
</ns8:GetGeneralMessageResponse>`

func (response *XMLGeneralMessageResponse) XMLGeneralMessage() []*XMLGeneralMessage {
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

func (visit *IDFGeneralMessageStructure) Messages() []*XMLMessage {
	if len(visit.messages) == 0 {
		nodes := visit.findNodes("Message")
		if nodes == nil {
			return visit.messages
		}
		for _, messageNode := range nodes {
			visit.messages = append(visit.messages, NewXMLMessage(messageNode))
		}
	}
	return visit.messages
}

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

func NewXMLGeneralMessage(node XMLNode) *XMLGeneralMessage {
	generalMessage := &XMLGeneralMessage{}
	generalMessage.node = node
	return generalMessage
}

func NewXMLMessage(node XMLNode) *XMLMessage {
	message := &XMLMessage{}
	message.node = node
	return message
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

func (visit *XMLGeneralMessage) Format() string {
	if visit.format == "" {
		node := visit.node.NativeNode()
		if node != nil {
			visit.format = node.Attr("formatRef")
		}
	}
	return visit.format
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

func (visit *IDFGeneralMessageStructure) GroupOfLinesRef() string {
	if visit.groupOfLinesRef == "" {
		visit.groupOfLinesRef = visit.findStringChildContent("GroupOfLinesRef")
	}
	return visit.groupOfLinesRef
}

func (visit *IDFGeneralMessageStructure) RouteRef() string {
	if visit.routeRef == "" {
		visit.routeRef = visit.findStringChildContent("RouteRef")
	}
	return visit.routeRef
}

func (visit *IDFGeneralMessageStructure) DestinationRef() string {
	if visit.destinationRef == "" {
		visit.destinationRef = visit.findStringChildContent("DestinationRef")
	}
	return visit.destinationRef
}

func (visit *IDFGeneralMessageStructure) JourneyPatternRef() string {
	if visit.journeyPatternRef == "" {
		visit.journeyPatternRef = visit.findStringChildContent("JourneyPatternRef")
	}
	return visit.journeyPatternRef
}

func (visit *IDFGeneralMessageStructure) StopPointRef() string {
	if visit.stopPointRef == "" {
		visit.stopPointRef = visit.findStringChildContent("StopPointRef")
	}
	return visit.stopPointRef
}

func (visit *IDFGeneralMessageStructure) LineRef() string {
	if visit.lineRef == "" {
		visit.lineRef = visit.findStringChildContent("LineRef")
	}
	return visit.lineRef
}

func (visit *IDFGeneralMessageStructure) createNewLineSection() IDFLineSectionStructure {
	visit.lineSection = IDFLineSectionStructure{}
	fmt.Println(visit.findNode("LineSection"))
	visit.lineSection.node = NewXMLNode(visit.findNode("LineSection"))
	return visit.lineSection
}

func (visit *IDFGeneralMessageStructure) LineSection() IDFLineSectionStructure {
	emptyStruct := IDFLineSectionStructure{}
	if visit.lineSection != emptyStruct {
		return visit.lineSection
	}
	return visit.createNewLineSection()
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
	var generalMessage = template.Must(template.New("generalMessage").Parse(generalMessageDeliveryTemplate))
	if err := generalMessage.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIGeneralMessageDelivery) BuildGeneralMessageXML() (string, error) {
	var buffer bytes.Buffer
	var generalMessageDelivery = template.Must(template.New("generalMessageDelivery").Parse(generalMessageTemplate))
	if err := generalMessageDelivery.Execute(&buffer, delivery); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
