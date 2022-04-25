package core

import (
	"fmt"
	"strings"
	"sync"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type SIRIProductionTimeTableSubscriptionBroadcaster struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	dataFrameGenerator             *idgen.IdentifierGenerator
	noDataFrameRefRewritingFrom    []string
	vjRemoteObjectidKinds          []string
	productionTimeTableBroadcaster SIRIProductionTimeTableBroadcaster
	toBroadcast                    map[SubscriptionId][]model.StopVisitId

	mutex *sync.Mutex //protect the map
}

type SIRIProductionTimetableSubscriptionBroadcasterFactory struct{}

func (factory *SIRIProductionTimetableSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIProductionTimeTableSubscriptionBroadcaster(partner)
}

func (factory *SIRIProductionTimetableSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func newSIRIProductionTimeTableSubscriptionBroadcaster(partner *Partner) *SIRIProductionTimeTableSubscriptionBroadcaster {
	connector := &SIRIProductionTimeTableSubscriptionBroadcaster{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
	connector.vjRemoteObjectidKinds = partner.VehicleJourneyRemoteObjectIDKindWithFallback(SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER)
	connector.noDataFrameRefRewritingFrom = partner.NoDataFrameRefRewritingFrom()
	connector.dataFrameGenerator = partner.IdentifierGenerator(idgen.DATA_FRAME_IDENTIFIER)
	connector.partner = partner
	connector.mutex = &sync.Mutex{}
	connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	connector.productionTimeTableBroadcaster = NewSIRIProductionTimeTableBroadcaster(connector)
	return connector
}

func (connector *SIRIProductionTimeTableSubscriptionBroadcaster) Start() {
	connector.productionTimeTableBroadcaster.Start()
}

func (connector *SIRIProductionTimeTableSubscriptionBroadcaster) Stop() {
	connector.productionTimeTableBroadcaster.Stop()
}

func (connector *SIRIProductionTimeTableSubscriptionBroadcaster) HandleSubscriptionRequest(request *siri.XMLSubscriptionRequest, message *audit.BigQueryMessage) (resps []siri.SIRIResponseStatus) {
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

		resources, unknownLineIds := connector.checkLines(ptt)
		if len(unknownLineIds) != 0 {
			logger.Log.Debugf("ProductionTimeTable subscription request Could not find line(s) with id : %v", strings.Join(unknownLineIds, ","))
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = fmt.Sprintf("Unknown Line(s) %v", strings.Join(unknownLineIds, ","))
		} else {
			rs.Status = true
			rs.ValidUntil = ptt.InitialTerminationTime()
		}

		resps = append(resps, rs)

		// We do not want to create a subscription that will fail
		if len(unknownLineIds) != 0 {
			continue
		}

		subIds = append(subIds, ptt.SubscriptionIdentifier())

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(ptt.SubscriptionIdentifier())
		if !ok {
			sub = connector.Partner().Subscriptions().New("ProductionTimeTableBroadcast")
			sub.SubscriberRef = ptt.SubscriberRef()
			sub.SetExternalId(ptt.SubscriptionIdentifier())
			sub.SetSubscriptionOption("MessageIdentifier", request.MessageIdentifier())
		}

		for _, r := range resources {
			line, ok := connector.Partner().Model().Lines().FindByObjectId(*r.Reference.ObjectId)
			if !ok {
				continue
			}

			// Init StopVisits LastChange
			connector.addLineStopVisits(sub, r, line.Id())

			sub.AddNewResource(*r)
		}
		sub.Save()
	}
	message.Type = "ProductionTimetableSubscriptionRequest"
	message.SubscriptionIdentifiers = subIds
	message.Lines = lineIds

	return resps
}

func (connector *SIRIProductionTimeTableSubscriptionBroadcaster) checkLines(ptt *siri.XMLProductionTimetableSubscriptionRequestEntry) (resources []*SubscribedResource, lineIds []string) {
	// check for subscription to all lines
	if len(ptt.Lines()) == 0 {
		var lv []string
		//find all lines corresponding to the remoteObjectidKind
		for _, line := range connector.Partner().Model().Lines().FindAll() {
			lineObjectID, ok := line.ObjectID(connector.remoteObjectidKind)
			if ok {
				lv = append(lv, lineObjectID.Value())
				continue
			}
		}

		for _, lineValue := range lv {
			lineObjectID := model.NewObjectID(connector.remoteObjectidKind, lineValue)
			ref := model.Reference{
				ObjectId: &lineObjectID,
				Type:     "Line",
			}
			r := NewResource(ref)
			r.SubscribedAt = connector.Clock().Now()
			r.SubscribedUntil = ptt.InitialTerminationTime()
			resources = append(resources, &r)
		}
		return resources, lineIds
	}

	for _, lineId := range ptt.Lines() {

		lineObjectID := model.NewObjectID(connector.remoteObjectidKind, lineId)
		_, ok := connector.Partner().Model().Lines().FindByObjectId(lineObjectID)

		if !ok {
			lineIds = append(lineIds, lineId)
			continue
		}

		ref := model.Reference{
			ObjectId: &lineObjectID,
			Type:     "Line",
		}

		r := NewResource(ref)
		r.SubscribedAt = connector.Clock().Now()
		r.SubscribedUntil = ptt.InitialTerminationTime()
		resources = append(resources, &r)
	}
	return resources, lineIds
}

func (connector *SIRIProductionTimeTableSubscriptionBroadcaster) addLineStopVisits(sub *Subscription, res *SubscribedResource, lineId model.LineId) {
	sas := connector.partner.Model().StopAreas().FindByLineId(lineId)
	for i := range sas {
		svs := connector.partner.Model().ScheduledStopVisits().FindByStopAreaId(sas[i].Id())
		for i := range svs {
			connector.addStopVisit(sub.Id(), svs[i].Id())
		}
	}
}

func (connector *SIRIProductionTimeTableSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIProductionTimeTableSubscriptionBroadcaster) HandleBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.checkEvent(model.StopVisitId(event.ModelId))
}

func (connector *SIRIProductionTimeTableSubscriptionBroadcaster) checkEvent(svId model.StopVisitId) {
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

	lineObj, ok := line.ObjectID(connector.remoteObjectidKind)
	if !ok {
		return
	}

	subs := connector.Partner().Subscriptions().FindByResourceId(lineObj.String(), "ProductionTimeTableBroadcast")

	for _, sub := range subs {
		r := sub.Resource(lineObj)
		if r == nil || r.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := r.LastState(string(svId))
		if ok && !lastState.(*ls.ProductionTimeTableLastChange).Haschanged(&sv) {
			continue
		}

		if !ok {
			r.SetLastState(string(sv.Id()), ls.NewProductionTimeTableLastChange(&sv, sub))
		}

		connector.addStopVisit(sub.Id(), svId)
	}
}
