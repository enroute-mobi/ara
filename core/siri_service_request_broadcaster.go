package core

import (
	"fmt"
	"strconv"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type ServiceRequestBroadcaster interface {
	HandleRequests(request *siri.XMLSiriServiceRequest) (*siri.SIRIServiceResponse, error)
}

type SIRIServiceRequestBroadcaster struct {
	model.ClockConsumer
	model.UUIDConsumer

	siriConnector
}

type SIRIServiceRequestBroadcasterFactory struct{}

func NewSIRIServiceRequestBroadcaster(partner *Partner) *SIRIServiceRequestBroadcaster {
	siriServiceRequestBroadcaster := &SIRIServiceRequestBroadcaster{}
	siriServiceRequestBroadcaster.partner = partner
	return siriServiceRequestBroadcaster
}

func (connector *SIRIServiceRequestBroadcaster) HandleRequests(request *siri.XMLSiriServiceRequest) (*siri.SIRIServiceResponse, error) {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := make(audit.LogStashEvent)
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLSiriServiceRequest(logStashEvent, request)

	response := new(siri.SIRIServiceResponse)
	response.ProducerRef = connector.Partner().Setting("remote_credential")
	if response.ProducerRef == "" {
		response.ProducerRef = "Edwig"
	}
	response.ResponseMessageIdentifier = connector.SIRIPartner().NewMessageIdentifier()
	response.Status = true
	response.RequestMessageRef = request.MessageIdentifier()
	response.ResponseTimestamp = connector.Clock().Now()

	stopMonitoringConnector := NewSIRIStopMonitoringRequestBroadcaster(connector.partner)

	for _, stopMonitoringRequest := range request.StopMonitoringRequests() {
		logXMLSiriServiceStopMonitoringRequest(logStashEvent, stopMonitoringRequest)

		objectidKind := connector.Partner().Setting("remote_objectid_kind")
		objectid := model.NewObjectID(objectidKind, stopMonitoringRequest.MonitoringRef())
		stopArea, ok := tx.Model().StopAreas().FindByObjectId(objectid)

		if !ok {
			return nil, fmt.Errorf("StopArea not found")
		}

		delivery := stopMonitoringConnector.getStopMonitoringDelivery(tx, stopArea, stopMonitoringRequest.MessageIdentifier())

		// logSIRIStopMonitoringDelivery(logStashEvent, delivery)

		response.Deliveries = append(response.Deliveries, delivery)
	}

	logSIRIServiceResponse(logStashEvent, response)

	return response, nil
}

func (factory *SIRIServiceRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting("remote_objectid_kind")
	ok = ok && apiPartner.ValidatePresenceOfSetting("local_credential")
	return ok
}

func (factory *SIRIServiceRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIServiceRequestBroadcaster(partner)
}

func logXMLSiriServiceRequest(logStashEvent audit.LogStashEvent, request *siri.XMLSiriServiceRequest) {
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["requestorRef"] = request.RequestorRef()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIServiceResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIServiceResponse) {
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

func logXMLSiriServiceStopMonitoringRequest(logStashEvent audit.LogStashEvent, request *siri.XMLSiriServiceStopMonitoringRequest) {
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["monitoringRef"] = request.MonitoringRef()
	logStashEvent["stopVisitTypes"] = request.StopVisitTypes()
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["requestXML"] = request.RawXML()
}

// Empty for now
// func logSIRIStopMonitoringDelivery(logStashEvent audit.LogStashEvent, delivery siri.SIRIStopMonitoringDelivery) {
// }
