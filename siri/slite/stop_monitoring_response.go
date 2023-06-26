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
	Order                 *int      `json:"Order,omitempty"`
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
	Monitored               *bool                   `json:"Monitored"`
	MonitoredCall           MonitoredCall           `json:"MonitoredCall,omitempty"`
}
type MonitoredStopVisit struct {
	RecordedAtTime          time.Time               `json:"RecordedAtTime,omitempty"`
	ItemIdentifier          *string                 `json:"ItemIdentifier"`
	MonitoringRef           string                  `json:"MonitoringRef,omitempty"`
	StopPointRef            *string                 `json:"StopPointRef"`
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

func (msv *MonitoredStopVisit) HasOrder() bool {
	return msv.MonitoredVehicleJourney.MonitoredCall.Order != nil
}

// When Monitored is not defined, it should be true by default
// see ARA-1240 "Special cases"
func (msv *MonitoredStopVisit) GetMonitored() bool {
	if monitored := msv.MonitoredVehicleJourney.Monitored; monitored != nil {
		return *monitored
	}

	return true
}

// When StopPointRef is not defined, we should use MonitoringRef value.
// see ARA-1240 "Special cases"
func (msv *MonitoredStopVisit) GetStopPointRef() string {
	if stopPointRef := msv.StopPointRef; stopPointRef != nil {
		return *stopPointRef
	}

	return msv.MonitoringRef
}

// When ItemIdentifier is not defined, we should use
// DatedVehicleJourneyRef + Order to create a default value.
// see ARA-1240 "Special cases"
func (msv *MonitoredStopVisit) GetItemIdentifier() string {
	if itemIdentifier := msv.ItemIdentifier; itemIdentifier != nil {
		return *itemIdentifier
	}

	return fmt.Sprintf("%s-%s",
		msv.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef,
		strconv.Itoa(*msv.MonitoredVehicleJourney.MonitoredCall.Order),
	)
}
