package core

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

func Test_SIRIGeneralMessageSubscriptionCollector(t *testing.T) {
	request := &sxml.XMLSubscriptionRequest{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength <= 0 {
			t.Errorf("Request ContentLength should be zero")
		}
		body, _ := io.ReadAll(r.Body)
		var err error
		request, err = sxml.NewXMLSubscriptionRequestFromContent(body)
		if err != nil {
			t.Fatalf("%v", err)
		}
	}))
	defer ts.Close()

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)

	partners := NewPartnerManager(referential)

	partner := partners.New("slug")
	settings := map[string]string{
		"local_url":            "http://example.com/test/siri",
		"remote_url":           ts.URL,
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

	connector := NewSIRIGeneralMessageSubscriptionCollector(partner)
	connector.SetGeneralMessageSubscriber(NewFakeGeneralMessageSubscriber(connector))

	connector.RequestSituationUpdate(SITUATION_UPDATE_REQUEST_LINE, lineCode)
	connector.Start()

	if expected := "http://example.com/test/siri"; request.ConsumerAddress() != expected {
		t.Errorf("Wrong ConsumerAddress:\n got: %v\nwant: %v", request.ConsumerAddress(), expected)
	}

	if len(request.XMLSubscriptionGMEntries()) != 1 {
		t.Errorf("Wrong XMLSubscriptionEntries:\n got: %v\nwant: 1", len(request.XMLSubscriptionGMEntries()))
	}
}

func Test_SIRIGeneralMessageDeleteSubscriptionRequest(t *testing.T) {
	request := &sxml.XMLDeleteSubscriptionRequest{}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		request, _ = sxml.NewXMLDeleteSubscriptionRequestFromContent(body)
	}))
	defer ts.Close()

	referentials := NewMemoryReferentials()
	referential := referentials.New(ReferentialSlug("referential"))
	referential.model = model.NewMemoryModel()
	referentials.Save(referential)
	partners := NewPartnerManager(referential)

	partner := partners.New("slug")
	settings := map[string]string{
		"local_url":                          "http://example.com/test/siri",
		"remote_url":                         ts.URL,
		"remote_code_space":               "test_kind",
		"generators.subscription_identifier": "Subscription::%{id}::LOC",
	}
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, settings)
	partner.subscriptionManager = NewMemorySubscriptions(partner)
	partners.Save(partner)

	file, _ := os.Open("testdata/notify-general-message.xml")
	content, _ := io.ReadAll(file)

	connector := NewSIRIGeneralMessageSubscriptionCollector(partner)
	connector.deletedSubscriptions = NewDeletedSubscriptions()

	notify, _ := sxml.NewXMLNotifyGeneralMessageFromContent(content)

	connector.HandleNotifyGeneralMessage(notify)

	if expected := "6ba7b814-9dad-11d1-0-00c04fd430c8"; request.SubscriptionRef() != expected {
		t.Errorf("Wrong SubscriptionRef want : %v  got %v :", expected, request.SubscriptionRef())
	}
}
