package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/edwig/audit"
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
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

type resourceToRequest struct {
	subId    SubscriptionId
	objectId model.ObjectID
	kind     string
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
		logger.Log.Debugf("GeneralMessageSubscriber visit without GeneralMessageCollect subscriptions")
		return
	}

	// LineRef for Logstash
	lineRefList := []string{}
	stopPointRefList := []string{}

	resourcesToRequest := make(map[string]*resourceToRequest)
	for _, subscription := range subscriptions {
		for _, resource := range subscription.ResourcesByObjectIDCopy() {
			if resource.SubscribedAt.IsZero() && resource.RetryCount <= 10 {
				messageIdentifier := subscriber.connector.Partner().IdentifierGenerator(MESSAGE_IDENTIFIER).NewMessageIdentifier()
				logger.Log.Debugf("send request for subscription with id : %v", subscription.id)
				resourcesToRequest[messageIdentifier] = &resourceToRequest{
					subId:    subscription.id,
					objectId: *(resource.Reference.ObjectId),
					kind:     resource.Reference.Type,
				}
			}
		}
	}

	if len(resourcesToRequest) == 0 {
		return
	}

	logStashEvent := subscriber.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	gmRequest := &siri.SIRIGeneralMessageSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Address(),
		MessageIdentifier: subscriber.connector.Partner().IdentifierGenerator(MESSAGE_IDENTIFIER).NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	for messageIdentifier, requestedResource := range resourcesToRequest {
		entry := &siri.SIRIGeneralMessageSubscriptionRequestEntry{
			SubscriberRef:          subscriber.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier: string(requestedResource.subId),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}
		entry.MessageIdentifier = messageIdentifier
		entry.RequestTimestamp = subscriber.Clock().Now()
		switch requestedResource.kind {
		case "Line":
			entry.LineRef = []string{requestedResource.objectId.Value()}
			lineRefList = append(lineRefList, requestedResource.objectId.Value())
		case "StopArea":
			entry.StopPointRef = []string{requestedResource.objectId.Value()}
			stopPointRefList = append(stopPointRefList, requestedResource.objectId.Value())
		}

		if b, _ := strconv.ParseBool(subscriber.connector.partner.Setting(GENEREAL_MESSAGE_REQUEST_2)); b {
			entry.XsdInWsdl = true
		}

		gmRequest.Entries = append(gmRequest.Entries, entry)
	}

	logStashEvent["lineRefs"] = strings.Join(lineRefList, ", ")
	logStashEvent["stopPointRefs"] = strings.Join(stopPointRefList, ", ")
	logSIRIGeneralMessageSubscriptionRequest(logStashEvent, gmRequest)

	startTime := subscriber.Clock().Now()
	response, err := subscriber.connector.SIRIPartner().SOAPClient().GeneralMessageSubscription(gmRequest)
	logStashEvent["responseTime"] = subscriber.Clock().Since(startTime).String()
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = fmt.Sprintf("Error during GeneralMessageSubscriptionRequest: %v", err)
		subscriber.incrementRetryCountFromMap(resourcesToRequest)
		return
	}

	logStashEvent["responseXML"] = response.RawXML()

	for _, responseStatus := range response.ResponseStatus() {
		requestedResource, ok := resourcesToRequest[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		delete(resourcesToRequest, responseStatus.RequestMessageRef()) // See #4691

		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedResource.subId)
		if !ok { // Should never happen
			logger.Log.Debugf("Response for unknown subscription %v", requestedResource.subId)
			continue
		}
		resource := subscription.Resource(requestedResource.objectId)
		if resource == nil { // Should never happen
			logger.Log.Debugf("Response for unknown subscription resource %v", requestedResource.objectId.String())
			continue
		}

		if !responseStatus.Status() {
			logger.Log.Debugf("Subscription status false for line %v: %v %v ", requestedResource.objectId.Value(), responseStatus.ErrorType(), responseStatus.ErrorText())
			resource.RetryCount++
			continue
		}
		resource.SubscribedAt = subscriber.Clock().Now()
		resource.RetryCount = 0
	}
	// Should not happen but see #4691
	if len(resourcesToRequest) == 0 {
		return
	}
	subscriber.incrementRetryCountFromMap(resourcesToRequest)
}

func (subscriber *GMSubscriber) incrementRetryCountFromMap(resourcesToRequest map[string]*resourceToRequest) {
	for _, requestedResource := range resourcesToRequest {
		subscription, ok := subscriber.connector.partner.Subscriptions().Find(requestedResource.subId)
		if !ok { // Should never happen
			continue
		}
		resource := subscription.Resource(requestedResource.objectId)
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
	logStashEvent["siriType"] = "GeneralMessageSubscriptionRequest"
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
