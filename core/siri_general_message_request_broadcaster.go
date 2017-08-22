package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type GeneralMessageRequestBroadcaster interface {
	Situations(request *siri.XMLGetGeneralMessage) (*siri.SIRIGeneralMessageResponse, error)
}

type SIRIGeneralMessageRequestBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer
	siriConnector
}

type SIRIGeneralMessageRequestBroadcasterFactory struct{}

func NewSIRIGeneralMessageRequestBroadcaster(partner *Partner) *SIRIGeneralMessageRequestBroadcaster {
	siriGeneralMessageRequestBroadcaster := &SIRIGeneralMessageRequestBroadcaster{}
	siriGeneralMessageRequestBroadcaster.partner = partner
	return siriGeneralMessageRequestBroadcaster
}

func (connector *SIRIGeneralMessageRequestBroadcaster) Situations(request *siri.XMLGetGeneralMessage) (*siri.SIRIGeneralMessageResponse, error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLGeneralMessageRequest(logStashEvent, &request.XMLGeneralMessageRequest)
	logStashEvent["requestorRef"] = request.RequestorRef()

	response := &siri.SIRIGeneralMessageResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
	}

	response.SIRIGeneralMessageDelivery = connector.getGeneralMessageDelivery(tx, logStashEvent, &request.XMLGeneralMessageRequest)

	logSIRIGeneralMessageDelivery(logStashEvent, response.SIRIGeneralMessageDelivery)
	logSIRIGeneralMessageResponse(logStashEvent, response)

	return response, nil
}

func (connector *SIRIGeneralMessageRequestBroadcaster) getGeneralMessageDelivery(tx *model.Transaction, logStashEvent audit.LogStashEvent, request *siri.XMLGeneralMessageRequest) siri.SIRIGeneralMessageDelivery {
	referenceGenerator := connector.SIRIPartner().IdentifierGenerator("reference_identifier")

	delivery := siri.SIRIGeneralMessageDelivery{
		RequestMessageRef: request.MessageIdentifier(),
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	// Prepare Id Array
	var messageArray []string

	for _, situation := range tx.Model().Situations().FindAll() {
		if situation.Channel == "Commercial" || situation.ValidUntil.Before(connector.Clock().Now()) {
			continue
		}

		var infoMessageIdentifier string
		objectid, present := situation.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_GENERAL_MESSAGE_REQUEST_BROADCASTER))
		if present {
			infoMessageIdentifier = objectid.Value()
		} else {
			objectid, present = situation.ObjectID("_default")
			if !present {
				continue
			}
			infoMessageIdentifier = referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "InfoMessage", Default: objectid.Value()})
		}

		messageArray = append(messageArray, infoMessageIdentifier)

		siriGeneralMessage := &siri.SIRIGeneralMessage{
			ItemIdentifier:        referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "Item", Default: connector.NewUUID()}),
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

		delivery.GeneralMessages = append(delivery.GeneralMessages, siriGeneralMessage)
	}

	logStashEvent["MessageIds"] = strings.Join(messageArray, ", ")

	return delivery
}

func (connector *SIRIGeneralMessageRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "GeneralMessageRequestBroadcaster"
	return event
}

func (factory *SIRIGeneralMessageRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRIGeneralMessageRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageRequestBroadcaster(partner)
}

func logXMLGeneralMessageRequest(logStashEvent audit.LogStashEvent, request *siri.XMLGeneralMessageRequest) {
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIGeneralMessageDelivery(logStashEvent audit.LogStashEvent, delivery siri.SIRIGeneralMessageDelivery) {
	logStashEvent["requestMessageRef"] = delivery.RequestMessageRef
	logStashEvent["responseTimestamp"] = delivery.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(delivery.Status)
	if !delivery.Status {
		logStashEvent["errorType"] = delivery.ErrorType
		if delivery.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber)
		}
		logStashEvent["errorText"] = delivery.ErrorText
	}
}

func logSIRIGeneralMessageResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIGeneralMessageResponse) {
	logStashEvent["address"] = response.Address
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
