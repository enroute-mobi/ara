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

	//stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
	toBroadcast map[SubscriptionId][]model.StopVisitId
	mutex       *sync.Mutex //protect the map
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
		}
		resps = append(resps, rs)
	}
	return resps
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
	return SubscriptionId(0), true
}
