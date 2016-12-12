package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type StopVisitUpdateSubscriber func(model.StopVisitUpdateEvent)

type CollectManagerInterface interface {
	UpdateStopArea(request *StopAreaUpdateRequest)
}

type CollectManager struct {
	partners              Partners
	stopVisitUpdateEvents chan model.StopVisitUpdateEvent
}

type TestCollectManager struct {
	Done   chan bool
	Events []*model.StopAreaUpdateEvent
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

func NewCollectManager(partners Partners) CollectManagerInterface {
	return &CollectManager{
		partners:              partners,
		stopVisitUpdateEvents: make(chan model.StopVisitUpdateEvent),
	}
}

func (manager *CollectManager) HandleStopVisitUpdateEvent(stopVisitUpdateSubscriber StopVisitUpdateSubscriber) {
	for stopVisitUpdateEvent := range manager.stopVisitUpdateEvents {
		stopVisitUpdateSubscriber(stopVisitUpdateEvent)
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
		manager.stopVisitUpdateEvents <- *stopVisitUpdateEvent
	}
}

func (manager *CollectManager) bestPartner(request *StopAreaUpdateRequest) *Partner {
	var testPartner *Partner
	for _, partner := range manager.partners.FindAll() {
		if partner.isConnectorDefined(SIRI_STOP_MONITORING_REQUEST_COLLECTOR) {
			return partner
		} else if partner.isConnectorDefined(TEST_STOP_MONITORING_REQUEST_COLLECTOR) {
			testPartner = partner
		}
	}
	if testPartner != nil {
		return testPartner
	}
	return nil
}

func (manager *CollectManager) requestStopAreaUpdate(partner *Partner, request *StopAreaUpdateRequest) (*model.StopAreaUpdateEvent, error) {
	logger.Log.Debugf("RequestStopAreaUpdate %#v", request)

	event, err := partner.StopMonitoringRequestCollector().RequestStopAreaUpdate(request)
	if err != nil {
		return nil, err
	}
	return event, nil
}
