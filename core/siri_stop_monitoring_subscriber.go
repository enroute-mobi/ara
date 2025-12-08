package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/state"
	"golang.org/x/exp/maps"
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
	collectSubscriber := NewCollectSubcriber(subscriber.connector, StopMonitoringCollect)
	subscriptionRequests := collectSubscriber.GetSubscriptionRequest()

	if len(subscriptionRequests) == 0 {
		return
	}

	subscriptionIds := maps.Keys(subscriptionRequests)
	batchSize := subscriber.connector.Partner().StopMonitoringMaxSubscriptionPerRequest()
	batches := make([][]SubscriptionId, 0, (len(subscriptionIds)+batchSize-1)/batchSize)

	for batchSize < len(subscriptionIds) {
		subscriptionIds, batches = subscriptionIds[batchSize:], append(batches, subscriptionIds[0:batchSize:batchSize])
	}

	batches = append(batches, subscriptionIds)

	for i := range batches {
		message := subscriber.newBQEvent()
		defer audit.CurrentBigQuery(string(subscriber.connector.Partner().Referential().Slug())).WriteEvent(message)

		siriStopMonitoringSubscriptionRequest := &siri.SIRIStopMonitoringSubscriptionRequest{
			ConsumerAddress:   subscriber.connector.Partner().Address(),
			MessageIdentifier: subscriber.connector.Partner().NewMessageIdentifier(),
			RequestorRef:      subscriber.connector.Partner().RequestorRef(),
			RequestTimestamp:  subscriber.Clock().Now(),
		}

		var subscriptionIdsToLog []string
		stopAreasToLog := []string{}
		for _, subscriptionId := range batches[i] {
			for _, m := range subscriptionRequests[subscriptionId].modelsToRequest {
				entry := &siri.SIRIStopMonitoringSubscriptionRequestEntry{
					SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
					SubscriptionIdentifier: string(subscriptionId),
					InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
				}
				entry.MessageIdentifier = subscriptionRequests[subscriptionId].requestMessageRef
				entry.RequestTimestamp = subscriber.Clock().Now()

				switch m.kind {
				case "StopArea":
					entry.MonitoringRef = m.code.Value()
					stopAreasToLog = append(stopAreasToLog, entry.MonitoringRef)
					siriStopMonitoringSubscriptionRequest.Entries = append(siriStopMonitoringSubscriptionRequest.Entries, entry)
				}
			}
			subscriptionIdsToLog = append(subscriptionIdsToLog, string(subscriptionId))
		}

		message.RequestIdentifier = siriStopMonitoringSubscriptionRequest.MessageIdentifier
		message.RequestRawMessage, _ = siriStopMonitoringSubscriptionRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
		message.RequestSize = int64(len(message.RequestRawMessage))
		message.StopAreas = stopAreasToLog
		message.SubscriptionIdentifiers = subscriptionIdsToLog

		startTime := subscriber.Clock().Now()
		response, err := subscriber.connector.Partner().SIRIClient().StopMonitoringSubscription(siriStopMonitoringSubscriptionRequest)
		message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
		if err != nil {
			logger.Log.Debugf("Error while subscribing: %v", err)
			e := fmt.Sprintf("Error during StopMonitoringSubscriptionRequest: %v", err)

			collectSubscriber.IncrementRetryCountFromMap(subscriptionRequests)

			message.Status = "Error"
			message.ErrorDetails = e

			continue
		}

		message.ResponseRawMessage = response.RawXML()
		message.ResponseSize = int64(len(message.ResponseRawMessage))
		message.ResponseIdentifier = response.ResponseMessageIdentifier()

		collectSubscriber.HandleResponse(subscriptionRequests, message, response)

		if len(subscriptionRequests) == 0 {
			return
		}
	}
	collectSubscriber.IncrementRetryCountFromMap(subscriptionRequests)
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
