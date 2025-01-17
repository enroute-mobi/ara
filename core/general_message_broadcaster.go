package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/state"

	"strings"
	"time"
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

	for subId, situationIds := range events {
		sub, ok := gmb.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			continue
		}

		notify := siri.SIRINotifyGeneralMessage{
			Address:                   gmb.connector.Partner().Address(),
			ProducerRef:               gmb.connector.Partner().ProducerRef(),
			ResponseMessageIdentifier: gmb.connector.Partner().NewResponseMessageIdentifier(),
			SubscriberRef:             sub.SubscriberRef,
			SubscriptionIdentifier:    sub.ExternalId(),
			RequestMessageRef:         sub.SubscriptionOption("MessageIdentifier"),
			Status:                    true,
			ResponseTimestamp:         gmb.Clock().Now(),
		}

		// Prepare Id Array

		builder := NewBroadcastGeneralMessageBuilder(gmb.connector.Partner(), SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		builder.InfoChannelRef = strings.Split(sub.SubscriptionOption("InfoChannelRef"), ",")
		if sub.SubscriptionOption("LineRef") != "" {
			builder.SetLineRef(strings.Split(sub.SubscriptionOption("LineRef"), ","))
		}
		if sub.SubscriptionOption("StopPointRef") != "" {
			builder.SetStopPointRef(strings.Split(sub.SubscriptionOption("StopPointRef"), ","))
		}

		for _, situationId := range situationIds {
			situation, ok := gmb.connector.Partner().Model().Situations().Find(situationId)
			if !ok {
				logger.Log.Debugf("Could not find situation : %v in general message broadcaster", situationId)
				continue
			}

			if situation.Progress == model.SituationProgressClosed {
				siriGeneralMessageCancellation := builder.BuildGeneralMessageCancellation(situation)
				if siriGeneralMessageCancellation == nil {
					continue
				}

				notify.GeneralMessageCancellations = append(notify.GeneralMessageCancellations, siriGeneralMessageCancellation)
			} else {
				siriGeneralMessage := builder.BuildGeneralMessage(situation)
				if siriGeneralMessage == nil {
					continue
				}
				notify.GeneralMessages = append(notify.GeneralMessages, siriGeneralMessage)
			}

		}

		if len(notify.GeneralMessages) != 0 || len(notify.GeneralMessageCancellations) != 0 {
			message := gmb.newBQEvent()

			gmb.logSIRIGeneralMessageNotify(message, &notify)
			t := gmb.Clock().Now()

			gmb.connector.Partner().SIRIClient().NotifyGeneralMessage(&notify)
			message.ProcessingTime = gmb.Clock().Since(t).Seconds()

			audit.CurrentBigQuery(string(gmb.connector.Partner().Referential().Slug())).WriteEvent(message)
		}
	}
}

func (gmb *GMBroadcaster) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "NotifyGeneralMessage",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(gmb.connector.Partner().Slug()),
		Status:    "OK",
	}
}

func (gmb *GMBroadcaster) logSIRIGeneralMessageNotify(message *audit.BigQueryMessage, response *siri.SIRINotifyGeneralMessage) {
	lineRefs := make(map[string]struct{})
	monitoringRefs := make(map[string]struct{})

	for _, message := range response.GeneralMessages {
		for _, affectedRef := range message.AffectedRefs {
			switch affectedRef.Kind {
			case "LineRef":
				lineRefs[affectedRef.Id] = struct{}{}
			case "StopPointRef", "DestinationRef":
				monitoringRefs[affectedRef.Id] = struct{}{}
			}
		}
		for _, affectedLineSection := range message.LineSections {
			lineRefs[affectedLineSection.LineRef] = struct{}{}
			monitoringRefs[affectedLineSection.FirstStop] = struct{}{}
			monitoringRefs[affectedLineSection.LastStop] = struct{}{}
		}
	}

	message.RequestIdentifier = response.RequestMessageRef
	message.ResponseIdentifier = response.ResponseMessageIdentifier
	message.SubscriptionIdentifiers = []string{response.SubscriptionIdentifier}

	message.StopAreas = GetModelReferenceSlice(monitoringRefs)
	message.Lines = GetModelReferenceSlice(lineRefs)

	if !response.Status {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	xml, err := response.BuildXML(gmb.connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.ResponseRawMessage = xml
	message.ResponseSize = int64(len(xml))
}
