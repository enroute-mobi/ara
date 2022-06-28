package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRISubscriptionResponse struct {
	Address           string
	ResponderRef      string
	RequestMessageRef string

	ResponseTimestamp  time.Time
	ServiceStartedTime time.Time

	ResponseStatus []SIRIResponseStatus
}

type SIRIResponseStatus struct {
	RequestMessageRef string
	SubscriberRef     string
	SubscriptionRef   string

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	ResponseTimestamp time.Time
	ValidUntil        time.Time
}

func (response *SIRISubscriptionResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("subscription_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
