package slite

import (
	"fmt"
	"strconv"
	"time"
)

type SIRILiteStopMonitoring struct {
	Siri Siri `json:"Siri"`
}

type FramedVehicleJourneyRef struct {
	DataFrameRef           string `json:"DataFrameRef,omitempty"`
	DatedVehicleJourneyRef string `json:"DatedVehicleJourneyRef,omitempty"`
}
type MonitoredCall struct {
	StopPointName         string    `json:"StopPointName,omitempty"`
	VehicleAtStop         bool      `json:"VehicleAtStop,omitempty"`
	DestinationDisplay    string    `json:"DestinationDisplay,omitempty"`
	ExpectedArrivalTime   time.Time `json:"ExpectedArrivalTime,omitempty"`
	ExpectedDepartureTime time.Time `json:"ExpectedDepartureTime,omitempty"`
	DepartureStatus       string    `json:"DepartureStatus,omitempty"`
	Order                 int       `json:"Order,omitempty"`
	AimedArrivalTime      time.Time `json:"AimedArrivalTime,omitempty"`
	ArrivalPlatformName   string    `json:"ArrivalPlatformName,omitempty"`
	AimedDepartureTime    time.Time `json:"AimedDepartureTime,omitempty"`
	ArrivalStatus         string    `json:"ArrivalStatus,omitempty"`
	ActualArrivalTime     time.Time `json:"ActualArrivalTime,omitempty"`
	ActualDepartureTime   time.Time `json:"ActualDepartureTime,omitempty"`
}
type MonitoredVehicleJourney struct {
	LineRef                 string                  `json:"LineRef,omitempty"`
	OperatorRef             string                  `json:"OperatorRef,omitempty"`
	FramedVehicleJourneyRef FramedVehicleJourneyRef `json:"FramedVehicleJourneyRef,omitempty"`
	DestinationRef          string                  `json:"DestinationRef,omitempty"`
	DestinationName         string                  `json:"DestinationName,omitempty"`
	JourneyNote             string                  `json:"JourneyNote,omitempty"`
	MonitoredCall           MonitoredCall           `json:"MonitoredCall,omitempty"`
}
type MonitoredStopVisit struct {
	RecordedAtTime          time.Time               `json:"RecordedAtTime,omitempty"`
	ItemIdentifier          string                  `json:"ItemIdentifier,omitempty"`
	MonitoringRef           string                  `json:"MonitoringRef,omitempty"`
	StopPointRef            string                  `json:"StopPointRef,omitempty"`
	MonitoredVehicleJourney MonitoredVehicleJourney `json:"MonitoredVehicleJourney,omitempty"`
}
type StopMonitoringDelivery struct {
	ResponseTimestamp  time.Time            `json:"ResponseTimestamp,omitempty"`
	Version            string               `json:"Version,omitempty"`
	Status             string               `json:"Status,omitempty"`
	MonitoredStopVisit []MonitoredStopVisit `json:"MonitoredStopVisit,omitempty"`
	ErrorCondition     ErrorCondition       `json:"ErrorCondition,omitempty"`
}

type ErrorInformation struct {
	ErrorText        string `json:"ErrorText,omitempty"`
	ErrorDescription string `json:"ErrorDescription,omitempty"`
}
type ErrorCondition struct {
	ErrorInformation ErrorInformation `json:"ErrorInformation,omitempty"`
}

type ServiceDelivery struct {
	ResponseTimestamp         time.Time                `json:"ResponseTimestamp,omitempty"`
	ProducerRef               string                   `json:"ProducerRef,omitempty"`
	ResponseMessageIdentifier string                   `json:"ResponseMessageIdentifier,omitempty"`
	StopMonitoringDelivery    []StopMonitoringDelivery `json:"StopMonitoringDelivery,omitempty"`
}
type Siri struct {
	ServiceDelivery ServiceDelivery `json:"ServiceDelivery,omitempty"`
}

func (msv *MonitoredStopVisit) GetStopPointRef() string {
	if msv.StopPointRef != "" {
		return msv.StopPointRef
	}

	return msv.MonitoringRef
}

func (msv *MonitoredStopVisit) GetItemIdentifier() string {
	if msv.ItemIdentifier != "" {
		return msv.ItemIdentifier
	}

	identifier := fmt.Sprintf("%s-%s",
		msv.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef,
		strconv.Itoa(msv.MonitoredVehicleJourney.MonitoredCall.Order),
	)

	return identifier
}
