package core

import (
	"fmt"
	"strings"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIGeneralMessageSubscriber interface {
	model.Stopable
	model.Startable
}

type GMSubscriber struct {
	model.ClockConsumer

	connector *SIRIGeneralMessageSubscriptionCollector
}

type GeneralMessageSubscriber struct {
	GMSubscriber

	stop chan struct{}
}

type FakeGeneralMessageSubscriber struct {
	GMSubscriber
}

type lineToRequest struct {
	subId  SubscriptionId
	lineId model.ObjectID
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
	subscriptions := subscriber.connector.partner.Subscriptions().FindSubscriptionsByKind("GeneralMessageCollect")
	if len(subscriptions) == 0 {
		logger.Log.Debugf("StopMonitoringSubscriber visit without GeneralMessageCollect subscriptions")
		return
	}

	// LineRef for Logstash
	lineRefList := []string{}

	linesToRequest := make(map[string]*lineToRequest)
	for _, subscription := range subscriptions {
		for _, resource := range subscription.ResourcesByObjectID() {
			if resource.SubscribedAt.IsZero() && resource.RetryCount <= 10 {
				messageIdentifier := subscriber.connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier()
				linesToRequest[messageIdentifier] = &lineToRequest{
					subId:  subscription.id,
					lineId: *(resource.Reference.ObjectId),
				}
			}
		}
	}

	if len(linesToRequest) == 0 {
		return
	}

	logStashEvent := subscriber.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	gmRequest := &siri.SIRIGeneralMessageSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	for messageIdentifier, requestedLine := range linesToRequest {
		entry := &siri.SIRIGeneralMessageSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier: string(requestedLine.subId),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = messageIdentifier
		entry.RequestTimestamp = subscriber.Clock().Now()
		entry.LineRef = []string{requestedLine.lineId.Value()}

		lineRefList = append(lineRefList, requestedLine.lineId.Value())
		gmRequest.Entries = append(gmRequest.Entries, entry)
	}

	logStashEvent["lineRefs"] = strings.Join(lineRefList, ", ")
	logSIRIGeneralMessageSubscriptionRequest(logStashEvent, gmRequest)

	response, err := subscriber.connector.SIRIPartner().SOAPClient().GeneralMessageSubscription(gmRequest)
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["response"] = fmt.Sprintf("Error during GeneralMessageSubscriptionRequest: %v", err)
		subscriber.incrementRetryCountFromMap(linesToRequest)
		return
	}

	logStashEvent["response"] = response.RawXML()

	for _, responseStatus := range response.ResponseStatus() {
		requestedLine, ok := linesToRequest[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		delete(linesToRequest, responseStatus.RequestMessageRef()) // See #4691

		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedLine.subId)
		if !ok { // Should never happen
			logger.Log.Debugf("Response for unknown subscription %v", requestedLine.subId)
			continue
		}
		resource := subscription.Resource(requestedLine.lineId)
		if resource == nil { // Should never happen
			logger.Log.Debugf("Response for unknown subscription resource %v", requestedLine.lineId.String())
			continue
		}

		if !responseStatus.Status() {
			logger.Log.Debugf("Subscription status false for line %v: %v %v ", requestedLine.lineId.Value(), responseStatus.ErrorType(), responseStatus.ErrorText())
			resource.RetryCount++
			continue
		}
		resource.SubscribedAt = subscriber.Clock().Now()
		resource.RetryCount = 0
	}
	// Should not happen but see #4691
	if len(linesToRequest) == 0 {
		return
	}
	subscriber.incrementRetryCountFromMap(linesToRequest)
}

func (subscriber *GMSubscriber) incrementRetryCountFromMap(linesToRequest map[string]*lineToRequest) {
	for _, requestedLine := range linesToRequest {
		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedLine.subId)
		if !ok { // Should never happen
			continue
		}
		resource := subscription.Resource(requestedLine.lineId)
		if resource == nil { // Should never happen
			continue
		}
		resource.RetryCount++
	}
}

func (smb *GMSubscriber) newLogStashEvent() audit.LogStashEvent {
	event := smb.connector.partner.NewLogStashEvent()
	event["connector"] = "GeneralMessageSubscriptionCollector"
	return event
}

func logSIRIGeneralMessageSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIGeneralMessageSubscriptionRequest) {
	logStashEvent["type"] = "GeneralMessageSubscriptionRequest"
	logStashEvent["consumerAddress"] = request.ConsumerAddress
	logStashEvent["messageIdentifier"] = request.MessageIdentifier
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	xml, err := request.BuildXML()
	if err != nil {
		logStashEvent["requestXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["requestXML"] = xml
}
