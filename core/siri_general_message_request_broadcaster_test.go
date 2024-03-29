package core

import (
	"io"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_SIRIGeneralMessageRequestBroadcaster_RequestSituation(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"local_url":                              "http://ara",
		"remote_code_space":                      "codeSpace",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())

	code := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:24:LOC")
	situation := referential.Model().Situations().New()
	period := &model.TimeRange{EndTime: referential.Clock().Now().Add(5 * time.Minute)}
	situation.ValidityPeriods = []*model.TimeRange{period}
	situation.Keywords = []string{"Perturbation"}
	situation.SetCode(code)

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	affectedStopArea := model.NewAffectedStopArea()
	affectedStopArea.StopAreaId = stopArea.Id()
	situation.Affects = append(situation.Affects, affectedStopArea)
	situation.Save()

	file, err := os.Open("testdata/generalmessage-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetGeneralMessageFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := connector.Situations(request, &audit.BigQueryMessage{})

	if response.Address != "http://ara" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://ara", response.Address)
	}
	if response.ProducerRef != "Ara" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Ara", response.ProducerRef)
	}
	if response.ResponseMessageIdentifier != "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
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
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"local_url":                              "http://ara",
		"remote_code_space":                      "codeSpace",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())

	code := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:24:LOC")
	situation := referential.Model().Situations().New()
	situation.Origin = "partner"
	period := &model.TimeRange{EndTime: referential.Clock().Now().Add(5 * time.Minute)}
	situation.ValidityPeriods = []*model.TimeRange{period}
	situation.SetCode(code)

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	affectedStopArea := model.NewAffectedStopArea()
	affectedStopArea.StopAreaId = stopArea.Id()
	situation.Affects = append(situation.Affects, affectedStopArea)
	situation.Save()

	file, err := os.Open("testdata/generalmessage-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetGeneralMessageFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := connector.Situations(request, &audit.BigQueryMessage{})

	if len(response.GeneralMessages) != 0 {
		t.Errorf("Response should have 0 GeneralMessage, got: %v", len(response.GeneralMessages))
	}
}

func Test_SIRIGeneralMessageRequestBroadcaster_RequestSituationWithFilter(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())

	settings := map[string]string{
		"local_url":                              "http://ara",
		"remote_code_space":                      "codeSpace",
		"generators.response_message_identifier": "Ara:ResponseMessage::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())

	line := referential.Model().Lines().New()
	line.SetCode(model.NewCode("codeSpace", "LineRef"))
	line.Save()

	code := model.NewCode("codeSpace", "NINOXE:StopPoint:SP:24:LOC")
	situation := referential.Model().Situations().New()
	period := &model.TimeRange{EndTime: referential.Clock().Now().Add(5 * time.Minute)}
	situation.ValidityPeriods = []*model.TimeRange{period}
	situation.Keywords = []string{"Perturbation"}
	situation.SetCode(code)

	stopArea := referential.Model().StopAreas().New()
	stopArea.SetCode(code)
	stopArea.Save()

	affectedStopArea := model.NewAffectedStopArea()
	affectedStopArea.StopAreaId = stopArea.Id()
	situation.Affects = append(situation.Affects, affectedStopArea)
	situation.Save()

	code2 := model.NewCode("codeSpace", "2")
	situation2 := referential.Model().Situations().New()
	situation2.ValidityPeriods = []*model.TimeRange{period}
	situation2.Keywords = []string{"Perturbation"}
	situation2.SetCode(code2)

	stopArea1 := referential.Model().StopAreas().New()
	code3 := model.NewCode("codeSpace", "DepartureStopArea")
	stopArea1.SetCode(code3)
	stopArea1.Save()

	affectedLine := model.NewAffectedLine()
	affectedLine.LineId = line.Id()
	affectedDestination := &model.AffectedDestination{StopAreaId: stopArea1.Id()}
	affectedLine.AffectedDestinations = append(affectedLine.AffectedDestinations, affectedDestination)
	situation2.Affects = append(situation2.Affects, affectedLine)
	situation2.Save()

	file, err := os.Open("testdata/generalmessage-request-lineref-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLGetGeneralMessageFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := connector.Situations(request, &audit.BigQueryMessage{})

	if response.Address != "http://ara" {
		t.Errorf("Response has wrong adress:\n got: %v\n want: http://ara", response.Address)
	}
	if response.ProducerRef != "Ara" {
		t.Errorf("Response has wrong producerRef:\n got: %v\n expected: Ara", response.ProducerRef)
	}
	if response.ResponseMessageIdentifier != "Ara:ResponseMessage::6ba7b814-9dad-11d1-0-00c04fd430c8:LOC" {
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
		ConnectorTypes: []string{"siri-general-message-request-broadcaster"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = s.NewEmptyPartnerSettings(partner.UUIDGenerator)
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have errors when local_credential and remote_code_space aren't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_code_space": "remote_code_space",
		"local_credential":  "local_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when local_credential and remote_code_space are set, got: %v", apiPartner.Errors)
	}
}

func Test_SIRIGeneralMessageRequestBroadcaster_RemoteCodeSpaceAbsent(t *testing.T) {
	partner := NewPartner()

	settings := map[string]string{
		"siri-general-message-request-broadcaster.remote_code_space": "",
		"remote_code_space": "CodeSpace2",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)

	if connector.partner.RemoteCodeSpace(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER) != "CodeSpace2" {
		t.Errorf("RemoteCodeSpace should be egals to CodeSpace2")
	}
}

func Test_SIRIGeneralMessageBroadcaster_RemoteCodeSpacePresent(t *testing.T) {
	partner := NewPartner()

	settings := map[string]string{
		"siri-general-message-request-broadcaster.remote_code_space": "CodeSpace1",
		"remote_code_space": "CodeSpace2",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIGeneralMessageRequestBroadcaster(partner)

	if connector.partner.RemoteCodeSpace(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER) != "CodeSpace1" {
		t.Errorf("RemoteCodeSpace should be egals to CodeSpace1")
	}
}
