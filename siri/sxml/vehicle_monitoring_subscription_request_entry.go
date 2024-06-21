package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLVehicleMonitoringSubscriptionRequestEntry struct {
	XMLVehicleMonitoringRequest

	subscriberRef          string
	subscriptionRef        string
	initialTerminationTime time.Time
}

func NewXMLVehicleMonitoringSubscriptionRequestEntry(node XMLNode) *XMLVehicleMonitoringSubscriptionRequestEntry {
	xmlVehicleMonitoringSubscriptionRequest := &XMLVehicleMonitoringSubscriptionRequestEntry{}
	xmlVehicleMonitoringSubscriptionRequest.node = node
	return xmlVehicleMonitoringSubscriptionRequest
}

func (request *XMLVehicleMonitoringSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent(siri_attributes.SubscriberRef)
	}
	return request.subscriberRef
}

func (request *XMLVehicleMonitoringSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
	}
	return request.subscriptionRef
}

func (request *XMLVehicleMonitoringSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent(siri_attributes.InitialTerminationTime)
	}
	return request.initialTerminationTime
}
