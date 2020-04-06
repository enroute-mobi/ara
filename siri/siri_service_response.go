package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/edwig/logger"
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

func (response *SIRIServiceResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "siri_service_response.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
