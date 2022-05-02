package core

import (
	"fmt"

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
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
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

		return &response, nil
	}

	if len(request.XMLSubscriptionSMEntries()) > 0 {
		smbc, ok := connector.Partner().Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("no StopMonitoringSubscriptionBroadcaster Connector")
		}

		response.ResponseStatus = smbc.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleSubscriptionRequest(request, message)

		return &response, nil
	}

	if len(request.XMLSubscriptionETTEntries()) > 0 {
		smbc, ok := connector.Partner().Connector(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("no EstimatedTimeTableSubscriptionBroadcaster Connector")
		}

		response.ResponseStatus = smbc.(*SIRIEstimatedTimeTableSubscriptionBroadcaster).HandleSubscriptionRequest(request, message)

		return &response, nil
	}

	return nil, fmt.Errorf("subscription not supported")
}

func (connector *SIRISubscriptionRequestDispatcher) CancelSubscription(r *siri.XMLDeleteSubscriptionRequest, message *audit.BigQueryMessage) *siri.SIRIDeleteSubscriptionResponse {
	message.RequestIdentifier = r.MessageIdentifier()

	currentTime := connector.Clock().Now()
	resp := &siri.SIRIDeleteSubscriptionResponse{
		ResponderRef:      connector.Partner().RequestorRef(),
		RequestMessageRef: r.MessageIdentifier(),
		ResponseTimestamp: currentTime,
	}

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

		responseStatus.ErrorType = "UnknownSubscriptionError"
		responseStatus.ErrorText = fmt.Sprintf("Subscription not found: '%s'", r.SubscriptionRef())

		message.Status = "Error"
		message.ErrorDetails = responseStatus.ErrorText

		return resp
	}

	responseStatus.Status = true

	sub.Delete()
	return resp
}

func (connector *SIRISubscriptionRequestDispatcher) HandleSubscriptionTerminatedNotification(r *siri.XMLSubscriptionTerminatedNotification) {
	connector.partner.Subscriptions().DeleteById(SubscriptionId(r.SubscriptionRef()))
}

func (connector *SIRISubscriptionRequestDispatcher) HandleNotifySubscriptionTerminated(r *siri.XMLNotifySubscriptionTerminated) {
	connector.partner.Subscriptions().DeleteById(SubscriptionId(r.SubscriptionRef()))
}
