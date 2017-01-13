package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type StopVisitUpdateSubscriber func(*model.StopVisitUpdateEvent)

type CollectManagerInterface interface {
	UpdateStopArea(request *StopAreaUpdateRequest)
	HandleStopVisitUpdateEvent(StopVisitUpdateSubscriber)
}

type CollectManager struct {
	partners                   Partners
	stopVisitUpdateSubscribers []StopVisitUpdateSubscriber
}

// TestCollectManager has a test StopVisitUpdateSubscriber method
type TestCollectManager struct {
	Done            chan bool
	Events          []*model.StopAreaUpdateEvent
	StopVisitEvents []*model.StopVisitUpdateEvent
}

func NewTestCollectManager() CollectManagerInterface {
	return &TestCollectManager{
		Done: make(chan bool, 1),
	}
}

func (manager *TestCollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	event := &model.StopAreaUpdateEvent{}
	manager.Events = append(manager.Events, event)

	manager.Done <- true
}

func (manager *TestCollectManager) TestStopVisitUpdateSubscriber(event *model.StopVisitUpdateEvent) {
	manager.StopVisitEvents = append(manager.StopVisitEvents, event)
}

func (manager *TestCollectManager) HandleStopVisitUpdateEvent(StopVisitUpdateSubscriber) {}

func NewCollectManager(partners Partners) CollectManagerInterface {
	return &CollectManager{
		partners:                   partners,
		stopVisitUpdateSubscribers: make([]StopVisitUpdateSubscriber, 0),
	}
}

func (manager *CollectManager) HandleStopVisitUpdateEvent(stopVisitUpdateSubscriber StopVisitUpdateSubscriber) {
	manager.stopVisitUpdateSubscribers = append(manager.stopVisitUpdateSubscribers, stopVisitUpdateSubscriber)
}

func (manager *CollectManager) broadcastStopVisitUpdateEvent(event *model.StopVisitUpdateEvent) {
	for _, stopVisitUpdateSubscriber := range manager.stopVisitUpdateSubscribers {
		stopVisitUpdateSubscriber(event)
	}
}

func (manager *CollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	partner := manager.bestPartner(request)
	if partner == nil {
		logger.Log.Debugf("Can't find a partner for StopArea %v", request.StopAreaId())
		return
	}

	event, err := manager.requestStopAreaUpdate(partner, request)
	if err != nil {
		logger.Log.Printf("Can't request stop area update : %v", err)
		return
	}

	for _, stopVisitUpdateEvent := range event.StopVisitUpdateEvents {
		manager.broadcastStopVisitUpdateEvent(stopVisitUpdateEvent)
	}
}

func (manager *CollectManager) bestPartner(request *StopAreaUpdateRequest) *Partner {
	var testPartner *Partner
	for _, partner := range manager.partners.FindAll() {
		_, ok := partner.Connector(SIRI_STOP_MONITORING_REQUEST_COLLECTOR)
		if ok && partner.OperationnalStatus() == OPERATIONNAL_STATUS_UP {
			return partner
		}
		_, ok = partner.Connector(TEST_STOP_MONITORING_REQUEST_COLLECTOR)
		if ok {
			testPartner = partner
		}
	}
	return testPartner
}

func (manager *CollectManager) requestStopAreaUpdate(partner *Partner, request *StopAreaUpdateRequest) (*model.StopAreaUpdateEvent, error) {
	logger.Log.Debugf("RequestStopAreaUpdate %v", request.StopAreaId())

	event, err := partner.StopMonitoringRequestCollector().RequestStopAreaUpdate(request)
	if err != nil {
		return nil, err
	}
	return event, nil
}
