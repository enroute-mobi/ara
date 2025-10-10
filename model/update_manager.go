package model

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model/hooks"
	"bitbucket.org/enroute-mobi/ara/model/model_types"
	"bitbucket.org/enroute-mobi/ara/siri/siri_attributes"
	"bitbucket.org/enroute-mobi/ara/uuid"

	"golang.org/x/exp/maps"
)

type UpdateManager struct {
	clock.ClockConsumer
	uuid.UUIDConsumer

	model Model
}

func NewUpdateManager(model Model) func(UpdateEvent) {
	manager := newUpdateManager(model)
	return manager.Update
}

// Test method
func newUpdateManager(model Model) *UpdateManager {
	return &UpdateManager{model: model}
}

func (manager *UpdateManager) Update(event UpdateEvent) {
	switch event.EventKind() {
	case STOP_AREA_EVENT:
		manager.updateStopArea(event.(*StopAreaUpdateEvent))
	case LINE_EVENT:
		manager.updateLine(event.(*LineUpdateEvent))
	case VEHICLE_JOURNEY_EVENT:
		manager.updateVehicleJourney(event.(*VehicleJourneyUpdateEvent))
	case STOP_VISIT_EVENT:
		manager.updateStopVisit(event.(*StopVisitUpdateEvent))
	case VEHICLE_EVENT:
		manager.updateVehicle(event.(*VehicleUpdateEvent))
	case STATUS_EVENT:
		manager.updateStatus(event.(*StatusUpdateEvent))
	case NOT_COLLECTED_EVENT:
		manager.updateNotCollected(event.(*NotCollectedUpdateEvent))
	case SITUATION_EVENT:
		manager.updateSituation(event.(*SituationUpdateEvent))
	case FACILITY_EVENT:
		manager.updateFacility(event.(*FacilityUpdateEvent))
	}
}

func (manager *UpdateManager) updateFacility(event *FacilityUpdateEvent) {
	facility, ok := manager.model.Facilities().FindByCode(event.Code)
	if !ok {
		facility = manager.model.Facilities().New()
		facility.Origin = event.Origin
		facility.SetCode(event.Code)
		facility.SetCode(NewCode(Default, event.Code.HashValue()))
	}

	if status, err := FacilityStatusFromString(event.Status); err == nil {
		facility.Status = *status
	}

	facility.Save()
}

func (manager *UpdateManager) updateSituation(event *SituationUpdateEvent) {
	situation, ok := manager.model.Situations().FindByCode(event.SituationCode)
	if ok &&
		situation.RecordedAt == event.RecordedAt &&
		situation.Version == event.Version {
		return
	}

	if !ok {
		situation = manager.model.Situations().New()
		situation.Origin = event.Origin
		situation.SetCode(event.SituationCode)
		situation.SetCode(NewCode(Default, event.SituationCode.HashValue()))
	}

	situation.RecordedAt = event.RecordedAt
	situation.Version = event.Version
	situation.ProducerRef = event.ProducerRef
	situation.ParticipantRef = event.ParticipantRef
	situation.InternalTags = event.InternalTags

	situation.Summary = event.Summary
	situation.Description = event.Description

	situation.VersionedAt = event.VersionedAt
	situation.ValidityPeriods = event.ValidityPeriods
	situation.PublicationWindows = event.PublicationWindows
	situation.Keywords = event.Keywords
	situation.ReportType = event.ReportType
	situation.AlertCause = event.AlertCause
	situation.Severity = event.Severity
	situation.Progress = event.Progress
	situation.Reality = event.Reality
	situation.Format = event.Format
	situation.Affects = event.Affects
	situation.Consequences = event.Consequences
	situation.PublishToWebActions = event.PublishToWebActions
	situation.PublishToMobileActions = event.PublishToMobileActions
	situation.PublishToDisplayActions = event.PublishToDisplayActions
	situation.InfoLinks = event.InfoLinks

	// Default is AfterCreate
	var h hooks.Type
	if ok {
		h = hooks.AfterSave
	}
	macros := manager.model.Macros().GetMacros(h, model_types.Situation)
	for i := range macros {
		macros[i].Update(situation)
	}
	controls := manager.model.Controls().GetControls(h, model_types.Situation)
	for i := range controls {
		controls[i].Control(situation)
	}

	situation.Save()
}

func (manager *UpdateManager) updateStopArea(event *StopAreaUpdateEvent) {
	if event.Code.Value() == "" { // Avoid creating a StopArea with an empty code
		return
	}

	stopArea, found := manager.model.StopAreas().FindByCode(event.Code)
	if !found {
		stopArea = manager.model.StopAreas().New()

		stopArea.SetCode(event.Code)
		stopArea.CollectSituations = true

		stopArea.Name = event.Name
		stopArea.CollectedAlways = event.CollectedAlways
		stopArea.Longitude = event.Longitude
		stopArea.Latitude = event.Latitude
	}

	if stopArea.ParentId == "" && event.ParentCode.Value() != "" {
		parentSA, _ := manager.model.StopAreas().FindByCode(event.ParentCode)
		stopArea.ParentId = parentSA.Id()
	}

	stopArea.Updated(manager.Clock().Now())

	// Default is AfterCreate
	var h hooks.Type
	if found {
		h = hooks.AfterSave
	}
	macros := manager.model.Macros().GetMacros(h, model_types.StopArea)
	for i := range macros {
		macros[i].Update(stopArea)
	}
	controls := manager.model.Controls().GetControls(h, model_types.StopArea)
	for i := range controls {
		controls[i].Control(stopArea)
	}

	manager.model.StopAreas().Save(stopArea)
	if event.Origin != "" {
		status, ok := stopArea.Origins.Origin(event.Origin)
		if !status || !ok {
			manager.updateStatus(NewStatusUpdateEvent(stopArea.Id(), event.Origin, true))
		}
	}
}

func (manager *UpdateManager) updateMonitoredStopArea(stopAreaId StopAreaId, partner string, status bool) {
	ascendants := manager.model.StopAreas().FindAscendants(stopAreaId)
	for i := range ascendants {
		stopArea := ascendants[i]
		stopArea.SetPartnerStatus(partner, status)
		manager.model.StopAreas().Save(stopArea)
	}
}

func (manager *UpdateManager) updateLine(event *LineUpdateEvent) {
	if event.Code.Value() == "" { // Avoid creating a Line with an empty code
		return
	}

	line, found := manager.model.Lines().FindByCode(event.Code)
	if !found {
		line = manager.model.Lines().New()

		line.SetCode(event.Code)
		line.SetCode(NewCode(Default, event.Code.HashValue()))

		line.CollectSituations = true

		line.Name = event.Name

		line.SetOrigin(event.Origin)
	}

	line.Updated(manager.Clock().Now())

	// Default is AfterCreate
	var h hooks.Type
	if found {
		h = hooks.AfterSave
	}
	macros := manager.model.Macros().GetMacros(h, model_types.Line)
	for i := range macros {
		macros[i].Update(line)
	}
	controls := manager.model.Controls().GetControls(h, model_types.Line)
	for i := range controls {
		controls[i].Control(line)
	}

	manager.model.Lines().Save(line)
}

func (manager *UpdateManager) updateVehicleJourney(event *VehicleJourneyUpdateEvent) {
	if event.Code.Value() == "" { // Avoid creating a VehicleJourney with an empty code
		return
	}

	if event.FromVehicleMonitoring {
		manager.updateVehicleJourneyFromVehicleMonitoring(event)
		return
	}

	vj, found := manager.model.VehicleJourneys().FindByCode(event.Code)
	if !found {
		// LineCode
		l, ok := manager.model.Lines().FindByCode(event.LineCode)
		if !ok {
			logger.Log.Debugf("VehicleJourney update event without corresponding line: %v", event.LineCode.String())
			return
		}

		vj = manager.model.VehicleJourneys().New()

		vj.SetCode(event.Code)
		vj.SetCode(NewCode(Default, event.Code.HashValue()))

		vj.Origin = event.Origin
		vj.Name = event.Attributes()[siri_attributes.VehicleJourneyName]
		vj.LineId = l.Id()
	}

	maps.Copy(vj.Attributes, event.Attributes())

	if vj.References.IsEmpty() {
		vj.References = event.References()
	}

	vj.References.SetCode("OriginRef", NewCode(event.Code.CodeSpace(), event.OriginRef))
	vj.OriginName = event.OriginName

	vj.References.SetCode("DestinationRef", NewCode(event.Code.CodeSpace(), event.DestinationRef))
	vj.DestinationName = event.DestinationName

	if event.Direction != "" { // Only used for Push collector
		vj.Attributes.Set(siri_attributes.DirectionName, event.Direction)
	}

	if event.Occupancy != Undefined {
		vj.Occupancy = event.Occupancy
	}
	vj.Monitored = event.Monitored
	if event.DirectionType != "" { // Do not override unknown DirectionType
		vj.DirectionType = event.DirectionType
	}

	// Default is AfterCreate
	var h hooks.Type
	if found {
		h = hooks.AfterSave
	}
	macros := manager.model.Macros().GetMacros(h, model_types.VehicleJourney)
	for i := range macros {
		macros[i].Update(vj)
	}
	controls := manager.model.Controls().GetControls(h, model_types.VehicleJourney)
	for i := range controls {
		controls[i].Control(vj)
	}

	manager.model.VehicleJourneys().Save(vj)
}

func (manager *UpdateManager) updateVehicleJourneyFromVehicleMonitoring(event *VehicleJourneyUpdateEvent) {
	vj, found := manager.model.VehicleJourneys().FindByCode(event.Code)
	if !found {
		return
	}

	if event.Occupancy != Undefined {
		vj.Occupancy = event.Occupancy
		manager.model.VehicleJourneys().Save(vj)
	}
}

func (manager *UpdateManager) updateStopVisit(event *StopVisitUpdateEvent) {
	if event.Code.Value() == "" { // Avoid creating a StopVisit with an empty code
		return
	}

	vj, ok := manager.model.VehicleJourneys().FindByCode(event.VehicleJourneyCode)
	if !ok {
		logger.Log.Debugf("StopVisit update event without corresponding vehicle journey: %v", event.VehicleJourneyCode.String())
		return
	}

	var sa *StopArea
	var sv *StopVisit
	var found bool
	if event.StopAreaCode.Value() == "" {
		sv, found = manager.model.StopVisits().FindByCode(event.Code)
		if !found {
			logger.Log.Debugf("Can't find Stopvisit from update event without stop area id")
			return
		}
		sa = sv.StopArea()
		if sa == nil {
			logger.Log.Printf("StopVisit in memory without a StopArea: %v", sv.Id())
			return
		}
	} else {
		sa, ok = manager.model.StopAreas().FindByCode(event.StopAreaCode)
		if !ok {
			logger.Log.Debugf("StopVisit update event without corresponding stop area: %v", event.StopAreaCode.String())
			return
		}

		sv, found = manager.model.StopVisits().FindByCode(event.Code)
		if !found {
			sv = manager.model.StopVisits().New()

			sv.SetCode(event.Code)
			sv.SetCode(NewCode(Default, event.Code.HashValue()))

			sv.StopAreaId = sa.Id()
			sv.VehicleJourneyId = vj.Id()

			sv.Origin = event.Origin
			sv.PassageOrder = event.PassageOrder
			sv.DataFrameRef = event.DataFrameRef
		}
	}

	if sv.Attributes.IsEmpty() {
		sv.Attributes = event.Attributes()
	}
	if sv.References.IsEmpty() {
		sv.References = event.References()
	}

	// Update StopArea Lines
	l := vj.Line()
	if l != nil {
		sa.LineIds.Add(l.Id())
		referent, ok := manager.model.StopAreas().Find(sa.ReferentId)
		if ok {
			referent.LineIds.Add(l.Id())
			manager.model.StopAreas().Save(referent)
		}
		manager.model.StopAreas().Save(sa)
	}

	if !event.RecordedAt.IsZero() {
		sv.RecordedAt = event.RecordedAt
	} else if !sv.Schedules.Include(event.Schedules) {
		sv.RecordedAt = manager.Clock().Now()
	}

	sv.Schedules.Merge(event.Schedules)
	if event.DepartureStatus != "" {
		sv.DepartureStatus = event.DepartureStatus
	}
	if event.ArrivalStatus != "" {
		sv.ArrivalStatus = event.ArrivalStatus
	}
	sv.VehicleAtStop = event.VehicleAtStop
	sv.Collected(manager.Clock().Now())

	if event.Monitored != vj.Monitored {
		vj.Monitored = event.Monitored
		manager.model.VehicleJourneys().Save(vj)
	}

	if event.Origin != "" {
		status, ok := sa.Origins.Origin(event.Origin)
		if status != event.Monitored || !ok {
			manager.updateMonitoredStopArea(sa.Id(), event.Origin, event.Monitored)
		}
	}

	// Default is AfterCreate
	var h hooks.Type
	if found {
		h = hooks.AfterSave
	}
	macros := manager.model.Macros().GetMacros(h, model_types.StopVisit)
	for i := range macros {
		macros[i].Update(sv)
	}
	controls := manager.model.Controls().GetControls(h, model_types.StopVisit)
	for i := range controls {
		controls[i].Control(sv)
	}

	manager.model.StopVisits().Save(sv)

	// VehicleJourney stop sequence
	if !vj.HasCompleteStopSequence {
		completeStopSequence := vj.model.ScheduledStopVisits().StopVisitsLenByVehicleJourney(vj.Id()) == vj.model.StopVisits().StopVisitsLenByVehicleJourney(vj.Id())
		if completeStopSequence {
			vj.HasCompleteStopSequence = true
			manager.model.VehicleJourneys().Save(vj)
		}
	}

	// long term historisation
	if sv.IsArchivable() {
		sva := &StopVisitArchiver{
			Model:     manager.model,
			StopVisit: sv,
		}
		sva.Archive()
	}

}

func (manager *UpdateManager) updateVehicle(event *VehicleUpdateEvent) {
	sa, _ := manager.model.StopAreas().FindByCode(event.StopAreaCode)
	vj, ok := manager.model.VehicleJourneys().FindByCode(event.VehicleJourneyCode)
	if !ok {
		logger.Log.Debugf("Vehicle update event without corresponding vehicle journey: %v", event.VehicleJourneyCode.String())
		return

	}

	line := vj.Line()

	vehicle, found := manager.model.Vehicles().FindByCode(event.Code)
	if !found {
		vehicle = manager.model.Vehicles().New()

		vehicle.SetCode(event.Code)
	}

	if event.NextStopPointOrder != 0 {
		if sv := manager.model.StopVisits().FindByVehicleJourneyIdAndStopVisitOrder(VehicleJourneyId(vj.Id()), event.NextStopPointOrder); sv != nil {
			vehicle.NextStopVisitId = sv.Id()
		}
	} else {
		if sa != nil {
			svIds := manager.model.StopVisits().FindByVehicleJourneyIdAndStopAreaId(VehicleJourneyId(vj.Id()), StopAreaId(sa.Id()))
			if len(svIds) == 1 {
				vehicle.NextStopVisitId = svIds[0]
			}
		}

	}

	vehicle.StopAreaId = sa.Id()
	vehicle.VehicleJourneyId = vj.Id()
	vehicle.DriverRef = event.DriverRef
	vehicle.Longitude = event.Longitude
	vehicle.Latitude = event.Latitude
	vehicle.Bearing = event.Bearing
	vehicle.LinkDistance = event.LinkDistance
	vehicle.Percentage = event.Percentage
	vehicle.ValidUntilTime = event.ValidUntilTime
	if event.RecordedAt.IsZero() {
		vehicle.RecordedAtTime = manager.Clock().Now()
	} else {
		vehicle.RecordedAtTime = event.RecordedAt
	}
	if event.Occupancy != Undefined {
		vehicle.Occupancy = event.Occupancy
	}

	if line != nil {
		vehicle.LineId = line.Id()
	}

	// Default is AfterCreate
	var h hooks.Type
	if found {
		h = hooks.AfterSave
	}
	macros := manager.model.Macros().GetMacros(h, model_types.Vehicle)
	for i := range macros {
		macros[i].Update(vehicle)
	}
	controls := manager.model.Controls().GetControls(h, model_types.Vehicle)
	for i := range controls {
		controls[i].Control(vehicle)
	}

	manager.model.Vehicles().Save(vehicle)
}

func (manager *UpdateManager) updateStatus(event *StatusUpdateEvent) {
	ascendants := manager.model.StopAreas().FindAscendants(event.StopAreaId)
	for i := range ascendants {
		stopArea := ascendants[i]
		stopArea.SetPartnerStatus(event.Partner, event.Status)
		manager.model.StopAreas().Save(stopArea)
	}
}

func (manager *UpdateManager) updateNotCollected(event *NotCollectedUpdateEvent) {
	stopVisit, found := manager.model.StopVisits().FindByCode(event.Code)
	if !found {
		logger.Log.Debugf("StopVisitNotCollectedEvent on unknown StopVisit: %#v", event)
		return
	}

	stopVisit.NotCollected(event.NotCollectedAt)
	manager.model.StopVisits().Save(stopVisit)
	if stopVisit.IsArchivable() {
		sva := &StopVisitArchiver{
			Model:     manager.model,
			StopVisit: stopVisit,
		}
		sva.Archive()
	}
	logger.Log.Debugf("StopVisit not Collected: %s (%v)", stopVisit.Id(), event.Code)
}
