package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func prepare_siriCheckStatusClient(t *testing.T, responseFilePath string) OperationnalStatus {
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
	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"remote_url": ts.URL,
		},
	}
	checkStatusClient := NewSIRICheckStatusClient(partner)

	status, err := checkStatusClient.Status()

	if err != nil {
		t.Fatal(err)
	}

	return status
}

func Test_SIRICheckStatusClient_Status_OK(t *testing.T) {
	status := prepare_siriCheckStatusClient(t, "testdata/checkstatus-response-soap.xml")
	if status != OPERATIONNAL_STATUS_UP {
		t.Errorf("Wrong status found:\n got: %v\n expected: 1", status)
	}
}

func Test_SIRICheckStatusClient_Status_KO(t *testing.T) {
	status := prepare_siriCheckStatusClient(t, "testdata/checkstatus-negative-response-soap.xml")
	if status != OPERATIONNAL_STATUS_DOWN {
		t.Errorf("Wrong status found:\n got: %v\n expected: 2", status)
	}
}

func Test_SIRICheckStatusClientFactory_Validate(t *testing.T) {
	partner := &Partner{
		Settings:       make(map[string]string),
		ConnectorTypes: []string{"siri-check-status-client"},
		connectors:     make(map[string]Connector),
	}
	apiPartner := partner.Definition()
	apiPartner.Validate()
	if len(apiPartner.Errors) != 2 {
		t.Errorf("apiPartner should have two errors when remote_url isn't set, got: %v", apiPartner.Errors)
	}

	apiPartner.Settings = map[string]string{
		"remote_url":        "remote_url",
		"remote_credential": "remote_credential",
	}
	apiPartner.Validate()
	if len(apiPartner.Errors) != 0 {
		t.Errorf("apiPartner shouldn't have any error when remote_url is set, got: %v", apiPartner.Errors)
	}
}
