package siri

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func getXMLCheckStatusResponse(t *testing.T) *XMLCheckStatusResponse {
	file, err := os.Open("testdata/checkstatus_response.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := NewXMLCheckStatusResponseFromContent(content)
	return response
}

func Test_XMLCheckStatusRequest_Address(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := "http://appli.chouette.mobi/siri_france/siri"; response.Address() != expected {
		t.Errorf("Wrong Address:\n got: %v\nwant: %v", response.Address(), expected)
	}
}

func Test_XMLCheckStatusRequest_ProducerRef(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := "NINOXE:default"; response.ProducerRef() != expected {
		t.Errorf("Wrong ProducerRef:\n got: %v\nwant: %v", response.ProducerRef(), expected)
	}
}

func Test_XMLCheckStatusRequest_RequestMessageRef(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := "CheckStatus:Test:0"; response.RequestMessageRef() != expected {
		t.Errorf("Wrong RequestMessageRef:\n got: %v\nwant: %v", response.RequestMessageRef(), expected)
	}
}

func Test_XMLCheckStatusRequest_ResponseMessageIdentifier(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := "e7a062c5-eb95-4e4e-bc4f-6792fa008c23"; response.ResponseMessageIdentifier() != expected {
		t.Errorf("Wrong ResponseMessageIdentifier:\n got: %v\nwant: %v", response.ResponseMessageIdentifier(), expected)
	}
}

func Test_XMLCheckStatusRequest_Status(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if response.Status() {
		t.Errorf("Wrong Status:\n got: %v\nwant: false", response.Status())
	}
}

func Test_XMLCheckStatusRequest_ErrorType(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := "OtherError"; response.ErrorType() != expected {
		t.Errorf("Wrong Error type:\n got: %v\nwant: %v", response.ErrorType(), expected)
	}

	if expected := 103; response.ErrorNumber() != expected {
		t.Errorf("Wrong Error number when calling ErrorType:\n got: %v\nwant: %v", response.ErrorNumber(), expected)
	}

	if expected := "UNAVAILABLE"; response.ErrorText() != expected {
		t.Errorf("Wrong Error text when calling ErrorType:\n got: %v\nwant: %v", response.ErrorText(), expected)
	}
}

func Test_XMLCheckStatusRequest_ErrorNumber(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := 103; response.ErrorNumber() != expected {
		t.Errorf("Wrong Error number:\n got: %v\nwant: %v", response.ErrorNumber(), expected)
	}
}

func Test_XMLCheckStatusRequest_ErrorText(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := "UNAVAILABLE"; response.ErrorText() != expected {
		t.Errorf("Wrong Error text:\n got: %v\nwant: %v", response.ErrorText(), expected)
	}
}

func Test_XMLCheckStatusRequest_ResponseTimestamp(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := time.Date(2016, time.September, 21, 18, 14, 46, 238000000, time.UTC); !response.ResponseTimestamp().Equal(expected) {
		t.Errorf("Wrong ResponseTimestamp:\n got: %v\nwant: %v", response.ResponseTimestamp(), expected)
	}
}

func Test_XMLCheckStatusRequest_ServiceStartedTime(t *testing.T) {
	response := getXMLCheckStatusResponse(t)
	if expected := time.Date(2016, time.September, 21, 1, 30, 22, 996000000, time.UTC); !response.ServiceStartedTime().Equal(expected) {
		t.Errorf("Wrong ServiceStartedTime:\n got: %v\nwant: %v", response.ServiceStartedTime(), expected)
	}
}

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
	request := NewSIRICheckStatusResponse("address", "producer", "ref", "identifier", false, "OtherError", 103, "text", responseTimestamp, serviceStartedTime)
	xml, err := request.BuildXML()
	if err != nil {
		t.Fatal(err)
	}
	if expectedXML != xml {
		t.Errorf("Wrong XML for Request:\n got:\n%v\nwant:\n%v", xml, expectedXML)
	}
}

func BenchmarkParseResponse(b *testing.B) {
	file, err := os.Open("testdata/checkstatus_response.xml")
	if err != nil {
		b.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		r, _ := NewXMLCheckStatusResponseFromContent(content)
		r.Address()
		r.ProducerRef()
		r.RequestMessageRef()
		r.ResponseMessageIdentifier()
		r.Status()
		r.ErrorType()
		r.ErrorNumber()
		r.ErrorText()
		r.ResponseTimestamp()
		r.ServiceStartedTime()
	}
}

func BenchmarkGenerateResponse(b *testing.B) {
	responseTimestamp := time.Date(2016, time.September, 21, 20, 14, 46, 0, time.UTC)
	serviceStartedTime := time.Date(2016, time.September, 21, 3, 30, 22, 0, time.UTC)

	for n := 0; n < b.N; n++ {
		r := NewSIRICheckStatusResponse("address", "producer", "ref", "identifier", false, "error", 103, "text", responseTimestamp, serviceStartedTime)
		r.BuildXML()
	}
}
