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

func Test_SIRILinesDiscoveryRequestBroadcaster_Lines(t *testing.T) {
	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space":          "test",
		"generators.message_identifier": "Ara:Message::%{uuid}:LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewSIRILinesDiscoveryRequestBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	line := referential.Model().Lines().New()
	line.Name = "line1"
	lineCode := model.NewCode("test", "1234")
	line.SetCode(lineCode)
	line.Save()

	line2 := referential.Model().Lines().New()
	line2.Name = "line2"
	line2Code := model.NewCode("test2", "1234")
	line2.SetCode(line2Code)
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
	if response.AnnotatedLines[0].LineRef != lineCode.Value() {
		t.Errorf("AnnotatedLines LineRef is wrong:\n got: %v\n want: %v", response.AnnotatedLines[0].LineRef, lineCode.Value())
	}
	if !response.AnnotatedLines[0].Monitored {
		t.Errorf("AnnotatedLines Monitored is false, should be true")
	}
}
