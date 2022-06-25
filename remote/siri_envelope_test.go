package remote

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"

	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/stretchr/testify/assert"
)

func Test_NewAutodetectSIRIEnvelope(t *testing.T) {
	assert := assert.New(t)

	// empty request
	req := strings.NewReader("")
	_, err := NewAutodetectSIRIEnvelope(req)
	expectedError := errors.New("empty body")

	assert.Equal(expectedError, err)

	// Invalid XML
	_, err1 := newAutodetectSIRIEnvelopeFrom("invalid-xml")
	expectedError1 := errors.New("failed to parse xml input")

	assert.Equal(expectedError1, err1)

	// SOAP request
	// when the request has a `Body` tag
	env2, err2 := newAutodetectSIRIEnvelopeFrom("checkstatus-request-soap")
	if err2 != nil {
		t.Errorf("cannot detect SIRI envelope: %s", err2)
	}
	expected2 := `<sw:CheckStatus xmlns:siri="http://www.siri.org.uk/siri" xmlns:ns3="http://www.ifopt.org.uk/acsb" xmlns:ns4="http://www.ifopt.org.uk/ifopt" xmlns:ns5="http://datex2.eu/schema/2_0RC1/2_0" xmlns:ns6="http://scma/siri" xmlns:sw="http://wsdl.siri.org.uk">
      <Request>
        <siri:RequestTimestamp>2017-01-01T12:00:00.000Z</siri:RequestTimestamp>
        <siri:RequestorRef>test2</siri:RequestorRef>
        <siri:MessageIdentifier>RATPDev:ResponseMessage::d3f94aa2-7b76-449b-aa18-50caf78f9dc7:LOC</siri:MessageIdentifier>
      </Request>
      <RequestExtension/>
    </sw:CheckStatus>`

	assert.Equal(expected2, env2.body.String())

	// raw request
	// when the request has no body
	_, err3 := newAutodetectSIRIEnvelopeFrom("checkstatus-raw-malformed")
	expectedError3 := errors.New("invalid raw xml: cannot find body")
	assert.Equal(expectedError3, err3)

	// raw request
	// when the request is valid with no special character between tags
	env4, err4 := newAutodetectSIRIEnvelopeFrom("checkstatus-raw-no-special-character")
	if err4 != nil {
		t.Errorf("cannot detect SIRI envelope: %s", err4)
	}
	expected4 := `<CheckStatusRequest>
  <RequestTimestamp>2017-01-01T12:00:00.000Z</RequestTimestamp>
  <RequestorRef>Ara</RequestorRef>
  <MessageIdentifier>6ba7b814-9dad-11d1-2-00c04fd430c8</MessageIdentifier>
</CheckStatusRequest>`

	assert.Equal(expected4, env4.body.String())

	// raw request
	// when the request is valid
	env5, err5 := newAutodetectSIRIEnvelopeFrom("checkstatus-request-raw")
	if err5 != nil {
		t.Errorf("cannot detect SIRI envelope: %s", err5)
	}
	expected5 := `<CheckStatusRequest>
      <RequestTimestamp>2017-01-01T12:00:30.000Z</RequestTimestamp>
      <RequestorRef>Ara</RequestorRef>
      <MessageIdentifier>6ba7b814-9dad-11d1-2-00c04fd430c8</MessageIdentifier>
  </CheckStatusRequest>`

	assert.Equal(expected5, env5.body.String())
}

func newAutodetectSIRIEnvelopeFrom(name string) (*SIRIEnvelope, error) {
	file, err := os.Open(fmt.Sprintf("testdata/%s.xml", name))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	envelope, err := NewAutodetectSIRIEnvelope(file)

	if err != nil {
		return nil, err
	}
	return envelope, nil
}
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
	response := sxml.NewXMLCheckStatusResponse(envelope.body)
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
	response := sxml.NewXMLCheckStatusResponse(envelope.body)
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
