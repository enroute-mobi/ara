package core

const (
	SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER     = "siri-stop-points-discovery-request-broadcaster"
	SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER           = "siri-general-message-request-broadcaster"
	SIRI_GENERAL_MESSAGE_REQUEST_COLLECTOR             = "siri-general-message-request-collector"
	TEST_GENERAL_MESSAGE_REQUEST_COLLECTOR             = "test-general-message-request-collector"
	SIRI_SERVICE_REQUEST_BROADCASTER                   = "siri-service-request-broadcaster"
	SIRI_STOP_MONITORING_REQUEST_COLLECTOR             = "siri-stop-monitoring-request-collector"
	TEST_STOP_MONITORING_REQUEST_COLLECTOR             = "test-stop-monitoring-request-collector"
	SIRI_STOP_MONITORING_REQUEST_BROADCASTER           = "siri-stop-monitoring-request-broadcaster"
	SIRI_STOP_MONITORING_DELIVERIES_RESPONSE_COLLECTOR = "siri-stop-monitoring-deliveries-response-collector"
	SIRI_CHECK_STATUS_CLIENT_TYPE                      = "siri-check-status-client"
	TEST_CHECK_STATUS_CLIENT_TYPE                      = "test-check-status-client"
	SIRI_CHECK_STATUS_SERVER_TYPE                      = "siri-check-status-server"
	TEST_VALIDATION_CONNECTOR                          = "test-validation-connector"
)

const (
	SIRI_PARTNER = "siri-partner"
)

type Connector interface{}

type BaseConnector struct {
	partner *Partner
}

func (connector *BaseConnector) Partner() *Partner {
	return connector.partner
}

type siriConnector struct {
	BaseConnector
}

func (connector *siriConnector) SIRIPartner() *SIRIPartner {
	if !connector.Partner().Context().IsDefined(SIRI_PARTNER) {
		connector.Partner().Context().SetValue(SIRI_PARTNER, NewSIRIPartner(connector.Partner()))
	}
	return connector.Partner().Context().Value(SIRI_PARTNER).(*SIRIPartner)
}

type ConnectorFactory interface {
	Validate(*APIPartner) bool
	CreateConnector(*Partner) Connector
}

func NewConnectorFactory(connectorType string) ConnectorFactory {
	switch connectorType {
	case SIRI_SERVICE_REQUEST_BROADCASTER:
		return &SIRIServiceRequestBroadcasterFactory{}
	case SIRI_STOP_MONITORING_REQUEST_COLLECTOR:
		return &SIRIStopMonitoringRequestCollectorFactory{}
	case TEST_STOP_MONITORING_REQUEST_COLLECTOR:
		return &TestStopMonitoringRequestCollectorFactory{}
	case SIRI_STOP_MONITORING_REQUEST_BROADCASTER:
		return &SIRIStopMonitoringRequestBroadcasterFactory{}
	case SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER:
		return &SIRIStopPointsDiscoveryRequestBroadcasterFactory{}
	case SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER:
		return &SIRIGeneralMessageRequestBroadcasterFactory{}
	case SIRI_GENERAL_MESSAGE_REQUEST_COLLECTOR:
		return &SIRIGeneralMessageRequestCollectorFactory{}
	case SIRI_CHECK_STATUS_CLIENT_TYPE:
		return &SIRICheckStatusClientFactory{}
	case TEST_CHECK_STATUS_CLIENT_TYPE:
		return &TestCheckStatusClientFactory{}
	case SIRI_CHECK_STATUS_SERVER_TYPE:
		return &SIRICheckStatusServerFactory{}
	case TEST_VALIDATION_CONNECTOR:
		return &TestValidationFactory{}
	case SIRI_STOP_MONITORING_DELIVERIES_RESPONSE_COLLECTOR:
		return &SIRIStopMonitoringSubscriptionCollectorFactory{}
	default:
		return nil
	}
}

type TestValidationFactory struct{}
type TestValidationConnector struct{}

func (factory *TestValidationFactory) Validate(apiPartner *APIPartner) bool {
	if apiPartner.Slug == PartnerSlug("InvalidSlug") {
		apiPartner.Errors.Add("slug", "Invalid format")
		return false
	}
	return true
}

func (factory *TestValidationFactory) CreateConnector(partner *Partner) Connector {
	return &TestValidationConnector{}
}
