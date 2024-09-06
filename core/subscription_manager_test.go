package core

import (
	"testing"

	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"github.com/stretchr/testify/assert"
)

func Test_FindOrCreateByKind_StopMonitoringAndVehicleMonitoringCollect(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()

	for _, kind := range []string{"StopMonitoringCollect", "VehicleMonitoringCollect"} {
		referential := referentials.New(ReferentialSlug("referential"))
		partner := referential.partners.New("test")
		partner.Subscriptions().FindOrCreateByKind(kind)

		subscriptions := partner.Subscriptions().FindAll()
		assert.Len(subscriptions, 1)
		assert.Equal(subscriptions[0].kind, kind)
		assert.Zero(subscriptions[0].ResourcesLen())
	}
}

func Test_FindOrCreateByKind_WithMaximumResources(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	partner := referential.partners.New("test")

	settings := map[string]string{
		"subscriptions.maximum_resources": "1",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	referential.partners.Save(partner)

	partner.Subscriptions().FindOrCreateByKind("kind")

	subscriptions := partner.Subscriptions().FindAll()
	assert.Len(subscriptions, 1)
	assert.Equal(subscriptions[0].kind, "kind")
	assert.Zero(subscriptions[0].ResourcesLen())
}

func Test_FindOrCreateByKind_WithExistingSubscription(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	partner := referential.partners.New("test")
	referential.partners.Save(partner)

	subscription := partner.Subscriptions().New("kind")
	subscription.Save()

	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
	}
	subscription.CreateAndAddNewResource(reference)

	foundSubscription := partner.Subscriptions().FindOrCreateByKind("kind")
	assert.Equal(foundSubscription.Id(), subscription.Id(), "should return already existing subscription")
}

func Test_FindOrCreateByKind_WithExistingSubscription_AlreadySubscribed(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	partner := referential.partners.New("test")
	referential.partners.Save(partner)

	subscription := partner.Subscriptions().New("kind")
	subscription.subscribed = true
	subscription.Save()

	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
	}
	subscription.CreateAndAddNewResource(reference)

	foundSubscription := partner.Subscriptions().FindOrCreateByKind("kind")
	assert.Len(partner.Subscriptions().FindAll(), 2, "a new subscription should be created")
	assert.NotEqual(foundSubscription.Id(), subscription.Id())
	assert.Zero(foundSubscription.ResourcesLen())
}

func Test_FindOrCreateByKind_WithExistingSubscription_WithResourceBelowMaximumResource(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	partner := referential.partners.New("test")
	settings := map[string]string{
		"subscriptions.maximum_resources": "3",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	referential.partners.Save(partner)

	subscription := partner.Subscriptions().New("kind")
	subscription.Save()

	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
	}
	subscription.CreateAndAddNewResource(reference)

	obj2 := model.NewCode("internal", "Value2")
	reference = model.Reference{
		Code: &obj2,
	}
	subscription.CreateAndAddNewResource(reference)

	foundSubscription := partner.Subscriptions().FindOrCreateByKind("kind")

	assert.Len(partner.Subscriptions().FindAll(), 1, "no new subscription should be created")
	assert.Equal(foundSubscription.Id(), subscription.Id(), "should return already existing subscription")
	assert.Equal(foundSubscription.ResourcesLen(), 2)
}

func Test_FindOrCreateByKind_WithExistingSubscription_WithResourceEqualToMaximumResource(t *testing.T) {
	assert := assert.New(t)

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	partner := referential.partners.New("test")
	settings := map[string]string{
		"subscriptions.maximum_resources": "2",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	referential.partners.Save(partner)

	subscription := partner.Subscriptions().New("kind")
	subscription.subscribed = true
	subscription.Save()

	obj := model.NewCode("internal", "Value")
	reference := model.Reference{
		Code: &obj,
	}
	subscription.CreateAndAddNewResource(reference)

	obj2 := model.NewCode("internal", "Value2")
	reference = model.Reference{
		Code: &obj2,
	}
	subscription.CreateAndAddNewResource(reference)

	foundSubscription := partner.Subscriptions().FindOrCreateByKind("kind")

	assert.Len(partner.Subscriptions().FindAll(), 2, "a new subscription should be created")
	assert.NotEqual(foundSubscription.Id(), subscription.Id())
	assert.Zero(foundSubscription.ResourcesLen())
}
