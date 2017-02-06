package api

import (
	"fmt"
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
			referential: handler.referential,
			xmlRequest:  siri.NewXMLCheckStatusRequest(envelope.Body()),
		}
	case "StopMonitoring":
		return &SIRIStopMonitoringRequestHandler{
			referential: handler.referential,
			xmlRequest:  siri.NewXMLStopMonitoringRequest(envelope.Body()),
		}
	}
	return nil
}

func (handler *SIRIHandler) siriError(errCode, errDescription string, response http.ResponseWriter) {
	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(fmt.Sprintf(`<S:Body>
  <S:Fault xmlns:ns4="http://www.w3.org/2003/05/soap-envelope">
    <faultcode>S:%s</faultcode>
    <faultstring>%s</faultstring>
  </S:Fault>
</S:Body>`, errCode, errDescription))

	soapEnvelope.WriteTo(response)
}

func (handler *SIRIHandler) serve(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/xml")

	if handler.referential == nil {
		handler.siriError("NotFound", "Referential not found", response)
		return
	}
	if request.Body == nil {
		handler.siriError("InvalidRequest", "Request Body is empty", response)
		return
	}
	envelope, err := siri.NewSOAPEnvelope(request.Body)
	if err != nil {
		handler.siriError("InvalidRequest", "Cannot read Request Body", response)
		return
	}

	// if envelope.BodyType() == "GetSiriService" {
	// 	// TODO
	// }

	requestHandler := handler.requestHandler(envelope)
	if requestHandler == nil {
		handler.siriError("NotSupported", fmt.Sprintf("SIRIRequest %v not supported", envelope.BodyType()), response)
		return
	}

	partner, ok := handler.referential.Partners().FindByLocalCredential(requestHandler.RequestorRef())
	if !ok {
		handler.siriError("UnknownCredential", "RequestorRef Unknown", response)
		return
	}
	connector, ok := partner.Connector(requestHandler.ConnectorType())
	if !ok {
		handler.siriError("NotFound", fmt.Sprintf("No Connectors for %v", envelope.BodyType()), response)
		return
	}

	xmlResponse := requestHandler.XMLResponse(connector)

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(xmlResponse)

	_, err = soapEnvelope.WriteTo(response)
	if err != nil {
		handler.siriError("InternalServiceError", fmt.Sprintf("Internal Error: %v", err), response)
	}
}
