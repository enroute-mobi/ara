package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SiriErrorResponse struct {
	response       http.ResponseWriter
	errCode        string
	errDescription string
	request        string
}

func siriError(errCode, errDescription, referentialSlug string, response http.ResponseWriter) {
	siriErrorWithRequest(errCode, errDescription, referentialSlug, "", response)
}

func siriErrorWithRequest(errCode, errDescription, referentialSlug, request string, response http.ResponseWriter) {
	SiriErrorResponse{
		response:       response,
		errCode:        errCode,
		errDescription: errDescription,
		request:        request,
	}.sendSiriError(referentialSlug)
}

func (siriError SiriErrorResponse) sendSiriError(referentialSlug string) {
	logger.Log.Debugf("Send SIRI error %v : %v", siriError.errCode, siriError.errDescription)

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(fmt.Sprintf(`
  <S:Fault>
    <faultcode>S:%s</faultcode>
    <faultstring>%s</faultstring>
  </S:Fault>`, siriError.errCode, siriError.errDescription))

	logStashEvent := make(audit.LogStashEvent)
	logStashEvent["status"] = "false"
	logStashEvent["siriType"] = "siriError"
	logStashEvent["responseXML"] = soapEnvelope.String()

	message := &audit.BigQueryMessage{
		Protocol:           "siri",
		Direction:          "send",
		Status:             "Error",
		Type:               "siri-error",
		ErrorDetails:       fmt.Sprintf("%v: %v", siriError.errCode, siriError.errDescription),
		ResponseRawMessage: soapEnvelope.String(),
	}

	if siriError.request != "" {
		logStashEvent["requestXML"] = siriError.request
		message.RequestRawMessage = siriError.request
	}

	soapEnvelope.WriteTo(siriError.response)
	message.ResponseSize = soapEnvelope.Length()

	audit.CurrentLogStash().WriteEvent(logStashEvent)
	audit.CurrentBigQuery(referentialSlug).WriteEvent(message)
}
