package siri

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
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

	Attributes map[string]string
	References map[string]string

	EstimatedCalls []*SIRIEstimatedCall
	RecordedCalls  []*SIRIRecordedCall
}

type SIRIEstimatedCall struct {
	ArrivalStatus      string
	DepartureStatus    string
	StopPointRef       string
	StopPointName      string
	DestinationDisplay string

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
	if delivery.ErrorType == "OtherError" {
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
	if err := templates.ExecuteTemplate(&buffer, "estimated_journey_version_frame_raw.template", frame); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}

	return strings.TrimSpace(buffer.String()), nil
}
