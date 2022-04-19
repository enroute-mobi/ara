package api

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRISubscribeRequestHandler struct {
	xmlRequest  *siri.XMLSubscriptionRequest
	referential *core.Referential
}

func (handler *SIRISubscribeRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRISubscribeRequestHandler) ConnectorType() string {
	return core.SIRI_SUBSCRIPTION_REQUEST_DISPATCHER
}

func (handler *SIRISubscribeRequestHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("SubscribeRequest %s\n", handler.xmlRequest.MessageIdentifier())
	t := handler.referential.Clock().Now()

	response, err := params.connector.(core.SubscriptionRequestDispatcher).Dispatch(handler.xmlRequest, params.message)
	if err != nil {
		SIRIError{
			errCode:         "NotFound",
			errDescription:  err.Error(),
			referentialSlug: string(handler.referential.Slug()),
			envelopeType:    params.envelopeType,
			request:         handler.xmlRequest.RawXML(),
			response:        params.rw,
		}.Send()
		return
	}

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

	params.message.ProcessingTime = handler.referential.Clock().Since(t).Seconds()
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ResponseRawMessage = xmlResponse
	params.message.ResponseSize = n
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
