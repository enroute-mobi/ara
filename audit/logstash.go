package audit

import (
	"encoding/json"
	"net"

	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type LogStashEvent map[string]string

type LogStash interface {
	model.Startable
	model.Stopable

	WriteEvent(event LogStashEvent) error
}

type NullLogStash struct {
}

func (logStash *NullLogStash) WriteEvent(event LogStashEvent) error {
	return nil
}

func (logStash *NullLogStash) Start() {}
func (logStash *NullLogStash) Stop()  {}

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
	address    string
	connection net.Conn
	events     chan LogStashEvent
	stop       chan struct{}
}

func NewTCPLogStash(address string) *TCPLogStash {
	return &TCPLogStash{
		address: address,
		events:  make(chan LogStashEvent, 5),
	}
}

type FakeLogStash struct {
	events []LogStashEvent
}

func NewFakeLogStash() *FakeLogStash {
	return &FakeLogStash{}
}

func (logStash *FakeLogStash) Start() {}
func (logStash *FakeLogStash) Stop()  {}

func (logStash *FakeLogStash) WriteEvent(event LogStashEvent) error {
	logStash.events = append(logStash.events, event)
	return nil
}

func (logStash *FakeLogStash) Events() []LogStashEvent {
	return logStash.events
}

func (logStash *TCPLogStash) WriteEvent(event LogStashEvent) error {
	select {
	case logStash.events <- event:
	default:
		logger.Log.Debugf("LogStash queue is full")
	}
	return nil
}

func (logStash *TCPLogStash) Start() {
	logStash.stop = make(chan struct{})
	go logStash.run()
}

func (logStash *TCPLogStash) Stop() {
	if logStash.stop != nil {
		close(logStash.stop)
	}
}

func (logStash *TCPLogStash) run() {
	logStash.connectLogstash()

	for {
		select {
		case <-logStash.stop:
			logStash.connection.Close()
			return
		case event := <-logStash.events:
			jsonBytes, err := json.Marshal(event)
			if err != nil {
				logger.Log.Debugf("Can't marshal LogStash event: %v", event)
				continue
			}
			jsonBytes = append(jsonBytes, '\n')

			logStash.send(jsonBytes)
		}
	}
}

func (logStash *TCPLogStash) send(jsonBytes []byte) {
	for {
		_, err := logStash.connection.Write(jsonBytes)
		if err != nil {
			logStash.connection.Close()
			logStash.connectLogstash()
		}
		return
	}
}

func (logStash *TCPLogStash) connectLogstash() {
	for {
		var err error
		logStash.connection, err = net.Dial("tcp", logStash.address)
		if err == nil {
			return
		}
		model.DefaultClock().Sleep(5)
	}
}
