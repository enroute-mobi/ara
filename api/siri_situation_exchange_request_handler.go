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

type SIRISituationExchangeRequestHandler struct {
	xmlRequest  *sxml.XMLGetSituationExchange
	referential *core.Referential
}

func (handler *SIRISituationExchangeRequestHandler) RequestorRef() string {
	return handler.xmlRequest.RequestorRef()
}

func (handler *SIRISituationExchangeRequestHandler) ConnectorType() string {
	return core.SIRI_SITUATION_EXCHANGE_REQUEST_BROADCASTER
}

func (handler *SIRISituationExchangeRequestHandler) Respond(params HandlerParams) {
	logger.Log.Debugf("Situation Exchange %s\n", handler.xmlRequest.MessageIdentifier())

	t := clock.DefaultClock().Now()

	response := params.connector.(*core.SIRISituationExchangeRequestBroadcaster).Situations(handler.xmlRequest, params.message)
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

	params.message.Type = audit.SITUATION_EXCHANGE_REQUEST
	params.message.RequestRawMessage = handler.xmlRequest.RawXML()
	params.message.ResponseRawMessage = xmlResponse
	params.message.ResponseSize = n
	params.message.ProcessingTime = clock.DefaultClock().Since(t).Seconds()
	audit.CurrentBigQuery(string(handler.referential.Slug())).WriteEvent(params.message)
}
