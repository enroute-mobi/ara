package core

import (
	"testing"
)

func Test_Factories_CreateConnector(t *testing.T) {
	partners := createTestPartnerManager()
	partner := partners.New("slug")

	partner.ConnectorTypes = []string{
		"siri-stop-monitoring-request-collector",
		"siri-stop-monitoring-request-broadcaster",
		"siri-check-status-client",
		"siri-check-status-server",
		"test-validation-connector",
		"siri-general-message-request-broadcaster",
		"siri-estimated-timetable-request-broadcaster",
		"siri-lines-discovery-request-broadcaster",
		"siri-situation-exchange-subscription-collector",
	}
	partner.RefreshConnectors()

	if _, ok := partner.Connector("siri-stop-monitoring-request-collector"); !ok {
		t.Error("siri-stop-monitoring-request-collector connector should be initialized")
	}
	if _, ok := partner.Connector("siri-stop-monitoring-request-broadcaster"); !ok {
		t.Error("siri-stop-monitoring-request-broadcaster connector should be initialized")
	}
	if _, ok := partner.Connector("siri-general-message-request-broadcaster"); !ok {
		t.Error("siri-general-message-request-broadcaster connector should be initialized")
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
	if _, ok := partner.Connector("test-validation-connector"); !ok {
		t.Error("test-validation-connector connector should be initialized")
	}
	if _, ok := partner.Connector("siri-estimated-timetable-request-broadcaster"); !ok {
		t.Error("siri-estimated-timetable-request-broadcaster connector should be initialized")
	}
	if _, ok := partner.Connector("siri-lines-discovery-request-broadcaster"); !ok {
		t.Error("siri-estimated-timetable-request-broadcaster connector should be initialized")
	}
	if _, ok := partner.Connector("siri-situation-exchange-subscription-collector"); !ok {
		t.Error("siri-situation-exchange-subscription-collector connector should be initialized")
	}
}
