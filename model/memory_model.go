package model

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type memoryManager struct {
	uuid.UUIDConsumer

	model Model
}

func (mm *memoryManager) SetModel(m Model) {
	mm.model = m
}

type memoryModel struct {
	lines               Lines
	lineGroups          LineGroups
	stopAreaGroups      StopAreaGroups
	vehicles            Vehicles
	stopAreas           StopAreas
	stopVisits          StopVisits
	scheduledStopVisits StopVisits
	vehicleJourneys     VehicleJourneys
	situations          Situations
	operators           Operators
	macros              Macros
	controls            Controls
	facilities          Facilities
	SMEventsChan        chan StopMonitoringBroadcastEvent
	GMEventsChan        chan SituationBroadcastEvent
	SXEventsChan        chan SituationBroadcastEvent
	VeEventChan         chan VehicleBroadcastEvent
	FMEventChan         chan FacilityBroadcastEvent
	referential         string
	date                Date
}

func NewMemoryModel(referential string) Model {
	model := &memoryModel{
		date:        NewDate(clock.DefaultClock().Now()),
		referential: referential,
	}

	model.refresh()

	return model
}

func NewTestMemoryModel(referential ...string) Model {
	model := &memoryModel{
		date: NewDate(clock.DefaultClock().Now()),
	}

	if len(referential) != 0 {
		model.referential = referential[0]
	}

	model.refresh()

	return model
}

func (model *memoryModel) refresh() {
	lines := NewMemoryLines()
	lines.SetModel(model)
	model.lines = lines

	situations := NewMemorySituations()
	situations.SetModel(model)
	model.situations = situations
	model.situations.SetBroadcaster(model.broadcastGMEvent, GMbroadcastEvent)
	model.situations.SetBroadcaster(model.broadcastSXEvent, SXbroadcastEvent)

	stopAreas := NewMemoryStopAreas()
	stopAreas.SetModel(model)
	model.stopAreas = stopAreas
	model.stopAreas.SetBroadcaster(model.broadcastSMEvent)

	stopVisits := NewMemoryStopVisits()
	stopVisits.SetModel(model)
	model.stopVisits = stopVisits
	model.stopVisits.SetBroadcaster(model.broadcastSMEvent)

	scheduledStopVisits := NewMemoryStopVisits()
	scheduledStopVisits.SetModel(model)
	model.scheduledStopVisits = scheduledStopVisits

	vehicleJourneys := NewMemoryVehicleJourneys()
	vehicleJourneys.SetModel(model)
	model.vehicleJourneys = vehicleJourneys

	operators := NewMemoryOperators()
	operators.SetModel(model)
	model.operators = operators

	lineGroups := NewMemoryLineGroups()
	lineGroups.SetModel(model)
	model.lineGroups = lineGroups

	stopAreaGroups := NewMemoryStopAreaGroups()
	stopAreaGroups.SetModel(model)
	model.stopAreaGroups = stopAreaGroups

	vehicles := NewMemoryVehicles()
	vehicles.SetModel(model)
	model.vehicles = vehicles
	model.vehicles.SetBroadcaster(model.broadcastVeEvent)

	macros := NewMacroManager()
	macros.SetModel(model)
	model.macros = macros

	facilities := NewMemoryFacilities()
	facilities.SetModel(model)
	model.facilities = facilities
	model.facilities.SetBroadcaster(model.broadcastFMEvent)

	model.controls = NewControlManager()
}

func (model *memoryModel) RefreshMacros() {
	model.macros = NewMacroManager()
	model.macros.Load(model.referential)
}

func (model *memoryModel) RefreshControls() {
	model.controls = NewControlManager()
	model.controls.Load(model.referential)
}

func (model *memoryModel) SetBroadcastSMChan(broadcastSMEventChan chan StopMonitoringBroadcastEvent) {
	model.SMEventsChan = broadcastSMEventChan
}

func (model *memoryModel) SetBroadcastGMChan(broadcastGMEventChan chan SituationBroadcastEvent) {
	model.GMEventsChan = broadcastGMEventChan
}

func (model *memoryModel) SetBroadcastSXChan(broadcastSXEventChan chan SituationBroadcastEvent) {
	model.SXEventsChan = broadcastSXEventChan
}

func (model *memoryModel) SetBroadcastVeChan(broadcastVeEventChan chan VehicleBroadcastEvent) {
	model.VeEventChan = broadcastVeEventChan
}

func (model *memoryModel) SetBroadcastFMChan(broadcastFMEventChan chan FacilityBroadcastEvent) {
	model.FMEventChan = broadcastFMEventChan
}

func (model *memoryModel) Referential() string {
	return model.referential
}

func (model *memoryModel) SetReferential(referential string) {
	model.referential = referential
}

func (model *memoryModel) broadcastSMEvent(event StopMonitoringBroadcastEvent) {
	select {
	case model.SMEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager StopMonitoringBroadcastEvent queue is full")
	}
}

func (model *memoryModel) broadcastVeEvent(event VehicleBroadcastEvent) {
	select {
	case model.VeEventChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager VehicleBroadcastEvent queue is full")
	}
}

func (model *memoryModel) broadcastGMEvent(event SituationBroadcastEvent) {
	select {
	case model.GMEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager GeneralMessage SituationBroadcastEvent queue is full")
	}
}

func (model *memoryModel) broadcastSXEvent(event SituationBroadcastEvent) {
	select {
	case model.SXEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager SituationExchangeBroadcastEvent queue is full")
	}
}

func (model *memoryModel) broadcastFMEvent(event FacilityBroadcastEvent) {
	select {
	case model.FMEventChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager FacilityBroadcastEvent queue is full")
	}
}

func (model *memoryModel) Reload() Model {
	model.refresh()
	model.date = NewDate(clock.DefaultClock().Now())
	model.Load()
	return model
}

func (model *memoryModel) Date() Date {
	return model.date
}

func (model *memoryModel) SetDate(d Date) {
	model.date = d
}

func (model *memoryModel) Situations() Situations {
	return model.situations
}

func (model *memoryModel) StopAreas() StopAreas {
	return model.stopAreas
}

func (model *memoryModel) StopVisits() StopVisits {
	return model.stopVisits
}

func (model *memoryModel) ScheduledStopVisits() StopVisits {
	return model.scheduledStopVisits
}

func (model *memoryModel) VehicleJourneys() VehicleJourneys {
	return model.vehicleJourneys
}

func (model *memoryModel) Lines() Lines {
	return model.lines
}

func (model *memoryModel) LineGroups() LineGroups {
	return model.lineGroups
}

func (model *memoryModel) StopAreaGroups() StopAreaGroups {
	return model.stopAreaGroups
}

func (model *memoryModel) Operators() Operators {
	return model.operators
}

func (model *memoryModel) Vehicles() Vehicles {
	return model.vehicles
}

func (model *memoryModel) Macros() Macros {
	return model.macros
}

func (model *memoryModel) Controls() Controls {
	return model.controls
}

func (model *memoryModel) Facilities() Facilities {
	return model.facilities
}

func (model *memoryModel) Load() error {
	err := model.stopAreas.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading StopAreas: %v", err)
	}
	err = model.stopAreaGroups.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading StopAreaGroups: %v", err)
	}
	err = model.lines.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading Lines: %v", err)
	}
	err = model.lineGroups.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading LineGroups: %v", err)
	}
	err = model.vehicleJourneys.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading VehicleJourneys: %v", err)
	}
	err = model.scheduledStopVisits.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading StopVisits: %v", err)
	}
	err = model.operators.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading Operators: %v", err)
	}
	err = model.facilities.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading Facilities: %v", err)
	}
	err = model.macros.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading Macros: %v", err)
	}
	err = model.controls.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading Controls: %v", err)
	}
	return nil
}
