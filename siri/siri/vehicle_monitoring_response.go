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
	DirectionName           string                   `json:",omitempty"`
	OriginRef               string                   `json:",omitempty"`
	OriginName              string                   `json:",omitempty"`
	DestinationRef          string                   `json:",omitempty"`
	DestinationName         string                   `json:",omitempty"`
	Monitored               bool
	Delay                   *time.Time `json:",omitempty"`
	Bearing                 float64
	VehicleLocation         *VehicleLocation
	Occupancy               string `json:",omitempty"`
	DriverRef               string `json:",omitempty"`
}

type FramedVehicleJourneyRef struct {
	DataFrameRef           string
	DatedVehicleJourneyRef string
}

type VehicleLocation struct {
	Coordinates string  `json:",omitempty"`
	Longitude   float64 `json:",omitempty"`
	Latitude    float64 `json:",omitempty"`
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

type VehicleActivities []*VehicleActivity

func (a VehicleActivities) Len() int      { return len(a) }
func (a VehicleActivities) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type SortByVehicleMonitoringRef struct {
	VehicleActivities
}

func (s SortByVehicleMonitoringRef) Less(i, j int) bool {
	return s.VehicleActivities[i].VehicleMonitoringRef < s.VehicleActivities[j].VehicleMonitoringRef
}
