package core

import (
	"strconv"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type SituationUpdateSubscriber func([]*model.SituationUpdateEvent)
type UpdateSubscriber func(model.UpdateEvent)

type CollectManagerInterface interface {
	HandlePartnerStatusChange(partner string, status bool)
	UpdateStopArea(request *StopAreaUpdateRequest)

	HandleUpdateEvent(UpdateSubscriber UpdateSubscriber)
	BroadcastUpdateEvent(event model.UpdateEvent)

	UpdateSituation(request *SituationUpdateRequest)
	HandleSituationUpdateEvent(SituationUpdateSubscriber)
	BroadcastSituationUpdateEvent(event []*model.SituationUpdateEvent)
}

type CollectManager struct {
	uuid.UUIDConsumer

	SituationUpdateSubscribers []SituationUpdateSubscriber
	UpdateSubscribers          []UpdateSubscriber
	referential                *Referential
}

// TestCollectManager has a test StopAreaUpdateSubscriber method
type TestCollectManager struct {
	Done         chan bool
	UpdateEvents []model.UpdateEvent
}

func NewTestCollectManager() CollectManagerInterface {
	return &TestCollectManager{
		Done: make(chan bool, 1),
	}
}

func (manager *TestCollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	event := &model.StopAreaUpdateEvent{}
	manager.UpdateEvents = append(manager.UpdateEvents, event)

	manager.Done <- true
}

func (manager *TestCollectManager) TestUpdateSubscriber(event model.UpdateEvent) {
	manager.UpdateEvents = append(manager.UpdateEvents, event)
}

func (manager *TestCollectManager) HandlePartnerStatusChange(partner string, status bool) {}

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
		SituationUpdateSubscribers: make([]SituationUpdateSubscriber, 0),
		UpdateSubscribers:          make([]UpdateSubscriber, 0),
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
		event := model.NewStatusUpdateEvent(stopAreaId, partner, status)
		manager.BroadcastUpdateEvent(event)
	}
}

func (manager *CollectManager) UpdateStopArea(request *StopAreaUpdateRequest) {
	stopArea, ok := manager.referential.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("Can't find StopArea %v in Collect Manager", request.StopAreaId())
		return
	}

	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		subscriptionCollector := partner.StopMonitoringSubscriptionCollector()
		requestCollector := partner.StopMonitoringRequestCollector()

		if subscriptionCollector == nil && requestCollector == nil {
			continue
		}

		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
			if b, _ := strconv.ParseBool(partner.Setting(COLLECT_SUBSCRIPTIONS_PERSISTENT)); !b || subscriptionCollector == nil {
				continue
			}
		}

		partnerKind := partner.Setting(REMOTE_OBJECTID_KIND)

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
			logger.Log.Debugf("RequestStopAreaUpdate %v", request.StopAreaId())

			if subscriptionCollector != nil {
				subscriptionCollector.RequestStopAreaUpdate(request)
				return
			}
			requestCollector.RequestStopAreaUpdate(request)
			return
		}
	}
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
		if b, _ := strconv.ParseBool(partner.Setting(COLLECT_FILTER_GENERAL_MESSAGES)); b {
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
		if b, _ := strconv.ParseBool(partner.Setting(COLLECT_FILTER_GENERAL_MESSAGES)); !b {
			continue
		}

		requestConnector := partner.GeneralMessageRequestCollector()
		subscriptionConnector := partner.GeneralMessageSubscriptionCollector()

		if requestConnector == nil && subscriptionConnector == nil {
			continue
		}

		partnerKind := partner.Setting(REMOTE_OBJECTID_KIND)

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
		if b, _ := strconv.ParseBool(partner.Setting(COLLECT_FILTER_GENERAL_MESSAGES)); !b {
			continue
		}

		requestConnector := partner.GeneralMessageRequestCollector()
		subscriptionConnector := partner.GeneralMessageSubscriptionCollector()

		if requestConnector == nil && subscriptionConnector == nil {
			continue
		}

		partnerKind := partner.Setting(REMOTE_OBJECTID_KIND)

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
