package core

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SubscriptionRequestDispatcher interface {
	Dispatch(*siri.XMLSubscriptionRequest) (*siri.SIRIStopMonitoringSubscriptionResponse, error)
	CancelSubscription(*siri.XMLTerminatedSubscriptionRequest) *siri.SIRITerminatedSubscriptionResponse
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
		Address:           connector.Partner().Setting("local_url"),
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

	return nil, fmt.Errorf("Subscription not supported")
}

func (connector *SIRISubscriptionRequestDispatcher) CancelSubscription(r *siri.XMLTerminatedSubscriptionRequest) *siri.SIRITerminatedSubscriptionResponse {
	resp := &siri.SIRITerminatedSubscriptionResponse{
		ResponseTimestamp: connector.Clock().Now(),
		ResponderRef:      connector.SIRIPartner().RequestorRef(),
		Status:            true,
		SubscriberRef:     connector.SIRIPartner().RequestorRef(),
		SubscriptionRef:   r.SubscriptionRef(),
	}

	reg := regexp.MustCompile(`\w+:Subscription::([\w+-?]+):LOC`)
	matches := reg.FindStringSubmatch(strings.TrimSpace(r.SubscriptionRef()))

	if len(matches) == 0 {
		logger.Log.Debugf("Partner %s sent a TerminateSubscriptionRequest response with a wrong message format: %s\n", connector.Partner().Slug(), r.SubscriptionRef())

		resp.Status = false
		return resp
	}
	subscriptionId := matches[1]

	sub, ok := connector.Partner().Subscriptions().FindByExternalId(subscriptionId)
	if !ok {
		logger.Log.Debugf("Could not Unsubscribe to unknow subscription %v", resp.SubscriptionRef)

		resp.Status = false
		return resp
	}

	sub.Delete()
	return resp
}
