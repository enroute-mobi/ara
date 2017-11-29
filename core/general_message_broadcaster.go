package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIGeneralMessageBroadcaster interface {
	model.Stopable
	model.Startable
}

type GMBroadcaster struct {
	model.ClockConsumer

	connector *SIRIGeneralMessageSubscriptionBroadcaster
}

type GeneralMessageBroadcaster struct {
	GMBroadcaster

	stop chan struct{}
}

type FakeGeneralMessageBroadcaster struct {
	GMBroadcaster

	model.ClockConsumer
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
			ResponseMessageIdentifier: gmb.connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
			SubscriberRef:             gmb.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier:    sub.ExternalId(),
			RequestMessageRef:         sub.SubscriptionOptions()["MessageIdentifier"],
			Status:                    true,
		}

		// Prepare Id Array
		var messageArray []string

		builder := NewBroadcastGeneralMessageBuilder(tx, gmb.connector.SIRIPartner(), SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		builder.InfoChannelRef = strings.Split(sub.SubscriptionOptions()["InfoChannelRef"], ",")
		if sub.SubscriptionOptions()["LineRef"] != "" {
			builder.SetLineRef(strings.Split(sub.SubscriptionOptions()["LineRef"], ","))
		}
		if sub.SubscriptionOptions()["StopPointRef"] != "" {
			builder.SetStopPointRef(strings.Split(sub.SubscriptionOptions()["StopPointRef"], ","))
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
			messageArray = append(messageArray, siriGeneralMessage.InfoMessageIdentifier)
			notify.GeneralMessages = append(notify.GeneralMessages, siriGeneralMessage)
		}
		if len(notify.GeneralMessages) != 0 {
			logStashEvent := gmb.newLogStashEvent()

			logSIRIGeneralMessageNotify(logStashEvent, &notify)
			audit.CurrentLogStash().WriteEvent(logStashEvent)

			err := gmb.connector.SIRIPartner().SOAPClient().NotifyGeneralMessage(&notify)
			if err != nil {
				event := gmb.newLogStashEvent()
				logSIRINotifyError(err.Error(), event)
				audit.CurrentLogStash().WriteEvent(event)
			}
		}
	}
}

func (gmb *GMBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := gmb.connector.partner.NewLogStashEvent()
	event["connector"] = "GeneralMessageSubscriptionBroadcaster"
	return event
}

func logSIRINotifyError(err string, logStashEvent audit.LogStashEvent) {
	logStashEvent["type"] = "NotifyError"
	logStashEvent["Error"] = err
}

func logSIRIGeneralMessageNotify(logStashEvent audit.LogStashEvent, response *siri.SIRINotifyGeneralMessage) {
	logStashEvent["type"] = "NotifyGeneralMessage"
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	if !response.Status {
		logStashEvent["errorType"] = response.ErrorType
		if response.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		}
		logStashEvent["errorText"] = response.ErrorText
	}
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
