package siri

import (
	"bytes"
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type SIRINotifyStopMonitoring struct {
	Address                   string
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string
	ResponseTimestamp         time.Time

	Deliveries []*SIRINotifyStopMonitoringDelivery
}

type SIRINotifyStopMonitoringDelivery struct {
	MonitoringRef          string
	RequestMessageRef      string
	SubscriberRef          string
	SubscriptionIdentifier string
	ResponseTimestamp      time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	MonitoredVehicleJourney *SIRIMonitoredVehicleJourney
	MonitoredStopVisits     []*SIRIMonitoredStopVisit
	CancelledStopVisits     []*SIRICancelledStopVisit
}

func (notify *SIRINotifyStopMonitoring) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("stop_monitoring_notify%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRINotifyStopMonitoringDelivery) ErrorString() string {
	return fmt.Sprintf("%v: %v", delivery.errorType(), delivery.ErrorText)
}

func (delivery *SIRINotifyStopMonitoringDelivery) errorType() string {
	if delivery.ErrorType == siri_attributes.OtherError {
		return fmt.Sprintf("%v %v", delivery.ErrorType, delivery.ErrorNumber)
	}
	return delivery.ErrorType
}

func (delivery *SIRINotifyStopMonitoringDelivery) BuildNotifyStopMonitoringDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "notify_stop_monitoring_delivery.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
