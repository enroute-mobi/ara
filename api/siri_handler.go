package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIRequestHandler interface {
	RequestorRef() string
	ConnectorType() string
	Respond(core.Connector, http.ResponseWriter, *audit.BigQueryMessage)
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
			xmlRequest:  siri.NewXMLCheckStatusRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "GetStopMonitoring":
		return &SIRIStopMonitoringRequestHandler{
			xmlRequest:  siri.NewXMLGetStopMonitoring(envelope.Body()),
			referential: handler.referential,
		}
	case "DeleteSubscription":
		return &SIRIDeleteSubscriptionRequestHandler{
			xmlRequest:  siri.NewXMLDeleteSubscriptionRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "Subscribe":
		return &SIRISubscribeRequestHandler{
			xmlRequest:  siri.NewXMLSubscriptionRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifyStopMonitoring":
		return &SIRIStopMonitoringRequestDeliveriesResponseHandler{
			xmlRequest:  siri.NewXMLNotifyStopMonitoring(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifyGeneralMessage":
		return &SIRIGeneralMessageRequestDeliveriesResponseHandler{
			xmlRequest:  siri.NewXMLNotifyGeneralMessage(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifySubscriptionTerminated":
		return &SIRINotifySubscriptionTerminatedHandler{
			xmlRequest:  siri.NewXMLNotifySubscriptionTerminated(envelope.Body()),
			referential: handler.referential,
		}
	case "SubscriptionTerminatedNotification":
		return &SIRISubscriptionTerminatedNotificationHandler{
			xmlRequest:  siri.NewXMLSubscriptionTerminatedNotification(envelope.Body()),
			referential: handler.referential,
		}
	case "StopPointsDiscovery":
		return &SIRIStopDiscoveryRequestHandler{
			xmlRequest:  siri.NewXMLStopPointsDiscoveryRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "LinesDiscovery":
		return &SIRILinesDiscoveryRequestHandler{
			xmlRequest:  siri.NewXMLLinesDiscoveryRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "GetSiriService":
		return &SIRIServiceRequestHandler{
			xmlRequest:  siri.NewXMLSiriServiceRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "GetGeneralMessage":
		return &SIRIGeneralMessageRequestHandler{
			xmlRequest:  siri.NewXMLGetGeneralMessage(envelope.Body()),
			referential: handler.referential,
		}
	case "GetEstimatedTimetable":
		return &SIRIEstimatedTimetableRequestHandler{
			xmlRequest:  siri.NewXMLGetEstimatedTimetable(envelope.Body()),
			referential: handler.referential,
		}
	}
	return nil
}

func (handler *SIRIHandler) serve(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/xml; charset=utf-8")

	if handler.referential == nil {
		siriError("NotFound", "Referential not found", "", response)
		return
	}

	// Check if request is gzip
	var requestReader io.Reader
	if request.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(request.Body)
		if err != nil {
			siriError("Client", "Can't unzip request", string(handler.referential.Slug()), response)
			return
		}
		defer gzipReader.Close()
		requestReader = gzipReader
	} else {
		requestReader = request.Body
	}

	envelope, err := siri.NewSOAPEnvelope(requestReader)
	if err != nil {
		siriError("Client", fmt.Sprintf("Invalid Request: %v", err), string(handler.referential.Slug()), response)
		return
	}

	requestHandler := handler.requestHandler(envelope)
	if requestHandler == nil {
		siriErrorWithRequest("NotSupported", fmt.Sprintf("SIRIRequest %v not supported", envelope.BodyType()), string(handler.referential.Slug()), envelope.Body().String(), response)
		return
	}

	if requestHandler.RequestorRef() == "" {
		siriErrorWithRequest("UnknownCredential", "Can't have empty credentials", string(handler.referential.Slug()), envelope.Body().String(), response)
		return
	}
	partner, ok := handler.referential.Partners().FindByCredential(requestHandler.RequestorRef())
	if !ok {
		siriErrorWithRequest("UnknownCredential", "RequestorRef Unknown", string(handler.referential.Slug()), envelope.Body().String(), response)
		return
	}
	connector, ok := partner.Connector(requestHandler.ConnectorType())
	if !ok {
		siriErrorWithRequest("NotFound", fmt.Sprintf("No Connectors for %v", envelope.BodyType()), string(handler.referential.Slug()), envelope.Body().String(), response)
		return
	}

	m := &audit.BigQueryMessage{
		Protocol:    "siri",
		Direction:   "received",
		Partner:     string(partner.Slug()),
		IPAddress:   request.RemoteAddr,
		RequestSize: request.ContentLength,
		Status:      "OK",
	}

	requestHandler.Respond(connector, response, m)
}
