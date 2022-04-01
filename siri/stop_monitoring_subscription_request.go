package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type XMLStopMonitoringSubscriptionRequestEntry struct {
	LightXMLStopMonitoringRequest

	subscriberRef          string
	subscriptionIdentifier string
	initialTerminationTime time.Time
}

type SIRIStopMonitoringSubscriptionRequest struct {
	ConsumerAddress   string
	MessageIdentifier string
	RequestorRef      string
	RequestTimestamp  time.Time

	Entries []*SIRIStopMonitoringSubscriptionRequestEntry
}

type SIRIStopMonitoringSubscriptionRequestEntry struct {
	SIRIStopMonitoringRequest

	SubscriberRef          string
	SubscriptionIdentifier string

	InitialTerminationTime time.Time
}

func NewXMLStopMonitoringSubscriptionRequestEntry(node XMLNode) *XMLStopMonitoringSubscriptionRequestEntry {
	xmlStopMonitoringSubscriptionRequestEntry := &XMLStopMonitoringSubscriptionRequestEntry{}
	xmlStopMonitoringSubscriptionRequestEntry.node = node
	return xmlStopMonitoringSubscriptionRequestEntry
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) SubscriberRef() string {
	if request.subscriberRef == "" {
		request.subscriberRef = request.findStringChildContent("SubscriberRef")
	}
	return request.subscriberRef
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) SubscriptionIdentifier() string {
	if request.subscriptionIdentifier == "" {
		request.subscriptionIdentifier = request.findStringChildContent("SubscriptionIdentifier")
	}
	return request.subscriptionIdentifier
}

func (request *XMLStopMonitoringSubscriptionRequestEntry) InitialTerminationTime() time.Time {
	if request.initialTerminationTime.IsZero() {
		request.initialTerminationTime = request.findTimeChildContent("InitialTerminationTime")
	}
	return request.initialTerminationTime
}

func (request *SIRIStopMonitoringSubscriptionRequest) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("stop_monitoring_subscription_request%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, request); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
