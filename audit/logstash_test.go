package audit

import (
	"io/ioutil"
	"net"
	"testing"
)

func Test_Logstash_WriteEvent(t *testing.T) {
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

func Test_TCPLogstash_WriteEvent(t *testing.T) {
	logStash := NewTCPLogStash("")
	event := make(LogStashEvent)
	event["foo"] = "bar"
	logStash.WriteEvent(event)
	event["foo"] = "bar2"
	logStash.WriteEvent(event)

	select {
	case <-logStash.events:
	default:
		t.Error("TCPLogStash should have one event after WriteEvent")
	}
	select {
	case <-logStash.events:
	default:
		t.Error("TCPLogStash should have two events after WriteEvent")
	}
}

func Test_TCPLogstash_Send(t *testing.T) {
	done := make(chan bool)
	go func() {
		logStash := NewTCPLogStash(":3000")
		logStash.Start()
		defer logStash.Stop()

		event := make(LogStashEvent)
		event["foo"] = "bar"
		logStash.WriteEvent(event)
		done <- true
	}()

	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	conn, err := l.Accept()
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	<-done

	buf, err := ioutil.ReadAll(conn)
	if err != nil {
		t.Fatal(err)
	}

	if msg := string(buf[:]); msg != "{\"foo\":\"bar\"}\n" {
		t.Fatalf("Unexpected message:\n Got: %s\n Expected: {\"foo\":\"bar\"}", msg)
	}
}
