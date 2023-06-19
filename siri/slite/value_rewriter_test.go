package slite

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func testRewriteValues(t *testing.T, input string) string {
	output, err := RewriteValues([]byte(input))
	if err != nil {
		t.Fatal(err)
		return ""
	}
	return string(output)
}

func testCompactedJSON(input string) string {
	var out bytes.Buffer
	json.Compact(&out, []byte(input))
	return out.String()
}

func Test_RewriteValues_RewriteKey(t *testing.T) {
	assert := assert.New(t)

	output := testRewriteValues(t, `{ "DestinationRef": { "value": "STIF:StopPoint:Q:22141:" } }`)
	assert.Equal(`{"DestinationRef":"STIF:StopPoint:Q:22141:"}`, output)
}

func Test_RewriteValues_RewriteArray(t *testing.T) {
	assert := assert.New(t)

	output := testRewriteValues(t, `{ "DestinationName": [ { "value": "Porte de Clignancourt" } ] }`)
	assert.Equal(`{"DestinationName":"Porte de Clignancourt"}`, output)
}

func Test_RewriteValues_SkipNormalMap(t *testing.T) {
	assert := assert.New(t)

	output := testRewriteValues(t, `{ "DestinationName": "Porte de Clignancourt" }`)
	assert.Equal(`{"DestinationName":"Porte de Clignancourt"}`, output)
}

func Test_RewriteValues_SkipNormalPayload(t *testing.T) {
	assert := assert.New(t)

	input := `{
		"Siri": {
			"ServiceDelivery": {
				"ProducerRef": "IVTR_HET",
				"ResponseMessageIdentifier": "IVTR_HET:ResponseMessage:17df5d24-c796-4bf6-bd4e-48ee41b19d17:LOC:",
				"ResponseTimestamp": "2023-05-31T09:27:01.999Z",
				"StopMonitoringDelivery": [
					{
						"MonitoredStopVisit": [
							{
								"ItemIdentifier": "RATP-SIV:Item::20230531.78.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoredVehicleJourney": {
									"DestinationName": "Porte de Clignancourt",
									"DestinationRef": "STIF:StopPoint:Q:22141:",
									"DirectionName": "PORTE DE CLIGNANCOURT",
									"FramedVehicleJourneyRef": {
										"DataFrameRef": "any",
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.78.R.C01374:LOC"
									},
									"JourneyNote": [],
									"LineRef": "STIF:Line::C01374:",
									"MonitoredCall": {
										"ArrivalStatus": "",
										"DepartureStatus": "onTime",
										"DestinationDisplay": "Porte de Clignancourt",
										"ExpectedArrivalTime": "2023-05-31T07:27:22.568Z",
										"ExpectedDepartureTime": "2023-05-31T07:27:22.568Z",
										"StopPointName": "Châtelet",
										"VehicleAtStop": false
									},
									"OperatorRef": "RATP-SIV:Operator::RATP.OCTAVE.4.4:",
									"TrainNumbers": {
										"TrainNumberRef": []
									},
									"VehicleJourneyName": []
								},
								"MonitoringRef": "STIF:StopPoint:Q:463158:",
								"RecordedAtTime": "2023-05-31T07:27:00.568Z"
							},
							{
								"ItemIdentifier": "RATP-SIV:Item::20230531.79.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoredVehicleJourney": {
									"DestinationName": "Porte de Clignancourt",
									"DestinationRef": "STIF:StopPoint:Q:22141:",
									"DirectionName": "PORTE DE CLIGNANCOURT",
									"FramedVehicleJourneyRef": {
										"DataFrameRef": "any",
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.79.R.C01374:LOC"
									},
									"JourneyNote": [],
									"LineRef": "STIF:Line::C01374:",
									"MonitoredCall": {
										"ArrivalStatus": "",
										"DepartureStatus": "onTime",
										"DestinationDisplay": "Porte de Clignancourt",
										"ExpectedArrivalTime": "2023-05-31T07:29:14.786Z",
										"ExpectedDepartureTime": "2023-05-31T07:29:14.786Z",
										"StopPointName": "Châtelet",
										"VehicleAtStop": false
									},
									"OperatorRef": "RATP-SIV:Operator::RATP.OCTAVE.4.4:",
									"TrainNumbers": {
										"TrainNumberRef": []
									},
									"VehicleJourneyName": []
								},
								"MonitoringRef": "STIF:StopPoint:Q:463158:",
								"RecordedAtTime": "2023-05-31T07:22:14.786Z"
							},
							{
								"ItemIdentifier": "RATP-SIV:Item::20230531.80.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoredVehicleJourney": {
									"DestinationName": "Porte de Clignancourt",
									"DestinationRef": "STIF:StopPoint:Q:22141:",
									"DirectionName": "PORTE DE CLIGNANCOURT",
									"FramedVehicleJourneyRef": {
										"DataFrameRef": "any",
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.80.R.C01374:LOC"
									},
									"JourneyNote": [],
									"LineRef": "STIF:Line::C01374:",
									"MonitoredCall": {
										"ArrivalStatus": "",
										"DepartureStatus": "onTime",
										"DestinationDisplay": "Porte de Clignancourt",
										"ExpectedArrivalTime": "2023-05-31T07:31:03.077Z",
										"ExpectedDepartureTime": "2023-05-31T07:31:03.077Z",
										"StopPointName": "Châtelet",
										"VehicleAtStop": false
									},
									"OperatorRef": "RATP-SIV:Operator::RATP.OCTAVE.4.4:",
									"TrainNumbers": {
										"TrainNumberRef": []
									},
									"VehicleJourneyName": []
								},
								"MonitoringRef": "STIF:StopPoint:Q:463158:",
								"RecordedAtTime": "2023-05-31T07:25:03.077Z"
							},
							{
								"ItemIdentifier": "RATP-SIV:Item::20230531.81.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoredVehicleJourney": {
									"DestinationName": "Porte de Clignancourt",
									"DestinationRef": "STIF:StopPoint:Q:22141:",
									"DirectionName": "PORTE DE CLIGNANCOURT",
									"FramedVehicleJourneyRef": {
										"DataFrameRef": "any",
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.81.R.C01374:LOC"
									},
									"JourneyNote": [],
									"LineRef": "STIF:Line::C01374:",
									"MonitoredCall": {
										"ArrivalStatus": "",
										"DepartureStatus": "onTime",
										"DestinationDisplay": "Porte de Clignancourt",
										"ExpectedArrivalTime": "2023-05-31T07:32:33.350Z",
										"ExpectedDepartureTime": "2023-05-31T07:32:33.350Z",
										"StopPointName": "Châtelet",
										"VehicleAtStop": false
									},
									"OperatorRef": "RATP-SIV:Operator::RATP.OCTAVE.4.4:",
									"TrainNumbers": {
										"TrainNumberRef": []
									},
									"VehicleJourneyName": []
								},
								"MonitoringRef": "STIF:StopPoint:Q:463158:",
								"RecordedAtTime": "2023-05-31T07:26:33.350Z"
							}
						],
						"ResponseTimestamp": "2023-05-31T09:27:02.012Z",
						"Status": "true",
						"Version": "2.0"
					}
				]
			}
		}
	}`

	output := testRewriteValues(t, input)
	assert.Equal(testCompactedJSON(input), output)
}

func Test_RewriteValues_RewriteWholePayload(t *testing.T) {
	assert := assert.New(t)

	input := `{
		"Siri": {
			"ServiceDelivery": {
				"ResponseTimestamp": "2023-05-31T09:27:01.999Z",
				"ProducerRef": "IVTR_HET",
				"ResponseMessageIdentifier": "IVTR_HET:ResponseMessage:17df5d24-c796-4bf6-bd4e-48ee41b19d17:LOC:",
				"StopMonitoringDelivery": [
					{
						"ResponseTimestamp": "2023-05-31T09:27:02.012Z",
						"Version": "2.0",
						"Status": "true",
						"MonitoredStopVisit": [
							{
								"RecordedAtTime": "2023-05-31T07:27:00.568Z",
								"ItemIdentifier": "RATP-SIV:Item::20230531.78.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoringRef": {
									"value": "STIF:StopPoint:Q:463158:"
								},
								"MonitoredVehicleJourney": {
									"LineRef": {
										"value": "STIF:Line::C01374:"
									},
									"OperatorRef": {
										"value": "RATP-SIV:Operator::RATP.OCTAVE.4.4:"
									},
									"FramedVehicleJourneyRef": {
										"DataFrameRef": {
											"value": "any"
										},
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.78.R.C01374:LOC"
									},
									"DirectionName": [
										{
											"value": "PORTE DE CLIGNANCOURT"
										}
									],
									"DestinationRef": {
										"value": "STIF:StopPoint:Q:22141:"
									},
									"DestinationName": [
										{
											"value": "Porte de Clignancourt"
										}
									],
									"VehicleJourneyName": [],
									"JourneyNote": [],
									"MonitoredCall": {
										"StopPointName": [
											{
												"value": "Châtelet"
											}
										],
										"VehicleAtStop": false,
										"DestinationDisplay": [
											{
												"value": "Porte de Clignancourt"
											}
										],
										"ExpectedArrivalTime": "2023-05-31T07:27:22.568Z",
										"ExpectedDepartureTime": "2023-05-31T07:27:22.568Z",
										"DepartureStatus": "onTime",
										"ArrivalStatus": ""
									},
									"TrainNumbers": {
										"TrainNumberRef": []
									}
								}
							},
							{
								"RecordedAtTime": "2023-05-31T07:22:14.786Z",
								"ItemIdentifier": "RATP-SIV:Item::20230531.79.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoringRef": {
									"value": "STIF:StopPoint:Q:463158:"
								},
								"MonitoredVehicleJourney": {
									"LineRef": {
										"value": "STIF:Line::C01374:"
									},
									"OperatorRef": {
										"value": "RATP-SIV:Operator::RATP.OCTAVE.4.4:"
									},
									"FramedVehicleJourneyRef": {
										"DataFrameRef": {
											"value": "any"
										},
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.79.R.C01374:LOC"
									},
									"DirectionName": [
										{
											"value": "PORTE DE CLIGNANCOURT"
										}
									],
									"DestinationRef": {
										"value": "STIF:StopPoint:Q:22141:"
									},
									"DestinationName": [
										{
											"value": "Porte de Clignancourt"
										}
									],
									"VehicleJourneyName": [],
									"JourneyNote": [],
									"MonitoredCall": {
										"StopPointName": [
											{
												"value": "Châtelet"
											}
										],
										"VehicleAtStop": false,
										"DestinationDisplay": [
											{
												"value": "Porte de Clignancourt"
											}
										],
										"ExpectedArrivalTime": "2023-05-31T07:29:14.786Z",
										"ExpectedDepartureTime": "2023-05-31T07:29:14.786Z",
										"DepartureStatus": "onTime",
										"ArrivalStatus": ""
									},
									"TrainNumbers": {
										"TrainNumberRef": []
									}
								}
							},
							{
								"RecordedAtTime": "2023-05-31T07:25:03.077Z",
								"ItemIdentifier": "RATP-SIV:Item::20230531.80.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoringRef": {
									"value": "STIF:StopPoint:Q:463158:"
								},
								"MonitoredVehicleJourney": {
									"LineRef": {
										"value": "STIF:Line::C01374:"
									},
									"OperatorRef": {
										"value": "RATP-SIV:Operator::RATP.OCTAVE.4.4:"
									},
									"FramedVehicleJourneyRef": {
										"DataFrameRef": {
											"value": "any"
										},
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.80.R.C01374:LOC"
									},
									"DirectionName": [
										{
											"value": "PORTE DE CLIGNANCOURT"
										}
									],
									"DestinationRef": {
										"value": "STIF:StopPoint:Q:22141:"
									},
									"DestinationName": [
										{
											"value": "Porte de Clignancourt"
										}
									],
									"VehicleJourneyName": [],
									"JourneyNote": [],
									"MonitoredCall": {
										"StopPointName": [
											{
												"value": "Châtelet"
											}
										],
										"VehicleAtStop": false,
										"DestinationDisplay": [
											{
												"value": "Porte de Clignancourt"
											}
										],
										"ExpectedArrivalTime": "2023-05-31T07:31:03.077Z",
										"ExpectedDepartureTime": "2023-05-31T07:31:03.077Z",
										"DepartureStatus": "onTime",
										"ArrivalStatus": ""
									},
									"TrainNumbers": {
										"TrainNumberRef": []
									}
								}
							},
							{
								"RecordedAtTime": "2023-05-31T07:26:33.350Z",
								"ItemIdentifier": "RATP-SIV:Item::20230531.81.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoringRef": {
									"value": "STIF:StopPoint:Q:463158:"
								},
								"MonitoredVehicleJourney": {
									"LineRef": {
										"value": "STIF:Line::C01374:"
									},
									"OperatorRef": {
										"value": "RATP-SIV:Operator::RATP.OCTAVE.4.4:"
									},
									"FramedVehicleJourneyRef": {
										"DataFrameRef": {
											"value": "any"
										},
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.81.R.C01374:LOC"
									},
									"DirectionName": [
										{
											"value": "PORTE DE CLIGNANCOURT"
										}
									],
									"DestinationRef": {
										"value": "STIF:StopPoint:Q:22141:"
									},
									"DestinationName": [
										{
											"value": "Porte de Clignancourt"
										}
									],
									"VehicleJourneyName": [],
									"JourneyNote": [],
									"MonitoredCall": {
										"StopPointName": [
											{
												"value": "Châtelet"
											}
										],
										"VehicleAtStop": false,
										"DestinationDisplay": [
											{
												"value": "Porte de Clignancourt"
											}
										],
										"ExpectedArrivalTime": "2023-05-31T07:32:33.350Z",
										"ExpectedDepartureTime": "2023-05-31T07:32:33.350Z",
										"DepartureStatus": "onTime",
										"ArrivalStatus": ""
									},
									"TrainNumbers": {
										"TrainNumberRef": []
									}
								}
							}
						]
					}
				]
			}
		}
	}`

	output := testRewriteValues(t, input)

	expected := `{
		"Siri": {
			"ServiceDelivery": {
				"ProducerRef": "IVTR_HET",
				"ResponseMessageIdentifier": "IVTR_HET:ResponseMessage:17df5d24-c796-4bf6-bd4e-48ee41b19d17:LOC:",
				"ResponseTimestamp": "2023-05-31T09:27:01.999Z",
				"StopMonitoringDelivery": [
					{
						"MonitoredStopVisit": [
							{
								"ItemIdentifier": "RATP-SIV:Item::20230531.78.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoredVehicleJourney": {
									"DestinationName": "Porte de Clignancourt",
									"DestinationRef": "STIF:StopPoint:Q:22141:",
									"DirectionName": "PORTE DE CLIGNANCOURT",
									"FramedVehicleJourneyRef": {
										"DataFrameRef": "any",
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.78.R.C01374:LOC"
									},
									"JourneyNote": [],
									"LineRef": "STIF:Line::C01374:",
									"MonitoredCall": {
										"ArrivalStatus": "",
										"DepartureStatus": "onTime",
										"DestinationDisplay": "Porte de Clignancourt",
										"ExpectedArrivalTime": "2023-05-31T07:27:22.568Z",
										"ExpectedDepartureTime": "2023-05-31T07:27:22.568Z",
										"StopPointName": "Châtelet",
										"VehicleAtStop": false
									},
									"OperatorRef": "RATP-SIV:Operator::RATP.OCTAVE.4.4:",
									"TrainNumbers": {
										"TrainNumberRef": []
									},
									"VehicleJourneyName": []
								},
								"MonitoringRef": "STIF:StopPoint:Q:463158:",
								"RecordedAtTime": "2023-05-31T07:27:00.568Z"
							},
							{
								"ItemIdentifier": "RATP-SIV:Item::20230531.79.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoredVehicleJourney": {
									"DestinationName": "Porte de Clignancourt",
									"DestinationRef": "STIF:StopPoint:Q:22141:",
									"DirectionName": "PORTE DE CLIGNANCOURT",
									"FramedVehicleJourneyRef": {
										"DataFrameRef": "any",
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.79.R.C01374:LOC"
									},
									"JourneyNote": [],
									"LineRef": "STIF:Line::C01374:",
									"MonitoredCall": {
										"ArrivalStatus": "",
										"DepartureStatus": "onTime",
										"DestinationDisplay": "Porte de Clignancourt",
										"ExpectedArrivalTime": "2023-05-31T07:29:14.786Z",
										"ExpectedDepartureTime": "2023-05-31T07:29:14.786Z",
										"StopPointName": "Châtelet",
										"VehicleAtStop": false
									},
									"OperatorRef": "RATP-SIV:Operator::RATP.OCTAVE.4.4:",
									"TrainNumbers": {
										"TrainNumberRef": []
									},
									"VehicleJourneyName": []
								},
								"MonitoringRef": "STIF:StopPoint:Q:463158:",
								"RecordedAtTime": "2023-05-31T07:22:14.786Z"
							},
							{
								"ItemIdentifier": "RATP-SIV:Item::20230531.80.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoredVehicleJourney": {
									"DestinationName": "Porte de Clignancourt",
									"DestinationRef": "STIF:StopPoint:Q:22141:",
									"DirectionName": "PORTE DE CLIGNANCOURT",
									"FramedVehicleJourneyRef": {
										"DataFrameRef": "any",
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.80.R.C01374:LOC"
									},
									"JourneyNote": [],
									"LineRef": "STIF:Line::C01374:",
									"MonitoredCall": {
										"ArrivalStatus": "",
										"DepartureStatus": "onTime",
										"DestinationDisplay": "Porte de Clignancourt",
										"ExpectedArrivalTime": "2023-05-31T07:31:03.077Z",
										"ExpectedDepartureTime": "2023-05-31T07:31:03.077Z",
										"StopPointName": "Châtelet",
										"VehicleAtStop": false
									},
									"OperatorRef": "RATP-SIV:Operator::RATP.OCTAVE.4.4:",
									"TrainNumbers": {
										"TrainNumberRef": []
									},
									"VehicleJourneyName": []
								},
								"MonitoringRef": "STIF:StopPoint:Q:463158:",
								"RecordedAtTime": "2023-05-31T07:25:03.077Z"
							},
							{
								"ItemIdentifier": "RATP-SIV:Item::20230531.81.R.C01374.PALS.IDFM.C01374.R.RATP.50026977:LOC",
								"MonitoredVehicleJourney": {
									"DestinationName": "Porte de Clignancourt",
									"DestinationRef": "STIF:StopPoint:Q:22141:",
									"DirectionName": "PORTE DE CLIGNANCOURT",
									"FramedVehicleJourneyRef": {
										"DataFrameRef": "any",
										"DatedVehicleJourneyRef": "RATP-SIV:VehicleJourney::20230531.81.R.C01374:LOC"
									},
									"JourneyNote": [],
									"LineRef": "STIF:Line::C01374:",
									"MonitoredCall": {
										"ArrivalStatus": "",
										"DepartureStatus": "onTime",
										"DestinationDisplay": "Porte de Clignancourt",
										"ExpectedArrivalTime": "2023-05-31T07:32:33.350Z",
										"ExpectedDepartureTime": "2023-05-31T07:32:33.350Z",
										"StopPointName": "Châtelet",
										"VehicleAtStop": false
									},
									"OperatorRef": "RATP-SIV:Operator::RATP.OCTAVE.4.4:",
									"TrainNumbers": {
										"TrainNumberRef": []
									},
									"VehicleJourneyName": []
								},
								"MonitoringRef": "STIF:StopPoint:Q:463158:",
								"RecordedAtTime": "2023-05-31T07:26:33.350Z"
							}
						],
						"ResponseTimestamp": "2023-05-31T09:27:02.012Z",
						"Status": "true",
						"Version": "2.0"
					}
				]
			}
		}
	}`

	assert.Equal(testCompactedJSON(expected), output)
}
