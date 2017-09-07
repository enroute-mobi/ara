package siri

import (
	"bytes"
	"html/template"
	"time"
)

type SIRIEstimatedTimeTable struct {
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

const estimatedTimeTablenotifyTemplate = `<ns1:NotifyEstimatedTimetable xmlns:ns1="http://wsdl.siri.org.uk">
 <ServiceDeliveryInfo xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri">
   <ns5:ResponseTimestamp>{{.ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns5:ResponseTimestamp>
   <ns5:ProducerRef>{{.ProducerRef}}</ns5:ProducerRef>
   <ns5:ResponseMessageIdentifier>{{.ResponseMessageIdentifier}}</ns5:ResponseMessageIdentifier>
 </ServiceDeliveryInfo>
 <Notification> {{range .Deliveries}}
   <ns3:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
      <ns3:ResponseTimestamp>{{.ResponseTimestamp.Format "2006-01-02T15:04:05.000Z07:00"}}</ns3:ResponseTimestamp>
      <ns5:SubscriberRef>{{.SubscriberRef}}</ns5:SubscriberRef>
      <ns5:SubscriptionRef>{{.SubscriptionRef}}</ns5:SubscriptionRef>
      <ns3:Status>true</ns3:Status>>{{ if not .Status }}
			<ns3:ErrorCondition>{{ if eq .ErrorType "OtherError" }}
				<ns3:OtherError number="{{.ErrorNumber}}">{{ else }}
				<ns3:{{.ErrorType}}>{{ end }}
					<ns3:ErrorText>{{.ErrorText}}</ns3:ErrorText>
				</ns3:{{.ErrorType}}>
      </ns3:ErrorCondition>{{ else }}{{ range .EstimatedJourneyVersionFrames }}
      <ns3:EstimatedJourneyVersionFrame>
        <ns3:RecordedAtTime>2017-01-01T12:00:00.000Z</ns3:RecordedAtTime>
        <ns3:EstimatedVehicleJourney>{{ range .EstimatedVehicleJourneys }}
          <ns3:LineRef>{{ .LineRef }}</ns3:LineRef>{{ if .Attributes.DirectionRef }}
          <ns3:DirectionRef>{{ .Attributes.DirectionRef }}</ns3:DirectionRef>{{ end }}
          <ns3:DatedVehicleJourneyRef>{{ .DatedVehicleJourneyRef }}</ns3:DatedVehicleJourneyRef>{{ if .References.OriginRef }}
          <ns3:OriginRef>{{ .References.OriginRef.ObjectId.Value }}</ns3:OriginRef>{{ end }}{{ if .References.DestinationRef }}
          <ns3:DestinationRef>{{ .References.DestinationRef.ObjectId.Value }}</ns3:DestinationRef>{{ end }}
          <ns3:EstimatedCalls>{{ range .EstimatedCalls }}
            <ns3:EstimatedCall>
              <ns3:StopPointRef>{{ .StopPointRef }}</ns3:StopPointRef>
              <ns3:Order>{{ .Order }}</ns3:Order>{{ if .StopPointName }}
              <ns3:StopPointName>{{ .StopPointName }}</ns3:StopPointName>{{ end }}
              <ns3:VehicleAtStop>{{ .VehicleAtStop }}</ns3:VehicleAtStop>{{ if .DestinationDisplay }}
              <ns3:DestinationDisplay>{{ .DestinationDisplay }}</ns3:DestinationDisplay>{{ end }}{{ if not .AimedArrivalTime.IsZero }}
              <ns3:AimedArrivalTime>{{ .AimedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:AimedArrivalTime>{{ end }}{{ if not .ExpectedArrivalTime.IsZero }}
              <ns3:ExpectedArrivalTime>{{ .ExpectedArrivalTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ExpectedArrivalTime>{{ end }}{{ if .ArrivalStatus }}
              <ns3:ArrivalStatus>{{ .ArrivalStatus }}</ns3:ArrivalStatus>{{end}}{{ if not .AimedDepartureTime.IsZero }}
              <ns3:AimedDepartureTime>{{ .AimedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:AimedDepartureTime>{{ end }}{{ if not .ExpectedDepartureTime.IsZero }}
              <ns3:ExpectedDepartureTime>{{ .ExpectedDepartureTime.Format "2006-01-02T15:04:05.000Z07:00" }}</ns3:ExpectedDepartureTime>{{ end }}{{ if .DepartureStatus }}
              <ns3:DepartureStatus>{{ .DepartureStatus }}</ns3:DepartureStatus>{{end}}
            </ns3:EstimatedCall>{{ end }}
          </ns3:EstimatedCalls>
      </ns3:EstimatedVehicleJourney>{{ end }}
    </ns3:EstimatedJourneyVersionFrame>{{ end }}{{ end }}
</ns3:EstimatedTimetableDelivery>>
 </Notification>
<NotifyExtension xmlns:ns2="http://www.ifopt.org.uk/acsb" xmlns:ns3="http://www.ifopt.org.uk/ifopt" xmlns:ns4="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns5="http://www.siri.org.uk/siri" xmlns:ns6="http://wsdl.siri.org.uk/siri"/>
</ns1:NotifyEstimatedTimetable>`

func (notify *SIRIEstimatedTimeTable) BuildXML() (string, error) {
	var buffer bytes.Buffer
	var notifyDelivery = template.Must(template.New("estimatedTimeTableNotify").Parse(estimatedTimeTablenotifyTemplate))
	if err := notifyDelivery.Execute(&buffer, notify); err != nil {
		return "", err
	}
	return buffer.String(), nil
}
