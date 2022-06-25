package siri_tests

import (
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
)

func Test_SIRICheckStatusResponse_BuildXML(t *testing.T) {
	expectedXML := `<sw:CheckStatusResponse xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<CheckStatusAnswerInfo>
		<siri:ResponseTimestamp>2016-09-21T20:14:46.000Z</siri:ResponseTimestamp>
		<siri:ProducerRef>producer</siri:ProducerRef>
		<siri:Address>address</siri:Address>
		<siri:ResponseMessageIdentifier>identifier</siri:ResponseMessageIdentifier>
		<siri:RequestMessageRef>ref</siri:RequestMessageRef>
	</CheckStatusAnswerInfo>
	<Answer>
		<siri:Status>false</siri:Status>
		<siri:ErrorCondition>
			<siri:OtherError number="103">
				<siri:ErrorText>text</siri:ErrorText>
			</siri:OtherError>
		</siri:ErrorCondition>
		<siri:ServiceStartedTime>2016-09-21T03:30:22.000Z</siri:ServiceStartedTime>
	</Answer>
	<AnswerExtension/>
</sw:CheckStatusResponse>`
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	serviceStartedTime := time.Date(2016, time.September, 21, 3, 30, 22, 0, time.UTC)
	request := siri.NewSIRICheckStatusResponse("address", "producer", "ref", "identifier", false, "OtherError", 103, "text", responseTimestamp, serviceStartedTime)
	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}

func BenchmarkGenerateResponse(b *testing.B) {
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	serviceStartedTime := time.Date(2016, time.September, 21, 3, 30, 22, 0, time.UTC)

	for n := 0; n < b.N; n++ {
		r := siri.NewSIRICheckStatusResponse("address", "producer", "ref", "identifier", false, "error", 103, "text", responseTimestamp, serviceStartedTime)
		r.BuildXML()
	}
}
