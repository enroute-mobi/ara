package sxml

import (
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLFacilityMonitoringSubscriptionRequestEntry struct {
	XMLFacilityMonitoringRequest

	subscriberRef          string
	subscriptionRef        string
	initialTerminationTime time.Time

	facilities []string
}

func NewXMLFacilityMonitoringSubscriptionRequestEntry(node XMLNode) *XMLFacilityMonitoringSubscriptionRequestEntry {
	xmlFacilityMonitoringSubscriptionRequest := &XMLFacilityMonitoringSubscriptionRequestEntry{}
	xmlFacilityMonitoringSubscriptionRequest.node = node
	return xmlFacilityMonitoringSubscriptionRequest
}

func (request *XMLFacilityMonitoringSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent(siri_attributes.SubscriberRef)
	}
	return request.subscriberRef
}

func (request *XMLFacilityMonitoringSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
	}
	return request.subscriptionRef
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
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent(siri_attributes.InitialTerminationTime)
	}
	return request.initialTerminationTime
}
