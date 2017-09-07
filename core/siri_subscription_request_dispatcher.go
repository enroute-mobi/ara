package core

import (
	"fmt"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SubscriptionRequestDispatcher interface {
	Dispatch(*siri.XMLSubscriptionRequest) (*siri.SIRIStopMonitoringSubscriptionResponse, error)
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

func (connector *SIRISubscriptionRequestDispatcher) Dispatch(request *siri.XMLSubscriptionRequest) (*siri.SIRIStopMonitoringSubscriptionResponse, error) {
	response := siri.SIRIStopMonitoringSubscriptionResponse{
		Address:           connector.Partner().Address(),
		ResponderRef:      connector.SIRIPartner().RequestorRef(),
		ResponseTimestamp: connector.Clock().Now(),
		RequestMessageRef: request.MessageIdentifier(),
	}

	if len(request.XMLSubscriptionGMEntries()) > 0 {
		gmbc, ok := connector.Partner().Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("No GeneralMessageSubscriptionBroadcaster Connector")
		}
		for _, sgm := range gmbc.(*SIRIGeneralMessageSubscriptionBroadcaster).HandleSubscriptionRequest(request) {
			response.ResponseStatus = append(response.ResponseStatus, sgm)
		}
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
		return &response, nil

	}

	return nil, fmt.Errorf("Subscription not supported")
}

func (connector *SIRISubscriptionRequestDispatcher) CancelSubscription(r *siri.XMLDeleteSubscriptionRequest) *siri.SIRIDeleteSubscriptionResponse {
	currentTime := connector.Clock().Now()
	resp := &siri.SIRIDeleteSubscriptionResponse{
		ResponderRef:      connector.SIRIPartner().RequestorRef(),
		RequestMessageRef: r.MessageIdentifier(),
		ResponseTimestamp: currentTime,
	}

	if r.CancelAll() {
		for _, subscription := range connector.Partner().Subscriptions().FindAll() {
			responseStatus := &siri.SIRITerminationResponseStatus{
				SubscriberRef:     r.RequestorRef(),
				SubscriptionRef:   subscription.ExternalId(),
				ResponseTimestamp: currentTime,
				Status:            true,
			}
			resp.ResponseStatus = append(resp.ResponseStatus, responseStatus)
		}
		connector.Partner().CancelSubscriptions()
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
