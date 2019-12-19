package siri

import (
	"bytes"
	"text/template"
	"time"

	"bitbucket.org/enroute-mobi/edwig/model"
)

type SIRIEstimatedTimeTableResponse struct {
	SIRIEstimatedTimetableDelivery

	Address                   string
	ProducerRef               string
	ResponseMessageIdentifier string
}

type SIRIEstimatedTimetableDelivery struct {
	RequestMessageRef string

	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	EstimatedJourneyVersionFrames []*SIRIEstimatedJourneyVersionFrame
}

type SIRIEstimatedJourneyVersionFrame struct {
	RecordedAtTime time.Time

	EstimatedVehicleJourneys []*SIRIEstimatedVehicleJourney
}

type SIRIEstimatedVehicleJourney struct {
	LineRef                string
	DatedVehicleJourneyRef string

	Attributes map[string]string
	References map[string]model.Reference

	EstimatedCalls []*SIRIEstimatedCall
}

type SIRIEstimatedCall struct {
	ArrivalStatus      string
	DepartureStatus    string
	StopPointRef       string
	StopPointName      string
	DestinationDisplay string

	VehicleAtStop bool

	Order int

	AimedArrivalTime    time.Time
	ExpectedArrivalTime time.Time

	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
}

const estimatedTimeTableResponseTemplate = `<sw:GetEstimatedTimetableResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
		<siri:ProducerRef>{{ .ProducerRef }}</siri:ProducerRef>{{ if .Address }}
		<siri:Address>{{ .Address }}</siri:Address>{{ end }}
		<siri:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		{{ .BuildEstimatedTimetableDeliveryXML }}
	</Answer>
	<AnswerExtension/>
</sw:GetEstimatedTimetableResponse>`

const estimatedTimetableDeliveryTemplate = `<siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:ResponseTimestamp>
			<siri:RequestMessageRef>{{ .RequestMessageRef }}</siri:RequestMessageRef>
			<siri:Status>{{ .Status }}</siri:Status>{{ if not .Status }}
			<siri:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<siri:OtherError number="{{.ErrorNumber}}">{{ else }}
				<siri:{{.ErrorType}}>{{ end }}
					<siri:ErrorText>{{.ErrorText}}</siri:ErrorText>
				</siri:{{.ErrorType}}>
			</siri:ErrorCondition>{{ end }}{{ if or .Status (eq .ErrorType "OtherError") }}{{ range .EstimatedJourneyVersionFrames }}
			{{ .BuildEstimatedJourneyVersionFrameXML }}{{ end }}{{ end }}
		</siri:EstimatedTimetableDelivery>`

const estimatedJourneyVersionFrameTemplate = `<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>{{ .RecordedAtTime.Format "2006-01-02T15:04:05.000Z07:00" }}</siri:RecordedAtTime>{{ range .EstimatedVehicleJourneys }}
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>{{ .LineRef }}</siri:LineRef>{{ if .Attributes.DirectionRef }}
					<siri:DirectionRef>{{ .Attributes.DirectionRef }}</siri:DirectionRef>{{ else }}
					<siri:DirectionRef/>{{ end }}{{ if .References.OperatorRef }}
					<siri:OperatorRef>{{ .References.OperatorRef.ObjectId.Value }}</siri:OperatorRef>{{ else }}
					<siri:OperatorRef/>{{ end }}
					<siri:DatedVehicleJourneyRef>{{ .DatedVehicleJourneyRef }}</siri:DatedVehicleJourneyRef>{{ if .References.OriginRef }}
					<siri:OriginRef>{{ .References.OriginRef.ObjectId.Value }}</siri:OriginRef>{{ end }}{{ if .References.DestinationRef }}
					<siri:DestinationRef>{{ .References.DestinationRef.ObjectId.Value }}</siri:DestinationRef>{{ end }}{{ if ne (len .EstimatedCalls) 0 }}
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
							<siri:DepartureStatus>{{ .DepartureStatus }}</siri:DepartureStatus>{{end}}
						</siri:EstimatedCall>{{ end }}
					</siri:EstimatedCalls>{{ end }}
				</siri:EstimatedVehicleJourney>{{ end }}
			</siri:EstimatedJourneyVersionFrame>`

func (response *SIRIEstimatedTimeTableResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(estimatedTimeTableResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (delivery *SIRIEstimatedTimetableDelivery) BuildEstimatedTimetableDeliveryXML() (string, error) {
	var buffer bytes.Buffer
	var estimatedTimetableDelivery = template.Must(template.New("estimatedTimetableDelivery").Parse(estimatedTimetableDeliveryTemplate))
	if err := estimatedTimetableDelivery.Execute(&buffer, delivery); err != nil {
		return "", err
	}
	return buffer.String(), nil
}

func (frame *SIRIEstimatedJourneyVersionFrame) BuildEstimatedJourneyVersionFrameXML() (string, error) {
	var buffer bytes.Buffer
	var estimatedJourneyVersionFrame = template.Must(template.New("estimatedJourneyVersionFrame").Parse(estimatedJourneyVersionFrameTemplate))
	if err := estimatedJourneyVersionFrame.Execute(&buffer, frame); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
