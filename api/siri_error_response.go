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

	buffer := remote.NewSIRIBuffer(e.envelopeType)
	buffer.WriteXML(fmt.Sprintf(`
  <S:Fault>
    <faultcode>S:%s</faultcode>
    <faultstring>%s</faultstring>
  </S:Fault>`, e.errCode, e.errDescription))

	message := &audit.BigQueryMessage{
		Protocol:     "siri",
		Direction:    "received",
		Status:       "Error",
		ErrorDetails: fmt.Sprintf("%v: %v", e.errCode, e.errDescription),
	}

	buffer.WriteTo(e.response)
	message.ResponseSize = buffer.Length()

	audit.CurrentBigQuery(e.referentialSlug).WriteEvent(message)
}
