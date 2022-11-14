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

type LineDirection struct {
	Id        model.LineId
	Direction string
}

type SIRIProductionTimetableBroadcaster interface {
	state.Stopable
	state.Startable
}

type PTTBroadcaster struct {
	clock.ClockConsumer

	connector *SIRIProductionTimetableSubscriptionBroadcaster
}

type ProductionTimetableBroadcaster struct {
	PTTBroadcaster

	stop chan struct{}
}

type FakeProductionTimetableBroadcaster struct {
	PTTBroadcaster
}

func NewFakeProductionTimetableBroadcaster(connector *SIRIProductionTimetableSubscriptionBroadcaster) SIRIProductionTimetableBroadcaster {
	broadcaster := &FakeProductionTimetableBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeProductionTimetableBroadcaster) Start() {
	broadcaster.prepareSIRIProductionTimetable()
}

func (broadcaster *FakeProductionTimetableBroadcaster) Stop() {}

func NewSIRIProductionTimetableBroadcaster(connector *SIRIProductionTimetableSubscriptionBroadcaster) SIRIProductionTimetableBroadcaster {
	broadcaster := &ProductionTimetableBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (ptt *ProductionTimetableBroadcaster) Start() {
	logger.Log.Debugf("Start ProductionTimetableBroadcaster")

	ptt.stop = make(chan struct{})
	go ptt.run()
}

func (ptt *ProductionTimetableBroadcaster) run() {
	c := ptt.Clock().After(5 * time.Second)

	for {
		select {
		case <-ptt.stop:
			logger.Log.Debugf("estimated time table broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRIProductionTimetableBroadcaster visit")

			ptt.prepareSIRIProductionTimetable()

			c = ptt.Clock().After(5 * time.Second)
		}
	}
}

func (ptt *ProductionTimetableBroadcaster) Stop() {
	if ptt.stop != nil {
		close(ptt.stop)
	}
}

func (ptt *PTTBroadcaster) prepareSIRIProductionTimetable() {
	ptt.connector.mutex.Lock()

	events := ptt.connector.toBroadcast
	ptt.connector.toBroadcast = make(map[SubscriptionId][]model.StopVisitId)

	ptt.connector.mutex.Unlock()

	currentTime := ptt.Clock().Now()

	for subId, stopVisits := range events {
		sub, ok := ptt.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			logger.Log.Debugf("PTT subscriptionBroadcast Could not find sub with id : %v", subId)
			continue
		}

		processedStopVisits := make(map[model.StopVisitId]struct{}) //Making sure not to send 2 times the same SV

		lines := make(map[LineDirection]*siri.SIRIDatedTimetableVersionFrame)
		vehicleJourneys := make(map[model.VehicleJourneyId]*siri.SIRIDatedVehicleJourney)

		delivery := &siri.SIRINotifyProductionTimetable{
			ProducerRef:            ptt.connector.Partner().ProducerRef(),
			SubscriptionIdentifier: sub.ExternalId(),
			ResponseTimestamp:      ptt.connector.Clock().Now(),
			Status:                 true,
			SortPayloadForTest:     ptt.connector.Partner().SortPaylodForTest(),
		}

		for _, stopVisitId := range stopVisits {
			// Check if resource is already in the map
			if _, ok := processedStopVisits[stopVisitId]; ok {
				continue
			}

			// Find the StopVisit
			stopVisit, ok := ptt.connector.Partner().Model().ScheduledStopVisits().Find(stopVisitId)
			if !ok {
				continue
			}

			// Handle StopPointRef
			stopArea, stopAreaId, ok := ptt.connector.stopPointRef(stopVisit.StopAreaId)
			if !ok {
				logger.Log.Printf("Ignore StopVisit %v without StopArea or with StopArea without correct ObjectID", stopVisit.Id())
				continue
			}

			// Find the VehicleJourney
			vehicleJourney, ok := ptt.connector.Partner().Model().VehicleJourneys().Find(stopVisit.VehicleJourneyId)
			if !ok {
				return
			}

			// Find the Line
			line, ok := ptt.connector.Partner().Model().Lines().Find(vehicleJourney.LineId)
			if !ok {
				continue
			}
			lineObjectId, ok := line.ObjectID(ptt.connector.remoteObjectidKind)
			if !ok {
				continue
			}

			// Find the Resource
			resource := sub.Resource(lineObjectId)
			if resource == nil {
				continue
			}

			// Get the DatedTimetableVersionFrame
			datedTTVersionFrame, ok := lines[LineDirection{Id: line.Id(), Direction: vehicleJourney.DirectionType}]
			if !ok {
				datedTTVersionFrame = &siri.SIRIDatedTimetableVersionFrame{
					LineRef:        lineObjectId.Value(),
					DirectionType:  ptt.connector.directionType(vehicleJourney.DirectionType),
					RecordedAtTime: currentTime,
					Attributes:     vehicleJourney.Attributes,
				}

				delivery.DatedTimetableVersionFrames = append(delivery.DatedTimetableVersionFrames, datedTTVersionFrame)
				lines[LineDirection{Id: line.Id(), Direction: vehicleJourney.DirectionType}] = datedTTVersionFrame
			}

			// Get the DatedVehicleJourney
			datedVehicleJourney, ok := vehicleJourneys[vehicleJourney.Id()]
			if !ok {
				// Handle vehicleJourney Objectid
				vehicleJourneyId, ok := vehicleJourney.ObjectIDWithFallback(ptt.connector.vjRemoteObjectidKinds)
				var datedVehicleJourneyRef string
				if ok {
					datedVehicleJourneyRef = vehicleJourneyId.Value()
				} else {
					defaultObjectID, ok := vehicleJourney.ObjectID("_default")
					if !ok {
						continue
					}
					referenceGenerator := ptt.connector.Partner().IdentifierGenerator(idgen.REFERENCE_IDENTIFIER)
					datedVehicleJourneyRef = referenceGenerator.NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultObjectID.Value()})
				}

				datedVehicleJourney = &siri.SIRIDatedVehicleJourney{
					DataFrameRef:           ptt.connector.dataFrameRef(),
					DatedVehicleJourneyRef: datedVehicleJourneyRef,
					PublishedLineName:      ptt.connector.publishedLineName(line),
					Attributes:             make(map[string]string),
					References:             make(map[string]string),
				}
				datedVehicleJourney.References["OperatorRef"] = ptt.connector.operatorRef(stopVisit)
				datedVehicleJourney.Attributes = vehicleJourney.Attributes

				datedTTVersionFrame.DatedVehicleJourneys = append(datedTTVersionFrame.DatedVehicleJourneys, datedVehicleJourney)
				vehicleJourneys[vehicleJourney.Id()] = datedVehicleJourney
			}

			// DatedCall
			datedCall := &siri.SIRIDatedCall{
				AimedArrivalTime:   stopVisit.Schedules.Schedule("aimed").ArrivalTime(),
				AimedDepartureTime: stopVisit.Schedules.Schedule("aimed").DepartureTime(),
				Order:              stopVisit.PassageOrder,
				StopPointRef:       stopAreaId,
				StopPointName:      stopArea.Name,
				DestinationDisplay: stopVisit.Attributes["DestinationDisplay"],
			}

			datedCall.UseVisitNumber = ptt.connector.useVisitNumber()

			datedVehicleJourney.DatedCalls = append(datedVehicleJourney.DatedCalls, datedCall)

			processedStopVisits[stopVisitId] = struct{}{}

			lastStateInterface, ok := resource.LastState(string(stopVisit.Id()))
			if !ok {
				resource.SetLastState(string(stopVisit.Id()), ls.NewProductionTimetableLastChange(stopVisit, sub))
			} else {
				lastState := lastStateInterface.(*ls.ProductionTimetableLastChange)
				lastState.UpdateState(stopVisit)
			}
		}
		ptt.sendDelivery(delivery)
	}
}
func (connector *SIRIProductionTimetableSubscriptionBroadcaster) useVisitNumber() bool {
	switch connector.Partner().PartnerSettings.SIRIPassageOrder() {
	case "visit_number":
		return true
	default:
		return false
	}
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) publishedLineName(line *model.Line) string {
	var pln string

	switch connector.partner.PartnerSettings.SIRILinePublishedName() {
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

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) directionType(direction string) (dir string) {
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

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) stopPointRef(stopAreaId model.StopAreaId) (*model.StopArea, string, bool) {
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

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) dataFrameRef() string {
	modelDate := connector.partner.Model().Date()
	return connector.dataFrameGenerator.NewIdentifier(idgen.IdentifierAttributes{Id: modelDate.String()})
}

func (connector *SIRIProductionTimetableSubscriptionBroadcaster) operatorRef(stopVisit *model.StopVisit) string {
	operatorRef, ok := stopVisit.Reference("OperatorRef")
	if !ok || operatorRef == (model.Reference{}) || operatorRef.ObjectId == nil {
		return ""
	}
	operator, ok := connector.Partner().Model().Operators().FindByObjectId(*operatorRef.ObjectId)
	if !ok {
		return operatorRef.ObjectId.Value()
	}
	obj, ok := operator.ObjectID(connector.remoteObjectidKind)
	if !ok {
		return operatorRef.ObjectId.Value()
	}

	return obj.Value()
}

func (ptt *PTTBroadcaster) sendDelivery(delivery *siri.SIRINotifyProductionTimetable) {
	message := ptt.newBQEvent()

	ptt.logSIRIProductionTimetableNotify(message, delivery)

	t := ptt.Clock().Now()

	err := ptt.connector.Partner().SIRIClient().NotifyProductionTimetable(delivery)
	message.ProcessingTime = ptt.Clock().Since(t).Seconds()
	if err != nil {
		message.Status = "Error"
		message.ErrorDetails = fmt.Sprintf("Error while sending ProductionTimetable notification: %v", err)
	}

	audit.CurrentBigQuery(string(ptt.connector.Partner().Referential().Slug())).WriteEvent(message)
}

func (ptt *PTTBroadcaster) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "NotifyProductionTimetable",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(ptt.connector.Partner().Slug()),
		Status:    "OK",
	}
}

func (ptt *PTTBroadcaster) logSIRIProductionTimetableNotify(message *audit.BigQueryMessage, response *siri.SIRINotifyProductionTimetable) {
	lineRefs := []string{}
	mr := make(map[string]struct{})
	for _, dttvf := range response.DatedTimetableVersionFrames {
		lineRefs = append(lineRefs, dttvf.LineRef)
		for _, vj := range dttvf.DatedVehicleJourneys {
			for _, ec := range vj.DatedCalls {
				mr[ec.StopPointRef] = struct{}{}
			}
		}
	}
	monitoringRefs := []string{}
	for k := range mr {
		monitoringRefs = append(monitoringRefs, k)
	}

	message.Lines = lineRefs
	message.StopAreas = monitoringRefs
	message.SubscriptionIdentifiers = []string{response.SubscriptionIdentifier}

	if !response.Status {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	xml, err := response.BuildXML(ptt.connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.ResponseRawMessage = xml
	message.ResponseSize = int64(len(xml))
}
