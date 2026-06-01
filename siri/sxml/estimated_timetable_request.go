package sxml

import (
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLGetEstimatedTimetable struct {
	XMLEstimatedTimetableRequest

	requestorRef *string
}

type XMLEstimatedTimetableRequest struct {
	LightRequestXMLStructure

	previewInterval time.Duration
	startTime       *time.Time

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
		nodes := request.findNodes(siri_attributes.LineRef)
		for _, node := range nodes {
			request.lines = append(request.lines, strings.TrimSpace(node.NativeNode().Content()))
		}
	}
	return request.lines
}

func (request *XMLGetEstimatedTimetable) RequestorRef() string {
	if request.requestorRef == nil {
		s := request.findStringChildContent(siri_attributes.RequestorRef)
		request.requestorRef = &s
	}
	return *request.requestorRef
}

func (request *XMLEstimatedTimetableRequest) PreviewInterval() time.Duration {
	if request.previewInterval == 0 {
		request.previewInterval = request.findDurationChildContent(siri_attributes.PreviewInterval)
	}
	return request.previewInterval
}

func (request *XMLEstimatedTimetableRequest) StartTime() time.Time {
	if request.startTime == nil {
		t := request.findTimeChildContent(siri_attributes.StartTime)
		request.startTime = &t
	}
	return *request.startTime
}
