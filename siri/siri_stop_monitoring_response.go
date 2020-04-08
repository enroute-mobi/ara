package siri

import (
	"bytes"
	"time"

	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/model"
)

type SIRIStopMonitoringResponse struct {
	SIRIStopMonitoringDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
}

type SIRIStopMonitoringDelivery struct {
	RequestMessageRef string
	Status            bool
	ErrorType         string
	ErrorNumber       int
	ErrorText         string
	ResponseTimestamp time.Time
	MonitoringRef     string

	MonitoredStopVisits []*SIRIMonitoredStopVisit
	CancelledStopVisits []*SIRICancelledStopVisit
}

type SIRICancelledStopVisit struct {
	RecordedAtTime         time.Time
	ItemRef                string
	MonitoringRef          string
	LineRef                string
	DataFrameRef           string
	DatedVehicleJourneyRef string
	PublishedLineName      string
}

type SIRIMonitoredStopVisit struct {
	ItemIdentifier         string
	MonitoringRef          string
	StopPointRef           string
	StopPointName          string
	DatedVehicleJourneyRef string
	LineRef                string
	PublishedLineName      string
	DepartureStatus        string
	ArrivalStatus          string
	VehicleJourneyName     string
	OriginName             string
	DestinationName        string
	StopAreaObjectId       string
	DataFrameRef           string

	VehicleAtStop bool
	Monitored     bool

	Order int

	AimedArrivalTime    time.Time
	ExpectedArrivalTime time.Time
	ActualArrivalTime   time.Time

	RecordedAt            time.Time
	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
	ActualDepartureTime   time.Time

	// Attributes
	Attributes map[string]map[string]string

	// Références
	References map[string]map[string]model.Reference
}

func (response *SIRIStopMonitoringResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "stop_monitoring_response.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIStopMonitoringDelivery) BuildStopMonitoringDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "stop_monitoring_delivery.template", delivery); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (stopVisit *SIRIMonitoredStopVisit) BuildMonitoredStopVisitXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "monitored_stop_visit.template", stopVisit); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (stopVisit *SIRICancelledStopVisit) BuildCancelledStopVisitXML() (string, error) {
	var buffer bytes.Buffer
	if err := templates.ExecuteTemplate(&buffer, "cancelled_stop_visit.template", stopVisit); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}

	return buffer.String(), nil
}
