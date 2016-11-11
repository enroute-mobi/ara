package model

import (
	"github.com/af83/edwig/siri"
)

type CheckStatusClient interface {
	Status() (OperationnalStatus, error)
}

type SIRICheckStatusClient struct {
	ClockConsumer

	partner *SIRIPartner
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
		return UNKNOWN, err
	}

	if response.Status() {
		return UP, nil
	} else {
		return DOWN, nil
	}
}
