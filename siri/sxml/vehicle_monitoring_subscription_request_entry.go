package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLVehicleMonitoringSubscriptionRequestEntry struct {
	XMLVehicleMonitoringRequest

	subscriberRef          *string
	subscriptionRef        *string
	initialTerminationTime *time.Time
}

func NewXMLVehicleMonitoringSubscriptionRequestEntry(node XMLNode) *XMLVehicleMonitoringSubscriptionRequestEntry {
	xmlVehicleMonitoringSubscriptionRequest := &XMLVehicleMonitoringSubscriptionRequestEntry{}
	xmlVehicleMonitoringSubscriptionRequest.node = node
	return xmlVehicleMonitoringSubscriptionRequest
}

func (request *XMLVehicleMonitoringSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriberRef)
		request.subscriberRef = &s
	}
	return *request.subscriberRef
}

func (request *XMLVehicleMonitoringSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
		request.subscriptionRef = &s
	}
	return *request.subscriptionRef
}

func (request *XMLVehicleMonitoringSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime == nil {
		t := request.findTimeChildContent(siri_attributes.InitialTerminationTime)
		request.initialTerminationTime = &t
	}
	return *request.initialTerminationTime
}
