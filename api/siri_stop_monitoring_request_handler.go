package api

import (
	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIStopMonitoringRequestHandler struct {
	referential *core.Referential
	xmlRequest  *siri.XMLStopMonitoringRequest
}

func (handler *SIRIStopMonitoringRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIStopMonitoringRequestHandler) ConnectorType() string {
	return "siri-stop-monitoring-request-collector"
}

func (handler *SIRIStopMonitoringRequestHandler) XMLResponse(connector core.Connector) string {
	logger.Log.Debugf("StopMonitoring %s\n", handler.xmlRequest.MessageIdentifier())

	response := new(siri.SIRIStopMonitoringResponse)
	response.Address = connector.(*core.SIRICheckStatusClient).Partner().Setting("Address")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = handler.xmlRequest.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.(*core.SIRICheckStatusClient).SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = model.DefaultClock().Now()

	return response.BuildXML()
}
