package core

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopPointsDiscoveryRequestBroadcaster interface {
	stopAreas(request *siri.XMLStopPointsDiscoveryRequest) (*siri.SIRIStopPointsDiscoveryResponse, error)
}

type SIRIStopPointsDiscoveryRequestBroadcaster struct {
	model.ClockConsumer

	siriConnector
}

type SIRIStopPointsDiscoveryRequestBroadcasterFactory struct{}

func NewSIRIStopDiscoveryRequestBroadcaster(partner *Partner) *SIRIStopPointsDiscoveryRequestBroadcaster {
	siriStopDiscoveryRequestBroadcaster := &SIRIStopPointsDiscoveryRequestBroadcaster{}
	siriStopDiscoveryRequestBroadcaster.partner = partner
	return siriStopDiscoveryRequestBroadcaster
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) StopAreas(request *siri.XMLStopPointsDiscoveryRequest) (*siri.SIRIStopPointsDiscoveryResponse, error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLStopPointDiscoveryRequest(logStashEvent, request)

	response := &siri.SIRIStopPointsDiscoveryResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		RequestMessageRef:         request.MessageIdentifier(),
		ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	var annotedStopPointArray []string

	objectIDKind := connector.partner.RemoteObjectIDKind(SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER)
	for _, stopArea := range tx.Model().StopAreas().FindAll() {
		if stopArea.Name == "" || stopArea.CollectedAlways == false {
			continue
		}

		objectID, ok := stopArea.ObjectID(objectIDKind)
		if !ok {
			continue
		}

		annotedStopPoint := &siri.SIRIAnnotatedStopPoint{
			StopName:     stopArea.Name,
			StopPointRef: objectID.Value(),
			Monitored:    true,
			TimingPoint:  true,
		}
		for _, line := range stopArea.Lines() {
			objectid, ok := line.ObjectID(objectIDKind)
			if ok {
				annotedStopPoint.Lines = append(annotedStopPoint.Lines, objectid.Value())
			} else {
				defaultObjectID, ok := line.ObjectID("_default")
				if ok {
					annotedStopPoint.Lines = append(annotedStopPoint.Lines, defaultObjectID.Value())
				}
			}
		}
		if len(annotedStopPoint.Lines) == 0 {
			continue
		}
		annotedStopPointArray = append(annotedStopPointArray, annotedStopPoint.StopPointRef)
		response.AnnotatedStopPoints = append(response.AnnotatedStopPoints, annotedStopPoint)
	}

	sort.Sort(siri.SIRIAnnotatedStopPointByStopPointRef(response.AnnotatedStopPoints))

	logStashEvent["annotedStopPoints"] = strings.Join(annotedStopPointArray, ", ")
	logSIRIStopPointDiscoveryResponse(logStashEvent, response)

	return response, nil
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopPointsDiscoveryRequestBroadcaster"
	return event
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopDiscoveryRequestBroadcaster(partner)
}

func logXMLStopPointDiscoveryRequest(logStashEvent audit.LogStashEvent, request *siri.XMLStopPointsDiscoveryRequest) {
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIStopPointDiscoveryResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIStopPointsDiscoveryResponse) {
	logStashEvent["address"] = response.Address
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
