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

type SIRIFacilityMonitoringSubscriptionBroadcaster struct {
	connector

	facilityMonitoringBroadcaster FacilityMonitoringBroadcaster
	toBroadcast                   map[SubscriptionId][]model.FacilityId

	mutex *sync.Mutex //protect the map
}

type SIRIFacilityMonitoringSubscriptionBroadcasterFactory struct{}

func (factory *SIRIFacilityMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIFacilityMonitoringSubscriptionBroadcaster(partner)
}

func (factory *SIRIFacilityMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func newSIRIFacilityMonitoringSubscriptionBroadcaster(partner *Partner) *SIRIFacilityMonitoringSubscriptionBroadcaster {
	connector := &SIRIFacilityMonitoringSubscriptionBroadcaster{}
	connector.remoteCodeSpace = partner.RemoteCodeSpace(SIRI_FACILITY_MONITORING_SUBSCRIPTION_BROADCASTER)
	connector.partner = partner
	connector.mutex = &sync.Mutex{}
	connector.toBroadcast = make(map[SubscriptionId][]model.FacilityId)

	connector.facilityMonitoringBroadcaster = NewSIRIFacilityMonitoringBroadcaster(connector)
	return connector
}

func (connector *SIRIFacilityMonitoringSubscriptionBroadcaster) HandleSubscriptionRequest(request *sxml.XMLSubscriptionRequest, message *audit.BigQueryMessage) (resps []siri.SIRIResponseStatus) {
	var facilityIds, subIds []string

	for _, fm := range request.XMLSubscriptionFMEntries() {
		rs := siri.SIRIResponseStatus{
			RequestMessageRef: fm.MessageIdentifier(),
			SubscriberRef:     fm.SubscriberRef(),
			SubscriptionRef:   fm.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		// for logging
		facilityIds = append(facilityIds, fm.FacilityRefs()...)
		sub, ok := connector.Partner().Subscriptions().FindByExternalId(fm.SubscriptionIdentifier())
		if ok {
			if sub.Kind() != FacilityMonitoringBroadcast {
				logger.Log.Debugf("FacilityMonitoring subscription request with a duplicated Id: %v", fm.SubscriptionIdentifier())
				rs.Status = false
				rs.ErrorType = "OtherError"
				rs.ErrorNumber = 2
				rs.ErrorText = fmt.Sprintf("[BAD_REQUEST] Subscription Id %v already exists", fm.SubscriptionIdentifier())

				resps = append(resps, rs)
				message.Status = "Error"
				continue
			}

			sub.Delete()
		}

		resources, unknownFacilityRefs := connector.checkFacilities(fm)
		if len(unknownFacilityRefs) != 0 {
			logger.Log.Debugf("FacilityMonitoring subscription request Could not find facility(ies) with id : %v", strings.Join(unknownFacilityRefs, ","))
			rs.Status = false
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = fmt.Sprintf("Unknown Facility(ies) %v", strings.Join(unknownFacilityRefs, ","))

			resps = append(resps, rs)
			message.Status = "Error"
			continue
		}

		rs.Status = true
		rs.ValidUntil = fm.InitialTerminationTime()
		resps = append(resps, rs)

		subIds = append(subIds, fm.SubscriptionIdentifier())

		sub = connector.Partner().Subscriptions().New(FacilityMonitoringBroadcast)
		sub.SubscriberRef = fm.SubscriberRef()
		sub.SetExternalId(fm.SubscriptionIdentifier())
		connector.fillOptions(sub, request)

		for _, r := range resources {
			facility, ok := connector.Partner().Model().Facilities().FindByCode(*r.Reference.Code)
			if !ok {
				continue
			}

			// Init StopVisits LastChange
			r.SetLastState(string(facility.Id()), ls.NewFacilityMonitoringLastChange(facility, sub))
			sub.AddNewResource(r)
		}
		sub.Save()
	}

	message.Type = audit.FACILITY_MONITORING_SUBSCRIPTION_REQUEST
	message.SubscriptionIdentifiers = subIds

	message.Facilities = facilityIds

	return resps
}

func (connector *SIRIFacilityMonitoringSubscriptionBroadcaster) checkFacilities(fm *sxml.XMLFacilityMonitoringSubscriptionRequestEntry) (resources []*SubscribedResource, unknownFacilityRefs []string) {
	// check for subscription to all facilities
	if len(fm.FacilityRefs()) == 0 {
		var facilityValues []string
		//find all facilities corresponding to the remoteCodeSpace
		for _, facility := range connector.Partner().Model().Facilities().FindAll() {
			facilityCode, ok := facility.Code(connector.remoteCodeSpace)
			if ok {
				facilityValues = append(facilityValues, facilityCode.Value())
				continue
			}
		}

		for _, facilityValue := range facilityValues {
			facilityCode := model.NewCode(connector.remoteCodeSpace, facilityValue)
			ref := model.Reference{
				Code: &facilityCode,
				Type: "Facility",
			}
			r := NewResource(ref)
			r.Subscribed(connector.Clock().Now())
			r.SubscribedUntil = fm.InitialTerminationTime()
			resources = append(resources, r)
		}
		return resources, unknownFacilityRefs
	}

	for _, facilityRef := range fm.FacilityRefs() {
		facilityCode := model.NewCode(connector.remoteCodeSpace, facilityRef)
		_, ok := connector.Partner().Model().Facilities().FindByCode(facilityCode)
		if !ok {
			unknownFacilityRefs = append(unknownFacilityRefs, facilityRef)
			continue
		}

		ref := model.Reference{
			Code: &facilityCode,
			Type: "Facility",
		}

		r := NewResource(ref)
		r.Subscribed(connector.Clock().Now())
		r.SubscribedUntil = fm.InitialTerminationTime()
		resources = append(resources, r)
	}
	return resources, unknownFacilityRefs
}

func (connector *SIRIFacilityMonitoringSubscriptionBroadcaster) Stop() {
	connector.facilityMonitoringBroadcaster.Stop()
}

func (connector *SIRIFacilityMonitoringSubscriptionBroadcaster) Start() {
	connector.facilityMonitoringBroadcaster.Start()
}

func (fmb *SIRIFacilityMonitoringSubscriptionBroadcaster) fillOptions(s *Subscription, request *sxml.XMLSubscriptionRequest) {
	s.SetSubscriptionOption("MessageIdentifier", request.MessageIdentifier())
}

func (connector *SIRIFacilityMonitoringSubscriptionBroadcaster) HandleBroadcastEvent(event *model.FacilityBroadcastEvent) {
	connector.checkEvent(model.FacilityId(event.ModelId))
}

func (connector *SIRIFacilityMonitoringSubscriptionBroadcaster) checkEvent(fId model.FacilityId) {
	facility, ok := connector.Partner().Model().Facilities().Find(fId)
	if !ok {
		return
	}

	facilityObj, ok := facility.Code(connector.remoteCodeSpace)
	if !ok {
		return
	}

	subs := connector.Partner().Subscriptions().FindByResourceId(facilityObj.String(), FacilityMonitoringBroadcast)

	for _, sub := range subs {
		r := sub.Resource(facilityObj)
		if r == nil || r.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := r.LastState(string(facility.Id()))
		if ok && !lastState.(*ls.FacilityMonitoringLastChange).HasChanged(facility) {
			continue
		}

		if !ok {
			r.SetLastState(string(facility.Id()), ls.NewFacilityMonitoringLastChange(facility, sub))
		}
		connector.addFacility(sub.Id(), facility.Id())
	}
}

func (connector *SIRIFacilityMonitoringSubscriptionBroadcaster) addFacility(subId SubscriptionId, fId model.FacilityId) {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], fId)

}

// START TEST

type TestSIRIFMSubscriptionBroadcasterFactory struct{}

type TestFMSubscriptionBroadcaster struct {
	connector

	events []*model.FacilityBroadcastEvent
}

func NewTestFMSubscriptionBroadcaster() *TestFMSubscriptionBroadcaster {
	connector := &TestFMSubscriptionBroadcaster{}
	return connector
}

func (connector *TestFMSubscriptionBroadcaster) HandleBroadcastEvent(event *model.FacilityBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIFMSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {} // Always valid

func (factory *TestSIRIFMSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestFMSubscriptionBroadcaster()
}
