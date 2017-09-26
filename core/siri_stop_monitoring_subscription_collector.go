package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopMonitoringSubscriptionCollector interface {
	model.Stopable
	model.Startable

	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
	HandleNotifyStopMonitoring(delivery *siri.XMLNotifyStopMonitoring)
	HandleTerminatedNotification(termination *siri.XMLStopMonitoringSubscriptionTerminatedResponse)
}

type SIRIStopMonitoringSubscriptionCollector struct {
	model.ClockConsumer
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
	connector := &SIRIStopMonitoringSubscriptionCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.stopAreaUpdateSubscriber = manager.BroadcastStopAreaUpdateEvent
	connector.stopMonitoringSubscriber = NewSIRIStopMonitoringSubscriber(connector)

	return connector
}

func (connector *SIRIStopMonitoringSubscriptionCollector) Stop() {
	connector.stopMonitoringSubscriber.Stop()
}

func (connector *SIRIStopMonitoringSubscriptionCollector) Start() {
	connector.stopMonitoringSubscriber.Start()
}

func (connector *SIRIStopMonitoringSubscriptionCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	stopArea, ok := tx.Model().StopAreas().Find(request.StopAreaId())
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

	subscription, _ := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoringCollect")

	// If we find the resource, we add time to SubscribedUntil if the subscription is active
	resource := subscription.Resource(stopAreaObjectid)
	if resource != nil {
		if !resource.SubscribedAt.IsZero() {
			resource.SubscribedUntil = resource.SubscribedUntil.Add(1 * time.Minute)
		}
		return
	}

	// Else we create a new resource
	ref := model.Reference{
		ObjectId: &stopAreaObjectid,
		Type:     "StopArea",
	}

	subscription.CreateAddNewResource(ref)
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

	logXMLSubscriptionTerminatedResponse(logStashEvent, response)

	subscriptionTerminated := response.XMLSubscriptionTerminateds()
	connector.setSubscriptionTerminatedEvents(subscriptionTerminated)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) setSubscriptionTerminatedEvents(terminations []*siri.XMLSubscriptionTerminated) {
	for _, termination := range terminations {
		connector.partner.Subscriptions().DeleteById(SubscriptionId(termination.SubscriptionRef()))
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) HandleNotifyStopMonitoring(notify *siri.XMLNotifyStopMonitoring) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)
	monitoringRefMap := make(map[string]struct{})

	logXMLStopMonitoringDelivery(logStashEvent, notify)

	stopAreaUpdateEvents := make(map[string]*model.StopAreaUpdateEvent)
	for _, delivery := range notify.StopMonitoringDeliveries() {

		subscriptionId := delivery.SubscriptionRef()

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))

		if !ok {
			logger.Log.Debugf("Partner %s sent a StopVisitNotify response to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			connector.cancelSubscription(delivery.SubscriptionRef())
			continue
		}
		if subscription.Kind() != "StopMonitoringCollect" {
			logger.Log.Debugf("Partner %s sent a StopVisitNotify response to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			continue
		}

		tx := connector.Partner().Referential().NewTransaction()
		defer tx.Close()

		connector.setStopVisitUpdateEvents(stopAreaUpdateEvents, delivery, tx, monitoringRefMap)
		connector.setStopVisitCancellationEvents(stopAreaUpdateEvents, delivery, tx, monitoringRefMap)
	}

	logStopVisitUpdateEventsFromMap(logStashEvent, stopAreaUpdateEvents)
	logMonitoringRefsFromMap(logStashEvent, monitoringRefMap)

	for _, event := range stopAreaUpdateEvents {
		connector.broadcastStopAreaUpdateEvent(event)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) cancelSubscription(subId string) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	request := &siri.SIRIDeleteSubscriptionRequest{
		RequestTimestamp:  connector.Clock().Now(),
		SubscriptionRef:   subId,
		RequestorRef:      connector.partner.ProducerRef(),
		MessageIdentifier: connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier(),
	}
	logSIRIDeleteSubscriptionRequest(logStashEvent, request)

	response, err := connector.SIRIPartner().SOAPClient().DeleteSubscription(request)
	if err != nil {
		logger.Log.Debugf("Error while terminating subcription with id : %v error : ", subId, err.Error())
		logStashEvent["status"] = "false"
		logStashEvent["response"] = fmt.Sprintf("Error during DeleteSubscription: %v", err)
		return
	}
	logXMLDeleteSubscriptionResponse(logStashEvent, response)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) setStopVisitUpdateEvents(events map[string]*model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringDelivery, tx *model.Transaction, monitoringRefMap map[string]struct{}) {
	xmlStopVisitEvents := xmlResponse.XMLMonitoredStopVisits()
	if len(xmlStopVisitEvents) == 0 {
		return
	}

	builder := newStopVisitUpdateEventBuilder(connector.partner)

	for _, xmlStopVisitEvent := range xmlStopVisitEvents {
		monitoringRefMap[xmlStopVisitEvent.StopPointRef()] = struct{}{}
		stopAreaObjectId := model.NewObjectID(connector.Partner().Setting("remote_objectid_kind"), xmlStopVisitEvent.StopPointRef())
		stopArea, ok := tx.Model().StopAreas().FindByObjectId(stopAreaObjectId)
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

func (connector *SIRIStopMonitoringSubscriptionCollector) setStopVisitCancellationEvents(events map[string]*model.StopAreaUpdateEvent, xmlResponse *siri.XMLStopMonitoringDelivery, tx *model.Transaction, monitoringRefMap map[string]struct{}) {
	xmlStopVisitCancellationEvents := xmlResponse.XMLMonitoredStopVisitCancellations()
	if len(xmlStopVisitCancellationEvents) == 0 {
		return
	}

	for _, xmlStopVisitCancellationEvent := range xmlStopVisitCancellationEvents {
		monitoringRefMap[xmlStopVisitCancellationEvent.MonitoringRef()] = struct{}{}

		stopAreaObjectId := model.NewObjectID(connector.Partner().Setting("remote_objectid_kind"), xmlStopVisitCancellationEvent.MonitoringRef())
		stopArea, ok := tx.Model().StopAreas().FindByObjectId(stopAreaObjectId)
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

func (connector *SIRIStopMonitoringSubscriptionCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringSubscriptionCollector"
	return event
}

func logSIRIDeleteSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIDeleteSubscriptionRequest) {
	logStashEvent["type"] = "DeleteStopMonitoringSubscription"
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	logStashEvent["subscriptionRef"] = request.SubscriptionRef
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["messageIdentifier"] = request.MessageIdentifier

	xml, err := request.BuildXML()
	if err != nil {
		logStashEvent["requestXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["requestXML"] = xml
}

func logXMLStopMonitoringDelivery(logStashEvent audit.LogStashEvent, notify *siri.XMLNotifyStopMonitoring) {
	logStashEvent["type"] = "NotifyStopMonitoringCollected"
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

func logXMLDeleteSubscriptionResponse(logStashEvent audit.LogStashEvent, response *siri.XMLDeleteSubscriptionResponse) {
	logStashEvent["responderRef"] = response.ResponderRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["responseXML"] = response.RawXML()

	var subscriptionIds []string
	for _, responseStatus := range response.ResponseStatus() {
		subscriptionIds = append(subscriptionIds, responseStatus.SubscriptionRef())
	}
	logStashEvent["subscriptionRefs"] = strings.Join(subscriptionIds, ", ")
}

func logXMLSubscriptionTerminatedResponse(logStashEvent audit.LogStashEvent, response *siri.XMLStopMonitoringSubscriptionTerminatedResponse) {
	logStashEvent["type"] = "StopMonitoringTerminatedSubscriptionCollected"
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
	logStashEvent["stopVisitUpdateEventIds"] = strings.Join(idArray, ", ")
	logStashEvent["stopVisitCancelledEventIds"] = strings.Join(cancelledIdArray, ", ")
}

func logMonitoringRefsFromMap(logStashEvent audit.LogStashEvent, refs map[string]struct{}) {
	refSlice := make([]string, len(refs))
	i := 0
	for monitoringRef := range refs {
		refSlice[i] = monitoringRef
		i++
	}
	logStashEvent["monitoringRefs"] = strings.Join(refSlice, ", ")
}
