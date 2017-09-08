package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLGetStopMonitoring(t *testing.T) *XMLGetStopMonitoring {
	file, err := os.Open("testdata/stopmonitoring-request-soap.xml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := NewXMLGetStopMonitoringFromContent(content)
	return request
}

func Test_XMLGetStopMonitoring_RequestorRef(t *testing.T) {
	request := getXMLGetStopMonitoring(t)

	if expected := "NINOXE:default"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}
}

func Test_XMLGetStopMonitoring_RequestTimestamp(t *testing.T) {
	request := getXMLGetStopMonitoring(t)

	if expected := time.Date(2016, time.September, 22, 7, 54, 52, 977000000, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}
}

func Test_XMLGetStopMonitoring_MessageIdentifier(t *testing.T) {
	request := getXMLGetStopMonitoring(t)

	if expected := "StopMonitoring:Test:0"; request.MessageIdentifier() != expected {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", request.MessageIdentifier(), expected)
	}
}

func Test_XMLGetStopMonitoring_MonitoringRef(t *testing.T) {
	request := getXMLGetStopMonitoring(t)

	if expected := "NINOXE:StopPoint:SP:24:LOC"; request.MonitoringRef() != expected {
		t.Errorf("Wrong MonitoringRef:\n got: %v\nwant: %v", request.MonitoringRef(), expected)
	}
}

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
	request := &SIRIStopMonitoringRequest{
		MessageIdentifier: "test",
		MonitoringRef:     "test",
		RequestorRef:      "test",
		RequestTimestamp:  date,
	}
	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}
