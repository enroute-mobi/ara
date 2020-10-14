package core

import (
	"fmt"
	"strconv"

	"bitbucket.org/enroute-mobi/edwig/audit"
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type GeneralMessageRequestCollector interface {
	RequestSituationUpdate(kind, requestedId string)
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

func (connector *SIRIGeneralMessageRequestCollector) RequestAllSituationsUpdate() {}

func (connector *SIRIGeneralMessageRequestCollector) RequestSituationUpdate(kind, requestedId string) {
	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	startTime := connector.Clock().Now()

	siriGeneralMessageRequest := &siri.SIRIGetGeneralMessageRequest{
		RequestorRef: connector.SIRIPartner().RequestorRef(),
	}
	siriGeneralMessageRequest.MessageIdentifier = connector.Partner().IdentifierGenerator(MESSAGE_IDENTIFIER).NewMessageIdentifier()
	siriGeneralMessageRequest.RequestTimestamp = connector.Clock().Now()

	// Check the request filter
	switch kind {
	case SITUATION_UPDATE_REQUEST_LINE:
		siriGeneralMessageRequest.LineRef = []string{requestedId}
		logStashEvent["lineRef"] = requestedId
	case SITUATION_UPDATE_REQUEST_STOP_AREA:
		siriGeneralMessageRequest.StopPointRef = []string{requestedId}
		logStashEvent["stopPointRef"] = requestedId
	}

	// Check the request version
	if b, _ := strconv.ParseBool(connector.partner.Setting(GENEREAL_MESSAGE_REQUEST_2)); b {
		siriGeneralMessageRequest.XsdInWsdl = true
	}

	logSIRIGeneralMessageRequest(logStashEvent, siriGeneralMessageRequest)

	xmlGeneralMessageResponse, err := connector.SIRIPartner().SOAPClient().SituationMonitoring(siriGeneralMessageRequest)
	logStashEvent["responseTime"] = connector.Clock().Since(startTime).String()
	if err != nil {
		logStashEvent["status"] = "false"
		logStashEvent["errorDescription"] = fmt.Sprintf("Error during GetGeneralMessage: %v", err)
		return
	}

	logXMLGeneralMessageResponse(logStashEvent, xmlGeneralMessageResponse)
	situationUpdateEvents := []*model.SituationUpdateEvent{}
	connector.setSituationUpdateEvents(&situationUpdateEvents, xmlGeneralMessageResponse)

	connector.broadcastSituationUpdateEvent(situationUpdateEvents)
}

func (connector *SIRIGeneralMessageRequestCollector) setSituationUpdateEvents(situationEvents *[]*model.SituationUpdateEvent, xmlResponse *siri.XMLGeneralMessageResponse) {
	builder := NewGeneralMessageUpdateEventBuilder(connector.partner)
	builder.SetGeneralMessageResponseUpdateEvents(situationEvents, xmlResponse)
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
	ok := apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	ok = ok && apiPartner.ValidatePresenceOfSetting(REMOTE_URL)
	ok = ok && apiPartner.ValidatePresenceOfSetting(REMOTE_CREDENTIAL)
	return ok
}

func (factory *SIRIGeneralMessageRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageRequestCollector(partner)
}

func logSIRIGeneralMessageRequest(logStashEvent audit.LogStashEvent, request *siri.SIRIGetGeneralMessageRequest) {
	logStashEvent["siriType"] = "GeneralMessageRequest"
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
