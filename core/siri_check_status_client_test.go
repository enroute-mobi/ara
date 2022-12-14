package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	s "bitbucket.org/enroute-mobi/ara/core/settings"
)

func prepare_siriCheckStatusClient(t *testing.T, responseFilePath string) PartnerStatus {
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

	// Create a CheckStatusClient
	referentials := NewMemoryReferentials()
	referential := referentials.New("slug")
	partner := referential.Partners().New("slug")
	partner.SetSetting(s.REMOTE_URL, ts.URL)
	checkStatusClient := NewSIRICheckStatusClient(partner)

	partnerStatus, err := checkStatusClient.Status()
	if err != nil {
		t.Fatal(err)
	}

	return partnerStatus
}

func Test_SIRICheckStatusClient_Status_OK(t *testing.T) {
	partnerStatus := prepare_siriCheckStatusClient(t, "testdata/checkstatus-response-soap.xml")
	if partnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UP {
		t.Errorf("Wrong status found:\n got: %v\n expected: up", partnerStatus.OperationnalStatus)
	}
}

func Test_SIRICheckStatusClient_Status_KO(t *testing.T) {
	partnerStatus := prepare_siriCheckStatusClient(t, "testdata/checkstatus-negative-response-soap.xml")
	if partnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_DOWN {
		t.Errorf("Wrong status found:\n got: %v\n expected: down", partnerStatus.OperationnalStatus)
	}
}

func Test_SIRICheckStatusClientFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		ConnectorTypes: []string{"siri-check-status-client"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator)
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if apiPartner.Errors.Empty() {
		t.Errorf("apiPartner should have two errors when remote_url isn't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_url":        "remote_url",
		"remote_credential": "remote_credential",
	}
	apiPartner.Validate()
	if !apiPartner.Errors.Empty() {
		t.Errorf("apiPartner shouldn't have any error when remote_url is set, got: %v", apiPartner.Errors)
	}
}
