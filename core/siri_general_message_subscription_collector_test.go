package core

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

func Test_SIRIGeneralMessageSubscriptionCollector(t *testing.T) {
	request := &siri.XMLSubscriptionRequest{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}
		body, _ := ioutil.ReadAll(r.Body)
		request, _ = siri.NewXMLSubscriptionRequestFromContent(body)
	}))
	defer ts.Close()

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)

	partner := &Partner{
		context: make(Context),
		Settings: map[string]string{
			"local_url":            "http://example.com/test/siri",
			"remote_url":           ts.URL,
			"remote_objectid_kind": "test_kind",
		},
		manager: partners,
	}
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	objectid := model.NewObjectID("test_kind", "value")
	situation := referential.Model().Situations().New()
	situation.SetObjectID(objectid)
	situation.Save()

	connector := NewSIRIGeneralMessageSubscriptionCollector(partner)

	situationUpdateEvent := NewSituationUpdateRequest(SituationUpdateRequestId(situation.Id()))
	connector.RequestSituationUpdate(situationUpdateEvent)

	if expected := "http://example.com/test/siri"; request.ConsumerAddress() != expected {
		t.Errorf("Wrong ConsumerAddress:\n got: %v\nwant: %v", request.ConsumerAddress(), expected)
	}

	if len(request.XMLSubscriptionGMEntries()) != 1 {
		t.Errorf("Wrong XMLSubscriptionEntries:\n got: %v\nwant: 1", len(request.XMLSubscriptionGMEntries()))
	}
}
