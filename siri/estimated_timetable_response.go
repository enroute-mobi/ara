package siri

import (
	"bytes"
	"text/template"
	"time"

	"github.com/af83/edwig/model"
)

type SIRIEstimatedTimeTableResponse struct {
	Address                   string
	ProducerRef               string
	RequestMessageRef         string
	ResponseMessageIdentifier string

	ResponseTimestamp time.Time

	Status      bool
	ErrorType   string
	ErrorNumber int
	ErrorText   string

	EstimatedJourneyVersionFrames []SIRIEstimatedJourneyVersionFrame
}

type SIRIEstimatedJourneyVersionFrame struct {
	RecordedAtTime time.Time

	EstimatedVehicleJourneys []SIRIEstimatedVehicleJourney
}

type SIRIEstimatedVehicleJourney struct {
	LineRef                string
	PublishedLineName      string
	DatedVehicleJourneyRef string

	Attributes map[string]string
	References map[string]model.Reference

	EstimatedCalls []SIRIEstimatedCall
}

type SIRIEstimatedCall struct {
	ArrivalStatus   string
	DepartureStatus string
	StopPointRef    string

	Order int

	AimedArrivalTime    time.Time
	ExpectedArrivalTime time.Time
	ActualArrivalTime   time.Time

	AimedDepartureTime    time.Time
	ExpectedDepartureTime time.Time
	ActualDepartureTime   time.Time
}

const estimatedTimeTableResponseTemplate = `<ns8:GetEstimatedTimetableResponse xmlns:ns3="http://www.siri.org.uk/siri"
																	 xmlns:ns4="http://www.ifopt.org.uk/acsb"
																	 xmlns:ns5="http://www.ifopt.org.uk/ifopt"
																	 xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
																	 xmlns:ns7="http://scma/siri"
																	 xmlns:ns8="http://wsdl.siri.org.uk"
																	 xmlns:ns9="http://wsdl.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
		<ns3:ProducerRef>{{ .ProducerRef }}</ns3:ProducerRef>
		<ns3:Address>{{ .Address }}</ns3:Address>
		<ns3:ResponseMessageIdentifier>{{ .ResponseMessageIdentifier }}</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<ns3:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>{{ .ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ResponseTimestamp>
			<ns3:RequestMessageRef>{{ .RequestMessageRef }}</ns3:RequestMessageRef>
			<ns3:Status>{{ .Status }}</ns3:Status>{{ if not .Status }}
			<ns3:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<ns3:OtherError number="{{.ErrorNumber}}">{{ else }}
				<ns3:{{.ErrorType}}>
					<ns3:ErrorText>{{.ErrorText}}</ns3:ErrorText>
				</ns3:{{.ErrorType}}>{{ end }}
			</ns3:ErrorCondition>{{ else }}{{ range .EstimatedJourneyVersionFrames }}
			<ns3:EstimatedJourneyVersionFrame>
				<ns3:RecordedAtTime>{{ .RecordedAtTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:RecordedAtTime>{{ range .EstimatedVehicleJourneys }}
				<ns3:EstimatedVehicleJourney>
					<ns3:LineRef>{{ .LineRef }}</ns3:LineRef>{{ if .Attributes.DirectionRef }}
					<ns3:DirectionRef>{{ .Attributes.DirectionRef }}</ns3:DirectionRef>{{ end }}
					<ns3:DatedVehicleJourneyRef>{{ .DatedVehicleJourneyRef }}</ns3:DatedVehicleJourneyRef>
					<ns3:PublishedLineName>{{ .PublishedLineName }}</ns3:PublishedLineName>{{ if .References.OriginRef }}
					<ns3:OriginRef>{{ .References.OriginRef.ObjectId.Value }}</ns3:OriginRef>{{ end }}{{ if .Attributes.OriginName }}
					<ns3:OriginName>{{ .Attributes.OriginName }}</ns3:OriginName>{{ end }}{{ if .References.DestinationRef }}
					<ns3:DestinationRef>{{ .References.DestinationRef.ObjectId.Value }}</ns3:DestinationRef>{{ end }}{{ if .Attributes.DestinationName }}
					<ns3:DestinationName>{{ .Attributes.DestinationName }}</ns3:DestinationName>{{ end }}
					<ns3:EstimatedCalls>{{ range .EstimatedCalls }}
						<ns3:EstimatedCall>
							<ns3:StopPointRef>{{ .StopPointRef }}</ns3:StopPointRef>
							<ns3:Order>{{ .Order }}</ns3:Order>{{ if not .AimedArrivalTime.IsZero }}
							<ns3:AimedArrivalTime>{{ .AimedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:AimedArrivalTime>{{ end }}{{ if not .ActualArrivalTime.IsZero }}
							<ns3:ActualArrivalTime>{{ .ActualArrivalTime.Format "2006-01-02T15:04:05.000Z07:00"}}</ns3:ActualArrivalTime>{{ end }}{{ if not .ExpectedArrivalTime.IsZero }}
							<ns3:ExpectedArrivalTime>{{ .ExpectedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ExpectedArrivalTime>{{ end }}{{ if .ArrivalStatus }}
							<ns3:ArrivalStatus>{{ .ArrivalStatus }}</ns3:ArrivalStatus>{{end}}{{ if not .AimedDepartureTime.IsZero }}
							<ns3:AimedDepartureTime>{{ .AimedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:AimedDepartureTime>{{ end }}{{ if not .ActualDepartureTime.IsZero }}
							<ns3:ActualDepartureTime>{{ .ActualDepartureTime.Format "2006-01-02T15:04:05.000Z07:00"}}</ns3:ActualDepartureTime>{{ end }}{{ if not .ExpectedDepartureTime.IsZero }}
							<ns3:ExpectedDepartureTime>{{ .ExpectedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ExpectedDepartureTime>{{ end }}{{ if .DepartureStatus }}
							<ns3:DepartureStatus>{{ .DepartureStatus }}</ns3:DepartureStatus>{{end}}
						</ns3:EstimatedCall>{{ end }}
					</ns3:EstimatedCalls>
				</ns3:EstimatedVehicleJourney>{{ end }}
			</ns3:EstimatedJourneyVersionFrame>{{ end }}{{ end }}
		</ns3:EstimatedTimetableDelivery>
	</Answer>
</ns8:GetEstimatedTimetableResponse>`

func (response *SIRIEstimatedTimeTableResponse) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var siriResponse = template.Must(template.New("siriResponse").Parse(estimatedTimeTableResponseTemplate))
	if err := siriResponse.Execute(&buffer, response); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
