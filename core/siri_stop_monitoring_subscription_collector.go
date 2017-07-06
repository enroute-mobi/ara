package core

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringSubscriptionCollector interface {
	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
	HandleNotifyStopMonitoring(delivery *siri.XMLNotifyStopMonitoring)
	HandleTerminatedNotification(termination *siri.XMLStopMonitoringSubscriptionTerminatedResponse)
}

type SIRIStopMonitoringSubscriptionCollector struct {
	// model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	stopMonitoringSubscriber SIRIStopMonitoringSubscriber
	stopAreaUpdateSubscriber StopAreaUpdateSubscriber
}

type SIRIStopMonitoringSubscriptionCollectorFactory struct{}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func NewSIRIStopMonitoringSubscriptionCollector(partner *Partner) *SIRIStopMonitoringSubscriptionCollector {
	siriStopMonitoringSubscriptionCollector := &SIRIStopMonitoringSubscriptionCollector{}
	siriStopMonitoringSubscriptionCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriStopMonitoringSubscriptionCollector.stopAreaUpdateSubscriber = manager.BroadcastStopAreaUpdateEvent

	return siriStopMonitoringSubscriptionCollector
}

func (connector *SIRIStopMonitoringSubscriptionCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	stopArea, ok := connector.Partner().Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("StopAreaUpdateRequest in StopMonitoring SubscriptionCollector for unknown StopArea %v", request.StopAreaId())
		return
	}

	objectidKind := connector.Partner().Setting("remote_objectid_kind")
	stopAreaObjectid, ok := stopArea.ObjectID(objectidKind)
	if !ok {
		logger.Log.Debugf("Requested stopArea %v doesn't have and objectId of kind %v", request.StopAreaId(), objectidKind)
		return
	}

	subscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoring")

	for _, sr := range subscription.resourcesByObjectID {
		if sr.Reference.ObjectId.String() == stopAreaObjectid.String() && !sr.SubscribedAt.IsZero() {
			sr.SubscribedUntil = sr.SubscribedUntil.Add(1 * time.Minute)
			return
		}
	}

	ref := model.Reference{
		ObjectId: &stopAreaObjectid,
		Id:       string(request.StopAreaId()),
		Type:     "StopArea",
	}

	subscription.CreateAddNewResource(ref)

	if connector.stopMonitoringSubscriber == nil {
		connector.stopMonitoringSubscriber = NewSIRIStopMonitoringSubscriber(connector)
		connector.stopMonitoringSubscriber.Run()
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) SetStopMonitoringSubscriber(stopMonitoringSubscriber SIRIStopMonitoringSubscriber) {
	connector.stopMonitoringSubscriber = stopMonitoringSubscriber
}

func (connector *SIRIStopMonitoringSubscriptionCollector) SetStopAreaUpdateSubscriber(stopAreaUpdateSubscriber StopAreaUpdateSubscriber) {
	connector.stopAreaUpdateSubscriber = stopAreaUpdateSubscriber
}

func (connector *SIRIStopMonitoringSubscriptionCollector) broadcastStopAreaUpdateEvent(event *model.StopAreaUpdateEvent) {
	if connector.stopAreaUpdateSubscriber != nil {
		connector.stopAreaUpdateSubscriber(event)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) HandleTerminatedNotification(response *siri.XMLStopMonitoringSubscriptionTerminatedResponse) {
	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logSIRISubscriptionTerminatedResponse(logStashEvent, response)

	subscriptionTerminated := response.XMLSubscriptionTerminateds()
	connector.setSubscriptionTerminatedEvents(subscriptionTerminated)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) setSubscriptionTerminatedEvents(terminations []*siri.XMLSubscriptionTerminated) {
	for _, termination := range terminations {
		sub, present := connector.partner.Subscriptions().Find(SubscriptionId(termination.SubscriptionRef()))
		if !present {
			continue
		}
		connector.partner.DeleteSubscription(sub)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) HandleNotifyStopMonitoring(notify *siri.XMLNotifyStopMonitoring) {
	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLStopMonitoringDelivery(logStashEvent, notify)

	stopAreaUpdateEvents := make(map[string]*model.StopAreaUpdateEvent)

	for _, delivery := range notify.StopMonitoringDeliveries() {

		reg := regexp.MustCompile(`\w+:Subscription::([\w+-?]+):LOC`)
		matches := reg.FindStringSubmatch(strings.TrimSpace(delivery.SubscriptionRef()))

		if len(matches) == 0 {
			logger.Log.Printf("Partner %s sent a StopVisitNotify response with a wrong message format: %s\n", connector.Partner().Slug(), delivery.SubscriptionRef())
			continue
		}
		subscriptionId := matches[1]
		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))

		if ok == false {
			logger.Log.Printf("Partner %s sent a StopVisitNotify response to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			continue
		}
		if subscription.Kind() != "StopMonitoring" {
			logger.Log.Printf("Partner %s sent a StopVisitNotify response to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			continue
		}

		connector.setStopVisitUpdateEvents(stopAreaUpdateEvents, delivery)
		connector.setStopVisitCancellationEvents(stopAreaUpdateEvents, delivery)
	}

	logStopVisitUpdateEventsFromMap(logStashEvent, stopAreaUpdateEvents)

	for _, event := range stopAreaUpdateEvents {
		connector.broadcastStopAreaUpdateEvent(event)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) setStopVisitUpdateEvents(events map[string]*model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringDelivery) {
	xmlStopVisitEvents := xmlResponse.XMLMonitoredStopVisits()
	if len(xmlStopVisitEvents) == 0 {
		return
	}

	builder := newStopVisitUpdateEventBuilder(connector.partner)

	for _, xmlStopVisitEvent := range xmlStopVisitEvents {
		stopAreaObjectId := model.NewObjectID(connector.Partner().Setting("remote_objectid_kind"), xmlStopVisitEvent.StopPointRef())
		stopArea, ok := connector.Partner().Model().StopAreas().FindByObjectId(stopAreaObjectId)
		if !ok {
			logger.Log.Debugf("StopVisitUpdateEvent for unknown StopArea %v", stopAreaObjectId.Value())
			continue
		}

		stopAreaUpdateEvent, ok := events[xmlStopVisitEvent.StopPointRef()]
		if !ok {
			stopAreaUpdateEvent = model.NewStopAreaUpdateEvent(connector.NewUUID(), stopArea.Id())
			events[xmlStopVisitEvent.StopPointRef()] = stopAreaUpdateEvent
		}
		builder.buildStopVisitUpdateEvent(stopAreaUpdateEvent, xmlStopVisitEvent)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) setStopVisitCancellationEvents(events map[string]*model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringDelivery) {
	xmlStopVisitCancellationEvents := xmlResponse.XMLMonitoredStopVisitCancellations()
	if len(xmlStopVisitCancellationEvents) == 0 {
		return
	}

	for _, xmlStopVisitCancellationEvent := range xmlStopVisitCancellationEvents {
		stopAreaObjectId := model.NewObjectID(connector.Partner().Setting("remote_objectid_kind"), xmlStopVisitCancellationEvent.MonitoringRef())
		stopArea, ok := connector.Partner().Model().StopAreas().FindByObjectId(stopAreaObjectId)
		if !ok {
			logger.Log.Debugf("StopVisitCancellationEvent for unknown StopArea %v", stopAreaObjectId.Value())
			continue
		}

		stopAreaUpdateEvent, ok := events[xmlStopVisitCancellationEvent.MonitoringRef()]
		if !ok {
			stopAreaUpdateEvent = model.NewStopAreaUpdateEvent(connector.NewUUID(), stopArea.Id())
			events[xmlStopVisitCancellationEvent.MonitoringRef()] = stopAreaUpdateEvent
		}
		stopVisitCancellationEvent := &model.StopVisitNotCollectedEvent{
			StopVisitObjectId: model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlStopVisitCancellationEvent.ItemRef()),
		}
		stopAreaUpdateEvent.StopVisitNotCollectedEvents = append(stopAreaUpdateEvent.StopVisitNotCollectedEvents, stopVisitCancellationEvent)
	}
}

func logXMLStopMonitoringDelivery(logStashEvent audit.LogStashEvent, notify *siri.XMLNotifyStopMonitoring) {
	logStashEvent["Connector"] = "StopMonitoringSubscriptionCollector"
	logStashEvent["address"] = notify.Address()
	logStashEvent["producerRef"] = notify.ProducerRef()
	logStashEvent["requestMessageRef"] = notify.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = notify.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = notify.ResponseTimestamp().String()
	logStashEvent["responseXML"] = notify.RawXML()
	logStashEvent["status"] = strconv.FormatBool(notify.Status())
	if !notify.Status() {
		logStashEvent["errorType"] = notify.ErrorType()
		logStashEvent["errorNumber"] = strconv.Itoa(notify.ErrorNumber())
		logStashEvent["errorText"] = notify.ErrorText()
		logStashEvent["errorDescription"] = notify.ErrorDescription()
	}
}

func logSIRISubscriptionTerminatedResponse(logStashEvent audit.LogStashEvent, response *siri.XMLStopMonitoringSubscriptionTerminatedResponse) {
	logStashEvent["Connector"] = "StopMonitoringSubscriptionRequestCollector"
	logStashEvent["address"] = response.Address()
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["responseXML"] = response.RawXML()

	logStashEvent["status"] = strconv.FormatBool(response.Status())
	if !response.Status() {
		logStashEvent["errorType"] = response.ErrorType()
		logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber())
		logStashEvent["errorText"] = response.ErrorText()
		logStashEvent["errorDescription"] = response.ErrorDescription()
	}
}

func logStopVisitUpdateEventsFromMap(logStashEvent audit.LogStashEvent, stopAreaUpdateEvents map[string]*model.StopAreaUpdateEvent) {
	var idArray []string
	var cancelledIdArray []string
	for _, stopAreaUpdateEvent := range stopAreaUpdateEvents {
		for _, stopVisitUpdateEvent := range stopAreaUpdateEvent.StopVisitUpdateEvents {
			idArray = append(idArray, stopVisitUpdateEvent.Id)
		}
		for _, stopVisitCancelledEvent := range stopAreaUpdateEvent.StopVisitNotCollectedEvents {
			cancelledIdArray = append(cancelledIdArray, stopVisitCancelledEvent.StopVisitObjectId.Value())
		}
	}
	logStashEvent["StopVisitUpdateEventIds"] = strings.Join(idArray, ", ")
	logStashEvent["StopVisitCancelledEventIds"] = strings.Join(cancelledIdArray, ", ")
}
