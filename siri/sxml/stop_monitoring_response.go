package sxml

import (
	"time"

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
	vehicleFeature              string
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
		nodes := response.findNodes("StopMonitoringDelivery")
		for _, node := range nodes {
			deliveries = append(deliveries, NewXMLStopMonitoringDelivery(node))
		}
		response.deliveries = deliveries
	}
	return response.deliveries
}

func (delivery *XMLStopMonitoringDelivery) MonitoringRef() string {
	if delivery.monitoringRef == "" {
		delivery.monitoringRef = delivery.findStringChildContent("MonitoringRef")
	}
	return delivery.monitoringRef
}

func (delivery *XMLStopMonitoringDelivery) XMLMonitoredStopVisits() []*XMLMonitoredStopVisit {
	if delivery.monitoredStopVisits == nil {
		stopVisits := []*XMLMonitoredStopVisit{}
		nodes := delivery.findNodes("MonitoredStopVisit")
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
		nodes := delivery.findNodes("MonitoredStopVisitCancellation")
		for _, node := range nodes {
			cancellations = append(cancellations, NewXMLCancelledStopVisit(node))
		}
		delivery.monitoredStopVisitCancellations = cancellations
	}
	return delivery.monitoredStopVisitCancellations
}

func (cancel *XMLMonitoredStopVisitCancellation) ItemRef() string {
	if cancel.itemRef == "" {
		cancel.itemRef = cancel.findStringChildContent("ItemRef")
	}
	return cancel.itemRef
}

func (cancel *XMLMonitoredStopVisitCancellation) MonitoringRef() string {
	if cancel.monitoringRef == "" {
		cancel.monitoringRef = cancel.findStringChildContent("MonitoringRef")
	}
	return cancel.monitoringRef
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
		sv.itemIdentifier = sv.findStringChildContent("ItemIdentifier")
	}
	return sv.itemIdentifier
}

func (sv *XMLMonitoredStopVisit) MonitoringRef() string {
	if sv.monitoringRef == "" {
		sv.monitoringRef = sv.findStringChildContent("MonitoringRef")
	}
	return sv.monitoringRef
}

func (sv *XMLMonitoredStopVisit) RecordedAt() time.Time {
	if sv.recordedAt.IsZero() {
		sv.recordedAt = sv.findTimeChildContent("RecordedAtTime")
	}
	return sv.recordedAt
}

func (vj *XMLMonitoredVehicleJourney) DatedVehicleJourneyRef() string {
	if vj.datedVehicleJourneyRef == "" {
		vj.datedVehicleJourneyRef = vj.findStringChildContent("DatedVehicleJourneyRef")
	}
	return vj.datedVehicleJourneyRef
}

func (vj *XMLMonitoredVehicleJourney) DataFrameRef() string {
	if vj.dataFrameRef == "" {
		vj.dataFrameRef = vj.findStringChildContent("DataFrameRef")
	}
	return vj.dataFrameRef
}

func (vj *XMLMonitoredVehicleJourney) LineRef() string {
	if vj.lineRef == "" {
		vj.lineRef = vj.findStringChildContent("LineRef")
	}
	return vj.lineRef
}

func (vj *XMLMonitoredVehicleJourney) PublishedLineName() string {
	if vj.publishedLineName == "" {
		vj.publishedLineName = vj.findStringChildContent("PublishedLineName")
	}
	return vj.publishedLineName
}

// Attributes
func (vj *XMLMonitoredVehicleJourney) Delay() string {
	if vj.delay == "" {
		vj.delay = vj.findStringChildContent("Delay")
	}
	return vj.delay
}

func (vj *XMLMonitoredVehicleJourney) ActualQuayName() string {
	if vj.actualQuayName == "" {
		vj.actualQuayName = vj.findStringChildContent("ActualQuayName")
	}
	return vj.actualQuayName
}

func (vj *XMLMonitoredVehicleJourney) AimedHeadwayInterval() string {
	if vj.aimedHeadwayInterval == "" {
		vj.aimedHeadwayInterval = vj.findStringChildContent("AimedHeadwayInterval")
	}
	return vj.aimedHeadwayInterval
}

func (vj *XMLMonitoredVehicleJourney) ArrivalPlatformName() string {
	if vj.arrivalPlatformName == "" {
		vj.arrivalPlatformName = vj.findStringChildContent("ArrivalPlatformName")
	}
	return vj.arrivalPlatformName
}

func (vj *XMLMonitoredVehicleJourney) ArrivalProximyTest() string {
	if vj.arrivalProximyTest == "" {
		vj.arrivalProximyTest = vj.findStringChildContent("ArrivalProximyTest")
	}
	return vj.arrivalProximyTest
}

func (vj *XMLMonitoredVehicleJourney) DepartureBoardingActivity() string {
	if vj.departureBoardingActivity == "" {
		vj.departureBoardingActivity = vj.findStringChildContent("DepartureBoardingActivity")
	}
	return vj.departureBoardingActivity
}

func (vj *XMLMonitoredVehicleJourney) DeparturePlatformName() string {
	if vj.departurePlatformName == "" {
		vj.departurePlatformName = vj.findStringChildContent("DeparturePlatformName")
	}
	return vj.departurePlatformName
}

func (vj *XMLMonitoredVehicleJourney) DistanceFromStop() string {
	if vj.distanceFromStop == "" {
		vj.distanceFromStop = vj.findStringChildContent("DistanceFromStop")
	}
	return vj.distanceFromStop
}

func (vj *XMLMonitoredVehicleJourney) ExpectedHeadwayInterval() string {
	if vj.expectedHeadwayInterval == "" {
		vj.expectedHeadwayInterval = vj.findStringChildContent("ExpectedHeadwayInterval")
	}
	return vj.expectedHeadwayInterval
}

func (vj *XMLMonitoredVehicleJourney) NumberOfStopsAway() string {
	if vj.numberOfStopsAway == "" {
		vj.numberOfStopsAway = vj.findStringChildContent("NumberOfStopsAway")
	}
	return vj.numberOfStopsAway
}

func (vj *XMLMonitoredVehicleJourney) PlatformTraversal() string {
	if vj.platformTraversal == "" {
		vj.platformTraversal = vj.findStringChildContent("PlatformTraversal")
	}
	return vj.platformTraversal
}

func (vj *XMLMonitoredVehicleJourney) DirectionName() string {
	if vj.directionName == "" {
		vj.directionName = vj.findStringChildContent("DirectionName")
	}
	return vj.directionName
}

func (vj *XMLMonitoredVehicleJourney) DestinationName() string {
	if vj.destinationName == "" {
		vj.destinationName = vj.findStringChildContent("DestinationName")
	}
	return vj.destinationName
}

func (vj *XMLMonitoredVehicleJourney) DirectionRef() string {
	if vj.directionRef == "" {
		vj.directionRef = vj.findStringChildContent("DirectionRef")
	}
	return vj.directionRef
}

func (vj *XMLMonitoredVehicleJourney) Bearing() string {
	if vj.bearing == "" {
		vj.bearing = vj.findStringChildContent("Bearing")
	}
	return vj.bearing
}

func (vj *XMLMonitoredVehicleJourney) InPanic() string {
	if vj.inPanic == "" {
		vj.inPanic = vj.findStringChildContent("InPanic")
	}
	return vj.inPanic
}

func (vj *XMLMonitoredVehicleJourney) SituationRef() string {
	if vj.situationRef == "" {
		vj.situationRef = vj.findStringChildContent("SituationRef")
	}
	return vj.situationRef
}

func (vj *XMLMonitoredVehicleJourney) InCongestion() string {
	if vj.inCongestion == "" {
		vj.inCongestion = vj.findStringChildContent("InCongestion")
	}
	return vj.inPanic
}

func (vj *XMLMonitoredVehicleJourney) HeadwayService() string {
	if vj.headwayService == "" {
		vj.headwayService = vj.findStringChildContent("HeadwayService")
	}
	return vj.headwayService
}

func (vj *XMLMonitoredVehicleJourney) FirstOrLastJourney() string {
	if vj.firstOrLastJourney == "" {
		vj.firstOrLastJourney = vj.findStringChildContent("FirstOrLastJourney")
	}
	return vj.firstOrLastJourney
}

func (vj *XMLMonitoredVehicleJourney) JourneyNote() string {
	if vj.journeyNote == "" {
		vj.journeyNote = vj.findStringChildContent("JourneyNote")
	}
	return vj.journeyNote
}

func (vj *XMLMonitoredVehicleJourney) JourneyPatternName() string {
	if vj.journeyPatternName == "" {
		vj.journeyPatternName = vj.findStringChildContent("JourneyPatternName")
	}
	return vj.journeyPatternName
}

func (vj *XMLMonitoredVehicleJourney) Monitored() bool {
	if !vj.monitored.Defined {
		vj.monitored.SetValue(vj.findBoolChildContent("Monitored"))
	}
	return vj.monitored.Value
}

func (vj *XMLMonitoredVehicleJourney) MonitoringError() string {
	if vj.monitoringError == "" {
		vj.monitoringError = vj.findStringChildContent("MonitoringError")
	}
	return vj.monitoringError
}

func (vj *XMLMonitoredVehicleJourney) Occupancy() string {
	if vj.occupancy == "" {
		vj.occupancy = vj.findStringChildContent("Occupancy")
	}
	return vj.occupancy
}

func (vj *XMLMonitoredVehicleJourney) OriginAimedDepartureTime() string {
	if vj.originAimedDepartureTime == "" {
		vj.originAimedDepartureTime = vj.findStringChildContent("OriginAimedDepartureTime")
	}
	return vj.originAimedDepartureTime
}

func (vj *XMLMonitoredVehicleJourney) DestinationAimedArrivalTime() string {
	if vj.destinationAimedArrivalTime == "" {
		vj.destinationAimedArrivalTime = vj.findStringChildContent("DestinationAimedArrivalTime")
	}
	return vj.destinationAimedArrivalTime
}

func (vj *XMLMonitoredVehicleJourney) OriginName() string {
	if vj.originName == "" {
		vj.originName = vj.findStringChildContent("OriginName")
	}
	return vj.originName
}

func (vj *XMLMonitoredVehicleJourney) ProductCategoryRef() string {
	if vj.productCategoryRef == "" {
		vj.productCategoryRef = vj.findStringChildContent("ProductCategoryRef")
	}
	return vj.productCategoryRef
}

func (vj *XMLMonitoredVehicleJourney) ServiceFeatureRef() string {
	if vj.serviceFeatureRef == "" {
		vj.serviceFeatureRef = vj.findStringChildContent("ServiceFeatureRef")
	}
	return vj.serviceFeatureRef
}

func (vj *XMLMonitoredVehicleJourney) TrainNumberRef() string {
	if vj.trainNumberRef == "" {
		vj.trainNumberRef = vj.findStringChildContent("TrainNumberRef")
	}
	return vj.trainNumberRef
}

func (vj *XMLMonitoredVehicleJourney) VehicleFeature() string {
	if vj.vehicleFeature == "" {
		vj.vehicleFeature = vj.findStringChildContent("VehicleFeature")
	}
	return vj.vehicleFeature
}

func (vj *XMLMonitoredVehicleJourney) VehicleJourneyName() string {
	if vj.vehicleJourneyName == "" {
		vj.vehicleJourneyName = vj.findStringChildContent("VehicleJourneyName")
	}
	return vj.vehicleJourneyName
}

func (vj *XMLMonitoredVehicleJourney) VehicleMode() string {
	if vj.vehicleMode == "" {
		vj.vehicleMode = vj.findStringChildContent("VehicleMode")
	}
	return vj.vehicleMode
}

func (vj *XMLMonitoredVehicleJourney) ViaPlaceName() string {
	if vj.viaPlaceName == "" {
		vj.viaPlaceName = vj.findStringChildContent("PlaceName")
	}
	return vj.viaPlaceName
}

// References

func (vj *XMLMonitoredVehicleJourney) OriginRef() string {
	if vj.originRef == "" {
		vj.originRef = vj.findStringChildContent("OriginRef")
	}
	return vj.originRef
}

func (vj *XMLMonitoredVehicleJourney) PlaceRef() string {
	if vj.placeRef == "" {
		vj.placeRef = vj.findStringChildContent("PlaceRef")
	}
	return vj.placeRef
}

func (vj *XMLMonitoredVehicleJourney) DestinationRef() string {
	if vj.destinationRef == "" {
		vj.destinationRef = vj.findStringChildContent("DestinationRef")
	}
	return vj.destinationRef
}

func (vj *XMLMonitoredVehicleJourney) JourneyPatternRef() string {
	if vj.journeyPatternRef == "" {
		vj.journeyPatternRef = vj.findStringChildContent("JourneyPatternRef")
	}
	return vj.journeyPatternRef
}

func (vj *XMLMonitoredVehicleJourney) RouteRef() string {
	if vj.routeRef == "" {
		vj.routeRef = vj.findStringChildContent("RouteRef")
	}
	return vj.routeRef
}

func (vj *XMLMonitoredVehicleJourney) OperatorRef() string {
	if vj.operatorRef == "" {
		vj.operatorRef = vj.findStringChildContent("OperatorRef")
	}
	return vj.operatorRef
}

func (vj *XMLMonitoredVehicleJourney) Coordinates() string {
	if vj.coordinates == "" {
		vj.coordinates = vj.findStringChildContent("Coordinates")
	}
	return vj.coordinates
}

func (vj *XMLMonitoredVehicleJourney) Longitude() string {
	if vj.longitude == "" {
		vj.longitude = vj.findStringChildContent("Longitude")
	}
	return vj.longitude
}

func (vj *XMLMonitoredVehicleJourney) Latitude() string {
	if vj.latitude == "" {
		vj.latitude = vj.findStringChildContent("Latitude")
	}
	return vj.latitude
}

func (vj *XMLMonitoredVehicleJourney) VehicleRef() string {
	if vj.vehicleRef == "" {
		vref := vj.findStringChildContent("VehicleRef")
		if vref != "" {
			vj.vehicleRef = vref
		} else {
			vj.vehicleRef = vj.findStringChildContent("VehicleMonitoringRef")
		}
	}
	return vj.vehicleRef
}

func (vj *XMLMonitoredVehicleJourney) DriverRef() string {
	if vj.driverRef == "" {
		vj.driverRef = vj.findStringChildContent("DriverRef")
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
