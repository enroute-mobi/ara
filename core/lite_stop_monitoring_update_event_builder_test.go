package core

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"testing"

	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/slite"
	"github.com/stretchr/testify/assert"
)

func getsmlite(t *testing.T, filePath string) *slite.SIRILiteStopMonitoring {
	file, err := os.Open(filePath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	dest := &slite.SIRILiteStopMonitoring{}
	err = json.NewDecoder(bytes.NewReader(content)).Decode(&dest)
	if err != nil {
		t.Fatal(err)
	}

	return dest
}

func Test_StopMonitoring_Without_Order(t *testing.T) {
	assert := assert.New(t)

	p := NewPartner()
	obj := model.NewCode("CodeSpace", "STIF:StopPoint:Q:41178:")

	builder := NewLiteStopMonitoringUpdateEventBuilder(p, obj)
	sm := getsmlite(t, "testdata/stopmonitoring-lite-without-order.json")

	for _, delivery := range sm.Siri.ServiceDelivery.StopMonitoringDelivery {
		if delivery.Status == "false" {
			continue
		}
		builder.SetUpdateEvents(delivery.MonitoredStopVisit)
	}

	updateEvents := builder.UpdateEvents()
	emptyUpdateEvents := *NewCollectUpdateEvents()

	assert.Equal(emptyUpdateEvents, updateEvents, "update event should be empty")
}
