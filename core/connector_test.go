package core

import "testing"

func Test_Factories_CreateConnector(t *testing.T) {
	partner := &Partner{
		Settings: make(map[string]string),
		ConnectorTypes: []string{
			"siri-stop-monitoring-request-collector",
			"siri-check-status-client",
			"test-validation-connector",
		},
		connectors: make(map[string]Connector),
	}
	partner.Settings = map[string]string{
		"remote_url":           "remote_url",
		"remote_objectid_kind": "remote_objectid_kind",
	}
	apiPartner := partner.Definition()
	apiPartner.Validate()
	partner.SetDefinition(apiPartner)

	if !partner.isConnectorDefined("siri-stop-monitoring-request-collector") {
		t.Error("siri-stop-monitoring-request-collector connector should be initialized")
	}
	if !partner.isConnectorDefined("siri-check-status-client") {
		t.Error("siri-check-status-client connector should be initialized")
	}
	if !partner.isConnectorDefined("test-validation-connector") {
		t.Error("test-validation-connector connector should be initialized")
	}
}
