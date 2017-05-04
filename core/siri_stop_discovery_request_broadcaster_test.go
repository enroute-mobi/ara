package core

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_SIRIStopPointDiscoveryRequestBroadcaster_StopAreas(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.Settings["remote_objectid_kind"] = "test"
	connector := NewSIRIStopDiscoveryRequestBroadcaster(partner)
	mid := NewFormatMessageIdentifierGenerator("Edwig:Message::%s:LOC")
	mid.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SIRIPartner().SetMessageIdentifierGenerator(mid)
	connector.SetClock(model.NewFakeClock())

	line := referential.Model().Lines().New()
	lineObjectId := model.NewObjectID("test", "1234")
	line.SetObjectID(lineObjectId)
	line.Save()

	firstStopArea := referential.Model().StopAreas().New()
	firstObjectID := model.NewObjectID("test", "NINOXE:StopPoint:SP:1:LOC")
	firstStopArea.SetObjectID(firstObjectID)
	firstStopArea.Name = "First"
	firstStopArea.LineIds = []model.LineId{line.Id()}
	firstStopArea.Save()

	secondStopArea := referential.Model().StopAreas().New()
	secondObjectID := model.NewObjectID("test", "NINOXE:StopPoint:SP:2:LOC")
	secondStopArea.SetObjectID(secondObjectID)
	secondStopArea.Name = "Second"
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
		t.Errorf("AnnotatedStopPoints lenght is wrong:\n got: %v\n want: 1", len(response.AnnotatedStopPoints))
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
