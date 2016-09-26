package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_XMLCheckStatusRequest_RequestorRef(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request := NewXMLCheckStatusRequestFromContent(content)
	if expected := "NINOXE:default"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef :\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}
}

func Test_XMLCheckStatusRequest_RequestTimestamp(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request := NewXMLCheckStatusRequestFromContent(content)
	if expected := time.Date(2016, time.September, 07, 9, 11, 25, 174000000, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp :\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}
}

func Test_XMLCheckStatusRequest_MessageIdentifier(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request := NewXMLCheckStatusRequestFromContent(content)
	if expected := "CheckStatus:Test:0"; request.MessageIdentifier() != expected {
		t.Errorf("Wrong MessageIdentifier :\n got: %v\nwant: %v", request.MessageIdentifier(), expected)
	}
}

func Test_SIRICheckStatusRequest_BuildXML(t *testing.T) {
	expectedXML := `<ns7:CheckStatus xmlns:ns2="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://scma/siri" xmlns:ns7="http://wsdl.siri.org.uk">
	<Request>
		<ns2:RequestTimestamp>2009-11-10T23:00:00.000Z</ns2:RequestTimestamp>
		<ns2:RequestorRef>test</ns2:RequestorRef>
		<ns2:MessageIdentifier>test</ns2:MessageIdentifier>
	</Request>
	<RequestExtension/>
</ns7:CheckStatus>`
	date := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	request := NewSIRICheckStatusRequest("test", date, "test")
	if expectedXML != request.BuildXML() {
		t.Errorf("Wrong XML for Request :\n got:\n%v\nwant:\n%v", request.BuildXML(), expectedXML)
	}
}
