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

	stopArea := referential.Model().StopAreas().New()
	objectID := model.NewObjectID("test", "NINOXE:StopPoint:SP:24:LOC")
	stopArea.SetObjectID(objectID)
	stopArea.Name = "Charle"
	stopArea.Save()

	file, err := os.Open("testdata/stoppointdiscovery-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := siri.NewXMLStopDiscoveryRequestFromContent(content)
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
	if len(response.AnnotatedStopPoints) != 1 {
		t.Errorf("AnnotatedStopPoints lenght is wrong:\n got: %v\n want: 1", len(response.AnnotatedStopPoints))
	}
	if response.AnnotatedStopPoints[0].StopName != "Charle" {
		t.Errorf("AnnotatedStopPoints StopName is wrong:\n got: %v\n want: Charle", response.AnnotatedStopPoints[0].StopName)
	}

	if response.AnnotatedStopPoints[0].StopPointRef != objectID.Value() {
		t.Errorf("AnnotatedStopPoints StopPointRef is wrong:\n got: %v\n want: %v", response.AnnotatedStopPoints[0].StopPointRef, objectID.Value())
	}
}
