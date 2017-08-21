package core

import (
	"sync"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringSubscriptionBroadcaster interface {
	model.Stopable
	model.Startable

	handleStopMonitoringBroadcastEvent(*model.StopMonitoringBroadcastEvent)
	HandleSubscriptionRequest([]*siri.XMLStopMonitoringSubscriptionRequestEntry)
}

type SIRIStopMonitoringSubscriptionBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
	toBroadcast               map[SubscriptionId][]model.StopVisitId
	mutex                     *sync.Mutex //protect the map
}

type SIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

func (factory *SIRIStopMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return newSIRIStopMonitoringSubscriptionBroadcaster(partner)
}

func (factory *SIRIStopMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
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

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) RemoteObjectIDKind() string {
	if connector.partner.Setting("siri-stop-monitoring-subscription-broadcaster.remote_objectid_kind") != "" {
		return connector.partner.Setting("siri-stop-monitoring-subscription-broadcaster.remote_objectid_kind")
	}
	return connector.partner.Setting("remote_objectid_kind")
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) Stop() {
	connector.stopMonitoringBroadcaster.Stop()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) Start() {
	connector.stopMonitoringBroadcaster.Start()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) handleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	switch event.ModelType {
	case "StopVisit":
		sv, ok := connector.Partner().Model().StopVisits().Find(model.StopVisitId(event.ModelId))
		subId, ok := connector.checkEvent(sv)
		if ok {
			connector.addStopVisit(subId, sv.Id())
		}
	case "VehicleJourney":
		for _, sv := range connector.Partner().Model().StopVisits().FindFollowingByVehicleJourneyId(model.VehicleJourneyId(event.ModelId)) {
			subId, ok := connector.checkEvent(sv)
			if ok {
				connector.addStopVisit(subId, sv.Id())
			}
		}
	default:
		return
	}
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) checkEvent(sv model.StopVisit) (SubscriptionId, bool) {
	subId := SubscriptionId(0) //just to return a correct type for errors

	stopArea, ok := connector.Partner().Model().StopAreas().Find(sv.StopAreaId)
	if !ok {
		return subId, false
	}

	obj, ok := stopArea.ObjectID(connector.RemoteObjectIDKind())
	if !ok {
		return subId, false
	}

	sub, ok := connector.partner.Subscriptions().FindByRessourceId(obj.String())
	if !ok {
		return subId, false
	}

	resources := sub.ResourcesByObjectID()

	resource, ok := resources[obj.String()]

	if !ok {
		return subId, false
	}

	lastState, ok := resource.LastStates[string(sv.Id())]

	if ok && !lastState.(*stopMonitoringLastChange).Haschanged(sv) {
		return subId, false
	}

	if !ok {
		smlc := &stopMonitoringLastChange{}
		smlc.SetSubscription(sub)
		resource.LastStates[string(sv.Id())] = smlc
	}

	return sub.Id(), true
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) HandleSubscriptionRequest(request *siri.XMLSubscriptionRequest) []siri.SIRIResponseStatus {
	sms := request.XMLSubscriptionSMEntries()

	resps := []siri.SIRIResponseStatus{}

	for _, sm := range sms {

		rs := siri.SIRIResponseStatus{
			RequestMessageRef: sm.MessageIdentifier(),
			SubscriberRef:     sm.SubscriberRef(),
			SubscriptionRef:   sm.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		objectid := model.NewObjectID(connector.RemoteObjectIDKind(), sm.MonitoringRef())
		sa, ok := connector.Partner().Model().StopAreas().FindByObjectId(objectid)
		if !ok {
			resps = append(resps, rs)
			continue
		}

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(sm.SubscriptionIdentifier())
		if !ok {
			sub = connector.Partner().Subscriptions().New("StopArea")
			sub.SetExternalId(sm.SubscriptionIdentifier())
		}

		ref := model.Reference{
			ObjectId: &objectid,
			Id:       string(sa.Id()),
			Type:     "StopArea",
		}

		r := sub.CreateAddNewResource(ref)
		r.SubscribedUntil = sm.InitialTerminationTime()

		connector.fillOptions(sub, r, request, sm)

		rs.Status = true
		rs.ValidUntil = sm.InitialTerminationTime()
		resps = append(resps, rs)

		connector.AddStopAreaStopVisits(sa, sub, r)
	}
	return resps
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) AddStopAreaStopVisits(sa model.StopArea, sub *Subscription, res *SubscribedResource) {
	svs := connector.Partner().Model().StopVisits().FindFollowingByStopAreaId(sa.Id())
	for _, sv := range svs {
		_, ok := sv.ObjectID(connector.RemoteObjectIDKind())
		if !ok {
			continue
		}
		smlc := &stopMonitoringLastChange{}
		smlc.SetSubscription(sub)
		res.LastStates[string(sv.Id())] = smlc
		connector.addStopVisit(sub.Id(), sv.Id())
	}
}

func (smsb *SIRIStopMonitoringSubscriptionBroadcaster) fillOptions(s *Subscription, r *SubscribedResource, request *siri.XMLSubscriptionRequest, sm *siri.XMLStopMonitoringSubscriptionRequestEntry) {
	ro := r.ResourcesOptions()
	ro["StopVisitTypes"] = sm.StopVisitTypes()

	so := s.SubscriptionOptions()

	so["IncrementalUpdates"] = request.IncrementalUpdates()
	so["MaximumStopVisits"] = request.MaximumStopVisits()
	so["ChangeBeforeUpdates"] = request.ChangeBeforeUpdates()
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

func (connector *TestStopMonitoringSubscriptionBroadcaster) handleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIStopMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestSIRIStopMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestStopMonitoringSubscriptionBroadcaster()
}

// END TEST
