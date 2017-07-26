package core

import (
	"fmt"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIPartner struct {
	model.UUIDConsumer

	partner *Partner

	soapClient *siri.SOAPClient
}

func NewSIRIPartner(partner *Partner) *SIRIPartner {
	return &SIRIPartner{partner: partner}
}

func (siriPartner *SIRIPartner) SOAPClient() *siri.SOAPClient {
	if siriPartner.soapClient == nil || siriPartner.soapClient.URL() != siriPartner.partner.Setting("remote_url") {
		siriUrl := siriPartner.partner.Setting("remote_url")
		logger.Log.Debugf("Create SIRI SOAPClient to %s", siriUrl)
		siriPartner.soapClient = siri.NewSOAPClient(siriUrl)
	}
	return siriPartner.soapClient
}

func (siriPartner *SIRIPartner) RequestorRef() string {
	return siriPartner.partner.Setting("remote_credential")
}

func (siriPartner *SIRIPartner) Partner() *Partner {
	return siriPartner.partner
}

func (siriPartner *SIRIPartner) IdentifierGenerator(generatorName string) *IdentifierGenerator {
	formatString := siriPartner.partner.Setting(fmt.Sprintf("generators.%v", generatorName))
	if formatString == "" {
		formatString, _ = defaultIdentifierGenerators[generatorName]
	}
	return NewIdentifierGeneratorWithUUID(formatString, siriPartner.UUIDConsumer)
}
