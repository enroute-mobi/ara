package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type XMLGeneralMessageSubscriptionRequestEntry struct {
	XMLGeneralMessageRequest

	subscriberRef          string
	subscriptionIdentifier string
	initialTerminationTime time.Time
}

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

func NewXMLGeneralMessageSubscriptionRequestEntry(node XMLNode) *XMLGeneralMessageSubscriptionRequestEntry {
	xmlGeneralMessageSubscriptionRequestEntry := &XMLGeneralMessageSubscriptionRequestEntry{}
	xmlGeneralMessageSubscriptionRequestEntry.node = node
	return xmlGeneralMessageSubscriptionRequestEntry
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionIdentifier == "" {
		request.subscriptionIdentifier = request.findStringChildContent("SubscriptionIdentifier")
	}
	return request.subscriptionIdentifier
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent("InitialTerminationTime")
	}
	return request.initialTerminationTime
}

func (request *SIRIGeneralMessageSubscriptionRequest) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "general_message_subscription_request.template", request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
