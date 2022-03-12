package api

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/remote"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIGeneralMessageRequestHandler struct {
	xmlRequest  *siri.XMLGetGeneralMessage
	referential *core.Referential
}

func (handler *SIRIGeneralMessageRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRIGeneralMessageRequestHandler) ConnectorType() string {
	return core.SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER
}

func (handler *SIRIGeneralMessageRequestHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("General Message %s\n", handler.xmlRequest.MessageIdentifier())

	t := clock.DefaultClock().Now()

	tmp := params.connector.(*core.SIRIGeneralMessageRequestBroadcaster)
	response, _ := tmp.Situations(handler.xmlRequest, params.message)
	xmlResponse, err := response.BuildXML()
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

	params.message.Type = "GeneralMessageRequest"
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ResponseRawMessage = xmlResponse
	params.message.ResponseSize = n
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
