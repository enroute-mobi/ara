package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLGeneralMessageSubscriptionRequestEntry struct {
	XMLGeneralMessageRequest

	subscriberRef          *string
	subscriptionIdentifier *string
	initialTerminationTime *time.Time
}

func NewXMLGeneralMessageSubscriptionRequestEntry(node XMLNode) *XMLGeneralMessageSubscriptionRequestEntry {
	xmlGeneralMessageSubscriptionRequestEntry := &XMLGeneralMessageSubscriptionRequestEntry{}
	xmlGeneralMessageSubscriptionRequestEntry.node = node
	return xmlGeneralMessageSubscriptionRequestEntry
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriberRef)
		request.subscriberRef = &s
	}
	return *request.subscriberRef
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionIdentifier == nil {
		s := request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
		request.subscriptionIdentifier = &s
	}
	return *request.subscriptionIdentifier
}

func (request *XMLGeneralMessageSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime == nil {
		t := request.findTimeChildContent(siri_attributes.InitialTerminationTime)
		request.initialTerminationTime = &t
	}
	return *request.initialTerminationTime
}
