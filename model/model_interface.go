package model

import "github.com/af83/edwig/logger"

type Model interface {
	Date() Date
	Lines() Lines
	Situations() Situations
	StopAreas() StopAreas
	StopVisits() StopVisits
	VehicleJourneys() VehicleJourneys
	Operators() Operators
}

type MemoryModel struct {
	stopAreas       *MemoryStopAreas
	stopVisits      *MemoryStopVisits
	vehicleJourneys *MemoryVehicleJourneys
	lines           *MemoryLines
	date            Date
	situations      *MemorySituations
	operators       *MemoryOperators

	SMEventsChan chan StopMonitoringBroadcastEvent
	GMEventsChan chan GeneralMessageBroadcastEvent
}

func NewMemoryModel() *MemoryModel {
	model := &MemoryModel{}

	model.date = NewDate(DefaultClock().Now())

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

	return model
}

func (model *MemoryModel) SetBroadcastSMChan(broadcastSMEventChan chan StopMonitoringBroadcastEvent) {
	model.SMEventsChan = broadcastSMEventChan
}

func (model *MemoryModel) SetBroadcastGMChan(broadcastGMEventChan chan GeneralMessageBroadcastEvent) {
	model.GMEventsChan = broadcastGMEventChan
}

func (model *MemoryModel) broadcastSMEvent(event StopMonitoringBroadcastEvent) {
	select {
	case model.SMEventsChan <- event:
	default:
		logger.Log.Debugf("Cannot send StopMonitoringBroadcastEvent to BrocasterManager")
	}
}

func (model *MemoryModel) broadcastGMEvent(event GeneralMessageBroadcastEvent) {
	select {
	case model.GMEventsChan <- event:
	default:
		logger.Log.Debugf("Cannot send GeneralMessageBroadcastEvent to BrocasterManager")
	}
}

func (model *MemoryModel) Clone() *MemoryModel {
	clone := NewMemoryModel()
	clone.stopAreas = model.stopAreas.Clone(clone)
	clone.lines = model.lines.Clone(clone)
	clone.date = NewDate(DefaultClock().Now())
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

func (model *MemoryModel) NewTransaction() *Transaction {
	return NewTransaction(model)
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
