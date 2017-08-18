package core

import (
	"fmt"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SubscriptionRequestDispatcher interface {
	Dispatch(*siri.XMLSubscriptionRequest) (*siri.SIRIStopMonitoringSubscriptionResponse, error)
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
