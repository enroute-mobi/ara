package model

const (
	SIRI_CHECK_STATUS_CLIENT_TYPE = "siri-check-status-client"
	TEST_CHECK_STATUS_CLIENT_TYPE = "test-check-status-client"
)

type Connector interface{}

type ConnectorFactory interface {
	Validate(*APIPartner) (string, bool)
	CreateConnector(*Partner) Connector
}

func NewConnectorFactory(connectorType string) ConnectorFactory {
	switch connectorType {
	case SIRI_CHECK_STATUS_CLIENT_TYPE:
		return &SIRICheckStatusClientFactory{}
	case TEST_CHECK_STATUS_CLIENT_TYPE:
		return &TestCheckStatusClientFactory{}
	default:
		return nil
	}
}
