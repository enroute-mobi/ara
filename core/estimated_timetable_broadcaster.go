package core

import (
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

			c = ett.Clock().After(5 * time.Second)
		}
	}
}

func (ett *SIRIEstimatedTimetableBroadcaster) Stop() {
	if ett.stop != nil {
		close(ett.stop)
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
				logger.Log.Printf("Ignore StopVisit %v without StopArea or with StopArea without correct Code", stopVisit.Id())
				continue
			}

			// Find the VehicleJourney
			vehicleJourney, ok := ett.connector.Partner().Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
			if !ok {
				continue
			}

			// Find the Line
			line, ok := ett.connector.Partner().Model().Lines().Find(vehicleJourney.LineId)
			if !ok {
				continue
			}
			lineCode, ok := line.Code(ett.connector.remoteCodeSpace)
			if !ok {
				continue
			}

			// Find the Resource
			resource := sub.Resource(lineCode)
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
				// Handle vehicleJourney Code
				vehicleJourneyId, ok := vehicleJourney.CodeWithFallback(ett.connector.vjRemoteCodeSpaces)
				var datedVehicleJourneyRef string
				if ok {
					datedVehicleJourneyRef = vehicleJourneyId.Value()
				} else {
					defaultCode, ok := vehicleJourney.Code("_default")
					if !ok {
						continue
					}
					datedVehicleJourneyRef = ett.connector.Partner().NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultCode.Value()})
				}

				estimatedVehicleJourney = &siri.SIRIEstimatedVehicleJourney{
					LineRef:                lineCode.Value(),
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

			// Get StopVist call
			// Broadcast full stopVisit sequence if needed
			if vehicleJourney.HasCompleteStopSequence && !ett.connector.Partner().Model().VehicleJourneys().FullVehicleJourneyExistBySubscriptionId(string(subId), vehicleJourney.Id()) {
				logger.Log.Printf("ETT VehicleJourney %v full StopVisit broadcast coming from StopVisit %v", vehicleJourney.Id(), stopVisitId)
				for _, sv := range ett.connector.Partner().Model().StopVisits().FindByVehicleJourneyId(vehicleJourney.Id()) {
					sa, saId, ok := ett.connector.stopPointRef(sv.StopAreaId)
					if !ok {
						logger.Log.Printf("Ignore StopVisit %v without StopArea or with StopArea without correct Code", stopVisit.Id())
						continue
					}
					ett.connector.buildCall(sv, sa, saId, estimatedVehicleJourney)
					processedStopVisits[sv.Id()] = struct{}{}

				}
				ett.connector.Partner().Model().VehicleJourneys().SetFullVehicleJourneyBySubscriptionId(string(subId), vehicleJourney.Id())
			} else {
				// or broadcast single stopVisit
				ett.connector.buildCall(stopVisit, stopArea, stopAreaId, estimatedVehicleJourney)
				processedStopVisits[stopVisitId] = struct{}{}
			}

			// Set IsCompleteStopSequence

			if vehicleJourney.HasCompleteStopSequence {
				expectedLen := ett.connector.Partner().Model().StopVisits().StopVisitsLenByVehicleJourney(vehicleJourney.Id())
				if len(estimatedVehicleJourney.RecordedCalls)+len(estimatedVehicleJourney.EstimatedCalls) == expectedLen {
					estimatedVehicleJourney.IsCompleteStopSequence = true
				}
			}

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

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) buildCall(sv *model.StopVisit, sa *model.StopArea, saId string, evj *siri.SIRIEstimatedVehicleJourney) {
	var useVisitNumber = connector.UseVisitNumber()

	if sv.IsRecordable() && connector.Partner().RecordedCallsDuration() != 0 {
		// recordedCall
		recordedCall := &siri.SIRIRecordedCall{
			ArrivalStatus:         string(sv.ArrivalStatus),
			DepartureStatus:       string(sv.DepartureStatus),
			AimedArrivalTime:      sv.Schedules.Schedule("aimed").ArrivalTime(),
			ExpectedArrivalTime:   sv.Schedules.Schedule("expected").ArrivalTime(),
			AimedDepartureTime:    sv.Schedules.Schedule("aimed").DepartureTime(),
			ExpectedDepartureTime: sv.Schedules.Schedule("expected").DepartureTime(),
			Order:                 sv.PassageOrder,
			StopPointRef:          saId,
			StopPointName:         sa.Name,
			DestinationDisplay:    sv.Attributes["DestinationDisplay"],
		}

		recordedCall.UseVisitNumber = useVisitNumber

		evj.RecordedCalls = append(evj.RecordedCalls, recordedCall)
	} else {
		// EstimatedCall
		estimatedCall := &siri.SIRIEstimatedCall{
			ArrivalStatus:         string(sv.ArrivalStatus),
			DepartureStatus:       string(sv.DepartureStatus),
			AimedArrivalTime:      sv.Schedules.Schedule("aimed").ArrivalTime(),
			ExpectedArrivalTime:   sv.Schedules.Schedule("expected").ArrivalTime(),
			AimedDepartureTime:    sv.Schedules.Schedule("aimed").DepartureTime(),
			ExpectedDepartureTime: sv.Schedules.Schedule("expected").DepartureTime(),
			Order:                 sv.PassageOrder,
			StopPointRef:          saId,
			StopPointName:         sa.Name,
			DestinationDisplay:    sv.Attributes["DestinationDisplay"],
			VehicleAtStop:         sv.VehicleAtStop,
		}

		estimatedCall.UseVisitNumber = useVisitNumber

		evj.EstimatedCalls = append(evj.EstimatedCalls, estimatedCall)
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

	if connector.partner.PreferReferentStopArea() {
		referent, ok := stopPointRef.Referent()
		if ok {
			referentCode, ok := referent.Code(connector.remoteCodeSpace)
			if ok {
				return referent, referentCode.Value(), true
			}
		}
	}

	stopPointRefCode, ok := stopPointRef.Code(connector.remoteCodeSpace)
	if ok {
		return stopPointRef, stopPointRefCode.Value(), true
	}

	referent, ok := stopPointRef.Referent()
	if ok {
		referentCode, ok := referent.Code(connector.remoteCodeSpace)
		if ok {
			return referent, referentCode.Value(), true
		}
	}
	parent, ok := stopPointRef.Parent()
	if ok {
		parentCode, ok := parent.Code(connector.remoteCodeSpace)
		if ok {
			return parent, parentCode.Value(), true
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
		if !ok || ref == (model.Reference{}) || ref.Code == nil {
			continue
		}
		if refType == "DestinationRef" && connector.noDestinationRefRewrite(vehicleJourney.Origin) {
			references[refType] = ref.Code.Value()
			continue
		}
		if foundStopArea, ok := connector.Partner().Model().StopAreas().FindByCode(*ref.Code); ok {
			obj, ok := foundStopArea.ReferentOrSelfCode(connector.remoteCodeSpace)
			if ok {
				references[refType] = obj.Value()
				continue
			}
		}
		defaultCode := model.NewCode(connector.remoteCodeSpace, connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "StopArea", Id: ref.GetSha1()}))
		references[refType] = defaultCode.Value()
	}

	// Handle OperatorRef
	operatorRef, ok := stopVisit.Reference("OperatorRef")
	if !ok || operatorRef == (model.Reference{}) || operatorRef.Code == nil {
		return references
	}
	operator, ok := connector.Partner().Model().Operators().FindByCode(*operatorRef.Code)
	if !ok {
		references["OperatorRef"] = operatorRef.Code.Value()
		return references
	}
	obj, ok := operator.Code(connector.remoteCodeSpace)
	if !ok {
		references["OperatorRef"] = operatorRef.Code.Value()
		return references
	}
	references["OperatorRef"] = obj.Value()
	return references
}

func (connector *SIRIEstimatedTimetableSubscriptionBroadcaster) dataFrameRef() string {
	modelDate := connector.partner.Model().Date()
	return connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "DataFrame", Id: modelDate.String()})
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
	lineRefs := make(map[string]struct{})
	vehicleJourneyRefs := make(map[string]struct{})
	monitoringRefs := make(map[string]struct{})

	for _, vjvf := range response.EstimatedJourneyVersionFrames {
		for _, vj := range vjvf.EstimatedVehicleJourneys {
			lineRefs[vj.LineRef] = struct{}{}
			vehicleJourneyRefs[vj.DatedVehicleJourneyRef] = struct{}{}
			for _, estimatedCall := range vj.EstimatedCalls {
				monitoringRefs[estimatedCall.StopPointRef] = struct{}{}
			}

			for _, recordedCall := range vj.RecordedCalls {
				monitoringRefs[recordedCall.StopPointRef] = struct{}{}
			}
		}
	}

	message.RequestIdentifier = response.RequestMessageRef
	message.ResponseIdentifier = response.ResponseMessageIdentifier

	message.StopAreas = GetModelReferenceSlice(monitoringRefs)
	message.Lines = GetModelReferenceSlice(lineRefs)
	message.VehicleJourneys = GetModelReferenceSlice(vehicleJourneyRefs)

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
