package remote

import (
	"bytes"
	"testing"
)

func Test_SIRIBuffer_Read(t *testing.T) {
	var r []byte = make([]byte, 512)
	expected := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
test
</S:Body>
</S:Envelope>`

	buffer := NewSIRIBuffer()
	buffer.WriteXML("test")

	readLength, err := buffer.Read(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(r[:readLength]) != expected {
		t.Errorf("Incorrect Read:\n got: %v\n want: %v", string(r[:readLength]), expected)
	}

	buffer = NewSIRIBuffer(RAW_SIRI_ENVELOPE)
	buffer.WriteXML("test")
	readLength, err = buffer.Read(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(r[:readLength]) != "test" {
		t.Errorf("Incorrect Read:\n got: %v\n want: test", string(r[:readLength]))
	}

}

func Test_SIRIBuffer_WriteTo(t *testing.T) {
	var buf bytes.Buffer
	expected := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
test
</S:Body>
</S:Envelope>`

	buffer := NewSIRIBuffer()
	buffer.WriteXML("test")

	_, err := buffer.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != expected {
		t.Errorf("Incorrect WriteTo:\n got: %v\n want: %v", buf.String(), expected)
	}

	buf = *new(bytes.Buffer)
	buffer = NewSIRIBuffer(RAW_SIRI_ENVELOPE)
	buffer.WriteXML("test")

	_, err = buffer.WriteTo(&buf)
	if err != nil {
		t.Fatal(err)
	}
	if buf.String() != "test" {
		t.Errorf("Incorrect WriteTo:\n got: %v\n want: test", buf.String())
	}

}
