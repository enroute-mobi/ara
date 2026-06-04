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

	monitoringRef *string

	monitoredStopVisits             []*XMLMonitoredStopVisit
	monitoredStopVisitCancellations []*XMLMonitoredStopVisitCancellation
}

type XMLMonitoredStopVisitCancellation struct {
	XMLStructure

	itemRef       *string
	monitoringRef *string
	recordedAt    *time.Time
}

type XMLMonitoredStopVisit struct {
	XMLMonitoredVehicleJourney

	stopVisitRawAttributes map[string]string

	itemIdentifier *string
	monitoringRef  *string
	recordedAt     *time.Time
}

type XMLMonitoredVehicleJourney struct {
	XMLCall

	rawAttributes map[string]string

	datedVehicleJourneyRef *string
	lineRef                *string
	vehicleJourneyName     *string
	publishedLineName      *string
	dataFrameRef           *string

	// Attributes
	situationRef                *string
	inCongestion                *string
	delay                       *string
	actualQuayName              *string
	aimedHeadwayInterval        *string
	arrivalPlatformName         *string
	arrivalProximityText        *string
	departureBoardingActivity   *string
	departurePlatformName       *string
	distanceFromStop            *string
	expectedHeadwayInterval     *string
	numberOfStopsAway           *string
	platformTraversal           *string
	directionName               *string
	destinationName             *string
	directionRef                *string
	firstOrLastJourney          *string
	headwayService              *string
	journeyNote                 *string
	journeyPatternName          *string
	monitored                   Bool
	monitoringError             *string
	occupancy                   *string
	originAimedDepartureTime    *string
	destinationAimedArrivalTime *string
	originName                  *string
	productCategoryRef          *string
	serviceFeatureRef           *string
	trainNumberRef              *string
	vehicleFeatureRef           *string
	vehicleMode                 *string
	operatorRef                 *string
	viaPlaceName                *string
	originRef                   *string
	placeRef                    *string
	destinationRef              *string
	journeyPatternRef           *string
	routeRef                    *string
	bearing                     *string
	inPanic                     *string

	// VehicleMonitoring attributes
	srsName     *string
	coordinates *string
	longitude   *string
	latitude    *string
	vehicleRef  *string
	driverRef   *string
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
	if delivery.monitoringRef == nil {
		s := delivery.findStringChildContent(siri_attributes.MonitoringRef)
		delivery.monitoringRef = &s
	}
	return *delivery.monitoringRef
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
	if cancel.itemRef == nil {
		s := cancel.findStringChildContent(siri_attributes.ItemRef)
		cancel.itemRef = &s
	}
	return *cancel.itemRef
}

func (cancel *XMLMonitoredStopVisitCancellation) MonitoringRef() string {
	if cancel.monitoringRef == nil {
		s := cancel.findStringChildContent(siri_attributes.MonitoringRef)
		cancel.monitoringRef = &s
	}
	return *cancel.monitoringRef
}

func (cancel *XMLMonitoredStopVisitCancellation) RecordedAt() time.Time {
	if cancel.recordedAt == nil {
		t := cancel.findTimeChildContent(siri_attributes.RecordedAtTime)
		cancel.recordedAt = &t
	}
	return *cancel.recordedAt
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
	if sv.itemIdentifier == nil {
		s := sv.findStringChildContent(siri_attributes.ItemIdentifier)
		sv.itemIdentifier = &s
	}
	return *sv.itemIdentifier
}

func (sv *XMLMonitoredStopVisit) MonitoringRef() string {
	if sv.monitoringRef == nil {
		s := sv.findStringChildContent(siri_attributes.MonitoringRef)
		sv.monitoringRef = &s
	}
	return *sv.monitoringRef
}

func (sv *XMLMonitoredStopVisit) RecordedAt() time.Time {
	if sv.recordedAt == nil {
		t := sv.findTimeChildContent(siri_attributes.RecordedAtTime)
		sv.recordedAt = &t
	}
	return *sv.recordedAt
}

func (vj *XMLMonitoredVehicleJourney) DatedVehicleJourneyRef() string {
	if vj.datedVehicleJourneyRef == nil {
		s := vj.findStringChildContent(siri_attributes.DatedVehicleJourneyRef)
		vj.datedVehicleJourneyRef = &s
	}
	return *vj.datedVehicleJourneyRef
}

func (vj *XMLMonitoredVehicleJourney) DataFrameRef() string {
	if vj.dataFrameRef == nil {
		s := vj.findStringChildContent(siri_attributes.DataFrameRef)
		vj.dataFrameRef = &s
	}
	return *vj.dataFrameRef
}

func (vj *XMLMonitoredVehicleJourney) LineRef() string {
	if vj.lineRef == nil {
		s := vj.findStringChildContent(siri_attributes.LineRef)
		vj.lineRef = &s
	}
	return *vj.lineRef
}

func (vj *XMLMonitoredVehicleJourney) PublishedLineName() string {
	if vj.publishedLineName == nil {
		s := vj.findStringChildContent(siri_attributes.PublishedLineName)
		vj.publishedLineName = &s
	}
	return *vj.publishedLineName
}

// Attributes
func (vj *XMLMonitoredVehicleJourney) Delay() string {
	if vj.delay == nil {
		s := vj.findStringChildContent(siri_attributes.Delay)
		vj.delay = &s
	}
	return *vj.delay
}

func (vj *XMLMonitoredVehicleJourney) ActualQuayName() string {
	if vj.actualQuayName == nil {
		s := vj.findStringChildContent(siri_attributes.ActualQuayName)
		vj.actualQuayName = &s
	}
	return *vj.actualQuayName
}

func (vj *XMLMonitoredVehicleJourney) AimedHeadwayInterval() string {
	if vj.aimedHeadwayInterval == nil {
		s := vj.findStringChildContent(siri_attributes.AimedHeadwayInterval)
		vj.aimedHeadwayInterval = &s
	}
	return *vj.aimedHeadwayInterval
}

func (vj *XMLMonitoredVehicleJourney) ArrivalPlatformName() string {
	if vj.arrivalPlatformName == nil {
		s := vj.findStringChildContent(siri_attributes.ArrivalPlatformName)
		vj.arrivalPlatformName = &s
	}
	return *vj.arrivalPlatformName
}

func (vj *XMLMonitoredVehicleJourney) ArrivalProximityText() string {
	if vj.arrivalProximityText == nil {
		s := vj.findStringChildContent(siri_attributes.ArrivalProximityText)
		vj.arrivalProximityText = &s
	}
	return *vj.arrivalProximityText
}

func (vj *XMLMonitoredVehicleJourney) DepartureBoardingActivity() string {
	if vj.departureBoardingActivity == nil {
		s := vj.findStringChildContent(siri_attributes.DepartureBoardingActivity)
		vj.departureBoardingActivity = &s
	}
	return *vj.departureBoardingActivity
}

func (vj *XMLMonitoredVehicleJourney) DeparturePlatformName() string {
	if vj.departurePlatformName == nil {
		s := vj.findStringChildContent(siri_attributes.DeparturePlatformName)
		vj.departurePlatformName = &s
	}
	return *vj.departurePlatformName
}

func (vj *XMLMonitoredVehicleJourney) DistanceFromStop() string {
	if vj.distanceFromStop == nil {
		s := vj.findStringChildContent(siri_attributes.DistanceFromStop)
		vj.distanceFromStop = &s
	}
	return *vj.distanceFromStop
}

func (vj *XMLMonitoredVehicleJourney) ExpectedHeadwayInterval() string {
	if vj.expectedHeadwayInterval == nil {
		s := vj.findStringChildContent(siri_attributes.ExpectedHeadwayInterval)
		vj.expectedHeadwayInterval = &s
	}
	return *vj.expectedHeadwayInterval
}

func (vj *XMLMonitoredVehicleJourney) NumberOfStopsAway() string {
	if vj.numberOfStopsAway == nil {
		s := vj.findStringChildContent(siri_attributes.NumberOfStopsAway)
		vj.numberOfStopsAway = &s
	}
	return *vj.numberOfStopsAway
}

func (vj *XMLMonitoredVehicleJourney) PlatformTraversal() string {
	if vj.platformTraversal == nil {
		s := vj.findStringChildContent(siri_attributes.PlatformTraversal)
		vj.platformTraversal = &s
	}
	return *vj.platformTraversal
}

func (vj *XMLMonitoredVehicleJourney) DirectionName() string {
	if vj.directionName == nil {
		s := vj.findStringChildContent(siri_attributes.DirectionName)
		vj.directionName = &s
	}
	return *vj.directionName
}

func (vj *XMLMonitoredVehicleJourney) DestinationName() string {
	if vj.destinationName == nil {
		s := vj.findStringChildContent(siri_attributes.DestinationName)
		vj.destinationName = &s
	}
	return *vj.destinationName
}

func (vj *XMLMonitoredVehicleJourney) DirectionRef() string {
	if vj.directionRef == nil {
		s := vj.findStringChildContent(siri_attributes.DirectionRef)
		vj.directionRef = &s
	}
	return *vj.directionRef
}

func (vj *XMLMonitoredVehicleJourney) Bearing() string {
	if vj.bearing == nil {
		s := vj.findStringChildContent(siri_attributes.Bearing)
		vj.bearing = &s
	}
	return *vj.bearing
}

func (vj *XMLMonitoredVehicleJourney) InPanic() string {
	if vj.inPanic == nil {
		s := vj.findStringChildContent(siri_attributes.InPanic)
		vj.inPanic = &s
	}
	return *vj.inPanic
}

func (vj *XMLMonitoredVehicleJourney) SituationRef() string {
	if vj.situationRef == nil {
		s := vj.findStringChildContent(siri_attributes.SituationRef)
		vj.situationRef = &s
	}
	return *vj.situationRef
}

func (vj *XMLMonitoredVehicleJourney) InCongestion() string {
	if vj.inCongestion == nil {
		s := vj.findStringChildContent(siri_attributes.InCongestion)
		vj.inCongestion = &s
	}
	return *vj.inPanic
}

func (vj *XMLMonitoredVehicleJourney) HeadwayService() string {
	if vj.headwayService == nil {
		s := vj.findStringChildContent(siri_attributes.HeadwayService)
		vj.headwayService = &s
	}
	return *vj.headwayService
}

func (vj *XMLMonitoredVehicleJourney) FirstOrLastJourney() string {
	if vj.firstOrLastJourney == nil {
		s := vj.findStringChildContent(siri_attributes.FirstOrLastJourney)
		vj.firstOrLastJourney = &s
	}
	return *vj.firstOrLastJourney
}

func (vj *XMLMonitoredVehicleJourney) JourneyNote() string {
	if vj.journeyNote == nil {
		s := vj.findStringChildContent(siri_attributes.JourneyNote)
		vj.journeyNote = &s
	}
	return *vj.journeyNote
}

func (vj *XMLMonitoredVehicleJourney) JourneyPatternName() string {
	if vj.journeyPatternName == nil {
		s := vj.findStringChildContent(siri_attributes.JourneyPatternName)
		vj.journeyPatternName = &s
	}
	return *vj.journeyPatternName
}

func (vj *XMLMonitoredVehicleJourney) Monitored() bool {
	if !vj.monitored.Defined {
		vj.monitored.SetValue(vj.findBoolChildContent(siri_attributes.Monitored))
	}
	return vj.monitored.Value
}

func (vj *XMLMonitoredVehicleJourney) MonitoringError() string {
	if vj.monitoringError == nil {
		s := vj.findStringChildContent(siri_attributes.MonitoringError)
		vj.monitoringError = &s
	}
	return *vj.monitoringError
}

func (vj *XMLMonitoredVehicleJourney) Occupancy() string {
	if vj.occupancy == nil {
		s := vj.findStringChildContent(siri_attributes.Occupancy)
		vj.occupancy = &s
	}
	return *vj.occupancy
}

func (vj *XMLMonitoredVehicleJourney) OriginAimedDepartureTime() string {
	if vj.originAimedDepartureTime == nil {
		s := vj.findStringChildContent(siri_attributes.OriginAimedDepartureTime)
		vj.originAimedDepartureTime = &s
	}
	return *vj.originAimedDepartureTime
}

func (vj *XMLMonitoredVehicleJourney) DestinationAimedArrivalTime() string {
	if vj.destinationAimedArrivalTime == nil {
		s := vj.findStringChildContent(siri_attributes.DestinationAimedArrivalTime)
		vj.destinationAimedArrivalTime = &s
	}
	return *vj.destinationAimedArrivalTime
}

func (vj *XMLMonitoredVehicleJourney) OriginName() string {
	if vj.originName == nil {
		s := vj.findStringChildContent(siri_attributes.OriginName)
		vj.originName = &s
	}
	return *vj.originName
}

func (vj *XMLMonitoredVehicleJourney) ProductCategoryRef() string {
	if vj.productCategoryRef == nil {
		s := vj.findStringChildContent(siri_attributes.ProductCategoryRef)
		vj.productCategoryRef = &s
	}
	return *vj.productCategoryRef
}

func (vj *XMLMonitoredVehicleJourney) ServiceFeatureRef() string {
	if vj.serviceFeatureRef == nil {
		s := vj.findStringChildContent(siri_attributes.ServiceFeatureRef)
		vj.serviceFeatureRef = &s
	}
	return *vj.serviceFeatureRef
}

func (vj *XMLMonitoredVehicleJourney) TrainNumberRef() string {
	if vj.trainNumberRef == nil {
		s := vj.findStringChildContent(siri_attributes.TrainNumberRef)
		vj.trainNumberRef = &s
	}
	return *vj.trainNumberRef
}

func (vj *XMLMonitoredVehicleJourney) VehicleFeatureRef() string {
	if vj.vehicleFeatureRef == nil {
		s := vj.findStringChildContent(siri_attributes.VehicleFeatureRef)
		vj.vehicleFeatureRef = &s
	}
	return *vj.vehicleFeatureRef
}

func (vj *XMLMonitoredVehicleJourney) VehicleJourneyName() string {
	if vj.vehicleJourneyName == nil {
		s := vj.findStringChildContent(siri_attributes.VehicleJourneyName)
		vj.vehicleJourneyName = &s
	}
	return *vj.vehicleJourneyName
}

func (vj *XMLMonitoredVehicleJourney) VehicleMode() string {
	if vj.vehicleMode == nil {
		s := vj.findStringChildContent(siri_attributes.VehicleMode)
		vj.vehicleMode = &s
	}
	return *vj.vehicleMode
}

func (vj *XMLMonitoredVehicleJourney) ViaPlaceName() string {
	if vj.viaPlaceName == nil {
		s := vj.findStringChildContent(siri_attributes.PlaceName)
		vj.viaPlaceName = &s
	}
	return *vj.viaPlaceName
}

// References

func (vj *XMLMonitoredVehicleJourney) OriginRef() string {
	if vj.originRef == nil {
		s := vj.findStringChildContent(siri_attributes.OriginRef)
		vj.originRef = &s
	}
	return *vj.originRef
}

func (vj *XMLMonitoredVehicleJourney) PlaceRef() string {
	if vj.placeRef == nil {
		s := vj.findStringChildContent(siri_attributes.PlaceRef)
		vj.placeRef = &s
	}
	return *vj.placeRef
}

func (vj *XMLMonitoredVehicleJourney) DestinationRef() string {
	if vj.destinationRef == nil {
		s := vj.findStringChildContent(siri_attributes.DestinationRef)
		vj.destinationRef = &s
	}
	return *vj.destinationRef
}

func (vj *XMLMonitoredVehicleJourney) JourneyPatternRef() string {
	if vj.journeyPatternRef == nil {
		s := vj.findStringChildContent(siri_attributes.JourneyPatternRef)
		vj.journeyPatternRef = &s
	}
	return *vj.journeyPatternRef
}

func (vj *XMLMonitoredVehicleJourney) RouteRef() string {
	if vj.routeRef == nil {
		s := vj.findStringChildContent(siri_attributes.RouteRef)
		vj.routeRef = &s
	}
	return *vj.routeRef
}

func (vj *XMLMonitoredVehicleJourney) OperatorRef() string {
	if vj.operatorRef == nil {
		s := vj.findStringChildContent(siri_attributes.OperatorRef)
		vj.operatorRef = &s
	}
	return *vj.operatorRef
}

func (vj *XMLMonitoredVehicleJourney) Coordinates() string {
	if vj.coordinates == nil {
		s := vj.findStringChildContent(siri_attributes.Coordinates)
		vj.coordinates = &s
	}
	return *vj.coordinates
}

func (vj *XMLMonitoredVehicleJourney) Longitude() string {
	if vj.longitude == nil {
		s := vj.findStringChildContent(siri_attributes.Longitude)
		vj.longitude = &s
	}
	return *vj.longitude
}

func (vj *XMLMonitoredVehicleJourney) Latitude() string {
	if vj.latitude == nil {
		s := vj.findStringChildContent(siri_attributes.Latitude)
		vj.latitude = &s
	}
	return *vj.latitude
}

func (vj *XMLMonitoredVehicleJourney) VehicleRef() string {
	if vj.vehicleRef == nil {
		vref := vj.findStringChildContent(siri_attributes.VehicleRef)
		if vref != "" {
			vj.vehicleRef = &vref
		} else {
			s := vj.findStringChildContent(siri_attributes.VehicleMonitoringRef)
			vj.vehicleRef = &s
		}
	}
	return *vj.vehicleRef
}

func (vj *XMLMonitoredVehicleJourney) DriverRef() string {
	if vj.driverRef == nil {
		s := vj.findStringChildContent(siri_attributes.DriverRef)
		vj.driverRef = &s
	}
	return *vj.driverRef
}

func (vj *XMLMonitoredVehicleJourney) SRSName() string {
	if vj.srsName == nil {
		s := vj.findChildAttribute("VehicleLocation", "srsName")
		vj.srsName = &s
	}
	return *vj.srsName
}

func (vj *XMLMonitoredVehicleJourney) RawAttributes() map[string]string {
	if vj.rawAttributes != nil {
		return vj.rawAttributes
	}
	attrs := make(map[string]string)
	if v := vj.Delay(); v != "" {
		attrs[siri_attributes.Delay] = v
	}
	if v := vj.Bearing(); v != "" {
		attrs[siri_attributes.Bearing] = v
	}
	if v := vj.InPanic(); v != "" {
		attrs[siri_attributes.InPanic] = v
	}
	if v := vj.InCongestion(); v != "" {
		attrs[siri_attributes.InCongestion] = v
	}
	if v := vj.SituationRef(); v != "" {
		attrs[siri_attributes.SituationRef] = v
	}
	if v := vj.DirectionName(); v != "" {
		attrs[siri_attributes.DirectionName] = v
	}
	if v := vj.FirstOrLastJourney(); v != "" {
		attrs[siri_attributes.FirstOrLastJourney] = v
	}
	if v := vj.HeadwayService(); v != "" {
		attrs[siri_attributes.HeadwayService] = v
	}
	if v := vj.JourneyNote(); v != "" {
		attrs[siri_attributes.JourneyNote] = v
	}
	if v := vj.JourneyPatternName(); v != "" {
		attrs[siri_attributes.JourneyPatternName] = v
	}
	if v := vj.MonitoringError(); v != "" {
		attrs[siri_attributes.MonitoringError] = v
	}
	if v := vj.OriginAimedDepartureTime(); v != "" {
		attrs[siri_attributes.OriginAimedDepartureTime] = v
	}
	if v := vj.DestinationAimedArrivalTime(); v != "" {
		attrs[siri_attributes.DestinationAimedArrivalTime] = v
	}
	if v := vj.ProductCategoryRef(); v != "" {
		attrs[siri_attributes.ProductCategoryRef] = v
	}
	if v := vj.ServiceFeatureRef(); v != "" {
		attrs[siri_attributes.ServiceFeatureRef] = v
	}
	if v := vj.TrainNumberRef(); v != "" {
		attrs[siri_attributes.TrainNumberRef] = v
	}
	if v := vj.VehicleFeatureRef(); v != "" {
		attrs[siri_attributes.VehicleFeatureRef] = v
	}
	if v := vj.VehicleMode(); v != "" {
		attrs[siri_attributes.VehicleMode] = v
	}
	if v := vj.ViaPlaceName(); v != "" {
		attrs[siri_attributes.ViaPlaceName] = v
	}
	if v := vj.VehicleJourneyName(); v != "" {
		attrs[siri_attributes.VehicleJourneyName] = v
	}
	vj.rawAttributes = attrs
	return attrs
}

func (sv *XMLMonitoredStopVisit) RawAttributes() map[string]string {
	if sv.stopVisitRawAttributes != nil {
		return sv.stopVisitRawAttributes
	}
	attrs := make(map[string]string)
	if v := sv.Delay(); v != "" {
		attrs[siri_attributes.Delay] = v
	}
	if v := sv.ActualQuayName(); v != "" {
		attrs[siri_attributes.ActualQuayName] = v
	}
	if v := sv.AimedHeadwayInterval(); v != "" {
		attrs[siri_attributes.AimedHeadwayInterval] = v
	}
	if v := sv.ArrivalPlatformName(); v != "" {
		attrs[siri_attributes.ArrivalPlatformName] = v
	}
	if v := sv.ArrivalProximityText(); v != "" {
		attrs[siri_attributes.ArrivalProximityText] = v
	}
	if v := sv.DepartureBoardingActivity(); v != "" {
		attrs[siri_attributes.DepartureBoardingActivity] = v
	}
	if v := sv.DeparturePlatformName(); v != "" {
		attrs[siri_attributes.DeparturePlatformName] = v
	}
	if v := sv.DestinationDisplay(); v != "" {
		attrs[siri_attributes.DestinationDisplay] = v
	}
	if v := sv.DistanceFromStop(); v != "" {
		attrs[siri_attributes.DistanceFromStop] = v
	}
	if v := sv.ExpectedHeadwayInterval(); v != "" {
		attrs[siri_attributes.ExpectedHeadwayInterval] = v
	}
	if v := sv.NumberOfStopsAway(); v != "" {
		attrs[siri_attributes.NumberOfStopsAway] = v
	}
	if v := sv.PlatformTraversal(); v != "" {
		attrs[siri_attributes.PlatformTraversal] = v
	}
	sv.stopVisitRawAttributes = attrs
	return attrs
}

// Test methods

func (vj *XMLMonitoredVehicleJourney) SetLongitude(s string) {
	vj.longitude = &s
}

func (vj *XMLMonitoredVehicleJourney) SetLatitude(s string) {
	vj.latitude = &s
}

func (vj *XMLMonitoredVehicleJourney) SetSRSName(s string) {
	vj.srsName = &s
}

func (vj *XMLMonitoredVehicleJourney) SetCoordinates(s string) {
	vj.coordinates = &s
}
