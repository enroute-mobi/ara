package siri_tests

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"github.com/stretchr/testify/assert"
)

func Test_SIRIEstimatedTimetableResponse_BuildXML(t *testing.T) {
	assert := assert.New(t)

	expectedXML := `<sw:GetEstimatedTimetableResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>producer</siri:ProducerRef>
		<siri:Address>address</siri:Address>
		<siri:ResponseMessageIdentifier>response</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>request</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef>request</siri:RequestMessageRef>
			<siri:Status>true</siri:Status>
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>2016-09-21T20:14:46.000Z</siri:RecordedAtTime>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>line1</siri:LineRef>
					<siri:DirectionRef>direction1</siri:DirectionRef>
					<siri:DatedVehicleJourneyRef>dvjref1</siri:DatedVehicleJourneyRef>
					<siri:OriginRef>origin1</siri:OriginRef>
					<siri:DestinationRef>destination1</siri:DestinationRef>
					<siri:DestinationName>destination 1</siri:DestinationName>
					<siri:OperatorRef>operator1</siri:OperatorRef>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopPoint1</siri:StopPointRef>
							<siri:Order>1</siri:Order>
							<siri:AimedArrivalTime>2016-09-21T20:14:46.000Z</siri:AimedArrivalTime>
							<siri:ArrivalStatus>astatus1</siri:ArrivalStatus>
							<siri:AimedDepartureTime>2016-09-21T20:14:46.000Z</siri:AimedDepartureTime>
							<siri:DepartureStatus>dstatus1</siri:DepartureStatus>
						</siri:EstimatedCall>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopPoint2</siri:StopPointRef>
							<siri:Order>2</siri:Order>
							<siri:AimedArrivalTime>2016-09-21T20:14:46.000Z</siri:AimedArrivalTime>
							<siri:ArrivalStatus>astatus2</siri:ArrivalStatus>
							<siri:AimedDepartureTime>2016-09-21T20:14:46.000Z</siri:AimedDepartureTime>
							<siri:DepartureStatus>dstatus2</siri:DepartureStatus>
						</siri:EstimatedCall>
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>
			</siri:EstimatedJourneyVersionFrame>
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>2016-09-21T20:14:46.000Z</siri:RecordedAtTime>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>line2</siri:LineRef>
					<siri:DirectionRef>direction2</siri:DirectionRef>
					<siri:DatedVehicleJourneyRef>dvjref2</siri:DatedVehicleJourneyRef>
					<siri:OriginRef>origin2</siri:OriginRef>
					<siri:DestinationRef>destination2</siri:DestinationRef>
					<siri:DestinationName>destination 2</siri:DestinationName>
					<siri:OperatorRef>operator2</siri:OperatorRef>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopPoint3</siri:StopPointRef>
							<siri:Order>3</siri:Order>
							<siri:AimedArrivalTime>2016-09-21T20:14:46.000Z</siri:AimedArrivalTime>
							<siri:ArrivalStatus>astatus3</siri:ArrivalStatus>
							<siri:AimedDepartureTime>2016-09-21T20:14:46.000Z</siri:AimedDepartureTime>
							<siri:DepartureStatus>dstatus3</siri:DepartureStatus>
						</siri:EstimatedCall>
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>line3</siri:LineRef>
					<siri:DirectionRef>direction3</siri:DirectionRef>
					<siri:DatedVehicleJourneyRef>dvjref3</siri:DatedVehicleJourneyRef>
					<siri:OriginRef>origin3</siri:OriginRef>
					<siri:DestinationRef>destination3</siri:DestinationRef>
					<siri:DestinationName>destination 3</siri:DestinationName>
					<siri:OperatorRef>operator3</siri:OperatorRef>
					<siri:EstimatedCalls>
						<siri:EstimatedCall>
							<siri:StopPointRef>stopPoint4</siri:StopPointRef>
							<siri:Order>4</siri:Order>
							<siri:AimedArrivalTime>2016-09-21T20:14:46.000Z</siri:AimedArrivalTime>
							<siri:ArrivalStatus>astatus4</siri:ArrivalStatus>
							<siri:AimedDepartureTime>2016-09-21T20:14:46.000Z</siri:AimedDepartureTime>
							<siri:DepartureStatus>dstatus4</siri:DepartureStatus>
						</siri:EstimatedCall>
					</siri:EstimatedCalls>
				</siri:EstimatedVehicleJourney>
			</siri:EstimatedJourneyVersionFrame>
		</siri:EstimatedTimetableDelivery>
	</Answer>
	<AnswerExtension/>
</sw:GetEstimatedTimetableResponse>`

	testTime := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	call1 := &siri.SIRIEstimatedCall{
		ArrivalStatus:      "astatus1",
		DepartureStatus:    "dstatus1",
		StopPointRef:       "stopPoint1",
		Order:              1,
		AimedArrivalTime:   testTime,
		AimedDepartureTime: testTime,
	}
	call2 := &siri.SIRIEstimatedCall{
		ArrivalStatus:      "astatus2",
		DepartureStatus:    "dstatus2",
		StopPointRef:       "stopPoint2",
		Order:              2,
		AimedArrivalTime:   testTime,
		AimedDepartureTime: testTime,
	}
	call3 := &siri.SIRIEstimatedCall{
		ArrivalStatus:      "astatus3",
		DepartureStatus:    "dstatus3",
		StopPointRef:       "stopPoint3",
		Order:              3,
		AimedArrivalTime:   testTime,
		AimedDepartureTime: testTime,
	}
	call4 := &siri.SIRIEstimatedCall{
		ArrivalStatus:      "astatus4",
		DepartureStatus:    "dstatus4",
		StopPointRef:       "stopPoint4",
		Order:              4,
		AimedArrivalTime:   testTime,
		AimedDepartureTime: testTime,
	}

	vehicleJourney1 := &siri.SIRIEstimatedVehicleJourney{
		LineRef:                "line1",
		DatedVehicleJourneyRef: "dvjref1",
		DirectionType:          "direction1",
		Attributes: map[string]string{
			"OriginName":      "origin 1",
			"DestinationName": "destination 1",
		},
		References: map[string]string{
			"OriginRef":      "origin1",
			"DestinationRef": "destination1",
			"OperatorRef":    "operator1",
		},
		EstimatedCalls: []*siri.SIRIEstimatedCall{call1, call2},
	}
	vehicleJourney2 := &siri.SIRIEstimatedVehicleJourney{
		LineRef:                "line2",
		DatedVehicleJourneyRef: "dvjref2",
		DirectionType:          "direction2",
		Attributes: map[string]string{
			"OriginName":      "origin 2",
			"DestinationName": "destination 2",
		},
		References: map[string]string{
			"OriginRef":      "origin2",
			"DestinationRef": "destination2",
			"OperatorRef":    "operator2",
		},
		EstimatedCalls: []*siri.SIRIEstimatedCall{call3},
	}
	vehicleJourney3 := &siri.SIRIEstimatedVehicleJourney{
		LineRef:                "line3",
		DatedVehicleJourneyRef: "dvjref3",
		DirectionType:          "direction3",
		Attributes: map[string]string{
			"OriginName":      "origin 3",
			"DestinationName": "destination 3",
		},
		References: map[string]string{
			"OriginRef":      "origin3",
			"DestinationRef": "destination3",
			"OperatorRef":    "operator3",
		},
		EstimatedCalls: []*siri.SIRIEstimatedCall{call4},
	}

	journeyVersion1 := &siri.SIRIEstimatedJourneyVersionFrame{
		RecordedAtTime:           testTime,
		EstimatedVehicleJourneys: []*siri.SIRIEstimatedVehicleJourney{vehicleJourney1},
	}
	journeyVersion2 := &siri.SIRIEstimatedJourneyVersionFrame{
		RecordedAtTime:           testTime,
		EstimatedVehicleJourneys: []*siri.SIRIEstimatedVehicleJourney{vehicleJourney2, vehicleJourney3},
	}

	response := &siri.SIRIEstimatedTimetableResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		ResponseMessageIdentifier: "response",
	}
	response.RequestMessageRef = "request"
	response.Status = true
	response.ResponseTimestamp = testTime
	response.EstimatedJourneyVersionFrames = []*siri.SIRIEstimatedJourneyVersionFrame{journeyVersion1, journeyVersion2}

	xml, err := response.BuildXML()
	assert.NoError(err)
	assert.Equal(expectedXML, xml)
}

func Test_SIRIEstimatedTimetableResponse_BuildXML_EmptyCalls(t *testing.T) {
	assert := assert.New(t)

	expectedXML := `<sw:GetEstimatedTimetableResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<ServiceDeliveryInfo>
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>producer</siri:ProducerRef>
		<siri:Address>address</siri:Address>
		<siri:ResponseMessageIdentifier>response</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>request</siri:RequestMessageRef>
	</ServiceDeliveryInfo>
	<Answer>
		<siri:EstimatedTimetableDelivery version="2.0:FR-IDF-2.4">
			<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
			<siri:RequestMessageRef>request</siri:RequestMessageRef>
			<siri:Status>true</siri:Status>
			<siri:EstimatedJourneyVersionFrame>
				<siri:RecordedAtTime>2016-09-21T20:14:46.000Z</siri:RecordedAtTime>
				<siri:EstimatedVehicleJourney>
					<siri:LineRef>line1</siri:LineRef>
					<siri:DirectionRef>direction1</siri:DirectionRef>
					<siri:DatedVehicleJourneyRef>dvjref1</siri:DatedVehicleJourneyRef>
					<siri:OriginRef>origin1</siri:OriginRef>
					<siri:DestinationRef>destination1</siri:DestinationRef>
					<siri:DestinationName>destination 1</siri:DestinationName>
					<siri:OperatorRef>operator1</siri:OperatorRef>
				</siri:EstimatedVehicleJourney>
			</siri:EstimatedJourneyVersionFrame>
		</siri:EstimatedTimetableDelivery>
	</Answer>
	<AnswerExtension/>
</sw:GetEstimatedTimetableResponse>`

	testTime := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)

	vehicleJourney1 := &siri.SIRIEstimatedVehicleJourney{
		LineRef:                "line1",
		DatedVehicleJourneyRef: "dvjref1",
		DirectionType:          "direction1",
		Attributes: map[string]string{
			"OriginName":      "origin 1",
			"DestinationName": "destination 1",
		},
		References: map[string]string{
			"OriginRef":      "origin1",
			"DestinationRef": "destination1",
			"OperatorRef":    "operator1",
		},
	}

	journeyVersion1 := &siri.SIRIEstimatedJourneyVersionFrame{
		RecordedAtTime:           testTime,
		EstimatedVehicleJourneys: []*siri.SIRIEstimatedVehicleJourney{vehicleJourney1},
	}

	response := &siri.SIRIEstimatedTimetableResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		ResponseMessageIdentifier: "response",
	}
	response.RequestMessageRef = "request"
	response.Status = true
	response.ResponseTimestamp = testTime
	response.EstimatedJourneyVersionFrames = []*siri.SIRIEstimatedJourneyVersionFrame{journeyVersion1}

	xml, err := response.BuildXML()
	assert.NoError(err)
	assert.Equal(expectedXML, xml)
}
