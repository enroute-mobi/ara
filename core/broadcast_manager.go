package core

import (
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/state"
)

type BroadcastManagerInterface interface {
	state.Startable
	state.Stopable

	GetStopMonitoringBroadcastEventChan() chan model.StopMonitoringBroadcastEvent
	GetGeneralMessageBroadcastEventChan() chan model.SituationBroadcastEvent
	GetVehicleBroadcastEventChan() chan model.VehicleBroadcastEvent
}

type BroadcastManager struct {
	Referential *Referential

	smbEventChan chan model.StopMonitoringBroadcastEvent
	gmbEventChan chan model.SituationBroadcastEvent
	vbEventChan  chan model.VehicleBroadcastEvent
	stop         chan struct{}
}

func NewBroadcastManager(referential *Referential) *BroadcastManager {
	return &BroadcastManager{
		Referential:  referential,
		smbEventChan: make(chan model.StopMonitoringBroadcastEvent, 2000),
		gmbEventChan: make(chan model.SituationBroadcastEvent, 2000),
		vbEventChan:  make(chan model.VehicleBroadcastEvent, 2000),
	}
}

func (manager *BroadcastManager) GetStopMonitoringBroadcastEventChan() chan model.StopMonitoringBroadcastEvent {
	return manager.smbEventChan
}

func (manager *BroadcastManager) GetGeneralMessageBroadcastEventChan() chan model.SituationBroadcastEvent {
	return manager.gmbEventChan
}

func (manager *BroadcastManager) GetVehicleBroadcastEventChan() chan model.VehicleBroadcastEvent {
	return manager.vbEventChan
}

func (manager *BroadcastManager) Start() {
	logger.Log.Debugf("BroadcastManager start")

	manager.stop = make(chan struct{})

	go manager.run()
}

func (manager *BroadcastManager) run() {
	for {
		select {
		case event := <-manager.smbEventChan:
			manager.smsbEvent_handler(event)
			manager.ettsbEvent_handler(event)
		case event := <-manager.gmbEventChan:
			manager.gmsbEvent_handler(event)
		case event := <-manager.vbEventChan:
			manager.vmEvent_handler(event)
		case <-manager.stop:
			logger.Log.Debugf("BroadcastManager Stop")
			return
		}
	}
}

func (manager *BroadcastManager) smsbEvent_handler(event model.StopMonitoringBroadcastEvent) {
	connectorTypes := []string{SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER, TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER}
	for _, partner := range manager.Referential.Partners().FindAllWithConnector(connectorTypes) {
		connector, ok := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleStopMonitoringBroadcastEvent(&event)
			continue
		}

		// TEST
		connector, ok = partner.Connector(TEST_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*TestStopMonitoringSubscriptionBroadcaster).HandleStopMonitoringBroadcastEvent(&event)
			continue
		}
	}
}

func (manager *BroadcastManager) ettsbEvent_handler(event model.StopMonitoringBroadcastEvent) {
	connectorTypes := []string{SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER, TEST_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER}
	for _, partner := range manager.Referential.Partners().FindAllWithConnector(connectorTypes) {
		connector, ok := partner.Connector(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*SIRIEstimatedTimetableSubscriptionBroadcaster).HandleBroadcastEvent(&event)
			continue
		}

		connector, ok = partner.Connector(TEST_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*TestETTSubscriptionBroadcaster).HandleBroadcastEvent(&event)
			continue
		}
	}
}

func (manager *BroadcastManager) vmEvent_handler(event model.VehicleBroadcastEvent) {
	connectorTypes := []string{SIRI_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER, TEST_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER}
	for _, partner := range manager.Referential.Partners().FindAllWithConnector(connectorTypes) {
		connector, ok := partner.Connector(SIRI_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*SIRIVehicleMonitoringSubscriptionBroadcaster).HandleBroadcastEvent(&event)
			continue
		}

		// TEST
		connector, ok = partner.Connector(TEST_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*TestVMSubscriptionBroadcaster).HandleBroadcastEvent(&event)
			continue
		}
	}
}

func (manager *BroadcastManager) gmsbEvent_handler(event model.SituationBroadcastEvent) {
	connectorTypes := []string{SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER, TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER}
	for _, partner := range manager.Referential.Partners().FindAllWithConnector(connectorTypes) {
		connector, ok := partner.Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*SIRIGeneralMessageSubscriptionBroadcaster).HandleGeneralMessageBroadcastEvent(&event)
			continue
		}

		// TEST
		connector, ok = partner.Connector(TEST_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		if ok {
			connector.(*TestGeneralMessageSubscriptionBroadcaster).HandleGeneralMessageBroadcastEvent(&event)
			continue
		}
	}
}

func (manager *BroadcastManager) Stop() {
	if manager.stop != nil {
		close(manager.stop)
	}
}
