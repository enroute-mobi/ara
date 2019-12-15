package siri

import (
	"time"
)

type VehicleMonitoringDelivery struct {
	Version           string
	ResponseTimestamp time.Time `json:",omitempty"`
	RequestMessageRef string    `json:",omitempty"`
	Status            bool
	ErrorCondition    *ErrorCondition `json:",omitempty"`
	VehicleActivity   []*VehicleActivity
}

type VehicleActivity struct {
	RecordedAtTime          time.Time                `json:",omitempty"`
	ValidUntilTime          time.Time                `json:",omitempty"`
	VehicleMonitoringRef    string                   `json:",omitempty"`
	MonitoredVehicleJourney *MonitoredVehicleJourney `json:",omitempty"`
}

type MonitoredVehicleJourney struct {
	LineRef                 string                   `json:",omitempty"`
	FramedVehicleJourneyRef *FramedVehicleJourneyRef `json:",omitempty"`
	PublishedLineName       string                   `json:",omitempty"`
	OriginRef               string                   `json:",omitempty"`
	OriginName              string                   `json:",omitempty"`
	DestinationRef          string                   `json:",omitempty"`
	DestinationName         string                   `json:",omitempty"`
	Monitored               bool
	Delay                   *time.Time `json:",omitempty"`
	Bearing                 float64
	VehicleLocation         *VehicleLocation
}

type FramedVehicleJourneyRef struct {
	DataFrameRef           string
	DatedVehicleJourneyRef string
}

type VehicleLocation struct {
	Longitude float64
	Latitude  float64
}

func NewSiriLiteVehicleMonitoringDelivery() *VehicleMonitoringDelivery {
	return &VehicleMonitoringDelivery{
		Version:         "2.0:FR-IDF-2.4",
		VehicleActivity: []*VehicleActivity{},
	}
}

func NewSiriLiteVehicleActivity() *VehicleActivity {
	mvj := &MonitoredVehicleJourney{
		VehicleLocation:         &VehicleLocation{},
		FramedVehicleJourneyRef: &FramedVehicleJourneyRef{},
	}
	return &VehicleActivity{
		MonitoredVehicleJourney: mvj,
	}
}
