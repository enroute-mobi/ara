package siri

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
	XMLStructure

	itemIdentifier         string
	monitoringRef          string
	stopPointRef           string
	stopPointName          string
	datedVehicleJourneyRef string
	lineRef                string
	vehicleJourneyName     string
	publishedLineName      string
	departureStatus        string
	arrivalStatus          string
	dataFrameRef           string
	recordedAt             time.Time
	order                  int

	aimedArrivalTime    time.Time
	expectedArrivalTime time.Time
	actualArrivalTime   time.Time

	aimedDepartureTime    time.Time
	expectedDepartureTime time.Time
	actualDepartureTime   time.Time

	// Attributes

	situationRef                string
	inCongestion                string
	delay                       string
	vehicleAtStop               Bool
	actualQuayName              string
	aimedHeadwayInterval        string
	arrivalPlatformName         string
	arrivalProximyTest          string
	departureBoardingActivity   string
	departurePlatformName       string
	destinationDisplay          string
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

func (visit *XMLMonitoredStopVisit) ItemIdentifier() string {
	if visit.itemIdentifier == "" {
		visit.itemIdentifier = visit.findStringChildContent("ItemIdentifier")
	}
	return visit.itemIdentifier
}

func (visit *XMLMonitoredStopVisit) MonitoringRef() string {
	if visit.monitoringRef == "" {
		visit.monitoringRef = visit.findStringChildContent("MonitoringRef")
	}
	return visit.monitoringRef
}

func (visit *XMLMonitoredStopVisit) StopPointRef() string {
	if visit.stopPointRef == "" {
		visit.stopPointRef = visit.findStringChildContent("StopPointRef")
	}
	return visit.stopPointRef
}

func (visit *XMLMonitoredStopVisit) StopPointName() string {
	if visit.stopPointName == "" {
		visit.stopPointName = visit.findStringChildContent("StopPointName")
	}
	return visit.stopPointName
}

func (visit *XMLMonitoredStopVisit) DatedVehicleJourneyRef() string {
	if visit.datedVehicleJourneyRef == "" {
		visit.datedVehicleJourneyRef = visit.findStringChildContent("DatedVehicleJourneyRef")
	}
	return visit.datedVehicleJourneyRef
}

func (visit *XMLMonitoredStopVisit) DataFrameRef() string {
	if visit.dataFrameRef == "" {
		visit.dataFrameRef = visit.findStringChildContent("DataFrameRef")
	}
	return visit.dataFrameRef
}

func (visit *XMLMonitoredStopVisit) LineRef() string {
	if visit.lineRef == "" {
		visit.lineRef = visit.findStringChildContent("LineRef")
	}
	return visit.lineRef
}

func (visit *XMLMonitoredStopVisit) PublishedLineName() string {
	if visit.publishedLineName == "" {
		visit.publishedLineName = visit.findStringChildContent("PublishedLineName")
	}
	return visit.publishedLineName
}

func (visit *XMLMonitoredStopVisit) DepartureStatus() string {
	if visit.departureStatus == "" {
		visit.departureStatus = visit.findStringChildContent("DepartureStatus")
	}
	return visit.departureStatus
}

func (visit *XMLMonitoredStopVisit) ArrivalStatus() string {
	if visit.arrivalStatus == "" {
		visit.arrivalStatus = visit.findStringChildContent("ArrivalStatus")
	}
	return visit.arrivalStatus
}

func (visit *XMLMonitoredStopVisit) RecordedAt() time.Time {
	if visit.recordedAt.IsZero() {
		visit.recordedAt = visit.findTimeChildContent("RecordedAtTime")
	}
	return visit.recordedAt
}

func (visit *XMLMonitoredStopVisit) Order() int {
	if visit.order == 0 {
		visit.order = visit.findIntChildContent("Order")
	}
	return visit.order
}

func (visit *XMLMonitoredStopVisit) AimedArrivalTime() time.Time {
	if visit.aimedArrivalTime.IsZero() {
		visit.aimedArrivalTime = visit.findTimeChildContent("AimedArrivalTime")
	}
	return visit.aimedArrivalTime
}

func (visit *XMLMonitoredStopVisit) ExpectedArrivalTime() time.Time {
	if visit.expectedArrivalTime.IsZero() {
		visit.expectedArrivalTime = visit.findTimeChildContent("ExpectedArrivalTime")
	}
	return visit.expectedArrivalTime
}

func (visit *XMLMonitoredStopVisit) ActualArrivalTime() time.Time {
	if visit.actualArrivalTime.IsZero() {
		visit.actualArrivalTime = visit.findTimeChildContent("ActualArrivalTime")
	}
	return visit.actualArrivalTime
}

func (visit *XMLMonitoredStopVisit) AimedDepartureTime() time.Time {
	if visit.aimedDepartureTime.IsZero() {
		visit.aimedDepartureTime = visit.findTimeChildContent("AimedDepartureTime")
	}
	return visit.aimedDepartureTime
}

func (visit *XMLMonitoredStopVisit) ExpectedDepartureTime() time.Time {
	if visit.expectedDepartureTime.IsZero() {
		visit.expectedDepartureTime = visit.findTimeChildContent("ExpectedDepartureTime")
	}
	return visit.expectedDepartureTime
}

func (visit *XMLMonitoredStopVisit) ActualDepartureTime() time.Time {
	if visit.actualDepartureTime.IsZero() {
		visit.actualDepartureTime = visit.findTimeChildContent("ActualDepartureTime")
	}
	return visit.actualDepartureTime
}

// Attributes
func (visit *XMLMonitoredStopVisit) Delay() string {
	if visit.delay == "" {
		visit.delay = visit.findStringChildContent("Delay")
	}
	return visit.delay
}

func (visit *XMLMonitoredStopVisit) ActualQuayName() string {
	if visit.actualQuayName == "" {
		visit.actualQuayName = visit.findStringChildContent("ActualQuayName")
	}
	return visit.actualQuayName
}

func (visit *XMLMonitoredStopVisit) AimedHeadwayInterval() string {
	if visit.aimedHeadwayInterval == "" {
		visit.aimedHeadwayInterval = visit.findStringChildContent("AimedHeadwayInterval")
	}
	return visit.aimedHeadwayInterval
}

func (visit *XMLMonitoredStopVisit) ArrivalPlatformName() string {
	if visit.arrivalPlatformName == "" {
		visit.arrivalPlatformName = visit.findStringChildContent("ArrivalPlatformName")
	}
	return visit.arrivalPlatformName
}

func (visit *XMLMonitoredStopVisit) ArrivalProximyTest() string {
	if visit.arrivalProximyTest == "" {
		visit.arrivalProximyTest = visit.findStringChildContent("ArrivalProximyTest")
	}
	return visit.arrivalProximyTest
}

func (visit *XMLMonitoredStopVisit) DepartureBoardingActivity() string {
	if visit.departureBoardingActivity == "" {
		visit.departureBoardingActivity = visit.findStringChildContent("DepartureBoardingActivity")
	}
	return visit.departureBoardingActivity
}

func (visit *XMLMonitoredStopVisit) DeparturePlatformName() string {
	if visit.departurePlatformName == "" {
		visit.departurePlatformName = visit.findStringChildContent("DeparturePlatformName")
	}
	return visit.departurePlatformName
}

func (visit *XMLMonitoredStopVisit) DestinationDisplay() string {
	if visit.destinationDisplay == "" {
		visit.destinationDisplay = visit.findStringChildContent("DestinationDisplay")
	}
	return visit.destinationDisplay
}

func (visit *XMLMonitoredStopVisit) DistanceFromStop() string {
	if visit.distanceFromStop == "" {
		visit.distanceFromStop = visit.findStringChildContent("DistanceFromStop")
	}
	return visit.distanceFromStop
}

func (visit *XMLMonitoredStopVisit) ExpectedHeadwayInterval() string {
	if visit.expectedHeadwayInterval == "" {
		visit.expectedHeadwayInterval = visit.findStringChildContent("ExpectedHeadwayInterval")
	}
	return visit.expectedHeadwayInterval
}

func (visit *XMLMonitoredStopVisit) NumberOfStopsAway() string {
	if visit.numberOfStopsAway == "" {
		visit.numberOfStopsAway = visit.findStringChildContent("NumberOfStopsAway")
	}
	return visit.numberOfStopsAway
}

func (visit *XMLMonitoredStopVisit) PlatformTraversal() string {
	if visit.platformTraversal == "" {
		visit.platformTraversal = visit.findStringChildContent("PlatformTraversal")
	}
	return visit.platformTraversal
}

func (visit *XMLMonitoredStopVisit) DirectionName() string {
	if visit.directionName == "" {
		visit.directionName = visit.findStringChildContent("DirectionName")
	}
	return visit.directionName
}

func (visit *XMLMonitoredStopVisit) DestinationName() string {
	if visit.destinationName == "" {
		visit.destinationName = visit.findStringChildContent("DestinationName")
	}
	return visit.destinationName
}

func (visit *XMLMonitoredStopVisit) DirectionRef() string {
	if visit.directionRef == "" {
		visit.directionRef = visit.findStringChildContent("DirectionRef")
	}
	return visit.directionRef
}

func (visit *XMLMonitoredStopVisit) Bearing() string {
	if visit.bearing == "" {
		visit.bearing = visit.findStringChildContent("Bearing")
	}
	return visit.bearing
}

func (visit *XMLMonitoredStopVisit) InPanic() string {
	if visit.inPanic == "" {
		visit.inPanic = visit.findStringChildContent("InPanic")
	}
	return visit.inPanic
}

func (visit *XMLMonitoredStopVisit) SituationRef() string {
	if visit.situationRef == "" {
		visit.situationRef = visit.findStringChildContent("SituationRef")
	}
	return visit.situationRef
}

func (visit *XMLMonitoredStopVisit) InCongestion() string {
	if visit.inCongestion == "" {
		visit.inCongestion = visit.findStringChildContent("InCongestion")
	}
	return visit.inPanic
}

func (visit *XMLMonitoredStopVisit) HeadwayService() string {
	if visit.headwayService == "" {
		visit.headwayService = visit.findStringChildContent("HeadwayService")
	}
	return visit.headwayService
}

func (visit *XMLMonitoredStopVisit) FirstOrLastJourney() string {
	if visit.firstOrLastJourney == "" {
		visit.firstOrLastJourney = visit.findStringChildContent("FirstOrLastJourney")
	}
	return visit.firstOrLastJourney
}

func (visit *XMLMonitoredStopVisit) JourneyNote() string {
	if visit.journeyNote == "" {
		visit.journeyNote = visit.findStringChildContent("JourneyNote")
	}
	return visit.journeyNote
}

func (visit *XMLMonitoredStopVisit) JourneyPatternName() string {
	if visit.journeyPatternName == "" {
		visit.journeyPatternName = visit.findStringChildContent("JourneyPatternName")
	}
	return visit.journeyPatternName
}

func (visit *XMLMonitoredStopVisit) VehicleAtStop() bool {
	if !visit.vehicleAtStop.Defined {
		visit.vehicleAtStop.Parse(visit.findStringChildContent("VehicleAtStop"))
	}
	return visit.vehicleAtStop.Value
}

func (visit *XMLMonitoredStopVisit) Monitored() bool {
	if !visit.monitored.Defined {
		visit.monitored.Parse(visit.findStringChildContent("Monitored"))
	}
	return visit.monitored.Value
}

func (visit *XMLMonitoredStopVisit) MonitoringError() string {
	if visit.monitoringError == "" {
		visit.monitoringError = visit.findStringChildContent("MonitoringError")
	}
	return visit.monitoringError
}

func (visit *XMLMonitoredStopVisit) Occupancy() string {
	if visit.occupancy == "" {
		visit.occupancy = visit.findStringChildContent("Occupancy")
	}
	return visit.occupancy
}

func (visit *XMLMonitoredStopVisit) OriginAimedDepartureTime() string {
	if visit.originAimedDepartureTime == "" {
		visit.originAimedDepartureTime = visit.findStringChildContent("OriginAimedDepartureTime")
	}
	return visit.originAimedDepartureTime
}

func (visit *XMLMonitoredStopVisit) DestinationAimedArrivalTime() string {
	if visit.destinationAimedArrivalTime == "" {
		visit.destinationAimedArrivalTime = visit.findStringChildContent("DestinationAimedArrivalTime")
	}
	return visit.destinationAimedArrivalTime
}

func (visit *XMLMonitoredStopVisit) OriginName() string {
	if visit.originName == "" {
		visit.originName = visit.findStringChildContent("OriginName")
	}
	return visit.originName
}

func (visit *XMLMonitoredStopVisit) ProductCategoryRef() string {
	if visit.productCategoryRef == "" {
		visit.productCategoryRef = visit.findStringChildContent("ProductCategoryRef")
	}
	return visit.productCategoryRef
}

func (visit *XMLMonitoredStopVisit) ServiceFeatureRef() string {
	if visit.serviceFeatureRef == "" {
		visit.serviceFeatureRef = visit.findStringChildContent("ServiceFeatureRef")
	}
	return visit.serviceFeatureRef
}

func (visit *XMLMonitoredStopVisit) TrainNumberRef() string {
	if visit.trainNumberRef == "" {
		visit.trainNumberRef = visit.findStringChildContent("TrainNumberRef")
	}
	return visit.trainNumberRef
}

func (visit *XMLMonitoredStopVisit) VehicleFeature() string {
	if visit.vehicleFeature == "" {
		visit.vehicleFeature = visit.findStringChildContent("VehicleFeature")
	}
	return visit.vehicleFeature
}

func (visit *XMLMonitoredStopVisit) VehicleJourneyName() string {
	if visit.vehicleJourneyName == "" {
		visit.vehicleJourneyName = visit.findStringChildContent("VehicleJourneyName")
	}
	return visit.vehicleJourneyName
}

func (visit *XMLMonitoredStopVisit) VehicleMode() string {
	if visit.vehicleMode == "" {
		visit.vehicleMode = visit.findStringChildContent("VehicleMode")
	}
	return visit.vehicleMode
}

func (visit *XMLMonitoredStopVisit) ViaPlaceName() string {
	if visit.viaPlaceName == "" {
		visit.viaPlaceName = visit.findStringChildContent("PlaceName")
	}
	return visit.viaPlaceName
}

// References

func (visit *XMLMonitoredStopVisit) OriginRef() string {
	if visit.originRef == "" {
		visit.originRef = visit.findStringChildContent("OriginRef")
	}
	return visit.originRef
}

func (visit *XMLMonitoredStopVisit) PlaceRef() string {
	if visit.placeRef == "" {
		visit.placeRef = visit.findStringChildContent("PlaceRef")
	}
	return visit.placeRef
}

func (visit *XMLMonitoredStopVisit) DestinationRef() string {
	if visit.destinationRef == "" {
		visit.destinationRef = visit.findStringChildContent("DestinationRef")
	}
	return visit.destinationRef
}

func (visit *XMLMonitoredStopVisit) JourneyPatternRef() string {
	if visit.journeyPatternRef == "" {
		visit.journeyPatternRef = visit.findStringChildContent("JourneyPatternRef")
	}
	return visit.journeyPatternRef
}

func (visit *XMLMonitoredStopVisit) RouteRef() string {
	if visit.routeRef == "" {
		visit.routeRef = visit.findStringChildContent("RouteRef")
	}
	return visit.routeRef
}

func (visit *XMLMonitoredStopVisit) OperatorRef() string {
	if visit.operatorRef == "" {
		visit.operatorRef = visit.findStringChildContent("OperatorRef")
	}
	return visit.operatorRef
}
