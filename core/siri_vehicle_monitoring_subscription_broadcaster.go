package core

import (
	"fmt"
	"strings"
	"sync"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIVehicleMonitoringSubscriptionBroadcaster struct {
	connector

	dataFrameGenerator           *idgen.IdentifierGenerator
	vjRemoteObjectidKinds        []string
	vehicleMonitoringBroadcaster VehicleMonitoringBroadcaster
	toBroadcast                  map[SubscriptionId][]model.VehicleId

	mutex *sync.Mutex //protect the map
}

type SIRIVehicleMonitoringSubscriptionBroadcasterFactory struct{}

func (factory *SIRIVehicleMonitoringSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	if _, ok := partner.Connector(SIRI_SUBSCRIPTION_REQUEST_DISPATCHER); !ok {
		partner.CreateSubscriptionRequestDispatcher()
	}
	return newSIRIVehicleMonitoringSubscriptionBroadcaster(partner)
}

func (factory *SIRIVehicleMonitoringSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func newSIRIVehicleMonitoringSubscriptionBroadcaster(partner *Partner) *SIRIVehicleMonitoringSubscriptionBroadcaster {
	connector := &SIRIVehicleMonitoringSubscriptionBroadcaster{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(SIRI_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER)
	connector.vjRemoteObjectidKinds = partner.VehicleJourneyRemoteObjectIDKindWithFallback(SIRI_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER)
	connector.dataFrameGenerator = partner.DataFrameIdentifierGenerator()
	connector.partner = partner
	connector.mutex = &sync.Mutex{}
	connector.toBroadcast = make(map[SubscriptionId][]model.VehicleId)

	connector.vehicleMonitoringBroadcaster = NewSIRIVehicleMonitoringBroadcaster(connector)
	return connector
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) HandleSubscriptionRequest(request *sxml.XMLSubscriptionRequest, message *audit.BigQueryMessage) (resps []siri.SIRIResponseStatus) {
	var lineIds, subIds []string

	for _, vm := range request.XMLSubscriptionVMEntries() {
		rs := siri.SIRIResponseStatus{
			RequestMessageRef: vm.MessageIdentifier(),
			SubscriberRef:     vm.SubscriberRef(),
			SubscriptionRef:   vm.SubscriptionIdentifier(),
			ResponseTimestamp: connector.Clock().Now(),
		}

		// for logging
		lineIds = append(lineIds, vm.Lines()...)

		sub, ok := connector.Partner().Subscriptions().FindByExternalId(vm.SubscriptionIdentifier())
		if ok {
			if sub.Kind() != VehicleMonitoringBroadcast {
				logger.Log.Debugf("VehicleMonitoring subscription request with a duplicated Id: %v", vm.SubscriptionIdentifier())
				rs.ErrorType = "OtherError"
				rs.ErrorNumber = 2
				rs.ErrorText = fmt.Sprintf("[BAD_REQUEST] Subscription Id %v already exists", vm.SubscriptionIdentifier())

				resps = append(resps, rs)
				message.Status = "Error"
				continue
			}

			sub.Delete()
		}

		resources, unknownLineIds := connector.checkLines(vm)
		if len(unknownLineIds) != 0 {
			logger.Log.Debugf("VehicleMonitoring subscription request Could not find line(s) with id : %v", strings.Join(unknownLineIds, ","))
			rs.ErrorType = "InvalidDataReferencesError"
			rs.ErrorText = fmt.Sprintf("Unknown Line(s) %v", strings.Join(unknownLineIds, ","))

			resps = append(resps, rs)
			message.Status = "Error"
			continue
		}

		rs.Status = true
		rs.ValidUntil = vm.InitialTerminationTime()
		resps = append(resps, rs)

		subIds = append(subIds, vm.SubscriptionIdentifier())

		sub = connector.Partner().Subscriptions().New(VehicleMonitoringBroadcast)
		sub.SubscriberRef = vm.SubscriberRef()
		sub.SetExternalId(vm.SubscriptionIdentifier())
		connector.fillOptions(sub, request)

		for _, r := range resources {
			line, ok := connector.Partner().Model().Lines().FindByObjectId(*r.Reference.ObjectId)
			if !ok {
				continue
			}

			// Init Vehicles LastChange
			connector.addLineVehicles(sub, r, line.Id())

			sub.AddNewResource(r)
		}
		sub.Save()
	}
	message.Type = audit.VEHICLE_MONITORING_SUBSCRIPTION_REQUEST
	message.SubscriptionIdentifiers = subIds
	message.Lines = lineIds

	return resps
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) addLineVehicles(sub *Subscription, res *SubscribedResource, lineId model.LineId) {
	vs := connector.partner.Model().Vehicles().FindByLineId(lineId)
	for i := range vs {
		// Init Vehicle LastChange
		res.SetLastState(string(vs[i].Id()), ls.NewVehicleMonitoringLastChange(vs[i], sub))
		connector.addVehicle(sub.Id(), vs[i].Id())
	}
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) checkLines(vm *sxml.XMLVehicleMonitoringSubscriptionRequestEntry) (resources []*SubscribedResource, lineIds []string) {
	// check for subscription to all lines
	if len(vm.Lines()) == 0 {
		var lv []string
		//find all lines corresponding to the remoteObjectidKind
		for _, line := range connector.Partner().Model().Lines().FindAll() {
			lineObjectID, ok := line.ObjectID(connector.remoteObjectidKind)
			if ok {
				lv = append(lv, lineObjectID.Value())
				continue
			}
		}

		for _, lineValue := range lv {
			lineObjectID := model.NewObjectID(connector.remoteObjectidKind, lineValue)
			ref := model.Reference{
				ObjectId: &lineObjectID,
				Type:     "Line",
			}
			r := NewResource(ref)
			r.Subscribed(connector.Clock().Now())
			r.SubscribedUntil = vm.InitialTerminationTime()
			resources = append(resources, r)
		}
		return resources, lineIds
	}

	for _, lineId := range vm.Lines() {

		lineObjectID := model.NewObjectID(connector.remoteObjectidKind, lineId)
		_, ok := connector.Partner().Model().Lines().FindByObjectId(lineObjectID)

		if !ok {
			lineIds = append(lineIds, lineId)
			continue
		}

		ref := model.Reference{
			ObjectId: &lineObjectID,
			Type:     "Line",
		}

		r := NewResource(ref)
		r.Subscribed(connector.Clock().Now())
		r.SubscribedUntil = vm.InitialTerminationTime()
		resources = append(resources, r)
	}
	return resources, lineIds
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) Stop() {
	connector.vehicleMonitoringBroadcaster.Stop()
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) Start() {
	connector.vehicleMonitoringBroadcaster.Start()
}

func (vmb *SIRIVehicleMonitoringSubscriptionBroadcaster) fillOptions(s *Subscription, request *sxml.XMLSubscriptionRequest) {
	changeBeforeUpdates := request.ChangeBeforeUpdates()
	if changeBeforeUpdates == "" {
		changeBeforeUpdates = "PT1M"
	}
	s.SetSubscriptionOption("ChangeBeforeUpdates", changeBeforeUpdates)
	s.SetSubscriptionOption("MessageIdentifier", request.MessageIdentifier())
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) HandleBroadcastEvent(event *model.VehicleBroadcastEvent) {
	switch event.ModelType {
	case "Vehicle":
		connector.checkEvent(model.VehicleId(event.ModelId))
	default:
		return
	}
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) checkEvent(vId model.VehicleId) {
	v, ok := connector.Partner().Model().Vehicles().Find(vId)
	if !ok {
		return
	}

	vj, ok := connector.Partner().Model().VehicleJourneys().Find(v.VehicleJourneyId)
	if !ok {
		return
	}

	line, ok := connector.Partner().Model().Lines().Find(vj.LineId)
	if !ok {
		return
	}

	lineObj, ok := line.ObjectID(connector.remoteObjectidKind)
	if !ok {
		return
	}

	subs := connector.Partner().Subscriptions().FindByResourceId(lineObj.String(), VehicleMonitoringBroadcast)

	for _, sub := range subs {
		r := sub.Resource(lineObj)
		if r == nil || r.SubscribedUntil.Before(connector.Clock().Now()) {
			continue
		}

		lastState, ok := r.LastState(string(vId))
		if ok && !lastState.(*ls.VehicleMonitoringLastChange).HasChanged(v) {
			continue
		}

		if !ok {
			r.SetLastState(string(v.Id()), ls.NewVehicleMonitoringLastChange(v, sub))
		}
		connector.addVehicle(sub.Id(), v.Id())
	}
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) addVehicle(subId SubscriptionId, vId model.VehicleId) {
	connector.mutex.Lock()
	defer connector.mutex.Unlock()

	connector.toBroadcast[SubscriptionId(subId)] = append(connector.toBroadcast[SubscriptionId(subId)], vId)

}

// START TEST

type TestSIRIVMSubscriptionBroadcasterFactory struct{}

type TestVMSubscriptionBroadcaster struct {
	connector

	events []*model.VehicleBroadcastEvent
}

func NewTestVMSubscriptionBroadcaster() *TestVMSubscriptionBroadcaster {
	connector := &TestVMSubscriptionBroadcaster{}
	return connector
}

func (connector *TestVMSubscriptionBroadcaster) HandleBroadcastEvent(event *model.VehicleBroadcastEvent) {
	connector.events = append(connector.events, event)
}

func (factory *TestSIRIVMSubscriptionBroadcasterFactory) Validate(apiPartner *APIPartner) {} // Always valid

func (factory *TestSIRIVMSubscriptionBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewTestVMSubscriptionBroadcaster()
}
