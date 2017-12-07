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

type saToRequest struct {
	subId SubscriptionId
	saId  model.ObjectID
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
	subscriptions := subscriber.connector.partner.Subscriptions().FindSubscriptionsByKind("StopMonitoringCollect")
	if len(subscriptions) == 0 {
		logger.Log.Debugf("StopMonitoringSubscriber visit without StopMonitoringCollect subscriptions")
		return
	}

	// MonitoringRef for Logstash
	monitoringRefList := []string{}

	stopAreasToRequest := make(map[string]*saToRequest)
	for _, subscription := range subscriptions {
		for _, resource := range subscription.ResourcesByObjectID() {
			if resource.SubscribedAt.IsZero() && resource.RetryCount <= 10 {
				messageIdentifier := subscriber.connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier()
				stopAreasToRequest[messageIdentifier] = &saToRequest{
					subId: subscription.id,
					saId:  *(resource.Reference.ObjectId),
				}
			}
		}
	}

	if len(stopAreasToRequest) == 0 {
		return
	}

	logStashEvent := subscriber.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	siriStopMonitoringSubscriptionRequest := &siri.SIRIStopMonitoringSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	for messageIdentifier, requestedSa := range stopAreasToRequest {
		entry := &siri.SIRIStopMonitoringSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier: string(requestedSa.subId),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = messageIdentifier
		entry.RequestTimestamp = subscriber.Clock().Now()
		entry.MonitoringRef = requestedSa.saId.Value()

		monitoringRefList = append(monitoringRefList, entry.MonitoringRef)
		siriStopMonitoringSubscriptionRequest.Entries = append(siriStopMonitoringSubscriptionRequest.Entries, entry)
	}

	logStashEvent["monitoringRef"] = strings.Join(monitoringRefList, ", ")
	logSIRIStopMonitoringSubscriptionRequest(logStashEvent, siriStopMonitoringSubscriptionRequest)

	response, err := subscriber.connector.SIRIPartner().SOAPClient().StopMonitoringSubscription(siriStopMonitoringSubscriptionRequest)
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["response"] = fmt.Sprintf("Error during StopMonitoringSubscriptionRequest: %v", err)
		subscriber.incrementRetryCountFromMap(stopAreasToRequest)
		return
	}

	logStashEvent["response"] = response.RawXML()

	for _, responseStatus := range response.ResponseStatus() {
		requestedSa, ok := stopAreasToRequest[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		delete(stopAreasToRequest, responseStatus.RequestMessageRef()) // See #4691

		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedSa.subId)
		if !ok { // Should never happen
			logger.Log.Debugf("Response for unknown subscription %v", requestedSa.subId)
			continue
		}
		resource := subscription.Resource(requestedSa.saId)
		if resource == nil { // Should never happen
			logger.Log.Debugf("Response for unknown subscription resource %v", requestedSa.saId.String())
			continue
		}

		if !responseStatus.Status() {
			logger.Log.Debugf("Subscription status false for stopArea %v: %v %v ", requestedSa.saId.Value(), responseStatus.ErrorType(), responseStatus.ErrorText())
			resource.RetryCount++
			continue
		}
		resource.SubscribedAt = subscriber.Clock().Now()
		resource.RetryCount = 0
	}
	// Should not happen but see #4691
	if len(stopAreasToRequest) == 0 {
		return
	}
	subscriber.incrementRetryCountFromMap(stopAreasToRequest)
}

func (subscriber *SMSubscriber) incrementRetryCountFromMap(stopAreasToRequest map[string]*saToRequest) {
	for _, requestedSa := range stopAreasToRequest {
		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedSa.subId)
		if !ok { // Should never happen
			continue
		}
		resource := subscription.Resource(requestedSa.saId)
		if resource == nil { // Should never happen
			continue
		}
		resource.RetryCount++
	}
}

func (subscriber *SMSubscriber) newLogStashEvent() audit.LogStashEvent {
	event := subscriber.connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringSubscriptionCollector"
	return event
}

func logSIRIStopMonitoringSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIStopMonitoringSubscriptionRequest) {
	logStashEvent["type"] = "StopMonitoringSubscriptionRequest"
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
