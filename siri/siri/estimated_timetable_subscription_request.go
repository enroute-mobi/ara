package siri

import (
	"bytes"
	"fmt"
	"sort"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIEstimatedTimetableSubscriptionRequest struct {
	ConsumerAddress   string
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time

	Entries []*SIRIEstimatedTimetableSubscriptionRequestEntry

	SortForTest bool
}

type SIRIEstimatedTimetableSubscriptionRequestEntry struct {
	SIRIEstimatedTimetableRequest

	SubscriberRef          string
	SubscriptionIdentifier string

	InitialTerminationTime time.Time
}

func (request *SIRIEstimatedTimetableSubscriptionRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("estimated_timetable_subscription_request%s.template", envType)

	if request.SortForTest {
		for _, entry := range request.Entries {
			sort.Strings(entry.Lines)
		}
	}

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
