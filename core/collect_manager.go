package core

import (
	"strconv"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type StopAreaUpdateSubscriber func(*model.LegacyStopAreaUpdateEvent)
type SituationUpdateSubscriber func([]*model.SituationUpdateEvent)
type UpdateSubscriber func(model.UpdateEvent)

type CollectManagerInterface interface {
	HandlePartnerStatusChange(partner string, status bool)
	UpdateStopArea(request *StopAreaUpdateRequest)

	// Legacy
	HandleLegacyStopAreaUpdateEvent(StopAreaUpdateSubscriber)
	BroadcastLegacyStopAreaUpdateEvent(event *model.LegacyStopAreaUpdateEvent)

	// New update events
	HandleUpdateEvent(UpdateSubscriber UpdateSubscriber)
	BroadcastUpdateEvent(event model.UpdateEvent)

	UpdateSituation(request *SituationUpdateRequest)
	HandleSituationUpdateEvent(SituationUpdateSubscriber)
	BroadcastSituationUpdateEvent(event []*model.SituationUpdateEvent)
}

type CollectManager struct {
	model.UUIDConsumer

	StopAreaUpdateSubscribers  []StopAreaUpdateSubscriber
	SituationUpdateSubscribers []SituationUpdateSubscriber
	UpdateSubscribers          []UpdateSubscriber
	referential                *Referential
}

// TestCollectManager has a test StopAreaUpdateSubscriber method
type TestCollectManager struct {
	Done            chan bool
	Events          []*model.LegacyStopAreaUpdateEvent
	StopVisitEvents []*model.LegacyStopVisitUpdateEvent
	UpdateEvents    []model.UpdateEvent
}

func NewTestCollectManager() CollectManagerInterface {
	return &TestCollectManager{
		Done: make(chan bool, 1),
	}
}

func (manager *TestCollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	event := &model.LegacyStopAreaUpdateEvent{}
	manager.Events = append(manager.Events, event)

	manager.Done <- true
}

func (manager *TestCollectManager) TestStopAreaUpdateSubscriber(event *model.LegacyStopAreaUpdateEvent) {
	manager.StopVisitEvents = append(manager.StopVisitEvents, event.LegacyStopVisitUpdateEvents...)
}

func (manager *TestCollectManager) HandlePartnerStatusChange(partner string, status bool) {}

// Legacy
func (manager *TestCollectManager) HandleLegacyStopAreaUpdateEvent(StopAreaUpdateSubscriber) {}
func (manager *TestCollectManager) BroadcastLegacyStopAreaUpdateEvent(event *model.LegacyStopAreaUpdateEvent) {
	manager.Events = append(manager.Events, event)
}

// New structure
func (manager *TestCollectManager) HandleUpdateEvent(UpdateSubscriber) {}
func (manager *TestCollectManager) BroadcastUpdateEvent(event model.UpdateEvent) {
	manager.UpdateEvents = append(manager.UpdateEvents, event)
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
		UpdateSubscribers:          make([]UpdateSubscriber, 0),
	}
}

func (manager *CollectManager) HandleLegacyStopAreaUpdateEvent(StopAreaUpdateSubscriber StopAreaUpdateSubscriber) {
	manager.StopAreaUpdateSubscribers = append(manager.StopAreaUpdateSubscribers, StopAreaUpdateSubscriber)
}

func (manager *CollectManager) BroadcastLegacyStopAreaUpdateEvent(event *model.LegacyStopAreaUpdateEvent) {
	for _, StopAreaUpdateSubscriber := range manager.StopAreaUpdateSubscribers {
		StopAreaUpdateSubscriber(event)
	}
}

func (manager *CollectManager) HandleUpdateEvent(UpdateSubscriber UpdateSubscriber) {
	manager.UpdateSubscribers = append(manager.UpdateSubscribers, UpdateSubscriber)
}

func (manager *CollectManager) BroadcastUpdateEvent(event model.UpdateEvent) {
	for _, UpdateSubscriber := range manager.UpdateSubscribers {
		UpdateSubscriber(event)
	}
}

func (manager *CollectManager) HandlePartnerStatusChange(partner string, status bool) {
	for _, stopAreaId := range manager.referential.Model().StopAreas().FindByOrigin(partner) {
		event := model.NewStopAreaMonitoredEvent(manager.NewUUID(), stopAreaId, partner, status)
		manager.BroadcastLegacyStopAreaUpdateEvent(event)
	}
}

func (manager *CollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	stopArea, ok := manager.referential.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("Can't find StopArea %v in Collect Manager", request.StopAreaId())
		return
	}
	partner := manager.bestPartner(stopArea)
	if partner == nil {
		// logger.Log.Debugf("Can't find a partner for StopArea %v in Collect Manager", request.StopAreaId())
		return
	}
	manager.requestStopAreaUpdate(partner, request)
}

func (manager *CollectManager) bestPartner(stopArea model.StopArea) *Partner {
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

func (manager *CollectManager) HandleSituationUpdateEvent(SituationUpdateSubscriber SituationUpdateSubscriber) {
	manager.SituationUpdateSubscribers = append(manager.SituationUpdateSubscribers, SituationUpdateSubscriber)
}

func (manager *CollectManager) BroadcastSituationUpdateEvent(event []*model.SituationUpdateEvent) {
	for _, SituationUpdateSubscriber := range manager.SituationUpdateSubscribers {
		SituationUpdateSubscriber(event)
	}
}

func (manager *CollectManager) UpdateSituation(request *SituationUpdateRequest) {
	switch request.Kind() {
	case SITUATION_UPDATE_REQUEST_ALL:
		manager.requestAllSituations()
	case SITUATION_UPDATE_REQUEST_LINE:
		manager.requestLineFilteredSituation(request.RequestedId())
	case SITUATION_UPDATE_REQUEST_STOP_AREA:
		manager.requestStopAreaFilteredSituation(request.RequestedId())
	default:
		logger.Log.Debugf("SituationUpdateRequest of unknown kind")
	}
}

func (manager *CollectManager) requestAllSituations() {
	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
			continue
		}
		if b, _ := strconv.ParseBool(partner.Setting("collect.filter_general_messages")); b {
			continue
		}

		requestConnector := partner.GeneralMessageRequestCollector()
		subscriptionConnector := partner.GeneralMessageSubscriptionCollector()
		if requestConnector == nil && subscriptionConnector == nil {
			continue
		}

		logger.Log.Debugf("RequestAllSituationsUpdate for Partner %v", partner.Slug())
		if subscriptionConnector != nil {
			subscriptionConnector.RequestAllSituationsUpdate()
			continue
		}
		requestConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_ALL, "")
	}
}

func (manager *CollectManager) requestLineFilteredSituation(requestedId string) {
	line, ok := manager.referential.Model().Lines().Find(model.LineId(requestedId))
	if !ok {
		logger.Log.Debugf("Can't find Line to request %v", requestedId)
		return
	}

	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
			continue
		}
		if b, _ := strconv.ParseBool(partner.Setting("collect.filter_general_messages")); !b {
			continue
		}

		requestConnector := partner.GeneralMessageRequestCollector()
		subscriptionConnector := partner.GeneralMessageSubscriptionCollector()

		if requestConnector == nil && subscriptionConnector == nil {
			continue
		}

		partnerKind := partner.Setting("remote_objectid_kind")

		lineObjectID, ok := line.ObjectID(partnerKind)
		if !ok {
			continue
		}

		if !partner.CanCollectLine(lineObjectID) {
			continue
		}

		logger.Log.Debugf("RequestSituationUpdate %v with Partner %v", lineObjectID.Value(), partner.Slug())
		if subscriptionConnector != nil {
			subscriptionConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, lineObjectID)
			return
		}
		requestConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, lineObjectID.Value())
		return
	}
	// logger.Log.Debugf("Can't find a partner to request filtered Situations for Line %v", requestedId)
}

func (manager *CollectManager) requestStopAreaFilteredSituation(requestedId string) {
	stopArea, ok := manager.referential.Model().StopAreas().Find(model.StopAreaId(requestedId))
	if !ok {
		logger.Log.Debugf("Can't find StopArea to request %v", requestedId)
		return
	}

	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
			continue
		}
		if b, _ := strconv.ParseBool(partner.Setting("collect.filter_general_messages")); !b {
			continue
		}

		requestConnector := partner.GeneralMessageRequestCollector()
		subscriptionConnector := partner.GeneralMessageSubscriptionCollector()

		if requestConnector == nil && subscriptionConnector == nil {
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

		if !partner.CanCollect(stopAreaObjectID, lineIds) {
			continue
		}

		logger.Log.Debugf("RequestSituationUpdate %v with Partner %v", stopAreaObjectID.Value(), partner.Slug())
		if subscriptionConnector != nil {
			subscriptionConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_STOP_AREA, stopAreaObjectID)
			return
		}
		requestConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_STOP_AREA, stopAreaObjectID.Value())
		return
	}
	// logger.Log.Debugf("Can't find a partner to request filtered Situations for StopArea %v", requestedId)
}
