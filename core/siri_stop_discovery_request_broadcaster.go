package core

import (
	"sort"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/siri"
)

type SIRIStopPointsDiscoveryRequestBroadcaster struct {
	clock.ClockConsumer

	connector
}

type SIRIStopPointsDiscoveryRequestBroadcasterFactory struct{}

func NewSIRIStopDiscoveryRequestBroadcaster(partner *Partner) *SIRIStopPointsDiscoveryRequestBroadcaster {
	connector := &SIRIStopPointsDiscoveryRequestBroadcaster{}
	connector.remoteObjectidKind = partner.RemoteObjectIDKind(SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER)
	connector.partner = partner
	return connector
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) StopAreas(request *siri.XMLStopPointsDiscoveryRequest, message *audit.BigQueryMessage) (*siri.SIRIStopPointsDiscoveryResponse, error) {
	response := &siri.SIRIStopPointsDiscoveryResponse{
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	annotedStopPointMap := make(map[string]struct{})

	sas := connector.partner.Model().StopAreas().FindAll()
	for i := range sas {
		if sas[i].Name == "" || !sas[i].CollectedAlways {
			continue
		}

		objectID, ok := sas[i].ReferentOrSelfObjectId(connector.remoteObjectidKind)
		if !ok || objectID.Value() == "" {
			continue
		}

		_, ok = annotedStopPointMap[objectID.Value()]
		if ok || !ok && sas[i].ReferentId != "" {
			continue
		}

		annotedStopPointMap[objectID.Value()] = struct{}{}

		annotedStopPoint := &siri.SIRIAnnotatedStopPoint{
			StopName:     sas[i].Name,
			StopPointRef: objectID.Value(),
			Monitored:    true,
			TimingPoint:  true,
		}
		lines := sas[i].Lines()
		for i := range lines {
			if lines[i].Origin() == string(connector.partner.Slug()) {
				continue
			}
			objectid, ok := lines[i].ObjectID(connector.remoteObjectidKind)
			if !ok {
				continue
			}
			annotedStopPoint.Lines = append(annotedStopPoint.Lines, objectid.Value())
		}
		if len(annotedStopPoint.Lines) == 0 && connector.partner.IgnoreStopWithoutLine() {
			continue
		}
		response.AnnotatedStopPoints = append(response.AnnotatedStopPoints, annotedStopPoint)
	}

	sort.Sort(siri.SIRIAnnotatedStopPointByStopPointRef(response.AnnotatedStopPoints))

	message.RequestIdentifier = request.MessageIdentifier()

	logAnnotatedStopPoints(annotedStopPointMap, message)

	return response, nil
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopDiscoveryRequestBroadcaster(partner)
}

func logAnnotatedStopPoints(annotedStopPointMap map[string]struct{}, message *audit.BigQueryMessage) {
	keys := make([]string, len(annotedStopPointMap))
	i := 0
	for key := range annotedStopPointMap {
		keys[i] = key
		i++
	}

	message.StopAreas = keys
}
