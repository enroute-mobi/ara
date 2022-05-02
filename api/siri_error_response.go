package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/remote"
)

type SIRIError struct {
	errCode         string
	errDescription  string
	referentialSlug string
	request         string
	envelopeType    string
	response        http.ResponseWriter
}

func (e SIRIError) Send() {
	logger.Log.Debugf("Send SIRI error %v : %v", e.errCode, e.errDescription)

	// Wrap soap and send response
	soapEnvelope := remote.NewSIRIBuffer(e.envelopeType)
	soapEnvelope.WriteXML(fmt.Sprintf(`
  <S:Fault>
    <faultcode>S:%s</faultcode>
    <faultstring>%s</faultstring>
  </S:Fault>`, e.errCode, e.errDescription))

	message := &audit.BigQueryMessage{
		Protocol:  "siri",
		Direction: "received",
		Status:    "Error",
		// Type:         "siri-error",
		ErrorDetails: fmt.Sprintf("%v: %v", e.errCode, e.errDescription),
		// ResponseRawMessage: soapEnvelope.String(),
	}

	// if e.request != "" {
	// 	message.RequestRawMessage = e.request
	// }

	soapEnvelope.WriteTo(e.response)
	message.ResponseSize = soapEnvelope.Length()

	audit.CurrentBigQuery(e.referentialSlug).WriteEvent(message)
}
