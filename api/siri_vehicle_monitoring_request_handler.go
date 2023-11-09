package api

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIVehicleMonitoringRequestHandler struct {
	xmlRequest  *sxml.XMLGetVehicleMonitoring
	referential *core.Referential
}

func (handler *SIRIVehicleMonitoringRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIVehicleMonitoringRequestHandler) ConnectorType() string {
	return core.SIRI_VEHICLE_MONITORING_REQUEST_BROADCASTER
}

func (handler *SIRIVehicleMonitoringRequestHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("Siri VehicleMonitoring %s", handler.xmlRequest.MessageIdentifier())

	t := clock.DefaultClock().Now()

	response := params.connector.(core.VehicleMonitoringRequestBroadcaster).RequestVehicles(handler.xmlRequest, params.message)
	xmlResponse, err := response.BuildXML(params.envelopeType)
	if err != nil {
		SIRIError{
			errCode:         "InternalServiceError",
			errDescription:  fmt.Sprintf("Internal Error: %v", err),
			referentialSlug: string(handler.referential.Slug()),
			envelopeType:    params.envelopeType,
			response:        params.rw,
		}.Send()
		return
	}

	// Wrap soap and send response
	buffer := remote.NewSIRIBuffer(params.envelopeType)
	buffer.WriteXML(xmlResponse)

	n, err := buffer.WriteTo(params.rw)
	if err != nil {
		SIRIError{
			errCode:         "InternalServiceError",
			errDescription:  fmt.Sprintf("Internal Error: %v", err),
			referentialSlug: string(handler.referential.Slug()),
			envelopeType:    params.envelopeType,
			response:        params.rw,
		}.Send()
		return
	}

	params.message.Type = audit.VEHICLE_MONITORING_REQUEST
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ResponseRawMessage = xmlResponse
	params.message.ResponseSize = n
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
