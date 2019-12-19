package core

import (
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/model"
)

type PartnersGuardian struct {
	model.ClockConsumer

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
	for {
		select {
		case <-guardian.stop:
			logger.Log.Debugf("Stop Partners Guardian")
			return
		case <-guardian.Clock().After(30 * time.Second):
			logger.Log.Debugf("Check partners status")
			for _, partner := range guardian.referential.Partners().FindAll() {
				go guardian.routineWork(partner)
			}
		}
	}
}

func (guardian *PartnersGuardian) routineWork(partner *Partner) {
	s := guardian.checkPartnerStatus(partner)
	if !s {
		return
	}

	guardian.checkSubscriptionsTerminatedTime(partner)
}

func (guardian *PartnersGuardian) checkPartnerStatus(partner *Partner) bool {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Printf("Recovered error %v in checkPartnerStatus for partner %#v", r, partner)
		}
	}()

	partnerStatus, _ := partner.CheckStatus()

	if partnerStatus.OperationnalStatus != partner.PartnerStatus.OperationnalStatus {
		logger.Log.Debugf("Partner %v status changed after a CheckStatus: was %v, now is %v", partner.Slug(), partner.PartnerStatus.OperationnalStatus, partnerStatus.OperationnalStatus)
		guardian.referential.CollectManager().HandlePartnerStatusChange(string(partner.Slug()), partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UP)
	}

	if partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UP && partnerStatus.ServiceStartedAt != partner.PartnerStatus.ServiceStartedAt {
		partner.PartnerStatus = partnerStatus
		partner.Subscriptions().CancelSubscriptions()
		return false
	}

	if partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UNKNOWN || partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_DOWN {
		partner.PartnerStatus.OperationnalStatus = partnerStatus.OperationnalStatus

		collectPersistent, _ := strconv.ParseBool(partner.Setting("collect.subscriptions.persistent"))
		if !collectPersistent {
			partner.Subscriptions().CancelCollectSubscriptions()
		}

		broadcastPersistent, _ := strconv.ParseBool(partner.Setting("broadcast.subscriptions.persistent"))
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
			if !value.SubscribedUntil.Before(guardian.Clock().Now()) || value.SubscribedAt.IsZero() {
				continue
			}
			sub.DeleteResource(key)
			logger.Log.Printf("Deleting ressource %v from subscription with id %v", key, sub.Id())
		}
		if sub.ResourcesLen() == 0 {
			sub.Delete()
		}
	}
}
