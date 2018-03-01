package core

import (
	"strconv"
	"sync"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringSubscriptionBroadcaster interface {
	model.Stopable
	model.Startable

	HandleStopMonitoringBroadcastEvent(*model.StopMonitoringBroadcastEvent)
	HandleSubscriptionRequest([]*siri.XMLStopMonitoringSubscriptionRequestEntry)
}

type SIRIStopMonitoringSubscriptionBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
	toBroadcast               map[SubscriptionId][]model.StopVisitId

	mutex *sync.Mutex //protect the map
}

type SIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

func (factory *SIRIStopMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIStopMonitoringSubscriptionBroadcaster(partner)
}

func (factory *SIRIStopMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func newSIRIStopMonitoringSubscriptionBroadcaster(partner *Partner) *SIRIStopMonitoringSubscriptionBroadcaster {
	siriStopMonitoringSubscriptionBroadcaster := &SIRIStopMonitoringSubscriptionBroadcaster{}
	siriStopMonitoringSubscriptionBroadcaster.partner = partner
	siriStopMonitoringSubscriptionBroadcaster.mutex = &sync.Mutex{}
	siriStopMonitoringSubscriptionBroadcaster.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	siriStopMonitoringSubscriptionBroadcaster.stopMonitoringBroadcaster = NewSIRIStopMonitoringBroadcaster(siriStopMonitoringSubscriptionBroadcaster)

	return siriStopMonitoringSubscriptionBroadcaster
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) Stop() {
	connector.stopMonitoringBroadcaster.Stop()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) Start() {
	connector.stopMonitoringBroadcaster.Start()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) HandleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	switch event.ModelType {
	case "StopVisit":
		sv, ok := tx.Model().StopVisits().Find(model.StopVisitId(event.ModelId))
		if ok {
			subsIds := connector.checkEvent(sv, tx)
			connector.addStopVisit(subsIds, sv.Id())
		}
	case "VehicleJourney":
		for _, sv := range tx.Model().StopVisits().FindFollowingByVehicleJourneyId(model.VehicleJourneyId(event.ModelId)) {
			subsIds := connector.checkEvent(sv, tx)
			connector.addStopVisit(subsIds, sv.Id())
		}
	default:
		return
	}
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) addStopVisit(subsIds []SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	for _, subId := range subsIds {
		connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	}
	connector.mutex.Unlock()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) checkEvent(sv model.StopVisit, tx *model.Transaction) []SubscriptionId {
	subscriptionIds := []SubscriptionId{}

	if sv.Origin == string(connector.Partner().Slug()) {
		return subscriptionIds
	}

	stopArea, ok := tx.Model().StopAreas().Find(sv.StopAreaId)
	if !ok {
		return subscriptionIds
	}

	obj, ok := stopArea.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER))
	if !ok {
		return subscriptionIds
	}

	subs := connector.partner.Subscriptions().FindByRessourceId(obj.String(), "StopMonitoringBroadcast")

	for _, sub := range subs {
		resource, ok := sub.ResourcesByObjectID()[obj.String()]
		if !ok || resource.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := resource.LastStates[string(sv.Id())]
		if ok && !lastState.(*stopMonitoringLastChange).Haschanged(sv) {
			continue
		}

		if !ok {
			smlc := &stopMonitoringLastChange{}
			smlc.InitState(&sv, sub)
			resource.LastStates[string(sv.Id())] = smlc
		}

		subscriptionIds = append(subscriptionIds, sub.Id())
	}

	return subscriptionIds
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) HandleSubscriptionRequest(request *siri.XMLSubscriptionRequest) []siri.SIRIResponseStatus {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	resps := []siri.SIRIResponseStatus{}

	for _, sm := range request.XMLSubscriptionSMEntries() {
		logStashEvent := connector.newLogStashEvent()
		logXMLStopMonitoringSubscriptionEntry(logStashEvent, sm)

		rs := siri.SIRIResponseStatus{
			RequestMessageRef: sm.MessageIdentifier(),
			SubscriberRef:     sm.SubscriberRef(),
			SubscriptionRef:   sm.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		objectid := model.NewObjectID(connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER), sm.MonitoringRef())
		sa, ok := tx.Model().StopAreas().FindByObjectId(objectid)
		if !ok {
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = "Stop Area not found"
			resps = append(resps, rs)

			logSIRIStopMonitoringSubscriptionResponseEntry(logStashEvent, &rs)
			audit.CurrentLogStash().WriteEvent(logStashEvent)

			continue
		}

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(sm.SubscriptionIdentifier())
		if !ok {
			sub = connector.Partner().Subscriptions().New("StopMonitoringBroadcast")
			sub.SetExternalId(sm.SubscriptionIdentifier())
		}

		ref := model.Reference{
			ObjectId: &objectid,
			Type:     "StopArea",
		}

		r := sub.CreateAddNewResource(ref)
		r.SubscribedAt = connector.Clock().Now()
		r.SubscribedUntil = sm.InitialTerminationTime()

		connector.fillOptions(sub, r, request, sm)

		rs.Status = true
		rs.ValidUntil = sm.InitialTerminationTime()
		resps = append(resps, rs)

		logSIRIStopMonitoringSubscriptionResponseEntry(logStashEvent, &rs)
		audit.CurrentLogStash().WriteEvent(logStashEvent)

		connector.addStopAreaStopVisits(sa, sub, r)
	}

	return resps
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) addStopAreaStopVisits(sa model.StopArea, sub *Subscription, res *SubscribedResource) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	for _, sv := range tx.Model().StopVisits().FindFollowingByStopAreaId(sa.Id()) {
		if _, ok := res.LastStates[string(sv.Id())]; ok {
			continue
		}
		smlc := &stopMonitoringLastChange{}
		smlc.InitState(&sv, sub)
		res.LastStates[string(sv.Id())] = smlc
		connector.addStopVisit([]SubscriptionId{sub.Id()}, sv.Id())
	}
}

func (smsb *SIRIStopMonitoringSubscriptionBroadcaster) fillOptions(s *Subscription, r *SubscribedResource, request *siri.XMLSubscriptionRequest, sm *siri.XMLStopMonitoringSubscriptionRequestEntry) {
	ro := r.ResourcesOptions()
	ro["StopVisitTypes"] = sm.StopVisitTypes()
	changeBeforeUpdates := request.ChangeBeforeUpdates()
	if changeBeforeUpdates == "" {
		changeBeforeUpdates = "PT1M"
	}

	so := s.SubscriptionOptions()

	so["IncrementalUpdates"] = request.IncrementalUpdates()
	so["MaximumStopVisits"] = request.MaximumStopVisits()
	so["ChangeBeforeUpdates"] = changeBeforeUpdates
	so["MessageIdentifier"] = request.MessageIdentifier()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringSubscriptionBroadcaster"
	return event
}

func logXMLStopMonitoringSubscriptionEntry(logStashEvent audit.LogStashEvent, request *siri.XMLStopMonitoringSubscriptionRequestEntry) {
	logStashEvent["type"] = "StopMonitoringSubscriptionEntry"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["monitoringRef"] = request.MonitoringRef()
	logStashEvent["stopVisitTypes"] = request.StopVisitTypes()
	logStashEvent["subscriberRef"] = request.SubscriberRef()
	logStashEvent["subscriptionIdentifier"] = request.SubscriptionIdentifier()
	logStashEvent["initialTerminationTime"] = request.InitialTerminationTime().String()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIStopMonitoringSubscriptionResponseEntry(logStashEvent audit.LogStashEvent, smEntry *siri.SIRIResponseStatus) {
	logStashEvent["requestMessageRef"] = smEntry.RequestMessageRef
	logStashEvent["subscriptionRef"] = smEntry.SubscriptionRef
	logStashEvent["responseTimestamp"] = smEntry.ResponseTimestamp.String()
	logStashEvent["validUntil"] = smEntry.ValidUntil.String()
	logStashEvent["status"] = strconv.FormatBool(smEntry.Status)
	if !smEntry.Status {
		logStashEvent["errorType"] = smEntry.ErrorType
		if smEntry.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(smEntry.ErrorNumber)
		}
		logStashEvent["errorText"] = smEntry.ErrorText
	}
}

// START TEST

type TestSIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

type TestStopMonitoringSubscriptionBroadcaster struct {
	model.UUIDConsumer

	events                    []*model.StopMonitoringBroadcastEvent
	stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
}

func NewTestStopMonitoringSubscriptionBroadcaster() *TestStopMonitoringSubscriptionBroadcaster {
	connector := &TestStopMonitoringSubscriptionBroadcaster{}
	return connector
}

func (connector *TestStopMonitoringSubscriptionBroadcaster) HandleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIStopMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestSIRIStopMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestStopMonitoringSubscriptionBroadcaster()
}

// END TEST
