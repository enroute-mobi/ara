package core

import (
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/monitoring"
	"cloud.google.com/go/civil"
)

type PartnersGuardian struct {
	clock.ClockConsumer

	stop        chan struct{}
	referential *Referential
}

func NewPartnersGuardian(referential *Referential) *PartnersGuardian {
	return &PartnersGuardian{referential: referential}
}

func (guardian *PartnersGuardian) Start() {
	logger.Log.Debugf("Start partners guardian")

	guardian.stop = make(chan struct{})
	go guardian.Run()
}

func (guardian *PartnersGuardian) Stop() {
	if guardian.stop != nil {
		close(guardian.stop)
	}
}

func (guardian *PartnersGuardian) Run() {
	partnerChannel := make(chan *Partner)
	guardian.listen(partnerChannel)
	for {
		select {
		case <-guardian.stop:
			close(partnerChannel)
			logger.Log.Debugf("Stop Partners Guardian")
			return
		case <-guardian.Clock().After(30 * time.Second):
			logger.Log.Debugf("Check partners status")
			for _, partner := range guardian.referential.Partners().FindAll() {
				partnerChannel <- partner
			}
		}
	}
}

func (guardian *PartnersGuardian) listen(partnerChannel <-chan *Partner) {
	go func() {
		for p := range partnerChannel {
			guardian.routineWork(p)
		}
	}()
}

func (guardian *PartnersGuardian) routineWork(partner *Partner) {
	defer monitoring.HandlePanic()

	s := guardian.checkPartnerStatus(partner)
	if s {
		guardian.checkSubscriptionsTerminatedTime(partner)
	}

	guardian.checkPartnerDiscovery(partner)
}

// Returns true if we need to check the subscriptions (false if subscriptions are deleted)
func (guardian *PartnersGuardian) checkPartnerStatus(partner *Partner) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Printf("Recovered error %v in checkPartnerStatus for partner %#v", r, partner)
		}
	}()

	partnerStatus, _ := partner.CheckStatus()

	// Do nothing if partner status is unknown
	if partner.PartnerStatus.OperationnalStatus != OPERATIONNAL_STATUS_UNKNOWN && partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UNKNOWN && partner.PartnerStatus.RetryCount < partner.MaximumChechstatusRetry() {
		partner.PartnerStatus.RetryCount += 1
		logger.Log.Debugf("Unknow Status, %v retry", partner.PartnerStatus.RetryCount)
		return partner.PartnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UP
	}

	if partnerStatus.OperationnalStatus != partner.PartnerStatus.OperationnalStatus {
		logger.Log.Debugf("Partner %v status changed after a CheckStatus: was %v, now is %v", partner.Slug(), partner.PartnerStatus.OperationnalStatus, partnerStatus.OperationnalStatus)
		guardian.referential.CollectManager().HandlePartnerStatusChange(string(partner.Slug()), partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UP)
		partnerEvent := &audit.BigQueryPartnerEvent{
			Timestamp:                guardian.Clock().Now(),
			Slug:                     string(partner.Slug()),
			PartnerUUID:              string(partner.Id()),
			PreviousStatus:           string(partner.PartnerStatus.OperationnalStatus),
			PreviousServiceStartedAt: civil.DateTimeOf(partner.PartnerStatus.ServiceStartedAt),
			NewStatus:                string(partnerStatus.OperationnalStatus),
			NewServiceStartedAt:      civil.DateTimeOf(partnerStatus.ServiceStartedAt),
		}

		audit.CurrentBigQuery(string(guardian.referential.Slug())).WriteEvent(partnerEvent)
	}

	if partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UP && partnerStatus.ServiceStartedAt != partner.PartnerStatus.ServiceStartedAt {
		partner.PartnerStatus = partnerStatus
		partner.Subscriptions().CancelSubscriptionsResourcesBefore(partnerStatus.ServiceStartedAt)
		partner.lastDiscovery = time.Time{} // Reset discoveries if distant partner reloaded
		return false
	}

	if partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UNKNOWN || partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_DOWN {
		partner.PartnerStatus.OperationnalStatus = partnerStatus.OperationnalStatus
		partner.PartnerStatus.RetryCount = 0
		partner.lastDiscovery = time.Time{} // Reset discoveries if distant partner is down

		collectPersistent := partner.PersistentCollect()
		if !collectPersistent {
			partner.Subscriptions().CancelCollectSubscriptions()
		}

		broadcastPersistent := partner.PersistentBroadcastSubscriptions()
		if !broadcastPersistent {
			partner.Subscriptions().CancelBroadcastSubscriptions()
		}

		return (collectPersistent || broadcastPersistent)
	}

	partner.PartnerStatus = partnerStatus
	return true
}

func (guardian *PartnersGuardian) checkSubscriptionsTerminatedTime(partner *Partner) {
	if partner.Subscriptions() == nil {
		return
	}

	for _, sub := range partner.Subscriptions().FindAll() {
		for key, value := range sub.ResourcesByObjectIDCopy() {
			if !value.SubscribedUntil.Before(guardian.Clock().Now()) || value.SubscribedAt().IsZero() {
				continue
			}
			sub.DeleteResource(key)
			logger.Log.Printf("%v from %v: Deleting ressource %v from subscription with id %v. SubscribedUntil %v befor Clock.Now %v ", partner.Slug(), partner.Referential().Slug(), key, sub.Id(), value.SubscribedUntil, guardian.Clock().Now())
		}
		if sub.ResourcesLen() == 0 {
			sub.Delete()
		}
	}
}

func (guardian *PartnersGuardian) checkPartnerDiscovery(partner *Partner) {
	if partner.OperationnalStatus() != OPERATIONNAL_STATUS_UP {
		return
	}

	if partner.LastDiscovery().IsZero() || partner.LastDiscovery().Before(guardian.Clock().Now().Add(partner.DiscoveryInterval())) {
		partner.Discover()
	}
}
