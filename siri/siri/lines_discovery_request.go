package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRILinesDiscoveryRequest struct {
	MessageIdentifier string
	RequestorRef      string

	RequestTimestamp time.Time
}

func NewSIRILinesDiscoveryRequest(messageIdentifier, requestorRef string, requestTimestamp time.Time) *SIRILinesDiscoveryRequest {
	return &SIRILinesDiscoveryRequest{
		MessageIdentifier: messageIdentifier,
		RequestorRef:      requestorRef,
		RequestTimestamp:  requestTimestamp,
	}
}

func (request *SIRILinesDiscoveryRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "" && envelopeType[0] != "soap" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("lines_discovery_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
