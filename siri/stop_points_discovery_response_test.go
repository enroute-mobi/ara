package siri

import (
	"testing"
	"time"
)

func Test_SIRIStopDiscoveryResponse_BuildXML(t *testing.T) {
	expectedXML := `<sw:StopPointsDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Answer version="2.0">
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:Address>address</siri:Address>
		<siri:ProducerRef>producer</siri:ProducerRef>
		<siri:RequestMessageRef>ref</siri:RequestMessageRef>
		<siri:ResponseMessageIdentifier>identifier</siri:ResponseMessageIdentifier>
		<siri:Status>true</siri:Status>
		<siri:AnnotatedStopPointRef>
			<siri:StopPointRef>NINOXE:StopPoint:BP:1:LOC</siri:StopPointRef>
			<siri:Monitored>true</siri:Monitored>
			<siri:StopName>Test 1</siri:StopName>
			<siri:Lines>
				<siri:LineRef>STIF:Line::C00272:</siri:LineRef>
			</siri:Lines>
		</siri:AnnotatedStopPointRef>
		<siri:AnnotatedStopPointRef>
			<siri:StopPointRef>NINOXE:StopPoint:BP:2:LOC</siri:StopPointRef>
			<siri:Monitored>true</siri:Monitored>
			<siri:StopName>Test 2</siri:StopName>
		</siri:AnnotatedStopPointRef>
		<siri:AnnotatedStopPointRef>
			<siri:StopPointRef>NINOXE:StopPoint:BP:3:LOC</siri:StopPointRef>
			<siri:Monitored>true</siri:Monitored>
			<siri:StopName>Test 3</siri:StopName>
		</siri:AnnotatedStopPointRef>
	</Answer>
	<AnswerExtension />
</sw:StopPointsDiscoveryResponse>`

	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	response := &SIRIStopPointsDiscoveryResponse{
		Address:                   "address",
		ProducerRef:               "producer",
		RequestMessageRef:         "ref",
		ResponseMessageIdentifier: "identifier",
		Status:            true,
		ResponseTimestamp: responseTimestamp,
	}

	annoted1 := &SIRIAnnotatedStopPoint{
		StopPointRef: "NINOXE:StopPoint:BP:1:LOC",
		StopName:     "Test 1",
		Monitored:    true,
		Lines:        []string{"STIF:Line::C00272:"},
	}
	annoted2 := &SIRIAnnotatedStopPoint{
		StopPointRef: "NINOXE:StopPoint:BP:2:LOC",
		StopName:     "Test 2",
		Monitored:    true,
	}
	annoted3 := &SIRIAnnotatedStopPoint{
		StopPointRef: "NINOXE:StopPoint:BP:3:LOC",
		StopName:     "Test 3",
		Monitored:    true,
	}

	response.AnnotatedStopPoints = []*SIRIAnnotatedStopPoint{annoted1, annoted2, annoted3}

	xml, err := response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:%v\nwant:%v", xml, expectedXML)
	}
}
