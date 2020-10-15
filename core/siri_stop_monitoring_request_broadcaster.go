package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/edwig/audit"
	"bitbucket.org/enroute-mobi/edwig/logger"
	"bitbucket.org/enroute-mobi/edwig/model"
	"bitbucket.org/enroute-mobi/edwig/siri"
)

type StopMonitoringRequestBroadcaster interface {
	RequestStopArea(request *siri.XMLGetStopMonitoring) *siri.SIRIStopMonitoringResponse
}

type SIRIStopMonitoringRequestBroadcaster struct {
	model.ClockConsumer

	siriConnector
}

type SIRIStopMonitoringRequestBroadcasterFactory struct{}

func NewSIRIStopMonitoringRequestBroadcaster(partner *Partner) *SIRIStopMonitoringRequestBroadcaster {
	siriStopMonitoringRequestBroadcaster := &SIRIStopMonitoringRequestBroadcaster{}
	siriStopMonitoringRequestBroadcaster.partner = partner
	return siriStopMonitoringRequestBroadcaster
}

func (connector *SIRIStopMonitoringRequestBroadcaster) getStopMonitoringDelivery(tx *model.Transaction, logStashEvent audit.LogStashEvent, request *siri.XMLStopMonitoringRequest) siri.SIRIStopMonitoringDelivery {
	objectidKind := connector.partner.RemoteObjectIDKind(SIRI_STOP_MONITORING_REQUEST_BROADCASTER)
	objectid := model.NewObjectID(objectidKind, request.MonitoringRef())
	stopArea, ok := tx.Model().StopAreas().FindByObjectId(objectid)
	if !ok {
		return siri.SIRIStopMonitoringDelivery{
			RequestMessageRef: request.MessageIdentifier(),
			Status:            false,
			ResponseTimestamp: connector.Clock().Now(),
			ErrorType:         "InvalidDataReferencesError",
			ErrorText:         fmt.Sprintf("StopArea not found: '%s'", objectid.Value()),
			MonitoringRef:     request.MonitoringRef(),
		}
	}

	if !stopArea.CollectedAlways {
		stopArea.CollectedUntil = connector.Clock().Now().Add(15 * time.Minute)
		logger.Log.Printf("StopArea %s will be collected until %v", stopArea.Id(), stopArea.CollectedUntil)
		stopArea.Save()
	}

	delivery := siri.SIRIStopMonitoringDelivery{
		RequestMessageRef: request.MessageIdentifier(),
		Status:            true,
		ResponseTimestamp: connector.Clock().Now(),
		MonitoringRef:     request.MonitoringRef(),
	}
	if !stopArea.Monitored {
		delivery.Status = false
		delivery.ErrorType = "OtherError"
		delivery.ErrorNumber = 1
		delivery.ErrorText = fmt.Sprintf("Erreur [PRODUCER_UNAVAILABLE] : %v indisponible", strings.Join(stopArea.Origins.PartnersKO(), ", "))
	}

	// Prepare StopVisit Selectors
	selectors := []model.StopVisitSelector{}
	if request.LineRef() != "" {
		lineSelectorObjectid := model.NewObjectID(objectidKind, request.LineRef())
		selectors = append(selectors, model.StopVisitSelectorByLine(lineSelectorObjectid))
	}
	if request.PreviewInterval() != 0 {
		duration := request.PreviewInterval()
		now := connector.Clock().Now()
		if !request.StartTime().IsZero() {
			now = request.StartTime()
		}
		selectors = append(selectors, model.StopVisitSelectorByTime(now, now.Add(duration)))
	}
	selector := model.CompositeStopVisitSelector(selectors)

	// Prepare Id Array
	var stopVisitArray []string

	// Initialize builder
	stopMonitoringBuilder := NewBroadcastStopMonitoringBuilder(tx, connector.Partner(), SIRI_STOP_MONITORING_REQUEST_BROADCASTER)
	stopMonitoringBuilder.StopVisitTypes = request.StopVisitTypes()
	stopMonitoringBuilder.MonitoringRef = request.MonitoringRef()

	// Find Descendants
	stopAreas := tx.Model().StopAreas().FindFamily(stopArea.Id())

	// Fill StopVisits
	for _, stopVisit := range tx.Model().StopVisits().FindFollowingByStopAreaIds(stopAreas) {
		if stopVisit.Origin == string(connector.Partner().Slug()) {
			continue
		}
		if request.MaximumStopVisits() > 0 && len(stopVisitArray) >= request.MaximumStopVisits() {
			break
		}
		if !selector(stopVisit) {
			continue
		}

		monitoredStopVisit := stopMonitoringBuilder.BuildMonitoredStopVisit(stopVisit)
		if monitoredStopVisit == nil {
			continue
		}
		stopVisitArray = append(stopVisitArray, monitoredStopVisit.ItemIdentifier)
		delivery.MonitoredStopVisits = append(delivery.MonitoredStopVisits, monitoredStopVisit)
	}

	logStashEvent["stopVisitIds"] = strings.Join(stopVisitArray, ", ")

	return delivery
}

func (connector *SIRIStopMonitoringRequestBroadcaster) RequestStopArea(request *siri.XMLGetStopMonitoring) *siri.SIRIStopMonitoringResponse {
	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	logStashEvent := connector.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	logXMLStopMonitoringRequest(logStashEvent, &request.XMLStopMonitoringRequest)
	logStashEvent["requestorRef"] = request.RequestorRef()

	response := &siri.SIRIStopMonitoringResponse{
		Address:                   connector.Partner().Address(),
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.Partner().IdentifierGenerator(RESPONSE_MESSAGE_IDENTIFIER).NewMessageIdentifier(),
	}

	response.SIRIStopMonitoringDelivery = connector.getStopMonitoringDelivery(tx, logStashEvent, &request.XMLStopMonitoringRequest)

	logSIRIStopMonitoringDelivery(logStashEvent, response.SIRIStopMonitoringDelivery)
	logSIRIStopMonitoringResponse(logStashEvent, response)

	return response
}

func (connector *SIRIStopMonitoringRequestBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := connector.partner.NewLogStashEvent()
	event["connector"] = "StopMonitoringRequestBroadcaster"
	return event
}

func (factory *SIRIStopMonitoringRequestBroadcasterFactory) Validate(apiPartner *APIPartner) bool {
	ok := apiPartner.ValidatePresenceOfSetting(REMOTE_OBJECTID_KIND)
	ok = ok && apiPartner.ValidatePresenceOfLocalCredentials()
	return ok
}

func (factory *SIRIStopMonitoringRequestBroadcasterFactory) CreateConnector(partner *Partner) Connector {
	return NewSIRIStopMonitoringRequestBroadcaster(partner)
}

func logXMLStopMonitoringRequest(logStashEvent audit.LogStashEvent, request *siri.XMLStopMonitoringRequest) {
	logStashEvent["siriType"] = "StopMonitoringResponse"
	logStashEvent["messageIdentifier"] = request.MessageIdentifier()
	logStashEvent["monitoringRef"] = request.MonitoringRef()
	logStashEvent["stopVisitTypes"] = request.StopVisitTypes()
	logStashEvent["lineRef"] = request.LineRef()
	logStashEvent["maximumStopVisits"] = strconv.Itoa(request.MaximumStopVisits())
	logStashEvent["requestTimestamp"] = request.RequestTimestamp().String()
	logStashEvent["startTime"] = request.StartTime().String()
	logStashEvent["previewInterval"] = request.PreviewInterval().String()
	logStashEvent["requestXML"] = request.RawXML()
}

func logSIRIStopMonitoringDelivery(logStashEvent audit.LogStashEvent, delivery siri.SIRIStopMonitoringDelivery) {
	logStashEvent["requestMessageRef"] = delivery.RequestMessageRef
	logStashEvent["responseTimestamp"] = delivery.ResponseTimestamp.String()
	logStashEvent["status"] = strconv.FormatBool(delivery.Status)
	if !delivery.Status {
		logStashEvent["errorType"] = delivery.ErrorType
		if delivery.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(delivery.ErrorNumber)
		}
		logStashEvent["errorText"] = delivery.ErrorText
	}
}

func logSIRIStopMonitoringResponse(logStashEvent audit.LogStashEvent, response *siri.SIRIStopMonitoringResponse) {
	logStashEvent["address"] = response.Address
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
