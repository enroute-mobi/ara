package model

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
	"bitbucket.org/enroute-mobi/ara/model/redisclient"
	"bitbucket.org/enroute-mobi/ara/uuid"
)

type hybridManager struct {
	uuid.UUIDConsumer

	model Model
}

func (mm *hybridManager) SetModel(m Model) {
	mm.model = m
}

type hybridModel struct {
	client redisclient.Client

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

func NewHybridModel(referential string, client redisclient.Client) Model {
	model := &hybridModel{
		client:      client,
		date:        NewDate(clock.DefaultClock().Now()),
		referential: referential,
	}

	model.refresh()

	return model
}

func NewTestHybridModel(codespaces []string, referential ...string) Model {
	model := &hybridModel{
		date:   NewDate(clock.DefaultClock().Now()),
		client: redisclient.TestClient,
	}

	if len(referential) != 0 {
		model.referential = referential[0]
	}

	model.refresh()

	return model
}

func (model *hybridModel) refresh() {
	lines := NewRedisLines(model.client)
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

func (model *hybridModel) RefreshMacros() {
	model.macros = NewMacroManager()
	model.macros.Load(model.referential)
}

func (model *hybridModel) RefreshControls() {
	model.controls = NewControlManager()
	model.controls.Load(model.referential)
}

func (model *hybridModel) SetBroadcastSMChan(broadcastSMEventChan chan StopMonitoringBroadcastEvent) {
	model.SMEventsChan = broadcastSMEventChan
}

func (model *hybridModel) SetBroadcastGMChan(broadcastGMEventChan chan SituationBroadcastEvent) {
	model.GMEventsChan = broadcastGMEventChan
}

func (model *hybridModel) SetBroadcastSXChan(broadcastSXEventChan chan SituationBroadcastEvent) {
	model.SXEventsChan = broadcastSXEventChan
}

func (model *hybridModel) SetBroadcastVeChan(broadcastVeEventChan chan VehicleBroadcastEvent) {
	model.VeEventChan = broadcastVeEventChan
}

func (model *hybridModel) SetBroadcastFMChan(broadcastFMEventChan chan FacilityBroadcastEvent) {
	model.FMEventChan = broadcastFMEventChan
}

func (model *hybridModel) Referential() string {
	return model.referential
}

func (model *hybridModel) SetReferential(referential string) {
	model.referential = referential
}

func (model *hybridModel) broadcastSMEvent(event StopMonitoringBroadcastEvent) {
	select {
	case model.SMEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager StopMonitoringBroadcastEvent queue is full")
	}
}

func (model *hybridModel) broadcastVeEvent(event VehicleBroadcastEvent) {
	select {
	case model.VeEventChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager VehicleBroadcastEvent queue is full")
	}
}

func (model *hybridModel) broadcastGMEvent(event SituationBroadcastEvent) {
	select {
	case model.GMEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager GeneralMessage SituationBroadcastEvent queue is full")
	}
}

func (model *hybridModel) broadcastSXEvent(event SituationBroadcastEvent) {
	select {
	case model.SXEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager SituationExchangeBroadcastEvent queue is full")
	}
}

func (model *hybridModel) broadcastFMEvent(event FacilityBroadcastEvent) {
	select {
	case model.FMEventChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager FacilityBroadcastEvent queue is full")
	}
}

func (model *hybridModel) Reload() Model {
	model.refresh()
	model.date = NewDate(clock.DefaultClock().Now())
	model.Load()
	return model
}

func (model *hybridModel) Date() Date {
	return model.date
}

func (model *hybridModel) SetDate(d Date) {
	model.date = d
}

func (model *hybridModel) Situations() Situations {
	return model.situations
}

func (model *hybridModel) StopAreas() StopAreas {
	return model.stopAreas
}

func (model *hybridModel) StopVisits() StopVisits {
	return model.stopVisits
}

func (model *hybridModel) ScheduledStopVisits() StopVisits {
	return model.scheduledStopVisits
}

func (model *hybridModel) VehicleJourneys() VehicleJourneys {
	return model.vehicleJourneys
}

func (model *hybridModel) Lines() Lines {
	return model.lines
}

func (model *hybridModel) LineGroups() LineGroups {
	return model.lineGroups
}

func (model *hybridModel) StopAreaGroups() StopAreaGroups {
	return model.stopAreaGroups
}

func (model *hybridModel) Operators() Operators {
	return model.operators
}

func (model *hybridModel) Vehicles() Vehicles {
	return model.vehicles
}

func (model *hybridModel) Macros() Macros {
	return model.macros
}

func (model *hybridModel) Controls() Controls {
	return model.controls
}

func (model *hybridModel) Facilities() Facilities {
	return model.facilities
}

func (model *hybridModel) Load() error {
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
