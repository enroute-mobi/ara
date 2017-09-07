package core

import (
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
)

type SIRIEstimatedTimeTableBroadcaster interface {
	model.Stopable
	model.Startable
}

type ETTBroadcaster struct {
	model.ClockConsumer

	connector *SIRIEstimatedTimeTableSubscriptionBroadcaster
}

type EstimatedTimeTableBroadcaster struct {
	ETTBroadcaster

	stop chan struct{}
}

type FakeEstimatedTimeTableBroadcaster struct {
	ETTBroadcaster

	model.ClockConsumer
}

func NewFakeEstimatedTimeTableBroadcaster(connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) SIRIEstimatedTimeTableBroadcaster {
	broadcaster := &FakeEstimatedTimeTableBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeEstimatedTimeTableBroadcaster) Start() {
	broadcaster.prepareSIRIEstimatedTimeTable()
}

func (broadcaster *FakeEstimatedTimeTableBroadcaster) Stop() {}

func NewSIRIEstimatedTimeTableBroadcaster(connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) SIRIEstimatedTimeTableBroadcaster {
	broadcaster := &EstimatedTimeTableBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (ett *EstimatedTimeTableBroadcaster) Start() {
	if ett.stop != nil {
		return
	}

	logger.Log.Debugf("Start EstimatedTimeTableBroadcaster")

	ett.stop = make(chan struct{})
	go ett.run()
}

func (ett *EstimatedTimeTableBroadcaster) run() {
	c := ett.Clock().After(5 * time.Second)

	for {
		select {
		case <-ett.stop:
			logger.Log.Debugf("estimated time table broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRIEstimatedTimeTableBroadcaster visit")

			ett.prepareSIRIEstimatedTimeTable()

			c = ett.Clock().After(5 * time.Second)
		}
	}
}

func (ett *EstimatedTimeTableBroadcaster) Stop() {
	if ett.stop != nil {
		close(ett.stop)
	}
}

func (ett *ETTBroadcaster) prepareSIRIEstimatedTimeTable() {
	ett.connector.mutex.Lock()

	events := ett.connector.toBroadcast
	ett.connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	ett.connector.mutex.Unlock()

	logStashEvent := ett.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	tx := ett.connector.Partner().Referential().NewTransaction()
	defer tx.Close()
}

func (smb *ETTBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := smb.connector.partner.NewLogStashEvent()
	event["connector"] = "EstimatedTimeTableSubscriptionBroadcaster"
	return event
}
