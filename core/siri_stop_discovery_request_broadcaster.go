package core

import (
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type StopPointsDiscoveryRequestBroadcaster interface {
	stopAreas(request *siri.XMLStopDiscoveryRequest) (*siri.SIRIStopPointsDiscoveryResponse, error)
}

type SIRIStopPointsDiscoveryRequestBroadcaster struct {
	model.ClockConsumer

	siriConnector
}

type SIRIStopPointsDiscoveryRequestBroadcasterFactory struct{}

const (
	SIRI_STOP_POINTS_DISCOVERY_REQUEST_BROADCASTER = "siri-stop-points-discovery-request-broadcaster"
)

func NewSIRIStopDiscoveryRequestBroadcaster(partner *Partner) *SIRIStopPointsDiscoveryRequestBroadcaster {
	siriStopDiscoveryRequestBroadcaster := &SIRIStopPointsDiscoveryRequestBroadcaster{}
	siriStopDiscoveryRequestBroadcaster.partner = partner
	return siriStopDiscoveryRequestBroadcaster
}

func (connector *SIRIStopPointsDiscoveryRequestBroadcaster) StopAreas(request *siri.XMLStopDiscoveryRequest) (*siri.SIRIStopPointsDiscoveryResponse, error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	response := &siri.SIRIStopPointsDiscoveryResponse{}

	// response.Address = connector.Partner().Setting("local_url")
	// response.ProducerRef = connector.Partner().Setting("remote_credential")
	// if response.ProducerRef == "" {
	// 	response.ProducerRef = "Edwig"
	// }

	response.Status = true
	response.ResponseTimestamp = connector.Clock().Now()

	for _, stopArea := range tx.Model().StopAreas().FindAll() {
		if stopArea.Name == "" || stopArea.References["StopPointRef"] == (model.Reference{}) || stopArea.References["StopPointRef"].ObjectId.Value() == "" {
			continue
		}
		annotedStopPoint := &siri.SIRIAnnotatedStopPoint{
			StopPointName: stopArea.Name,
			StopPointRef:  stopArea.References["StopPointRef"].ObjectId.Value(),
		}
		response.AnnotatedStopPoints = append(response.AnnotatedStopPoints, annotedStopPoint)
	}
	return response, nil
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRIStopPointsDiscoveryRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopDiscoveryRequestBroadcaster(partner)
}
