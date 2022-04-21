package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/api/rah"
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIRequestHandler interface {
	RequestorRef() string
	ConnectorType() string
	Respond(HandlerParams)
}

type SIRIHandler struct {
	rah.RemoteAddressHandler

	referential *core.Referential
}

type HandlerParams struct {
	connector    core.Connector
	rw           http.ResponseWriter
	message      *audit.BigQueryMessage
	envelopeType string
}

func NewSIRIHandler(referential *core.Referential) *SIRIHandler {
	return &SIRIHandler{referential: referential}
}

func (handler *SIRIHandler) requestHandler(envelope *remote.SIRIEnvelope) SIRIRequestHandler {
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
	case "DeleteSubscription", "TerminateSubscription":
		return &SIRIDeleteSubscriptionRequestHandler{
			xmlRequest:  siri.NewXMLDeleteSubscriptionRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "Subscribe", "Subscription":
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
		SIRIError{
			errCode:        "NotFound",
			errDescription: "Referential not found",
			response:       response,
		}.Send()
		return
	}

	// Check if request is gzip
	var requestReader io.Reader
	if request.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(request.Body)
		if err != nil {
			SIRIError{
				errCode:         "Client",
				errDescription:  "Can't unzip request",
				referentialSlug: string(handler.referential.Slug()),
				response:        response,
			}.Send()
			return
		}
		defer gzipReader.Close()
		requestReader = gzipReader
	} else {
		requestReader = request.Body
	}

	envelope, err := remote.NewAutodetectSIRIEnvelope(requestReader)
	if err != nil {
		SIRIError{
			errCode:         "Client",
			errDescription:  fmt.Sprintf("Invalid Request: %v", err),
			referentialSlug: string(handler.referential.Slug()),
			response:        response,
		}.Send()
		return
	}

	requestHandler := handler.requestHandler(envelope)
	if requestHandler == nil {
		SIRIError{
			errCode:         "NotSupported",
			errDescription:  fmt.Sprintf("SIRIRequest %v not supported", envelope.BodyType()),
			referentialSlug: string(handler.referential.Slug()),
			request:         envelope.Body().String(),
			response:        response,
		}.Send()
		return
	}

	if requestHandler.RequestorRef() == "" {
		SIRIError{
			errCode:         "UnknownCredential",
			errDescription:  "Can't have empty credentials",
			referentialSlug: string(handler.referential.Slug()),
			request:         envelope.Body().String(),
			response:        response,
		}.Send()
		return
	}
	partner, ok := handler.referential.Partners().FindByCredential(requestHandler.RequestorRef())
	if !ok {
		SIRIError{
			errCode:         "UnknownCredential",
			errDescription:  fmt.Sprintf("RequestorRef Unknown '%s'", requestHandler.RequestorRef()),
			referentialSlug: string(handler.referential.Slug()),
			request:         envelope.Body().String(),
			response:        response,
		}.Send()
		return
	}
	connector, ok := partner.Connector(requestHandler.ConnectorType())
	if !ok {
		SIRIError{
			errCode:         "NotFound",
			errDescription:  fmt.Sprintf("No Connectors for %v", envelope.BodyType()),
			referentialSlug: string(handler.referential.Slug()),
			request:         envelope.Body().String(),
			response:        response,
		}.Send()
		return
	}

	m := &audit.BigQueryMessage{
		Protocol:    "siri",
		Direction:   "received",
		Partner:     string(partner.Slug()),
		IPAddress:   handler.HandleRemoteAddress(request),
		RequestSize: request.ContentLength,
		Status:      "OK",
	}

	params := HandlerParams{
		connector:    connector,
		rw:           response,
		message:      m,
		envelopeType: partner.SIRIEnvelopeType(),
	}

	requestHandler.Respond(params)
}
