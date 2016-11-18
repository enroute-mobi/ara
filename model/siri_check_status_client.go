package model

import (
	"github.com/af83/edwig/siri"
)

type CheckStatusClient interface {
	Status() (OperationnalStatus, error)
}

const (
	SIRI_CHECK_STATUS_CLIENT_TYPE = "siri-check-status-client"
	TEST_CHECK_STATUS_CLIENT_TYPE = "test-check-status-client"
)

type TestCheckStatusClient struct {
	status OperationnalStatus
	Done   chan bool
}

type SIRICheckStatusClient struct {
	ClockConsumer

	partner *SIRIPartner
}

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

func NewSIRICheckStatusClient(partner *SIRIPartner) *SIRICheckStatusClient {
	return &SIRICheckStatusClient{partner: partner}
}

func (connector *SIRICheckStatusClient) Status() (OperationnalStatus, error) {
	request := &siri.SIRICheckStatusRequest{
		RequestorRef:      connector.partner.RequestorRef(),
		RequestTimestamp:  connector.Clock().Now(),
		MessageIdentifier: connector.partner.NewMessageIdentifier(),
	}

	response, err := connector.partner.SOAPClient().CheckStatus(request)
	if err != nil {
		return OPERATIONNAL_STATUS_UNKNOWN, err
	}

	if response.Status() {
		return OPERATIONNAL_STATUS_UP, nil
	} else {
		return OPERATIONNAL_STATUS_DOWN, nil
	}
}
