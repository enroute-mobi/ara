package siri

import (
	"bytes"
	"text/template"
	"time"
)

type SIRINotifyStopMonitoring struct {
	Address                   string
	RequestMessageRef         string
	ProducerRef               string
	ResponseMessageIdentifier string
	SubscriberRef             string
	SubscriptionIdentifier    string

	ResponseTimestamp time.Time
	Status            bool
	ErrorType         string
	ErrorNumber       int
	ErrorText         string

	MonitoredStopVisits []*SIRIMonitoredStopVisit
}

const stopMonitoringNotifyTemplate = `<ns1:NotifyStopMonitoring xmlns:ns1="http://wsdl.siri.org.uk">
	<ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
		<ns5:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns5:ResponseTimestamp>
		<ns5:ProducerRef>{{ .ProducerRef }}</ns5:ProducerRef>{{ if .Address }}
		<ns5:Address>{{ .Address }}</ns5:Address>{{ end }}
		<ns5:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</ns5:ResponseMessageIdentifier>
		<ns5:RequestMessageRef>{{ .RequestMessageRef }}</ns5:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Notification xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
		<ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
			<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
			<ns5:SubscriberRef>{{.SubscriberRef}}</ns5:SubscriberRef>
			<ns5:SubscriptionRef>{{.SubscriptionIdentifier}}</ns5:SubscriptionRef>
			<ns3:Status>{{ .Status }}</ns3:Status>{{ if not .Status }}
			<ns3:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<ns3:OtherError number="{{.ErrorNumber}}">{{ else }}
				<ns3:{{.ErrorType}}>{{ end }}
					<ns3:ErrorText>{{.ErrorText}}</ns3:ErrorText>
				</ns3:{{.ErrorType}}>
			</ns3:ErrorCondition>{{ else }}{{ range .MonitoredStopVisits }}
			<ns3:MonitoredStopVisit>
				<ns3:RecordedAtTime>{{ .RecordedAt.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:RecordedAtTime>
				<ns3:ItemIdentifier>{{ .ItemIdentifier }}</ns3:ItemIdentifier>
				<ns3:MonitoringRef>{{ .MonitoringRef }}</ns3:MonitoringRef>
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
		</ns3:StopMonitoringDelivery>
		</Notification>
		<NotifyExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
	</ns1:NotifyStopMonitoring>`

func (notify *SIRINotifyStopMonitoring) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var notifyDelivery = template.Must(template.New("stopMonitoringNotify").Parse(stopMonitoringNotifyTemplate))
	if err := notifyDelivery.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
