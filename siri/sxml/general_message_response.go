package sxml

import (
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGeneralMessageResponse struct {
	ResponseXMLStructureWithStatus

	xmlGeneralMessages []*XMLGeneralMessage
}

type XMLGeneralMessageCancellation struct {
	XMLStructure

	infoMessageIdentifier string
	recordedAtTime        time.Time
}

type XMLGeneralMessage struct {
	XMLStructure

	itemIdentifier        string
	infoMessageIdentifier string
	infoChannelRef        string
	formatRef             string

	infoMessageVersion Int

	recordedAtTime time.Time
	validUntilTime time.Time

	content interface{}
}

type IDFGeneralMessageStructure struct {
	XMLStructure

	lineRef        []string
	stopPointRef   []string
	destinationRef []string
	routeRef       []string

	lineSections []*IDFLineSectionStructure
	messages     []*XMLMessage
}

type XMLMessage struct {
	XMLStructure

	messageTexts        map[string]string
	messageType         string
	numberOfLines       Int
	numberOfCharPerLine Int
}

type IDFLineSectionStructure struct {
	XMLStructure

	firstStop string
	lastStop  string
	lineRef   string
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

func (response *XMLGeneralMessageResponse) ErrorString() string {
	return fmt.Sprintf("%v: %v", response.errorType(), response.ErrorText())
}

func (response *XMLGeneralMessageResponse) errorType() string {
	if response.ErrorType() == siri_attributes.OtherError {
		return fmt.Sprintf("%v %v", response.ErrorType(), response.ErrorNumber())
	}
	return response.ErrorType()
}

func (response *XMLGeneralMessageResponse) XMLGeneralMessages() []*XMLGeneralMessage {
	if len(response.xmlGeneralMessages) == 0 {
		nodes := response.findNodes(siri_attributes.GeneralMessage)
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
		visit.infoMessageIdentifier = visit.findStringChildContent(siri_attributes.InfoMessageIdentifier)
	}
	return visit.infoMessageIdentifier
}

func (visit *XMLGeneralMessageCancellation) RecordedAtTime() time.Time {
	if visit.recordedAtTime.IsZero() {
		visit.recordedAtTime = visit.findTimeChildContent("RecordedAtTime")
	}
	return visit.recordedAtTime
}

func (visit *XMLGeneralMessage) RecordedAtTime() time.Time {
	if visit.recordedAtTime.IsZero() {
		visit.recordedAtTime = visit.findTimeChildContent(siri_attributes.RecordedAtTime)
	}
	return visit.recordedAtTime
}

func (visit *XMLGeneralMessage) ValidUntilTime() time.Time {
	if visit.validUntilTime.IsZero() {
		visit.validUntilTime = visit.findTimeChildContent(siri_attributes.ValidUntilTime)
	}
	return visit.validUntilTime
}

func (visit *XMLGeneralMessage) ItemIdentifier() string {
	if visit.itemIdentifier == "" {
		visit.itemIdentifier = visit.findStringChildContent(siri_attributes.ItemIdentifier)
	}
	return visit.itemIdentifier
}

func (visit *XMLGeneralMessage) InfoMessageIdentifier() string {
	if visit.infoMessageIdentifier == "" {
		visit.infoMessageIdentifier = visit.findStringChildContent(siri_attributes.InfoMessageIdentifier)
	}
	return visit.infoMessageIdentifier
}

func (visit *XMLGeneralMessage) InfoMessageVersion() int {
	if !visit.infoMessageVersion.Defined {
		visit.infoMessageVersion.SetValueWithDefault(visit.findIntChildContent(siri_attributes.InfoMessageVersion), 1)
	}
	return visit.infoMessageVersion.Value
}

func (visit *XMLGeneralMessage) InfoChannelRef() string {
	if visit.infoChannelRef == "" {
		visit.infoChannelRef = visit.findStringChildContent(siri_attributes.InfoChannelRef)
	}
	return visit.infoChannelRef
}

func (visit *XMLGeneralMessage) FormatRef() string {
	if visit.formatRef == "" {
		visit.formatRef = visit.node.NativeNode().Attr("formatRef")
	}
	return visit.formatRef
}

func (visit *XMLGeneralMessage) createNewContent() IDFGeneralMessageStructure {
	content := IDFGeneralMessageStructure{}
	content.node = NewXMLNode(visit.findNode(siri_attributes.Content))
	return content
}

func (visit *XMLGeneralMessage) Content() interface{} {
	if visit.content != nil {
		return visit.content
	}
	visit.content = visit.createNewContent()
	return visit.content
}

func (visit *IDFGeneralMessageStructure) RouteRef() []string {
	if len(visit.routeRef) == 0 {
		nodes := visit.findNodes(siri_attributes.RouteRef)
		for _, routeRef := range nodes {
			visit.routeRef = append(visit.routeRef, strings.TrimSpace(routeRef.NativeNode().Content()))
		}
	}
	return visit.routeRef
}

func (visit *IDFGeneralMessageStructure) DestinationRef() []string {
	if len(visit.destinationRef) == 0 {
		nodes := visit.findNodes(siri_attributes.DestinationRef)
		for _, destinationRef := range nodes {
			visit.destinationRef = append(visit.destinationRef, strings.TrimSpace(destinationRef.NativeNode().Content()))
		}
	}
	return visit.destinationRef
}

func (visit *IDFGeneralMessageStructure) StopPointRef() []string {
	if len(visit.stopPointRef) == 0 {
		nodes := visit.findNodes(siri_attributes.StopPointRef)
		for _, stopPointRef := range nodes {
			visit.stopPointRef = append(visit.stopPointRef, strings.TrimSpace(stopPointRef.NativeNode().Content()))
		}
	}
	return visit.stopPointRef
}

func (visit *IDFGeneralMessageStructure) LineRefs() []string {
	if len(visit.lineRef) == 0 {
		nodes := visit.findDirectChildrenNodes("LineRef")
		for _, lineRef := range nodes {
			visit.lineRef = append(visit.lineRef, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return visit.lineRef
}

func (visit *IDFGeneralMessageStructure) LineSections() []*IDFLineSectionStructure {
	if len(visit.lineSections) == 0 {
		nodes := visit.findNodes(siri_attributes.LineSection)
		for _, lineNode := range nodes {
			visit.lineSections = append(visit.lineSections, NewXMLLineSection(lineNode))
		}
	}
	return visit.lineSections
}

func (visit *IDFGeneralMessageStructure) Messages() []*XMLMessage {
	if len(visit.messages) == 0 {
		nodes := visit.findNodes(siri_attributes.Message)
		for _, messageNode := range nodes {
			message := NewXMLMessage(messageNode)
			// shortMessage should be inserted first
			if message.MessageType() == "shortMessage" {
				visit.messages = append([]*XMLMessage{message}, visit.messages...)
			} else {
				visit.messages = append(visit.messages, message)
			}
		}
	}
	return visit.messages
}

func (visit *IDFLineSectionStructure) FirstStop() string {
	if visit.firstStop == "" {
		visit.firstStop = visit.findStringChildContent(siri_attributes.FirstStop)
	}
	return visit.firstStop
}

func (visit *IDFLineSectionStructure) LastStop() string {
	if visit.lastStop == "" {
		visit.lastStop = visit.findStringChildContent(siri_attributes.LastStop)
	}
	return visit.lastStop
}

func (visit *IDFLineSectionStructure) LineRef() string {
	if visit.lineRef == "" {
		visit.lineRef = visit.findStringChildContent(siri_attributes.LineRef)
	}
	return visit.lineRef
}

func (message *XMLMessage) MessageTexts() map[string]string {
	if message.messageTexts == nil {
		translations := FindTranslations(message.findNodes(siri_attributes.MessageText))
		if translations != nil {
			message.messageTexts = translations
		}
	}
	return message.messageTexts
}

func (message *XMLMessage) MessageType() string {
	if message.messageType == "" {
		message.messageType = message.findStringChildContent(siri_attributes.MessageType)
	}
	return message.messageType
}

func (message *XMLMessage) NumberOfLines() int {
	if !message.numberOfLines.Defined {
		message.numberOfLines.SetValue(message.findIntChildContent(siri_attributes.NumberOfLines))
	}
	return message.numberOfLines.Value
}

func (message *XMLMessage) NumberOfCharPerLine() int {
	if !message.numberOfCharPerLine.Defined {
		message.numberOfCharPerLine.SetValue(message.findIntChildContent(siri_attributes.NumberOfCharPerLine))
	}
	return message.numberOfCharPerLine.Value
}
