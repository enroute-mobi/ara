package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
	"golang.org/x/exp/maps"
)

type VehicleMonitoringSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestVehicleUpdate(request *VehicleUpdateRequest)
	HandleNotifyVehicleMonitoring(delivery *sxml.XMLNotifyVehicleMonitoring) *CollectedRefs
}

type SIRIVehicleMonitoringSubscriptionCollector struct {
	connector

	deletedSubscriptions        *DeletedSubscriptions
	vehicleMonitoringSubscriber SIRIVehicleMonitoringSubscriber
	updateSubscriber            UpdateSubscriber
}

type SIRIVehicleMonitoringSubscriptionCollectorFactory struct{}

func (factory *SIRIVehicleMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIVehicleMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIVehicleMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIVehicleMonitoringSubscriptionCollector(partner *Partner) *SIRIVehicleMonitoringSubscriptionCollector {
	connector := &SIRIVehicleMonitoringSubscriptionCollector{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace()
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent
	connector.vehicleMonitoringSubscriber = NewSIRIVehicleMonitoringSubscriber(connector)

	return connector
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) Stop() {
	connector.vehicleMonitoringSubscriber.Stop()
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) Start() {
	connector.deletedSubscriptions = NewDeletedSubscriptions()
	connector.vehicleMonitoringSubscriber.Start()
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) RequestVehicleUpdate(request *VehicleUpdateRequest) {
	line, ok := connector.partner.Model().Lines().Find(request.LineId())
	if !ok {
		logger.Log.Debugf("VehicleUpdateRequest in VehicleMonitoring SubscriptionCollector for unknown line %v", request.LineId())
		return
	}

	lineCode, ok := line.Code(connector.remoteCodeSpace)
	if !ok {
		logger.Log.Debugf("Requested line %v doesn't have a code with codeSpace %v", request.LineId(), connector.remoteCodeSpace)
		return
	}

	// Try to find a Subscription with the resource
	subscriptions := connector.partner.Subscriptions().FindByResourceId(lineCode.String(), "VehicleMonitoringCollect")
	if len(subscriptions) > 0 {
		for _, subscription := range subscriptions {
			resource := subscription.Resource(lineCode)
			if resource == nil { // Should never happen
				logger.Log.Debugf("Can't find resource in subscription after Subscriptions#FindByResourceId")
				return
			}
			if !resource.SubscribedAt().IsZero() {
				resource.SubscribedUntil = connector.Clock().Now().Add(2 * time.Minute)
			}
		}
		return
	}

	// Else we find or create a subscription to add the resource
	newSubscription := connector.partner.Subscriptions().FindOrCreateByKind("VehicleMonitoringCollect")
	ref := model.Reference{
		Code: &lineCode,
		Type:     "Line",
	}

	newSubscription.CreateAndAddNewResource(ref)
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) SetVehicleMonitoringSubscriber(vehicleMonitoringSubscriber SIRIVehicleMonitoringSubscriber) {
	connector.vehicleMonitoringSubscriber = vehicleMonitoringSubscriber
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) HandleNotifyVehicleMonitoring(notify *sxml.XMLNotifyVehicleMonitoring) (collectedRefs *CollectedRefs) {
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})
	var updateEvents VehicleMonitoringUpdateEvents

	collectedRefs = NewCollectedRefs()
	for _, delivery := range notify.VehicleMonitoringDeliveries() {
		subscriptionId := delivery.SubscriptionRef()

		if subscriptionId == "" {
			logger.Log.Debugf("Partner %s sent a NotifyVehicleMonitoring with an empty SubscriptionRef\n", connector.Partner().Slug())
			continue
		}

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Debugf("Partner %s sent a NotifyVehicleMonitoring to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			if !connector.deletedSubscriptions.AlreadySend(subscriptionId) {
				subToDelete[subscriptionId] = struct{}{}
			}
			continue
		}
		if subscription.Kind() != VehicleMonitoringCollect {
			logger.Log.Debugf("Partner %s sent a NotifyVehicleMonitoring to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind VehicleMonitoringCollect"
			continue
		}

		builder := NewVehicleMonitoringUpdateEventBuilder(connector.partner)
		builder.SetUpdateEvents(delivery.VehicleActivities())

		updateEvents = builder.UpdateEvents()

		maps.Copy(collectedRefs.LineRefs, updateEvents.LineRefs)
		maps.Copy(collectedRefs.MonitoringRefs, updateEvents.MonitoringRefs)
		maps.Copy(collectedRefs.VehicleJourneyRefs, updateEvents.VehicleJourneyRefs)
		maps.Copy(collectedRefs.VehicleRefs, updateEvents.VehicleRefs)

		connector.broadcastUpdateEvents(&updateEvents)
	}

	for subId := range subToDelete {
		CancelSubscription(subId, "VehicleMonitoringSubscriptionCollector", connector)
	}

	return collectedRefs
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) broadcastUpdateEvents(events *VehicleMonitoringUpdateEvents) {
	if connector.updateSubscriber == nil {
		return
	}
	for _, e := range events.StopAreas {
		connector.updateSubscriber(e)
	}
	for _, e := range events.Lines {
		connector.updateSubscriber(e)
	}
	for _, e := range events.VehicleJourneys {
		connector.updateSubscriber(e)
	}
	for _, e := range events.Vehicles {
		connector.updateSubscriber(e)
	}
}
