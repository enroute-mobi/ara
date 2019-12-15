package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
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
	urls := siri.SOAPClientUrls{
		Url:              siriPartner.partner.Setting("remote_url"),
		SubscriptionsUrl: siriPartner.partner.Setting("subscriptions.remote_url"),
		NotificationsUrl: siriPartner.partner.Setting("notifications.remote_url"),
	}
	if siriPartner.soapClient == nil || siriPartner.soapClient.SOAPClientUrls != urls {
		logger.Log.Debugf("Create SIRI SOAPClient to %s", urls.Url)
		siriPartner.soapClient = siri.NewSOAPClient(urls)
	}
	return siriPartner.soapClient
}

func (siriPartner *SIRIPartner) RequestorRef() string {
	return siriPartner.partner.ProducerRef()
}

func (siriPartner *SIRIPartner) SubscriberRef() string {
	return siriPartner.partner.Setting("local_credential")
}

func (siriPartner *SIRIPartner) Partner() *Partner {
	return siriPartner.partner
}

func (siriPartner *SIRIPartner) IdentifierGenerator(generatorName string) *IdentifierGenerator {
	formatString := siriPartner.partner.Setting(fmt.Sprintf("generators.%v", generatorName))
	if formatString == "" {
		formatString = defaultIdentifierGenerators[generatorName]
	}
	return NewIdentifierGeneratorWithUUID(formatString, siriPartner.UUIDConsumer)
}
