package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type CheckStatusClient interface {
	Status() (OperationnalStatus, error)
}

type TestCheckStatusClient struct {
	status OperationnalStatus
	Done   chan bool
}

type TestCheckStatusClientFactory struct{}

type SIRICheckStatusClient struct {
	model.ClockConsumer

	SIRIConnector
}

type SIRICheckStatusClientFactory struct{}

func NewTestCheckStatusClient() *TestCheckStatusClient {
	return &TestCheckStatusClient{
		status: OPERATIONNAL_STATUS_UP,
		Done:   make(chan bool, 1),
	}
}

func (connector *TestCheckStatusClient) Status() (OperationnalStatus, error) {
	connector.Done <- true
	return connector.status, nil
}

func (connector *TestCheckStatusClient) SetStatus(status OperationnalStatus) {
	connector.status = status
}

func (factory *TestCheckStatusClientFactory) Validate(apiPartner *APIPartner) bool {
	return true
}

func (factory *TestCheckStatusClientFactory) CreateConnector(partner *Partner) Connector {
	return NewTestCheckStatusClient()
}

func NewSIRICheckStatusClient(partner *Partner) *SIRICheckStatusClient {
	siriCheckStatusClient := &SIRICheckStatusClient{}
	siriCheckStatusClient.partner = partner
	return siriCheckStatusClient
}

func (connector *SIRICheckStatusClient) Status() (OperationnalStatus, error) {
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  connector.Clock().Now(),
		MessageIdentifier: connector.SIRIPartner().NewMessageIdentifier(),
	}

	response, err := connector.SIRIPartner().SOAPClient().CheckStatus(request)
	if err != nil {
		return OPERATIONNAL_STATUS_UNKNOWN, err
	}

	if response.Status() {
		return OPERATIONNAL_STATUS_UP, nil
	} else {
		return OPERATIONNAL_STATUS_DOWN, nil
	}
}

func (factory *SIRICheckStatusClientFactory) Validate(apiPartner *APIPartner) bool {
	ok := true
	if !apiPartner.IsSettingDefined("remote_url") {
		apiPartner.Errors = append(apiPartner.Errors, "SIRICheckStatusClient needs partner to have 'remote_url' setting defined")
		ok = false
	}
	if !apiPartner.IsSettingDefined("remote_credential") {
		apiPartner.Errors = append(apiPartner.Errors, "SIRICheckStatusClient needs partner to have 'remote_credential' setting defined")
		ok = false
	}
	return ok
}

func (factory *SIRICheckStatusClientFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRICheckStatusClient(partner)
}
