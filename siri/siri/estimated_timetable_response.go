package siri

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type SIRIEstimatedTimetableResponse struct {
	SIRIEstimatedTimetableDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
}

type SIRIEstimatedTimetableDelivery struct {
	RequestMessageRef string

	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	EstimatedJourneyVersionFrames []*SIRIEstimatedJourneyVersionFrame

	MonitoringRefs     map[string]struct{}
	VehicleJourneyRefs map[string]struct{}
}

type SIRIEstimatedJourneyVersionFrame struct {
	RecordedAtTime time.Time

	EstimatedVehicleJourneys []*SIRIEstimatedVehicleJourney
}

type SIRIEstimatedVehicleJourney struct {
	LineRef                string
	DatedVehicleJourneyRef string
	DataFrameRef           string
	DirectionType          string
	PublishedLineName      string

	IsCompleteStopSequence bool

	Attributes map[string]string
	References map[string]string

	EstimatedCalls []*SIRIEstimatedCall
	RecordedCalls  []*SIRIRecordedCall
}

type SIRIEstimatedCalls []*SIRIEstimatedCall

func (a SIRIEstimatedCalls) Len() int      { return len(a) }
func (a SIRIEstimatedCalls) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type SortEstimatedCallsByStopVisitOrder struct{ SIRIEstimatedCalls }

func (s SortEstimatedCallsByStopVisitOrder) Less(i, j int) bool {
	return s.SIRIEstimatedCalls[i].Order < s.SIRIEstimatedCalls[j].Order
}

type SIRIRecordedCalls []*SIRIRecordedCall

func (a SIRIRecordedCalls) Len() int      { return len(a) }
func (a SIRIRecordedCalls) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type SortRecordedCallsByStopVisitOrder struct{ SIRIRecordedCalls }

func (s SortRecordedCallsByStopVisitOrder) Less(i, j int) bool {
	return s.SIRIRecordedCalls[i].Order < s.SIRIRecordedCalls[j].Order
}

type SIRIEstimatedCall struct {
	ArrivalStatus      string
	DepartureStatus    string
	StopPointRef       string
	StopPointName      string
	DestinationDisplay string

	Occupancy     string
	VehicleAtStop bool

	Order          int
	UseVisitNumber bool

	AimedArrivalTime    time.Time
	ExpectedArrivalTime time.Time

	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
}

type SIRIRecordedCall struct {
	ArrivalStatus      string
	DepartureStatus    string
	StopPointRef       string
	StopPointName      string
	DestinationDisplay string

	Order          int
	UseVisitNumber bool

	AimedArrivalTime    time.Time
	ExpectedArrivalTime time.Time

	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
}

func (response *SIRIEstimatedTimetableResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("estimated_timetable_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIEstimatedTimetableDelivery) ErrorString() string {
	return fmt.Sprintf("%v: %v", delivery.errorType(), delivery.ErrorText)
}

func (delivery *SIRIEstimatedTimetableDelivery) errorType() string {
	if delivery.ErrorType == siri_attributes.OtherError {
		return fmt.Sprintf("%v %v", delivery.ErrorType, delivery.ErrorNumber)
	}
	return delivery.ErrorType
}

func (delivery *SIRIEstimatedTimetableDelivery) BuildEstimatedTimetableDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "estimated_timetable_delivery.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (frame *SIRIEstimatedJourneyVersionFrame) BuildEstimatedJourneyVersionFrameXML() (string, error) {
	var buffer bytes.Buffer

	for _, evj := range frame.EstimatedVehicleJourneys {
		sort.Sort(SortEstimatedCallsByStopVisitOrder{evj.EstimatedCalls})
		sort.Sort(SortRecordedCallsByStopVisitOrder{evj.RecordedCalls})
	}

	if err := templates.ExecuteTemplate(&buffer, "estimated_journey_version_frame.template", frame); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (delivery *SIRIEstimatedTimetableDelivery) BuildEstimatedTimetableDeliveryXMLRaw() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "estimated_timetable_delivery_raw.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (frame *SIRIEstimatedJourneyVersionFrame) BuildEstimatedJourneyVersionFrameXMLRaw() (string, error) {
	var buffer bytes.Buffer

	for _, evj := range frame.EstimatedVehicleJourneys {
		sort.Sort(SortEstimatedCallsByStopVisitOrder{evj.EstimatedCalls})
		sort.Sort(SortRecordedCallsByStopVisitOrder{evj.RecordedCalls})
	}

	if err := templates.ExecuteTemplate(&buffer, "estimated_journey_version_frame_raw.template", frame); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}

	return strings.TrimSpace(buffer.String()), nil
}
