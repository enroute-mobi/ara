package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/state"
	"fmt"
	"time"
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
	collectSubscriber := NewCollectSubcriber(subscriber.connector, EstimatedTimetableCollect)
	subscriptionRequests := collectSubscriber.GetSubscriptionRequest()

	if len(subscriptionRequests) == 0 {
		return
	}

	message := subscriber.newBQEvent()
	defer audit.CurrentBigQuery(string(subscriber.connector.Partner().Referential().Slug())).WriteEvent(message)

	siriEstimatedTimetableSubscriptionRequest := &siri.SIRIEstimatedTimetableSubscriptionRequest{
		ConsumerAddress:    subscriber.connector.Partner().Address(),
		MessageIdentifier:  subscriber.connector.Partner().NewMessageIdentifier(),
		RequestorRef:       subscriber.connector.Partner().RequestorRef(),
		RequestTimestamp:   subscriber.Clock().Now(),
		SortPayloadForTest: subscriber.connector.Partner().SortPaylodForTest(),
	}

	subIds := []string{}
	linesToLog := []string{}
	for subId, subscriptionRequest := range subscriptionRequests {
		entry := &siri.SIRIEstimatedTimetableSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
			SubscriptionIdentifier: string(subId),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = subscriptionRequest.requestMessageRef
		entry.RequestTimestamp = subscriber.Clock().Now()
		for _, m := range subscriptionRequest.modelsToRequest {
			switch m.kind {
			case "Line":
				entry.Lines = append(entry.Lines, m.code.Value())
			}
		}
		linesToLog = append(linesToLog, entry.Lines...)
		subIds = append(subIds, string(subId))
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

		collectSubscriber.IncrementRetryCountFromMap(subscriptionRequests)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
	message.ResponseIdentifier = response.ResponseMessageIdentifier()

	collectSubscriber.HandleResponse(subscriptionRequests, message, response)

	if len(subscriptionRequests) == 0 {
		return
	}

	collectSubscriber.IncrementRetryCountFromMap(subscriptionRequests)
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
