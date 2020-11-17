package core

import (
	"io/ioutil"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_SIRIStopPointDiscoveryRequestBroadcaster_StopAreas(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "test"
	partner.Settings["generators.message_identifier"] = "Ara:Message::%{uuid}:LOC"
	connector := NewSIRIStopDiscoveryRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.Model().Lines().New()
	lineObjectId := model.NewObjectID("test", "1234")
	line.SetObjectID(lineObjectId)
	line.Save()

	line2 := referential.Model().Lines().New()
	lineObjectId2 := model.NewObjectID("test", "5678")
	line2.SetObjectID(lineObjectId2)
	line2.Save()

	line3 := referential.Model().Lines().New()
	lineObjectId3 := model.NewObjectID("_default", "5678")
	line3.SetObjectID(lineObjectId3)
	line3.Save()

	line4 := referential.Model().Lines().New()
	lineObjectId4 := model.NewObjectID("test", "91011")
	line4.SetObjectID(lineObjectId4)
	line4.SetOrigin("partner")
	line4.Save()

	firstStopArea := referential.Model().StopAreas().New()
	firstObjectID := model.NewObjectID("test", "NINOXE:StopPoint:SP:1:LOC")
	firstStopArea.SetObjectID(firstObjectID)
	firstStopArea.Name = "First"
	firstStopArea.LineIds = []model.LineId{line.Id(), line3.Id(), line4.Id()}
	firstStopArea.Save()

	secondStopArea := referential.Model().StopAreas().New()
	secondObjectID := model.NewObjectID("test", "NINOXE:StopPoint:SP:2:LOC")
	secondStopArea.SetObjectID(secondObjectID)
	secondStopArea.Name = "Second"
	secondStopArea.LineIds = []model.LineId{line2.Id()}
	secondStopArea.Save()

	file, err := os.Open("testdata/stoppointdiscovery-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLStopPointsDiscoveryRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, err := connector.StopAreas(request)
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
	if response.AnnotatedStopPoints[0].StopPointRef != firstObjectID.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[0].StopPointRef, firstObjectID.Value())
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
	if response.AnnotatedStopPoints[1].StopPointRef != secondObjectID.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[1].StopPointRef, secondObjectID.Value())
	}
}

func Test_SIRIStopPointDiscoveryRequestBroadcaster_StopAreasWithParent(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "test"
	partner.Settings["generators.message_identifier"] = "Ara:Message::%{uuid}:LOC"
	connector := NewSIRIStopDiscoveryRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.Model().Lines().New()
	lineObjectId := model.NewObjectID("test", "1234")
	line.SetObjectID(lineObjectId)
	line.Save()

	firstStopArea := referential.Model().StopAreas().New()
	firstObjectID := model.NewObjectID("test_incorrect", "NINOXE:StopPoint:SP:1:LOC")
	firstStopArea.SetObjectID(firstObjectID)
	firstStopArea.Name = "First"
	firstStopArea.LineIds = []model.LineId{line.Id()}
	firstStopArea.Save()

	secondStopArea := referential.Model().StopAreas().New()
	secondObjectID := model.NewObjectID("test", "NINOXE:StopPoint:SP:2:LOC")
	secondStopArea.SetObjectID(secondObjectID)
	secondStopArea.Name = "Second"
	secondStopArea.LineIds = []model.LineId{line.Id()}
	secondStopArea.Save()

	thirdStopArea := referential.Model().StopAreas().New()
	thirdObjectID := model.NewObjectID("test", "NINOXE:StopPoint:SP:3:LOC")
	thirdStopArea.SetObjectID(thirdObjectID)
	thirdStopArea.ReferentId = secondStopArea.Id()
	thirdStopArea.Name = "Third"
	thirdStopArea.LineIds = []model.LineId{line.Id()}
	thirdStopArea.Save()

	fourthStopArea := referential.Model().StopAreas().New()
	fourthObjectID := model.NewObjectID("test", "NINOXE:StopPoint:SP:4:LOC")
	fourthStopArea.SetObjectID(fourthObjectID)
	fourthStopArea.ReferentId = firstStopArea.Id()
	fourthStopArea.Name = "Fourth"
	fourthStopArea.LineIds = []model.LineId{line.Id()}
	fourthStopArea.Save()

	file, err := os.Open("testdata/stoppointdiscovery-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLStopPointsDiscoveryRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, err := connector.StopAreas(request)
	if err != nil {
		t.Fatal(err)
	}

	if response.Status != true {
		t.Errorf("Response status wrong:\n got: %v\n want: true", response.Status)
	}
	if len(response.AnnotatedStopPoints) != 2 {
		t.Fatalf("AnnotatedStopPoints lenght is wrong:\n got: %v\n want: 2\n%v", len(response.AnnotatedStopPoints), response.AnnotatedStopPoints)
	}

	if response.AnnotatedStopPoints[0].StopPointRef != secondObjectID.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef 1 is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[0].StopPointRef, firstObjectID.Value())
	}
	if response.AnnotatedStopPoints[1].StopPointRef != fourthObjectID.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef 2 is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[1].StopPointRef, secondObjectID.Value())
	}
}
