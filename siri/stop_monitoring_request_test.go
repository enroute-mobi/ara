package siri

import (
	"testing"
	"time"
)

func Test_SIRIStopMonitoringRequest_BuildXML(t *testing.T) {
	expectedXML := `<ns7:GetStopMonitoring xmlns:ns2="http://www.siri.org.uk/siri"
													 xmlns:ns3="http://www.ifopt.org.uk/acsb"
													 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
													 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
													 xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
	<ServiceRequestInfo>
		<ns2:RequestTimestamp>2009-11-10T23:00:00.000Z</ns2:RequestTimestamp>
		<ns2:RequestorRef>test</ns2:RequestorRef>
		<ns2:MessageIdentifier>test</ns2:MessageIdentifier>
	</ServiceRequestInfo>
	<Request version="2.0:FR-IDF-2.4">
		<ns2:RequestTimestamp>2009-11-10T23:00:00.000Z</ns2:RequestTimestamp>
		<ns2:MessageIdentifier>test</ns2:MessageIdentifier>
		<ns2:StartTime>2009-11-10T23:00:00.000Z</ns2:StartTime>
		<ns2:MonitoringRef>test</ns2:MonitoringRef>
		<ns2:StopVisitTypes>all</ns2:StopVisitTypes>
	</Request>
	<RequestExtension />
</ns7:GetStopMonitoring>`
	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := &SIRIStopMonitoringRequest{
		MessageIdentifier: "test",
		MonitoringRef:     "test",
		RequestorRef:      "test",
		RequestTimestamp:  date,
	}
	if expectedXML != request.BuildXML() {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", request.BuildXML(), expectedXML)
	}
}
