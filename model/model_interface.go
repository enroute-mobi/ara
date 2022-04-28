package model

import (
	"bitbucket.org/enroute-mobi/ara/clock"
	"bitbucket.org/enroute-mobi/ara/logger"
)

type ModelId string

type ModelInstance interface {
	ObjectIDConsumerInterface

	modelId() ModelId
}

type Model interface {
	Date() Date
	Referential() string
	Lines() Lines
	Situations() Situations
	StopAreas() StopAreas
	StopVisits() StopVisits
	VehicleJourneys() VehicleJourneys
	Operators() Operators
	Vehicles() Vehicles
}

type MemoryModel struct {
	date        Date
	referential string

	stopAreas       *MemoryStopAreas
	stopVisits      *MemoryStopVisits
	vehicleJourneys *MemoryVehicleJourneys
	lines           *MemoryLines
	vehicles        *MemoryVehicles
	situations      *MemorySituations
	operators       *MemoryOperators

	SMEventsChan chan StopMonitoringBroadcastEvent
	GMEventsChan chan GeneralMessageBroadcastEvent
}

// Optionnal argument for tests
func NewMemoryModel(referential ...string) *MemoryModel {
	model := &MemoryModel{
		date: NewDate(clock.DefaultClock().Now()),
	}

	if len(referential) != 0 {
		model.referential = referential[0]
	}

	lines := NewMemoryLines()
	lines.model = model
	model.lines = lines

	situations := NewMemorySituations()
	situations.model = model
	model.situations = situations
	model.situations.broadcastEvent = model.broadcastGMEvent

	stopAreas := NewMemoryStopAreas()
	stopAreas.model = model
	model.stopAreas = stopAreas
	model.stopAreas.broadcastEvent = model.broadcastSMEvent

	stopVisits := NewMemoryStopVisits()
	stopVisits.model = model
	model.stopVisits = stopVisits
	model.stopVisits.broadcastEvent = model.broadcastSMEvent

	vehicleJourneys := NewMemoryVehicleJourneys()
	vehicleJourneys.model = model
	model.vehicleJourneys = vehicleJourneys

	operators := NewMemoryOperators()
	operators.model = model
	model.operators = operators

	vehicles := NewMemoryVehicles()
	vehicles.model = model
	model.vehicles = vehicles

	return model
}

func (model *MemoryModel) SetBroadcastSMChan(broadcastSMEventChan chan StopMonitoringBroadcastEvent) {
	model.SMEventsChan = broadcastSMEventChan
}

func (model *MemoryModel) SetBroadcastGMChan(broadcastGMEventChan chan GeneralMessageBroadcastEvent) {
	model.GMEventsChan = broadcastGMEventChan
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

func (model *MemoryModel) broadcastGMEvent(event GeneralMessageBroadcastEvent) {
	select {
	case model.GMEventsChan <- event:
	default:
		logger.Log.Debugf("BrocasterManager GeneralMessageBroadcastEvent queue is full")
	}
}

func (model *MemoryModel) Reload(referentialSlug string) *MemoryModel {
	model = NewMemoryModel()
	model.date = NewDate(clock.DefaultClock().Now())
	model.stopAreas.Load(referentialSlug)
	model.lines.Load(referentialSlug)
	model.operators.Load(referentialSlug)
	model.vehicleJourneys.Load(referentialSlug)
	return model
}

func (model *MemoryModel) Clone() *MemoryModel {
	clone := NewMemoryModel()
	clone.stopAreas = model.stopAreas.Clone(clone)
	clone.lines = model.lines.Clone(clone)
	clone.date = NewDate(clock.DefaultClock().Now())
	return clone
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

// TEMP: See what to do with errors
func (model *MemoryModel) Load(referentialSlug string) error {
	err := model.stopAreas.Load(referentialSlug)
	if err != nil {
		logger.Log.Debugf("Error while loading StopAreas: %v", err)
	}
	err = model.lines.Load(referentialSlug)
	if err != nil {
		logger.Log.Debugf("Error while loading Lines: %v", err)
	}
	err = model.vehicleJourneys.Load(referentialSlug)
	if err != nil {
		logger.Log.Debugf("Error while loading VehicleJourneys: %v", err)
	}
	err = model.stopVisits.Load(referentialSlug)
	if err != nil {
		logger.Log.Debugf("Error while loading StopVisits: %v", err)
	}
	err = model.operators.Load(referentialSlug)
	if err != nil {
		logger.Log.Debugf("Error while loading Operators: %v", err)
	}
	return nil
}
