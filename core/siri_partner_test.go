package core

import "testing"

func Test_SIRIPartner_SOAPClient(t *testing.T) {
	partner := &Partner{
		slug:     "partner",
		Settings: make(map[string]string),
	}
	siriPartner := NewSIRIPartner(partner)
	siriPartner.SOAPClient()
	if siriPartner.soapClient == nil {
		t.Error("siriPartner.SOAPClient() should set SIRIPartner soapClient")
	}

	partner.Settings["remote_url"] = "remote_url"
	siriPartner.SOAPClient()
	if siriPartner.soapClient.Url != "remote_url" {
		t.Error("SIRIPartner should have created a new SoapClient when partner setting changes")
	}

	partner.Settings["subscriptions.remote_url"] = "sub_remote_url"
	siriPartner.SOAPClient()
	if siriPartner.soapClient.SubscriptionsUrl != "sub_remote_url" {
		t.Error("SIRIPartner should have created a new SoapClient when partner setting changes")
	}
}

func Test_SIRIPartner_RequestorRef(t *testing.T) {
	partner := &Partner{
		slug: "partner",
		Settings: map[string]string{
			"remote_credential": "ara",
		},
	}
	siriPartner := NewSIRIPartner(partner)
	if siriPartner.RequestorRef() != "ara" {
		t.Errorf("Wrong SIRIPartner RequestorRef:\n got: %s\n want: \"ara\"", siriPartner.RequestorRef())
	}

}

// func Test_SIRIPartner_NewMessageIdentifier(t *testing.T) {
// 	partner := &Partner{
// 		slug: "partner",
// 	}
// 	siriPartner := NewSIRIPartner(partner)

// 	// Set MessageIdentifierGenerator
// 	midGenerator := NewFormatMessageIdentifierGenerator("Ara:Message::%s:LOC")
// 	midGenerator.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
// 	siriPartner.SetMessageIdentifierGenerator(midGenerator)

// 	mid := siriPartner.IdentifierGenerator("message_identifier").NewMessageIdentifier()
// 	if expected := "Ara:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; mid != expected {
// 		t.Errorf("Wrong MessageIdentifier:\n got: %s\n want: %s", mid, expected)
// 	}
// }
