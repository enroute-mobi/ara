package siri

import "testing"

func Test_WrapSoap(t *testing.T) {
	expected := `<?xml version='1.0' encoding='utf-8'?>
<S:Envelope xmlns:S="http://schemas.xmlsoap.org/soap/envelope/" xmlns:SOAP-ENV="http://schemas.xmlsoap.org/soap/envelope/">
<S:Body>
test
</S:Body>
</S:Envelope>`
	if WrapSoap("test") != expected {
		t.Errorf("Error when wraping soap:\n got: %v\nwant: %v", WrapSoap("test"), expected)
	}
}
