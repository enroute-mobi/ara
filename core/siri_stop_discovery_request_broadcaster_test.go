package core

import (
	"io"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_SIRIStopPointDiscoveryRequestBroadcaster_StopAreas(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space":          "test",
		"generators.message_identifier": "Ara:Message::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIStopDiscoveryRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.Model().Lines().New()
	lineCode := model.NewCode("test", "1234")
	line.SetCode(lineCode)
	line.Save()

	line2 := referential.Model().Lines().New()
	lineCode2 := model.NewCode("test", "5678")
	line2.SetCode(lineCode2)
	line2.Save()

	line3 := referential.Model().Lines().New()
	lineCode3 := model.NewCode("_default", "5678")
	line3.SetCode(lineCode3)
	line3.Save()

	line4 := referential.Model().Lines().New()
	lineCode4 := model.NewCode("test", "91011")
	line4.SetCode(lineCode4)
	line4.SetOrigin("partner")
	line4.Save()

	firstStopArea := referential.Model().StopAreas().New()
	firstCode := model.NewCode("test", "NINOXE:StopPoint:SP:1:LOC")
	firstStopArea.SetCode(firstCode)
	firstStopArea.Name = "First"
	firstStopArea.LineIds = []model.LineId{line.Id(), line3.Id(), line4.Id()}
	firstStopArea.Save()

	secondStopArea := referential.Model().StopAreas().New()
	secondCode := model.NewCode("test", "NINOXE:StopPoint:SP:2:LOC")
	secondStopArea.SetCode(secondCode)
	secondStopArea.Name = "Second"
	secondStopArea.LineIds = []model.LineId{line2.Id()}
	secondStopArea.Save()

	file, err := os.Open("testdata/stoppointdiscovery-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLStopPointsDiscoveryRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, err := connector.StopAreas(request, &audit.BigQueryMessage{})
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != true {
		t.Errorf("Response status wrong:\n got: %v\n want: true", response.Status)
	}
	time := connector.Clock().Now()
	if !response.ResponseTimestamp.Equal(time) {
		t.Errorf("Response has wrong responseTimestamp:\n got: %v\n expected: 2016-09-22 08:01:20.227 +0200 CEST", response.ResponseTimestamp)
	}
	if len(response.AnnotatedStopPoints) != 2 {
		t.Fatalf("AnnotatedStopPoints lenght is wrong:\n got: %v\n want: 2", len(response.AnnotatedStopPoints))
	}

	if response.AnnotatedStopPoints[0].StopName != "First" {
		t.Errorf("AnnotatedStopPoints StopName is wrong:\n got: %v\n want: First", response.AnnotatedStopPoints[0].StopName)
	}
	if response.AnnotatedStopPoints[0].StopPointRef != firstCode.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[0].StopPointRef, firstCode.Value())
	}
	if !response.AnnotatedStopPoints[0].Monitored {
		t.Errorf("AnnotatedStopPoints Monitored is false, should be true")
	}
	if !response.AnnotatedStopPoints[0].TimingPoint {
		t.Errorf("AnnotatedStopPoints TimingPoint is false, should be true")
	}
	if len(response.AnnotatedStopPoints[0].Lines) != 1 || response.AnnotatedStopPoints[0].Lines[0] != "1234" {
		t.Errorf("AnnotatedStopPoints Lines is wrong:\n got: %v\n want: [1234]", response.AnnotatedStopPoints[0].Lines)
	}

	if response.AnnotatedStopPoints[1].StopName != "Second" {
		t.Errorf("AnnotatedStopPoints StopName is wrong:\n got: %v\n want: Second", response.AnnotatedStopPoints[1].StopName)
	}
	if response.AnnotatedStopPoints[1].StopPointRef != secondCode.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[1].StopPointRef, secondCode.Value())
	}
}

func Test_SIRIStopPointDiscoveryRequestBroadcaster_StopAreasWithParent(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space":          "test",
		"generators.message_identifier": "Ara:Message::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRIStopDiscoveryRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.Model().Lines().New()
	lineCode := model.NewCode("test", "1234")
	line.SetCode(lineCode)
	line.Save()

	firstStopArea := referential.Model().StopAreas().New()
	firstCode := model.NewCode("test_incorrect", "NINOXE:StopPoint:SP:1:LOC")
	firstStopArea.SetCode(firstCode)
	firstStopArea.Name = "First"
	firstStopArea.LineIds = []model.LineId{line.Id()}
	firstStopArea.Save()

	secondStopArea := referential.Model().StopAreas().New()
	secondCode := model.NewCode("test", "NINOXE:StopPoint:SP:2:LOC")
	secondStopArea.SetCode(secondCode)
	secondStopArea.Name = "Second"
	secondStopArea.LineIds = []model.LineId{line.Id()}
	secondStopArea.Save()

	thirdStopArea := referential.Model().StopAreas().New()
	thirdCode := model.NewCode("test", "NINOXE:StopPoint:SP:3:LOC")
	thirdStopArea.SetCode(thirdCode)
	thirdStopArea.ReferentId = secondStopArea.Id()
	thirdStopArea.Name = "Third"
	thirdStopArea.LineIds = []model.LineId{line.Id()}
	thirdStopArea.Save()

	fourthStopArea := referential.Model().StopAreas().New()
	fourthCode := model.NewCode("test", "NINOXE:StopPoint:SP:4:LOC")
	fourthStopArea.SetCode(fourthCode)
	fourthStopArea.ReferentId = firstStopArea.Id()
	fourthStopArea.Name = "Fourth"
	fourthStopArea.LineIds = []model.LineId{line.Id()}
	fourthStopArea.Save()

	file, err := os.Open("testdata/stoppointdiscovery-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLStopPointsDiscoveryRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, err := connector.StopAreas(request, &audit.BigQueryMessage{})
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != true {
		t.Errorf("Response status wrong:\n got: %v\n want: true", response.Status)
	}

	if len(response.AnnotatedStopPoints) != 2 {
		t.Fatalf("AnnotatedStopPoints lenght is wrong:\n got: %v\n want: 2\n%v", len(response.AnnotatedStopPoints), response.AnnotatedStopPoints)
	}

	if response.AnnotatedStopPoints[0].StopPointRef != secondCode.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef 1 is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[0].StopPointRef, firstCode.Value())
	}

	if response.AnnotatedStopPoints[1].StopPointRef != fourthCode.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef 1 is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[0].StopPointRef, firstCode.Value())
	}
}
