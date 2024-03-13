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

type SIRISituationExchangeSubscriptionBroadcaster struct {
	connector

	situationExchangeBroadcaster SIRISituationExchangeBroadcaster
	toBroadcast                  map[SubscriptionId][]model.SituationId
	mutex                        *sync.Mutex //protect the map
}

type SIRISituationExchangeSubscriptionBroadcasterFactory struct{}

func (factory *SIRISituationExchangeSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRISituationExchangeSubscriptionBroadcaster(partner)
}

func (factory *SIRISituationExchangeSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func newSIRISituationExchangeSubscriptionBroadcaster(partner *Partner) *SIRISituationExchangeSubscriptionBroadcaster {
	siriSituationExchangeSubscriptionBroadcaster := &SIRISituationExchangeSubscriptionBroadcaster{}
	siriSituationExchangeSubscriptionBroadcaster.partner = partner
	siriSituationExchangeSubscriptionBroadcaster.mutex = &sync.Mutex{}
	siriSituationExchangeSubscriptionBroadcaster.toBroadcast = make(map[SubscriptionId][]model.SituationId)

	siriSituationExchangeSubscriptionBroadcaster.situationExchangeBroadcaster = NewSIRISituationExchangeBroadcaster(siriSituationExchangeSubscriptionBroadcaster)

	return siriSituationExchangeSubscriptionBroadcaster
}

func (connector *SIRISituationExchangeSubscriptionBroadcaster) Stop() {
	if connector.situationExchangeBroadcaster != nil {
		connector.situationExchangeBroadcaster.Stop()
	}
}

func (connector *SIRISituationExchangeSubscriptionBroadcaster) Start() {
	if connector.situationExchangeBroadcaster == nil {
		connector.situationExchangeBroadcaster = NewSIRISituationExchangeBroadcaster(connector)
	}
	connector.situationExchangeBroadcaster.Start()
}

func (connector *SIRISituationExchangeSubscriptionBroadcaster) HandleSituationExchangeBroadcastEvent(event *model.SituationBroadcastEvent) {
	connector.checkEvent(event.SituationId)
}

func (connector *SIRISituationExchangeSubscriptionBroadcaster) addSituation(subId SubscriptionId, svId model.SituationId) {
	connector.mutex.Lock()
	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], svId)
	connector.mutex.Unlock()
}

func (connector *SIRISituationExchangeSubscriptionBroadcaster) checkEvent(sId model.SituationId) {
	situation, ok := connector.partner.Model().Situations().Find(sId)
	if !ok || situation.Origin == string(connector.partner.Slug()) {
		return
	}

	obj := model.NewCode("SituationResource", "Situation")
	subs := connector.partner.Subscriptions().FindSubscriptionsByKind(SituationExchangeBroadcast)

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

func (connector *SIRISituationExchangeSubscriptionBroadcaster) HandleSubscriptionRequest(request *sxml.XMLSubscriptionRequest, message *audit.BigQueryMessage) []siri.SIRIResponseStatus {
	resps := []siri.SIRIResponseStatus{}

	var subIds []string

	for _, sx := range request.XMLSubscriptionSXEntries() {
		rs := siri.SIRIResponseStatus{
			RequestMessageRef: sx.MessageIdentifier(),
			SubscriberRef:     sx.SubscriberRef(),
			SubscriptionRef:   sx.SubscriptionIdentifier(),
			Status:            true,
			ResponseTimestamp: connector.Clock().Now(),
			ValidUntil:        sx.InitialTerminationTime(),
		}

		subIds = append(subIds, sx.SubscriptionIdentifier())

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(sx.SubscriptionIdentifier())
		if ok {
			if sub.Kind() != SituationExchangeBroadcast {
				logger.Log.Debugf("SituationExchange subscription request with a duplicated Id: %v", sx.SubscriptionIdentifier())
				rs.Status = false
				rs.ErrorType = "OtherError"
				rs.ErrorNumber = 2
				rs.ErrorText = fmt.Sprintf("[BAD_REQUEST] Subscription Id %v already exists", sx.SubscriptionIdentifier())
				resps = append(resps, rs)
				message.Status = "Error"
				continue
			}

			sub.Delete()
		}

		sub = connector.Partner().Subscriptions().New(SituationExchangeBroadcast)
		sub.SubscriberRef = sx.SubscriberRef()
		sub.SetExternalId(sx.SubscriptionIdentifier())

		resps = append(resps, rs)

		sub.SetSubscriptionOption("LineRef", strings.Join(sx.LineRefs(), ","))
		sub.SetSubscriptionOption("StopPointRef", strings.Join(sx.StopPointRefs(), ","))
		obj := model.NewCode("SituationResource", "Situation")
		r := sub.Resource(obj)
		if r == nil {
			ref := model.Reference{
				Code: &obj,
				Type: "Situation",
			}
			r = sub.CreateAndAddNewResource(ref)
			r.Subscribed(connector.Clock().Now())
			r.SubscribedUntil = sx.InitialTerminationTime()
		}

		sub.Save()

		connector.addSituations(sub, r)
	}

	message.Type = audit.SITUATION_EXCHANGE_SUBSCRIPTION_REQUEST
	message.SubscriptionIdentifiers = subIds

	return resps
}

func (connector *SIRISituationExchangeSubscriptionBroadcaster) addSituations(sub *Subscription, r *SubscribedResource) {
	situations := connector.partner.Model().Situations().FindAll()
	for i := range situations {
		if situations[i].GMValidUntil().Before(connector.Clock().Now()) {
			continue
		}

		r.SetLastState(string(situations[i].Id()), ls.NewSituationLastChange(&situations[i], sub))
		connector.addSituation(sub.Id(), situations[i].Id())
	}
}
