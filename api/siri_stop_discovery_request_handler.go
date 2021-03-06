package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIStopDiscoveryRequestHandler struct {
	xmlRequest *siri.XMLStopPointsDiscoveryRequest
}

func (handler *SIRIStopDiscoveryRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIStopDiscoveryRequestHandler) ConnectorType() string {
	return core.SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER
}

func (handler *SIRIStopDiscoveryRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("StopDiscovery %s\n", handler.xmlRequest.MessageIdentifier())

	tmp := connector.(*core.SIRIStopPointsDiscoveryRequestBroadcaster)
	response, _ := tmp.StopAreas(handler.xmlRequest)
	xmlResponse, err := response.BuildXML()
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
		return
	}

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	_, err = soapEnvelope.WriteTo(rw)
	if err != nil {
		siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), rw)
		return
	}
}
