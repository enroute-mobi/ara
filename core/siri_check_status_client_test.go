package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/core/partners"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"github.com/stretchr/testify/assert"
)

func prepareSiriCheckStatusClient(t *testing.T, responseFilePath string) (partners.Status, error) {
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

	settings := map[string]string{s.REMOTE_URL: ts.URL}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)

	checkStatusClient := NewSIRICheckStatusClient(partner)

	partnerStatus, err := checkStatusClient.Status()

	return partnerStatus, err
}

func Test_SIRICheckStatusClient_Status_OK(t *testing.T) {
	assert := assert.New(t)
	partnerStatus, err := prepareSiriCheckStatusClient(t, "testdata/checkstatus-response-soap.xml")

	assert.Nil(err)
	assert.Equal(partnerStatus.OperationnalStatus, partners.OperationnalStatusUp)
}

func Test_SIRICheckStatusClient_Status_KO(t *testing.T) {
	assert := assert.New(t)
	partnerStatus, err := prepareSiriCheckStatusClient(t, "testdata/checkstatus-negative-response-soap.xml")

	assert.Nil(err)
	assert.Equal(partnerStatus.OperationnalStatus, partners.OperationnalStatusDown)
}

func Test_SIRICheckStatusClient_Status_Not_Successful(t *testing.T) {
	assert := assert.New(t)
	partnerStatus, err := prepareSiriCheckStatusClient(t, "testdata/checkstatus-500.html")

	assert.Error(err, "SIRI CRITICAL: HTTP Content-Type text/html; charset=utf-8")
	assert.Equal(partnerStatus.OperationnalStatus, partners.OperationnalStatusDown)
}

func Test_SIRICheckStatusClientFactory_Validate(t *testing.T) {
	partner := &Partner{
		slug:           "partner",
		ConnectorTypes: []string{"siri-check-status-client"},
		connectors:     make(map[string]Connector),
		manager:        NewPartnerManager(nil),
	}
	partner.PartnerSettings = s.NewEmptyPartnerSettings(partner.UUIDGenerator)
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
