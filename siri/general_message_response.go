package siri

import (
	"bytes"
	rxml "encoding/xml"
	"fmt"
	"text/template"
	"time"

	"github.com/af83/edwig/model"
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
	infoMessageVersion    string
	infoChannelRef        string
	content               interface{}
}

type IDFGeneralMessageStructure struct {
	XMLStructure

	messages          []model.Message
	lineRef           string
	stopPointRef      string
	journeyPatternRef string
	destinationRef    string
	routeRef          string
	groupOfLinesRef   string
	lineSection       IDFLineSectionStructure
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

	Messages []*model.Message

	FirstStop string
	LastStop  string
	LineRef   string
}

const generalMessageTemplate = `<ns3:GeneralMessageDelivery version="2.0:FR-IDF-2.4">
				  <ns3:ResponseTimestamp>2017-03-29T16:48:00.039+02:00</ns3:ResponseTimestamp>
					<ns3:Status>{{.Status}}</ns3:Status>{{range .GeneralMessages}}
					<ns3:GeneralMessage formatRef="FRANCE">
						<ns3:RecordedAtTime>{{ .RecordedAtTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:RecordedAtTime>
						<ns3:InfoMessageIdentifier>{{.InfoMessageIdentifier}}</ns3:InfoMessageIdentifier>
						<ns3:InfoMessageVersion>{{.InfoMessageVersion}}</ns3:InfoMessageVersion>
						<ns3:InfoChannelRef>{{.InfoChannelRef}}</ns3:InfoChannelRef>
						<ns3:ValidUntilTime>{{ .ValidUntilTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ValidUntilTime>
						<ns3:Content xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
						xsi:type="ns9:IDFLineSectionStructure">{{range .Messages}}
							<Message>
								<MessageType>{{ .Type }}</MessageType>
								<MessageText>{{ .Content }}</MessageText>
								<NumberOfLines>{{ .NumberOfLines }}</NumberOfLines>
								<NumberOfCharPerLine>{{ .NumberOfCharPerLine }}</NumberOfCharPerLine>
							</Message>{{end}}
							<LineSection>
								<FirstStop>{{.FirstStop}}</FirstStop>
							  <LastStop>{{.LastStop}}</LastStop>
							  <LineRef>{{.LineRef}}</LineRef>
							</LineSection>
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
			response.xmlGeneralMessages = append(response.xmlGeneralMessages, NewIDFGeneralMessageStructure(generalMessage))
		}
	}
	return response.xmlGeneralMessages
}

func (visit *IDFGeneralMessageStructure) Messages() []model.Message {
	if len(visit.messages) == 0 {
		nodes := visit.findNodes("Message")
		if nodes == nil {
			return visit.messages
		}
		unmashallMessage := model.Message{}
		for _, message := range nodes {
			rxml.Unmarshal([]byte(message.NativeNode().String()), unmashallMessage)
			visit.messages = append(visit.messages, unmashallMessage)
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

func NewIDFGeneralMessageStructure(node XMLNode) *XMLGeneralMessage {
	generalMessage := &XMLGeneralMessage{}
	generalMessage.node = node
	return generalMessage
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

func (visit *XMLGeneralMessage) InfoMessageVersion() string {
	if visit.infoMessageVersion == "" {
		visit.infoMessageVersion = visit.findStringChildContent("InfoMessageVersion")
	}
	return visit.infoMessageVersion
}

func (visit *XMLGeneralMessage) InfoChannelRef() string {
	if visit.infoChannelRef == "" {
		visit.infoChannelRef = visit.findStringChildContent("InfoChannelRef")
	}
	return visit.infoChannelRef
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
