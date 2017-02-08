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
	Respond(core.SIRIConnector, http.ResponseWriter)
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

func siriError(errCode, errDescription string, response http.ResponseWriter) {
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
		siriError("NotFound", "Referential not found", response)
		return
	}
	if request.Body == nil {
		siriError("InvalidRequest", "Request Body is empty", response)
		return
	}
	envelope, err := siri.NewSOAPEnvelope(request.Body)
	if err != nil {
		siriError("InvalidRequest", "Cannot read Request Body", response)
		return
	}

	// if envelope.BodyType() == "GetSiriService" {
	// 	// TODO
	// }

	requestHandler := handler.requestHandler(envelope)
	if requestHandler == nil {
		siriError("NotSupported", fmt.Sprintf("SIRIRequest %v not supported", envelope.BodyType()), response)
		return
	}

	partner, ok := handler.referential.Partners().FindByLocalCredential(requestHandler.RequestorRef())
	if !ok {
		siriError("UnknownCredential", "RequestorRef Unknown", response)
		return
	}
	connector, ok := partner.SIRIConnector(requestHandler.ConnectorType())
	if !ok {
		siriError("NotFound", fmt.Sprintf("No Connectors for %v", envelope.BodyType()), response)
		return
	}

	requestHandler.Respond(connector, response)
}
