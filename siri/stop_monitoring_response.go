package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/af83/edwig/model"
	"github.com/jbowtie/gokogiri"
	"github.com/jbowtie/gokogiri/xml"
)

type XMLStopMonitoringDelivery struct {
	XMLStructure

	subscriptionRef                 string
	subscriberRef                   string
	monitoredStopVisits             []*XMLMonitoredStopVisit
	monitoredStopVisitCancellations []*XMLMonitoredStopVisitCancellation
}

func NewXMLStopMonitoringDelivery(node XMLNode) *XMLStopMonitoringDelivery {
	delivery := &XMLStopMonitoringDelivery{}
	delivery.node = node
	return delivery
}

type XMLStopMonitoringResponse struct {
	ResponseXMLStructure

	// Can't include XMLStopMonitoringDelivery :(
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
	MonitoringRef          string
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

const stopMonitoringDeliveryTemplate = `<siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
			<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
			<siri:Status>{{ .Status }}</siri:Status>{{ if not .Status }}
			<siri:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<siri:OtherError number="{{.ErrorNumber}}">{{ else }}
				<siri:{{.ErrorType}}>{{ end }}
					<siri:ErrorText>{{.ErrorText}}</siri:ErrorText>
				</siri:{{.ErrorType}}>
			</siri:ErrorCondition>{{ else }}{{ range .MonitoredStopVisits }}
			<siri:MonitoredStopVisit>
				<siri:RecordedAtTime>{{ .RecordedAt.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:RecordedAtTime>
				<siri:ItemIdentifier>{{ .ItemIdentifier }}</siri:ItemIdentifier>
				<siri:MonitoringRef>{{ .MonitoringRef }}</siri:MonitoringRef>
				<siri:MonitoredVehicleJourney>{{ if .LineRef }}
					<siri:LineRef>{{ .LineRef }}</siri:LineRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DirectionRef }}
					<siri:DirectionRef>{{ .Attributes.VehicleJourneyAttributes.DirectionRef }}</siri:DirectionRef>{{ end }}{{ if or .DatedVehicleJourneyRef .DataFrameRef }}
					<siri:FramedVehicleJourneyRef>{{ if .DataFrameRef }}
						<siri:DataFrameRef>{{ .DataFrameRef }}</siri:DataFrameRef>{{ end }}{{ if .DatedVehicleJourneyRef }}
						<siri:DatedVehicleJourneyRef>{{ .DatedVehicleJourneyRef }}</siri:DatedVehicleJourneyRef>{{ end }}
					</siri:FramedVehicleJourneyRef>{{ end }}{{ if .References.VehicleJourney.JourneyPatternRef }}
					<siri:JourneyPatternRef>{{.References.VehicleJourney.JourneyPatternRef.ObjectId.Value}}</siri:JourneyPatternRef>{{ end }}{{ if .Attributes.VehicleJourneyAttributes.JourneyPatternName }}
					<siri:JourneyPatternName>{{ .Attributes.VehicleJourneyAttributes.JourneyPatternName }}</siri:JourneyPatternName>{{ end }}{{ if .Attributes.VehicleJourneyAttributes.VehicleMode }}
					<siri:VehicleMode>{{ .Attributes.VehicleJourneyAttributes.VehicleMode }}</siri:VehicleMode>{{ end }}{{ if .PublishedLineName }}
					<siri:PublishedLineName>{{ .PublishedLineName }}</siri:PublishedLineName>{{ end }}{{ if .References.VehicleJourney.RouteRef}}
					<siri:RouteRef>{{.References.VehicleJourney.RouteRef.ObjectId.Value }}</siri:RouteRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DirectionName}}
					<siri:DirectionName>{{.Attributes.VehicleJourneyAttributes.DirectionName}}</siri:DirectionName>{{end}}{{ if .References.StopVisitReferences.OperatorRef}}
					<siri:OperatorRef>{{.References.StopVisitReferences.OperatorRef.ObjectId.Value}}</siri:OperatorRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.ProductCategoryRef}}
					<siri:ProductCategoryRef>{{.Attributes.VehicleJourneyAttributes.ProductCategoryRef}}</siri:ProductCategoryRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.ServiceFeatureRef}}
					<siri:ServiceFeatureRef>{{.Attributes.VehicleJourneyAttributes.ServiceFeatureRef}}</siri:ServiceFeatureRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.VehicleFeatureRef}}
					<siri:VehicleFeatureRef>{{.Attributes.VehicleJourneyAttributes.VehicleFeatureRef}}</siri:VehicleFeatureRef>{{end}}{{ if .References.VehicleJourney.OriginRef.ObjectId.Value}}
					<siri:OriginRef>{{ .References.VehicleJourney.OriginRef.ObjectId.Value }}</siri:OriginRef>{{ end }}{{ if .Attributes.VehicleJourneyAttributes.OriginName }}
					<siri:OriginName>{{ .Attributes.VehicleJourneyAttributes.OriginName }}</siri:OriginName>{{ end }}{{ if or .Attributes.VehicleJourneyAttributes.ViaPlaceName .References.VehicleJourney.PlaceRef }}
					<siri:Via>{{ if .Attributes.VehicleJourneyAttributes.ViaPlaceName }}
						<siri:PlaceName>{{ .Attributes.VehicleJourneyAttributes.ViaPlaceName }}</siri:PlaceName>{{end}}{{ if .References.VehicleJourney.PlaceRef}}
					  <siri:PlaceRef>{{.References.VehicleJourney.PlaceRef.ObjectId.Value}}</siri:PlaceRef>{{ end }}
					</siri:Via>{{ end }}{{ if .References.VehicleJourney.DestinationRef.ObjectId.Value }}
					<siri:DestinationRef>{{ .References.VehicleJourney.DestinationRef.ObjectId.Value }}</siri:DestinationRef>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DestinationName}}
					<siri:DestinationName>{{ .Attributes.VehicleJourneyAttributes.DestinationName }}</siri:DestinationName>{{end}}{{ if .VehicleJourneyName }}
					<siri:VehicleJourneyName>{{ .VehicleJourneyName }}</siri:VehicleJourneyName>{{end}}{{ if .Attributes.VehicleJourneyAttributes.JourneyNote}}
					<siri:JourneyNote>{{.Attributes.VehicleJourneyAttributes.JourneyNote}}</siri:JourneyNote>{{end}}{{ if .Attributes.VehicleJourneyAttributes.HeadwayService}}
					<siri:HeadwayService>{{.Attributes.VehicleJourneyAttributes.HeadwayService}}</siri:HeadwayService>{{end}}{{ if .Attributes.VehicleJourneyAttributes.OriginAimedDepartureTime}}
					<siri:OriginAimedDepartureTime>{{.Attributes.VehicleJourneyAttributes.OriginAimedDepartureTime}}</siri:OriginAimedDepartureTime>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DestinationAimedArrivalTime}}
					<siri:DestinationAimedArrivalTime>{{.Attributes.VehicleJourneyAttributes.DestinationAimedArrivalTime}}</siri:DestinationAimedArrivalTime>{{end}}{{ if .Attributes.VehicleJourneyAttributes.FirstOrLastJourney}}
					<siri:FirstOrLastJourney>{{.Attributes.VehicleJourneyAttributes.FirstOrLastJourney}}</siri:FirstOrLastJourney>{{end}}{{ if .Attributes.VehicleJourneyAttributes.Monitored }}
					<siri:Monitored>{{.Attributes.VehicleJourneyAttributes.Monitored}}</siri:Monitored>{{end}}{{ if .Attributes.VehicleJourneyAttributes.MonitoringError}}
					<siri:MonitoringError>{{.Attributes.VehicleJourneyAttributes.MonitoringError}}</siri:MonitoringError>{{end}}{{ if .Attributes.VehicleJourneyAttributes.Occupancy }}
					<siri:Occupancy>{{.Attributes.VehicleJourneyAttributes.Occupancy}}</siri:Occupancy>{{end}}{{if .Attributes.VehicleJourneyAttributes.Delay}}
					<siri:Delay>{{.Attributes.VehicleJourneyAttributes.Delay}}</siri:Delay>{{end}}{{if .Attributes.VehicleJourneyAttributes.Bearing}}
					<siri:Bearing>{{.Attributes.VehicleJourneyAttributes.Bearing}}</siri:Bearing>{{end}}{{ if .Attributes.VehicleJourneyAttributes.InPanic }}
					<siri:InPanic>{{.Attributes.VehicleJourneyAttributes.InPanic}}</siri:InPanic>{{end}}{{ if .Attributes.VehicleJourneyAttributes.InCongestion }}
					<siri:InCongestion>{{.Attributes.VehicleJourneyAttributes.InCongestion}}</siri:InCongestion>{{end}}{{ if .Attributes.VehicleJourneyAttributes.TrainNumberRef }}
					<siri:TrainNumber>
						<siri:TrainNumberRef>{{ .Attributes.VehicleJourneyAttributes.TrainNumberRef }}</siri:TrainNumberRef>
					</siri:TrainNumber>{{ end }}{{ if .Attributes.VehicleJourneyAttributes.SituationRef }}
					<siri:SituationRef>{{.Attributes.VehicleJourneyAttributes.SituationRef}}</siri:SituationRef>{{end}}
					<siri:MonitoredCall>{{if .StopPointRef}}
						<siri:StopPointRef>{{ .StopPointRef }}</siri:StopPointRef>{{end}}{{if .Order}}
						<siri:Order>{{ .Order }}</siri:Order>{{end}}{{if .StopPointName}}
						<siri:StopPointName>{{ .StopPointName }}</siri:StopPointName>{{end}}
						<siri:VehicleAtStop>{{ .VehicleAtStop }}</siri:VehicleAtStop>{{if .Attributes.StopVisitAttributes.PlatformTraversal }}
						<siri:PlatformTraversal>{{ .Attributes.StopVisitAttributes.PlatformTraversal }}</siri:PlatformTraversal>{{end}}{{if .Attributes.StopVisitAttributes.DestinationDisplay }}
						<siri:DestinationDisplay>{{ .Attributes.StopVisitAttributes.DestinationDisplay }}</siri:DestinationDisplay>{{end}}{{ if not .AimedArrivalTime.IsZero }}
						<siri:AimedArrivalTime>{{ .AimedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:AimedArrivalTime>{{ end }}{{ if not .ActualArrivalTime.IsZero }}
						<siri:ActualArrivalTime>{{ .ActualArrivalTime.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:ActualArrivalTime>{{ end }}{{ if not .ExpectedArrivalTime.IsZero }}
						<siri:ExpectedArrivalTime>{{ .ExpectedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ExpectedArrivalTime>{{ end }}{{ if .ArrivalStatus }}
						<siri:ArrivalStatus>{{ .ArrivalStatus }}</siri:ArrivalStatus>{{end}}{{if .Attributes.StopVisitAttributes.ArrivalProximyTest }}
						<siri:ArrivalProximyTest>{{ .Attributes.StopVisitAttributes.ArrivalProximyTest }}</siri:ArrivalProximyTest>{{end}}{{if .Attributes.StopVisitAttributes.ArrivalPlatformName }}
						<siri:ArrivalPlatformName>{{ .Attributes.StopVisitAttributes.ArrivalPlatformName }}</siri:ArrivalPlatformName>{{end}}{{ if .Attributes.StopVisitAttributes.ActualQuayName }}
						<siri:ArrivalStopAssignment>
							<siri:ActualQuayName>{{ .Attributes.StopVisitAttributes.ActualQuayName }}</siri:ActualQuayName>
						</siri:ArrivalStopAssignment>{{ end }}{{ if not .AimedDepartureTime.IsZero }}
						<siri:AimedDepartureTime>{{ .AimedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:AimedDepartureTime>{{ end }}{{ if not .ActualDepartureTime.IsZero }}
						<siri:ActualDepartureTime>{{ .ActualDepartureTime.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:ActualDepartureTime>{{ end }}{{ if not .ExpectedDepartureTime.IsZero }}
						<siri:ExpectedDepartureTime>{{ .ExpectedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ExpectedDepartureTime>{{ end }}{{ if .DepartureStatus }}
						<siri:DepartureStatus>{{ .DepartureStatus }}</siri:DepartureStatus>{{end}}{{ if .Attributes.StopVisitAttributes.DeparturePlatformName }}
						<siri:DeparturePlatformName>{{ .Attributes.StopVisitAttributes.DeparturePlatformName }}</siri:DeparturePlatformName>{{end}}{{ if .Attributes.StopVisitAttributes.DepartureBoardingActivity }}
						<siri:DepartureBoardingActivity>{{ .Attributes.StopVisitAttributes.DepartureBoardingActivity }}</siri:DepartureBoardingActivity>{{end}}{{ if .Attributes.StopVisitAttributes.AimedHeadwayInterval }}
						<siri:AimedHeadwayInterval>{{ .Attributes.StopVisitAttributes.AimedHeadwayInterval }}</siri:AimedHeadwayInterval>{{end}}{{ if .Attributes.StopVisitAttributes.ExpectedHeadwayInterval }}
						<siri:ExpectedHeadwayInterval>{{ .Attributes.StopVisitAttributes.ExpectedHeadwayInterval }}</siri:ExpectedHeadwayInterval>{{end}}{{ if .Attributes.StopVisitAttributes.DistanceFromStop }}
						<siri:DistanceFromStop>{{ .Attributes.StopVisitAttributes.DistanceFromStop }}</siri:DistanceFromStop>{{end}}{{ if .Attributes.StopVisitAttributes.NumberOfStopsAway }}
						<siri:NumberOfStopsAway>{{ .Attributes.StopVisitAttributes.NumberOfStopsAway }}</siri:NumberOfStopsAway>{{end}}
					</siri:MonitoredCall>
				</siri:MonitoredVehicleJourney>
			</siri:MonitoredStopVisit>{{ end }}{{ end }}
		</siri:StopMonitoringDelivery>`

const stopMonitoringResponseTemplate = `<sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
		<siri:ProducerRef>{{ .ProducerRef }}</siri:ProducerRef>{{ if .Address }}
		<siri:Address>{{ .Address }}</siri:Address>{{ end }}
		<siri:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		{{ .BuildStopMonitoringDeliveryXML }}
	</Answer>
	<AnswerExtension/>
</sw:GetStopMonitoringResponse>`

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

func (response *XMLStopMonitoringResponse) XMLMonitoredStopVisitCancellations() []*XMLMonitoredStopVisitCancellation {
	if response.monitoredStopVisitCancellations == nil {
		cancellations := []*XMLMonitoredStopVisitCancellation{}
		nodes := response.findNodes("MonitoredStopVisitCancellation")
		if nodes != nil {
			for _, node := range nodes {
				cancellations = append(cancellations, NewXMLCancelledStopVisit(node))
			}
		}
		response.monitoredStopVisitCancellations = cancellations
	}
	return response.monitoredStopVisitCancellations
}

func (delivery *XMLStopMonitoringDelivery) XMLMonitoredStopVisitCancellations() []*XMLMonitoredStopVisitCancellation {
	if delivery.monitoredStopVisitCancellations == nil {
		cancellations := []*XMLMonitoredStopVisitCancellation{}
		nodes := delivery.findNodes("MonitoredStopVisitCancellation")
		if nodes != nil {
			for _, node := range nodes {
				cancellations = append(cancellations, NewXMLCancelledStopVisit(node))
			}
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

func (response *XMLStopMonitoringResponse) XMLMonitoredStopVisits() []*XMLMonitoredStopVisit {
	if response.monitoredStopVisits == nil {
		stopVisits := []*XMLMonitoredStopVisit{}
		nodes := response.findNodes("MonitoredStopVisit")
		if nodes != nil {
			for _, node := range nodes {
				stopVisits = append(stopVisits, NewXMLMonitoredStopVisit(node))
			}
		}
		response.monitoredStopVisits = stopVisits
	}
	return response.monitoredStopVisits
}

func (delivery *XMLStopMonitoringDelivery) XMLMonitoredStopVisits() []*XMLMonitoredStopVisit {
	if delivery.monitoredStopVisits == nil {
		stopVisits := []*XMLMonitoredStopVisit{}
		nodes := delivery.findNodes("MonitoredStopVisit")
		if nodes != nil {
			for _, node := range nodes {
				stopVisits = append(stopVisits, NewXMLMonitoredStopVisit(node))
			}
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
