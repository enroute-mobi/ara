package core

import (
	"testing"

	"github.com/af83/edwig/model"
)

func Test_SIRIPartner_SOAPClient(t *testing.T) {
	partner := &Partner{
		slug: "partner",
	}
	connector := NewSIRIPartner(partner)
	connector.SOAPClient()
	if connector.soapClient == nil {
		t.Error("connector.SOAPClient() should set SIRIPartner soapClient")
	}
}

func Test_SIRIPartner_RequestorRef(t *testing.T) {
	partner := &Partner{
		slug: "partner",
		Settings: map[string]string{
			"remote_credential": "edwig",
		},
	}
	connector := NewSIRIPartner(partner)
	if connector.RequestorRef() != "edwig" {
		t.Errorf("Wrong SIRIPartner RequestorRef:\n got: %s\n want: \"edwig\"", connector.RequestorRef())
	}

}

func Test_SIRIPartner_NewMessageIdentifier(t *testing.T) {
	partner := &Partner{
		slug: "partner",
	}
	connector := NewSIRIPartner(partner)

	// Set MessageIdentifierGenerator
	midGenerator := NewFormatMessageIdentifierGenerator("Edwig:Message::%s:LOC")
	midGenerator.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetMessageIdentifierGenerator(midGenerator)

	mid := connector.NewMessageIdentifier()
	if expected := "Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; mid != expected {
		t.Errorf("Wrong MessageIdentifier:\n got: %s\n want: %s", mid, expected)
	}
}
