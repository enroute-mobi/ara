package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopMonitoringSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
	HandleNotifyStopMonitoring(delivery *siri.XMLNotifyStopMonitoring)
}

type SIRIStopMonitoringSubscriptionCollector struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	stopMonitoringSubscriber SIRIStopMonitoringSubscriber
	updateSubscriber         UpdateSubscriber
}

type SIRIStopMonitoringSubscriptionCollectorFactory struct{}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	apiPartner.ValidatePresenceOfSetting(REMOTE_URL)
	apiPartner.ValidatePresenceOfSetting(REMOTE_CREDENTIAL)
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIStopMonitoringSubscriptionCollector(partner *Partner) *SIRIStopMonitoringSubscriptionCollector {
	connector := &SIRIStopMonitoringSubscriptionCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.updateSubscriber = manager.BroadcastUpdateEvent
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

	objectidKind := connector.Partner().Setting(REMOTE_OBJECTID_KIND)
	stopAreaObjectid, ok := stopArea.ObjectID(objectidKind)
	if !ok {
		logger.Log.Debugf("Requested stopArea %v doesn't have and objectId of kind %v", request.StopAreaId(), objectidKind)
		return
	}

	// Try to find a Subscription with the resource
	subscriptions := connector.partner.Subscriptions().FindByResourceId(stopAreaObjectid.String(), "StopMonitoringCollect")
	if len(subscriptions) > 0 {
		for _, subscription := range subscriptions {
			resource := subscription.Resource(stopAreaObjectid)
			if resource == nil { // Should never happen
				logger.Log.Debugf("Can't find resource in subscription after Subscriptions#FindByResourceId")
				return
			}
			if !resource.SubscribedAt.IsZero() {
				resource.SubscribedUntil = connector.Clock().Now().Add(2 * time.Minute)
			}
		}
		return
	}

	// Else we find or create a subscription to add the resource
	newSubscription := connector.partner.Subscriptions().FindOrCreateByKind("StopMonitoringCollect")
	ref := model.Reference{
		ObjectId: &stopAreaObjectid,
		Type:     "StopArea",
	}

	newSubscription.CreateAddNewResource(ref)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) SetStopMonitoringSubscriber(stopMonitoringSubscriber SIRIStopMonitoringSubscriber) {
	connector.stopMonitoringSubscriber = stopMonitoringSubscriber
}

func (connector *SIRIStopMonitoringSubscriptionCollector) HandleNotifyStopMonitoring(notify *siri.XMLNotifyStopMonitoring) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)
	monitoringRefMap := make(map[string]struct{})
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})

	logXMLNotifyStopMonitoring(logStashEvent, notify)

	for _, delivery := range notify.StopMonitoringDeliveries() {
		if connector.Partner().LogSubscriptionStopMonitoringDeliveries() {
			deliveryLogStashEvent := connector.newLogStashEvent()
			logXMLSubscriptionStopMonitoringDelivery(deliveryLogStashEvent, notify.ResponseMessageIdentifier(), delivery)
			audit.CurrentLogStash().WriteEvent(deliveryLogStashEvent)
		}

		subscriptionId := delivery.SubscriptionRef()

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Debugf("Partner %s sent a StopVisitNotify response to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			subToDelete[delivery.SubscriptionRef()] = struct{}{}
			continue
		}
		if subscription.Kind() != "StopMonitoringCollect" {
			logger.Log.Debugf("Partner %s sent a StopVisitNotify response to a subscription with kind: %s\n", connector.Partner().Slug(), subscription.Kind())
			subscriptionErrors[subscriptionId] = "Subscription of id %s is not a subscription of kind StopMonitoringCollect"
			continue
		}

		originStopAreaObjectId := model.ObjectID{}
		resource := subscription.UniqueResource()
		if resource != nil {
			originStopAreaObjectId = *resource.Reference.ObjectId
		} else if delivery.MonitoringRef() != "" {
			originStopAreaObjectId = model.NewObjectID(connector.Partner().Setting(REMOTE_OBJECTID_KIND), delivery.MonitoringRef())
		}

		builder := NewStopMonitoringUpdateEventBuilder(connector.partner, originStopAreaObjectId)
		builder.SetUpdateEvents(delivery.XMLMonitoredStopVisits())
		builder.SetStopVisitCancellationEvents(delivery)
		updateEvents := builder.UpdateEvents()

		// Copy MonitoringRefs for global log
		for k := range updateEvents.MonitoringRefs {
			monitoringRefMap[k] = struct{}{}
		}

		connector.broadcastUpdateEvents(&updateEvents)
	}

	logMonitoringRefsFromMap(logStashEvent, monitoringRefMap)
	if len(subscriptionErrors) != 0 {
		logSubscriptionErrorsFromMap(logStashEvent, subscriptionErrors)
	}

	for subId := range subToDelete {
		connector.cancelSubscription(subId)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) cancelSubscription(subId string) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	request := &siri.SIRIDeleteSubscriptionRequest{
		RequestTimestamp:  connector.Clock().Now(),
		SubscriptionRef:   subId,
		RequestorRef:      connector.partner.ProducerRef(),
		MessageIdentifier: connector.Partner().IdentifierGenerator(MESSAGE_IDENTIFIER).NewMessageIdentifier(),
	}
	logSIRIDeleteSubscriptionRequest(logStashEvent, message, request, "StopMonitoringSubscriptionCollector")

	startTime := connector.Clock().Now()
	response, err := connector.Partner().SOAPClient().DeleteSubscription(request)

	responseTime := connector.Clock().Since(startTime)
	logStashEvent["responseTime"] = responseTime.String()
	message.ProcessingTime = responseTime.Seconds()

	if err != nil {
		logger.Log.Debugf("Error while terminating subcription with id : %v error : %v", subId, err.Error())
		e := fmt.Sprintf("Error during DeleteSubscription: %v", err)
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = e
		message.Status = "Error"
		message.ErrorDetails = e
		return
	}
	logXMLDeleteSubscriptionResponse(logStashEvent, message, response)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) broadcastUpdateEvents(events *StopMonitoringUpdateEvents) {
	if connector.updateSubscriber == nil {
		return
	}
	for _, e := range events.StopAreas {
		connector.updateSubscriber(e)
	}
	for _, e := range events.Lines {
		connector.updateSubscriber(e)
	}
	for _, e := range events.VehicleJourneys {
		connector.updateSubscriber(e)
	}
	for _, es := range events.StopVisits { // Stopvisits are map[MonitoringRef]map[ItemIdentifier]event
		for _, e := range es {
			connector.updateSubscriber(e)
		}
	}
	for _, e := range events.Cancellations {
		connector.updateSubscriber(e)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringSubscriptionCollector"
	return event
}

func (connector *SIRIStopMonitoringSubscriptionCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func logSIRIDeleteSubscriptionRequest(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, request *siri.SIRIDeleteSubscriptionRequest, subType string) {
	logStashEvent["siriType"] = "DeleteSubscriptionRequest" // This function is also used on GM delete subscription
	logStashEvent["connector"] = subType
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

	message.Type = "DeleteSubscriptionRequest"
	message.RequestIdentifier = request.MessageIdentifier
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))

	message.SubscriptionIdentifiers = []string{request.SubscriptionRef}
}

func logXMLNotifyStopMonitoring(logStashEvent audit.LogStashEvent, notify *siri.XMLNotifyStopMonitoring) {
	logStashEvent["siriType"] = "CollectedNotifyStopMonitoring"
	logStashEvent["producerRef"] = notify.ProducerRef()
	logStashEvent["requestMessageRef"] = notify.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = notify.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = notify.ResponseTimestamp().String()
	logStashEvent["responseXML"] = notify.RawXML()

	status := "true"
	errorCount := 0
	for _, delivery := range notify.StopMonitoringDeliveries() {
		if !delivery.Status() {
			status = "false"
			errorCount++
		}
	}
	logStashEvent["status"] = status
	logStashEvent["errorCount"] = strconv.Itoa(errorCount)
}

func logXMLSubscriptionStopMonitoringDelivery(logStashEvent audit.LogStashEvent, parent string, delivery *siri.XMLNotifyStopMonitoringDelivery) {
	logStashEvent["siriType"] = "CollectedNotifyStopMonitoringDelivery"
	logStashEvent["parentMessageIdentifier"] = parent
	logStashEvent["monitoringRef"] = delivery.MonitoringRef()
	logStashEvent["requestMessageRef"] = delivery.RequestMessageRef()
	logStashEvent["subscriberRef"] = delivery.SubscriberRef()
	logStashEvent["subscriptionRef"] = delivery.SubscriptionRef()
	logStashEvent["responseTimestamp"] = delivery.ResponseTimestamp().String()

	logStashEvent["status"] = strconv.FormatBool(delivery.Status())
	if !delivery.Status() {
		logStashEvent["errorType"] = delivery.ErrorType()
		if delivery.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber())
		}
		logStashEvent["errorText"] = delivery.ErrorText()
		logStashEvent["errorDescription"] = delivery.ErrorDescription()
	}
}

func logXMLDeleteSubscriptionResponse(logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage, response *siri.XMLDeleteSubscriptionResponse) {
	logStashEvent["responderRef"] = response.ResponderRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["responseXML"] = response.RawXML()

	var subscriptionIds []string
	for _, responseStatus := range response.ResponseStatus() {
		subscriptionIds = append(subscriptionIds, responseStatus.SubscriptionRef())
	}
	logStashEvent["subscriptionRefs"] = strings.Join(subscriptionIds, ", ")

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
	// TODO no ResponseMessageIdentifier() method in XMLDeleteSubscriptionResponse
	// message.ResponseIdentifier = response.ResponseMessageIdentifier()
}

func logMonitoringRefsFromMap(logStashEvent audit.LogStashEvent, refs map[string]struct{}) {
	logMonitoringRefs(logStashEvent, nil, refs)
}

func logSubscriptionErrorsFromMap(logStashEvent audit.LogStashEvent, errors map[string]string) {
	errSlice := make([]string, len(errors))
	i := 0
	for subId, err := range errors {
		errSlice[i] = fmt.Sprintf(err, subId)
		i++
	}
	logStashEvent["notificationErrors"] = strings.Join(errSlice, ", ")
}
