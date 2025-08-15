package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_SIRIFacilityMonitoringSubscriptionCollector(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	request := &sxml.XMLSubscriptionRequest{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}
		body, _ := io.ReadAll(r.Body)
		var err error
		request, err = sxml.NewXMLSubscriptionRequestFromContent(body)
		require.NoError(err)
	}))
	defer ts.Close()

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewTestMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)

	partner := partners.New("slug")
	settings := map[string]string{
		"local_url":         "http://example.com/test/siri",
		"remote_url":        ts.URL,
		"remote_code_space": "test_kind",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	code := model.NewCode("test_kind", "value")
	facility := referential.Model().Facilities().New()
	facility.SetCode(code)
	facility.Save()

	connector := NewSIRIFacilityMonitoringSubscriptionCollector(partner)
	connector.SetFacilityMonitoringSubscriber(NewFakeFacilityMonitoringSubscriber(connector))

	facilityUpdateRequest := NewFacilityUpdateRequest(facility.Id())
	connector.RequestFacilityUpdate(facilityUpdateRequest)

	connector.facilityMonitoringSubscriber.Start()

	assert.Len(request.XMLSubscriptionFMEntries(), 1)
	assert.ElementsMatch(
		request.XMLSubscriptionFMEntries()[0].FacilityRefs(),
		[]string{"value"},
	)
}
