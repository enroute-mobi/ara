package siri_tests

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLCheckStatusRequest(t *testing.T) *sxml.XMLCheckStatusRequest {
	file, err := os.Open("testdata/checkstatus_request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := sxml.NewXMLCheckStatusRequestFromContent(content)
	return request
}

func Test_XMLCheckStatusRequest_RequestorRef(t *testing.T) {
	request := getXMLCheckStatusRequest(t)
	if expected := "NINOXE:default"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}
}

func Test_XMLCheckStatusRequest_RequestTimestamp(t *testing.T) {
	request := getXMLCheckStatusRequest(t)
	if expected := time.Date(2016, time.September, 7, 9, 11, 25, 174000000, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}
}

func Test_XMLCheckStatusRequest_MessageIdentifier(t *testing.T) {
	request := getXMLCheckStatusRequest(t)
	if expected := "CheckStatus:Test:0"; request.MessageIdentifier() != expected {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", request.MessageIdentifier(), expected)
	}
}
