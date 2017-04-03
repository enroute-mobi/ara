package siri

import (
	"fmt"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGeneralMessageResponse struct {
	ResponseXMLStructure

	recordedAtTime        time.Time
	validUntilTime        time.Time
	itemIdentifier        string
	infoMessageIdentifier string
	infoMessageVersion    string
	infoChannelRef        string
	contents              interface{}
}

type IDFGeneralMessageStructure struct {
	XMLStructure

	numberOfLines       int
	numberOfCharPerLine int
	messageType         string
	messageText         string
	lineRef             string
	stopPointRef        string
	journeyPatternRef   string
	destinationRef      string
	routeRef            string
	groupOfLinesRef     string
	lineSection         IDFLineSectionStructure
}

type IDFLineSectionStructure struct {
	XMLStructure

	firstStop string
	lastStop  string
	lineRef   string
}

func (visit *XMLGeneralMessageResponse) Contents() interface{} {
	if visit.contents != nil {
		return visit.contents
	}
	nodes := visit.findNodes("Content")
	if nodes == nil {
		return nil
	}
	contents := []*IDFGeneralMessageStructure{}
	for _, content := range nodes {
		contents = append(contents, NewIDFGeneralMessageStructure(content))
	}
	visit.contents = contents
	return visit.contents
}

func NewIDFGeneralMessageStructure(node XMLNode) *IDFGeneralMessageStructure {
	content := &IDFGeneralMessageStructure{}
	content.node = node
	return content
}

func (visit *XMLGeneralMessageResponse) RecordedAtTime() time.Time {
	if visit.recordedAtTime.IsZero() {
		visit.recordedAtTime = visit.findTimeChildContent("RecordedAtTime")
	}
	return visit.recordedAtTime
}

func (visit *XMLGeneralMessageResponse) ValidUntilTime() time.Time {
	if visit.validUntilTime.IsZero() {
		visit.validUntilTime = visit.findTimeChildContent("ValidUntilTime")
	}
	return visit.validUntilTime
}

func (visit *XMLGeneralMessageResponse) ItemIdentifier() string {
	if visit.itemIdentifier == "" {
		visit.itemIdentifier = visit.findStringChildContent("ItemIdentifier")
	}
	return visit.itemIdentifier
}

func (visit *XMLGeneralMessageResponse) InfoMessageIdentifier() string {
	if visit.infoMessageIdentifier == "" {
		visit.infoMessageIdentifier = visit.findStringChildContent("InfoMessageIdentifier")
	}
	return visit.infoMessageIdentifier
}

func (visit *XMLGeneralMessageResponse) InfoChannelRef() string {
	if visit.infoChannelRef == "" {
		visit.infoChannelRef = visit.findStringChildContent("InfoChannelRef")
	}
	return visit.infoChannelRef
}

func (visit *XMLGeneralMessageResponse) InfoMessageVersion() string {
	if visit.infoMessageVersion == "" {
		visit.infoMessageVersion = visit.findStringChildContent("InfoMessageVersion")
	}
	return visit.infoMessageVersion
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

func (visit *IDFGeneralMessageStructure) MessageText() string {
	if visit.messageText == "" {
		visit.messageText = visit.findStringChildContent("MessageText")
	}
	return visit.messageText
}

func (visit *IDFGeneralMessageStructure) MessageType() string {
	if visit.messageType == "" {
		visit.messageType = visit.findStringChildContent("MessageType")
	}
	return visit.messageType
}

func (visit *IDFGeneralMessageStructure) NumberOfLines() int {
	if visit.numberOfLines == 0 {
		visit.numberOfLines = visit.findIntChildContent("NumberOfLines")
	}
	return visit.numberOfLines
}

func (visit *IDFGeneralMessageStructure) NumberOfCharPerLine() int {
	if visit.numberOfCharPerLine == 0 {
		visit.numberOfCharPerLine = visit.findIntChildContent("NumberOfCharPerLine")
	}
	return visit.numberOfCharPerLine
}

func (visit *IDFGeneralMessageStructure) createNewLineSection() IDFLineSectionStructure {
	visit.lineSection = IDFLineSectionStructure{}
	visit.lineSection.node = NewXMLNode(visit.findNode("LineSection"))
	fmt.Println(visit.lineSection.node)
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
