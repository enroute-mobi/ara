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

type SIRIStopPointsDiscoveryRequestBroadcaster struct {
	clock.ClockConsumer

	siriConnector
}

type SIRIStopPointsDiscoveryRequestBroadcasterFactory struct{}

func NewSIRIStopDiscoveryRequestBroadcaster(partner *Partner) *SIRIStopPointsDiscoveryRequestBroadcaster {
	siriStopDiscoveryRequestBroadcaster := &SIRIStopPointsDiscoveryRequestBroadcaster{}
	siriStopDiscoveryRequestBroadcaster.partner = partner
	return siriStopDiscoveryRequestBroadcaster
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) StopAreas(request *siri.XMLStopPointsDiscoveryRequest, message *audit.BigQueryMessage) (*siri.SIRIStopPointsDiscoveryResponse, error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLStopPointDiscoveryRequest(logStashEvent, request)

	response := &siri.SIRIStopPointsDiscoveryResponse{
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	annotedStopPointMap := make(map[string]struct{})

	objectIDKind := connector.partner.RemoteObjectIDKind(SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER)
	for _, stopArea := range tx.Model().StopAreas().FindAll() {
		if stopArea.Name == "" || !stopArea.CollectedAlways {
			continue
		}

		objectID, ok := stopArea.ReferentOrSelfObjectId(objectIDKind)
		if !ok || objectID.Value() == "" {
			continue
		}
		_, ok = annotedStopPointMap[objectID.Value()]
		if ok {
			continue
		}
		annotedStopPointMap[objectID.Value()] = struct{}{}

		annotedStopPoint := &siri.SIRIAnnotatedStopPoint{
			StopName:     stopArea.Name,
			StopPointRef: objectID.Value(),
			Monitored:    true,
			TimingPoint:  true,
		}
		for _, line := range stopArea.Lines() {
			if line.Origin() == string(connector.partner.Slug()) {
				continue
			}
			objectid, ok := line.ObjectID(objectIDKind)
			if !ok {
				continue
			}
			annotedStopPoint.Lines = append(annotedStopPoint.Lines, objectid.Value())
		}
		if len(annotedStopPoint.Lines) == 0 && connector.ignoreStopWithoutLine() {
			continue
		}
		response.AnnotatedStopPoints = append(response.AnnotatedStopPoints, annotedStopPoint)
	}

	sort.Sort(siri.SIRIAnnotatedStopPointByStopPointRef(response.AnnotatedStopPoints))

	message.RequestIdentifier = request.MessageIdentifier()

	logAnnotatedStopPoints(annotedStopPointMap, logStashEvent, message)
	logSIRIStopPointDiscoveryResponse(logStashEvent, response)

	return response, nil
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) ignoreStopWithoutLine() bool {
	return connector.partner.Setting(IGNORE_STOP_WITHOUT_LINE) != "false"
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopPointsDiscoveryRequestBroadcaster"
	return event
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	ok = ok && apiPartner.ValidatePresenceOfLocalCredentials()
	return ok
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopDiscoveryRequestBroadcaster(partner)
}

func logAnnotatedStopPoints(annotedStopPointMap map[string]struct{}, logStashEvent audit.LogStashEvent, message *audit.BigQueryMessage) {
	keys := make([]string, len(annotedStopPointMap))
	i := 0
	for key := range annotedStopPointMap {
		keys[i] = key
		i++
	}

	logStashEvent["annotedStopPoints"] = strings.Join(keys, ", ")
	message.StopAreas = keys
}

func logXMLStopPointDiscoveryRequest(logStashEvent audit.LogStashEvent, request *siri.XMLStopPointsDiscoveryRequest) {
	logStashEvent["siriType"] = "StopPointsDiscoveryResponse"
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIStopPointDiscoveryResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIStopPointsDiscoveryResponse) {
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
