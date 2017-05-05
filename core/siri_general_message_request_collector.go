package core

import (
	"fmt"
	"strconv"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type GeneralMessageRequestCollector interface {
	RequestSituationUpdate(request *SituationUpdateRequest) ([]*model.SituationUpdateEvent, error)
}

type SIRIGeneralMessageRequestCollectorFactory struct{}

type SIRIGeneralMessageRequestCollector struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector
}

func NewSIRIGeneralMessageRequestCollector(partner *Partner) *SIRIGeneralMessageRequestCollector {
	siriGeneralMessageRequestCollector := &SIRIGeneralMessageRequestCollector{}
	siriGeneralMessageRequestCollector.partner = partner
	return siriGeneralMessageRequestCollector
}

func (connector *SIRIGeneralMessageRequestCollector) RequestSituationUpdate(request *SituationUpdateRequest) ([]*model.SituationUpdateEvent, error) {
	logStashEvent := make(audit.LogStashEvent)
	startTime := connector.Clock().Now()

	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	siriGeneralMessageRequest := &siri.SIRIGeneralMessageRequest{
		MessageIdentifier: connector.SIRIPartner().NewMessageIdentifier(),
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  connector.Clock().Now(),
	}

	logSIRIGeneralMessageRequest(logStashEvent, siriGeneralMessageRequest)

	xmlGeneralMessageResponse, err := connector.SIRIPartner().SOAPClient().SituationMonitoring(siriGeneralMessageRequest)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["response"] = fmt.Sprintf("Error during CheckStatus: %v", err)
		return nil, err
	}

	logXMLGeneralMessageResponse(logStashEvent, xmlGeneralMessageResponse)
	situationUpdateEvents := []*model.SituationUpdateEvent{}
	connector.setSituationUpdateEvents(&situationUpdateEvents, xmlGeneralMessageResponse)

	return situationUpdateEvents, nil
}

func (connector *SIRIGeneralMessageRequestCollector) setSituationUpdateEvents(situationEvents *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageResponse) {
	xmlGeneralMessages := xmlResponse.XMLGeneralMessage()
	if len(xmlGeneralMessages) == 0 {
		return
	}

	for _, xmlGeneralMessage := range xmlGeneralMessages {
		situationEvent := &model.SituationUpdateEvent{
			CreatedAt:         connector.Clock().Now(),
			RecordedAt:        xmlGeneralMessage.RecordedAtTime(),
			SituationObjectID: model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlGeneralMessage.InfoMessageIdentifier()),
			Version:           int64(xmlGeneralMessage.InfoMessageVersion()),
			ProducerRef:       xmlResponse.ProducerRef(),
		}
		situationEvent.SetId(model.SituationUpdateRequestId(connector.NewUUID()))
		if xmlGeneralMessage.Content() != nil {
			content := xmlGeneralMessage.Content().(siri.IDFGeneralMessageStructure)
			for _, xmlMessage := range content.Messages() {
				message := &model.Message{
					Content:             xmlMessage.MessageText(),
					Type:                xmlMessage.MessageType(),
					NumberOfLines:       xmlMessage.NumberOfLines(),
					NumberOfCharPerLine: xmlMessage.NumberOfCharPerLine(),
				}
				situationEvent.SituationAttributes.Messages = append(situationEvent.SituationAttributes.Messages, message)
			}
		}
		situationEvent.SituationAttributes.Format = xmlGeneralMessage.FormatRef()
		situationEvent.SituationAttributes.Channel = xmlGeneralMessage.InfoChannelRef()
		situationEvent.SituationAttributes.ValidUntil = xmlGeneralMessage.ValidUntilTime()
		*situationEvents = append(*situationEvents, situationEvent)
	}
}

func (factory *SIRIGeneralMessageRequestCollectorFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_url")
	ok = ok && apiPartner.ValidatePresenceOfSetting("remote_credential")
	return ok
}

func (factory *SIRIGeneralMessageRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageRequestCollector(partner)
}

func logSIRIGeneralMessageRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIGeneralMessageRequest) {
	logStashEvent["Connector"] = "GeneralMessageRequestCollector"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier
	logStashEvent["requestorRef"] = request.RequestorRef
	logStashEvent["requestTimestamp"] = request.RequestTimestamp.String()
	xml, err := request.BuildXML()
	if err != nil {
		logStashEvent["requestXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["requestXML"] = xml
}

func logXMLGeneralMessageResponse(logStashEvent audit.LogStashEvent, response *siri.XMLGeneralMessageResponse) {
	logStashEvent["address"] = response.Address()
	logStashEvent["producerRef"] = response.ProducerRef()
	logStashEvent["requestMessageRef"] = response.RequestMessageRef()
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier()
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp().String()
	logStashEvent["responseXML"] = response.RawXML()
	logStashEvent["status"] = strconv.FormatBool(response.Status())
}
