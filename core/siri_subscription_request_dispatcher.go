package core

import (
	"fmt"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SubscriptionRequestDispatcher interface {
	Dispatch(*siri.XMLSubscriptionRequest) (*siri.SIRISubscriptionResponse, error)
	CancelSubscription(*siri.XMLDeleteSubscriptionRequest) *siri.SIRIDeleteSubscriptionResponse
}

type SIRISubscriptionRequestDispatcherFactory struct{}

type SIRISubscriptionRequestDispatcher struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	xmlRequest siri.XMLSubscriptionRequest
}

func (factory *SIRISubscriptionRequestDispatcherFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func (factory *SIRISubscriptionRequestDispatcherFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRISubscriptionRequestDispatcher(partner)
}

func NewSIRISubscriptionRequestDispatcher(partner *Partner) *SIRISubscriptionRequestDispatcher {
	siriSubscriptionRequest := &SIRISubscriptionRequestDispatcher{}
	siriSubscriptionRequest.partner = partner

	return siriSubscriptionRequest
}

func (connector *SIRISubscriptionRequestDispatcher) Dispatch(request *siri.XMLSubscriptionRequest) (*siri.SIRISubscriptionResponse, error) {
	logStashEvent := connector.newLogStashEvent()

	logXMLSubscriptionRequest(logStashEvent, request)

	response := siri.SIRISubscriptionResponse{
		Address:            connector.Partner().Address(),
		ResponderRef:       connector.SIRIPartner().RequestorRef(),
		ResponseTimestamp:  connector.Clock().Now(),
		RequestMessageRef:  request.MessageIdentifier(),
		ServiceStartedTime: connector.Partner().Referential().StartedAt(),
	}

	if len(request.XMLSubscriptionGMEntries()) > 0 {
		gmbc, ok := connector.Partner().Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("No GeneralMessageSubscriptionBroadcaster Connector")
		}
		for _, sgm := range gmbc.(*SIRIGeneralMessageSubscriptionBroadcaster).HandleSubscriptionRequest(request) {
			response.ResponseStatus = append(response.ResponseStatus, sgm)
		}
		logSIRISubscriptionResponse(logStashEvent, &response, "GeneralMessageSubscriptionBroadcaster")
		logStashEvent["siriType"] = "GeneralMessageSubscriptionRequest"
		audit.CurrentLogStash().WriteEvent(logStashEvent)
		return &response, nil
	}

	if len(request.XMLSubscriptionSMEntries()) > 0 {
		smbc, ok := connector.Partner().Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("No StopMonitoringSubscriptionBroadcaster Connector")
		}
		for _, smr := range smbc.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleSubscriptionRequest(request) {
			response.ResponseStatus = append(response.ResponseStatus, smr)
		}
		logSIRISubscriptionResponse(logStashEvent, &response, "StopMonitoringSubscriptionBroadcaster")
		logStashEvent["siriType"] = "StopMonitoringSubscriptionRequest"
		audit.CurrentLogStash().WriteEvent(logStashEvent)
		return &response, nil
	}

	if len(request.XMLSubscriptionETTEntries()) > 0 {
		smbc, ok := connector.Partner().Connector(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("No EstimatedTimeTableSubscriptionBroadcaster Connector")
		}
		for _, smr := range smbc.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).HandleSubscriptionRequest(request) {
			response.ResponseStatus = append(response.ResponseStatus, smr)
		}
		logSIRISubscriptionResponse(logStashEvent, &response, "EstimatedTimeTableSubscriptionBroadcaster")
		logStashEvent["siriType"] = "EstimatedTimetableSubscriptionRequest"
		audit.CurrentLogStash().WriteEvent(logStashEvent)
		return &response, nil
	}

	return nil, fmt.Errorf("Subscription not supported")
}

func (connector *SIRISubscriptionRequestDispatcher) CancelSubscription(r *siri.XMLDeleteSubscriptionRequest) *siri.SIRIDeleteSubscriptionResponse {
	logStashEvent := connector.newLogStashEvent()

	logXMLCancelSubscriptionRequest(logStashEvent, r)

	currentTime := connector.Clock().Now()
	resp := &siri.SIRIDeleteSubscriptionResponse{
		ResponderRef:      connector.SIRIPartner().RequestorRef(),
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

	resp.ResponseStatus = append(resp.ResponseStatus, responseStatus)

	sub, ok := connector.Partner().Subscriptions().FindByExternalId(r.SubscriptionRef())
	if !ok {
		logger.Log.Debugf("Could not Unsubscribe to unknow subscription %v", r.SubscriptionRef())

		responseStatus.ErrorType = "InvalidDataReferencesError"
		responseStatus.ErrorText = fmt.Sprintf("Subscription not found: '%s'", r.SubscriptionRef())
		return resp
	}

	responseStatus.Status = true

	sub.Delete()
	return resp
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
