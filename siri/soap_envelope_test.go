package siri

import (
	"os"
	"runtime"
	"testing"
)

// Test the persistance of request body
func Test_SOAPEnvelope_Finalizer(t *testing.T) {
	// Create a new SOAPEnvelope
	file, err := os.Open("testdata/checkstatus-soap-response.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	envelope, err := NewSOAPEnvelope(file)
	if err != nil {
		t.Fatal(err)
	}

	// Create a finalizer and a channel to be sure the finalizer has ended
	var finalized bool
	done := make(chan bool, 1)
	finalizer := func(sa *SOAPEnvelope) {
		finalized = true
		done <- true
	}
	runtime.SetFinalizer(envelope, finalizer)

	// Create a CheckStatusResponse and destroy envelope
	response := NewXMLCheckStatusResponse(envelope.body)
	envelope = nil
	runtime.GC()
	<-done

	if !finalized {
		t.Errorf("SOAPEnvelope should be destroyed by GC")
	}
	if response.node == nil {
		t.Errorf("Xml Node shouldn't be destroyed with SOAPEnvelope")
	}
}

func Test_SOAPEnvelope_BodyType(t *testing.T) {
	file, err := os.Open("testdata/checkstatus-soap-response.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	envelope, err := NewSOAPEnvelope(file)
	if err != nil {
		t.Fatal(err)
	}

	if expected := "CheckStatusResponse"; envelope.BodyType() != expected {
		t.Errorf("Wrong BodyType:\n got: %v\n want: %v", envelope.BodyType(), expected)
	}
}
