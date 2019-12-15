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

type SIRIStopMonitoringBroadcaster interface {
	model.Stopable
	model.Startable
}

type SMBroadcaster struct {
	model.ClockConsumer

	connector *SIRIStopMonitoringSubscriptionBroadcaster
}

type StopMonitoringBroadcaster struct {
	SMBroadcaster

	stop chan struct{}
}

type FakeStopMonitoringBroadcaster struct {
	SMBroadcaster

	model.ClockConsumer
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
				ResponseMessageIdentifier: smb.connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
				ResponseTimestamp:         smb.connector.Clock().Now(),
			}

			delivery := &siri.SIRINotifyStopMonitoringDelivery{
				SubscriberRef:          smb.connector.SIRIPartner().SubscriberRef(),
				SubscriptionIdentifier: sub.ExternalId(),
				RequestMessageRef:      sub.SubscriptionOption("MessageIdentifier"),
				ResponseTimestamp:      smb.connector.Clock().Now(),
				Status:                 false,
				ErrorType:              "OtherError",
				ErrorNumber:            1,
				ErrorText:              fmt.Sprintf("Erreur [PRODUCER_UNAVAILABLE] : %v indisponible", producer),
			}

			notification.Deliveries = []*siri.SIRINotifyStopMonitoringDelivery{delivery}

			logStashEvent := smb.newLogStashEvent()
			logSIRINotMonitoredNotify(logStashEvent, notification)
			audit.CurrentLogStash().WriteEvent(logStashEvent)

			smb.sendNotification(notification)
		}
	}
}

func (smb *SMBroadcaster) prepareSIRIStopMonitoringNotify() {
	smb.connector.mutex.Lock()

	events := smb.connector.toBroadcast
	smb.connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	smb.connector.mutex.Unlock()

	tx := smb.connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	for key, stopVisits := range events {
		sub, ok := smb.connector.Partner().Subscriptions().Find(key)
		if !ok {
			continue
		}

		// Initialize builder
		stopMonitoringBuilder := NewBroadcastStopMonitoringBuilder(tx, smb.connector.SIRIPartner(), SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		stopMonitoringBuilder.StopVisitTypes = sub.SubscriptionOption("StopVisitTypes")

		// maximumStopVisits, _ := strconv.Atoi(sub.SubscriptionOption("MaximumStopVisits"))
		monitoredStopVisits := make(map[model.StopVisitId]struct{}) //Making sure not to send 2 times the same SV

		notification := &siri.SIRINotifyStopMonitoring{
			Address:                   smb.connector.Partner().Address(),
			ProducerRef:               smb.connector.Partner().ProducerRef(),
			RequestMessageRef:         sub.SubscriptionOption("MessageIdentifier"),
			ResponseMessageIdentifier: smb.connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
			ResponseTimestamp:         smb.connector.Clock().Now(),
		}
		deliveries := make(map[string]*siri.SIRINotifyStopMonitoringDelivery)

		for _, stopVisitId := range stopVisits {
			// Check if resource is already in the map
			if _, ok := monitoredStopVisits[stopVisitId]; ok {
				continue
			}

			// Find the StopVisit
			stopVisit, ok := tx.Model().StopVisits().Find(stopVisitId)
			if !ok {
				continue
			}

			// Find the Resource
			monitoringRef, resource, ok := smb.findResource(stopVisit.StopAreaId, sub, tx)
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

			// See what to do about the MaximumStopVisits #10333
			// // Refresh delivery
			// if maximumStopVisits != 0 && (len(delivery.MonitoredStopVisits)+len(delivery.CancelledStopVisits)) >= maximumStopVisits {
			// 	smb.sendNotification(delivery)
			// 	delivery.MonitoredStopVisits = []*siri.SIRIMonitoredStopVisit{}
			// 	delivery.CancelledStopVisits = []*siri.SIRICancelledStopVisit{}
			// }

			// Get the Resource lastState for the StopVisit
			lastStateInterface, _ := resource.LastState(string(stopVisitId))
			lastState, ok := lastStateInterface.(*stopMonitoringLastChange)
			if !ok {
				continue
			}
			lastState.UpdateState(&stopVisit)
		}

		for _, delivery := range deliveries {
			if len(delivery.MonitoredStopVisits) != 0 || len(delivery.CancelledStopVisits) != 0 {
				notification.Deliveries = append(notification.Deliveries, delivery)
			}
		}
		if len(notification.Deliveries) != 0 {
			logStashEvent := smb.newLogStashEvent()
			logSIRIStopMonitoringNotify(logStashEvent, notification)
			audit.CurrentLogStash().WriteEvent(logStashEvent)

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
			SubscriberRef:          smb.connector.SIRIPartner().SubscriberRef(),
			SubscriptionIdentifier: sub.ExternalId(),
			Status:                 true,
		}
		deliveries[monitoringRef] = delivery
	}
	return
}

func (smb *SMBroadcaster) findResource(stopAreaId model.StopAreaId, sub *Subscription, tx *model.Transaction) (string, *SubscribedResource, bool) {
	for _, objectid := range tx.Model().StopAreas().FindAscendantsWithObjectIdKind(stopAreaId, smb.connector.Partner().RemoteObjectIDKind(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)) {
		resource := sub.Resource(objectid)
		if resource != nil {
			return objectid.Value(), resource, true
		}
	}
	return "", nil, false
}

func (smb *SMBroadcaster) handledStopVisitAppend(stopVisit model.StopVisit, delivery *siri.SIRINotifyStopMonitoringDelivery, stopMonitoringBuilder *BroadcastStopMonitoringBuilder) bool {

	if stopVisit.ArrivalStatus == model.STOP_VISIT_ARRIVAL_CANCELLED || stopVisit.DepartureStatus == model.STOP_VISIT_DEPARTURE_CANCELLED || stopVisit.DepartureStatus == model.STOP_VISIT_DEPARTURE_DEPARTED {
		return smb.handleCancelledStopVisit(stopVisit, delivery, stopMonitoringBuilder)
	} else {
		return smb.handleMonitoredStopVisit(stopVisit, delivery, stopMonitoringBuilder)
	}
}

func (smb *SMBroadcaster) handleCancelledStopVisit(stopVisit model.StopVisit, delivery *siri.SIRINotifyStopMonitoringDelivery, stopMonitoringBuilder *BroadcastStopMonitoringBuilder) bool {
	cancelledStopVisit := stopMonitoringBuilder.BuildCancelledStopVisit(stopVisit)
	if cancelledStopVisit == nil {
		return false
	}

	delivery.CancelledStopVisits = append(delivery.CancelledStopVisits, cancelledStopVisit)
	return true
}

func (smb *SMBroadcaster) handleMonitoredStopVisit(stopVisit model.StopVisit, delivery *siri.SIRINotifyStopMonitoringDelivery, stopMonitoringBuilder *BroadcastStopMonitoringBuilder) bool {
	monitoredStopVisit := stopMonitoringBuilder.BuildMonitoredStopVisit(stopVisit)
	if monitoredStopVisit == nil {
		return false
	}
	delivery.MonitoredStopVisits = append(delivery.MonitoredStopVisits, monitoredStopVisit)
	return true
}

func (smb *SMBroadcaster) sendNotification(notify *siri.SIRINotifyStopMonitoring) {
	err := smb.connector.SIRIPartner().SOAPClient().NotifyStopMonitoring(notify)
	if err != nil {
		logger.Log.Debugf("Error in StopMonitoringBroadcaster while attempting to send a notification: %v", err)
		event := smb.newLogStashEvent()
		logSIRINotifyError(err.Error(), notify.ResponseMessageIdentifier, event)
		audit.CurrentLogStash().WriteEvent(event)
	}
}

func (smb *SMBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := smb.connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringSubscriptionBroadcaster"
	return event
}

func logSIRIStopMonitoringNotify(logStashEvent audit.LogStashEvent, notification *siri.SIRINotifyStopMonitoring) {
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

	logStashEvent["siriType"] = "NotifyStopMonitoring"
	logStashEvent["producerRef"] = notification.ProducerRef
	logStashEvent["requestMessageRef"] = notification.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = notification.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = notification.ResponseTimestamp.String()
	logStashEvent["subscriberRef"] = notification.Deliveries[0].SubscriberRef
	logStashEvent["subscriptionIdentifier"] = notification.Deliveries[0].SubscriptionIdentifier
	logStashEvent["monitoringRefs"] = strings.Join(monitoringRefs, ",")
	logStashEvent["cancelledMonitoringRefs"] = strings.Join(cancelledMonitoringRefs, ",")
	logStashEvent["status"] = "true"

	xml, err := notification.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}

func logSIRINotMonitoredNotify(logStashEvent audit.LogStashEvent, notification *siri.SIRINotifyStopMonitoring) {
	logStashEvent["siriType"] = "NotifyStopMonitoring"
	logStashEvent["producerRef"] = notification.ProducerRef
	logStashEvent["requestMessageRef"] = notification.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = notification.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = notification.ResponseTimestamp.String()

	delivery := notification.Deliveries[0]
	logStashEvent["subscriberRef"] = delivery.SubscriberRef
	logStashEvent["subscriptionIdentifier"] = delivery.SubscriptionIdentifier
	logStashEvent["status"] = strconv.FormatBool(delivery.Status)
	logStashEvent["errorType"] = delivery.ErrorType
	logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber)
	logStashEvent["errorText"] = delivery.ErrorText

	xml, err := notification.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
