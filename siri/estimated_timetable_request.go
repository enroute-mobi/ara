package siri

import (
	"strings"
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetEstimatedTimetable struct {
	XMLEstimatedTimetableRequest

	requestorRef string
}

type XMLEstimatedTimetableRequest struct {
	XMLStructure

	messageIdentifier string

	requestTimestamp time.Time

	previewInterval time.Duration
	startTime       time.Time

	lines []string
}

func NewXMLGetEstimatedTimetable(node xml.Node) *XMLGetEstimatedTimetable {
	xmlGetEstimatedTimetable := &XMLGetEstimatedTimetable{}
	xmlGetEstimatedTimetable.node = NewXMLNode(node)
	return xmlGetEstimatedTimetable
}

func NewXMLGetEstimatedTimetableFromContent(content []byte) (*XMLGetEstimatedTimetable, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLGetEstimatedTimetable(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLEstimatedTimetableRequest) Lines() []string {
	if len(request.lines) == 0 {
		nodes := request.findNodes("LineRef")
		for _, node := range nodes {
			request.lines = append(request.lines, strings.TrimSpace(node.NativeNode().Content()))
		}
	}
	return request.lines
}

func (request *XMLGetEstimatedTimetable) RequestorRef() string {
	if request.requestorRef == "" {
		request.requestorRef = request.findStringChildContent("RequestorRef")
	}
	return request.requestorRef
}

func (request *XMLEstimatedTimetableRequest) MessageIdentifier() string {
	if request.messageIdentifier == "" {
		request.messageIdentifier = request.findStringChildContent("MessageIdentifier")
	}
	return request.messageIdentifier
}

func (request *XMLEstimatedTimetableRequest) RequestTimestamp() time.Time {
	if request.requestTimestamp.IsZero() {
		request.requestTimestamp = request.findTimeChildContent("RequestTimestamp")
	}
	return request.requestTimestamp
}

func (request *XMLEstimatedTimetableRequest) PreviewInterval() time.Duration {
	if request.previewInterval == 0 {
		request.previewInterval = request.findDurationChildContent("PreviewInterval")
	}
	return request.previewInterval
}

func (request *XMLEstimatedTimetableRequest) StartTime() time.Time {
	if request.startTime.IsZero() {
		request.startTime = request.findTimeChildContent("StartTime")
	}
	return request.startTime
}
