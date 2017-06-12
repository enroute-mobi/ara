package api

import (
	"fmt"
	"net/http"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/siri"
)

type SiriError struct {
	response       http.ResponseWriter
	errCode        string
	errDescription string
	request        string
}

func siriError(errCode, errDescription string, response http.ResponseWriter) {
	siriErrorWithRequest(errCode, errDescription, "", response)
}

func siriErrorWithRequest(errCode, errDescription, request string, response http.ResponseWriter) {
	SiriError{
		response:       response,
		errCode:        errCode,
		errDescription: errDescription,
		request:        request,
	}.sendSiriError()
}

func (siriError SiriError) sendSiriError() {
	// Wrap soap and send response
	soapEnvelope := siri.NewSOAPEnvelopeBuffer()
	soapEnvelope.WriteXML(fmt.Sprintf(`
  <S:Fault xmlns:ns4="http://www.w3.org/2003/05/soap-envelope">
    <faultcode>S:%s</faultcode>
    <faultstring>%s</faultstring>
  </S:Fault>`, siriError.errCode, siriError.errDescription))

	logStashEvent := make(audit.LogStashEvent)
	logStashEvent["SIRIError"] = soapEnvelope.String()
	if siriError.request != "" {
		logStashEvent["Request"] = siriError.request
	}
	audit.CurrentLogStash().WriteEvent(logStashEvent)

	soapEnvelope.WriteTo(siriError.response)
}
