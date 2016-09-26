package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func Test_XMLCheckStatusRequest_ProducerRef(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_response.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response := NewXMLCheckStatusResponseFromContent(content)
	if expected := "NINOXE:default"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef :\n got: %v\nwant: %v", response.ProducerRef(), expected)
	}
}

func Test_XMLCheckStatusRequest_RequestMessageRef(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_response.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response := NewXMLCheckStatusResponseFromContent(content)
	if expected := "CheckStatus:Test:0"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef :\n got: %v\nwant: %v", response.RequestMessageRef(), expected)
	}
}

func Test_XMLCheckStatusRequest_ResponseMessageIdentifier(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_response.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response := NewXMLCheckStatusResponseFromContent(content)
	if expected := "e7a062c5-eb95-4e4e-bc4f-6792fa008c23"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Wrong ResponseMessageIdentifier :\n got: %v\nwant: %v", response.ResponseMessageIdentifier(), expected)
	}
}

func Test_XMLCheckStatusRequest_Status(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_response.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response := NewXMLCheckStatusResponseFromContent(content)
	if expected := true; response.Status() != expected {
		t.Errorf("Wrong Status :\n got: %v\nwant: %v", response.Status(), expected)
	}
}

func Test_XMLCheckStatusRequest_ResponseTimestamp(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_response.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response := NewXMLCheckStatusResponseFromContent(content)
	if expected := time.Date(2016, time.September, 21, 18, 14, 46, 238000000, time.UTC); !response.ResponseTimestamp().Equal(expected) {
		t.Errorf("Wrong ResponseTimestamp :\n got: %v\nwant: %v", response.ResponseTimestamp(), expected)
	}
}

func Test_XMLCheckStatusRequest_ServiceStartedTime(t *testing.T) {
	file, err := os.Open("testdata/checkstatus_response.xml")
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response := NewXMLCheckStatusResponseFromContent(content)
	if expected := time.Date(2016, time.September, 21, 1, 30, 22, 996000000, time.UTC); !response.ServiceStartedTime().Equal(expected) {
		t.Errorf("Wrong ServiceStartedTime :\n got: %v\nwant: %v", response.ServiceStartedTime(), expected)
	}
}

func Test_SIRICheckStatusResponse_BuildXML(t *testing.T) {
	expectedXML := `<ns7:CheckStatusResponse xmlns:ns2="http://www.siri.org.uk/siri"
												 xmlns:ns3="http://www.ifopt.org.uk/acsb"
												 xmlns:ns4="http://www.ifopt.org.uk/ifopt"
												 xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0"
												 xmlns:ns6="http://scma/siri"
												 xmlns:ns7="http://wsdl.siri.org.uk"
												 xmlns:ns8="http://wsdl.siri.org.uk/siri">
	<CheckStatusAnswerInfo>
		<ns2:ResponseTimestamp>
		2016-09-21T20:14:46.000Z</ns2:ResponseTimestamp>
		<ns2:ProducerRef>test</ns2:ProducerRef>
		<ns2:Address>
		http://appli.chouette.mobi/siri_france/siri</ns2:Address>
		<ns2:ResponseMessageIdentifier>
		test</ns2:ResponseMessageIdentifier>
		<ns2:RequestMessageRef>
		test</ns2:RequestMessageRef>
	</CheckStatusAnswerInfo>
	<Answer>
		<ns2:Status>true</ns2:Status>
		<ns2:ServiceStartedTime>
		2016-09-21T03:30:22.000Z</ns2:ServiceStartedTime>
	</Answer>
	<AnswerExtension />
</ns7:CheckStatusResponse>`
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	serviceStartedTime := time.Date(2016, time.September, 21, 3, 30, 22, 0, time.UTC)
	request := NewSIRICheckStatusResponse("test", "test", "test", true, responseTimestamp, serviceStartedTime)
	if expectedXML != request.BuildXML() {
		t.Errorf("Wrong XML for Request :\n got:\n%v\nwant:\n%v", request.BuildXML(), expectedXML)
	}
}
