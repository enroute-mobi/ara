package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/state"
)

type SIRISituationExchangeSubscriber interface {
	state.Stopable
	state.Startable
}

type SXSubscriber struct {
	clock.ClockConsumer

	connector *SIRISituationExchangeSubscriptionCollector
}

type SituationExchangeSubscriber struct {
	SXSubscriber

	stop chan struct{}
}

func NewSIRISituationExchangeSubscriber(connector *SIRISituationExchangeSubscriptionCollector) SIRISituationExchangeSubscriber {
	subscriber := &SituationExchangeSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *SituationExchangeSubscriber) Start() {
	logger.Log.Debugf("Start SituationExchangeSubscriber")

	subscriber.stop = make(chan struct{})
	go subscriber.run()
}

func (subscriber *SituationExchangeSubscriber) run() {
	c := subscriber.Clock().After(5 * time.Second)

	for {
		select {
		case <-subscriber.stop:
			return
		case <-c:
			logger.Log.Debugf("SIRISituationExchangeSubscriber visit")

			subscriber.prepareSIRISituationExchangeSubscriptionRequest()

			c = subscriber.Clock().After(5 * time.Second)
		}
	}
}

func (subscriber *SituationExchangeSubscriber) Stop() {
	if subscriber.stop != nil {
		close(subscriber.stop)
	}
}

func (subscriber *SXSubscriber) prepareSIRISituationExchangeSubscriptionRequest() {
	subscriptions := subscriber.connector.partner.Subscriptions().FindSubscriptionsByKind(SituationExchangeCollect)
	if len(subscriptions) == 0 {
		logger.Log.Debugf("SituationExchangeSubscriber visit without SituationExchangeCollect subscriptions")
		return
	}

	lineRefList := []string{}
	stopPointRefList := []string{}

	resourcesToRequest := make(map[string]*resourceToRequest)
	for _, subscription := range subscriptions {
		for _, resource := range subscription.ResourcesByCodeCopy() {
			if resource.SubscribedAt().IsZero() && resource.RetryCount <= 10 {
				messageIdentifier := subscriber.connector.Partner().NewMessageIdentifier()
				logger.Log.Debugf("send request for subscription with id : %v", subscription.id)
				resourcesToRequest[messageIdentifier] = &resourceToRequest{
					subId: subscription.id,
					code:  *(resource.Reference.Code),
					kind:  resource.Reference.Type,
				}
			}
		}
	}

	if len(resourcesToRequest) == 0 {
		return
	}

	message := subscriber.newBQEvent()
	defer audit.CurrentBigQuery(string(subscriber.connector.Partner().Referential().Slug())).WriteEvent(message)

	sxRequest := &siri.SIRISituationExchangeSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.Partner().NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.Partner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	var subIDs []string
	for messageIdentifier, requestedResource := range resourcesToRequest {
		entry := &siri.SIRISituationExchangeSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
			SubscriptionIdentifier: string(requestedResource.subId),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = messageIdentifier
		entry.RequestTimestamp = subscriber.Clock().Now()
		subIDs = append(subIDs, entry.SubscriptionIdentifier)
		switch requestedResource.kind {
		case "Line":
			entry.LineRef = []string{requestedResource.code.Value()}
			lineRefList = append(lineRefList, requestedResource.code.Value())
		case "StopArea":
			entry.StopPointRef = []string{requestedResource.code.Value()}
			stopPointRefList = append(stopPointRefList, requestedResource.code.Value())
		}

		sxRequest.Entries = append(sxRequest.Entries, entry)
	}

	message.RequestIdentifier = sxRequest.MessageIdentifier
	message.RequestRawMessage, _ = sxRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
	message.RequestSize = int64(len(message.RequestRawMessage))
	message.StopAreas = stopPointRefList
	message.Lines = lineRefList
	message.SubscriptionIdentifiers = subIDs

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.Partner().SIRIClient().SituationExchangeSubscription(sxRequest)
	message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		e := fmt.Sprintf("Error during SituationExchangeSubscriptionRequest: %v", err)
		subscriber.incrementRetryCountFromMap(resourcesToRequest)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))

	for _, responseStatus := range response.ResponseStatus() {
		requestedResource, ok := resourcesToRequest[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		delete(resourcesToRequest, responseStatus.RequestMessageRef()) // See #4691

		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedResource.subId)
		if !ok { // Should never happen
			logger.Log.Debugf("Response for unknown subscription %v", requestedResource.subId)
			continue
		}
		resource := subscription.Resource(requestedResource.code)
		if resource == nil { // Should never happen
			logger.Log.Debugf("Response for unknown subscription resource %v", requestedResource.code.String())
			continue
		}

		if !responseStatus.Status() {
			logger.Log.Debugf("Subscription status false for %v %v: %v %v ",
				requestedResource.kind,
				requestedResource.code.Value(),
				responseStatus.ErrorType(),
				responseStatus.ErrorText(),
			)
			resource.RetryCount++
			message.Status = "Error"
			continue
		}
		resource.Subscribed(subscriber.Clock().Now())
		resource.RetryCount = 0
	}
	// Should not happen but see #4691
	if len(resourcesToRequest) == 0 {
		return
	}
	subscriber.incrementRetryCountFromMap(resourcesToRequest)
}

func (subscriber *SXSubscriber) incrementRetryCountFromMap(resourcesToRequest map[string]*resourceToRequest) {
	for _, requestedResource := range resourcesToRequest {
		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedResource.subId)
		if !ok { // Should never happen
			continue
		}
		resource := subscription.Resource(requestedResource.code)
		if resource == nil { // Should never happen
			continue
		}
		resource.RetryCount++
	}
}

func (subscriber *SXSubscriber) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.SITUATION_EXCHANGE_SUBSCRIPTION_REQUEST,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(subscriber.connector.partner.Slug()),
		Status:    "OK",
	}
}
