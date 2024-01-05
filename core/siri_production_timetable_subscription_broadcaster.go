package core

import (
	"fmt"
	"strings"
	"sync"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIProductionTimetableSubscriptionBroadcaster struct {
	connector

	dataFrameGenerator             *idgen.IdentifierGenerator
	noDataFrameRefRewritingFrom    []string
	vjRemoteCodeSpaces          []string
	productionTimetableBroadcaster SIRIProductionTimetableBroadcaster
	toBroadcast                    map[SubscriptionId][]model.StopVisitId

	mutex *sync.Mutex //protect the map
}

type SIRIProductionTimetableSubscriptionBroadcasterFactory struct{}

func (factory *SIRIProductionTimetableSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIProductionTimetableSubscriptionBroadcaster(partner)
}

func (factory *SIRIProductionTimetableSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func newSIRIProductionTimetableSubscriptionBroadcaster(partner *Partner) *SIRIProductionTimetableSubscriptionBroadcaster {
	connector := &SIRIProductionTimetableSubscriptionBroadcaster{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace(SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER)
	connector.vjRemoteCodeSpaces = partner.VehicleJourneyRemoteCodeSpaceWithFallback(SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER)
	connector.noDataFrameRefRewritingFrom = partner.NoDataFrameRefRewritingFrom()
	connector.dataFrameGenerator = partner.DataFrameIdentifierGenerator()
	connector.partner = partner
	connector.mutex = &sync.Mutex{}
	connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	connector.productionTimetableBroadcaster = NewSIRIProductionTimetableBroadcaster(connector)
	return connector
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) Start() {
	connector.productionTimetableBroadcaster.Start()
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) Stop() {
	connector.productionTimetableBroadcaster.Stop()
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) HandleSubscriptionRequest(request *sxml.XMLSubscriptionRequest, message *audit.BigQueryMessage) (resps []siri.SIRIResponseStatus) {
	var lineIds, subIds []string

	for _, ptt := range request.XMLSubscriptionPTTEntries() {
		rs := siri.SIRIResponseStatus{
			RequestMessageRef: ptt.MessageIdentifier(),
			SubscriberRef:     ptt.SubscriberRef(),
			SubscriptionRef:   ptt.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		// for logging
		lineIds = append(lineIds, ptt.Lines()...)

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(ptt.SubscriptionIdentifier())
		if ok {
			if sub.Kind() != ProductionTimetableBroadcast {
				logger.Log.Debugf("ProductionTimetable subscription request with a duplicated Id: %v", ptt.SubscriptionIdentifier())
				rs.ErrorType = "OtherError"
				rs.ErrorNumber = 2
				rs.ErrorText = fmt.Sprintf("[BAD_REQUEST] Subscription Id %v already exists", ptt.SubscriptionIdentifier())

				resps = append(resps, rs)
				message.Status = "Error"
				continue
			}

			sub.Delete()
		}

		resources, unknownLineIds := connector.checkLines(ptt)
		if len(unknownLineIds) != 0 {
			logger.Log.Debugf("ProductionTimetable subscription request Could not find line(s) with id : %v", strings.Join(unknownLineIds, ","))
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = fmt.Sprintf("Unknown Line(s) %v", strings.Join(unknownLineIds, ","))

			resps = append(resps, rs)
			message.Status = "Error"
			continue
		}

		rs.Status = true
		rs.ValidUntil = ptt.InitialTerminationTime()
		resps = append(resps, rs)

		subIds = append(subIds, ptt.SubscriptionIdentifier())

		sub = connector.Partner().Subscriptions().New(ProductionTimetableBroadcast)
		sub.SubscriberRef = ptt.SubscriberRef()
		sub.SetExternalId(ptt.SubscriptionIdentifier())
		sub.SetSubscriptionOption("MessageIdentifier", request.MessageIdentifier())

		for _, r := range resources {
			line, ok := connector.Partner().Model().Lines().FindByCode(*r.Reference.Code)
			if !ok {
				continue
			}

			// Init StopVisits LastChange
			connector.addLineStopVisits(sub, r, line.Id())

			sub.AddNewResource(r)
		}
		sub.Save()
	}
	message.Type = audit.PRODUCTION_TIMETABLE_SUBSCRIPTION_REQUEST
	message.SubscriptionIdentifiers = subIds
	message.Lines = lineIds

	return resps
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) checkLines(ptt *sxml.XMLProductionTimetableSubscriptionRequestEntry) (resources []*SubscribedResource, lineIds []string) {
	// check for subscription to all lines
	if len(ptt.Lines()) == 0 {
		var lv []string
		//find all lines corresponding to the remoteCodeSpace
		for _, line := range connector.Partner().Model().Lines().FindAll() {
			lineCode, ok := line.Code(connector.remoteCodeSpace)
			if ok {
				lv = append(lv, lineCode.Value())
				continue
			}
		}

		for _, lineValue := range lv {
			lineCode := model.NewCode(connector.remoteCodeSpace, lineValue)
			ref := model.Reference{
				Code: &lineCode,
				Type:     "Line",
			}
			r := NewResource(ref)
			r.Subscribed(connector.Clock().Now())
			r.SubscribedUntil = ptt.InitialTerminationTime()
			resources = append(resources, r)
		}
		return resources, lineIds
	}

	for _, lineId := range ptt.Lines() {

		lineCode := model.NewCode(connector.remoteCodeSpace, lineId)
		_, ok := connector.Partner().Model().Lines().FindByCode(lineCode)

		if !ok {
			lineIds = append(lineIds, lineId)
			continue
		}

		ref := model.Reference{
			Code: &lineCode,
			Type:     "Line",
		}

		r := NewResource(ref)
		r.Subscribed(connector.Clock().Now())
		r.SubscribedUntil = ptt.InitialTerminationTime()
		resources = append(resources, r)
	}
	return resources, lineIds
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) addLineStopVisits(sub *Subscription, res *SubscribedResource, lineId model.LineId) {
	sas := connector.partner.Model().StopAreas().FindByLineId(lineId)
	for i := range sas {
		svs := connector.partner.Model().ScheduledStopVisits().FindByStopAreaId(sas[i].Id())
		for i := range svs {
			connector.addStopVisit(sub.Id(), svs[i].Id())
		}
	}
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) HandleBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.checkEvent(model.StopVisitId(event.ModelId))
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) checkEvent(svId model.StopVisitId) {
	sv, ok := connector.Partner().Model().ScheduledStopVisits().Find(svId)
	if !ok {
		return
	}

	vj, ok := connector.Partner().Model().VehicleJourneys().Find(sv.VehicleJourneyId)
	if !ok {
		return
	}

	line, ok := connector.Partner().Model().Lines().Find(vj.LineId)
	if !ok {
		return
	}

	lineObj, ok := line.Code(connector.remoteCodeSpace)
	if !ok {
		return
	}

	subs := connector.Partner().Subscriptions().FindByResourceId(lineObj.String(), ProductionTimetableBroadcast)

	for _, sub := range subs {
		r := sub.Resource(lineObj)
		if r == nil || r.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := r.LastState(string(svId))
		if ok && !lastState.(*ls.ProductionTimetableLastChange).Haschanged(sv) {
			continue
		}

		if !ok {
			r.SetLastState(string(sv.Id()), ls.NewProductionTimetableLastChange(sv, sub))
		}

		connector.addStopVisit(sub.Id(), svId)
	}
}
