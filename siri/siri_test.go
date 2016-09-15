package siri

import (
	"io/ioutil"
	"os"
	"testing"
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

	request := NewXMLCheckStatusRequest(content)
	if expected := "NINOXE:default"; request.RequestorRef() != expected {
		t.Errorf("Wrong RequestorRef :\n got: %v\nwant: %v", request.RequestorRef(), expected)
	}
}
