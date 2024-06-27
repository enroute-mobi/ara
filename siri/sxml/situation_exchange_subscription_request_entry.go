package sxml

import (
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type XMLSituationExchangeSubscriptionRequestEntry struct {
	XMLSituationExchangeRequest

	subscriberRef          string
	subscriptionIdentifier string
	initialTerminationTime time.Time
}

func NewXMLSituationExchangeSubscriptionRequestEntry(node XMLNode) *XMLSituationExchangeSubscriptionRequestEntry {
	xmlSituationExchangeSubscriptionRequestEntry := &XMLSituationExchangeSubscriptionRequestEntry{}
	xmlSituationExchangeSubscriptionRequestEntry.node = node
	return xmlSituationExchangeSubscriptionRequestEntry
}

func (request *XMLSituationExchangeSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent(siri_attributes.SubscriberRef)
	}
	return request.subscriberRef
}

func (request *XMLSituationExchangeSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionIdentifier == "" {
		request.subscriptionIdentifier = request.findStringChildContent(siri_attributes.SubscriptionIdentifier)
	}
	return request.subscriptionIdentifier
}

func (request *XMLSituationExchangeSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent(siri_attributes.InitialTerminationTime)
	}
	return request.initialTerminationTime
}

func (request *XMLSituationExchangeSubscriptionRequestEntry) LineRefs() []string {
	if len(request.lineRefs) == 0 {
		nodes := request.findNodes(siri_attributes.LineRef)
		for _, lineRef := range nodes {
			request.lineRefs = append(request.lineRefs, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return request.lineRefs
}

func (request *XMLSituationExchangeSubscriptionRequestEntry) StopPointRefs() []string {
	if len(request.stopPointRefs) == 0 {
		nodes := request.findNodes(siri_attributes.StopPointRef)
		for _, lineRef := range nodes {
			request.stopPointRefs = append(request.stopPointRefs, strings.TrimSpace(lineRef.NativeNode().Content()))
		}
	}
	return request.stopPointRefs
}
