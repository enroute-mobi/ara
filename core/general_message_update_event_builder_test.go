package core

import (
	"io/ioutil"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func Test_GeneralMessageUpdateEventBuilder_BuildGeneralMessageUpdateEvent(t *testing.T) {
	file, err := os.Open("testdata/long-general-message-response.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	response, _ := sxml.NewXMLGeneralMessageResponseFromContent(content)

	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	partner := referential.Partners().New("slug")
	partner.SetSetting("remote_objectid_kind", "remote_objectid_kind")
	builder := NewGeneralMessageUpdateEventBuilder(partner)

	events := &[]*model.SituationUpdateEvent{}

	builder.buildGeneralMessageUpdateEvent(events, response.XMLGeneralMessages()[0], "producer")

	if len(*events) != 1 {
		t.Fatalf("One event should have been created, got %v", len(*events))
	}

	event := (*events)[0]
	if event.SituationAttributes.Format != "FRANCE" {
		t.Errorf("Wrong Format, expected: FRANCE, got: %v", event.SituationAttributes.Format)
	}
	if event.SituationAttributes.Channel != "Commercial" {
		t.Errorf("Wrong Channel, expected: Commercial, got: %v", event.SituationAttributes.Channel)
	}
	if len(event.SituationAttributes.References) != 12 {
		t.Fatalf("Wrong number of References, expected: 12, got: %v", len(event.SituationAttributes.References))
	}
	if event.SituationAttributes.References[0].ObjectId.Value() != "lineRef1" {
		t.Errorf("Wrong first LineRef: %v", event.SituationAttributes.References[0])
	}
	if event.SituationAttributes.References[2].ObjectId.Value() != "stopPointRef1" {
		t.Errorf("Wrong first StopPointRef: %v", event.SituationAttributes.References[2])
	}
	if event.SituationAttributes.References[4].ObjectId.Value() != "journeyPatternRef1" {
		t.Errorf("Wrong first JourneyPatternRef: %v", event.SituationAttributes.References[4])
	}
	if event.SituationAttributes.References[6].ObjectId.Value() != "destinationRef1" {
		t.Errorf("Wrong first DestinationRef: %v", event.SituationAttributes.References[6])
	}
	if event.SituationAttributes.References[8].ObjectId.Value() != "routeRef1" {
		t.Errorf("Wrong first RouteRef: %v", event.SituationAttributes.References[8])
	}
	if event.SituationAttributes.References[10].ObjectId.Value() != "groupOfLineRef1" {
		t.Errorf("Wrong first GroupOfLinesRef: %v", event.SituationAttributes.References[10])
	}

	if len(event.SituationAttributes.LineSections) != 2 {
		t.Fatalf("Wrong number of LineSections, expected: 2, got: %v", len(event.SituationAttributes.LineSections))
	}
	firstLineSection := *event.SituationAttributes.LineSections[0]
	if ref, _ := firstLineSection.Get("FirstStop"); ref.ObjectId.Value() != "firstStop1" {
		t.Errorf("Wrong first LineSection FirstStop: %v", ref)
	}
	if ref, _ := firstLineSection.Get("LastStop"); ref.ObjectId.Value() != "lastStop1" {
		t.Errorf("Wrong first LineSection LastStop: %v", ref)
	}
	if ref, _ := firstLineSection.Get("LineRef"); ref.ObjectId.Value() != "lineSectionRef1" {
		t.Errorf("Wrong first LineSection LineRef: %v", ref)
	}
}
