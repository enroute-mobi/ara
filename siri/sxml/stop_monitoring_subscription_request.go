package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLStopMonitoringSubscriptionRequestEntry struct {
	LightXMLStopMonitoringRequest

	subscriberRef          string
	subscriptionIdentifier string
	initialTerminationTime time.Time
}

func NewXMLStopMonitoringSubscriptionRequestEntry(node XMLNode) *XMLStopMonitoringSubscriptionRequestEntry {
	xmlStopMonitoringSubscriptionRequestEntry := &XMLStopMonitoringSubscriptionRequestEntry{}
	xmlStopMonitoringSubscriptionRequestEntry.node = node
	return xmlStopMonitoringSubscriptionRequestEntry
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent(siri_attributes.SubscriberRef)
	}
	return request.subscriberRef
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionIdentifier == "" {
		request.subscriptionIdentifier = request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
	}
	return request.subscriptionIdentifier
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent(siri_attributes.InitialTerminationTime)
	}
	return request.initialTerminationTime
}
