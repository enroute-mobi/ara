package siri

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
)

type SIRIVehicleMonitoringResponse struct {
	ResponseTimestamp         time.Time `json:",omitempty"`
	ProducerRef               string    `json:",omitempty"`
	ResponseMessageIdentifier string    `json:",omitempty"`
	RequestMessageRef         string    `json:",omitempty"`
	Address                   string    `json:".omitempty"`

	SIRIVehicleMonitoringDelivery
}

type SIRIVehicleMonitoringDelivery struct {
	Version           string
	ResponseTimestamp time.Time `json:",omitempty"`
	RequestMessageRef string    `json:",omitempty"`
	Status            bool
	ErrorCondition    *ErrorCondition `json:",omitempty"`
	VehicleActivity   []*SIRIVehicleActivity

	LineRefs           map[string]struct{} `json:"-"`
	VehicleJourneyRefs map[string]struct{} `json:"-"`
	VehicleRefs        map[string]struct{} `json:"-"`
}

type SIRIVehicleActivity struct {
	RecordedAtTime       time.Time                 `json:",omitempty"`
	ValidUntilTime       time.Time                 `json:",omitempty"`
	VehicleMonitoringRef string                    `json:",omitempty"`
	ProgressBetweenStops *SIRIProgressBetweenStops `json:",omitempty"`

	MonitoredVehicleJourney *SIRIMonitoredVehicleJourney `json:",omitempty"`
}

type SIRIMonitoredVehicleJourney struct {
	LineRef                 string                       `json:",omitempty"`
	FramedVehicleJourneyRef *SIRIFramedVehicleJourneyRef `json:",omitempty"`
	PublishedLineName       string                       `json:",omitempty"`
	DirectionName           string                       `json:",omitempty"`
	DirectionType           string                       `json:",omitempty"`
	OriginRef               string                       `json:",omitempty"`
	OriginName              string                       `json:",omitempty"`
	DestinationRef          string                       `json:",omitempty"`
	DestinationName         string                       `json:",omitempty"`
	Monitored               bool
	Delay                   *time.Time `json:",omitempty"`
	Bearing                 float64
	VehicleLocation         *SIRIVehicleLocation
	Occupancy               string         `json:",omitempty"`
	DriverRef               string         `json:",omitempty"`
	JourneyPatternRef       string         `json:",omitempty"`
	JourneyPatternName      string         `json:",omitempty"`
	VehicleRef              string         `json:",omitempty"`
	MonitoredCall           *MonitoredCall `json:",omitempty"`
}

type MonitoredCall struct {
	StopPointRef          string    `json:",omitempty"`
	StopPointName         string    `json:",omitempty"`
	VehicleAtStop         bool      `json:",omitempty"`
	DestinationDisplay    string    `json:",omitempty"`
	ExpectedArrivalTime   time.Time `json:",omitempty"`
	ExpectedDepartureTime time.Time `json:",omitempty"`
	DepartureStatus       string    `json:",omitempty"`
	Order                 *int      `json:",omitempty"`
	AimedArrivalTime      time.Time `json:",omitempty"`
	ArrivalPlatformName   string    `json:",omitempty"`
	AimedDepartureTime    time.Time `json:",omitempty"`
	ArrivalStatus         string    `json:",omitempty"`
	ActualArrivalTime     time.Time `json:",omitempty"`
	ActualDepartureTime   time.Time `json:",omitempty"`
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

type SIRIProgressBetweenStops struct {
	LinkDistance float64 `json:",omitempty"`
	Percentage   float64 `json:",omitempty"`
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

func (response *SIRIVehicleMonitoringResponse) BuildXML(envelopeType ...string) (string, error) {
	var buffer bytes.Buffer
	var envType string
	var templateName string

	if len(envelopeType) != 0 && envelopeType[0] != "soap" && envelopeType[0] != "" {
		envType = "_" + envelopeType[0]
	}

	templateName = fmt.Sprintf("vehicle_monitoring_response%s.template", envType)

	if err := templates.ExecuteTemplate(&buffer, templateName, response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return buffer.String(), nil
}

func (response *SIRIVehicleMonitoringDelivery) BuildVehicleMonitoringDeliveryXML() (string, error) {
	var buffer bytes.Buffer

	if err := templates.ExecuteTemplate(&buffer, "vehicle_monitoring_delivery.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (response *SIRIVehicleActivity) BuildVehicleActivityXML() (string, error) {
	var buffer bytes.Buffer

	if err := templates.ExecuteTemplate(&buffer, "vehicle_activity.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (response *SIRIMonitoredVehicleJourney) BuildMonitoredVehicleJourneyXML() (string, error) {
	var buffer bytes.Buffer

	if err := templates.ExecuteTemplate(&buffer, "monitored_vehicle_journey.template", response); err != nil {
		logger.Log.Debugf("Error while executing template: %v", err)
		return "", err
	}
	return strings.TrimSpace(buffer.String()), nil
}

func (delivery *SIRIVehicleMonitoringDelivery) ErrorString() string {
	return fmt.Sprintf("%v: %v", delivery.errorType(), delivery.ErrorCondition.ErrorText)
}

func (delivery *SIRIVehicleMonitoringDelivery) errorType() string {
	if delivery.ErrorCondition.ErrorType == siri_attributes.OtherError {
		return fmt.Sprintf("%v %v", delivery.ErrorCondition.ErrorType, delivery.ErrorCondition.ErrorType)
	}
	return delivery.ErrorCondition.ErrorType
}
