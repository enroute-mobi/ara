package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type LinesDiscoveryRequestCollector interface {
	state.Startable

	RequestLines()
}

type SIRILinesDiscoveryRequestCollector struct {
	connector

	lineUpdateSubscriber UpdateSubscriber
}

type SIRILinesDiscoveryRequestCollectorFactory struct{}

func (factory *SIRILinesDiscoveryRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRILinesDiscoveryRequestCollector(partner)
}

func (factory *SIRILinesDiscoveryRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func NewSIRILinesDiscoveryRequestCollector(partner *Partner) *SIRILinesDiscoveryRequestCollector {
	connector := &SIRILinesDiscoveryRequestCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.lineUpdateSubscriber = manager.BroadcastUpdateEvent

	return connector
}

func (connector *SIRILinesDiscoveryRequestCollector) Start() {
	connector.remoteObjectidKind = connector.partner.RemoteObjectIDKind()
}

func (connector *SIRILinesDiscoveryRequestCollector) SetSubscriber(subscriber UpdateSubscriber) {
	connector.lineUpdateSubscriber = subscriber
}

func (connector *SIRILinesDiscoveryRequestCollector) broadcastUpdateEvent(event model.UpdateEvent) {
	if connector.lineUpdateSubscriber != nil {
		connector.lineUpdateSubscriber(event)
	}
}

func (connector *SIRILinesDiscoveryRequestCollector) RequestLines() {
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	request := &siri.SIRILinesDiscoveryRequest{
		MessageIdentifier: connector.Partner().NewMessageIdentifier(),
		RequestorRef:      connector.Partner().RequestorRef(),
		RequestTimestamp:  startTime,
	}

	connector.logSIRILinesDiscoveryRequest(message, request)

	response, err := connector.Partner().SIRIClient().LineDiscovery(request)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during LinesDiscovery: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLLinesDiscoveryResponse(message, response)

	if !response.Status() {
		return
	}

	lineRefs := []string{}
	partner := string(connector.Partner().Slug())

	for _, annotatedLine := range response.AnnotatedLineRefs() {
		lineRefs = append(lineRefs, annotatedLine.LineRef())
		event := model.NewLineUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(connector.remoteObjectidKind, annotatedLine.LineRef())
		event.Name = annotatedLine.LineName()

		connector.broadcastUpdateEvent(event)
	}

	connector.partner.RegisterDiscoveredLines(lineRefs)
	message.Lines = lineRefs
}

func (connector *SIRILinesDiscoveryRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      audit.LINES_DISCOVERY_REQUEST,
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (connector *SIRILinesDiscoveryRequestCollector) logSIRILinesDiscoveryRequest(message *audit.BigQueryMessage, request *siri.SIRILinesDiscoveryRequest) {
	message.RequestIdentifier = request.MessageIdentifier

	xml, err := request.BuildXML(connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}

	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLLinesDiscoveryResponse(message *audit.BigQueryMessage, response *sxml.XMLLinesDiscoveryResponse) {
	if !response.Status() {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
