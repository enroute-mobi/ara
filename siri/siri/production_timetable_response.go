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
	Order              int
	UseVisitNumber     bool

	AimedArrivalTime   time.Time
	AimedDepartureTime time.Time
}

type SIRIDatedCalls []*SIRIDatedCall

func (a SIRIDatedCalls) Len() int      { return len(a) }
func (a SIRIDatedCalls) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type SortByStopPointOrder struct{ SIRIDatedCalls }

func (s SortByStopPointOrder) Less(i, j int) bool {
	return s.SIRIDatedCalls[i].Order < s.SIRIDatedCalls[j].Order
}

type SIRIDatedTimetableVersionFrames []*SIRIDatedTimetableVersionFrame

func (a SIRIDatedTimetableVersionFrames) Len() int      { return len(a) }
func (a SIRIDatedTimetableVersionFrames) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type SortByDirectionType struct {
	SIRIDatedTimetableVersionFrames
}

func (s SortByDirectionType) Less(i, j int) bool {
	return s.SIRIDatedTimetableVersionFrames[i].DirectionType < s.SIRIDatedTimetableVersionFrames[j].DirectionType
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
