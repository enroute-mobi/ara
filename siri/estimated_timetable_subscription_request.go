package siri

import "time"

type XMLEstimatedTimetableSubscriptionRequestEntry struct {
	XMLEstimatedTimetableRequest

	subscriberRef          string
	subscriptionRef        string
	initialTerminationTime time.Time
}

func NewXMLEstimatedTimetableSubscriptionRequestEntry(node XMLNode) *XMLEstimatedTimetableSubscriptionRequestEntry {
	xmlEstimatedTimetableSubscriptionRequest := &XMLEstimatedTimetableSubscriptionRequestEntry{}
	xmlEstimatedTimetableSubscriptionRequest.node = node
	return xmlEstimatedTimetableSubscriptionRequest
}

func (request *XMLEstimatedTimetableSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLEstimatedTimetableSubscriptionRequestEntry) SubscriptionRef() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent("SubscriptionRef")
	}
	return request.subscriptionRef
}

func (request *XMLEstimatedTimetableSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent("InitialTerminationTime")
	}
	return request.initialTerminationTime
}
