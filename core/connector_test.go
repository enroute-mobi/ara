package core

import "testing"

func Test_Factories_CreateConnector(t *testing.T) {
	partner := &Partner{
		Settings: make(map[string]string),
		ConnectorTypes: []string{
			"siri-service-request-broadcaster",
			"siri-stop-monitoring-request-collector",
			"siri-stop-monitoring-request-broadcaster",
			"siri-check-status-client",
			"siri-check-status-server",
			"test-validation-connector",
		},
		connectors: make(map[string]Connector),
		manager:    NewPartnerManager(nil),
	}
	partner.RefreshConnectors()

	if _, ok := partner.Connector("siri-service-request-broadcaster"); !ok {
		t.Error("siri-service-request-broadcaster connector should be initialized")
	}
	if _, ok := partner.Connector("siri-stop-monitoring-request-collector"); !ok {
		t.Error("siri-stop-monitoring-request-collector connector should be initialized")
	}
	if _, ok := partner.Connector("siri-stop-monitoring-request-broadcaster"); !ok {
		t.Error("siri-stop-monitoring-request-broadcaster connector should be initialized")
	}
	if _, ok := partner.Connector("siri-check-status-client"); !ok {
		t.Error("siri-check-status-client connector should be initialized")
	}
	if _, ok := partner.Connector("siri-check-status-server"); !ok {
		t.Error("siri-check-status-server connector should be initialized")
	}
	if _, ok := partner.Connector("test-validation-connector"); !ok {
		t.Error("test-validation-connector connector should be initialized")
	}
}
