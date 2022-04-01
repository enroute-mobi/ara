package remote

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"bitbucket.org/enroute-mobi/ara/siri"
)

func newSIRIEnvelopeFrom(name, envelopeType string) (*SIRIEnvelope, error) {
	file, err := os.Open(fmt.Sprintf("testdata/%s.xml", name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envelope, err := NewSIRIEnvelope(file, envelopeType)
	if err != nil {
		return nil, err
	}
	return envelope, nil
}

func Test_SOAPEnvelope_WithoutSOAP(t *testing.T) {
	_, err := newSIRIEnvelopeFrom("checkstatus-response", SOAP_SIRI_ENVELOPE)
	if err == nil {
		t.Error("Attempting to create a SOAP envelope without soap should return an error")
	}
}

// Test the persistance of request body
func Test_SOAPEnvelope_Finalizer(t *testing.T) {
	envelope, err := newSIRIEnvelopeFrom("checkstatus-response-soap", SOAP_SIRI_ENVELOPE)
	if err != nil {
		t.Fatal(err)
	}

	// Create a finalizer and a channel to be sure the finalizer has ended
	var finalized bool
	done := make(chan bool, 1)
	finalizer := func(sa *SIRIEnvelope) {
		finalized = true
		done <- true
	}
	runtime.SetFinalizer(envelope, finalizer)

	// Create a CheckStatusResponse and destroy envelope
	response := siri.NewXMLCheckStatusResponse(envelope.body)
	envelope = nil
	runtime.GC()
	<-done

	if !finalized {
		t.Errorf("SIRIEnvelope should be destroyed by GC")
	}
	if response.Node() == nil {
		t.Errorf("Xml Node shouldn't be destroyed with SIRIEnvelope")
	}
}

func Test_SOAPEnvelope_BodyType(t *testing.T) {
	envelope, err := newSIRIEnvelopeFrom("checkstatus-response-soap", SOAP_SIRI_ENVELOPE)
	if err != nil {
		t.Fatal(err)
	}

	if expected := "CheckStatusResponse"; envelope.BodyType() != expected {
		t.Errorf("Wrong BodyType:\n got: %v\n want: %v", envelope.BodyType(), expected)
	}
}

// Test the persistance of request body
func Test_RawEnvelope_Finalizer(t *testing.T) {
	envelope, err := newSIRIEnvelopeFrom("checkstatus-response", RAW_SIRI_ENVELOPE)
	if err != nil {
		t.Fatal(err)
	}

	// Create a finalizer and a channel to be sure the finalizer has ended
	var finalized bool
	done := make(chan bool, 1)
	finalizer := func(sa *SIRIEnvelope) {
		finalized = true
		done <- true
	}
	runtime.SetFinalizer(envelope, finalizer)

	// Create a CheckStatusResponse and destroy envelope
	response := siri.NewXMLCheckStatusResponse(envelope.body)
	envelope = nil
	runtime.GC()
	<-done

	if !finalized {
		t.Errorf("SIRIEnvelope should be destroyed by GC")
	}
	if response.Node() == nil {
		t.Errorf("Xml Node shouldn't be destroyed with SIRIEnvelope")
	}
}

func Test_RawEnvelope_BodyType(t *testing.T) {
	envelope, err := newSIRIEnvelopeFrom("checkstatus-response-raw", RAW_SIRI_ENVELOPE)
	if err != nil {
		t.Fatal(err)
	}

	if expected := "CheckStatusResponse"; envelope.BodyType() != expected {
		t.Errorf("Wrong BodyType:\n got: %v\n want: %v", envelope.BodyType(), expected)
	}
}
