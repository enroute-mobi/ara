package remote

import (
	"bytes"
	"testing"
)

func Test_SIRIBuffer_Read(t *testing.T) {
	r := make([]byte, 512)
	expectedSOAP := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
test
</S:Body>
</S:Envelope>`

	expectedRaw := `<?xml version='1.0' encoding='utf-8'?>
test`

	buffer := NewSIRIBuffer(SOAP_SIRI_ENVELOPE)
	buffer.WriteXML("test")

	readLength, err := buffer.Read(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(r[:readLength]) != expectedSOAP {
		t.Errorf("Incorrect Read:\n got: %v\n want: %v", string(r[:readLength]), expectedSOAP)
	}

	buffer = NewSIRIBuffer(RAW_SIRI_ENVELOPE)
	buffer.WriteXML("test")
	readLength, err = buffer.Read(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(r[:readLength]) != expectedRaw {
		t.Errorf("Incorrect Read:\n got: %v\n want: %v", string(r[:readLength]), expectedRaw)
	}

}

func Test_SIRIBuffer_WriteTo(t *testing.T) {
	var buf bytes.Buffer
	expectedSOAP := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
test
</S:Body>
</S:Envelope>`

	expectedRaw := `<?xml version='1.0' encoding='utf-8'?>
test`

	buffer := NewSIRIBuffer(SOAP_SIRI_ENVELOPE)
	buffer.WriteXML("test")

	_, err := buffer.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != expectedSOAP {
		t.Errorf("Incorrect WriteTo:\n got: %v\n want: %v", buf.String(), expectedSOAP)
	}

	buf = *new(bytes.Buffer)
	buffer = NewSIRIBuffer(RAW_SIRI_ENVELOPE)
	buffer.WriteXML("test")

	_, err = buffer.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != expectedRaw {
		t.Errorf("Incorrect WriteTo:\n got: %v\n want: test", buf.String())
	}

}
