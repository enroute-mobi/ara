package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type LinesDiscoveryRequestBroadcaster interface {
	Lines(*siri.XMLLinesDiscoveryRequest, *audit.BigQueryMessage) (*siri.SIRILinesDiscoveryResponse, error)
}

type SIRILinesDiscoveryRequestBroadcaster struct {
	clock.ClockConsumer

	connector
}

type SIRILinesDiscoveryRequestBroadcasterFactory struct{}

func NewSIRILinesDiscoveryRequestBroadcaster(partner *Partner) *SIRILinesDiscoveryRequestBroadcaster {
	connector := &SIRILinesDiscoveryRequestBroadcaster{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(SIRI_LINES_DISCOVERY_REQUEST_BROADCASTER)
	connector.partner = partner
	return connector
}

func (connector *SIRILinesDiscoveryRequestBroadcaster) Lines(request *siri.XMLLinesDiscoveryRequest, message *audit.BigQueryMessage) (*siri.SIRILinesDiscoveryResponse, error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLLineDiscoveryRequest(logStashEvent, request)

	response := &siri.SIRILinesDiscoveryResponse{
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	var annotedLineArray []string

	for _, line := range tx.Model().Lines().FindAll() {
		if line.Name == "" {
			continue
		}

		objectID, ok := line.ObjectID(connector.remoteObjectidKind)
		if !ok {
			continue
		}

		annotedLine := &siri.SIRIAnnotatedLine{
			LineName:  line.Name,
			LineRef:   objectID.Value(),
			Monitored: true,
		}
		annotedLineArray = append(annotedLineArray, annotedLine.LineRef)
		response.AnnotatedLines = append(response.AnnotatedLines, annotedLine)
	}

	sort.Sort(siri.SIRIAnnotatedLineByLineRef(response.AnnotatedLines))

	message.RequestIdentifier = request.MessageIdentifier()
	message.Lines = annotedLineArray

	logStashEvent["annotedLines"] = strings.Join(annotedLineArray, ", ")
	logSIRILineDiscoveryResponse(logStashEvent, response)

	return response, nil
}

func (connector *SIRILinesDiscoveryRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "LinesDiscoveryRequestBroadcaster"
	return event
}

func (factory *SIRILinesDiscoveryRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRILinesDiscoveryRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRILinesDiscoveryRequestBroadcaster(partner)
}

func logXMLLineDiscoveryRequest(logStashEvent audit.LogStashEvent, request *siri.XMLLinesDiscoveryRequest) {
	logStashEvent["siriType"] = "LinesDiscoveryResponse"
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRILineDiscoveryResponse(logStashEvent audit.LogStashEvent, response *siri.SIRILinesDiscoveryResponse) {
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
