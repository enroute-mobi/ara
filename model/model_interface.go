package model

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
)

type ModelId string

type ModelInstance interface {
	CodeConsumerInterface

	modelId() ModelId
}

type Model interface {
	Date() Date
	Referential() string
	Lines() Lines
	Situations() Situations
	StopAreas() StopAreas
	StopVisits() StopVisits
	ScheduledStopVisits() StopVisits
	VehicleJourneys() VehicleJourneys
	Operators() Operators
	Vehicles() Vehicles
	Macros() Macros
}

type MemoryModel struct {
	lines               *MemoryLines
	vehicles            *MemoryVehicles
	stopAreas           *MemoryStopAreas
	stopVisits          *MemoryStopVisits
	scheduledStopVisits *MemoryStopVisits
	vehicleJourneys     *MemoryVehicleJourneys
	situations          *MemorySituations
	operators           *MemoryOperators
	macros              *MacroManager
	SMEventsChan        chan StopMonitoringBroadcastEvent
	GMEventsChan        chan SituationBroadcastEvent
	SXEventsChan        chan SituationBroadcastEvent
	VeEventChan         chan VehicleBroadcastEvent
	referential         string
	date                Date
}

func NewMemoryModel(referential string) *MemoryModel {
	model := &MemoryModel{
		date:        NewDate(clock.DefaultClock().Now()),
		referential: referential,
	}

	model.refresh()

	return model
}

func NewTestMemoryModel(referential ...string) *MemoryModel {
	model := &MemoryModel{
		date: NewDate(clock.DefaultClock().Now()),
	}

	if len(referential) != 0 {
		model.referential = referential[0]
	}

	model.refresh()

	return model
}

func (model *MemoryModel) refresh() {
	lines := NewMemoryLines()
	lines.model = model
	model.lines = lines

	situations := NewMemorySituations()
	situations.model = model
	model.situations = situations
	model.situations.GMbroadcastEvent = model.broadcastGMEvent
	model.situations.SXbroadcastEvent = model.broadcastSXEvent

	stopAreas := NewMemoryStopAreas()
	stopAreas.model = model
	model.stopAreas = stopAreas
	model.stopAreas.broadcastEvent = model.broadcastSMEvent

	stopVisits := NewMemoryStopVisits()
	stopVisits.model = model
	model.stopVisits = stopVisits
	model.stopVisits.broadcastEvent = model.broadcastSMEvent

	scheduledStopVisits := NewMemoryStopVisits()
	scheduledStopVisits.model = model
	model.scheduledStopVisits = scheduledStopVisits

	vehicleJourneys := NewMemoryVehicleJourneys()
	vehicleJourneys.model = model
	model.vehicleJourneys = vehicleJourneys

	operators := NewMemoryOperators()
	operators.model = model
	model.operators = operators

	vehicles := NewMemoryVehicles()
	vehicles.model = model
	model.vehicles = vehicles
	model.vehicles.broadcastEvent = model.broadcastVeEvent

	model.macros = NewMacroManager()
}

func (model *MemoryModel) RefreshMacros() {
	model.macros = NewMacroManager()
	model.macros.Load(model.referential)
}

func (model *MemoryModel) SetBroadcastSMChan(broadcastSMEventChan chan StopMonitoringBroadcastEvent) {
	model.SMEventsChan = broadcastSMEventChan
}

func (model *MemoryModel) SetBroadcastGMChan(broadcastGMEventChan chan SituationBroadcastEvent) {
	model.GMEventsChan = broadcastGMEventChan
}

func (model *MemoryModel) SetBroadcastSXChan(broadcastSXEventChan chan SituationBroadcastEvent) {
	model.SXEventsChan = broadcastSXEventChan
}

func (model *MemoryModel) SetBroadcastVeChan(broadcastVeEventChan chan VehicleBroadcastEvent) {
	model.VeEventChan = broadcastVeEventChan
}

func (model *MemoryModel) Referential() string {
	return model.referential
}

func (model *MemoryModel) SetReferential(referential string) {
	model.referential = referential
}

func (model *MemoryModel) broadcastSMEvent(event StopMonitoringBroadcastEvent) {
	select {
	case model.SMEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager StopMonitoringBroadcastEvent queue is full")
	}
}

func (model *MemoryModel) broadcastVeEvent(event VehicleBroadcastEvent) {
	select {
	case model.VeEventChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager VehicleBroadcastEvent queue is full")
	}
}

func (model *MemoryModel) broadcastGMEvent(event SituationBroadcastEvent) {
	select {
	case model.GMEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager GeneralMessage SituationBroadcastEvent queue is full")
	}
}

func (model *MemoryModel) broadcastSXEvent(event SituationBroadcastEvent) {
	select {
	case model.SXEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager SituationExchangeBroadcastEvent queue is full")
	}
}

func (model *MemoryModel) Reload() *MemoryModel {
	model.refresh()
	model.date = NewDate(clock.DefaultClock().Now())
	model.Load()
	return model
}

func (model *MemoryModel) Date() Date {
	return model.date
}

func (model *MemoryModel) Situations() Situations {
	return model.situations
}

func (model *MemoryModel) StopAreas() StopAreas {
	return model.stopAreas
}

func (model *MemoryModel) StopVisits() StopVisits {
	return model.stopVisits
}

func (model *MemoryModel) ScheduledStopVisits() StopVisits {
	return model.scheduledStopVisits
}

func (model *MemoryModel) VehicleJourneys() VehicleJourneys {
	return model.vehicleJourneys
}

func (model *MemoryModel) Lines() Lines {
	return model.lines
}

func (model *MemoryModel) Operators() Operators {
	return model.operators
}

func (model *MemoryModel) Vehicles() Vehicles {
	return model.vehicles
}

func (model *MemoryModel) Macros() Macros {
	return model.macros
}

func (model *MemoryModel) Load() error {
	err := model.stopAreas.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading StopAreas: %v", err)
	}
	err = model.lines.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading Lines: %v", err)
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
	err = model.macros.Load(model.referential)
	if err != nil {
		logger.Log.Debugf("Error while loading Macros: %v", err)
	}
	return nil
}
