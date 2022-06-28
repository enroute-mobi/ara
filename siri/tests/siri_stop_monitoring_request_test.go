package siri_tests

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
)

func Test_SIRIStopMonitoringRequest_BuildXML(t *testing.T) {
	expectedXML := `<sw:GetStopMonitoring xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceRequestInfo>
		<siri:RequestTimestamp>2009-11-10T23:00:00.000Z</siri:RequestTimestamp>
		<siri:RequestorRef>test</siri:RequestorRef>
		<siri:MessageIdentifier>test</siri:MessageIdentifier>
	</ServiceRequestInfo>
	<Request version="2.0:FR-IDF-2.4">
		<siri:RequestTimestamp>2009-11-10T23:00:00.000Z</siri:RequestTimestamp>
		<siri:MessageIdentifier>test</siri:MessageIdentifier>
		<siri:MonitoringRef>test</siri:MonitoringRef>
		<siri:StopVisitTypes>all</siri:StopVisitTypes>
	</Request>
	<RequestExtension />
</sw:GetStopMonitoring>`
	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &siri.SIRIGetStopMonitoringRequest{
		RequestorRef: "test",
	}
	request.MessageIdentifier = "test"
	request.MonitoringRef = "test"
	request.RequestTimestamp = date
	request.StopVisitTypes = "all"

	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}
