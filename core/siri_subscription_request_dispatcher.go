package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SubscriptionRequestDispatcher interface {
	Dispatch(*sxml.XMLSubscriptionRequest, *audit.BigQueryMessage) (*siri.SIRISubscriptionResponse, error)
	CancelSubscription(*sxml.XMLDeleteSubscriptionRequest, *audit.BigQueryMessage) *siri.SIRIDeleteSubscriptionResponse
	HandleSubscriptionTerminatedNotification(*sxml.XMLSubscriptionTerminatedNotification)
	HandleNotifySubscriptionTerminated(*sxml.XMLNotifySubscriptionTerminated)
}

type SIRISubscriptionRequestDispatcherFactory struct{}

type SIRISubscriptionRequestDispatcher struct {
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

func (connector *SIRISubscriptionRequestDispatcher) Dispatch(request *sxml.XMLSubscriptionRequest, message *audit.BigQueryMessage) (*siri.SIRISubscriptionResponse, error) {
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
			return nil, fmt.Errorf("no EstimatedTimetableSubscriptionBroadcaster Connector")
		}

		response.ResponseStatus = smbc.(*SIRIEstimatedTimetableSubscriptionBroadcaster).HandleSubscriptionRequest(request, message)

		return &response, nil
	}

	if len(request.XMLSubscriptionPTTEntries()) > 0 {
		ptbc, ok := connector.Partner().Connector(SIRI_PRODUCTION_TIMETABLE_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("no ProductionTableSubscriptionBroadcaster Connector")
		}

		response.ResponseStatus = ptbc.(*SIRIProductionTimetableSubscriptionBroadcaster).HandleSubscriptionRequest(request, message)

		return &response, nil
	}

	if len(request.XMLSubscriptionVMEntries()) > 0 {
		vmbc, ok := connector.Partner().Connector(SIRI_VEHICLE_MONITORING_SUBSCRIPTION_BROADCASTER)
		if !ok {
			return nil, fmt.Errorf("no VehicleMonitoringSubscriptionBroadcaster Connector")
		}

		response.ResponseStatus = vmbc.(*SIRIVehicleMonitoringSubscriptionBroadcaster).HandleSubscriptionRequest(request, message)

		return &response, nil
	}

	return nil, fmt.Errorf("subscription not supported")
}

func (connector *SIRISubscriptionRequestDispatcher) CancelSubscription(r *sxml.XMLDeleteSubscriptionRequest, message *audit.BigQueryMessage) *siri.SIRIDeleteSubscriptionResponse {
	message.RequestIdentifier = r.MessageIdentifier()

	currentTime := connector.Clock().Now()
	resp := &siri.SIRIDeleteSubscriptionResponse{
		ResponderRef:      connector.Partner().RequestorRef(),
		RequestMessageRef: r.MessageIdentifier(),
		ResponseTimestamp: currentTime,
	}

	ignoreTerminate := connector.Partner().IgnoreTerminateSubscriptionsRequest()

	if r.CancelAll() {
		for _, subscription := range connector.Partner().Subscriptions().FindBroadcastSubscriptions() {
			responseStatus := &siri.SIRITerminationResponseStatus{
				SubscriberRef:     r.RequestorRef(),
				SubscriptionRef:   subscription.ExternalId(),
				ResponseTimestamp: currentTime,
				Status:            true,
			}

			if ignoreTerminate {
				responseStatus.ErrorType = "CapabilityNotSupportedError"
				responseStatus.ErrorText = "Subscription Termination is disabled for this Subscriber"
				responseStatus.Status = false
			}

			resp.ResponseStatus = append(resp.ResponseStatus, responseStatus)
		}

		if ignoreTerminate {
			message.Status = "Error"
			message.ErrorDetails = "Subscription Termination is disabled for this Subscriber"
		} else {
			connector.Partner().CancelBroadcastSubscriptions()
		}

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

	if ignoreTerminate {
		logger.Log.Debugf("Subscription Termination is disabled for partner %s", connector.Partner().Id())

		responseStatus.ErrorType = "CapabilityNotSupportedError"
		responseStatus.ErrorText = "Subscription Termination is disabled for this Subscriber"

		message.Status = "Error"
		message.ErrorDetails = responseStatus.ErrorText

		return resp
	}

	responseStatus.Status = true

	sub.Delete()
	return resp
}

func (connector *SIRISubscriptionRequestDispatcher) HandleSubscriptionTerminatedNotification(r *sxml.XMLSubscriptionTerminatedNotification) {
	connector.partner.Subscriptions().DeleteById(SubscriptionId(r.SubscriptionRef()))
}

func (connector *SIRISubscriptionRequestDispatcher) HandleNotifySubscriptionTerminated(r *sxml.XMLNotifySubscriptionTerminated) {
	connector.partner.Subscriptions().DeleteById(SubscriptionId(r.SubscriptionRef()))
}
