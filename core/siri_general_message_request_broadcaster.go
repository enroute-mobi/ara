package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type GeneralMessageRequestBroadcaster interface {
	Situations(*siri.XMLGetGeneralMessage, *audit.BigQueryMessage) (*siri.SIRIGeneralMessageResponse, error)
}

type SIRIGeneralMessageRequestBroadcaster struct {
	clock.ClockConsumer
	uuid.UUIDConsumer
	connector
}

type SIRIGeneralMessageRequestBroadcasterFactory struct{}

func NewSIRIGeneralMessageRequestBroadcaster(partner *Partner) *SIRIGeneralMessageRequestBroadcaster {
	siriGeneralMessageRequestBroadcaster := &SIRIGeneralMessageRequestBroadcaster{}
	siriGeneralMessageRequestBroadcaster.partner = partner
	return siriGeneralMessageRequestBroadcaster
}

func (connector *SIRIGeneralMessageRequestBroadcaster) Situations(request *siri.XMLGetGeneralMessage, message *audit.BigQueryMessage) (*siri.SIRIGeneralMessageResponse, error) {
	response := &siri.SIRIGeneralMessageResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
	}

	response.SIRIGeneralMessageDelivery = connector.getGeneralMessageDelivery(&request.XMLGeneralMessageRequest)

	if !response.SIRIGeneralMessageDelivery.Status {
		message.Status = "Error"
		message.ErrorDetails = response.SIRIGeneralMessageDelivery.ErrorString()
	}
	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	return response, nil
}

func (connector *SIRIGeneralMessageRequestBroadcaster) getGeneralMessageDelivery(request *siri.XMLGeneralMessageRequest) siri.SIRIGeneralMessageDelivery {
	delivery := siri.SIRIGeneralMessageDelivery{
		RequestMessageRef: request.MessageIdentifier(),
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	builder := NewBroadcastGeneralMessageBuilder(connector.Partner(), SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER)
	builder.InfoChannelRef = request.InfoChannelRef()
	builder.SetLineRef(request.LineRef())
	builder.SetStopPointRef(request.StopPointRef())

	situations := connector.partner.Model().Situations().FindAll()
	for i := range situations {
		siriGeneralMessage := builder.BuildGeneralMessage(situations[i])
		if siriGeneralMessage == nil {
			continue
		}
		delivery.GeneralMessages = append(delivery.GeneralMessages, siriGeneralMessage)
	}

	return delivery
}

func (factory *SIRIGeneralMessageRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIGeneralMessageRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageRequestBroadcaster(partner)
}
