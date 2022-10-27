package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIGetEstimatedTimetableRequest struct {
	SIRIEstimatedTimetableRequest

	RequestorRef string
}

type SIRIEstimatedTimetableRequest struct {
	MessageIdentifier string

	Lines []string

	RequestTimestamp time.Time
}

func NewSIRIGetEstimatedTimetableRequest(
	messageIdentifier string,
	lines []string,
	requestorRef string,
	requestTimestamp time.Time) *SIRIGetEstimatedTimetableRequest {
	request := &SIRIGetEstimatedTimetableRequest{
		RequestorRef: requestorRef,
	}
	request.MessageIdentifier = messageIdentifier
	request.Lines = lines
	request.RequestTimestamp = requestTimestamp
	return request
}

func (request *SIRIGetEstimatedTimetableRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("get_estimated_timetable_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIEstimatedTimetableRequest) BuildEstimatedTimetableRequestXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "estimated_timetable_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
