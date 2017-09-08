package core

import (
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
	if ett.stop != nil {
		return
	}

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

			c = ett.Clock().After(5 * time.Second)
		}
	}
}

func (ett *EstimatedTimeTableBroadcaster) Stop() {
	if ett.stop != nil {
		close(ett.stop)
	}
}

func (ett *ETTBroadcaster) prepareSIRIEstimatedTimeTable() {
	connector := ett.connector

	connector.mutex.Lock()

	events := connector.toBroadcast
	connector.toBroadcast = make(map[SubscriptionId][]model.LineId)

	connector.mutex.Unlock()

	logStashEvent := ett.newLogStashEvent()
	defer audit.CurrentLogStash().WriteEvent(logStashEvent)

	tx := connector.Partner().Referential().NewTransaction()
	defer tx.Close()

	currentTime := connector.Clock().Now()

	notify := &siri.SIRINotifyEstimatedTimeTable{
		ResponseTimestamp:         currentTime,
		ProducerRef:               connector.Partner().ProducerRef(),
		ResponseMessageIdentifier: connector.SIRIPartner().IdentifierGenerator("response_message_identifier").NewMessageIdentifier(),
	}

	for subId, lines := range events {
		delivery := ett.getEstimatedTimetableDelivery(tx, lines)
		delivery.SubscriptionIdentifier = string(subId)
		delivery.SubscriberRef = connector.SIRIPartner().RequestorRef()
		notify.Deliveries = append(notify.Deliveries, &delivery)
	}
	connector.SIRIPartner().SOAPClient().NotifyEstimatedTimeTable(notify)
}

func (ett *ETTBroadcaster) getEstimatedTimetableDelivery(tx *model.Transaction, lines []model.LineId) siri.SIRIEstimatedTimetableSubscriptionDelivery {
	connector := ett.connector
	currentTime := connector.Clock().Now()
	sentlines := make(map[model.LineId]bool)

	delivery := siri.SIRIEstimatedTimetableSubscriptionDelivery{
		ResponseTimestamp: currentTime,
		Status:            true,
	}

	for _, lineId := range lines {
		if _, ok := sentlines[lineId]; ok {
			continue
		}

		sentlines[lineId] = true
		lineObjectId := model.NewObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER), string(lineId))
		line, ok := tx.Model().Lines().FindByObjectId(lineObjectId)
		if !ok {
			continue
		}

		journeyFrame := &siri.SIRIEstimatedJourneyVersionFrame{
			RecordedAtTime: currentTime,
		}

		// SIRIEstimatedVehicleJourney
		for _, vehicleJourney := range tx.Model().VehicleJourneys().FindByLineId(line.Id()) {
			// Handle vehicleJourney Objectid
			vehicleJourneyId, ok := vehicleJourney.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER))
			var datedVehicleJourneyRef string
			if ok {
				datedVehicleJourneyRef = vehicleJourneyId.Value()
			} else {
				defaultObjectID, ok := vehicleJourney.ObjectID("_default")
				if !ok {
					continue
				}
				referenceGenerator := connector.SIRIPartner().IdentifierGenerator("reference_identifier")
				datedVehicleJourneyRef = referenceGenerator.NewIdentifier(IdentifierAttributes{Type: "VehicleJourney", Default: defaultObjectID.Value()})
			}

			estimatedVehicleJourney := &siri.SIRIEstimatedVehicleJourney{
				LineRef:                lineObjectId.Value(),
				DatedVehicleJourneyRef: datedVehicleJourneyRef,
				Attributes:             make(map[string]string),
				References:             make(map[string]model.Reference),
			}
			estimatedVehicleJourney.References = connector.getEstimatedVehicleJourneyReferences(vehicleJourney, tx)
			estimatedVehicleJourney.Attributes = vehicleJourney.Attributes

			// SIRIEstimatedCall
			for _, stopVisit := range tx.Model().StopVisits().FindFollowingByVehicleJourneyId(vehicleJourney.Id()) {
				// Handle StopPointRef
				stopArea, ok := tx.Model().StopAreas().Find(stopVisit.StopAreaId)
				if !ok {
					continue
				}
				stopAreaId, ok := stopArea.ObjectID(connector.partner.RemoteObjectIDKind(SIRI_ESTIMATED_TIMETABLE_REQUEST_BROADCASTER))
				if !ok {
					continue
				}

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
			}

			journeyFrame.EstimatedVehicleJourneys = append(journeyFrame.EstimatedVehicleJourneys, estimatedVehicleJourney)
		}

		delivery.EstimatedJourneyVersionFrames = append(delivery.EstimatedJourneyVersionFrames, journeyFrame)
	}

	return delivery
}

func (connector *SIRIEstimatedTimeTableSubscriptionBroadcaster) getEstimatedVehicleJourneyReferences(vehicleJourney model.VehicleJourney, tx *model.Transaction) (references map[string]model.Reference) {
	for _, refType := range []string{"OriginRef", "DestinationRef"} {
		ref, ok := vehicleJourney.Reference(refType)
		if !ok || ref == (model.Reference{}) {
			continue
		}
		if foundStopArea, ok := tx.Model().StopAreas().Find(model.StopAreaId(ref.Id)); ok {
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
	return
}

func (smb *ETTBroadcaster) newLogStashEvent() audit.LogStashEvent {
	event := smb.connector.partner.NewLogStashEvent()
	event["connector"] = "EstimatedTimeTableSubscriptionBroadcaster"
	return event
}
