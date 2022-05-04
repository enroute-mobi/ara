package siri

import (
	"strings"
	"time"
)

type XMLProductionTimetableRequest struct {
	LightRequestXMLStructure

	previewInterval time.Duration
	startTime       time.Time

	lines []string
}

func (request *XMLProductionTimetableRequest) Lines() []string {
	if len(request.lines) == 0 {
		nodes := request.findNodes("LineRef")
		for _, node := range nodes {
			request.lines = append(request.lines, strings.TrimSpace(node.NativeNode().Content()))
		}
	}
	return request.lines
}

func (request *XMLProductionTimetableRequest) PreviewInterval() time.Duration {
	if request.previewInterval == 0 {
		request.previewInterval = request.findDurationChildContent("PreviewInterval")
	}
	return request.previewInterval
}

func (request *XMLProductionTimetableRequest) StartTime() time.Time {
	if request.startTime.IsZero() {
		request.startTime = request.findTimeChildContent("StartTime")
	}
	return request.startTime
}
