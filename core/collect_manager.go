package core

import (
	"context"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

type UpdateSubscriber func(model.UpdateEvent)

type CollectManagerInterface interface {
	HandlePartnerStatusChange(partner string, status bool)

	UpdateStopArea(request *StopAreaUpdateRequest)
	UpdateLine(ctx context.Context, request *LineUpdateRequest)
	UpdateVehicle(ctx context.Context, request *VehicleUpdateRequest)
	UpdateSituation(request *SituationUpdateRequest)
	UpdateFacility(request *FacilityUpdateRequest)

	HandleUpdateEvent(UpdateSubscriber UpdateSubscriber)
	BroadcastUpdateEvent(event model.UpdateEvent)
}

type CollectManager struct {
	uuid.UUIDConsumer

	UpdateSubscribers []UpdateSubscriber
	referential       *Referential
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
func (manager *TestCollectManager) UpdateLine(context.Context, *LineUpdateRequest)       {}
func (manager *TestCollectManager) UpdateVehicle(context.Context, *VehicleUpdateRequest) {}
func (manager *TestCollectManager) UpdateFacility(*FacilityUpdateRequest)                {}

// TEST END

func NewCollectManager(referential *Referential) CollectManagerInterface {
	return &CollectManager{
		referential:       referential,
		UpdateSubscribers: make([]UpdateSubscriber, 0),
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
	localLogger := NewStopAreaLogger(manager.referential, stopArea)

	if !ok {
		localLogger.Printf("Can't find StopArea %v in Collect Manager", request.StopAreaId())
		return
	}

	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		subscriptionCollector := partner.StopMonitoringSubscriptionCollector()
		SMrequestCollector := partner.StopMonitoringRequestCollector()
		SMLiterequestCollector := partner.LiteStopMonitoringRequestCollector()

		if subscriptionCollector == nil && SMrequestCollector == nil && SMLiterequestCollector == nil {
			localLogger.Printf("No Collector for Partner %s", partner.Slug())
			continue
		}

		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP && !partner.PersistentCollect() {
			localLogger.Printf("Partner %s isn't up", partner.Slug())
			continue
		}

		partnerCodeSpace := partner.RemoteCodeSpace()

		stopAreaCode, ok := stopArea.Code(partnerCodeSpace)
		if !ok {
			localLogger.Printf("No Code matching Partner CodeSpace (%s)", partnerCodeSpace)
			continue
		}

		lineIds := make(map[string]struct{})
		for _, lineId := range stopArea.LineIds {
			line, ok := manager.referential.Model().Lines().Find(lineId)
			if !ok {
				continue
			}
			lineCode, ok := line.Code(partnerCodeSpace)
			if !ok {
				continue
			}
			lineIds[lineCode.Value()] = struct{}{}
		}

		if partner.CanCollect(stopAreaCode.Value(), lineIds) {
			localLogger.Printf("RequestStopAreaUpdate %v", request.StopAreaId())

			if subscriptionCollector != nil {
				subscriptionCollector.RequestStopAreaUpdate(request)
				return
			} else if SMrequestCollector != nil {
				SMrequestCollector.RequestStopAreaUpdate(request)
				return
			} else if SMLiterequestCollector != nil {
				SMLiterequestCollector.RequestStopAreaUpdate(request)
				return
			}
		} else {
			localLogger.Printf("Partner %s can't collect StopArea", partner.Slug())
			localLogger.Printf("Partner %s collect settings: %#v", partner.Slug(), partner.CollectSettings())
		}

	}
}

func (manager *CollectManager) UpdateLine(ctx context.Context, request *LineUpdateRequest) {
	child, _ := tracer.StartSpanFromContext(ctx, "update_line")
	defer child.Finish()
	line, ok := manager.referential.Model().Lines().Find(request.LineId())
	if !ok {
		logger.Log.Debugf("Can't find Line %v in Collect Manager", request.LineId())
		return
	}

	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		subscriptionCollector := partner.EstimatedTimetableSubscriptionCollector()
		// requestCollector := partner.EstimatedTimetableRequestCollector()

		// if subscriptionCollector == nil && requestCollector == nil {
		if subscriptionCollector == nil {
			continue
		}

		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP && !partner.PersistentCollect() {
			logger.Log.Debugf("Partner %s isn't up", partner.Slug())
			continue
		}

		partnerCodeSpace := partner.RemoteCodeSpace()

		lineCode, ok := line.Code(partnerCodeSpace)
		if !ok {
			continue
		}

		if !partner.CanCollectLine(lineCode.Value()) {
			continue
		}
		logger.Log.Debugf("RequestLineUpdate with LineId %v", request.LineId())

		// if subscriptionCollector != nil {
		// 	subscriptionCollector.RequestLineUpdate(request)
		// 	return
		// }
		// requestCollector.RequestLineUpdate(request)
		// return
		subscriptionCollector.RequestLineUpdate(request)
		return
	}
}

func (manager *CollectManager) UpdateVehicle(ctx context.Context, request *VehicleUpdateRequest) {
	child, _ := tracer.StartSpanFromContext(ctx, "update_vehicle")
	defer child.Finish()
	line, ok := manager.referential.Model().Lines().Find(request.LineId())
	if !ok {
		logger.Log.Debugf("Can't find Line %v in Collect Manager", request.LineId())
		return
	}

	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		subscriptionCollector := partner.VehicleMonitoringSubscriptionCollector()
		requestCollector := partner.VehicleMonitoringRequestCollector()

		if subscriptionCollector == nil && requestCollector == nil {
			continue
		}

		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP && !partner.PersistentCollect() {
			logger.Log.Debugf("Partner %s isn't up", partner.Slug())
			continue
		}

		partnerCodeSpace := partner.RemoteCodeSpace()

		lineCode, ok := line.Code(partnerCodeSpace)
		if !ok {
			continue
		}

		if !partner.CanCollectLine(lineCode.Value()) {
			continue
		}
		logger.Log.Debugf("RequestVehicleUpdate with LineId %v", request.LineId())

		if subscriptionCollector != nil {
			subscriptionCollector.RequestVehicleUpdate(request)
			return
		}
		requestCollector.RequestVehicleUpdate(request)
		return
	}
}

func (manager *CollectManager) UpdateFacility(request *FacilityUpdateRequest) {
	for _, partner := range manager.referential.Partners().FindAllByCollectPriority() {
		if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
			continue
		}

		requestConnector := partner.FacilityMonitoringRequestCollector()

		logger.Log.Debugf("Request FacilityUpdate for Partner %v", partner.Slug())
		if requestConnector != nil {
			requestConnector.RequestFacilityUpdate(request)
		}
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
		if partner.CollectFilteredSituations() {
			continue
		}

		requestConnector := partner.SituationExchangeRequestCollector()
		gmSubscriptionConnector := partner.GeneralMessageSubscriptionCollector()
		sxSubscriptionConnector := partner.SituationExchangeSubscriptionCollector()
		if requestConnector == nil &&
			gmSubscriptionConnector == nil &&
			sxSubscriptionConnector == nil {
			continue
		}

		logger.Log.Debugf("RequestAllSituationsUpdate for Partner %v", partner.Slug())
		if gmSubscriptionConnector != nil {
			gmSubscriptionConnector.RequestAllSituationsUpdate()
			continue
		}

		if sxSubscriptionConnector != nil {
			sxSubscriptionConnector.RequestAllSituationsUpdate()
			continue
		}
		if requestConnector != nil {
			requestConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_ALL, "")

		}
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
		if !partner.CollectFilteredSituations() {
			continue
		}

		requestConnector := partner.SituationExchangeRequestCollector()
		gmSubscriptionConnector := partner.GeneralMessageSubscriptionCollector()
		sxSubscriptionConnector := partner.SituationExchangeSubscriptionCollector()

		if requestConnector == nil &&
			gmSubscriptionConnector == nil &&
			sxSubscriptionConnector == nil {
			continue
		}

		partnerCodeSpace := partner.RemoteCodeSpace()

		lineCode, ok := line.Code(partnerCodeSpace)
		if !ok {
			continue
		}

		if !partner.CanCollectLine(lineCode.Value()) {
			continue
		}

		logger.Log.Debugf("RequestSituationUpdate %v with Partner %v", lineCode.Value(), partner.Slug())
		if gmSubscriptionConnector != nil {
			gmSubscriptionConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, lineCode)
			return
		}

		if sxSubscriptionConnector != nil {
			sxSubscriptionConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, lineCode)
			return
		}

		requestConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, lineCode.Value())
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
		if !partner.CollectFilteredSituations() {
			continue
		}

		requestConnector := partner.SituationExchangeRequestCollector()
		gmSubscriptionConnector := partner.GeneralMessageSubscriptionCollector()
		sxSubscriptionConnector := partner.SituationExchangeSubscriptionCollector()

		if requestConnector == nil &&
			gmSubscriptionConnector == nil &&
			sxSubscriptionConnector == nil {
			continue
		}

		partnerCodeSpace := partner.RemoteCodeSpace()

		stopAreaCode, ok := stopArea.Code(partnerCodeSpace)
		if !ok {
			continue
		}

		lineIds := make(map[string]struct{})
		for _, lineId := range stopArea.LineIds {
			line, ok := manager.referential.Model().Lines().Find(lineId)
			if !ok {
				continue
			}
			lineCode, ok := line.Code(partnerCodeSpace)
			if !ok {
				continue
			}
			lineIds[lineCode.Value()] = struct{}{}
		}

		if !partner.CanCollect(stopAreaCode.Value(), lineIds) {
			continue
		}

		logger.Log.Debugf("RequestSituationUpdate %v with Partner %v", stopAreaCode.Value(), partner.Slug())
		if gmSubscriptionConnector != nil {
			gmSubscriptionConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_STOP_AREA, stopAreaCode)
			return
		}
		if sxSubscriptionConnector != nil {
			sxSubscriptionConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_STOP_AREA, stopAreaCode)
			return
		}
		requestConnector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_STOP_AREA, stopAreaCode.Value())
		return
	}
	// logger.Log.Debugf("Can't find a partner to request filtered Situations for StopArea %v", requestedId)
}
