package core

import (
	"strconv"
	"strings"
	"sync"
	"time"

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
	toBroadcast                   map[SubscriptionId][]model.LineId
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
	siriEstimatedTimeTableSubscriptionBroadcaster.toBroadcast = make(map[SubscriptionId][]model.LineId)

	siriEstimatedTimeTableSubscriptionBroadcaster.estimatedTimeTableBroadcaster = NewSIRIEstimatedTimeTableBroadcaster(siriEstimatedTimeTableSubscriptionBroadcaster)
	return siriEstimatedTimeTableSubscriptionBroadcaster
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) HandleSubscriptionRequest(request *siri.XMLSubscriptionRequest) []siri.SIRIResponseStatus {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logSIRIEstimatedTimeTableBroadcasterSubscriptionRequest(logStashEvent, request)

	ettEntries := request.XMLSubscriptionETTEntries()

	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	resps := []siri.SIRIResponseStatus{}

	for _, ett := range ettEntries {
		logStashEvent := connector.newLogStashEvent()
		logSIRIEstimatedTimeTableBroadcasterEntries(logStashEvent, ett)
		audit.CurrentLogStash().WriteEvent(logStashEvent)

		resources := []SubscribedResource{}
		rs := siri.SIRIResponseStatus{
			Status:            true,
			RequestMessageRef: ett.MessageIdentifier(),
			SubscriberRef:     ett.SubscriberRef(),
			SubscriptionRef:   ett.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
			ValidUntil:        ett.InitialTerminationTime(),
		}

		for _, lineId := range ett.Lines() {
			lineObjectId := model.NewObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER), lineId)
			_, ok := connector.Partner().Model().Lines().FindByObjectId(lineObjectId)

			if !ok {
				rs.Status = false
				rs.ErrorType = "InvalidDataReferencesError"
				rs.ErrorText = "Could not find Line " + lineId
				rs.ValidUntil = time.Time{}

				resources = []SubscribedResource{}
				resps = append(resps, rs)

				logSIRIEstimatedTimeTableBroadcasterSubscriptionResponse(logStashEvent, rs)
				audit.CurrentLogStash().WriteEvent(logStashEvent)
				logger.Log.Debugf("EstimatedTimeTable subscription request Could not find line with id : %v", lineId)

				break
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
		resps = append(resps, rs)

		logSIRIEstimatedTimeTableBroadcasterSubscriptionResponse(logStashEvent, rs)
		audit.CurrentLogStash().WriteEvent(logStashEvent)

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(ett.SubscriptionIdentifier())
		if !ok {
			sub = connector.Partner().Subscriptions().New("EstimatedTimeTableBroadcast")
			sub.SetExternalId(ett.SubscriptionIdentifier())
		}

		for _, r := range resources {
			line, _ := connector.Partner().Model().Lines().FindByObjectId(*r.Reference.ObjectId)
			connector.addLine(sub.Id(), line.Id())

			sub.AddNewResource(r)
			connector.fillOptions(sub, request)
		}

		resources = []SubscribedResource{}
		sub.Save()
	}
	return resps
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

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) HandleStopVisitBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	if event.ModelType != "StopVisit" {
		return
	}

	connector.checkEvent(model.StopVisitId(event.ModelId), tx)
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) addLine(subId SubscriptionId, lineId model.LineId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], lineId)
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
		if r == nil {
			continue
		}

		lastState, ok := r.LastStates[string(sv.Id())]
		if !ok {
			r.LastStates[string(sv.Id())] = &estimatedTimeTableLastChange{}
			lastState, _ = r.LastStates[string(sv.Id())]
			lastState.(*estimatedTimeTableLastChange).InitState(&sv, sub)
		}

		if ok = lastState.(*estimatedTimeTableLastChange).Haschanged(&sv); !ok {
			continue
		}
		connector.addLine(sub.Id(), line.Id())
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "EstimatedTimeTableSubscriptionBroadcaster"
	return event
}

func logSIRIEstimatedTimeTableBroadcasterSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.XMLSubscriptionRequest) {
	logStashEvent["type"] = "EstimatedTimeTableSubscriptions"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()

	xml := request.RawXML()
	logStashEvent["requestXML"] = xml
}

func logSIRIEstimatedTimeTableBroadcasterSubscriptionResponse(logStashEvent audit.LogStashEvent, response siri.SIRIResponseStatus) {
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

func logSIRIEstimatedTimeTableBroadcasterEntries(logStashEvent audit.LogStashEvent, ettEntries *siri.XMLEstimatedTimetableSubscriptionRequestEntry) {
	logStashEvent["type"] = "EstimatedTimeTableSubscription"
	logStashEvent["subscriberRef"] = ettEntries.SubscriberRef()
	logStashEvent["subscriptionRef"] = ettEntries.SubscriptionIdentifier()
	logStashEvent["LineRef"] = strings.Join(ettEntries.Lines(), ",")

	xml := ettEntries.RawXML()
	logStashEvent["requestXML"] = xml
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
