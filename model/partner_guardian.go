package model

import (
	"time"

	"github.com/af83/edwig/logger"
)

type PartnersGuardian struct {
	ClockConsumer

	partners Partners
}

func NewPartnersGuardian(partners Partners) *PartnersGuardian {
	return &PartnersGuardian{partners: partners}
}

func (guardian *PartnersGuardian) Start() {
	logger.Log.Debugf("Start partners guardian")

	go guardian.Run()
}

func (guardian *PartnersGuardian) Run() {
	ticker := time.NewTicker(30 * time.Second)
	for _ = range ticker.C {
		logger.Log.Debugf("Check partners status")
		for _, partner := range guardian.partners.FindAll() {
			go guardian.checkPartnerStatus(partner)
		}
	}
}

func (guardian *PartnersGuardian) checkPartnerStatus(partner *Partner) {
	defer func() {
		if r := recover(); r != nil {
			logger.Log.Printf("Recovered error %v in checkPartnerStatus for partner %#v", r, partner)
		}
	}()

	partner.CheckStatus()
}
