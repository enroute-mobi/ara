package core

import (
	"fmt"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringSubscriber interface {
	model.Stopable

	Run()
}

type SMSubscriber struct {
	model.ClockConsumer

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

func (subscriber *FakeStopMonitoringSubscriber) Run() {
	subscriber.prepareSIRIStopMonitoringSubscriptionRequest()
}

func (subscriber *FakeStopMonitoringSubscriber) Stop() {}

func NewSIRIStopMonitoringSubscriber(connector *SIRIStopMonitoringSubscriptionCollector) SIRIStopMonitoringSubscriber {
	subscriber := &StopMonitoringSubscriber{}
	subscriber.connector = connector
	return subscriber
}

func (subscriber *StopMonitoringSubscriber) Run() {
	if subscriber.stop != nil {
		return
	}

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
	subscription := subscriber.connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")

	stopAreasToRequest := make(map[string]*model.ObjectID)
	for _, resource := range subscription.ResourcesByObjectID() {
		if resource.SubscribedAt.IsZero() {
			messageIdentifier := subscriber.connector.SIRIPartner().NewMessageIdentifier()
			stopAreasToRequest[messageIdentifier] = resource.Reference.ObjectId
		}
	}

	if len(stopAreasToRequest) == 0 {
		return
	}

	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	siriStopMonitoringSubscriptionRequest := &siri.SIRIStopMonitoringSubscriptionRequest{
		ConsumerAddress:   subscriber.connector.Partner().Setting("local_url"),
		MessageIdentifier: subscriber.connector.SIRIPartner().NewMessageIdentifier(),
		RequestorRef:      subscriber.connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  subscriber.Clock().Now(),
	}

	for messageIdentifier, stopAreaObjectid := range stopAreasToRequest {
		entry := &siri.SIRIStopMonitoringSubscriptionRequestEntry{
			MessageIdentifier:      messageIdentifier,
			RequestTimestamp:       subscriber.Clock().Now(),
			MonitoringRef:          stopAreaObjectid.Value(),
			SubscriberRef:          subscriber.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier: fmt.Sprintf("Edwig:Subscription::%v:LOC", subscription.Id()),
			InitialTerminationTime: subscriber.Clock().Now().Add(48 * time.Hour),
		}

		siriStopMonitoringSubscriptionRequest.Entries = append(siriStopMonitoringSubscriptionRequest.Entries, entry)
	}

	logSIRIStopMonitoringSubscriptionRequest(logStashEvent, siriStopMonitoringSubscriptionRequest)
	// logStashEvent["StopAreasIds"] = strings.Join(stopAreasToRequest, ", ")

	response, err := subscriber.connector.SIRIPartner().SOAPClient().StopMonitoringSubscription(siriStopMonitoringSubscriptionRequest)
	if err != nil {
		logger.Log.Debugf("Error while subscribing: %v", err)
		return
	}

	logStashEvent["response"] = response.RawXML()

	for _, responseStatus := range response.ResponseStatus() {
		stopAreaObjectid, ok := stopAreasToRequest[responseStatus.RequestMessageRef()]
		if !ok {
			logger.Log.Debugf("ResponseStatus RequestMessageRef unknown: %v", responseStatus.RequestMessageRef())
			continue
		}
		if !responseStatus.Status() {
			logger.Log.Debugf("Subscription status false for stopArea %v: %v %v ", stopAreaObjectid.Value(), responseStatus.ErrorType(), responseStatus.ErrorText())
			continue
		}
		resource := subscription.Resource(*stopAreaObjectid)
		if resource != nil {
			logger.Log.Debugf("Response for unknown subscription resource %v", stopAreaObjectid.String())
		}
		resource.SubscribedAt = subscriber.Clock().Now()
	}

	// if response.Status() == true {
	// 	for _, stopAreaObjectid := range stopAreasToRequest {
	// resource := subscription.Resource(*stopAreaObjectid)
	// if resource != nil {
	// 	resource.SubscribedAt = subscriber.Clock().Now()
	// }
	// 	}
	// }
}

func logSIRIStopMonitoringSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIStopMonitoringSubscriptionRequest) {
	logStashEvent["Connector"] = "SIRIStopMonitoringSubscriber"
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
