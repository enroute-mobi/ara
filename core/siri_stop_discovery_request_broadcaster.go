package core

import (
	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
	"bitbucket.org/enroute-mobi/ara/state"

	"golang.org/x/exp/maps"
	"slices"
	"sort"
)

type SIRIStopPointsDiscoveryRequestBroadcaster struct {
	state.Startable

	connector
}

type SIRIStopPointsDiscoveryRequestBroadcasterFactory struct{}

func NewSIRIStopDiscoveryRequestBroadcaster(partner *Partner) *SIRIStopPointsDiscoveryRequestBroadcaster {
	connector := &SIRIStopPointsDiscoveryRequestBroadcaster{}

	connector.partner = partner
	return connector
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) Start() {
	connector.remoteCodeSpace = connector.partner.RemoteCodeSpace(SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER)
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) StopAreas(request *sxml.XMLStopPointsDiscoveryRequest, message *audit.BigQueryMessage) (*siri.SIRIStopPointsDiscoveryResponse, error) {
	response := &siri.SIRIStopPointsDiscoveryResponse{
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	annotedStopPointMap := make(map[string]struct{})

	sas := connector.partner.Model().StopAreas().FindAllValues()

	// Replace Particulars by Referent
	referents := make(map[model.StopAreaId]model.StopArea)
	sas = slices.DeleteFunc(sas, func(stopArea model.StopArea) bool {
		referent, ok := stopArea.Referent()
		if ok {
			_, ok = referent.Code(connector.remoteCodeSpace)
			return ok
		}

		referents[referent.Id()] = *referent
		return false
	})

	sas = slices.Concat(sas, maps.Values(referents))

	for i := range sas {
		if sas[i].Name == "" || !sas[i].CollectedAlways {
			continue
		}

		code, ok := sas[i].Code(connector.remoteCodeSpace)
		if !ok {
			continue
		}
		if ok && code.Value() == "" {
			continue
		}

		annotedStopPointMap[code.Value()] = struct{}{}

		annotedStopPoint := &siri.SIRIAnnotatedStopPoint{
			StopName:     sas[i].Name,
			StopPointRef: code.Value(),
			Monitored:    true,
			TimingPoint:  true,
		}

		lines := sas[i].Lines()
		stopAreaLines := make(map[string]struct{})
		for i := range lines {
			if lines[i].Origin() == string(connector.partner.Slug()) {
				continue
			}
			code := model.Code{}
			referent, ok := lines[i].Referent()

			if ok {
				code, _ = referent.Code(connector.remoteCodeSpace)
			}

			if code.IsEmpty() {
				code, _ = lines[i].Code(connector.remoteCodeSpace)
			}

			if code.IsEmpty() {
				continue
			}
			stopAreaLines[code.Value()] = struct{}{}

		}
		annotedStopPoint.Lines = append(annotedStopPoint.Lines, maps.Keys(stopAreaLines)...)

		if len(annotedStopPoint.Lines) == 0 && connector.partner.IgnoreStopWithoutLine() {
			continue
		}

		sort.Slice(annotedStopPoint.Lines, func(i, j int) bool {
			return annotedStopPoint.Lines[i] < annotedStopPoint.Lines[j]
		})

		response.AnnotatedStopPoints = append(response.AnnotatedStopPoints, annotedStopPoint)

	}

	sort.Sort(siri.SIRIAnnotatedStopPointByStopPointRef(response.AnnotatedStopPoints))

	message.RequestIdentifier = request.MessageIdentifier()

	logAnnotatedStopPoints(annotedStopPointMap, message)

	return response, nil
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteCodeSpace()
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
