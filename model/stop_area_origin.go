package model

import (
	"encoding/json"
	"sync"
)

type StopAreaOrigins struct {
	sync.RWMutex

	partners map[string]bool
}

func NewStopAreaOrigins() StopAreaOrigins {
	return StopAreaOrigins{partners: make(map[string]bool)}
}

func (origins *StopAreaOrigins) MarshalJSON() ([]byte, error) {
	origins.RLock()
	aux := make(map[string]bool)
	for partner, status := range origins.partners {
		aux[partner] = status
	}
	origins.RUnlock()
	return json.Marshal(aux)
}

func (origins *StopAreaOrigins) SetOriginsFromMap(originMap map[string]bool) {
	origins.Lock()
	origins.partners = make(map[string]bool)
	for partner, status := range originMap {
		origins.partners[partner] = status
	}
	origins.Unlock()
}

func (origins *StopAreaOrigins) Copy() *StopAreaOrigins {
	cpy := NewStopAreaOrigins()

	origins.RLock()
	for key, value := range origins.partners {
		cpy.partners[key] = value
	}
	origins.RUnlock()
	return &cpy
}

func (origins *StopAreaOrigins) Origin(partner string) (status bool, present bool) {
	origins.RLock()
	status, present = origins.partners[partner]
	origins.RUnlock()
	return
}

func (origins *StopAreaOrigins) AllOrigin() map[string]bool {
	a := make(map[string]bool)
	origins.RLock()
	for k, v := range origins.partners {
		a[k] = v
	}
	origins.RUnlock()
	return a
}

func (origins *StopAreaOrigins) NewOrigin(partner string) {
	origins.SetPartnerStatus(partner, true)
}

func (origins *StopAreaOrigins) SetPartnerStatus(partner string, status bool) {
	origins.Lock()
	origins.partners[partner] = status
	origins.Unlock()
}

func (origins *StopAreaOrigins) Monitored() (monitored bool) {
	origins.RLock()
	monitored = true
	for _, status := range origins.partners {
		if !status {
			monitored = false
			break
		}
	}
	origins.RUnlock()
	return
}

func (origins *StopAreaOrigins) PartnersKO() (partnersSlice []string) {
	origins.RLock()
	for partner, status := range origins.partners {
		if !status {
			partnersSlice = append(partnersSlice, partner)
		}
	}
	origins.RUnlock()
	return
}

func (origins *StopAreaOrigins) PartnersLost(comparedOrigins *StopAreaOrigins) (partnersLost []string, r bool) {
	if comparedOrigins == nil {
		return
	}
	origins.RLock()
	comparedOrigins.RLock()
	for partner, status := range origins.partners {
		comparedStatus, _ := comparedOrigins.partners[partner]
		if status && !comparedStatus {
			partnersLost = append(partnersLost, partner)
			r = true
		}
	}
	origins.RUnlock()
	comparedOrigins.RUnlock()
	return
}
