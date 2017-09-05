package siri

import (
	"time"

	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLEstimatedTimetableSubscriptionRequest struct {
	XMLEstimatedTimetableRequest

	subscriberRef          string
	subscriptionRef        string
	initialTerminationTime time.Time
}

func NewXMLEstimatedTimetableSubscriptionRequest(node xml.Node) *XMLEstimatedTimetableSubscriptionRequest {
	xmlEstimatedTimetableSubscriptionRequest := &XMLEstimatedTimetableSubscriptionRequest{}
	xmlEstimatedTimetableSubscriptionRequest.node = NewXMLNode(node)
	return xmlEstimatedTimetableSubscriptionRequest
}

func NewXMLEstimatedTimetableSubscriptionRequestFromContent(content []byte) (*XMLEstimatedTimetableSubscriptionRequest, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	request := NewXMLEstimatedTimetableSubscriptionRequest(doc.Root().XmlNode)
	return request, nil
}

func (request *XMLEstimatedTimetableSubscriptionRequest) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLEstimatedTimetableSubscriptionRequest) SubscriptionRef() string {
	if request.subscriptionRef == "" {
		request.subscriptionRef = request.findStringChildContent("SubscriptionRef")
	}
	return request.subscriptionRef
}

func (request *XMLEstimatedTimetableSubscriptionRequest) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent("InitialTerminationTime")
	}
	return request.initialTerminationTime
}
