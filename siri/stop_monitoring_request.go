package siri

import (
	"bytes"
	"html/template"
	"time"

	"github.com/af83/edwig/logger"
)

type SIRIStopMonitoringRequest struct {
	MessageIdentifier string
	MonitoringRef     string
	RequestorRef      string
	RequestTimestamp  time.Time
}

const StopMonitoringRequestTemplate = `<ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
													 xmlns:ns3="http://www.ifopt.org.uk/acsb"
													 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
													 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
													 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
	<ServiceRequestInfo>
		<ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z"}}</ns2:RequestTimestamp>
		<ns2:RequestorRef>{{.RequestorRef}}</ns2:RequestorRef>
		<ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
	</ServiceRequestInfo>
	<Request version="2.0:FR-IDF-2.4">
		<ns2:RequestTimestamp>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z"}}</ns2:RequestTimestamp>
		<ns2:MessageIdentifier>{{.MessageIdentifier}}</ns2:MessageIdentifier>
		<ns2:StartTime>{{.RequestTimestamp.Format "2006-01-02T15:04:05.000Z"}}</ns2:StartTime>
		<ns2:MonitoringRef>{{.MonitoringRef}}</ns2:MonitoringRef>
		<ns2:StopVisitTypes>all</ns2:StopVisitTypes>
	</Request>
	<RequestExtension />
</ns7:GetStopMonitoring>`

func (request *SIRIStopMonitoringRequest) BuildXML() string {
	var buffer bytes.Buffer
	var siriRequest = template.Must(template.New("siriRequest").Parse(StopMonitoringRequestTemplate))
	if err := siriRequest.Execute(&buffer, request); err != nil {
		logger.Log.Panicf("Error while using request template: %v", err)
	}
	return buffer.String()
}
