package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

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
		request.subscriberRef = request.findStringChildContent(siri_attributes.SubscriberRef)
	}
	return request.subscriberRef
}

func (request *XMLEstimatedTimetableSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
	}
	return request.subscriptionRef
}

func (request *XMLEstimatedTimetableSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent(siri_attributes.InitialTerminationTime)
	}
	return request.initialTerminationTime
}
