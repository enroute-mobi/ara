package sxml

import (
	"time"
)

type XMLProductionTimetableSubscriptionRequestEntry struct {
	XMLProductionTimetableRequest

	subscriberRef          string
	subscriptionRef        string
	initialTerminationTime time.Time
}

func NewXMLProductionTimetableSubscriptionRequestEntry(node XMLNode) *XMLProductionTimetableSubscriptionRequestEntry {
	xmlProductionTimetableSubscriptionRequest := &XMLProductionTimetableSubscriptionRequestEntry{}
	xmlProductionTimetableSubscriptionRequest.node = node
	return xmlProductionTimetableSubscriptionRequest
}

func (request *XMLProductionTimetableSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLProductionTimetableSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent("SubscriptionIdentifier")
	}
	return request.subscriptionRef
}

func (request *XMLProductionTimetableSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent("InitialTerminationTime")
	}
	return request.initialTerminationTime
}
