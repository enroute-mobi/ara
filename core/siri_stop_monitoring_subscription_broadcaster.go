package core

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIStopMonitoringSubscriptionBroadcaster struct {
	connector

	stopMonitoringBroadcaster SIRIStopMonitoringBroadcaster
	toBroadcast               map[SubscriptionId][]model.StopVisitId
	notMonitored              map[SubscriptionId]map[string]struct{}

	mutex *sync.Mutex //protect the map
}

type SIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

func (factory *SIRIStopMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIStopMonitoringSubscriptionBroadcaster(partner)
}

func (factory *SIRIStopMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func newSIRIStopMonitoringSubscriptionBroadcaster(partner *Partner) *SIRIStopMonitoringSubscriptionBroadcaster {
	connector := &SIRIStopMonitoringSubscriptionBroadcaster{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
	connector.partner = partner
	connector.mutex = &sync.Mutex{}
	connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)
	connector.notMonitored = make(map[SubscriptionId]map[string]struct{})

	connector.stopMonitoringBroadcaster = NewSIRIStopMonitoringBroadcaster(connector)

	return connector
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) Stop() {
	connector.stopMonitoringBroadcaster.Stop()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) Start() {
	connector.stopMonitoringBroadcaster.Start()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) HandleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	switch event.ModelType {
	case "StopVisit":
		sv, ok := connector.partner.Model().StopVisits().Find(model.StopVisitId(event.ModelId))
		if ok {
			subsIds := connector.checkEvent(sv)
			if len(subsIds) != 0 {
				connector.addStopVisit(subsIds, sv.Id())
			}
		}
	case "VehicleJourney":
		svs := connector.partner.Model().StopVisits().FindFollowingByVehicleJourneyId(model.VehicleJourneyId(event.ModelId))
		for i := range svs {
			subsIds := connector.checkEvent(svs[i])
			if len(subsIds) != 0 {
				connector.addStopVisit(subsIds, svs[i].Id())
			}
		}
	case "StopArea":
		sa, ok := connector.partner.Model().StopAreas().Find(model.StopAreaId(event.ModelId))
		if ok {
			connector.checkStopAreaEvent(sa)
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

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) checkEvent(sv *model.StopVisit) (subscriptionIds []SubscriptionId) {
	if sv.Origin == string(connector.Partner().Slug()) {
		return
	}

	vj, _ := connector.partner.Model().VehicleJourneys().Find(sv.VehicleJourneyId)

	for _, stopAreaCode := range connector.partner.Model().StopAreas().FindAscendantsWithCodeSpace(sv.StopAreaId, connector.remoteCodeSpace) {
		subs := connector.partner.Subscriptions().FindByResourceId(stopAreaCode.String(), StopMonitoringBroadcast)

		for _, sub := range subs {
			resource := sub.Resource(stopAreaCode)
			if resource == nil || resource.SubscribedUntil.Before(connector.Clock().Now()) {
				continue
			}

			// Handle LineRef filter
			if lineRef, ok := connector.lineRef(sub); ok && lineRef != vj.LineId {
				continue
			}

			lastState, ok := resource.LastState(string(sv.Id()))
			if ok && !lastState.(*ls.StopMonitoringLastChange).Haschanged(sv) {
				continue
			}

			if !ok {
				resource.SetLastState(string(sv.Id()), ls.NewStopMonitoringLastChange(sv, sub))
			}

			subscriptionIds = append(subscriptionIds, sub.Id())
		}
	}

	return
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) checkStopAreaEvent(stopArea *model.StopArea) {
	obj, ok := stopArea.Code(connector.remoteCodeSpace)
	if !ok {
		return
	}

	connector.mutex.Lock()

	subs := connector.partner.Subscriptions().FindByResourceId(obj.String(), StopMonitoringBroadcast)
	for _, sub := range subs {
		resource := sub.Resource(obj)
		if resource == nil || resource.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := resource.LastState(string(stopArea.Id()))
		if ok {
			partners, ok := lastState.(*ls.StopAreaLastChange).Haschanged(stopArea)
			if ok {
				nm, ok := connector.notMonitored[sub.Id()]
				if !ok {
					nm = make(map[string]struct{})
					connector.notMonitored[sub.Id()] = nm
				}
				for _, partner := range partners {
					if partner != string(connector.partner.Slug()) {
						nm[partner] = struct{}{}
					}
				}
			}
			lastState.(*ls.StopAreaLastChange).UpdateState(stopArea)
		} else { // Should not happen
			resource.SetLastState(string(stopArea.Id()), ls.NewStopAreaLastChange(stopArea, sub))
		}
	}

	connector.mutex.Unlock()
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) HandleSubscriptionRequest(request *sxml.XMLSubscriptionRequest, message *audit.BigQueryMessage) (resps []siri.SIRIResponseStatus) {
	var monitoringRefs, subIds []string

	for _, sm := range request.XMLSubscriptionSMEntries() {
		rs := siri.SIRIResponseStatus{
			RequestMessageRef: sm.MessageIdentifier(),
			SubscriberRef:     sm.SubscriberRef(),
			SubscriptionRef:   sm.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		monitoringRefs = append(monitoringRefs, sm.MonitoringRef())

		code := model.NewCode(connector.remoteCodeSpace, sm.MonitoringRef())
		sa, ok := connector.partner.Model().StopAreas().FindByCode(code)
		if !ok {
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = fmt.Sprintf("StopArea not found: '%s'", code.Value())
			resps = append(resps, rs)

			message.Status = "Error"
			continue
		}

		subIds = append(subIds, sm.SubscriptionIdentifier())

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(sm.SubscriptionIdentifier())
		if ok {
			if sub.Kind() != StopMonitoringBroadcast {
				logger.Log.Debugf("StopMonitoring subscription request with a duplicated Id: %v", sm.SubscriptionIdentifier())
				rs.ErrorType = "OtherError"
				rs.ErrorNumber = 2
				rs.ErrorText = fmt.Sprintf("[BAD_REQUEST] Subscription Id %v already exists", sm.SubscriptionIdentifier())
				resps = append(resps, rs)

				message.Status = "Error"
				continue
			}

			sub.Delete()
		}

		sub = connector.Partner().Subscriptions().New(StopMonitoringBroadcast)
		sub.SubscriberRef = sm.SubscriberRef()
		sub.SetExternalId(sm.SubscriptionIdentifier())

		ref := model.Reference{
			Code: &code,
			Type:     "StopArea",
		}

		r := sub.CreateAndAddNewResource(ref)
		r.Subscribed(connector.Clock().Now())
		r.SubscribedUntil = sm.InitialTerminationTime()

		connector.fillOptions(sub, r, request, sm)
		if sm.LineRef() != "" {
			sub.SetSubscriptionOption("LineRef", fmt.Sprintf("%s:%s", connector.remoteCodeSpace, sm.LineRef()))
		}

		rs.Status = true
		rs.ValidUntil = sm.InitialTerminationTime()
		resps = append(resps, rs)

		// Init SA LastChange
		r.SetLastState(string(sa.Id()), ls.NewStopAreaLastChange(sa, sub))
		// Init StopVisits LastChange
		connector.addStopAreaStopVisits(sa, sub, r)

		sub.Save()
	}

	message.Type = audit.STOP_MONITORING_SUBSCRIPTION_REQUEST
	message.SubscriptionIdentifiers = subIds
	message.StopAreas = monitoringRefs

	return
}

func (connector *SIRIStopMonitoringSubscriptionBroadcaster) addStopAreaStopVisits(sa *model.StopArea, sub *Subscription, res *SubscribedResource) {
	for _, saId := range connector.partner.Model().StopAreas().FindFamily(sa.Id()) {
		svs := connector.partner.Model().StopVisits().FindFollowingByStopAreaId(saId)
		for i := range svs {
			if _, ok := res.LastState(string(svs[i].Id())); ok {
				continue
			}

			// Handle LineRef filter
			vj, _ := connector.partner.Model().VehicleJourneys().Find(svs[i].VehicleJourneyId)
			if lineRef, ok := connector.lineRef(sub); ok && lineRef != vj.LineId {
				continue
			}

			res.SetLastState(string(svs[i].Id()), ls.NewStopMonitoringLastChange(svs[i], sub))
			connector.addStopVisit([]SubscriptionId{sub.Id()}, svs[i].Id())
		}
	}
}

// WIP Need to do something about this method Refs #6338
func (smsb *SIRIStopMonitoringSubscriptionBroadcaster) fillOptions(s *Subscription, r *SubscribedResource, request *sxml.XMLSubscriptionRequest, sm *sxml.XMLStopMonitoringSubscriptionRequestEntry) {
	changeBeforeUpdates := request.ChangeBeforeUpdates()
	if changeBeforeUpdates == "" {
		changeBeforeUpdates = "PT1M"
	}

	s.SetSubscriptionOption("StopVisitTypes", sm.StopVisitTypes())
	s.SetSubscriptionOption("IncrementalUpdates", request.IncrementalUpdates())
	s.SetSubscriptionOption("MaximumStopVisits", strconv.Itoa(sm.MaximumStopVisits()))
	s.SetSubscriptionOption("ChangeBeforeUpdates", changeBeforeUpdates)
	s.SetSubscriptionOption("MessageIdentifier", request.MessageIdentifier())
}

// Returns the LineId of the line defined in the LineRef subscription option
// If LineRef isn't defined or with an incorrect format, returns false
func (connector *SIRIStopMonitoringSubscriptionBroadcaster) lineRef(sub *Subscription) (model.LineId, bool) {
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
		return "", true
	}
	return line.Id(), true
}

// START TEST

type TestSIRIStopMonitoringSubscriptionBroadcasterFactory struct{}

type TestStopMonitoringSubscriptionBroadcaster struct {
	connector

	events []*model.StopMonitoringBroadcastEvent
}

func NewTestStopMonitoringSubscriptionBroadcaster() *TestStopMonitoringSubscriptionBroadcaster {
	connector := &TestStopMonitoringSubscriptionBroadcaster{}
	return connector
}

func (connector *TestStopMonitoringSubscriptionBroadcaster) HandleStopMonitoringBroadcastEvent(event *model.StopMonitoringBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIStopMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
} // Always valid

func (factory *TestSIRIStopMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestStopMonitoringSubscriptionBroadcaster()
}

// END TEST
