package core

import (
	"io"
	"os"
	"testing"

	s "bitbucket.org/enroute-mobi/ara/core/settings"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"github.com/jbowtie/gokogiri/xml"
)

func BenchmarkHandleNotifyStopMonitoring(b *testing.B) {
	_, referential := newTestReferential(b)

	partners := NewPartnerManager(referential)
	partner := partners.New("slug")
	partner.PartnerSettings = s.NewPartnerSettings(partner.UUIDGenerator, map[string]string{
		"remote_code_space": "internal",
	})

	connector := NewSIRIStopMonitoringSubscriptionCollector(partner)
	connector.deletedSubscriptions = NewDeletedSubscriptions()

	f, err := os.Open("testdata/notify-stop-monitoring-large.xml")
	if err != nil {
		b.Fatal(err)
	}
	defer f.Close()
	content, err := io.ReadAll(f)
	if err != nil {
		b.Fatal(err)
	}

	doc, err := xml.Parse(content, xml.DefaultEncodingBytes, nil, xml.StrictParseOption, xml.DefaultEncodingBytes)
	if err != nil {
		b.Fatal(err)
	}
	notify := sxml.NewXMLNotifyStopMonitoring(doc.Root())

	subs := partner.Subscriptions().(*MemorySubscriptions)
	for _, delivery := range notify.StopMonitoringDeliveries() {
		subId := delivery.SubscriptionRef()
		if subId == "" {
			continue
		}
		sub := &Subscription{
			id:                  SubscriptionId(subId),
			kind:                StopMonitoringCollect,
			resourcesByCode:     make(map[string]*SubscribedResource),
			subscriptionOptions: make(map[string]string),
		}
		subs.Save(sub)
	}

	connector.HandleNotifyStopMonitoring(notify)

	if got := len(referential.Model().StopVisits().FindAll()); got != 888 {
		b.Fatalf("expected 888 StopVisits, got %d", got)
	}
	if got := len(referential.Model().StopAreas().FindAll()); got != 173 {
		b.Fatalf("expected 173 StopAreas, got %d", got)
	}

	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		connector.HandleNotifyStopMonitoring(notify)
	}
}
