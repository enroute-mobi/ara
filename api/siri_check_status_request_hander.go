package api

import (
	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRICheckStatusRequestHandler struct {
	referential *core.Referential
	xmlRequest  *siri.XMLCheckStatusRequest
}

func (handler *SIRICheckStatusRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRICheckStatusRequestHandler) ConnectorType() string {
	return "siri-check-status-client"
}

func (handler *SIRICheckStatusRequestHandler) XMLResponse(connector core.Connector) string {
	logger.Log.Debugf("CheckStatus %s\n", handler.xmlRequest.MessageIdentifier())

	response := new(siri.SIRICheckStatusResponse)
	response.Address = connector.(*core.SIRICheckStatusClient).Partner().Setting("Address")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = handler.xmlRequest.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.(*core.SIRICheckStatusClient).SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = model.DefaultClock().Now()
	response.ServiceStartedTime = handler.referential.StartedAt()

	return response.BuildXML()
}
