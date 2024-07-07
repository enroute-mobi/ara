package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

type fakeSituationBroadcaster struct {
	Events []model.UpdateEvent
}

func (fb *fakeSituationBroadcaster) FakeBroadcaster(event model.UpdateEvent) {
	fb.Events = append(fb.Events, event)
}

func prepare_SIRIGeneralMessageRequestCollector(t *testing.T, responseFilePath string) []model.UpdateEvent {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request Content Length should be zero")
		}
		file, err := os.Open(responseFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))
	defer ts.Close()

	partners := createTestPartnerManager()
	partner := partners.New("slug")
	settings := map[string]string{
		"remote_url":        ts.URL,
		"remote_code_space": "test kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partners.Save(partner)

	situation := partners.Model().Situations().New()
	code := model.NewCode("test kind", "test value")
	situation.SetCode(code)
	partners.Model().Situations().Save(&situation)

	line := partners.Model().Lines().New()
	lineCode := model.NewCode("test kind", "line value")
	line.SetCode(lineCode)
	partners.Model().Lines().Save(line)

	siriGeneralMessageRequestCollector := NewSIRIGeneralMessageRequestCollector(partner)

	fs := fakeSituationBroadcaster{}
	siriGeneralMessageRequestCollector.SetUpdateSubscriber(fs.FakeBroadcaster)
	siriGeneralMessageRequestCollector.SetClock(clock.NewFakeClock())
	siriGeneralMessageRequestCollector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, "line value")

	return fs.Events
}

func Test_SIRIGeneralMessageCollectorFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		ConnectorTypes: []string{"siri-general-message-request-collector"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = s.NewEmptyPartnerSettings(partner.UUIDGenerator)
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have three errors when remote_url and remote_code_space aren't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_url":        "remote_url",
		"remote_code_space": "remote_code_space",
		"remote_credential": "remote_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when remote_url is set, got: %v", apiPartner.Errors)
	}
}

func Test_SIRIGeneralMessageRequestCollector_RequestSituationUpdate(t *testing.T) {
	assert := assert.New(t)
	updateEvents := prepare_SIRIGeneralMessageRequestCollector(t, "testdata/generalmessage-response-soap.xml")
	if updateEvents == nil {
		t.Error("RequestSituationUpdate should not return nil")
	}

	if len(updateEvents) != 2 {
		t.Errorf("RequestSituationUpdate should have 2 SituationUpdateEvents, got: %v", len(updateEvents))
	}
	situationEvent, _ := updateEvents[0].(*model.SituationUpdateEvent)

	if expected := clock.FAKE_CLOCK_INITIAL_DATE; situationEvent.CreatedAt != expected {
		t.Errorf("Wrong Created_At for situationEvent:\n expected: %v\n got: %v", expected, situationEvent.CreatedAt)
	}
	if expected, _ := time.Parse(time.RFC3339, "2017-03-29T03:30:06.000+02:00"); !situationEvent.RecordedAt.Equal(expected) {
		t.Errorf("Wrong RecorderAt for situationEvent:\n expected: %v\n got: %v", expected, situationEvent.RecordedAt)
	}

	if expected := 1; situationEvent.Version != expected {
		t.Errorf("Wrong Version for situationEvent:\n expected: %v\n got: %v", expected, situationEvent.Version)
	}

	assert.ElementsMatch([]string{"Commercial"}, situationEvent.Keywords)

	expected, _ := time.Parse(time.RFC3339, "2017-03-29T20:30:06.000+02:00")
	assert.Equal(expected, situationEvent.ValidityPeriods[0].EndTime)

	if expected := "NINOXE:default"; situationEvent.ProducerRef != "NINOXE:default" {
		t.Errorf("Wrong ProducerRef for situationEvent:\n expected: %v\n got: %v", expected, situationEvent.ProducerRef)
	}

	assert.Equal("Un deuxième message #etouais", situationEvent.Description.DefaultValue)
}
