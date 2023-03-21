package siri

import (
	"time"
)

type SIRIVehicleMonitoringDelivery struct {
	Version           string
	ResponseTimestamp time.Time `json:",omitempty"`
	RequestMessageRef string    `json:",omitempty"`
	Status            bool
	ErrorCondition    *ErrorCondition `json:",omitempty"`
	VehicleActivity   []*SIRIVehicleActivity
}

type SIRIVehicleActivity struct {
	RecordedAtTime          time.Time                    `json:",omitempty"`
	ValidUntilTime          time.Time                    `json:",omitempty"`
	VehicleMonitoringRef    string                       `json:",omitempty"`
	MonitoredVehicleJourney *SIRIMonitoredVehicleJourney `json:",omitempty"`
}

type SIRIMonitoredVehicleJourney struct {
	LineRef                 string                       `json:",omitempty"`
	FramedVehicleJourneyRef *SIRIFramedVehicleJourneyRef `json:",omitempty"`
	PublishedLineName       string                       `json:",omitempty"`
	DirectionName           string                       `json:",omitempty"`
	OriginRef               string                       `json:",omitempty"`
	OriginName              string                       `json:",omitempty"`
	DestinationRef          string                       `json:",omitempty"`
	DestinationName         string                       `json:",omitempty"`
	Monitored               bool
	Delay                   *time.Time `json:",omitempty"`
	Bearing                 float64
	VehicleLocation         *SIRIVehicleLocation
	Occupancy               string `json:",omitempty"`
	DriverRef               string `json:",omitempty"`
}

type SIRIFramedVehicleJourneyRef struct {
	DataFrameRef           string
	DatedVehicleJourneyRef string
}

type SIRIVehicleLocation struct {
	Coordinates string  `json:",omitempty"`
	Longitude   float64 `json:",omitempty"`
	Latitude    float64 `json:",omitempty"`
}

func NewSiriLiteVehicleMonitoringDelivery() *SIRIVehicleMonitoringDelivery {
	return &SIRIVehicleMonitoringDelivery{
		Version:         "2.0:FR-IDF-2.4",
		VehicleActivity: []*SIRIVehicleActivity{},
	}
}

func NewSiriLiteVehicleActivity() *SIRIVehicleActivity {
	mvj := &SIRIMonitoredVehicleJourney{
		VehicleLocation:         &SIRIVehicleLocation{},
		FramedVehicleJourneyRef: &SIRIFramedVehicleJourneyRef{},
	}
	return &SIRIVehicleActivity{
		MonitoredVehicleJourney: mvj,
	}
}

type VehicleActivities []*SIRIVehicleActivity

func (a VehicleActivities) Len() int      { return len(a) }
func (a VehicleActivities) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type SortByVehicleMonitoringRef struct {
	VehicleActivities
}

func (s SortByVehicleMonitoringRef) Less(i, j int) bool {
	return s.VehicleActivities[i].VehicleMonitoringRef < s.VehicleActivities[j].VehicleMonitoringRef
}
