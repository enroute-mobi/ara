package api

import (
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/siri"
)

type SIRIRequestHandler interface {
	RequestorRef() string
	ConnectorType() string
	XMLResponse(core.Connector) string
}

type SIRIHandler struct {
	referential *core.Referential
}

func NewSIRIHandler(referential *core.Referential) *SIRIHandler {
	return &SIRIHandler{referential: referential}
}

func (handler *SIRIHandler) requestHandler(envelope *siri.SOAPEnvelope) SIRIRequestHandler {
	switch envelope.BodyType() {
	case "CheckStatus":
		return &SIRICheckStatusRequestHandler{
			xmlRequest: siri.NewXMLCheckStatusRequest(envelope.Body()),
		}
	case "StopMonitoring":
		return &SIRIStopMonitoringRequestHandler{}
	}
	return nil
}

func (handler *SIRIHandler) siriError(err string, response http.ResponseWriter) {

}

func (handler *SIRIHandler) serve(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/xml")

	if handler.referential == nil {
		// http.Error(response, "Referential not found", 500)
		handler.siriError("...", response)
		return
	}
	if request.Body == nil {
		// http.Error(response, "Invalid request: Empty body", 400)
		handler.siriError("...", response)
		return
	}
	envelope, err := siri.NewSOAPEnvelope(request.Body)
	if err != nil {
		// http.Error(response, "Invalid request: can't read content", 400)
		handler.siriError("...", response)
		return
	}

	// if envelope.BodyType() == "GetSiriService" {
	// 	// TODO
	// }

	requestHandler := handler.requestHandler(envelope)
	if requestHandler == nil {
		// http.Error(response, fmt.Sprintf("Cannot handle SIRI request %v", envelope.BodyType()), 500)
		handler.siriError("...", response)
		return
	}

	partner, ok := handler.referential.Partners().FindByLocalCredential(requestHandler.RequestorRef())
	if !ok {
		handler.siriError("...", response)
	}
	connector, ok := partner.Connector(requestHandler.ConnectorType())
	if !ok {
		handler.siriError("...", response)
	}

	xmlResponse := requestHandler.XMLResponse(connector)

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	_, err = soapEnvelope.WriteTo(response)
	if err != nil {
		handler.siriError("...", response)
	}
}
