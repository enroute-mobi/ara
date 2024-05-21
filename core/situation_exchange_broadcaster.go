package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/state"
	"golang.org/x/exp/maps"
)

type SIRISituationExchangeBroadcaster interface {
	state.Stopable
	state.Startable
}

type SXBroadcaster struct {
	clock.ClockConsumer

	connector *SIRISituationExchangeSubscriptionBroadcaster
}

type SituationExchangeBroadcaster struct {
	SXBroadcaster

	stop chan struct{}
}

type FakeSituationExchangeBroadcaster struct {
	SXBroadcaster

	clock.ClockConsumer
}

func NewFakeSituationExchangeBroadcaster(connector *SIRISituationExchangeSubscriptionBroadcaster) SIRISituationExchangeBroadcaster {
	broadcaster := &FakeSituationExchangeBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeSituationExchangeBroadcaster) Start() {
	broadcaster.prepareSIRISituationExchangeNotify()
}

func (broadcaster *FakeSituationExchangeBroadcaster) Stop() {}

func NewSIRISituationExchangeBroadcaster(connector *SIRISituationExchangeSubscriptionBroadcaster) SIRISituationExchangeBroadcaster {
	broadcaster := &SituationExchangeBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (sxb *SituationExchangeBroadcaster) Start() {
	logger.Log.Debugf("Start SituationExchangeBroadcaster")

	sxb.stop = make(chan struct{})
	go sxb.run()
}

func (sxb *SituationExchangeBroadcaster) run() {
	c := sxb.Clock().After(5 * time.Second)

	for {
		select {
		case <-sxb.stop:
			logger.Log.Debugf("situation exchange broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRISituationExchangeBroadcaster visit")

			sxb.prepareSIRISituationExchangeNotify()

			c = sxb.Clock().After(5 * time.Second)
		}
	}
}

func (sxb *SituationExchangeBroadcaster) Stop() {
	if sxb.stop != nil {
		close(sxb.stop)
	}
}

func (sxb *SXBroadcaster) prepareSIRISituationExchangeNotify() {
	sxb.connector.mutex.Lock()

	events := sxb.connector.toBroadcast
	sxb.connector.toBroadcast = make(map[SubscriptionId][]model.SituationId)

	sxb.connector.mutex.Unlock()

	notify := siri.SIRINotifySituationExchange{
		Address:                   sxb.connector.Partner().Address(),
		ProducerRef:               sxb.connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: sxb.connector.Partner().NewResponseMessageIdentifier(),
		ResponseTimestamp:         sxb.Clock().Now(),
		Status:                    true,
	}

	for subId, situationIds := range events {
		sub, ok := sxb.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			logger.Log.Debugf("SX subscriptionBroadcast Could not find subscription with id : %v", subId)
			continue
		}
		processedSituations := make(map[model.SituationId]struct{}) //Making sure not to send 2 times the same Situation

		// Initialize builder
		builder := NewBroadcastSituationExchangeBuilder(sxb.connector.Partner(), SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_BROADCASTER)

		delivery := &siri.SIRISituationExchangeDelivery{
			ResponseTimestamp:      builder.Clock().Now(),
			Status:                 true,
			SubscriberRef:          sub.SubscriberRef,
			SubscriptionIdentifier: sub.ExternalId(),
			LineRefs:               make(map[string]struct{}),
			MonitoringRefs:         make(map[string]struct{}),
		}

		for _, situationId := range situationIds {
			if _, ok := processedSituations[situationId]; ok {
				continue
			}

			// Find the Situation
			situation, ok := sxb.connector.Partner().Model().Situations().Find(situationId)
			if !ok {
				logger.Log.Debugf("Could not find situation : %v in situation exchange broadcaster", situationId)
				continue
			}

			builder.BuildSituationExchange(situation, delivery)
		}

		if len(delivery.Situations) == 0 {
			continue
		}

		notify.SituationExchangeDeliveries = append(notify.SituationExchangeDeliveries, delivery)
	}

	if len(notify.SituationExchangeDeliveries) != 0 {
		message := sxb.newBQEvent()

		sxb.logSIRISituationExchangeNotify(message, &notify)
		t := sxb.Clock().Now()

		sxb.connector.Partner().SIRIClient().NotifySituationExchange(&notify)
		message.ProcessingTime = sxb.Clock().Since(t).Seconds()

		audit.CurrentBigQuery(string(sxb.connector.Partner().Referential().Slug())).WriteEvent(message)
	}
}

func (sxb *SXBroadcaster) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.NOTIFY_SITUATION_EXCHANGE,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(sxb.connector.Partner().Slug()),
		Status:    "OK",
	}
}

func (sxb *SXBroadcaster) logSIRISituationExchangeNotify(message *audit.BigQueryMessage, response *siri.SIRINotifySituationExchange) {
	lineRefs := make(map[string]struct{})
	monitoringRefs := make(map[string]struct{})

	for _, delivery := range response.SituationExchangeDeliveries {
		maps.Copy(lineRefs, delivery.LineRefs)
		maps.Copy(monitoringRefs, delivery.MonitoringRefs)
		message.SubscriptionIdentifiers = append(message.SubscriptionIdentifiers, delivery.SubscriptionIdentifier)
	}

	message.RequestIdentifier = response.RequestMessageRef
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	message.StopAreas = GetModelReferenceSlice(monitoringRefs)
	message.Lines = GetModelReferenceSlice(lineRefs)

	if !response.Status {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	xml, err := response.BuildXML(sxb.connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.ResponseRawMessage = xml
	message.ResponseSize = int64(len(xml))
}
