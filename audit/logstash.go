package audit

import (
	"encoding/json"
	"net"

	"github.com/af83/edwig/logger"
)

type LogStashEvent map[string]string

type LogStash interface {
	WriteEvent(event LogStashEvent) error
}

type NullLogStash struct {
}

func (logStash *NullLogStash) WriteEvent(event LogStashEvent) error {
	return nil
}

func NewNullLogStash() LogStash {
	return &NullLogStash{}
}

var currentLogStash LogStash = NewNullLogStash()

func CurrentLogStash() LogStash {
	return currentLogStash
}

func SetCurrentLogstash(logStash LogStash) {
	currentLogStash = logStash
}

type TCPLogStash struct {
	address string
}

func NewTCPLogStash(address string) *TCPLogStash {
	return &TCPLogStash{address: address}
}

type FakeLogStash struct {
	events []LogStashEvent
}

func NewFakeLogStash() *FakeLogStash {
	return &FakeLogStash{}
}

func (logStash *FakeLogStash) WriteEvent(event LogStashEvent) error {
	logStash.events = append(logStash.events, event)
	return nil
}

func (logStash *FakeLogStash) Events() []LogStashEvent {
	return logStash.events
}

func (logStash *TCPLogStash) WriteEvent(datas LogStashEvent) error {
	conn, err := net.Dial("tcp", logStash.address)
	if err != nil {
		return err
	}
	defer conn.Close()
	jsonBytes, err := json.Marshal(datas)
	if err != nil {
		return err
	}
	jsonBytes = append(jsonBytes, '\n')

	logger.Log.Debugf("Sending data to logstash: %s", string(jsonBytes))

	_, err = conn.Write(jsonBytes)
	if err != nil {
		return err
	}

	return nil
}
