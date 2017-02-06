package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLStopMonitoringResponse(t *testing.T) *XMLStopMonitoringResponse {
	file, err := os.Open("testdata/stopmonitoring-response-soap.xml")
	defer file.Close()
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := NewXMLStopMonitoringResponseFromContent(content)
	return response
}

func Test_XMLStopMonitoringResponse_Address(t *testing.T) {
	response := getXMLStopMonitoringResponse(t)
	if expected := "http://appli.chouette.mobi/siri_france/siri"; response.Address() != expected {
		t.Errorf("Wrong Address:\n got: %v\nwant: %v", response.Address(), expected)
	}
}

func Test_XMLStopMonitoringResponse_ProducerRef(t *testing.T) {
	response := getXMLStopMonitoringResponse(t)
	if expected := "NINOXE:default"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef:\n got: %v\nwant: %v", response.ProducerRef(), expected)
	}
}

func Test_XMLStopMonitoringResponse_RequestMessageRef(t *testing.T) {
	response := getXMLStopMonitoringResponse(t)
	if expected := "StopMonitoring:Test:0"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef:\n got: %v\nwant: %v", response.RequestMessageRef(), expected)
	}
}

func Test_XMLStopMonitoringResponse_ResponseMessageIdentifier(t *testing.T) {
	response := getXMLStopMonitoringResponse(t)
	if expected := "fd0c67ac-2d3a-4ee5-9672-5f3f160cbd26"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Wrong ResponseMessageIdentifier:\n got: %v\nwant: %v", response.ResponseMessageIdentifier(), expected)
	}
}

func Test_XMLStopMonitoringResponse_ResponseTimestamp(t *testing.T) {
	response := getXMLStopMonitoringResponse(t)
	if expected := time.Date(2016, time.September, 22, 6, 01, 20, 227000000, time.UTC); !response.ResponseTimestamp().Equal(expected) {
		t.Errorf("Wrong ResponseTimestamp:\n got: %v\nwant: %v", response.ResponseTimestamp(), expected)
	}
}

func Test_XMLStopMonitoringResponse_XMLMonitoredStopVisit(t *testing.T) {
	response := getXMLStopMonitoringResponse(t)
	monitoredStopVisits := response.XMLMonitoredStopVisits()

	if len(monitoredStopVisits) != 2 {
		t.Errorf("Incorrect number of MonitoredStopVisit, expected 2 got %d", len(monitoredStopVisits))
	}
}

func Test_XMLMonitoredStopVisit(t *testing.T) {
	response := getXMLStopMonitoringResponse(t)
	monitoredStopVisit := response.XMLMonitoredStopVisits()[0]

	if expected := "NINOXE:VehicleJourney:201-NINOXE:StopPoint:SP:24:LOC-3"; monitoredStopVisit.ItemIdentifier() != expected {
		t.Errorf("Incorrect ItemIdentifier for stopVisit:\n expected: %v\n got: %v", expected, monitoredStopVisit.ItemIdentifier())
	}
	if expected := "NINOXE:StopPoint:Q:50:LOC"; monitoredStopVisit.StopPointRef() != expected {
		t.Errorf("Incorrect StopPointRef for stopVisit:\n expected: %v\n got: %v", expected, monitoredStopVisit.StopPointRef())
	}
	if expected := "Elf Sylvain - Métro (R)"; monitoredStopVisit.StopPointName() != expected {
		t.Errorf("Incorrect StopPointName for stopVisit:\n expected: %v\n got: %v", expected, monitoredStopVisit.StopPointName())
	}
	if expected := "NINOXE:VehicleJourney:201"; monitoredStopVisit.DatedVehicleJourneyRef() != expected {
		t.Errorf("Incorrect DatedVehicleJourneyRef for stopVisit:\n expected: %v\n got: %v", expected, monitoredStopVisit.DatedVehicleJourneyRef())
	}
	if expected := "NINOXE:Line:3:LOC"; monitoredStopVisit.LineRef() != expected {
		t.Errorf("Incorrect LineRef for stopVisit:\n expected: %v\n got: %v", expected, monitoredStopVisit.LineRef())
	}
	if expected := "Ligne 3 Metro"; monitoredStopVisit.PublishedLineName() != expected {
		t.Errorf("Incorrect PublishedLineName for stopVisit:\n expected: %v\n got: %v", expected, monitoredStopVisit.PublishedLineName())
	}
	if expected := ""; monitoredStopVisit.DepartureStatus() != expected {
		t.Errorf("Incorrect DepartureStatus for stopVisit:\n expected: \"%v\"\n got: \"%v\"", expected, monitoredStopVisit.DepartureStatus())
	}
	if expected := "arrived"; monitoredStopVisit.ArrivalStatus() != expected {
		t.Errorf("Incorrect ArrivalStatus for stopVisit:\n expected: \"%v\"\n got: \"%v\"", expected, monitoredStopVisit.ArrivalStatus())
	}
	if expected := 4; monitoredStopVisit.Order() != expected {
		t.Errorf("Incorrect Order for stopVisit:\n expected: \"%v\"\n got: \"%v\"", expected, monitoredStopVisit.Order())
	}
	if expected := time.Date(2016, time.September, 22, 5, 54, 0, 000000000, time.UTC); !monitoredStopVisit.AimedArrivalTime().Equal(expected) {
		t.Errorf("Incorrect AimedArrivalTime for stopVisit:\n expected: %v\n got: %v", expected, monitoredStopVisit.AimedArrivalTime())
	}
	if !monitoredStopVisit.ExpectedArrivalTime().IsZero() {
		t.Errorf("Incorrect ExpectedArrivalTime for stopVisit, should be zero got: %v", monitoredStopVisit.ExpectedArrivalTime())
	}
	if expected := time.Date(2016, time.September, 22, 5, 54, 0, 000000000, time.UTC); !monitoredStopVisit.ActualArrivalTime().Equal(expected) {
		t.Errorf("Incorrect ActualArrivalTime for stopVisit:\n expected: %v\n got: %v", expected, monitoredStopVisit.ActualArrivalTime())
	}
	if !monitoredStopVisit.AimedDepartureTime().IsZero() {
		t.Errorf("Incorrect AimedDepartureTime for stopVisit, should be zero got: %v", monitoredStopVisit.AimedDepartureTime())
	}
	if !monitoredStopVisit.ExpectedDepartureTime().IsZero() {
		t.Errorf("Incorrect ExpectedDepartureTime for stopVisit, should be zero got: %v", monitoredStopVisit.ExpectedDepartureTime())
	}
	if !monitoredStopVisit.ActualDepartureTime().IsZero() {
		t.Errorf("Incorrect ActualDepartureTime for stopVisit, should be zero got: %v", monitoredStopVisit.ActualDepartureTime())
	}
}

func Test_SIRIStopMonitoringResponse_BuildXML(t *testing.T) {
	expectedXML := `<ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
                               xmlns:ns4="http://www.ifopt.org.uk/acsb"
                               xmlns:ns5="http://www.ifopt.org.uk/ifopt"
                               xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0"
                               xmlns:ns7="http://scma/siri"
                               xmlns:ns8="http://wsdl.siri.org.uk"
                               xmlns:ns9="http://wsdl.siri.org.uk/siri">
  <ServiceDeliveryInfo>
    <ns3:ResponseTimestamp>2016-09-21T20:14:46.000Z</ns3:ResponseTimestamp>
    <ns3:ProducerRef>producer</ns3:ProducerRef>
    <ns3:Address>address</ns3:Address>
    <ns3:ResponseMessageIdentifier>identifier</ns3:ResponseMessageIdentifier>
    <ns3:RequestMessageRef>ref</ns3:RequestMessageRef>
  </ServiceDeliveryInfo>
  <Answer>
    <ns3:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
      <ns3:ResponseTimestamp>2016-09-21T20:14:46.000Z</ns3:ResponseTimestamp>
      <ns3:RequestMessageRef>ref</ns3:RequestMessageRef>
      <ns3:Status>true</ns3:Status>
    </ns3:StopMonitoringDelivery>
  </Answer>
  <AnswerExtension />
</ns8:GetStopMonitoringResponse>`
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	request := NewSIRIStopMonitoringResponse("address", "producer", "ref", "identifier", true, responseTimestamp)
	if expectedXML != request.BuildXML() {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", request.BuildXML(), expectedXML)
	}
}
