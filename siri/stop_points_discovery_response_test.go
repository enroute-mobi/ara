package siri

import (
	"testing"
	"time"
)

func Test_SIRIStopDiscoveryResponse_BuildXML(t *testing.T) {
	expectedXML := `
<ns8:StopPointsDiscoveryResponse xmlns:ns8="http://wsdl.siri.org.uk" xmlns:ns3="http://www.siri.org.uk/siri" xmlns:ns4="http://www.ifopt.org.uk/acsb" xmlns:ns5="http://www.ifopt.org.uk/ifopt" xmlns:ns6="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns7="http://scma/siri" xmlns:ns9="http://wsdl.siri.org.uk/siri">
	<Answer version="2.0">
		<ns3:ResponseTimestamp>2016-09-21T20:14:46.000Z</ns3:ResponseTimestamp>
		<ns3:Address>address</ns3:Address>
		<ns3:ProducerRef>producer</ns3:ProducerRef>
		<ns3:RequestMessageRef>ref</ns3:RequestMessageRef>
		<ns3:ResponseMessageIdentifier>identifier</ns3:ResponseMessageIdentifier>
		<ns3:Status>true</ns3:Status>
		<ns3:AnnotatedStopPointRef>
			<ns3:StopPointRef>NINOXE:StopPoint:BP:1:LOC</ns3:StopPointRef>
			<ns3:Monitored>true</ns3:Monitored>
			<ns3:StopName>Test 1</ns3:StopName>
			<ns3:Lines>
				<ns3:LineRef>STIF:Line::C00272:</ns3:LineRef>
			</ns3:Lines>
		</ns3:AnnotatedStopPointRef>
		<ns3:AnnotatedStopPointRef>
			<ns3:StopPointRef>NINOXE:StopPoint:BP:2:LOC</ns3:StopPointRef>
			<ns3:Monitored>true</ns3:Monitored>
			<ns3:StopName>Test 2</ns3:StopName>
		</ns3:AnnotatedStopPointRef>
		<ns3:AnnotatedStopPointRef>
			<ns3:StopPointRef>NINOXE:StopPoint:BP:3:LOC</ns3:StopPointRef>
			<ns3:Monitored>true</ns3:Monitored>
			<ns3:StopName>Test 3</ns3:StopName>
		</ns3:AnnotatedStopPointRef>
	</Answer>
	<AnswerExtension />
</ns8:StopPointsDiscoveryResponse>`

	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	request := &SIRIStopPointsDiscoveryResponse{
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

	request.AnnotatedStopPoints = []*SIRIAnnotatedStopPoint{annoted1, annoted2, annoted3}

	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:%v\nwant:%v", xml, expectedXML)
	}
}
