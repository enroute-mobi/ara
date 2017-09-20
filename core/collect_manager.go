package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type StopAreaUpdateSubscriber func(*model.StopAreaUpdateEvent)
type SituationUpdateSubscriber func([]*model.SituationUpdateEvent)

type CollectManagerInterface interface {
	UpdateStopArea(request *StopAreaUpdateRequest)
	HandleStopAreaUpdateEvent(StopAreaUpdateSubscriber)
	BroadcastStopAreaUpdateEvent(event *model.StopAreaUpdateEvent)

	UpdateSituation(request *SituationUpdateRequest)
	HandleSituationUpdateEvent(SituationUpdateSubscriber)
	BroadcastSituationUpdateEvent(event []*model.SituationUpdateEvent)
}

type CollectManager struct {
	StopAreaUpdateSubscribers  []StopAreaUpdateSubscriber
	SituationUpdateSubscribers []SituationUpdateSubscriber
	referential                *Referential
}

// TestCollectManager has a test StopAreaUpdateSubscriber method
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

func (manager *TestCollectManager) TestStopAreaUpdateSubscriber(event *model.StopAreaUpdateEvent) {
	for _, stopVisitUpdateEvent := range event.StopVisitUpdateEvents {
		manager.StopVisitEvents = append(manager.StopVisitEvents, stopVisitUpdateEvent)
	}
}

func (manager *TestCollectManager) HandleStopAreaUpdateEvent(StopAreaUpdateSubscriber) {}
func (manager *TestCollectManager) BroadcastStopAreaUpdateEvent(event *model.StopAreaUpdateEvent) {
	manager.Events = append(manager.Events, event)
}

func (manager *TestCollectManager) UpdateSituation(*SituationUpdateRequest)              {}
func (manager *TestCollectManager) HandleSituationUpdateEvent(SituationUpdateSubscriber) {}
func (manager *TestCollectManager) BroadcastSituationUpdateEvent(event []*model.SituationUpdateEvent) {
}

// TEST END

func NewCollectManager(referential *Referential) CollectManagerInterface {
	return &CollectManager{
		referential:                referential,
		StopAreaUpdateSubscribers:  make([]StopAreaUpdateSubscriber, 0),
		SituationUpdateSubscribers: make([]SituationUpdateSubscriber, 0),
	}
}

func (manager *CollectManager) HandleStopAreaUpdateEvent(StopAreaUpdateSubscriber StopAreaUpdateSubscriber) {
	manager.StopAreaUpdateSubscribers = append(manager.StopAreaUpdateSubscribers, StopAreaUpdateSubscriber)
}

func (manager *CollectManager) BroadcastStopAreaUpdateEvent(event *model.StopAreaUpdateEvent) {
	for _, StopAreaUpdateSubscriber := range manager.StopAreaUpdateSubscribers {
		StopAreaUpdateSubscriber(event)
	}
}

func (manager *CollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	partner := manager.bestPartner(request)
	if partner == nil {
		logger.Log.Debugf("Can't find a partner for StopArea %v", request.StopAreaId())
		return
	}
	manager.requestStopAreaUpdate(partner, request)
}

func (manager *CollectManager) bestPartner(request *StopAreaUpdateRequest) *Partner {

	stopArea, ok := manager.referential.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		return nil
	}

	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
			continue
		}
		_, connectorPresent := partner.Connector(SIRI_STOP_MONITORING_REQUEST_COLLECTOR)
		_, testConnectorPresent := partner.Connector(TEST_STOP_MONITORING_REQUEST_COLLECTOR)
		_, subscriptionPresent := partner.Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_COLLECTOR)

		if !(connectorPresent || testConnectorPresent || subscriptionPresent) {
			continue
		}

		partnerKind := partner.Setting("remote_objectid_kind")

		stopAreaObjectID, ok := stopArea.ObjectID(partnerKind)
		if !ok {
			continue
		}

		lineIds := make(map[string]struct{})
		for _, lineId := range stopArea.LineIds {
			line, ok := manager.referential.Model().Lines().Find(lineId)
			if !ok {
				continue
			}
			lineObjectID, ok := line.ObjectID(partnerKind)
			if !ok {
				continue
			}
			lineIds[lineObjectID.Value()] = struct{}{}
		}

		if partner.CanCollect(stopAreaObjectID, lineIds) {
			return partner
		}
	}
	return nil
}

func (manager *CollectManager) requestStopAreaUpdate(partner *Partner, request *StopAreaUpdateRequest) {
	logger.Log.Debugf("RequestStopAreaUpdate %v", request.StopAreaId())

	if collect := partner.StopMonitoringSubscriptionCollector(); collect != nil {
		collect.RequestStopAreaUpdate(request)
		return
	}
	partner.StopMonitoringRequestCollector().RequestStopAreaUpdate(request)
}

func (manager *CollectManager) BroadcastSituationUpdateEvent(event []*model.SituationUpdateEvent) {
	for _, SituationUpdateSubscriber := range manager.SituationUpdateSubscribers {
		SituationUpdateSubscriber(event)
	}
}

func (manager *CollectManager) requestSituationUpdate(partner *Partner, request *SituationUpdateRequest) {
	logger.Log.Debugf("RequestSituationUpdate %v", request.Id())

	if collect := partner.GeneralMessageSubscriptionCollector(); collect != nil {
		collect.RequestSituationUpdate(request)
		return
	}
	partner.GeneralMessageRequestCollector().RequestSituationUpdate(request)
}

func (manager *CollectManager) HandleSituationUpdateEvent(SituationUpdateSubscriber SituationUpdateSubscriber) {
	manager.SituationUpdateSubscribers = append(manager.SituationUpdateSubscribers, SituationUpdateSubscriber)
}

func (manager *CollectManager) UpdateSituation(request *SituationUpdateRequest) {
	for _, partner := range manager.referential.Partners().FindAll() {
		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
			continue
		}
		if connector := partner.GeneralMessageSubscriptionCollector(); connector != nil {
			logger.Log.Debugf("RequestSituationUpdate %v", request.Id())
			connector.RequestSituationUpdate(request)
			continue
		}
		if connector := partner.GeneralMessageRequestCollector(); connector != nil {
			logger.Log.Debugf("RequestSituationUpdate %v", request.Id())
			connector.RequestSituationUpdate(request)
			continue
		}
	}
}
