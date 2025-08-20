package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/state"
)

type FacilityMonitoringBroadcaster interface {
	state.Stopable
	state.Startable
}

type FMBroadcaster struct {
	clock.ClockConsumer

	connector *SIRIFacilityMonitoringSubscriptionBroadcaster
}

type SIRIFacilityMonitoringBroadcaster struct {
	FMBroadcaster

	stop chan struct{}
}

type FakeSIRIFacilityMonitoringBroadcaster struct {
	FMBroadcaster
}

func NewFakeSIRIFacilityMonitoringBroadcaster(connector *SIRIFacilityMonitoringSubscriptionBroadcaster) FacilityMonitoringBroadcaster {
	broadcaster := &FakeSIRIFacilityMonitoringBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeSIRIFacilityMonitoringBroadcaster) Start() {
	broadcaster.prepareSIRIFacilityMonitoring()
}

func (broadcaster *FakeSIRIFacilityMonitoringBroadcaster) Stop() {}

func NewSIRIFacilityMonitoringBroadcaster(connector *SIRIFacilityMonitoringSubscriptionBroadcaster) FacilityMonitoringBroadcaster {
	broadcaster := &SIRIFacilityMonitoringBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (fm *SIRIFacilityMonitoringBroadcaster) Start() {
	logger.Log.Debugf("Start SIRIFacilityMonitoringBroadcaster")

	fm.stop = make(chan struct{})
	go fm.run()
}

func (fm *SIRIFacilityMonitoringBroadcaster) Stop() {
	if fm.stop != nil {
		close(fm.stop)
	}
}

func (fm *SIRIFacilityMonitoringBroadcaster) run() {
	c := fm.Clock().After(5 * time.Second)

	for {
		select {
		case <-fm.stop:
			logger.Log.Debugf("facility monitoring broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRIFacilityMonitoringBroadcaster visit")

			fm.prepareSIRIFacilityMonitoring()

			c = fm.Clock().After(5 * time.Second)
		}
	}
}

func (fm *FMBroadcaster) prepareSIRIFacilityMonitoring() {
	fm.connector.mutex.Lock()
	defer fm.connector.mutex.Unlock()

	events := fm.connector.toBroadcast
	fm.connector.toBroadcast = make(map[SubscriptionId][]model.FacilityId)

	for subId, facilityIds := range events {
		sub, ok := fm.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			logger.Log.Debugf("FM subscriptionBroadcast Could not find sub with id : %v", subId)
			continue
		}

		processedFacilities := make(map[model.FacilityId]struct{}) //Making sure not to send 2 times the same Facility
		delivery := &siri.SIRINotifyFacilityMonitoring{
			Address:                   fm.connector.Partner().Address(),
			ProducerRef:               fm.connector.Partner().ProducerRef(),
			ResponseMessageIdentifier: fm.connector.Partner().NewResponseMessageIdentifier(),
			SubscriberRef:             sub.SubscriberRef,
			SubscriptionIdentifier:    sub.ExternalId(),
			ResponseTimestamp:         fm.connector.Clock().Now(),
			Status:                    true,
		}

		for _, facilityId := range facilityIds {
			if _, ok := processedFacilities[facilityId]; ok {
				continue
			}

			// Find the Facility
			facility, ok := fm.connector.Partner().Model().Facilities().Find(facilityId)
			if !ok {
				continue
			}
			facilityCode, ok := facility.Code(fm.connector.remoteCodeSpace)
			if !ok {
				continue
			}

			// Find the Resource
			resource := sub.Resource(facilityCode)
			if resource == nil {
				continue
			}

			condition := &siri.SIRIFacilityCondition{
				FacilityRef:    facilityCode.Value(),
				FacilityStatus: string(facility.Status),
			}

			delivery.FacilityConditions = append(delivery.FacilityConditions, condition)

			lastStateInterface, ok := resource.LastState(string(facility.Id()))
			if !ok {
				resource.SetLastState(string(facility.Id()), ls.NewFacilityMonitoringLastChange(facility, sub))
			} else {
				lastStateInterface.(*ls.FacilityMonitoringLastChange).UpdateState(facility)
			}

			processedFacilities[facilityId] = struct{}{}
		}
		fm.sendDelivery(delivery)
	}

}

func (fm *FMBroadcaster) sendDelivery(delivery *siri.SIRINotifyFacilityMonitoring) {
	message := fm.newBQEvent()

	fm.logSIRIFacilityMonitoring(message, delivery)

	t := fm.Clock().Now()

	fm.connector.Partner().SIRIClient().NotifyFacilityMonitoring(delivery)
	message.ProcessingTime = fm.Clock().Since(t).Seconds()

	audit.CurrentBigQuery(string(fm.connector.Partner().Referential().Slug())).WriteEvent(message)
}

func (fm *FMBroadcaster) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.NOTIFY_FACILITY_MONITORING,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(fm.connector.Partner().Slug()),
		Status:    "OK",
	}
}

func (fm *FMBroadcaster) logSIRIFacilityMonitoring(message *audit.BigQueryMessage, response *siri.SIRINotifyFacilityMonitoring) {
	facilityRefs := make(map[string]struct{})

	for _, fc := range response.FacilityConditions {
		facilityRefs[fc.FacilityRef] = struct{}{}
	}

	message.ResponseIdentifier = response.ResponseMessageIdentifier
	message.SubscriptionIdentifiers = []string{response.SubscriptionIdentifier}

	if !response.Status {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	xml, err := response.BuildXML(fm.connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.ResponseRawMessage = xml
	message.ResponseSize = int64(len(xml))
}
