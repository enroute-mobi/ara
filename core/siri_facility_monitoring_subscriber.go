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

type SIRIFacilityMonitoringSubscriber interface {
	state.Stopable
	state.Startable
}

type FMSubscriber struct {
	clock.ClockConsumer

	connector *SIRIFacilityMonitoringSubscriptionCollector
}

type FacilityMonitoringSubscriber struct {
	FMSubscriber

	stop chan struct{}
}

type FakeFacilityMonitoringSubscriber struct {
	FMSubscriber
}

func NewFakeFacilityMonitoringSubscriber(connector *SIRIFacilityMonitoringSubscriptionCollector) SIRIFacilityMonitoringSubscriber {
	subscriber := &FakeFacilityMonitoringSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *FakeFacilityMonitoringSubscriber) Start() {
	subscriber.prepareSIRIFacilityMonitoringSubscriptionRequest()
}

func (subscriber *FakeFacilityMonitoringSubscriber) Stop() {}

func NewSIRIFacilityMonitoringSubscriber(connector *SIRIFacilityMonitoringSubscriptionCollector) SIRIFacilityMonitoringSubscriber {
	subscriber := &FacilityMonitoringSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *FacilityMonitoringSubscriber) Start() {
	logger.Log.Debugf("Start FacilityMonitoringSubscriber")

	subscriber.stop = make(chan struct{})
	go subscriber.run()
}

func (subscriber *FacilityMonitoringSubscriber) run() {
	c := subscriber.Clock().After(5 * time.Second)

	for {
		select {
		case <-subscriber.stop:
			return
		case <-c:
			logger.Log.Debugf("SIRIFacilityMonitoringSubscriber visit")

			subscriber.prepareSIRIFacilityMonitoringSubscriptionRequest()

			c = subscriber.Clock().After(5 * time.Second)
		}
	}
}

func (subscriber *FacilityMonitoringSubscriber) Stop() {
	if subscriber.stop != nil {
		close(subscriber.stop)
	}
}

func (subscriber *FMSubscriber) prepareSIRIFacilityMonitoringSubscriptionRequest() {
	collectSubscriber := NewCollectSubcriber(subscriber.connector, FacilityMonitoringCollect)
	subscriptionRequests := collectSubscriber.GetSubscriptionRequest()

	if len(subscriptionRequests) == 0 {
		return
	}

	message := subscriber.newBQEvent()
	defer audit.CurrentBigQuery(string(subscriber.connector.Partner().Referential().Slug())).WriteEvent(message)

	siriFacilityMonitoringSubscriptionRequest := &siri.SIRIFacilityMonitoringSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.Partner().NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.Partner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	subIds := []string{}
	for subId, subscriptionRequest := range subscriptionRequests {
		for _, m := range subscriptionRequest.modelsToRequest {
			entry := &siri.SIRIFacilityMonitoringSubscriptionRequestEntry{
				SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
				SubscriptionIdentifier: string(subId),
				InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
			}
			entry.MessageIdentifier = subscriptionRequest.requestMessageRef
			entry.RequestTimestamp = subscriber.Clock().Now()

			switch m.kind {
			case "Facility":
				entry.FacilityRef = m.code.Value()
				subIds = append(subIds, string(subId))
			}
			siriFacilityMonitoringSubscriptionRequest.Entries = append(siriFacilityMonitoringSubscriptionRequest.Entries, entry)
		}
	}

	message.RequestIdentifier = siriFacilityMonitoringSubscriptionRequest.MessageIdentifier
	message.RequestRawMessage, _ = siriFacilityMonitoringSubscriptionRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
	message.RequestSize = int64(len(message.RequestRawMessage))
	message.SubscriptionIdentifiers = subIds

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.Partner().SIRIClient().FacilityMonitoringSubscription(siriFacilityMonitoringSubscriptionRequest)
	message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		e := fmt.Sprintf("Error during FacilityMonitoringSubscriptionRequest: %v", err)

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

func (subscriber *FMSubscriber) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.FACILITY_MONITORING_SUBSCRIPTION_REQUEST,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(subscriber.connector.partner.Slug()),
		Status:    "OK",
	}
}
