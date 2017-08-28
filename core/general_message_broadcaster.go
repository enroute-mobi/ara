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

func (gmb *GMBroadcaster) prepareSIRIGeneralMessageNotify() {
	gmb.connector.mutex.Lock()

	events := gmb.connector.toBroadcast
	gmb.connector.toBroadcast = make(map[SubscriptionId][]model.SituationId)

	gmb.connector.mutex.Unlock()

	tx := gmb.connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	referenceGenerator := gmb.connector.SIRIPartner().IdentifierGenerator("reference_identifier")

	for subId, situationIds := range events {
		logStashEvent := make(audit.LogStashEvent)
		defer audit.CurrentLogStash().WriteEvent(logStashEvent)

		sub, ok := gmb.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			continue
		}

		notify := siri.SIRINotifyGeneralMessage{
			Address:                   gmb.connector.Partner().Address(),
			ProducerRef:               gmb.connector.Partner().ProducerRef(),
			ResponseMessageIdentifier: gmb.connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
			SubscriberRef:             gmb.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier:    sub.ExternalId(),
			Status:                    true,
		}

		for _, situationId := range situationIds {
			situation, ok := tx.Model().Situations().Find(situationId)
			if !ok {
				logger.Log.Debugf("Could not find situation : %v in general message broadcaster", situationId)
				continue
			}

			if situation.Channel == "Commercial" || situation.ValidUntil.Before(gmb.connector.Clock().Now()) {
				continue
			}

			var infoMessageIdentifier string
			objectid, present := situation.ObjectID(gmb.connector.partner.RemoteObjectIDKind(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER))
			if present {
				infoMessageIdentifier = objectid.Value()
			} else {
				objectid, present = situation.ObjectID("_default")
				if !ok {
					continue
				}
				infoMessageIdentifier = referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "InfoMessage", Default: objectid.Value()})
			}

			siriGeneralMessage := &siri.SIRIGeneralMessage{
				ItemIdentifier:        referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "Item", Default: gmb.connector.NewUUID()}),
				InfoMessageIdentifier: infoMessageIdentifier,
				InfoChannelRef:        situation.Channel,
				InfoMessageVersion:    situation.Version,
				ValidUntilTime:        situation.ValidUntil,
				RecordedAtTime:        situation.RecordedAt,
				FormatRef:             "STIF-IDF",
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

			notify.GeneralMessages = append(notify.GeneralMessages, siriGeneralMessage)
		}
		gmb.connector.SIRIPartner().SOAPClient().NotifyGeneralMessage(&notify)
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
		if response.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		}
		logStashEvent["errorText"] = response.ErrorText
	}
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
