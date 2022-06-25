package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type CheckStatusServer interface {
	CheckStatus(*sxml.XMLCheckStatusRequest, *audit.BigQueryMessage) (*siri.SIRICheckStatusResponse, error)
}

type SIRICheckStatusServer struct {
	clock.ClockConsumer

	connector
}

type SIRICheckStatusServerFactory struct{}

func NewSIRICheckStatusServer(partner *Partner) *SIRICheckStatusServer {
	siriCheckStatusServer := &SIRICheckStatusServer{}
	siriCheckStatusServer.partner = partner
	return siriCheckStatusServer
}

func (connector *SIRICheckStatusServer) CheckStatus(request *sxml.XMLCheckStatusRequest, message *audit.BigQueryMessage) (*siri.SIRICheckStatusResponse, error) {
	response := &siri.SIRICheckStatusResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		RequestMessageRef:         request.MessageIdentifier(),
		ResponseMessageIdentifier: connector.Partner().NewResponseMessageIdentifier(),
		Status:                    true,
		ResponseTimestamp:         connector.Clock().Now(),
		ServiceStartedTime:        connector.Partner().StartedAt(),
	}

	message.RequestIdentifier = request.MessageIdentifier()
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	return response, nil
}

func (factory *SIRICheckStatusServerFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRICheckStatusServerFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusServer(partner)
}
