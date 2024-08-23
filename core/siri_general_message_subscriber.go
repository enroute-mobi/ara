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

type SIRIGeneralMessageSubscriber interface {
	state.Stopable
	state.Startable
}

type GMSubscriber struct {
	clock.ClockConsumer

	connector *SIRIGeneralMessageSubscriptionCollector
}

type GeneralMessageSubscriber struct {
	GMSubscriber

	stop chan struct{}
}

type FakeGeneralMessageSubscriber struct {
	GMSubscriber
}

func NewFakeGeneralMessageSubscriber(connector *SIRIGeneralMessageSubscriptionCollector) SIRIGeneralMessageSubscriber {
	subscriber := &FakeGeneralMessageSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *FakeGeneralMessageSubscriber) Start() {
	subscriber.prepareSIRIGeneralMessageSubscriptionRequest()
}

func (subscriber *FakeGeneralMessageSubscriber) Stop() {}

func NewSIRIGeneralMessageSubscriber(connector *SIRIGeneralMessageSubscriptionCollector) SIRIGeneralMessageSubscriber {
	subscriber := &GeneralMessageSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *GeneralMessageSubscriber) Start() {
	logger.Log.Debugf("Start GeneralMessageSubscriber")

	subscriber.stop = make(chan struct{})
	go subscriber.run()
}

func (subscriber *GeneralMessageSubscriber) run() {
	c := subscriber.Clock().After(5 * time.Second)

	for {
		select {
		case <-subscriber.stop:
			return
		case <-c:
			logger.Log.Debugf("SIRIGeneralMessageSubscriber visit")

			subscriber.prepareSIRIGeneralMessageSubscriptionRequest()

			c = subscriber.Clock().After(5 * time.Second)
		}
	}
}

func (subscriber *GeneralMessageSubscriber) Stop() {
	if subscriber.stop != nil {
		close(subscriber.stop)
	}
}

func (subscriber *GMSubscriber) prepareSIRIGeneralMessageSubscriptionRequest() {
	collectSubscriber := NewCollectSubcriber(subscriber.connector, GeneralMessageCollect)
	subscriptionRequests := collectSubscriber.GetSubscriptionRequest()

	if len(subscriptionRequests) == 0 {
		return
	}

	message := subscriber.newBQEvent()
	defer audit.CurrentBigQuery(string(subscriber.connector.Partner().Referential().Slug())).WriteEvent(message)

	gmRequest := &siri.SIRIGeneralMessageSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.Partner().NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.Partner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	subIds := []string{}
	linesToLog := []string{}
	stopAreasToLog := []string{}
	for subId, subscriptionRequest := range subscriptionRequests {
		for _, m := range subscriptionRequest.modelsToRequest {
			entry := &siri.SIRIGeneralMessageSubscriptionRequestEntry{
				SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
				SubscriptionIdentifier: string(subId),
				InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
			}
			entry.MessageIdentifier = subscriptionRequest.requestMessageRef
			entry.RequestTimestamp = subscriber.Clock().Now()

			switch m.kind {
			case "Line":
				entry.LineRef = append(entry.LineRef, m.code.Value())
			case "StopArea":
				entry.StopPointRef = append(entry.StopPointRef, m.code.Value())
			}

			if subscriber.connector.Partner().GeneralMessageRequestVersion22() {
				entry.XsdInWsdl = true
			}

			linesToLog = append(linesToLog, entry.LineRef...)
			stopAreasToLog = append(stopAreasToLog, entry.StopPointRef...)
			gmRequest.Entries = append(gmRequest.Entries, entry)

		}
		subIds = append(subIds, string(subId))
	}

	message.RequestIdentifier = gmRequest.MessageIdentifier
	message.RequestRawMessage, _ = gmRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
	message.RequestSize = int64(len(message.RequestRawMessage))
	message.StopAreas = stopAreasToLog
	message.Lines = linesToLog
	message.SubscriptionIdentifiers = subIds

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.Partner().SIRIClient().GeneralMessageSubscription(gmRequest)
	message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		e := fmt.Sprintf("Error during GeneralMessageSubscriptionRequest: %v", err)
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

func (subscriber *GMSubscriber) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "GeneralMessageSubscriptionRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(subscriber.connector.partner.Slug()),
		Status:    "OK",
	}
}
