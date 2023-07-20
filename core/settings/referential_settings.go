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
	DEFAULT_MODEL_REFRESH_TIME = 50_000_000_000
	MINIMUM_MODEL_REFRESH_TIME = 30_000_000_000
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

var loggerObjectId = regexp.MustCompile(`^([^:]+):(.*)$`)

func (rs *ReferentialSettings) LoggerVerboseStopAreas() []model.ObjectID {
	rs.m.RLock()
	setting, ok := rs.s[LOGGER_VERBOSE_STOP_AREAS]
	rs.m.RUnlock()

	if !ok {
		return []model.ObjectID{}
	}

	parsedSetting := loggerObjectId.FindStringSubmatch(setting)
	if len(parsedSetting) == 0 {
		return []model.ObjectID{}
	}

	kind := parsedSetting[1]
	value := parsedSetting[2]

	return []model.ObjectID{model.NewObjectID(kind, value)}
}

func (rs *ReferentialSettings) ModelRefreshTime() (d time.Duration) {
	rs.m.RLock()
	mp, ok := rs.s[MODEL_REFRESH_TIME]
	rs.m.RUnlock()
	if !ok {
		return time.Duration(DEFAULT_MODEL_REFRESH_TIME)
	}

	d, _ = time.ParseDuration(mp)
	if minDuration := time.Duration(MINIMUM_MODEL_REFRESH_TIME); d < minDuration {
		d = minDuration
	}
	return d
}
