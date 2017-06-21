package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"github.com/af83/edwig/core"
	"github.com/af83/edwig/siri"
	"github.com/af83/edwig/version"
)

type SIRIRequestHandler interface {
	RequestorRef() string
	ConnectorType() string
	Respond(core.Connector, http.ResponseWriter)
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
	case "GetStopMonitoring":
		return &SIRIStopMonitoringRequestHandler{
			xmlRequest: siri.NewXMLStopMonitoringRequest(envelope.Body()),
		}
	case "NotifyStopMonitoring":
		return &SIRIStopMonitoringDeliveriesResponseHandler{
			xmlRequest: siri.NewXMLNotifyStopMonitoring(envelope.Body()),
		}
	case "SubscriptionTerminatedNotification":
		return &SIRIStopMonitoringSubscriptionTerminatedResponseHandler{
			xmlRequest: siri.NewXMLStopMonitoringSubscriptionTerminatedResponse(envelope.Body()),
		}
	case "StopPointsDiscovery":
		return &SIRIStopDiscoveryRequestHandler{
			xmlRequest: siri.NewXMLStopPointsDiscoveryRequest(envelope.Body()),
		}
	case "GetSiriService":
		return &SIRIServiceRequestHandler{
			xmlRequest: siri.NewXMLSiriServiceRequest(envelope.Body()),
		}
	case "GetGeneralMessage":
		return &SIRIGeneralMessageRequestHandler{
			xmlRequest: siri.NewXMLGeneralMessageRequest(envelope.Body()),
		}
	}
	return nil
}

func (handler *SIRIHandler) serve(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/xml; charset=utf-8")
	response.Header().Set("Server", version.ApplicationName())

	if handler.referential == nil {
		siriError("NotFound", "Referential not found", response)
		return
	}

	// Check if request is gzip
	var requestReader io.Reader
	if request.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(request.Body)
		if err != nil {
			siriError("Client", "Can't unzip request", response)
			return
		}
		defer gzipReader.Close()
		requestReader = gzipReader
	} else {
		requestReader = request.Body
	}

	envelope, err := siri.NewSOAPEnvelope(requestReader)
	if err != nil {
		siriError("Client", fmt.Sprintf("Invalid Request: %v", err), response)
		return
	}

	requestHandler := handler.requestHandler(envelope)
	if requestHandler == nil {
		siriErrorWithRequest("NotSupported", fmt.Sprintf("SIRIRequest %v not supported", envelope.BodyType()), envelope.Body().String(), response)
		return
	}

	partner, ok := handler.referential.Partners().FindByLocalCredential(requestHandler.RequestorRef())
	if !ok {
		siriErrorWithRequest("UnknownCredential", "RequestorRef Unknown", envelope.Body().String(), response)
		return
	}
	connector, ok := partner.Connector(requestHandler.ConnectorType())
	if !ok {
		siriErrorWithRequest("NotFound", fmt.Sprintf("No Connectors for %v", envelope.BodyType()), envelope.Body().String(), response)
		return
	}

	requestHandler.Respond(connector, response)
}
