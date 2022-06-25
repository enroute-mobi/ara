package sxml

import (
	"time"
)

type XMLGeneralMessageSubscriptionRequestEntry struct {
	XMLGeneralMessageRequest

	subscriberRef          string
	subscriptionIdentifier string
	initialTerminationTime time.Time
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
