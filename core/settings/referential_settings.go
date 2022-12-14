package settings

import (
	"strconv"
	"time"
)

const (
	MODEL_RELOAD_AT   = "model.reload_at"
	MODEL_PERSISTENCE = "model.persistence"
)

type ReferentialSettings struct {
	Settings
}

func NewReferentialSettings() (rs ReferentialSettings) {
	rs = ReferentialSettings{
		Settings: NewSettings(),
	}
	return
}

func (rs *ReferentialSettings) NextReloadAtSetting() (hour, minute int) {
	rs.m.RLock()
	r, _ := rs.s[MODEL_RELOAD_AT]
	rs.m.RUnlock()

	if len(r) != 5 {
		return 4, 0
	}
	hour, _ = strconv.Atoi(r[0:2])
	minute, _ = strconv.Atoi(r[3:5])
	return
}

func (rs *ReferentialSettings) ModelPersistenceDuration() (d time.Duration, ok bool) {
	rs.m.RLock()
	mp, ok := rs.s[MODEL_PERSISTENCE]
	rs.m.RUnlock()
	if !ok {
		return
	}
	d, _ = time.ParseDuration(mp)
	if d < 0 {
		d = 0
	}
	return -d, true
}
