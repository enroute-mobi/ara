package core

import (
	"fmt"
	"sort"
	"strconv"

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

	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)
	logXMLStopPointDiscoveryRequest(logStashEvent, request)

	response := &siri.SIRIStopPointsDiscoveryResponse{}

	response.Address = connector.Partner().Setting("local_url")
	response.ProducerRef = connector.Partner().Setting("remote_credential")
	response.RequestMessageRef = request.MessageIdentifier()
	response.ResponseMessageIdentifier = connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier()
	response.Status = true
	response.ResponseTimestamp = connector.Clock().Now()

	if response.ProducerRef == "" {
		response.ProducerRef = "Edwig"
	}
	objectIDKind := connector.RemoteObjectIDKind()

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

		response.AnnotatedStopPoints = append(response.AnnotatedStopPoints, annotedStopPoint)
	}

	sort.Sort(siri.SIRIAnnotatedStopPointByStopPointRef(response.AnnotatedStopPoints))

	logSIRIStopPointDiscoveryResponse(logStashEvent, response)
	return response, nil
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) RemoteObjectIDKind() string {
	if connector.partner.Setting("siri-stop-points-discovery-request-broadcaster.remote_objectid_kind") != "" {
		return connector.partner.Setting("siri-stop-points-discovery-request-broadcaster.remote_objectid_kind")
	}
	return connector.Partner().Setting("remote_objectid_kind")
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
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIStopPointDiscoveryResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIStopPointsDiscoveryResponse) {
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
