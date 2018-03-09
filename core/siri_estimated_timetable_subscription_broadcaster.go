package core

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type EstimatedTimeTableSubscriptionBroadcaster interface {
	model.Stopable
	model.Startable

	HandleStopMonitoringBroadcastEvent(*model.StopMonitoringBroadcastEvent)
	HandleSubscriptionRequest([]*siri.XMLEstimatedTimetableSubscriptionRequestEntry) []siri.SIRIResponseStatus
}

type SIRIEstimatedTimeTableSubscriptionBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

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

func (factory *SIRIEstimatedTimetableSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func newSIRIEstimatedTimeTableSubscriptionBroadcaster(partner *Partner) *SIRIEstimatedTimeTableSubscriptionBroadcaster {
	siriEstimatedTimeTableSubscriptionBroadcaster := &SIRIEstimatedTimeTableSubscriptionBroadcaster{}
	siriEstimatedTimeTableSubscriptionBroadcaster.partner = partner
	siriEstimatedTimeTableSubscriptionBroadcaster.mutex = &sync.Mutex{}
	siriEstimatedTimeTableSubscriptionBroadcaster.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	siriEstimatedTimeTableSubscriptionBroadcaster.estimatedTimeTableBroadcaster = NewSIRIEstimatedTimeTableBroadcaster(siriEstimatedTimeTableSubscriptionBroadcaster)
	return siriEstimatedTimeTableSubscriptionBroadcaster
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) HandleSubscriptionRequest(request *siri.XMLSubscriptionRequest) (resps []siri.SIRIResponseStatus) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	for _, ett := range request.XMLSubscriptionETTEntries() {
		logStashEvent := connector.newLogStashEvent()
		logSIRIEstimatedTimeTableSubscriptionEntry(logStashEvent, ett)

		rs := siri.SIRIResponseStatus{
			RequestMessageRef: ett.MessageIdentifier(),
			SubscriberRef:     ett.SubscriberRef(),
			SubscriptionRef:   ett.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		resources, lineIds := connector.checkLines(ett)
		if len(lineIds) != 0 {
			logger.Log.Debugf("EstimatedTimeTable subscription request Could not find line(s) with id : %v", strings.Join(lineIds, ","))
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = fmt.Sprintf("Unknown Line(s) %v", strings.Join(lineIds, ","))
		} else {
			rs.Status = true
			rs.ValidUntil = ett.InitialTerminationTime()
		}

		resps = append(resps, rs)

		logSIRIEstimatedTimeTableSubscriptionResponseEntry(logStashEvent, &rs)
		audit.CurrentLogStash().WriteEvent(logStashEvent)

		if len(lineIds) != 0 {
			continue
		}

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(ett.SubscriptionIdentifier())
		if !ok {
			sub = connector.Partner().Subscriptions().New("EstimatedTimeTableBroadcast")
			sub.SetExternalId(ett.SubscriptionIdentifier())
			connector.fillOptions(sub, request)
		}

		for _, r := range resources {
			line, ok := connector.Partner().Model().Lines().FindByObjectId(*r.Reference.ObjectId)
			if !ok {
				continue
			}

			sub.AddNewResource(r)
			connector.addLineStopVisits(sub.Id(), line.Id())
		}
		sub.Save()
	}
	return resps
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) addLineStopVisits(subId SubscriptionId, lineId model.LineId) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	for _, sa := range tx.Model().StopAreas().FindByLineId(lineId) {
		for _, sv := range tx.Model().StopVisits().FindFollowingByStopAreaId(sa.Id()) {
			connector.addStopVisit(subId, sv.Id())
		}
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) checkLines(ett *siri.XMLEstimatedTimetableSubscriptionRequestEntry) (resources []SubscribedResource, lineIds []string) {
	for _, lineId := range ett.Lines() {
		lineObjectId := model.NewObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER), lineId)
		_, ok := connector.Partner().Model().Lines().FindByObjectId(lineObjectId)

		if !ok {
			lineIds = append(lineIds, lineId)
			continue
		}

		ref := model.Reference{
			ObjectId: &lineObjectId,
			Type:     "Line",
		}

		r := NewResource(ref)
		r.SubscribedAt = connector.Clock().Now()
		r.SubscribedUntil = ett.InitialTerminationTime()
		resources = append(resources, r)
	}
	return resources, lineIds
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) Stop() {
	connector.estimatedTimeTableBroadcaster.Stop()
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) Start() {
	connector.estimatedTimeTableBroadcaster.Start()
}

func (ettb *SIRIEstimatedTimeTableSubscriptionBroadcaster) fillOptions(s *Subscription, request *siri.XMLSubscriptionRequest) {
	so := s.SubscriptionOptions()
	changeBeforeUpdates := request.ChangeBeforeUpdates()
	if changeBeforeUpdates == "" {
		changeBeforeUpdates = "PT1M"
	}
	so["ChangeBeforeUpdates"] = changeBeforeUpdates
	so["MessageIdentifier"] = request.MessageIdentifier()
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
			subsIds := connector.checkStopAreaEvent(sa, tx)
			if len(subsIds) != 0 {
				connector.addProducerUnavailable(subsIds, sa.Origin)
			}
		}
	default:
		return
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) addProducerUnavailable(subsIds []SubscriptionId, producer string) {
	connector.mutex.Lock()
	for _, subId := range subsIds {
		nm, ok := connector.notMonitored[SubscriptionId(subId)]
		if !ok {
			nm = make(map[string]struct{})
			connector.notMonitored[SubscriptionId(subId)] = nm
		}
		nm[producer] = struct{}{}
	}
	connector.mutex.Unlock()
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

	lineObj, ok := line.ObjectID(connector.Partner().RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER))
	if !ok {
		return
	}

	subs := connector.Partner().Subscriptions().FindByRessourceId(lineObj.String(), "EstimatedTimeTableBroadcast")

	for _, sub := range subs {
		r := sub.Resource(lineObj)
		if r == nil || r.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := r.LastStates[string(svId)]
		if ok && !lastState.(*estimatedTimeTableLastChange).Haschanged(&sv) {
			continue
		}

		if !ok {
			ettlc := &estimatedTimeTableLastChange{}
			ettlc.InitState(&sv, sub)
			r.LastStates[string(sv.Id())] = ettlc
		}

		connector.addStopVisit(sub.Id(), svId)
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) checkStopAreaEvent(stopArea model.StopArea, tx *model.Transaction) []SubscriptionId {
	subscriptionIds := []SubscriptionId{}

	obj, ok := stopArea.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER))
	if !ok {
		return subscriptionIds
	}

	subs := connector.partner.Subscriptions().FindByRessourceId(obj.String(), "EstimatedTimeTableBroadcast")

	for _, sub := range subs {
		resource, ok := sub.ResourcesByObjectID()[obj.String()]
		if !ok || resource.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := resource.LastStates[string(stopArea.Id())]
		if ok {
			if lastState.(*stopAreaLastChange).Haschanged(stopArea) && !stopArea.Monitored {
				subscriptionIds = append(subscriptionIds, sub.Id())
			}
			lastState.(*stopAreaLastChange).UpdateState(&stopArea)
		}

		if !ok {
			salc := &stopAreaLastChange{}
			salc.InitState(&stopArea, sub)
			resource.LastStates[string(stopArea.Id())] = salc
		}
	}

	return subscriptionIds
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "EstimatedTimeTableSubscriptionBroadcaster"
	return event
}

func logSIRIEstimatedTimeTableSubscriptionEntry(logStashEvent audit.LogStashEvent, ettEntry *siri.XMLEstimatedTimetableSubscriptionRequestEntry) {
	logStashEvent["type"] = "EstimatedTimeTableSubscriptionEntry"
	logStashEvent["LineRef"] = strings.Join(ettEntry.Lines(), ",")
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
	model.UUIDConsumer

	events []*model.StopMonitoringBroadcastEvent
}

func NewTestETTSubscriptionBroadcaster() *TestETTSubscriptionBroadcaster {
	connector := &TestETTSubscriptionBroadcaster{}
	return connector
}

func (connector *TestETTSubscriptionBroadcaster) HandleBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIETTSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestSIRIETTSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestETTSubscriptionBroadcaster()
}
