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
			"remote_credential": "edwig",
		},
	}
	siriPartner := NewSIRIPartner(partner)
	if siriPartner.RequestorRef() != "edwig" {
		t.Errorf("Wrong SIRIPartner RequestorRef:\n got: %s\n want: \"edwig\"", siriPartner.RequestorRef())
	}

}

// func Test_SIRIPartner_NewMessageIdentifier(t *testing.T) {
// 	partner := &Partner{
// 		slug: "partner",
// 	}
// 	siriPartner := NewSIRIPartner(partner)

// 	// Set MessageIdentifierGenerator
// 	midGenerator := NewFormatMessageIdentifierGenerator("Edwig:Message::%s:LOC")
// 	midGenerator.SetUUIDGenerator(model.NewFakeUUIDGenerator())
// 	siriPartner.SetMessageIdentifierGenerator(midGenerator)

// 	mid := siriPartner.IdentifierGenerator("message_identifier").NewMessageIdentifier()
// 	if expected := "Edwig:Message::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC"; mid != expected {
// 		t.Errorf("Wrong MessageIdentifier:\n got: %s\n want: %s", mid, expected)
// 	}
// }

func Test_SIRIPartner_IdentifierGenerator(t *testing.T) {
	partner := &Partner{
		slug:     "partner",
		Settings: make(map[string]string),
	}
	siriPartner := NewSIRIPartner(partner)

	midGenerator := siriPartner.IdentifierGenerator("message_identifier")
	if expected := "%{uuid}"; midGenerator.formatString != expected {
		t.Errorf("siriPartner message_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = siriPartner.IdentifierGenerator("response_message_identifier")
	if expected := "%{uuid}"; midGenerator.formatString != expected {
		t.Errorf("siriPartner response_message_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = siriPartner.IdentifierGenerator("data_frame_identifier")
	if expected := "%{id}"; midGenerator.formatString != expected {
		t.Errorf("siriPartner data_frame_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = siriPartner.IdentifierGenerator("reference_identifier")
	if expected := "%{type}:%{default}"; midGenerator.formatString != expected {
		t.Errorf("siriPartner reference_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = siriPartner.IdentifierGenerator("reference_stop_area_identifier")
	if expected := "%{default}"; midGenerator.formatString != expected {
		t.Errorf("siriPartner reference_stop_area_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}

	partner.Settings = map[string]string{
		"generators.message_identifier":             "mid",
		"generators.response_message_identifier":    "rmid",
		"generators.data_frame_identifier":          "dfid",
		"generators.reference_identifier":           "rid",
		"generators.reference_stop_area_identifier": "rsaid",
	}

	midGenerator = siriPartner.IdentifierGenerator("message_identifier")
	if expected := "mid"; midGenerator.formatString != expected {
		t.Errorf("siriPartner message_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = siriPartner.IdentifierGenerator("response_message_identifier")
	if expected := "rmid"; midGenerator.formatString != expected {
		t.Errorf("siriPartner response_message_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = siriPartner.IdentifierGenerator("data_frame_identifier")
	if expected := "dfid"; midGenerator.formatString != expected {
		t.Errorf("siriPartner data_frame_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = siriPartner.IdentifierGenerator("reference_identifier")
	if expected := "rid"; midGenerator.formatString != expected {
		t.Errorf("siriPartner reference_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
	midGenerator = siriPartner.IdentifierGenerator("reference_stop_area_identifier")
	if expected := "rsaid"; midGenerator.formatString != expected {
		t.Errorf("siriPartner reference_stop_area_identifier IdentifierGenerator should be %v, got: %v ", expected, midGenerator.formatString)
	}
}
