package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/af83/edwig/model"
)

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
	MonitoringRef     string

	MonitoredStopVisits []*SIRIMonitoredStopVisit
	CancelledStopVisits []*SIRICancelledStopVisit
}

type SIRICancelledStopVisit struct {
	RecordedAtTime         time.Time
	ItemRef                string
	MonitoringRef          string
	LineRef                string
	DataFrameRef           string
	DatedVehicleJourneyRef string
	PublishedLineName      string
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
	OriginName             string
	DestinationName        string
	StopAreaObjectId       string
	DataFrameRef           string

	VehicleAtStop bool
	Monitored     bool

	Order int

	AimedArrivalTime    time.Time
	ExpectedArrivalTime time.Time
	ActualArrivalTime   time.Time

	RecordedAt            time.Time
	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
	ActualDepartureTime   time.Time

	// Attributes
	Attributes map[string]map[string]string

	// Références
	References map[string]map[string]model.Reference
}

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

const stopMonitoringDeliveryTemplate = `<siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
			<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
			<siri:MonitoringRef>{{ .MonitoringRef }}</siri:MonitoringRef>
			<siri:Status>{{ .Status }}</siri:Status>{{ if not .Status }}
			<siri:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<siri:OtherError number="{{.ErrorNumber}}">{{ else }}
				<siri:{{.ErrorType}}>{{ end }}
					<siri:ErrorText>{{.ErrorText}}</siri:ErrorText>
				</siri:{{.ErrorType}}>
			</siri:ErrorCondition>{{ end }}{{ if or .Status (eq .ErrorType "OtherError") }}{{ range .MonitoredStopVisits }}
			{{ .BuildMonitoredStopVisitXML }}{{ end }}{{ end }}
		</siri:StopMonitoringDelivery>`

const monitoredStopVisitTemplate = `<siri:MonitoredStopVisit>
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
					<siri:OriginRef>{{ .References.VehicleJourney.OriginRef.ObjectId.Value }}</siri:OriginRef>{{ end }}{{ if .OriginName }}
					<siri:OriginName>{{ .OriginName }}</siri:OriginName>{{ end }}{{ if or .Attributes.VehicleJourneyAttributes.ViaPlaceName .References.VehicleJourney.PlaceRef }}
					<siri:Via>{{ if .Attributes.VehicleJourneyAttributes.ViaPlaceName }}
						<siri:PlaceName>{{ .Attributes.VehicleJourneyAttributes.ViaPlaceName }}</siri:PlaceName>{{end}}{{ if .References.VehicleJourney.PlaceRef}}
					  <siri:PlaceRef>{{.References.VehicleJourney.PlaceRef.ObjectId.Value}}</siri:PlaceRef>{{ end }}
					</siri:Via>{{ end }}{{ if .References.VehicleJourney.DestinationRef.ObjectId.Value }}
					<siri:DestinationRef>{{ .References.VehicleJourney.DestinationRef.ObjectId.Value }}</siri:DestinationRef>{{end}}{{ if .DestinationName}}
					<siri:DestinationName>{{ .DestinationName }}</siri:DestinationName>{{end}}{{ if .VehicleJourneyName }}
					<siri:VehicleJourneyName>{{ .VehicleJourneyName }}</siri:VehicleJourneyName>{{end}}{{ if .Attributes.VehicleJourneyAttributes.JourneyNote}}
					<siri:JourneyNote>{{.Attributes.VehicleJourneyAttributes.JourneyNote}}</siri:JourneyNote>{{end}}{{ if .Attributes.VehicleJourneyAttributes.HeadwayService}}
					<siri:HeadwayService>{{.Attributes.VehicleJourneyAttributes.HeadwayService}}</siri:HeadwayService>{{end}}{{ if .Attributes.VehicleJourneyAttributes.OriginAimedDepartureTime}}
					<siri:OriginAimedDepartureTime>{{.Attributes.VehicleJourneyAttributes.OriginAimedDepartureTime}}</siri:OriginAimedDepartureTime>{{end}}{{ if .Attributes.VehicleJourneyAttributes.DestinationAimedArrivalTime}}
					<siri:DestinationAimedArrivalTime>{{.Attributes.VehicleJourneyAttributes.DestinationAimedArrivalTime}}</siri:DestinationAimedArrivalTime>{{end}}{{ if .Attributes.VehicleJourneyAttributes.FirstOrLastJourney}}
					<siri:FirstOrLastJourney>{{.Attributes.VehicleJourneyAttributes.FirstOrLastJourney}}</siri:FirstOrLastJourney>{{end}}
					<siri:Monitored>{{.Monitored}}</siri:Monitored>{{ if .Attributes.VehicleJourneyAttributes.MonitoringError}}
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
			</siri:MonitoredStopVisit>`

const cancelledStopVisitTemplate = `<siri:MonitoredStopVisitCancellation>
			<siri:RecordedAtTime>{{ .RecordedAtTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:RecordedAtTime>
			<siri:ItemRef>{{ .ItemRef }}</siri:ItemRef>
			<siri:MonitoringRef>{{ .MonitoringRef }}</siri:MonitoringRef>
			<siri:LineRef>{{ .LineRef }}</siri:LineRef>
			<siri:PublishedLineName>{{.PublishedLineName}}</siri:PublishedLineName>
			<siri:VehicleJourneyRef>
			  <siri:DataFrameRef>{{.DataFrameRef}}</siri:DataFrameRef>
			  <siri:DatedVehicleJourneyRef>{{.DatedVehicleJourneyRef}}</siri:DatedVehicleJourneyRef>
			</siri:VehicleJourneyRef>
		</siri:MonitoredStopVisitCancellation>`

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

func (stopVisit *SIRIMonitoredStopVisit) BuildMonitoredStopVisitXML() (string, error) {
	var buffer bytes.Buffer
	var stopMonitoringDelivery = template.Must(template.New("monitoredStopVisit").Parse(monitoredStopVisitTemplate))
	if err := stopMonitoringDelivery.Execute(&buffer, stopVisit); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (stopVisit *SIRICancelledStopVisit) BuildCancelledStopVisitXML() (string, error) {
	var buffer bytes.Buffer
	var stopMonitoringDelivery = template.Must(template.New("cancelledStopVisit").Parse(cancelledStopVisitTemplate))
	if err := stopMonitoringDelivery.Execute(&buffer, stopVisit); err != nil {
		return "", err
	}

	return buffer.String(), nil
}
