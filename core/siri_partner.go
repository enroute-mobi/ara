package core

import (
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/siri"
)

type SIRIPartner struct {
	MessageIdentifierConsumer

	partner *Partner

	soapClient *siri.SOAPClient
}

func NewSIRIPartner(partner *Partner) *SIRIPartner {
	return &SIRIPartner{partner: partner}
}

func (connector *SIRIPartner) SOAPClient() *siri.SOAPClient {
	if connector.soapClient == nil || connector.soapClient.URL() != connector.partner.Setting("remote_url") {
		siriUrl := connector.partner.Setting("remote_url")
		logger.Log.Debugf("Create SIRI SOAPClient to %s", siriUrl)
		connector.soapClient = siri.NewSOAPClient(siriUrl)
	}
	return connector.soapClient
}

func (connector *SIRIPartner) RequestorRef() string {
	return connector.partner.Setting("remote_credential")
}

func (connector *SIRIPartner) Partner() *Partner {
	return connector.partner
}
