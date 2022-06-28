package siri

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIDeleteSubscriptionResponse struct {
	ResponderRef      string
	RequestMessageRef string
	ResponseTimestamp time.Time

	ResponseStatus []*SIRITerminationResponseStatus
}

type SIRITerminationResponseStatus struct {
	SubscriberRef     string
	SubscriptionRef   string
	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber string
	ErrorText   string
}

func (notify *SIRIDeleteSubscriptionResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("delete_subscription_response%s.template", envType)

	// order SubscriptionRef lexicographically
	sort.Slice(notify.ResponseStatus, func(i, j int) bool {
		return strings.ToLower(notify.ResponseStatus[i].SubscriptionRef) <
			strings.ToLower(notify.ResponseStatus[j].SubscriptionRef)
	})

	if err := templates.ExecuteTemplate(&buffer, templateName, notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
