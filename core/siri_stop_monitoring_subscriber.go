package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringSubscriber interface {
	model.Stopable
	model.Startable
}

type SMSubscriber struct {
	model.ClockConsumer

	connector *SIRIStopMonitoringSubscriptionCollector
}

type StopMonitoringSubscriber struct {
	SMSubscriber

	stop chan struct{}
}

type FakeStopMonitoringSubscriber struct {
	SMSubscriber
}

func NewFakeStopMonitoringSubscriber(connector *SIRIStopMonitoringSubscriptionCollector) SIRIStopMonitoringSubscriber {
	subscriber := &FakeStopMonitoringSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *FakeStopMonitoringSubscriber) Start() {
	subscriber.prepareSIRIStopMonitoringSubscriptionRequest()
}

func (subscriber *FakeStopMonitoringSubscriber) Stop() {}

func NewSIRIStopMonitoringSubscriber(connector *SIRIStopMonitoringSubscriptionCollector) SIRIStopMonitoringSubscriber {
	subscriber := &StopMonitoringSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *StopMonitoringSubscriber) Start() {
	logger.Log.Debugf("Start StopMonitoringSubscriber")

	subscriber.stop = make(chan struct{})
	go subscriber.run()
}

func (subscriber *StopMonitoringSubscriber) run() {
	c := subscriber.Clock().After(5 * time.Second)

	for {
		select {
		case <-subscriber.stop:
			return
		case <-c:
			logger.Log.Debugf("SIRIStopMonitoringSubscriber visit")

			subscriber.prepareSIRIStopMonitoringSubscriptionRequest()

			c = subscriber.Clock().After(5 * time.Second)
		}
	}
}

func (subscriber *StopMonitoringSubscriber) Stop() {
	if subscriber.stop != nil {
		close(subscriber.stop)
	}
}

func (subscriber *SMSubscriber) prepareSIRIStopMonitoringSubscriptionRequest() {
	subscription, _ := subscriber.connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")
	monitoringRefList := []string{}

	stopAreasToRequest := make(map[string]*model.ObjectID)
	for _, resource := range subscription.ResourcesByObjectID() {
		if resource.SubscribedAt.IsZero() && resource.RetryCount <= 10 {
			messageIdentifier := subscriber.connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier()
			stopAreasToRequest[messageIdentifier] = resource.Reference.ObjectId
		}
	}

	if len(stopAreasToRequest) == 0 {
		return
	}

	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	siriStopMonitoringSubscriptionRequest := &siri.SIRIStopMonitoringSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	for messageIdentifier, stopAreaObjectid := range stopAreasToRequest {
		entry := &siri.SIRIStopMonitoringSubscriptionRequestEntry{
			MessageIdentifier:      messageIdentifier,
			RequestTimestamp:       subscriber.Clock().Now(),
			MonitoringRef:          stopAreaObjectid.Value(),
			SubscriberRef:          subscriber.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier: string(subscription.Id()),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}

		monitoringRefList = append(monitoringRefList, stopAreaObjectid.Value())
		siriStopMonitoringSubscriptionRequest.Entries = append(siriStopMonitoringSubscriptionRequest.Entries, entry)
	}

	logStashEvent["MonitoringRef"] = strings.Join(monitoringRefList, ", ")
	logSIRIStopMonitoringSubscriptionRequest(logStashEvent, siriStopMonitoringSubscriptionRequest, monitoringRefList)

	response, err := subscriber.connector.SIRIPartner().SOAPClient().StopMonitoringSubscription(siriStopMonitoringSubscriptionRequest)
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		for _, stopAreaObjectid := range stopAreasToRequest {
			resource := subscription.Resource(*stopAreaObjectid)
			resource.RetryCount++
		}
		return
	}

	logStashEvent["response"] = response.RawXML()

	for _, responseStatus := range response.ResponseStatus() {
		stopAreaObjectid, ok := stopAreasToRequest[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		resource := subscription.Resource(*stopAreaObjectid)
		if resource == nil { // Should never happen
			logger.Log.Debugf("Response for unknown subscription resource %v", stopAreaObjectid.String())
			continue
		}
		if !responseStatus.Status() {
			logger.Log.Debugf("Subscription status false for stopArea %v: %v %v ", stopAreaObjectid.Value(), responseStatus.ErrorType(), responseStatus.ErrorText())
			resource.RetryCount++
			continue
		}
		resource.SubscribedAt = subscriber.Clock().Now()
		resource.RetryCount = 0
	}
}

func logSIRIStopMonitoringSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIStopMonitoringSubscriptionRequest, monitorinRefList []string) {
	logStashEvent["Connector"] = "SIRIStopMonitoringSubscriptionCollector"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	xml, err := request.BuildXML()
	if err != nil {
		logStashEvent["requestXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["requestXML"] = xml
}
