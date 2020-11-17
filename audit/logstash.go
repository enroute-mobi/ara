package audit

import (
	"encoding/json"
	"net"
	"time"

	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/state"
)

type LogStashEvent map[string]string

type LogStash interface {
	state.Startable
	state.Stopable

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
		if err == nil {
			return
		}
		logStash.connection.Close()
		logStash.connectLogstash()
	}
}

func (logStash *TCPLogStash) connectLogstash() {
	for {
		clock.DefaultClock().Sleep(5 * time.Second)
		var err error
		logStash.connection, err = net.DialTimeout("tcp", logStash.address, 5*time.Second)
		if err == nil {
			return
		}
	}
}
