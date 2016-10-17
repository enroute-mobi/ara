package siri

import (
	"runtime"
	"testing"
)

func testSOAPEnvelope(name string) (*SOAPEnvelope, error) {
	// Create a new SOAPEnvelope
	file, err := testSOAPFile(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envelope, err := NewSOAPEnvelope(file)
	if err != nil {
		return nil, err
	}
	return envelope, nil
}

// Test the persistance of request body
func Test_SOAPEnvelope_Finalizer(t *testing.T) {
	envelope, err := testSOAPEnvelope("checkstatus-response")
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
	envelope, err := testSOAPEnvelope("checkstatus-response")
	if err != nil {
		t.Fatal(err)
	}

	if expected := "CheckStatusResponse"; envelope.BodyType() != expected {
		t.Errorf("Wrong BodyType:\n got: %v\n want: %v", envelope.BodyType(), expected)
	}
}
