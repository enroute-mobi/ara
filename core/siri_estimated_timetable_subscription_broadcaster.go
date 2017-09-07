package core

import (
	"sync"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type EstimatedTimeTableSubscriptionBroadcaster interface {
	HandleStopMonitoringBroadcastEvent(*model.StopMonitoringBroadcastEvent)
	HandleSubscriptionRequest([]*siri.XMLEstimatedTimetableSubscriptionRequestEntry) []siri.SIRIResponseStatus
}

type SIRIEstimatedTimeTableSubscriptionBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	estimatedTimeTableBroadcaster SIRIEstimatedTimeTableBroadcaster
	toBroadcast                   map[SubscriptionId][]model.StopVisitId
	mutex                         *sync.Mutex //protect the map
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

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) HandleSubscriptionRequest(request *siri.XMLSubscriptionRequest) []siri.SIRIResponseStatus {
	ettEntries := request.XMLSubscriptionETTEntries()

	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	resps := []siri.SIRIResponseStatus{}

	for _, ett := range ettEntries {

		rs := siri.SIRIResponseStatus{
			RequestMessageRef: ett.MessageIdentifier(),
			SubscriberRef:     ett.SubscriberRef(),
			SubscriptionRef:   ett.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
			Status:            true,
		}

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(ett.SubscriptionIdentifier())
		if !ok {
			sub = connector.Partner().Subscriptions().New("EstimatedTimeTable")
			sub.SetExternalId(ett.SubscriptionIdentifier())
		}

		for _, lineId := range ett.Lines() {
			_, ok := connector.Partner().Model().Lines().Find(model.LineId(lineId))
			if !ok {
				logger.Log.Debugf("EstimatedTimeTable subscription request Could not find line with id : %v", lineId)
				continue
			}
			lineObjectId := model.NewObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER), lineId)
			ref := model.Reference{
				ObjectId: &lineObjectId,
				Id:       lineId,
				Type:     "line",
			}

			r := sub.CreateAddNewResource(ref)
			r.SubscribedUntil = ett.InitialTerminationTime()
			connector.fillOptions(sub, request)

			rs.ValidUntil = ett.InitialTerminationTime()
		}
		resps = append(resps, rs)
	}
	return resps
}

func (ettb *SIRIEstimatedTimeTableSubscriptionBroadcaster) fillOptions(s *Subscription, request *siri.XMLSubscriptionRequest) {
	so := s.SubscriptionOptions()
	so["ChangeBeforeUpdates"] = request.ChangeBeforeUpdates()
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) HandleStopVisitBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	if event.ModelType != "StopVisit" {
		return
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) addStopVisit(subId SubscriptionId, svId model.StopVisitId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) checkEvent(sv model.StopVisit, tx *model.Transaction) (SubscriptionId, bool) {
	subId := SubscriptionId(0)

	vj, ok := connector.Partner().Model().VehicleJourneys().Find(sv.VehicleJourneyId)
	if !ok {
		return subId, false
	}

	line, ok := connector.Partner().Model().Lines().Find(vj.LineId)
	if !ok {
		return subId, false
	}

	lineObj, ok := line.ObjectID(connector.Partner().RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER))
	if !ok {
		return subId, false
	}

	sub, ok := connector.Partner().Subscriptions().FindByRessourceId(lineObj.String())
	if !ok {
		return subId, false
	}

	r := sub.Resource(lineObj)
	if r == nil {
		return subId, false
	}

	lastState, ok := r.LastStates[string(sv.Id())]

	if !ok {
		lastState.(*estimatedTimeTable).InitState(&sv, sub)
	}

	return subId, lastState.(*estimatedTimeTable).Haschanged(&sv)
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

func (connector *TestETTSubscriptionBroadcaster) HandleStopVisitBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIETTSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestSIRIETTSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestETTSubscriptionBroadcaster()
}
