package siri_tests

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLStopMonitoringResponse(t *testing.T) *sxml.XMLStopMonitoringResponse {
	file, err := os.Open("testdata/stopmonitoring-response-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLStopMonitoringResponseFromContent(content)
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
	monitoredStopVisits := response.StopMonitoringDeliveries()[0].XMLMonitoredStopVisits()

	if len(monitoredStopVisits) != 2 {
		t.Errorf("Incorrect number of MonitoredStopVisit, expected 2 got %d", len(monitoredStopVisits))
	}
}

func Test_XMLMonitoredStopVisit(t *testing.T) {
	response := getXMLStopMonitoringResponse(t)
	monitoredStopVisit := response.StopMonitoringDeliveries()[0].XMLMonitoredStopVisits()[0]
	othetMonitoredStopVisit := response.StopMonitoringDeliveries()[0].XMLMonitoredStopVisits()[1]

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
	if expected := 5; othetMonitoredStopVisit.Order() != expected {
		t.Errorf("Incorrect Order for stopVisit:\n expected: \"%v\"\n got: \"%v\"", expected, othetMonitoredStopVisit.Order())
	}
	if expected := "NINOXE:Company:15563880:LOC"; monitoredStopVisit.OperatorRef() != expected {
		t.Errorf("Incorrect OperatorRef for stopVisit:\n expected: \"%v\"\n got: \"%v\"", expected, monitoredStopVisit.OperatorRef())
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
