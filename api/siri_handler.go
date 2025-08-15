package api

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"bitbucket.org/enroute-mobi/ara/api/rah"
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
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
			xmlRequest:  sxml.NewXMLCheckStatusRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "GetStopMonitoring":
		return &SIRIStopMonitoringRequestHandler{
			xmlRequest:  sxml.NewXMLGetStopMonitoring(envelope.Body()),
			referential: handler.referential,
		}
	case "DeleteSubscription", "TerminateSubscription":
		return &SIRIDeleteSubscriptionRequestHandler{
			xmlRequest:  sxml.NewXMLDeleteSubscriptionRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "Subscribe", "Subscription":
		return &SIRISubscribeRequestHandler{
			xmlRequest:  sxml.NewXMLSubscriptionRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifyStopMonitoring":
		return &SIRIStopMonitoringRequestDeliveriesResponseHandler{
			xmlRequest:  sxml.NewXMLNotifyStopMonitoring(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifyGeneralMessage":
		return &SIRIGeneralMessageRequestDeliveriesResponseHandler{
			xmlRequest:  sxml.NewXMLNotifyGeneralMessage(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifySituationExchange":
		return &SIRISituationExchangeDeliveriesResponseHandler{
			xmlRequest:  sxml.NewXMLNotifySituationExchange(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifyEstimatedTimetable", "EstimatedTimetableDelivery":
		return &SIRIEstimatedTimetableRequestDeliveriesResponseHandler{
			xmlRequest:  sxml.NewXMLNotifyEstimatedTimetable(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifyVehicleMonitoring", "VehicleMonitoringDelivery":
		return &SIRIVehicleMonitoringRequestDeliveriesResponseHandler{
			xmlRequest:  sxml.NewXMLNotifyVehicleMonitoring(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifyFacilityMonitoring", "FacilityMonitoringDelivery":
		return &SIRIFacilityMonitoringRequestDeliveriesResponseHandler{
			xmlRequest:  sxml.NewXMLNotifyFacilityMonitoring(envelope.Body()),
			referential: handler.referential,
		}
	case "NotifySubscriptionTerminated":
		return &SIRINotifySubscriptionTerminatedHandler{
			xmlRequest:  sxml.NewXMLNotifySubscriptionTerminated(envelope.Body()),
			referential: handler.referential,
		}
	case "SubscriptionTerminatedNotification":
		return &SIRISubscriptionTerminatedNotificationHandler{
			xmlRequest:  sxml.NewXMLSubscriptionTerminatedNotification(envelope.Body()),
			referential: handler.referential,
		}
	case "StopPointsDiscovery", "StopPoints":
		return &SIRIStopDiscoveryRequestHandler{
			xmlRequest:  sxml.NewXMLStopPointsDiscoveryRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "LinesDiscovery", "Lines":
		return &SIRILinesDiscoveryRequestHandler{
			xmlRequest:  sxml.NewXMLLinesDiscoveryRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "GetSiriService":
		return &SIRIServiceRequestHandler{
			xmlRequest:  sxml.NewXMLSiriServiceRequest(envelope.Body()),
			referential: handler.referential,
		}
	case "GetGeneralMessage":
		return &SIRIGeneralMessageRequestHandler{
			xmlRequest:  sxml.NewXMLGetGeneralMessage(envelope.Body()),
			referential: handler.referential,
		}
	case "GetEstimatedTimetable":
		return &SIRIEstimatedTimetableRequestHandler{
			xmlRequest:  sxml.NewXMLGetEstimatedTimetable(envelope.Body()),
			referential: handler.referential,
		}
	case "GetVehicleMonitoring":
		return &SIRIVehicleMonitoringRequestHandler{
			xmlRequest:  sxml.NewXMLGetVehicleMonitoring(envelope.Body()),
			referential: handler.referential,
		}
	case "GetSituationExchange":
		return &SIRISituationExchangeRequestHandler{
			xmlRequest:  sxml.NewXMLGetSituationExchange(envelope.Body()),
			referential: handler.referential,
		}
	case "GetFacilityMonitoring", "FacilityMonitoring":
		return &SIRIFacilityMonitoringRequestHandler{
			xmlRequest:  sxml.NewXMLGetFacilityMonitoring(envelope.Body()),
			referential: handler.referential,
		}
	}

	return nil
}

func (handler *SIRIHandler) serve(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "text/xml; charset=utf-8")

	m := &audit.BigQueryMessage{
		Protocol:    "siri",
		Direction:   "received",
		IPAddress:   handler.HandleRemoteAddress(request),
		RequestSize: request.ContentLength,
	}

	if handler.referential == nil {
		SIRIError{
			errCode:        "NotFound",
			errDescription: "Referential not found",
			request:        dumpRequest(request),
			response:       response,
			message:        m,
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
				request:         dumpRequest(request),
				response:        response,
				message:         m,
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
			request:         dumpRequest(request),
			response:        response,
			message:         m,
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
			message:         m,
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
			message:         m,
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
			message:         m,
		}.Send()
		return
	}

	m.Partner = string(partner.Slug())
	if !partner.Allow(handler.HandleRemoteAddress(request)) {
		http.Error(response, "Too many requests", http.StatusTooManyRequests)
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
			message:         m,
		}.Send()
		return
	}

	m.Status = "OK"

	params := HandlerParams{
		connector:    connector,
		rw:           response,
		message:      m,
		envelopeType: partner.SIRIEnvelopeType(),
	}

	requestHandler.Respond(params)
}

func dumpRequest(r *http.Request) string {
	d, err := httputil.DumpRequest(r, true)
	if err != nil {
		return ""
	}
	return string(d)
}
