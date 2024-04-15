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
	connector := &SIRISituationExchangeSubscriptionBroadcaster{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace(SIRI_SITUATION_EXCHANGE_SUBSCRIPTION_BROADCASTER)
	connector.partner = partner
	connector.mutex = &sync.Mutex{}
	connector.toBroadcast = make(map[SubscriptionId][]model.SituationId)

	connector.situationExchangeBroadcaster = NewSIRISituationExchangeBroadcaster(connector)

	return connector
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
		connector.addfilteredSituations(situation, sub, resource)
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

		if len(sx.LineRefs()) != 0 {
			for _, xmlLineRef := range sx.LineRefs() {
				sub.SetSubscriptionOption("LineRef", fmt.Sprintf("%s:%s", connector.remoteCodeSpace, xmlLineRef))
			}

		}

		if len(sx.StopPointRefs()) != 0 {
			for _, xmlStopPointRef := range sx.StopPointRefs() {
				sub.SetSubscriptionOption("StopPointRefRef", fmt.Sprintf("%s:%s", connector.remoteCodeSpace, xmlStopPointRef))
			}

		}

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
		connector.addfilteredSituations(situations[i], sub, r)
	}
}

func (connector *SIRISituationExchangeSubscriptionBroadcaster) addfilteredSituations(situation model.Situation, sub *Subscription, r *SubscribedResource) {
	if situation.GMValidUntil().Before(connector.Clock().Now()) {
		return
	}

	if sub.SubscriptionOption("LineRef") == "" && sub.SubscriptionOption("StopPointRef") == "" {
		r.SetLastState(string(situation.Id()), ls.NewSituationLastChange(&situation, sub))
		connector.addSituation(sub.Id(), situation.Id())
		return
	}
	// Filtered subscription
	for _, affect := range situation.Affects {
		if affect.GetType() == model.SituationTypeLine {
			if lineRef, ok := connector.lineRef(sub); ok && model.ModelId(lineRef) == affect.GetId() {
				r.SetLastState(string(situation.Id()), ls.NewSituationLastChange(&situation, sub))
				connector.addSituation(sub.Id(), situation.Id())
				continue
			}
		}
		if affect.GetType() == model.SituationTypeStopArea {
			if stopPointRef, ok := connector.stopPointRef(sub); ok && model.ModelId(stopPointRef) == affect.GetId() {
				r.SetLastState(string(situation.Id()), ls.NewSituationLastChange(&situation, sub))
				connector.addSituation(sub.Id(), situation.Id())
				continue
			}
		}
	}
}

// Returns the LineId of the line defined in the LineRef subscription option
// If LineRef isn't defined or with an incorrect format, returns false
func (connector *SIRISituationExchangeSubscriptionBroadcaster) lineRef(sub *Subscription) (model.LineId, bool) {
	lineRef := sub.SubscriptionOption("LineRef")
	if lineRef == "" {
		return "", false
	}
	kindValue := strings.SplitN(lineRef, ":", 2)
	if len(kindValue) != 2 { // Should not happen but we don't want an index out of range panic
		logger.Log.Debugf("The LineRef Setting hasn't been stored in the correct format: %v", lineRef)
		return "", false
	}
	line, ok := connector.partner.Model().Lines().FindByCode(model.NewCode(kindValue[0], kindValue[1]))
	if !ok {
		return "", false
	}
	return line.Id(), true
}

// Returns the StopAreaId of the stopArea defined in the StopPointRef subscription option
// If StopPointRef isn't defined or with an incorrect format, returns false
func (connector *SIRISituationExchangeSubscriptionBroadcaster) stopPointRef(sub *Subscription) (model.StopAreaId, bool) {
	stopPointRef := sub.SubscriptionOption("StopPointRef")
	if stopPointRef == "" {
		return "", false
	}
	kindValue := strings.SplitN(stopPointRef, ":", 2)
	if len(kindValue) != 2 { // Should not happen but we don't want an index out of range panic
		logger.Log.Debugf("The StopPointRef Setting hasn't been stored in the correct format: %v", stopPointRef)
		return "", false
	}
	stopArea, ok := connector.partner.Model().StopAreas().FindByCode(model.NewCode(kindValue[0], kindValue[1]))
	if !ok {
		return "", false
	}
	return stopArea.Id(), true
}
