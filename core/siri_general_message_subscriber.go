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
	subscription, _ := subscriber.connector.partner.Subscriptions().FindOrCreateByKind("GeneralMessageCollect")
	lineRefList := []string{}

	linesToRequest := make(map[string]*model.ObjectID)
	for _, resource := range subscription.ResourcesByObjectID() {
		if resource.SubscribedAt.IsZero() {
			messageIdentifier := subscriber.connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier()
			linesToRequest[messageIdentifier] = resource.Reference.ObjectId
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

	for messageIdentifier, lineObjectid := range linesToRequest {
		entry := &siri.SIRIGeneralMessageSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier: string(subscription.Id()),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = messageIdentifier
		entry.RequestTimestamp = subscriber.Clock().Now()
		entry.LineRef = []string{lineObjectid.Value()}

		lineRefList = append(lineRefList, lineObjectid.Value())
		gmRequest.Entries = append(gmRequest.Entries, entry)
	}

	logStashEvent["lineRefs"] = strings.Join(lineRefList, ", ")
	logSIRIGeneralMessageSubscriptionRequest(logStashEvent, gmRequest)

	response, err := subscriber.connector.SIRIPartner().SOAPClient().GeneralMessageSubscription(gmRequest)
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["response"] = fmt.Sprintf("Error during GeneralMessageSubscriptionRequest: %v", err)
		for _, lineObjectid := range linesToRequest {
			resource := subscription.Resource(*lineObjectid)
			resource.RetryCount++
		}
		return
	}

	logStashEvent["response"] = response.RawXML()

	for _, responseStatus := range response.ResponseStatus() {
		lineObjectid, ok := linesToRequest[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		resource := subscription.Resource(*lineObjectid)
		if resource == nil { // Should never happen
			logger.Log.Debugf("Response for unknown subscription resource %v", lineObjectid.String())
			continue
		}
		if !responseStatus.Status() {
			logger.Log.Debugf("Subscription status false for line %v: %v %v ", lineObjectid.Value(), responseStatus.ErrorType(), responseStatus.ErrorText())
			resource.RetryCount++
			continue
		}
		resource.SubscribedAt = subscriber.Clock().Now()
		resource.RetryCount = 0
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
