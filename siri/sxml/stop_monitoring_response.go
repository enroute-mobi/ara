package sxml

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringResponse struct {
	ResponseXMLStructure

	deliveries []*XMLStopMonitoringDelivery
}

type XMLStopMonitoringDelivery struct {
	DeliveryXMLStructure

	monitoringRef string

	monitoredStopVisits             []*XMLMonitoredStopVisit
	monitoredStopVisitCancellations []*XMLMonitoredStopVisitCancellation
}

type XMLMonitoredStopVisitCancellation struct {
	XMLStructure

	itemRef       string
	monitoringRef string
	recordedAt    time.Time
}

type XMLMonitoredStopVisit struct {
	XMLMonitoredVehicleJourney

	itemIdentifier string
	monitoringRef  string
	recordedAt     time.Time
}

type XMLMonitoredVehicleJourney struct {
	XMLCall

	datedVehicleJourneyRef string
	lineRef                string
	vehicleJourneyName     string
	publishedLineName      string
	dataFrameRef           string

	// Attributes
	situationRef                string
	inCongestion                string
	delay                       string
	actualQuayName              string
	aimedHeadwayInterval        string
	arrivalPlatformName         string
	arrivalProximyTest          string
	departureBoardingActivity   string
	departurePlatformName       string
	distanceFromStop            string
	expectedHeadwayInterval     string
	numberOfStopsAway           string
	platformTraversal           string
	directionName               string
	destinationName             string
	directionRef                string
	firstOrLastJourney          string
	headwayService              string
	journeyNote                 string
	journeyPatternName          string
	monitored                   Bool
	monitoringError             string
	occupancy                   string
	originAimedDepartureTime    string
	destinationAimedArrivalTime string
	originName                  string
	productCategoryRef          string
	serviceFeatureRef           string
	trainNumberRef              string
	vehicleFeatureRef           string
	vehicleMode                 string
	operatorRef                 string
	viaPlaceName                string
	originRef                   string
	placeRef                    string
	destinationRef              string
	journeyPatternRef           string
	routeRef                    string
	bearing                     string
	inPanic                     string

	// VehicleMonitoring attributes
	srsName     string
	coordinates string
	longitude   string
	latitude    string
	vehicleRef  string
	driverRef   string
}

func NewXMLStopMonitoringResponse(node xml.Node) *XMLStopMonitoringResponse {
	xmlStopMonitoringResponse := &XMLStopMonitoringResponse{}
	xmlStopMonitoringResponse.node = NewXMLNode(node)
	return xmlStopMonitoringResponse
}

func NewXMLStopMonitoringResponseFromContent(content []byte) (*XMLStopMonitoringResponse, error) {
	doc, err := gokogiri.ParseXml(content)
	if err != nil {
		return nil, err
	}
	response := NewXMLStopMonitoringResponse(doc.Root().XmlNode)
	return response, nil
}

func NewXMLStopMonitoringDelivery(node XMLNode) *XMLStopMonitoringDelivery {
	delivery := &XMLStopMonitoringDelivery{}
	delivery.node = node
	return delivery
}

func (response *XMLStopMonitoringResponse) StopMonitoringDeliveries() []*XMLStopMonitoringDelivery {
	if response.deliveries == nil {
		deliveries := []*XMLStopMonitoringDelivery{}
		nodes := response.findNodes(siri_attributes.StopMonitoringDelivery)
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLStopMonitoringDelivery(node))
		}
		response.deliveries = deliveries
	}
	return response.deliveries
}

func (delivery *XMLStopMonitoringDelivery) MonitoringRef() string {
	if delivery.monitoringRef == "" {
		delivery.monitoringRef = delivery.findStringChildContent(siri_attributes.MonitoringRef)
	}
	return delivery.monitoringRef
}

func (delivery *XMLStopMonitoringDelivery) XMLMonitoredStopVisits() []*XMLMonitoredStopVisit {
	if delivery.monitoredStopVisits == nil {
		stopVisits := []*XMLMonitoredStopVisit{}
		nodes := delivery.findNodes(siri_attributes.MonitoredStopVisit)
		for _, node := range nodes {
			stopVisits = append(stopVisits, NewXMLMonitoredStopVisit(node))
		}
		delivery.monitoredStopVisits = stopVisits
	}
	return delivery.monitoredStopVisits
}

func (delivery *XMLStopMonitoringDelivery) XMLMonitoredStopVisitCancellations() []*XMLMonitoredStopVisitCancellation {
	if delivery.monitoredStopVisitCancellations == nil {
		cancellations := []*XMLMonitoredStopVisitCancellation{}
		nodes := delivery.findNodes(siri_attributes.MonitoredStopVisitCancellation)
		for _, node := range nodes {
			cancellations = append(cancellations, NewXMLCancelledStopVisit(node))
		}
		delivery.monitoredStopVisitCancellations = cancellations
	}
	return delivery.monitoredStopVisitCancellations
}

func (cancel *XMLMonitoredStopVisitCancellation) ItemRef() string {
	if cancel.itemRef == "" {
		cancel.itemRef = cancel.findStringChildContent(siri_attributes.ItemRef)
	}
	return cancel.itemRef
}

func (cancel *XMLMonitoredStopVisitCancellation) MonitoringRef() string {
	if cancel.monitoringRef == "" {
		cancel.monitoringRef = cancel.findStringChildContent(siri_attributes.MonitoringRef)
	}
	return cancel.monitoringRef
}

func (cancel *XMLMonitoredStopVisitCancellation) RecordedAt() time.Time {
	if cancel.recordedAt.IsZero() {
		cancel.recordedAt = cancel.findTimeChildContent(siri_attributes.RecordedAtTime)
	}
	return cancel.recordedAt
}

func NewXMLCancelledStopVisit(node XMLNode) *XMLMonitoredStopVisitCancellation {
	cancelledStopVisit := &XMLMonitoredStopVisitCancellation{}
	cancelledStopVisit.node = node
	return cancelledStopVisit
}

func NewXMLMonitoredStopVisit(node XMLNode) *XMLMonitoredStopVisit {
	stopVisit := &XMLMonitoredStopVisit{}
	stopVisit.node = node
	return stopVisit
}

func (sv *XMLMonitoredStopVisit) ItemIdentifier() string {
	if sv.itemIdentifier == "" {
		sv.itemIdentifier = sv.findStringChildContent(siri_attributes.ItemIdentifier)
	}
	return sv.itemIdentifier
}

func (sv *XMLMonitoredStopVisit) MonitoringRef() string {
	if sv.monitoringRef == "" {
		sv.monitoringRef = sv.findStringChildContent(siri_attributes.MonitoringRef)
	}
	return sv.monitoringRef
}

func (sv *XMLMonitoredStopVisit) RecordedAt() time.Time {
	if sv.recordedAt.IsZero() {
		sv.recordedAt = sv.findTimeChildContent(siri_attributes.RecordedAtTime)
	}
	return sv.recordedAt
}

func (vj *XMLMonitoredVehicleJourney) DatedVehicleJourneyRef() string {
	if vj.datedVehicleJourneyRef == "" {
		vj.datedVehicleJourneyRef = vj.findStringChildContent(siri_attributes.DatedVehicleJourneyRef)
	}
	return vj.datedVehicleJourneyRef
}

func (vj *XMLMonitoredVehicleJourney) DataFrameRef() string {
	if vj.dataFrameRef == "" {
		vj.dataFrameRef = vj.findStringChildContent(siri_attributes.DataFrameRef)
	}
	return vj.dataFrameRef
}

func (vj *XMLMonitoredVehicleJourney) LineRef() string {
	if vj.lineRef == "" {
		vj.lineRef = vj.findStringChildContent(siri_attributes.LineRef)
	}
	return vj.lineRef
}

func (vj *XMLMonitoredVehicleJourney) PublishedLineName() string {
	if vj.publishedLineName == "" {
		vj.publishedLineName = vj.findStringChildContent(siri_attributes.PublishedLineName)
	}
	return vj.publishedLineName
}

// Attributes
func (vj *XMLMonitoredVehicleJourney) Delay() string {
	if vj.delay == "" {
		vj.delay = vj.findStringChildContent(siri_attributes.Delay)
	}
	return vj.delay
}

func (vj *XMLMonitoredVehicleJourney) ActualQuayName() string {
	if vj.actualQuayName == "" {
		vj.actualQuayName = vj.findStringChildContent(siri_attributes.ActualQuayName)
	}
	return vj.actualQuayName
}

func (vj *XMLMonitoredVehicleJourney) AimedHeadwayInterval() string {
	if vj.aimedHeadwayInterval == "" {
		vj.aimedHeadwayInterval = vj.findStringChildContent(siri_attributes.AimedHeadwayInterval)
	}
	return vj.aimedHeadwayInterval
}

func (vj *XMLMonitoredVehicleJourney) ArrivalPlatformName() string {
	if vj.arrivalPlatformName == "" {
		vj.arrivalPlatformName = vj.findStringChildContent(siri_attributes.ArrivalPlatformName)
	}
	return vj.arrivalPlatformName
}

func (vj *XMLMonitoredVehicleJourney) ArrivalProximyTest() string {
	if vj.arrivalProximyTest == "" {
		vj.arrivalProximyTest = vj.findStringChildContent(siri_attributes.ArrivalProximyTest)
	}
	return vj.arrivalProximyTest
}

func (vj *XMLMonitoredVehicleJourney) DepartureBoardingActivity() string {
	if vj.departureBoardingActivity == "" {
		vj.departureBoardingActivity = vj.findStringChildContent(siri_attributes.DepartureBoardingActivity)
	}
	return vj.departureBoardingActivity
}

func (vj *XMLMonitoredVehicleJourney) DeparturePlatformName() string {
	if vj.departurePlatformName == "" {
		vj.departurePlatformName = vj.findStringChildContent(siri_attributes.DeparturePlatformName)
	}
	return vj.departurePlatformName
}

func (vj *XMLMonitoredVehicleJourney) DistanceFromStop() string {
	if vj.distanceFromStop == "" {
		vj.distanceFromStop = vj.findStringChildContent(siri_attributes.DistanceFromStop)
	}
	return vj.distanceFromStop
}

func (vj *XMLMonitoredVehicleJourney) ExpectedHeadwayInterval() string {
	if vj.expectedHeadwayInterval == "" {
		vj.expectedHeadwayInterval = vj.findStringChildContent(siri_attributes.ExpectedHeadwayInterval)
	}
	return vj.expectedHeadwayInterval
}

func (vj *XMLMonitoredVehicleJourney) NumberOfStopsAway() string {
	if vj.numberOfStopsAway == "" {
		vj.numberOfStopsAway = vj.findStringChildContent(siri_attributes.NumberOfStopsAway)
	}
	return vj.numberOfStopsAway
}

func (vj *XMLMonitoredVehicleJourney) PlatformTraversal() string {
	if vj.platformTraversal == "" {
		vj.platformTraversal = vj.findStringChildContent(siri_attributes.PlatformTraversal)
	}
	return vj.platformTraversal
}

func (vj *XMLMonitoredVehicleJourney) DirectionName() string {
	if vj.directionName == "" {
		vj.directionName = vj.findStringChildContent(siri_attributes.DirectionName)
	}
	return vj.directionName
}

func (vj *XMLMonitoredVehicleJourney) DestinationName() string {
	if vj.destinationName == "" {
		vj.destinationName = vj.findStringChildContent(siri_attributes.DestinationName)
	}
	return vj.destinationName
}

func (vj *XMLMonitoredVehicleJourney) DirectionRef() string {
	if vj.directionRef == "" {
		vj.directionRef = vj.findStringChildContent(siri_attributes.DirectionRef)
	}
	return vj.directionRef
}

func (vj *XMLMonitoredVehicleJourney) Bearing() string {
	if vj.bearing == "" {
		vj.bearing = vj.findStringChildContent(siri_attributes.Bearing)
	}
	return vj.bearing
}

func (vj *XMLMonitoredVehicleJourney) InPanic() string {
	if vj.inPanic == "" {
		vj.inPanic = vj.findStringChildContent(siri_attributes.InPanic)
	}
	return vj.inPanic
}

func (vj *XMLMonitoredVehicleJourney) SituationRef() string {
	if vj.situationRef == "" {
		vj.situationRef = vj.findStringChildContent(siri_attributes.SituationRef)
	}
	return vj.situationRef
}

func (vj *XMLMonitoredVehicleJourney) InCongestion() string {
	if vj.inCongestion == "" {
		vj.inCongestion = vj.findStringChildContent(siri_attributes.InCongestion)
	}
	return vj.inPanic
}

func (vj *XMLMonitoredVehicleJourney) HeadwayService() string {
	if vj.headwayService == "" {
		vj.headwayService = vj.findStringChildContent(siri_attributes.HeadwayService)
	}
	return vj.headwayService
}

func (vj *XMLMonitoredVehicleJourney) FirstOrLastJourney() string {
	if vj.firstOrLastJourney == "" {
		vj.firstOrLastJourney = vj.findStringChildContent(siri_attributes.FirstOrLastJourney)
	}
	return vj.firstOrLastJourney
}

func (vj *XMLMonitoredVehicleJourney) JourneyNote() string {
	if vj.journeyNote == "" {
		vj.journeyNote = vj.findStringChildContent(siri_attributes.JourneyNote)
	}
	return vj.journeyNote
}

func (vj *XMLMonitoredVehicleJourney) JourneyPatternName() string {
	if vj.journeyPatternName == "" {
		vj.journeyPatternName = vj.findStringChildContent(siri_attributes.JourneyPatternName)
	}
	return vj.journeyPatternName
}

func (vj *XMLMonitoredVehicleJourney) Monitored() bool {
	if !vj.monitored.Defined {
		vj.monitored.SetValue(vj.findBoolChildContent(siri_attributes.Monitored))
	}
	return vj.monitored.Value
}

func (vj *XMLMonitoredVehicleJourney) MonitoringError() string {
	if vj.monitoringError == "" {
		vj.monitoringError = vj.findStringChildContent(siri_attributes.MonitoringError)
	}
	return vj.monitoringError
}

func (vj *XMLMonitoredVehicleJourney) Occupancy() string {
	if vj.occupancy == "" {
		vj.occupancy = vj.findStringChildContent(siri_attributes.Occupancy)
	}
	return vj.occupancy
}

func (vj *XMLMonitoredVehicleJourney) OriginAimedDepartureTime() string {
	if vj.originAimedDepartureTime == "" {
		vj.originAimedDepartureTime = vj.findStringChildContent(siri_attributes.OriginAimedDepartureTime)
	}
	return vj.originAimedDepartureTime
}

func (vj *XMLMonitoredVehicleJourney) DestinationAimedArrivalTime() string {
	if vj.destinationAimedArrivalTime == "" {
		vj.destinationAimedArrivalTime = vj.findStringChildContent(siri_attributes.DestinationAimedArrivalTime)
	}
	return vj.destinationAimedArrivalTime
}

func (vj *XMLMonitoredVehicleJourney) OriginName() string {
	if vj.originName == "" {
		vj.originName = vj.findStringChildContent(siri_attributes.OriginName)
	}
	return vj.originName
}

func (vj *XMLMonitoredVehicleJourney) ProductCategoryRef() string {
	if vj.productCategoryRef == "" {
		vj.productCategoryRef = vj.findStringChildContent(siri_attributes.ProductCategoryRef)
	}
	return vj.productCategoryRef
}

func (vj *XMLMonitoredVehicleJourney) ServiceFeatureRef() string {
	if vj.serviceFeatureRef == "" {
		vj.serviceFeatureRef = vj.findStringChildContent(siri_attributes.ServiceFeatureRef)
	}
	return vj.serviceFeatureRef
}

func (vj *XMLMonitoredVehicleJourney) TrainNumberRef() string {
	if vj.trainNumberRef == "" {
		vj.trainNumberRef = vj.findStringChildContent(siri_attributes.TrainNumberRef)
	}
	return vj.trainNumberRef
}

func (vj *XMLMonitoredVehicleJourney) VehicleFeatureRef() string {
	if vj.vehicleFeatureRef == "" {
		vj.vehicleFeatureRef = vj.findStringChildContent(siri_attributes.VehicleFeatureRef)
	}
	return vj.vehicleFeatureRef
}

func (vj *XMLMonitoredVehicleJourney) VehicleJourneyName() string {
	if vj.vehicleJourneyName == "" {
		vj.vehicleJourneyName = vj.findStringChildContent(siri_attributes.VehicleJourneyName)
	}
	return vj.vehicleJourneyName
}

func (vj *XMLMonitoredVehicleJourney) VehicleMode() string {
	if vj.vehicleMode == "" {
		vj.vehicleMode = vj.findStringChildContent(siri_attributes.VehicleMode)
	}
	return vj.vehicleMode
}

func (vj *XMLMonitoredVehicleJourney) ViaPlaceName() string {
	if vj.viaPlaceName == "" {
		vj.viaPlaceName = vj.findStringChildContent(siri_attributes.PlaceName)
	}
	return vj.viaPlaceName
}

// References

func (vj *XMLMonitoredVehicleJourney) OriginRef() string {
	if vj.originRef == "" {
		vj.originRef = vj.findStringChildContent(siri_attributes.OriginRef)
	}
	return vj.originRef
}

func (vj *XMLMonitoredVehicleJourney) PlaceRef() string {
	if vj.placeRef == "" {
		vj.placeRef = vj.findStringChildContent(siri_attributes.PlaceRef)
	}
	return vj.placeRef
}

func (vj *XMLMonitoredVehicleJourney) DestinationRef() string {
	if vj.destinationRef == "" {
		vj.destinationRef = vj.findStringChildContent(siri_attributes.DestinationRef)
	}
	return vj.destinationRef
}

func (vj *XMLMonitoredVehicleJourney) JourneyPatternRef() string {
	if vj.journeyPatternRef == "" {
		vj.journeyPatternRef = vj.findStringChildContent(siri_attributes.JourneyPatternRef)
	}
	return vj.journeyPatternRef
}

func (vj *XMLMonitoredVehicleJourney) RouteRef() string {
	if vj.routeRef == "" {
		vj.routeRef = vj.findStringChildContent(siri_attributes.RouteRef)
	}
	return vj.routeRef
}

func (vj *XMLMonitoredVehicleJourney) OperatorRef() string {
	if vj.operatorRef == "" {
		vj.operatorRef = vj.findStringChildContent(siri_attributes.OperatorRef)
	}
	return vj.operatorRef
}

func (vj *XMLMonitoredVehicleJourney) Coordinates() string {
	if vj.coordinates == "" {
		vj.coordinates = vj.findStringChildContent(siri_attributes.Coordinates)
	}
	return vj.coordinates
}

func (vj *XMLMonitoredVehicleJourney) Longitude() string {
	if vj.longitude == "" {
		vj.longitude = vj.findStringChildContent(siri_attributes.Longitude)
	}
	return vj.longitude
}

func (vj *XMLMonitoredVehicleJourney) Latitude() string {
	if vj.latitude == "" {
		vj.latitude = vj.findStringChildContent(siri_attributes.Latitude)
	}
	return vj.latitude
}

func (vj *XMLMonitoredVehicleJourney) VehicleRef() string {
	if vj.vehicleRef == "" {
		vref := vj.findStringChildContent(siri_attributes.VehicleRef)
		if vref != "" {
			vj.vehicleRef = vref
		} else {
			vj.vehicleRef = vj.findStringChildContent(siri_attributes.VehicleMonitoringRef)
		}
	}
	return vj.vehicleRef
}

func (vj *XMLMonitoredVehicleJourney) DriverRef() string {
	if vj.driverRef == "" {
		vj.driverRef = vj.findStringChildContent(siri_attributes.DriverRef)
	}
	return vj.driverRef
}

func (vj *XMLMonitoredVehicleJourney) SRSName() string {
	if vj.srsName == "" {
		vj.srsName = vj.findChildAttribute("VehicleLocation", "srsName")
	}
	return vj.srsName
}

// Test methods

func (vj *XMLMonitoredVehicleJourney) SetLongitude(s string) {
	vj.longitude = s
}

func (vj *XMLMonitoredVehicleJourney) SetLatitude(s string) {
	vj.latitude = s
}

func (vj *XMLMonitoredVehicleJourney) SetSRSName(s string) {
	vj.srsName = s
}

func (vj *XMLMonitoredVehicleJourney) SetCoordinates(s string) {
	vj.coordinates = s
}
