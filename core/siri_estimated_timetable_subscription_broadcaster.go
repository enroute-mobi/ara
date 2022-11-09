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

type SIRIEstimatedTimetableSubscriptionBroadcaster struct {
	connector

	dataFrameGenerator            *idgen.IdentifierGenerator
	vjRemoteObjectidKinds         []string
	estimatedTimetableBroadcaster EstimatedTimetableBroadcaster
	toBroadcast                   map[SubscriptionId][]model.StopVisitId
	notMonitored                  map[SubscriptionId]map[string]struct{}

	mutex *sync.Mutex //protect the map
}

type SIRIEstimatedTimetableSubscriptionBroadcasterFactory struct{}

func (factory *SIRIEstimatedTimetableSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIEstimatedTimetableSubscriptionBroadcaster(partner)
}

func (factory *SIRIEstimatedTimetableSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func newSIRIEstimatedTimetableSubscriptionBroadcaster(partner *Partner) *SIRIEstimatedTimetableSubscriptionBroadcaster {
	connector := &SIRIEstimatedTimetableSubscriptionBroadcaster{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
	connector.vjRemoteObjectidKinds = partner.VehicleJourneyRemoteObjectIDKindWithFallback(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
	connector.dataFrameGenerator = partner.IdentifierGenerator(idgen.DATA_FRAME_IDENTIFIER)
	connector.partner = partner
	connector.mutex = &sync.Mutex{}
	connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	connector.estimatedTimetableBroadcaster = NewSIRIEstimatedTimetableBroadcaster(connector)
	return connector
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) HandleSubscriptionRequest(request *sxml.XMLSubscriptionRequest, message *audit.BigQueryMessage) (resps []siri.SIRIResponseStatus) {
	var lineIds, subIds []string

	for _, ett := range request.XMLSubscriptionETTEntries() {
		rs := siri.SIRIResponseStatus{
			RequestMessageRef: ett.MessageIdentifier(),
			SubscriberRef:     ett.SubscriberRef(),
			SubscriptionRef:   ett.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		// for logging
		lineIds = append(lineIds, ett.Lines()...)

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(ett.SubscriptionIdentifier())
		if ok && sub.Kind() != EstimatedTimetableBroadcast {
			logger.Log.Debugf("EstimatedTimetable subscription request with a duplicated Id: %v", ett.SubscriptionIdentifier())
			rs.ErrorType = "OtherError"
			rs.ErrorNumber = 2
			rs.ErrorText = fmt.Sprintf("[BAD_REQUEST] Subscription Id %v already exists", ett.SubscriptionIdentifier())

			resps = append(resps, rs)
			message.Status = "Error"
			continue
		}

		resources, unknownLineIds := connector.checkLines(ett)
		if len(unknownLineIds) != 0 {
			logger.Log.Debugf("EstimatedTimetable subscription request Could not find line(s) with id : %v", strings.Join(unknownLineIds, ","))
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = fmt.Sprintf("Unknown Line(s) %v", strings.Join(unknownLineIds, ","))

			resps = append(resps, rs)
			message.Status = "Error"
			continue
		}

		rs.Status = true
		rs.ValidUntil = ett.InitialTerminationTime()
		resps = append(resps, rs)

		subIds = append(subIds, ett.SubscriptionIdentifier())

		if !ok {
			sub = connector.Partner().Subscriptions().New(EstimatedTimetableBroadcast)
			sub.SubscriberRef = ett.SubscriberRef()
			sub.SetExternalId(ett.SubscriptionIdentifier())
			connector.fillOptions(sub, request)
		}

		for _, r := range resources {
			line, ok := connector.Partner().Model().Lines().FindByObjectId(*r.Reference.ObjectId)
			if !ok {
				continue
			}

			// Init StopVisits LastChange
			connector.addLineStopVisits(sub, r, line.Id())

			sub.AddNewResource(r)
		}
		sub.Save()
	}
	message.Type = "EstimatedTimetableSubscriptionRequest"
	message.SubscriptionIdentifiers = subIds
	message.Lines = lineIds

	return resps
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) addLineStopVisits(sub *Subscription, res *SubscribedResource, lineId model.LineId) {
	sas := connector.partner.Model().StopAreas().FindByLineId(lineId)
	for i := range sas {
		// Init SA LastChange
		res.SetLastState(string(sas[i].Id()), ls.NewStopAreaLastChange(sas[i], sub))
		svs := connector.partner.Model().StopVisits().FindFollowingByStopAreaId(sas[i].Id())
		for i := range svs {
			connector.addStopVisit(sub.Id(), svs[i].Id())
		}
	}
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) checkLines(ett *sxml.XMLEstimatedTimetableSubscriptionRequestEntry) (resources []*SubscribedResource, lineIds []string) {
	// check for subscription to all lines
	if len(ett.Lines()) == 0 {
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
			r.Subscribed(connector.Clock().Now())
			r.SubscribedUntil = ett.InitialTerminationTime()
			resources = append(resources, r)
		}
		return resources, lineIds
	}

	for _, lineId := range ett.Lines() {

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
		r.Subscribed(connector.Clock().Now())
		r.SubscribedUntil = ett.InitialTerminationTime()
		resources = append(resources, r)
	}
	return resources, lineIds
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) Stop() {
	connector.estimatedTimetableBroadcaster.Stop()
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) Start() {
	connector.estimatedTimetableBroadcaster.Start()
}

func (ettb *SIRIEstimatedTimetableSubscriptionBroadcaster) fillOptions(s *Subscription, request *sxml.XMLSubscriptionRequest) {
	changeBeforeUpdates := request.ChangeBeforeUpdates()
	if changeBeforeUpdates == "" {
		changeBeforeUpdates = "PT1M"
	}
	s.SetSubscriptionOption("ChangeBeforeUpdates", changeBeforeUpdates)
	s.SetSubscriptionOption("MessageIdentifier", request.MessageIdentifier())
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) HandleBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	switch event.ModelType {
	case "StopVisit":
		connector.checkEvent(model.StopVisitId(event.ModelId))
	case "StopArea":
		sa, ok := connector.partner.Model().StopAreas().Find(model.StopAreaId(event.ModelId))
		if ok {
			connector.checkStopAreaEvent(sa)
		}
	default:
		return
	}
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) checkEvent(svId model.StopVisitId) {
	sv, ok := connector.Partner().Model().StopVisits().Find(svId)
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

	subs := connector.Partner().Subscriptions().FindByResourceId(lineObj.String(), EstimatedTimetableBroadcast)

	for _, sub := range subs {
		r := sub.Resource(lineObj)
		if r == nil || r.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := r.LastState(string(svId))
		if ok && !lastState.(*ls.EstimatedTimetableLastChange).Haschanged(sv) {
			continue
		}

		if !ok {
			r.SetLastState(string(sv.Id()), ls.NewEstimatedTimetableLastChange(sv, sub))
		}
		connector.addFilteredStopVisit(sub.Id(), sv)
	}
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) addFilteredStopVisit(subId SubscriptionId, sv *model.StopVisit) {
	// ignore stopVist before the RecordedCallsDuration if any
	if connector.Partner().RecordedCallsDuration() != 0 &&
		sv.ReferenceTime().Before(connector.Clock().Now().Add(-connector.Partner().RecordedCallsDuration())) {
		return
	}

	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], sv.Id())
	connector.mutex.Unlock()
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) checkStopAreaEvent(stopArea *model.StopArea) {
	obj, ok := stopArea.ObjectID(connector.remoteObjectidKind)
	if !ok {
		return
	}

	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	subs := connector.partner.Subscriptions().FindByResourceId(obj.String(), EstimatedTimetableBroadcast)
	for _, sub := range subs {
		resource := sub.Resource(obj)
		if resource == nil || resource.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := resource.LastState(string(stopArea.Id()))
		if ok {
			partners, ok := lastState.(*ls.StopAreaLastChange).Haschanged(stopArea)
			if ok {
				nm, ok := connector.notMonitored[sub.Id()]
				if !ok {
					nm = make(map[string]struct{})
					connector.notMonitored[sub.Id()] = nm
				}
				for _, partner := range partners {
					nm[partner] = struct{}{}
				}
			}
			lastState.(*ls.StopAreaLastChange).UpdateState(stopArea)
		} else { // Should not happen
			resource.SetLastState(string(stopArea.Id()), ls.NewStopAreaLastChange(stopArea, sub))
		}
	}
}

// START TEST

type TestSIRIETTSubscriptionBroadcasterFactory struct{}

type TestETTSubscriptionBroadcaster struct {
	connector

	events []*model.StopMonitoringBroadcastEvent
}

func NewTestETTSubscriptionBroadcaster() *TestETTSubscriptionBroadcaster {
	connector := &TestETTSubscriptionBroadcaster{}
	return connector
}

func (connector *TestETTSubscriptionBroadcaster) HandleBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIETTSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {} // Always valid

func (factory *TestSIRIETTSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestETTSubscriptionBroadcaster()
}
