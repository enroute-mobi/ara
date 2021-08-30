package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type SubscriptionRequestDispatcher interface {
	Dispatch(*siri.XMLSubscriptionRequest, *audit.BigQueryMessage) (*siri.SIRISubscriptionResponse, error)
	CancelSubscription(*siri.XMLDeleteSubscriptionRequest, *audit.BigQueryMessage) *siri.SIRIDeleteSubscriptionResponse
	HandleSubscriptionTerminatedNotification(*siri.XMLSubscriptionTerminatedNotification)
	HandleNotifySubscriptionTerminated(*siri.XMLNotifySubscriptionTerminated)
}

type SIRISubscriptionRequestDispatcherFactory struct{}

type SIRISubscriptionRequestDispatcher struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	connector
}

func (factory *SIRISubscriptionRequestDispatcherFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	apiPartner.ValidatePresenceOfSetting(REMOTE_URL)
	apiPartner.ValidatePresenceOfSetting(REMOTE_CREDENTIAL)
}

func (factory *SIRISubscriptionRequestDispatcherFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRISubscriptionRequestDispatcher(partner)
}

func NewSIRISubscriptionRequestDispatcher(partner *Partner) *SIRISubscriptionRequestDispatcher {
	siriSubscriptionRequest := &SIRISubscriptionRequestDispatcher{}
	siriSubscriptionRequest.partner = partner

	return siriSubscriptionRequest
}

func (connector *SIRISubscriptionRequestDispatcher) Dispatch(request *siri.XMLSubscriptionRequest, message *audit.BigQueryMessage) (*siri.SIRISubscriptionResponse, error) {
	logStashEvent := connector.newLogStashEvent()

	logXMLSubscriptionRequest(logStashEvent, request)

	response := siri.SIRISubscriptionResponse{
		Address:            connector.Partner().Address(),
		ResponderRef:       connector.Partner().RequestorRef(),
		ResponseTimestamp:  connector.Clock().Now(),
		RequestMessageRef:  request.MessageIdentifier(),
		ServiceStartedTime: connector.Partner().Referential().StartedAt(),
	}

	message.RequestIdentifier = request.MessageIdentifier()

	if len(request.XMLSubscriptionGMEntries()) > 0 {
		gmbc, ok := connector.Partner().Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("no GeneralMessageSubscriptionBroadcaster Connector")
		}

		response.ResponseStatus = gmbc.(*SIRIGeneralMessageSubscriptionBroadcaster).HandleSubscriptionRequest(request, message)

		logSIRISubscriptionResponse(logStashEvent, &response, "GeneralMessageSubscriptionBroadcaster")
		logStashEvent["siriType"] = "GeneralMessageSubscriptionRequest"
		audit.CurrentLogStash().WriteEvent(logStashEvent)
		return &response, nil
	}

	if len(request.XMLSubscriptionSMEntries()) > 0 {
		smbc, ok := connector.Partner().Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("no StopMonitoringSubscriptionBroadcaster Connector")
		}

		response.ResponseStatus = smbc.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleSubscriptionRequest(request, message)

		logSIRISubscriptionResponse(logStashEvent, &response, "StopMonitoringSubscriptionBroadcaster")
		logStashEvent["siriType"] = "StopMonitoringSubscriptionRequest"
		audit.CurrentLogStash().WriteEvent(logStashEvent)
		return &response, nil
	}

	if len(request.XMLSubscriptionETTEntries()) > 0 {
		smbc, ok := connector.Partner().Connector(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("no EstimatedTimeTableSubscriptionBroadcaster Connector")
		}

		response.ResponseStatus = smbc.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).HandleSubscriptionRequest(request, message)

		logSIRISubscriptionResponse(logStashEvent, &response, "EstimatedTimeTableSubscriptionBroadcaster")
		logStashEvent["siriType"] = "EstimatedTimetableSubscriptionRequest"
		audit.CurrentLogStash().WriteEvent(logStashEvent)
		return &response, nil
	}

	return nil, fmt.Errorf("subscription not supported")
}

func (connector *SIRISubscriptionRequestDispatcher) CancelSubscription(r *siri.XMLDeleteSubscriptionRequest, message *audit.BigQueryMessage) *siri.SIRIDeleteSubscriptionResponse {
	logStashEvent := connector.newLogStashEvent()

	message.RequestIdentifier = r.MessageIdentifier()

	logXMLCancelSubscriptionRequest(logStashEvent, r)

	currentTime := connector.Clock().Now()
	resp := &siri.SIRIDeleteSubscriptionResponse{
		ResponderRef:      connector.Partner().RequestorRef(),
		RequestMessageRef: r.MessageIdentifier(),
		ResponseTimestamp: currentTime,
	}

	defer func() {
		logSIRICancelSubscriptionResponse(logStashEvent, resp)
		audit.CurrentLogStash().WriteEvent(logStashEvent)
	}()

	if r.CancelAll() {
		for _, subscription := range connector.Partner().Subscriptions().FindBroadcastSubscriptions() {
			responseStatus := &siri.SIRITerminationResponseStatus{
				SubscriberRef:     r.RequestorRef(),
				SubscriptionRef:   subscription.ExternalId(),
				ResponseTimestamp: currentTime,
				Status:            true,
			}
			resp.ResponseStatus = append(resp.ResponseStatus, responseStatus)
		}
		connector.Partner().CancelBroadcastSubscriptions()
		return resp
	}

	responseStatus := &siri.SIRITerminationResponseStatus{
		SubscriberRef:     r.RequestorRef(),
		SubscriptionRef:   r.SubscriptionRef(),
		ResponseTimestamp: currentTime,
	}

	message.SubscriptionIdentifiers = []string{r.SubscriptionRef()}

	resp.ResponseStatus = append(resp.ResponseStatus, responseStatus)

	sub, ok := connector.Partner().Subscriptions().FindByExternalId(r.SubscriptionRef())
	if !ok {
		logger.Log.Debugf("Could not Unsubscribe to unknow subscription %v", r.SubscriptionRef())

		responseStatus.ErrorType = "InvalidDataReferencesError"
		responseStatus.ErrorText = fmt.Sprintf("Subscription not found: '%s'", r.SubscriptionRef())

		message.Status = "Error"
		message.ErrorDetails = responseStatus.ErrorText

		return resp
	}

	responseStatus.Status = true

	sub.Delete()
	return resp
}

func (connector *SIRISubscriptionRequestDispatcher) HandleSubscriptionTerminatedNotification(response *siri.XMLSubscriptionTerminatedNotification) {
	logStashEvent := make(audit.LogStashEvent)

	logXMLSubscriptionTerminatedNotification(logStashEvent, response)

	connector.partner.Subscriptions().DeleteById(SubscriptionId(response.SubscriptionRef()))
	audit.CurrentLogStash().WriteEvent(logStashEvent)
}

func (connector *SIRISubscriptionRequestDispatcher) HandleNotifySubscriptionTerminated(r *siri.XMLNotifySubscriptionTerminated) {
	logStashEvent := connector.newLogStashEvent()

	logXMLNotifySubscriptionTerminated(logStashEvent, r)

	connector.partner.Subscriptions().DeleteById(SubscriptionId(r.SubscriptionRef()))
	audit.CurrentLogStash().WriteEvent(logStashEvent)
}

func (connector *SIRISubscriptionRequestDispatcher) newLogStashEvent() audit.LogStashEvent {
	return connector.partner.NewLogStashEvent()
}

func logXMLSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.XMLSubscriptionRequest) {
	logStashEvent["consumerAddress"] = request.ConsumerAddress()
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRISubscriptionResponse(logStashEvent audit.LogStashEvent, response *siri.SIRISubscriptionResponse, connector string) {
	logStashEvent["connector"] = connector
	logStashEvent["address"] = response.Address
	logStashEvent["responderRef"] = response.ResponderRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["serviceStartedTime"] = response.ServiceStartedTime.String()
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}

func logXMLCancelSubscriptionRequest(logStashEvent audit.LogStashEvent, request *siri.XMLDeleteSubscriptionRequest) {
	logStashEvent["siriType"] = "DeleteSubscriptionResponse"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	if request.CancelAll() {
		logStashEvent["subscriptionToCancel"] = "All"
	} else {
		logStashEvent["subscriptionToCancel"] = request.SubscriptionRef()
	}
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRICancelSubscriptionResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIDeleteSubscriptionResponse) {
	logStashEvent["responderRef"] = response.ResponderRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}

func logXMLNotifySubscriptionTerminated(logStashEvent audit.LogStashEvent, notify *siri.XMLNotifySubscriptionTerminated) {
	logStashEvent["siriType"] = "NotifySubscriptionTerminated"
	logStashEvent["address"] = notify.Address()
	logStashEvent["producerRef"] = notify.ProducerRef()
	logStashEvent["requestMessageRef"] = notify.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = notify.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = notify.ResponseTimestamp().String()
	logStashEvent["subscriberRef"] = notify.SubscriberRef()
	logStashEvent["subscriptionRef"] = notify.SubscriptionRef()
	logStashEvent["responseXML"] = notify.RawXML()
}

func logXMLSubscriptionTerminatedNotification(logStashEvent audit.LogStashEvent, response *siri.XMLSubscriptionTerminatedNotification) {
	logStashEvent["siriType"] = "TerminatedSubscriptionNotification"
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["subscriberRef"] = response.SubscriberRef()
	logStashEvent["subscriptionRef"] = response.SubscriptionRef()
	logStashEvent["responseXML"] = response.RawXML()

	logStashEvent["errorType"] = response.ErrorType()
	if response.ErrorType() == "OtherError" {
		logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber())
	}
	logStashEvent["errorDescription"] = response.ErrorDescription()
}
