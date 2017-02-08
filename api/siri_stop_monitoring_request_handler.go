package api

import (
	"fmt"
	"net/http"

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

func (handler *SIRIStopMonitoringRequestHandler) Respond(connector core.SIRIConnector, rw http.ResponseWriter) {
	logger.Log.Debugf("StopMonitoring %s\n", handler.xmlRequest.MessageIdentifier())

	// tx := handler.referential.NewTransaction()
	// defer tx.Close()

	// objectid := model.NewObjectID(connector.(core.SIRIConnector).Partner().Setting("remote_objectid_kind"), handler.xmlRequest.MonitoringRef())
	// stopArea, ok := tx.Model().StopAreas().FindByObjectId(objectid)
	// if !ok {
	// 	siriError("NotFound", "StopArea not found", rw)
	// 	return
	// }

	response := new(siri.SIRIStopMonitoringResponse)
	response.Address = connector.(core.SIRIConnector).Partner().Setting("Address")
	response.ProducerRef = "Edwig"
	response.RequestMessageRef = handler.xmlRequest.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.(core.SIRIConnector).SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = model.DefaultClock().Now()

	// Fill StopVisits

	xmlResponse := response.BuildXML()

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	_, err := soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
	}
}
