package remote

import (
	"bytes"
	"runtime"
	"testing"

	"bitbucket.org/enroute-mobi/ara/siri"
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
	response := siri.NewXMLCheckStatusResponse(envelope.body)
	envelope = nil
	runtime.GC()
	<-done

	if !finalized {
		t.Errorf("SOAPEnvelope should be destroyed by GC")
	}
	if response.Node() == nil {
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

func Test_SOAPEnvelopeBuffer_Read(t *testing.T) {
	var r []byte = make([]byte, 512)
	expected := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
test
</S:Body>
</S:Envelope>`

	buffer := NewSOAPEnvelopeBuffer()
	buffer.WriteXML("test")

	readLength, err := buffer.Read(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(r[:readLength]) != expected {
		t.Errorf("Incorrect Read:\n got: %v\n want: %v", string(r[:readLength]), expected)
	}
}

func Test_SOAPEnvelopeBuffer_WriteTo(t *testing.T) {
	var buf bytes.Buffer
	expected := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
test
</S:Body>
</S:Envelope>`

	buffer := NewSOAPEnvelopeBuffer()
	buffer.WriteXML("test")

	_, err := buffer.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != expected {
		t.Errorf("Incorrect WriteTo:\n got: %v\n want: %v", buf.String(), expected)
	}

}
