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

func Test_SIRISituationExchangeSubscriptionCollector(t *testing.T) {
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
	situation := referential.Model().Situations().New()
	situation.SetCode(code)
	situation.Save()

	line := partners.Model().Lines().New()
	lineCode := model.NewCode("test_kind", "line value")
	line.SetCode(lineCode)
	partners.Model().Lines().Save(line)

	line2 := partners.Model().Lines().New()
	lineCode2 := model.NewCode("test_kind", "line value2")
	line2.SetCode(lineCode2)
	partners.Model().Lines().Save(line2)

	connector := NewSIRISituationExchangeSubscriptionCollector(partner)
	connector.SetSituationExchangeSubscriber(NewFakeSituationExchangeSubscriber(connector))

	connector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, lineCode)
	connector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, lineCode2)

	connector.Start()

	assert.Len(request.XMLSubscriptionSXEntries(), 1)
	assert.ElementsMatch(
		request.XMLSubscriptionSXEntries()[0].LineRef(),
		[]string{"line value", "line value2"},
	)
}
