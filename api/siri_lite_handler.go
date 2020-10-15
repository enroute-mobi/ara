package api

import (
	"net/http"

	"bitbucket.org/enroute-mobi/edwig/core"
)

type SIRILiteRequestHandler interface {
	ConnectorType() string
	Respond(core.Connector, http.ResponseWriter)
}

type SIRILiteHandler struct {
	referential *core.Referential
	token       string
}

func NewSIRILiteHandler(referential *core.Referential, token string) *SIRILiteHandler {
	return &SIRILiteHandler{
		referential: referential,
		token:       token,
	}
}

func (handler *SIRILiteHandler) requestHandler(requestData *SIRIRequestData) SIRILiteRequestHandler {
	switch requestData.Request {
	case "vehicle-monitoring":
		return &SIRILiteVehicleMonitoringRequestHandler{
			requestUrl: requestData.Url,
			filters:    requestData.Filters,
		}
	}
	// case "CheckStatus":
	// case "GetStopMonitoring":
	// case "DeleteSubscription":
	// case "Subscribe":
	// case "NotifyStopMonitoring":
	// case "NotifyGeneralMessage":
	// case "SubscriptionTerminatedNotification":
	// case "StopPointsDiscovery":
	// case "LinesDiscovery":
	// case "GetSiriService":
	// case "GetGeneralMessage":
	// case "GetEstimatedTimetable":
	// }
	return nil
}

func (handler *SIRILiteHandler) serve(response http.ResponseWriter, request *http.Request, requestData *SIRIRequestData) {
	if handler.token == "" {
		http.Error(response, "No Authorization Token", http.StatusUnauthorized)
		return
	}

	if handler.referential == nil {
		http.Error(response, "Referential not found", http.StatusNotFound)
		return
	}

	// Find Partner by authorization Key
	partner, ok := handler.referential.Partners().FindByCredential(handler.token)
	if !ok {
		http.Error(response, "Invalid Authorization Token", http.StatusForbidden)
		return
	}

	requestHandler := handler.requestHandler(requestData)
	if requestHandler == nil {
		http.Error(response, "The SIRI Lite request doesn’t match a defined broadcast", http.StatusNotFound)
		return
	}

	connector, ok := partner.Connector(requestHandler.ConnectorType())
	if !ok {
		http.Error(response, "The Partner don't support this SIRI Lite request", http.StatusNotFound)
		return
	}

	response.Header().Set("Content-Type", "application/json")

	requestHandler.Respond(connector, response)
}
