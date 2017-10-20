package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/af83/edwig/model"
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
	if expected := model.STOP_VISIT_ARRIVAL_ARRIVED; model.StopVisitArrivalStatus(monitoredStopVisit.ArrivalStatus()) != expected {
		t.Errorf("Incorrect ArrivalStatus for stopVisit:\n expected: \"%v\"\n got: \"%v\"", expected, monitoredStopVisit.ArrivalStatus())
	}
	if expected := 4; monitoredStopVisit.Order() != expected {
		t.Errorf("Incorrect Order for stopVisit:\n expected: \"%v\"\n got: \"%v\"", expected, monitoredStopVisit.Order())
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

func Test_SIRIStopMonitoringResponse_BuildXML(t *testing.T) {
	expectedXML := `<sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>producer</siri:ProducerRef>
		<siri:Address>address</siri:Address>
		<siri:ResponseMessageIdentifier>identifier</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>ref</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef>ref</siri:RequestMessageRef>
			<siri:MonitoringRef>MonitoringRef</siri:MonitoringRef>
			<siri:Status>true</siri:Status>
		</siri:StopMonitoringDelivery>
	</Answer>
	<AnswerExtension/>
</sw:GetStopMonitoringResponse>`

	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	response := &SIRIStopMonitoringResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		ResponseMessageIdentifier: "identifier",
	}
	response.RequestMessageRef = "ref"
	response.Status = true
	response.ResponseTimestamp = responseTimestamp
	response.MonitoringRef = "MonitoringRef"
	xml, err := response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}

	expectedXML = `<sw:GetStopMonitoringResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>producer</siri:ProducerRef>
		<siri:Address>address</siri:Address>
		<siri:ResponseMessageIdentifier>identifier</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>ref</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<siri:StopMonitoringDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef>ref</siri:RequestMessageRef>
			<siri:MonitoringRef>MonitoringRef</siri:MonitoringRef>
			<siri:Status>true</siri:Status>
			<siri:MonitoredStopVisit>
				<siri:RecordedAtTime>2015-09-21T20:14:46.000Z</siri:RecordedAtTime>
				<siri:ItemIdentifier>itemId</siri:ItemIdentifier>
				<siri:MonitoringRef>monitoringRef</siri:MonitoringRef>
				<siri:MonitoredVehicleJourney>
					<siri:LineRef>lineRef</siri:LineRef>
					<siri:FramedVehicleJourneyRef>
						<siri:DataFrameRef>2016-09-21</siri:DataFrameRef>
						<siri:DatedVehicleJourneyRef>vehicleJourney#ObjectID</siri:DatedVehicleJourneyRef>
					</siri:FramedVehicleJourneyRef>
					<siri:PublishedLineName>lineName</siri:PublishedLineName>
					<siri:OperatorRef>OperatorRef</siri:OperatorRef>
					<siri:DestinationRef>NINOXE:StopPoint:SP:62:LOC</siri:DestinationRef>
					<siri:VehicleJourneyName>NameOfVj</siri:VehicleJourneyName>
					<siri:Monitored>true</siri:Monitored>
					<siri:Delay>30</siri:Delay>
					<siri:MonitoredCall>
						<siri:StopPointRef>stopPointRef</siri:StopPointRef>
						<siri:Order>1</siri:Order>
						<siri:StopPointName>stopPointName</siri:StopPointName>
						<siri:VehicleAtStop>true</siri:VehicleAtStop>
						<siri:AimedArrivalTime>2017-09-21T20:14:46.000Z</siri:AimedArrivalTime>
						<siri:ActualArrivalTime>2018-09-21T20:14:46.000Z</siri:ActualArrivalTime>
						<siri:ArrivalStatus>arrStatus</siri:ArrivalStatus>
						<siri:AimedDepartureTime>2019-09-21T20:14:46.000Z</siri:AimedDepartureTime>
						<siri:ExpectedDepartureTime>2020-09-21T20:14:46.000Z</siri:ExpectedDepartureTime>
						<siri:DepartureStatus>depStatus</siri:DepartureStatus>
					</siri:MonitoredCall>
				</siri:MonitoredVehicleJourney>
			</siri:MonitoredStopVisit>
		</siri:StopMonitoringDelivery>
	</Answer>
	<AnswerExtension/>
</sw:GetStopMonitoringResponse>`
	siriMonitoredStopVisit := &SIRIMonitoredStopVisit{
		ItemIdentifier:     "itemId",
		MonitoringRef:      "monitoringRef",
		StopPointRef:       "stopPointRef",
		StopPointName:      "stopPointName",
		LineRef:            "lineRef",
		PublishedLineName:  "lineName",
		DepartureStatus:    "depStatus",
		ArrivalStatus:      "arrStatus",
		VehicleJourneyName: "NameOfVj",
		VehicleAtStop:      true,
		Order:              1,
		Monitored:          true,
		RecordedAt:         time.Date(2015, time.September, 21, 20, 14, 46, 0, time.UTC),
		DataFrameRef:       "2016-09-21",
		AimedArrivalTime:   time.Date(2017, time.September, 21, 20, 14, 46, 0, time.UTC),
		// ExpectedArrivalTime: time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC),
		ActualArrivalTime:     time.Date(2018, time.September, 21, 20, 14, 46, 0, time.UTC),
		AimedDepartureTime:    time.Date(2019, time.September, 21, 20, 14, 46, 0, time.UTC),
		ExpectedDepartureTime: time.Date(2020, time.September, 21, 20, 14, 46, 0, time.UTC),
		// ActualDepartureTime: time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC),
		Attributes: make(map[string]map[string]string),
		References: make(map[string]map[string]model.Reference),
	}

	operatorRefObjId := model.NewObjectID("intenal", "OperatorRef")
	destinationRefObjId := model.NewObjectID("intenal", "NINOXE:StopPoint:SP:62:LOC")

	siriMonitoredStopVisit.Attributes["StopVisitAttributes"] = make(map[string]string)
	siriMonitoredStopVisit.References["VehicleJourney"] = make(map[string]model.Reference)
	siriMonitoredStopVisit.References["StopVisitReferences"] = make(map[string]model.Reference)
	siriMonitoredStopVisit.Attributes["VehicleJourneyAttributes"] = make(map[string]string)
	siriMonitoredStopVisit.Attributes["VehicleJourneyAttributes"]["Delay"] = "30"
	siriMonitoredStopVisit.DatedVehicleJourneyRef = "vehicleJourney#ObjectID"
	siriMonitoredStopVisit.References["StopVisitReferences"]["OperatorRef"] = model.Reference{ObjectId: &operatorRefObjId}
	siriMonitoredStopVisit.References["VehicleJourney"]["DestinationRef"] = model.Reference{ObjectId: &destinationRefObjId}

	response.MonitoredStopVisits = []*SIRIMonitoredStopVisit{siriMonitoredStopVisit}
	xml, err = response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}
