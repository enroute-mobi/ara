package core

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type SIRIEstimatedTimeTableSubscriptionBroadcaster struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	vjRemoteObjectidKinds         []string
	estimatedTimeTableBroadcaster SIRIEstimatedTimeTableBroadcaster
	toBroadcast                   map[SubscriptionId][]model.StopVisitId
	notMonitored                  map[SubscriptionId]map[string]struct{}

	mutex *sync.Mutex //protect the map
}

type SIRIEstimatedTimetableSubscriptionBroadcasterFactory struct{}

func (factory *SIRIEstimatedTimetableSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIEstimatedTimeTableSubscriptionBroadcaster(partner)
}

func (factory *SIRIEstimatedTimetableSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func newSIRIEstimatedTimeTableSubscriptionBroadcaster(partner *Partner) *SIRIEstimatedTimeTableSubscriptionBroadcaster {
	connector := &SIRIEstimatedTimeTableSubscriptionBroadcaster{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
	connector.vjRemoteObjectidKinds = partner.VehicleJourneyRemoteObjectIDKindWithFallback(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
	connector.partner = partner
	connector.mutex = &sync.Mutex{}
	connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	connector.estimatedTimeTableBroadcaster = NewSIRIEstimatedTimeTableBroadcaster(connector)
	return connector
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) HandleSubscriptionRequest(request *siri.XMLSubscriptionRequest, message *audit.BigQueryMessage) (resps []siri.SIRIResponseStatus) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	var lineIds, subIds []string

	for _, ett := range request.XMLSubscriptionETTEntries() {
		logStashEvent := connector.newLogStashEvent()
		logSIRIEstimatedTimeTableSubscriptionEntry(logStashEvent, ett)

		rs := siri.SIRIResponseStatus{
			RequestMessageRef: ett.MessageIdentifier(),
			SubscriberRef:     ett.SubscriberRef(),
			SubscriptionRef:   ett.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		// for logging
		lineIds = append(lineIds, ett.Lines()...)

		resources, unknownLineIds := connector.checkLines(ett)
		if len(unknownLineIds) != 0 {
			logger.Log.Debugf("EstimatedTimeTable subscription request Could not find line(s) with id : %v", strings.Join(unknownLineIds, ","))
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = fmt.Sprintf("Unknown Line(s) %v", strings.Join(unknownLineIds, ","))
		} else {
			rs.Status = true
			rs.ValidUntil = ett.InitialTerminationTime()
		}

		resps = append(resps, rs)

		logSIRIEstimatedTimeTableSubscriptionResponseEntry(logStashEvent, &rs)
		audit.CurrentLogStash().WriteEvent(logStashEvent)

		// We do not want to create a subscription that will fail
		if len(unknownLineIds) != 0 {
			continue
		}

		subIds = append(subIds, ett.SubscriptionIdentifier())

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(ett.SubscriptionIdentifier())
		if !ok {
			sub = connector.Partner().Subscriptions().New("EstimatedTimeTableBroadcast")
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

			sub.AddNewResource(*r)
		}
		sub.Save()
	}
	message.Type = "EstimatedTimetableSubscriptionRequest"
	message.SubscriptionIdentifiers = subIds
	message.Lines = lineIds

	return resps
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) addLineStopVisits(sub *Subscription, res *SubscribedResource, lineId model.LineId) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	sas := tx.Model().StopAreas().FindByLineId(lineId)
	for i := range sas {
		// Init SA LastChange
		salc := &stopAreaLastChange{}
		salc.InitState(&sas[i], sub)
		res.SetLastState(string(sas[i].Id()), salc)
		svs := tx.Model().StopVisits().FindFollowingByStopAreaId(sas[i].Id())
		for i := range svs {
			connector.addStopVisit(sub.Id(), svs[i].Id())
		}
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) checkLines(ett *siri.XMLEstimatedTimetableSubscriptionRequestEntry) (resources []*SubscribedResource, lineIds []string) {
	var lids []string
	// check for subscription to all lines
	if len(ett.Lines()) == 0 {
		var lv []string
		//find all lines corresponding to the remoteObjectidKind
		for _, line := range connector.Partner().Model().Lines().FindAll() {
			for _, objectid := range line.ObjectIDs() {
				if objectid.Kind() == connector.remoteObjectidKind {
					lv = append(lv, objectid.Value())
					break
				}
			}
		}

		for _, lineValue := range lv {
			lineObjectID := model.NewObjectID(connector.remoteObjectidKind, lineValue)
			ref := model.Reference{
				ObjectId: &lineObjectID,
				Type:     "line",
			}
			r := NewResource(ref)
			r.SubscribedAt = connector.Clock().Now()
			r.SubscribedUntil = ett.InitialTerminationTime()
			resources = append(resources, &r)
		}
		return resources, lids
	}

	for _, lineId := range ett.Lines() {

		lineObjectID := model.NewObjectID(connector.remoteObjectidKind, lineId)
		_, ok := connector.Partner().Model().Lines().FindByObjectId(lineObjectID)

		if !ok {
			lids = append(lids, lineId)
			continue
		}

		ref := model.Reference{
			ObjectId: &lineObjectID,
			Type:     "Line",
		}

		r := NewResource(ref)
		r.SubscribedAt = connector.Clock().Now()
		r.SubscribedUntil = ett.InitialTerminationTime()
		resources = append(resources, &r)
	}
	return resources, lids
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) Stop() {
	connector.estimatedTimeTableBroadcaster.Stop()
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) Start() {
	connector.estimatedTimeTableBroadcaster.Start()
}

func (ettb *SIRIEstimatedTimeTableSubscriptionBroadcaster) fillOptions(s *Subscription, request *siri.XMLSubscriptionRequest) {
	changeBeforeUpdates := request.ChangeBeforeUpdates()
	if changeBeforeUpdates == "" {
		changeBeforeUpdates = "PT1M"
	}
	s.SetSubscriptionOption("ChangeBeforeUpdates", changeBeforeUpdates)
	s.SetSubscriptionOption("MessageIdentifier", request.MessageIdentifier())
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) HandleBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	switch event.ModelType {
	case "StopVisit":
		connector.checkEvent(model.StopVisitId(event.ModelId), tx)
	case "StopArea":
		sa, ok := tx.Model().StopAreas().Find(model.StopAreaId(event.ModelId))
		if ok {
			connector.checkStopAreaEvent(&sa, tx)
		}
	default:
		return
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) checkEvent(svId model.StopVisitId, tx *model.Transaction) {
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

	subs := connector.Partner().Subscriptions().FindByResourceId(lineObj.String(), "EstimatedTimeTableBroadcast")

	for _, sub := range subs {
		r := sub.Resource(lineObj)
		if r == nil || r.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := r.LastState(string(svId))
		if ok && !lastState.(*estimatedTimeTableLastChange).Haschanged(&sv) {
			continue
		}

		if !ok {
			ettlc := &estimatedTimeTableLastChange{}
			ettlc.InitState(&sv, sub)
			r.SetLastState(string(sv.Id()), ettlc)
		}

		connector.addStopVisit(sub.Id(), svId)
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) checkStopAreaEvent(stopArea *model.StopArea, tx *model.Transaction) {
	obj, ok := stopArea.ObjectID(connector.remoteObjectidKind)
	if !ok {
		return
	}

	connector.mutex.Lock()

	subs := connector.partner.Subscriptions().FindByResourceId(obj.String(), "EstimatedTimeTableBroadcast")
	for _, sub := range subs {
		resource := sub.Resource(obj)
		if resource == nil || resource.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := resource.LastState(string(stopArea.Id()))
		if ok {
			partners, ok := lastState.(*stopAreaLastChange).Haschanged(stopArea)
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
			lastState.(*stopAreaLastChange).UpdateState(stopArea)
		} else { // Should not happen
			salc := &stopAreaLastChange{}
			salc.InitState(stopArea, sub)
			resource.SetLastState(string(stopArea.Id()), salc)
		}
	}

	connector.mutex.Unlock()
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "EstimatedTimeTableSubscriptionBroadcaster"
	return event
}

func logSIRIEstimatedTimeTableSubscriptionEntry(logStashEvent audit.LogStashEvent, ettEntry *siri.XMLEstimatedTimetableSubscriptionRequestEntry) {
	logStashEvent["siriType"] = "EstimatedTimeTableSubscriptionEntry"
	logStashEvent["lineRefs"] = strings.Join(ettEntry.Lines(), ",")
	logStashEvent["messageIdentifier"] = ettEntry.MessageIdentifier()
	logStashEvent["subscriberRef"] = ettEntry.SubscriberRef()
	logStashEvent["subscriptionIdentifier"] = ettEntry.SubscriptionIdentifier()
	logStashEvent["initialTerminationTime"] = ettEntry.InitialTerminationTime().String()
	logStashEvent["requestTimestamp"] = ettEntry.RequestTimestamp().String()
	logStashEvent["requestXML"] = ettEntry.RawXML()
}

func logSIRIEstimatedTimeTableSubscriptionResponseEntry(logStashEvent audit.LogStashEvent, response *siri.SIRIResponseStatus) {
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["subscriptionRef"] = response.SubscriptionRef
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["validUntil"] = response.ValidUntil.String()
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	if !response.Status {
		logStashEvent["errorType"] = response.ErrorType
		if response.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		}
		logStashEvent["errorText"] = response.ErrorText
	}
}

// START TEST

type TestSIRIETTSubscriptionBroadcasterFactory struct{}

type TestETTSubscriptionBroadcaster struct {
	uuid.UUIDConsumer

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
