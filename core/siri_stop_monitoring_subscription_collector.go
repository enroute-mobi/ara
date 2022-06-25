package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type StopMonitoringSubscriptionCollector interface {
	state.Stopable
	state.Startable

	RequestStopAreaUpdate(request *StopAreaUpdateRequest)
	HandleNotifyStopMonitoring(delivery *sxml.XMLNotifyStopMonitoring)
}

type SIRIStopMonitoringSubscriptionCollector struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector

	deletedSubscriptions     *DeletedSubscriptions
	stopMonitoringSubscriber SIRIStopMonitoringSubscriber
	updateSubscriber         UpdateSubscriber
}

type SIRIStopMonitoringSubscriptionCollectorFactory struct{}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringSubscriptionCollector(partner)
}

func (factory *SIRIStopMonitoringSubscriptionCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func NewSIRIStopMonitoringSubscriptionCollector(partner *Partner) *SIRIStopMonitoringSubscriptionCollector {
	connector := &SIRIStopMonitoringSubscriptionCollector{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind()
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
	connector.deletedSubscriptions = NewDeletedSubscriptions()
	connector.stopMonitoringSubscriber.Start()
}

func (connector *SIRIStopMonitoringSubscriptionCollector) RequestStopAreaUpdate(request *StopAreaUpdateRequest) {
	stopArea, ok := connector.partner.Model().StopAreas().Find(request.StopAreaId())
	if !ok {
		logger.Log.Debugf("StopAreaUpdateRequest in StopMonitoring SubscriptionCollector for unknown StopArea %v", request.StopAreaId())
		return
	}

	stopAreaObjectid, ok := stopArea.ObjectID(connector.remoteObjectidKind)
	if !ok {
		logger.Log.Debugf("Requested stopArea %v doesn't have and objectId of kind %v", request.StopAreaId(), connector.remoteObjectidKind)
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

func (connector *SIRIStopMonitoringSubscriptionCollector) HandleNotifyStopMonitoring(notify *sxml.XMLNotifyStopMonitoring) {
	monitoringRefMap := make(map[string]struct{})
	subscriptionErrors := make(map[string]string)
	subToDelete := make(map[string]struct{})

	for _, delivery := range notify.StopMonitoringDeliveries() {
		subscriptionId := delivery.SubscriptionRef()

		subscription, ok := connector.Partner().Subscriptions().Find(SubscriptionId(subscriptionId))
		if !ok {
			logger.Log.Debugf("Partner %s sent a StopVisitNotify response to a non existant subscription of id: %s\n", connector.Partner().Slug(), subscriptionId)
			subscriptionErrors[subscriptionId] = "Non existant subscription of id %s"
			if !connector.deletedSubscriptions.AlreadySend(subscriptionId) {
				subToDelete[subscriptionId] = struct{}{}
			}
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
			originStopAreaObjectId = model.NewObjectID(connector.remoteObjectidKind, delivery.MonitoringRef())
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

	for subId := range subToDelete {
		connector.cancelSubscription(subId)
	}
}

func (connector *SIRIStopMonitoringSubscriptionCollector) cancelSubscription(subId string) {
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	request := &siri.SIRIDeleteSubscriptionRequest{
		RequestTimestamp:  connector.Clock().Now(),
		SubscriptionRef:   subId,
		RequestorRef:      connector.partner.ProducerRef(),
		MessageIdentifier: connector.Partner().NewMessageIdentifier(),
	}
	logSIRIDeleteSubscriptionRequest(message, request, "StopMonitoringSubscriptionCollector", connector.Partner().SIRIEnvelopeType())

	startTime := connector.Clock().Now()
	response, err := connector.Partner().SIRIClient().DeleteSubscription(request)

	responseTime := connector.Clock().Since(startTime)
	message.ProcessingTime = responseTime.Seconds()

	if err != nil {
		logger.Log.Debugf("Error while terminating subcription with id : %v error : %v", subId, err.Error())
		e := fmt.Sprintf("Error during DeleteSubscription: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}
	logXMLDeleteSubscriptionResponse(message, response)
}

func (connector *SIRIStopMonitoringSubscriptionCollector) broadcastUpdateEvents(events *CollectUpdateEvents) {
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

func (connector *SIRIStopMonitoringSubscriptionCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func logSIRIDeleteSubscriptionRequest(message *audit.BigQueryMessage, request *siri.SIRIDeleteSubscriptionRequest, subType, envelopeType string) {
	message.Type = "DeleteSubscriptionRequest"
	message.RequestIdentifier = request.MessageIdentifier
	message.SubscriptionIdentifiers = []string{request.SubscriptionRef}

	xml, err := request.BuildXML(envelopeType)
	if err != nil {
		return
	}
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLDeleteSubscriptionResponse(message *audit.BigQueryMessage, response *sxml.XMLDeleteSubscriptionResponse) {
	var i int
	for _, responseStatus := range response.ResponseStatus() {
		if !responseStatus.Status() {
			i++
		}
	}

	if i > 0 {
		message.Status = "Error"
		message.ErrorDetails = fmt.Sprintf("%d ResponseStatus returned false", i)
	}
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
