package core

import (
	"time"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type PartnersGuardian struct {
	model.ClockConsumer

	stop     chan struct{}
	partners Partners
}

func NewPartnersGuardian(partners Partners) *PartnersGuardian {
	return &PartnersGuardian{partners: partners}
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
			for _, partner := range guardian.partners.FindAll() {
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

	if partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_UNKNOWN || partnerStatus.OperationnalStatus == OPERATIONNAL_STATUS_DOWN || partnerStatus.ServiceStartedAt != partner.PartnerStatus.ServiceStartedAt {
		partner.PartnerStatus = partnerStatus

		if partner.Setting("subscription.persistent") == "true" {
			partner.Subscriptions().CancelCollectSubscriptions()
			return false
		}

		partner.CancelSubscriptions()

		return false
	}

	partner.PartnerStatus = partnerStatus
	return true
}

func (guardian *PartnersGuardian) checkSubscriptionsTerminatedTime(partner *Partner) {
	if partner.Subscriptions() == nil {
		return
	}

	for _, sub := range partner.Subscriptions().FindAll() {
		for key, value := range sub.ResourcesByObjectID() {
			if !value.SubscribedUntil.Before(guardian.Clock().Now()) || value.SubscribedAt.IsZero() {
				continue
			}
			sub.DeleteResource(key)
			logger.Log.Printf("Deleting ressource %v from subscription with id %v", key, sub.Id())
		}
		if len(sub.ResourcesByObjectID()) == 0 {
			sub.Delete()
		}
	}
}
