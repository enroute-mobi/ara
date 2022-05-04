package siri

import (
	"testing"
	"time"
)

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
		ItemIdentifier:        "itemId",
		MonitoringRef:         "monitoringRef",
		StopPointRef:          "stopPointRef",
		StopPointName:         "stopPointName",
		LineRef:               "lineRef",
		PublishedLineName:     "lineName",
		DepartureStatus:       "depStatus",
		ArrivalStatus:         "arrStatus",
		VehicleJourneyName:    "NameOfVj",
		VehicleAtStop:         true,
		Order:                 1,
		Monitored:             true,
		RecordedAt:            time.Date(2015, time.September, 21, 20, 14, 46, 0, time.UTC),
		DataFrameRef:          "2016-09-21",
		AimedArrivalTime:      time.Date(2017, time.September, 21, 20, 14, 46, 0, time.UTC),
		ActualArrivalTime:     time.Date(2018, time.September, 21, 20, 14, 46, 0, time.UTC),
		AimedDepartureTime:    time.Date(2019, time.September, 21, 20, 14, 46, 0, time.UTC),
		ExpectedDepartureTime: time.Date(2020, time.September, 21, 20, 14, 46, 0, time.UTC),
		Attributes:            make(map[string]map[string]string),
		References:            make(map[string]map[string]string),
	}

	siriMonitoredStopVisit.Attributes["StopVisitAttributes"] = make(map[string]string)
	siriMonitoredStopVisit.References["VehicleJourney"] = make(map[string]string)
	siriMonitoredStopVisit.References["StopVisitReferences"] = make(map[string]string)
	siriMonitoredStopVisit.Attributes["VehicleJourneyAttributes"] = make(map[string]string)
	siriMonitoredStopVisit.Attributes["VehicleJourneyAttributes"]["Delay"] = "30"
	siriMonitoredStopVisit.DatedVehicleJourneyRef = "vehicleJourney#ObjectID"
	siriMonitoredStopVisit.References["StopVisitReferences"]["OperatorRef"] = "OperatorRef"
	siriMonitoredStopVisit.References["VehicleJourney"]["DestinationRef"] = "NINOXE:StopPoint:SP:62:LOC"

	response.MonitoredStopVisits = []*SIRIMonitoredStopVisit{siriMonitoredStopVisit}
	xml, err = response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}
