package core

import (
	"io"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

func Test_SIRILinesDiscoveryRequestBroadcaster_Lines(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetSetting("remote_objectid_kind", "test")
	partner.SetSetting("generators.message_identifier", "Ara:Message::%{uuid}:LOC")
	connector := NewSIRILinesDiscoveryRequestBroadcaster(partner)
	connector.Partner().SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	connector.SetClock(clock.NewFakeClock())

	line := referential.Model().Lines().New()
	line.Name = "line1"
	lineObjectId := model.NewObjectID("test", "1234")
	line.SetObjectID(lineObjectId)
	line.Save()

	line2 := referential.Model().Lines().New()
	line2.Name = "line2"
	line2ObjectId := model.NewObjectID("test2", "1234")
	line2.SetObjectID(line2ObjectId)
	line2.Save()

	file, err := os.Open("testdata/stoppointdiscovery-request-soap.xml")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}
	request, err := sxml.NewXMLLinesDiscoveryRequestFromContent(content)
	if err != nil {
		t.Fatal(err)
	}

	response, err := connector.Lines(request, &audit.BigQueryMessage{})
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
	if len(response.AnnotatedLines) != 1 {
		t.Fatalf("AnnotatedLines lenght is wrong:\n got: %v\n want: 1", len(response.AnnotatedLines))
	}

	if response.AnnotatedLines[0].LineName != "line1" {
		t.Errorf("AnnotatedLines LineName is wrong:\n got: %v\n want: line1", response.AnnotatedLines[0].LineName)
	}
	if response.AnnotatedLines[0].LineRef != lineObjectId.Value() {
		t.Errorf("AnnotatedLines LineRef is wrong:\n got: %v\n want: %v", response.AnnotatedLines[0].LineRef, lineObjectId.Value())
	}
	if !response.AnnotatedLines[0].Monitored {
		t.Errorf("AnnotatedLines Monitored is false, should be true")
	}
}
