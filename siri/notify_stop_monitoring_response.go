package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/edwig/logger"
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

	MonitoredStopVisits []*SIRIMonitoredStopVisit
	CancelledStopVisits []*SIRICancelledStopVisit
}

func (notify *SIRINotifyStopMonitoring) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "stop_monitoring_notify.template", notify); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRINotifyStopMonitoringDelivery) BuildNotifyStopMonitoringDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "notify_stop_monitoring_delivery.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}
