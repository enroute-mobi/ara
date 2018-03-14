package siri

import (
	"testing"
	"time"
)

func Test_SIRILinesDiscoveryResponse_BuildXML(t *testing.T) {
	expectedXML := `<sw:LinesDiscoveryResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Answer version="2.0">
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:Status>true</siri:Status>
		<siri:AnnotatedLineRef>
			<siri:LineRef>NINOXE:Line:BP:1:LOC</siri:LineRef>
			<siri:LineName>Test 1</siri:LineName>
			<siri:Monitored>true</siri:Monitored>
		</siri:AnnotatedLineRef>
		<siri:AnnotatedLineRef>
			<siri:LineRef>NINOXE:Line:BP:2:LOC</siri:LineRef>
			<siri:LineName>Test 2</siri:LineName>
			<siri:Monitored>true</siri:Monitored>
		</siri:AnnotatedLineRef>
		<siri:AnnotatedLineRef>
			<siri:LineRef>NINOXE:Line:BP:3:LOC</siri:LineRef>
			<siri:LineName>Test 3</siri:LineName>
			<siri:Monitored>true</siri:Monitored>
		</siri:AnnotatedLineRef>
	</Answer>
	<AnswerExtension/>
</sw:LinesDiscoveryResponse>`

	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	response := &SIRILinesDiscoveryResponse{
		Status:            true,
		ResponseTimestamp: responseTimestamp,
	}

	annoted1 := &SIRIAnnotatedLine{
		LineRef:   "NINOXE:Line:BP:1:LOC",
		LineName:  "Test 1",
		Monitored: true,
	}
	annoted2 := &SIRIAnnotatedLine{
		LineRef:   "NINOXE:Line:BP:2:LOC",
		LineName:  "Test 2",
		Monitored: true,
	}
	annoted3 := &SIRIAnnotatedLine{
		LineRef:   "NINOXE:Line:BP:3:LOC",
		LineName:  "Test 3",
		Monitored: true,
	}

	response.AnnotatedLines = []*SIRIAnnotatedLine{annoted1, annoted2, annoted3}

	xml, err := response.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:%v\nwant:%v", xml, expectedXML)
	}
}
