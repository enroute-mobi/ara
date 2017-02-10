package api

import (
	"fmt"
	"net/http"

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

func (handler *SIRICheckStatusRequestHandler) Respond(connector core.SIRIConnector, rw http.ResponseWriter) {
	logger.Log.Debugf("CheckStatus %s\n", handler.xmlRequest.MessageIdentifier())

	response := new(siri.SIRICheckStatusResponse)
	response.Address = connector.(core.SIRIConnector).Partner().Setting("Address")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = handler.xmlRequest.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.(core.SIRIConnector).SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = model.DefaultClock().Now()
	response.ServiceStartedTime = handler.referential.StartedAt()

	xmlResponse := response.BuildXML()

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	_, err := soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
	}
}
