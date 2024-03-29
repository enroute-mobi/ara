package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/state"
)

type SIRIStopMonitoringSubscriber interface {
	state.Stopable
	state.Startable
}

type SMSubscriber struct {
	clock.ClockConsumer

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
	saId  model.Code
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
	subscriptions := subscriber.connector.partner.Subscriptions().FindSubscriptionsByKind(StopMonitoringCollect)
	if len(subscriptions) == 0 {
		logger.Log.Debugf("StopMonitoringSubscriber visit without StopMonitoringCollect subscriptions")
		return
	}

	monitoringRefList := []string{}

	stopAreasToRequest := make(map[string]*saToRequest)
	for _, subscription := range subscriptions {
		for _, resource := range subscription.ResourcesByCodeCopy() {
			if resource.SubscribedAt().IsZero() && resource.RetryCount <= 10 {
				messageIdentifier := subscriber.connector.Partner().NewMessageIdentifier()
				stopAreasToRequest[messageIdentifier] = &saToRequest{
					subId: subscription.id,
					saId:  *(resource.Reference.Code),
				}
			}
		}
	}

	if len(stopAreasToRequest) == 0 {
		return
	}

	message := subscriber.newBQEvent()
	defer audit.CurrentBigQuery(string(subscriber.connector.Partner().Referential().Slug())).WriteEvent(message)

	siriStopMonitoringSubscriptionRequest := &siri.SIRIStopMonitoringSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.Partner().NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.Partner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	var subIds []string
	for messageIdentifier, requestedSa := range stopAreasToRequest {
		entry := &siri.SIRIStopMonitoringSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
			SubscriptionIdentifier: string(requestedSa.subId),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = messageIdentifier
		entry.RequestTimestamp = subscriber.Clock().Now()
		entry.MonitoringRef = requestedSa.saId.Value()

		monitoringRefList = append(monitoringRefList, entry.MonitoringRef)
		subIds = append(subIds, entry.SubscriptionIdentifier)
		siriStopMonitoringSubscriptionRequest.Entries = append(siriStopMonitoringSubscriptionRequest.Entries, entry)
	}

	message.RequestIdentifier = siriStopMonitoringSubscriptionRequest.MessageIdentifier
	message.RequestRawMessage, _ = siriStopMonitoringSubscriptionRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
	message.RequestSize = int64(len(message.RequestRawMessage))
	message.StopAreas = monitoringRefList
	message.SubscriptionIdentifiers = subIds

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.Partner().SIRIClient().StopMonitoringSubscription(siriStopMonitoringSubscriptionRequest)
	message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		e := fmt.Sprintf("Error during StopMonitoringSubscriptionRequest: %v", err)

		subscriber.incrementRetryCountFromMap(stopAreasToRequest)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
	message.ResponseIdentifier = response.ResponseMessageIdentifier()

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
			message.Status = "Error"
			continue
		}
		resource.Subscribed(subscriber.Clock().Now())
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

func (subscriber *SMSubscriber) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "StopMonitoringSubscriptionRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(subscriber.connector.partner.Slug()),
		Status:    "OK",
	}
}
