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

type lineToRequest struct {
	subID SubscriptionId
	lID   model.ObjectID
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
	subscriptions := subscriber.connector.partner.Subscriptions().FindSubscriptionsByKind(VehicleMonitoringCollect)
	if len(subscriptions) == 0 {
		logger.Log.Debugf("VehicleMonitoringSubscriber visit without VehicleMonitoringCollect subscriptions")
		return
	}

	linesList := []string{}

	linesToRequest := make(map[string]*lineToRequest)
	for _, subscription := range subscriptions {
		for _, resource := range subscription.ResourcesByObjectIDCopy() {
			if resource.SubscribedAt().IsZero() && resource.RetryCount <= 10 {
				messageIdentifier := subscriber.connector.Partner().NewMessageIdentifier()
				linesToRequest[messageIdentifier] = &lineToRequest{
					subID: subscription.id,
					lID:   *(resource.Reference.ObjectId),
				}
			}
		}
	}

	if len(linesToRequest) == 0 {
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

	var subIDs []string
	for messageIdentifier, requestedLines := range linesToRequest {
		entry := &siri.SIRIVehicleMonitoringSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.Partner().RequestorRef(),
			SubscriptionIdentifier: string(requestedLines.subID),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = messageIdentifier
		entry.RequestTimestamp = subscriber.Clock().Now()
		entry.LineRef = requestedLines.lID.Value()

		linesList = append(linesList, entry.LineRef)
		subIDs = append(subIDs, entry.SubscriptionIdentifier)
		siriVehicleMonitoringSubscriptionRequest.Entries = append(siriVehicleMonitoringSubscriptionRequest.Entries, entry)
	}

	message.RequestIdentifier = siriVehicleMonitoringSubscriptionRequest.MessageIdentifier
	message.RequestRawMessage, _ = siriVehicleMonitoringSubscriptionRequest.BuildXML(subscriber.connector.Partner().SIRIEnvelopeType())
	message.RequestSize = int64(len(message.RequestRawMessage))
	message.Lines = linesList
	message.SubscriptionIdentifiers = subIDs

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.Partner().SIRIClient().VehicleMonitoringSubscription(siriVehicleMonitoringSubscriptionRequest)
	message.ProcessingTime = subscriber.Clock().Since(startTime).Seconds()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		e := fmt.Sprintf("Error during VehicleMonitoringSubscriptionRequest: %v", err)

		subscriber.incrementRetryCountFromMap(linesToRequest)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
	message.ResponseIdentifier = response.ResponseMessageIdentifier()

	for _, responseStatus := range response.ResponseStatus() {
		requestedLines, ok := linesToRequest[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		delete(linesToRequest, responseStatus.RequestMessageRef()) // See #4691

		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedLines.subID)
		if !ok { // Should never happen
			logger.Log.Debugf("Response for unknown subscription %v", requestedLines.subID)
			continue
		}
		resource := subscription.Resource(requestedLines.lID)
		if resource == nil { // Should never happen
			logger.Log.Debugf("Response for unknown subscription resource %v", requestedLines.lID.String())
			continue
		}

		if !responseStatus.Status() {
			logger.Log.Debugf("Subscription status false for line %v: %v %v ", requestedLines.lID.Value(), responseStatus.ErrorType(), responseStatus.ErrorText())
			resource.RetryCount++
			message.Status = "Error"
			continue
		}
		resource.Subscribed(subscriber.Clock().Now())
		resource.RetryCount = 0
	}
	// Should not happen but see #4691
	if len(linesToRequest) == 0 {
		return
	}
	subscriber.incrementRetryCountFromMap(linesToRequest)
}

func (subscriber *VMSubscriber) incrementRetryCountFromMap(linesToRequest map[string]*lineToRequest) {
	for _, requestedLines := range linesToRequest {
		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedLines.subID)
		if !ok { // Should never happen
			continue
		}
		resource := subscription.Resource(requestedLines.lID)
		if resource == nil { // Should never happen
			continue
		}
		resource.RetryCount++
	}
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
