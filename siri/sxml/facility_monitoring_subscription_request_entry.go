package sxml

import (
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLFacilityMonitoringSubscriptionRequestEntry struct {
	XMLFacilityMonitoringRequest

	subscriberRef          *string
	subscriptionRef        *string
	initialTerminationTime *time.Time

	facilities []string
}

func NewXMLFacilityMonitoringSubscriptionRequestEntry(node XMLNode) *XMLFacilityMonitoringSubscriptionRequestEntry {
	xmlFacilityMonitoringSubscriptionRequest := &XMLFacilityMonitoringSubscriptionRequestEntry{}
	xmlFacilityMonitoringSubscriptionRequest.node = node
	return xmlFacilityMonitoringSubscriptionRequest
}

func (request *XMLFacilityMonitoringSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriberRef)
		request.subscriberRef = &s
	}
	return *request.subscriberRef
}

func (request *XMLFacilityMonitoringSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionRef == nil {
		s := request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
		request.subscriptionRef = &s
	}
	return *request.subscriptionRef
}

func (request *XMLFacilityMonitoringSubscriptionRequestEntry) FacilityRefs() []string {
	if len(request.facilities) == 0 {
		nodes := request.findNodes(siri_attributes.FacilityRef)
		for _, node := range nodes {
			request.facilities = append(request.facilities, strings.TrimSpace(node.NativeNode().Content()))
		}
	}
	return request.facilities
}

func (request *XMLFacilityMonitoringSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime == nil {
		t := request.findTimeChildContent(siri_attributes.InitialTerminationTime)
		request.initialTerminationTime = &t
	}
	return *request.initialTerminationTime
}
