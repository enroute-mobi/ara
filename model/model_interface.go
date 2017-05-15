package model

type Model interface {
	Date() Date
	Lines() Lines
	Situations() Situations
	StopAreas() StopAreas
	StopVisits() StopVisits
	VehicleJourneys() VehicleJourneys
	// ...
}

type MemoryModel struct {
	stopAreas       *MemoryStopAreas
	stopVisits      StopVisits
	vehicleJourneys VehicleJourneys
	lines           *MemoryLines
	date            Date
	situations      Situations
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

	stopAreas := NewMemoryStopAreas()
	stopAreas.model = model
	model.stopAreas = stopAreas

	stopVisits := NewMemoryStopVisits()
	stopVisits.model = model
	model.stopVisits = stopVisits

	vehicleJourneys := NewMemoryVehicleJourneys()
	vehicleJourneys.model = model
	model.vehicleJourneys = vehicleJourneys

	return model
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

func (model *MemoryModel) NewTransaction() *Transaction {
	return NewTransaction(model)
}

// TEMP: See what to do with errors
func (model *MemoryModel) Load(referentialId string) error {
	model.stopAreas.Load(referentialId)
	model.lines.Load(referentialId)
	return nil
}
