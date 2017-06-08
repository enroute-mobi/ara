package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/af83/edwig/model"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringResponse struct {
	ResponseXMLStructure

	monitoredStopVisits []*XMLMonitoredStopVisit
}

type XMLMonitoredStopVisit struct {
	XMLStructure

	// parent *XMLStopMonitoringResponse

	itemIdentifier         string
	stopPointRef           string
	stopPointName          string
	datedVehicleJourneyRef string
	lineRef                string
	vehicleJourneyName     string
	publishedLineName      string
	departureStatus        string
	arrivalStatus          string
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
	monitored                   string
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

type SIRIStopMonitoringResponse struct {
	SIRIStopMonitoringDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
}

type SIRIStopMonitoringDelivery struct {
	RequestMessageRef string
	Status            bool
	ErrorType         string
	ErrorNumber       int
	ErrorText         string
	ResponseTimestamp time.Time

	MonitoredStopVisits []*SIRIMonitoredStopVisit
}

type SIRIMonitoredStopVisit struct {
	ItemIdentifier         string
	StopPointRef           string
	StopPointName          string
	DatedVehicleJourneyRef string
	LineRef                string
	PublishedLineName      string
	DepartureStatus        string
	ArrivalStatus          string
	VehicleJourneyName     string
	StopAreaObjectId       string

	VehicleAtStop bool

	Order int

	AimedArrivalTime    time.Time
	ExpectedArrivalTime time.Time
	ActualArrivalTime   time.Time

	DataFrameRef          string
	RecordedAt            time.Time
	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
	ActualDepartureTime   time.Time

	// Attributes
	Attributes map[string]map[string]string

	// Références
	References map[string]map[string]model.Reference
}

const stopMonitoringDeliveryTemplate = `<ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
			<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
			<ns3:Status>{{ .Status }}</ns3:Status>{{ if not .Status }}
			<ns3:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<ns3:OtherError number="{{.ErrorNumber}}">{{ else }}
				<ns3:{{.ErrorType}}>
					<ns3:ErrorText>{{.ErrorText}}</ns3:ErrorText>
				</ns3:{{.ErrorType}}>{{ end }}
			</ns3:ErrorCondition>{{ else }}{{ range .MonitoredStopVisits }}
			<ns3:MonitoredStopVisit>
				<ns3:RecordedAtTime>{{ .RecordedAt.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:RecordedAtTime>
				<ns3:ItemIdentifier>{{ .ItemIdentifier }}</ns3:ItemIdentifier>
				<ns3:MonitoringRef>{{ .StopPointRef }}</ns3:MonitoringRef>
				<ns3:MonitoredVehicleJourney>{{ if .LineRef }}
					<ns3:LineRef>{{ .LineRef }}</ns3:LineRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DirectionRef }}
					<ns3:DirectionRef>{{ .Attributes.VehicleJourneyAttributes.DirectionRef }}</ns3:DirectionRef>{{ end }}{{ if or .DatedVehicleJourneyRef .DataFrameRef }}
					<ns3:FramedVehicleJourneyRef>{{ if .DataFrameRef }}
						<ns3:DataFrameRef>{{ .DataFrameRef }}</ns3:DataFrameRef>{{ end }}{{ if .DatedVehicleJourneyRef }}
						<ns3:DatedVehicleJourneyRef>{{ .DatedVehicleJourneyRef }}</ns3:DatedVehicleJourneyRef>{{ end }}
					</ns3:FramedVehicleJourneyRef>{{ end }}{{ if .References.VehicleJourney.JourneyPatternRef }}
					<ns3:JourneyPatternRef>{{.References.VehicleJourney.JourneyPatternRef.ObjectId.Value}}</ns3:JourneyPatternRef>{{ end }}{{ if .Attributes.VehicleJourneyAttributes.JourneyPatternName }}
					<ns3:JourneyPatternName>{{ .Attributes.VehicleJourneyAttributes.JourneyPatternName }}</ns3:JourneyPatternName>{{ end }}{{ if .Attributes.VehicleJourneyAttributes.VehicleMode }}
					<ns3:VehicleMode>{{ .Attributes.VehicleJourneyAttributes.VehicleMode }}</ns3:VehicleMode>{{ end }}{{ if .PublishedLineName }}
					<ns3:PublishedLineName>{{ .PublishedLineName }}</ns3:PublishedLineName>{{ end }}{{ if .References.VehicleJourney.RouteRef}}
					<ns3:RouteRef>{{.References.VehicleJourney.RouteRef.ObjectId.Value }}</ns3:RouteRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DirectionName}}
					<ns3:DirectionName>{{.Attributes.VehicleJourneyAttributes.DirectionName}}</ns3:DirectionName>{{end}}{{ if .References.StopVisitReferences.OperatorRef}}
					<ns3:OperatorRef>{{.References.StopVisitReferences.OperatorRef.ObjectId.Value}}</ns3:OperatorRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.ProductCategoryRef}}
					<ns3:ProductCategoryRef>{{.Attributes.VehicleJourneyAttributes.ProductCategoryRef}}</ns3:ProductCategoryRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.ServiceFeatureRef}}
					<ns3:ServiceFeatureRef>{{.Attributes.VehicleJourneyAttributes.ServiceFeatureRef}}</ns3:ServiceFeatureRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.VehicleFeatureRef}}
					<ns3:VehicleFeatureRef>{{.Attributes.VehicleJourneyAttributes.VehicleFeatureRef}}</ns3:VehicleFeatureRef>{{end}}{{ if .References.VehicleJourney.OriginRef.ObjectId.Value}}
					<ns3:OriginRef>{{ .References.VehicleJourney.OriginRef.ObjectId.Value }}</ns3:OriginRef>{{ end }}{{ if .Attributes.VehicleJourneyAttributes.OriginName }}
					<ns3:OriginName>{{ .Attributes.VehicleJourneyAttributes.OriginName }}</ns3:OriginName>{{ end }}{{ if or .Attributes.VehicleJourneyAttributes.ViaPlaceName .References.VehicleJourney.PlaceRef }}
					<ns3:Via>{{ if .Attributes.VehicleJourneyAttributes.ViaPlaceName }}
						<ns3:PlaceName>{{ .Attributes.VehicleJourneyAttributes.ViaPlaceName }}</ns3:PlaceName>{{end}}{{ if .References.VehicleJourney.PlaceRef}}
					  <ns3:PlaceRef>{{.References.VehicleJourney.PlaceRef.ObjectId.Value}}</ns3:PlaceRef>{{ end }}
					</ns3:Via>{{ end }}{{ if .References.VehicleJourney.DestinationRef.ObjectId.Value }}
					<ns3:DestinationRef>{{ .References.VehicleJourney.DestinationRef.ObjectId.Value }}</ns3:DestinationRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DestinationName}}
					<ns3:DestinationName>{{ .Attributes.VehicleJourneyAttributes.DestinationName }}</ns3:DestinationName>{{end}}{{ if .VehicleJourneyName }}
					<ns3:VehicleJourneyName>{{ .VehicleJourneyName }}</ns3:VehicleJourneyName>{{end}}{{ if .Attributes.VehicleJourneyAttributes.JourneyNote}}
					<ns3:JourneyNote>{{.Attributes.VehicleJourneyAttributes.JourneyNote}}</ns3:JourneyNote>{{end}}{{ if .Attributes.VehicleJourneyAttributes.HeadwayService}}
					<ns3:HeadwayService>{{.Attributes.VehicleJourneyAttributes.HeadwayService}}</ns3:HeadwayService>{{end}}{{ if .Attributes.VehicleJourneyAttributes.OriginAimedDepartureTime}}
					<ns3:OriginAimedDepartureTime>{{.Attributes.VehicleJourneyAttributes.OriginAimedDepartureTime}}</ns3:OriginAimedDepartureTime>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DestinationAimedArrivalTime}}
					<ns3:DestinationAimedArrivalTime>{{.Attributes.VehicleJourneyAttributes.DestinationAimedArrivalTime}}</ns3:DestinationAimedArrivalTime>{{end}}{{ if .Attributes.VehicleJourneyAttributes.FirstOrLastJourney}}
					<ns3:FirstOrLastJourney>{{.Attributes.VehicleJourneyAttributes.FirstOrLastJourney}}</ns3:FirstOrLastJourney>{{end}}{{ if .Attributes.VehicleJourneyAttributes.Monitored }}
					<ns3:Monitored>{{.Attributes.VehicleJourneyAttributes.Monitored}}</ns3:Monitored>{{end}}{{ if .Attributes.VehicleJourneyAttributes.MonitoringError}}
					<ns3:MonitoringError>{{.Attributes.VehicleJourneyAttributes.MonitoringError}}</ns3:MonitoringError>{{end}}{{ if .Attributes.VehicleJourneyAttributes.Occupancy }}
					<ns3:Occupancy>{{.Attributes.VehicleJourneyAttributes.Occupancy}}</ns3:Occupancy>{{end}}{{if .Attributes.VehicleJourneyAttributes.Delay}}
					<ns3:Delay>{{.Attributes.VehicleJourneyAttributes.Delay}}</ns3:Delay>{{end}}{{if .Attributes.VehicleJourneyAttributes.Bearing}}
					<ns3:Bearing>{{.Attributes.VehicleJourneyAttributes.Bearing}}</ns3:Bearing>{{end}}{{ if .Attributes.VehicleJourneyAttributes.InPanic }}
					<ns3:InPanic>{{.Attributes.VehicleJourneyAttributes.InPanic}}</ns3:InPanic>{{end}}{{ if .Attributes.VehicleJourneyAttributes.InCongestion }}
					<ns3:InCongestion>{{.Attributes.VehicleJourneyAttributes.InCongestion}}</ns3:InCongestion>{{end}}{{ if .Attributes.VehicleJourneyAttributes.TrainNumberRef }}
					<ns3:TrainNumber>
						<ns3:TrainNumberRef>{{ .Attributes.VehicleJourneyAttributes.TrainNumberRef }}</ns3:TrainNumberRef>
					</ns3:TrainNumber>{{ end }}{{ if .Attributes.VehicleJourneyAttributes.SituationRef }}
					<ns3:SituationRef>{{.Attributes.VehicleJourneyAttributes.SituationRef}}</ns3:SituationRef>{{end}}
					<ns3:MonitoredCall>{{if .StopPointRef}}
						<ns3:StopPointRef>{{ .StopPointRef }}</ns3:StopPointRef>{{end}}{{if .Order}}
						<ns3:Order>{{ .Order }}</ns3:Order>{{end}}{{if .StopPointName}}
						<ns3:StopPointName>{{ .StopPointName }}</ns3:StopPointName>{{end}}
						<ns3:VehicleAtStop>{{ .VehicleAtStop }}</ns3:VehicleAtStop>{{if .Attributes.StopVisitAttributes.PlatformTraversal }}
						<ns3:PlatformTraversal>{{ .Attributes.StopVisitAttributes.PlatformTraversal }}</ns3:PlatformTraversal>{{end}}{{if .Attributes.StopVisitAttributes.DestinationDisplay }}
						<ns3:DestinationDisplay>{{ .Attributes.StopVisitAttributes.DestinationDisplay }}</ns3:DestinationDisplay>{{end}}{{ if not .AimedArrivalTime.IsZero }}
						<ns3:AimedArrivalTime>{{ .AimedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:AimedArrivalTime>{{ end }}{{ if not .ActualArrivalTime.IsZero }}
						<ns3:ActualArrivalTime>{{ .ActualArrivalTime.Format "2006-01-02T15:04:05.000Z07:00"}}</ns3:ActualArrivalTime>{{ end }}{{ if not .ExpectedArrivalTime.IsZero }}
						<ns3:ExpectedArrivalTime>{{ .ExpectedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ExpectedArrivalTime>{{ end }}{{ if .ArrivalStatus }}
						<ns3:ArrivalStatus>{{ .ArrivalStatus }}</ns3:ArrivalStatus>{{end}}{{if .Attributes.StopVisitAttributes.ArrivalProximyTest }}
						<ns3:ArrivalProximyTest>{{ .Attributes.StopVisitAttributes.ArrivalProximyTest }}</ns3:ArrivalProximyTest>{{end}}{{if .Attributes.StopVisitAttributes.ArrivalPlatformName }}
						<ns3:ArrivalPlatformName>{{ .Attributes.StopVisitAttributes.ArrivalPlatformName }}</ns3:ArrivalPlatformName>{{end}}{{ if .Attributes.StopVisitAttributes.ActualQuayName }}
						<ns3:ArrivalStopAssignment>
							<ns3:ActualQuayName>{{ .Attributes.StopVisitAttributes.ActualQuayName }}</ns3:ActualQuayName>
						</ns3:ArrivalStopAssignment>{{ end }}{{ if not .AimedDepartureTime.IsZero }}
						<ns3:AimedDepartureTime>{{ .AimedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:AimedDepartureTime>{{ end }}{{ if not .ActualDepartureTime.IsZero }}
						<ns3:ActualDepartureTime>{{ .ActualDepartureTime.Format "2006-01-02T15:04:05.000Z07:00"}}</ns3:ActualDepartureTime>{{ end }}{{ if not .ExpectedDepartureTime.IsZero }}
						<ns3:ExpectedDepartureTime>{{ .ExpectedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ExpectedDepartureTime>{{ end }}{{ if .DepartureStatus }}
						<ns3:DepartureStatus>{{ .DepartureStatus }}</ns3:DepartureStatus>{{end}}{{ if .Attributes.StopVisitAttributes.DeparturePlatformName }}
						<ns3:DeparturePlatformName>{{ .Attributes.StopVisitAttributes.DeparturePlatformName }}</ns3:DeparturePlatformName>{{end}}{{ if .Attributes.StopVisitAttributes.DepartureBoardingActivity }}
						<ns3:DepartureBoardingActivity>{{ .Attributes.StopVisitAttributes.DepartureBoardingActivity }}</ns3:DepartureBoardingActivity>{{end}}{{ if .Attributes.StopVisitAttributes.AimedHeadwayInterval }}
						<ns3:AimedHeadwayInterval>{{ .Attributes.StopVisitAttributes.AimedHeadwayInterval }}</ns3:AimedHeadwayInterval>{{end}}{{ if .Attributes.StopVisitAttributes.ExpectedHeadwayInterval }}
						<ns3:ExpectedHeadwayInterval>{{ .Attributes.StopVisitAttributes.ExpectedHeadwayInterval }}</ns3:ExpectedHeadwayInterval>{{end}}{{ if .Attributes.StopVisitAttributes.DistanceFromStop }}
						<ns3:DistanceFromStop>{{ .Attributes.StopVisitAttributes.DistanceFromStop }}</ns3:DistanceFromStop>{{end}}{{ if .Attributes.StopVisitAttributes.NumberOfStopsAway }}
						<ns3:NumberOfStopsAway>{{ .Attributes.StopVisitAttributes.NumberOfStopsAway }}</ns3:NumberOfStopsAway>{{end}}
					</ns3:MonitoredCall>
				</ns3:MonitoredVehicleJourney>
			</ns3:MonitoredStopVisit>{{ end }}{{ end }}
		</ns3:StopMonitoringDelivery>`

const stopMonitoringResponseTemplate = `<ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
															 xmlns:ns4="http://www.ifopt.org.uk/acsb"
															 xmlns:ns5="http://www.ifopt.org.uk/ifopt"
															 xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
															 xmlns:ns7="http://scma/siri"
															 xmlns:ns8="http://wsdl.siri.org.uk"
															 xmlns:ns9="http://wsdl.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
		<ns3:ProducerRef>{{ .ProducerRef }}</ns3:ProducerRef>{{ if .Address }}
		<ns3:Address>{{ .Address }}</ns3:Address>{{ end }}
		<ns3:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		{{ .BuildStopMonitoringDeliveryXML }}
	</Answer>
	<AnswerExtension/>
</ns8:GetStopMonitoringResponse>`

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

func (response *XMLStopMonitoringResponse) XMLMonitoredStopVisits() []*XMLMonitoredStopVisit {
	if len(response.monitoredStopVisits) == 0 {
		nodes := response.findNodes("MonitoredStopVisit")
		if nodes == nil {
			return response.monitoredStopVisits
		}
		for _, stopVisitNode := range nodes {
			response.monitoredStopVisits = append(response.monitoredStopVisits, NewXMLMonitoredStopVisit(stopVisitNode))
		}
	}
	return response.monitoredStopVisits
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

func (visit *XMLMonitoredStopVisit) Monitored() string {
	if visit.monitored == "" {
		visit.monitored = visit.findStringChildContent("Monitored")
	}
	return visit.monitored
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

func (response *SIRIStopMonitoringResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(stopMonitoringResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIStopMonitoringDelivery) BuildStopMonitoringDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	var stopMonitoringDelivery = template.Must(template.New("stopMonitoringDelivery").Parse(stopMonitoringDeliveryTemplate))
	if err := stopMonitoringDelivery.Execute(&buffer, delivery); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
