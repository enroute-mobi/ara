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
	if expected := "arrived"; monitoredStopVisit.ArrivalStatus() != expected {
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
	<AnswerExtension/>
</ns8:GetStopMonitoringResponse>`
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	request := &SIRIStopMonitoringResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		RequestMessageRef:         "ref",
		ResponseMessageIdentifier: "identifier",
		Status:            true,
		ResponseTimestamp: responseTimestamp,
	}
	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}

	expectedXML = `<ns8:GetStopMonitoringResponse xmlns:ns3="http://www.siri.org.uk/siri"
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
			<ns3:MonitoredStopVisit>
				<ns3:RecordedAtTime>2015-09-21T20:14:46.000Z</ns3:RecordedAtTime>
				<ns3:ItemIdentifier>itemId</ns3:ItemIdentifier>
				<ns3:MonitoringRef>stopPointRef</ns3:MonitoringRef>
				<ns3:MonitoredVehicleJourney>
					<ns3:LineRef>lineRef</ns3:LineRef>
					<ns3:FramedVehicleJourneyRef>
						<ns3:DataFrameRef>2016-09-21</ns3:DataFrameRef>
						<ns3:DatedVehicleJourneyRef>vehicleJourneyRef</ns3:DatedVehicleJourneyRef>
					</ns3:FramedVehicleJourneyRef>
					<ns3:PublishedLineName>lineName</ns3:PublishedLineName>
					<ns3:OperatorRef>OperatorRef</ns3:OperatorRef>
					<ns3:DestinationRef>NINOXE:StopPoint:SP:62:LOC</ns3:DestinationRef>
					<ns3:VehicleJourneyName>NameOfVj</ns3:VehicleJourneyName>
					<ns3:Delay>30</ns3:Delay>
					<ns3:MonitoredCall>
						<ns3:StopPointRef>stopPointRef</ns3:StopPointRef>
						<ns3:Order>1</ns3:Order>
						<ns3:StopPointName>stopPointName</ns3:StopPointName>
						<ns3:VehicleAtStop>true</ns3:VehicleAtStop>
						<ns3:AimedArrivalTime>2017-09-21T20:14:46.000Z</ns3:AimedArrivalTime>
						<ns3:ActualArrivalTime>2018-09-21T20:14:46.000Z</ns3:ActualArrivalTime>
						<ns3:ArrivalStatus>arrStatus</ns3:ArrivalStatus>
						<ns3:AimedDepartureTime>2019-09-21T20:14:46.000Z</ns3:AimedDepartureTime>
						<ns3:ExpectedDepartureTime>2020-09-21T20:14:46.000Z</ns3:ExpectedDepartureTime>
						<ns3:DepartureStatus>depStatus</ns3:DepartureStatus>
					</ns3:MonitoredCall>
				</ns3:MonitoredVehicleJourney>
			</ns3:MonitoredStopVisit>
		</ns3:StopMonitoringDelivery>
	</Answer>
	<AnswerExtension/>
</ns8:GetStopMonitoringResponse>`
	siriMonitoredStopVisit := &SIRIMonitoredStopVisit{
		ItemIdentifier:     "itemId",
		StopPointRef:       "stopPointRef",
		StopPointName:      "stopPointName",
		LineRef:            "lineRef",
		PublishedLineName:  "lineName",
		DepartureStatus:    "depStatus",
		ArrivalStatus:      "arrStatus",
		VehicleJourneyName: "NameOfVj",
		VehicleAtStop:      true,
		Order:              1,
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

	destinationRefObjId := model.NewObjectID("intenal", "NINOXE:StopPoint:SP:62:LOC")
	datedVehicleJourneyRefObjId := model.NewObjectID("intenal", "vehicleJourneyRef")
	operatorRefObjId := model.NewObjectID("intenal", "OperatorRef")

	siriMonitoredStopVisit.Attributes["StopVisitAttributes"] = make(map[string]string)
	siriMonitoredStopVisit.References["VehicleJourney"] = make(map[string]model.Reference)
	siriMonitoredStopVisit.References["StopVisitReferences"] = make(map[string]model.Reference)
	siriMonitoredStopVisit.Attributes["VehicleJourneyAttributes"] = make(map[string]string)

	siriMonitoredStopVisit.Attributes["VehicleJourneyAttributes"]["Delay"] = "30"
	siriMonitoredStopVisit.References["VehicleJourney"]["DestinationRef"] = model.Reference{ObjectId: &destinationRefObjId, Id: "42"}
	siriMonitoredStopVisit.References["VehicleJourney"]["DatedVehicleJourneyRef"] = model.Reference{ObjectId: &datedVehicleJourneyRefObjId, Id: "42"}
	siriMonitoredStopVisit.References["StopVisitReferences"]["OperatorRef"] = model.Reference{ObjectId: &operatorRefObjId, Id: "42"}

	request.MonitoredStopVisits = []*SIRIMonitoredStopVisit{siriMonitoredStopVisit}
	xml, err = request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}
