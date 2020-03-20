package core

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
)

func Test_SIRIGeneralMessageRequestBroadcaster_RequestSituation(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"

	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	objectid := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	situation := referential.Model().Situations().New()
	situation.ValidUntil = referential.Clock().Now().Add(5 * time.Minute)
	situation.SetObjectID(objectid)
	routeReference := model.NewReference(model.NewObjectID("internal", "value"))
	routeReference.Type = "RouteRef"
	situation.References = append(situation.References, routeReference)
	situation.Save()

	file, err := os.Open("testdata/generalmessage-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetGeneralMessageFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := connector.Situations(request)

	if response.Address != "http://edwig" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://edwig", response.Address)
	}
	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: RATPDev:Message::ade15433-06a6-4f7b-b331-2c1080a5d279:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
	if len(response.GeneralMessages) != 1 {
		t.Errorf("Response should have 1 GeneralMessage, got: %v", len(response.GeneralMessages))
	}
}

func Test_SIRIGeneralMessageRequestBroadcaster_RequestSituationWithSameOrigin(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"

	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	objectid := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	situation := referential.Model().Situations().New()
	situation.Origin = "partner"
	situation.ValidUntil = referential.Clock().Now().Add(5 * time.Minute)
	situation.SetObjectID(objectid)
	routeReference := model.NewReference(model.NewObjectID("internal", "value"))
	routeReference.Type = "RouteRef"
	situation.References = append(situation.References, routeReference)
	situation.Save()

	file, err := os.Open("testdata/generalmessage-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetGeneralMessageFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := connector.Situations(request)

	if len(response.GeneralMessages) != 0 {
		t.Errorf("Response should have 0 GeneralMessage, got: %v", len(response.GeneralMessages))
	}
}

func Test_SIRIGeneralMessageRequestBroadcaster_RequestSituationWithFilter(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	partner.Settings["generators.response_message_identifier"] = "Edwig:ResponseMessage::%{uuid}:LOC"

	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SetClock(model.NewFakeClock())

	line := referential.Model().Lines().New()
	line.SetObjectID(model.NewObjectID("objectidKind", "LineRef"))
	line.Save()

	objectid := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	situation := referential.Model().Situations().New()
	situation.ValidUntil = referential.Clock().Now().Add(5 * time.Minute)
	situation.SetObjectID(objectid)
	routeReference := model.NewReference(model.NewObjectID("internal", "value"))
	routeReference.Type = "RouteRef"
	situation.References = append(situation.References, routeReference)
	situation.Save()

	objectid2 := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:25:LOC")
	situation2 := referential.Model().Situations().New()
	situation2.ValidUntil = referential.Clock().Now().Add(5 * time.Minute)
	situation2.SetObjectID(objectid2)
	lineReference := model.NewReference(model.NewObjectID("objectidKind", "LineRef"))
	lineReference.Type = "LineRef"
	situation2.References = append(situation.References, lineReference)
	situation2.Save()

	file, err := os.Open("testdata/generalmessage-request-lineref-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLGetGeneralMessageFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := connector.Situations(request)

	if response.Address != "http://edwig" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://edwig", response.Address)
	}
	if response.ProducerRef != "Edwig" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Edwig", response.ProducerRef)
	}
	if response.ResponseMessageIdentifier != "Edwig:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
		t.Errorf("Response has wesponseMessageIdentifier:\n got: %v\n expected: RATPDev:Message::ade15433-06a6-4f7b-b331-2c1080a5d279:LOC", response.ResponseMessageIdentifier)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
	if len(response.GeneralMessages) != 1 {
		t.Errorf("Response should have 1 GeneralMessage, got: %v", len(response.GeneralMessages))
	}
}

func Test_SIRIGeneralMessageRequestBroadcasterFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		Settings:       make(map[string]string),
		ConnectorTypes: []string{"siri-general-message-request-broadcaster"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have errors when local_credential and remote_objectid_kind aren't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_objectid_kind": "remote_objectid_kind",
		"local_credential":     "local_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when local_credential and remote_objectid_kind are set, got: %v", apiPartner.Errors)
	}
}

func Test_SIRIGeneralMessageRequestBroadcaster_RemoteObjectIDKindAbsent(t *testing.T) {
	partner := &Partner{}
	partner.Settings = make(map[string]string)

	partner.Settings["siri-general-message-request-broadcaster.remote_objectid_kind"] = ""
	partner.Settings["remote_objectid_kind"] = "Kind2"

	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)

	if connector.partner.RemoteObjectIDKind(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER) != "Kind2" {
		t.Errorf("RemoteObjectIDKind should be egals to Kind2")
	}
}

func Test_SIRIGeneralMessageBroadcaster_RemoteObjectIDKindPresent(t *testing.T) {
	partner := &Partner{}
	partner.Settings = make(map[string]string)

	partner.Settings["siri-general-message-request-broadcaster.remote_objectid_kind"] = "Kind1"
	partner.Settings["remote_objectid_kind"] = "Kind2"

	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)

	if connector.partner.RemoteObjectIDKind(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER) != "Kind1" {
		t.Errorf("RemoteObjectIDKind should be egals to Kind1")
	}
}
