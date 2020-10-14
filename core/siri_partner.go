package core

import (
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type SIRIPartner struct {
	partner *Partner

	soapClient *siri.SOAPClient
}

func NewSIRIPartner(partner *Partner) *SIRIPartner {
	return &SIRIPartner{partner: partner}
}

func (siriPartner *SIRIPartner) SOAPClient() *siri.SOAPClient {
	urls := siri.SOAPClientUrls{
		Url:              siriPartner.partner.Setting(REMOTE_URL),
		SubscriptionsUrl: siriPartner.partner.Setting(SUBSCRIPTIONS_REMOTE_URL),
		NotificationsUrl: siriPartner.partner.Setting(NOTIFICATIONS_REMOTE_URL),
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
	return siriPartner.partner.Setting(LOCAL_CREDENTIAL)
}

func (siriPartner *SIRIPartner) Partner() *Partner {
	return siriPartner.partner
}
