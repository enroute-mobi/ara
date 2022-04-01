package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
)

type SIRIServiceResponse struct {
	ProducerRef               string
	ResponseMessageIdentifier string
	RequestMessageRef         string
	Status                    bool
	// ErrorType                 string
	// ErrorNumber               int
	// ErrorText                 string

	ResponseTimestamp time.Time

	StopMonitoringDeliveries     []*SIRIStopMonitoringDelivery
	GeneralMessageDeliveries     []*SIRIGeneralMessageDelivery
	EstimatedTimetableDeliveries []*SIRIEstimatedTimetableDelivery
}

func (response *SIRIServiceResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("siri_service_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
