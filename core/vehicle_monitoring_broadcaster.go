package core

import (
	"strings"
	"time"

	"bitbucket.org/enroute-mobi/ara/audit"
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/core/idgen"
	"bitbucket.org/enroute-mobi/ara/core/ls"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model"
	"bitbucket.org/enroute-mobi/ara/model/schedules"
	"bitbucket.org/enroute-mobi/ara/siri/siri"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/state"
)

type VehicleMonitoringBroadcaster interface {
	state.Stopable
	state.Startable
}

type VMBroadcaster struct {
	clock.ClockConsumer

	connector *SIRIVehicleMonitoringSubscriptionBroadcaster
}

type SIRIVehicleMonitoringBroadcaster struct {
	VMBroadcaster

	stop chan struct{}
}

type FakeSIRIVehicleMonitoringBroadcaster struct {
	VMBroadcaster
}

func NewFakeSIRIVehicleMonitoringBroadcaster(connector *SIRIVehicleMonitoringSubscriptionBroadcaster) VehicleMonitoringBroadcaster {
	broadcaster := &FakeSIRIVehicleMonitoringBroadcaster{}
	broadcaster.connector = connector
	return broadcaster
}

func (broadcaster *FakeSIRIVehicleMonitoringBroadcaster) Start() {
	broadcaster.prepareSIRIVehicleMonitoring()
}

func (broadcaster *FakeSIRIVehicleMonitoringBroadcaster) Stop() {}

func NewSIRIVehicleMonitoringBroadcaster(connector *SIRIVehicleMonitoringSubscriptionBroadcaster) VehicleMonitoringBroadcaster {
	broadcaster := &SIRIVehicleMonitoringBroadcaster{}
	broadcaster.connector = connector

	return broadcaster
}

func (vm *SIRIVehicleMonitoringBroadcaster) Start() {
	logger.Log.Debugf("Start SIRIVehicleMonitoringBroadcaster")

	vm.stop = make(chan struct{})
	go vm.run()
}

func (vm *SIRIVehicleMonitoringBroadcaster) Stop() {
	if vm.stop != nil {
		close(vm.stop)
	}
}

func (vm *SIRIVehicleMonitoringBroadcaster) run() {
	c := vm.Clock().After(5 * time.Second)

	for {
		select {
		case <-vm.stop:
			logger.Log.Debugf("vehicle monitoring broadcaster routine stop")

			return
		case <-c:
			logger.Log.Debugf("SIRIVehicleMonitoringBroadcaster visit")

			vm.prepareSIRIVehicleMonitoring()

			c = vm.Clock().After(5 * time.Second)
		}
	}
}

func (vm *VMBroadcaster) prepareSIRIVehicleMonitoring() {
	vm.connector.mutex.Lock()
	defer vm.connector.mutex.Unlock()

	events := vm.connector.toBroadcast
	vm.connector.toBroadcast = make(map[SubscriptionId][]model.VehicleId)

	for subId, vehicleIds := range events {
		sub, ok := vm.connector.Partner().Subscriptions().Find(subId)
		if !ok {
			logger.Log.Debugf("VM subscriptionBroadcast Could not find sub with id : %v", subId)
			continue
		}

		processedVehicles := make(map[model.VehicleId]struct{}) //Making sure not to send 2 times the same Vehicle
		delivery := &siri.SIRINotifyVehicleMonitoring{
			Address:                   vm.connector.Partner().Address(),
			ProducerRef:               vm.connector.Partner().ProducerRef(),
			ResponseMessageIdentifier: vm.connector.Partner().NewResponseMessageIdentifier(),
			SubscriberRef:             sub.SubscriberRef,
			SubscriptionIdentifier:    sub.ExternalId(),
			ResponseTimestamp:         vm.connector.Clock().Now(),
			Status:                    true,
			RequestMessageRef:         sub.SubscriptionOption("MessageIdentifier"),
			SortPayloadForTest:        vm.connector.Partner().SortPaylodForTest(),
		}

		for _, vehicleId := range vehicleIds {
			if _, ok := processedVehicles[vehicleId]; ok {
				continue
			}

			// Find the Vehicle
			vehicle, ok := vm.connector.Partner().Model().Vehicles().Find(vehicleId)
			if !ok {
				continue
			}
			vehicleCode, ok := vehicle.Code(vm.connector.remoteCodeSpace)
			if !ok {
				continue
			}

			// Find the VehicleJourney
			vj := vehicle.VehicleJourney()
			if vj == nil {
				continue
			}

			// Handle vj Code
			vjId, ok := vj.CodeWithFallback(vm.connector.vjRemoteCodeSpaces)
			var datedVehicleJourneyRef string
			if ok {
				datedVehicleJourneyRef = vjId.Value()
			} else {
				defaultCode, ok := vj.Code(model.Default)
				if !ok {
					continue
				}
				datedVehicleJourneyRef = vm.connector.Partner().NewIdentifier(idgen.IdentifierAttributes{Type: "VehicleJourney", Id: defaultCode.Value()})
			}

			// Find the Line
			line, ok := vm.connector.Partner().Model().Lines().Find(vj.LineId)
			if !ok {
				continue
			}
			lineCode, ok := line.Code(vm.connector.remoteCodeSpace)
			if !ok {
				continue
			}
			lineRef := lineCode.Value()

			// Find the Resource
			resource := sub.Resource(lineCode)
			if resource == nil {
				continue
			}

			refs := vj.References.Copy()

			activity := &siri.SIRIVehicleActivity{
				RecordedAtTime:       vehicle.RecordedAtTime,
				ValidUntilTime:       vehicle.ValidUntilTime,
				VehicleMonitoringRef: vehicleCode.Value(),
				ProgressBetweenStops: vm.connector.handleProgressBetweenStops(vehicle),
			}

			monitoredVehicleJourney := &siri.SIRIMonitoredVehicleJourney{
				LineRef:            lineRef,
				PublishedLineName:  line.Name,
				DirectionName:      vj.Attributes[siri_attributes.DirectionName],
				DirectionType:      vj.DirectionType,
				OriginName:         vj.OriginName,
				DestinationName:    vj.DestinationName,
				Monitored:          vj.Monitored,
				Bearing:            vehicle.Bearing,
				DriverRef:          vehicle.DriverRef,
				Occupancy:          vehicle.Occupancy,
				OriginRef:          vm.connector.handleRef(siri_attributes.OriginRef, vj.Origin, refs),
				DestinationRef:     vm.connector.handleRef(siri_attributes.DestinationRef, vj.Origin, refs),
				JourneyPatternRef:  vm.connector.handleJourneyPatternRef(refs),
				JourneyPatternName: vm.connector.handleJourneyPatternName(refs),
				VehicleLocation:    vm.connector.handleVehicleLocation(vehicle),
			}

			if vehicle.NextStopVisitId != model.StopVisitId("") {
				nextStopVisit, ok := vm.connector.Partner().Model().StopVisits().Find(vehicle.NextStopVisitId)
				if ok {
					stopArea, stopAreaCode, ok := vm.connector.stopPointRef(nextStopVisit.StopArea().Id())
					if ok {
						monitoredCall := &siri.MonitoredCall{
							StopPointRef:          stopAreaCode,
							StopPointName:         stopArea.Name,
							VehicleAtStop:         nextStopVisit.VehicleAtStop,
							DestinationDisplay:    nextStopVisit.Attributes["DestinationDisplay"],
							ExpectedArrivalTime:   nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Expected}),
							ExpectedDepartureTime: nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Expected}),
							DepartureStatus:       string(nextStopVisit.DepartureStatus),
							Order:                 &nextStopVisit.PassageOrder,
							AimedArrivalTime:      nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Aimed}),
							AimedDepartureTime:    nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Aimed}),
							ArrivalStatus:         string(nextStopVisit.ArrivalStatus),
							ActualArrivalTime:     nextStopVisit.Schedules.ArrivalTimeFromKind([]schedules.StopVisitScheduleType{schedules.Actual}),
							ActualDepartureTime:   nextStopVisit.Schedules.DepartureTimeFromKind([]schedules.StopVisitScheduleType{schedules.Actual}),
						}
						monitoredVehicleJourney.MonitoredCall = monitoredCall
					}
				}
			}

			framedVehicleJourneyRef := &siri.SIRIFramedVehicleJourneyRef{}
			modelDate := vm.connector.partner.Model().Date()
			framedVehicleJourneyRef.DataFrameRef =
				vm.connector.Partner().NewIdentifier(idgen.IdentifierAttributes{Type: "DataFrame", Id: modelDate.String()})
			framedVehicleJourneyRef.DatedVehicleJourneyRef = datedVehicleJourneyRef

			monitoredVehicleJourney.FramedVehicleJourneyRef = framedVehicleJourneyRef
			activity.MonitoredVehicleJourney = monitoredVehicleJourney
			delivery.VehicleActivities = append(delivery.VehicleActivities, activity)

			lastStateInterface, ok := resource.LastState(string(vehicle.Id()))
			if !ok {
				resource.SetLastState(string(vehicle.Id()), ls.NewVehicleMonitoringLastChange(vehicle, sub))
			} else {
				lastStateInterface.(*ls.VehicleMonitoringLastChange).UpdateState(vehicle)
			}

			processedVehicles[vehicleId] = struct{}{}
		}

		vm.sendDelivery(delivery)
	}

}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) stopPointRef(stopAreaId model.StopAreaId) (*model.StopArea, string, bool) {
	stopPointRef, ok := connector.partner.Model().StopAreas().Find(stopAreaId)
	if !ok {
		return &model.StopArea{}, "", false
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
	return &model.StopArea{}, "", false
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) handleRef(refType, origin string, references model.References) string {
	reference, ok := references.Get(refType)
	if !ok || reference.Code == nil || (refType == siri_attributes.DestinationRef && connector.noDestinationRefRewritingFrom(origin)) {
		return ""
	}
	return connector.resolveStopAreaRef(reference)
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) noDestinationRefRewritingFrom(origin string) bool {
	ndrrf := connector.Partner().NoDestinationRefRewritingFrom()
	for _, o := range ndrrf {
		if origin == strings.TrimSpace(o) {
			return true
		}
	}
	return false
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) resolveStopAreaRef(reference model.Reference) string {
	stopArea, ok := connector.partner.Model().StopAreas().FindByCode(*reference.Code)
	if ok {
		obj, ok := stopArea.ReferentOrSelfCode(connector.remoteCodeSpace)
		if ok {
			return obj.Value()
		}
	}
	return connector.partner.NewIdentifier(idgen.IdentifierAttributes{Type: "StopArea", Id: reference.GetSha1()})
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) handleJourneyPatternRef(refs model.References) string {
	journeyPatternRef, ok := refs.Get(siri_attributes.JourneyPatternRef)
	if ok {
		if connector.remoteCodeSpace == journeyPatternRef.Code.CodeSpace() {
			return journeyPatternRef.Code.Value()
		}
	}

	return ""
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) handleJourneyPatternName(refs model.References) string {
	journeyPatternName, ok := refs.Get(siri_attributes.JourneyPatternName)
	if ok {
		if connector.remoteCodeSpace == journeyPatternName.Code.CodeSpace() {
			return journeyPatternName.Code.Value()
		}
	}

	return ""
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) handleVehicleLocation(v *model.Vehicle) *siri.SIRIVehicleLocation {
	var lat = v.Latitude
	var lon = v.Longitude
	if lat != 0. || lon != 0. {
		vehicleLocation := &siri.SIRIVehicleLocation{
			Longitude: lon,
			Latitude:  lat,
		}
		return vehicleLocation
	}
	return nil
}

func (connector *SIRIVehicleMonitoringSubscriptionBroadcaster) handleProgressBetweenStops(v *model.Vehicle) *siri.SIRIProgressBetweenStops {
	var dist = v.LinkDistance
	var percent = v.Percentage
	if dist != 0. || percent != 0. {
		progressBetweenStops := &siri.SIRIProgressBetweenStops{
			LinkDistance: dist,
			Percentage:   percent,
		}
		return progressBetweenStops
	}
	return nil
}

func (vm *VMBroadcaster) sendDelivery(delivery *siri.SIRINotifyVehicleMonitoring) {
	message := vm.newBQEvent()

	vm.logSIRIVehicleMonitoring(message, delivery)

	t := vm.Clock().Now()

	vm.connector.Partner().SIRIClient().NotifyVehicleMonitoring(delivery)
	message.ProcessingTime = vm.Clock().Since(t).Seconds()

	audit.CurrentBigQuery(string(vm.connector.Partner().Referential().Slug())).WriteEvent(message)
}

func (vm *VMBroadcaster) newBQEvent() *audit.BigQueryMessage {
	return &audit.BigQueryMessage{
		Type:      "NotifyVehicleMonitoring",
		Protocol:  "siri",
		Direction: "sent",
		Partner:   string(vm.connector.Partner().Slug()),
		Status:    "OK",
	}
}

func (vm *VMBroadcaster) logSIRIVehicleMonitoring(message *audit.BigQueryMessage, response *siri.SIRINotifyVehicleMonitoring) {
	lineRefs := make(map[string]struct{})
	vehicleJourneyRefs := make(map[string]struct{})
	vehicleRefs := make(map[string]struct{})

	for _, va := range response.VehicleActivities {
		vehicleRefs[va.VehicleMonitoringRef] = struct{}{}

		line := va.MonitoredVehicleJourney.LineRef
		if _, ok := lineRefs[line]; ok {
			continue
		}

		lineRefs[line] = struct{}{}
		vehicleJourneyRefs[va.MonitoredVehicleJourney.FramedVehicleJourneyRef.DatedVehicleJourneyRef] = struct{}{}
	}

	message.ResponseIdentifier = response.ResponseMessageIdentifier

	message.Lines = GetModelReferenceSlice(lineRefs)
	message.VehicleJourneys = GetModelReferenceSlice(vehicleJourneyRefs)
	message.Vehicles = GetModelReferenceSlice(vehicleRefs)
	message.SubscriptionIdentifiers = []string{response.SubscriptionIdentifier}

	if !response.Status {
		message.Status = "Error"
		message.ErrorDetails = response.ErrorString()
	}
	xml, err := response.BuildXML(vm.connector.Partner().SIRIEnvelopeType())
	if err != nil {
		return
	}
	message.ResponseRawMessage = xml
	message.ResponseSize = int64(len(xml))
}
