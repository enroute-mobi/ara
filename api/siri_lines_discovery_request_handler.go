package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRILinesDiscoveryRequestHandler struct {
	xmlRequest *siri.XMLLinesDiscoveryRequest
}

func (handler *SIRILinesDiscoveryRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRILinesDiscoveryRequestHandler) ConnectorType() string {
	return core.SIRI_LINES_DISCOVERY_REQUEST_BROADCASTER
}

func (handler *SIRILinesDiscoveryRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("LinesDiscovery %s\n", handler.xmlRequest.MessageIdentifier())

	response, _ := connector.(core.LinesDiscoveryRequestBroadcaster).Lines(handler.xmlRequest)
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
