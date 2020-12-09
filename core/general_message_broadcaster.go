package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/state"
)

type SIRIGeneralMessageBroadcaster interface {
	state.Stopable
	state.Startable
}

type GMBroadcaster struct {
	clock.ClockConsumer

	connector *SIRIGeneralMessageSubscriptionBroadcaster
}

type GeneralMessageBroadcaster struct {
	GMBroadcaster

	stop chan struct{}
}

type FakeGeneralMessageBroadcaster struct {
	GMBroadcaster

	clock.ClockConsumer
}

func NewFakeGeneralMessageBroadcaster(connector *SIRIGeneralMessageSubscriptionBroadcaster) SIRIGeneralMessageBroadcaster {
	broadcaster := &FakeGeneralMessageBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeGeneralMessageBroadcaster) Start() {
	broadcaster.prepareSIRIGeneralMessageNotify()
}

func (broadcaster *FakeGeneralMessageBroadcaster) Stop() {}

func NewSIRIGeneralMessageBroadcaster(connector *SIRIGeneralMessageSubscriptionBroadcaster) SIRIGeneralMessageBroadcaster {
	broadcaster := &GeneralMessageBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (gmb *GeneralMessageBroadcaster) Start() {
	logger.Log.Debugf("Start GeneralMessageBroadcaster")

	gmb.stop = make(chan struct{})
	go gmb.run()
}

func (gmb *GeneralMessageBroadcaster) run() {
	c := gmb.Clock().After(5 * time.Second)

	for {
		select {
		case <-gmb.stop:
			logger.Log.Debugf("general message broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRIGeneralMessageBroadcaster visit")

			gmb.prepareSIRIGeneralMessageNotify()

			c = gmb.Clock().After(5 * time.Second)
		}
	}
}

func (gmb *GeneralMessageBroadcaster) Stop() {
	if gmb.stop != nil {
		close(gmb.stop)
	}
}

func (gmb *GMBroadcaster) prepareSIRIGeneralMessageNotify() {
	gmb.connector.mutex.Lock()

	events := gmb.connector.toBroadcast
	gmb.connector.toBroadcast = make(map[SubscriptionId][]model.SituationId)

	gmb.connector.mutex.Unlock()

	tx := gmb.connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	for subId, situationIds := range events {
		sub, ok := gmb.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			continue
		}

		notify := siri.SIRINotifyGeneralMessage{
			Address:                   gmb.connector.Partner().Address(),
			ProducerRef:               gmb.connector.Partner().ProducerRef(),
			ResponseMessageIdentifier: gmb.connector.Partner().IdentifierGenerator(RESPONSE_MESSAGE_IDENTIFIER).NewMessageIdentifier(),
			SubscriberRef:             gmb.connector.SIRIPartner().SubscriberRef(),
			SubscriptionIdentifier:    sub.ExternalId(),
			RequestMessageRef:         sub.SubscriptionOption("MessageIdentifier"),
			Status:                    true,
			ResponseTimestamp:         gmb.Clock().Now(),
		}

		// Prepare Id Array
		// var messageArray []string

		builder := NewBroadcastGeneralMessageBuilder(tx, gmb.connector.Partner(), SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		builder.InfoChannelRef = strings.Split(sub.SubscriptionOption("InfoChannelRef"), ",")
		if sub.SubscriptionOption("LineRef") != "" {
			builder.SetLineRef(strings.Split(sub.SubscriptionOption("LineRef"), ","))
		}
		if sub.SubscriptionOption("StopPointRef") != "" {
			builder.SetStopPointRef(strings.Split(sub.SubscriptionOption("StopPointRef"), ","))
		}

		for _, situationId := range situationIds {
			situation, ok := tx.Model().Situations().Find(situationId)
			if !ok {
				logger.Log.Debugf("Could not find situation : %v in general message broadcaster", situationId)
				continue
			}

			siriGeneralMessage := builder.BuildGeneralMessage(situation)
			if siriGeneralMessage == nil {
				continue
			}
			// messageArray = append(messageArray, siriGeneralMessage.InfoMessageIdentifier)
			notify.GeneralMessages = append(notify.GeneralMessages, siriGeneralMessage)
		}
		if len(notify.GeneralMessages) != 0 {
			logStashEvent := gmb.newLogStashEvent()
			message := gmb.newBQEvent()

			logSIRIGeneralMessageNotify(logStashEvent, message, &notify)
			audit.CurrentLogStash().WriteEvent(logStashEvent)
			t := gmb.Clock().Now()

			err := gmb.connector.SIRIPartner().SOAPClient().NotifyGeneralMessage(&notify)
			message.ProcessingTime = gmb.Clock().Since(t).Seconds()
			if err != nil {
				event := gmb.newLogStashEvent()
				logSIRINotifyError(err.Error(), notify.ResponseMessageIdentifier, event)
				audit.CurrentLogStash().WriteEvent(event)
			}

			audit.CurrentBigQuery().WriteEvent(message)
		}
	}
}

func (gmb *GMBroadcaster) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "NotifyGeneralMessage",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(gmb.connector.partner.Slug()),
		Status:    "OK",
	}
}

func (gmb *GMBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := gmb.connector.partner.NewLogStashEvent()
	event["connector"] = "GeneralMessageSubscriptionBroadcaster"
	return event
}

func logSIRINotifyError(err, rmi string, logStashEvent audit.LogStashEvent) {
	logStashEvent["siriType"] = "NotifyError"
	logStashEvent["errorDescription"] = err
	logStashEvent["status"] = "false"
	logStashEvent["responseMessageIdentifier"] = rmi
}

func logSIRIGeneralMessageNotify(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, response *siri.SIRINotifyGeneralMessage) {
	message.RequestIdentifier = response.RequestMessageRef
	message.ResponseIdentifier = response.ResponseMessageIdentifier
	message.SubscriptionIdentifiers = []string{response.SubscriptionIdentifier}

	logStashEvent["siriType"] = "NotifyGeneralMessage"
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["subscriberRef"] = response.SubscriberRef
	logStashEvent["subscriptionIdentifier"] = response.SubscriptionIdentifier
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	if !response.Status {
		message.Status = "Error"
		logStashEvent["errorType"] = response.ErrorType
		if response.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		}
		logStashEvent["errorText"] = response.ErrorText
		message.ErrorDetails = response.ErrorString()
	}
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
	message.ResponseRawMessage = xml
	message.ResponseSize = int64(len(xml))
}
