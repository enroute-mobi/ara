package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRINotifySituationExchange struct {
	Address                   string
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string
	SubscriberRef             string
	SubscriptionIdentifier    string

	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	SituationExchangeDeliveries []*SIRISituationExchangeDelivery
}

func (notify *SIRINotifySituationExchange) ErrorString() string {
	return fmt.Sprintf("%v: %v", notify.errorType(), notify.ErrorText)
}

func (notify *SIRINotifySituationExchange) errorType() string {
	if notify.ErrorType == "OtherError" {
		return fmt.Sprintf("%v %v", notify.ErrorType, notify.ErrorNumber)
	}
	return notify.ErrorType
}

func (notify *SIRINotifySituationExchange) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("situation_exchange_notify%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
