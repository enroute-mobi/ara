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
		ResponseMessageIdentifier: connector.SIRIPartner().NewResponseMessageIdentifier(),
	}

	if response.ProducerRef == "" {
		response.ProducerRef = "Edwig"
	}

	response.Status = true
	response.ResponseTimestamp = connector.Clock().Now()

	objectidKind := connector.RemoteObjectIDKind()

	for _, situation := range tx.Model().Situations().FindAll() {
		xmlGeneralMessage := &siri.SIRIGeneralMessage{}
		objectid, _ := situation.ObjectID(objectidKind)
		xmlGeneralMessage.Messages = situation.Messages
		//xmlGeneralMessage.ItemIdentifier = situation.ItemIdentifier
		xmlGeneralMessage.InfoMessageIdentifier = objectid.Value()
		xmlGeneralMessage.InfoChannelRef = situation.Channel
		xmlGeneralMessage.InfoMessageVersion = situation.Version
		xmlGeneralMessage.ValidUntilTime = situation.ValidUntil
		xmlGeneralMessage.RecordedAtTime = situation.RecordedAt
		response.GeneralMessages = append(response.GeneralMessages, xmlGeneralMessage)
	}
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
