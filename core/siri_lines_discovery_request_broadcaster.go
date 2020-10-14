package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"bitbucket.org/enroute-mobi/edwig/audit"
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type LinesDiscoveryRequestBroadcaster interface {
	Lines(request *siri.XMLLinesDiscoveryRequest) (*siri.SIRILinesDiscoveryResponse, error)
}

type SIRILinesDiscoveryRequestBroadcaster struct {
	model.ClockConsumer

	siriConnector
}

type SIRILinesDiscoveryRequestBroadcasterFactory struct{}

func NewSIRILinesDiscoveryRequestBroadcaster(partner *Partner) *SIRILinesDiscoveryRequestBroadcaster {
	siriLinesDiscoveryRequestBroadcaster := &SIRILinesDiscoveryRequestBroadcaster{}
	siriLinesDiscoveryRequestBroadcaster.partner = partner
	return siriLinesDiscoveryRequestBroadcaster
}

func (connector *SIRILinesDiscoveryRequestBroadcaster) Lines(request *siri.XMLLinesDiscoveryRequest) (*siri.SIRILinesDiscoveryResponse, error) {
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

	objectIDKind := connector.partner.RemoteObjectIDKind(SIRI_LINES_DISCOVERY_REQUEST_BROADCASTER)
	for _, line := range tx.Model().Lines().FindAll() {
		if line.Name == "" {
			continue
		}

		objectID, ok := line.ObjectID(objectIDKind)
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

	logStashEvent["annotedLines"] = strings.Join(annotedLineArray, ", ")
	logSIRILineDiscoveryResponse(logStashEvent, response)

	return response, nil
}

func (connector *SIRILinesDiscoveryRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "LinesDiscoveryRequestBroadcaster"
	return event
}

func (factory *SIRILinesDiscoveryRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfLocalCredentials()
	return ok
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
