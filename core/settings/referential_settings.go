package settings

import (
	"regexp"
	"strconv"
	"time"

	"bitbucket.org/enroute-mobi/ara/model"
)

const (
	MODEL_RELOAD_AT            = "model.reload_at"
	MODEL_PERSISTENCE          = "model.persistence"
	MODEL_REFRESH_TIME         = "model.refresh_time"
	LOGGER_VERBOSE_STOP_AREAS  = "logger.verbose.stop_areas"
	DEFAULT_MODEL_REFRESH_TIME = 50 * time.Second
	MINIMUM_MODEL_REFRESH_TIME = 30 * time.Second
	DEFAULT_MODEL_PERSISTENCE  = 3 * time.Hour
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
	r := rs.s[MODEL_RELOAD_AT]
	rs.m.RUnlock()

	if len(r) != 5 {
		return 4, 0
	}
	hour, _ = strconv.Atoi(r[0:2])
	minute, _ = strconv.Atoi(r[3:5])
	return
}

func (rs *ReferentialSettings) ModelPersistenceDuration() (d time.Duration) {
	rs.m.RLock()
	mp, ok := rs.s[MODEL_PERSISTENCE]
	rs.m.RUnlock()
	if !ok {
		return -DEFAULT_MODEL_PERSISTENCE
	}

	d, _ = time.ParseDuration(mp)
	if d < 0 {
		d = DEFAULT_MODEL_PERSISTENCE
	}
	return -d
}

var loggerCode = regexp.MustCompile(`^([^:]+):(.*)$`)

func (rs *ReferentialSettings) LoggerVerboseStopAreas() []model.Code {
	rs.m.RLock()
	setting, ok := rs.s[LOGGER_VERBOSE_STOP_AREAS]
	rs.m.RUnlock()

	if !ok {
		return []model.Code{}
	}

	parsedSetting := loggerCode.FindStringSubmatch(setting)
	if len(parsedSetting) == 0 {
		return []model.Code{}
	}

	kind := parsedSetting[1]
	value := parsedSetting[2]

	return []model.Code{model.NewCode(kind, value)}
}

func (rs *ReferentialSettings) ModelRefreshTime() (d time.Duration) {
	rs.m.RLock()
	mp, ok := rs.s[MODEL_REFRESH_TIME]
	rs.m.RUnlock()
	if !ok {
		return DEFAULT_MODEL_REFRESH_TIME
	}

	d, _ = time.ParseDuration(mp)
	if d < MINIMUM_MODEL_REFRESH_TIME {
		d = MINIMUM_MODEL_REFRESH_TIME
	}
	return d
}
