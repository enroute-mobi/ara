package model

import (
	"github.com/af83/edwig/siri"
)

type SIRIPartner struct {
	MessageIdentifierConsumer

	partner *Partner
}

func NewSIRIPartner(partner *Partner) *SIRIPartner {
	return &SIRIPartner{partner: partner}
}

func (connector *SIRIPartner) SOAPClient() *siri.SOAPClient {
	return nil
}

func (connector *SIRIPartner) RequestorRef() string {
	return ""
}
