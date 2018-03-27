package core

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/af83/edwig/audit"
	"github.com/af83/edwig/logger"
	"github.com/af83/edwig/model"
	"github.com/af83/edwig/siri"
)

type SIRIEstimatedTimeTableBroadcaster interface {
	model.Stopable
	model.Startable
}

type ETTBroadcaster struct {
	model.ClockConsumer

	connector *SIRIEstimatedTimeTableSubscriptionBroadcaster
}

type EstimatedTimeTableBroadcaster struct {
	ETTBroadcaster

	stop chan struct{}
}

type FakeEstimatedTimeTableBroadcaster struct {
	ETTBroadcaster

	model.ClockConsumer
}

func NewFakeEstimatedTimeTableBroadcaster(connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) SIRIEstimatedTimeTableBroadcaster {
	broadcaster := &FakeEstimatedTimeTableBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeEstimatedTimeTableBroadcaster) Start() {
	broadcaster.prepareSIRIEstimatedTimeTable()
}

func (broadcaster *FakeEstimatedTimeTableBroadcaster) Stop() {}

func NewSIRIEstimatedTimeTableBroadcaster(connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) SIRIEstimatedTimeTableBroadcaster {
	broadcaster := &EstimatedTimeTableBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (ett *EstimatedTimeTableBroadcaster) Start() {
	logger.Log.Debugf("Start EstimatedTimeTableBroadcaster")

	ett.stop = make(chan struct{})
	go ett.run()
}

func (ett *EstimatedTimeTableBroadcaster) run() {
	c := ett.Clock().After(5 * time.Second)

	for {
		select {
		case <-ett.stop:
			logger.Log.Debugf("estimated time table broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRIEstimatedTimeTableBroadcaster visit")

			ett.prepareSIRIEstimatedTimeTable()
			ett.prepareNotMonitored()

			c = ett.Clock().After(5 * time.Second)
		}
	}
}

func (ett *EstimatedTimeTableBroadcaster) Stop() {
	if ett.stop != nil {
		close(ett.stop)
	}
}

func (ett *ETTBroadcaster) prepareNotMonitored() {
	ett.connector.mutex.Lock()

	notMonitored := ett.connector.notMonitored
	ett.connector.notMonitored = make(map[SubscriptionId]map[string]struct{})

	ett.connector.mutex.Unlock()

	for subId, producers := range notMonitored {
		sub, ok := ett.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			continue
		}

		for producer := range producers {
			delivery := &siri.SIRINotifyEstimatedTimeTable{
				Address:                   ett.connector.Partner().Address(),
				ProducerRef:               ett.connector.Partner().ProducerRef(),
				ResponseMessageIdentifier: ett.connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
				SubscriberRef:             ett.connector.SIRIPartner().RequestorRef(),
				SubscriptionIdentifier:    sub.ExternalId(),
				ResponseTimestamp:         ett.connector.Clock().Now(),
				Status:                    false,
				ErrorType:                 "OtherError",
				ErrorNumber:               1,
				ErrorText:                 fmt.Sprintf("Erreur [PRODUCER_UNAVAILABLE] : %v indisponible", producer),
				RequestMessageRef:         sub.SubscriptionOptions()["MessageIdentifier"],
			}

			ett.sendDelivery(delivery)
		}
	}
}

func (ett *ETTBroadcaster) prepareSIRIEstimatedTimeTable() {
	ett.connector.mutex.Lock()

	events := ett.connector.toBroadcast
	ett.connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	ett.connector.mutex.Unlock()

	tx := ett.connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	currentTime := ett.Clock().Now()

	for subId, stopVisits := range events {
		sub, ok := ett.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			logger.Log.Debugf("ETT subscriptionBroadcast Could not find sub with id : %v", subId)
			continue
		}

		processedStopVisits := make(map[model.StopVisitId]struct{}) //Making sure not to send 2 times the same SV
		lines := make(map[model.LineId]*siri.SIRIEstimatedJourneyVersionFrame)
		vehicleJourneys := make(map[model.VehicleJourneyId]*siri.SIRIEstimatedVehicleJourney)

		delivery := &siri.SIRINotifyEstimatedTimeTable{
			Address:                   ett.connector.Partner().Address(),
			ProducerRef:               ett.connector.Partner().ProducerRef(),
			ResponseMessageIdentifier: ett.connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
			SubscriberRef:             ett.connector.SIRIPartner().RequestorRef(),
			SubscriptionIdentifier:    sub.ExternalId(),
			ResponseTimestamp:         ett.connector.Clock().Now(),
			Status:                    true,
			RequestMessageRef:         sub.SubscriptionOptions()["MessageIdentifier"],
		}

		for _, stopVisitId := range stopVisits {
			// Check if resource is already in the map
			if _, ok := processedStopVisits[stopVisitId]; ok {
				continue
			}

			// Find the StopVisit
			stopVisit, ok := tx.Model().StopVisits().Find(stopVisitId)
			if !ok {
				continue
			}

			// Handle StopPointRef
			stopArea, ok := tx.Model().StopAreas().Find(stopVisit.StopAreaId)
			if !ok {
				continue
			}
			stopAreaId, ok := stopArea.ObjectID(ett.connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER))
			if !ok {
				continue
			}

			// Find the VehicleJourney
			vehicleJourney, ok := tx.Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
			if !ok {
				return
			}

			// Find the Line
			line, ok := tx.Model().Lines().Find(vehicleJourney.LineId)
			if !ok {
				continue
			}
			lineObjectId, ok := line.ObjectID(ett.connector.Partner().RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER))
			if !ok {
				continue
			}

			// Find the Resource
			resource := sub.Resource(lineObjectId)
			if resource == nil {
				continue
			}

			// Get the EstimatedJourneyVersionFrame
			journeyFrame, ok := lines[line.Id()]
			if !ok {
				journeyFrame = &siri.SIRIEstimatedJourneyVersionFrame{
					RecordedAtTime: currentTime,
				}

				delivery.EstimatedJourneyVersionFrames = append(delivery.EstimatedJourneyVersionFrames, journeyFrame)
				lines[line.Id()] = journeyFrame
			}

			// Get the EstiatedVehicleJourney
			estimatedVehicleJourney, ok := vehicleJourneys[vehicleJourney.Id()]
			if !ok {
				// Handle vehicleJourney Objectid
				vehicleJourneyId, ok := vehicleJourney.ObjectID(ett.connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_SUBSCRIPTION_BROADCASTER))
				var datedVehicleJourneyRef string
				if ok {
					datedVehicleJourneyRef = vehicleJourneyId.Value()
				} else {
					defaultObjectID, ok := vehicleJourney.ObjectID("_default")
					if !ok {
						continue
					}
					referenceGenerator := ett.connector.SIRIPartner().IdentifierGenerator("reference_identifier")
					datedVehicleJourneyRef = referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "VehicleJourney", Default: defaultObjectID.Value()})
				}

				estimatedVehicleJourney = &siri.SIRIEstimatedVehicleJourney{
					LineRef:                lineObjectId.Value(),
					DatedVehicleJourneyRef: datedVehicleJourneyRef,
					Attributes:             make(map[string]string),
					References:             make(map[string]model.Reference),
				}
				estimatedVehicleJourney.References = ett.connector.getEstimatedVehicleJourneyReferences(&vehicleJourney, &stopVisit, tx)
				estimatedVehicleJourney.Attributes = vehicleJourney.Attributes

				journeyFrame.EstimatedVehicleJourneys = append(journeyFrame.EstimatedVehicleJourneys, estimatedVehicleJourney)
				vehicleJourneys[vehicleJourney.Id()] = estimatedVehicleJourney
			}

			// EstimatedCall
			estimatedCall := &siri.SIRIEstimatedCall{
				ArrivalStatus:         string(stopVisit.ArrivalStatus),
				DepartureStatus:       string(stopVisit.DepartureStatus),
				AimedArrivalTime:      stopVisit.Schedules.Schedule("aimed").ArrivalTime(),
				ExpectedArrivalTime:   stopVisit.Schedules.Schedule("expected").ArrivalTime(),
				AimedDepartureTime:    stopVisit.Schedules.Schedule("aimed").DepartureTime(),
				ExpectedDepartureTime: stopVisit.Schedules.Schedule("expected").DepartureTime(),
				Order:              stopVisit.PassageOrder,
				StopPointRef:       stopAreaId.Value(),
				StopPointName:      stopArea.Name,
				DestinationDisplay: stopVisit.Attributes["DestinationDisplay"],
				VehicleAtStop:      stopVisit.VehicleAtStop,
			}

			estimatedVehicleJourney.EstimatedCalls = append(estimatedVehicleJourney.EstimatedCalls, estimatedCall)

			processedStopVisits[stopVisitId] = struct{}{}

			lastStateInterface, ok := resource.LastState(string(stopVisit.Id()))
			if !ok {
				ettlc := &estimatedTimeTableLastChange{}
				ettlc.InitState(&stopVisit, sub)
				resource.SetLastState(string(stopVisit.Id()), ettlc)
			} else {
				lastState := lastStateInterface.(*estimatedTimeTableLastChange)
				lastState.UpdateState(&stopVisit)
			}
		}
		ett.sendDelivery(delivery)
	}
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) getEstimatedVehicleJourneyReferences(vehicleJourney *model.VehicleJourney, stopVisit *model.StopVisit, tx *model.Transaction) map[string]model.Reference {
	references := make(map[string]model.Reference)

	for _, refType := range []string{"OriginRef", "DestinationRef"} {
		ref, ok := vehicleJourney.Reference(refType)
		if !ok || ref == (model.Reference{}) || ref.ObjectId == nil {
			continue
		}
		if foundStopArea, ok := tx.Model().StopAreas().FindByObjectId(*ref.ObjectId); ok {
			obj, ok := foundStopArea.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER))
			if ok {
				references[refType] = *model.NewReference(obj)
				continue
			}
		}
		generator := connector.SIRIPartner().IdentifierGenerator("reference_stop_area_identifier")
		defaultObjectID := model.NewObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER), generator.NewIdentifier(IdentifierAttributes{Default: ref.GetSha1()}))
		references[refType] = *model.NewReference(defaultObjectID)
	}

	// Handle OperatorRef
	operatorRef, ok := stopVisit.Reference("OperatorRef")
	if !ok || operatorRef == (model.Reference{}) || operatorRef.ObjectId == nil {
		return references
	}
	operator, ok := tx.Model().Operators().FindByObjectId(*operatorRef.ObjectId)
	if !ok {
		references["OperatorRef"] = *model.NewReference(*operatorRef.ObjectId)
		return references
	}
	obj, ok := operator.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER))
	if !ok {
		references["OperatorRef"] = *model.NewReference(*operatorRef.ObjectId)
		return references
	}
	references["OperatorRef"] = *model.NewReference(obj)
	return references
}

func (ett *ETTBroadcaster) sendDelivery(delivery *siri.SIRINotifyEstimatedTimeTable) {
	logStashEvent := ett.newLogStashEvent()
	logSIRIEstimatedTimeTableNotify(logStashEvent, delivery)
	audit.CurrentLogStash().WriteEvent(logStashEvent)

	err := ett.connector.SIRIPartner().SOAPClient().NotifyEstimatedTimeTable(delivery)
	if err != nil {
		event := ett.newLogStashEvent()
		logSIRINotifyError(err.Error(), event)
		audit.CurrentLogStash().WriteEvent(event)
	}
}

func (smb *ETTBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := smb.connector.partner.NewLogStashEvent()
	event["connector"] = "EstimatedTimeTableSubscriptionBroadcaster"
	return event
}

func logSIRIEstimatedTimeTableNotify(logStashEvent audit.LogStashEvent, response *siri.SIRINotifyEstimatedTimeTable) {
	lineRefs := []string{}
	for _, vjvf := range response.EstimatedJourneyVersionFrames {
		for _, vj := range vjvf.EstimatedVehicleJourneys {
			lineRefs = append(lineRefs, vj.LineRef)
		}
	}

	logStashEvent["siriType"] = "NotifyEstimatedTimetable"
	logStashEvent["producerRef"] = response.ProducerRef
	logStashEvent["requestMessageRef"] = response.RequestMessageRef
	logStashEvent["responseMessageIdentifier"] = response.ResponseMessageIdentifier
	logStashEvent["responseTimestamp"] = response.ResponseTimestamp.String()
	logStashEvent["subscriberRef"] = response.SubscriberRef
	logStashEvent["subscriptionIdentifier"] = response.SubscriptionIdentifier
	logStashEvent["lineRefs"] = strings.Join(lineRefs, ",")
	logStashEvent["status"] = strconv.FormatBool(response.Status)
	if !response.Status {
		logStashEvent["errorType"] = response.ErrorType
		if response.ErrorType == "OtherError" {
			logStashEvent["errorNumber"] = strconv.Itoa(response.ErrorNumber)
		}
		logStashEvent["errorText"] = response.ErrorText
	}
	xml, err := response.BuildXML()
	if err != nil {
		logStashEvent["responseXML"] = fmt.Sprintf("%v", err)
		return
	}
	logStashEvent["responseXML"] = xml
}
