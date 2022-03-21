package core

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type VehicleMonitoringSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestVehicleUpdate(request *VehicleUpdateRequest)
	// HandleNotifyVehicleMonitoring(delivery *siri.XMLNotifyVehicleMonitoring)
}

type SIRIVehicleMonitoringSubscriptionCollector struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	// vehicleMonitoringSubscriber SIRIVehicleMonitoringSubscriber
	// updateSubscriber            UpdateSubscriber
}

type SIRIVehicleMonitoringSubscriptionCollectorFactory struct{}

func (factory *SIRIVehicleMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIVehicleMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIVehicleMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIVehicleMonitoringSubscriptionCollector(partner *Partner) *SIRIVehicleMonitoringSubscriptionCollector {
	connector := &SIRIVehicleMonitoringSubscriptionCollector{}
	connector.partner = partner
	// manager := partner.Referential().CollectManager()
	// connector.updateSubscriber = manager.BroadcastUpdateEvent
	// connector.vehicleMonitoringSubscriber = NewSIRIVehicleMonitoringSubscriber(connector)

	return connector
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) Stop() {
	// connector.vehicleMonitoringSubscriber.Stop()
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) Start() {
	// connector.vehicleMonitoringSubscriber.Start()
}

func (connector *SIRIVehicleMonitoringSubscriptionCollector) RequestVehicleUpdate(request *VehicleUpdateRequest) {
}

// func (connector *SIRIVehicleMonitoringSubscriptionCollector) SetVehicleMonitoringSubscriber(vehicleMonitoringSubscriber SIRIVehicleMonitoringSubscriber) {
// 	connector.vehicleMonitoringSubscriber = vehicleMonitoringSubscriber
// }

// func (connector *SIRIVehicleMonitoringSubscriptionCollector) HandleNotifyVehicleMonitoring(notify *siri.XMLNotifyVehicleMonitoring) {
// }

// func (connector *SIRIVehicleMonitoringSubscriptionCollector) broadcastUpdateEvents(events *VehicleMonitoringUpdateEvents) {
// 	if connector.updateSubscriber == nil {
// 		return
// 	}
// 	for _, e := range events.Vehicles {
// 		connector.updateSubscriber(e)
// 	}
// 	for _, e := range events.StopAreas {
// 		connector.updateSubscriber(e)
// 	}
// 	for _, e := range events.Lines {
// 		connector.updateSubscriber(e)
// 	}
// 	for _, e := range events.VehicleJourneys {
// 		connector.updateSubscriber(e)
// 	}
// 	for _, es := range events.StopVisits { // Stopvisits are map[MonitoringRef]map[ItemIdentifier]event
// 		for _, e := range es {
// 			connector.updateSubscriber(e)
// 		}
// 	}
// }

// func (connector *SIRIVehicleMonitoringSubscriptionCollector) newBQEvent() *audit.BigQueryMessage {
// 	return &audit.BigQueryMessage{
// 		Protocol:  "siri",
// 		Direction: "sent",
// 		Partner:   string(connector.partner.Slug()),
// 		Status:    "OK",
// 	}
// }
