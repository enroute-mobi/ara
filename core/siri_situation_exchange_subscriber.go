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

type FakeSituationExchangeSubscriber struct {
	SXSubscriber
}

func NewFakeSituationExchangeSubscriber(connector *SIRISituationExchangeSubscriptionCollector) SIRISituationExchangeSubscriber {
	subscriber := &FakeSituationExchangeSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *FakeSituationExchangeSubscriber) Start() {
	subscriber.prepareSIRISituationExchangeSubscriptionRequest()
}

func (subscriber *FakeSituationExchangeSubscriber) Stop() {}

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
	collectSubscriber := NewCollectSubcriber(subscriber.connector, SituationExchangeCollect)
	subscriptionRequests := collectSubscriber.GetSubscriptionRequest()

	if len(subscriptionRequests) == 0 {
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

	subIds := []string{}
	linesToLog := []string{}
	stopAreasToLog := []string{}
	for subId, subscriptionRequest := range subscriptionRequests {
		entry := &siri.SIRISituationExchangeSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
			SubscriptionIdentifier: string(subId),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = subscriptionRequest.requestMessageRef
		entry.RequestTimestamp = subscriber.Clock().Now()
		for _, m := range subscriptionRequest.modelsToRequest {
			switch m.kind {
			case "Line":
				entry.LineRef = append(entry.LineRef, m.code.Value())
			case "StopArea":
				entry.StopPointRef = append(entry.StopPointRef, m.code.Value())
			}
		}

		linesToLog = append(linesToLog, entry.LineRef...)
		stopAreasToLog = append(stopAreasToLog, entry.StopPointRef...)
		sxRequest.Entries = append(sxRequest.Entries, entry)
		subIds = append(subIds, string(subId))
	}

	message.RequestIdentifier = sxRequest.MessageIdentifier
	message.RequestRawMessage, _ = sxRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
	message.RequestSize = int64(len(message.RequestRawMessage))
	message.StopAreas = stopAreasToLog
	message.Lines = linesToLog
	message.SubscriptionIdentifiers = subIds

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.Partner().SIRIClient().SituationExchangeSubscription(sxRequest)
	message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		e := fmt.Sprintf("Error during SituationExchangeSubscriptionRequest: %v", err)
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

func (subscriber *SXSubscriber) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.SITUATION_EXCHANGE_SUBSCRIPTION_REQUEST,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(subscriber.connector.partner.Slug()),
		Status:    "OK",
	}
}
