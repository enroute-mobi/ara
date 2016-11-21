package core

const (
	SIRI_CHECK_STATUS_CLIENT_TYPE = "siri-check-status-client"
	TEST_CHECK_STATUS_CLIENT_TYPE = "test-check-status-client"
	TEST_VALIDATION_CONNECTOR     = "test-validation-connector"
)

type Connector interface{}

type ConnectorFactory interface {
	Validate(*APIPartner) bool
	CreateConnector(*Partner) Connector
}

func NewConnectorFactory(connectorType string) ConnectorFactory {
	switch connectorType {
	case SIRI_CHECK_STATUS_CLIENT_TYPE:
		return &SIRICheckStatusClientFactory{}
	case TEST_CHECK_STATUS_CLIENT_TYPE:
		return &TestCheckStatusClientFactory{}
	case TEST_VALIDATION_CONNECTOR:
		return &TestValidationFactory{}
	default:
		return nil
	}
}

type TestValidationFactory struct{}
type TestValidationConnector struct{}

func (factory *TestValidationFactory) Validate(apiPartner *APIPartner) bool {
	if apiPartner.Slug == PartnerSlug("InvalidSlug") {
		apiPartner.Errors = append(apiPartner.Errors, "Partner have an invalid slug")
		return false
	}
	return true
}

func (factory *TestValidationFactory) CreateConnector(partner *Partner) Connector {
	return &TestValidationFactory{}
}
