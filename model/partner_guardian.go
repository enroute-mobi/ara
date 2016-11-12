package model

import "time"

type PartnersGuardian struct {
	ClockConsumer

	partners Partners
}

func NewPartnersGuardian(partners Partners) *PartnersGuardian {
	return &PartnersGuardian{partners: partners}
}

func (guardian *PartnersGuardian) Run() {
	ticker := time.NewTicker(30 * time.Second)
	for _ = range ticker.C {
		for _, partner := range guardian.partners.FindAll() {
			go partner.CheckStatus()
		}
	}
}
