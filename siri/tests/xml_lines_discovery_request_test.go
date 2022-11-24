package siri_tests

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func getXMLLinesDiscoveryRequest(t *testing.T) *sxml.XMLLinesDiscoveryRequest {
	file, err := os.Open("testdata/lines-discovery-request.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	request, _ := sxml.NewXMLLinesDiscoveryRequestFromContent(content)
	return request
}

func Test_XMLLinesDiscoveryRequest_RequestorRef(t *testing.T) {
	request := getXMLLinesDiscoveryRequest(t)

	if expected := "test"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef:\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}
}

func Test_XMLLinesDiscoveryRequest_MessageIdentifier(t *testing.T) {
	request := getXMLLinesDiscoveryRequest(t)

	if expected := "STIF:Message::2345Fsdfrg35df:LOC"; request.MessageIdentifier() != expected {
		t.Errorf("Wrong MessageIdentifier:\n got: %v\nwant: %v", request.MessageIdentifier(), expected)
	}
}

func Test_XMLLinesDiscoveryRequest_RequestTimestamp(t *testing.T) {
	request := getXMLLinesDiscoveryRequest(t)

	if expected := time.Date(2017, time.March, 03, 11, 28, 00, 359000000, time.UTC); request.RequestTimestamp() != expected {
		t.Errorf("Wrong RequestTimestamp:\n got: %v\nwant: %v", request.RequestTimestamp(), expected)
	}
}
