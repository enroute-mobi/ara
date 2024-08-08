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

type SIRIVehicleMonitoringSubscriber interface {
	state.Stopable
	state.Startable
}

type VMSubscriber struct {
	clock.ClockConsumer

	connector *SIRIVehicleMonitoringSubscriptionCollector
}

type VehicleMonitoringSubscriber struct {
	VMSubscriber

	stop chan struct{}
}

type FakeVehicleMonitoringSubscriber struct {
	VMSubscriber
}

func NewFakeVehicleMonitoringSubscriber(connector *SIRIVehicleMonitoringSubscriptionCollector) SIRIVehicleMonitoringSubscriber {
	subscriber := &FakeVehicleMonitoringSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *FakeVehicleMonitoringSubscriber) Start() {
	subscriber.prepareSIRIVehicleMonitoringSubscriptionRequest()
}

func (subscriber *FakeVehicleMonitoringSubscriber) Stop() {}

func NewSIRIVehicleMonitoringSubscriber(connector *SIRIVehicleMonitoringSubscriptionCollector) SIRIVehicleMonitoringSubscriber {
	subscriber := &VehicleMonitoringSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *VehicleMonitoringSubscriber) Start() {
	logger.Log.Debugf("Start VehicleMonitoringSubscriber")

	subscriber.stop = make(chan struct{})
	go subscriber.run()
}

func (subscriber *VehicleMonitoringSubscriber) run() {
	c := subscriber.Clock().After(5 * time.Second)

	for {
		select {
		case <-subscriber.stop:
			return
		case <-c:
			logger.Log.Debugf("SIRIVehicleMonitoringSubscriber visit")

			subscriber.prepareSIRIVehicleMonitoringSubscriptionRequest()

			c = subscriber.Clock().After(5 * time.Second)
		}
	}
}

func (subscriber *VehicleMonitoringSubscriber) Stop() {
	if subscriber.stop != nil {
		close(subscriber.stop)
	}
}

func (subscriber *VMSubscriber) prepareSIRIVehicleMonitoringSubscriptionRequest() {
	collectSubscriber := NewCollectSubcriber(subscriber.connector, VehicleMonitoringCollect)
	subscriptionRequests := collectSubscriber.GetSubscriptionRequest()

	if len(subscriptionRequests) == 0 {
		return
	}

	message := subscriber.newBQEvent()
	defer audit.CurrentBigQuery(string(subscriber.connector.Partner().Referential().Slug())).WriteEvent(message)

	siriVehicleMonitoringSubscriptionRequest := &siri.SIRIVehicleMonitoringSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.Partner().NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.Partner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	subIds := []string{}
	linesToLog := []string{}
	for subId, subscriptionRequest := range subscriptionRequests {
		for _, m := range subscriptionRequest.modelsToRequest {
			entry := &siri.SIRIVehicleMonitoringSubscriptionRequestEntry{
				SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
				SubscriptionIdentifier: string(subId),
				InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
			}
			entry.MessageIdentifier = subscriptionRequest.requestMessageRef
			entry.RequestTimestamp = subscriber.Clock().Now()

			switch m.kind {
			case "Line":
				entry.LineRef = m.code.Value()
				linesToLog = append(linesToLog, entry.LineRef)
				subIds = append(subIds, string(subId))
			}
			siriVehicleMonitoringSubscriptionRequest.Entries = append(siriVehicleMonitoringSubscriptionRequest.Entries, entry)
		}
	}

	message.RequestIdentifier = siriVehicleMonitoringSubscriptionRequest.MessageIdentifier
	message.RequestRawMessage, _ = siriVehicleMonitoringSubscriptionRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
	message.RequestSize = int64(len(message.RequestRawMessage))
	message.Lines = linesToLog
	message.SubscriptionIdentifiers = subIds

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.Partner().SIRIClient().VehicleMonitoringSubscription(siriVehicleMonitoringSubscriptionRequest)
	message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		e := fmt.Sprintf("Error during VehicleMonitoringSubscriptionRequest: %v", err)

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

func (subscriber *VMSubscriber) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "VehicleMonitoringSubscriptionRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(subscriber.connector.partner.Slug()),
		Status:    "OK",
	}
}
