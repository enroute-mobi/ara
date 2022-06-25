package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRICheckStatusRequest struct {
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time
}

func NewSIRICheckStatusRequest(
	RequestorRef string,
	RequestTimestamp time.Time,
	MessageIdentifier string) *SIRICheckStatusRequest {
	return &SIRICheckStatusRequest{
		RequestorRef:      RequestorRef,
		RequestTimestamp:  RequestTimestamp,
		MessageIdentifier: MessageIdentifier,
	}
}

// TODO : Handle errors
func (request *SIRICheckStatusRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("check_status_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)

		return "", err
	}
	return buffer.String(), nil
}
