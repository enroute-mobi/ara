package api

import (
	"encoding/json"
	"net/http"
	"net/url"

	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRILiteVehicleMonitoringRequestHandler struct {
	requestUrl string
	filters    url.Values
}

func (handler *SIRILiteVehicleMonitoringRequestHandler) ConnectorType() string {
	return core.SIRI_LITE_VEHICLE_MONITORING_REQUEST_BROADCASTER
}

func (handler *SIRILiteVehicleMonitoringRequestHandler) Respond(connector core.Connector, rw http.ResponseWriter) {
	logger.Log.Debugf("Siri Lite VehicleMonitoring %s", handler.requestUrl)

	response := connector.(core.VehicleMonitoringRequestBroadcaster).RequestVehicles(handler.requestUrl, handler.filters)

	jsonBytes, err := json.Marshal(response)
	if err != nil {
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
		logger.Log.Debugf("Internal error while Marshaling a SiriLite response in vehicle monitoring handler: %v", err)
		return
	}
	_, err = rw.Write(jsonBytes)
	if err != nil {
		logger.Log.Debugf("Internal error while writing a SiriLite response in vehicle monitoring handler: %v", err)
		http.Error(rw, "Internal Server Error", http.StatusInternalServerError)
	}

}
