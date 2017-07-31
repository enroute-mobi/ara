package siri

import (
	"testing"
	"time"

	"github.com/af83/edwig/model"
)

func Test_SIRIEstimatedTimeTableResponse_BuildXML(t *testing.T) {
	expectedXML := `<ns8:GetEstimatedTimetableResponse xmlns:ns3="http://www.siri.org.uk/siri"
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
		<ns3:ResponseMessageIdentifier>response</ns3:ResponseMessageIdentifier>
		<ns3:RequestMessageRef>request</ns3:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<ns3:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<ns3:ResponseTimestamp>2016-09-21T20:14:46.000Z</ns3:ResponseTimestamp>
			<ns3:RequestMessageRef>request</ns3:RequestMessageRef>
			<ns3:Status>true</ns3:Status>
			<ns3:EstimatedJourneyVersionFrame>
				<ns3:RecordedAtTime>2016-09-21T20:14:46.000Z</ns3:RecordedAtTime>
				<ns3:EstimatedVehicleJourney>
					<ns3:LineRef>line1</ns3:LineRef>
					<ns3:DirectionRef>direction1</ns3:DirectionRef>
					<ns3:DatedVehicleJourneyRef>dvjref1</ns3:DatedVehicleJourneyRef>
					<ns3:PublishedLineName>line 1</ns3:PublishedLineName>
					<ns3:OriginRef>origin1</ns3:OriginRef>
					<ns3:OriginName>origin 1</ns3:OriginName>
					<ns3:DestinationRef>destination1</ns3:DestinationRef>
					<ns3:DestinationName>destination 1</ns3:DestinationName>
					<ns3:EstimatedCalls>
						<ns3:EstimatedCall>
							<ns3:StopPointRef>stopPoint1</ns3:StopPointRef>
							<ns3:Order>1</ns3:Order>
							<ns3:ActualArrivalTime>2016-09-21T20:14:46.000Z</ns3:ActualArrivalTime>
							<ns3:ArrivalStatus>astatus1</ns3:ArrivalStatus>
							<ns3:ActualDepartureTime>2016-09-21T20:14:46.000Z</ns3:ActualDepartureTime>
							<ns3:DepartureStatus>dstatus1</ns3:DepartureStatus>
						</ns3:EstimatedCall>
						<ns3:EstimatedCall>
							<ns3:StopPointRef>stopPoint2</ns3:StopPointRef>
							<ns3:Order>2</ns3:Order>
							<ns3:ActualArrivalTime>2016-09-21T20:14:46.000Z</ns3:ActualArrivalTime>
							<ns3:ArrivalStatus>astatus2</ns3:ArrivalStatus>
							<ns3:ActualDepartureTime>2016-09-21T20:14:46.000Z</ns3:ActualDepartureTime>
							<ns3:DepartureStatus>dstatus2</ns3:DepartureStatus>
						</ns3:EstimatedCall>
					</ns3:EstimatedCalls>
				</ns3:EstimatedVehicleJourney>
			</ns3:EstimatedJourneyVersionFrame>
			<ns3:EstimatedJourneyVersionFrame>
				<ns3:RecordedAtTime>2016-09-21T20:14:46.000Z</ns3:RecordedAtTime>
				<ns3:EstimatedVehicleJourney>
					<ns3:LineRef>line2</ns3:LineRef>
					<ns3:DirectionRef>direction2</ns3:DirectionRef>
					<ns3:DatedVehicleJourneyRef>dvjref2</ns3:DatedVehicleJourneyRef>
					<ns3:PublishedLineName>line 2</ns3:PublishedLineName>
					<ns3:OriginRef>origin2</ns3:OriginRef>
					<ns3:OriginName>origin 2</ns3:OriginName>
					<ns3:DestinationRef>destination2</ns3:DestinationRef>
					<ns3:DestinationName>destination 2</ns3:DestinationName>
					<ns3:EstimatedCalls>
						<ns3:EstimatedCall>
							<ns3:StopPointRef>stopPoint3</ns3:StopPointRef>
							<ns3:Order>3</ns3:Order>
							<ns3:ActualArrivalTime>2016-09-21T20:14:46.000Z</ns3:ActualArrivalTime>
							<ns3:ArrivalStatus>astatus3</ns3:ArrivalStatus>
							<ns3:ActualDepartureTime>2016-09-21T20:14:46.000Z</ns3:ActualDepartureTime>
							<ns3:DepartureStatus>dstatus3</ns3:DepartureStatus>
						</ns3:EstimatedCall>
					</ns3:EstimatedCalls>
				</ns3:EstimatedVehicleJourney>
				<ns3:EstimatedVehicleJourney>
					<ns3:LineRef>line3</ns3:LineRef>
					<ns3:DirectionRef>direction3</ns3:DirectionRef>
					<ns3:DatedVehicleJourneyRef>dvjref3</ns3:DatedVehicleJourneyRef>
					<ns3:PublishedLineName>line 3</ns3:PublishedLineName>
					<ns3:OriginRef>origin3</ns3:OriginRef>
					<ns3:OriginName>origin 3</ns3:OriginName>
					<ns3:DestinationRef>destination3</ns3:DestinationRef>
					<ns3:DestinationName>destination 3</ns3:DestinationName>
					<ns3:EstimatedCalls>
						<ns3:EstimatedCall>
							<ns3:StopPointRef>stopPoint4</ns3:StopPointRef>
							<ns3:Order>4</ns3:Order>
							<ns3:ActualArrivalTime>2016-09-21T20:14:46.000Z</ns3:ActualArrivalTime>
							<ns3:ArrivalStatus>astatus4</ns3:ArrivalStatus>
							<ns3:ActualDepartureTime>2016-09-21T20:14:46.000Z</ns3:ActualDepartureTime>
							<ns3:DepartureStatus>dstatus4</ns3:DepartureStatus>
						</ns3:EstimatedCall>
					</ns3:EstimatedCalls>
				</ns3:EstimatedVehicleJourney>
			</ns3:EstimatedJourneyVersionFrame>
		</ns3:EstimatedTimetableDelivery>
	</Answer>
</ns8:GetEstimatedTimetableResponse>`

	testTime := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	call1 := SIRIEstimatedCall{
		ArrivalStatus:       "astatus1",
		DepartureStatus:     "dstatus1",
		StopPointRef:        "stopPoint1",
		Order:               1,
		ActualArrivalTime:   testTime,
		ActualDepartureTime: testTime,
	}
	call2 := SIRIEstimatedCall{
		ArrivalStatus:       "astatus2",
		DepartureStatus:     "dstatus2",
		StopPointRef:        "stopPoint2",
		Order:               2,
		ActualArrivalTime:   testTime,
		ActualDepartureTime: testTime,
	}
	call3 := SIRIEstimatedCall{
		ArrivalStatus:       "astatus3",
		DepartureStatus:     "dstatus3",
		StopPointRef:        "stopPoint3",
		Order:               3,
		ActualArrivalTime:   testTime,
		ActualDepartureTime: testTime,
	}
	call4 := SIRIEstimatedCall{
		ArrivalStatus:       "astatus4",
		DepartureStatus:     "dstatus4",
		StopPointRef:        "stopPoint4",
		Order:               4,
		ActualArrivalTime:   testTime,
		ActualDepartureTime: testTime,
	}

	// LineRef           string
	// PublishedLineName string

	// Attributes map[string]string

	vehicleJourney1 := SIRIEstimatedVehicleJourney{
		LineRef:                "line1",
		PublishedLineName:      "line 1",
		DatedVehicleJourneyRef: "dvjref1",
		Attributes: map[string]string{
			"DirectionRef":    "direction1",
			"OriginName":      "origin 1",
			"DestinationName": "destination 1",
		},
		References: map[string]model.Reference{
			"OriginRef":      *model.NewReference(model.NewObjectID("kind", "origin1")),
			"DestinationRef": *model.NewReference(model.NewObjectID("kind", "destination1")),
		},
		EstimatedCalls: []SIRIEstimatedCall{call1, call2},
	}
	vehicleJourney2 := SIRIEstimatedVehicleJourney{
		LineRef:                "line2",
		PublishedLineName:      "line 2",
		DatedVehicleJourneyRef: "dvjref2",
		Attributes: map[string]string{
			"DirectionRef":    "direction2",
			"OriginName":      "origin 2",
			"DestinationName": "destination 2",
		},
		References: map[string]model.Reference{
			"OriginRef":      *model.NewReference(model.NewObjectID("kind", "origin2")),
			"DestinationRef": *model.NewReference(model.NewObjectID("kind", "destination2")),
		},
		EstimatedCalls: []SIRIEstimatedCall{call3},
	}
	vehicleJourney3 := SIRIEstimatedVehicleJourney{
		LineRef:                "line3",
		PublishedLineName:      "line 3",
		DatedVehicleJourneyRef: "dvjref3",
		Attributes: map[string]string{
			"DirectionRef":    "direction3",
			"OriginName":      "origin 3",
			"DestinationName": "destination 3",
		},
		References: map[string]model.Reference{
			"OriginRef":      *model.NewReference(model.NewObjectID("kind", "origin3")),
			"DestinationRef": *model.NewReference(model.NewObjectID("kind", "destination3")),
		},
		EstimatedCalls: []SIRIEstimatedCall{call4},
	}

	journeyVersion1 := SIRIEstimatedJourneyVersionFrame{
		RecordedAtTime:           testTime,
		EstimatedVehicleJourneys: []SIRIEstimatedVehicleJourney{vehicleJourney1},
	}
	journeyVersion2 := SIRIEstimatedJourneyVersionFrame{
		RecordedAtTime:           testTime,
		EstimatedVehicleJourneys: []SIRIEstimatedVehicleJourney{vehicleJourney2, vehicleJourney3},
	}

	response := &SIRIEstimatedTimeTableResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		RequestMessageRef:         "request",
		ResponseMessageIdentifier: "response",
		ResponseTimestamp:         testTime,
		Status:                    true,
		EstimatedJourneyVersionFrames: []SIRIEstimatedJourneyVersionFrame{journeyVersion1, journeyVersion2},
	}

	xml, err := response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}
