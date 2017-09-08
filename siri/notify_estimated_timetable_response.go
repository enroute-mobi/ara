package siri

import (
	"bytes"
	"html/template"
	"time"
)

type SIRINotifyEstimatedTimeTable struct {
	ProducerRef               string
	ResponseMessageIdentifier string
	ResponseTimestamp         time.Time
	Deliveries                []*SIRIEstimatedTimetableSubscriptionDelivery
}

type SIRIEstimatedTimetableSubscriptionDelivery struct {
	ResponseTimestamp      time.Time
	SubscriberRef          string
	SubscriptionIdentifier string

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	EstimatedJourneyVersionFrames []*SIRIEstimatedJourneyVersionFrame
}

const estimatedTimeTablenotifyTemplate = `<sw:NotifyEstimatedTimetable xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>{{.ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:ResponseTimestamp>
		<siri:ProducerRef>{{.ProducerRef}}</siri:ProducerRef>
		<siri:ResponseMessageIdentifier>{{.ResponseMessageIdentifier}}</siri:ResponseMessageIdentifier>
	</ServiceDeliveryInfo>
	<Notification>{{range .Deliveries}}
		<siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>{{.ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</siri:ResponseTimestamp>
			<siri:SubscriberRef>{{.SubscriberRef}}</siri:SubscriberRef>
			<siri:SubscriptionRef>{{.SubscriptionIdentifier}}</siri:SubscriptionRef>
			<siri:Status>{{ .Status }}</siri:Status>{{ if not .Status }}
						<siri:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
							<siri:OtherError number="{{.ErrorNumber}}">{{ else }}
							<siri:{{.ErrorType}}>{{ end }}
								<siri:ErrorText>{{.ErrorText}}</siri:ErrorText>
							</siri:{{.ErrorType}}>
						</siri:ErrorCondition>{{ else }}{{ range .EstimatedJourneyVersionFrames }}
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>2017-01-01T12:00:00.000Z</siri:RecordedAtTime>{{ range .EstimatedVehicleJourneys }}
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>{{ .LineRef }}</siri:LineRef>{{ if .Attributes.DirectionRef }}
					<siri:DirectionRef>{{ .Attributes.DirectionRef }}</siri:DirectionRef>{{ end }}
					<siri:DatedVehicleJourneyRef>{{ .DatedVehicleJourneyRef }}</siri:DatedVehicleJourneyRef>{{ if .References.OriginRef }}
					<siri:OriginRef>{{ .References.OriginRef.ObjectId.Value }}</siri:OriginRef>{{ end }}{{ if .References.DestinationRef }}
					<siri:DestinationRef>{{ .References.DestinationRef.ObjectId.Value }}</siri:DestinationRef>{{ end }}
					<siri:EstimatedCalls>{{ range .EstimatedCalls }}
						<siri:EstimatedCall>
							<siri:StopPointRef>{{ .StopPointRef }}</siri:StopPointRef>
							<siri:Order>{{ .Order }}</siri:Order>{{ if .StopPointName }}
							<siri:StopPointName>{{ .StopPointName }}</siri:StopPointName>{{ end }}
							<siri:VehicleAtStop>{{ .VehicleAtStop }}</siri:VehicleAtStop>{{ if .DestinationDisplay }}
							<siri:DestinationDisplay>{{ .DestinationDisplay }}</siri:DestinationDisplay>{{ end }}{{ if not .AimedArrivalTime.IsZero }}
							<siri:AimedArrivalTime>{{ .AimedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:AimedArrivalTime>{{ end }}{{ if not .ExpectedArrivalTime.IsZero }}
							<siri:ExpectedArrivalTime>{{ .ExpectedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ExpectedArrivalTime>{{ end }}{{ if .ArrivalStatus }}
							<siri:ArrivalStatus>{{ .ArrivalStatus }}</siri:ArrivalStatus>{{end}}{{ if not .AimedDepartureTime.IsZero }}
							<siri:AimedDepartureTime>{{ .AimedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:AimedDepartureTime>{{ end }}{{ if not .ExpectedDepartureTime.IsZero }}
							<siri:ExpectedDepartureTime>{{ .ExpectedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ExpectedDepartureTime>{{ end }}{{ if .DepartureStatus }}
							<siri:DepartureStatus>{{ .DepartureStatus }}</siri:DepartureStatus>{{ end }}
						</siri:EstimatedCall>{{ end }}
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>{{ end }}
			</siri:EstimatedJourneyVersionFrame>{{ end }}{{ end }}
		</siri:EstimatedTimetableDelivery>{{end}}
	</Notification>
</sw:NotifyEstimatedTimetable>`

func (notify *SIRINotifyEstimatedTimeTable) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var notifyDelivery = template.Must(template.New("estimatedTimeTableNotify").Parse(estimatedTimeTablenotifyTemplate))
	if err := notifyDelivery.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
