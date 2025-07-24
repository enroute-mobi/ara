package core

import (
	"testing"

	"bitbucket.org/enroute-mobi/ara/clock"
	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/gtfs"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/uuid"
	"github.com/stretchr/testify/assert"
)

func Test_ServiceAlertsBroadcaster_HandleGtfs_WithEmptyAffectsSituations(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New("referential")
	partner := referential.Partners().New("partner")
	partner.SetUUIDGenerator(uuid.NewFakeUUIDGenerator())
	settings := map[string]string{
		"remote_code_space": "codeSpace",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	connector := NewServiceAlertsBroadcaster(partner)
	connector.SetClock(clock.NewFakeClock())
	connector.Start()

	code := model.NewCode("codeSpace", "saId")
	situation := referential.Model().Situations().New()
	situation.SetCode(code)
	situation.Save()

	gtfsFeed := &gtfs.FeedMessage{}
	connector.HandleGtfs(gtfsFeed)

	assert.Len(gtfsFeed.Entity, 1, "Entity should not be nil")
	assert.NotNil(gtfsFeed.Entity[0].Alert, "Should send Alert in Entity with Situation without Affect")
	assert.Nil(gtfsFeed.Entity[0].Alert.InformedEntity)
}
