package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/af83/edwig/model"
)

func prepare_SIRIStopMonitoringRequestCollector(t *testing.T, responseFilePath string) *model.StopAreaUpdateEvent {
	// Create a test http server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}
		file, err := os.Open(responseFilePath)
		if err != nil {
			t.Fatal(err)
		}
		defer file.Close()
		io.Copy(w, file)
	}))
	defer ts.Close()

	// Create a SIRIStopMonitoringRequestCollector
	partners := NewPartnerManager(model.NewMemoryModel())
	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"remote_url":           ts.URL,
			"remote_objectid_kind": "test kind",
		},
	}
	partners.Save(partner)

	// Create StopArea with ObjectId
	stopArea := partners.Model().StopAreas().New()
	objectid := model.NewObjectID("test kind", "test value")
	stopArea.SetObjectID(objectid)
	partners.Model().StopAreas().Save(&stopArea)

	siriStopMonitoringRequestCollector := NewSIRIStopMonitoringRequestCollector(partner)
	stopAreaUpdateRequest := NewStopAreaUpdateRequest(stopArea.Id())
	stopAreaUpdateEvent, err := siriStopMonitoringRequestCollector.RequestStopAreaUpdate(stopAreaUpdateRequest)
	if err != nil {
		t.Fatal(err)
	}

	return stopAreaUpdateEvent
}

// WIP
func Test_SIRIStopMonitoringRequestCollector_RequestStopAreaUpdate(t *testing.T) {
	stopAreaUpdateEvent := prepare_SIRIStopMonitoringRequestCollector(t, "testdata/stopmonitoring-response-soap.xml")
	if stopAreaUpdateEvent == nil {
		t.Error("RequestStopAreaUpdate should not return nil")
	}
}

// Test Factory Validate
func Test_SIRIStopMonitoringRequestCollectorFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		Settings:       make(map[string]string),
		ConnectorTypes: []string{"siri-stop-monitoring-request-collector"},
		connectors:     make(map[string]Connector),
	}
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have three errors when remote_url and remote_objectid_kind aren't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_url":           "remote_url",
		"remote_objectid_kind": "remote_objectid_kind",
		"remote_credential":    "remote_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when remote_url is set, got: %v", apiPartner.Errors)
	}
}
