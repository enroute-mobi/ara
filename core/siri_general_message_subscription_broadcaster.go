package core

import (
	"fmt"
	"strings"
	"sync"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIGeneralMessageSubscriptionBroadcaster struct {
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
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
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

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) HandleGeneralMessageBroadcastEvent(event *model.SituationBroadcastEvent) {
	connector.checkEvent(event.SituationId)
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) addSituation(subId SubscriptionId, svId model.SituationId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) checkEvent(sId model.SituationId) {
	situation, ok := connector.partner.Model().Situations().Find(sId)
	if !ok || situation.Origin == string(connector.partner.Slug()) {
		return
	}

	obj := model.NewCode("SituationResource", "Situation")
	subs := connector.partner.Subscriptions().FindSubscriptionsByKind(GeneralMessageBroadcast)

	for _, sub := range subs {
		resource := sub.Resource(obj)
		if resource == nil || resource.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := resource.LastState(string(situation.Id()))

		if ok && !lastState.(*ls.SituationLastChange).Haschanged(&situation) {
			continue
		}

		if !ok {
			resource.SetLastState(string(situation.Id()), ls.NewSituationLastChange(&situation, sub))
		}
		connector.addSituation(sub.Id(), sId)
	}
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) HandleSubscriptionRequest(request *sxml.XMLSubscriptionRequest, message *audit.BigQueryMessage) []siri.SIRIResponseStatus {
	resps := []siri.SIRIResponseStatus{}

	var subIds []string

	for _, gm := range request.XMLSubscriptionGMEntries() {
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
		if ok {
			if sub.Kind() != GeneralMessageBroadcast {
				logger.Log.Debugf("GeneralMessage subscription request with a duplicated Id: %v", gm.SubscriptionIdentifier())
				rs.Status = false
				rs.ErrorType = "OtherError"
				rs.ErrorNumber = 2
				rs.ErrorText = fmt.Sprintf("[BAD_REQUEST] Subscription Id %v already exists", gm.SubscriptionIdentifier())
				resps = append(resps, rs)
				message.Status = "Error"
				continue
			}

			sub.Delete()
		}

		sub = connector.Partner().Subscriptions().New(GeneralMessageBroadcast)
		sub.SubscriberRef = gm.SubscriberRef()
		sub.SetExternalId(gm.SubscriptionIdentifier())

		resps = append(resps, rs)

		sub.SetSubscriptionOption("InfoChannelRef", strings.Join(gm.InfoChannelRef(), ","))
		sub.SetSubscriptionOption("LineRef", strings.Join(gm.LineRef(), ","))
		sub.SetSubscriptionOption("StopPointRef", strings.Join(gm.StopPointRef(), ","))
		sub.SetSubscriptionOption("MessageIdentifier", gm.MessageIdentifier())

		obj := model.NewCode("SituationResource", "Situation")
		r := sub.Resource(obj)
		if r == nil {
			ref := model.Reference{
				Code: &obj,
				Type: "Situation",
			}
			r = sub.CreateAndAddNewResource(ref)
			r.Subscribed(connector.Clock().Now())
			r.SubscribedUntil = gm.InitialTerminationTime()
		}

		sub.Save()

		connector.addSituations(sub, r)
	}

	message.Type = audit.GENERAL_MESSAGE_SUBSCRIPTION_REQUEST
	message.SubscriptionIdentifiers = subIds

	return resps
}

func (connector *SIRIGeneralMessageSubscriptionBroadcaster) addSituations(sub *Subscription, r *SubscribedResource) {
	situations := connector.partner.Model().Situations().FindAll()
	for i := range situations {
		if situations[i].GMValidUntil().Before(connector.Clock().Now()) {
			continue
		}

		r.SetLastState(string(situations[i].Id()), ls.NewSituationLastChange(&situations[i], sub))
		connector.addSituation(sub.Id(), situations[i].Id())
	}
}

// Start Test

type TestSIRIGeneralMessageSubscriptionBroadcasterFactory struct{}

type TestGeneralMessageSubscriptionBroadcaster struct {
	connector

	events []*model.SituationBroadcastEvent
	// generalMessageBroadcaster SIRIGeneralMessageBroadcaster
}

func NewTestGeneralMessageSubscriptionBroadcaster() *TestGeneralMessageSubscriptionBroadcaster {
	connector := &TestGeneralMessageSubscriptionBroadcaster{}
	return connector
}

func (connector *TestGeneralMessageSubscriptionBroadcaster) HandleGeneralMessageBroadcastEvent(event *model.SituationBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIGeneralMessageSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
} // Always valid

func (factory *TestSIRIGeneralMessageSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestGeneralMessageSubscriptionBroadcaster()
}

// END OF TEST
