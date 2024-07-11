package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type SIRIGeneralMessageRequestCollectorFactory struct{}

type SIRIGeneralMessageRequestCollector struct {
	connector

	updateSubscriber UpdateSubscriber
}

func NewSIRIGeneralMessageRequestCollector(partner *Partner) *SIRIGeneralMessageRequestCollector {
	siriGeneralMessageRequestCollector := &SIRIGeneralMessageRequestCollector{}
	siriGeneralMessageRequestCollector.partner = partner
	manager := partner.Referential().CollectManager()
	siriGeneralMessageRequestCollector.updateSubscriber = manager.BroadcastUpdateEvent

	return siriGeneralMessageRequestCollector
}

func (connector *SIRIGeneralMessageRequestCollector) RequestAllSituationsUpdate() {}

func (connector *SIRIGeneralMessageRequestCollector) RequestSituationUpdate(kind, requestedId string) {
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	siriGeneralMessageRequest := &siri.SIRIGetGeneralMessageRequest{
		RequestorRef: connector.Partner().RequestorRef(),
	}
	siriGeneralMessageRequest.MessageIdentifier = connector.Partner().NewMessageIdentifier()
	siriGeneralMessageRequest.RequestTimestamp = connector.Clock().Now()

	// Check the request filter
	switch kind {
	case SITUATION_UPDATE_REQUEST_LINE:
		siriGeneralMessageRequest.LineRef = []string{requestedId}
	case SITUATION_UPDATE_REQUEST_STOP_AREA:
		siriGeneralMessageRequest.StopPointRef = []string{requestedId}
	}

	// Check the request version
	if connector.partner.GeneralMessageRequestVersion22() {
		siriGeneralMessageRequest.XsdInWsdl = true
	}

	connector.logSIRIGeneralMessageRequest(message, siriGeneralMessageRequest)

	xmlGeneralMessageResponse, err := connector.Partner().SIRIClient().SituationMonitoring(siriGeneralMessageRequest)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during GetGeneralMessage: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLGeneralMessageResponse(message, xmlGeneralMessageResponse)
	updateEvents := NewCollectUpdateEvents()

	builder := NewGeneralMessageUpdateEventBuilder(connector.partner)
	builder.SetGeneralMessageResponseUpdateEvents(updateEvents, xmlGeneralMessageResponse)

	// Log VehicleRefs
	message.Lines = GetModelReferenceSlice(builder.LineRefs)
	message.StopAreas = GetModelReferenceSlice(builder.MonitoringRefs)

	connector.broadcastSituationUpdateEvent(updateEvents)
}

func (connector *SIRIGeneralMessageRequestCollector) SetUpdateSubscriber(updateSubscriber UpdateSubscriber) {
	connector.updateSubscriber = updateSubscriber
}

func (connector *SIRIGeneralMessageRequestCollector) broadcastSituationUpdateEvent(event *CollectUpdateEvents) {
	if connector.updateSubscriber == nil {
		return

	}
	for _, e := range event.Situations {
		connector.updateSubscriber(e)
	}
}

func (connector *SIRIGeneralMessageRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.GENERAL_MESSAGE_REQUEST,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (factory *SIRIGeneralMessageRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func (factory *SIRIGeneralMessageRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIGeneralMessageRequestCollector(partner)
}

func (connector *SIRIGeneralMessageRequestCollector) logSIRIGeneralMessageRequest(message *audit.BigQueryMessage, request *siri.SIRIGetGeneralMessageRequest) {
	xml, err := request.BuildXML(connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}

	message.RequestIdentifier = request.MessageIdentifier
	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLGeneralMessageResponse(message *audit.BigQueryMessage, response *sxml.XMLGeneralMessageResponse) {
	message.ResponseIdentifier = response.ResponseMessageIdentifier()

	if !response.Status() {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
