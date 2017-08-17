package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRIEstimatedTimetableRequestHandler struct {
	xmlRequest *siri.XMLGetEstimatedTimetable
	Partner    core.Partner
}

func (handler *SIRIEstimatedTimetableRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIEstimatedTimetableRequestHandler) ConnectorType() string {
	return core.SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER
}

func (handler *SIRIEstimatedTimetableRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("Estimated Timetable %s\n", handler.xmlRequest.MessageIdentifier())

	response := connector.(core.EstimatedTimetableBroadcaster).RequestLine(handler.xmlRequest)
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
