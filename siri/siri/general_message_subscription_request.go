package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIGeneralMessageSubscriptionRequest struct {
	ConsumerAddress   string
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time

	Entries []*SIRIGeneralMessageSubscriptionRequestEntry
}

type SIRIGeneralMessageSubscriptionRequestEntry struct {
	SIRIGeneralMessageRequest

	SubscriberRef          string
	SubscriptionIdentifier string

	InitialTerminationTime time.Time
}

func (request *SIRIGeneralMessageSubscriptionRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("general_message_subscription_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
