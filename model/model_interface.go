package model

type Model interface {
	StopAreas() StopAreas
	StopVisits() StopVisits
	VehicleJourneys() VehicleJourneys
	Lines() Lines
	Date() Date
	Reset() error
	// ...
}

type MemoryModel struct {
	stopAreas       StopAreas
	stopVisits      StopVisits
	vehicleJourneys VehicleJourneys
	lines           Lines
	date            Date
}

func NewMemoryModel() *MemoryModel {
	model := &MemoryModel{}

	stopAreas := NewMemoryStopAreas()
	stopAreas.model = model
	model.stopAreas = stopAreas

	stopVisits := NewMemoryStopVisits()
	stopVisits.model = model
	model.stopVisits = stopVisits

	vehicleJourneys := NewMemoryVehicleJourneys()
	vehicleJourneys.model = model
	model.vehicleJourneys = vehicleJourneys

	lines := NewMemoryLines()
	lines.model = model
	model.lines = lines

	model.date = NewDate(DefaultClock().Now())

	return model
}

func (model *MemoryModel) Reset() error {
	tx := NewTransaction(model)
	defer tx.Close()

	for _, vehicleJourney := range tx.Model().VehicleJourneys().FindAll() {
		tx.Model().VehicleJourneys().Delete(&vehicleJourney)
	}
	for _, stopVisit := range tx.Model().StopVisits().FindAll() {
		tx.Model().StopVisits().Delete(&stopVisit)
	}
	for _, line := range tx.Model().Lines().FindAll() {
		tx.Model().Lines().Delete(&line)
	}

	model.date = NewDate(DefaultClock().Now())
	return nil
}

func (model *MemoryModel) Date() Date {
	return model.date
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
