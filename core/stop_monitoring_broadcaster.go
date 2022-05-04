package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/state"
)

type SIRIStopMonitoringBroadcaster interface {
	state.Stopable
	state.Startable
}

type SMBroadcaster struct {
	clock.ClockConsumer

	connector *SIRIStopMonitoringSubscriptionBroadcaster
}

type StopMonitoringBroadcaster struct {
	SMBroadcaster

	stop chan struct{}
}

type FakeStopMonitoringBroadcaster struct {
	SMBroadcaster

	clock.ClockConsumer
}

func NewFakeStopMonitoringBroadcaster(connector *SIRIStopMonitoringSubscriptionBroadcaster) SIRIStopMonitoringBroadcaster {
	broadcaster := &FakeStopMonitoringBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeStopMonitoringBroadcaster) Start() {
	broadcaster.prepareSIRIStopMonitoringNotify()
}

func (broadcaster *FakeStopMonitoringBroadcaster) Stop() {}

func NewSIRIStopMonitoringBroadcaster(connector *SIRIStopMonitoringSubscriptionBroadcaster) SIRIStopMonitoringBroadcaster {
	broadcaster := &StopMonitoringBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (smb *StopMonitoringBroadcaster) Start() {
	logger.Log.Debugf("Start StopMonitoringBroadcaster")

	smb.stop = make(chan struct{})
	go smb.run()
}

func (smb *StopMonitoringBroadcaster) run() {
	c := smb.Clock().After(5 * time.Second)

	for {
		select {
		case <-smb.stop:
			logger.Log.Debugf("stop monitoring broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRIStopMonitoringBroadcaster visit")

			smb.prepareSIRIStopMonitoringNotify()
			smb.prepareNotMonitored()

			c = smb.Clock().After(5 * time.Second)
		}
	}
}

func (smb *StopMonitoringBroadcaster) Stop() {
	if smb.stop != nil {
		close(smb.stop)
	}
}

func (smb *SMBroadcaster) prepareNotMonitored() {
	smb.connector.mutex.Lock()

	notMonitored := smb.connector.notMonitored
	smb.connector.notMonitored = make(map[SubscriptionId]map[string]struct{})

	smb.connector.mutex.Unlock()

	for subId, producers := range notMonitored {
		sub, ok := smb.connector.Partner().Subscriptions().Find(subId)
		if !ok || len(producers) == 0 {
			continue
		}

		for producer := range producers {
			notification := &siri.SIRINotifyStopMonitoring{
				Address:                   smb.connector.Partner().Address(),
				ProducerRef:               smb.connector.Partner().ProducerRef(),
				RequestMessageRef:         sub.SubscriptionOption("MessageIdentifier"),
				ResponseMessageIdentifier: smb.connector.Partner().NewResponseMessageIdentifier(),
				ResponseTimestamp:         smb.connector.Clock().Now(),
			}

			delivery := &siri.SIRINotifyStopMonitoringDelivery{
				SubscriberRef:          sub.SubscriberRef,
				SubscriptionIdentifier: sub.ExternalId(),
				RequestMessageRef:      sub.SubscriptionOption("MessageIdentifier"),
				ResponseTimestamp:      smb.connector.Clock().Now(),
				Status:                 false,
				ErrorType:              "OtherError",
				ErrorNumber:            1,
				ErrorText:              fmt.Sprintf("Erreur [PRODUCER_UNAVAILABLE] : %v indisponible", producer),
			}

			notification.Deliveries = []*siri.SIRINotifyStopMonitoringDelivery{delivery}

			smb.sendNotification(notification)
		}
	}
}

func (smb *SMBroadcaster) prepareSIRIStopMonitoringNotify() {
	smb.connector.mutex.Lock()

	events := smb.connector.toBroadcast
	smb.connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	smb.connector.mutex.Unlock()

	for key, stopVisits := range events {
		sub, ok := smb.connector.Partner().Subscriptions().Find(key)
		if !ok {
			continue
		}

		// Initialize builder
		stopMonitoringBuilder := NewBroadcastStopMonitoringBuilder(smb.connector.Partner(), SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		stopMonitoringBuilder.StopVisitTypes = sub.SubscriptionOption("StopVisitTypes")

		// maximumStopVisits, _ := strconv.Atoi(sub.SubscriptionOption("MaximumStopVisits"))
		monitoredStopVisits := make(map[model.StopVisitId]struct{}) //Making sure not to send 2 times the same SV

		notification := &siri.SIRINotifyStopMonitoring{
			Address:                   smb.connector.Partner().Address(),
			ProducerRef:               smb.connector.Partner().ProducerRef(),
			RequestMessageRef:         sub.SubscriptionOption("MessageIdentifier"),
			ResponseMessageIdentifier: smb.connector.Partner().NewResponseMessageIdentifier(),
			ResponseTimestamp:         smb.connector.Clock().Now(),
		}
		deliveries := make(map[string]*siri.SIRINotifyStopMonitoringDelivery)

		for _, stopVisitId := range stopVisits {
			// Check if resource is already in the map
			if _, ok := monitoredStopVisits[stopVisitId]; ok {
				continue
			}

			// Find the StopVisit
			stopVisit, ok := smb.connector.Partner().Model().StopVisits().Find(stopVisitId)
			if !ok {
				continue
			}

			// Find the Resource
			monitoringRef, resource, ok := smb.findResource(stopVisit.StopAreaId, sub)
			if !ok {
				continue
			}

			// Get the delivery
			delivery := smb.getDelivery(deliveries, sub, monitoringRef)

			// Get the monitoredStopVisit
			stopMonitoringBuilder.MonitoringRef = monitoringRef
			if !smb.handledStopVisitAppend(stopVisit, delivery, stopMonitoringBuilder) {
				continue
			}

			monitoredStopVisits[stopVisitId] = struct{}{}

			// Get the Resource lastState for the StopVisit
			lastStateInterface, _ := resource.LastState(string(stopVisitId))
			lastState, ok := lastStateInterface.(*stopMonitoringLastChange)
			if !ok {
				continue
			}
			lastState.UpdateState(stopVisit)
		}

		for _, delivery := range deliveries {
			if len(delivery.MonitoredStopVisits) != 0 || len(delivery.CancelledStopVisits) != 0 {
				notification.Deliveries = append(notification.Deliveries, delivery)
			}
		}
		if len(notification.Deliveries) != 0 {
			smb.sendNotification(notification)
		}
	}
}

func (smb *SMBroadcaster) getDelivery(deliveries map[string]*siri.SIRINotifyStopMonitoringDelivery, sub *Subscription, monitoringRef string) (delivery *siri.SIRINotifyStopMonitoringDelivery) {
	delivery, ok := deliveries[monitoringRef]
	if !ok {
		delivery = &siri.SIRINotifyStopMonitoringDelivery{
			MonitoringRef:          monitoringRef,
			RequestMessageRef:      sub.SubscriptionOption("MessageIdentifier"),
			ResponseTimestamp:      smb.connector.Clock().Now(),
			SubscriberRef:          sub.SubscriberRef,
			SubscriptionIdentifier: sub.ExternalId(),
			Status:                 true,
		}
		deliveries[monitoringRef] = delivery
	}
	return
}

func (smb *SMBroadcaster) findResource(stopAreaId model.StopAreaId, sub *Subscription) (string, *SubscribedResource, bool) {
	for _, objectid := range smb.connector.Partner().Model().StopAreas().FindAscendantsWithObjectIdKind(stopAreaId, smb.connector.remoteObjectidKind) {
		resource := sub.Resource(objectid)
		if resource != nil {
			return objectid.Value(), resource, true
		}
	}
	return "", nil, false
}

func (smb *SMBroadcaster) handledStopVisitAppend(stopVisit *model.StopVisit, delivery *siri.SIRINotifyStopMonitoringDelivery, stopMonitoringBuilder *BroadcastStopMonitoringBuilder) bool {

	if stopVisit.ArrivalStatus == model.STOP_VISIT_ARRIVAL_CANCELLED || stopVisit.DepartureStatus == model.STOP_VISIT_DEPARTURE_CANCELLED || stopVisit.DepartureStatus == model.STOP_VISIT_DEPARTURE_DEPARTED {
		return smb.handleCancelledStopVisit(stopVisit, delivery, stopMonitoringBuilder)
	} else {
		return smb.handleMonitoredStopVisit(stopVisit, delivery, stopMonitoringBuilder)
	}
}

func (smb *SMBroadcaster) handleCancelledStopVisit(stopVisit *model.StopVisit, delivery *siri.SIRINotifyStopMonitoringDelivery, stopMonitoringBuilder *BroadcastStopMonitoringBuilder) bool {
	cancelledStopVisit := stopMonitoringBuilder.BuildCancelledStopVisit(stopVisit)
	if cancelledStopVisit == nil {
		return false
	}

	delivery.CancelledStopVisits = append(delivery.CancelledStopVisits, cancelledStopVisit)
	return true
}

func (smb *SMBroadcaster) handleMonitoredStopVisit(stopVisit *model.StopVisit, delivery *siri.SIRINotifyStopMonitoringDelivery, stopMonitoringBuilder *BroadcastStopMonitoringBuilder) bool {
	monitoredStopVisit := stopMonitoringBuilder.BuildMonitoredStopVisit(stopVisit)
	if monitoredStopVisit == nil {
		return false
	}
	delivery.MonitoredStopVisits = append(delivery.MonitoredStopVisits, monitoredStopVisit)
	return true
}

func (smb *SMBroadcaster) sendNotification(notify *siri.SIRINotifyStopMonitoring) {
	message := smb.newBQEvent()

	logSIRIStopMonitoringNotify(message, notify)

	t := smb.Clock().Now()

	err := smb.connector.Partner().SIRIClient().NotifyStopMonitoring(notify)
	message.ProcessingTime = smb.Clock().Since(t).Seconds()
	if err != nil {
		logger.Log.Debugf("Error in StopMonitoringBroadcaster while attempting to send a notification: %v", err)
	}

	audit.CurrentBigQuery(string(smb.connector.Partner().Referential().Slug())).WriteEvent(message)
}

func (smb *SMBroadcaster) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "NotifyStopMonitoring",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(smb.connector.partner.Slug()),
		Status:    "OK",
	}
}

func logSIRIStopMonitoringNotify(message *audit.BigQueryMessage, notification *siri.SIRINotifyStopMonitoring) {
	monitoringRefs := []string{}
	cancelledMonitoringRefs := []string{}

	for _, delivery := range notification.Deliveries {
		for _, sv := range delivery.MonitoredStopVisits {
			monitoringRefs = append(monitoringRefs, sv.MonitoringRef)
		}
		for _, sv := range delivery.CancelledStopVisits {
			cancelledMonitoringRefs = append(cancelledMonitoringRefs, sv.MonitoringRef)
		}
	}

	message.RequestIdentifier = notification.RequestMessageRef
	message.ResponseIdentifier = notification.ResponseMessageIdentifier
	message.StopAreas = append(monitoringRefs, cancelledMonitoringRefs...)

	delivery := notification.Deliveries[0]

	message.SubscriptionIdentifiers = []string{delivery.SubscriptionIdentifier}

	if !delivery.Status {
		message.Status = "Error"
		message.ErrorDetails = delivery.ErrorString()
	}

	xml, err := notification.BuildXML()
	if err != nil {
		return
	}
	message.ResponseRawMessage = xml
	message.ResponseSize = int64(len(xml))
}
