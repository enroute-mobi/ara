package core

import (
	"sort"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/sxml"
)

type LinesDiscoveryRequestBroadcaster interface {
	Lines(*sxml.XMLLinesDiscoveryRequest, *audit.BigQueryMessage) (*siri.SIRILinesDiscoveryResponse, error)
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

func (connector *SIRILinesDiscoveryRequestBroadcaster) Lines(request *sxml.XMLLinesDiscoveryRequest, message *audit.BigQueryMessage) (*siri.SIRILinesDiscoveryResponse, error) {
	response := &siri.SIRILinesDiscoveryResponse{
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
	}

	var annotedLineArray []string

	lines := connector.partner.Model().Lines().FindAll()
	for i := range lines {
		if lines[i].Name == "" {
			continue
		}

		objectID, ok := lines[i].ObjectID(connector.remoteObjectidKind)
		if !ok {
			continue
		}

		annotedLine := &siri.SIRIAnnotatedLine{
			LineName:  lines[i].Name,
			LineRef:   objectID.Value(),
			Monitored: true,
		}
		annotedLineArray = append(annotedLineArray, annotedLine.LineRef)
		response.AnnotatedLines = append(response.AnnotatedLines, annotedLine)
	}

	sort.Sort(siri.SIRIAnnotatedLineByLineRef(response.AnnotatedLines))

	message.RequestIdentifier = request.MessageIdentifier()
	message.Lines = annotedLineArray

	return response, nil
}

func (factory *SIRILinesDiscoveryRequestBroadcasterFactory) Validate(apiPartner *APIPartner) {
	apiPartner.ValidatePresenceOfRemoteObjectIdKind()
	apiPartner.ValidatePresenceOfLocalCredentials()
}

func (factory *SIRILinesDiscoveryRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRILinesDiscoveryRequestBroadcaster(partner)
}
