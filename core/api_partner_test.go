package core

import (
	"encoding/json"
	"testing"

	e "bitbucket.org/enroute-mobi/ara/core/api_errors"
	ps "bitbucket.org/enroute-mobi/ara/core/partner_settings"
)

func Test_APIPartner_SetFactories(t *testing.T) {
	partner := &APIPartner{
		ConnectorTypes: []string{"unexistant-factory", "test-check-status-client"},
		factories:      make(map[string]ConnectorFactory),
	}
	partner.setFactories()

	if len(partner.factories) != 1 {
		t.Errorf("Factories should have been successfully created by setFactories")
	}
}

func Test_APIPartner_Validate(t *testing.T) {
	partners := createTestPartnerManager() // in core/partner_test.go
	// Check empty Slug
	apiPartner := &APIPartner{
		manager: partners,
	}
	valid := apiPartner.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiPartner.Errors) != 1 {
		t.Errorf("apiPartner Errors should not be empty")
	}
	if len(apiPartner.Errors.Get("Slug")) != 1 || apiPartner.Errors.Get("Slug")[0] != e.ERROR_BLANK {
		t.Errorf("apiPartner should have Error for Slug, got %v", apiPartner.Errors)
	}

	// Check wrong format slug
	apiPartner.Slug = "Wrong_format"
	valid = apiPartner.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiPartner.Errors) != 1 {
		t.Errorf("apiPartner Errors should not be empty")
	}
	if len(apiPartner.Errors.Get("Slug")) != 1 || apiPartner.Errors.Get("Slug")[0] != e.ERROR_SLUG_FORMAT {
		t.Errorf("apiPartner should have Error for Slug, got %v", apiPartner.Errors)
	}

	jsonBytes, _ := json.Marshal(apiPartner)
	expected := `{"Slug":"Wrong_format","Errors":{"Slug":["Invalid format: only lowercase alphanumeric characters and _"]}}`
	if string(jsonBytes) != expected {
		t.Fatalf("Invalid JSON, expected %v, got %v", expected, string(jsonBytes))
	}

	// Check Already Used Slug and local_credential
	partner := partners.New("slug")
	partner.SetSetting("local_credential", "cred")
	partners.Save(partner)

	apiPartner = &APIPartner{
		Slug:     "slug",
		Settings: map[string]string{"local_credential": "cred"},
		manager:  partners,
	}
	valid = apiPartner.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiPartner.Errors) != 2 {
		t.Errorf("apiPartner Errors should not be empty")
	}
	if len(apiPartner.Errors.Get("Slug")) != 1 || apiPartner.Errors.Get("Slug")[0] != e.ERROR_UNIQUE {
		t.Errorf("apiPartner should have Error for Slug, got %v", apiPartner.Errors)
	}
	if len(apiPartner.Errors.GetSettingError(ps.LOCAL_CREDENTIAL)) != 1 || apiPartner.Errors.GetSettingError(ps.LOCAL_CREDENTIAL)[0] != e.ERROR_UNIQUE {
		t.Errorf("apiPartner should have Error for local_credential, got %v", apiPartner.Errors)
	}

	jsonBytes, _ = json.Marshal(apiPartner)
	expected = `{"Slug":"slug","Settings":{"local_credential":"cred"},"Errors":{"Settings":{"local_credential":["Is already in use"]},"Slug":["Is already in use"]}}`
	if string(jsonBytes) != expected {
		t.Fatalf("Invalid JSON, expected %v, got %v", expected, string(jsonBytes))
	}

	// Check ok
	apiPartner = &APIPartner{
		Slug:     "slug_2",
		Settings: map[string]string{"local_credential": "cred2"},
		manager:  partners,
	}
	valid = apiPartner.Validate()

	if !valid {
		t.Errorf("Validate should return true")
	}
	if len(apiPartner.Errors) != 0 {
		t.Errorf("apiPartner Errors should be empty, got %v", apiPartner.Errors)
	}

	// Check settings errors
	apiPartner = &APIPartner{
		Slug:           "",
		Settings:       map[string]string{},
		ConnectorTypes: []string{SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER},
		manager:        partners,
		factories:      make(map[string]ConnectorFactory),
	}
	valid = apiPartner.Validate()

	if valid {
		t.Errorf("Validate should return false")
	}
	if len(apiPartner.Errors) != 2 {
		t.Errorf("apiPartner Errors should not be empty, got %v", apiPartner.Errors)
	}
	if len(apiPartner.Errors.GetSettings()) != 2 {
		t.Errorf("apiPartner Setting Errors should have 2 errors, got %v", apiPartner.Errors.GetSettings())
	}
}
