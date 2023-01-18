package core

import (
	"fmt"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/state"
)

type EstimatedTimetableBroadcaster interface {
	state.Stopable
	state.Startable
}

type ETTBroadcaster struct {
	clock.ClockConsumer

	connector *SIRIEstimatedTimetableSubscriptionBroadcaster
}

type SIRIEstimatedTimetableBroadcaster struct {
	ETTBroadcaster

	stop chan struct{}
}

type FakeSIRIEstimatedTimetableBroadcaster struct {
	ETTBroadcaster
}

func NewFakeSIRIEstimatedTimetableBroadcaster(connector *SIRIEstimatedTimetableSubscriptionBroadcaster) EstimatedTimetableBroadcaster {
	broadcaster := &FakeSIRIEstimatedTimetableBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeSIRIEstimatedTimetableBroadcaster) Start() {
	broadcaster.prepareSIRIEstimatedTimetable()
}

func (broadcaster *FakeSIRIEstimatedTimetableBroadcaster) Stop() {}

func NewSIRIEstimatedTimetableBroadcaster(connector *SIRIEstimatedTimetableSubscriptionBroadcaster) EstimatedTimetableBroadcaster {
	broadcaster := &SIRIEstimatedTimetableBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (ett *SIRIEstimatedTimetableBroadcaster) Start() {
	logger.Log.Debugf("Start SIRIEstimatedTimetableBroadcaster")

	ett.stop = make(chan struct{})
	go ett.run()
}

func (ett *SIRIEstimatedTimetableBroadcaster) run() {
	c := ett.Clock().After(5 * time.Second)

	for {
		select {
		case <-ett.stop:
			logger.Log.Debugf("estimated time table broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRISIRIEstimatedTimetableBroadcaster visit")

			ett.prepareSIRIEstimatedTimetable()
			ett.prepareNotMonitored()

			c = ett.Clock().After(5 * time.Second)
		}
	}
}

func (ett *SIRIEstimatedTimetableBroadcaster) Stop() {
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
			delivery := &siri.SIRINotifyEstimatedTimetable{
				Address:                   ett.connector.Partner().Address(),
				ProducerRef:               ett.connector.Partner().ProducerRef(),
				ResponseMessageIdentifier: ett.connector.Partner().NewResponseMessageIdentifier(),
				SubscriberRef:             sub.SubscriberRef,
				SubscriptionIdentifier:    sub.ExternalId(),
				ResponseTimestamp:         ett.connector.Clock().Now(),
				Status:                    false,
				ErrorType:                 "OtherError",
				ErrorNumber:               1,
				ErrorText:                 fmt.Sprintf("Erreur [PRODUCER_UNAVAILABLE] : %v indisponible", producer),
				RequestMessageRef:         sub.SubscriptionOption("MessageIdentifier"),
			}

			ett.sendDelivery(delivery)
		}
	}
}

func (ett *ETTBroadcaster) prepareSIRIEstimatedTimetable() {
	ett.connector.mutex.Lock()

	events := ett.connector.toBroadcast
	ett.connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	ett.connector.mutex.Unlock()

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

		delivery := &siri.SIRINotifyEstimatedTimetable{
			Address:                   ett.connector.Partner().Address(),
			ProducerRef:               ett.connector.Partner().ProducerRef(),
			ResponseMessageIdentifier: ett.connector.Partner().NewResponseMessageIdentifier(),
			SubscriberRef:             sub.SubscriberRef,
			SubscriptionIdentifier:    sub.ExternalId(),
			ResponseTimestamp:         ett.connector.Clock().Now(),
			Status:                    true,
			RequestMessageRef:         sub.SubscriptionOption("MessageIdentifier"),
		}

		for _, stopVisitId := range stopVisits {
			// Check if resource is already in the map
			if _, ok := processedStopVisits[stopVisitId]; ok {
				continue
			}

			// Find the StopVisit
			stopVisit, ok := ett.connector.Partner().Model().StopVisits().Find(stopVisitId)
			if !ok {
				continue
			}

			// Handle StopPointRef
			stopArea, stopAreaId, ok := ett.connector.stopPointRef(stopVisit.StopAreaId)
			if !ok {
				logger.Log.Printf("Ignore StopVisit %v without StopArea or with StopArea without correct ObjectID", stopVisit.Id())
				continue
			}

			// Find the VehicleJourney
			vehicleJourney, ok := ett.connector.Partner().Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
			if !ok {
				return
			}

			// Find the Line
			line, ok := ett.connector.Partner().Model().Lines().Find(vehicleJourney.LineId)
			if !ok {
				continue
			}
			lineObjectId, ok := line.ObjectID(ett.connector.remoteObjectidKind)
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
				vehicleJourneyId, ok := vehicleJourney.ObjectIDWithFallback(ett.connector.vjRemoteObjectidKinds)
				var datedVehicleJourneyRef string
				if ok {
					datedVehicleJourneyRef = vehicleJourneyId.Value()
				} else {
					defaultObjectID, ok := vehicleJourney.ObjectID("_default")
					if !ok {
						continue
					}
					referenceGenerator := ett.connector.Partner().IdentifierGenerator(idgen.REFERENCE_IDENTIFIER)
					datedVehicleJourneyRef = referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultObjectID.Value()})
				}

				estimatedVehicleJourney = &siri.SIRIEstimatedVehicleJourney{
					LineRef:                lineObjectId.Value(),
					DirectionType:          ett.connector.directionType(vehicleJourney.DirectionType),
					DatedVehicleJourneyRef: datedVehicleJourneyRef,
					DataFrameRef:           ett.connector.dataFrameRef(),
					PublishedLineName:      ett.connector.publishedLineName(line),
					Attributes:             make(map[string]string),
					References:             make(map[string]string),
				}
				estimatedVehicleJourney.References = ett.connector.getEstimatedVehicleJourneyReferences(vehicleJourney, stopVisit)
				estimatedVehicleJourney.Attributes = vehicleJourney.Attributes

				journeyFrame.EstimatedVehicleJourneys = append(journeyFrame.EstimatedVehicleJourneys, estimatedVehicleJourney)
				vehicleJourneys[vehicleJourney.Id()] = estimatedVehicleJourney
			}

			var useVisitNumber = ett.connector.UseVisitNumber()

			if stopVisit.IsRecordable() && ett.connector.Partner().RecordedCallsDuration() != 0 {
				// recordedCall
				recordedCall := &siri.SIRIRecordedCall{
					ArrivalStatus:         string(stopVisit.ArrivalStatus),
					DepartureStatus:       string(stopVisit.DepartureStatus),
					AimedArrivalTime:      stopVisit.Schedules.Schedule("aimed").ArrivalTime(),
					ExpectedArrivalTime:   stopVisit.Schedules.Schedule("expected").ArrivalTime(),
					AimedDepartureTime:    stopVisit.Schedules.Schedule("aimed").DepartureTime(),
					ExpectedDepartureTime: stopVisit.Schedules.Schedule("expected").DepartureTime(),
					Order:                 stopVisit.PassageOrder,
					StopPointRef:          stopAreaId,
					StopPointName:         stopArea.Name,
					DestinationDisplay:    stopVisit.Attributes["DestinationDisplay"],
				}

				recordedCall.UseVisitNumber = useVisitNumber

				estimatedVehicleJourney.RecordedCalls = append(estimatedVehicleJourney.RecordedCalls, recordedCall)
			} else {
				// EstimatedCall
				estimatedCall := &siri.SIRIEstimatedCall{
					ArrivalStatus:         string(stopVisit.ArrivalStatus),
					DepartureStatus:       string(stopVisit.DepartureStatus),
					AimedArrivalTime:      stopVisit.Schedules.Schedule("aimed").ArrivalTime(),
					ExpectedArrivalTime:   stopVisit.Schedules.Schedule("expected").ArrivalTime(),
					AimedDepartureTime:    stopVisit.Schedules.Schedule("aimed").DepartureTime(),
					ExpectedDepartureTime: stopVisit.Schedules.Schedule("expected").DepartureTime(),
					Order:                 stopVisit.PassageOrder,
					StopPointRef:          stopAreaId,
					StopPointName:         stopArea.Name,
					DestinationDisplay:    stopVisit.Attributes["DestinationDisplay"],
					VehicleAtStop:         stopVisit.VehicleAtStop,
				}

				estimatedCall.UseVisitNumber = useVisitNumber

				estimatedVehicleJourney.EstimatedCalls = append(estimatedVehicleJourney.EstimatedCalls, estimatedCall)
			}

			max := max(ett.connector.Partner().Model().StopVisits().StopVisitsLenByVehicleJourney(vehicleJourney.Id()), ett.connector.Partner().Model().ScheduledStopVisits().StopVisitsLenByVehicleJourney(vehicleJourney.Id()))
			if len(estimatedVehicleJourney.RecordedCalls)+len(estimatedVehicleJourney.EstimatedCalls) == max {
				estimatedVehicleJourney.IsCompleteStopSequence = true
			}

			processedStopVisits[stopVisitId] = struct{}{}

			lastStateInterface, ok := resource.LastState(string(stopVisit.Id()))
			if !ok {
				resource.SetLastState(string(stopVisit.Id()), ls.NewEstimatedTimetableLastChange(stopVisit, sub))
			} else {
				lastStateInterface.(*ls.EstimatedTimetableLastChange).UpdateState(stopVisit)
			}
		}
		ett.sendDelivery(delivery)
	}
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) UseVisitNumber() bool {
	switch connector.Partner().PartnerSettings.SIRIPassageOrder() {
	case "visit_number":
		return true
	default:
		return false
	}
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) directionType(direction string) (dir string) {
	in, out, err := connector.partner.PartnerSettings.SIRIDirectionType()
	if err {
		return direction
	}

	switch direction {
	case model.VEHICLE_DIRECTION_INBOUND:
		dir = in
	case model.VEHICLE_DIRECTION_OUTBOUND:
		dir = out
	default:
		dir = direction
	}

	return dir
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) stopPointRef(stopAreaId model.StopAreaId) (*model.StopArea, string, bool) {
	stopPointRef, ok := connector.Partner().Model().StopAreas().Find(stopAreaId)
	if !ok {
		return &model.StopArea{}, "", false
	}
	stopPointRefObjectId, ok := stopPointRef.ObjectID(connector.remoteObjectidKind)
	if ok {
		return stopPointRef, stopPointRefObjectId.Value(), true
	}
	referent, ok := stopPointRef.Referent()
	if ok {
		referentObjectId, ok := referent.ObjectID(connector.remoteObjectidKind)
		if ok {
			return referent, referentObjectId.Value(), true
		}
	}
	parent, ok := stopPointRef.Parent()
	if ok {
		parentObjectId, ok := parent.ObjectID(connector.remoteObjectidKind)
		if ok {
			return parent, parentObjectId.Value(), true
		}
	}
	return &model.StopArea{}, "", false
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) publishedLineName(line *model.Line) string {
	var pln string

	switch connector.Partner().PartnerSettings.SIRILinePublishedName() {
	case "number":
		if line.Number != "" {
			pln = line.Number
		} else {
			pln = line.Name
		}
	default:
		pln = line.Name
	}

	return pln
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) getEstimatedVehicleJourneyReferences(vehicleJourney *model.VehicleJourney, stopVisit *model.StopVisit) map[string]string {
	references := make(map[string]string)

	for _, refType := range []string{"OriginRef", "DestinationRef"} {
		ref, ok := vehicleJourney.Reference(refType)
		if !ok || ref == (model.Reference{}) || ref.ObjectId == nil {
			continue
		}
		if refType == "DestinationRef" && connector.noDestinationRefRewrite(vehicleJourney.Origin) {
			references[refType] = ref.ObjectId.Value()
			continue
		}
		if foundStopArea, ok := connector.Partner().Model().StopAreas().FindByObjectId(*ref.ObjectId); ok {
			obj, ok := foundStopArea.ReferentOrSelfObjectId(connector.remoteObjectidKind)
			if ok {
				references[refType] = obj.Value()
				continue
			}
		}
		generator := connector.Partner().IdentifierGenerator(idgen.REFERENCE_STOP_AREA_IDENTIFIER)
		defaultObjectID := model.NewObjectID(connector.remoteObjectidKind, generator.NewIdentifier(idgen.IdentifierAttributes{Id: ref.GetSha1()}))
		references[refType] = defaultObjectID.Value()
	}

	// Handle OperatorRef
	operatorRef, ok := stopVisit.Reference("OperatorRef")
	if !ok || operatorRef == (model.Reference{}) || operatorRef.ObjectId == nil {
		return references
	}
	operator, ok := connector.Partner().Model().Operators().FindByObjectId(*operatorRef.ObjectId)
	if !ok {
		references["OperatorRef"] = operatorRef.ObjectId.Value()
		return references
	}
	obj, ok := operator.ObjectID(connector.remoteObjectidKind)
	if !ok {
		references["OperatorRef"] = operatorRef.ObjectId.Value()
		return references
	}
	references["OperatorRef"] = obj.Value()
	return references
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) dataFrameRef() string {
	modelDate := connector.partner.Model().Date()
	return connector.dataFrameGenerator.NewIdentifier(idgen.IdentifierAttributes{Id: modelDate.String()})
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) noDestinationRefRewrite(origin string) bool {
	for _, o := range connector.Partner().NoDestinationRefRewritingFrom() {
		if origin == o {
			return true
		}
	}
	return false
}

func (ett *ETTBroadcaster) sendDelivery(delivery *siri.SIRINotifyEstimatedTimetable) {
	message := ett.newBQEvent()

	ett.logSIRIEstimatedTimetableNotify(message, delivery)

	t := ett.Clock().Now()

	ett.connector.Partner().SIRIClient().NotifyEstimatedTimetable(delivery)
	message.ProcessingTime = ett.Clock().Since(t).Seconds()

	audit.CurrentBigQuery(string(ett.connector.Partner().Referential().Slug())).WriteEvent(message)
}

func (ett *ETTBroadcaster) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "NotifyEstimatedTimetable",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(ett.connector.Partner().Slug()),
		Status:    "OK",
	}
}

func (ett *ETTBroadcaster) logSIRIEstimatedTimetableNotify(message *audit.BigQueryMessage, response *siri.SIRINotifyEstimatedTimetable) {
	lineRefs := []string{}
	mr := make(map[string]struct{})
	for _, vjvf := range response.EstimatedJourneyVersionFrames {
		for _, vj := range vjvf.EstimatedVehicleJourneys {
			lineRefs = append(lineRefs, vj.LineRef)
			for _, ec := range vj.EstimatedCalls {
				mr[ec.StopPointRef] = struct{}{}
			}
		}
	}
	monitoringRefs := []string{}
	for k := range mr {
		monitoringRefs = append(monitoringRefs, k)
	}

	message.RequestIdentifier = response.RequestMessageRef
	message.ResponseIdentifier = response.ResponseMessageIdentifier
	message.Lines = lineRefs
	message.StopAreas = monitoringRefs
	message.SubscriptionIdentifiers = []string{response.SubscriptionIdentifier}

	if !response.Status {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	xml, err := response.BuildXML(ett.connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.ResponseRawMessage = xml
	message.ResponseSize = int64(len(xml))
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
