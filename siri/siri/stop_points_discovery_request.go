package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIStopPointsDiscoveryRequest struct {
	MessageIdentifier string
	RequestorRef      string

	RequestTimestamp time.Time
}

func NewSIRIStopPointsDiscoveryRequest(messageIdentifier, requestorRef string, requestTimestamp time.Time) *SIRIStopPointsDiscoveryRequest {
	return &SIRIStopPointsDiscoveryRequest{
		MessageIdentifier: messageIdentifier,
		RequestorRef:      requestorRef,
		RequestTimestamp:  requestTimestamp,
	}
}

func (request *SIRIStopPointsDiscoveryRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "" && envelopeType[0] != "soap" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("stop_points_discovery_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
