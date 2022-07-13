package siri

import (
	"bytes"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIDatedTimetableVersionFrame struct {
	LineRef        string
	DirectionType  string
	RecordedAtTime time.Time

	Attributes map[string]string

	DatedVehicleJourneys []*SIRIDatedVehicleJourney
}

type SIRIDatedVehicleJourney struct {
	DataFrameRef           string
	DatedVehicleJourneyRef string
	PublishedLineName      string

	Attributes map[string]string
	References map[string]string

	DatedCalls []*SIRIDatedCall
}

type SIRIDatedCall struct {
	StopPointRef       string
	StopPointName      string
	DestinationDisplay string
	VisitNumber        string

	AimedArrivalTime   time.Time
	AimedDepartureTime time.Time
}

func (frame *SIRIDatedTimetableVersionFrame) BuildDatedTimetableVersionFrameXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "dated_timetable_version_frame.template", frame); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (frame *SIRIDatedTimetableVersionFrame) BuildDatedTimetableVersionFrameXMLRaw() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "dated_timetable_version_frame_raw.template", frame); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}

	return strings.TrimSpace(buffer.String()), nil
}
