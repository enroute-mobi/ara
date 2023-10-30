package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIGetGeneralMessageRequest struct {
	SIRIGeneralMessageRequest

	RequestorRef string
}

type SIRIGeneralMessageRequest struct {
	XsdInWsdl bool

	MessageIdentifier string

	RequestTimestamp time.Time

	InfoChannelRef []string

	LineRef           []string
	StopPointRef      []string
	JourneyPatternRef []string
	DestinationRef    []string
	RouteRef          []string
}

func NewSIRIGeneralMessageRequest(
	messageIdentifier,
	requestorRef string,
	requestTimestamp time.Time) *SIRIGetGeneralMessageRequest {
	request := &SIRIGetGeneralMessageRequest{
		RequestorRef: requestorRef,
	}
	request.MessageIdentifier = messageIdentifier
	request.RequestTimestamp = requestTimestamp
	return request
}

func (request *SIRIGetGeneralMessageRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("get_general_message_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (request *SIRIGeneralMessageRequest) BuildGeneralMessageRequestXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "general_message_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
