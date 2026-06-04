package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLEstimatedTimetableSubscriptionRequestEntry struct {
	XMLEstimatedTimetableRequest

	subscriberRef          *string
	subscriptionRef        *string
	initialTerminationTime *time.Time
}

func NewXMLEstimatedTimetableSubscriptionRequestEntry(node XMLNode) *XMLEstimatedTimetableSubscriptionRequestEntry {
	xmlEstimatedTimetableSubscriptionRequest := &XMLEstimatedTimetableSubscriptionRequestEntry{}
	xmlEstimatedTimetableSubscriptionRequest.node = node
	return xmlEstimatedTimetableSubscriptionRequest
}

func (request *XMLEstimatedTimetableSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriberRef)
		request.subscriberRef = &s
	}
	return *request.subscriberRef
}

func (request *XMLEstimatedTimetableSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
		request.subscriptionRef = &s
	}
	return *request.subscriptionRef
}

func (request *XMLEstimatedTimetableSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime == nil {
		t := request.findTimeChildContent(siri_attributes.InitialTerminationTime)
		request.initialTerminationTime = &t
	}
	return *request.initialTerminationTime
}
