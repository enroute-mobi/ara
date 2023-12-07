package siri

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIGetSituationExchangeRequest struct {
	SIRISituationExchangeRequest

	RequestorRef string
}

type SIRISituationExchangeRequest struct {
	MessageIdentifier string

	RequestTimestamp time.Time

	InfoChannelRef []string

	LineRef      []string
	StopPointRef []string
}

func NewSIRISituationExchangeRequest(
	messageIdentifier,
	requestorRef string,
	requestTimestamp time.Time) *SIRIGetSituationExchangeRequest {
	request := &SIRIGetSituationExchangeRequest{
		RequestorRef: requestorRef,
	}
	request.MessageIdentifier = messageIdentifier
	request.RequestTimestamp = requestTimestamp
	return request
}

func (request *SIRIGetSituationExchangeRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("get_situation_exchange_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (request *SIRISituationExchangeRequest) BuildSituationExchangeRequestXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "situation_exchange_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}
