package core

import (
	"fmt"
	"strconv"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIGeneralMessageBroadcaster interface {
	model.Stopable
	model.Startable
}

type GMBroadcaster struct {
	model.ClockConsumer

	connector *SIRIGeneralMessageSubscriptionBroadcaster
}

type GeneralMessageBroadcaster struct {
	GMBroadcaster

	stop chan struct{}
}

type FakeGeneralMessageBroadcaster struct {
	GMBroadcaster

	model.ClockConsumer
}

func NewFakeGeneralMessageBroadcaster(connector *SIRIGeneralMessageSubscriptionBroadcaster) SIRIGeneralMessageBroadcaster {
	broadcaster := &FakeGeneralMessageBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeGeneralMessageBroadcaster) Start() {
	broadcaster.prepareSIRIGeneralMessageNotify()
}

func (broadcaster *FakeGeneralMessageBroadcaster) Stop() {}

func NewSIRIGeneralMessageBroadcaster(connector *SIRIGeneralMessageSubscriptionBroadcaster) SIRIGeneralMessageBroadcaster {
	broadcaster := &GeneralMessageBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (gmb *GeneralMessageBroadcaster) Start() {
	if gmb.stop != nil {
		return
	}

	logger.Log.Debugf("Start GeneralMessageBroadcaster")

	gmb.stop = make(chan struct{})
	go gmb.run()
}

func (gmb *GeneralMessageBroadcaster) run() {
	c := gmb.Clock().After(5 * time.Second)

	for {
		select {
		case <-gmb.stop:
			logger.Log.Debugf("general message broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRIGeneralMessageBroadcaster visit")

			gmb.prepareSIRIGeneralMessageNotify()

			c = gmb.Clock().After(5 * time.Second)
		}
	}
}

func (gmb *GeneralMessageBroadcaster) Stop() {
	if gmb.stop != nil {
		close(gmb.stop)
	}
}

func (gmb *GMBroadcaster) RemoteObjectIDKind() string {
	if gmb.connector.partner.Setting("siri-general-message-request-broadcaster.remote_objectid_kind") != "" {
		return gmb.connector.partner.Setting("siri-general-message-request-broadcaster.remote_objectid_kind")
	}
	return gmb.connector.partner.Setting("remote_objectid_kind")
}

func (gmb *GMBroadcaster) prepareSIRIGeneralMessageNotify() {
	connector := gmb.connector

	gmb.connector.mutex.Lock()

	events := gmb.connector.toBroadcast
	gmb.connector.toBroadcast = make(map[SubscriptionId][]model.SituationId)

	gmb.connector.mutex.Unlock()

	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	for subId, situationIds := range events {
		logStashEvent := make(audit.LogStashEvent)
		defer audit.CurrentLogStash().WriteEvent(logStashEvent)

		notify := siri.SIRINotifyGeneralMessage{
			Address:                   connector.Partner().Setting("local_url"),
			ProducerRef:               connector.Partner().Setting("remote_credential"),
			ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
			SubscriberRef:             gmb.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier:    fmt.Sprintf("Edwig:Subscription::%v:LOC", subId),
			Status:                    true,
		}

		if notify.ProducerRef == "" {
			notify.ProducerRef = "Edwig"
		}

		for _, situationId := range situationIds {
			situation, ok := connector.Partner().Model().Situations().Find(situationId)
			if !ok {
				continue
			}

			// if situation.Channel == "Commercial" || situation.ValidUntil.Before(connector.Clock().Now()) {
			// 	continue
			// }

			siriGeneralMessage := &siri.SIRIGeneralMessage{}
			objectid, present := situation.ObjectID(gmb.RemoteObjectIDKind())
			if !present {
				objectid, _ = situation.ObjectID("_default")
			}

			for _, message := range situation.Messages {
				siriMessage := &siri.SIRIMessage{
					Content:             message.Content,
					Type:                message.Type,
					NumberOfLines:       message.NumberOfLines,
					NumberOfCharPerLine: message.NumberOfCharPerLine,
				}
				siriGeneralMessage.Messages = append(siriGeneralMessage.Messages, siriMessage)
			}

			siriGeneralMessage.ItemIdentifier = fmt.Sprintf("RATPDev:Item::%s:LOC", connector.NewUUID())
			siriGeneralMessage.InfoMessageIdentifier = fmt.Sprintf("Edwig:InfoMessage::%s:LOC", objectid.Value())
			siriGeneralMessage.InfoChannelRef = situation.Channel
			siriGeneralMessage.InfoMessageVersion = situation.Version
			siriGeneralMessage.ValidUntilTime = situation.ValidUntil
			siriGeneralMessage.RecordedAtTime = situation.RecordedAt
			siriGeneralMessage.FormatRef = "STIF-IDF"

			notify.GeneralMessages = append(notify.GeneralMessages, siriGeneralMessage)
		}

		//fmt.Println(notify.BuildXML())
		connector.SIRIPartner().SOAPClient().NotifyGeneralMessage(&notify)
		logSIRIGeneralMessageNotify(logStashEvent, &notify)
	}
}

func logSIRIGeneralMessageNotify(logStashEvent audit.LogStashEvent, response *siri.SIRINotifyGeneralMessage) {
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	if !response.Status {
		logStashEvent["errorType"] = response.ErrorType
		logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		logStashEvent["errorText"] = response.ErrorText
	}
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
