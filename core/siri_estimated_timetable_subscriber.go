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

type SIRIEstimatedTimetableSubscriber interface {
	state.Stopable
	state.Startable
}

type ETTSubscriber struct {
	clock.ClockConsumer

	connector *SIRIEstimatedTimetableSubscriptionCollector
}

type EstimatedTimetableSubscriber struct {
	ETTSubscriber

	stop chan struct{}
}

func NewSIRIEstimatedTimetableSubscriber(connector *SIRIEstimatedTimetableSubscriptionCollector) SIRIEstimatedTimetableSubscriber {
	subscriber := &EstimatedTimetableSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *EstimatedTimetableSubscriber) Start() {
	logger.Log.Debugf("Start EstimatedTimetableSubscriber")

	subscriber.stop = make(chan struct{})
	go subscriber.run()
}

func (subscriber *EstimatedTimetableSubscriber) run() {
	c := subscriber.Clock().After(5 * time.Second)

	for {
		select {
		case <-subscriber.stop:
			return
		case <-c:
			logger.Log.Debugf("SIRIEstimatedTimetableSubscriber visit")

			subscriber.prepareSIRIEstimatedTimetableSubscriptionRequest()

			c = subscriber.Clock().After(5 * time.Second)
		}
	}
}

func (subscriber *EstimatedTimetableSubscriber) Stop() {
	if subscriber.stop != nil {
		close(subscriber.stop)
	}
}

func (subscriber *ETTSubscriber) prepareSIRIEstimatedTimetableSubscriptionRequest() {
	subscriptions := subscriber.connector.partner.Subscriptions().FindSubscriptionsByKind(EstimatedTimetableCollect)
	if len(subscriptions) == 0 {
		logger.Log.Debugf("EstimatedTimetableSubscriber visit without EstimatedTimetableCollect subscriptions")
		return
	}

	linesToLog := []string{}
	requestMessageRefToSub := make(map[string]string)
	subToRequestMessageRef := make(map[string]string)

	linesToRequest := make(map[string][]string)
	for _, subscription := range subscriptions {
		for _, resource := range subscription.ResourcesByObjectIDCopy() {
			if resource.SubscribedAt().IsZero() && resource.RetryCount <= 10 {
				mid := subscriber.connector.Partner().NewMessageIdentifier()
				if len(linesToRequest[string(subscription.id)]) == 0 {
					requestMessageRefToSub[mid] = string(subscription.id)
					subToRequestMessageRef[string(subscription.id)] = mid
				}
				linesToRequest[string(subscription.id)] = append(linesToRequest[string(subscription.id)], resource.Reference.ObjectId.Value())
			}
		}
	}

	if len(linesToRequest) == 0 {
		return
	}

	message := subscriber.newBQEvent()
	defer audit.CurrentBigQuery(string(subscriber.connector.Partner().Referential().Slug())).WriteEvent(message)

	siriEstimatedTimetableSubscriptionRequest := &siri.SIRIEstimatedTimetableSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.Partner().NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.Partner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
		SortForTest:       subscriber.connector.Partner().SortForTests(),
	}

	var subIds []string
	for subscription, requestedLines := range linesToRequest {
		entry := &siri.SIRIEstimatedTimetableSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
			SubscriptionIdentifier: subscription,
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = subToRequestMessageRef[subscription]
		entry.RequestTimestamp = subscriber.Clock().Now()
		entry.Lines = requestedLines

		linesToLog = append(linesToLog, entry.Lines...)
		subIds = append(subIds, subscription)
		siriEstimatedTimetableSubscriptionRequest.Entries = append(siriEstimatedTimetableSubscriptionRequest.Entries, entry)
	}

	message.RequestIdentifier = siriEstimatedTimetableSubscriptionRequest.MessageIdentifier
	message.RequestRawMessage, _ = siriEstimatedTimetableSubscriptionRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
	message.RequestSize = int64(len(message.RequestRawMessage))
	message.Lines = linesToLog
	message.SubscriptionIdentifiers = subIds

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.Partner().SIRIClient().EstimatedTimetableSubscription(siriEstimatedTimetableSubscriptionRequest)
	message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		e := fmt.Sprintf("Error during EstimatedTimetableSubscriptionRequest: %v", err)

		subscriber.incrementRetryCountFromMap(linesToRequest)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
	message.ResponseIdentifier = response.ResponseMessageIdentifier()

	for _, responseStatus := range response.ResponseStatus() {
		subId, ok := requestMessageRefToSub[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		_, ok = linesToRequest[subId]
		if !ok { // Should never happen
			logger.Log.Debugf("Error in ETT Subscription Collector, no lines to request for subscription %v", subId)
			continue
		}

		subscription, ok := subscriber.connector.partner.Subscriptions().Find(SubscriptionId(subId))
		if !ok { // Should never happen
			logger.Log.Debugf("Response for unknown subscription %v", subId)
			continue
		}
		for _, line := range linesToRequest[subId] {
			resource := subscription.Resource(model.NewObjectID(subscriber.connector.remoteObjectidKind, line))
			if resource == nil { // Should never happen
				logger.Log.Debugf("Response for unknown subscription resource %v", line)
				continue
			}

			if !responseStatus.Status() {
				logger.Log.Debugf("Subscription status false for line %v: %v %v ", line, responseStatus.ErrorType(), responseStatus.ErrorText())
				resource.RetryCount++
				message.Status = "Error"
				continue
			}
			resource.Subscribed(subscriber.Clock().Now())
			resource.RetryCount = 0
		}
		delete(linesToRequest, subId) // See #4691
	}
	// Should not happen but see #4691
	if len(linesToRequest) == 0 {
		return
	}
	subscriber.incrementRetryCountFromMap(linesToRequest)
}

func (subscriber *ETTSubscriber) incrementRetryCountFromMap(linesToRequest map[string][]string) {
	for subId, requestedLines := range linesToRequest {
		subscription, ok := subscriber.connector.partner.Subscriptions().Find(SubscriptionId(subId))
		if !ok { // Should never happen
			continue
		}
		for _, l := range requestedLines {
			resource := subscription.Resource(model.NewObjectID(subscriber.connector.remoteObjectidKind, l))
			if resource == nil { // Should never happen
				continue
			}
			resource.RetryCount++

		}
	}
}

func (subscriber *ETTSubscriber) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "EstimatedTimetableSubscriptionRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(subscriber.connector.partner.Slug()),
		Status:    "OK",
	}
}
