package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SubscriptionRequest interface {
	Dispatch(*siri.XMLSubscriptionRequest) *siri.SIRIStopMonitoringSubscriptionResponse
}

type SIRISubscriptionRequestFactory struct{}

type SIRISubscriptionRequest struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	xmlRequest siri.XMLSubscriptionRequest
}

func (factory *SIRISubscriptionRequestFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func (factory *SIRISubscriptionRequestFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRISubscriptionRequest(partner)
}

func NewSIRISubscriptionRequest(partner *Partner) *SIRISubscriptionRequest {
	siriSubscriptionRequest := &SIRISubscriptionRequest{}
	siriSubscriptionRequest.partner = partner

	return siriSubscriptionRequest
}

func (connector *SIRISubscriptionRequest) Dispatch(request *siri.XMLSubscriptionRequest) *siri.SIRIStopMonitoringSubscriptionResponse {
	response := siri.SIRIStopMonitoringSubscriptionResponse{
		Address:           connector.Partner().Setting("local_url"),
		ResponderRef:      connector.SIRIPartner().RequestorRef(),
		ResponseTimestamp: connector.Clock().Now(),
		RequestMessageRef: request.MessageIdentifier(),
	}

	gms := request.XMLSubscriptionGMEntries()
	sms := request.XMLSubscriptionSMEntries()

	gmbc, ok := connector.Partner().Connector(SIRI_GENERAL_MESSAGE_SUBSCRIPTION_BROADCASTER)

	if ok {
		for _, sgm := range gmbc.(*SIRIGeneralMessageSubscriptionBroadcaster).HandleSubscriptionRequest(gms) {
			response.ResponseStatus = append(response.ResponseStatus, sgm)
		}
	}

	smbc, ok := connector.Partner().Connector(SIRI_STOP_MONITORING_SUBSCRIPTION_BROADCASTER)
	if ok {
		for _, smr := range smbc.(*SIRIStopMonitoringSubscriptionBroadcaster).HandleSubscriptionRequest(sms) {
			response.ResponseStatus = append(response.ResponseStatus, smr)
		}
	}
	return &response
}
