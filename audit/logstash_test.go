package audit

import (
	"testing"

	"github.com/af83/edwig/logger"
)

func Test_Logstash_WriteEvent(t *testing.T) {
	logger.Log.Debug = true

	logstash := NewFakeLogStash()

	logstashDatas := make(map[string]string)
	logstashDatas["messageIdentifier"] = "identifier"
	logstashDatas["requestorRef"] = "requestor"
	logstashDatas["requestTimestamp"] = "timestamp"

	err := logstash.WriteEvent(logstashDatas)
	if err != nil {
		t.Errorf("Error when sending to logstash: %v", err)
	}

	if len(logstash.Events()) != 1 {
		t.Errorf("Logstash should have one event, got %d", len(logstash.Events()))
	}
}
