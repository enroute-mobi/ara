package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIEstimatedTimeTableResponse struct {
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

	Attributes map[string]string
	References map[string]string

	EstimatedCalls []*SIRIEstimatedCall
}

type SIRIEstimatedCall struct {
	ArrivalStatus      string
	DepartureStatus    string
	StopPointRef       string
	StopPointName      string
	DestinationDisplay string

	VehicleAtStop bool

	Order int

	AimedArrivalTime    time.Time
	ExpectedArrivalTime time.Time

	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
}

func (response *SIRIEstimatedTimeTableResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "estimated_timetable_response.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIEstimatedTimetableDelivery) BuildEstimatedTimetableDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "estimated_timetable_delivery.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (frame *SIRIEstimatedJourneyVersionFrame) BuildEstimatedJourneyVersionFrameXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "estimated_journey_version_frame.template", frame); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
