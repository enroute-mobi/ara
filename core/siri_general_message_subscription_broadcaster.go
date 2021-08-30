package core

import (
	"strconv"
	"strings"
	"sync"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type SIRIGeneralMessageSubscriptionBroadcaster struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	generalMessageBroadcaster SIRIGeneralMessageBroadcaster
	toBroadcast               map[SubscriptionId][]model.SituationId
	mutex                     *sync.Mutex //protect the map
}

type SIRIGeneralMessageSubscriptionBroadcasterFactory struct{}

func (factory *SIRIGeneralMessageSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIGeneralMessageSubscriptionBroadcaster(partner)
}

func (factory *SIRIGeneralMessageSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	apiPartner.ValidatePresenceOfSetting(REMOTE_URL)
	apiPartner.ValidatePresenceOfSetting(REMOTE_CREDENTIAL)
	apiPartner.ValidatePresenceOfSetting(LOCAL_CREDENTIAL)
}

func newSIRIGeneralMessageSubscriptionBroadcaster(partner *Partner) *SIRIGeneralMessageSubscriptionBroadcaster {
	siriGeneralMessageSubscriptionBroadcaster := &SIRIGeneralMessageSubscriptionBroadcaster{}
	siriGeneralMessageSubscriptionBroadcaster.partner = partner
	siriGeneralMessageSubscriptionBroadcaster.mutex = &sync.Mutex{}
	siriGeneralMessageSubscriptionBroadcaster.toBroadcast = make(map[SubscriptionId][]model.SituationId)

	siriGeneralMessageSubscriptionBroadcaster.generalMessageBroadcaster = NewSIRIGeneralMessageBroadcaster(siriGeneralMessageSubscriptionBroadcaster)

	return siriGeneralMessageSubscriptionBroadcaster
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) Stop() {
	if connector.generalMessageBroadcaster != nil {
		connector.generalMessageBroadcaster.Stop()
	}
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) Start() {
	if connector.generalMessageBroadcaster == nil {
		connector.generalMessageBroadcaster = NewSIRIGeneralMessageBroadcaster(connector)
	}
	connector.generalMessageBroadcaster.Start()
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) HandleGeneralMessageBroadcastEvent(event *model.GeneralMessageBroadcastEvent) {
	connector.checkEvent(event.SituationId)
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) addSituation(subId SubscriptionId, svId model.SituationId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) checkEvent(sId model.SituationId) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	situation, ok := tx.Model().Situations().Find(sId)
	if !ok || situation.Origin == string(connector.partner.Slug()) {
		return
	}

	obj := model.NewObjectID("SituationResource", "Situation")
	subs := connector.partner.Subscriptions().FindSubscriptionsByKind("GeneralMessageBroadcast")

	for _, sub := range subs {
		resource := sub.Resource(obj)
		if resource == nil || resource.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := resource.LastState(string(situation.Id()))

		if ok && !lastState.(*generalMessageLastChange).Haschanged(&situation) {
			continue
		}

		if !ok {
			gmlc := &generalMessageLastChange{}
			gmlc.InitState(&situation, sub)
			resource.SetLastState(string(situation.Id()), gmlc)
		}
		connector.addSituation(sub.Id(), sId)
	}
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) HandleSubscriptionRequest(request *siri.XMLSubscriptionRequest, message *audit.BigQueryMessage) []siri.SIRIResponseStatus {
	resps := []siri.SIRIResponseStatus{}

	var subIds []string

	for _, gm := range request.XMLSubscriptionGMEntries() {
		logStashEvent := connector.newLogStashEvent()
		logXMLGeneralMessageSubscriptionEntry(logStashEvent, gm)

		rs := siri.SIRIResponseStatus{
			RequestMessageRef: gm.MessageIdentifier(),
			SubscriberRef:     gm.SubscriberRef(),
			SubscriptionRef:   gm.SubscriptionIdentifier(),
			Status:            true,
			ResponseTimestamp: connector.Clock().Now(),
			ValidUntil:        gm.InitialTerminationTime(),
		}

		subIds = append(subIds, gm.SubscriptionIdentifier())

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(gm.SubscriptionIdentifier())
		if !ok {
			sub = connector.Partner().Subscriptions().New("GeneralMessageBroadcast")
			sub.SubscriberRef = gm.SubscriberRef()
			sub.SetExternalId(gm.SubscriptionIdentifier())
		}

		sub.SetSubscriptionOption("InfoChannelRef", strings.Join(gm.InfoChannelRef(), ","))
		sub.SetSubscriptionOption("LineRef", strings.Join(gm.LineRef(), ","))
		sub.SetSubscriptionOption("StopPointRef", strings.Join(gm.StopPointRef(), ","))
		sub.SetSubscriptionOption("MessageIdentifier", gm.MessageIdentifier())

		obj := model.NewObjectID("SituationResource", "Situation")
		r := sub.Resource(obj)
		if r == nil {
			ref := model.Reference{
				ObjectId: &obj,
				Type:     "Situation",
			}
			r = sub.CreateAddNewResource(ref)
			r.SubscribedAt = connector.Clock().Now()
			r.SubscribedUntil = gm.InitialTerminationTime()
		}

		sub.Save()

		connector.addSituations(sub, r)
		resps = append(resps, rs)

		logSIRIGeneralMessageSubscriptionResponseEntry(logStashEvent, &rs)
		audit.CurrentLogStash().WriteEvent(logStashEvent)
	}

	message.Type = "StopMonitoringSubscriptionRequest"
	message.SubscriptionIdentifiers = subIds

	return resps
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) addSituations(sub *Subscription, r *SubscribedResource) {
	for _, situation := range connector.partner.Model().Situations().FindAll() {
		if situation.ValidUntil.Before(connector.Clock().Now()) {
			continue
		}

		gmlc := &generalMessageLastChange{}
		gmlc.InitState(&situation, sub)
		r.SetLastState(string(situation.Id()), gmlc)
		connector.addSituation(sub.Id(), situation.Id())
	}
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "GeneralMessageSubscriptionBroadcaster"
	return event
}

func logXMLGeneralMessageSubscriptionEntry(logStashEvent audit.LogStashEvent, request *siri.XMLGeneralMessageSubscriptionRequestEntry) {
	logStashEvent["siriType"] = "GeneralMessageSubscriptionEntry"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["infoChannelRefs"] = strings.Join(request.InfoChannelRef(), ", ")
	logStashEvent["stopPointRefs"] = strings.Join(request.StopPointRef(), ", ")
	logStashEvent["lineRefs"] = strings.Join(request.LineRef(), ", ")
	logStashEvent["subscriberRef"] = request.SubscriberRef()
	logStashEvent["subscriptionIdentifier"] = request.SubscriptionIdentifier()
	logStashEvent["initialTerminationTime"] = request.InitialTerminationTime().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIGeneralMessageSubscriptionResponseEntry(logStashEvent audit.LogStashEvent, gmEntry *siri.SIRIResponseStatus) {
	logStashEvent["requestMessageRef"] = gmEntry.RequestMessageRef
	logStashEvent["subscriptionRef"] = gmEntry.SubscriptionRef
	logStashEvent["responseTimestamp"] = gmEntry.ResponseTimestamp.String()
	logStashEvent["validUntil"] = gmEntry.ValidUntil.String()
	logStashEvent["status"] = strconv.FormatBool(gmEntry.Status)
	if !gmEntry.Status {
		logStashEvent["errorType"] = gmEntry.ErrorType
		if gmEntry.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(gmEntry.ErrorNumber)
		}
		logStashEvent["errorText"] = gmEntry.ErrorText
	}
}

// Start Test

type TestSIRIGeneralMessageSubscriptionBroadcasterFactory struct{}

type TestGeneralMessageSubscriptionBroadcaster struct {
	uuid.UUIDConsumer

	events []*model.GeneralMessageBroadcastEvent
	// generalMessageBroadcaster SIRIGeneralMessageBroadcaster
}

func NewTestGeneralMessageSubscriptionBroadcaster() *TestGeneralMessageSubscriptionBroadcaster {
	connector := &TestGeneralMessageSubscriptionBroadcaster{}
	return connector
}

func (connector *TestGeneralMessageSubscriptionBroadcaster) HandleGeneralMessageBroadcastEvent(event *model.GeneralMessageBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIGeneralMessageSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
} // Always valid

func (factory *TestSIRIGeneralMessageSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestGeneralMessageSubscriptionBroadcaster()
}

// END OF TEST
