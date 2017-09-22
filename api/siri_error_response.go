package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SiriErrorResponse struct {
	response       http.ResponseWriter
	errCode        string
	errDescription string
	request        string
}

func siriError(errCode, errDescription string, response http.ResponseWriter) {
	siriErrorWithRequest(errCode, errDescription, "", response)
}

func siriErrorWithRequest(errCode, errDescription, request string, response http.ResponseWriter) {
	SiriErrorResponse{
		response:       response,
		errCode:        errCode,
		errDescription: errDescription,
		request:        request,
	}.sendSiriError()
}

func (siriError SiriErrorResponse) sendSiriError() {
	logger.Log.Debugf("Send SIRI error %v : %v", siriError.errCode, siriError.errDescription)

	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(fmt.Sprintf(`
  <S:Fault>
    <faultcode>S:%s</faultcode>
    <faultstring>%s</faultstring>
  </S:Fault>`, siriError.errCode, siriError.errDescription))

	logStashEvent := make(audit.LogStashEvent)
	logStashEvent["siriError"] = soapEnvelope.String()
	if siriError.request != "" {
		logStashEvent["request"] = siriError.request
	}
	audit.CurrentLogStash().WriteEvent(logStashEvent)

	soapEnvelope.WriteTo(siriError.response)
}
