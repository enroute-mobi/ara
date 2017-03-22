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
	partner.Settings["local_url"] = "http://edwig"
	partner.Settings["remote_objectid_kind"] = "objectidKind"
	connector := NewSIRIStopDiscoveryRequestBroadcaster(partner)
	mid := NewFormatMessageIdentifierGenerator("Edwig:Message::%s:LOC")
	mid.SetUUIDGenerator(model.NewFakeUUIDGenerator())
	connector.SIRIPartner().SetMessageIdentifierGenerator(mid)
	connector.SetClock(model.NewFakeClock())

	objectid := model.NewObjectID("objectidKind", "NINOXE:StopPoint:SP:24:LOC")
	refObj := model.NewObjectID("internal", "NINOXE:StopPoint:SP:16:LOC")
	stopArea := referential.Model().StopAreas().New()
	stopArea.SetObjectID(objectid)
	stopArea.Name = "Charle"
	stopArea.References = make(model.References)
	stopArea.References["StopPointRef"] = model.Reference{ObjectId: &refObj, Id: ""}
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
	if response.AnnotatedStopPoints[0].StopPointName != "Charle" {
		t.Errorf("AnnotatedStopPoints StopPointName is wrong:\n got: %v\n want: Charle", response.AnnotatedStopPoints[0].StopPointName)
	}

	if response.AnnotatedStopPoints[0].StopPointRef != "NINOXE:StopPoint:SP:16:LOC" {
		t.Errorf("AnnotatedStopPoints lenght is wrong:\n got: %v\n want: NINOXE:StopPoint:SP:16:LOC", response.AnnotatedStopPoints[0].StopPointRef)
	}
}
