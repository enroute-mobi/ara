package siri_tests

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func Test_SIRICheckStatusRequest_BuildXML(t *testing.T) {
	expectedXML := `<sw:CheckStatus xmlns:sw="http://wsdl.siri.org.uk" xmlns:siri="http://www.siri.org.uk/siri">
	<Request>
		<siri:RequestTimestamp>2009-11-10T23:00:00.000Z</siri:RequestTimestamp>
		<siri:RequestorRef>test</siri:RequestorRef>
		<siri:MessageIdentifier>test</siri:MessageIdentifier>
	</Request>
	<RequestExtension/>
</sw:CheckStatus>`
	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := siri.NewSIRICheckStatusRequest("test", date, "test")
	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}

func BenchmarkParseRequest(b *testing.B) {
	file, err := os.Open("testdata/checkstatus_request.xml")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		b.Fatal(err)
	}
	for n := 0; n < b.N; n++ {
		r, _ := sxml.NewXMLCheckStatusRequestFromContent(content)
		r.MessageIdentifier()
		r.RequestorRef()
		r.RequestTimestamp()
	}
}

func BenchmarkGenerateRequest(b *testing.B) {
	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)

	for n := 0; n < b.N; n++ {
		r := siri.NewSIRICheckStatusRequest("test", date, "test")
		r.BuildXML()
	}
}
