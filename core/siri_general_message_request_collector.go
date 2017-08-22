package core

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type GeneralMessageRequestCollector interface {
	RequestSituationUpdate(request *SituationUpdateRequest)
}

type SIRIGeneralMessageRequestCollectorFactory struct{}

type SIRIGeneralMessageRequestCollector struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector

	situationUpdateSubscriber SituationUpdateSubscriber
}

func NewSIRIGeneralMessageRequestCollector(partner *Partner) *SIRIGeneralMessageRequestCollector {
	siriGeneralMessageRequestCollector := &SIRIGeneralMessageRequestCollector{}
	siriGeneralMessageRequestCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriGeneralMessageRequestCollector.situationUpdateSubscriber = manager.BroadcastSituationUpdateEvent

	return siriGeneralMessageRequestCollector
}

func (connector *SIRIGeneralMessageRequestCollector) RequestSituationUpdate(request *SituationUpdateRequest) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	startTime := connector.Clock().Now()

	siriGeneralMessageRequest := &siri.SIRIGeneralMessageRequest{
		MessageIdentifier: connector.SIRIPartner().IdentifierGenerator("message_identifier").NewMessageIdentifier(),
		RequestorRef:      connector.SIRIPartner().RequestorRef(),
		RequestTimestamp:  connector.Clock().Now(),
	}

	logSIRIGeneralMessageRequest(logStashEvent, siriGeneralMessageRequest)

	xmlGeneralMessageResponse, err := connector.SIRIPartner().SOAPClient().SituationMonitoring(siriGeneralMessageRequest)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["response"] = fmt.Sprintf("Error during GetGeneralMessage: %v", err)
		return
	}

	logXMLGeneralMessageResponse(logStashEvent, xmlGeneralMessageResponse)
	situationUpdateEvents := []*model.SituationUpdateEvent{}
	connector.setSituationUpdateEvents(&situationUpdateEvents, xmlGeneralMessageResponse)

	logSituationUpdateEvents(logStashEvent, situationUpdateEvents)

	connector.broadcastSituationUpdateEvent(situationUpdateEvents)
}

func (connector *SIRIGeneralMessageRequestCollector) setSituationUpdateEvents(situationEvents *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageResponse) {
	xmlGeneralMessages := xmlResponse.XMLGeneralMessages()
	if len(xmlGeneralMessages) == 0 {
		return
	}

	for _, xmlGeneralMessage := range xmlGeneralMessages {
		situationEvent := &model.SituationUpdateEvent{
			CreatedAt:         connector.Clock().Now(),
			RecordedAt:        xmlGeneralMessage.RecordedAtTime(),
			SituationObjectID: model.NewObjectID(connector.partner.Setting("remote_objectid_kind"), xmlGeneralMessage.InfoMessageIdentifier()),
			Version:           xmlGeneralMessage.InfoMessageVersion(),
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

func (connector *SIRIGeneralMessageRequestCollector) SetSituationUpdateSubscriber(situationUpdateSubscriber SituationUpdateSubscriber) {
	connector.situationUpdateSubscriber = situationUpdateSubscriber
}

func (connector *SIRIGeneralMessageRequestCollector) broadcastSituationUpdateEvent(event []*model.SituationUpdateEvent) {
	if connector.situationUpdateSubscriber != nil {
		connector.situationUpdateSubscriber(event)
	}
}

func (connector *SIRIGeneralMessageRequestCollector) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "GeneralMessageRequestCollector"
	return event
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
	logStashEvent["status"] = strconv.FormatBool(response.Status())
	if !response.Status() {
		logStashEvent["errorType"] = response.ErrorType()
		if response.ErrorType() == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber())
		}
		logStashEvent["errorText"] = response.ErrorText()
		logStashEvent["errorDescription"] = response.ErrorDescription()
	}
	logStashEvent["responseXML"] = response.RawXML()
}

func logSituationUpdateEvents(logStashEvent audit.LogStashEvent, situations []*model.SituationUpdateEvent) {
	var updateEvents []string
	for _, situationUpdateEvent := range situations {
		updateEvents = append(updateEvents, situationUpdateEvent.SituationObjectID.Value())
	}
	logStashEvent["SituationUpdateEvents"] = strings.Join(updateEvents, ", ")
}
