package core

import (
	"fmt"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type GeneralMessageRequestBroadcaster interface {
	Situations(request *siri.XMLGeneralMessageRequest) (*siri.SIRIGeneralMessageResponse, error)
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

func (connector *SIRIGeneralMessageRequestBroadcaster) Situations(request *siri.XMLGeneralMessageRequest) (*siri.SIRIGeneralMessageResponse, error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)
	logXMLGeneralMessageRequest(logStashEvent, request)

	response := &siri.SIRIGeneralMessageResponse{
		Address:                   connector.Partner().Setting("local_url"),
		ProducerRef:               connector.Partner().Setting("remote_credential"),
		ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
	}

	if response.ProducerRef == "" {
		response.ProducerRef = "Edwig"
	}

	response.Status = true
	response.ResponseTimestamp = connector.Clock().Now()
	response.RequestMessageRef = request.MessageIdentifier()

	for _, situation := range tx.Model().Situations().FindAll() {
		if situation.Channel == "Commercial" || situation.ValidUntil.Before(connector.Clock().Now()) {
			continue
		}
		siriGeneralMessage := &siri.SIRIGeneralMessage{}
		objectid, present := situation.ObjectID(connector.RemoteObjectIDKind())
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

		response.GeneralMessages = append(response.GeneralMessages, siriGeneralMessage)
	}

	logSIRIGeneralMessageResponse(logStashEvent, response)
	return response, nil
}

func (connector *SIRIGeneralMessageRequestBroadcaster) RemoteObjectIDKind() string {
	if connector.partner.Setting("siri-general-message-request-broadcaster.remote_objectid_kind") != "" {
		return connector.partner.Setting("siri-general-message-request-broadcaster.remote_objectid_kind")
	}
	return connector.partner.Setting("remote_objectid_kind")
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
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIGeneralMessageResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIGeneralMessageResponse) {
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
