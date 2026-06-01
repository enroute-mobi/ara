package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLProductionTimetableSubscriptionRequestEntry struct {
	XMLProductionTimetableRequest

	subscriberRef          *string
	subscriptionRef        *string
	initialTerminationTime *time.Time
}

func NewXMLProductionTimetableSubscriptionRequestEntry(node XMLNode) *XMLProductionTimetableSubscriptionRequestEntry {
	xmlProductionTimetableSubscriptionRequest := &XMLProductionTimetableSubscriptionRequestEntry{}
	xmlProductionTimetableSubscriptionRequest.node = node
	return xmlProductionTimetableSubscriptionRequest
}

func (request *XMLProductionTimetableSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriberRef)
		request.subscriberRef = &s
	}
	return *request.subscriberRef
}

func (request *XMLProductionTimetableSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
		request.subscriptionRef = &s
	}
	return *request.subscriptionRef
}

func (request *XMLProductionTimetableSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime == nil {
		t := request.findTimeChildContent(siri_attributes.InitialTerminationTime)
		request.initialTerminationTime = &t
	}
	return *request.initialTerminationTime
}
