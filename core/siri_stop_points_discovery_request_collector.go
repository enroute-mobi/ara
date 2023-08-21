package core

import (
	"fmt"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"
)

type StopPointsDiscoveryRequestCollector interface {
	state.Startable

	RequestStopPoints()
}

type SIRIStopPointsDiscoveryRequestCollector struct {
	connector

	stopAreaUpdateSubscriber UpdateSubscriber
}

type SIRIStopPointsDiscoveryRequestCollectorFactory struct{}

func (factory *SIRIStopPointsDiscoveryRequestCollectorFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopPointsDiscoveryRequestCollector(partner)
}

func (factory *SIRIStopPointsDiscoveryRequestCollectorFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfRemoteCredentials()
}

func NewSIRIStopPointsDiscoveryRequestCollector(partner *Partner) *SIRIStopPointsDiscoveryRequestCollector {
	connector := &SIRIStopPointsDiscoveryRequestCollector{}
	connector.partner = partner
	manager := partner.Referential().CollectManager()
	connector.stopAreaUpdateSubscriber = manager.BroadcastUpdateEvent

	return connector
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) Start() {
	connector.remoteObjectidKind = connector.partner.RemoteObjectIDKind()
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) SetSubscriber(subscriber UpdateSubscriber) {
	connector.stopAreaUpdateSubscriber = subscriber
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) broadcastUpdateEvent(event model.UpdateEvent) {
	if connector.stopAreaUpdateSubscriber != nil {
		connector.stopAreaUpdateSubscriber(event)
	}
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) RequestStopPoints() {
	message := connector.newBQEvent()
	defer audit.CurrentBigQuery(string(connector.Partner().Referential().Slug())).WriteEvent(message)

	startTime := connector.Clock().Now()

	request := &siri.SIRIStopPointsDiscoveryRequest{
		MessageIdentifier: connector.Partner().NewMessageIdentifier(),
		RequestorRef:      connector.Partner().RequestorRef(),
		RequestTimestamp:  startTime,
	}

	connector.logSIRIStopPointsDiscoveryRequest(message, request)

	response, err := connector.Partner().SIRIClient().StopDiscovery(request)
	message.ProcessingTime = connector.Clock().Since(startTime).Seconds()
	if err != nil {
		e := fmt.Sprintf("Error during StopDiscovery: %v", err)

		message.Status = "Error"
		message.ErrorDetails = e
		return
	}

	logXMLStopPointsDiscoveryResponse(message, response)

	if !response.Status() {
		return
	}

	stopPointRefs := []string{}
	idKind := connector.remoteObjectidKind
	partner := string(connector.Partner().Slug())

	for _, annotatedStopPoint := range response.AnnotatedStopPointRefs() {
		stopPointRefs = append(stopPointRefs, annotatedStopPoint.StopPointRef())
		event := model.NewStopAreaUpdateEvent()

		event.Origin = partner
		event.ObjectId = model.NewObjectID(idKind, annotatedStopPoint.StopPointRef())
		event.Name = annotatedStopPoint.StopName()
		event.CollectedAlways = true

		connector.broadcastUpdateEvent(event)
	}

	connector.partner.RegisterDiscoveredStopAreas(stopPointRefs)
	message.StopAreas = stopPointRefs
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "StopPointsDiscoveryRequest",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(connector.partner.Slug()),
		Status:    "OK",
	}
}

func (connector *SIRIStopPointsDiscoveryRequestCollector) logSIRIStopPointsDiscoveryRequest(message *audit.BigQueryMessage, request *siri.SIRIStopPointsDiscoveryRequest) {
	message.RequestIdentifier = request.MessageIdentifier

	xml, err := request.BuildXML(connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}

	message.RequestRawMessage = xml
	message.RequestSize = int64(len(xml))
}

func logXMLStopPointsDiscoveryResponse(message *audit.BigQueryMessage, response *sxml.XMLStopPointsDiscoveryResponse) {
	if !response.Status() {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}

	message.ResponseRawMessage = response.RawXML()
	message.ResponseSize = int64(len(message.ResponseRawMessage))
}
