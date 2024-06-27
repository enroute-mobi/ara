package sxml

import (
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLProductionTimetableRequest struct {
	LightRequestXMLStructure

	previewInterval time.Duration
	startTime       time.Time

	lines []string
}

func (request *XMLProductionTimetableRequest) Lines() []string {
	if len(request.lines) == 0 {
		nodes := request.findNodes(siri_attributes.LineRef)
		for _, node := range nodes {
			request.lines = append(request.lines, strings.TrimSpace(node.NativeNode().Content()))
		}
	}
	return request.lines
}

func (request *XMLProductionTimetableRequest) PreviewInterval() time.Duration {
	if request.previewInterval == 0 {
		request.previewInterval = request.findDurationChildContent(siri_attributes.PreviewInterval)
	}
	return request.previewInterval
}

func (request *XMLProductionTimetableRequest) StartTime() time.Time {
	if request.startTime.IsZero() {
		request.startTime = request.findTimeChildContent(siri_attributes.StartTime)
	}
	return request.startTime
}
